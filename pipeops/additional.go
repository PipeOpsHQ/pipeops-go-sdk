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

// NotificationService handles notification related methods.
type NotificationService struct {
	client *Client
}

// Notification represents a notification.
type Notification struct {
	ID        string     `json:"id,omitempty"`
	UUID      string     `json:"uuid,omitempty"`
	Type      string     `json:"type,omitempty"`
	Title     string     `json:"title,omitempty"`
	Message   string     `json:"message,omitempty"`
	Read      bool       `json:"read,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
}

// NotificationsResponse represents notifications response.
type NotificationsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Notifications []Notification `json:"notifications"`
	} `json:"data"`
}

// ListNotifications lists all user notifications.
func (s *NotificationService) ListNotifications(ctx context.Context) (*NotificationsResponse, *http.Response, error) {
	u := "notifications"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	notificationsResp := new(NotificationsResponse)
	resp, err := s.client.Do(ctx, req, notificationsResp)
	if err != nil {
		return nil, resp, err
	}

	return notificationsResp, resp, nil
}

// MarkAsRead marks a notification as read.
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationUUID string) (*http.Response, error) {
	u := fmt.Sprintf("notifications/%s/read", notificationUUID)

	req, err := s.client.NewRequest(http.MethodPut, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// MarkAllAsRead marks all notifications as read.
func (s *NotificationService) MarkAllAsRead(ctx context.Context) (*http.Response, error) {
	u := "notifications/read-all"

	req, err := s.client.NewRequest(http.MethodPut, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// DeleteNotification deletes a notification.
func (s *NotificationService) DeleteNotification(ctx context.Context, notificationUUID string) (*http.Response, error) {
	u := fmt.Sprintf("notifications/%s", notificationUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// TemplateService handles template related methods.
type TemplateService struct {
	client *Client
}

// Template represents a project template.
type Template struct {
	ID          string `json:"id,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Framework   string `json:"framework,omitempty"`
	Repository  string `json:"repository,omitempty"`
}

// TemplatesResponse represents templates response.
type TemplatesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Templates []Template `json:"templates"`
	} `json:"data"`
}

// ListTemplates lists available project templates.
func (s *TemplateService) ListTemplates(ctx context.Context) (*TemplatesResponse, *http.Response, error) {
	u := "templates"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	templatesResp := new(TemplatesResponse)
	resp, err := s.client.Do(ctx, req, templatesResp)
	if err != nil {
		return nil, resp, err
	}

	return templatesResp, resp, nil
}

// GetTemplate gets a template by UUID.
func (s *TemplateService) GetTemplate(ctx context.Context, templateUUID string) (*TemplatesResponse, *http.Response, error) {
	u := fmt.Sprintf("templates/%s", templateUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	templateResp := new(TemplatesResponse)
	resp, err := s.client.Do(ctx, req, templateResp)
	if err != nil {
		return nil, resp, err
	}

	return templateResp, resp, nil
}

// IntegrationService handles integration related methods.
type IntegrationService struct {
	client *Client
}

// Integration represents a third-party integration.
type Integration struct {
	ID        string     `json:"id,omitempty"`
	UUID      string     `json:"uuid,omitempty"`
	Name      string     `json:"name,omitempty"`
	Type      string     `json:"type,omitempty"`
	Enabled   bool       `json:"enabled,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
}

// IntegrationsResponse represents integrations response.
type IntegrationsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Integrations []Integration `json:"integrations"`
	} `json:"data"`
}

// ListIntegrations lists all integrations.
func (s *IntegrationService) ListIntegrations(ctx context.Context) (*IntegrationsResponse, *http.Response, error) {
	u := "integrations"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	integrationsResp := new(IntegrationsResponse)
	resp, err := s.client.Do(ctx, req, integrationsResp)
	if err != nil {
		return nil, resp, err
	}

	return integrationsResp, resp, nil
}

// ConnectIntegration connects an integration.
func (s *IntegrationService) ConnectIntegration(ctx context.Context, integrationType string) (*http.Response, error) {
	u := fmt.Sprintf("integrations/%s/connect", integrationType)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// DisconnectIntegration disconnects an integration.
func (s *IntegrationService) DisconnectIntegration(ctx context.Context, integrationUUID string) (*http.Response, error) {
	u := fmt.Sprintf("integrations/%s/disconnect", integrationUUID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// HealthCheckService handles health check related methods.
type HealthCheckService struct {
	client *Client
}

// HealthCheckResponse represents health check response.
type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Health map[string]interface{} `json:"health"`
	} `json:"data"`
}

// CheckAPIHealth checks API health.
func (s *HealthCheckService) CheckAPIHealth(ctx context.Context) (*HealthCheckResponse, *http.Response, error) {
	u := "health"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	healthResp := new(HealthCheckResponse)
	resp, err := s.client.Do(ctx, req, healthResp)
	if err != nil {
		return nil, resp, err
	}

	return healthResp, resp, nil
}

// CheckDatabaseHealth checks database health.
func (s *HealthCheckService) CheckDatabaseHealth(ctx context.Context) (*HealthCheckResponse, *http.Response, error) {
	u := "health/database"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	healthResp := new(HealthCheckResponse)
	resp, err := s.client.Do(ctx, req, healthResp)
	if err != nil {
		return nil, resp, err
	}

	return healthResp, resp, nil
}

// BackupService handles backup and restore related methods.
type BackupService struct {
client *Client
}

// Backup represents a backup.
type Backup struct {
ID        string     `json:"id,omitempty"`
UUID      string     `json:"uuid,omitempty"`
ProjectID string     `json:"project_id,omitempty"`
Type      string     `json:"type,omitempty"`
Status    string     `json:"status,omitempty"`
Size      int64      `json:"size,omitempty"`
CreatedAt *Timestamp `json:"created_at,omitempty"`
}

// BackupsResponse represents backups response.
type BackupsResponse struct {
Status  string `json:"status"`
Message string `json:"message"`
Data    struct {
Backups []Backup `json:"backups"`
} `json:"data"`
}

// CreateBackup creates a new backup.
func (s *BackupService) CreateBackup(ctx context.Context, projectUUID string) (*http.Response, error) {
u := fmt.Sprintf("backups/projects/%s", projectUUID)

req, err := s.client.NewRequest(http.MethodPost, u, nil)
if err != nil {
return nil, err
}

resp, err := s.client.Do(ctx, req, nil)
return resp, err
}

// ListBackups lists all backups for a project.
func (s *BackupService) ListBackups(ctx context.Context, projectUUID string) (*BackupsResponse, *http.Response, error) {
u := fmt.Sprintf("backups/projects/%s", projectUUID)

req, err := s.client.NewRequest(http.MethodGet, u, nil)
if err != nil {
return nil, nil, err
}

backupsResp := new(BackupsResponse)
resp, err := s.client.Do(ctx, req, backupsResp)
if err != nil {
return nil, resp, err
}

return backupsResp, resp, nil
}

// RestoreBackup restores a backup.
func (s *BackupService) RestoreBackup(ctx context.Context, backupUUID string) (*http.Response, error) {
u := fmt.Sprintf("backups/%s/restore", backupUUID)

req, err := s.client.NewRequest(http.MethodPost, u, nil)
if err != nil {
return nil, err
}

resp, err := s.client.Do(ctx, req, nil)
return resp, err
}

// DeleteBackup deletes a backup.
func (s *BackupService) DeleteBackup(ctx context.Context, backupUUID string) (*http.Response, error) {
u := fmt.Sprintf("backups/%s", backupUUID)

req, err := s.client.NewRequest(http.MethodDelete, u, nil)
if err != nil {
return nil, err
}

resp, err := s.client.Do(ctx, req, nil)
return resp, err
}

// SecurityScanService handles security scanning related methods.
type SecurityScanService struct {
client *Client
}

// ScanResult represents a security scan result.
type ScanResult struct {
ID           string     `json:"id,omitempty"`
UUID         string     `json:"uuid,omitempty"`
ProjectID    string     `json:"project_id,omitempty"`
Severity     string     `json:"severity,omitempty"`
Vulnerabilities int     `json:"vulnerabilities,omitempty"`
Status       string     `json:"status,omitempty"`
ScannedAt    *Timestamp `json:"scanned_at,omitempty"`
}

// ScanResultsResponse represents scan results response.
type ScanResultsResponse struct {
Status  string `json:"status"`
Message string `json:"message"`
Data    struct {
Results []ScanResult `json:"results"`
} `json:"data"`
}

// ScanProject initiates a security scan for a project.
func (s *SecurityScanService) ScanProject(ctx context.Context, projectUUID string) (*http.Response, error) {
u := fmt.Sprintf("security/scan/projects/%s", projectUUID)

req, err := s.client.NewRequest(http.MethodPost, u, nil)
if err != nil {
return nil, err
}

resp, err := s.client.Do(ctx, req, nil)
return resp, err
}

// GetScanResults retrieves scan results for a project.
func (s *SecurityScanService) GetScanResults(ctx context.Context, projectUUID string) (*ScanResultsResponse, *http.Response, error) {
u := fmt.Sprintf("security/scan/projects/%s/results", projectUUID)

req, err := s.client.NewRequest(http.MethodGet, u, nil)
if err != nil {
return nil, nil, err
}

resultsResp := new(ScanResultsResponse)
resp, err := s.client.Do(ctx, req, resultsResp)
if err != nil {
return nil, resp, err
}

return resultsResp, resp, nil
}

// LogService handles centralized logging related methods.
type LogService struct {
client *Client
}

// LogQuery represents a log query request.
type LogQuery struct {
Query     string `json:"query,omitempty"`
StartTime string `json:"start_time,omitempty"`
EndTime   string `json:"end_time,omitempty"`
Limit     int    `json:"limit,omitempty"`
}

// LogsResponse represents logs response.
type LogsResponse struct {
Status  string `json:"status"`
Message string `json:"message"`
Data    struct {
Logs []map[string]interface{} `json:"logs"`
} `json:"data"`
}

// QueryLogs queries logs across projects.
func (s *LogService) QueryLogs(ctx context.Context, req *LogQuery) (*LogsResponse, *http.Response, error) {
u := "logs/query"

httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
if err != nil {
return nil, nil, err
}

logsResp := new(LogsResponse)
resp, err := s.client.Do(ctx, httpReq, logsResp)
if err != nil {
return nil, resp, err
}

return logsResp, resp, nil
}

// StreamLogs streams logs in real-time.
func (s *LogService) StreamLogs(ctx context.Context, projectUUID string) (*http.Response, error) {
u := fmt.Sprintf("logs/stream/projects/%s", projectUUID)

req, err := s.client.NewRequest(http.MethodGet, u, nil)
if err != nil {
return nil, err
}

resp, err := s.client.Do(ctx, req, nil)
return resp, err
}

// AuditLogService handles audit log related methods.
type AuditLogService struct {
client *Client
}

// AuditLog represents an audit log entry.
type AuditLog struct {
ID         string     `json:"id,omitempty"`
UUID       string     `json:"uuid,omitempty"`
UserID     string     `json:"user_id,omitempty"`
Action     string     `json:"action,omitempty"`
Resource   string     `json:"resource,omitempty"`
IPAddress  string     `json:"ip_address,omitempty"`
CreatedAt  *Timestamp `json:"created_at,omitempty"`
}

// AuditLogsResponse represents audit logs response.
type AuditLogsResponse struct {
Status  string `json:"status"`
Message string `json:"message"`
Data    struct {
Logs []AuditLog `json:"logs"`
} `json:"data"`
}

// ListAuditLogs lists audit logs.
func (s *AuditLogService) ListAuditLogs(ctx context.Context) (*AuditLogsResponse, *http.Response, error) {
u := "audit/logs"

req, err := s.client.NewRequest(http.MethodGet, u, nil)
if err != nil {
return nil, nil, err
}

logsResp := new(AuditLogsResponse)
resp, err := s.client.Do(ctx, req, logsResp)
if err != nil {
return nil, resp, err
}

return logsResp, resp, nil
}

// GetAuditLog gets a specific audit log entry.
func (s *AuditLogService) GetAuditLog(ctx context.Context, logUUID string) (*http.Response, error) {
u := fmt.Sprintf("audit/logs/%s", logUUID)

req, err := s.client.NewRequest(http.MethodGet, u, nil)
if err != nil {
return nil, err
}

resp, err := s.client.Do(ctx, req, nil)
return resp, err
}

// AlertService handles alert and monitoring related methods.
type AlertService struct {
client *Client
}

// Alert represents an alert.
type Alert struct {
ID        string     `json:"id,omitempty"`
UUID      string     `json:"uuid,omitempty"`
Type      string     `json:"type,omitempty"`
Severity  string     `json:"severity,omitempty"`
Message   string     `json:"message,omitempty"`
Resolved  bool       `json:"resolved,omitempty"`
CreatedAt *Timestamp `json:"created_at,omitempty"`
}

// AlertsResponse represents alerts response.
type AlertsResponse struct {
Status  string `json:"status"`
Message string `json:"message"`
Data    struct {
Alerts []Alert `json:"alerts"`
} `json:"data"`
}

// CreateAlertRequest represents create alert request.
type CreateAlertRequest struct {
Type      string `json:"type"`
Threshold int    `json:"threshold"`
ProjectID string `json:"project_id,omitempty"`
}

// CreateAlert creates a new alert rule.
func (s *AlertService) CreateAlert(ctx context.Context, req *CreateAlertRequest) (*http.Response, error) {
u := "alerts"

httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
if err != nil {
return nil, err
}

resp, err := s.client.Do(ctx, httpReq, nil)
return resp, err
}

// ListAlerts lists all alerts.
func (s *AlertService) ListAlerts(ctx context.Context) (*AlertsResponse, *http.Response, error) {
u := "alerts"

req, err := s.client.NewRequest(http.MethodGet, u, nil)
if err != nil {
return nil, nil, err
}

alertsResp := new(AlertsResponse)
resp, err := s.client.Do(ctx, req, alertsResp)
if err != nil {
return nil, resp, err
}

return alertsResp, resp, nil
}

// ResolveAlert resolves an alert.
func (s *AlertService) ResolveAlert(ctx context.Context, alertUUID string) (*http.Response, error) {
u := fmt.Sprintf("alerts/%s/resolve", alertUUID)

req, err := s.client.NewRequest(http.MethodPost, u, nil)
if err != nil {
return nil, err
}

resp, err := s.client.Do(ctx, req, nil)
return resp, err
}

// DeleteAlert deletes an alert rule.
func (s *AlertService) DeleteAlert(ctx context.Context, alertUUID string) (*http.Response, error) {
u := fmt.Sprintf("alerts/%s", alertUUID)

req, err := s.client.NewRequest(http.MethodDelete, u, nil)
if err != nil {
return nil, err
}

resp, err := s.client.Do(ctx, req, nil)
return resp, err
}
