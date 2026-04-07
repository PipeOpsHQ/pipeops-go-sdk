package pipeops

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ServerService handles communication with the server related
// methods of the PipeOps API.
type ServerService struct {
	client *Client
}

// Server represents a PipeOps server.
type Server struct {
	ID          string     `json:"id,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	Name        string     `json:"name,omitempty"`
	Provider    string     `json:"provider,omitempty"`
	Region      string     `json:"region,omitempty"`
	Status      string     `json:"status,omitempty"`
	WorkspaceID string     `json:"workspace_id,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`
}

// ServersResponse represents a list of servers response.
type ServersResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Servers []Server `json:"servers"`
	} `json:"data"`
}

// ServerResponse represents a single server response.
type ServerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Server Server `json:"server"`
	} `json:"data"`
}

// List lists all servers in a cluster.
func (s *ServerService) List(ctx context.Context, workspaceUUID string) (*ServersResponse, *http.Response, error) {
	if workspaceUUID == "" {
		return nil, nil, errors.New("workspace UUID cannot be empty")
	}

	u, err := addOptions("cluster", &clusterWorkspaceOptions{WorkspaceUUID: workspaceUUID})
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	rawResp := new(clusterListResponse)
	resp, err := s.client.Do(ctx, req, rawResp)
	if err != nil {
		return nil, resp, err
	}

	serversResp := &ServersResponse{
		Status:  coalesceNonEmpty(rawResp.Status, statusFromSuccess(rawResp.Success)),
		Message: rawResp.Message,
	}
	for _, cluster := range rawResp.Data.Clusters {
		serversResp.Data.Servers = append(serversResp.Data.Servers, clusterToServer(cluster))
	}

	return serversResp, resp, nil
}

// Get fetches a server by UUID.
func (s *ServerService) Get(ctx context.Context, clusterUUID, workspaceUUID string) (*ServerResponse, *http.Response, error) {
	if clusterUUID == "" {
		return nil, nil, errors.New("cluster UUID cannot be empty")
	}
	if workspaceUUID == "" {
		return nil, nil, errors.New("workspace UUID cannot be empty")
	}

	u, err := addOptions(fmt.Sprintf("cluster/%s", clusterUUID), &clusterWorkspaceOptions{WorkspaceUUID: workspaceUUID})
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	rawResp := new(clusterFetchResponse)
	resp, err := s.client.Do(ctx, req, rawResp)
	if err != nil {
		return nil, resp, err
	}

	clusters, err := parseClusterFetchItems(rawResp.Data)
	if err != nil {
		return nil, resp, err
	}
	if len(clusters) == 0 {
		return nil, resp, errors.New("no cluster data returned")
	}

	serverResp := &ServerResponse{
		Status:  coalesceNonEmpty(rawResp.Status, statusFromSuccess(rawResp.Success)),
		Message: rawResp.Message,
	}
	serverResp.Data.Server = clusterToServer(clusters[0])

	return serverResp, resp, nil
}

// CreateServerRequest represents a request to create a server.
type CreateServerRequest struct {
	ServerName   string `json:"server_name,omitempty"`
	ServerRegion string `json:"server_region,omitempty"`
	ServerType   string `json:"server_type,omitempty"`
	ServerCloud  string `json:"server_cloud,omitempty"`

	Name      string `json:"-"`
	Region    string `json:"-"`
	Port      string `json:"-"`
	IPAddress string `json:"-"`
	Provider  string `json:"-"`
}

// Create creates a new server in a cluster.
func (s *ServerService) Create(ctx context.Context, clusterUUID string, req *CreateServerRequest) (*ServerResponse, *http.Response, error) {
	if req == nil {
		return nil, nil, errors.New("create server request cannot be nil")
	}

	_ = clusterUUID

	u := "server/create"
	payload := &createServerPayload{
		ServerName:   coalesceNonEmpty(req.ServerName, req.Name),
		ServerRegion: coalesceNonEmpty(req.ServerRegion, req.Region),
		ServerType:   req.ServerType,
		ServerCloud:  coalesceNonEmpty(req.ServerCloud, req.Provider),
	}

	if payload.ServerName == "" {
		return nil, nil, errors.New("server name is required")
	}

	httpReq, err := s.client.NewRequest(http.MethodPost, u, payload)
	if err != nil {
		return nil, nil, err
	}

	serverResp := new(ServerResponse)
	resp, err := s.client.Do(ctx, httpReq, serverResp)
	if err != nil {
		return nil, resp, err
	}

	return serverResp, resp, nil
}

// Delete deletes a server from a cluster.
func (s *ServerService) Delete(ctx context.Context, clusterUUID, serverUUID string) (*http.Response, error) {
	if clusterUUID == "" {
		return nil, errors.New("cluster UUID cannot be empty")
	}

	u := fmt.Sprintf("api/v1/clusters/%s", clusterUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err == nil || !isNotFound(err) {
		return resp, err
	}

	if workspaceUUID, _, wsErr := firstWorkspaceUUID(ctx, s.client); wsErr == nil && workspaceUUID != "" {
		withWorkspace, addErr := addOptions(fmt.Sprintf("cluster/%s", clusterUUID), &clusterWorkspaceOptions{WorkspaceUUID: workspaceUUID})
		if addErr == nil {
			req, reqErr := s.client.NewRequest(http.MethodDelete, withWorkspace, nil)
			if reqErr == nil {
				resp, err = s.client.Do(ctx, req, nil)
				if err == nil || !isNotFound(err) {
					return resp, err
				}
			}
		}
	}

	if serverUUID == "" {
		return resp, err
	}

	u = fmt.Sprintf("clusters/%s/servers/%s", clusterUUID, serverUUID)
	req, reqErr := s.client.NewRequest(http.MethodDelete, u, nil)
	if reqErr != nil {
		return resp, err
	}

	return s.client.Do(ctx, req, nil)
}

type clusterWorkspaceOptions struct {
	WorkspaceUUID string `url:"workspace_uuid"`
}

type clusterFetchResponse struct {
	Success bool            `json:"success,omitempty"`
	Status  string          `json:"status,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type clusterListResponse struct {
	Success bool   `json:"success,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Clusters []clusterListItem `json:"clusters,omitempty"`
	} `json:"data,omitempty"`
}

type clusterListItem struct {
	Cluster struct {
		ID               jsonID `json:"id,omitempty"`
		IDAlt            jsonID `json:"ID,omitempty"`
		UUID             string `json:"uuid,omitempty"`
		UUIDAlt          string `json:"UUID,omitempty"`
		Name             string `json:"name,omitempty"`
		NameAlt          string `json:"Name,omitempty"`
		CloudProvider    string `json:"cloudProvider,omitempty"`
		CloudProviderAlt string `json:"CloudProvider,omitempty"`
		Region           string `json:"region,omitempty"`
		RegionAlt        string `json:"Region,omitempty"`
		Status           string `json:"status,omitempty"`
		StatusAlt        string `json:"Status,omitempty"`
		WorkspaceID      jsonID `json:"workspace_id,omitempty"`
		WorkspaceIDAlt   jsonID `json:"WorkspaceID,omitempty"`
	} `json:"Cluster,omitempty"`
	IsActive bool `json:"IsActive,omitempty"`
	InUse    bool `json:"InUse,omitempty"`
}

type clusterFetchItem struct {
	ID               jsonID `json:"id,omitempty"`
	IDAlt            jsonID `json:"ID,omitempty"`
	UUID             string `json:"uuid,omitempty"`
	UUIDAlt          string `json:"UUID,omitempty"`
	Name             string `json:"name,omitempty"`
	NameAlt          string `json:"Name,omitempty"`
	CloudProvider    string `json:"cloudProvider,omitempty"`
	CloudProviderAlt string `json:"CloudProvider,omitempty"`
	Region           string `json:"region,omitempty"`
	RegionAlt        string `json:"Region,omitempty"`
	Status           string `json:"status,omitempty"`
	StatusAlt        string `json:"Status,omitempty"`
	WorkspaceID      jsonID `json:"workspace_id,omitempty"`
	WorkspaceIDAlt   jsonID `json:"WorkspaceID,omitempty"`
	IsActive         *bool  `json:"IsActive,omitempty"`
	IsActiveAlt      *bool  `json:"is_active,omitempty"`
	InUse            *bool  `json:"InUse,omitempty"`
	InUseAlt         *bool  `json:"in_use,omitempty"`
}

type createServerPayload struct {
	ServerName   string `json:"server_name"`
	ServerRegion string `json:"server_region,omitempty"`
	ServerType   string `json:"server_type,omitempty"`
	ServerCloud  string `json:"server_cloud,omitempty"`
}

func clusterToServer(cluster clusterListItem) Server {
	status := coalesceNonEmpty(cluster.Cluster.Status, cluster.Cluster.StatusAlt)
	if status == "" {
		if cluster.IsActive {
			status = "active"
		} else {
			status = "inactive"
		}
	}

	return Server{
		ID:          coalesceNonEmpty(cluster.Cluster.ID.String(), cluster.Cluster.IDAlt.String()),
		UUID:        coalesceNonEmpty(cluster.Cluster.UUID, cluster.Cluster.UUIDAlt),
		Name:        coalesceNonEmpty(cluster.Cluster.Name, cluster.Cluster.NameAlt),
		Provider:    coalesceNonEmpty(cluster.Cluster.CloudProvider, cluster.Cluster.CloudProviderAlt),
		Region:      coalesceNonEmpty(cluster.Cluster.Region, cluster.Cluster.RegionAlt),
		Status:      status,
		WorkspaceID: coalesceNonEmpty(cluster.Cluster.WorkspaceID.String(), cluster.Cluster.WorkspaceIDAlt.String()),
	}
}

func parseClusterFetchItems(data json.RawMessage) ([]clusterListItem, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}

	var wrapped struct {
		Clusters    []clusterListItem `json:"clusters,omitempty"`
		ClustersAlt []clusterListItem `json:"Clusters,omitempty"`
		Cluster     *clusterFetchItem `json:"cluster,omitempty"`
		ClusterAlt  *clusterFetchItem `json:"Cluster,omitempty"`
	}
	if err := json.Unmarshal(data, &wrapped); err == nil {
		switch {
		case len(wrapped.Clusters) > 0:
			return wrapped.Clusters, nil
		case len(wrapped.ClustersAlt) > 0:
			return wrapped.ClustersAlt, nil
		case wrapped.Cluster != nil && !wrapped.Cluster.isEmpty():
			return []clusterListItem{wrapped.Cluster.toListItem()}, nil
		case wrapped.ClusterAlt != nil && !wrapped.ClusterAlt.isEmpty():
			return []clusterListItem{wrapped.ClusterAlt.toListItem()}, nil
		}
	}

	var direct clusterFetchItem
	if err := json.Unmarshal(data, &direct); err != nil {
		return nil, err
	}
	if direct.isEmpty() {
		return nil, nil
	}
	return []clusterListItem{direct.toListItem()}, nil
}

func (c clusterFetchItem) isEmpty() bool {
	return coalesceNonEmpty(
		c.ID.String(),
		c.IDAlt.String(),
		c.UUID,
		c.UUIDAlt,
		c.Name,
		c.NameAlt,
		c.CloudProvider,
		c.CloudProviderAlt,
		c.Region,
		c.RegionAlt,
		c.Status,
		c.StatusAlt,
		c.WorkspaceID.String(),
		c.WorkspaceIDAlt.String(),
	) == ""
}

func (c clusterFetchItem) toListItem() clusterListItem {
	var item clusterListItem
	item.Cluster.ID = jsonID{value: coalesceNonEmpty(c.ID.String(), c.IDAlt.String())}
	item.Cluster.UUID = coalesceNonEmpty(c.UUID, c.UUIDAlt)
	item.Cluster.Name = coalesceNonEmpty(c.Name, c.NameAlt)
	item.Cluster.CloudProvider = coalesceNonEmpty(c.CloudProvider, c.CloudProviderAlt)
	item.Cluster.Region = coalesceNonEmpty(c.Region, c.RegionAlt)
	item.Cluster.Status = coalesceNonEmpty(c.Status, c.StatusAlt)
	item.Cluster.WorkspaceID = jsonID{value: coalesceNonEmpty(c.WorkspaceID.String(), c.WorkspaceIDAlt.String())}
	if c.IsActive != nil {
		item.IsActive = *c.IsActive
	} else if c.IsActiveAlt != nil {
		item.IsActive = *c.IsActiveAlt
	}
	if c.InUse != nil {
		item.InUse = *c.InUse
	} else if c.InUseAlt != nil {
		item.InUse = *c.InUseAlt
	}
	return item
}

func statusFromSuccess(success bool) string {
	if success {
		return "success"
	}
	return "error"
}

func coalesceNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func isNotFound(err error) bool {
	apiErr, ok := err.(*ErrorResponse)
	if !ok || apiErr.Response == nil {
		return false
	}
	return apiErr.Response.StatusCode == http.StatusNotFound
}

// ServiceToken represents a service account token.
type ServiceToken struct {
	ID          string     `json:"id,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	Name        string     `json:"name,omitempty"`
	Token       string     `json:"token,omitempty"`
	Description string     `json:"description,omitempty"`
	ExpiresAt   *Timestamp `json:"expires_at,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
}

// ServiceTokenRequest represents a request to create a service token.
type ServiceTokenRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"` // in days
}

// ServiceTokenResponse represents a service token response.
type ServiceTokenResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Token ServiceToken `json:"token"`
	} `json:"data"`
}

// ServiceTokensResponse represents a list of service tokens response.
type ServiceTokensResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Tokens []ServiceToken `json:"tokens"`
	} `json:"data"`
}

// CreateServiceToken creates a new service account token.
func (s *ServerService) CreateServiceToken(ctx context.Context, req *ServiceTokenRequest) (*ServiceTokenResponse, *http.Response, error) {
	u := "api/v1/service-account-tokens"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	tokenResp := new(ServiceTokenResponse)
	resp, err := s.client.Do(ctx, httpReq, tokenResp)
	if err != nil {
		return nil, resp, err
	}

	return tokenResp, resp, nil
}

// ListServiceTokens lists all service account tokens.
func (s *ServerService) ListServiceTokens(ctx context.Context) (*ServiceTokensResponse, *http.Response, error) {
	u := "api/v1/service-account-tokens"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	tokensResp := new(ServiceTokensResponse)
	resp, err := s.client.Do(ctx, req, tokensResp)
	if err != nil {
		return nil, resp, err
	}

	return tokensResp, resp, nil
}

// GetServiceToken gets a service token by UUID.
func (s *ServerService) GetServiceToken(ctx context.Context, tokenUUID string) (*ServiceTokenResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/service-account-tokens/%s", tokenUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	tokenResp := new(ServiceTokenResponse)
	resp, err := s.client.Do(ctx, req, tokenResp)
	if err != nil {
		return nil, resp, err
	}

	return tokenResp, resp, nil
}

// UpdateServiceTokenRequest represents a request to update a service token.
type UpdateServiceTokenRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateServiceToken updates a service token.
func (s *ServerService) UpdateServiceToken(ctx context.Context, tokenUUID string, req *UpdateServiceTokenRequest) (*ServiceTokenResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/service-account-tokens/%s", tokenUUID)

	httpReq, err := s.client.NewRequest(http.MethodPatch, u, req)
	if err != nil {
		return nil, nil, err
	}

	tokenResp := new(ServiceTokenResponse)
	resp, err := s.client.Do(ctx, httpReq, tokenResp)
	if err != nil {
		return nil, resp, err
	}

	return tokenResp, resp, nil
}

// RevokeServiceToken revokes a service token.
func (s *ServerService) RevokeServiceToken(ctx context.Context, tokenUUID string) (*http.Response, error) {
	u := fmt.Sprintf("api/v1/service-account-tokens/%s", tokenUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// ClusterConnectionResponse represents cluster connection information.
type ClusterConnectionResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Connection map[string]interface{} `json:"connection"`
	} `json:"data"`
}

// GetClusterConnection gets connection information for a cluster.
func (s *ServerService) GetClusterConnection(ctx context.Context, clusterUUID string) (*ClusterConnectionResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/clusters/%s/connection", clusterUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	connResp := new(ClusterConnectionResponse)
	resp, err := s.client.Do(ctx, req, connResp)
	if err != nil {
		return nil, resp, err
	}

	return connResp, resp, nil
}

// AgentRegisterRequest represents an agent registration request.
type AgentRegisterRequest struct {
	ClusterName string                 `json:"cluster_name"`
	ServerSpecs map[string]interface{} `json:"server_specs,omitempty"`
}

// AgentRegisterResponse represents an agent registration response.
type AgentRegisterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ClusterUUID string `json:"cluster_uuid"`
		Token       string `json:"token"`
	} `json:"data"`
}

// RegisterAgent registers a new agent/cluster.
func (s *ServerService) RegisterAgent(ctx context.Context, req *AgentRegisterRequest) (*AgentRegisterResponse, *http.Response, error) {
	u := "api/v1/clusters/agent/register"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	registerResp := new(AgentRegisterResponse)
	resp, err := s.client.Do(ctx, httpReq, registerResp)
	if err != nil {
		return nil, resp, err
	}

	return registerResp, resp, nil
}

// AgentHeartbeatRequest represents an agent heartbeat request.
type AgentHeartbeatRequest struct {
	Status      string                 `json:"status"`
	Metrics     map[string]interface{} `json:"metrics,omitempty"`
	LastUpdated string                 `json:"last_updated,omitempty"`
}

// AgentHeartbeat sends a heartbeat for an agent.
func (s *ServerService) AgentHeartbeat(ctx context.Context, clusterUUID string, req *AgentHeartbeatRequest) (*http.Response, error) {
	u := fmt.Sprintf("api/v1/clusters/agent/%s/heartbeat", clusterUUID)

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// TunnelInfoResponse represents tunnel information response.
type TunnelInfoResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		TunnelInfo map[string]interface{} `json:"tunnel_info"`
	} `json:"data"`
}

// GetTunnelInfo gets tunnel information for a cluster.
func (s *ServerService) GetTunnelInfo(ctx context.Context, clusterUUID string) (*TunnelInfoResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/clusters/agent/%s/tunnel-info", clusterUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	tunnelResp := new(TunnelInfoResponse)
	resp, err := s.client.Do(ctx, req, tunnelResp)
	if err != nil {
		return nil, resp, err
	}

	return tunnelResp, resp, nil
}

// CostAllocationResponse represents cost allocation response.
type CostAllocationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Costs map[string]interface{} `json:"costs"`
	} `json:"data"`
}

// GetClusterCostAllocation gets cost allocation for a cluster.
func (s *ServerService) GetClusterCostAllocation(ctx context.Context, clusterUUID string) (*CostAllocationResponse, *http.Response, error) {
	u := fmt.Sprintf("cluster/%s/cost/allocation/compute", clusterUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	costResp := new(CostAllocationResponse)
	resp, err := s.client.Do(ctx, req, costResp)
	if err != nil {
		return nil, resp, err
	}

	return costResp, resp, nil
}

// UpdateAgentStatusRequest represents agent status update.
type UpdateAgentStatusRequest struct {
	Status string `json:"status"`
}

// UpdateAgentStatus updates agent status.
func (s *ServerService) UpdateAgentStatus(ctx context.Context, clusterUUID string, req *UpdateAgentStatusRequest) (*http.Response, error) {
	u := fmt.Sprintf("api/v1/clusters/agent/%s/status", clusterUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// GetAgentConfig retrieves agent configuration.
func (s *ServerService) GetAgentConfig(ctx context.Context, clusterUUID string) (*http.Response, error) {
	u := fmt.Sprintf("api/v1/clusters/agent/%s/config", clusterUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// SyncAgentConfig syncs agent configuration.
func (s *ServerService) SyncAgentConfig(ctx context.Context, clusterUUID string) (*http.Response, error) {
	u := fmt.Sprintf("api/v1/clusters/agent/%s/sync", clusterUUID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// GetAgentLogs retrieves agent logs.
func (s *ServerService) GetAgentLogs(ctx context.Context, clusterUUID string) (*http.Response, error) {
	u := fmt.Sprintf("api/v1/clusters/agent/%s/logs", clusterUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// GetAgentMetrics retrieves agent metrics.
func (s *ServerService) GetAgentMetrics(ctx context.Context, clusterUUID string) (*http.Response, error) {
	u := fmt.Sprintf("api/v1/clusters/agent/%s/metrics", clusterUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// DeregisterAgent deregisters an agent.
func (s *ServerService) DeregisterAgent(ctx context.Context, clusterUUID string) (*http.Response, error) {
	u := fmt.Sprintf("api/v1/clusters/agent/%s/deregister", clusterUUID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// PollAgent polls for agent tasks.
func (s *ServerService) PollAgent(ctx context.Context, clusterUUID string) (*http.Response, error) {
	u := fmt.Sprintf("api/v1/clusters/agent/%s/poll", clusterUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// GetAgentTunnelStatus gets the tunnel status for an agent.
func (s *ServerService) GetAgentTunnelStatus(ctx context.Context, agentID string) (*http.Response, error) {
	u := fmt.Sprintf("api/agents/%s/tunnel/status", agentID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}
