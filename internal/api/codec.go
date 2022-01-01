package api

import (
	"encoding/json"
)

type status struct {
	Success string `json:"success"`
	Failure string `json:"failure"`
}

func (out *Status) UnmarshalJSON(data []byte) error {
	var s status

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if out == nil {
		out = new(Status)
	}

	out.Success = s.Success == "1" || s.Failure == ""
	out.Failure = s.Failure == "1"

	return nil
}
