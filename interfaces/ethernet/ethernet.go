package ethernet

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/frankgreco/edge-sdk-go/internal/api"
	"github.com/frankgreco/edge-sdk-go/types"
)

type Client interface {
	Get(context.Context, string) (*types.Ethernet, error)
	AttachFirewallRuleset(context.Context, string, *types.FirewallAttachment) (*types.FirewallAttachment, error)
	DetachFirewallRuleset(context.Context, string) error
}

type client struct {
	apiClient api.Client
}

func New(httpClient *http.Client, host string) Client {
	return &client{
		apiClient: api.New(httpClient, host),
	}
}

func (c *client) Get(ctx context.Context, id string) (*types.Ethernet, error) {
	op, err := c.apiClient.Get(ctx)
	if err != nil {
		return nil, err
	}
	return toEthernet(id, op)
}

func (c *client) AttachFirewallRuleset(ctx context.Context, id string, firewall *types.FirewallAttachment) (*types.FirewallAttachment, error) {
	op, err := c.apiClient.Post(ctx, &api.Operation{
		Set: &api.Set{
			Resources: api.Resources{
				Interfaces: &types.Interfaces{
					Ethernet: map[string]*types.Ethernet{
						id: {
							Firewall: firewall,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	ethernet, err := toEthernet(id, op)
	if err != nil {
		return nil, err
	}
	return ethernet.Firewall, nil // TODO: Potentially return error if no firewalls are attached.
}

func (c *client) DetachFirewallRuleset(ctx context.Context, id string) error {
	return nil
}

func toEthernet(id string, op *api.Operation) (*types.Ethernet, error) {
	if op == nil || op.Get == nil || op.Get.Interfaces == nil || op.Get.Interfaces.Ethernet == nil {
		return nil, errors.New("No ethernet interfaces exist.")
	}

	ethernet, ok := op.Get.Interfaces.Ethernet[id]
	if !ok {
		return nil, fmt.Errorf("The ethernet interface %s does not exist.", id)
	}

	ethernet.Firewall.Interface = id

	return ethernet, nil
}
