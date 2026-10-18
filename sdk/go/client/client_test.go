package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
