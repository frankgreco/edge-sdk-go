package edge

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/frankgreco/edge-sdk-go/firewall"
	"github.com/frankgreco/edge-sdk-go/interfaces"
)

type Client struct {
	Firewall   firewall.Client
	Interfaces *interfaces.Client
}

func Login(host string, insecure bool, username, password string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
		Jar: jar,
	}

	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)

	req, err := http.NewRequest(http.MethodPost, host, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &Client{
		Firewall:   firewall.New(httpClient, host),
		Interfaces: interfaces.New(httpClient, host),
	}, nil
}
