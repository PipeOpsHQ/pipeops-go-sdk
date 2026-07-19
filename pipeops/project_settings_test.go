package pipeops

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectServiceUpdateDeploySettings_ThinBody(t *testing.T) {
	t.Parallel()

	auto := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/project/settings/deploy/p1" {
			t.Fatalf("path = %s, want /project/settings/deploy/p1", r.URL.Path)
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var body map[string]interface{}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if _, ok := body["branch"]; ok {
			t.Fatalf("thin body should omit branch, got %#v", body)
		}
		if _, ok := body["repository"]; ok {
			t.Fatalf("thin body should omit repository, got %#v", body)
		}
		if v, ok := body["autoDeployEnabled"].(bool); !ok || v {
			t.Fatalf("autoDeployEnabled = %#v, want false", body["autoDeployEnabled"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"message":"ok","data":{"autoDeployEnabled":false,"branch":"main"}}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	resp, _, err := client.Projects.UpdateDeploySettings(context.Background(), "p1", &DeploySettingsRequest{
		AutoDeployEnabled: &auto,
		WorkspaceUUID:     "ws-1",
	})
	if err != nil {
		t.Fatalf("UpdateDeploySettings: %v", err)
	}
	if resp == nil || resp.Data.Branch != "main" {
		t.Fatalf("response: %+v", resp)
	}
}

func TestProjectServiceUpdateDeploySettings_RequiresUUID(t *testing.T) {
	t.Parallel()
	client, err := NewClient("http://example.invalid")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	_, _, err = client.Projects.UpdateDeploySettings(context.Background(), "  ", nil)
	if err == nil {
		t.Fatal("expected empty UUID error")
	}
}

func TestProjectServiceUpdateSecurityPolicy_PartialBody(t *testing.T) {
	t.Parallel()

	enabled := true
	maxCritical := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/project/settings/security-policy/p1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var body map[string]interface{}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(body) != 2 {
			t.Fatalf("expected only enabled+maxCritical keys, got %#v", body)
		}
		if v, ok := body["enabled"].(bool); !ok || !v {
			t.Fatalf("enabled = %#v", body["enabled"])
		}
		if v, ok := body["maxCritical"].(float64); !ok || v != 0 {
			t.Fatalf("maxCritical = %#v", body["maxCritical"])
		}
		if _, ok := body["maxHigh"]; ok {
			t.Fatalf("maxHigh should be omitted: %#v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"message":"ok","data":{"securityPolicy":{"enabled":true,"maxCritical":0}}}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	resp, _, err := client.Projects.UpdateSecurityPolicy(context.Background(), "p1", &SecurityPolicyRequest{
		Enabled:       &enabled,
		MaxCritical:   &maxCritical,
		WorkspaceUUID: "ws-1",
	})
	if err != nil {
		t.Fatalf("UpdateSecurityPolicy: %v", err)
	}
	if resp == nil || !resp.Success {
		t.Fatalf("response: %+v", resp)
	}
}
