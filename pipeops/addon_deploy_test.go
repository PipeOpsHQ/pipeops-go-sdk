package pipeops

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddOnServiceDeploy_SendsNestedAndThinShape(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/addons/deploy" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		var body map[string]interface{}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatal(err)
		}
		if body["Workspace"] != "ws-1" || body["Server"] != "cluster-1" {
			t.Fatalf("workspace/server: %#v", body)
		}
		if body["id"] != "addon-1" {
			t.Fatalf("id alias: %#v", body["id"])
		}
		dep, ok := body["Deployment"].(map[string]interface{})
		if !ok || dep["ID"] != "addon-1" {
			t.Fatalf("Deployment: %#v", body["Deployment"])
		}
		if body["Environment"] != "env-1" {
			t.Fatalf("Environment: %#v", body["Environment"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"message":"deploying","data":{"deployment":{"UID":"dep-1"}}}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp, _, err := client.AddOns.Deploy(context.Background(), &DeployAddOnRequest{
		ID:          "addon-1",
		Server:      "cluster-1",
		Workspace:   "ws-1",
		Environment: "env-1",
		Config:      map[string]interface{}{"Env": map[string]string{"A": "1"}},
	})
	if err != nil {
		t.Fatalf("Deploy: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
}

func TestAddOnServiceDeploy_RequiresServer(t *testing.T) {
	t.Parallel()
	client, err := NewClient("https://api.pipeops.test")
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = client.AddOns.Deploy(context.Background(), &DeployAddOnRequest{
		ID:        "addon-1",
		Workspace: "ws-1",
	})
	if err == nil {
		t.Fatal("expected server required error")
	}
}
