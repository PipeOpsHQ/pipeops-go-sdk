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
		var err error
		u, err = addOptions(u, opts)
		if err != nil {
			return nil, nil, err
		}
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
		var err error
		u, err = addOptions(u, opts)
		if err != nil {
			return nil, nil, err
		}
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

// TailLogs tails logs for a project (streams recent logs).
func (s *ProjectService) TailLogs(ctx context.Context, projectUUID string, opts *LogsOptions) (*LogsResponse, *http.Response, error) {
	u := fmt.Sprintf("project/logs/%s", projectUUID)
	if opts != nil {
		var err error
		u, err = addOptions(u, opts)
		if err != nil {
			return nil, nil, err
		}
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

// SearchLogs searches logs for a project.
func (s *ProjectService) SearchLogs(ctx context.Context, projectUUID string, opts *LogsOptions) (*LogsResponse, *http.Response, error) {
	u := fmt.Sprintf("project/logs/%s", projectUUID)
	if opts != nil {
		var err error
		u, err = addOptions(u, opts)
		if err != nil {
			return nil, nil, err
		}
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

// Observability and Metrics

// CPUMetricsRequest represents CPU metrics request.
type CPUMetricsRequest struct {
	ProjectUUID string `json:"project_uuid"`
	StartTime   string `json:"start_time,omitempty"`
	EndTime     string `json:"end_time,omitempty"`
}

// GetCPUMetrics retrieves CPU metrics for a project.
func (s *ProjectService) GetCPUMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, *http.Response, error) {
	u := "observability/app/cpu"

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

// GetStorageMetrics retrieves storage metrics for a project.
func (s *ProjectService) GetStorageMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, *http.Response, error) {
	u := "observability/app/storage"

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

// GetMemoryMetrics retrieves memory metrics for a project.
func (s *ProjectService) GetMemoryMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, *http.Response, error) {
	u := "observability/app/memory"

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

// GetNetworkIOMetrics retrieves network I/O metrics for a project.
func (s *ProjectService) GetNetworkIOMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, *http.Response, error) {
	u := "observability/app/network-io"

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

// GetControlPlaneMetrics retrieves control plane metrics.
func (s *ProjectService) GetControlPlaneMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, *http.Response, error) {
	u := "observability/control-plane"

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

// GetMetricsOverview retrieves metrics overview for a project.
func (s *ProjectService) GetMetricsOverview(ctx context.Context, req *MetricsRequest) (*MetricsResponse, *http.Response, error) {
	u := "observability/app/overview"

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

// Network Policies

// NetworkPolicy represents a network policy.
type NetworkPolicy struct {
	ID          string     `json:"id,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Rules       []string   `json:"rules,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
}

// NetworkPolicyRequest represents a network policy request.
type NetworkPolicyRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Rules       []string `json:"rules,omitempty"`
}

// NetworkPolicyResponse represents a network policy response.
type NetworkPolicyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Policy NetworkPolicy `json:"policy"`
	} `json:"data"`
}

// NetworkPoliciesResponse represents network policies response.
type NetworkPoliciesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Policies []NetworkPolicy `json:"policies"`
	} `json:"data"`
}

// CreateNetworkPolicy creates a network policy for a project.
func (s *ProjectService) CreateNetworkPolicy(ctx context.Context, projectUUID string, req *NetworkPolicyRequest) (*NetworkPolicyResponse, *http.Response, error) {
	u := fmt.Sprintf("project/settings/%s/network-policy", projectUUID)

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	policyResp := new(NetworkPolicyResponse)
	resp, err := s.client.Do(ctx, httpReq, policyResp)
	if err != nil {
		return nil, resp, err
	}

	return policyResp, resp, nil
}

// UpdateNetworkPolicy updates a network policy.
func (s *ProjectService) UpdateNetworkPolicy(ctx context.Context, projectUUID, policyUUID string, req *NetworkPolicyRequest) (*NetworkPolicyResponse, *http.Response, error) {
	u := fmt.Sprintf("project/settings/%s/network-policy/%s", projectUUID, policyUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	policyResp := new(NetworkPolicyResponse)
	resp, err := s.client.Do(ctx, httpReq, policyResp)
	if err != nil {
		return nil, resp, err
	}

	return policyResp, resp, nil
}

// ListNetworkPolicies lists network policies for a project.
func (s *ProjectService) ListNetworkPolicies(ctx context.Context, projectUUID string) (*NetworkPoliciesResponse, *http.Response, error) {
	u := fmt.Sprintf("project/settings/%s/network-policy", projectUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	policiesResp := new(NetworkPoliciesResponse)
	resp, err := s.client.Do(ctx, req, policiesResp)
	if err != nil {
		return nil, resp, err
	}

	return policiesResp, resp, nil
}

// Network Settings

// NetworkSettingsRequest represents network settings update request.
type NetworkSettingsRequest struct {
	Port int `json:"port"`
}

// NetworkSettingsResponse represents network settings response.
type NetworkSettingsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Settings map[string]interface{} `json:"settings"`
	} `json:"data"`
}

// UpdateNetworkingPort updates the networking port for a project.
func (s *ProjectService) UpdateNetworkingPort(ctx context.Context, projectUUID string, req *NetworkSettingsRequest) (*NetworkSettingsResponse, *http.Response, error) {
	u := fmt.Sprintf("project/settings/network/%s", projectUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	settingsResp := new(NetworkSettingsResponse)
	resp, err := s.client.Do(ctx, httpReq, settingsResp)
	if err != nil {
		return nil, resp, err
	}

	return settingsResp, resp, nil
}

// GenerateDomainFromNetworkPort generates a domain from network port.
func (s *ProjectService) GenerateDomainFromNetworkPort(ctx context.Context, projectUUID string) (*DomainResponse, *http.Response, error) {
	u := fmt.Sprintf("project/settings/network-name/%s", projectUUID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, nil, err
	}

	domainResp := new(DomainResponse)
	resp, err := s.client.Do(ctx, req, domainResp)
	if err != nil {
		return nil, resp, err
	}

	return domainResp, resp, nil
}

// GetNetworkSettings retrieves network settings for a project.
func (s *ProjectService) GetNetworkSettings(ctx context.Context, projectUUID string) (*NetworkSettingsResponse, *http.Response, error) {
	u := fmt.Sprintf("project/settings/network/%s", projectUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	settingsResp := new(NetworkSettingsResponse)
	resp, err := s.client.Do(ctx, req, settingsResp)
	if err != nil {
		return nil, resp, err
	}

	return settingsResp, resp, nil
}

// GitHub/GitLab Integration

// GitHubOrgsResponse represents GitHub organizations response.
type GitHubOrgsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Organizations []map[string]interface{} `json:"organizations"`
	} `json:"data"`
}

// GetGitHubOrgs retrieves GitHub organizations.
func (s *ProjectService) GetGitHubOrgs(ctx context.Context) (*GitHubOrgsResponse, *http.Response, error) {
	u := "project/github/organisations"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	orgsResp := new(GitHubOrgsResponse)
	resp, err := s.client.Do(ctx, req, orgsResp)
	if err != nil {
		return nil, resp, err
	}

	return orgsResp, resp, nil
}

// GitLabOrgReposRequest represents GitLab org repos request.
type GitLabOrgReposRequest struct {
	OrgID string `json:"org_id"`
}

// GitLabReposResponse represents GitLab repos response.
type GitLabReposResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Repos []map[string]interface{} `json:"repos"`
	} `json:"data"`
}

// GetGitLabOrgRepos retrieves GitLab organization repos.
func (s *ProjectService) GetGitLabOrgRepos(ctx context.Context, req *GitLabOrgReposRequest) (*GitLabReposResponse, *http.Response, error) {
	u := "project/gitlab/organisations/repos"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	reposResp := new(GitLabReposResponse)
	resp, err := s.client.Do(ctx, httpReq, reposResp)
	if err != nil {
		return nil, resp, err
	}

	return reposResp, resp, nil
}

// MigrateProject migrates a project to different server/workspace.
func (s *ProjectService) MigrateProject(ctx context.Context, projectUUID, serverUUID, workspaceUUID string) (*http.Response, error) {
	u := fmt.Sprintf("project/migrate/%s/server/%s/workspace/%s", projectUUID, serverUUID, workspaceUUID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// RuntimeLogsResponse represents runtime logs response.
type RuntimeLogsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Logs []string `json:"logs"`
	} `json:"data"`
}

// GetRuntimeLogs retrieves runtime logs for a project pod.
func (s *ProjectService) GetRuntimeLogs(ctx context.Context, projectUUID, podName string) (*RuntimeLogsResponse, *http.Response, error) {
	u := fmt.Sprintf("project/runtime-logs/%s/%s", projectUUID, podName)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	logsResp := new(RuntimeLogsResponse)
	resp, err := s.client.Do(ctx, req, logsResp)
	if err != nil {
		return nil, resp, err
	}

	return logsResp, resp, nil
}

// PodsResponse represents pods response.
type PodsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Pods []map[string]interface{} `json:"pods"`
	} `json:"data"`
}

// GetPodsFromLabel retrieves pods from label for a project.
func (s *ProjectService) GetPodsFromLabel(ctx context.Context, projectUUID string) (*PodsResponse, *http.Response, error) {
	u := fmt.Sprintf("project/pod-label/%s", projectUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	podsResp := new(PodsResponse)
	resp, err := s.client.Do(ctx, req, podsResp)
	if err != nil {
		return nil, resp, err
	}

	return podsResp, resp, nil
}

// CheckDockerfileRequest represents dockerfile check request.
type CheckDockerfileRequest struct {
	Provider   string `json:"provider"`
	Workspace  string `json:"workspace"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
}

// CheckDockerfileResponse represents dockerfile check response.
type CheckDockerfileResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Exists bool `json:"exists"`
	} `json:"data"`
}

// CheckDockerfile checks if Dockerfile exists in repository.
func (s *ProjectService) CheckDockerfile(ctx context.Context, provider, workspace, repo, branch string) (*CheckDockerfileResponse, *http.Response, error) {
	u := fmt.Sprintf("project/check-dockerfile/%s/%s/%s/%s", provider, workspace, repo, branch)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	checkResp := new(CheckDockerfileResponse)
	resp, err := s.client.Do(ctx, req, checkResp)
	if err != nil {
		return nil, resp, err
	}

	return checkResp, resp, nil
}

// LinkProvider initiates linking a Git provider.
func (s *ProjectService) LinkProvider(ctx context.Context, provider string) (*http.Response, error) {
	u := fmt.Sprintf("project/link/%s", provider)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// LinkProviderCallback handles provider link callback.
func (s *ProjectService) LinkProviderCallback(ctx context.Context, provider, uuid string) (*http.Response, error) {
	u := fmt.Sprintf("project/link/%s/callback/%s", provider, uuid)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// JobEventResponse represents job event response.
type JobEventResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Event map[string]interface{} `json:"event"`
	} `json:"data"`
}

// GetJobEvent retrieves job event for a project.
func (s *ProjectService) GetJobEvent(ctx context.Context, projectUUID, internalProjectName string) (*JobEventResponse, *http.Response, error) {
	u := fmt.Sprintf("project/job/event/%s/%s", projectUUID, internalProjectName)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	eventResp := new(JobEventResponse)
	resp, err := s.client.Do(ctx, req, eventResp)
	if err != nil {
		return nil, resp, err
	}

	return eventResp, resp, nil
}

// ValidatePort validates if a port is available.
func (s *ProjectService) ValidatePort(ctx context.Context, environment, port string) (*http.Response, error) {
	u := fmt.Sprintf("project/port-validator/%s/%s", environment, port)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// CheckDomainSSLRequest represents domain SSL check request.
type CheckDomainSSLRequest struct {
	Domain string `json:"domain"`
}

// CheckDomainSSL checks domain SSL configuration.
func (s *ProjectService) CheckDomainSSL(ctx context.Context, req *CheckDomainSSLRequest) (*http.Response, error) {
	u := "project/domain/check-ssl"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// SetProjectDomainName sets the project domain name.
func (s *ProjectService) SetProjectDomainName(ctx context.Context, projectUUID string, req *DomainRequest) (*http.Response, error) {
	u := fmt.Sprintf("project/settings/name/%s", projectUUID)

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// DeleteCustomDomain deletes a custom domain from a project.
func (s *ProjectService) DeleteCustomDomain(ctx context.Context, projectUUID string) (*http.Response, error) {
	u := fmt.Sprintf("project/%s/custom-domain", projectUUID)

	req, err := s.client.NewRequest(http.MethodPatch, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// RepoSearchRequest represents repository search request.
type RepoSearchRequest struct {
	Query string `json:"query"`
}

// RepoSearchResponse represents repository search response.
type RepoSearchResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Repos []map[string]interface{} `json:"repos"`
	} `json:"data"`
}

// SearchRepos searches for repositories.
func (s *ProjectService) SearchRepos(ctx context.Context, req *RepoSearchRequest) (*RepoSearchResponse, *http.Response, error) {
	u := "project/github/repo-search"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	searchResp := new(RepoSearchResponse)
	resp, err := s.client.Do(ctx, httpReq, searchResp)
	if err != nil {
		return nil, resp, err
	}

	return searchResp, resp, nil
}

// ProjectNamesResponse represents project names response.
type ProjectNamesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Names []string `json:"names"`
	} `json:"data"`
}

// GetProjectNames retrieves user's project names.
func (s *ProjectService) GetProjectNames(ctx context.Context) (*ProjectNamesResponse, *http.Response, error) {
	u := "project/fetch-names"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	namesResp := new(ProjectNamesResponse)
	resp, err := s.client.Do(ctx, req, namesResp)
	if err != nil {
		return nil, resp, err
	}

	return namesResp, resp, nil
}

// CheckProjectName checks if a project name is available.
func (s *ProjectService) CheckProjectName(ctx context.Context) (*http.Response, error) {
	u := "project/check-project-name"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}
