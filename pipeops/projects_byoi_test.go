package pipeops

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectService_DeployFromImage_UsesBYOIRoute(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/project/deploy-from-image" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/project/deploy-from-image")
		}

		var body map[string]any
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body error: %v", err)
		}
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Fatalf("unmarshal body error: %v", err)
		}
		if got := body["container_image"]; got != "docker.io/library/nginx" {
			t.Fatalf("container_image = %#v, want %q", got, "docker.io/library/nginx")
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"success":true,"message":"BYOI deployment initiated successfully","data":{"project_uuid":"p1","project_name":"nginx-demo","container_image":"docker.io/library/nginx:latest","image_tag":"latest","status":"pending","domain":"https://nginx-demo.example.com","build_sha":"sha123"}}`)); err != nil {
			t.Fatalf("write response error: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Projects.DeployFromImage(context.Background(), &DeployFromImageRequest{
		Name:            "nginx-demo",
		ContainerImage:  "docker.io/library/nginx",
		ImageTag:        "latest",
		Port:            80,
		Replicas:        1,
		VCPU:            0.5,
		Memory:          DeployFromImageMemory{Value: 512, Unit: "MB"},
		ClusterUUID:     "c1",
		EnvironmentUUID: "e1",
		WorkspaceUUID:   "w1",
	})
	if err != nil {
		t.Fatalf("DeployFromImage error: %v", err)
	}
	if resp.Data.ProjectUUID != "p1" {
		t.Fatalf("project_uuid = %q, want %q", resp.Data.ProjectUUID, "p1")
	}
}
