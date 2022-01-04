package api

import (
	"errors"
	"strings"

	"github.com/frankgreco/edge-sdk-go/types"
)

type Get struct {
	Resources
}

type Set struct {
	Resources
	Status `json:"-"`
}

type Delete struct {
	Resources
	Status `json:"-"`
}

type Resources struct {
	Firewall   *types.Firewall   `json:"firewall,omitempty"`
	Interfaces *types.Interfaces `json:"interfaces,omitempty"`
}

type Commit struct {
	Status `json:"-"`
}

type Status struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"-"`
	Failure bool   `json:"-"`
}

type Save struct {
	Status `json:"-"`
}

type Operation struct {
	Success   bool    `json:"success,omitempty"`
	Get       *Get    `json:"GET,omitempty"`
	Set       *Set    `json:"SET,omitempty"`
	Delete    *Delete `json:"DELETE,omitempty"`
	Commit    *Commit `json:"COMMIT,omitempty"`
	Save      *Save   `json:"SAVE,omitempty"`
	justError bool
}

func (op Operation) Failed() error {
	err := []string{}
	failed := false

	if !op.Success {
		failed = true
	}

	if op.Set != nil && op.Set.Failure {
		if op.Set.Error != "" {
			err = append(err, op.Set.Error)
		}
		failed = true
	}

	if op.Commit != nil && op.Commit.Failure {
		if op.Commit.Error != "" {
			err = append(err, op.Commit.Error)
		}
		failed = true
	}

	if op.Delete != nil && op.Delete.Failure {
		if op.Delete.Error != "" {
			err = append(err, op.Delete.Error)
		}
		failed = true
	}

	if !failed {
		return nil
	}

	if len(err) == 0 {
		err = append(err, "The operation failed for a unknown reason.")
	}

	return errors.New(strings.Join(err, ", "))
}
