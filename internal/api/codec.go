package api

import (
	"encoding/json"
)

type status struct {
	Success string `json:"success,omitempty"`
	Failure string `json:"failure,omitempty"`
}

func (out *Status) UnmarshalJSON(data []byte) error {
	var s status

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if out == nil {
		out = new(Status)
	}

	out.Success = s.Success == "0"
	out.Failure = s.Failure == "1"

	return nil
}
