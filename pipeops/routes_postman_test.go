package pipeops

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectService_UsesPostmanRoutes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		run    func(ctx context.Context, client *Client) error
		method string
		path   string
	}{
		{
			name:   "List",
			method: http.MethodGet,
			path:   "/project/fetch-names",
			run: func(ctx context.Context, client *Client) error {
				projects, _, err := client.Projects.List(ctx, nil)
				if err != nil {
					return err
				}
				if len(projects.Data.Projects) != 1 {
					return fmt.Errorf("projects len = %d, want %d", len(projects.Data.Projects), 1)
				}
				if got := projects.Data.Projects[0].ID.String(); got != "1487" {
					return fmt.Errorf("project id = %q, want %q", got, "1487")
				}
				return nil
			},
		},
		{
			name:   "Create",
			method: http.MethodPost,
			path:   "/project/create",
			run: func(ctx context.Context, client *Client) error {
				worker := false
				_, _, err := client.Projects.Create(ctx, &CreateProjectRequest{
					Name:               "test",
					Username:           "acme",
					Source:             "github",
					Repository:         "https://example.com/repo.git",
					Branch:             "main",
					CommitURL:          "https://example.com/repo/commit/sha",
					CommitSha:          "sha",
					RepositoryLanguage: "nodejs",
					Environment:        "development",
					EnvironmentUUID:    "e1",
					ClusterUUID:        "s1",
					WorkspaceUUID:      "w1",
					BuildSettings: CreateProjectBuildSettings{
						BuildMethod:  "nodejs",
						BuildCommand: "npm run build",
						RunCommand:   "npm start",
						Worker:       &worker,
					},
					NetworkSettings: []CreateProjectNetworkSetting{
						{Port: 3000, Protocol: "HTTP"},
					},
					EnvVariables: []CreateProjectEnvVar{},
				})
				return err
			},
		},
		{
			name:   "Delete",
			method: http.MethodDelete,
			path:   "/project/delete/p1",
			run: func(ctx context.Context, client *Client) error {
				_, err := client.Projects.Delete(ctx, "p1")
				return err
			},
		},
		{
			name:   "BulkDelete",
			method: http.MethodDelete,
			path:   "/project/delete/bulk",
			run: func(ctx context.Context, client *Client) error {
				_, err := client.Projects.BulkDelete(ctx, &BulkDeleteRequest{ProjectUUIDs: []string{"p1", "p2"}})
				return err
			},
		},
		{
			name:   "CPUMetrics",
			method: http.MethodGet,
			path:   "/observability/app/cpu",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Projects.GetCPUMetrics(ctx, &MetricsRequest{App: "project", WorkspaceUUID: "w1"})
				return err
			},
		},
		{
			name:   "GetLogs",
			method: http.MethodGet,
			path:   "/project/logs/p1",
			run: func(ctx context.Context, client *Client) error {
				logs, _, err := client.Projects.GetLogs(ctx, "p1", &LogsOptions{
					WorkspaceUUID: "w1",
					StartTime:     "s1",
					EndTime:       "e1",
					Limit:         10,
				})
				if err != nil {
					return err
				}
				if len(logs.Data.Logs) != 1 {
					return fmt.Errorf("logs len = %d, want %d", len(logs.Data.Logs), 1)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Fatalf("method = %s, want %s", r.Method, tt.method)
				}
				if r.URL.Path != tt.path {
					t.Fatalf("path = %s, want %s", r.URL.Path, tt.path)
				}

				w.Header().Set("Content-Type", "application/json")
				switch tt.name {
				case "List":
					if _, writeErr := w.Write([]byte(`{"data":{"projects":[{"UUID":"p1","Name":"proj","ID":1487}]},"message":"ok","success":true}`)); writeErr != nil {
						t.Fatalf("write response error: %v", writeErr)
					}
				case "Create":
					if _, writeErr := w.Write([]byte(`{"status":"success","message":"ok","data":{"project":{"uuid":"p1"}}}`)); writeErr != nil {
						t.Fatalf("write response error: %v", writeErr)
					}
				case "CPUMetrics":
					if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
						t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
					}
					if got := r.URL.Query().Get("app"); got != "project" {
						t.Fatalf("app = %q, want %q", got, "project")
					}
					if _, writeErr := w.Write([]byte(`{"status":"success","message":"ok","data":{"metrics":{}}}`)); writeErr != nil {
						t.Fatalf("write response error: %v", writeErr)
					}
				case "GetLogs":
					if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
						t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
					}
					if got := r.URL.Query().Get("app"); got != "project" {
						t.Fatalf("app = %q, want %q", got, "project")
					}
					if got := r.URL.Query().Get("start"); got != "s1" {
						t.Fatalf("start = %q, want %q", got, "s1")
					}
					if got := r.URL.Query().Get("end"); got != "e1" {
						t.Fatalf("end = %q, want %q", got, "e1")
					}
					if got := r.URL.Query().Get("limit"); got != "10" {
						t.Fatalf("limit = %q, want %q", got, "10")
					}
					if _, writeErr := w.Write([]byte(`{"success":true,"message":"ok","data":[{"values":[]} ]}`)); writeErr != nil {
						t.Fatalf("write response error: %v", writeErr)
					}
				default:
					w.WriteHeader(http.StatusNoContent)
				}
			}))
			t.Cleanup(server.Close)

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("NewClient error: %v", err)
			}

			if err := tt.run(context.Background(), client); err != nil {
				t.Fatalf("call error: %v", err)
			}
		})
	}
}

func TestProjectService_List_UsesProjectFetchEndpoint_WhenWorkspaceUUIDProvided(t *testing.T) {
	t.Parallel()

	const wantWorkspaceUUID = "w1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/project/fetch" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/project/fetch")
		}
		if got := r.URL.Query().Get("workspace_uuid"); got != wantWorkspaceUUID {
			t.Fatalf("workspace_uuid = %q, want %q", got, wantWorkspaceUUID)
		}
		if got := r.URL.Query().Get("limit"); got != "30" {
			t.Fatalf("limit = %q, want %q", got, "30")
		}

		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write([]byte(`{"data":{"projects":[{"UUID":"p1","Name":"proj","ID":1487}]},"message":"ok","success":true}`)); writeErr != nil {
			t.Fatalf("write response error: %v", writeErr)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	projects, _, err := client.Projects.List(context.Background(), &ProjectListOptions{WorkspaceUUID: wantWorkspaceUUID, Limit: 30})
	if err != nil {
		t.Fatalf("Projects.List error: %v", err)
	}
	if len(projects.Data.Projects) != 1 {
		t.Fatalf("projects len = %d, want %d", len(projects.Data.Projects), 1)
	}
	if got := projects.Data.Projects[0].ID.String(); got != "1487" {
		t.Fatalf("project id = %q, want %q", got, "1487")
	}
}

func TestProjectService_Get_FallsBackToWorkspaceUUIDQuery_OnNotFound(t *testing.T) {
	t.Parallel()

	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		switch calls {
		case 1:
			// First call is to fetch workspace
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/workspace" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/workspace")
			}
			w.Header().Set("Content-Type", "application/json")
			if _, writeErr := w.Write([]byte(`{"data":[{"UUID":"w1"}],"message":"ok","success":true}`)); writeErr != nil {
				t.Fatalf("write response error: %v", writeErr)
			}
			return
		case 2:
			// Second call is project fetch with workspace_uuid
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/project/fetch/p1" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/project/fetch/p1")
			}
			if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
				t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
			}
			w.Header().Set("Content-Type", "application/json")
			if _, writeErr := w.Write([]byte(`{"data":{"project":{"uuid":"p1"}},"message":"ok","success":true}`)); writeErr != nil {
				t.Fatalf("write response error: %v", writeErr)
			}
			return
		default:
			t.Fatalf("unexpected call %d", calls)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, _, err = client.Projects.Get(context.Background(), "p1")
	if err != nil {
		t.Fatalf("Projects.Get error: %v", err)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want %d", calls, 2)
	}
}

func TestProjectService_TailLogs_DefaultsToTailMode(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/project/logs/p1" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/project/logs/p1")
		}
		if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
			t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
		}
		if got := r.URL.Query().Get("app"); got != "project" {
			t.Fatalf("app = %q, want %q", got, "project")
		}
		if got := r.URL.Query().Get("log"); got != "tail" {
			t.Fatalf("log = %q, want %q", got, "tail")
		}

		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write([]byte(`{"success":true,"message":"ok","data":[{"values":[]} ]}`)); writeErr != nil {
			t.Fatalf("write response error: %v", writeErr)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	logs, _, err := client.Projects.TailLogs(context.Background(), "p1", &LogsOptions{WorkspaceUUID: "w1"})
	if err != nil {
		t.Fatalf("Projects.TailLogs error: %v", err)
	}
	if len(logs.Data.Logs) != 1 {
		t.Fatalf("logs len = %d, want %d", len(logs.Data.Logs), 1)
	}
}

func TestProjectService_TailLogs_ResolvesWorkspaceUUID_WhenOptionsNil(t *testing.T) {
	t.Parallel()

	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		switch calls {
		case 1:
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
		case 2:
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/project/logs/p1" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/project/logs/p1")
			}
			if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
				t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
			}
			if got := r.URL.Query().Get("app"); got != "project" {
				t.Fatalf("app = %q, want %q", got, "project")
			}
			if got := r.URL.Query().Get("log"); got != "tail" {
				t.Fatalf("log = %q, want %q", got, "tail")
			}

			w.Header().Set("Content-Type", "application/json")
			if _, writeErr := w.Write([]byte(`{"success":true,"message":"ok","data":[{"values":[]} ]}`)); writeErr != nil {
				t.Fatalf("write response error: %v", writeErr)
			}
		default:
			t.Fatalf("unexpected call %d (%s %s?%s)", calls, r.Method, r.URL.Path, r.URL.RawQuery)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	logs, _, err := client.Projects.TailLogs(context.Background(), "p1", nil)
	if err != nil {
		t.Fatalf("Projects.TailLogs error: %v", err)
	}
	if len(logs.Data.Logs) != 1 {
		t.Fatalf("logs len = %d, want %d", len(logs.Data.Logs), 1)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want %d", calls, 2)
	}
}

func TestProjectService_SettingsRoutes_IncludeWorkspaceUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		run       func(ctx context.Context, client *Client) error
		method    string
		path      string
		checkBody func(t *testing.T, body []byte)
		response  func(t *testing.T, w http.ResponseWriter)
	}{
		{
			name:   "UpdateEnvVariables",
			method: http.MethodPost,
			path:   "/project/settings/env/p1",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Projects.UpdateEnvVariables(ctx, "p1", &EnvVariablesRequest{})
				return err
			},
			response: func(t *testing.T, w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"status":"success","message":"ok","data":{"env_variables":[]}}`)); err != nil {
					t.Fatalf("write response error: %v", err)
				}
			},
		},
		{
			name:   "GetEnvVariables",
			method: http.MethodGet,
			path:   "/project/settings/env/p1",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Projects.GetEnvVariables(ctx, "p1")
				return err
			},
			response: func(t *testing.T, w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"status":"success","message":"ok","data":{"env_variables":[]}}`)); err != nil {
					t.Fatalf("write response error: %v", err)
				}
			},
		},
		{
			name:   "CreateNetworkPolicy",
			method: http.MethodPost,
			path:   "/project/settings/p1/network-policy",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Projects.CreateNetworkPolicy(ctx, "p1", &NetworkPolicyRequest{Name: "np"})
				return err
			},
			response: func(t *testing.T, w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"status":"success","message":"ok","data":{"policy":{}}}`)); err != nil {
					t.Fatalf("write response error: %v", err)
				}
			},
		},
		{
			name:   "UpdateNetworkPolicy",
			method: http.MethodPut,
			path:   "/project/settings/p1/network-policy/np1",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Projects.UpdateNetworkPolicy(ctx, "p1", "np1", &NetworkPolicyRequest{Name: "np"})
				return err
			},
			response: func(t *testing.T, w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"status":"success","message":"ok","data":{"policy":{}}}`)); err != nil {
					t.Fatalf("write response error: %v", err)
				}
			},
		},
		{
			name:   "ListNetworkPolicies",
			method: http.MethodGet,
			path:   "/project/settings/p1/network-policy",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Projects.ListNetworkPolicies(ctx, "p1")
				return err
			},
			response: func(t *testing.T, w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"status":"success","message":"ok","data":{"policies":[]}}`)); err != nil {
					t.Fatalf("write response error: %v", err)
				}
			},
		},
		{
			name:   "UpdateNetworkingPort",
			method: http.MethodPut,
			path:   "/project/settings/network/p1",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Projects.UpdateNetworkingPort(ctx, "p1", &NetworkSettingsRequest{Port: 8080})
				return err
			},
			response: func(t *testing.T, w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"status":"success","message":"ok","data":{"settings":{}}}`)); err != nil {
					t.Fatalf("write response error: %v", err)
				}
			},
		},
		{
			name:   "GetNetworkSettings",
			method: http.MethodGet,
			path:   "/project/settings/network/p1",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Projects.GetNetworkSettings(ctx, "p1")
				return err
			},
			response: func(t *testing.T, w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"status":"success","message":"ok","data":{"settings":{}}}`)); err != nil {
					t.Fatalf("write response error: %v", err)
				}
			},
		},
		{
			name:   "GenerateDomainFromNetworkPort",
			method: http.MethodPost,
			path:   "/project/settings/network-name/p1",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Projects.GenerateDomainFromNetworkPort(ctx, "p1")
				return err
			},
			response: func(t *testing.T, w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"status":"success","message":"ok","data":{"domain":"example.com"}}`)); err != nil {
					t.Fatalf("write response error: %v", err)
				}
			},
		},
		{
			name:   "UpdateDomain",
			method: http.MethodPost,
			path:   "/project/settings/name/p1",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Projects.UpdateDomain(ctx, "p1", &DomainRequest{Domain: "example.com"})
				return err
			},
			checkBody: func(t *testing.T, body []byte) {
				var payload map[string]any
				if err := json.Unmarshal(body, &payload); err != nil {
					t.Fatalf("unmarshal body error: %v", err)
				}
				if got := payload["customDomainName"]; got != "example.com" {
					t.Fatalf("customDomainName = %#v, want %q", got, "example.com")
				}
			},
			response: func(t *testing.T, w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"success":true,"message":"ok","data":{}}`)); err != nil {
					t.Fatalf("write response error: %v", err)
				}
			},
		},
		{
			name:   "DeleteCustomDomain",
			method: http.MethodPatch,
			path:   "/project/p1/custom-domain",
			run: func(ctx context.Context, client *Client) error {
				_, err := client.Projects.DeleteCustomDomain(ctx, "p1")
				return err
			},
			response: func(t *testing.T, w http.ResponseWriter) {
				w.WriteHeader(http.StatusNoContent)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				switch calls {
				case 1:
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
				case 2:
					if r.Method != tt.method {
						t.Fatalf("method = %s, want %s", r.Method, tt.method)
					}
					if r.URL.Path != tt.path {
						t.Fatalf("path = %s, want %s", r.URL.Path, tt.path)
					}
					if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
						t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
					}

					if tt.checkBody != nil {
						bodyBytes, err := io.ReadAll(r.Body)
						if err != nil {
							t.Fatalf("read body error: %v", err)
						}
						tt.checkBody(t, bodyBytes)
					}

					if tt.response == nil {
						t.Fatalf("test response handler missing")
					}
					tt.response(t, w)
				default:
					t.Fatalf("unexpected call %d (%s %s?%s)", calls, r.Method, r.URL.Path, r.URL.RawQuery)
				}
			}))
			t.Cleanup(server.Close)

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("NewClient error: %v", err)
			}

			if err := tt.run(context.Background(), client); err != nil {
				t.Fatalf("call error: %v", err)
			}
			if calls != 2 {
				t.Fatalf("calls = %d, want %d", calls, 2)
			}
		})
	}
}

func TestProjectService_List_FallsBackToLegacyProjectsEndpoint_OnNotFound(t *testing.T) {
	t.Parallel()

	var calledFetchNames bool
	var calledLegacy bool
	var calledWorkspace bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/project/fetch-names":
			calledFetchNames = true
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte(`{"message":"not found"}`)); err != nil {
				t.Fatalf("write response error: %v", err)
			}
			return
		case r.Method == http.MethodGet && r.URL.Path == "/projects":
			calledLegacy = true
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte(`{"data":{"projects":[{"UUID":"p1","Name":"proj","ID":1487}]},"message":"ok","success":true}`)); err != nil {
				t.Fatalf("write response error: %v", err)
			}
			return
		case r.Method == http.MethodGet && (r.URL.Path == "/workspace" || r.URL.Path == "/workspace/fetch/w1"):
			calledWorkspace = true
			t.Fatalf("unexpected workspace fallback request: %s %s", r.Method, r.URL.Path)
			return
		default:
			t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	projects, _, err := client.Projects.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Projects.List error: %v", err)
	}
	if !calledFetchNames || !calledLegacy || calledWorkspace {
		t.Fatalf("fallback calls: fetchNames=%v legacy=%v workspace=%v", calledFetchNames, calledLegacy, calledWorkspace)
	}
	if len(projects.Data.Projects) != 1 {
		t.Fatalf("projects len = %d, want %d", len(projects.Data.Projects), 1)
	}
	if got := projects.Data.Projects[0].ID.String(); got != "1487" {
		t.Fatalf("project id = %q, want %q", got, "1487")
	}
}

func TestProjectService_List_FallsBackToWorkspaceRoutes_OnNotFound(t *testing.T) {
	t.Parallel()

	var calledFetchNames bool
	var calledLegacy bool
	var calledWorkspaceList bool
	var calledWorkspaceFetch bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/project/fetch-names":
			calledFetchNames = true
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte(`{"message":"not found"}`)); err != nil {
				t.Fatalf("write response error: %v", err)
			}
			return
		case r.Method == http.MethodGet && r.URL.Path == "/projects":
			calledLegacy = true
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte(`{"message":"not found"}`)); err != nil {
				t.Fatalf("write response error: %v", err)
			}
			return
		case r.Method == http.MethodGet && r.URL.Path == "/workspace":
			calledWorkspaceList = true
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte(`{"data":[{"UUID":"w1"}],"message":"ok","success":true}`)); err != nil {
				t.Fatalf("write response error: %v", err)
			}
			return
		case r.Method == http.MethodGet && r.URL.Path == "/workspace/fetch/w1":
			calledWorkspaceFetch = true
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte(`{"data":{"workspace":{"UUID":"w1","Projects":[{"UUID":"p1","Name":"proj","ID":1487}]}},"message":"ok","success":true}`)); err != nil {
				t.Fatalf("write response error: %v", err)
			}
			return
		default:
			t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	projects, _, err := client.Projects.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Projects.List error: %v", err)
	}
	if !calledFetchNames || !calledLegacy || !calledWorkspaceList || !calledWorkspaceFetch {
		t.Fatalf("fallback calls: fetchNames=%v legacy=%v workspaceList=%v workspaceFetch=%v", calledFetchNames, calledLegacy, calledWorkspaceList, calledWorkspaceFetch)
	}
	if len(projects.Data.Projects) != 1 {
		t.Fatalf("projects len = %d, want %d", len(projects.Data.Projects), 1)
	}
	if got := projects.Data.Projects[0].ID.String(); got != "1487" {
		t.Fatalf("project id = %q, want %q", got, "1487")
	}
}

func TestTeamService_UsesPostmanRoutes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		run    func(ctx context.Context, client *Client) error
		method string
		path   string
	}{
		{
			name:   "Get",
			method: http.MethodGet,
			path:   "/team/fetch/t1",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Teams.Get(ctx, "t1")
				return err
			},
		},
		{
			name:   "Delete",
			method: http.MethodDelete,
			path:   "/team/t1/delete",
			run: func(ctx context.Context, client *Client) error {
				_, err := client.Teams.Delete(ctx, "t1")
				return err
			},
		},
		{
			name:   "AcceptInvite",
			method: http.MethodPost,
			path:   "/team/accept-invite",
			run: func(ctx context.Context, client *Client) error {
				_, err := client.Teams.AcceptInvitation(ctx, "invite-123")
				return err
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Fatalf("method = %s, want %s", r.Method, tt.method)
				}
				if r.URL.Path != tt.path {
					t.Fatalf("path = %s, want %s", r.URL.Path, tt.path)
				}

				if tt.name == "AcceptInvite" {
					bodyBytes, err := io.ReadAll(r.Body)
					if err != nil {
						t.Fatalf("read body error: %v", err)
					}
					var payload map[string]any
					if err := json.Unmarshal(bodyBytes, &payload); err != nil {
						t.Fatalf("unmarshal body error: %v", err)
					}
					if payload["invite_id"] != "invite-123" {
						t.Fatalf("invite_id = %#v, want %q", payload["invite_id"], "invite-123")
					}
				}

				w.Header().Set("Content-Type", "application/json")
				if _, writeErr := w.Write([]byte(`{"status":"success","message":"ok","data":{}}`)); writeErr != nil {
					t.Fatalf("write response error: %v", writeErr)
				}
			}))
			t.Cleanup(server.Close)

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("NewClient error: %v", err)
			}

			if err := tt.run(context.Background(), client); err != nil {
				t.Fatalf("call error: %v", err)
			}
		})
	}
}

func TestTeamService_List_UsesWorkspaceScopedRoute(t *testing.T) {
	t.Parallel()

	const wantWorkspaceUUID = "w1"

	var workspaceCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/workspace":
			workspaceCalls++
			w.Header().Set("Content-Type", "application/json")
			if _, writeErr := w.Write([]byte(`{"data":[{"UUID":"w1"}],"message":"ok","success":true}`)); writeErr != nil {
				t.Fatalf("write response error: %v", writeErr)
			}
			return
		case r.Method == http.MethodGet && r.URL.Path == "/team/fetch":
			if got := r.URL.Query().Get("workspace_uuid"); got != wantWorkspaceUUID {
				t.Fatalf("workspace_uuid = %q, want %q", got, wantWorkspaceUUID)
			}
			w.Header().Set("Content-Type", "application/json")
			if _, writeErr := w.Write([]byte(`{"status":"success","message":"ok","data":{"teams":[]}}`)); writeErr != nil {
				t.Fatalf("write response error: %v", writeErr)
			}
			return
		default:
			t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, _, err = client.Teams.List(context.Background())
	if err != nil {
		t.Fatalf("Teams.List error: %v", err)
	}
	if workspaceCalls == 0 {
		t.Fatalf("expected workspace UUID lookup before listing teams")
	}
}

func TestEnvironmentService_UsesPostmanRoutes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		run    func(ctx context.Context, client *Client) error
		method string
		path   string
	}{
		{
			name:   "Get",
			method: http.MethodGet,
			path:   "/environment/fetch/e1",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Environments.Get(ctx, "e1")
				return err
			},
		},
		{
			name:   "Create",
			method: http.MethodPost,
			path:   "/environment/create",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Environments.Create(ctx, &CreateEnvironmentRequest{
					Name:          "env",
					WorkspaceUUID: "w1",
					ClusterUUID:   "c1",
				})
				return err
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Fatalf("method = %s, want %s", r.Method, tt.method)
				}
				if r.URL.Path != tt.path {
					t.Fatalf("path = %s, want %s", r.URL.Path, tt.path)
				}

				if tt.name == "Create" {
					if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
						t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
					}
				}

				w.Header().Set("Content-Type", "application/json")
				switch tt.name {
				case "List":
					if _, writeErr := w.Write([]byte(`{"status":"success","message":"ok","data":{"environments":[]}}`)); writeErr != nil {
						t.Fatalf("write response error: %v", writeErr)
					}
				case "Get", "Create":
					if _, writeErr := w.Write([]byte(`{"status":"success","message":"ok","data":{"environment":{"uuid":"e1"}}}`)); writeErr != nil {
						t.Fatalf("write response error: %v", writeErr)
					}
				default:
					if _, writeErr := w.Write([]byte(`{"status":"success","message":"ok"}`)); writeErr != nil {
						t.Fatalf("write response error: %v", writeErr)
					}
				}
			}))
			t.Cleanup(server.Close)

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("NewClient error: %v", err)
			}

			if err := tt.run(context.Background(), client); err != nil {
				t.Fatalf("call error: %v", err)
			}
		})
	}
}

func TestEnvironmentService_List_UsesWorkspaceScopedRoute(t *testing.T) {
	t.Parallel()

	const wantWorkspaceUUID = "w1"

	var workspaceCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/workspace":
			workspaceCalls++
			w.Header().Set("Content-Type", "application/json")
			if _, writeErr := w.Write([]byte(`{"data":[{"UUID":"w1"}],"message":"ok","success":true}`)); writeErr != nil {
				t.Fatalf("write response error: %v", writeErr)
			}
			return
		case r.Method == http.MethodGet && r.URL.Path == "/environment/fetch":
			if got := r.URL.Query().Get("workspace_uuid"); got != wantWorkspaceUUID {
				t.Fatalf("workspace_uuid = %q, want %q", got, wantWorkspaceUUID)
			}
			w.Header().Set("Content-Type", "application/json")
			if _, writeErr := w.Write([]byte(`{"status":"success","message":"ok","data":{"environments":[]}}`)); writeErr != nil {
				t.Fatalf("write response error: %v", writeErr)
			}
			return
		default:
			t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, _, err = client.Environments.List(context.Background())
	if err != nil {
		t.Fatalf("Environments.List error: %v", err)
	}
	if workspaceCalls == 0 {
		t.Fatalf("expected workspace UUID lookup before listing environments")
	}
}

func TestWebhookService_UsesPostmanRoutes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		run    func(ctx context.Context, client *Client) error
		method string
		path   string
		query  string
	}{
		{
			name:   "Create",
			method: http.MethodPost,
			path:   "/customer-webhook/create",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Webhooks.Create(ctx, &CreateWebhookRequest{
					URL:         "https://example.com/hook",
					Events:      []string{"deployment"},
					Description: "test",
				})
				return err
			},
		},
		{
			name:   "List",
			method: http.MethodGet,
			path:   "/customer-webhook/fetch",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Webhooks.List(ctx)
				return err
			},
		},
		{
			name:   "UpdateEnable",
			method: http.MethodPut,
			path:   "/customer-webhook/wh1",
			query:  "action=enable",
			run: func(ctx context.Context, client *Client) error {
				enabled := true
				_, _, err := client.Webhooks.Update(ctx, "wh1", &UpdateWebhookRequest{Active: &enabled})
				return err
			},
		},
		{
			name:   "Delete",
			method: http.MethodDelete,
			path:   "/customer-webhook/wh1",
			run: func(ctx context.Context, client *Client) error {
				_, err := client.Webhooks.Delete(ctx, "wh1")
				return err
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Fatalf("method = %s, want %s", r.Method, tt.method)
				}
				if r.URL.Path != tt.path {
					t.Fatalf("path = %s, want %s", r.URL.Path, tt.path)
				}
				if tt.query != "" && r.URL.RawQuery != tt.query {
					t.Fatalf("query = %q, want %q", r.URL.RawQuery, tt.query)
				}

				w.Header().Set("Content-Type", "application/json")
				if _, writeErr := w.Write([]byte(`{"status":"success","message":"ok","data":{}}`)); writeErr != nil {
					t.Fatalf("write response error: %v", writeErr)
				}
			}))
			t.Cleanup(server.Close)

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("NewClient error: %v", err)
			}

			if err := tt.run(context.Background(), client); err != nil {
				t.Fatalf("call error: %v", err)
			}
		})
	}
}

func TestMiscAndBillingAndAddons_RouteFixes(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/workspace":
			w.Header().Set("Content-Type", "application/json")
			if _, writeErr := w.Write([]byte(`{"data":[{"UUID":"w1"}],"message":"ok","success":true}`)); writeErr != nil {
				t.Fatalf("write response error: %v", writeErr)
			}
			return
		case r.Method == http.MethodGet && r.URL.Path == "/partners/participants/verify":
			if got := r.URL.Query().Get("verification_code"); got != "abc+123" {
				t.Fatalf("verification_code = %q, want %q", got, "abc+123")
			}
			w.Header().Set("Content-Type", "application/json")
			if _, writeErr := w.Write([]byte(`{"status":"success","message":"ok","data":{"valid":true}}`)); writeErr != nil {
				t.Fatalf("write response error: %v", writeErr)
			}
			return
		case r.Method == http.MethodPost && r.URL.Path == "/addons/domains/a1":
			if got := r.URL.Query().Get("workspace"); got != "w1" {
				t.Fatalf("workspace = %q, want %q", got, "w1")
			}
			w.WriteHeader(http.StatusNoContent)
			return
		case r.Method == http.MethodPut && r.URL.Path == "/billing/workspace/cards/c1":
			if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
				t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
			}
			w.Header().Set("Content-Type", "application/json")
			if _, writeErr := w.Write([]byte(`{"status":"success","message":"ok","data":{"card":{"uuid":"c1"}}}`)); writeErr != nil {
				t.Fatalf("write response error: %v", writeErr)
			}
			return
		default:
			t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, _, err = client.PartnerParticipants.VerifyProgramCode(context.Background(), "abc+123")
	if err != nil {
		t.Fatalf("VerifyProgramCode error: %v", err)
	}

	_, err = client.AddOns.AddDomain(context.Background(), "a1", &DomainRequest{Domain: "example.com"})
	if err != nil {
		t.Fatalf("AddDomain error: %v", err)
	}

	_, _, err = client.Billing.UpdateCard(context.Background(), "c1", &AddCardRequest{Token: "tok"})
	if err != nil {
		t.Fatalf("UpdateCard error: %v", err)
	}
}
