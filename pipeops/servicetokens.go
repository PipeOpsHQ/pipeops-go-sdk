package pipeops

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
	TokenPrefix string     `json:"token_prefix,omitempty"`
	WorkspaceID string     `json:"workspace_id,omitempty"`
	Permissions []string   `json:"permissions,omitempty"`
	Scopes      []string   `json:"scopes,omitempty"`
	ExpiresAt   *Timestamp `json:"expires_at,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`
	LastUsedAt  *Timestamp `json:"last_used_at,omitempty"`
	IsActive    bool       `json:"is_active,omitempty"`
}

// ServiceAccountTokenRequest represents a request to create a service account token.
type ServiceAccountTokenRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	Permissions   []string `json:"permissions,omitempty"`
	ExpiresAt     string   `json:"expires_at,omitempty"`
	WorkspaceUUID string   `json:"workspace_uuid,omitempty"` // required by controller
	Preset        string   `json:"preset,omitempty"`         // e.g. sandbox, mcp, sdk
}

// ServiceAccountTokenUpdateRequest represents a request to update a service account token.
type ServiceAccountTokenUpdateRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

// ServiceTokenWorkspaceOptions scopes service-account-token routes.
// CheckWorkspaceIntegrationsAccess requires workspace or workspace_uuid.
type ServiceTokenWorkspaceOptions struct {
	WorkspaceUUID string `url:"workspace_uuid,omitempty"`
}

// ServiceAccountTokenResponse represents the response from service token operations.
// Create returns a flat data object with token as a string; get/list nest under token.
type ServiceAccountTokenResponse struct {
	Success bool   `json:"success,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Token ServiceAccountToken `json:"token,omitempty"`
	} `json:"data"`
}

// UnmarshalJSON accepts both create (flat data.token string) and get (data.token object) shapes.
func (r *ServiceAccountTokenResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		Success bool            `json:"success"`
		Status  string          `json:"status"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.Success = raw.Success
	r.Status = raw.Status
	r.Message = raw.Message
	if len(raw.Data) == 0 || string(raw.Data) == "null" {
		return nil
	}

	// Nested: { "token": { "uuid": "...", ... } }
	var nested struct {
		Token ServiceAccountToken `json:"token"`
	}
	if err := json.Unmarshal(raw.Data, &nested); err == nil && nested.Token.UUID != "" {
		r.Data.Token = nested.Token
		return nil
	}
	// Nested with empty uuid but present name (list item style)
	if nested.Token.Name != "" || nested.Token.Token != "" {
		r.Data.Token = nested.Token
		return nil
	}

	// Flat create: { "id": "...", "token": "sat_...", "name": "..." }
	var flat struct {
		ID          string          `json:"id"`
		UUID        string          `json:"uuid"`
		Name        string          `json:"name"`
		Token       json.RawMessage `json:"token"`
		TokenPrefix string          `json:"token_prefix"`
		WorkspaceID string          `json:"workspace_id"`
		Scopes      []string        `json:"scopes"`
		Permissions []string        `json:"permissions"`
		ExpiresAt   *Timestamp      `json:"expires_at"`
		CreatedAt   *Timestamp      `json:"created_at"`
		IsActive    bool            `json:"is_active"`
	}
	if err := json.Unmarshal(raw.Data, &flat); err != nil {
		return err
	}
	tok := ServiceAccountToken{
		UUID:        firstNonEmpty(flat.UUID, flat.ID),
		Name:        flat.Name,
		TokenPrefix: flat.TokenPrefix,
		WorkspaceID: flat.WorkspaceID,
		Scopes:      flat.Scopes,
		Permissions: flat.Permissions,
		ExpiresAt:   flat.ExpiresAt,
		CreatedAt:   flat.CreatedAt,
		IsActive:    flat.IsActive,
	}
	if len(flat.Token) > 0 {
		var secret string
		if err := json.Unmarshal(flat.Token, &secret); err == nil {
			tok.Token = secret
		} else {
			var nestedTok ServiceAccountToken
			if err := json.Unmarshal(flat.Token, &nestedTok); err == nil {
				if tok.UUID == "" {
					tok.UUID = nestedTok.UUID
				}
				if tok.Name == "" {
					tok.Name = nestedTok.Name
				}
				tok.Token = nestedTok.Token
			}
		}
	}
	r.Data.Token = tok
	return nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// ServiceAccountTokenListResponse represents a list of service account tokens.
type ServiceAccountTokenListResponse struct {
	Success bool   `json:"success,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Tokens []ServiceAccountToken `json:"tokens,omitempty"`
		Total  int                   `json:"total,omitempty"`
	} `json:"data"`
}

func withServiceTokenWorkspace(path string, opts *ServiceTokenWorkspaceOptions) string {
	if opts == nil || strings.TrimSpace(opts.WorkspaceUUID) == "" {
		return path
	}
	sep := "?"
	if strings.Contains(path, "?") {
		sep = "&"
	}
	return path + sep + "workspace_uuid=" + url.QueryEscape(strings.TrimSpace(opts.WorkspaceUUID))
}

// CreateServiceAccountToken creates a new service account token.
// Controller requires workspace_uuid on the body (and often on the query).
func (s *ServiceTokenService) CreateServiceAccountToken(ctx context.Context, req *ServiceAccountTokenRequest) (*ServiceAccountTokenResponse, *http.Response, error) {
	u := "api/v1/service-account-tokens"
	if req != nil && strings.TrimSpace(req.WorkspaceUUID) != "" {
		u = u + "?workspace_uuid=" + url.QueryEscape(strings.TrimSpace(req.WorkspaceUUID))
	} else if ws, _, err := firstWorkspaceUUID(ctx, s.client); err == nil && ws != "" {
		if req == nil {
			req = &ServiceAccountTokenRequest{}
		}
		if req.WorkspaceUUID == "" {
			req.WorkspaceUUID = ws
		}
		u = u + "?workspace_uuid=" + url.QueryEscape(ws)
	}

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

// ListServiceAccountTokens lists service account tokens for a workspace.
func (s *ServiceTokenService) ListServiceAccountTokens(ctx context.Context, opts *ServiceTokenWorkspaceOptions) (*ServiceAccountTokenListResponse, *http.Response, error) {
	u := withServiceTokenWorkspace("api/v1/service-account-tokens", opts)

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
func (s *ServiceTokenService) GetServiceAccountToken(ctx context.Context, tokenUUID string, opts *ServiceTokenWorkspaceOptions) (*ServiceAccountTokenResponse, *http.Response, error) {
	u := withServiceTokenWorkspace(fmt.Sprintf("api/v1/service-account-tokens/%s", tokenUUID), opts)

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
func (s *ServiceTokenService) UpdateServiceAccountToken(ctx context.Context, tokenUUID string, req *ServiceAccountTokenUpdateRequest, opts *ServiceTokenWorkspaceOptions) (*ServiceAccountTokenResponse, *http.Response, error) {
	u := withServiceTokenWorkspace(fmt.Sprintf("api/v1/service-account-tokens/%s", tokenUUID), opts)

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
func (s *ServiceTokenService) RevokeServiceAccountToken(ctx context.Context, tokenUUID string, opts *ServiceTokenWorkspaceOptions) (*http.Response, error) {
	u := withServiceTokenWorkspace(fmt.Sprintf("api/v1/service-account-tokens/%s", tokenUUID), opts)

	httpReq, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}
