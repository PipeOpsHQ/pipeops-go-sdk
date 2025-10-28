package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// WebhookService handles communication with the webhook related
// methods of the PipeOps API.
type WebhookService struct {
	client *Client
}

// Webhook represents a webhook configuration.
type Webhook struct {
	ID          string     `json:"id,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	URL         string     `json:"url,omitempty"`
	Events      []string   `json:"events,omitempty"`
	Secret      string     `json:"secret,omitempty"`
	Active      bool       `json:"active,omitempty"`
	Description string     `json:"description,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`
}

// WebhooksResponse represents a list of webhooks response.
type WebhooksResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Webhooks []Webhook `json:"webhooks"`
	} `json:"data"`
}

// CreateWebhookRequest represents a request to create a webhook.
type CreateWebhookRequest struct {
	URL         string   `json:"url"`
	Events      []string `json:"events"`
	Secret      string   `json:"secret,omitempty"`
	Description string   `json:"description,omitempty"`
}

// Create creates a new webhook.
func (s *WebhookService) Create(ctx context.Context, req *CreateWebhookRequest) (*WebhookResponse, *http.Response, error) {
	u := "webhook/customer/webhook/create"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	webhookResp := new(WebhookResponse)
	resp, err := s.client.Do(ctx, httpReq, webhookResp)
	if err != nil {
		return nil, resp, err
	}

	return webhookResp, resp, nil
}

// List lists all webhooks.
func (s *WebhookService) List(ctx context.Context) (*WebhooksResponse, *http.Response, error) {
	u := "webhook/customer/webhooks"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	webhooksResp := new(WebhooksResponse)
	resp, err := s.client.Do(ctx, req, webhooksResp)
	if err != nil {
		return nil, resp, err
	}

	return webhooksResp, resp, nil
}

// Get fetches a webhook by UUID.
func (s *WebhookService) Get(ctx context.Context, webhookUUID string) (*WebhookResponse, *http.Response, error) {
	u := fmt.Sprintf("webhook/customer/webhook/%s", webhookUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
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

// UpdateWebhookRequest represents a request to update a webhook.
type UpdateWebhookRequest struct {
	URL         string   `json:"url,omitempty"`
	Events      []string `json:"events,omitempty"`
	Active      *bool    `json:"active,omitempty"`
	Description string   `json:"description,omitempty"`
}

// Update updates a webhook.
func (s *WebhookService) Update(ctx context.Context, webhookUUID string, req *UpdateWebhookRequest) (*WebhookResponse, *http.Response, error) {
	u := fmt.Sprintf("webhook/customer/webhook/%s", webhookUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	webhookResp := new(WebhookResponse)
	resp, err := s.client.Do(ctx, httpReq, webhookResp)
	if err != nil {
		return nil, resp, err
	}

	return webhookResp, resp, nil
}

// Delete deletes a webhook.
func (s *WebhookService) Delete(ctx context.Context, webhookUUID string) (*http.Response, error) {
	u := fmt.Sprintf("webhook/customer/webhook/%s", webhookUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}
