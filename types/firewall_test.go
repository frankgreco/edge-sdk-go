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
