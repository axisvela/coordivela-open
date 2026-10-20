// Package example implements a minimal RuntimeAdapter v1 for contract tests.
package example

import (
	"context"
	"sync"
	"time"

	"github.com/axisvela/coordivela-open/sdk/go/v3/adapter"
)

// Adapter is a deterministic in-memory RuntimeAdapter v1 implementation.
type Adapter struct {
	mu   sync.Mutex
	runs map[string]chan adapter.Event
}

// New returns a deterministic example adapter.
func New() *Adapter {
	return &Adapter{runs: make(map[string]chan adapter.Event)}
}

// Name implements adapter.Adapter.
func (a *Adapter) Name() string { return "example-runtime" }

// Capabilities implements adapter.Adapter.
func (a *Adapter) Capabilities() adapter.Capabilities {
	return adapter.Capabilities{
		Interactive:  true,
		ApprovalHook: true,
		Cancel:       true,
		HTTP:         false,
	}
}

// Start implements adapter.Adapter.
func (a *Adapter) Start(_ context.Context, _ adapter.StartRequest) (string, error) {
	const runID = "example-run"
	events := make(chan adapter.Event, 3)
	now := time.Now().UTC()
	events <- adapter.Event{RunID: runID, Kind: adapter.EventStarted, At: now}
	events <- adapter.Event{RunID: runID, Kind: adapter.EventOutputActivity, At: now, Message: "example output"}
	events <- adapter.Event{RunID: runID, Kind: adapter.EventCompleted, At: now}
	close(events)
	a.mu.Lock()
	a.runs[runID] = events
	a.mu.Unlock()
	return runID, nil
}

// SendInput implements adapter.Adapter.
func (a *Adapter) SendInput(context.Context, string, string) error { return nil }

// RespondApproval implements adapter.Adapter.
func (a *Adapter) RespondApproval(context.Context, string, string, adapter.ApprovalDecision) error {
	return nil
}

// Cancel implements adapter.Adapter.
func (a *Adapter) Cancel(context.Context, string) error { return nil }

// Snapshot implements adapter.Adapter.
func (a *Adapter) Snapshot(_ context.Context, runID string) (adapter.Snapshot, error) {
	return adapter.Snapshot{RunID: runID, State: "completed"}, nil
}

// Subscribe implements adapter.Adapter.
func (a *Adapter) Subscribe(_ context.Context, runID string) (<-chan adapter.Event, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if events, ok := a.runs[runID]; ok {
		return events, nil
	}
	events := make(chan adapter.Event)
	close(events)
	return events, nil
}

// Resume implements adapter.Adapter.
func (a *Adapter) Resume(context.Context, string) error { return nil }
