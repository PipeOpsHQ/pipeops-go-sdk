package pipeops

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateEnvVariables_SendsCamelCaseAndMergeQuery(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/project/settings/env/p1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("merge"); got != "true" {
			t.Fatalf("merge = %q, want true", got)
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		var body map[string]interface{}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatal(err)
		}
		if _, ok := body["env_variables"]; ok {
			t.Fatalf("must not send snake_case env_variables: %s", raw)
		}
		envs, ok := body["envVariables"].([]interface{})
		if !ok || len(envs) != 1 {
			t.Fatalf("envVariables = %#v", body["envVariables"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"message":"ok","data":[{"key":"PORT","value":"8080"}]}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp, _, err := client.Projects.UpdateEnvVariables(context.Background(), "p1", &EnvVariablesRequest{
		EnvVariables:  []EnvVariable{{Key: "PORT", Value: "8080"}},
		Merge:         true,
		WorkspaceUUID: "ws-1", // skip workspace discovery GETs
	})
	if err != nil {
		t.Fatalf("UpdateEnvVariables: %v", err)
	}
	if len(resp.Data.EnvVariables) != 1 || resp.Data.EnvVariables[0].Value != "8080" {
		t.Fatalf("response data: %+v", resp.Data.EnvVariables)
	}
}

func TestEnvVariablesData_UnmarshalBareArray(t *testing.T) {
	t.Parallel()
	var d EnvVariablesData
	if err := json.Unmarshal([]byte(`[{"key":"A","value":"1"}]`), &d); err != nil {
		t.Fatal(err)
	}
	if len(d.EnvVariables) != 1 || d.EnvVariables[0].Key != "A" {
		t.Fatalf("%+v", d.EnvVariables)
	}
}
