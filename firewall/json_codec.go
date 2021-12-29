package firewall

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	enable = "enable"
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
		Port string `json:"port"`
		*Alias
	}{
		Port:  s.port(),
		Alias: (*Alias)(s),
	})
}

func (s *Source) UnmarshalJSON(data []byte) (err error) {
	type Alias Source
	aux := &struct {
		Port string `json:"port"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.Port = new(Port)
	s.FromPort, s.ToPort, err = ports(aux.Port)
	return err
}

func (d *Destination) MarshalJSON() ([]byte, error) {
	type Alias Destination
	return json.Marshal(&struct {
		Port string `json:"port"`
		*Alias
	}{
		Port:  d.port(),
		Alias: (*Alias)(d),
	})
}

func (d *Destination) UnmarshalJSON(data []byte) (err error) {
	type Alias Destination
	aux := &struct {
		Port string `json:"port"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	d.Port = new(Port)
	d.FromPort, d.ToPort, err = ports(aux.Port)
	return err
}

func (r *Rule) MarshalJSON() ([]byte, error) {
	type Alias Rule
	if r.isTerraform {
		return json.Marshal(&struct{
			Priority int `json:"priority"`
			*Alias
		}{
			Priority: r.Priority,
			Alias:    (*Alias)(r),
		})
	}
	return json.Marshal(&struct{
		*Alias
	}{
		Alias: (*Alias)(r),
	})
	return json.Marshal(r)
}

func (rs *Ruleset) MarshalJSON() ([]byte, error) {
	if rs != nil && rs.isTerraform {
		for _, rule := range rs.Rules {
			rule.Terraform()
		}
	}

	type Alias Ruleset
	if rs.isTerraform {
		return json.Marshal(&struct{
			Rules []*Rule `json:"rule"`
			*Alias
		}{
			Rules: rs.Rules,
			Alias: (*Alias)(rs),
		})
	}
	return json.Marshal(&struct{
		RulesMap map[string]*Rule `json:"rule,omitempty"`
		*Alias
	}{
		RulesMap: rs.buildMap(),
		Alias:    (*Alias)(rs),
	})
}

func (rs *Ruleset) UnmarshalJSON(data []byte) (err error) {
	type Alias Ruleset
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

func (r *Rule) UnmarshalJSON(data []byte) error {
	type Alias Rule
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Protocol == "" {
		r.Protocol = "*"
	}

	return nil
}
