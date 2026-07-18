package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// GitOpsService handles GitOps application configuration APIs.
//
// Controller routes (JWT session):
//
//	POST   /api/v1/gitops/applications
//	GET    /api/v1/gitops/applications
//	GET    /api/v1/gitops/applications/:uuid
//	PUT    /api/v1/gitops/applications/:uuid
//	DELETE /api/v1/gitops/applications/:uuid
//	POST   /api/v1/gitops/applications/:uuid/sync
//	GET    /api/v1/gitops/applications/:uuid/sync-status
//	GET    /api/v1/gitops/applications/:uuid/diff
//	GET    /api/v1/gitops/applications/:uuid/history
type GitOpsService struct {
	client *Client
}

// CreateGitOpsConfigRequest is POST /api/v1/gitops/applications.
type CreateGitOpsConfigRequest struct {
	Name          string `json:"name"`
	ProjectID     *uint  `json:"project_id,omitempty"`
	EnvironmentID *uint  `json:"environment_id,omitempty"`

	RepoURL        string `json:"repo_url"`
	Branch         string `json:"branch,omitempty"`
	Path           string `json:"path,omitempty"`
	TargetRevision string `json:"target_revision,omitempty"`
	ManifestType   string `json:"manifest_type,omitempty"` // pipeops | kubernetes

	SyncPolicy *GitOpsSyncPolicyRequest `json:"sync_policy,omitempty"`

	HealthCheckEnabled  *bool `json:"health_check_enabled,omitempty"`
	HealthCheckInterval int   `json:"health_check_interval,omitempty"`
}

// UpdateGitOpsConfigRequest is PUT /api/v1/gitops/applications/:uuid.
type UpdateGitOpsConfigRequest struct {
	Name           string `json:"name,omitempty"`
	Branch         string `json:"branch,omitempty"`
	Path           string `json:"path,omitempty"`
	TargetRevision string `json:"target_revision,omitempty"`

	SyncPolicy *GitOpsSyncPolicyRequest `json:"sync_policy,omitempty"`

	HealthCheckEnabled  *bool `json:"health_check_enabled,omitempty"`
	HealthCheckInterval *int  `json:"health_check_interval,omitempty"`
}

// GitOpsSyncPolicyRequest is the sync policy payload on create/update.
type GitOpsSyncPolicyRequest struct {
	Automated   *GitOpsAutomatedSyncRequest `json:"automated,omitempty"`
	SyncOptions []string                    `json:"sync_options,omitempty"`
	Retry       *GitOpsRetryStrategyRequest `json:"retry,omitempty"`
}

// GitOpsAutomatedSyncRequest configures auto-sync.
type GitOpsAutomatedSyncRequest struct {
	Prune      bool `json:"prune"`
	SelfHeal   bool `json:"self_heal"`
	AllowEmpty bool `json:"allow_empty"`
}

// GitOpsRetryStrategyRequest configures retry backoff for failed syncs.
type GitOpsRetryStrategyRequest struct {
	Limit              int `json:"limit"`
	BackoffDuration    int `json:"backoff_duration"`
	BackoffFactor      int `json:"backoff_factor"`
	BackoffMaxDuration int `json:"backoff_max_duration"`
}

// TriggerGitOpsSyncRequest is POST /api/v1/gitops/applications/:uuid/sync.
type TriggerGitOpsSyncRequest struct {
	Revision string `json:"revision,omitempty"`
	Prune    bool   `json:"prune,omitempty"`
	DryRun   bool   `json:"dry_run,omitempty"`
}

// GitOpsSyncPolicy is the policy stored on a config.
type GitOpsSyncPolicy struct {
	Automated   *GitOpsAutomatedSync `json:"automated,omitempty"`
	SyncOptions []string             `json:"sync_options,omitempty"`
	Retry       *GitOpsRetryStrategy `json:"retry,omitempty"`
}

// GitOpsAutomatedSync is stored automated sync settings.
type GitOpsAutomatedSync struct {
	Prune      bool `json:"prune"`
	SelfHeal   bool `json:"self_heal"`
	AllowEmpty bool `json:"allow_empty"`
}

// GitOpsRetryStrategy is stored retry settings.
type GitOpsRetryStrategy struct {
	Limit              int `json:"limit"`
	BackoffDuration    int `json:"backoff_duration"`
	BackoffFactor      int `json:"backoff_factor"`
	BackoffMaxDuration int `json:"backoff_max_duration"`
}

// GitOpsConfig is a GitOps application configuration.
type GitOpsConfig struct {
	UUID                string           `json:"uuid,omitempty"`
	Name                string           `json:"name,omitempty"`
	ProjectID           *uint            `json:"project_id,omitempty"`
	ProjectName         string           `json:"project_name,omitempty"`
	EnvironmentID       *uint            `json:"environment_id,omitempty"`
	EnvironmentName     string           `json:"environment_name,omitempty"`
	RepoURL             string           `json:"repo_url,omitempty"`
	Branch              string           `json:"branch,omitempty"`
	Path                string           `json:"path,omitempty"`
	TargetRevision      string           `json:"target_revision,omitempty"`
	SyncPolicy          GitOpsSyncPolicy `json:"sync_policy,omitempty"`
	HealthCheckEnabled  bool             `json:"health_check_enabled,omitempty"`
	HealthCheckInterval int              `json:"health_check_interval,omitempty"`
	LastSyncedCommit    string           `json:"last_synced_commit,omitempty"`
	LastSyncedAt        string           `json:"last_synced_at,omitempty"`
	SyncStatus          string           `json:"sync_status,omitempty"`
	SyncMessage         string           `json:"sync_message,omitempty"`
	HealthStatus        string           `json:"health_status,omitempty"`
	HealthMessage       string           `json:"health_message,omitempty"`
	CreatedAt           string           `json:"created_at,omitempty"`
	UpdatedAt           string           `json:"updated_at,omitempty"`
}

// GitOpsResourceChange is a single resource change in a diff.
type GitOpsResourceChange struct {
	Kind     string      `json:"kind,omitempty"`
	Name     string      `json:"name,omitempty"`
	Field    string      `json:"field,omitempty"`
	OldValue interface{} `json:"old_value,omitempty"`
	NewValue interface{} `json:"new_value,omitempty"`
}

// GitOpsDiffSnapshot captures added/modified/removed resources.
type GitOpsDiffSnapshot struct {
	Added    []GitOpsResourceChange `json:"added,omitempty"`
	Modified []GitOpsResourceChange `json:"modified,omitempty"`
	Removed  []GitOpsResourceChange `json:"removed,omitempty"`
}

// GitOpsSyncHistoryEntry is one sync history row.
type GitOpsSyncHistoryEntry struct {
	ID            uint                `json:"id,omitempty"`
	CommitSHA     string              `json:"commit_sha,omitempty"`
	CommitMessage string              `json:"commit_message,omitempty"`
	CommitAuthor  string              `json:"commit_author,omitempty"`
	SyncStatus    string              `json:"sync_status,omitempty"`
	SyncMessage   string              `json:"sync_message,omitempty"`
	StartedAt     string              `json:"started_at,omitempty"`
	FinishedAt    string              `json:"finished_at,omitempty"`
	DurationMs    int                 `json:"duration_ms,omitempty"`
	TriggeredBy   string              `json:"triggered_by,omitempty"`
	DiffSnapshot  *GitOpsDiffSnapshot `json:"diff_snapshot,omitempty"`
	CreatedAt     string              `json:"created_at,omitempty"`
}

// GitOpsListOptions filters list/history pagination.
type GitOpsListOptions struct {
	Page  int `url:"page,omitempty"`
	Limit int `url:"limit,omitempty"`
}

// GitOpsConfigResponse is a single config envelope.
type GitOpsConfigResponse struct {
	Success bool         `json:"success,omitempty"`
	Message string       `json:"message,omitempty"`
	Data    GitOpsConfig `json:"data"`
}

// GitOpsListResponse is GET /api/v1/gitops/applications.
type GitOpsListResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Items      []GitOpsConfig `json:"items"`
		Total      int64          `json:"total"`
		Page       int            `json:"page"`
		Limit      int            `json:"limit"`
		TotalPages int            `json:"total_pages"`
	} `json:"data"`
}

// GitOpsSyncTriggerResponse is POST .../sync.
type GitOpsSyncTriggerResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Status   string `json:"status,omitempty"`
		Revision string `json:"revision,omitempty"`
		DryRun   bool   `json:"dry_run,omitempty"`
	} `json:"data"`
}

// GitOpsSyncStatusResponse is GET .../sync-status.
type GitOpsSyncStatusResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		SyncStatus       string  `json:"sync_status,omitempty"`
		SyncMessage      string  `json:"sync_message,omitempty"`
		LastSyncedCommit string  `json:"last_synced_commit,omitempty"`
		LastSyncedAt     *string `json:"last_synced_at,omitempty"`
		HealthStatus     string  `json:"health_status,omitempty"`
		HealthMessage    string  `json:"health_message,omitempty"`
	} `json:"data"`
}

// GitOpsDiffResponse is GET .../diff.
type GitOpsDiffResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		CurrentCommit string              `json:"current_commit,omitempty"`
		TargetCommit  string              `json:"target_commit,omitempty"`
		Diff          *GitOpsDiffSnapshot `json:"diff,omitempty"`
		SyncRequired  bool                `json:"sync_required,omitempty"`
	} `json:"data"`
}

// GitOpsSyncHistoryResponse is GET .../history.
type GitOpsSyncHistoryResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Items      []GitOpsSyncHistoryEntry `json:"items"`
		Total      int64                    `json:"total"`
		Page       int                      `json:"page"`
		Limit      int                      `json:"limit"`
		TotalPages int                      `json:"total_pages"`
	} `json:"data"`
}

// Create creates a new GitOps application configuration.
// POST /api/v1/gitops/applications
func (s *GitOpsService) Create(ctx context.Context, body *CreateGitOpsConfigRequest) (*GitOpsConfigResponse, *http.Response, error) {
	u := "api/v1/gitops/applications"

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}

	out := new(GitOpsConfigResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// List returns paginated GitOps applications for the session workspace.
// GET /api/v1/gitops/applications?page=&limit=
func (s *GitOpsService) List(ctx context.Context, opts *GitOpsListOptions) (*GitOpsListResponse, *http.Response, error) {
	u := "api/v1/gitops/applications"
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(GitOpsListResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Get returns one GitOps application by UUID.
// GET /api/v1/gitops/applications/:uuid
func (s *GitOpsService) Get(ctx context.Context, uuid string) (*GitOpsConfigResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/gitops/applications/%s", uuid)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(GitOpsConfigResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Update updates a GitOps application configuration.
// PUT /api/v1/gitops/applications/:uuid
func (s *GitOpsService) Update(ctx context.Context, uuid string, body *UpdateGitOpsConfigRequest) (*GitOpsConfigResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/gitops/applications/%s", uuid)

	req, err := s.client.NewRequest(http.MethodPut, u, body)
	if err != nil {
		return nil, nil, err
	}

	out := new(GitOpsConfigResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Delete removes a GitOps application configuration.
// DELETE /api/v1/gitops/applications/:uuid
func (s *GitOpsService) Delete(ctx context.Context, uuid string) (*http.Response, error) {
	u := fmt.Sprintf("api/v1/gitops/applications/%s", uuid)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// TriggerSync starts a manual sync for a GitOps application.
// POST /api/v1/gitops/applications/:uuid/sync
func (s *GitOpsService) TriggerSync(ctx context.Context, uuid string, body *TriggerGitOpsSyncRequest) (*GitOpsSyncTriggerResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/gitops/applications/%s/sync", uuid)

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}

	out := new(GitOpsSyncTriggerResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// GetSyncStatus returns the current sync/health status.
// GET /api/v1/gitops/applications/:uuid/sync-status
func (s *GitOpsService) GetSyncStatus(ctx context.Context, uuid string) (*GitOpsSyncStatusResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/gitops/applications/%s/sync-status", uuid)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(GitOpsSyncStatusResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// GetDiff returns the git vs live state diff for an application.
// GET /api/v1/gitops/applications/:uuid/diff
func (s *GitOpsService) GetDiff(ctx context.Context, uuid string) (*GitOpsDiffResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/gitops/applications/%s/diff", uuid)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(GitOpsDiffResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// GetHistory returns paginated sync history for an application.
// GET /api/v1/gitops/applications/:uuid/history?page=&limit=
func (s *GitOpsService) GetHistory(ctx context.Context, uuid string, opts *GitOpsListOptions) (*GitOpsSyncHistoryResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/gitops/applications/%s/history", uuid)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(GitOpsSyncHistoryResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}
