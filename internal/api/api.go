package api

import (
	"github.com/frankgreco/edge-sdk-go/types"
)

type Get struct {
	Resources
}

type Set struct {
	Resources
	*Status
}

type Delete struct {
	Resources
	*Status
}

type Resources struct {
	Firewall   *types.Firewall   `json:"firewall,omitempty"`
	Interfaces *types.Interfaces `json:"interfaces,omitempty"`
}

type Commit struct {
	*Status
}

type Status struct {
	Success bool
	Failure bool
}

type Save struct {
	Success string `json:"success,omitempty"`
}

type Operation struct {
	Success bool    `json:"success,omitempty"`
	Get     *Get    `json:"GET,omitempty"`
	Set     *Set    `json:"SET,omitempty"`
	Delete  *Delete `json:"DELETE,omitempty"`
	Commit  *Commit `json:"COMMIT,omitempty"`
	Save    *Save   `json:"SAVE,omitempty"`
}

func (op Operation) Failed() bool {
	if !op.Success {
		return true
	}

	if op.Set != nil && op.Set.Status != nil && op.Set.Failure {
		return true
	}

	if op.Commit != nil && op.Commit.Status != nil && op.Commit.Failure {
		return true
	}

	if op.Delete != nil && op.Delete.Status != nil && op.Delete.Failure {
		return true
	}

	return false
}
