package ethernet

import (
	"encoding/json"
)

type apiFirewallDetails struct {
	Name string `json:"name,omitempty"`
}

type apiFirewall struct {
	In    *apiFirewallDetails `json:"in,omitempty"`
	Out   *apiFirewallDetails `json:"out,omitempty"`
	Local *apiFirewallDetails `json:"local,omitempty"`
}

func (f *Firewall) MarshalJSON() ([]byte, error) {
	af := new(apiFirewall)

	if f.In != "" {
		af.In = &apiFirewallDetails{
			Name: f.In,
		}
	}
	if f.Out != "" {
		af.Out = &apiFirewallDetails{
			Name: f.Out,
		}
	}
	if f.Local != "" {
		af.Local = &apiFirewallDetails{
			Name: f.Local,
		}
	}

	return json.Marshal(af)
}

func (f *Firewall) UnmarshalJSON(data []byte) (err error) {
	var ap apiFirewall

	if err := json.Unmarshal(data, &ap); err != nil {
		return err
	}

	if ap.In != nil {
		f.In = ap.In.Name
	}
	if ap.Out != nil {
		f.Out = ap.Out.Name
	}
	if ap.Local != nil {
		f.Local = ap.Local.Name
	}

	return nil
}
