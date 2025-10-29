package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// ServerService handles communication with the server/cluster related
// methods of the PipeOps API.
type ServerService struct {
	client *Client
}

// Server represents a PipeOps server/cluster.
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

// List lists all servers.
func (s *ServerService) List(ctx context.Context) (*ServersResponse, *http.Response, error) {
	u := "server"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	serversResp := new(ServersResponse)
	resp, err := s.client.Do(ctx, req, serversResp)
	if err != nil {
		return nil, resp, err
	}

	return serversResp, resp, nil
}

// Get fetches a server by UUID.
func (s *ServerService) Get(ctx context.Context, serverUUID string) (*ServerResponse, *http.Response, error) {
	u := fmt.Sprintf("server/%s", serverUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	serverResp := new(ServerResponse)
	resp, err := s.client.Do(ctx, req, serverResp)
	if err != nil {
		return nil, resp, err
	}

	return serverResp, resp, nil
}

// CreateServerRequest represents a request to create a server.
type CreateServerRequest struct {
	Name         string `json:"name"`
	Provider     string `json:"provider"`
	Region       string `json:"region"`
	WorkspaceID  string `json:"workspace_id"`
	InstanceType string `json:"instance_type,omitempty"`
}

// Create creates a new server.
func (s *ServerService) Create(ctx context.Context, req *CreateServerRequest) (*ServerResponse, *http.Response, error) {
	u := "server/create"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
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

// Delete deletes a server.
func (s *ServerService) Delete(ctx context.Context, serverUUID string) (*http.Response, error) {
	u := fmt.Sprintf("api/v1/clusters/%s", serverUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
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
