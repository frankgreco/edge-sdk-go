package firewall

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	tokenKey = "X-CSRF-TOKEN"
)

type Port struct {
	FromPort int `json:"-"`
	ToPort   int `json:"-"`
}

type Source struct {
	Address string `json:"address" tfsdk:"address"`
	MAC     string `tfsdk:"mac"`
	*Port
}

type Destination struct {
	Address string `json:"address"`
	*Port
}

type State struct {
	Established bool `json:"established" tfsdk:"established"`
	Invalid     bool `json:"invalid" tfsdk:"invalid"`
	New         bool `json:"new" tfsdk:"new"`
	Related     bool `json:"related" tfsdk:"related"`
}

type Rule struct {
	Priority    int          `json:"-" tfsdk:"priority"`
	Description string       `json:"description" tfsdk:"description"`
	Action      string       `json:"action" tfsdk:"action"`
	Protocol    string       `json:"protocol" tfsdk:"protocol"`
	Source      *Source      `json:"source" tfsdk:"source"`
	Destination *Destination `json:"destination" tfsdk:"destination"`
	State       *State       `json:"state" tfsdk:"state"`
	isTerraform bool
}

type Ruleset struct {
	Name          string  `json:"-" tfsdk:"name"`
	Description   string  `json:"description,omitempty" tfsdk:"description"`
	DefaultAction string  `json:"default-action,omitempty" tfsdk:"default_action"`
	Rules         []*Rule `json:"-" tfsdk:"rule"` // Omitting the json tag due to custom marshal/unmarshal methods.
	isTerraform   bool
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
	return toRuleset(name, resp.Body)
}

func (c *client) CreateRuleset(ctx context.Context, p *Ruleset) (*Ruleset, error) {
	return c.post(ctx, p.Name, Operation{
		Set: &Set{
			Firewall: &Firewall{
				Rulesets: map[string]*Ruleset{
					p.Name: p,
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
	return toRuleset(name, resp.Body)
}

func toRuleset(name string, reader io.Reader) (*Ruleset, error) {
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

	ruleset.Name = name
	return ruleset, nil

	// return ruleset.toTerraform(name)
}

func (s *Source) port() string {
	return port(s.FromPort, s.ToPort)
}

func (d *Destination) port() string {
	return port(d.FromPort, d.ToPort)
}

func (rs *Ruleset) Terraform() {
	(*rs).isTerraform = true
}

func (r *Rule) Terraform() {
	(*r).isTerraform = true
}

// // consider having
// // type ruleMap map[string]*Rule
// // and having a MarshalJSON for that instead.
func (rs *Ruleset) buildMap() map[string]*Rule {
	if rs == nil || len(rs.Rules) == 0 {
		return nil
	}

	m := map[string]*Rule{}
	for _, rule := range rs.Rules {
		m[strconv.Itoa(rule.Priority)] = rule
	}
	return m
}
