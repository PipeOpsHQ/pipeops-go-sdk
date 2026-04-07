package pipeops

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerService_List_UsesClusterWorkspaceEndpoint(t *testing.T) {
	t.Parallel()

	wantWorkspaceUUID := "workspace-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/cluster" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/cluster")
		}
		if got := r.URL.Query().Get("workspace_uuid"); got != wantWorkspaceUUID {
			t.Fatalf("workspace_uuid = %q, want %q", got, wantWorkspaceUUID)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"success": true,
			"message": "ok",
			"data": {
				"clusters": [
					{
						"Cluster": {
							"uuid": "c1",
							"name": "cluster-1",
							"cloudProvider": "aws",
							"region": "us-east-1",
							"status": "ready"
						},
						"IsActive": true,
						"InUse": true
					}
				]
			}
		}`)); err != nil {
			t.Fatalf("write response error: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Servers.List(context.Background(), wantWorkspaceUUID)
	if err != nil {
		t.Fatalf("Servers.List error: %v", err)
	}
	if resp.Status != "success" {
		t.Fatalf("Status = %q, want %q", resp.Status, "success")
	}
	if len(resp.Data.Servers) != 1 {
		t.Fatalf("len(Servers) = %d, want 1", len(resp.Data.Servers))
	}

	got := resp.Data.Servers[0]
	if got.UUID != "c1" || got.Name != "cluster-1" || got.Provider != "aws" || got.Region != "us-east-1" || got.Status != "ready" {
		t.Fatalf("server = %+v, want mapped cluster fields", got)
	}
}

func TestServerService_Get_AcceptsSingleClusterObject(t *testing.T) {
	t.Parallel()

	wantClusterUUID := "c1"
	wantWorkspaceUUID := "workspace-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/cluster/"+wantClusterUUID {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/cluster/"+wantClusterUUID)
		}
		if got := r.URL.Query().Get("workspace_uuid"); got != wantWorkspaceUUID {
			t.Fatalf("workspace_uuid = %q, want %q", got, wantWorkspaceUUID)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"success": true,
			"message": "ok",
			"data": {
				"id": 42,
				"uuid": "c1",
				"name": "cluster-1",
				"cloudProvider": "aws",
				"region": "us-east-1",
				"status": "ready",
				"WorkspaceID": 99
			}
		}`)); err != nil {
			t.Fatalf("write response error: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Servers.Get(context.Background(), wantClusterUUID, wantWorkspaceUUID)
	if err != nil {
		t.Fatalf("Servers.Get error: %v", err)
	}
	if resp.Status != "success" {
		t.Fatalf("Status = %q, want %q", resp.Status, "success")
	}

	got := resp.Data.Server
	if got.ID != "42" || got.UUID != "c1" || got.Name != "cluster-1" || got.Provider != "aws" || got.Region != "us-east-1" || got.Status != "ready" || got.WorkspaceID != "99" {
		t.Fatalf("server = %+v, want mapped single-cluster object fields", got)
	}
}

func TestServerService_Get_UsesClusterWorkspaceEndpoint(t *testing.T) {
	t.Parallel()

	wantClusterUUID := "c1"
	wantWorkspaceUUID := "workspace-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/cluster/"+wantClusterUUID {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/cluster/"+wantClusterUUID)
		}
		if got := r.URL.Query().Get("workspace_uuid"); got != wantWorkspaceUUID {
			t.Fatalf("workspace_uuid = %q, want %q", got, wantWorkspaceUUID)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"success": true,
			"message": "ok",
			"data": {
				"clusters": [
					{
						"Cluster": {
							"uuid": "c1",
							"name": "cluster-1",
							"cloudProvider": "aws",
							"region": "us-east-1",
							"status": "ready"
						},
						"IsActive": true
					}
				]
			}
		}`)); err != nil {
			t.Fatalf("write response error: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Servers.Get(context.Background(), wantClusterUUID, wantWorkspaceUUID)
	if err != nil {
		t.Fatalf("Servers.Get error: %v", err)
	}
	if resp.Status != "success" {
		t.Fatalf("Status = %q, want %q", resp.Status, "success")
	}

	got := resp.Data.Server
	if got.UUID != "c1" || got.Name != "cluster-1" || got.Provider != "aws" || got.Region != "us-east-1" || got.Status != "ready" {
		t.Fatalf("server = %+v, want mapped cluster fields", got)
	}
}

func TestServerService_Create_UsesServerCreateEndpoint(t *testing.T) {
	t.Parallel()

	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/server/create" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/server/create")
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body error: %v", err)
		}
		if err := json.Unmarshal(bodyBytes, &gotBody); err != nil {
			t.Fatalf("unmarshal body error: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"status": "success",
			"message": "created",
			"data": {
				"server": {
					"uuid": "s1",
					"name": "my-server",
					"provider": "aws",
					"region": "us",
					"status": "active"
				}
			}
		}`)); err != nil {
			t.Fatalf("write response error: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Servers.Create(context.Background(), "ignored-cluster-uuid", &CreateServerRequest{
		Name:     "my-server",
		Provider: "aws",
		Region:   "us",
	})
	if err != nil {
		t.Fatalf("Servers.Create error: %v", err)
	}

	if got := gotBody["server_name"]; got != "my-server" {
		t.Fatalf("server_name = %#v, want %q", got, "my-server")
	}
	if got := gotBody["server_cloud"]; got != "aws" {
		t.Fatalf("server_cloud = %#v, want %q", got, "aws")
	}
	if got := gotBody["server_region"]; got != "us" {
		t.Fatalf("server_region = %#v, want %q", got, "us")
	}

	if resp.Data.Server.UUID != "s1" {
		t.Fatalf("resp.Data.Server.UUID = %q, want %q", resp.Data.Server.UUID, "s1")
	}
}

func TestServerService_Delete_FallsBackOn404(t *testing.T) {
	t.Parallel()

	call := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		call++
		switch call {
		case 1:
			if r.Method != http.MethodDelete {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodDelete)
			}
			if r.URL.Path != "/api/v1/clusters/c1" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/api/v1/clusters/c1")
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte(`{"message":"not found"}`)); err != nil {
				t.Fatalf("write response error: %v", err)
			}
		case 2:
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/workspace" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/workspace")
			}
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte(`{"data":[{"UUID":"w1"}],"message":"ok","success":true}`)); err != nil {
				t.Fatalf("write response error: %v", err)
			}
		case 3:
			if r.Method != http.MethodDelete {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodDelete)
			}
			if r.URL.Path != "/cluster/c1" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/cluster/c1")
			}
			if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
				t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte(`{"message":"not found"}`)); err != nil {
				t.Fatalf("write response error: %v", err)
			}
		case 4:
			if r.Method != http.MethodDelete {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodDelete)
			}
			if r.URL.Path != "/clusters/c1/servers/s1" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/clusters/c1/servers/s1")
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected call %d", call)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.Servers.Delete(context.Background(), "c1", "s1")
	if err != nil {
		t.Fatalf("Servers.Delete error: %v", err)
	}
}
