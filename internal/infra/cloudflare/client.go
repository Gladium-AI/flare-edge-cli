package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

type Route struct {
	ID      string `json:"id"`
	Pattern string `json:"pattern"`
	Script  string `json:"script,omitempty"`
}

type DomainRecord struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	Service  string `json:"service,omitempty"`
}

type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewClient(token string) *Client {
	return &Client{
		baseURL:    "https://api.cloudflare.com/client/v4",
		httpClient: &http.Client{},
		token:      token,
	}
}

func (c *Client) FindZoneID(ctx context.Context, accountID, zone string) (string, error) {
	values := url.Values{}
	values.Set("name", zone)
	if accountID != "" {
		values.Set("account.id", accountID)
	}
	var payload response[[]Zone]
	if err := c.do(ctx, http.MethodGet, "/zones?"+values.Encode(), nil, &payload); err != nil {
		return "", err
	}
	if len(payload.Result) == 0 {
		return "", fmt.Errorf("zone %q not found", zone)
	}
	return payload.Result[0].ID, nil
}

func (c *Client) ListRoutes(ctx context.Context, zoneID string) ([]Route, error) {
	var payload response[[]Route]
	if err := c.do(ctx, http.MethodGet, "/zones/"+zoneID+"/workers/routes", nil, &payload); err != nil {
		return nil, err
	}
	return payload.Result, nil
}

func (c *Client) DeleteRoute(ctx context.Context, zoneID, routeID string) error {
	return c.do(ctx, http.MethodDelete, "/zones/"+zoneID+"/workers/routes/"+routeID, nil, nil)
}

func (c *Client) ListDomainRecords(ctx context.Context, accountID, hostname string) ([]DomainRecord, error) {
	values := url.Values{}
	values.Set("hostname", hostname)
	var payload response[[]DomainRecord]
	if err := c.do(ctx, http.MethodGet, "/accounts/"+accountID+"/workers/domains/records?"+values.Encode(), nil, &payload); err != nil {
		return nil, err
	}
	return payload.Result, nil
}

func (c *Client) DeleteDomainRecord(ctx context.Context, accountID, recordID string) error {
	return c.do(ctx, http.MethodDelete, "/accounts/"+accountID+"/workers/domains/records/"+recordID, nil, nil)
}

func (c *Client) do(ctx context.Context, method, path string, body io.Reader, target any) error {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s %s: status %d: %s", method, path, resp.StatusCode, string(data))
	}
	if target == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode %s %s: %w", method, path, err)
	}
	return nil
}

type response[T any] struct {
	Result  T          `json:"result"`
	Success bool       `json:"success"`
	Errors  []apiError `json:"errors"`
}

type apiError struct {
	Message string `json:"message"`
}
