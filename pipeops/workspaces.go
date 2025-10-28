package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// WorkspaceService handles communication with the workspace related
// methods of the PipeOps API.
type WorkspaceService struct {
	client *Client
}

// Workspace represents a PipeOps workspace.
type Workspace struct {
	ID          string     `json:"id,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	OwnerID     string     `json:"owner_id,omitempty"`
	TeamID      string     `json:"team_id,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`
}

// WorkspacesResponse represents a list of workspaces response.
type WorkspacesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Workspaces []Workspace `json:"workspaces"`
	} `json:"data"`
}

// WorkspaceResponse represents a single workspace response.
type WorkspaceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Workspace Workspace `json:"workspace"`
	} `json:"data"`
}

// CreateWorkspaceRequest represents a request to create a workspace.
type CreateWorkspaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	TeamID      string `json:"team_id,omitempty"`
}

// Create creates a new workspace.
func (s *WorkspaceService) Create(ctx context.Context, req *CreateWorkspaceRequest) (*WorkspaceResponse, *http.Response, error) {
	u := "workspace"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	workspaceResp := new(WorkspaceResponse)
	resp, err := s.client.Do(ctx, httpReq, workspaceResp)
	if err != nil {
		return nil, resp, err
	}

	return workspaceResp, resp, nil
}

// List lists all workspaces for the authenticated user.
func (s *WorkspaceService) List(ctx context.Context) (*WorkspacesResponse, *http.Response, error) {
	u := "workspace"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	workspacesResp := new(WorkspacesResponse)
	resp, err := s.client.Do(ctx, req, workspacesResp)
	if err != nil {
		return nil, resp, err
	}

	return workspacesResp, resp, nil
}

// Get fetches a workspace by UUID.
func (s *WorkspaceService) Get(ctx context.Context, workspaceUUID string) (*WorkspaceResponse, *http.Response, error) {
	u := fmt.Sprintf("workspace/fetch/%s", workspaceUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	workspaceResp := new(WorkspaceResponse)
	resp, err := s.client.Do(ctx, req, workspaceResp)
	if err != nil {
		return nil, resp, err
	}

	return workspaceResp, resp, nil
}
