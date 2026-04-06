package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// ExternalRegistryService handles communication with BYOI/external registry endpoints.
type ExternalRegistryService struct {
	client *Client
}

// ExternalRegistry represents an external container registry configuration.
type ExternalRegistry struct {
	ID          int        `json:"id,omitempty"`
	UID         string     `json:"uid,omitempty"`
	Name        string     `json:"name,omitempty"`
	Type        string     `json:"type,omitempty"`
	RegistryURL string     `json:"registry_url,omitempty"`
	Username    string     `json:"username,omitempty"`
	Region      string     `json:"region,omitempty"`
	AccountID   string     `json:"account_id,omitempty"`
	IsActive    bool       `json:"is_active,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`
}

// CreateExternalRegistryRequest represents a request to create an external registry.
type CreateExternalRegistryRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	RegistryURL string `json:"registry_url,omitempty"`
	Region      string `json:"region,omitempty"`
	AccountID   string `json:"account_id,omitempty"`
}

// ExternalRegistryListOptions specifies optional parameters for listing registries.
type ExternalRegistryListOptions struct {
	Page     int `url:"page,omitempty"`
	PageSize int `url:"page_size,omitempty"`
}

// DockerHubListOptions specifies optional parameters for Docker Hub list/search operations.
type DockerHubListOptions struct {
	Page     int `url:"page,omitempty"`
	PageSize int `url:"page_size,omitempty"`
}

// DockerHubSearchOptions specifies optional parameters for public image search.
type DockerHubSearchOptions struct {
	Query    string `url:"q,omitempty"`
	Page     int    `url:"page,omitempty"`
	PageSize int    `url:"page_size,omitempty"`
}

// ExternalRegistryResponse represents a single external registry response.
type ExternalRegistryResponse struct {
	Success bool             `json:"success,omitempty"`
	Status  string           `json:"status,omitempty"`
	Message string           `json:"message,omitempty"`
	Data    ExternalRegistry `json:"data"`
}

// ExternalRegistryListResponse represents a list of external registries response.
type ExternalRegistryListResponse struct {
	Success bool   `json:"success,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Registries []ExternalRegistry `json:"registries"`
		Total      int                `json:"total,omitempty"`
		Page       int                `json:"page,omitempty"`
		PageSize   int                `json:"page_size,omitempty"`
	} `json:"data"`
}

// DockerHubRepository represents a repository returned by registry APIs.
type DockerHubRepository struct {
	Name             string `json:"name,omitempty"`
	Namespace        string `json:"namespace,omitempty"`
	FullName         string `json:"full_name,omitempty"`
	Description      string `json:"description,omitempty"`
	ShortDescription string `json:"short_description,omitempty"`
	IsPrivate        bool   `json:"is_private,omitempty"`
	StarCount        int    `json:"star_count,omitempty"`
	PullCount        int    `json:"pull_count,omitempty"`
	LastUpdated      string `json:"last_updated,omitempty"`
}

// DockerHubTag represents a Docker image tag.
type DockerHubTag struct {
	Name        string `json:"name,omitempty"`
	FullSize    int64  `json:"full_size,omitempty"`
	LastUpdated string `json:"last_updated,omitempty"`
	Digest      string `json:"digest,omitempty"`
}

// DockerHubRepositoriesResponse represents a repositories response.
type DockerHubRepositoriesResponse struct {
	Success bool   `json:"success,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Repositories []DockerHubRepository `json:"repositories,omitempty"`
		Results      []DockerHubRepository `json:"results,omitempty"`
		Total        int                   `json:"total,omitempty"`
		Page         int                   `json:"page,omitempty"`
		PageSize     int                   `json:"page_size,omitempty"`
		HasMore      bool                  `json:"has_more,omitempty"`
	} `json:"data"`
}

// DockerHubTagsResponse represents an image tags response.
type DockerHubTagsResponse struct {
	Success bool   `json:"success,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Tags     []DockerHubTag `json:"tags"`
		Total    int            `json:"total,omitempty"`
		Page     int            `json:"page,omitempty"`
		PageSize int            `json:"page_size,omitempty"`
		HasMore  bool           `json:"has_more,omitempty"`
	} `json:"data"`
}

func withWorkspaceQuery(path, workspaceUUID string) (string, error) {
	if workspaceUUID == "" {
		return "", fmt.Errorf("workspace_uuid is required")
	}
	return addOptions(path, &struct {
		WorkspaceUUID string `url:"workspace_uuid"`
	}{WorkspaceUUID: workspaceUUID})
}

// Create creates a new external registry in a workspace.
func (s *ExternalRegistryService) Create(ctx context.Context, workspaceUUID string, req *CreateExternalRegistryRequest) (*ExternalRegistryResponse, *http.Response, error) {
	u, err := withWorkspaceQuery("api/v1/external-registry", workspaceUUID)
	if err != nil {
		return nil, nil, err
	}

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	registryResp := new(ExternalRegistryResponse)
	resp, err := s.client.Do(ctx, httpReq, registryResp)
	if err != nil {
		return nil, resp, err
	}

	return registryResp, resp, nil
}

// List lists external registries for a workspace.
func (s *ExternalRegistryService) List(ctx context.Context, workspaceUUID string, opts *ExternalRegistryListOptions) (*ExternalRegistryListResponse, *http.Response, error) {
	if workspaceUUID == "" {
		return nil, nil, fmt.Errorf("workspace_uuid is required")
	}

	queryOpts := &struct {
		WorkspaceUUID string `url:"workspace_uuid"`
		Page          int    `url:"page,omitempty"`
		PageSize      int    `url:"page_size,omitempty"`
	}{
		WorkspaceUUID: workspaceUUID,
	}
	if opts != nil {
		queryOpts.Page = opts.Page
		queryOpts.PageSize = opts.PageSize
	}

	u, err := addOptions("api/v1/external-registry", queryOpts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	listResp := new(ExternalRegistryListResponse)
	resp, err := s.client.Do(ctx, req, listResp)
	if err != nil {
		return nil, resp, err
	}

	return listResp, resp, nil
}

// Get gets an external registry by UID.
func (s *ExternalRegistryService) Get(ctx context.Context, registryUID string) (*ExternalRegistryResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/external-registry/%s", registryUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	registryResp := new(ExternalRegistryResponse)
	resp, err := s.client.Do(ctx, req, registryResp)
	if err != nil {
		return nil, resp, err
	}

	return registryResp, resp, nil
}

// Delete deletes an external registry by UID.
func (s *ExternalRegistryService) Delete(ctx context.Context, registryUID string) (*http.Response, error) {
	u := fmt.Sprintf("api/v1/external-registry/%s", registryUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// ListDockerHubImages lists repositories for an authenticated Docker Hub registry.
func (s *ExternalRegistryService) ListDockerHubImages(ctx context.Context, registryUID string, opts *DockerHubListOptions) (*DockerHubRepositoriesResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/external-registry/%s/dockerhub/images", registryUID)
	var err error
	if opts != nil {
		u, err = addOptions(u, opts)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	reposResp := new(DockerHubRepositoriesResponse)
	resp, err := s.client.Do(ctx, req, reposResp)
	if err != nil {
		return nil, resp, err
	}

	return reposResp, resp, nil
}

// ListDockerHubTags lists tags for a repository in an authenticated Docker Hub registry.
func (s *ExternalRegistryService) ListDockerHubTags(ctx context.Context, registryUID, namespace, repository string, opts *DockerHubListOptions) (*DockerHubTagsResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/external-registry/%s/dockerhub/%s/%s/tags", registryUID, namespace, repository)
	var err error
	if opts != nil {
		u, err = addOptions(u, opts)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	tagsResp := new(DockerHubTagsResponse)
	resp, err := s.client.Do(ctx, req, tagsResp)
	if err != nil {
		return nil, resp, err
	}

	return tagsResp, resp, nil
}

// SearchPublicDockerHubImages searches public Docker Hub images without a registry configuration.
func (s *ExternalRegistryService) SearchPublicDockerHubImages(ctx context.Context, opts *DockerHubSearchOptions) (*DockerHubRepositoriesResponse, *http.Response, error) {
	u := "api/v1/external-registry/dockerhub/search"
	var err error
	if opts != nil {
		u, err = addOptions(u, opts)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	searchResp := new(DockerHubRepositoriesResponse)
	resp, err := s.client.Do(ctx, req, searchResp)
	if err != nil {
		return nil, resp, err
	}

	return searchResp, resp, nil
}

// ListPublicDockerHubTags lists tags for a public Docker Hub image.
func (s *ExternalRegistryService) ListPublicDockerHubTags(ctx context.Context, namespace, repository string, opts *DockerHubListOptions) (*DockerHubTagsResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v1/external-registry/dockerhub/%s/%s/tags", namespace, repository)
	var err error
	if opts != nil {
		u, err = addOptions(u, opts)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	tagsResp := new(DockerHubTagsResponse)
	resp, err := s.client.Do(ctx, req, tagsResp)
	if err != nil {
		return nil, resp, err
	}

	return tagsResp, resp, nil
}
