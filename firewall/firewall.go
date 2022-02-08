package firewall

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/frankgreco/edge-sdk-go/internal/api"
	"github.com/frankgreco/edge-sdk-go/internal/utils"
	"github.com/frankgreco/edge-sdk-go/types"

	patcher "github.com/evanphx/json-patch"
	"github.com/mattbaird/jsonpatch"
)

type Client interface {
	GetRuleset(context.Context, string) (*types.Ruleset, error)
	CreateRuleset(context.Context, *types.Ruleset) (*types.Ruleset, error)
	UpdateRuleset(context.Context, *types.Ruleset, []jsonpatch.JsonPatchOperation) (*types.Ruleset, error)
	DeleteRuleset(context.Context, string) error

	CreateAddressGroup(context.Context, *types.AddressGroup) (*types.AddressGroup, error)
	GetAddressGroup(context.Context, string) (*types.AddressGroup, error)
	UpdateAddressGroup(context.Context, *types.AddressGroup, []jsonpatch.JsonPatchOperation) (*types.AddressGroup, error)
	DeleteAddressGroup(context.Context, string) error

	CreatePortGroup(context.Context, *types.PortGroup) (*types.PortGroup, error)
	GetPortGroup(context.Context, string) (*types.PortGroup, error)
	UpdatePortGroup(context.Context, *types.PortGroup, []jsonpatch.JsonPatchOperation) (*types.PortGroup, error)
	DeletePortGroup(context.Context, string) error
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
	p.SetCodecMode(types.CodecModeRemote)
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

	in := &api.Operation{
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

	var del *api.Delete
	{
		shouldDelete := false

		del = &api.Delete{
			Resources: api.Resources{
				Firewall: &types.Firewall{
					Rulesets: map[string]*types.Ruleset{
						current.Name: new(types.Ruleset),
					},
				},
			},
		}

		if len(discard) > 0 {
			del.Resources.Firewall.Rulesets[current.Name].Rules = discard
			shouldDelete = true
		}

		if l := current.DefaultLogging; l != nil && *l && (rs.DefaultLogging == nil || !*rs.DefaultLogging) {
			del.Resources.Firewall.Rulesets[current.Name].DefaultLogging = current.DefaultLogging
			shouldDelete = true
		}

		if (current.Description != nil) && (rs.Description == nil || *rs.Description == "") {
			del.Resources.Firewall.Rulesets[current.Name].Description = current.Description
			shouldDelete = true
		}

		if !shouldDelete {
			del = nil
		}
	}

	if del != nil {
		in.Delete = del
		in.Delete.Firewall.Rulesets[current.Name].SetOpMode(types.OpModeDelete)
	}

	if _, err := c.apiClient.Post(ctx, in); err != nil {
		return nil, err
	}
	return c.GetRuleset(ctx, current.Name)
}

func (c *client) CreateAddressGroup(ctx context.Context, g *types.AddressGroup) (*types.AddressGroup, error) {
	_, err := c.apiClient.Post(ctx, &api.Operation{
		Set: &api.Set{
			Resources: api.Resources{
				Firewall: &types.Firewall{
					Groups: &types.Groups{
						Address: map[string]*types.AddressGroup{
							g.Name: g,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return c.GetAddressGroup(ctx, g.Name)
}

func (c *client) GetAddressGroup(ctx context.Context, name string) (*types.AddressGroup, error) {
	op, err := c.apiClient.Get(ctx)
	if err != nil {
		return nil, err
	}
	return toAddressGroup(name, op)
}

func (c *client) UpdateAddressGroup(ctx context.Context, current *types.AddressGroup, patches []jsonpatch.JsonPatchOperation) (*types.AddressGroup, error) {
	var group types.AddressGroup
	if err := utils.Patch(current, &group, patches); err != nil {
		return nil, err
	}

	in := &api.Operation{
		Set: &api.Set{
			Resources: api.Resources{
				Firewall: &types.Firewall{
					Groups: &types.Groups{
						Address: map[string]*types.AddressGroup{
							current.Name: &group,
						},
					},
				},
			},
		},
	}

	var del *api.Delete
	{
		shouldDelete := false

		del = &api.Delete{
			Resources: api.Resources{
				Firewall: &types.Firewall{
					Groups: &types.Groups{
						Address: map[string]*types.AddressGroup{
							current.Name: new(types.AddressGroup),
						},
					},
				},
			},
		}

		if staleCIDRs := utils.StringSliceDiff(group.Cidrs, current.Cidrs); len(staleCIDRs) > 0 {
			del.Firewall.Groups.Address[current.Name].Cidrs = staleCIDRs
			shouldDelete = true
		}

		if group.Description == nil || *group.Description == "" {
			del.Firewall.Groups.Address[current.Name].Description = current.Description
			shouldDelete = true
		}

		if !shouldDelete {
			del = nil
		}
	}

	if del != nil {
		in.Delete = del
	}

	if _, err := c.apiClient.Post(ctx, in); err != nil {
		return nil, err
	}
	return c.GetAddressGroup(ctx, current.Name)
}

func (c *client) DeleteAddressGroup(ctx context.Context, name string) error {
	_, err := c.apiClient.Post(ctx, &api.Operation{
		Delete: &api.Delete{
			Resources: api.Resources{
				Firewall: &types.Firewall{
					Groups: &types.Groups{
						Address: map[string]*types.AddressGroup{
							name: nil,
						},
					},
				},
			},
		},
	})
	return err
}

func (c *client) CreatePortGroup(ctx context.Context, g *types.PortGroup) (*types.PortGroup, error) {
	_, err := c.apiClient.Post(ctx, &api.Operation{
		Set: &api.Set{
			Resources: api.Resources{
				Firewall: &types.Firewall{
					Groups: &types.Groups{
						Port: map[string]*types.PortGroup{
							g.Name: g,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return c.GetPortGroup(ctx, g.Name)
}

func (c *client) GetPortGroup(ctx context.Context, name string) (*types.PortGroup, error) {
	op, err := c.apiClient.Get(ctx)
	if err != nil {
		return nil, err
	}
	return toPortGroup(name, op)
}

func (c *client) UpdatePortGroup(ctx context.Context, current *types.PortGroup, patches []jsonpatch.JsonPatchOperation) (*types.PortGroup, error) {
	if len(patches) == 0 {
		return current, nil
	}

	var group types.PortGroup
	if err := utils.Patch(current, &group, patches); err != nil {
		return nil, err
	}

	in := &api.Operation{
		Set: &api.Set{
			Resources: api.Resources{
				Firewall: &types.Firewall{
					Groups: &types.Groups{
						Port: map[string]*types.PortGroup{
							current.Name: &group,
						},
					},
				},
			},
		},
	}

	var del *api.Delete
	{
		shouldDelete := false

		del = &api.Delete{
			Resources: api.Resources{
				Firewall: &types.Firewall{
					Groups: &types.Groups{
						Port: map[string]*types.PortGroup{
							current.Name: new(types.PortGroup),
						},
					},
				},
			},
		}

		if stalePorts := utils.IntSliceDiff(group.Ports, current.Ports); len(stalePorts) > 0 {
			del.Firewall.Groups.Port[current.Name].Ports = stalePorts
			shouldDelete = true
		}

		if group.Description == nil || *group.Description == "" {
			del.Firewall.Groups.Port[current.Name].Description = current.Description
			shouldDelete = true
		}

		if !shouldDelete {
			del = nil
		}
	}

	if del != nil {
		in.Delete = del
	}

	if _, err := c.apiClient.Post(ctx, in); err != nil {
		return nil, err
	}
	return c.GetPortGroup(ctx, current.Name)
}

func (c *client) DeletePortGroup(ctx context.Context, name string) error {
	_, err := c.apiClient.Post(ctx, &api.Operation{
		Delete: &api.Delete{
			Resources: api.Resources{
				Firewall: &types.Firewall{
					Groups: &types.Groups{
						Port: map[string]*types.PortGroup{
							name: nil,
						},
					},
				},
			},
		},
	})
	return err
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
	// ruleset.ID = ruleset.Name
	return ruleset, nil
}

func toAddressGroup(name string, op *api.Operation) (*types.AddressGroup, error) {
	if op == nil || op.Get == nil || op.Get.Firewall == nil || op.Get.Firewall.Groups == nil || op.Get.Firewall.Groups.Address == nil {
		return nil, errors.New("No address groups exist.")
	}

	group, ok := op.Get.Firewall.Groups.Address[name]
	if !ok || group == nil {
		return nil, errors.New("The address group does not exist.")
	}

	group.Name = name
	// group.ID = group.Name
	return group, nil
}

func toPortGroup(name string, op *api.Operation) (*types.PortGroup, error) {
	if op == nil || op.Get == nil || op.Get.Firewall == nil || op.Get.Firewall.Groups == nil || op.Get.Firewall.Groups.Port == nil {
		return nil, errors.New("No port groups exist.")
	}

	group, ok := op.Get.Firewall.Groups.Port[name]
	if !ok || group == nil {
		return nil, errors.New("The port group does not exist.")
	}

	group.Name = name
	// group.ID = group.Name
	return group, nil
}
