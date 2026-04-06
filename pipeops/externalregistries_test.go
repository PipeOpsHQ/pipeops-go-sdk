package pipeops

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExternalRegistryService_Create_List_Get_Delete(t *testing.T) {
	t.Parallel()

	step := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		step++
		w.Header().Set("Content-Type", "application/json")
		switch step {
		case 1:
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
			}
			if r.URL.Path != "/api/v1/external-registry" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/api/v1/external-registry")
			}
			if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
				t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
			}
			var body map[string]any
			bytes, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read body error: %v", err)
			}
			if err := json.Unmarshal(bytes, &body); err != nil {
				t.Fatalf("unmarshal body error: %v", err)
			}
			if got := body["name"]; got != "Docker Hub" {
				t.Fatalf("name = %#v, want %q", got, "Docker Hub")
			}
			_, _ = w.Write([]byte(`{"success":true,"data":{"id":1,"uid":"reg1","name":"Docker Hub","type":"dockerhub"}}`))
		case 2:
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/api/v1/external-registry" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/api/v1/external-registry")
			}
			if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
				t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
			}
			if got := r.URL.Query().Get("page"); got != "2" {
				t.Fatalf("page = %q, want %q", got, "2")
			}
			_, _ = w.Write([]byte(`{"success":true,"data":{"registries":[{"uid":"reg1","name":"Docker Hub"}],"total":1,"page":2,"page_size":10}}`))
		case 3:
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/api/v1/external-registry/reg1" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/api/v1/external-registry/reg1")
			}
			_, _ = w.Write([]byte(`{"success":true,"data":{"uid":"reg1","name":"Docker Hub"}}`))
		case 4:
			if r.Method != http.MethodDelete {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodDelete)
			}
			if r.URL.Path != "/api/v1/external-registry/reg1" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/api/v1/external-registry/reg1")
			}
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected call %d", step)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	createResp, _, err := client.ExternalRegistries.Create(context.Background(), "w1", &CreateExternalRegistryRequest{
		Name:     "Docker Hub",
		Type:     "dockerhub",
		Username: "user",
		Password: "pass",
	})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if createResp.Data.UID != "reg1" {
		t.Fatalf("uid = %q, want %q", createResp.Data.UID, "reg1")
	}

	listResp, _, err := client.ExternalRegistries.List(context.Background(), "w1", &ExternalRegistryListOptions{Page: 2, PageSize: 10})
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(listResp.Data.Registries) != 1 {
		t.Fatalf("len(registries) = %d, want 1", len(listResp.Data.Registries))
	}

	getResp, _, err := client.ExternalRegistries.Get(context.Background(), "reg1")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if getResp.Data.Name != "Docker Hub" {
		t.Fatalf("name = %q, want %q", getResp.Data.Name, "Docker Hub")
	}

	if _, err := client.ExternalRegistries.Delete(context.Background(), "reg1"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}

func TestExternalRegistryService_BrowseDockerHub(t *testing.T) {
	t.Parallel()

	step := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		step++
		w.Header().Set("Content-Type", "application/json")
		switch step {
		case 1:
			if r.URL.Path != "/api/v1/external-registry/reg1/dockerhub/images" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/api/v1/external-registry/reg1/dockerhub/images")
			}
			if got := r.URL.Query().Get("page_size"); got != "20" {
				t.Fatalf("page_size = %q, want %q", got, "20")
			}
			_, _ = w.Write([]byte(`{"success":true,"data":{"repositories":[{"name":"nginx"}]}}`))
		case 2:
			if r.URL.Path != "/api/v1/external-registry/reg1/dockerhub/library/nginx/tags" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/api/v1/external-registry/reg1/dockerhub/library/nginx/tags")
			}
			_, _ = w.Write([]byte(`{"success":true,"data":{"tags":[{"name":"latest"}]}}`))
		case 3:
			if r.URL.Path != "/api/v1/external-registry/dockerhub/search" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/api/v1/external-registry/dockerhub/search")
			}
			if got := r.URL.Query().Get("q"); got != "nginx" {
				t.Fatalf("q = %q, want %q", got, "nginx")
			}
			_, _ = w.Write([]byte(`{"success":true,"data":{"results":[{"name":"nginx"}]}}`))
		case 4:
			if r.URL.Path != "/api/v1/external-registry/dockerhub/library/nginx/tags" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/api/v1/external-registry/dockerhub/library/nginx/tags")
			}
			_, _ = w.Write([]byte(`{"success":true,"data":{"tags":[{"name":"latest"}]}}`))
		default:
			t.Fatalf("unexpected call %d", step)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	repos, _, err := client.ExternalRegistries.ListDockerHubImages(context.Background(), "reg1", &DockerHubListOptions{PageSize: 20})
	if err != nil {
		t.Fatalf("ListDockerHubImages error: %v", err)
	}
	if len(repos.Data.Repositories) != 1 {
		t.Fatalf("len(repositories) = %d, want 1", len(repos.Data.Repositories))
	}

	tags, _, err := client.ExternalRegistries.ListDockerHubTags(context.Background(), "reg1", "library", "nginx", nil)
	if err != nil {
		t.Fatalf("ListDockerHubTags error: %v", err)
	}
	if len(tags.Data.Tags) != 1 {
		t.Fatalf("len(tags) = %d, want 1", len(tags.Data.Tags))
	}

	search, _, err := client.ExternalRegistries.SearchPublicDockerHubImages(context.Background(), &DockerHubSearchOptions{Query: "nginx"})
	if err != nil {
		t.Fatalf("SearchPublicDockerHubImages error: %v", err)
	}
	if len(search.Data.Results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(search.Data.Results))
	}

	publicTags, _, err := client.ExternalRegistries.ListPublicDockerHubTags(context.Background(), "library", "nginx", nil)
	if err != nil {
		t.Fatalf("ListPublicDockerHubTags error: %v", err)
	}
	if len(publicTags.Data.Tags) != 1 {
		t.Fatalf("len(tags) = %d, want 1", len(publicTags.Data.Tags))
	}
}
