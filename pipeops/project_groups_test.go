package pipeops

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectGroupServicePaths(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := NewClient(server.URL, WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}

	// Stub workspace list so firstWorkspaceUUID succeeds when opts omit workspace.
	mux.HandleFunc("/workspace", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    []map[string]string{{"UUID": "ws-1", "uuid": "ws-1"}},
		})
	})

	ws := &ProjectGroupWorkspaceOptions{WorkspaceUUID: "ws-1"}

	t.Run("List", func(t *testing.T) {
		mux.HandleFunc("/project-groups", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				body, _ := io.ReadAll(r.Body)
				var req map[string]interface{}
				_ = json.Unmarshal(body, &req)
				if req["name"] != "plane" {
					t.Fatalf("create body = %s", body)
				}
				if r.URL.Query().Get("workspace_uuid") != "ws-1" {
					t.Fatalf("workspace = %s", r.URL.RawQuery)
				}
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"success": true,
					"data":    map[string]string{"uuid": "pg-1", "name": "plane"},
				})
				return
			}
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s", r.Method)
			}
			if r.URL.Query().Get("workspace_uuid") != "ws-1" {
				t.Fatalf("workspace = %s", r.URL.RawQuery)
			}
			if r.URL.Query().Get("limit") != "25" {
				t.Fatalf("limit = %q", r.URL.Query().Get("limit"))
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"groups": []map[string]string{{"uuid": "pg-1", "name": "plane"}},
					"total":  1, "limit": 25, "offset": 0,
				},
			})
		})

		list, _, err := client.ProjectGroups.List(context.Background(), &ProjectGroupListOptions{
			WorkspaceUUID: "ws-1",
			Limit:         25,
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(list.Data.Groups) != 1 || list.Data.Groups[0].UUID != "pg-1" {
			t.Fatalf("groups = %+v", list.Data.Groups)
		}

		created, _, err := client.ProjectGroups.Create(context.Background(), &CreateProjectGroupRequest{Name: "plane"}, ws)
		if err != nil {
			t.Fatal(err)
		}
		if created.Data.UUID != "pg-1" {
			t.Fatalf("create uuid = %q", created.Data.UUID)
		}
	})

	t.Run("GetUpdateDelete", func(t *testing.T) {
		mux.HandleFunc("/project-groups/pg-2", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("workspace_uuid") != "ws-1" {
				t.Fatalf("workspace = %s", r.URL.RawQuery)
			}
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"success": true,
					"data":    map[string]string{"uuid": "pg-2", "name": "g2"},
				})
			case http.MethodPatch:
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"success": true,
					"data":    map[string]string{"uuid": "pg-2", "name": "renamed"},
				})
			case http.MethodDelete:
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"success": true,
					"message": "deleted",
				})
			default:
				t.Fatalf("method = %s", r.Method)
			}
		})

		got, _, err := client.ProjectGroups.Get(context.Background(), "pg-2", ws)
		if err != nil {
			t.Fatal(err)
		}
		if got.Data.UUID != "pg-2" {
			t.Fatalf("get = %+v", got.Data)
		}

		name := "renamed"
		upd, _, err := client.ProjectGroups.Update(context.Background(), "pg-2", &UpdateProjectGroupRequest{Name: &name}, ws)
		if err != nil {
			t.Fatal(err)
		}
		if upd.Data.Name != "renamed" {
			t.Fatalf("update = %+v", upd.Data)
		}

		if _, err := client.ProjectGroups.Delete(context.Background(), "pg-2", ws); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Members", func(t *testing.T) {
		mux.HandleFunc("/project-groups/pg-3/members", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"attached_member_uuids": []string{"proj-1"},
					"group_uuid":            "pg-3",
				},
			})
		})
		mux.HandleFunc("/project-groups/pg-3/members/project/proj-1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Fatalf("method = %s", r.Method)
			}
			if r.URL.Query().Get("workspace_uuid") != "ws-1" {
				t.Fatalf("workspace = %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "detached"})
		})

		att, _, err := client.ProjectGroups.AttachMember(context.Background(), "pg-3", &AttachProjectGroupMemberRequest{
			MemberType: "project",
			MemberUUID: "proj-1",
		}, ws)
		if err != nil {
			t.Fatal(err)
		}
		if len(att.Data.AttachedMemberUUIDs) != 1 {
			t.Fatalf("attach = %+v", att.Data)
		}

		if _, err := client.ProjectGroups.DetachMember(context.Background(), "pg-3", "project", "proj-1", &ProjectGroupDetachOptions{
			WorkspaceUUID: "ws-1",
		}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("TopologyEnvConnectRedeploy", func(t *testing.T) {
		mux.HandleFunc("/project-groups/pg-4/topology", func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"group": map[string]string{"uuid": "pg-4"},
					"nodes": []map[string]string{{"member_uuid": "proj-1", "name": "api"}},
				},
			})
		})
		mux.HandleFunc("/project-groups/pg-4/env", func(w http.ResponseWriter, r *http.Request) {
			vars := []map[string]string{{"key": "FOO", "value": "bar"}}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    map[string]interface{}{"variables": vars},
			})
		})
		mux.HandleFunc("/project-groups/pg-4/env/inject", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"written_keys":     []string{"FOO"},
					"projects_touched": []string{"proj-1"},
				},
			})
		})
		mux.HandleFunc("/project-groups/pg-4/connections", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"written_keys": []string{"DATABASE_URL"},
					"message":      "connected",
				},
			})
		})
		mux.HandleFunc("/project-groups/pg-4/redeploy-apps", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"queued":  []string{"proj-1"},
					"failed":  []string{},
					"message": "Queued redeploy for 1 app(s)",
				},
			})
		})

		topo, _, err := client.ProjectGroups.GetTopology(context.Background(), "pg-4", ws)
		if err != nil {
			t.Fatal(err)
		}
		if topo.Data.Group.UUID != "pg-4" || len(topo.Data.Nodes) != 1 {
			t.Fatalf("topology = %+v", topo.Data)
		}

		env, _, err := client.ProjectGroups.GetSharedEnv(context.Background(), "pg-4", ws)
		if err != nil {
			t.Fatal(err)
		}
		if len(env.Data.Variables) != 1 || env.Data.Variables[0].Key != "FOO" {
			t.Fatalf("env = %+v", env.Data)
		}

		put, _, err := client.ProjectGroups.PutSharedEnv(context.Background(), "pg-4", &UpsertProjectGroupSharedEnvRequest{
			Variables: []ProjectGroupSharedEnvVar{{Key: "FOO", Value: "bar"}},
		}, ws)
		if err != nil {
			t.Fatal(err)
		}
		if len(put.Data.Variables) != 1 {
			t.Fatalf("put env = %+v", put.Data)
		}

		inj, _, err := client.ProjectGroups.InjectSharedEnv(context.Background(), "pg-4", &InjectProjectGroupSharedEnvRequest{
			Overwrite: true,
		}, ws)
		if err != nil {
			t.Fatal(err)
		}
		if len(inj.Data.WrittenKeys) != 1 {
			t.Fatalf("inject = %+v", inj.Data)
		}

		conn, _, err := client.ProjectGroups.ConnectServices(context.Background(), "pg-4", &ConnectProjectGroupServicesRequest{
			ConsumerType: "project",
			ConsumerUUID: "proj-1",
			ProviderType: "addon_deployment",
			ProviderUUID: "addon-1",
		}, ws)
		if err != nil {
			t.Fatal(err)
		}
		if len(conn.Data.WrittenKeys) != 1 {
			t.Fatalf("connect = %+v", conn.Data)
		}

		rd, _, err := client.ProjectGroups.RedeployApps(context.Background(), "pg-4", ws)
		if err != nil {
			t.Fatal(err)
		}
		if len(rd.Data.Queued) != 1 {
			t.Fatalf("redeploy = %+v", rd.Data)
		}
	})

	t.Run("ResolveAndCandidates", func(t *testing.T) {
		mux.HandleFunc("/project-groups/resolve", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("member_type") != "project" || r.URL.Query().Get("member_uuid") != "proj-9" {
				t.Fatalf("query = %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]string{
					"group_uuid":  "pg-9",
					"member_type": "project",
					"member_uuid": "proj-9",
				},
			})
		})
		mux.HandleFunc("/project-groups/candidates", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("group_uuid") != "pg-9" {
				t.Fatalf("group_uuid = %q", r.URL.Query().Get("group_uuid"))
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"projects": []map[string]string{{"member_uuid": "proj-a", "name": "api"}},
					"addons":   []map[string]string{{"member_uuid": "addon-a", "name": "pg"}},
				},
			})
		})

		res, _, err := client.ProjectGroups.ResolveMember(context.Background(), &ProjectGroupResolveOptions{
			WorkspaceUUID: "ws-1",
			MemberType:    "project",
			MemberUUID:    "proj-9",
		})
		if err != nil {
			t.Fatal(err)
		}
		if res.Data.GroupUUID != "pg-9" {
			t.Fatalf("resolve = %+v", res.Data)
		}

		cands, _, err := client.ProjectGroups.ListCandidates(context.Background(), &ProjectGroupCandidatesOptions{
			WorkspaceUUID: "ws-1",
			GroupUUID:     "pg-9",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(cands.Data.Projects) != 1 || len(cands.Data.Addons) != 1 {
			t.Fatalf("candidates = %+v", cands.Data)
		}
	})
}
