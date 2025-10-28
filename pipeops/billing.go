package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// BillingService handles communication with the billing related
// methods of the PipeOps API.
type BillingService struct {
	client *Client
}

// Card represents a billing card.
type Card struct {
	ID        string     `json:"id,omitempty"`
	UUID      string     `json:"uuid,omitempty"`
	Last4     string     `json:"last4,omitempty"`
	Brand     string     `json:"brand,omitempty"`
	ExpMonth  int        `json:"exp_month,omitempty"`
	ExpYear   int        `json:"exp_year,omitempty"`
	IsDefault bool       `json:"is_default,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
}

// CardsResponse represents a list of cards response.
type CardsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Cards []Card `json:"cards"`
	} `json:"data"`
}

// CardResponse represents a single card response.
type CardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Card Card `json:"card"`
	} `json:"data"`
}

// AddCardRequest represents a request to add a payment card.
type AddCardRequest struct {
	Token string `json:"token"` // Payment provider token
}

// AddCard adds a new payment card.
func (s *BillingService) AddCard(ctx context.Context, req *AddCardRequest) (*CardResponse, *http.Response, error) {
	u := "billing/cards"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	cardResp := new(CardResponse)
	resp, err := s.client.Do(ctx, httpReq, cardResp)
	if err != nil {
		return nil, resp, err
	}

	return cardResp, resp, nil
}

// ListCards lists all payment cards.
func (s *BillingService) ListCards(ctx context.Context) (*CardsResponse, *http.Response, error) {
	u := "billing/cards"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	cardsResp := new(CardsResponse)
	resp, err := s.client.Do(ctx, req, cardsResp)
	if err != nil {
		return nil, resp, err
	}

	return cardsResp, resp, nil
}

// DeleteCard deletes a payment card.
func (s *BillingService) DeleteCard(ctx context.Context, cardUUID string) (*http.Response, error) {
	u := fmt.Sprintf("billing/cards/%s", cardUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// Subscription represents a billing subscription.
type Subscription struct {
	ID        string     `json:"id,omitempty"`
	UUID      string     `json:"uuid,omitempty"`
	PlanID    string     `json:"plan_id,omitempty"`
	PlanName  string     `json:"plan_name,omitempty"`
	Status    string     `json:"status,omitempty"`
	StartDate *Timestamp `json:"start_date,omitempty"`
	EndDate   *Timestamp `json:"end_date,omitempty"`
	Amount    float64    `json:"amount,omitempty"`
	Currency  string     `json:"currency,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
	UpdatedAt *Timestamp `json:"updated_at,omitempty"`
}

// SubscriptionsResponse represents a list of subscriptions response.
type SubscriptionsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Subscriptions []Subscription `json:"subscriptions"`
	} `json:"data"`
}

// SubscriptionResponse represents a single subscription response.
type SubscriptionResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Subscription Subscription `json:"subscription"`
	} `json:"data"`
}

// SubscribeRequest represents a subscription request.
type SubscribeRequest struct {
	PlanID string `json:"plan_id"`
}

// Subscribe creates a new subscription.
func (s *BillingService) Subscribe(ctx context.Context, req *SubscribeRequest) (*SubscriptionResponse, *http.Response, error) {
	u := "billing/subscribe"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	subResp := new(SubscriptionResponse)
	resp, err := s.client.Do(ctx, httpReq, subResp)
	if err != nil {
		return nil, resp, err
	}

	return subResp, resp, nil
}

// ListSubscriptions lists all subscriptions.
func (s *BillingService) ListSubscriptions(ctx context.Context) (*SubscriptionsResponse, *http.Response, error) {
	u := "billing/subscriptions"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	subsResp := new(SubscriptionsResponse)
	resp, err := s.client.Do(ctx, req, subsResp)
	if err != nil {
		return nil, resp, err
	}

	return subsResp, resp, nil
}

// GetSubscription gets a subscription by UUID.
func (s *BillingService) GetSubscription(ctx context.Context, subscriptionUUID string) (*SubscriptionResponse, *http.Response, error) {
	u := fmt.Sprintf("billing/subscriptions/%s", subscriptionUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	subResp := new(SubscriptionResponse)
	resp, err := s.client.Do(ctx, req, subResp)
	if err != nil {
		return nil, resp, err
	}

	return subResp, resp, nil
}

// CancelSubscription cancels a subscription.
func (s *BillingService) CancelSubscription(ctx context.Context, subscriptionUUID string) (*http.Response, error) {
	u := fmt.Sprintf("billing/subscriptions/%s/cancel", subscriptionUUID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// Invoice represents a billing invoice.
type Invoice struct {
	ID        string     `json:"id,omitempty"`
	UUID      string     `json:"uuid,omitempty"`
	Amount    float64    `json:"amount,omitempty"`
	Currency  string     `json:"currency,omitempty"`
	Status    string     `json:"status,omitempty"`
	DueDate   *Timestamp `json:"due_date,omitempty"`
	PaidAt    *Timestamp `json:"paid_at,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
}

// InvoicesResponse represents a list of invoices response.
type InvoicesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Invoices []Invoice `json:"invoices"`
	} `json:"data"`
}

// ListInvoices lists all invoices.
func (s *BillingService) ListInvoices(ctx context.Context) (*InvoicesResponse, *http.Response, error) {
	u := "billing/invoices"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	invoicesResp := new(InvoicesResponse)
	resp, err := s.client.Do(ctx, req, invoicesResp)
	if err != nil {
		return nil, resp, err
	}

	return invoicesResp, resp, nil
}

// Usage represents billing usage information.
type Usage struct {
	ResourceType string  `json:"resource_type,omitempty"`
	Amount       float64 `json:"amount,omitempty"`
	Unit         string  `json:"unit,omitempty"`
	Cost         float64 `json:"cost,omitempty"`
	Period       string  `json:"period,omitempty"`
}

// UsageResponse represents usage information response.
type UsageResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Usage []Usage `json:"usage"`
		Total float64 `json:"total,omitempty"`
	} `json:"data"`
}

// GetUsage gets current billing usage.
func (s *BillingService) GetUsage(ctx context.Context) (*UsageResponse, *http.Response, error) {
	u := "billing/usage"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	usageResp := new(UsageResponse)
	resp, err := s.client.Do(ctx, req, usageResp)
	if err != nil {
		return nil, resp, err
	}

	return usageResp, resp, nil
}
