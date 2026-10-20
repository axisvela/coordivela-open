package conformance

import (
	"context"
	"testing"
	"time"

	"github.com/axisvela/coordivela-open/sdk/go/v3/adapter"
)

type fakeAdapter struct {
	caps   adapter.Capabilities
	events chan adapter.Event
}

func (f *fakeAdapter) Name() string { return "fake" }
func (f *fakeAdapter) Capabilities() adapter.Capabilities {
	return f.caps
}
func (f *fakeAdapter) Start(context.Context, adapter.StartRequest) (string, error) {
	f.events = make(chan adapter.Event, 2)
	now := time.Now().UTC()
	f.events <- adapter.Event{RunID: "run-1", Kind: adapter.EventStarted, At: now}
	f.events <- adapter.Event{RunID: "run-1", Kind: adapter.EventCompleted, At: now}
	close(f.events)
	return "run-1", nil
}
func (f *fakeAdapter) SendInput(context.Context, string, string) error { return nil }
func (f *fakeAdapter) RespondApproval(context.Context, string, string, adapter.ApprovalDecision) error {
	return nil
}
func (f *fakeAdapter) Cancel(context.Context, string) error { return nil }
func (f *fakeAdapter) Snapshot(context.Context, string) (adapter.Snapshot, error) {
	return adapter.Snapshot{RunID: "run-1", State: "executing"}, nil
}
func (f *fakeAdapter) Subscribe(context.Context, string) (<-chan adapter.Event, error) {
	return f.events, nil
}
func (f *fakeAdapter) Resume(context.Context, string) error { return nil }

func TestRunAcceptsConformingAdapter(t *testing.T) {
	Run(t, &fakeAdapter{}, Options{WaitForTerminal: true, Timeout: time.Second})
}
