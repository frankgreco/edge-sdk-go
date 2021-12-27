package firewall

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/frankgreco/edge-sdk-go/internal/api"
	"github.com/frankgreco/edge-sdk-go/types"

	patcher "github.com/evanphx/json-patch"
	"github.com/mattbaird/jsonpatch"
)

type Client interface {
	GetRuleset(context.Context, string) (*types.Ruleset, error)
	CreateRuleset(context.Context, *types.Ruleset) (*types.Ruleset, error)
	UpdateRuleset(context.Context, *types.Ruleset, []jsonpatch.JsonPatchOperation) (*types.Ruleset, error)
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

func (c *client) GetRuleset(ctx context.Context, name string) (*types.Ruleset, error) {
	op, err := c.apiClient.Get(ctx)
	if err != nil {
		return nil, err
	}
	return toRuleset(name, op)
}

func (c *client) CreateRuleset(ctx context.Context, p *types.Ruleset) (*types.Ruleset, error) {
	_, err := c.apiClient.Post(ctx, &api.Operation{
		Set: &api.Set{
			Resources: api.Resources{
				Firewall: &types.Firewall{
					Rulesets: map[string]*types.Ruleset{
						p.Name: p,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return c.GetRuleset(ctx, p.Name)
}

func (c *client) DeleteRuleset(ctx context.Context, name string) error {
	_, err := c.apiClient.Post(ctx, &api.Operation{
		Delete: &api.Delete{
			Resources: api.Resources{
				Firewall: &types.Firewall{
					Rulesets: map[string]*types.Ruleset{
						name: nil,
					},
				},
			},
		},
	})
	return err
}

func (c *client) UpdateRuleset(ctx context.Context, current *types.Ruleset, patches []jsonpatch.JsonPatchOperation) (*types.Ruleset, error) {
	current.SetCodecMode(types.CodecModeLocal)

	patchData, err := json.Marshal(patches)
	if err != nil {
		return nil, err
	}

	patchObj, err := patcher.DecodePatch(patchData)
	if err != nil {
		return nil, err
	}

	currentData, err := json.Marshal(current)
	if err != nil {
		return nil, err
	}

	modifiedData, err := patchObj.Apply(currentData)
	if err != nil {
		return nil, err
	}

	var rs types.Ruleset
	{
		rs.SetCodecMode(types.CodecModeLocal)
		if err := json.Unmarshal(modifiedData, &rs); err != nil {
			return nil, err
		}
	}

	var discard []*types.Rule
	{
		modifiedRuleKeys := map[int]bool{}
		for _, rule := range rs.Rules {
			modifiedRuleKeys[rule.Priority] = true
		}
		for _, rule := range current.Rules {
			if _, ok := modifiedRuleKeys[rule.Priority]; !ok {
				discard = append(discard, &types.Rule{
					Priority: rule.Priority,
				})
			}
		}
	}

	rs.SetCodecMode(types.CodecModeRemote)

	var in *api.Operation
	{
		in = &api.Operation{
			Set: &api.Set{
				Resources: api.Resources{
					Firewall: &types.Firewall{
						Rulesets: map[string]*types.Ruleset{
							current.Name: &rs,
						},
					},
				},
			},
		}

		if len(discard) > 0 {
			in.Delete = &api.Delete{
				Resources: api.Resources{
					Firewall: &types.Firewall{
						Rulesets: map[string]*types.Ruleset{
							current.Name: {
								Rules: discard,
							},
						},
					},
				},
			}
			in.Delete.Firewall.Rulesets[current.Name].SetOpMode(types.OpModeDelete)
		}
	}

	if _, err := c.apiClient.Post(ctx, in); err != nil {
		return nil, err
	}
	return c.GetRuleset(ctx, current.Name)
}

func toRuleset(name string, op *api.Operation) (*types.Ruleset, error) {
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
