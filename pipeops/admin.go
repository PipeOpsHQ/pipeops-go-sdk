package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// AdminService handles communication with the admin related
// methods of the PipeOps API.
type AdminService struct {
	client *Client
}

// AdminUser represents a user in admin context.
type AdminUser struct {
	ID            string     `json:"id,omitempty"`
	UUID          string     `json:"uuid,omitempty"`
	Email         string     `json:"email,omitempty"`
	FirstName     string     `json:"first_name,omitempty"`
	LastName      string     `json:"last_name,omitempty"`
	IsActive      bool       `json:"is_active,omitempty"`
	EmailVerified bool       `json:"email_verified,omitempty"`
	Role          string     `json:"role,omitempty"`
	CreatedAt     *Timestamp `json:"created_at,omitempty"`
	UpdatedAt     *Timestamp `json:"updated_at,omitempty"`
}

// AdminUsersResponse represents a list of users in admin context.
type AdminUsersResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Users []AdminUser `json:"users"`
		Total int         `json:"total,omitempty"`
	} `json:"data"`
}

// AdminUserResponse represents a single user in admin context.
type AdminUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		User AdminUser `json:"user"`
	} `json:"data"`
}

// AdminListUsersOptions specifies options for listing users.
type AdminListUsersOptions struct {
	Page   int    `url:"page,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Role   string `url:"role,omitempty"`
	Active *bool  `url:"active,omitempty"`
}

// ListUsers lists all users (admin only).
func (s *AdminService) ListUsers(ctx context.Context, opts *AdminListUsersOptions) (*AdminUsersResponse, *http.Response, error) {
	u := "admin/users"
	if opts != nil {
		var err error
		u, err = addOptions(u, opts)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	usersResp := new(AdminUsersResponse)
	resp, err := s.client.Do(ctx, req, usersResp)
	if err != nil {
		return nil, resp, err
	}

	return usersResp, resp, nil
}

// GetUser gets a user by UUID (admin only).
func (s *AdminService) GetUser(ctx context.Context, userUUID string) (*AdminUserResponse, *http.Response, error) {
	u := fmt.Sprintf("admin/users/%s", userUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	userResp := new(AdminUserResponse)
	resp, err := s.client.Do(ctx, req, userResp)
	if err != nil {
		return nil, resp, err
	}

	return userResp, resp, nil
}

// UpdateUserRequest represents a request to update a user (admin).
type UpdateUserRequest struct {
	IsActive      *bool  `json:"is_active,omitempty"`
	EmailVerified *bool  `json:"email_verified,omitempty"`
	Role          string `json:"role,omitempty"`
}

// UpdateUser updates a user (admin only).
func (s *AdminService) UpdateUser(ctx context.Context, userUUID string, req *UpdateUserRequest) (*AdminUserResponse, *http.Response, error) {
	u := fmt.Sprintf("admin/users/%s", userUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	userResp := new(AdminUserResponse)
	resp, err := s.client.Do(ctx, httpReq, userResp)
	if err != nil {
		return nil, resp, err
	}

	return userResp, resp, nil
}

// DeleteUser deletes a user (admin only).
func (s *AdminService) DeleteUser(ctx context.Context, userUUID string) (*http.Response, error) {
	u := fmt.Sprintf("admin/users/%s", userUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// AdminStats represents admin dashboard statistics.
type AdminStats struct {
	TotalUsers       int     `json:"total_users,omitempty"`
	ActiveUsers      int     `json:"active_users,omitempty"`
	TotalProjects    int     `json:"total_projects,omitempty"`
	TotalDeployments int     `json:"total_deployments,omitempty"`
	TotalRevenue     float64 `json:"total_revenue,omitempty"`
	MonthlyRevenue   float64 `json:"monthly_revenue,omitempty"`
}

// AdminStatsResponse represents admin statistics response.
type AdminStatsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Stats AdminStats `json:"stats"`
	} `json:"data"`
}

// GetStats retrieves admin dashboard statistics.
func (s *AdminService) GetStats(ctx context.Context) (*AdminStatsResponse, *http.Response, error) {
	u := "admin/stats"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	statsResp := new(AdminStatsResponse)
	resp, err := s.client.Do(ctx, req, statsResp)
	if err != nil {
		return nil, resp, err
	}

	return statsResp, resp, nil
}

// Plan represents a subscription plan.
type Plan struct {
	ID          string   `json:"id,omitempty"`
	UUID        string   `json:"uuid,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Price       float64  `json:"price,omitempty"`
	Currency    string   `json:"currency,omitempty"`
	Interval    string   `json:"interval,omitempty"`
	Features    []string `json:"features,omitempty"`
	Active      bool     `json:"active,omitempty"`
}

// PlansResponse represents a list of plans response.
type PlansResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Plans []Plan `json:"plans"`
	} `json:"data"`
}

// ListPlans lists all subscription plans (admin only).
func (s *AdminService) ListPlans(ctx context.Context) (*PlansResponse, *http.Response, error) {
	u := "admin/plans"

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

// WaitlistProgramRequest represents a request to create a waitlist program.
type WaitlistProgramRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// WaitlistProgramResponse represents waitlist program response.
type WaitlistProgramResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Program map[string]interface{} `json:"program"`
	} `json:"data"`
}

// CreateWaitlistProgram creates a waitlist program (admin only).
func (s *AdminService) CreateWaitlistProgram(ctx context.Context, req *WaitlistProgramRequest) (*WaitlistProgramResponse, *http.Response, error) {
	u := "admin/create/waitlist-program"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	programResp := new(WaitlistProgramResponse)
	resp, err := s.client.Do(ctx, httpReq, programResp)
	if err != nil {
		return nil, resp, err
	}

	return programResp, resp, nil
}

// BulkAddWaitlistRequest represents a request to bulk add users to waitlist.
type BulkAddWaitlistRequest struct {
	Emails     []string `json:"emails"`
	ProgramUID string   `json:"program_uid"`
}

// BulkAddToWaitlist adds multiple users to a waitlist (admin only).
func (s *AdminService) BulkAddToWaitlist(ctx context.Context, req *BulkAddWaitlistRequest) (*http.Response, error) {
	u := "admin/bulk-add-waitlist"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// BulkRemoveFromWaitlist removes multiple users from a waitlist (admin only).
func (s *AdminService) BulkRemoveFromWaitlist(ctx context.Context, req *BulkAddWaitlistRequest) (*http.Response, error) {
	u := "admin/bulk-add-waitlist"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// BulkExtendSubscriptionRequest represents a request to extend subscriptions.
type BulkExtendSubscriptionRequest struct {
	UserUIDs []string `json:"user_uids"`
	Days     int      `json:"days"`
}

// BulkExtendSubscription extends subscriptions for multiple users (admin only).
func (s *AdminService) BulkExtendSubscription(ctx context.Context, req *BulkExtendSubscriptionRequest) (*http.Response, error) {
	u := "admin/bulk-waitlist/extend-subscription"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// ProgramParticipantsResponse represents program participants response.
type ProgramParticipantsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Participants []map[string]interface{} `json:"participants"`
	} `json:"data"`
}

// GetProgramParticipants gets all participants in a program (admin only).
func (s *AdminService) GetProgramParticipants(ctx context.Context, programUID string) (*ProgramParticipantsResponse, *http.Response, error) {
	u := fmt.Sprintf("admin/programs/%s/participants", programUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	participantsResp := new(ProgramParticipantsResponse)
	resp, err := s.client.Do(ctx, req, participantsResp)
	if err != nil {
		return nil, resp, err
	}

	return participantsResp, resp, nil
}

// CreatePlanRequest represents a request to create a plan.
type CreatePlanRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Price       float64  `json:"price"`
	Currency    string   `json:"currency"`
	Interval    string   `json:"interval"`
	Features    []string `json:"features,omitempty"`
}

// CreatePlan creates a new plan (admin only).
func (s *AdminService) CreatePlan(ctx context.Context, req *CreatePlanRequest) (*http.Response, error) {
	u := "admin/plans/create"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// UpdatePlanRequest represents a request to update a plan.
type UpdatePlanRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Price       float64  `json:"price,omitempty"`
	Features    []string `json:"features,omitempty"`
	Active      *bool    `json:"active,omitempty"`
}

// UpdatePlan updates a plan (admin only).
func (s *AdminService) UpdatePlan(ctx context.Context, planUUID string, req *UpdatePlanRequest) (*http.Response, error) {
	u := fmt.Sprintf("admin/plans/update/%s", planUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// DeletePlan deletes a plan (admin only).
func (s *AdminService) DeletePlan(ctx context.Context, planID string) (*http.Response, error) {
	u := fmt.Sprintf("admin/plans/delete/%s", planID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// SubscribeUserRequest represents a request to subscribe a user.
type SubscribeUserRequest struct {
	UserUUID string `json:"user_uuid"`
	PlanUUID string `json:"plan_uuid"`
}

// SubscribeUser creates a subscription for a user (admin only).
func (s *AdminService) SubscribeUser(ctx context.Context, req *SubscribeUserRequest) (*http.Response, error) {
	u := "admin/billing/subscribe/user"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// PauseSubscriptionRequest represents a request to pause a subscription.
type PauseSubscriptionRequest struct {
	Reason string `json:"reason,omitempty"`
}

// PauseSubscription pauses a subscription (admin only).
func (s *AdminService) PauseSubscription(ctx context.Context, subID string, req *PauseSubscriptionRequest) (*http.Response, error) {
	u := fmt.Sprintf("admin/subscriptions/%s/pause", subID)

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// GetAuditLogs retrieves audit logs (admin only).
func (s *AdminService) GetAuditLogs(ctx context.Context) (*http.Response, error) {
	u := "admin/audit-logs"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// GetSystemHealth retrieves system health status (admin only).
func (s *AdminService) GetSystemHealth(ctx context.Context) (*http.Response, error) {
	u := "admin/system/health"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// BroadcastRequest represents a broadcast message request.
type BroadcastRequest struct {
	Message string   `json:"message"`
	Users   []string `json:"users,omitempty"`
}

// BroadcastMessage broadcasts a message to users (admin only).
func (s *AdminService) BroadcastMessage(ctx context.Context, req *BroadcastRequest) (*http.Response, error) {
	u := "admin/broadcast"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// ImpersonateUserRequest represents a user impersonation request.
type ImpersonateUserRequest struct {
	UserUUID string `json:"user_uuid"`
}

// ImpersonateUser impersonates a user (admin only).
func (s *AdminService) ImpersonateUser(ctx context.Context, req *ImpersonateUserRequest) (*http.Response, error) {
	u := "admin/impersonate"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}
