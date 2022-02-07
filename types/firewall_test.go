package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDestinationFromPort(t *testing.T) {
	for _, test := range []struct {
		name     string
		expected *Destination
		port     string
	}{
		{
			name: "happy path",
			expected: &Destination{
				Port: &PortRange{
					From: 80,
					To:   80,
				},
			},
			port: "80",
		},
	} {
		destination := new(Destination)
		require.NoError(t, destination.fromPort(test.port), test.name)
		require.Equal(t, test.expected, destination, test.name)
	}
}

func TestRulesetJSONMarshal(t *testing.T) {
	for _, test := range []struct {
		name     string
		rs       *Ruleset
		expected string
	}{
		{
			name: "default logging true",
			rs: &Ruleset{
				Name:           "name",
				Description:    strptr("description"),
				DefaultAction:  "accept",
				DefaultLogging: boolptr(true),
				Rules: []*Rule{
					{
						Priority:    10,
						Description: strptr("rule description"),
						Action:      "accept",
						Protocol:    "all",
						Log:         boolptr(false),
					},
				},
				codecMode: CodecModeLocal,
			},
			expected: `{"enable-default-log":null,"rule":[{"priority":10,"log":"disable","description":"rule description","action":"accept","protocol":"all","source":null,"destination":null,"state":null}],"description":"description","default-action":"accept"}`,
		},
		{
			name: "default logging false",
			rs: &Ruleset{
				Name:           "name",
				Description:    strptr("description"),
				DefaultAction:  "accept",
				DefaultLogging: boolptr(false),
				Rules: []*Rule{
					{
						Priority:    10,
						Description: strptr("rule description"),
						Action:      "accept",
						Protocol:    "all",
						Log:         boolptr(false),
					},
				},
				codecMode: CodecModeLocal,
			},
			expected: `{"rule":[{"priority":10,"log":"disable","description":"rule description","action":"accept","protocol":"all","source":null,"destination":null,"state":null}],"description":"description","default-action":"accept"}`,
		},
		{
			name: "default logging false with remote codec",
			rs: &Ruleset{
				Name:           "name",
				Description:    strptr("description"),
				DefaultAction:  "accept",
				DefaultLogging: boolptr(false),
				Rules: []*Rule{
					{
						Priority:    10,
						Description: strptr("rule description"),
						Action:      "accept",
						Protocol:    "all",
						Log:         boolptr(false),
					},
				},
				codecMode: CodecModeRemote,
			},
			expected: `{"rule":{"10":{"log":"disable","description":"rule description","action":"accept","protocol":"all","source":null,"destination":null,"state":null}},"description":"description","default-action":"accept"}`,
		},
	} {
		data, err := test.rs.MarshalJSON()
		require.NoError(t, err, test.name)
		require.Equal(t, test.expected, string(data), test.name)
	}
}

func TestRulesetJSONUnmarshal(t *testing.T) {
	for _, test := range []struct {
		name     string
		expected *Ruleset
		data     string
		codec    CodecMode
	}{
		{
			name: "happy path",
			expected: &Ruleset{
				Description:    strptr("description"),
				DefaultAction:  "accept",
				DefaultLogging: boolptr(true),
			},
			codec: CodecModeLocal,
			data:  `{"description":"description","default-action":"accept", "enable-default-log": null}`,
		},
		{
			name: "happy path",
			expected: &Ruleset{
				Description:    strptr("description"),
				DefaultAction:  "accept",
				DefaultLogging: boolptr(false),
			},
			codec: CodecModeLocal,
			data:  `{"description":"description","default-action":"accept"}`,
		},
		{
			name: "happy path",
			expected: &Ruleset{
				Description:    strptr("description"),
				DefaultAction:  "accept",
				DefaultLogging: boolptr(true),
			},
			codec: CodecModeRemote,
			data:  `{"description":"description","default-action":"accept", "enable-default-log": null}`,
		},
		{
			name: "happy path",
			expected: &Ruleset{
				Description:    strptr("description"),
				DefaultAction:  "accept",
				DefaultLogging: boolptr(false),
			},
			codec: CodecModeRemote,
			data:  `{"description":"description","default-action":"accept"}`,
		},
	} {
		rs := new(Ruleset)
		require.NoError(t, rs.UnmarshalJSON([]byte(test.data)), test.name)
		require.Equal(t, test.expected, rs, test.name)
	}
}

func strptr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
