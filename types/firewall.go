package types

import (
	"fmt"
	"strconv"
	"strings"
)

type Port struct {
	FromPort int `json:"-"`
	ToPort   int `json:"-"`
}

type AddressGroup struct {
	Name        string   `json:"-" tfsdk:"name"`
	Description *string  `json:"description,omitempty" tfsdk:"description"`
	Cidrs       []string `json:"address,omitempty" tfsdk:"cidrs"`
}

type Source struct {
	AddressGroup string `json:"-" tfsdk:"address_group"`
	MAC          string `tfsdk:"mac"`
	*Port
}

type Destination struct {
	AddressGroup string `json:"-" tfsdk:"address_group"`
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
	codecMode   CodecMode
}

type Ruleset struct {
	Name          string  `json:"-" tfsdk:"name"`
	Description   string  `json:"description,omitempty" tfsdk:"description"`
	DefaultAction string  `json:"default-action,omitempty" tfsdk:"default_action"`
	Rules         []*Rule `json:"-" tfsdk:"rule"` // Omitting the json tag due to custom marshal/unmarshal methods.
	codecMode     CodecMode
	opMode        OpMode
}

type Groups struct {
	Address map[string]*AddressGroup `json:"address-group"`
}

type Firewall struct {
	Rulesets map[string]*Ruleset `json:"name,omitempty"`
	Groups   *Groups             `json:"group,omitempty"`
}

func (rs *Ruleset) SetCodecMode(c CodecMode) {
	(*rs).codecMode = c
}

func (rs *Ruleset) SetOpMode(m OpMode) {
	(*rs).opMode = m
}

func (r *Rule) SetCodecMode(c CodecMode) {
	(*r).codecMode = c
}

func (s *Source) port() string {
	return port(s.FromPort, s.ToPort)
}

func (d *Destination) port() string {
	return port(d.FromPort, d.ToPort)
}

func port(from, to int) string {
	if from == to {
		return strconv.Itoa(from)
	}
	return fmt.Sprintf("%d-%d", from, to)
}

func ports(p string) (int, int, error) {
	// if p == "" {
	// 	return 0, 0, nil
	// }

	if !strings.Contains(p, "-") {
		i, err := strconv.Atoi(p)
		if err != nil {
			return 0, 0, fmt.Errorf("The port is malformed: %s", err.Error())
		}
		return i, i, nil
	}

	arr := strings.Split(p, "-")
	from, err := strconv.Atoi(arr[0])
	if err != nil {
		return 0, 0, fmt.Errorf("The \"from\" port is malformed: %s", err.Error())
	}
	to, err := strconv.Atoi(arr[1])
	if err != nil {
		return 0, 0, fmt.Errorf("The \"to\" port is malformed: %s", err.Error())
	}

	return from, to, nil
}
