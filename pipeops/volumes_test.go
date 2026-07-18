package pipeops

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVolumeServicePaths(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := NewClient(server.URL, WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}

	// Stub workspace list so firstWorkspaceUUID succeeds.
	mux.HandleFunc("/workspace", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    []map[string]string{{"UUID": "ws-1", "uuid": "ws-1"}},
		})
	})

	t.Run("List", func(t *testing.T) {
		mux.HandleFunc("/volumes", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s", r.Method)
			}
			if got := r.URL.Query().Get("workspace_uuid"); got != "ws-1" && r.URL.Query().Get("workspace") == "" {
				// List may use workspace_uuid from opts or first workspace
				if r.URL.Query().Get("workspace_uuid") == "" && r.URL.Query().Get("workspace") == "" {
					t.Fatalf("expected workspace query, got %s", r.URL.RawQuery)
				}
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"volumes": []map[string]string{{"uuid": "vol-1", "status": "mounted"}},
					"summary": map[string]int{"mounted": 1, "unattached": 0},
					"total":   1,
				},
			})
		})
		resp, _, err := client.Volumes.List(context.Background(), &VolumeListOptions{WorkspaceUUID: "ws-1"})
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data.Volumes) != 1 || resp.Data.Volumes[0].UUID != "vol-1" {
			t.Fatalf("volumes = %+v", resp.Data.Volumes)
		}
	})

	t.Run("Get", func(t *testing.T) {
		mux.HandleFunc("/volumes/vol-2", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    map[string]string{"uuid": "vol-2", "status": "unattached"},
			})
		})
		resp, _, err := client.Volumes.Get(context.Background(), "vol-2", &VolumeListOptions{WorkspaceUUID: "ws-1"})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Data.UUID != "vol-2" {
			t.Fatalf("uuid = %q", resp.Data.UUID)
		}
	})

	t.Run("Remount", func(t *testing.T) {
		mux.HandleFunc("/volumes/vol-3/remount", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"volume":  map[string]string{"uuid": "vol-3", "status": "mounted"},
					"message": "remount scheduled",
				},
			})
		})
		resp, _, err := client.Volumes.Remount(context.Background(), "vol-3", &RemountVolumeRequest{
			TargetType: "project",
			TargetUUID: "proj-1",
		}, &VolumeListOptions{WorkspaceUUID: "ws-1"})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Data.Volume.UUID != "vol-3" {
			t.Fatalf("volume = %+v", resp.Data.Volume)
		}
	})

	t.Run("Export", func(t *testing.T) {
		mux.HandleFunc("/volumes/vol-4/export", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost && r.Method != http.MethodGet {
				t.Fatalf("method = %s", r.Method)
			}
			status := "pending"
			if r.Method == http.MethodGet {
				status = "ready"
			}
			code := http.StatusOK
			if r.Method == http.MethodPost {
				code = http.StatusAccepted
			}
			w.WriteHeader(code)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    map[string]string{"uuid": "exp-1", "status": status},
			})
		})
		start, _, err := client.Volumes.StartExport(context.Background(), "vol-4", &VolumeListOptions{WorkspaceUUID: "ws-1"})
		if err != nil {
			t.Fatal(err)
		}
		if start.Data.Status != "pending" {
			t.Fatalf("start status = %q", start.Data.Status)
		}
		got, _, err := client.Volumes.GetExport(context.Background(), "vol-4", &VolumeListOptions{WorkspaceUUID: "ws-1"})
		if err != nil {
			t.Fatal(err)
		}
		if got.Data.Status != "ready" {
			t.Fatalf("get status = %q", got.Data.Status)
		}
	})
}

func TestBackupServiceRetired(t *testing.T) {
	client, err := NewClient("https://api.pipeops.test", WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Backups.CreateBackup(context.Background(), "p1"); err == nil {
		t.Fatal("expected deprecated error")
	}
	if _, _, err := client.Backups.ListBackups(context.Background(), "p1"); err == nil {
		t.Fatal("expected deprecated error")
	}
}

func TestAddOnBackupPaths(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := NewClient(server.URL, WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}
	mux.HandleFunc("/workspace", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    []map[string]string{{"UUID": "ws-1"}},
		})
	})
	mux.HandleFunc("/addons/deployments/dep-1/backups", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"addon_uid": "dep-1",
				"snapshots": []map[string]string{{"id": "snap-1"}},
			},
		})
	})
	mux.HandleFunc("/addons/deployments/dep-1/backups/export", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    map[string]string{"export_id": "exp-1", "status": "pending"},
		})
	})

	list, _, err := client.AddOns.ListAddonBackups(context.Background(), "dep-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Data.Snapshots) != 1 {
		t.Fatalf("snapshots = %+v", list.Data.Snapshots)
	}
	exp, _, err := client.AddOns.StartAddonBackupExport(context.Background(), "dep-1", &AddonBackupExportRequest{
		SnapshotID: "snap-1",
		Path:       "dump.sql",
	})
	if err != nil {
		t.Fatal(err)
	}
	if exp.Data.ExportID != "exp-1" {
		t.Fatalf("export = %+v", exp.Data)
	}
}
