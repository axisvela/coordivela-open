// Package client is the Apache-2.0 Go SDK for the COORDIVELA coordination API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client calls the COORDIVELA HTTP API with trusted-context headers.
type Client struct {
	BaseURL   string
	Tenant    string
	Workspace string
	Principal string
	HTTP      *http.Client
}

// New returns a client. The base URL must include the scheme and host, e.g.
// http://localhost:1976.
func New(baseURL, tenant, workspace string) *Client {
	return &Client{
		BaseURL:   strings.TrimRight(baseURL, "/"),
		Tenant:    tenant,
		Workspace: workspace,
		Principal: "coordivela-sdk",
		HTTP:      http.DefaultClient,
	}
}

// WorkItem is one coordinated unit of development work.
type WorkItem struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenantId"`
	WorkspaceID string `json:"workspaceId"`
	ProgramID   string `json:"programId,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	State       string `json:"state"`
	Priority    int    `json:"priority,omitempty"`
	OwnerKind   string `json:"ownerKind,omitempty"`
	OwnerRef    string `json:"ownerRef,omitempty"`
	Version     int64  `json:"version"`
}

// CreateWorkItem is the request body for creating work.
type CreateWorkItem struct {
	ID          string   `json:"id"`
	ProgramID   string   `json:"programId,omitempty"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	StallAfter  string   `json:"stallAfter,omitempty"`
	Acceptance  []string `json:"acceptance,omitempty"`
}

// ListWorkItems returns work items in the configured workspace.
func (c *Client) ListWorkItems(ctx context.Context) ([]WorkItem, error) {
	var out struct {
		WorkItems []WorkItem `json:"workItems"`
	}
	if err := c.do(ctx, http.MethodGet, "/v1/work-items", nil, &out); err != nil {
		return nil, err
	}
	return out.WorkItems, nil
}

// CreateWorkItem creates a work item.
func (c *Client) CreateWorkItem(ctx context.Context, item CreateWorkItem) (*WorkItem, error) {
	var out WorkItem
	if err := c.do(ctx, http.MethodPost, "/v1/work-items", item, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWorkItem returns one work item.
func (c *Client) GetWorkItem(ctx context.Context, id string) (*WorkItem, error) {
	var out WorkItem
	if err := c.do(ctx, http.MethodGet, "/v1/work-items/"+id, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("X-Axis-Tenant", c.Tenant)
	req.Header.Set("X-Axis-Workspace", c.Workspace)
	req.Header.Set("X-Axis-Principal", c.Principal)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("coordivela: %s: %s", resp.Status, strings.TrimSpace(string(raw)))
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
