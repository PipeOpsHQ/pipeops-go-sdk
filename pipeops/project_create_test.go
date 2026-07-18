package pipeops

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectService_Create_ControllerContract(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := NewClient(server.URL, WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}

	mux.HandleFunc("/workspace", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    []map[string]string{{"UUID": "ws-1", "uuid": "ws-1"}},
		}); err != nil {
			t.Fatal(err)
		}
	})

	var gotBody map[string]interface{}
	mux.HandleFunc("/project/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(raw, &gotBody); err != nil {
			t.Fatal(err)
		}
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "ok",
			"data": map[string]interface{}{
				"project": map[string]interface{}{
					"UUID": "proj-1",
					"Name": "my-app",
				},
			},
		}); err != nil {
			t.Fatal(err)
		}
	})

	worker := false
	enable := false
	suspended := false
	resp, _, err := client.Projects.Create(context.Background(), &CreateProjectRequest{
		Name:               "my-app",
		Username:           "acme",
		Source:             "github",
		Repository:         "https://github.com/acme/app",
		CommitURL:          "https://github.com/acme/app/commit/abc",
		CommitSha:          "abc123",
		RepositoryLanguage: "nodejs",
		Branch:             "main",
		Environment:        "development",
		EnvironmentUUID:    "env-1",
		ClusterUUID:        "cluster-1",
		// WorkspaceUUID intentionally omitted → filled from /workspace
		BuildSettings: CreateProjectBuildSettings{
			Type:         "user",
			BuildMethod:  "nodejs",
			BuildCommand: "npm run build",
			RunCommand:   "npm start",
			Worker:       &worker,
		},
		JobDetails: CreateProjectJobDetails{
			Enable:    &enable,
			Suspended: &suspended,
		},
		NetworkSettings: []CreateProjectNetworkSetting{
			{Port: 3000, Protocol: "HTTP"},
		},
		EnvVariables: []CreateProjectEnvVar{
			{Key: "PORT", Value: "3000"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data.Project.UUID != "proj-1" {
		t.Fatalf("uuid = %q", resp.Data.Project.UUID)
	}

	// Assert controller field names (not the old server_id / environment_id / build_command shape).
	for _, key := range []string{
		"name", "username", "source", "repository", "commitURL", "commitSha",
		"repositoryLanguage", "branch", "environment", "environment_uuid",
		"clusterUUID", "workspace_uuid", "buildSettings", "envVariables", "networkSettings",
	} {
		if _, ok := gotBody[key]; !ok {
			t.Fatalf("missing JSON key %q in body: %v", key, gotBody)
		}
	}
	if _, ok := gotBody["server_id"]; ok {
		t.Fatal("legacy server_id must not be sent")
	}
	if _, ok := gotBody["environment_id"]; ok {
		t.Fatal("legacy environment_id must not be sent")
	}
	if gotBody["workspace_uuid"] != "ws-1" {
		t.Fatalf("workspace_uuid = %v, want ws-1", gotBody["workspace_uuid"])
	}
	if gotBody["clusterUUID"] != "cluster-1" {
		t.Fatalf("clusterUUID = %v", gotBody["clusterUUID"])
	}
	bs, ok := gotBody["buildSettings"].(map[string]interface{})
	if !ok {
		t.Fatalf("buildSettings type = %T", gotBody["buildSettings"])
	}
	if bs["buildCommand"] != "npm run build" || bs["runCommand"] != "npm start" {
		t.Fatalf("buildSettings = %+v", bs)
	}
	nets, ok := gotBody["networkSettings"].([]interface{})
	if !ok || len(nets) != 1 {
		t.Fatalf("networkSettings = %+v", gotBody["networkSettings"])
	}
	net0, ok := nets[0].(map[string]interface{})
	if !ok {
		t.Fatalf("network entry type = %T", nets[0])
	}
	if net0["Port"] != float64(3000) || net0["Protocol"] != "HTTP" {
		t.Fatalf("network entry = %+v", net0)
	}
}
