package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	tokenKey = "X-CSRF-TOKEN"
)

type Client interface {
	Post(context.Context, *Operation) (*Operation, error)
	Get(context.Context) (*Operation, error)
}

type client struct {
	httpClient *http.Client
	baseURL    string
}

func New(httpClient *http.Client, baseURL string) Client {
	return &client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *client) Get(context.Context) (*Operation, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/edge/get.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return toOperation(resp.Body)
}

func (c *client) Post(ctx context.Context, in *Operation) (*Operation, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, string(data))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/edge/batch.json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for _, cookie := range c.httpClient.Jar.Cookies(req.URL) {
		if cookie.Name == tokenKey {
			req.Header.Set(tokenKey, cookie.Value)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return toOperation(resp.Body)
}

func toOperation(reader io.Reader) (*Operation, error) {
	var out Operation
	{
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(os.Stderr, string(data))

		if err := json.Unmarshal([]byte(data), &out); err != nil {
			return nil, err
		}
	}

	if out.Failed() {
		return nil, errors.New("The operation was not successfull.")
	}
	return &out, nil
}
