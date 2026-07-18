package pipeops

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitOpsServicePaths(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := NewClient(server.URL, WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Create", func(t *testing.T) {
		mux.HandleFunc("/api/v1/gitops/applications", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				// List also hits this path with GET; only assert POST here via method branch.
				if r.Method == http.MethodGet {
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"success": true,
						"data": map[string]interface{}{
							"items": []map[string]string{{"uuid": "go-1", "name": "app"}},
							"total": 1, "page": 1, "limit": 20, "total_pages": 1,
						},
					})
					return
				}
				t.Fatalf("method = %s", r.Method)
			}
			body, _ := io.ReadAll(r.Body)
			var req map[string]interface{}
			if err := json.Unmarshal(body, &req); err != nil {
				t.Fatalf("body: %v", err)
			}
			if req["name"] != "my-app" || req["repo_url"] != "https://github.com/acme/app" {
				t.Fatalf("body = %s", body)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "created",
				"data":    map[string]string{"uuid": "go-1", "name": "my-app", "repo_url": "https://github.com/acme/app"},
			})
		})

		resp, _, err := client.GitOps.Create(context.Background(), &CreateGitOpsConfigRequest{
			Name:    "my-app",
			RepoURL: "https://github.com/acme/app",
			Branch:  "main",
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Data.UUID != "go-1" {
			t.Fatalf("uuid = %q", resp.Data.UUID)
		}
	})

	t.Run("List", func(t *testing.T) {
		// Handler registered in Create subtest covers GET on same path for sequential runs;
		// re-register a dedicated list assertion if needed via fresh path state.
		resp, _, err := client.GitOps.List(context.Background(), &GitOpsListOptions{Page: 1, Limit: 20})
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data.Items) != 1 || resp.Data.Items[0].UUID != "go-1" {
			t.Fatalf("items = %+v", resp.Data.Items)
		}
	})

	t.Run("Get", func(t *testing.T) {
		mux.HandleFunc("/api/v1/gitops/applications/go-2", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    map[string]string{"uuid": "go-2", "name": "other"},
			})
		})
		resp, _, err := client.GitOps.Get(context.Background(), "go-2")
		if err != nil {
			t.Fatal(err)
		}
		if resp.Data.UUID != "go-2" {
			t.Fatalf("uuid = %q", resp.Data.UUID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		mux.HandleFunc("/api/v1/gitops/applications/go-3", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Fatalf("method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    map[string]string{"uuid": "go-3", "branch": "develop"},
			})
		})
		resp, _, err := client.GitOps.Update(context.Background(), "go-3", &UpdateGitOpsConfigRequest{
			Branch: "develop",
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Data.Branch != "develop" {
			t.Fatalf("branch = %q", resp.Data.Branch)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		mux.HandleFunc("/api/v1/gitops/applications/go-4", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Fatalf("method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "deleted",
			})
		})
		if _, err := client.GitOps.Delete(context.Background(), "go-4"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Sync", func(t *testing.T) {
		mux.HandleFunc("/api/v1/gitops/applications/go-5/sync", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s", r.Method)
			}
			w.WriteHeader(http.StatusAccepted)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    map[string]interface{}{"status": "Syncing", "revision": "abc", "dry_run": false},
			})
		})
		resp, _, err := client.GitOps.TriggerSync(context.Background(), "go-5", &TriggerGitOpsSyncRequest{
			Revision: "abc",
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Data.Status != "Syncing" {
			t.Fatalf("status = %q", resp.Data.Status)
		}
	})

	t.Run("SyncStatus", func(t *testing.T) {
		mux.HandleFunc("/api/v1/gitops/applications/go-6/sync-status", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]string{
					"sync_status":   "Synced",
					"health_status": "Healthy",
				},
			})
		})
		resp, _, err := client.GitOps.GetSyncStatus(context.Background(), "go-6")
		if err != nil {
			t.Fatal(err)
		}
		if resp.Data.SyncStatus != "Synced" {
			t.Fatalf("sync_status = %q", resp.Data.SyncStatus)
		}
	})

	t.Run("Diff", func(t *testing.T) {
		mux.HandleFunc("/api/v1/gitops/applications/go-7/diff", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"current_commit": "aaa",
					"target_commit":  "bbb",
					"sync_required":  true,
				},
			})
		})
		resp, _, err := client.GitOps.GetDiff(context.Background(), "go-7")
		if err != nil {
			t.Fatal(err)
		}
		if !resp.Data.SyncRequired || resp.Data.CurrentCommit != "aaa" {
			t.Fatalf("diff = %+v", resp.Data)
		}
	})

	t.Run("History", func(t *testing.T) {
		mux.HandleFunc("/api/v1/gitops/applications/go-8/history", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s", r.Method)
			}
			if r.URL.Query().Get("page") != "2" {
				t.Fatalf("page = %q", r.URL.Query().Get("page"))
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"items": []map[string]interface{}{
						{"id": 1, "commit_sha": "deadbeef", "sync_status": "Synced"},
					},
					"total": 1, "page": 2, "limit": 10, "total_pages": 1,
				},
			})
		})
		resp, _, err := client.GitOps.GetHistory(context.Background(), "go-8", &GitOpsListOptions{Page: 2, Limit: 10})
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data.Items) != 1 || resp.Data.Items[0].CommitSHA != "deadbeef" {
			t.Fatalf("history = %+v", resp.Data.Items)
		}
	})
}
