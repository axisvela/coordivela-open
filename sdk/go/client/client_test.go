package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateAndListWorkItems(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Axis-Tenant") != "acme" {
			http.Error(w, "missing tenant", http.StatusBadRequest)
			return
		}
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/work-items":
			if r.Header.Get("Idempotency-Key") == "" {
				http.Error(w, "missing idempotency key", http.StatusBadRequest)
				return
			}
			var body CreateWorkItem
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(WorkItem{ID: body.ID, Title: body.Title, State: "proposed", Version: 1})
		case r.Method == http.MethodGet && r.URL.Path == "/v1/work-items":
			_ = json.NewEncoder(w).Encode(map[string]any{"workItems": []WorkItem{{ID: "wi-1", Title: "test", State: "proposed", Version: 1}}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := New(srv.URL, "acme", "default")
	created, err := c.CreateWorkItem(context.Background(), CreateWorkItem{ID: "wi-1", Title: "test"})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID != "wi-1" || created.State != "proposed" {
		t.Fatalf("created = %+v", created)
	}
	items, err := c.ListWorkItems(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != "wi-1" {
		t.Fatalf("items = %+v", items)
	}
}

func TestCreateWorkItemWithKey(t *testing.T) {
	const key = "goal-42-wi-1"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Idempotency-Key"); got != key {
			http.Error(w, "unexpected idempotency key: "+got, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(WorkItem{ID: "wi-1", Title: "test", State: "proposed", Version: 1})
	}))
	defer srv.Close()

	c := New(srv.URL, "acme", "default")
	if _, err := c.CreateWorkItemWithKey(context.Background(), key, CreateWorkItem{ID: "wi-1", Title: "test"}); err != nil {
		t.Fatal(err)
	}
}

func TestAddDependencySendsIdempotencyKey(t *testing.T) {
	var gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/work-items/wi-1/dependencies" {
			http.NotFound(w, r)
			return
		}
		gotKey = r.Header.Get("Idempotency-Key")
		var body DependencyRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if body.DependsOnID != "wi-0" {
			http.Error(w, "unexpected dependency", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Dependency{WorkItemID: "wi-1", DependsOnID: "wi-0", Kind: "blocks"})
	}))
	defer srv.Close()

	c := New(srv.URL, "acme", "default")
	dep, err := c.AddDependency(context.Background(), "wi-1", DependencyRequest{DependsOnID: "wi-0"})
	if err != nil {
		t.Fatal(err)
	}
	if gotKey == "" {
		t.Fatal("missing Idempotency-Key")
	}
	if dep.DependsOnID != "wi-0" || dep.Kind != "blocks" {
		t.Fatalf("dependency = %+v", dep)
	}
}

func TestSaveScheduleUsesCallerKey(t *testing.T) {
	const key = "schedule-save-1"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/schedules" {
			http.NotFound(w, r)
			return
		}
		if got := r.Header.Get("Idempotency-Key"); got != key {
			http.Error(w, "unexpected idempotency key: "+got, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Schedule{ID: "sched-1", Kind: "interval", Enabled: true})
	}))
	defer srv.Close()

	c := New(srv.URL, "acme", "default")
	schedule, err := c.SaveSchedule(context.Background(), SaveScheduleRequest{
		ID: "sched-1", Kind: "interval", Expression: "5m", Enabled: true, IdempotencyKey: key,
	})
	if err != nil {
		t.Fatal(err)
	}
	if schedule.ID != "sched-1" || schedule.Kind != "interval" {
		t.Fatalf("schedule = %+v", schedule)
	}
}

func TestStreamEvents(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/events" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "id: 1\nevent: work.created\ndata: {\"id\":1,\"type\":\"work.created\",\"tenantId\":\"acme\",\"workItemId\":\"wi-1\",\"at\":\"2026-10-20T00:00:00Z\"}\n\n")
	}))
	defer srv.Close()

	c := New(srv.URL, "acme", "default")
	resp, err := c.StreamEvents(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "event: work.created") {
		t.Fatalf("stream = %q", raw)
	}
}

func TestSymbivelaAttentionProjection(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/projections/symbivela-attention" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"attention": []SymbivelaAttention{{AttentionID: "att-1", WorkspaceID: "default", Kind: "approval", Status: "open", Priority: 2}},
		})
	}))
	defer srv.Close()
	c := New(srv.URL, "acme", "default")
	items, err := c.SymbivelaAttention(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].AttentionID != "att-1" || items[0].Priority != 2 {
		t.Fatalf("items = %+v", items)
	}
}

func TestConflictOverrideSDK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Idempotency-Key") == "" {
			http.Error(w, "missing idempotency key", http.StatusBadRequest)
			return
		}
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/work-items/wi-1/conflict-override":
			_ = json.NewEncoder(w).Encode(WorkItem{
				ID: "wi-1", State: "ready", Version: 2,
				ConflictOverride: &ConflictOverride{Reason: "planned overlap", ActorRef: "alice"},
			})
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/work-items/wi-1/conflict-override":
			if r.URL.Query().Get("expectedVersion") != "2" {
				http.Error(w, "bad version", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(WorkItem{ID: "wi-1", State: "ready", Version: 3})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := New(srv.URL, "acme", "default")
	item, err := c.OverrideConflicts(context.Background(), "wi-1", ConflictOverrideRequest{
		Reason: "planned overlap", TTLSeconds: 300, ExpectedVersion: 1, IdempotencyKey: "override-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if item.ConflictOverride == nil || item.ConflictOverride.Reason != "planned overlap" {
		t.Fatalf("item = %+v", item)
	}
	cleared, err := c.ClearConflictOverride(context.Background(), "wi-1", item.Version, "override-2")
	if err != nil {
		t.Fatal(err)
	}
	if cleared.Version != 3 {
		t.Fatalf("cleared = %+v", cleared)
	}
}

func TestLeaseSDK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/work-items/wi-1/leases":
			_ = json.NewEncoder(w).Encode(map[string]any{"leases": []Lease{{ID: "lease-1", Scope: "backend/internal", FencingToken: 2}}})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/work-items/wi-1/leases":
			if r.Header.Get("Idempotency-Key") == "" {
				http.Error(w, "missing idempotency key", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(Lease{ID: "lease-2", Scope: "frontend/src", FencingToken: 3})
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/work-items/wi-1/leases/lease-2":
			if r.URL.Query().Get("fencingToken") != "3" {
				http.Error(w, "bad token", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "released"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := New(srv.URL, "acme", "default")
	leases, err := c.ListLeases(context.Background(), "wi-1")
	if err != nil || len(leases) != 1 || leases[0].Scope != "backend/internal" {
		t.Fatalf("leases = %+v, err = %v", leases, err)
	}
	lease, err := c.AcquireLease(context.Background(), "wi-1", AcquireLeaseRequest{
		OwnerRef: "collaborator", Scope: "frontend/src", TTLSeconds: 300, IdempotencyKey: "lease-key",
	})
	if err != nil || lease.ID != "lease-2" {
		t.Fatalf("lease = %+v, err = %v", lease, err)
	}
	if err := c.ReleaseLease(context.Background(), "wi-1", lease.ID, lease.FencingToken, "release-key"); err != nil {
		t.Fatal(err)
	}
}
