package pipeops

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuditLogService_ListProject(t *testing.T) {
	t.Parallel()

	const projectUUID = "proj-abc"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		wantPath := "/project/audit-logs/" + projectUUID
		if r.URL.Path != wantPath {
			t.Fatalf("path = %s, want %s", r.URL.Path, wantPath)
		}
		if got := r.URL.Query().Get("limit"); got != "25" {
			t.Fatalf("limit = %q, want 25", got)
		}
		if got := r.URL.Query().Get("offset"); got != "5" {
			t.Fatalf("offset = %q, want 5", got)
		}
		if got := r.URL.Query().Get("action"); got != "project.redeploy,project.env.update" {
			t.Fatalf("action = %q", got)
		}
		if got := r.URL.Query().Get("actor_type"); got != "user" {
			t.Fatalf("actor_type = %q", got)
		}
		if got := r.URL.Query().Get("category"); got != "lifecycle" {
			t.Fatalf("category = %q", got)
		}
		if got := r.URL.Query().Get("from"); !strings.HasPrefix(got, "2026-01-01") {
			t.Fatalf("from = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "fetched project audit logs successfully",
			"data": []map[string]interface{}{
				{
					"uuid":         "log-1",
					"action":       "project.redeploy",
					"action_label": "Redeployed",
					"category":     "lifecycle",
					"status":       "success",
					"summary":      "Redeployed faulty-art",
					"project_uuid": projectUUID,
					"project_name": "faulty-art",
					"created_at":   "2026-08-10T12:00:00Z",
					"actor": map[string]interface{}{
						"type": "user",
						"uuid": "user-1",
						"name": "Ada",
					},
				},
			},
			"pagination": map[string]interface{}{"total": 1, "limit": 25, "offset": 5},
		})
	}))
	defer srv.Close()

	client, err := NewClient(srv.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	out, resp, err := client.AuditLogs.ListProject(context.Background(), projectUUID, &ProjectAuditLogListOptions{
		Action:    JoinAuditActions("project.redeploy", "project.env.update"),
		ActorType: "user",
		Category:  "lifecycle",
		From:      "2026-01-01T00:00:00Z",
		Limit:     25,
		Offset:    5,
	})
	if err != nil {
		t.Fatalf("ListProject: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if !out.Success || len(out.Data) != 1 {
		t.Fatalf("out = %+v", out)
	}
	if out.Data[0].Action != "project.redeploy" || out.Data[0].Actor.Name != "Ada" {
		t.Fatalf("entry = %+v", out.Data[0])
	}
	if out.Pagination.Total != 1 || out.Pagination.Limit != 25 {
		t.Fatalf("pagination = %+v", out.Pagination)
	}
}

func TestAuditLogService_ListProjectRequiresUUID(t *testing.T) {
	t.Parallel()
	client, err := NewClient("https://api.pipeops.test")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	_, _, err = client.AuditLogs.ListProject(context.Background(), "  ", nil)
	if err == nil || !strings.Contains(err.Error(), "project_uuid is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestAuditLogService_ListWorkspace(t *testing.T) {
	t.Parallel()

	const workspaceUUID = "ws-xyz"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/project/workspace-audit-logs" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("workspace_uuid"); got != workspaceUUID {
			t.Fatalf("workspace_uuid = %q", got)
		}
		if got := r.URL.Query().Get("project_uuid"); got != "proj-1" {
			t.Fatalf("project_uuid = %q", got)
		}
		if got := r.URL.Query().Get("search"); got != "redeploy" {
			t.Fatalf("search = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"success": true,
			"message": "fetched workspace audit logs successfully",
			"data": [{
				"uuid": "log-2",
				"action": "project.env.update",
				"action_label": "Updated environment variables",
				"category": "settings",
				"status": "success",
				"project_uuid": "proj-1",
				"project_name": "api",
				"actor": {"type": "agent", "label": "cortex"}
			}],
			"pagination": {"total": 12, "limit": 20, "offset": 0}
		}`))
	}))
	defer srv.Close()

	client, err := NewClient(srv.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	out, _, err := client.AuditLogs.ListWorkspace(context.Background(), &WorkspaceAuditLogListOptions{
		WorkspaceUUID: workspaceUUID,
		ProjectUUID:   "proj-1",
		Search:        "redeploy",
		Limit:         20,
	})
	if err != nil {
		t.Fatalf("ListWorkspace: %v", err)
	}
	if len(out.Data) != 1 || out.Data[0].Actor.Type != "agent" {
		t.Fatalf("out = %+v", out)
	}
	if out.Pagination.Total != 12 {
		t.Fatalf("total = %d", out.Pagination.Total)
	}
}

func TestAuditLogService_ListWorkspaceRequiresWorkspace(t *testing.T) {
	t.Parallel()

	// No workspaces returned → must error instead of calling a broken path.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/workspace" || r.URL.Path == "/workspaces" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"success":true,"data":[]}`))
			return
		}
		t.Fatalf("unexpected path %s", r.URL.Path)
	}))
	defer srv.Close()

	client, err := NewClient(srv.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	_, _, err = client.AuditLogs.ListWorkspace(context.Background(), &WorkspaceAuditLogListOptions{})
	if err == nil || !strings.Contains(err.Error(), "workspace_uuid is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestJoinAuditActions(t *testing.T) {
	t.Parallel()
	if got := JoinAuditActions(" project.redeploy ", "", "project.pause"); got != "project.redeploy,project.pause" {
		t.Fatalf("got %q", got)
	}
}
