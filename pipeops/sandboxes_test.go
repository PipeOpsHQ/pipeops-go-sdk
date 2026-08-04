package pipeops

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSandbox_Unmarshal_PascalAndSnake(t *testing.T) {
	t.Parallel()
	var s Sandbox
	if err := json.Unmarshal([]byte(`{"id":"c1","Name":"dev","Image":"ubuntu","Role":"standard","Status":"running"}`), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.Name != "dev" || s.Image != "ubuntu" || s.Role != "standard" || s.Status != "running" {
		t.Fatalf("sandbox = %+v", s)
	}
}

func TestSandboxSession_Unmarshal_TokenAlias(t *testing.T) {
	t.Parallel()
	var s SandboxSession
	if err := json.Unmarshal([]byte(`{"container_id":"c1","Token":"rexec_secret","base_url":"https://rexec.sh"}`), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.Token != "rexec_secret" || s.ContainerID != "c1" {
		t.Fatalf("session = %+v", s)
	}
}

func TestSandboxService_List_SendsWorkspaceQuery(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/api/v1/sandboxes" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("workspace_uuid"); got != "ws-1" {
			t.Fatalf("workspace_uuid = %q", got)
		}
		if got := r.URL.Query().Get("workspace"); got != "ws-1" {
			t.Fatalf("workspace = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"success":true,"message":"sandboxes retrieved","data":[{"id":"c1","name":"dev","status":"running"}],"meta":{"count":1}}`)); err != nil {
			t.Fatalf("write: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	out, _, err := client.Sandboxes.List(context.Background(), &SandboxWorkspaceOptions{WorkspaceUUID: "ws-1"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(out.Data) != 1 || out.Data[0].ID != "c1" {
		t.Fatalf("data = %+v", out.Data)
	}
}

func TestSandboxService_Create_Get_Lifecycle(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/sandboxes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var body CreateSandboxRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if body.Name != "my-box" || body.Image != "ubuntu" {
				t.Fatalf("body = %+v", body)
			}
			if r.URL.Query().Get("workspace_uuid") != "ws-1" {
				t.Fatalf("workspace = %q", r.URL.Query().Get("workspace_uuid"))
			}
			w.WriteHeader(http.StatusCreated)
			if _, err := w.Write([]byte(`{"success":true,"message":"sandbox created","data":{"id":"c9","name":"my-box","image":"ubuntu","status":"creating"}}`)); err != nil {
				t.Fatalf("write: %v", err)
			}
		default:
			t.Fatalf("unexpected %s", r.Method)
		}
	})
	mux.HandleFunc("/api/v1/sandboxes/c9", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if _, err := w.Write([]byte(`{"success":true,"data":{"id":"c9","Name":"my-box","Status":"running"}}`)); err != nil {
			t.Fatalf("write: %v", err)
		}
	})
	mux.HandleFunc("/api/v1/sandboxes/c9/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if _, err := w.Write([]byte(`{"success":true,"message":"sandbox started"}`)); err != nil {
			t.Fatalf("write: %v", err)
		}
	})
	mux.HandleFunc("/api/v1/sandboxes/c9/stop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if _, err := w.Write([]byte(`{"success":true,"message":"sandbox stopped"}`)); err != nil {
			t.Fatalf("write: %v", err)
		}
	})
	mux.HandleFunc("/api/v1/sandboxes/c9/session", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if _, err := w.Write([]byte(`{"success":true,"data":{"container_id":"c9","token":"ephemeral_tok","base_url":"https://rexec.sh","expires_in_seconds":900,"token_source":"ephemeral"}}`)); err != nil {
			t.Fatalf("write: %v", err)
		}
	})
	mux.HandleFunc("/api/v1/sandboxes/c9/exec", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		var body ExecSandboxRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Command != "echo hi" {
			t.Fatalf("command = %q", body.Command)
		}
		if r.URL.Query().Get("workspace_uuid") != "ws-1" {
			t.Fatalf("workspace = %q", r.URL.Query().Get("workspace_uuid"))
		}
		if _, err := w.Write([]byte(`{"success":true,"message":"command executed","data":{"sandbox_id":"c9","stdout":"hi\n","output":"hi\n","exit_code":0,"command":"echo hi"}}`)); err != nil {
			t.Fatalf("write: %v", err)
		}
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	opts := &SandboxWorkspaceOptions{WorkspaceUUID: "ws-1"}

	created, _, err := client.Sandboxes.Create(context.Background(), opts, &CreateSandboxRequest{
		Name:  "my-box",
		Image: "ubuntu",
		Role:  "standard",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.Data.ID != "c9" {
		t.Fatalf("created id = %q", created.Data.ID)
	}

	got, _, err := client.Sandboxes.Get(context.Background(), "c9", opts)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Data.Status != "running" || got.Data.Name != "my-box" {
		t.Fatalf("get = %+v", got.Data)
	}

	if _, _, err := client.Sandboxes.Start(context.Background(), "c9", opts); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if _, _, err := client.Sandboxes.Stop(context.Background(), "c9", opts); err != nil {
		t.Fatalf("Stop: %v", err)
	}

	sess, _, err := client.Sandboxes.CreateSession(context.Background(), "c9", opts)
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	if sess.Data.Token != "ephemeral_tok" || sess.Data.ExpiresIn != 900 {
		t.Fatalf("session = %+v", sess.Data)
	}

	execOut, _, err := client.Sandboxes.Exec(context.Background(), "c9", opts, &ExecSandboxRequest{Command: "echo hi"})
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if execOut.Data.ExitCode != 0 || execOut.Data.Output != "hi\n" {
		t.Fatalf("exec = %+v", execOut.Data)
	}
	if _, _, err := client.Sandboxes.Exec(context.Background(), "c9", opts, &ExecSandboxRequest{}); err == nil {
		t.Fatal("expected missing command error")
	}
}

func TestSandboxService_Restart(t *testing.T) {
	t.Parallel()
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if _, err := w.Write([]byte(`{"success":true,"message":"ok"}`)); err != nil {
			t.Fatalf("write: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if _, _, err := client.Sandboxes.Restart(context.Background(), "c1", &SandboxWorkspaceOptions{WorkspaceUUID: "ws"}); err != nil {
		t.Fatalf("Restart: %v", err)
	}
	if len(calls) != 2 || !strings.HasSuffix(calls[0], "/stop") || !strings.HasSuffix(calls[1], "/start") {
		t.Fatalf("calls = %v", calls)
	}
}

func TestSandboxService_MintAPIToken_AndBinding(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/sandboxes/api-token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte(`{"success":true,"data":{"token":"rexec_once","base_url":"https://rexec.sh","token_prefix":"rexec_onc"}}`)); err != nil {
			t.Fatalf("write: %v", err)
		}
	})
	mux.HandleFunc("/api/v1/sandboxes/rexec-binding", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if _, err := w.Write([]byte(`{"success":true,"data":{"workspace_uuid":"ws-1","configured":true,"enabled":true,"token_prefix":"rexec_abc","source":"workspace"}}`)); err != nil {
				t.Fatalf("write: %v", err)
			}
		case http.MethodPut:
			var body UpsertRexecBindingRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if body.Token != "rexec_new" {
				t.Fatalf("token = %q", body.Token)
			}
			if _, err := w.Write([]byte(`{"success":true,"data":{"configured":true,"enabled":true,"token_prefix":"rexec_new"}}`)); err != nil {
				t.Fatalf("write: %v", err)
			}
		case http.MethodDelete:
			if _, err := w.Write([]byte(`{"success":true,"message":"rexec binding removed"}`)); err != nil {
				t.Fatalf("write: %v", err)
			}
		default:
			t.Fatalf("method = %s", r.Method)
		}
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	opts := &SandboxWorkspaceOptions{WorkspaceUUID: "ws-1"}

	tok, _, err := client.Sandboxes.MintAPIToken(context.Background(), opts, &MintRexecAPITokenRequest{Name: "cli"})
	if err != nil {
		t.Fatalf("MintAPIToken: %v", err)
	}
	if tok.Data.Token != "rexec_once" {
		t.Fatalf("token = %q", tok.Data.Token)
	}

	binding, _, err := client.Sandboxes.GetRexecBinding(context.Background(), opts)
	if err != nil {
		t.Fatalf("GetRexecBinding: %v", err)
	}
	if !binding.Data.Configured || binding.Data.Source != "workspace" {
		t.Fatalf("binding = %+v", binding.Data)
	}
	if _, _, err := client.Sandboxes.UpsertRexecBinding(context.Background(), opts, &UpsertRexecBindingRequest{Token: "rexec_new"}); err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	if _, _, err := client.Sandboxes.DeleteRexecBinding(context.Background(), opts); err != nil {
		t.Fatalf("Delete binding: %v", err)
	}
}

func TestSandboxService_UsageDaily(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/sandboxes/usage/daily" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("from") != "2026-08-01" || r.URL.Query().Get("to") != "2026-08-03" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		if _, err := w.Write([]byte(`{"success":true,"data":[{"workspace_uuid":"ws-1","created_count":2,"session_count":5}]}`)); err != nil {
			t.Fatalf("write: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	from := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 8, 3, 0, 0, 0, 0, time.UTC)
	out, _, err := client.Sandboxes.UsageDaily(context.Background(), &SandboxWorkspaceOptions{WorkspaceUUID: "ws-1"}, from, to)
	if err != nil {
		t.Fatalf("UsageDaily: %v", err)
	}
	if len(out.Data) != 1 || out.Data[0].CreatedCount != 2 {
		t.Fatalf("data = %+v", out.Data)
	}
}

func TestSandboxService_RequiresSandboxID(t *testing.T) {
	t.Parallel()
	client, err := NewClient("https://api.pipeops.test")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if _, _, err := client.Sandboxes.Get(context.Background(), "", &SandboxWorkspaceOptions{WorkspaceUUID: "ws"}); err == nil {
		t.Fatal("expected error for empty id")
	}
	if _, _, err := client.Sandboxes.UpsertRexecBinding(context.Background(), &SandboxWorkspaceOptions{WorkspaceUUID: "ws"}, &UpsertRexecBindingRequest{}); err == nil {
		t.Fatal("expected error for empty token")
	}
}
