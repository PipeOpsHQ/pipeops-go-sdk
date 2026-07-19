package pipeops

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// AddOnService handles communication with the add-on related
// methods of the PipeOps API.
type AddOnService struct {
	client *Client
}

// AddOn represents a PipeOps add-on.
type AddOn struct {
	ID          string     `json:"id,omitempty"`
	UID         string     `json:"UID,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	Name        string     `json:"Name,omitempty"`
	Description string     `json:"Description,omitempty"`
	Category    string     `json:"Category,omitempty"`
	Version     string     `json:"version,omitempty"`
	Icon        string     `json:"icon,omitempty"`
	ImageURL    string     `json:"ImageURL,omitempty"`
	Status      string     `json:"SubmissionStatus,omitempty"`
	IsFeatured  bool       `json:"IsFeatured,omitempty"`
	IsVerified  bool       `json:"IsVerified,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`
}

// AddOnsResponse represents a list of add-ons response.
type AddOnsResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Data    []AddOn `json:"data"`
}

// AddOnResponse represents a single add-on response.
type AddOnResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    AddOn  `json:"data"`
}

// List lists all available add-ons.
// ListAddOnsOptions specifies optional parameters for listing addons.
type ListAddOnsOptions struct {
	Page          int    `url:"page,omitempty"`
	Limit         int    `url:"limit,omitempty"`
	Size          int    `url:"size,omitempty"`
	Category      string `url:"category,omitempty"`
	Search        string `url:"s,omitempty"`
	Featured      *bool  `url:"featured,omitempty"`
	WorkspaceUUID string `url:"workspace,omitempty"`
}

// List lists all available add-ons.
func (s *AddOnService) List(ctx context.Context, opts ...*ListAddOnsOptions) (*AddOnsResponse, *http.Response, error) {
	u := "addons"
	if len(opts) > 0 && opts[0] != nil {
		var err error
		u, err = addOptions(u, opts[0])
		if err != nil {
			return nil, nil, err
		}
	}

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

// Search searches available add-ons using the same filters as List.
func (s *AddOnService) Search(ctx context.Context, query string, opts ...*ListAddOnsOptions) (*AddOnsResponse, *http.Response, error) {
	listOpts := &ListAddOnsOptions{Search: query}
	if len(opts) > 0 && opts[0] != nil {
		*listOpts = *opts[0]
		listOpts.Search = query
	}
	return s.List(ctx, listOpts)
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
	UID               string     `json:"UID,omitempty"`
	Name              string     `json:"Name,omitempty"`
	DeploymentName    string     `json:"DeploymentName,omitempty"`
	DeploymentURL     string     `json:"DeploymentURL,omitempty"`
	Category          string     `json:"Category,omitempty"`
	Status            string     `json:"Status,omitempty"`
	StatusMessage     string     `json:"StatusMessage,omitempty"`
	Environment       string     `json:"Environment,omitempty"`
	ImageURL          string     `json:"ImageURL,omitempty"`
	Version           string     `json:"Version,omitempty"`
	CurrentVersion    string     `json:"current_version,omitempty"`
	UpgradableVersion string     `json:"upgradable_version,omitempty"`
	UpgradeAvailable  bool       `json:"upgrade_available,omitempty"`
	CreatedAt         *Timestamp `json:"CreatedAt,omitempty"`
	UpdatedAt         *Timestamp `json:"UpdatedAt,omitempty"`
}

// AddOnDeploymentsResponse represents a list of add-on deployments response.
type AddOnDeploymentsResponse struct {
	Data []AddOnDeployment `json:"data"`
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
// Prefer-client on the control plane fills Config from the marketplace catalog
// when omitted; thin clients only need addon ID + workspace + server.
type DeployAddOnRequest struct {
	// ID is the marketplace addon UID (also accepted as Deployment.ID).
	ID string `json:"id,omitempty"`
	// Server is the cluster UUID.
	Server string `json:"Server,omitempty"`
	// Workspace is the workspace UUID.
	Workspace string `json:"Workspace,omitempty"`
	// Environment is the environment UUID. When empty, control plane picks the
	// first environment on the target cluster (prefer-client).
	Environment string `json:"Environment,omitempty"`
	// ProjectID is optional placement hint (reserved / future).
	ProjectID string `json:"project_id,omitempty"`
	// Tag is the image/version tag override.
	Tag string `json:"Tag,omitempty"`
	// Config is optional partial deployment config. Gaps are filled from catalog.
	Config map[string]interface{} `json:"config,omitempty"`
}

// deployAddonsWireBody matches POST /addons/deploy (dashboard + prefer-client API).
type deployAddonsWireBody struct {
	Workspace   string                 `json:"Workspace"`
	Server      string                 `json:"Server"`
	Environment string                 `json:"Environment,omitempty"`
	Deployment  deployAddonsWireDeploy `json:"Deployment"`
	// Thin aliases also accepted by control-plane UnmarshalJSON.
	ID     string                 `json:"id,omitempty"`
	Config map[string]interface{} `json:"config,omitempty"`
}

type deployAddonsWireDeploy struct {
	ID     string                 `json:"ID,omitempty"`
	Tag    string                 `json:"Tag,omitempty"`
	Config map[string]interface{} `json:"Config,omitempty"`
}

// Deploy deploys an add-on via POST /addons/deploy.
// Sends both nested dashboard shape and thin aliases so older/newer controllers accept it.
// Prefer-client fills missing Config/Environment from catalog and cluster defaults.
func (s *AddOnService) Deploy(ctx context.Context, req *DeployAddOnRequest) (*AddOnDeploymentResponse, *http.Response, error) {
	if req == nil {
		return nil, nil, errors.New("deploy request cannot be nil")
	}
	id := strings.TrimSpace(req.ID)
	if id == "" {
		return nil, nil, errors.New("addon id is required")
	}
	workspace := strings.TrimSpace(req.Workspace)
	if workspace == "" {
		if ws, _, wsErr := firstWorkspaceUUID(ctx, s.client); wsErr == nil {
			workspace = ws
		}
	}
	if workspace == "" {
		return nil, nil, errors.New("workspace is required")
	}
	server := strings.TrimSpace(req.Server)
	if server == "" {
		return nil, nil, errors.New("server (cluster UUID) is required")
	}

	body := &deployAddonsWireBody{
		Workspace:   workspace,
		Server:      server,
		Environment: strings.TrimSpace(req.Environment),
		ID:          id,
		Config:      req.Config,
		Deployment: deployAddonsWireDeploy{
			ID:     id,
			Tag:    strings.TrimSpace(req.Tag),
			Config: req.Config,
		},
	}

	u := "addons/deploy"
	httpReq, err := s.client.NewRequest(http.MethodPost, u, body)
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

// ListDeploymentsOptions specifies optional parameters for listing deployments.
type ListDeploymentsOptions struct {
	WorkspaceUUID string `url:"workspace,omitempty"`
}

// ListDeployments lists all add-on deployments for a workspace.
func (s *AddOnService) ListDeployments(ctx context.Context, opts ...*ListDeploymentsOptions) (*AddOnDeploymentsResponse, *http.Response, error) {
	u := "addons/deployments/overview"

	workspaceUUID := ""
	if len(opts) > 0 && opts[0] != nil {
		workspaceUUID = opts[0].WorkspaceUUID
	}
	// Overview requires workspace scoping; resolve default when callers omit it
	// (team members / SA dual-auth). Prefer explicit opts over first workspace.
	if workspaceUUID == "" {
		if ws, _, wsErr := firstWorkspaceUUID(ctx, s.client); wsErr == nil {
			workspaceUUID = ws
		}
	}
	if workspaceUUID != "" {
		u = u + "?workspace=" + workspaceUUID
	}

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

type AddOnCategoriesData struct {
	Categories []AddOnCategory `json:"categories,omitempty"`
}

func (d *AddOnCategoriesData) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || string(trimmed) == "null" {
		return nil
	}

	switch trimmed[0] {
	case '[':
		var categories []AddOnCategory
		if err := json.Unmarshal(trimmed, &categories); err != nil {
			return err
		}
		d.Categories = categories
		return nil
	case '{':
		var wrapped struct {
			Categories []AddOnCategory `json:"categories,omitempty"`
		}
		if err := json.Unmarshal(trimmed, &wrapped); err == nil && wrapped.Categories != nil {
			d.Categories = wrapped.Categories
			return nil
		}

		var single AddOnCategory
		if err := json.Unmarshal(trimmed, &single); err != nil {
			return err
		}
		if single.ID == "" && single.UUID == "" && single.Name == "" && single.Description == "" && single.Icon == "" {
			return nil
		}
		d.Categories = []AddOnCategory{single}
		return nil
	default:
		return fmt.Errorf("unexpected add-on categories data: %s", string(trimmed))
	}
}

// AddOnCategoriesResponse represents a list of add-on categories response.
type AddOnCategoriesResponse struct {
	Success bool                `json:"success,omitempty"`
	Status  string              `json:"status,omitempty"`
	Message string              `json:"message"`
	Data    AddOnCategoriesData `json:"data,omitempty"`
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
// Prefer ListDeployments for typed deployment rows; this helper remains for
// callers that expect the generic overview envelope.
func (s *AddOnService) GetDeploymentOverview(ctx context.Context) (*DeploymentOverviewResponse, *http.Response, error) {
	u := "addons/deployments/overview"
	if workspaceUUID, _, wsErr := firstWorkspaceUUID(ctx, s.client); wsErr == nil && workspaceUUID != "" {
		u = u + "?workspace=" + workspaceUUID
	}

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
	u := fmt.Sprintf("addons/domains/%s", addonUUID)
	if workspaceUUID, _, wsErr := firstWorkspaceUUID(ctx, s.client); wsErr == nil {
		if withWorkspace, err := addOptions(u, &addonWorkspaceOptions{Workspace: workspaceUUID}); err == nil {
			u = withWorkspace
		}
	}

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// --- Addon backup export (data-layer addons) ---
// Controller: GET/POST /addons/deployments/:id/backups*

// AddonBackupSnapshot is one snapshot row in the Backups tab.
type AddonBackupSnapshot struct {
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	Time           string   `json:"time,omitempty"`
	TotalSizeBytes int64    `json:"total_size_bytes,omitempty"`
	Useful         bool     `json:"useful,omitempty"`
	SizeUnknown    bool     `json:"size_unknown,omitempty"`
	TypeChip       string   `json:"type_chip,omitempty"`
	Paths          []string `json:"paths,omitempty"`
	Warning        string   `json:"warning,omitempty"`
	Hostname       string   `json:"hostname,omitempty"`
}

// AddonBackupListResponse is GET /addons/deployments/:id/backups.
type AddonBackupListResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		AddonUID   string                `json:"addon_uid,omitempty"`
		Namespace  string                `json:"namespace,omitempty"`
		ServerCode string                `json:"server_code,omitempty"`
		Snapshots  []AddonBackupSnapshot `json:"snapshots,omitempty"`
	} `json:"data"`
}

// AddonBackupExportRequest is POST /addons/deployments/:id/backups/export body.
type AddonBackupExportRequest struct {
	SnapshotID string `json:"snapshot_id"`
	Path       string `json:"path,omitempty"`
	Format     string `json:"format,omitempty"` // auto | sql | rdb | archive
}

// AddonBackupExportStatus is create/get export status.
type AddonBackupExportStatus struct {
	ExportID     string `json:"export_id,omitempty"`
	Status       string `json:"status,omitempty"`
	DownloadURL  string `json:"download_url,omitempty"`
	Filename     string `json:"filename,omitempty"`
	ContentType  string `json:"content_type,omitempty"`
	SizeBytes    int64  `json:"size_bytes,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	SnapshotID   string `json:"snapshot_id,omitempty"`
	Path         string `json:"path,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
}

// AddonBackupExportResponse wraps export status.
type AddonBackupExportResponse struct {
	Success bool                    `json:"success,omitempty"`
	Message string                  `json:"message,omitempty"`
	Data    AddonBackupExportStatus `json:"data"`
}

// ListAddonBackups lists snapshots for an addon deployment.
// GET /addons/deployments/:id/backups
func (s *AddOnService) ListAddonBackups(ctx context.Context, deploymentUID string) (*AddonBackupListResponse, *http.Response, error) {
	u := fmt.Sprintf("addons/deployments/%s/backups", deploymentUID)
	u = withAddonWorkspaceQuery(ctx, s.client, u)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	out := new(AddonBackupListResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// StartAddonBackupExport starts an async backup export for a snapshot path.
// POST /addons/deployments/:id/backups/export
func (s *AddOnService) StartAddonBackupExport(ctx context.Context, deploymentUID string, body *AddonBackupExportRequest) (*AddonBackupExportResponse, *http.Response, error) {
	u := fmt.Sprintf("addons/deployments/%s/backups/export", deploymentUID)
	u = withAddonWorkspaceQuery(ctx, s.client, u)

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}
	out := new(AddonBackupExportResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// GetAddonBackupExport polls export status.
// GET /addons/deployments/:id/backups/exports/:export_id
func (s *AddOnService) GetAddonBackupExport(ctx context.Context, deploymentUID, exportID string) (*AddonBackupExportResponse, *http.Response, error) {
	u := fmt.Sprintf("addons/deployments/%s/backups/exports/%s", deploymentUID, exportID)
	u = withAddonWorkspaceQuery(ctx, s.client, u)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	out := new(AddonBackupExportResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// DownloadAddonBackupExport returns the download response (follow DownloadURL or stream).
// GET /addons/deployments/:id/backups/exports/:export_id/download
func (s *AddOnService) DownloadAddonBackupExport(ctx context.Context, deploymentUID, exportID string) (*http.Response, error) {
	u := fmt.Sprintf("addons/deployments/%s/backups/exports/%s/download", deploymentUID, exportID)
	u = withAddonWorkspaceQuery(ctx, s.client, u)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

func withAddonWorkspaceQuery(ctx context.Context, client *Client, path string) string {
	if workspaceUUID, _, wsErr := firstWorkspaceUUID(ctx, client); wsErr == nil {
		if withWorkspace, err := addOptions(path, &addonWorkspaceOptions{Workspace: workspaceUUID}); err == nil {
			return withWorkspace
		}
	}
	return path
}

type addonWorkspaceOptions struct {
	Workspace string `url:"workspace"`
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

// Admin Add-On Endpoints

// GetSubmittedAddOns retrieves submitted add-ons (admin only).
func (s *AddOnService) GetSubmittedAddOns(ctx context.Context) (*MySubmissionsResponse, *http.Response, error) {
	u := "admin/addons/submissions"

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

// ReviewAddOnRequest represents an add-on review request.
type ReviewAddOnRequest struct {
	Status   string `json:"status"` // "approved" or "rejected"
	Comments string `json:"comments,omitempty"`
}

// ReviewAddOnApprove approves an add-on submission (admin only).
func (s *AddOnService) ReviewAddOnApprove(ctx context.Context, addonUUID string, req *ReviewAddOnRequest) (*http.Response, error) {
	u := fmt.Sprintf("admin/addons/%s/review", addonUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// PublishAddOn publishes an approved add-on (admin only).
func (s *AddOnService) PublishAddOn(ctx context.Context, addonUUID string) (*http.Response, error) {
	u := fmt.Sprintf("admin/addons/%s/publish", addonUUID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// UnpublishAddOn unpublishes an add-on (admin only).
func (s *AddOnService) UnpublishAddOn(ctx context.Context, addonUUID string) (*http.Response, error) {
	u := fmt.Sprintf("admin/addons/%s/unpublish", addonUUID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// DeleteAddOn deletes an add-on (admin only).
func (s *AddOnService) DeleteAddOn(ctx context.Context, addonUUID string) (*http.Response, error) {
	u := fmt.Sprintf("admin/addons/%s", addonUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}
