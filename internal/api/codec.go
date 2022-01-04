package api

import (
	"encoding/json"
)

type operation struct {
	Success bool    `json:"success,omitempty"`
	Set     *Status `json:"SET,omitempty"`
	Delete  *Status `json:"DELETE,omitempty"`
	Commit  *Status `json:"COMMIT,omitempty"`
	Save    *Status `json:"SAVE,omitempty"`
}

func (o *Operation) UnmarshalJSON(data []byte) error {
	if !o.justError {
		type Alias Operation
		return json.Unmarshal(data, &struct {
			*Alias
		}{
			Alias: (*Alias)(o),
		})
	}

	var op operation
	if err := json.Unmarshal(data, &op); err != nil {
		return err
	}
	o.Success = op.Success
	if op.Set != nil {
		o.Set = &Set{Status: *op.Set}
	}
	if op.Delete != nil {
		o.Delete = &Delete{Status: *op.Delete}
	}
	if op.Commit != nil {
		o.Commit = &Commit{Status: *op.Commit}
	}
	if op.Save != nil {
		o.Save = &Save{Status: *op.Save}
	}

	return nil
}

func (s *Status) UnmarshalJSON(data []byte) error {
	type Alias Status
	aux := &struct {
		Success string `json:"success"`
		Failure string `json:"failure"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	s.Success = aux.Success == "1" || aux.Failure == ""
	s.Failure = aux.Failure == "1"

	return nil
}
