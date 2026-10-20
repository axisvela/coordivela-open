// Package conformance validates RuntimeAdapter v1 implementations.
package conformance

import (
	"context"
	"time"

	"github.com/axisvela/coordivela-open/sdk/go/v3/adapter"
)

// TestingT is the minimal testing surface used by the conformance kit.
type TestingT interface {
	Helper()
	Fatalf(format string, args ...any)
	Errorf(format string, args ...any)
}

// Options configures a conformance run.
type Options struct {
	Start           adapter.StartRequest
	Timeout         time.Duration
	WaitForTerminal bool
	CheckInput      bool
	InputText       string
	CheckCancel     bool
	CheckResume     bool
}

// Run executes the RuntimeAdapter v1 conformance checks.
func Run(t TestingT, a adapter.Adapter, opts Options) {
	t.Helper()
	if a == nil {
		t.Fatalf("adapter is nil")
	}
	if a.Name() == "" {
		t.Fatalf("adapter Name() is empty")
	}
	caps := a.Capabilities()
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	runID, err := a.Start(ctx, opts.Start)
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	if runID == "" {
		t.Fatalf("Start returned an empty run id")
	}
	snapshot, err := a.Snapshot(ctx, runID)
	if err != nil {
		t.Fatalf("Snapshot: %v", err)
	}
	if snapshot.State == "" {
		t.Fatalf("Snapshot.State is empty")
	}
	if snapshot.RunID != "" && snapshot.RunID != runID {
		t.Fatalf("Snapshot.RunID = %q, want %q", snapshot.RunID, runID)
	}

	events, err := a.Subscribe(ctx, runID)
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	first, ok := waitEvent(t, events, timeout)
	if !ok {
		t.Fatalf("Subscribe produced no event within %s", timeout)
	}
	validateEvent(t, first, runID)
	if opts.WaitForTerminal {
		for !adapter.IsTerminalEvent(first.Kind) {
			next, ok := waitEvent(t, events, timeout)
			if !ok {
				t.Fatalf("event stream ended before a terminal event; last kind = %s", first.Kind)
			}
			validateEvent(t, next, runID)
			first = next
		}
	}

	if opts.CheckInput && caps.Interactive {
		text := opts.InputText
		if text == "" {
			text = "continue"
		}
		if err := a.SendInput(ctx, runID, text); err != nil {
			t.Errorf("SendInput: %v", err)
		}
	}
	if opts.CheckResume && caps.SessionResume {
		if err := a.Resume(ctx, runID); err != nil {
			t.Errorf("Resume: %v", err)
		}
	}
	if opts.CheckCancel && caps.Cancel {
		if err := a.Cancel(ctx, runID); err != nil {
			t.Errorf("Cancel: %v", err)
		}
	}
}

func waitEvent(t TestingT, events <-chan adapter.Event, timeout time.Duration) (adapter.Event, bool) {
	t.Helper()
	select {
	case event, ok := <-events:
		if !ok {
			return adapter.Event{}, false
		}
		return event, true
	case <-time.After(timeout):
		return adapter.Event{}, false
	}
}

func validateEvent(t TestingT, event adapter.Event, runID string) {
	t.Helper()
	if event.RunID != "" && event.RunID != runID {
		t.Errorf("event RunID = %q, want %q", event.RunID, runID)
	}
	if !adapter.IsKnownEventKind(event.Kind) {
		t.Errorf("unknown event kind %q", event.Kind)
	}
	if event.At.IsZero() {
		t.Errorf("event %q has no timestamp", event.Kind)
	}
	if event.Kind == adapter.EventWaitingForApproval && event.RequestID == "" {
		t.Errorf("waiting_for_approval event has no request id")
	}
}
