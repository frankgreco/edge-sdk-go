package interfaces

import (
	"github.com/frankgreco/edge-sdk-go/api/interfaces/ethernet"
)

type Interfaces struct {
	Ethernet map[string]*ethernet.Ethernet `json:"ethernet,omitempty"`
}
