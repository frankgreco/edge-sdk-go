package firewall

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/frankgreco/edge-sdk-go/api/firewall"
	"github.com/frankgreco/edge-sdk-go/internal/api"

	patcher "github.com/evanphx/json-patch"
	"github.com/mattbaird/jsonpatch"
)

type Client interface {
	GetRuleset(context.Context, string) (*firewall.Ruleset, error)
	CreateRuleset(context.Context, *firewall.Ruleset) (*firewall.Ruleset, error)
	UpdateRuleset(context.Context, *firewall.Ruleset, []jsonpatch.JsonPatchOperation) (*firewall.Ruleset, error)
	DeleteRuleset(context.Context, string) error
}

type client struct {
	apiClient api.Client
}

func New(httpClient *http.Client, host string) Client {
	return &client{
		apiClient: api.New(httpClient, host),
	}
}

func (c *client) GetRuleset(ctx context.Context, name string) (*firewall.Ruleset, error) {
	op, err := c.apiClient.Get(ctx)
	if err != nil {
		return nil, err
	}
	return toRuleset(name, op)
}

func (c *client) CreateRuleset(ctx context.Context, p *firewall.Ruleset) (*firewall.Ruleset, error) {
	op, err := c.apiClient.Post(ctx, &api.Operation{
		Set: &api.Set{
			Resources: api.Resources{
				Firewall: &firewall.Firewall{
					Rulesets: map[string]*firewall.Ruleset{
						p.Name: p,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return toRuleset(p.Name, op)
}

func (c *client) DeleteRuleset(ctx context.Context, name string) error {
	op, err := c.apiClient.Post(ctx, &api.Operation{
		Delete: &api.Delete{
			Resources: api.Resources{
				Firewall: &firewall.Firewall{
					Rulesets: map[string]*firewall.Ruleset{
						name: nil,
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = toRuleset(name, op)
	return err
}

func (c *client) UpdateRuleset(ctx context.Context, current *firewall.Ruleset, patches []jsonpatch.JsonPatchOperation) (*firewall.Ruleset, error) {
	patchData, err := json.Marshal(patches)
	if err != nil {
		return nil, err
	}

	patchObj, err := patcher.DecodePatch(patchData)
	if err != nil {
		return nil, err
	}

	current.Terraform()
	currentData, err := json.Marshal(current)
	if err != nil {
		return nil, err
	}

	modifiedData, err := patchObj.Apply(currentData)
	if err != nil {
		return nil, err
	}

	var rs firewall.Ruleset
	(&rs).Terraform()
	if err := json.Unmarshal(modifiedData, &rs); err != nil {
		return nil, err
	}

	op, err := c.apiClient.Post(ctx, &api.Operation{
		Set: &api.Set{
			Resources: api.Resources{
				Firewall: &firewall.Firewall{
					Rulesets: map[string]*firewall.Ruleset{
						current.Name: &rs,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return toRuleset(current.Name, op)
}

func toRuleset(name string, op *api.Operation) (*firewall.Ruleset, error) {
	if op == nil || op.Get == nil || op.Get.Firewall == nil {
		return nil, errors.New("A firewall does not exist.")
	}

	if op.Get.Firewall.Rulesets == nil {
		return nil, errors.New("there are no rulesets")
	}

	ruleset, ok := op.Get.Firewall.Rulesets[name]
	if !ok || ruleset == nil {
		return nil, errors.New("The ruleset does not exist.")
	}

	ruleset.Name = name
	return ruleset, nil
}
