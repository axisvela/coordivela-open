// Package adapter defines the public RuntimeAdapter v1 contract.
package adapter

import (
	"context"
	"errors"
	"time"
)

// Version is the published RuntimeAdapter contract version.
const Version = "1.0.0"

// EventKind is a normalized runtime observation.
type EventKind string

const (
	EventStarted            EventKind = "started"
	EventOutputActivity     EventKind = "output_activity"
	EventWaitingForInput    EventKind = "waiting_for_input"
	EventWaitingForApproval EventKind = "waiting_for_approval"
	EventIdle               EventKind = "idle"
	EventHeartbeat          EventKind = "heartbeat"
	EventStalled            EventKind = "stalled"
	EventExited             EventKind = "exited"
	EventError              EventKind = "error"
	EventCompleted          EventKind = "completed"
)

// IsKnownEventKind reports whether kind belongs to the v1 event vocabulary.
func IsKnownEventKind(kind EventKind) bool {
	switch kind {
	case EventStarted, EventOutputActivity, EventWaitingForInput,
		EventWaitingForApproval, EventIdle, EventHeartbeat, EventStalled,
		EventExited, EventError, EventCompleted:
		return true
	default:
		return false
	}
}

// IsTerminalEvent reports whether kind ends a run's event stream.
func IsTerminalEvent(kind EventKind) bool {
	return kind == EventExited || kind == EventError || kind == EventCompleted
}

// ApprovalDecision is a human or policy approval outcome.
type ApprovalDecision string

const (
	ApprovalApprove ApprovalDecision = "approve"
	ApprovalDeny    ApprovalDecision = "deny"
)

// Capabilities describes what an adapter can actually do.
type Capabilities struct {
	Interactive        bool `json:"interactive"`
	StructuredStream   bool `json:"structuredStream"`
	ApprovalHook       bool `json:"approvalHook"`
	SessionResume      bool `json:"sessionResume"`
	Cancel             bool `json:"cancel"`
	PTY                bool `json:"pty"`
	HTTP               bool `json:"http"`
	WorkspaceIsolation bool `json:"workspaceIsolation"`
	EvidenceHooks      bool `json:"evidenceHooks"`
}

// StartRequest asks an adapter to start one coding run.
type StartRequest struct {
	WorkItemID string            `json:"workItemId"`
	Workspace  string            `json:"workspace"`
	Prompt     string            `json:"prompt"`
	Model      string            `json:"model,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// Snapshot is a point-in-time adapter observation.
type Snapshot struct {
	RunID           string    `json:"runId"`
	State           string    `json:"state"`
	LastOutputAt    time.Time `json:"lastOutputAt,omitempty"`
	PendingInput    bool      `json:"pendingInput"`
	PendingApproval bool      `json:"pendingApproval"`
	RequestID       string    `json:"requestId,omitempty"`
	ExitCode        *int      `json:"exitCode,omitempty"`
	Output          string    `json:"output,omitempty"`
}

// Event is one normalized runtime observation.
type Event struct {
	RunID     string            `json:"runId"`
	Kind      EventKind         `json:"kind"`
	At        time.Time         `json:"at"`
	Message   string            `json:"message,omitempty"`
	RequestID string            `json:"requestId,omitempty"`
	Risk      string            `json:"risk,omitempty"`
	Data      map[string]string `json:"data,omitempty"`
}

// Error taxonomy for v1.
var (
	ErrUnsupported     = errors.New("runtime adapter: operation unsupported")
	ErrUnavailable     = errors.New("runtime adapter: unavailable")
	ErrUnauthenticated = errors.New("runtime adapter: unauthenticated")
	ErrInvalidRequest  = errors.New("runtime adapter: invalid request")
	ErrConflict        = errors.New("runtime adapter: conflict")
)

// Adapter controls and observes one runtime family.
type Adapter interface {
	Name() string
	Capabilities() Capabilities
	Start(ctx context.Context, req StartRequest) (string, error)
	SendInput(ctx context.Context, runID, text string) error
	RespondApproval(ctx context.Context, runID, requestID string, decision ApprovalDecision) error
	Cancel(ctx context.Context, runID string) error
	Snapshot(ctx context.Context, runID string) (Snapshot, error)
	Subscribe(ctx context.Context, runID string) (<-chan Event, error)
	Resume(ctx context.Context, runID string) error
}
