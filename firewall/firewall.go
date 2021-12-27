package firewall

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	tokenKey = "X-CSRF-TOKEN"
)

type Source struct {
	Address  string `json:"address" tfsdk:"address"`
	FromPort int    `json:"-" tfsdk:"from_port"`
	ToPort   int    `json:"-" tfsdk:"to_port"`
	Port     string `json:"port" tfsdk:"-"`
	MAC      string `tfsdk:"mac"`
}

type Destination struct {
	Address  string `json:"address" tfsdk:"address"`
	FromPort int    `json:"-" tfsdk:"from_port"`
	ToPort   int    `json:"-" tfsdk:"to_port"`
	Port     string `json:"port" tfsdk:"-"`
}

type Rule struct {
	Priority    int          `json:"-" tfsdk:"priority"`
	Description string       `json:"description" tfsdk:"description"`
	Action      string       `json:"action" tfsdk:"action"`
	Protocol    string       `json:"protocol" tfsdk:"protocol"`
	Source      *Source      `json:"source" tfsdk:"source"`
	Destination *Destination `json:"destination" tfsdk:"destination"`
}

type Ruleset struct {
	Name          string           `json:"-" tfsdk:"name"`
	Description   string           `json:"description" tfsdk:"description"`
	DefaultAction string           `json:"default-action" tfsdk:"default_action"`
	Rules         []*Rule          `json:"-" tfsdk:"rule"`
	RulesMap      map[string]*Rule `json:"rule,omitempty" tfsdk:"-"`
}

type Firewall struct {
	Rulesets map[string]*Ruleset `json:"name"`
}

type Get struct {
	Firewall *Firewall `json:"firewall"`
}

type Set struct {
	Firewall *Firewall `json:"firewall"`
	Success  string    `json:"success,omitempty"`
	Failure  string    `json:"failure,omitempty"`
}

type Delete struct {
	Firewall *Firewall `json:"firewall,omitempty"`
	Success  string    `json:"success,omitempty"`
	Failure  string    `json:"failure,omitempty"`
}

type Commit struct {
	Success string `json:"success,omitempty"`
	Failure string `json:"failure,omitempty"`
}

type Save struct {
	Success string `json:"success,omitempty"`
}

type Operation struct {
	Success bool    `json:"success,omitempty"`
	Get     *Get    `json:"GET,omitempty"`
	Set     *Set    `json:"SET,omitempty"`
	Delete  *Delete `json:"DELETE,omitempty"`
	Commit  *Commit `json:"COMMIT,omitempty"`
	Save    *Save   `json:"SAVE,omitempty"`
}

type Client interface {
	GetRuleset(context.Context, string) (*Ruleset, error)
	CreateRuleset(context.Context, *Ruleset) (*Ruleset, error)
	DeleteRuleset(context.Context, string) error
}

type client struct {
	httpClient *http.Client
	baseURL    string
}

func New(httpClient *http.Client, baseURL string) Client {
	return &client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *client) GetRuleset(ctx context.Context, name string) (*Ruleset, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/edge/get.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return getRuleset(name, resp.Body)
}

func (c *client) CreateRuleset(ctx context.Context, p *Ruleset) (*Ruleset, error) {
	return c.post(ctx, p.Name, Operation{
		Set: &Set{
			Firewall: &Firewall{
				Rulesets: map[string]*Ruleset{
					p.Name: p.fromTerraform(),
				},
			},
		},
	})
}

func (c *client) DeleteRuleset(ctx context.Context, name string) error {
	_, err := c.post(ctx, name, Operation{
		Delete: &Delete{
			Firewall: &Firewall{
				Rulesets: map[string]*Ruleset{
					name: nil,
				},
			},
		},
	})
	return err
}

func (c *client) post(ctx context.Context, name string, in Operation) (*Ruleset, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/edge/batch.json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for _, cookie := range c.httpClient.Jar.Cookies(req.URL) {
		if cookie.Name == tokenKey {
			req.Header.Set(tokenKey, cookie.Value)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return getRuleset(name, resp.Body)
}

func getRuleset(name string, reader io.Reader) (*Ruleset, error) {
	var op Operation

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(data), &op); err != nil {
		return nil, err
	}

	if !op.Success {
		return nil, errors.New("The operation was not successfull.")
	}

	if op.Get == nil || op.Get.Firewall == nil {
		return nil, errors.New("A firewall does not exist.")
	}

	if op.Get.Firewall.Rulesets == nil {
		return nil, errors.New("there are no rulesets")
	}

	ruleset, ok := op.Get.Firewall.Rulesets[name]
	if !ok || ruleset == nil {
		return nil, errors.New("The ruleset does not exist.")
	}

	return ruleset.toTerraform(name)
}

func (rs *Ruleset) fromTerraform() *Ruleset {
	normalizePort := func(from, to int) string {
		if from == to {
			return strconv.Itoa(from)
		}
		return fmt.Sprintf("%d-%d", from, to)
	}

	if len(rs.Rules) > 0 {
		rs.RulesMap = map[string]*Rule{}
	}

	for _, rule := range rs.Rules {
		rs.RulesMap[strconv.Itoa(rule.Priority)] = rule

		if rule.Destination != nil {
			rule.Destination.Port = normalizePort(rule.Destination.FromPort, rule.Destination.ToPort)
		}

		if rule.Source != nil {
			rule.Source.Port = normalizePort(rule.Source.FromPort, rule.Source.ToPort)
		}

	}
	return rs
}

func (rs *Ruleset) toTerraform(name string) (_ *Ruleset, err error) {
	construct := func(port string) (int, int, error) {
		if !strings.Contains(port, "-") {
			i, err := strconv.Atoi(port)
			if err != nil {
				return 0, 0, errors.New("Malformed port.")
			}
			return i, i, nil
		}

		arr := strings.Split(port, "-")
		from, err := strconv.Atoi(arr[0])
		if err != nil {
			return 0, 0, errors.New("Malformed port.")
		}
		to, err := strconv.Atoi(arr[1])
		if err != nil {
			return 0, 0, errors.New("Malformed port.")
		}

		return from, to, nil
	}

	rs.Name = name

	for k, v := range rs.RulesMap {
		i, err := strconv.Atoi(k)
		if err != nil {
			return nil, errors.New("Malformed rule priority.")
		}
		v.Priority = i

		if v.Destination != nil && v.Destination.Port != "" {
			v.Destination.FromPort, v.Destination.ToPort, err = construct(v.Destination.Port)
			if err != nil {
				return nil, err
			}
		}

		if v.Source != nil && v.Source.Port != "" {
			v.Source.FromPort, v.Source.ToPort, err = construct(v.Source.Port)
			if err != nil {
				return nil, err
			}
		}

		rs.Rules = append(rs.Rules, v)
	}

	return rs, nil
}
