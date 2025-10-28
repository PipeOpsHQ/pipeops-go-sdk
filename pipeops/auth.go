package pipeops

import (
	"context"
	"net/http"
)

// AuthService handles communication with the authentication related
// methods of the PipeOps API.
type AuthService struct {
	client *Client
}

// LoginRequest represents a login request.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a login response.
type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
		User  User   `json:"user"`
	} `json:"data"`
}

// User represents a PipeOps user.
type User struct {
	ID            string     `json:"id,omitempty"`
	UUID          string     `json:"uuid,omitempty"`
	Email         string     `json:"email,omitempty"`
	FirstName     string     `json:"first_name,omitempty"`
	LastName      string     `json:"last_name,omitempty"`
	IsActive      bool       `json:"is_active,omitempty"`
	EmailVerified bool       `json:"email_verified,omitempty"`
	CreatedAt     *Timestamp `json:"created_at,omitempty"`
	UpdatedAt     *Timestamp `json:"updated_at,omitempty"`
}

// Login authenticates a user with email and password.
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, *http.Response, error) {
	u := "auth/login"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	loginResp := new(LoginResponse)
	resp, err := s.client.Do(ctx, httpReq, loginResp)
	if err != nil {
		return nil, resp, err
	}

	return loginResp, resp, nil
}

// SignupRequest represents a signup request.
type SignupRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// SignupResponse represents a signup response.
type SignupResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		User User `json:"user"`
	} `json:"data"`
}

// Signup creates a new user account.
func (s *AuthService) Signup(ctx context.Context, req *SignupRequest) (*SignupResponse, *http.Response, error) {
	u := "auth/signup"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	signupResp := new(SignupResponse)
	resp, err := s.client.Do(ctx, httpReq, signupResp)
	if err != nil {
		return nil, resp, err
	}

	return signupResp, resp, nil
}

// PasswordResetRequest represents a password reset request.
type PasswordResetRequest struct {
	Email string `json:"email"`
}

// PasswordResetResponse represents a password reset response.
type PasswordResetResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// RequestPasswordReset sends a password reset email.
func (s *AuthService) RequestPasswordReset(ctx context.Context, req *PasswordResetRequest) (*PasswordResetResponse, *http.Response, error) {
	u := "auth/reset_password/send"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	resetResp := new(PasswordResetResponse)
	resp, err := s.client.Do(ctx, httpReq, resetResp)
	if err != nil {
		return nil, resp, err
	}

	return resetResp, resp, nil
}

// ChangePasswordRequest represents a password change request.
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ChangePasswordResponse represents a password change response.
type ChangePasswordResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ChangePassword changes the user's password.
func (s *AuthService) ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, *http.Response, error) {
	u := "auth/change_password"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	changeResp := new(ChangePasswordResponse)
	resp, err := s.client.Do(ctx, httpReq, changeResp)
	if err != nil {
		return nil, resp, err
	}

	return changeResp, resp, nil
}
