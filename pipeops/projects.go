package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// ProjectService handles communication with the project related
// methods of the PipeOps API.
type ProjectService struct {
	client *Client
}

// Project represents a PipeOps project.
type Project struct {
	ID            string     `json:"id,omitempty"`
	UUID          string     `json:"uuid,omitempty"`
	Name          string     `json:"name,omitempty"`
	Description   string     `json:"description,omitempty"`
	Status        string     `json:"status,omitempty"`
	ServerID      string     `json:"server_id,omitempty"`
	EnvironmentID string     `json:"environment_id,omitempty"`
	WorkspaceID   string     `json:"workspace_id,omitempty"`
	Repository    string     `json:"repository,omitempty"`
	Branch        string     `json:"branch,omitempty"`
	BuildCommand  string     `json:"build_command,omitempty"`
	StartCommand  string     `json:"start_command,omitempty"`
	Port          int        `json:"port,omitempty"`
	Framework     string     `json:"framework,omitempty"`
	CreatedAt     *Timestamp `json:"created_at,omitempty"`
	UpdatedAt     *Timestamp `json:"updated_at,omitempty"`
}

// ProjectsResponse represents a list of projects response.
type ProjectsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Projects []Project `json:"projects"`
	} `json:"data"`
}

// ProjectResponse represents a single project response.
type ProjectResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Project Project `json:"project"`
	} `json:"data"`
}

// ProjectListOptions specifies the optional parameters to the
// ProjectService.List method.
type ProjectListOptions struct {
	WorkspaceID string `url:"workspace_id,omitempty"`
	ServerID    string `url:"server_id,omitempty"`
	Page        int    `url:"page,omitempty"`
	Limit       int    `url:"limit,omitempty"`
}

// List lists all projects.
func (s *ProjectService) List(ctx context.Context, opts *ProjectListOptions) (*ProjectsResponse, *http.Response, error) {
	u := "project"
	if opts != nil {
		u, _ = addOptions(u, opts)
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectsResp := new(ProjectsResponse)
	resp, err := s.client.Do(ctx, req, projectsResp)
	if err != nil {
		return nil, resp, err
	}

	return projectsResp, resp, nil
}

// Get fetches a project by UUID.
func (s *ProjectService) Get(ctx context.Context, projectUUID string) (*ProjectResponse, *http.Response, error) {
	u := fmt.Sprintf("project/%s", projectUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectResp := new(ProjectResponse)
	resp, err := s.client.Do(ctx, req, projectResp)
	if err != nil {
		return nil, resp, err
	}

	return projectResp, resp, nil
}

// CreateProjectRequest represents a request to create a project.
type CreateProjectRequest struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description,omitempty"`
	ServerID      string                 `json:"server_id"`
	EnvironmentID string                 `json:"environment_id"`
	Repository    string                 `json:"repository"`
	Branch        string                 `json:"branch"`
	BuildCommand  string                 `json:"build_command,omitempty"`
	StartCommand  string                 `json:"start_command,omitempty"`
	Port          int                    `json:"port,omitempty"`
	Framework     string                 `json:"framework,omitempty"`
	EnvVars       map[string]interface{} `json:"env_vars,omitempty"`
}

// Create creates a new project.
func (s *ProjectService) Create(ctx context.Context, req *CreateProjectRequest) (*ProjectResponse, *http.Response, error) {
	u := "project"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	projectResp := new(ProjectResponse)
	resp, err := s.client.Do(ctx, httpReq, projectResp)
	if err != nil {
		return nil, resp, err
	}

	return projectResp, resp, nil
}

// UpdateProjectRequest represents a request to update a project.
type UpdateProjectRequest struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	BuildCommand string `json:"build_command,omitempty"`
	StartCommand string `json:"start_command,omitempty"`
	Port         int    `json:"port,omitempty"`
}

// Update updates a project.
func (s *ProjectService) Update(ctx context.Context, projectUUID string, req *UpdateProjectRequest) (*ProjectResponse, *http.Response, error) {
	u := fmt.Sprintf("project/%s", projectUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	projectResp := new(ProjectResponse)
	resp, err := s.client.Do(ctx, httpReq, projectResp)
	if err != nil {
		return nil, resp, err
	}

	return projectResp, resp, nil
}

// Delete deletes a project.
func (s *ProjectService) Delete(ctx context.Context, projectUUID string) (*http.Response, error) {
	u := fmt.Sprintf("project/%s", projectUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// LogsOptions specifies options for retrieving project logs.
type LogsOptions struct {
	StartTime string `url:"start_time,omitempty"`
	EndTime   string `url:"end_time,omitempty"`
	Limit     int    `url:"limit,omitempty"`
	Search    string `url:"search,omitempty"`
}

// LogsResponse represents project logs response.
type LogsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Logs []map[string]interface{} `json:"logs"`
	} `json:"data"`
}

// GetLogs retrieves logs for a project.
func (s *ProjectService) GetLogs(ctx context.Context, projectUUID string, opts *LogsOptions) (*LogsResponse, *http.Response, error) {
	u := fmt.Sprintf("project/logs/%s", projectUUID)
	if opts != nil {
		u, _ = addOptions(u, opts)
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	logsResp := new(LogsResponse)
	resp, err := s.client.Do(ctx, req, logsResp)
	if err != nil {
		return nil, resp, err
	}

	return logsResp, resp, nil
}

// GitHubBranchesRequest represents a request to fetch GitHub branches.
type GitHubBranchesRequest struct {
	Repository string `json:"repository"`
}

// GitHubBranchesResponse represents GitHub branches response.
type GitHubBranchesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Branches []string `json:"branches"`
	} `json:"data"`
}

// GetGitHubBranches fetches branches from a GitHub repository.
func (s *ProjectService) GetGitHubBranches(ctx context.Context, req *GitHubBranchesRequest) (*GitHubBranchesResponse, *http.Response, error) {
	u := "project/github/branches"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	branchesResp := new(GitHubBranchesResponse)
	resp, err := s.client.Do(ctx, httpReq, branchesResp)
	if err != nil {
		return nil, resp, err
	}

	return branchesResp, resp, nil
}

// DomainRequest represents a request to add/update a project domain.
type DomainRequest struct {
	Domain string `json:"domain"`
}

// DomainResponse represents domain response.
type DomainResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Domain string `json:"domain"`
	} `json:"data"`
}

// UpdateDomain updates the domain for a project.
func (s *ProjectService) UpdateDomain(ctx context.Context, projectUUID string, req *DomainRequest) (*DomainResponse, *http.Response, error) {
	u := fmt.Sprintf("project/%s/domain", projectUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	domainResp := new(DomainResponse)
	resp, err := s.client.Do(ctx, httpReq, domainResp)
	if err != nil {
		return nil, resp, err
	}

	return domainResp, resp, nil
}

// EnvVariablesRequest represents a request to update environment variables.
type EnvVariablesRequest struct {
	EnvVariables []EnvVariable `json:"env_variables"`
}

// EnvVariablesResponse represents environment variables response.
type EnvVariablesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		EnvVariables []EnvVariable `json:"env_variables"`
	} `json:"data"`
}

// UpdateEnvVariables updates environment variables for a project.
func (s *ProjectService) UpdateEnvVariables(ctx context.Context, projectUUID string, req *EnvVariablesRequest) (*EnvVariablesResponse, *http.Response, error) {
	u := fmt.Sprintf("project/settings/env/%s", projectUUID)

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	envResp := new(EnvVariablesResponse)
	resp, err := s.client.Do(ctx, httpReq, envResp)
	if err != nil {
		return nil, resp, err
	}

	return envResp, resp, nil
}

// GetEnvVariables retrieves environment variables for a project.
func (s *ProjectService) GetEnvVariables(ctx context.Context, projectUUID string) (*EnvVariablesResponse, *http.Response, error) {
	u := fmt.Sprintf("project/settings/env/%s", projectUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	envResp := new(EnvVariablesResponse)
	resp, err := s.client.Do(ctx, req, envResp)
	if err != nil {
		return nil, resp, err
	}

	return envResp, resp, nil
}

// Deploy triggers a deployment for a project.
func (s *ProjectService) Deploy(ctx context.Context, projectUUID string) (*http.Response, error) {
	u := fmt.Sprintf("project/%s/deploy", projectUUID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// Restart restarts a project.
func (s *ProjectService) Restart(ctx context.Context, projectUUID string) (*http.Response, error) {
	u := fmt.Sprintf("project/%s/restart", projectUUID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// Stop stops a project.
func (s *ProjectService) Stop(ctx context.Context, projectUUID string) (*http.Response, error) {
	u := fmt.Sprintf("project/%s/stop", projectUUID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// MetricsRequest represents a metrics request.
type MetricsRequest struct {
	ProjectUUID string `json:"project_uuid"`
	StartTime   string `json:"start_time,omitempty"`
	EndTime     string `json:"end_time,omitempty"`
}

// MetricsResponse represents metrics response.
type MetricsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Metrics map[string]interface{} `json:"metrics"`
	} `json:"data"`
}

// GetMetrics retrieves metrics for a project.
func (s *ProjectService) GetMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, *http.Response, error) {
	u := "observability/project/summary"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	metricsResp := new(MetricsResponse)
	resp, err := s.client.Do(ctx, httpReq, metricsResp)
	if err != nil {
		return nil, resp, err
	}

	return metricsResp, resp, nil
}

// BulkDeleteRequest represents a request to delete multiple projects.
type BulkDeleteRequest struct {
	ProjectUUIDs []string `json:"project_uuids"`
}

// BulkDelete deletes multiple projects.
func (s *ProjectService) BulkDelete(ctx context.Context, req *BulkDeleteRequest) (*http.Response, error) {
	u := "project/delete"

	httpReq, err := s.client.NewRequest(http.MethodDelete, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// CostsResponse represents project costs response.
type CostsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Costs map[string]interface{} `json:"costs"`
	} `json:"data"`
}

// GetCosts retrieves costs for a project.
func (s *ProjectService) GetCosts(ctx context.Context, projectUUID string) (*CostsResponse, *http.Response, error) {
	u := fmt.Sprintf("project/costs/%s/billing", projectUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	costsResp := new(CostsResponse)
	resp, err := s.client.Do(ctx, req, costsResp)
	if err != nil {
		return nil, resp, err
	}

	return costsResp, resp, nil
}
