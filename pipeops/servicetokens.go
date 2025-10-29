package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// ServiceTokenService handles communication with service account token related
// methods of the PipeOps API.
type ServiceTokenService struct {
	client *Client
}

// ServiceAccountToken represents a service account token.
type ServiceAccountToken struct {
	UUID        string     `json:"uuid,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Token       string     `json:"token,omitempty"`
	WorkspaceID string     `json:"workspace_id,omitempty"`
	Permissions []string   `json:"permissions,omitempty"`
	ExpiresAt   *Timestamp `json:"expires_at,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`
	LastUsedAt  *Timestamp `json:"last_used_at,omitempty"`
	IsActive    bool       `json:"is_active,omitempty"`
}

// ServiceAccountTokenRequest represents a request to create a service account token.
type ServiceAccountTokenRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	ExpiresAt   string   `json:"expires_at,omitempty"`
}

// ServiceAccountTokenUpdateRequest represents a request to update a service account token.
type ServiceAccountTokenUpdateRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

// ServiceAccountTokenResponse represents the response from service token operations.
type ServiceAccountTokenResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Token ServiceAccountToken `json:"token,omitempty"`
	} `json:"data"`
}

// ServiceAccountTokenListResponse represents a list of service account tokens.
type ServiceAccountTokenListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Tokens []ServiceAccountToken `json:"tokens,omitempty"`
		Total  int                   `json:"total,omitempty"`
	} `json:"data"`
}

// CreateServiceAccountToken creates a new service account token.
func (s *ServiceTokenService) CreateServiceAccountToken(ctx context.Context, req *ServiceAccountTokenRequest) (*ServiceAccountTokenResponse, *http.Response, error) {
	u := "api/v1/service-account-tokens"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	var tokenResp ServiceAccountTokenResponse
	resp, err := s.client.Do(ctx, httpReq, &tokenResp)
	if err != nil {
		return nil, resp, err
	}

	return &tokenResp, resp, nil
}

// ListServiceAccountTokens lists all service account tokens.
func (s *ServiceTokenService) ListServiceAccountTokens(ctx context.Context) (*ServiceAccountTokenListResponse, *http.Response, error) {
	u := "api/v1/service-account-tokens"

	httpReq, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var listResp ServiceAccountTokenListResponse
	resp, err := s.client.Do(ctx, httpReq, &listResp)
	if err != nil {
		return nil, resp, err
	}

	return &listResp, resp, nil
}

// GetServiceAccountToken gets details of a specific service account token.
func (s *ServiceTokenService) GetServiceAccountToken(ctx context.Context, tokenUUID string) (*ServiceAccountTokenResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/service-account-tokens/%s", tokenUUID)

	httpReq, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tokenResp ServiceAccountTokenResponse
	resp, err := s.client.Do(ctx, httpReq, &tokenResp)
	if err != nil {
		return nil, resp, err
	}

	return &tokenResp, resp, nil
}

// UpdateServiceAccountToken updates a service account token.
func (s *ServiceTokenService) UpdateServiceAccountToken(ctx context.Context, tokenUUID string, req *ServiceAccountTokenUpdateRequest) (*ServiceAccountTokenResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/service-account-tokens/%s", tokenUUID)

	httpReq, err := s.client.NewRequest(http.MethodPatch, u, req)
	if err != nil {
		return nil, nil, err
	}

	var tokenResp ServiceAccountTokenResponse
	resp, err := s.client.Do(ctx, httpReq, &tokenResp)
	if err != nil {
		return nil, resp, err
	}

	return &tokenResp, resp, nil
}

// RevokeServiceAccountToken revokes (deletes) a service account token.
func (s *ServiceTokenService) RevokeServiceAccountToken(ctx context.Context, tokenUUID string) (*http.Response, error) {
	u := fmt.Sprintf("api/v1/service-account-tokens/%s", tokenUUID)

	httpReq, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}
