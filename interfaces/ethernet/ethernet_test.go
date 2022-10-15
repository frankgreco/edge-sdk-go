package ethernet

import (
	"encoding/json"
	"github.com/frankgreco/edge-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToEthernetMapMarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		id       string
		firewall *types.FirewallAttachment
		expected string
	}{
		{
			name:     "no vlan and no firewall",
			id:       "eth0",
			expected: `{"eth0":{}}`,
		},
		{
			name:     "vlan and no firewall",
			id:       "eth0.20",
			expected: `{"eth0":{"vif":{"20":{}}}}`,
		},
		{
			name:     "no vlan and a firewall",
			id:       "eth0",
			firewall: &types.FirewallAttachment{In: strptr("eth0")},
			expected: `{"eth0":{"firewall":{"in":{"name":"eth0"}}}}`,
		},
		{
			name:     "vlan and a firewall",
			id:       "eth0.20",
			firewall: &types.FirewallAttachment{In: strptr("eth0.20")},
			expected: `{"eth0":{"vif":{"20":{"firewall":{"in":{"name":"eth0.20"}}}}}}`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			out, err := json.Marshal(toEthernetMap(tt.id, tt.firewall))
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(out))
		})
	}
}

func strptr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
