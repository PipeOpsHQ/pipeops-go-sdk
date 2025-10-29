// Package pipeops provides a Go client library for the PipeOps Control Plane API.
package pipeops

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL   = "https://api.pipeops.io"
	defaultUserAgent = "pipeops-go-sdk/1.0.0"
	defaultTimeout   = 30 * time.Second
)

// Client manages communication with the PipeOps API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent used when communicating with the PipeOps API.
	UserAgent string

	// Authentication token for API requests.
	token string

	// Services used for talking to different parts of the PipeOps API.
	Auth                *AuthService
	OAuth               *OAuthService
	Projects            *ProjectService
	Servers             *ServerService
	Environments        *EnvironmentService
	Teams               *TeamService
	Workspaces          *WorkspaceService
	Billing             *BillingService
	AddOns              *AddOnService
	Webhooks            *WebhookService
	Users               *UserService
	Admin               *AdminService
	CloudProviders      *CloudProviderService
	Events              *EventService
	Survey              *SurveyService
	Partners            *PartnerService
	Misc                *MiscService
	DeploymentWebhooks  *DeploymentWebhookService
	Campaign            *CampaignService
	Coupons             *CouponService
	Services            *ServiceService
	PartnerAgreements   *PartnerAgreementService
	PartnerParticipants *PartnerParticipantService
	Profile             *ProfileService
	MCPRegistry         *MCPRegistryService
	OpenCost            *OpenCostService
	Notifications       *NotificationService
	Templates           *TemplateService
	Integrations        *IntegrationService
	HealthCheck         *HealthCheckService
	Backups             *BackupService
	SecurityScan        *SecurityScanService
	Logs                *LogService
	AuditLogs           *AuditLogService
	Alerts              *AlertService
	ServiceTokens       *ServiceTokenService
}

// NewClient returns a new PipeOps API client.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		parsedURL, err = url.Parse(defaultBaseURL)
		if err != nil {
			panic(fmt.Sprintf("invalid default base URL %q: %v", defaultBaseURL, err))
		}
	}

	c := &Client{
		client:    &http.Client{Timeout: defaultTimeout},
		BaseURL:   parsedURL,
		UserAgent: defaultUserAgent,
	}

	// Initialize services
	c.Auth = &AuthService{client: c}
	c.OAuth = &OAuthService{client: c}
	c.Projects = &ProjectService{client: c}
	c.Servers = &ServerService{client: c}
	c.Environments = &EnvironmentService{client: c}
	c.Teams = &TeamService{client: c}
	c.Workspaces = &WorkspaceService{client: c}
	c.Billing = &BillingService{client: c}
	c.AddOns = &AddOnService{client: c}
	c.Webhooks = &WebhookService{client: c}
	c.Users = &UserService{client: c}
	c.Admin = &AdminService{client: c}
	c.CloudProviders = &CloudProviderService{client: c}
	c.Events = &EventService{client: c}
	c.Survey = &SurveyService{client: c}
	c.Partners = &PartnerService{client: c}
	c.Misc = &MiscService{client: c}
	c.DeploymentWebhooks = &DeploymentWebhookService{client: c}
	c.Campaign = &CampaignService{client: c}
	c.Coupons = &CouponService{client: c}
	c.Services = &ServiceService{client: c}
	c.PartnerAgreements = &PartnerAgreementService{client: c}
	c.PartnerParticipants = &PartnerParticipantService{client: c}
	c.Profile = &ProfileService{client: c}
	c.MCPRegistry = &MCPRegistryService{client: c}
	c.OpenCost = &OpenCostService{client: c}
	c.Notifications = &NotificationService{client: c}
	c.Templates = &TemplateService{client: c}
	c.Integrations = &IntegrationService{client: c}
	c.HealthCheck = &HealthCheckService{client: c}
	c.Backups = &BackupService{client: c}
	c.SecurityScan = &SecurityScanService{client: c}
	c.Logs = &LogService{client: c}
	c.AuditLogs = &AuditLogService{client: c}
	c.Alerts = &AlertService{client: c}
	c.ServiceTokens = &ServiceTokenService{client: c}

	return c
}

// SetToken sets the authentication token for API requests.
func (c *Client) SetToken(token string) {
	c.token = token
}

// SetHTTPClient sets a custom HTTP client.
func (c *Client) SetHTTPClient(client *http.Client) {
	c.client = client
}

// NewRequest creates an API request.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		c.BaseURL.Path += "/"
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	return req, nil
}

// Do sends an API request and returns the API response.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (resp *http.Response, err error) {
	if ctx == nil {
		return nil, fmt.Errorf("context must be non-nil")
	}

	req = req.WithContext(ctx)

	resp, err = c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if err = CheckResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			if _, copyErr := io.Copy(w, resp.Body); copyErr != nil && err == nil {
				err = copyErr
			}
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return resp, err
}

// ErrorResponse represents an error response from the PipeOps API.
type ErrorResponse struct {
	Response *http.Response
	Message  string `json:"message"`
	Status   string `json:"status"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}

// CheckResponse checks the API response for errors.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		if unmarshalErr := json.Unmarshal(data, errorResponse); unmarshalErr != nil {
			errorResponse.Message = strings.TrimSpace(string(data))
		}
	}

	// Try to get a meaningful error message
	if errorResponse.Message == "" {
		errorResponse.Message = r.Status
	}

	return errorResponse
}

// addOptions adds the parameters in opts as URL query parameters to s.
func addOptions(s string, opts interface{}) (string, error) {
	v, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	u.RawQuery = v.Encode()
	return u.String(), nil
}

// Common types used across the API

// Timestamp represents a time that can be unmarshalled from a JSON string
type Timestamp struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "null" || str == "" {
		return nil
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.999Z07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
	}

	var err error
	for _, format := range formats {
		t.Time, err = time.Parse(format, str)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("could not parse time: %s", str)
}
