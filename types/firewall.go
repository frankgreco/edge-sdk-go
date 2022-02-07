package types

import (
	"fmt"
	"strconv"
	"strings"
)

type AddressGroup struct {
	Name        string   `json:"-" tfsdk:"name"`
	Description *string  `json:"description,omitempty" tfsdk:"description"`
	Cidrs       []string `json:"address,omitempty" tfsdk:"cidrs"`
}

type PortRange struct {
	From int `tfsdk:"from"`
	To   int `tfsdk:"to"`
}

type PortGroup struct {
	Name        string       `json:"-" tfsdk:"name"`
	Description *string      `json:"description,omitempty" tfsdk:"description"`
	Ports       []int        `json:"-" tfsdk:"ports"`
	Ranges      []*PortRange `json:"-" tfsdk:"port_ranges"`
}

type Source struct {
	Address      *string    `json:"address,omitempty" tfsdk:"address"`
	AddressGroup *string    `json:"-" tfsdk:"address_group"`
	PortGroup    *string    `json:"-" tfsdk:"port_group"`
	Port         *PortRange `json:"-" tfsdk:"port"`
	MAC          *string    `json:"mac,omitempty" tfsdk:"mac"`
}

type Destination struct {
	Address      *string    `json:"address,omitempty" tfsdk:"address"`
	AddressGroup *string    `json:"-" tfsdk:"address_group"`
	PortGroup    *string    `json:"-" tfsdk:"port_group"`
	Port         *PortRange `json:"-" tfsdk:"port"`
}

type State struct {
	Established *bool `json:"established" tfsdk:"established"`
	Invalid     *bool `json:"invalid" tfsdk:"invalid"`
	New         *bool `json:"new" tfsdk:"new"`
	Related     *bool `json:"related" tfsdk:"related"`
}

type Rule struct {
	Priority    int          `json:"-" tfsdk:"priority"`
	Description *string      `json:"description,omitempty" tfsdk:"description"`
	Action      string       `json:"action" tfsdk:"action"`
	Protocol    string       `json:"protocol" tfsdk:"protocol"`
	Source      *Source      `json:"source" tfsdk:"source"`
	Destination *Destination `json:"destination" tfsdk:"destination"`
	State       *State       `json:"state" tfsdk:"state"`
	Log         *bool        `json:"-" tfsdk:"log"`
	codecMode   CodecMode
}

type Ruleset struct {
	Name           string  `json:"-" tfsdk:"name"`
	Description    *string `json:"description,omitempty" tfsdk:"description"`
	DefaultAction  string  `json:"default-action,omitempty" tfsdk:"default_action"`
	DefaultLogging *bool   `json:"-" tfsdk:"default_logging"`
	Rules          []*Rule `json:"-" tfsdk:"rule"` // Omitting the json tag due to custom marshal/unmarshal methods.
	codecMode      CodecMode
	opMode         OpMode
}

type Groups struct {
	Address map[string]*AddressGroup `json:"address-group,omitempty"`
	Port    map[string]*PortGroup    `json:"port-group,omitempty"`
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

func (s *Source) toPort() string {
	return s.Port.toPort()
}

func (d *Destination) toPort() string {
	return d.Port.toPort()
}

func (p *PortRange) toPort() string {
	if p == nil {
		return ""
	}

	if p.From == p.To {
		return strconv.Itoa(p.From)
	}
	return fmt.Sprintf("%d-%d", p.From, p.To)
}

func (s *Source) fromPort(port string) error {
	portRange, err := fromPort(port)
	if err != nil {
		return err
	}
	s.Port = portRange
	return nil
}

func (d *Destination) fromPort(port string) error {
	portRange, err := fromPort(port)
	if err != nil {
		return err
	}
	d.Port = portRange
	return nil
}

func fromPort(port string) (p *PortRange, err error) {
	if port == "" {
		return nil, nil
	}

	var from, to int

	if !strings.Contains(port, "-") {
		from, err = strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("Could not turn %s into a valid port: %s", port, err.Error())
		}
		to = from
	} else {
		fromTo := strings.Split(port, "-")
		from, err = strconv.Atoi(fromTo[0])
		if err != nil {
			return nil, fmt.Errorf("The \"from\" port is malformed for port range %s: %s", port, err.Error())
		}
		to, err = strconv.Atoi(fromTo[1])
		if err != nil {
			return nil, fmt.Errorf("The \"to\" port is malformed for port range %s: %s", port, err.Error())
		}
	}

	return &PortRange{
		From: from,
		To:   to,
	}, nil
}

func (g *PortGroup) toPorts() []string {
	if g == nil || (len(g.Ports) == 0 && len(g.Ranges) == 0) {
		return nil
	}

	ports := []string{}
	for _, port := range g.Ports {
		ports = append(ports, strconv.Itoa(port))
	}

	for _, portRange := range g.Ranges {
		if portRange == nil {
			continue
		}
		if portRange.From != portRange.To {
			ports = append(ports, fmt.Sprintf("%s-%s", strconv.Itoa(portRange.From), strconv.Itoa(portRange.To)))
		}
	}

	return ports
}

func (g *PortGroup) fromPorts(ports []string) error {
	for _, port := range ports {
		if !strings.Contains(port, "-") {
			i, err := strconv.Atoi(port)
			if err != nil {
				return fmt.Errorf("Could not turn %s into a valid port: %s", port, err.Error())
			}
			if g.Ports == nil {
				g.Ports = []int{}
			}
			g.Ports = append(g.Ports, i)
			continue
		}
		fromTo := strings.Split(port, "-")
		from, err := strconv.Atoi(fromTo[0])
		if err != nil {
			return fmt.Errorf("The \"from\" port is malformed for port range %s: %s", port, err.Error())
		}
		to, err := strconv.Atoi(fromTo[1])
		if err != nil {
			return fmt.Errorf("The \"to\" port is malformed for port range %s: %s", port, err.Error())
		}
		if g.Ranges == nil {
			g.Ranges = []*PortRange{}
		}
		g.Ranges = append(g.Ranges, &PortRange{
			From: from,
			To:   to,
		})
	}

	return nil
}
