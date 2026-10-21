package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type PlanStep struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status,omitempty"`
}

type CoordinationPlan struct {
	TenantID     string     `json:"tenantId"`
	WorkspaceID  string     `json:"workspaceId"`
	WorkItemID   string     `json:"workItemId"`
	Version      int64      `json:"version"`
	Objective    string     `json:"objective"`
	Steps        []PlanStep `json:"steps,omitempty"`
	Repositories []string   `json:"repositories,omitempty"`
	Paths        []string   `json:"paths,omitempty"`
	Dependencies []string   `json:"dependencies,omitempty"`
	EvidenceRefs []string   `json:"evidenceRefs,omitempty"`
	CreatedAt    string     `json:"createdAt,omitempty"`
	UpdatedAt    string     `json:"updatedAt,omitempty"`
}

type SaveCoordinationPlanRequest struct {
	FencingToken   int64      `json:"fencingToken"`
	Version        int64      `json:"version"`
	Objective      string     `json:"objective"`
	Steps          []PlanStep `json:"steps,omitempty"`
	Repositories   []string   `json:"repositories,omitempty"`
	Paths          []string   `json:"paths,omitempty"`
	Dependencies   []string   `json:"dependencies,omitempty"`
	EvidenceRefs   []string   `json:"evidenceRefs,omitempty"`
	IdempotencyKey string     `json:"-"`
}

type CoordinationClaim struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenantId"`
	WorkspaceID  string `json:"workspaceId"`
	WorkItemID   string `json:"workItemId"`
	Repository   string `json:"repository"`
	Path         string `json:"path"`
	Mode         string `json:"mode"`
	LeaseID      string `json:"leaseId"`
	FencingToken int64  `json:"fencingToken"`
	ExpiresAt    string `json:"expiresAt"`
	ReleasedAt   string `json:"releasedAt,omitempty"`
	CreatedAt    string `json:"createdAt"`
}

type AddCoordinationClaimRequest struct {
	FencingToken   int64  `json:"fencingToken"`
	Repository     string `json:"repository"`
	Path           string `json:"path"`
	Mode           string `json:"mode"`
	IdempotencyKey string `json:"-"`
}

type CoordinationFinding struct {
	ID          string         `json:"id"`
	TenantID    string         `json:"tenantId"`
	WorkspaceID string         `json:"workspaceId"`
	WorkItemID  string         `json:"workItemId"`
	Kind        string         `json:"kind"`
	Severity    string         `json:"severity"`
	Status      string         `json:"status"`
	Evidence    map[string]any `json:"evidence,omitempty"`
	CreatedAt   string         `json:"createdAt"`
	UpdatedAt   string         `json:"updatedAt"`
}

type DirectiveOutcome struct {
	Kind       string `json:"kind"`
	Note       string `json:"note,omitempty"`
	RecordedBy string `json:"recordedBy,omitempty"`
	RecordedAt string `json:"recordedAt"`
}

type CoordinationDirective struct {
	ID               string            `json:"id"`
	TenantID         string            `json:"tenantId"`
	WorkspaceID      string            `json:"workspaceId"`
	WorkItemID       string            `json:"workItemId"`
	FindingID        string            `json:"findingId"`
	Kind             string            `json:"kind"`
	Status           string            `json:"status"`
	Message          string            `json:"message"`
	IdempotencyKey   string            `json:"idempotencyKey"`
	AuthorizationRef string            `json:"authorizationRef,omitempty"`
	DispatchedAt     string            `json:"dispatchedAt,omitempty"`
	Outcome          *DirectiveOutcome `json:"outcome,omitempty"`
	CreatedAt        string            `json:"createdAt"`
	UpdatedAt        string            `json:"updatedAt"`
}

type DirectiveOutcomeRequest struct {
	Outcome        string `json:"outcome"`
	Note           string `json:"note,omitempty"`
	IdempotencyKey string `json:"-"`
}

type coordinationClaimList struct {
	Claims []CoordinationClaim `json:"claims"`
}

type coordinationFindingList struct {
	Findings []CoordinationFinding `json:"findings"`
}

type coordinationDirectiveList struct {
	Directives []CoordinationDirective `json:"directives"`
}

func (c *Client) GetCoordinationPlan(ctx context.Context, workItemID string) (*CoordinationPlan, error) {
	var out CoordinationPlan
	if err := c.do(ctx, http.MethodGet, coordinationPath(workItemID, "plan"), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SaveCoordinationPlan(ctx context.Context, workItemID string, req SaveCoordinationPlanRequest) (*CoordinationPlan, error) {
	var out CoordinationPlan
	if err := c.do(ctx, http.MethodPost, coordinationPath(workItemID, "plan"), req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) AddCoordinationClaim(ctx context.Context, workItemID string, req AddCoordinationClaimRequest) (*CoordinationClaim, error) {
	var out CoordinationClaim
	if err := c.do(ctx, http.MethodPost, coordinationPath(workItemID, "claims"), req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListCoordinationClaims(ctx context.Context, workItemID string) ([]CoordinationClaim, error) {
	var out coordinationClaimList
	if err := c.do(ctx, http.MethodGet, coordinationPath(workItemID, "claims"), nil, &out); err != nil {
		return nil, err
	}
	return out.Claims, nil
}

func (c *Client) ReleaseCoordinationClaim(ctx context.Context, workItemID, claimID string, fencingToken int64) error {
	path := coordinationPath(workItemID, "claims") + "/" + url.PathEscape(claimID) +
		"?fencingToken=" + strconv.FormatInt(fencingToken, 10)
	return c.do(ctx, http.MethodDelete, path, nil, &statusResponse{})
}

func (c *Client) ListCoordinationFindings(ctx context.Context, workItemID string) ([]CoordinationFinding, error) {
	var out coordinationFindingList
	if err := c.do(ctx, http.MethodGet, coordinationPath(workItemID, "findings"), nil, &out); err != nil {
		return nil, err
	}
	return out.Findings, nil
}

func (c *Client) ListCoordinationDirectives(ctx context.Context, workItemID string) ([]CoordinationDirective, error) {
	var out coordinationDirectiveList
	if err := c.do(ctx, http.MethodGet, coordinationPath(workItemID, "directives"), nil, &out); err != nil {
		return nil, err
	}
	return out.Directives, nil
}

func (c *Client) RecordCoordinationDirectiveOutcome(ctx context.Context, workItemID, directiveID string, req DirectiveOutcomeRequest) (*CoordinationDirective, error) {
	path := coordinationPath(workItemID, "directives") + "/" + url.PathEscape(directiveID) + "/outcome"
	var out CoordinationDirective
	if err := c.do(ctx, http.MethodPost, path, req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

func coordinationPath(workItemID, suffix string) string {
	return "/v1/work-items/" + pathEscape(workItemID) + "/coordination/" + suffix
}
