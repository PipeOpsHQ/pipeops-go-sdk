package pipeops

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectService_ListDeployments_UsesControllerRouteAndQuery(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/project/get-deployments/p1" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/project/get-deployments/p1")
		}
		if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
			t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
		}
		if got := r.URL.Query().Get("filterBy"); got != "git" {
			t.Fatalf("filterBy = %q, want %q", got, "git")
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Fatalf("page = %q, want %q", got, "2")
		}
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Fatalf("limit = %q, want %q", got, "5")
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"success":true,"message":"ok","data":[{"SHA":"abc123","CommitMessage":"ship it"}],"meta":{"total_pages":4,"current_page":2,"next_page":3,"current_count":1}}`)); err != nil {
			t.Fatalf("write response error: %v", err)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Projects.ListDeployments(context.Background(), "p1", &ProjectDeploymentListOptions{WorkspaceUUID: "w1", FilterBy: "git", Page: 2, Limit: 5})
	if err != nil {
		t.Fatalf("Projects.ListDeployments error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("deployments len = %d, want %d", len(resp.Data), 1)
	}
	if got := resp.Data[0]["SHA"]; got != "abc123" {
		t.Fatalf("SHA = %v, want %q", got, "abc123")
	}
	if resp.Meta.CurrentPage != 2 {
		t.Fatalf("current_page = %d, want %d", resp.Meta.CurrentPage, 2)
	}
}

func TestProjectService_ListDeploymentHistory_UsesControllerRouteAndPagination(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/project/deployment/p1" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/project/deployment/p1")
		}
		if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
			t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
		}
		if got := r.URL.Query().Get("page"); got != "3" {
			t.Fatalf("page = %q, want %q", got, "3")
		}
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Fatalf("limit = %q, want %q", got, "10")
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"success":true,"message":"ok","data":[{"UUID":"d1","CommitSha":"abc123","Status":"deployed","DurationSeconds":12}],"meta":{"total_pages":5,"current_page":3,"next_page":4,"current_count":1}}`)); err != nil {
			t.Fatalf("write response error: %v", err)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Projects.ListDeploymentHistory(context.Background(), "p1", &ProjectDeploymentHistoryOptions{WorkspaceID: "w1", Page: 3, Limit: 10})
	if err != nil {
		t.Fatalf("Projects.ListDeploymentHistory error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("history len = %d, want %d", len(resp.Data), 1)
	}
	if got := resp.Data[0]["UUID"]; got != "d1" {
		t.Fatalf("UUID = %v, want %q", got, "d1")
	}
	if resp.Meta.NextPage != 4 {
		t.Fatalf("next_page = %d, want %d", resp.Meta.NextPage, 4)
	}
}
