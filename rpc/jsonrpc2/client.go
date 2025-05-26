package jsonrpc2

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"

	"github.com/goccy/go-json"
)

type Client struct {
	url        string
	httpClient *http.Client
}

func NewClient(url string, httpClient *http.Client) (*Client, error) {
	if _, err := urlpkg.Parse(url); err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	return &Client{
		url:        url,
		httpClient: httpClient,
	}, nil
}

func (c *Client) Call(ctx context.Context, method string, params any, result any) error {
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal params: %w", err)
	}
	req := NewRequest(-1, method, paramsJSON)
	reqJSON, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(reqJSON))
	if err != nil {
		return fmt.Errorf("new http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, httpResp.Body)
		_ = httpResp.Body.Close()
	}()
	var resp Response
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return fmt.Errorf("unmarshal response: status code %d: %w", httpResp.StatusCode, err)
	}
	if resp.Error != nil {
		return fmt.Errorf("rpc error: status code %d: %w", httpResp.StatusCode, resp.Error)
	}
	if err := json.Unmarshal(resp.Result, result); err != nil {
		return fmt.Errorf("unmarshal result: %w", err)
	}
	return nil
}
