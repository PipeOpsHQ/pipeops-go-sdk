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

// AddOnSubmissionRequest represents an add-on submission request.
type AddOnSubmissionRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Version     string                 `json:"version"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// SubmitAddOn submits a new add-on for review.
func (s *AddOnService) SubmitAddOn(ctx context.Context, req *AddOnSubmissionRequest) (*AddOnResponse, *http.Response, error) {
	u := "addons/submit"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	addOnResp := new(AddOnResponse)
	resp, err := s.client.Do(ctx, httpReq, addOnResp)
	if err != nil {
		return nil, resp, err
	}

	return addOnResp, resp, nil
}

// MySubmissionsResponse represents user's add-on submissions response.
type MySubmissionsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Submissions []AddOn `json:"submissions"`
	} `json:"data"`
}

// GetMySubmissions retrieves user's add-on submissions.
func (s *AddOnService) GetMySubmissions(ctx context.Context) (*MySubmissionsResponse, *http.Response, error) {
	u := "addons/my-submissions"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	submissionsResp := new(MySubmissionsResponse)
	resp, err := s.client.Do(ctx, req, submissionsResp)
	if err != nil {
		return nil, resp, err
	}

	return submissionsResp, resp, nil
}

// UpdateDeploymentRequest represents a request to update an add-on deployment.
type UpdateDeploymentRequest struct {
	Config map[string]interface{} `json:"config,omitempty"`
	Status string                 `json:"status,omitempty"`
}

// UpdateDeployment updates an add-on deployment configuration.
func (s *AddOnService) UpdateDeployment(ctx context.Context, deploymentUUID string, req *UpdateDeploymentRequest) (*AddOnDeploymentResponse, *http.Response, error) {
	u := fmt.Sprintf("addons/deployments/%s", deploymentUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
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

// SyncDeployment syncs an add-on deployment.
func (s *AddOnService) SyncDeployment(ctx context.Context, deploymentUID string) (*http.Response, error) {
	u := fmt.Sprintf("addons/deployments/%s/sync", deploymentUID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// DeploymentOverviewResponse represents deployment overview response.
type DeploymentOverviewResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Overview map[string]interface{} `json:"overview"`
	} `json:"data"`
}

// GetDeploymentOverview retrieves deployment overview.
func (s *AddOnService) GetDeploymentOverview(ctx context.Context) (*DeploymentOverviewResponse, *http.Response, error) {
	u := "addons/deployments/overview"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	overviewResp := new(DeploymentOverviewResponse)
	resp, err := s.client.Do(ctx, req, overviewResp)
	if err != nil {
		return nil, resp, err
	}

	return overviewResp, resp, nil
}

// DeploymentSessionResponse represents deployment session response.
type DeploymentSessionResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Session map[string]interface{} `json:"session"`
	} `json:"data"`
}

// GetDeploymentSession retrieves deployment session information.
func (s *AddOnService) GetDeploymentSession(ctx context.Context, sessionID string) (*DeploymentSessionResponse, *http.Response, error) {
	u := fmt.Sprintf("addons/deployments/sessions/%s", sessionID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	sessionResp := new(DeploymentSessionResponse)
	resp, err := s.client.Do(ctx, req, sessionResp)
	if err != nil {
		return nil, resp, err
	}

	return sessionResp, resp, nil
}

// DeploymentConfigsResponse represents deployment configs response.
type DeploymentConfigsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Configs map[string]interface{} `json:"configs"`
	} `json:"data"`
}

// ViewDeploymentConfigs views deployment configurations.
func (s *AddOnService) ViewDeploymentConfigs(ctx context.Context, addonUUID string) (*DeploymentConfigsResponse, *http.Response, error) {
	u := fmt.Sprintf("addons/deployments/%s/view/configs", addonUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	configsResp := new(DeploymentConfigsResponse)
	resp, err := s.client.Do(ctx, req, configsResp)
	if err != nil {
		return nil, resp, err
	}

	return configsResp, resp, nil
}

// AddDomain adds a domain to an add-on.
func (s *AddOnService) AddDomain(ctx context.Context, addonUUID string, req *DomainRequest) (*http.Response, error) {
	u := fmt.Sprintf("addons/%s/domain", addonUUID)

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// BulkDeleteDeploymentsRequest represents a request to bulk delete deployments.
type BulkDeleteDeploymentsRequest struct {
	DeploymentUIDs []string `json:"deployment_uids"`
}

// BulkDeleteDeployments deletes multiple add-on deployments.
func (s *AddOnService) BulkDeleteDeployments(ctx context.Context, req *BulkDeleteDeploymentsRequest) (*http.Response, error) {
	u := "addons/deployments/bulk"

	httpReq, err := s.client.NewRequest(http.MethodDelete, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}
