package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// EnvironmentService handles communication with the environment related
// methods of the PipeOps API.
type EnvironmentService struct {
	client *Client
}

// Environment represents a PipeOps environment.
type Environment struct {
	ID          string     `json:"id,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	Name        string     `json:"name,omitempty"`
	WorkspaceID string     `json:"workspace_id,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`
}

// EnvironmentsResponse represents a list of environments response.
type EnvironmentsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Environments []Environment `json:"environments"`
	} `json:"data"`
}

// EnvironmentResponse represents a single environment response.
type EnvironmentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Environment Environment `json:"environment"`
	} `json:"data"`
}

// List lists all environments.
func (s *EnvironmentService) List(ctx context.Context) (*EnvironmentsResponse, *http.Response, error) {
	u := "environment"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	envsResp := new(EnvironmentsResponse)
	resp, err := s.client.Do(ctx, req, envsResp)
	if err != nil {
		return nil, resp, err
	}

	return envsResp, resp, nil
}

// Get fetches an environment by UUID.
func (s *EnvironmentService) Get(ctx context.Context, envUUID string) (*EnvironmentResponse, *http.Response, error) {
	u := fmt.Sprintf("environment/%s", envUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	envResp := new(EnvironmentResponse)
	resp, err := s.client.Do(ctx, req, envResp)
	if err != nil {
		return nil, resp, err
	}

	return envResp, resp, nil
}

// CreateEnvironmentRequest represents a request to create an environment.
type CreateEnvironmentRequest struct {
	Name         string        `json:"name"`
	WorkspaceID  string        `json:"workspace_id"`
	EnvVariables []EnvVariable `json:"env_variables,omitempty"`
}

// EnvVariable represents an environment variable.
type EnvVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Create creates a new environment.
func (s *EnvironmentService) Create(ctx context.Context, req *CreateEnvironmentRequest) (*EnvironmentResponse, *http.Response, error) {
	u := "environment/create"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	envResp := new(EnvironmentResponse)
	resp, err := s.client.Do(ctx, httpReq, envResp)
	if err != nil {
		return nil, resp, err
	}

	return envResp, resp, nil
}

// UpdateEnvironmentRequest represents a request to update an environment.
type UpdateEnvironmentRequest struct {
	Name string `json:"name,omitempty"`
}

// Update updates an environment.
func (s *EnvironmentService) Update(ctx context.Context, envUUID string, req *UpdateEnvironmentRequest) (*EnvironmentResponse, *http.Response, error) {
	u := fmt.Sprintf("environment/%s/update", envUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	envResp := new(EnvironmentResponse)
	resp, err := s.client.Do(ctx, httpReq, envResp)
	if err != nil {
		return nil, resp, err
	}

	return envResp, resp, nil
}

// Delete deletes an environment.
func (s *EnvironmentService) Delete(ctx context.Context, envUUID string) (*http.Response, error) {
	u := fmt.Sprintf("environment/%s", envUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// SetEnvironmentVariablesRequest represents a request to set environment variables.
type SetEnvironmentVariablesRequest struct {
	EnvVariables []EnvVariable `json:"env_variables"`
}

// SetEnvVariables sets environment variables for an environment.
func (s *EnvironmentService) SetEnvVariables(ctx context.Context, envUUID string, req *SetEnvironmentVariablesRequest) (*http.Response, error) {
	u := fmt.Sprintf("environment/%s/set-environment-env", envUUID)

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}
