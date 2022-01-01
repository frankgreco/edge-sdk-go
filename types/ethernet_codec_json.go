package types

import (
	"encoding/json"
	"strconv"
)

type apiFirewallDetails struct {
	Name string `json:"name"`
}

type apiFirewall struct {
	In    *apiFirewallDetails `json:"in"`
	Out   *apiFirewallDetails `json:"out"`
	Local *apiFirewallDetails `json:"local"`
}

func (f *FirewallAttachment) MarshalJSON() ([]byte, error) {
	af := new(apiFirewall)

	if f.In != nil && *f.In != "" {
		af.In = &apiFirewallDetails{
			Name: *f.In,
		}
	}
	if f.Out != nil && *f.Out != "" {
		af.Out = &apiFirewallDetails{
			Name: *f.Out,
		}
	}
	if f.Local != nil && *f.Local != "" {
		af.Local = &apiFirewallDetails{
			Name: *f.Local,
		}
	}

	return json.Marshal(af)
}

func (f *FirewallAttachment) UnmarshalJSON(data []byte) (err error) {
	var ap apiFirewall

	if err := json.Unmarshal(data, &ap); err != nil {
		return err
	}

	if ap.In != nil && ap.In.Name != "" {
		f.In = &ap.In.Name
	}
	if ap.Out != nil && ap.Out.Name != "" {
		f.Out = &ap.Out.Name
	}
	if ap.Local != nil && ap.Local.Name != "" {
		f.Local = &ap.Local.Name
	}

	return nil
}

func (o *DHCPOptions) MarshalJSON() ([]byte, error) {
	type Alias DHCPOptions
	return json.Marshal(&struct {
		DefaultRouteDistance string `json:"default-route-distance,omitempty"`
		*Alias
	}{
		DefaultRouteDistance: strconv.Itoa(o.DefaultRouteDistance),
		Alias:                (*Alias)(o),
	})
}

func (o *DHCPOptions) UnmarshalJSON(data []byte) (err error) {
	type Alias DHCPOptions
	aux := &struct {
		DefaultRouteDistance string `json:"default-route-distance,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	i, err := strconv.Atoi(aux.DefaultRouteDistance)
	if err != nil {
		return err
	}
	o.DefaultRouteDistance = i
	return nil
}
