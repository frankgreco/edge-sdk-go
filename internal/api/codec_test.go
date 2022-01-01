package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
			name: "should be successfull if fields are ampty",
			expected: &Status{
				Success: true,
				Failure: false,
			},
			json: `{"success": "", "failure": ""}`,
		},
	} {
		status := new(Status)
		require.NoError(t, status.UnmarshalJSON([]byte(test.json)), test.name)
		require.Equal(t, test.expected, status, test.name)
	}
}
