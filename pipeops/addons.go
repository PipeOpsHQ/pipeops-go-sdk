package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// AddOnService handles communication with the add-on related
// methods of the PipeOps API.
type AddOnService struct {
	client *Client
}

// AddOn represents a PipeOps add-on.
type AddOn struct {
	ID          string     `json:"id,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Category    string     `json:"category,omitempty"`
	Version     string     `json:"version,omitempty"`
	Icon        string     `json:"icon,omitempty"`
	Status      string     `json:"status,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`
}

// AddOnsResponse represents a list of add-ons response.
type AddOnsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AddOns []AddOn `json:"addons"`
	} `json:"data"`
}

// AddOnResponse represents a single add-on response.
type AddOnResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AddOn AddOn `json:"addon"`
	} `json:"data"`
}

// List lists all available add-ons.
func (s *AddOnService) List(ctx context.Context) (*AddOnsResponse, *http.Response, error) {
	u := "addons"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	addOnsResp := new(AddOnsResponse)
	resp, err := s.client.Do(ctx, req, addOnsResp)
	if err != nil {
		return nil, resp, err
	}

	return addOnsResp, resp, nil
}

// Get fetches an add-on by UUID.
func (s *AddOnService) Get(ctx context.Context, addonUUID string) (*AddOnResponse, *http.Response, error) {
	u := fmt.Sprintf("addons/%s", addonUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	addOnResp := new(AddOnResponse)
	resp, err := s.client.Do(ctx, req, addOnResp)
	if err != nil {
		return nil, resp, err
	}

	return addOnResp, resp, nil
}

// AddOnDeployment represents a deployed add-on instance.
type AddOnDeployment struct {
	ID        string                 `json:"id,omitempty"`
	UUID      string                 `json:"uuid,omitempty"`
	AddOnID   string                 `json:"addon_id,omitempty"`
	AddOnName string                 `json:"addon_name,omitempty"`
	ProjectID string                 `json:"project_id,omitempty"`
	ServerID  string                 `json:"server_id,omitempty"`
	Status    string                 `json:"status,omitempty"`
	Config    map[string]interface{} `json:"config,omitempty"`
	CreatedAt *Timestamp             `json:"created_at,omitempty"`
	UpdatedAt *Timestamp             `json:"updated_at,omitempty"`
}

// AddOnDeploymentsResponse represents a list of add-on deployments response.
type AddOnDeploymentsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Deployments []AddOnDeployment `json:"deployments"`
	} `json:"data"`
}

// AddOnDeploymentResponse represents a single add-on deployment response.
type AddOnDeploymentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Deployment AddOnDeployment `json:"deployment"`
	} `json:"data"`
}

// DeployAddOnRequest represents a request to deploy an add-on.
type DeployAddOnRequest struct {
	AddOnUUID string                 `json:"addon_uuid"`
	ProjectID string                 `json:"project_id,omitempty"`
	ServerID  string                 `json:"server_id,omitempty"`
	Config    map[string]interface{} `json:"config,omitempty"`
}

// Deploy deploys an add-on.
func (s *AddOnService) Deploy(ctx context.Context, req *DeployAddOnRequest) (*AddOnDeploymentResponse, *http.Response, error) {
	u := "addons/deploy"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	deployResp := new(AddOnDeploymentResponse)
	resp, err := s.client.Do(ctx, httpReq, deployResp)
	if err != nil {
		return nil, resp, err
	}

	return deployResp, resp, nil
}

// ListDeployments lists all add-on deployments.
func (s *AddOnService) ListDeployments(ctx context.Context) (*AddOnDeploymentsResponse, *http.Response, error) {
	u := "addons/deployments"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	deploymentsResp := new(AddOnDeploymentsResponse)
	resp, err := s.client.Do(ctx, req, deploymentsResp)
	if err != nil {
		return nil, resp, err
	}

	return deploymentsResp, resp, nil
}

// GetDeployment fetches an add-on deployment by UUID.
func (s *AddOnService) GetDeployment(ctx context.Context, deploymentUUID string) (*AddOnDeploymentResponse, *http.Response, error) {
	u := fmt.Sprintf("addons/deployments/%s", deploymentUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	deployResp := new(AddOnDeploymentResponse)
	resp, err := s.client.Do(ctx, req, deployResp)
	if err != nil {
		return nil, resp, err
	}

	return deployResp, resp, nil
}

// DeleteDeployment deletes an add-on deployment.
func (s *AddOnService) DeleteDeployment(ctx context.Context, deploymentUUID string) (*http.Response, error) {
	u := fmt.Sprintf("addons/deployments/%s", deploymentUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// AddOnCategory represents an add-on category.
type AddOnCategory struct {
	ID          string `json:"id,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

// AddOnCategoriesResponse represents a list of add-on categories response.
type AddOnCategoriesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Categories []AddOnCategory `json:"categories"`
	} `json:"data"`
}

// ListCategories lists all add-on categories.
func (s *AddOnService) ListCategories(ctx context.Context) (*AddOnCategoriesResponse, *http.Response, error) {
	u := "addons/categories"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	categoriesResp := new(AddOnCategoriesResponse)
	resp, err := s.client.Do(ctx, req, categoriesResp)
	if err != nil {
		return nil, resp, err
	}

	return categoriesResp, resp, nil
}
