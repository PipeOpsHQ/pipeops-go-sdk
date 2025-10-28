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
