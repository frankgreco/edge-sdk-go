package interfaces

import (
	"net/http"

	"github.com/frankgreco/edge-sdk-go/interfaces/ethernet"
)

type Client struct {
	Ethernet ethernet.Client
	// Loopback
	// VTI
	// ...
}

func New(httpClient *http.Client, baseURL string) *Client {
	return &Client{
		Ethernet: ethernet.New(httpClient, baseURL),
	}
}