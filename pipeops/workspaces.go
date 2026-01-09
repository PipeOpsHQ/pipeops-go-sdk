package pipeops

import (
	"context"
	"encoding/json"
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

func (w *Workspace) UnmarshalJSON(data []byte) error {
	type workspaceWire struct {
		ID jsonID `json:"id,omitempty"`

		UUID        string `json:"uuid,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		OwnerID     string `json:"owner_id,omitempty"`
		TeamID      string `json:"team_id,omitempty"`

		CreatedAt    *Timestamp `json:"created_at,omitempty"`
		UpdatedAt    *Timestamp `json:"updated_at,omitempty"`
		CreatedAtAlt *Timestamp `json:"CreatedAt,omitempty"`
		UpdatedAtAlt *Timestamp `json:"UpdatedAt,omitempty"`
	}

	var tmp workspaceWire
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	w.ID = tmp.ID.String()
	w.UUID = tmp.UUID
	w.Name = tmp.Name
	w.Description = tmp.Description
	w.OwnerID = tmp.OwnerID
	w.TeamID = tmp.TeamID

	w.CreatedAt = tmp.CreatedAt
	if w.CreatedAt == nil {
		w.CreatedAt = tmp.CreatedAtAlt
	}

	w.UpdatedAt = tmp.UpdatedAt
	if w.UpdatedAt == nil {
		w.UpdatedAt = tmp.UpdatedAtAlt
	}

	return nil
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

	rawResp := new(workspaceListEnvelope)
	resp, err := s.client.Do(ctx, httpReq, rawResp)
	if err != nil {
		return nil, resp, err
	}

	workspaceResp := &WorkspaceResponse{
		Status:  coalesceNonEmpty(rawResp.Status, statusFromSuccess(rawResp.Success)),
		Message: rawResp.Message,
	}

	var workspace Workspace
	if len(rawResp.Data) != 0 {
		if err := json.Unmarshal(rawResp.Data, &workspace); err != nil {
			var wrapped struct {
				Workspace Workspace `json:"workspace,omitempty"`
			}
			if err := json.Unmarshal(rawResp.Data, &wrapped); err != nil {
				return nil, resp, err
			}
			workspace = wrapped.Workspace
		}
	}

	workspaceResp.Data.Workspace = workspace
	return workspaceResp, resp, nil
}

// List lists all workspaces for the authenticated user.
func (s *WorkspaceService) List(ctx context.Context) (*WorkspacesResponse, *http.Response, error) {
	u := "workspace"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	rawResp := new(workspaceListEnvelope)
	resp, err := s.client.Do(ctx, req, rawResp)
	if err != nil {
		return nil, resp, err
	}

	workspacesResp := &WorkspacesResponse{
		Status:  coalesceNonEmpty(rawResp.Status, statusFromSuccess(rawResp.Success)),
		Message: rawResp.Message,
	}

	if len(rawResp.Data) == 0 {
		return workspacesResp, resp, nil
	}

	var workspaces []Workspace
	if err := json.Unmarshal(rawResp.Data, &workspaces); err == nil {
		workspacesResp.Data.Workspaces = workspaces
		return workspacesResp, resp, nil
	}

	var wrapped struct {
		Workspaces []Workspace `json:"workspaces,omitempty"`
	}
	if err := json.Unmarshal(rawResp.Data, &wrapped); err != nil {
		return nil, resp, err
	}
	workspacesResp.Data.Workspaces = wrapped.Workspaces

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

// UpdateWorkspaceRequest represents a request to update a workspace.
type UpdateWorkspaceRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Update updates a workspace.
func (s *WorkspaceService) Update(ctx context.Context, workspaceUUID string, req *UpdateWorkspaceRequest) (*WorkspaceResponse, *http.Response, error) {
	u := fmt.Sprintf("workspace/%s", workspaceUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
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

// Delete deletes a workspace.
func (s *WorkspaceService) Delete(ctx context.Context, workspaceUUID string) (*http.Response, error) {
	u := fmt.Sprintf("workspace/%s", workspaceUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// SetBillingEmailRequest represents a request to set billing email.
type SetBillingEmailRequest struct {
	Email string `json:"email"`
}

// SetBillingEmail sets the billing email for a workspace.
func (s *WorkspaceService) SetBillingEmail(ctx context.Context, workspaceUUID string, req *SetBillingEmailRequest) (*http.Response, error) {
	u := fmt.Sprintf("workspace/%s/add-billing-email", workspaceUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}
