package pipeops

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWorkspace_Unmarshal_AllowsNumericID(t *testing.T) {
	t.Parallel()

	var workspace Workspace
	if err := json.Unmarshal([]byte(`{"ID":0,"UUID":"w1","Name":"ws"}`), &workspace); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}
	if workspace.ID != "0" {
		t.Fatalf("workspace.ID = %q, want %q", workspace.ID, "0")
	}
	if workspace.UUID != "w1" {
		t.Fatalf("workspace.UUID = %q, want %q", workspace.UUID, "w1")
	}
	if workspace.Name != "ws" {
		t.Fatalf("workspace.Name = %q, want %q", workspace.Name, "ws")
	}
}

func TestWorkspaceService_List_Unmarshal_AllowsNumericID(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/workspace" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/workspace")
		}

		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write([]byte(`{"data":[{"ID":0,"UUID":"w1","Name":"ws"}],"message":"ok","success":true}`)); writeErr != nil {
			t.Fatalf("write response error: %v", writeErr)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	workspaces, _, err := client.Workspaces.List(context.Background())
	if err != nil {
		t.Fatalf("Workspaces.List error: %v", err)
	}
	if len(workspaces.Data.Workspaces) != 1 {
		t.Fatalf("workspaces len = %d, want %d", len(workspaces.Data.Workspaces), 1)
	}
	if got := workspaces.Data.Workspaces[0].ID; got != "0" {
		t.Fatalf("workspace.ID = %q, want %q", got, "0")
	}
}
