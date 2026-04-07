package pipeops

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectService_ListProviderOrganizations_UsesProviderRoute(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/project/github/organisations" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/project/github/organisations")
		}
		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write([]byte(`{"success":true,"message":"ok","data":[{"org_name":"acme"}]}`)); writeErr != nil {
			t.Fatalf("write response error: %v", writeErr)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Projects.ListProviderOrganizations(context.Background(), "github")
	if err != nil {
		t.Fatalf("ListProviderOrganizations error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("organizations len = %d, want %d", len(resp.Data), 1)
	}
}

func TestProjectService_ListProviderOrganizationRepos_UsesOrgNameAndPage(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/project/gitlab/organisations/repos" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/project/gitlab/organisations/repos")
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Fatalf("page = %q, want %q", got, "2")
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload error: %v", err)
		}
		if got := payload["org_name"]; got != "acme" {
			t.Fatalf("org_name = %q, want %q", got, "acme")
		}
		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write([]byte(`{"success":true,"message":"ok","data":[{"repo_fullname":"acme/web"}],"meta_data":{"next_page":3}}`)); writeErr != nil {
			t.Fatalf("write response error: %v", writeErr)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Projects.ListProviderOrganizationRepos(context.Background(), "gitlab", &ProviderOrganizationReposRequest{OrgName: "acme"}, &ProviderCollectionOptions{Page: 2})
	if err != nil {
		t.Fatalf("ListProviderOrganizationRepos error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("repos len = %d, want %d", len(resp.Data), 1)
	}
}

func TestProjectService_ListProviderBranches_UsesRepoFullnameVisibilityAndSearch(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/project/github/branches" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/project/github/branches")
		}
		if got := r.URL.Query().Get("search"); got != "release" {
			t.Fatalf("search = %q, want %q", got, "release")
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload error: %v", err)
		}
		if got := payload["repo_fullname"]; got != "acme/web" {
			t.Fatalf("repo_fullname = %q, want %q", got, "acme/web")
		}
		if got := payload["visibility"]; got != "private" {
			t.Fatalf("visibility = %q, want %q", got, "private")
		}
		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write([]byte(`{"success":true,"message":"ok","data":[{"name":"release/v1"}]}`)); writeErr != nil {
			t.Fatalf("write response error: %v", writeErr)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Projects.ListProviderBranches(context.Background(), "github", &ProviderBranchesRequest{RepoFullname: "acme/web", Visibility: "private"}, &ProviderBranchesOptions{Search: "release"})
	if err != nil {
		t.Fatalf("ListProviderBranches error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("branches len = %d, want %d", len(resp.Data), 1)
	}
}

func TestProjectService_GetGitHubBranches_MapsBranchNames(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/project/github/branches" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/project/github/branches")
		}
		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write([]byte(`{"success":true,"message":"ok","data":[{"name":"main"},{"name":"develop"}]}`)); writeErr != nil {
			t.Fatalf("write response error: %v", writeErr)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Projects.GetGitHubBranches(context.Background(), &GitHubBranchesRequest{RepoFullname: "acme/web"})
	if err != nil {
		t.Fatalf("GetGitHubBranches error: %v", err)
	}
	if len(resp.Data.Branches) != 2 {
		t.Fatalf("branches len = %d, want %d", len(resp.Data.Branches), 2)
	}
	if got := resp.Data.Branches[0]; got != "main" {
		t.Fatalf("first branch = %q, want %q", got, "main")
	}
}

func TestProjectService_SearchProviderRepositories_UsesRepositoryNameOrgNameAndPage(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/project/bitbucket/repo-search" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/project/bitbucket/repo-search")
		}
		if got := r.URL.Query().Get("page"); got != "3" {
			t.Fatalf("page = %q, want %q", got, "3")
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload error: %v", err)
		}
		if got := payload["org_name"]; got != "acme" {
			t.Fatalf("org_name = %q, want %q", got, "acme")
		}
		if got := payload["repository_name"]; got != "web" {
			t.Fatalf("repository_name = %q, want %q", got, "web")
		}
		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write([]byte(`{"success":true,"message":"ok","data":[{"repo_fullname":"acme/web"}]}`)); writeErr != nil {
			t.Fatalf("write response error: %v", writeErr)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Projects.SearchProviderRepositories(context.Background(), "bitbucket", &ProviderRepoSearchRequest{OrgName: "acme", RepositoryName: "web"}, &ProviderCollectionOptions{Page: 3})
	if err != nil {
		t.Fatalf("SearchProviderRepositories error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("repos len = %d, want %d", len(resp.Data), 1)
	}
}

func TestProjectService_LinkProviderWithRedirect_UsesPostBody(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/project/link/github" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/project/link/github")
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload error: %v", err)
		}
		if got := payload["redirectPath"]; got != "/apps/new" {
			t.Fatalf("redirectPath = %q, want %q", got, "/apps/new")
		}
		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write([]byte(`{"redirectUrl":"https://github.com/login/oauth/authorize","provider":"github"}`)); writeErr != nil {
			t.Fatalf("write response error: %v", writeErr)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Projects.LinkProviderWithRedirect(context.Background(), "github", &LinkProviderRequest{RedirectPath: "/apps/new"})
	if err != nil {
		t.Fatalf("LinkProviderWithRedirect error: %v", err)
	}
	if got := resp.Provider; got != "github" {
		t.Fatalf("provider = %q, want %q", got, "github")
	}
}

func TestProjectService_CheckRepositoryDockerfile_EscapesPathSegments(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if got := r.URL.EscapedPath(); got != "/project/check-dockerfile/github/acme/web/feature%2Fmain" {
			t.Fatalf("escaped path = %s, want %s", got, "/project/check-dockerfile/github/acme/web/feature%2Fmain")
		}
		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write([]byte(`{"status":"success","message":"ok","data":{"exists":true}}`)); writeErr != nil {
			t.Fatalf("write response error: %v", writeErr)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Projects.CheckRepositoryDockerfile(context.Background(), "github", "acme", "web", "feature/main")
	if err != nil {
		t.Fatalf("CheckRepositoryDockerfile error: %v", err)
	}
	if !resp.Data.Exists {
		t.Fatal("expected Dockerfile to exist")
	}
}

func TestAddOnService_List_UsesSearchAndFilterQuery(t *testing.T) {
	t.Parallel()

	featured := true
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/addons" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/addons")
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Fatalf("page = %q, want %q", got, "2")
		}
		if got := r.URL.Query().Get("limit"); got != "25" {
			t.Fatalf("limit = %q, want %q", got, "25")
		}
		if got := r.URL.Query().Get("category"); got != "databases" {
			t.Fatalf("category = %q, want %q", got, "databases")
		}
		if got := r.URL.Query().Get("s"); got != "redis" {
			t.Fatalf("s = %q, want %q", got, "redis")
		}
		if got := r.URL.Query().Get("featured"); got != "true" {
			t.Fatalf("featured = %q, want %q", got, "true")
		}
		if got := r.URL.Query().Get("workspace"); got != "w1" {
			t.Fatalf("workspace = %q, want %q", got, "w1")
		}
		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write([]byte(`{"success":true,"message":"ok","data":[{"id":"addon-1"}]}`)); writeErr != nil {
			t.Fatalf("write response error: %v", writeErr)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.AddOns.List(context.Background(), &ListAddOnsOptions{
		Page:          2,
		Limit:         25,
		Category:      "databases",
		Search:        "redis",
		Featured:      &featured,
		WorkspaceUUID: "w1",
	})
	if err != nil {
		t.Fatalf("AddOns.List error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("addons len = %d, want %d", len(resp.Data), 1)
	}
}

func TestAddOnService_ListCategories_DecodesArrayData(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/addons/categories" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/addons/categories")
		}
		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write([]byte(`{"success":true,"message":"ok","data":[{"id":"cat-1","name":"Databases"}]}`)); writeErr != nil {
			t.Fatalf("write response error: %v", writeErr)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.AddOns.ListCategories(context.Background())
	if err != nil {
		t.Fatalf("AddOns.ListCategories error: %v", err)
	}
	if len(resp.Data.Categories) != 1 {
		t.Fatalf("categories len = %d, want %d", len(resp.Data.Categories), 1)
	}
	if got := resp.Data.Categories[0].Name; got != "Databases" {
		t.Fatalf("category name = %q, want %q", got, "Databases")
	}
}
