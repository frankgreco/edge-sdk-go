package firewall

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
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

func (d *Destination) FromTerraform5Value(v tftypes.Value) error {
	if !v.IsKnown() {
		return errors.New("The provided value is unknown. This is an issue with the Terraform SDK.")
	}

	if v.IsNull() {
		return nil
	}

	d.Port = new(Port)

	// medium == "the 'medium in which terraform is using to plumb values to us'"
	medium := map[string]tftypes.Value{}
	if err := v.As(&medium); err != nil {
		return err
	}

	if err := medium["address"].As(&d.Address); err != nil {
		return err
	}

	var fromPort int
	{
		port := big.NewFloat(-42)
		if err := medium["from_port"].As(&port); err != nil {
			return err
		}
		i, _ := port.Int64()
		fromPort = int(i)
	}
	d.Port.FromPort = fromPort

	var toPort int
	{
		port := big.NewFloat(-42)
		if err := medium["to_port"].As(&port); err != nil {
			return err
		}
		i, _ := port.Int64()
		toPort = int(i)
	}
	d.Port.ToPort = toPort

	return nil
}

func (d *Destination) ToTerraform5Value() (interface{}, error) {
	if d == nil {
		return nil, nil
	}

	var fromPort, toPort *int
	if d.Port != nil {
		fromPort = &d.Port.FromPort
		toPort = &d.Port.ToPort
	}

	return map[string]tftypes.Value{
		"address":   tftypes.NewValue(tftypes.String, d.Address),
		"from_port": tftypes.NewValue(tftypes.Number, fromPort),
		"to_port":   tftypes.NewValue(tftypes.Number, toPort),
	}, nil
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
