package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	enable  = "enable"
	disable = "disable"
)

type CodecMode int

const (
	CodecModeRemote CodecMode = iota
	CodecModeLocal
)

type apiState struct {
	Established string `json:"established"`
	Invalid     string `json:"invalid"`
	New         string `json:"new"`
	Related     string `json:"related"`
}

func (s *State) MarshalJSON() ([]byte, error) {
	value := func(b *bool) string {
		if b == nil || !*b {
			return disable
		}
		return enable
	}

	return json.Marshal(&apiState{
		Established: value(s.Established),
		Invalid:     value(s.Invalid),
		New:         value(s.New),
		Related:     value(s.Related),
	})
}

func (s *State) UnmarshalJSON(data []byte) error {
	t := true

	value := func(b string) *bool {
		if b == enable {
			return &t
		}
		return nil
	}

	var state apiState

	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	s.Established = value(state.Established)
	s.Invalid = value(state.Invalid)
	s.New = value(state.New)
	s.Related = value(state.Related)

	return nil
}

func (s *Source) MarshalJSON() ([]byte, error) {
	var g *group
	{
		if s.AddressGroup != nil || s.PortGroup != nil {
			g = &group{
				Address: s.AddressGroup,
				Port:    s.PortGroup,
			}
		}
	}

	type Alias Source
	return json.Marshal(&struct {
		Port  string `json:"port,omitempty"`
		Group *group `json:"group,omitempty"`
		*Alias
	}{
		Port:  s.toPort(),
		Group: g,
		Alias: (*Alias)(s),
	})
}

func (s *Source) UnmarshalJSON(data []byte) (err error) {
	type Alias Source
	aux := &struct {
		Port  string `json:"port"`
		Group *group `json:"group,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Group != nil {
		s.AddressGroup = aux.Group.Address
		s.PortGroup = aux.Group.Port
	}

	if err := s.fromPort(aux.Port); err != nil {
		return fmt.Errorf("Error setting source ports %s from json `%s`: %s", aux.Port, string(data), err.Error())
	}

	return nil
}

func (d *Destination) UnmarshalJSON(data []byte) (err error) {
	type Alias Destination
	aux := &struct {
		Port  string `json:"port,omitempty"`
		Group *group `json:"group,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Group != nil {
		d.AddressGroup = aux.Group.Address
		d.PortGroup = aux.Group.Port
	}

	if err := d.fromPort(aux.Port); err != nil {
		return fmt.Errorf("Error setting destination ports %s from json `%s`: %s", aux.Port, string(data), err.Error())
	}

	return nil
}

func (d *Destination) MarshalJSON() ([]byte, error) {
	var g *group
	{
		if d.AddressGroup != nil || d.PortGroup != nil {
			g = &group{
				Address: d.AddressGroup,
				Port:    d.PortGroup,
			}
		}
	}

	type Alias Destination
	return json.Marshal(&struct {
		Port  string `json:"port,omitempty"`
		Group *group `json:"group,omitempty"`
		*Alias
	}{
		Port:  d.toPort(),
		Group: g,
		Alias: (*Alias)(d),
	})
}

type group struct {
	Address *string `json:"address-group,omitempty"`
	Port    *string `json:"port-group,omitempty"`
}

func (g *PortGroup) UnmarshalJSON(data []byte) (err error) {
	type Alias PortGroup
	aux := &struct {
		Ports []string `json:"port,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(g),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.Ports) > 0 {
		if err := g.fromPorts(aux.Ports); err != nil {
			return fmt.Errorf("Error setting ports on port group from json `%s`: %s", string(data), err.Error())
		}
	}

	return nil
}

func (g PortGroup) MarshalJSON() ([]byte, error) {
	return (&g).marshalJSON()
}

func (g *PortGroup) marshalJSON() ([]byte, error) {
	type Alias PortGroup
	return json.Marshal(&struct {
		Ports []string `json:"port,omitempty"`
		*Alias
	}{
		Ports: g.toPorts(),
		Alias: (*Alias)(g),
	})
}

func (rs Ruleset) MarshalJSON() ([]byte, error) {
	return (&rs).marshalJSON()
}

func (rs *Ruleset) marshalJSON() ([]byte, error) {
	for _, rule := range rs.Rules {
		rule.SetCodecMode(rs.codecMode)
	}

	var isDelete bool
	{
		if rs.opMode == OpModeDelete {
			isDelete = true
		}
	}

	var n *null
	{
		if l := rs.DefaultLogging; l != nil && *l {
			n = &null{true}
		}
	}

	var data interface{}
	{
		type Alias Ruleset
		if rs.codecMode == CodecModeLocal {
			data = &struct {
				DefaultLogging *null   `json:"enable-default-log,omitempty"`
				Rules          []*Rule `json:"rule,omitempty"`
				*Alias
			}{
				DefaultLogging: n,
				Rules:          rs.Rules,
				Alias:          (*Alias)(rs),
			}
		} else {
			data = &struct {
				DefaultLogging *null            `json:"enable-default-log,omitempty"`
				RulesMap       map[string]*Rule `json:"rule,omitempty"`
				*Alias
			}{
				DefaultLogging: n,
				RulesMap:       buildMap(rs, isDelete),
				Alias:          (*Alias)(rs),
			}
		}
	}
	return json.Marshal(data)
}

func (rs *Ruleset) UnmarshalJSON(data []byte) (err error) {
	type Alias Ruleset

	if rs.codecMode == CodecModeLocal {
		aux := &struct {
			DefaultLogging null    `json:"enable-default-log"`
			Rules          []*Rule `json:"rule,omitempty"`
			*Alias
		}{
			Alias: (*Alias)(rs),
		}
		if err := json.Unmarshal(data, &aux); err != nil {
			return err
		}

		rs.DefaultLogging = boolptr(aux.DefaultLogging.val)
		rs.Rules = aux.Rules
		return nil
	}

	aux := &struct {
		DefaultLogging null             `json:"enable-default-log"`
		RulesMap       map[string]*Rule `json:"rule,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(rs),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	rs.DefaultLogging = boolptr(aux.DefaultLogging.val)

	for k, v := range aux.RulesMap {
		i, err := strconv.Atoi(k)
		if err != nil {
			return fmt.Errorf("malformed rule priority: %v", k)
		}
		v.Priority = i
		rs.Rules = append(rs.Rules, v)
	}
	return nil
}

func (r *Rule) MarshalJSON() ([]byte, error) {
	var data interface{}
	{
		type Alias Rule
		if r.codecMode == CodecModeLocal {
			data = &struct {
				Priority int    `json:"priority"`
				Log      string `json:"log,omitempty"`
				*Alias
			}{
				Priority: r.Priority,
				Log:      toEnableDisable(r.Log),
				Alias:    (*Alias)(r),
			}
		} else {
			if r.Protocol == "*" {
				r.Protocol = ""
			}
			data = &struct {
				Log string `json:"log,omitempty"`
				*Alias
			}{
				Log:   toEnableDisable(r.Log),
				Alias: (*Alias)(r),
			}
		}
	}
	return json.Marshal(data)
}

func (r *Rule) UnmarshalJSON(data []byte) error {
	type Alias Rule
	aux := &struct {
		Priority int    `json:"priority"`
		Log      string `json:"log,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	r.Priority = aux.Priority

	if aux.Protocol == "" {
		r.Protocol = "*"
	}
	r.Log = toBoolPtr(aux.Log)
	return nil
}

func toBoolPtr(val string) *bool {
	t, f := true, false

	if val == "" || val == disable {
		return &f
	}
	return &t
}

func toEnableDisable(b *bool) string {
	if b == nil || !*b {
		return disable
	}
	return enable
}

// // consider having
// // type ruleMap map[string]*Rule
// // and having a MarshalJSON for that instead.
func buildMap(rs *Ruleset, isDelete bool) map[string]*Rule {
	if rs == nil || len(rs.Rules) == 0 {
		return nil
	}

	m := map[string]*Rule{}
	for _, rule := range rs.Rules {
		if isDelete {
			m[strconv.Itoa(rule.Priority)] = nil
		} else {
			m[strconv.Itoa(rule.Priority)] = rule
		}
	}
	return m
}

func boolptr(b bool) *bool {
	return &b
}
