package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSourceUnmarshalJSON(t *testing.T) {
	for _, test := range []struct {
		name     string
		expected *Source
		json     string
	}{
		{
			name: "null json value",
			expected: &Source{
				AddressGroup: "example",
				Port: &Port{
					FromPort: 80,
					ToPort:   80,
				},
			},
			json: `{"group": {"address-group": "example"}, "port": "80"}`,
		},
	} {
		source := new(Source)
		require.NoError(t, source.UnmarshalJSON([]byte(test.json)), test.name)
		require.Equal(t, test.expected, source, test.name)
	}
}

func TestRuleUnmarshalJSON(t *testing.T) {
	for _, test := range []struct {
		name     string
		expected *Rule
		json     string
	}{
		{
			name: "happy path",
			expected: &Rule{
				Description: "rule description 5",
				Action:      "drop",
				Protocol:    "tcp",
				Destination: &Destination{
					AddressGroup: "example",
					Port: &Port{
						FromPort: 80,
						ToPort:   80,
					},
				},
			},
			json: `{"action": "drop", "description": "rule description 5", "destination": {"group": {"address-group": "example"}, "port": "80"}, "protocol": "tcp", "source": null, "state": null}`,
		},
	} {
		rule := new(Rule)
		require.NoError(t, rule.UnmarshalJSON([]byte(test.json)), test.name)
		require.Equal(t, test.expected, rule, test.name)
	}
}
