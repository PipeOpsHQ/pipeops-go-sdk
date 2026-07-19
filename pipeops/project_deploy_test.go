package pipeops

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func writeProjectDeployResponse(t *testing.T, w http.ResponseWriter, body string) {
	t.Helper()
	if _, err := w.Write([]byte(body)); err != nil {
		t.Errorf("write response: %v", err)
	}
}

func TestProjectServiceDeployUsesThinRedeployContract(t *testing.T) {
	t.Parallel()

	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			t.Fatalf("deploy method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/project/redeploy/p1" {
			t.Fatalf("deploy path = %s, want /project/redeploy/p1", r.URL.Path)
		}
		if got := r.URL.Query().Get("action"); got != "deploy" {
			t.Fatalf("action = %q, want deploy", got)
		}

		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var body map[string]interface{}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatalf("decode redeploy body: %v", err)
		}
		// Prefer-client: empty/thin body — no snapshot fields required.
		if _, ok := body["name"]; ok {
			t.Fatalf("thin redeploy should omit name, got %#v", body["name"])
		}
		if _, ok := body["networkSettings"]; ok {
			t.Fatalf("thin redeploy should omit networkSettings, got %#v", body["networkSettings"])
		}
		if _, ok := body["configuration"]; ok {
			t.Fatalf("thin redeploy should omit configuration, got %#v", body["configuration"])
		}

		w.WriteHeader(http.StatusAccepted)
		writeProjectDeployResponse(t, w, `{"success":true,"message":"Deployment queued"}`)
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, err := client.Projects.Deploy(context.Background(), "p1")
	if err != nil {
		t.Fatalf("Projects.Deploy error: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusAccepted)
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want 1 (no snapshot fetch)", requests)
	}
}

func TestProjectServiceDeployScopesRequestsToWorkspace(t *testing.T) {
	t.Parallel()

	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.Header().Set("Content-Type", "application/json")
		if got := r.URL.Query().Get("workspace_uuid"); got != "workspace-1" {
			t.Fatalf("workspace_uuid = %q, want workspace-1", got)
		}
		if r.URL.Path != "/project/redeploy/p1" {
			t.Fatalf("deploy path = %s, want /project/redeploy/p1", r.URL.Path)
		}
		if got := r.URL.Query().Get("action"); got != "deploy" {
			t.Fatalf("action = %q, want deploy", got)
		}
		if got := r.URL.Query().Get("no_cache"); got != "true" {
			t.Fatalf("no_cache = %q, want true", got)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got := body["workspace_uuid"]; got != "workspace-1" {
			t.Fatalf("body workspace_uuid = %v, want workspace-1", got)
		}

		writeProjectDeployResponse(t, w, `{"success":true}`)
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	if _, err := client.Projects.Deploy(context.Background(), "p1", &ProjectDeployOptions{
		WorkspaceUUID: "workspace-1",
		NoCache:       true,
	}); err != nil {
		t.Fatalf("Projects.Deploy error: %v", err)
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
}

func TestProjectServiceDeployRejectsEmptyProjectUUID(t *testing.T) {
	t.Parallel()

	client, err := NewClient("https://api.pipeops.test")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.Projects.Deploy(context.Background(), "  ")
	if err == nil {
		t.Fatal("Projects.Deploy error = nil, want empty UUID error")
	}
	if got, want := err.Error(), "project UUID cannot be empty"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
}
