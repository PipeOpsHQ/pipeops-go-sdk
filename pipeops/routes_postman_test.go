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
				if got := projects.Data.Projects[0].ID; got != "1487" {
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
				_, _, err := client.Projects.Create(ctx, &CreateProjectRequest{
					Name:          "test",
					ServerID:      "s1",
					EnvironmentID: "e1",
					Repository:    "https://example.com/repo.git",
					Branch:        "main",
				})
				return err
			},
		},
		{
			name:   "Get",
			method: http.MethodGet,
			path:   "/project/fetch/p1",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Projects.Get(ctx, "p1")
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
				case "Get":
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

func TestTeamService_UsesPostmanRoutes(t *testing.T) {
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
			path:   "/team/fetch",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Teams.List(ctx)
				return err
			},
		},
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

func TestEnvironmentService_UsesPostmanRoutes(t *testing.T) {
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
			path:   "/environment/fetch",
			run: func(ctx context.Context, client *Client) error {
				_, _, err := client.Environments.List(ctx)
				return err
			},
		},
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
			w.WriteHeader(http.StatusNoContent)
			return
		case r.Method == http.MethodPut && r.URL.Path == "/billing/workspace/cards/c1":
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
