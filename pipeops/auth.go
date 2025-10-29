package pipeops

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	if req == nil {
		return nil, nil, errors.New("login request cannot be nil")
	}
	if err := validateEmail(req.Email); err != nil {
		return nil, nil, fmt.Errorf("invalid email: %w", err)
	}
	if req.Password == "" {
		return nil, nil, errors.New("password cannot be empty")
	}

	u := "auth/login"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create login request: %w", err)
	}

	loginResp := new(LoginResponse)
	resp, err := s.client.Do(ctx, httpReq, loginResp)
	if err != nil {
		return nil, resp, fmt.Errorf("login failed: %w", err)
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
	if req == nil {
		return nil, nil, errors.New("signup request cannot be nil")
	}
	if err := validateEmail(req.Email); err != nil {
		return nil, nil, fmt.Errorf("invalid email: %w", err)
	}
	if len(req.Password) < 6 {
		return nil, nil, errors.New("password must be at least 6 characters")
	}

	u := "auth/signup"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create signup request: %w", err)
	}

	signupResp := new(SignupResponse)
	resp, err := s.client.Do(ctx, httpReq, signupResp)
	if err != nil {
		return nil, resp, fmt.Errorf("signup failed: %w", err)
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

// VerifyLoginRequest represents a login verification request.
type VerifyLoginRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// VerifyLogin verifies a login with 2FA code.
func (s *AuthService) VerifyLogin(ctx context.Context, req *VerifyLoginRequest) (*LoginResponse, *http.Response, error) {
	u := "auth/verify_login"

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

// ActivateEmailRequest represents an email activation request.
type ActivateEmailRequest struct {
	Token string `json:"token"`
}

// ActivateEmail activates a user's email.
func (s *AuthService) ActivateEmail(ctx context.Context, req *ActivateEmailRequest) (*http.Response, error) {
	u := "auth/activate_email"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// OAuthSignup initiates OAuth signup with a provider.
func (s *AuthService) OAuthSignup(ctx context.Context, provider string) (*http.Response, error) {
	u := "auth/" + provider + "/signup"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// OAuthCallback handles OAuth callback.
func (s *AuthService) OAuthCallback(ctx context.Context, provider string) (*LoginResponse, *http.Response, error) {
	u := "auth/" + provider + "/callback"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	loginResp := new(LoginResponse)
	resp, err := s.client.Do(ctx, req, loginResp)
	if err != nil {
		return nil, resp, err
	}

	return loginResp, resp, nil
}

// ResetPasswordRequest represents a password reset request.
type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// ResetPassword resets password with a token.
func (s *AuthService) ResetPassword(ctx context.Context, req *ResetPasswordRequest) (*http.Response, error) {
	u := "auth/reset_password"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// VerifyPasswordResetToken verifies a password reset token.
func (s *AuthService) VerifyPasswordResetToken(ctx context.Context, token string) (*http.Response, error) {
	u := "auth/reset_password/verify/" + token

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// validateEmail performs basic email validation.
func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("email must contain @")
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return errors.New("invalid email format")
	}
	return nil
}
