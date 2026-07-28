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

func TestWorkspace_Unmarshal_BillingEmail(t *testing.T) {
	t.Parallel()

	var workspace Workspace
	if err := json.Unmarshal([]byte(`{"UUID":"w1","Name":"ws","BillingEmail":"billing@example.com"}`), &workspace); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}
	if workspace.BillingEmail != "billing@example.com" {
		t.Fatalf("BillingEmail = %q, want billing@example.com", workspace.BillingEmail)
	}

	var snake Workspace
	if err := json.Unmarshal([]byte(`{"uuid":"w1","billing_email":"snake@example.com"}`), &snake); err != nil {
		t.Fatalf("json.Unmarshal snake error: %v", err)
	}
	if snake.BillingEmail != "snake@example.com" {
		t.Fatalf("BillingEmail = %q, want snake@example.com", snake.BillingEmail)
	}
}

func TestWorkspaceService_SetBillingEmail_UsesPascalCaseJSON(t *testing.T) {
	t.Parallel()

	var body map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/workspace/w1/add-billing-email" {
			t.Fatalf("path = %s, want /workspace/w1/add-billing-email", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"success":true,"message":"ok"}`)); err != nil {
			t.Fatalf("write: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	if _, err := client.Workspaces.SetBillingEmail(context.Background(), "w1", &SetBillingEmailRequest{
		Email: "billing@example.com",
	}); err != nil {
		t.Fatalf("SetBillingEmail error: %v", err)
	}
	if got, _ := body["BillingEmail"].(string); got != "billing@example.com" {
		t.Fatalf("body = %#v, want BillingEmail=billing@example.com", body)
	}
	if _, hasEmail := body["email"]; hasEmail {
		t.Fatalf("body still has lowercase email key: %#v", body)
	}
}
