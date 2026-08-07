package pipeops

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SandboxService talks to the PipeOps sandboxes BFF (Rexec proxy).
//
// Controller mounts the same handlers at /sandboxes and /api/v1/sandboxes.
// This client uses /api/v1/sandboxes (SDK-oriented alias).
//
// Auth: user JWT session or workspace service account (sat_*) with api:read/write
// or preset "sandbox". Workspace is required for multi-tenant correctness
// (query workspace_uuid / workspace, or SA-bound workspace).
//
// This is not a direct Rexec client — PipeOps proxies with workspace/platform
// credentials. Use MintAPIToken only when you need a raw rexec_* for external tools.
type SandboxService struct {
	client *Client
}

const sandboxesBase = "api/v1/sandboxes"

// Sandbox is a Rexec container as returned by the BFF.
type Sandbox struct {
	ID        string            `json:"id,omitempty"`
	UUID      string            `json:"uuid,omitempty"`
	Name      string            `json:"name,omitempty"`
	Image     string            `json:"image,omitempty"`
	Role      string            `json:"role,omitempty"`
	Status    string            `json:"status,omitempty"`
	CreatedAt *Timestamp        `json:"created_at,omitempty"`
	UpdatedAt *Timestamp        `json:"updated_at,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// UnmarshalJSON accepts snake_case and PascalCase field aliases from the BFF.
func (s *Sandbox) UnmarshalJSON(data []byte) error {
	type wire struct {
		ID        string            `json:"id,omitempty"`
		UUID      string            `json:"uuid,omitempty"`
		Name      string            `json:"name,omitempty"`
		NamePC    string            `json:"Name,omitempty"`
		Image     string            `json:"image,omitempty"`
		ImagePC   string            `json:"Image,omitempty"`
		Role      string            `json:"role,omitempty"`
		RolePC    string            `json:"Role,omitempty"`
		Status    string            `json:"status,omitempty"`
		StatusPC  string            `json:"Status,omitempty"`
		CreatedAt *Timestamp        `json:"created_at,omitempty"`
		UpdatedAt *Timestamp        `json:"updated_at,omitempty"`
		Labels    map[string]string `json:"labels,omitempty"`
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	s.ID = w.ID
	s.UUID = w.UUID
	s.Name = coalesceNonEmpty(w.Name, w.NamePC)
	s.Image = coalesceNonEmpty(w.Image, w.ImagePC)
	s.Role = coalesceNonEmpty(w.Role, w.RolePC)
	s.Status = coalesceNonEmpty(w.Status, w.StatusPC)
	s.CreatedAt = w.CreatedAt
	s.UpdatedAt = w.UpdatedAt
	s.Labels = w.Labels
	return nil
}

// CreateSandboxRequest creates a sandbox (empty body uses server defaults).
type CreateSandboxRequest struct {
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
	Role  string `json:"role,omitempty"`
}

// SandboxSession is a short-lived terminal/embed grant from POST .../session.
type SandboxSession struct {
	ContainerID string `json:"container_id,omitempty"`
	Token       string `json:"token,omitempty"`
	BaseURL     string `json:"base_url,omitempty"`
	ExpiresIn   int    `json:"expires_in_seconds,omitempty"`
	TokenSource string `json:"token_source,omitempty"` // ephemeral | workspace | platform
	GrantID     string `json:"grant_id,omitempty"`
}

// UnmarshalJSON also accepts Token PascalCase alias from the BFF.
func (s *SandboxSession) UnmarshalJSON(data []byte) error {
	type wire struct {
		ContainerID string `json:"container_id,omitempty"`
		Token       string `json:"token,omitempty"`
		TokenPC     string `json:"Token,omitempty"`
		BaseURL     string `json:"base_url,omitempty"`
		ExpiresIn   int    `json:"expires_in_seconds,omitempty"`
		TokenSource string `json:"token_source,omitempty"`
		GrantID     string `json:"grant_id,omitempty"`
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	s.ContainerID = w.ContainerID
	s.Token = coalesceNonEmpty(w.Token, w.TokenPC)
	s.BaseURL = w.BaseURL
	s.ExpiresIn = w.ExpiresIn
	s.TokenSource = w.TokenSource
	s.GrantID = w.GrantID
	return nil
}

// MintRexecAPITokenRequest is POST /api/v1/sandboxes/api-token.
type MintRexecAPITokenRequest struct {
	Name          string `json:"name,omitempty"`
	ExpiresInDays *int   `json:"expires_in_days,omitempty"` // default 90, max 365
}

// MintRexecAPITokenResult is returned once; store the token client-side.
type MintRexecAPITokenResult struct {
	Token       string     `json:"token,omitempty"`
	TokenID     string     `json:"token_id,omitempty"`
	TokenPrefix string     `json:"token_prefix,omitempty"`
	Name        string     `json:"name,omitempty"`
	Scopes      []string   `json:"scopes,omitempty"`
	BaseURL     string     `json:"base_url,omitempty"`
	ExpiresAt   *Timestamp `json:"expires_at,omitempty"`
	UsageHint   string     `json:"usage_hint,omitempty"`
}

// RexecBinding is a safe view of workspace Rexec credentials (no secret).
type RexecBinding struct {
	WorkspaceUUID string     `json:"workspace_uuid,omitempty"`
	BaseURL       string     `json:"base_url,omitempty"`
	TokenPrefix   string     `json:"token_prefix,omitempty"`
	Enabled       bool       `json:"enabled"`
	Configured    bool       `json:"configured"`
	LastUsedAt    *Timestamp `json:"last_used_at,omitempty"`
	UpdatedAt     *Timestamp `json:"updated_at,omitempty"`
	Source        string     `json:"source,omitempty"` // workspace | platform
}

// UpsertRexecBindingRequest sets an optional workspace-owned Rexec API token.
type UpsertRexecBindingRequest struct {
	Token   string `json:"token"`
	BaseURL string `json:"base_url,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
}

// SandboxUsageDaily is a billing rollup row for a workspace day.
type SandboxUsageDaily struct {
	WorkspaceUUID        string     `json:"workspace_uuid,omitempty"`
	Day                  *Timestamp `json:"day,omitempty"`
	CreatedCount         int64      `json:"created_count,omitempty"`
	StartedCount         int64      `json:"started_count,omitempty"`
	StoppedCount         int64      `json:"stopped_count,omitempty"`
	DeletedCount         int64      `json:"deleted_count,omitempty"`
	SessionCount         int64      `json:"session_count,omitempty"`
	TotalDurationSeconds int64      `json:"total_duration_seconds,omitempty"`
	UniqueContainers     int64      `json:"unique_containers,omitempty"`
}

// SandboxWorkspaceOptions scopes sandbox calls to a workspace.
// Prefer WorkspaceUUID; Workspace is accepted as a controller alias.
type SandboxWorkspaceOptions struct {
	WorkspaceUUID string `url:"workspace_uuid,omitempty"`
	Workspace     string `url:"workspace,omitempty"`
}

// SandboxListResponse is GET /api/v1/sandboxes.
type SandboxListResponse struct {
	Success bool      `json:"success,omitempty"`
	Message string    `json:"message,omitempty"`
	Data    []Sandbox `json:"data"`
	Meta    struct {
		Count int `json:"count"`
	} `json:"meta"`
}

// SandboxResponse is GET/POST single-sandbox envelopes.
type SandboxResponse struct {
	Success bool    `json:"success,omitempty"`
	Message string  `json:"message,omitempty"`
	Data    Sandbox `json:"data"`
}

// SandboxSessionResponse is POST .../session.
type SandboxSessionResponse struct {
	Success bool           `json:"success,omitempty"`
	Message string         `json:"message,omitempty"`
	Data    SandboxSession `json:"data"`
}

// MintRexecAPITokenResponse is POST .../api-token.
type MintRexecAPITokenResponse struct {
	Success bool                    `json:"success,omitempty"`
	Message string                  `json:"message,omitempty"`
	Data    MintRexecAPITokenResult `json:"data"`
}

// RexecBindingResponse is GET/PUT .../rexec-binding.
type RexecBindingResponse struct {
	Success bool         `json:"success,omitempty"`
	Message string       `json:"message,omitempty"`
	Data    RexecBinding `json:"data"`
}

// SandboxUsageDailyResponse is GET .../usage/daily.
type SandboxUsageDailyResponse struct {
	Success bool                `json:"success,omitempty"`
	Message string              `json:"message,omitempty"`
	Data    []SandboxUsageDaily `json:"data"`
}

// MessageOnlyResponse is start/stop/delete success without a body payload.
type MessageOnlyResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

// ExecSandboxRequest is POST /api/v1/sandboxes/:id/exec.
// Prefer Command for shell strings; Cmd is argv when you need no shell.
type ExecSandboxRequest struct {
	Command        string   `json:"command,omitempty"`
	Cmd            []string `json:"cmd,omitempty"`
	WorkDir        string   `json:"workdir,omitempty"`
	Env            []string `json:"env,omitempty"`
	User           string   `json:"user,omitempty"`
	TimeoutSeconds int      `json:"timeout_seconds,omitempty"` // default 60, max 300
}

// ExecSandboxResult is the captured output from a non-interactive sandbox exec.
type ExecSandboxResult struct {
	SandboxID string   `json:"sandbox_id,omitempty"`
	Stdout    string   `json:"stdout,omitempty"`
	Stderr    string   `json:"stderr,omitempty"`
	Output    string   `json:"output,omitempty"`
	ExitCode  int      `json:"exit_code"`
	Command   string   `json:"command,omitempty"`
	Cmd       []string `json:"cmd,omitempty"`
	Truncated bool     `json:"truncated,omitempty"`
}

// ExecSandboxResponse is the BFF envelope for POST .../exec.
type ExecSandboxResponse struct {
	Success bool              `json:"success,omitempty"`
	Message string            `json:"message,omitempty"`
	Data    ExecSandboxResult `json:"data"`
}

// SandboxFileInfo is one directory entry from ListFiles.
type SandboxFileInfo struct {
	Name  string `json:"name,omitempty"`
	Path  string `json:"path,omitempty"`
	Size  int64  `json:"size,omitempty"`
	Mode  string `json:"mode,omitempty"`
	IsDir bool   `json:"is_dir"`
}

// SandboxFileList is GET /api/v1/sandboxes/:id/files data.
type SandboxFileList struct {
	SandboxID string            `json:"sandbox_id,omitempty"`
	Path      string            `json:"path,omitempty"`
	Files     []SandboxFileInfo `json:"files"`
	Count     int               `json:"count"`
}

// SandboxFileListResponse is the BFF envelope for list files.
type SandboxFileListResponse struct {
	Success bool            `json:"success,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    SandboxFileList `json:"data"`
}

// SandboxFileContent is GET /api/v1/sandboxes/:id/files/content data.
// Encoding is "utf-8" for text or "base64" for binary.
type SandboxFileContent struct {
	SandboxID string `json:"sandbox_id,omitempty"`
	Path      string `json:"path,omitempty"`
	Content   string `json:"content,omitempty"`
	Encoding  string `json:"encoding,omitempty"`
	Size      int    `json:"size"`
	Truncated bool   `json:"truncated,omitempty"`
}

// SandboxFileContentResponse is the BFF envelope for read file.
type SandboxFileContentResponse struct {
	Success bool               `json:"success,omitempty"`
	Message string             `json:"message,omitempty"`
	Data    SandboxFileContent `json:"data"`
}

// List lists sandboxes for a workspace.
// GET /api/v1/sandboxes?workspace_uuid=
func (s *SandboxService) List(ctx context.Context, opts *SandboxWorkspaceOptions) (*SandboxListResponse, *http.Response, error) {
	u, err := withSandboxWorkspaceQuery(sandboxesBase, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	out := new(SandboxListResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Get returns one sandbox by id.
// GET /api/v1/sandboxes/:id?workspace_uuid=
func (s *SandboxService) Get(ctx context.Context, sandboxID string, opts *SandboxWorkspaceOptions) (*SandboxResponse, *http.Response, error) {
	if strings.TrimSpace(sandboxID) == "" {
		return nil, nil, errors.New("sandbox id is required")
	}
	u, err := withSandboxWorkspaceQuery(fmt.Sprintf("%s/%s", sandboxesBase, url.PathEscape(sandboxID)), opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	out := new(SandboxResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Create creates a sandbox. Empty req is allowed (server defaults).
// POST /api/v1/sandboxes?workspace_uuid=
func (s *SandboxService) Create(ctx context.Context, opts *SandboxWorkspaceOptions, body *CreateSandboxRequest) (*SandboxResponse, *http.Response, error) {
	u, err := withSandboxWorkspaceQuery(sandboxesBase, opts)
	if err != nil {
		return nil, nil, err
	}
	if body == nil {
		body = &CreateSandboxRequest{}
	}
	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}
	out := new(SandboxResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Start starts a stopped sandbox.
// POST /api/v1/sandboxes/:id/start?workspace_uuid=
func (s *SandboxService) Start(ctx context.Context, sandboxID string, opts *SandboxWorkspaceOptions) (*MessageOnlyResponse, *http.Response, error) {
	return s.postAction(ctx, sandboxID, "start", opts)
}

// Stop stops a running sandbox.
// POST /api/v1/sandboxes/:id/stop?workspace_uuid=
func (s *SandboxService) Stop(ctx context.Context, sandboxID string, opts *SandboxWorkspaceOptions) (*MessageOnlyResponse, *http.Response, error) {
	return s.postAction(ctx, sandboxID, "stop", opts)
}

// Restart stops then starts a sandbox (dashboard convenience).
func (s *SandboxService) Restart(ctx context.Context, sandboxID string, opts *SandboxWorkspaceOptions) (*MessageOnlyResponse, *http.Response, error) {
	if _, resp, err := s.Stop(ctx, sandboxID, opts); err != nil {
		return nil, resp, err
	}
	return s.Start(ctx, sandboxID, opts)
}

// Delete deletes a sandbox.
// DELETE /api/v1/sandboxes/:id?workspace_uuid=
func (s *SandboxService) Delete(ctx context.Context, sandboxID string, opts *SandboxWorkspaceOptions) (*MessageOnlyResponse, *http.Response, error) {
	if strings.TrimSpace(sandboxID) == "" {
		return nil, nil, errors.New("sandbox id is required")
	}
	u, err := withSandboxWorkspaceQuery(fmt.Sprintf("%s/%s", sandboxesBase, url.PathEscape(sandboxID)), opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, nil, err
	}
	out := new(MessageOnlyResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Exec runs a non-interactive command inside a running sandbox.
// POST /api/v1/sandboxes/:id/exec?workspace_uuid=
// Body requires Command (shell string) and/or Cmd (argv).
func (s *SandboxService) Exec(ctx context.Context, sandboxID string, opts *SandboxWorkspaceOptions, body *ExecSandboxRequest) (*ExecSandboxResponse, *http.Response, error) {
	if strings.TrimSpace(sandboxID) == "" {
		return nil, nil, errors.New("sandbox id is required")
	}
	if body == nil {
		return nil, nil, errors.New("exec body is required")
	}
	if strings.TrimSpace(body.Command) == "" && len(body.Cmd) == 0 {
		return nil, nil, errors.New("command or cmd is required")
	}
	u, err := withSandboxWorkspaceQuery(fmt.Sprintf("%s/%s/exec", sandboxesBase, url.PathEscape(sandboxID)), opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}
	out := new(ExecSandboxResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListFiles lists a directory inside a running sandbox.
// GET /api/v1/sandboxes/:id/files?workspace_uuid=&path=
// Empty path defaults to /home/user on the server.
func (s *SandboxService) ListFiles(ctx context.Context, sandboxID, path string, opts *SandboxWorkspaceOptions) (*SandboxFileListResponse, *http.Response, error) {
	if strings.TrimSpace(sandboxID) == "" {
		return nil, nil, errors.New("sandbox id is required")
	}
	u, err := withSandboxWorkspaceQuery(fmt.Sprintf("%s/%s/files", sandboxesBase, url.PathEscape(sandboxID)), opts)
	if err != nil {
		return nil, nil, err
	}
	if p := strings.TrimSpace(path); p != "" {
		if strings.Contains(u, "?") {
			u += "&path=" + url.QueryEscape(p)
		} else {
			u += "?path=" + url.QueryEscape(p)
		}
	}
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	out := new(SandboxFileListResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ReadFile reads a file from a running sandbox (UTF-8 text or base64).
// GET /api/v1/sandboxes/:id/files/content?workspace_uuid=&path=
func (s *SandboxService) ReadFile(ctx context.Context, sandboxID, path string, opts *SandboxWorkspaceOptions) (*SandboxFileContentResponse, *http.Response, error) {
	if strings.TrimSpace(sandboxID) == "" {
		return nil, nil, errors.New("sandbox id is required")
	}
	if strings.TrimSpace(path) == "" {
		return nil, nil, errors.New("path is required")
	}
	u, err := withSandboxWorkspaceQuery(fmt.Sprintf("%s/%s/files/content", sandboxesBase, url.PathEscape(sandboxID)), opts)
	if err != nil {
		return nil, nil, err
	}
	if strings.Contains(u, "?") {
		u += "&path=" + url.QueryEscape(path)
	} else {
		u += "?path=" + url.QueryEscape(path)
	}
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	out := new(SandboxFileContentResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// CreateSession mints a short-lived terminal/session grant for a sandbox.
// POST /api/v1/sandboxes/:id/session?workspace_uuid=
func (s *SandboxService) CreateSession(ctx context.Context, sandboxID string, opts *SandboxWorkspaceOptions) (*SandboxSessionResponse, *http.Response, error) {
	if strings.TrimSpace(sandboxID) == "" {
		return nil, nil, errors.New("sandbox id is required")
	}
	u, err := withSandboxWorkspaceQuery(fmt.Sprintf("%s/%s/session", sandboxesBase, url.PathEscape(sandboxID)), opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, nil, err
	}
	out := new(SandboxSessionResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// MintAPIToken mints a long-lived Rexec API token (rexec_*). Shown once.
// POST /api/v1/sandboxes/api-token?workspace_uuid=
func (s *SandboxService) MintAPIToken(ctx context.Context, opts *SandboxWorkspaceOptions, body *MintRexecAPITokenRequest) (*MintRexecAPITokenResponse, *http.Response, error) {
	u, err := withSandboxWorkspaceQuery(sandboxesBase+"/api-token", opts)
	if err != nil {
		return nil, nil, err
	}
	if body == nil {
		body = &MintRexecAPITokenRequest{}
	}
	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}
	out := new(MintRexecAPITokenResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// GetRexecBinding returns workspace Rexec credential status (no secret).
// GET /api/v1/sandboxes/rexec-binding?workspace_uuid=
func (s *SandboxService) GetRexecBinding(ctx context.Context, opts *SandboxWorkspaceOptions) (*RexecBindingResponse, *http.Response, error) {
	u, err := withSandboxWorkspaceQuery(sandboxesBase+"/rexec-binding", opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	out := new(RexecBindingResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// UpsertRexecBinding sets a workspace-owned Rexec API token (BYOS).
// PUT /api/v1/sandboxes/rexec-binding?workspace_uuid=
func (s *SandboxService) UpsertRexecBinding(ctx context.Context, opts *SandboxWorkspaceOptions, body *UpsertRexecBindingRequest) (*RexecBindingResponse, *http.Response, error) {
	if body == nil || strings.TrimSpace(body.Token) == "" {
		return nil, nil, errors.New("token is required")
	}
	u, err := withSandboxWorkspaceQuery(sandboxesBase+"/rexec-binding", opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodPut, u, body)
	if err != nil {
		return nil, nil, err
	}
	out := new(RexecBindingResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// DeleteRexecBinding removes the workspace Rexec binding.
// DELETE /api/v1/sandboxes/rexec-binding?workspace_uuid=
func (s *SandboxService) DeleteRexecBinding(ctx context.Context, opts *SandboxWorkspaceOptions) (*MessageOnlyResponse, *http.Response, error) {
	u, err := withSandboxWorkspaceQuery(sandboxesBase+"/rexec-binding", opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, nil, err
	}
	out := new(MessageOnlyResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// UsageDaily returns usage rollups for a workspace day range (inclusive).
// GET /api/v1/sandboxes/usage/daily?workspace_uuid=&from=&to=
// from/to use YYYY-MM-DD. Zero times omit the corresponding query param.
func (s *SandboxService) UsageDaily(ctx context.Context, opts *SandboxWorkspaceOptions, from, to time.Time) (*SandboxUsageDailyResponse, *http.Response, error) {
	u, err := withSandboxWorkspaceQuery(sandboxesBase+"/usage/daily", opts)
	if err != nil {
		return nil, nil, err
	}
	values := url.Values{}
	if !from.IsZero() {
		values.Set("from", from.Format("2006-01-02"))
	}
	if !to.IsZero() {
		values.Set("to", to.Format("2006-01-02"))
	}
	if encoded := values.Encode(); encoded != "" {
		if strings.Contains(u, "?") {
			u += "&" + encoded
		} else {
			u += "?" + encoded
		}
	}
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	out := new(SandboxUsageDailyResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

func (s *SandboxService) postAction(ctx context.Context, sandboxID, action string, opts *SandboxWorkspaceOptions) (*MessageOnlyResponse, *http.Response, error) {
	if strings.TrimSpace(sandboxID) == "" {
		return nil, nil, errors.New("sandbox id is required")
	}
	u, err := withSandboxWorkspaceQuery(fmt.Sprintf("%s/%s/%s", sandboxesBase, url.PathEscape(sandboxID), action), opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, nil, err
	}
	out := new(MessageOnlyResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// withSandboxWorkspaceQuery appends workspace query params.
// Requires an explicit workspace (UUID) for user JWT multi-tenant safety.
// SA-bound tokens may still work if the controller resolves workspace from the token
// when query is empty — callers should still pass workspace when known.
func withSandboxWorkspaceQuery(path string, opts *SandboxWorkspaceOptions) (string, error) {
	if opts == nil {
		opts = &SandboxWorkspaceOptions{}
	}
	ws := strings.TrimSpace(opts.WorkspaceUUID)
	if ws == "" {
		ws = strings.TrimSpace(opts.Workspace)
	}
	// Prefer both keys for dashboard/controller compatibility when set.
	queryOpts := &SandboxWorkspaceOptions{}
	if ws != "" {
		queryOpts.WorkspaceUUID = ws
		queryOpts.Workspace = ws
	}
	return addOptions(path, queryOpts)
}
