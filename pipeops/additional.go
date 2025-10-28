package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// DeploymentWebhookService handles deployment webhook related
// methods of the PipeOps API.
type DeploymentWebhookService struct {
	client *Client
}

// WebhookPayload represents a deployment webhook payload.
type WebhookPayload struct {
	Repository string                 `json:"repository,omitempty"`
	Branch     string                 `json:"branch,omitempty"`
	Commit     string                 `json:"commit,omitempty"`
	Author     string                 `json:"author,omitempty"`
	Message    string                 `json:"message,omitempty"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
}

// WebhookResponse represents webhook response.
type WebhookResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GitHubWebhook handles GitHub deployment webhook.
func (s *DeploymentWebhookService) GitHubWebhook(ctx context.Context, payload *WebhookPayload) (*WebhookResponse, *http.Response, error) {
	u := "deployment-webhook/github"

	req, err := s.client.NewRequest(http.MethodPost, u, payload)
	if err != nil {
		return nil, nil, err
	}

	webhookResp := new(WebhookResponse)
	resp, err := s.client.Do(ctx, req, webhookResp)
	if err != nil {
		return nil, resp, err
	}

	return webhookResp, resp, nil
}

// GitLabWebhook handles GitLab deployment webhook.
func (s *DeploymentWebhookService) GitLabWebhook(ctx context.Context, payload *WebhookPayload) (*WebhookResponse, *http.Response, error) {
	u := "deployment-webhook/gitlab"

	req, err := s.client.NewRequest(http.MethodPost, u, payload)
	if err != nil {
		return nil, nil, err
	}

	webhookResp := new(WebhookResponse)
	resp, err := s.client.Do(ctx, req, webhookResp)
	if err != nil {
		return nil, resp, err
	}

	return webhookResp, resp, nil
}

// BitbucketWebhook handles Bitbucket deployment webhook.
func (s *DeploymentWebhookService) BitbucketWebhook(ctx context.Context, payload *WebhookPayload) (*WebhookResponse, *http.Response, error) {
	u := "deployment-webhook/bitbucket"

	req, err := s.client.NewRequest(http.MethodPost, u, payload)
	if err != nil {
		return nil, nil, err
	}

	webhookResp := new(WebhookResponse)
	resp, err := s.client.Do(ctx, req, webhookResp)
	if err != nil {
		return nil, resp, err
	}

	return webhookResp, resp, nil
}

// CampaignService handles campaign related methods of the PipeOps API.
type CampaignService struct {
	client *Client
}

// Campaign represents a campaign.
type Campaign struct {
	ID          string     `json:"id,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	StartDate   *Timestamp `json:"start_date,omitempty"`
	EndDate     *Timestamp `json:"end_date,omitempty"`
	Status      string     `json:"status,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
}

// CampaignRequest represents a campaign request.
type CampaignRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
}

// CampaignResponse represents a campaign response.
type CampaignResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Campaign Campaign `json:"campaign"`
	} `json:"data"`
}

// CampaignsResponse represents campaigns response.
type CampaignsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Campaigns []Campaign `json:"campaigns"`
	} `json:"data"`
}

// Create creates a new campaign.
func (s *CampaignService) Create(ctx context.Context, req *CampaignRequest) (*CampaignResponse, *http.Response, error) {
	u := "campaign/create"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	campaignResp := new(CampaignResponse)
	resp, err := s.client.Do(ctx, httpReq, campaignResp)
	if err != nil {
		return nil, resp, err
	}

	return campaignResp, resp, nil
}

// List lists all campaigns.
func (s *CampaignService) List(ctx context.Context) (*CampaignsResponse, *http.Response, error) {
	u := "campaign"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	campaignsResp := new(CampaignsResponse)
	resp, err := s.client.Do(ctx, req, campaignsResp)
	if err != nil {
		return nil, resp, err
	}

	return campaignsResp, resp, nil
}

// Get gets a campaign by UUID.
func (s *CampaignService) Get(ctx context.Context, campaignUUID string) (*CampaignResponse, *http.Response, error) {
	u := "campaign/" + campaignUUID

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	campaignResp := new(CampaignResponse)
	resp, err := s.client.Do(ctx, req, campaignResp)
	if err != nil {
		return nil, resp, err
	}

	return campaignResp, resp, nil
}

// Update updates a campaign.
func (s *CampaignService) Update(ctx context.Context, campaignUUID string, req *CampaignRequest) (*CampaignResponse, *http.Response, error) {
	u := "campaign/" + campaignUUID

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	campaignResp := new(CampaignResponse)
	resp, err := s.client.Do(ctx, httpReq, campaignResp)
	if err != nil {
		return nil, resp, err
	}

	return campaignResp, resp, nil
}

// Delete deletes a campaign.
func (s *CampaignService) Delete(ctx context.Context, campaignUUID string) (*http.Response, error) {
	u := "campaign/" + campaignUUID

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// Start starts a campaign.
func (s *CampaignService) Start(ctx context.Context, campaignUUID string) (*http.Response, error) {
	u := "campaign/" + campaignUUID + "/start"

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// Stop stops a campaign.
func (s *CampaignService) Stop(ctx context.Context, campaignUUID string) (*http.Response, error) {
	u := "campaign/" + campaignUUID + "/stop"

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// CouponService handles coupon related methods of the PipeOps API.
type CouponService struct {
	client *Client
}

// Coupon represents a coupon.
type Coupon struct {
	ID        string     `json:"id,omitempty"`
	UUID      string     `json:"uuid,omitempty"`
	Code      string     `json:"code,omitempty"`
	Discount  float64    `json:"discount,omitempty"`
	Type      string     `json:"type,omitempty"`
	ExpiresAt *Timestamp `json:"expires_at,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
}

// CouponRequest represents a coupon request.
type CouponRequest struct {
	Code      string  `json:"code"`
	Discount  float64 `json:"discount"`
	Type      string  `json:"type,omitempty"`
	ExpiresAt string  `json:"expires_at,omitempty"`
}

// CouponResponse represents a coupon response.
type CouponResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Coupon Coupon `json:"coupon"`
	} `json:"data"`
}

// Create creates a new coupon for an agreement.
func (s *CouponService) Create(ctx context.Context, agreementUUID string, req *CouponRequest) (*CouponResponse, *http.Response, error) {
	u := "coupons/agreements/" + agreementUUID

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	couponResp := new(CouponResponse)
	resp, err := s.client.Do(ctx, httpReq, couponResp)
	if err != nil {
		return nil, resp, err
	}

	return couponResp, resp, nil
}

// Get gets a coupon by UUID and agreement.
func (s *CouponService) Get(ctx context.Context, couponUUID, agreementUUID string) (*CouponResponse, *http.Response, error) {
	u := "coupons/" + couponUUID + "/agreements/" + agreementUUID

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	couponResp := new(CouponResponse)
	resp, err := s.client.Do(ctx, req, couponResp)
	if err != nil {
		return nil, resp, err
	}

	return couponResp, resp, nil
}

// ServiceService handles service related methods of the PipeOps API.
type ServiceService struct {
	client *Client
}

// CreateDatabaseRequest represents a database creation request.
type CreateDatabaseRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Version  string `json:"version,omitempty"`
	ServerID string `json:"server_id"`
}

// CreateDatabase creates a new database service.
func (s *ServiceService) CreateDatabase(ctx context.Context, req *CreateDatabaseRequest) (*http.Response, error) {
	u := "service/create-database"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// MCPRegistryService handles MCP registry related methods.
type MCPRegistryService struct {
	client *Client
}

// MCPServersResponse represents MCP servers response.
type MCPServersResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Servers []map[string]interface{} `json:"servers"`
	} `json:"data"`
}

// GetMCPServers retrieves MCP registry servers.
func (s *MCPRegistryService) GetMCPServers(ctx context.Context) (*MCPServersResponse, *http.Response, error) {
	u := "servers"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	serversResp := new(MCPServersResponse)
	resp, err := s.client.Do(ctx, req, serversResp)
	if err != nil {
		return nil, resp, err
	}

	return serversResp, resp, nil
}

// OpenCostService handles open cost related methods.
type OpenCostService struct {
	client *Client
}

// ClusterCostResponse represents cluster cost response.
type ClusterCostResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Cost map[string]interface{} `json:"cost"`
	} `json:"data"`
}

// GetClusterComputeCost gets total cost for carpenter enabled server.
func (s *OpenCostService) GetClusterComputeCost(ctx context.Context, clusterUUID string) (*ClusterCostResponse, *http.Response, error) {
	u := fmt.Sprintf("cluster/%s/cost/allocation/compute", clusterUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	costResp := new(ClusterCostResponse)
	resp, err := s.client.Do(ctx, req, costResp)
	if err != nil {
		return nil, resp, err
	}

	return costResp, resp, nil
}

// GetProjectsCost gets cluster projects cost metrics.
func (s *OpenCostService) GetProjectsCost(ctx context.Context) (*ClusterCostResponse, *http.Response, error) {
	u := "projects/cost/allocation/compute"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	costResp := new(ClusterCostResponse)
	resp, err := s.client.Do(ctx, req, costResp)
	if err != nil {
		return nil, resp, err
	}

	return costResp, resp, nil
}

// GetNovaServerCost gets total cost calculation for nova server.
func (s *OpenCostService) GetNovaServerCost(ctx context.Context) (*ClusterCostResponse, *http.Response, error) {
	u := "cluster/cost/allocation/compute"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	costResp := new(ClusterCostResponse)
	resp, err := s.client.Do(ctx, req, costResp)
	if err != nil {
		return nil, resp, err
	}

	return costResp, resp, nil
}
