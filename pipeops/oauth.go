package pipeops

import (
	"context"
	"net/http"
	"net/url"
)

// OAuthService handles communication with the OAuth 2.0 related
// methods of the PipeOps API.
type OAuthService struct {
	client *Client
}

// AuthorizeOptions represents OAuth authorization request parameters.
type AuthorizeOptions struct {
	ClientID     string `url:"client_id"`
	RedirectURI  string `url:"redirect_uri"`
	ResponseType string `url:"response_type"` // "code" for authorization code flow
	Scope        string `url:"scope,omitempty"`
	State        string `url:"state,omitempty"`
}

// Authorize initiates the OAuth 2.0 authorization code flow.
// This redirects the user to the authorization endpoint where they can grant access.
// Returns the authorization URL that the user should be redirected to.
func (s *OAuthService) Authorize(opts *AuthorizeOptions) (string, error) {
	u := "oauth/authorize"

	authURL, err := addOptions(u, opts)
	if err != nil {
		return "", err
	}

	// Parse to get full URL
	fullURL, err := s.client.BaseURL.Parse(authURL)
	if err != nil {
		return "", err
	}

	return fullURL.String(), nil
}

// TokenRequest represents an OAuth token exchange request.
type TokenRequest struct {
	GrantType    string `json:"grant_type"`     // "authorization_code" or "refresh_token"
	Code         string `json:"code,omitempty"` // authorization code from callback
	RedirectURI  string `json:"redirect_uri,omitempty"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token,omitempty"` // for refresh token grant
}

// TokenResponse represents an OAuth token response.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// ExchangeCodeForToken exchanges an authorization code for an access token.
func (s *OAuthService) ExchangeCodeForToken(ctx context.Context, req *TokenRequest) (*TokenResponse, *http.Response, error) {
	u := "oauth/token"

	// Build form data
	data := url.Values{}
	data.Set("grant_type", req.GrantType)
	if req.Code != "" {
		data.Set("code", req.Code)
	}
	if req.RedirectURI != "" {
		data.Set("redirect_uri", req.RedirectURI)
	}
	data.Set("client_id", req.ClientID)
	data.Set("client_secret", req.ClientSecret)
	if req.RefreshToken != "" {
		data.Set("refresh_token", req.RefreshToken)
	}

	httpReq, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, nil, err
	}

	// Set form content type and body
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Body = &readCloser{data.Encode()}

	tokenResp := new(TokenResponse)
	resp, err := s.client.Do(ctx, httpReq, tokenResp)
	if err != nil {
		return nil, resp, err
	}

	return tokenResp, resp, nil
}

// UserInfo represents OAuth user information.
type UserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	Name          string `json:"name,omitempty"`
	GivenName     string `json:"given_name,omitempty"`
	FamilyName    string `json:"family_name,omitempty"`
	Picture       string `json:"picture,omitempty"`
}

// UserInfoResponse represents the OAuth userinfo response.
type UserInfoResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Data    UserInfo `json:"data"`
}

// GetUserInfo retrieves user information using an OAuth access token.
// The access token should be set on the client using SetToken().
func (s *OAuthService) GetUserInfo(ctx context.Context) (*UserInfoResponse, *http.Response, error) {
	u := "oauth/userinfo"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	userInfo := new(UserInfoResponse)
	resp, err := s.client.Do(ctx, req, userInfo)
	if err != nil {
		return nil, resp, err
	}

	return userInfo, resp, nil
}

// ConsentResponse represents the OAuth consent page response.
type ConsentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GetConsent retrieves the OAuth consent page (optional endpoint).
func (s *OAuthService) GetConsent(ctx context.Context) (*ConsentResponse, *http.Response, error) {
	u := "oauth/consent"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	consent := new(ConsentResponse)
	resp, err := s.client.Do(ctx, req, consent)
	if err != nil {
		return nil, resp, err
	}

	return consent, resp, nil
}

// Helper type for form body
type readCloser struct {
	content string
}

func (r *readCloser) Read(p []byte) (n int, err error) {
	n = copy(p, r.content)
	r.content = r.content[n:]
	if len(r.content) == 0 {
		err = nil
	}
	return
}

func (r *readCloser) Close() error {
	return nil
}
