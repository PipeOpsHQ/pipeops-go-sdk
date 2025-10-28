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
