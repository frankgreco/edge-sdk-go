package api

import (
	"testing"

	"github.com/frankgreco/edge-sdk-go/types"

	"github.com/stretchr/testify/require"
)

func TestOperationUnmarshalJSON(t *testing.T) {
	for _, test := range []struct {
		name      string
		expected  *Operation
		json      string
		justError bool
	}{
		{
			name:      "just error",
			justError: true,
			expected: &Operation{
				justError: true,
				Set: &Set{
					Status: Status{
						Failure: false,
						Success: true,
					},
				},
				Commit: &Commit{
					Status: Status{
						Failure: true,
						Success: false,
						Error:   "Configuration system temporarily locked due to another commit in progress\n",
					},
				},
				Save: &Save{
					Status: Status{
						Failure: false,
						Success: true,
					},
				},
				Success: true,
			},
			json: `{"SET": {"failure": "0", "success": "1"}, "SESSION_ID": "ef1585d01f324c02bcfd4fb0f03e55db", "GET": {"firewall": {"group": {"address-group": {"router": ""}}}}, "COMMIT": {"error": "Configuration system temporarily locked due to another commit in progress\n", "failure": "1", "success": "0"}, "SAVE": {"success": "1"}, "success": true}`,
		},
		{
			name: "no error",
			expected: &Operation{
				Get: &Get{
					Resources: Resources{
						Firewall: &types.Firewall{
							Groups: &types.Groups{
								Address: map[string]*types.AddressGroup{
									"router": {
										Description: strptr("router interface addresses"),
										Cidrs: []string{
											"192.168.2.1",
											"192.168.3.1",
											"192.168.4.1",
										},
									},
								},
							},
						},
					},
				},
				Success: true,
			},
			json: `{"GET": {"firewall": {"group": {"address-group": {"router": {"address": ["192.168.2.1","192.168.3.1","192.168.4.1"], "description": "router interface addresses"}}}}}, "success": true}`,
		},
	} {
		operation := new(Operation)
		if test.justError {
			operation.justError = true
		}
		require.NoError(t, operation.UnmarshalJSON([]byte(test.json)), test.name)
		require.Equal(t, test.expected, operation, test.name)
	}
}

func TestCommitUnmarshalJSON(t *testing.T) {
	for _, test := range []struct {
		name     string
		expected *Commit
		json     string
	}{
		{
			name: "has error",
			expected: &Commit{
				Status: Status{
					Success: false,
					Failure: true,
					Error:   "Configuration system temporarily locked due to another commit in progress\n",
				},
			},
			json: `{"error": "Configuration system temporarily locked due to another commit in progress\n", "failure": "1", "success": "0"}`,
		},
	} {
		commit := new(Commit)
		require.NoError(t, commit.UnmarshalJSON([]byte(test.json)), test.name)
		require.Equal(t, test.expected, commit, test.name)
	}
}

func TestStatusUnmarshalJSON(t *testing.T) {
	for _, test := range []struct {
		name     string
		expected *Status
		json     string
	}{
		{
			name: "should be successfull if every field is set as expected",
			expected: &Status{
				Success: true,
				Failure: false,
			},
			json: `{"success": "1", "failure": "0"}`,
		},
		{
			name: "should be successfull if only success is set",
			expected: &Status{
				Success: true,
				Failure: false,
			},
			json: `{"success": "1"}`,
		},
		{
			name: "should be successfull if no fields are set",
			expected: &Status{
				Success: true,
				Failure: false,
			},
			json: `{}`,
		},
		{
			name: "should be successfull if fields are empty",
			expected: &Status{
				Success: true,
				Failure: false,
			},
			json: `{"success": "", "failure": ""}`,
		},
		{
			name: "should be failed if fields are set",
			expected: &Status{
				Success: false,
				Failure: true,
			},
			json: `{"success": "0", "failure": "1"}`,
		},
	} {
		status := new(Status)
		require.NoError(t, status.UnmarshalJSON([]byte(test.json)), test.name)
		require.Equal(t, test.expected, status, test.name)
	}
}

func strptr(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}
