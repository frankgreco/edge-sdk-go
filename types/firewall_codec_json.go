package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	enable = "enable"
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
	value := func(b bool) string {
		if b {
			return "enable"
		}
		return ""
	}

	return json.Marshal(&apiState{
		Established: value(s.Established),
		Invalid:     value(s.Invalid),
		New:         value(s.New),
		Related:     value(s.Related),
	})
}

func (s *State) UnmarshalJSON(data []byte) error {
	var state apiState

	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	s.Established = state.Established == enable
	s.Invalid = state.Invalid == enable
	s.New = state.New == enable
	s.Related = state.Related == enable

	return nil
}

func (s *Source) MarshalJSON() ([]byte, error) {
	type Alias Source
	return json.Marshal(&struct {
		Port  string `json:"port"`
		Group *group `json:"group,omitempty"`
		*Alias
	}{
		Port: s.port(),
		Group: &group{
			Address: s.AddressGroup,
		},
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
	}

	if aux.Port != "" {
		s.Port = new(Port)
		s.FromPort, s.ToPort, err = ports(aux.Port)
		if err != nil {
			return fmt.Errorf("Error setting source ports from json `%s`: %s", string(data), err.Error())
		}
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
	}

	if aux.Port != "" {
		d.Port = new(Port)
		d.FromPort, d.ToPort, err = ports(aux.Port)
		if err != nil {
			return fmt.Errorf("Error setting destination ports from json `%s`: %s", string(data), err.Error())
		}
	}

	return nil
}

type group struct {
	Address string `json:"address-group,omitempty"`
}

func (d *Destination) MarshalJSON() ([]byte, error) {
	type Alias Destination
	return json.Marshal(&struct {
		Port  string `json:"port,omitempty"`
		Group *group `json:"group,omitempty"`
		*Alias
	}{
		Port: d.port(),
		Group: &group{
			Address: d.AddressGroup,
		},
		Alias: (*Alias)(d),
	})
}

func (rs *Ruleset) MarshalJSON() ([]byte, error) {
	for _, rule := range rs.Rules {
		rule.SetCodecMode(rs.codecMode)
	}

	var isDelete bool
	{
		if rs.opMode == OpModeDelete {
			isDelete = true
		}
	}

	var data interface{}
	{
		type Alias Ruleset
		if rs.codecMode == CodecModeLocal {
			data = &struct {
				Rules []*Rule `json:"rule,omitempty"`
				*Alias
			}{
				Rules: rs.Rules,
				Alias: (*Alias)(rs),
			}
		} else {
			data = &struct {
				RulesMap map[string]*Rule `json:"rule,omitempty"`
				*Alias
			}{
				RulesMap: buildMap(rs, isDelete),
				Alias:    (*Alias)(rs),
			}
		}
	}
	return json.Marshal(data)
}

func (rs *Ruleset) UnmarshalJSON(data []byte) (err error) {
	type Alias Ruleset

	if rs.codecMode == CodecModeLocal {
		aux := &struct {
			Rules []*Rule `json:"rule,omitempty"`
			*Alias
		}{
			Alias: (*Alias)(rs),
		}
		if err := json.Unmarshal(data, &aux); err != nil {
			return err
		}
		rs.Rules = aux.Rules
		return nil
	}

	aux := &struct {
		RulesMap map[string]*Rule `json:"rule,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(rs),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

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
				Priority int `json:"priority"`
				*Alias
			}{
				Priority: r.Priority,
				Alias:    (*Alias)(r),
			}
		} else {
			if r.Protocol == "*" {
				r.Protocol = ""
			}
			data = &struct {
				*Alias
			}{
				Alias: (*Alias)(r),
			}
		}
	}
	return json.Marshal(data)
}

func (r *Rule) UnmarshalJSON(data []byte) error {
	type Alias Rule
	aux := &struct {
		Priority int `json:"priority"`
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
	return nil
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
