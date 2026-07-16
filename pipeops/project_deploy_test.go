package pipeops

import (
	"context"
	"encoding/json"
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

func TestProjectServiceDeployUsesRedeployContract(t *testing.T) {
	t.Parallel()

	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.Header().Set("Content-Type", "application/json")

		switch requests {
		case 1:
			if r.Method != http.MethodGet {
				t.Fatalf("fetch method = %s, want GET", r.Method)
			}
			if r.URL.Path != "/project/fetch/p1" {
				t.Fatalf("fetch path = %s, want /project/fetch/p1", r.URL.Path)
			}
			writeProjectDeployResponse(t, w, `{
				"success": true,
				"data": {
					"project": {
						"Name": "ora-landing",
						"Username": "PipeOpsHQ",
						"Source": "github",
						"Repository": "https://github.com/PipeOpsHQ/ora-website",
						"Branch": "main",
						"Environment": "production",
						"ClusterUUID": "cluster-1",
						"RawLanguage": "Node",
						"BuildMethod": "railpack",
						"BuildCommand": "npm run build",
						"BuildPath": ".",
						"RunCommand": "npm start",
						"BuildVersion": "20",
						"BuildDirectory": "",
						"DockerPath": "",
						"Worker": false,
						"Replicas": 1,
						"PostStartCommand": "",
						"Job": false,
						"JobSuspended": false,
						"JobRunCommand": "",
						"JobRunInterval": "",
						"Configuration": {"settings": {"zdd": {"enabled": true}}},
						"Kind": "application"
					},
					"deployment": {
						"CommitSha": "abc123",
						"CommitURL": "https://example.test/commit/abc123"
					}
				}
			}`)
		case 2:
			if r.Method != http.MethodGet {
				t.Fatalf("network settings method = %s, want GET", r.Method)
			}
			if r.URL.Path != "/project/settings/network/p1" {
				t.Fatalf("network settings path = %s, want /project/settings/network/p1", r.URL.Path)
			}
			writeProjectDeployResponse(t, w, `{
				"success": true,
				"data": [{
					"UUID": "network-1",
					"Port": 3000,
					"Protocol": "HTTP",
					"AutoHTTPS": true,
					"Public": true,
					"Domains": []
				}]
			}`)
		case 3:
			if r.Method != http.MethodPost {
				t.Fatalf("deploy method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/project/redeploy/p1" {
				t.Fatalf("deploy path = %s, want /project/redeploy/p1", r.URL.Path)
			}
			if got := r.URL.Query().Get("action"); got != "deploy" {
				t.Fatalf("action = %q, want deploy", got)
			}

			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode redeploy body: %v", err)
			}
			for field, want := range map[string]string{
				"name":               "ora-landing",
				"username":           "PipeOpsHQ",
				"source":             "github",
				"repository":         "https://github.com/PipeOpsHQ/ora-website",
				"branch":             "main",
				"environment":        "production",
				"clusterUUID":        "cluster-1",
				"commitSha":          "abc123",
				"commitURL":          "https://example.test/commit/abc123",
				"repositoryLanguage": "railpack",
			} {
				got, ok := body[field].(string)
				if !ok || got != want {
					t.Fatalf("%s = %q, want %q", field, got, want)
				}
			}
			if _, ok := body["configuration"].(map[string]interface{}); !ok {
				t.Fatalf("configuration = %#v, want object", body["configuration"])
			}
			networkSettings, ok := body["networkSettings"].([]interface{})
			if !ok || len(networkSettings) != 1 {
				t.Fatalf("networkSettings = %#v, want one item", body["networkSettings"])
			}
			network, ok := networkSettings[0].(map[string]interface{})
			if !ok || network["UUID"] != "network-1" || network["Port"] != float64(3000) {
				t.Fatalf("networkSettings[0] = %#v, want fetched network", networkSettings[0])
			}
			buildSettings, ok := body["buildSettings"].(map[string]interface{})
			if !ok {
				t.Fatalf("buildSettings = %#v, want object", body["buildSettings"])
			}
			got, ok := buildSettings["buildMethod"].(string)
			if !ok || got != "railpack" {
				t.Fatalf("buildMethod = %q, want railpack", got)
			}
			w.WriteHeader(http.StatusAccepted)
			writeProjectDeployResponse(t, w, `{"success":true,"message":"Deployment queued"}`)
		default:
			t.Fatalf("unexpected request %d: %s %s", requests, r.Method, r.URL.String())
		}
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
	if requests != 3 {
		t.Fatalf("requests = %d, want 3", requests)
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

		switch requests {
		case 1:
			if r.URL.Path != "/project/fetch/p1" {
				t.Fatalf("fetch path = %s, want /project/fetch/p1", r.URL.Path)
			}
			writeProjectDeployResponse(t, w, `{
				"data": {
					"project": {
						"Name": "ora-landing",
						"Configuration": {"settings": {}}
					},
					"deployment": {}
				}
			}`)
		case 2:
			if r.URL.Path != "/project/settings/network/p1" {
				t.Fatalf("network settings path = %s, want /project/settings/network/p1", r.URL.Path)
			}
			writeProjectDeployResponse(t, w, `{"data":[{"UUID":"network-1","Port":3000,"Protocol":"HTTP"}]}`)
		case 3:
			if r.URL.Path != "/project/redeploy/p1" {
				t.Fatalf("deploy path = %s, want /project/redeploy/p1", r.URL.Path)
			}
			if got := r.URL.Query().Get("action"); got != "deploy" {
				t.Fatalf("action = %q, want deploy", got)
			}
			if got := r.URL.Query().Get("no_cache"); got != "true" {
				t.Fatalf("no_cache = %q, want true", got)
			}
			writeProjectDeployResponse(t, w, `{"success":true}`)
		default:
			t.Fatalf("unexpected request %d: %s", requests, r.URL.String())
		}
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
	if requests != 3 {
		t.Fatalf("requests = %d, want 3", requests)
	}
}

func TestProjectServiceDeployRejectsIncompleteSnapshot(t *testing.T) {
	t.Parallel()

	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		w.Header().Set("Content-Type", "application/json")
		writeProjectDeployResponse(t, w, `{"data":{"project":{"Name":"ora-landing"}}}`)
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.Projects.Deploy(context.Background(), "p1")
	if err == nil {
		t.Fatal("Projects.Deploy error = nil, want missing configuration error")
	}
	if got, want := err.Error(), "project snapshot is missing configuration"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
}

func TestProjectServiceDeployRejectsMissingNetworkSettings(t *testing.T) {
	t.Parallel()

	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.Header().Set("Content-Type", "application/json")

		switch requests {
		case 1:
			if r.URL.Path != "/project/fetch/p1" {
				t.Fatalf("fetch path = %s, want /project/fetch/p1", r.URL.Path)
			}
			writeProjectDeployResponse(t, w, `{
				"data": {
					"project": {
						"Name": "ora-landing",
						"Configuration": {"Settings": {}}
					}
				}
			}`)
		case 2:
			if r.URL.Path != "/project/settings/network/p1" {
				t.Fatalf("network settings path = %s, want /project/settings/network/p1", r.URL.Path)
			}
			writeProjectDeployResponse(t, w, `{"data":[]}`)
		default:
			t.Fatalf("unexpected request %d: %s", requests, r.URL.String())
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.Projects.Deploy(context.Background(), "p1")
	if err == nil {
		t.Fatal("Projects.Deploy error = nil, want missing network settings error")
	}
	if got, want := err.Error(), "project snapshot is missing network settings"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
	if requests != 2 {
		t.Fatalf("requests = %d, want 2", requests)
	}
}
