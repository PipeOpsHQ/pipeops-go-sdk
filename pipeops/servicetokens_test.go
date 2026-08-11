package pipeops

import (
	"encoding/json"
	"testing"
)

func TestServiceAccountTokenResponse_UnmarshalCreateFlatToken(t *testing.T) {
	raw := []byte(`{
		"success": true,
		"message": "created",
		"data": {
			"id": "tok-1",
			"name": "readonly",
			"token": "sat_secret_once",
			"token_prefix": "sat_sec",
			"workspace_id": "ws-1",
			"scopes": ["api:read"]
		}
	}`)
	var resp ServiceAccountTokenResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Data.Token.UUID != "tok-1" {
		t.Fatalf("uuid = %q", resp.Data.Token.UUID)
	}
	if resp.Data.Token.Token != "sat_secret_once" {
		t.Fatalf("token = %q", resp.Data.Token.Token)
	}
	if resp.Data.Token.Name != "readonly" {
		t.Fatalf("name = %q", resp.Data.Token.Name)
	}
}

func TestServiceAccountTokenResponse_UnmarshalNestedToken(t *testing.T) {
	raw := []byte(`{
		"status": "success",
		"data": {
			"token": {
				"uuid": "tok-2",
				"name": "nested",
				"workspace_id": "ws-2"
			}
		}
	}`)
	var resp ServiceAccountTokenResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Data.Token.UUID != "tok-2" || resp.Data.Token.Name != "nested" {
		t.Fatalf("token = %+v", resp.Data.Token)
	}
}

func TestDeploymentSessionResponse_UnmarshalArray(t *testing.T) {
	raw := []byte(`{
		"success": true,
		"message": "addon deployment overview",
		"data": [
			{"UID": "dep-1", "Name": "redis"},
			{"UID": "dep-2", "Name": "postgres"}
		]
	}`)
	var resp DeploymentSessionResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Deployments) != 2 {
		t.Fatalf("deployments = %d, want 2", len(resp.Deployments))
	}
	if resp.Deployments[0]["UID"] != "dep-1" {
		t.Fatalf("first = %+v", resp.Deployments[0])
	}
}
