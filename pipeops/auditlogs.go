package pipeops

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// AuditLogService lists project- and workspace-scoped activity (who did what).
//
// Controller routes (user JWT / team session; project read permission for project path):
//
//	GET /project/audit-logs/:uuid?limit=&offset=&action=&actor_type=&category=&search=&from=&to=
//	GET /project/workspace-audit-logs?workspace_uuid=&project_uuid=&limit=&offset=&...
//
// These are the console “audit log” surfaces for historical actions (deploy,
// env change, domain, pause/resume, agent/webhook deploys, etc.). They are
// distinct from admin-only staff audit tables.
type AuditLogService struct {
	client *Client
}

// ProjectAuditActor is the actor summary returned on list items (email omitted for PII).
type ProjectAuditActor struct {
	Type      string `json:"type,omitempty"` // user | webhook | system | service_account | agent
	UUID      string `json:"uuid,omitempty"`
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Label     string `json:"label,omitempty"`
}

// ProjectAuditLog is one project-scoped historical action.
type ProjectAuditLog struct {
	UUID         string                 `json:"uuid,omitempty"`
	Action       string                 `json:"action,omitempty"`       // e.g. project.redeploy
	ActionLabel  string                 `json:"action_label,omitempty"` // human label
	Category     string                 `json:"category,omitempty"`     // lifecycle | settings | deployment | security | access
	Status       string                 `json:"status,omitempty"`       // success | failure | attempted
	Summary      string                 `json:"summary,omitempty"`
	ProjectUUID  string                 `json:"project_uuid,omitempty"`
	ProjectName  string                 `json:"project_name,omitempty"`
	ResourceType string                 `json:"resource_type,omitempty"`
	ResourceUUID string                 `json:"resource_uuid,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    *Timestamp             `json:"created_at,omitempty"`
	Actor        ProjectAuditActor      `json:"actor,omitempty"`
}

// AuditLogPagination is the list envelope pagination block.
type AuditLogPagination struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

// ProjectAuditLogListResponse is GET /project/audit-logs/:uuid.
type ProjectAuditLogListResponse struct {
	Success    bool               `json:"success,omitempty"`
	Message    string             `json:"message,omitempty"`
	Data       []ProjectAuditLog  `json:"data"`
	Pagination AuditLogPagination `json:"pagination"`
}

// WorkspaceAuditLogListResponse is GET /project/workspace-audit-logs.
type WorkspaceAuditLogListResponse struct {
	Success    bool               `json:"success,omitempty"`
	Message    string             `json:"message,omitempty"`
	Data       []ProjectAuditLog  `json:"data"`
	Pagination AuditLogPagination `json:"pagination"`
}

// ProjectAuditLogListOptions filters GET /project/audit-logs/:uuid.
type ProjectAuditLogListOptions struct {
	// Action is a single action or comma-separated list (project.redeploy,project.env.update).
	Action        string `url:"action,omitempty"`
	ActorUserUUID string `url:"actor_user_uuid,omitempty"`
	ActorType     string `url:"actor_type,omitempty"` // user | webhook | system | service_account | agent
	Category      string `url:"category,omitempty"`   // lifecycle | settings | deployment | security | access
	Search        string `url:"search,omitempty"`     // free-text over summary / names
	// From / To are RFC3339 timestamps (controller parses with time.RFC3339).
	From   string `url:"from,omitempty"`
	To     string `url:"to,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

// WorkspaceAuditLogListOptions filters GET /project/workspace-audit-logs.
// WorkspaceUUID is required by the controller (unless workspace is already in session context).
type WorkspaceAuditLogListOptions struct {
	WorkspaceUUID string `url:"workspace_uuid,omitempty"`
	// ProjectUUID optionally narrows the workspace feed to one project.
	ProjectUUID   string `url:"project_uuid,omitempty"`
	Action        string `url:"action,omitempty"`
	ActorUserUUID string `url:"actor_user_uuid,omitempty"`
	ActorType     string `url:"actor_type,omitempty"`
	Category      string `url:"category,omitempty"`
	Search        string `url:"search,omitempty"`
	From          string `url:"from,omitempty"`
	To            string `url:"to,omitempty"`
	Limit         int    `url:"limit,omitempty"`
	Offset        int    `url:"offset,omitempty"`
}

// ListProject returns historical actions for one project.
// GET /project/audit-logs/:projectUUID
func (s *AuditLogService) ListProject(ctx context.Context, projectUUID string, opts *ProjectAuditLogListOptions) (*ProjectAuditLogListResponse, *http.Response, error) {
	projectUUID = strings.TrimSpace(projectUUID)
	if projectUUID == "" {
		return nil, nil, fmt.Errorf("project_uuid is required")
	}
	u := fmt.Sprintf("project/audit-logs/%s", projectUUID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectAuditLogListResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListWorkspace returns historical actions across projects in a workspace.
// GET /project/workspace-audit-logs?workspace_uuid=
//
// If opts.WorkspaceUUID is empty, the SDK attempts firstWorkspaceUUID for
// prefer-client single-workspace accounts.
func (s *AuditLogService) ListWorkspace(ctx context.Context, opts *WorkspaceAuditLogListOptions) (*WorkspaceAuditLogListResponse, *http.Response, error) {
	if opts == nil {
		opts = &WorkspaceAuditLogListOptions{}
	}
	if strings.TrimSpace(opts.WorkspaceUUID) == "" {
		if ws, _, err := firstWorkspaceUUID(ctx, s.client); err == nil && ws != "" {
			opts.WorkspaceUUID = ws
		}
	}
	if strings.TrimSpace(opts.WorkspaceUUID) == "" {
		return nil, nil, fmt.Errorf("workspace_uuid is required")
	}

	u, err := addOptions("project/workspace-audit-logs", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(WorkspaceAuditLogListResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAuditLogs is a convenience alias for ListWorkspace (workspace-wide feed).
// Prefer ListWorkspace or ListProject for explicit scope.
//
// Deprecated: use ListWorkspace or ListProject. Kept so existing AuditLogs
// callers that expected a list method still compile; they previously hit a
// non-existent /audit/logs path.
func (s *AuditLogService) ListAuditLogs(ctx context.Context) (*WorkspaceAuditLogListResponse, *http.Response, error) {
	return s.ListWorkspace(ctx, nil)
}

// JoinAuditActions joins multiple action codes for the action= query param.
func JoinAuditActions(actions ...string) string {
	var parts []string
	for _, a := range actions {
		a = strings.TrimSpace(a)
		if a != "" {
			parts = append(parts, a)
		}
	}
	return strings.Join(parts, ",")
}
