// Package client is the Apache-2.0 Go SDK for the COORDIVELA coordination API.
package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	ID               string            `json:"id"`
	TenantID         string            `json:"tenantId"`
	WorkspaceID      string            `json:"workspaceId"`
	Profile          string            `json:"profile,omitempty"`
	Repository       string            `json:"repository,omitempty"`
	Paths            []string          `json:"paths,omitempty"`
	ProgramID        string            `json:"programId,omitempty"`
	Title            string            `json:"title"`
	Description      string            `json:"description,omitempty"`
	State            string            `json:"state"`
	Priority         int               `json:"priority,omitempty"`
	OwnerKind        string            `json:"ownerKind,omitempty"`
	OwnerRef         string            `json:"ownerRef,omitempty"`
	LeaseID          string            `json:"leaseId,omitempty"`
	Lease            *Lease            `json:"lease,omitempty"`
	Leases           []Lease           `json:"leases,omitempty"`
	Runtime          *RuntimeRef       `json:"runtime,omitempty"`
	ExternalRefs     map[string]string `json:"externalRefs,omitempty"`
	Assignments      []Assignment      `json:"assignments,omitempty"`
	ConflictOverride *ConflictOverride `json:"conflictOverride,omitempty"`
	VerificationRef  string            `json:"verificationRef,omitempty"`
	Version          int64             `json:"version"`
	CreatedAt        string            `json:"createdAt,omitempty"`
	UpdatedAt        string            `json:"updatedAt,omitempty"`
}

// Lease is a fencing-token lease held by a work-item owner.
type Lease struct {
	ID           string `json:"id"`
	WorkItemID   string `json:"workItemId"`
	OwnerRef     string `json:"ownerRef"`
	Scope        string `json:"scope,omitempty"`
	CanComplete  bool   `json:"canComplete,omitempty"`
	FencingToken int64  `json:"fencingToken"`
	AcquiredAt   string `json:"acquiredAt"`
	ExpiresAt    string `json:"expiresAt"`
	HeartbeatAt  string `json:"heartbeatAt"`
	ReleasedAt   string `json:"releasedAt,omitempty"`
}

// RuntimeRef identifies an adapter run bound to a work item.
type RuntimeRef struct {
	Adapter       string `json:"adapter"`
	RunID         string `json:"runId"`
	AgentID       string `json:"agentId,omitempty"`
	SessionID     string `json:"sessionId,omitempty"`
	CorrelationID string `json:"correlationId"`
	StartedAt     string `json:"startedAt"`
}

// Assignment attaches a human or agent to a work item.
type Assignment struct {
	ID           string `json:"id"`
	WorkItemID   string `json:"workItemId"`
	AssigneeKind string `json:"assigneeKind"`
	AssigneeRef  string `json:"assigneeRef"`
	AssignedBy   string `json:"assignedBy,omitempty"`
	AssignedAt   string `json:"assignedAt"`
	ReleasedAt   string `json:"releasedAt,omitempty"`
	Reason       string `json:"reason,omitempty"`
}

// Dependency is one edge in the work graph.
type Dependency struct {
	WorkItemID  string `json:"workItemId"`
	DependsOnID string `json:"dependsOnId"`
	Kind        string `json:"kind"`
}

// Schedule is a durable coordination trigger. It dispatches work; it does not
// execute work itself.
type Schedule struct {
	ID            string `json:"id"`
	TenantID      string `json:"tenantId"`
	WorkspaceID   string `json:"workspaceId"`
	Profile       string `json:"profile,omitempty"`
	WorkItemID    string `json:"workItemId,omitempty"`
	Kind          string `json:"kind"`
	Expression    string `json:"expression,omitempty"`
	EventType     string `json:"eventType,omitempty"`
	Timezone      string `json:"timezone,omitempty"`
	MisfirePolicy string `json:"misfirePolicy,omitempty"`
	Enabled       bool   `json:"enabled"`
	NextRun       string `json:"nextRun,omitempty"`
	LastRun       string `json:"lastRun,omitempty"`
	LeaseUntil    string `json:"leaseUntil,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
}

// StreamEvent is one coordination event from the SSE stream or history.
type StreamEvent struct {
	ID          uint64         `json:"id"`
	Type        string         `json:"type"`
	TenantID    string         `json:"tenantId"`
	WorkspaceID string         `json:"workspaceId,omitempty"`
	Profile     string         `json:"profile,omitempty"`
	WorkItemID  string         `json:"workItemId"`
	At          string         `json:"at"`
	Data        map[string]any `json:"data,omitempty"`
}

// CreateWorkItem is the request body for creating work.
type CreateWorkItem struct {
	ID          string   `json:"id"`
	Profile     string   `json:"profile,omitempty"`
	Repository  string   `json:"repository,omitempty"`
	Paths       []string `json:"paths,omitempty"`
	ProgramID   string   `json:"programId,omitempty"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	StallAfter  string   `json:"stallAfter,omitempty"`
	Acceptance  []string `json:"acceptance,omitempty"`
}

// PatchWorkItemRequest updates descriptive fields with optimistic versioning.
type PatchWorkItemRequest struct {
	Title           *string `json:"title,omitempty"`
	Description     *string `json:"description,omitempty"`
	Priority        *int    `json:"priority,omitempty"`
	ExpectedVersion int64   `json:"expectedVersion"`
	IdempotencyKey  string  `json:"-"`
}

// DependencyRequest adds one outgoing dependency edge.
type DependencyRequest struct {
	DependsOnID    string `json:"dependsOnId"`
	Kind           string `json:"kind,omitempty"`
	IdempotencyKey string `json:"-"`
}

// ExternalRefRequest links an external system reference.
type ExternalRefRequest struct {
	System         string `json:"system"`
	Ref            string `json:"ref"`
	FencingToken   int64  `json:"fencingToken,omitempty"`
	IdempotencyKey string `json:"-"`
}

// VerificationRequest submits a release candidate to TEKMOVELA.
type VerificationRequest struct {
	ResultRef      string `json:"resultRef,omitempty"`
	SubmissionRef  string `json:"submissionRef,omitempty"`
	IdempotencyKey string `json:"-"`
}

// AssignmentRequest attaches another human or agent.
type AssignmentRequest struct {
	AssigneeKind   string `json:"assigneeKind,omitempty"`
	AssigneeRef    string `json:"assigneeRef"`
	AssignedBy     string `json:"assignedBy,omitempty"`
	Reason         string `json:"reason,omitempty"`
	IdempotencyKey string `json:"-"`
}

// AssignRequest assigns an owner and acquires a lease.
type AssignRequest struct {
	OwnerKind       string `json:"ownerKind"`
	OwnerRef        string `json:"ownerRef"`
	ExpectedVersion int64  `json:"expectedVersion"`
	IdempotencyKey  string `json:"-"`
}

// TokenRequest carries a fencing token for lifecycle actions.
type TokenRequest struct {
	FencingToken   int64  `json:"fencingToken"`
	IdempotencyKey string `json:"-"`
}

// ProgressRequest reports progress and extends the lease.
type ProgressRequest struct {
	FencingToken   int64  `json:"fencingToken"`
	Kind           string `json:"kind"`
	IdempotencyKey string `json:"-"`
}

// InterventionRequest records a coordination intervention.
type InterventionRequest struct {
	Action         string `json:"action"`
	Reason         string `json:"reason,omitempty"`
	IdempotencyKey string `json:"-"`
}

// ReassignRequest replaces the owner and acquires a new lease.
type ReassignRequest struct {
	OwnerKind      string `json:"ownerKind"`
	OwnerRef       string `json:"ownerRef"`
	FencingToken   int64  `json:"fencingToken"`
	Authorized     bool   `json:"authorized,omitempty"`
	IdempotencyKey string `json:"-"`
}

// CompleteRequest completes work through configured gates.
type CompleteRequest struct {
	FencingToken    int64  `json:"fencingToken"`
	ResultRef       string `json:"resultRef,omitempty"`
	SubmissionRef   string `json:"submissionRef,omitempty"`
	VerificationRef string `json:"verificationRef,omitempty"`
	IdempotencyKey  string `json:"-"`
}

// RunRequest starts a runtime adapter run.
type RunRequest struct {
	AdapterName    string            `json:"adapterName"`
	Workspace      string            `json:"workspace,omitempty"`
	Prompt         string            `json:"prompt,omitempty"`
	Model          string            `json:"model,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	FencingToken   int64             `json:"fencingToken"`
	IdempotencyKey string            `json:"-"`
}

// SaveScheduleRequest defines or replaces a schedule.
type SaveScheduleRequest struct {
	ID             string `json:"id"`
	Profile        string `json:"profile,omitempty"`
	WorkItemID     string `json:"workItemId,omitempty"`
	Kind           string `json:"kind"`
	Expression     string `json:"expression,omitempty"`
	EventType      string `json:"eventType,omitempty"`
	Timezone       string `json:"timezone,omitempty"`
	MisfirePolicy  string `json:"misfirePolicy,omitempty"`
	Enabled        bool   `json:"enabled"`
	IdempotencyKey string `json:"-"`
}

// ToggleScheduleRequest enables or disables a schedule.
type ToggleScheduleRequest struct {
	Enabled        bool   `json:"enabled"`
	IdempotencyKey string `json:"-"`
}

// RunDueRequest bounds the number of claimed schedules.
type RunDueRequest struct {
	Limit          int    `json:"limit,omitempty"`
	IdempotencyKey string `json:"-"`
}

// RunEventRequest claims schedules bound to an event type.
type RunEventRequest struct {
	EventType      string `json:"eventType"`
	IdempotencyKey string `json:"-"`
}

type workItemList struct {
	WorkItems []WorkItem `json:"workItems"`
}

type dependencyList struct {
	Dependencies []Dependency `json:"dependencies"`
}

type conflictList struct {
	Conflicts []WorkItem `json:"conflicts"`
}

type historyResponse struct {
	Events []StreamEvent `json:"events"`
}

type assignmentResponse struct {
	WorkItem  *WorkItem  `json:"workItem"`
	Conflicts []WorkItem `json:"conflicts"`
}

type workItemResponse struct {
	WorkItem *WorkItem `json:"workItem"`
}

type verificationResponse struct {
	ReleaseCandidateRef string `json:"releaseCandidateRef"`
}

type runResponse struct {
	RunID string `json:"runId"`
}

type stallResponse struct {
	Stalled bool `json:"stalled"`
}

type statusResponse struct {
	Status string `json:"status"`
}

type scheduleList struct {
	Schedules []Schedule `json:"schedules"`
}

type triggeredSchedules struct {
	Triggered []Schedule `json:"triggered"`
}

type programProjection struct {
	WorkItems []WorkItem     `json:"workItems"`
	States    map[string]int `json:"states"`
}

type workItemProjection struct {
	WorkItem     *WorkItem    `json:"workItem"`
	Dependencies []Dependency `json:"dependencies"`
}

// ListWorkItems returns work items in the configured workspace.
func (c *Client) ListWorkItems(ctx context.Context) ([]WorkItem, error) {
	var out workItemList
	if err := c.do(ctx, http.MethodGet, "/v1/work-items", nil, &out); err != nil {
		return nil, err
	}
	return out.WorkItems, nil
}

// CreateWorkItem creates a work item with a generated idempotency key.
func (c *Client) CreateWorkItem(ctx context.Context, item CreateWorkItem) (*WorkItem, error) {
	return c.CreateWorkItemWithKey(ctx, newIdempotencyKey(), item)
}

// CreateWorkItemWithKey creates a work item with a caller-supplied idempotency key.
func (c *Client) CreateWorkItemWithKey(ctx context.Context, key string, item CreateWorkItem) (*WorkItem, error) {
	if strings.TrimSpace(key) == "" {
		key = newIdempotencyKey()
	}
	var out WorkItem
	if err := c.do(ctx, http.MethodPost, "/v1/work-items", item, &out, key); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWorkItem returns one work item.
func (c *Client) GetWorkItem(ctx context.Context, id string) (*WorkItem, error) {
	var out WorkItem
	if err := c.do(ctx, http.MethodGet, "/v1/work-items/"+pathEscape(id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PatchWorkItem updates title, description, or priority.
func (c *Client) PatchWorkItem(ctx context.Context, id string, req PatchWorkItemRequest) (*WorkItem, error) {
	var out WorkItem
	if err := c.do(ctx, http.MethodPatch, "/v1/work-items/"+pathEscape(id), req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

// History returns bounded coordination history for a work item.
func (c *Client) History(ctx context.Context, id string, limit int) ([]StreamEvent, error) {
	path := "/v1/work-items/" + pathEscape(id) + "/history"
	if limit > 0 {
		path += "?limit=" + strconv.Itoa(limit)
	}
	var out historyResponse
	if err := c.do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return out.Events, nil
}

// Conflicts returns active work items with repository or path conflicts.
func (c *Client) Conflicts(ctx context.Context, id string) ([]WorkItem, error) {
	var out conflictList
	if err := c.do(ctx, http.MethodGet, "/v1/work-items/"+pathEscape(id)+"/conflicts", nil, &out); err != nil {
		return nil, err
	}
	return out.Conflicts, nil
}

// ListDependencies returns outgoing dependency edges.
func (c *Client) ListDependencies(ctx context.Context, id string) ([]Dependency, error) {
	var out dependencyList
	if err := c.do(ctx, http.MethodGet, "/v1/work-items/"+pathEscape(id)+"/dependencies", nil, &out); err != nil {
		return nil, err
	}
	return out.Dependencies, nil
}

// AddDependency adds one outgoing dependency edge.
func (c *Client) AddDependency(ctx context.Context, id string, req DependencyRequest) (*Dependency, error) {
	var out Dependency
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/dependencies", req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveDependency removes one outgoing dependency edge.
func (c *Client) RemoveDependency(ctx context.Context, id, dependsOnID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/work-items/"+pathEscape(id)+"/dependencies/"+pathEscape(dependsOnID), nil, &statusResponse{})
}

// LinkExternalRef links an external system reference.
func (c *Client) LinkExternalRef(ctx context.Context, id string, req ExternalRefRequest) (*WorkItem, error) {
	var out WorkItem
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/external-refs", req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

// SubmitVerification submits a release candidate and returns its reference.
func (c *Client) SubmitVerification(ctx context.Context, id string, req VerificationRequest) (string, error) {
	var out verificationResponse
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/verification", req, &out, req.IdempotencyKey); err != nil {
		return "", err
	}
	return out.ReleaseCandidateRef, nil
}

// AddAssignment attaches another human or agent to a work item.
func (c *Client) AddAssignment(ctx context.Context, id string, req AssignmentRequest) (*WorkItem, []WorkItem, error) {
	var out assignmentResponse
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/assignments", req, &out, req.IdempotencyKey); err != nil {
		return nil, nil, err
	}
	return out.WorkItem, out.Conflicts, nil
}

// ReleaseAssignment releases one assignment.
func (c *Client) ReleaseAssignment(ctx context.Context, id, assignmentID string) (*WorkItem, error) {
	var out workItemResponse
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/assignments/"+pathEscape(assignmentID)+"/release", nil, &out); err != nil {
		return nil, err
	}
	return out.WorkItem, nil
}

// ListLeases returns active leases for a work item.
func (c *Client) ListLeases(ctx context.Context, id string) ([]Lease, error) {
	var out struct {
		Leases []Lease `json:"leases"`
	}
	if err := c.do(ctx, http.MethodGet, "/v1/work-items/"+pathEscape(id)+"/leases", nil, &out); err != nil {
		return nil, err
	}
	return out.Leases, nil
}

// AcquireLease adds a collaborator lease with a resource scope.
func (c *Client) AcquireLease(ctx context.Context, id string, req AcquireLeaseRequest) (*Lease, error) {
	var out Lease
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/leases", req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReleaseLease releases a collaborator lease.
func (c *Client) ReleaseLease(ctx context.Context, id, leaseID string, fencingToken int64, idempotencyKey ...string) error {
	key := ""
	if len(idempotencyKey) > 0 {
		key = idempotencyKey[0]
	}
	path := "/v1/work-items/" + pathEscape(id) + "/leases/" + pathEscape(leaseID) + "?fencingToken=" + strconv.FormatInt(fencingToken, 10)
	return c.do(ctx, http.MethodDelete, path, nil, &statusResponse{}, key)
}

// MarkReady marks work ready after dependency checks.
func (c *Client) MarkReady(ctx context.Context, id string) (*WorkItem, error) {
	var out WorkItem
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/ready", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Assign assigns an owner and acquires a lease.
func (c *Client) Assign(ctx context.Context, id string, req AssignRequest) (*WorkItem, error) {
	var out WorkItem
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/assign", req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

// OverrideConflicts records an authorized path-conflict override.
func (c *Client) OverrideConflicts(ctx context.Context, id string, req ConflictOverrideRequest) (*WorkItem, error) {
	var out WorkItem
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/conflict-override", req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

// ClearConflictOverride removes a path-conflict override.
func (c *Client) ClearConflictOverride(ctx context.Context, id string, expectedVersion int64, idempotencyKey ...string) (*WorkItem, error) {
	var out WorkItem
	path := "/v1/work-items/" + pathEscape(id) + "/conflict-override?expectedVersion=" + strconv.FormatInt(expectedVersion, 10)
	key := ""
	if len(idempotencyKey) > 0 {
		key = idempotencyKey[0]
	}
	if err := c.do(ctx, http.MethodDelete, path, nil, &out, key); err != nil {
		return nil, err
	}
	return &out, nil
}

// Start starts execution with a fencing token.
func (c *Client) Start(ctx context.Context, id string, fencingToken int64) (*WorkItem, error) {
	return c.lifecycle(ctx, id, "start", TokenRequest{FencingToken: fencingToken})
}

// Progress reports progress and extends the lease.
func (c *Client) Progress(ctx context.Context, id string, req ProgressRequest) (*WorkItem, error) {
	var out WorkItem
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/progress", req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

// Heartbeat extends the lease without a progress event.
func (c *Client) Heartbeat(ctx context.Context, id string, fencingToken int64) (*WorkItem, error) {
	return c.lifecycle(ctx, id, "heartbeat", TokenRequest{FencingToken: fencingToken})
}

// RecordIntervention records a coordination intervention.
func (c *Client) RecordIntervention(ctx context.Context, id string, req InterventionRequest) error {
	return c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/interventions", req, &statusResponse{}, req.IdempotencyKey)
}

// StallCheck evaluates staleness on demand.
func (c *Client) StallCheck(ctx context.Context, id string) (bool, error) {
	var out stallResponse
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/stall-check", nil, &out); err != nil {
		return false, err
	}
	return out.Stalled, nil
}

// Reassign replaces the owner and acquires a new lease.
func (c *Client) Reassign(ctx context.Context, id string, req ReassignRequest) (*WorkItem, error) {
	var out WorkItem
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/reassign", req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

// Complete completes work through configured gates.
func (c *Client) Complete(ctx context.Context, id string, req CompleteRequest) (*WorkItem, error) {
	var out WorkItem
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/complete", req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

// StartRun starts a runtime adapter run and returns its run id.
func (c *Client) StartRun(ctx context.Context, id string, req RunRequest) (string, error) {
	var out runResponse
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/run", req, &out, req.IdempotencyKey); err != nil {
		return "", err
	}
	return out.RunID, nil
}

// CancelRun cancels the active runtime run.
func (c *Client) CancelRun(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/cancel-run", nil, &statusResponse{})
}

// ListSchedules returns schedules in the configured workspace.
func (c *Client) ListSchedules(ctx context.Context) ([]Schedule, error) {
	var out scheduleList
	if err := c.do(ctx, http.MethodGet, "/v1/schedules", nil, &out); err != nil {
		return nil, err
	}
	return out.Schedules, nil
}

// SaveSchedule defines or replaces a schedule.
func (c *Client) SaveSchedule(ctx context.Context, req SaveScheduleRequest) (*Schedule, error) {
	var out Schedule
	if err := c.do(ctx, http.MethodPost, "/v1/schedules", req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

// RunDueSchedules claims due schedules and emits triggers.
func (c *Client) RunDueSchedules(ctx context.Context, req RunDueRequest) ([]Schedule, error) {
	var out triggeredSchedules
	if err := c.do(ctx, http.MethodPost, "/v1/schedules/run-due", req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return out.Triggered, nil
}

// RunEventSchedules claims schedules bound to an event type.
func (c *Client) RunEventSchedules(ctx context.Context, req RunEventRequest) ([]Schedule, error) {
	var out triggeredSchedules
	if err := c.do(ctx, http.MethodPost, "/v1/schedules/run-event", req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return out.Triggered, nil
}

// CompleteSchedule computes the next run and releases the schedule lease.
func (c *Client) CompleteSchedule(ctx context.Context, id string) (*Schedule, error) {
	var out Schedule
	if err := c.do(ctx, http.MethodPost, "/v1/schedules/"+pathEscape(id)+"/complete", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReleaseSchedule releases a schedule lease.
func (c *Client) ReleaseSchedule(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodPost, "/v1/schedules/"+pathEscape(id)+"/release", nil, &statusResponse{})
}

// ToggleSchedule enables or disables a schedule.
func (c *Client) ToggleSchedule(ctx context.Context, id string, req ToggleScheduleRequest) (*Schedule, error) {
	var out Schedule
	if err := c.do(ctx, http.MethodPost, "/v1/schedules/"+pathEscape(id)+"/toggle", req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

// TriggerSchedule triggers a schedule manually.
func (c *Client) TriggerSchedule(ctx context.Context, id string) (*Schedule, error) {
	var out Schedule
	if err := c.do(ctx, http.MethodPost, "/v1/schedules/"+pathEscape(id)+"/trigger", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteSchedule removes a schedule.
func (c *Client) DeleteSchedule(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/schedules/"+pathEscape(id), nil, &statusResponse{})
}

// MyWork returns the current principal's active work queue.
func (c *Client) MyWork(ctx context.Context) ([]WorkItem, error) {
	var out workItemList
	if err := c.do(ctx, http.MethodGet, "/v1/projections/my-work", nil, &out); err != nil {
		return nil, err
	}
	return out.WorkItems, nil
}

// Attention returns blocked, review, or stalled work.
func (c *Client) Attention(ctx context.Context) ([]WorkItem, error) {
	var out workItemList
	if err := c.do(ctx, http.MethodGet, "/v1/projections/attention", nil, &out); err != nil {
		return nil, err
	}
	return out.WorkItems, nil
}

// SymbivelaAttention is one read-only SYMBIVELA attention item.
type SymbivelaAttention struct {
	AttentionID string `json:"attention_id"`
	WorkspaceID string `json:"workspace_id"`
	Kind        string `json:"kind"`
	Status      string `json:"status"`
	Priority    int    `json:"priority"`
}

// SymbivelaPlanSummary is the plan section of a SYMBIVELA work projection.
type SymbivelaPlanSummary struct {
	Version      int64    `json:"version"`
	Objective    string   `json:"objective"`
	Repositories []string `json:"repositories,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	UpdatedAt    string   `json:"updatedAt"`
}

// SymbivelaClaimSummary is one claim in a SYMBIVELA work projection.
type SymbivelaClaimSummary struct {
	ID         string `json:"id"`
	Repository string `json:"repository"`
	Path       string `json:"path"`
	Mode       string `json:"mode"`
}

// SymbivelaFindingSummary is one finding in a SYMBIVELA work projection.
type SymbivelaFindingSummary struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Severity string `json:"severity"`
	Status   string `json:"status"`
}

// SymbivelaWorkProjection is the read-only per-work-item SYMBIVELA projection.
type SymbivelaWorkProjection struct {
	WorkItemID       string                    `json:"workItemId"`
	GoalRef          string                    `json:"goalRef"`
	WorkspaceID      string                    `json:"workspaceId"`
	Repository       string                    `json:"repository"`
	State            string                    `json:"state"`
	OwnerKind        string                    `json:"ownerKind"`
	OwnerRef         string                    `json:"ownerRef"`
	Version          int64                     `json:"version"`
	UpdatedAt        string                    `json:"updatedAt"`
	CoordinationPlan *SymbivelaPlanSummary     `json:"coordinationPlan,omitempty"`
	Claims           []SymbivelaClaimSummary   `json:"claims,omitempty"`
	Findings         []SymbivelaFindingSummary `json:"findings,omitempty"`
	EvidenceRefs     []string                  `json:"evidenceRefs,omitempty"`
}

// AttachOpenCodeRequest binds an existing OpenCode session to a work item.
type AttachOpenCodeRequest struct {
	SessionID      string `json:"sessionId"`
	FencingToken   int64  `json:"fencingToken"`
	IdempotencyKey string `json:"-"`
}

// SymbivelaAttention returns read-only SYMBIVELA attention for the client scope.
func (c *Client) SymbivelaAttention(ctx context.Context) ([]SymbivelaAttention, error) {
	var out struct {
		Attention []SymbivelaAttention `json:"attention"`
	}
	if err := c.do(ctx, http.MethodGet, "/v1/projections/symbivela-attention", nil, &out); err != nil {
		return nil, err
	}
	return out.Attention, nil
}

// Program returns a program projection.
func (c *Client) Program(ctx context.Context, programID string) ([]WorkItem, map[string]int, error) {
	var out programProjection
	if err := c.do(ctx, http.MethodGet, "/v1/projections/programs/"+pathEscape(programID), nil, &out); err != nil {
		return nil, nil, err
	}
	return out.WorkItems, out.States, nil
}

// WorkItemProjection returns one work item with its outgoing dependencies.
func (c *Client) WorkItemProjection(ctx context.Context, id string) (*WorkItem, []Dependency, error) {
	var out workItemProjection
	if err := c.do(ctx, http.MethodGet, "/v1/projections/work-items/"+pathEscape(id), nil, &out); err != nil {
		return nil, nil, err
	}
	return out.WorkItem, out.Dependencies, nil
}

// SymbivelaWorkItemProjection returns the read-only SYMBIVELA projection for one work item.
func (c *Client) SymbivelaWorkItemProjection(ctx context.Context, id string) (*SymbivelaWorkProjection, error) {
	var out SymbivelaWorkProjection
	if err := c.do(ctx, http.MethodGet, "/v1/projections/work-items/"+pathEscape(id)+"/symbivela", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AttachOpenCodeSession binds an existing OpenCode session to one work item.
func (c *Client) AttachOpenCodeSession(ctx context.Context, id string, req AttachOpenCodeRequest) (*WorkItem, error) {
	var out workItemResponse
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/attach-opencode", req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return out.WorkItem, nil
}

// StreamEvents opens the SSE coordination stream. The caller must close the
// response body. Pass a non-zero lastEventID to request bounded replay.
func (c *Client) StreamEvents(ctx context.Context, lastEventID uint64) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/v1/events", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("X-Axis-Tenant", c.Tenant)
	req.Header.Set("X-Axis-Workspace", c.Workspace)
	req.Header.Set("X-Axis-Principal", c.Principal)
	if lastEventID > 0 {
		req.Header.Set("Last-Event-ID", strconv.FormatUint(lastEventID, 10))
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		resp.Body.Close()
		return nil, fmt.Errorf("coordivela: %s: %s", resp.Status, strings.TrimSpace(string(raw)))
	}
	return resp, nil
}

func (c *Client) lifecycle(ctx context.Context, id, action string, req TokenRequest) (*WorkItem, error) {
	var out WorkItem
	if err := c.do(ctx, http.MethodPost, "/v1/work-items/"+pathEscape(id)+"/"+action, req, &out, req.IdempotencyKey); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) do(ctx context.Context, method, path string, body, out any, idempotencyKey ...string) error {
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
	if method != http.MethodGet {
		key := ""
		if len(idempotencyKey) > 0 {
			key = strings.TrimSpace(idempotencyKey[0])
		}
		if key == "" {
			key = newIdempotencyKey()
		}
		req.Header.Set("Idempotency-Key", key)
	}
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

func pathEscape(value string) string {
	return url.PathEscape(value)
}

func newIdempotencyKey() string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return fmt.Sprintf("sdk-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("sdk-%x", raw)
}

// AcquireLeaseRequest asks for a collaborator lease scope.
type AcquireLeaseRequest struct {
	OwnerRef       string `json:"ownerRef"`
	Scope          string `json:"scope"`
	TTLSeconds     int    `json:"ttlSeconds,omitempty"`
	CanComplete    bool   `json:"canComplete,omitempty"`
	IdempotencyKey string `json:"-"`
}

// ConflictOverrideRequest records an authorized path-conflict override.
type ConflictOverrideRequest struct {
	Reason          string `json:"reason"`
	TTLSeconds      int    `json:"ttlSeconds,omitempty"`
	ExpectedVersion int64  `json:"expectedVersion"`
	IdempotencyKey  string `json:"-"`
}

// ConflictOverride records an active path-conflict override.
type ConflictOverride struct {
	Reason    string `json:"reason"`
	ActorRef  string `json:"actorRef"`
	CreatedAt string `json:"createdAt"`
	ExpiresAt string `json:"expiresAt,omitempty"`
}
