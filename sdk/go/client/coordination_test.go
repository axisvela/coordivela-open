package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCoordinationPlanAndClaims(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.Method+" "+r.URL.RequestURI())
		if r.Header.Get("X-Axis-Tenant") != "acme" || r.Header.Get("X-Axis-Workspace") != "ws-1" {
			t.Errorf("headers=%v", r.Header)
		}
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/coordination/plan"):
			_ = json.NewEncoder(w).Encode(CoordinationPlan{Objective: "ship"})
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/coordination/plan"):
			if r.Header.Get("Idempotency-Key") == "" {
				t.Error("missing idempotency key")
			}
			var req SaveCoordinationPlanRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(CoordinationPlan{Objective: req.Objective, Version: 1})
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/coordination/claims"):
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(CoordinationClaim{ID: "claim-1", Path: "backend/runtime"})
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/coordination/claims"):
			_ = json.NewEncoder(w).Encode(map[string]any{"claims": []CoordinationClaim{{ID: "claim-1"}}})
		case r.Method == http.MethodDelete:
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "released"})
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(server.URL, "acme", "ws-1")
	ctx := context.Background()
	if _, err := client.GetCoordinationPlan(ctx, "wi-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.SaveCoordinationPlan(ctx, "wi-1", SaveCoordinationPlanRequest{Objective: "ship", FencingToken: 2}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.AddCoordinationClaim(ctx, "wi-1", AddCoordinationClaimRequest{Repository: "repo", Path: "backend/runtime", Mode: "write", FencingToken: 2}); err != nil {
		t.Fatal(err)
	}
	claims, err := client.ListCoordinationClaims(ctx, "wi-1")
	if err != nil || len(claims) != 1 {
		t.Fatalf("claims=%+v err=%v", claims, err)
	}
	if err := client.ReleaseCoordinationClaim(ctx, "wi-1", "claim-1", 2); err != nil {
		t.Fatal(err)
	}
	if len(paths) != 5 {
		t.Fatalf("paths=%v", paths)
	}
}

func TestCoordinationFindingsAndDirectives(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/coordination/findings"):
			_ = json.NewEncoder(w).Encode(map[string]any{"findings": []CoordinationFinding{{ID: "finding-1", Kind: "path_conflict"}}})
		case strings.HasSuffix(r.URL.Path, "/coordination/directives"):
			_ = json.NewEncoder(w).Encode(map[string]any{"directives": []CoordinationDirective{{ID: "directive-1", Kind: "avoid_path"}}})
		case strings.HasSuffix(r.URL.Path, "/outcome"):
			var req DirectiveOutcomeRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(CoordinationDirective{ID: "directive-1", Status: "applied", Outcome: &DirectiveOutcome{Kind: req.Outcome}})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(server.URL, "acme", "ws-1")
	ctx := context.Background()
	findings, err := client.ListCoordinationFindings(ctx, "wi-1")
	if err != nil || len(findings) != 1 || findings[0].Kind != "path_conflict" {
		t.Fatalf("findings=%+v err=%v", findings, err)
	}
	directives, err := client.ListCoordinationDirectives(ctx, "wi-1")
	if err != nil || len(directives) != 1 {
		t.Fatalf("directives=%+v err=%v", directives, err)
	}
	updated, err := client.RecordCoordinationDirectiveOutcome(ctx, "wi-1", "directive-1", DirectiveOutcomeRequest{Outcome: "effective"})
	if err != nil || updated.Status != "applied" || updated.Outcome == nil {
		t.Fatalf("updated=%+v err=%v", updated, err)
	}
}
