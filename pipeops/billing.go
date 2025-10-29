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

// UpdateCard updates a payment card.
func (s *BillingService) UpdateCard(ctx context.Context, cardUUID string, req *AddCardRequest) (*CardResponse, *http.Response, error) {
	u := fmt.Sprintf("billing/workspace/cards%s", cardUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
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

// ListWorkspaceCards lists workspace payment cards.
func (s *BillingService) ListWorkspaceCards(ctx context.Context) (*CardsResponse, *http.Response, error) {
	u := "billing/workspace/cards"

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

// GetActiveCard retrieves the active workspace billing card.
func (s *BillingService) GetActiveCard(ctx context.Context) (*CardResponse, *http.Response, error) {
	u := "billing/workspace/cards/active"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	cardResp := new(CardResponse)
	resp, err := s.client.Do(ctx, req, cardResp)
	if err != nil {
		return nil, resp, err
	}

	return cardResp, resp, nil
}

// GetUsagePlanProviders retrieves usage plan providers.
func (s *BillingService) GetUsagePlanProviders(ctx context.Context) (*http.Response, error) {
	u := "billing/usage-plan-providers"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
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

// BalanceResponse represents account balance response.
type BalanceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Balance  float64 `json:"balance"`
		Currency string  `json:"currency"`
	} `json:"data"`
}

// GetBalance retrieves the current account balance.
func (s *BillingService) GetBalance(ctx context.Context) (*BalanceResponse, *http.Response, error) {
	u := "billing/balance"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	balanceResp := new(BalanceResponse)
	resp, err := s.client.Do(ctx, req, balanceResp)
	if err != nil {
		return nil, resp, err
	}

	return balanceResp, resp, nil
}

// CreditRequest represents a credit add request.
type CreditRequest struct {
	Amount float64 `json:"amount"`
}

// AddCredit adds credit to the account.
func (s *BillingService) AddCredit(ctx context.Context, req *CreditRequest) (*BalanceResponse, *http.Response, error) {
	u := "billing/credit"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	balanceResp := new(BalanceResponse)
	resp, err := s.client.Do(ctx, httpReq, balanceResp)
	if err != nil {
		return nil, resp, err
	}

	return balanceResp, resp, nil
}

// BillingHistoryResponse represents billing history response.
type BillingHistoryResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		History []map[string]interface{} `json:"history"`
	} `json:"data"`
}

// GetHistory retrieves billing history.
func (s *BillingService) GetHistory(ctx context.Context) (*BillingHistoryResponse, *http.Response, error) {
	u := "billing/history"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	historyResp := new(BillingHistoryResponse)
	resp, err := s.client.Do(ctx, req, historyResp)
	if err != nil {
		return nil, resp, err
	}

	return historyResp, resp, nil
}

// SetActiveCardRequest represents a request to set active billing card.
type SetActiveCardRequest struct {
	CardUUID string `json:"card_uuid"`
}

// SetActiveCard sets the active billing card.
func (s *BillingService) SetActiveCard(ctx context.Context, cardUUID string) (*http.Response, error) {
	u := fmt.Sprintf("billing/workspace/cards%s", cardUUID)

	req, err := s.client.NewRequest(http.MethodPut, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// CreateFreeServerRequest represents a request to create a free server.
type CreateFreeServerRequest struct {
	Provider string `json:"provider"`
	Region   string `json:"region"`
}

// CreateFreeServer creates a free trial server.
func (s *BillingService) CreateFreeServer(ctx context.Context, req *CreateFreeServerRequest) (*http.Response, error) {
	u := "billing/create_free_server"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// StartTrialRequest represents a request to start a trial.
type StartTrialRequest struct {
	PlanID string `json:"plan_id"`
}

// StartTrial starts a free trial.
func (s *BillingService) StartTrial(ctx context.Context, req *StartTrialRequest) (*http.Response, error) {
	u := "billing/start-trial"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// PortalResponse represents billing portal response.
type PortalResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		PortalURL string `json:"portal_url"`
	} `json:"data"`
}

// GetPortalURL retrieves the billing portal URL.
func (s *BillingService) GetPortalURL(ctx context.Context) (*PortalResponse, *http.Response, error) {
	u := "billing/portal"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	portalResp := new(PortalResponse)
	resp, err := s.client.Do(ctx, req, portalResp)
	if err != nil {
		return nil, resp, err
	}

	return portalResp, resp, nil
}

// DeploymentQuotaTopupRequest represents a deployment quota topup request.
type DeploymentQuotaTopupRequest struct {
	Amount int `json:"amount"`
}

// DeploymentQuotaTopup adds deployment quota.
func (s *BillingService) DeploymentQuotaTopup(ctx context.Context, req *DeploymentQuotaTopupRequest) (*http.Response, error) {
	u := "billing/deployment-quota/topup"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// GetWorkspaceSubscription retrieves subscription for a workspace.
func (s *BillingService) GetWorkspaceSubscription(ctx context.Context, workspaceUUID string) (*SubscriptionResponse, *http.Response, error) {
	u := fmt.Sprintf("billing/subscriptions/workspace/%s/current", workspaceUUID)

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

// GetTeamSeatSubscription retrieves team seat subscription.
func (s *BillingService) GetTeamSeatSubscription(ctx context.Context) (*SubscriptionResponse, *http.Response, error) {
	u := "billing/subscriptions/workspace/team-seat"

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

// GetCurrentSubscription retrieves the current subscription.
func (s *BillingService) GetCurrentSubscription(ctx context.Context) (*SubscriptionResponse, *http.Response, error) {
	u := "billing/subscriptions/current"

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

// GetPlans retrieves available billing plans.
func (s *BillingService) GetPlans(ctx context.Context) (*PlansResponse, *http.Response, error) {
	u := "billing/plans"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	plansResp := new(PlansResponse)
	resp, err := s.client.Do(ctx, req, plansResp)
	if err != nil {
		return nil, resp, err
	}

	return plansResp, resp, nil
}

// ResetSubscription resets a user's subscription (admin only).
func (s *BillingService) ResetSubscription(ctx context.Context, userUUID string) (*http.Response, error) {
	u := fmt.Sprintf("billing/subscriptions/reset/user/%s", userUUID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// GetWorkspaceCards retrieves cards for a workspace.
func (s *BillingService) GetWorkspaceCards(ctx context.Context) (*CardsResponse, *http.Response, error) {
	u := "billing/workspace/cards"

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

// CreateWorkspaceBilling creates workspace billing configuration.
func (s *BillingService) CreateWorkspaceBilling(ctx context.Context) (*http.Response, error) {
	u := "billing/workspace"

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// RefundRequest represents a refund request.
type RefundRequest struct {
	InvoiceUUID string  `json:"invoice_uuid"`
	Amount      float64 `json:"amount,omitempty"`
	Reason      string  `json:"reason,omitempty"`
}

// ProcessRefund processes a billing refund (admin only).
func (s *BillingService) ProcessRefund(ctx context.Context, req *RefundRequest) (*http.Response, error) {
	u := "billing/refund"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// ApplyDiscount applies a discount to a subscription.
type ApplyDiscountRequest struct {
	SubscriptionUUID string  `json:"subscription_uuid"`
	DiscountPercent  float64 `json:"discount_percent"`
	Duration         int     `json:"duration,omitempty"` // months
}

// ApplyDiscount applies a discount (admin only).
func (s *BillingService) ApplyDiscount(ctx context.Context, req *ApplyDiscountRequest) (*http.Response, error) {
	u := "billing/discount/apply"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// GetBillingReports retrieves billing reports (admin only).
func (s *BillingService) GetBillingReports(ctx context.Context) (*http.Response, error) {
	u := "billing/reports"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// ExportInvoices exports invoices to CSV/PDF.
type ExportInvoicesRequest struct {
	Format    string `json:"format"` // "csv" or "pdf"
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

// ExportInvoices exports invoices.
func (s *BillingService) ExportInvoices(ctx context.Context, req *ExportInvoicesRequest) (*http.Response, error) {
	u := "billing/invoices/export"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// UpdatePaymentMethod updates the default payment method.
type UpdatePaymentMethodRequest struct {
	CardUUID string `json:"card_uuid"`
}

// UpdatePaymentMethod updates payment method.
func (s *BillingService) UpdatePaymentMethod(ctx context.Context, req *UpdatePaymentMethodRequest) (*http.Response, error) {
	u := "billing/payment-method/update"

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}
