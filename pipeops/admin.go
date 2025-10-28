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
		u, _ = addOptions(u, opts)
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
