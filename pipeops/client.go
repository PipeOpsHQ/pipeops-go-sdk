// Package pipeops provides a Go client library for the PipeOps Control Plane API.
package pipeops

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL      = "https://api.pipeops.io"
	defaultUserAgent    = "pipeops-go-sdk/1.0.0"
	defaultTimeout      = 30 * time.Second
	defaultMaxRetries   = 3
	defaultRetryWaitMin = 100 * time.Millisecond
	defaultRetryWaitMax = 5 * time.Second
)

// RetryConfig configures retry behavior for failed requests.
type RetryConfig struct {
	MaxRetries   int
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
	RetryPolicy  RetryPolicy
}

// RetryPolicy determines if a request should be retried.
type RetryPolicy func(ctx context.Context, resp *http.Response, err error) (bool, error)

// Logger is an interface for logging SDK operations.
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// defaultLogger is a no-op logger.
type defaultLogger struct{}

func (l *defaultLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (l *defaultLogger) Info(msg string, keysAndValues ...interface{})  {}
func (l *defaultLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (l *defaultLogger) Error(msg string, keysAndValues ...interface{}) {}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client) error

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

	// Retry configuration
	retryConfig *RetryConfig

	// Logger for debug output
	logger Logger

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

// newHTTPClient creates a properly configured HTTP client for production use.
func newHTTPClient(timeout time.Duration) *http.Client {
	transport := &http.Transport{
		// Connection pool settings for high concurrency
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
		IdleConnTimeout:     90 * time.Second,

		// Timeouts for different phases
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

		// Enable HTTP/2
		ForceAttemptHTTP2: true,
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
	}
}

// defaultRetryPolicy implements a sensible retry policy for HTTP requests.
func defaultRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// Don't retry if context is done
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	// Retry on network errors
	if err != nil {
		return true, nil
	}

	// Retry on specific status codes
	if resp.StatusCode == 0 || resp.StatusCode == 429 || resp.StatusCode >= 500 {
		return true, nil
	}

	return false, nil
}

// NewClient returns a new PipeOps API client with optional configuration.
// If baseURL is empty, the default API URL is used.
// Returns an error if the provided baseURL is invalid.
func NewClient(baseURL string, opts ...ClientOption) (*Client, error) {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL %q: %w", baseURL, err)
	}

	// Ensure base URL has trailing slash
	if !strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path += "/"
	}

	c := &Client{
		client:    newHTTPClient(defaultTimeout),
		BaseURL:   parsedURL,
		UserAgent: defaultUserAgent,
		retryConfig: &RetryConfig{
			MaxRetries:   defaultMaxRetries,
			RetryWaitMin: defaultRetryWaitMin,
			RetryWaitMax: defaultRetryWaitMax,
			RetryPolicy:  defaultRetryPolicy,
		},
		logger: &defaultLogger{},
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
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

	return c, nil
}

// ClientOption functions

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) error {
		if client == nil {
			return errors.New("HTTP client cannot be nil")
		}
		c.client = client
		return nil
	}
}

// WithTimeout sets the timeout for API requests.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		if timeout <= 0 {
			return errors.New("timeout must be positive")
		}
		c.client = newHTTPClient(timeout)
		return nil
	}
}

// WithRetryConfig sets custom retry configuration.
func WithRetryConfig(config *RetryConfig) ClientOption {
	return func(c *Client) error {
		if config == nil {
			return errors.New("retry config cannot be nil")
		}
		if config.MaxRetries < 0 {
			return errors.New("max retries must be non-negative")
		}
		c.retryConfig = config
		return nil
	}
}

// WithMaxRetries sets the maximum number of retry attempts.
func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *Client) error {
		if maxRetries < 0 {
			return errors.New("max retries must be non-negative")
		}
		c.retryConfig.MaxRetries = maxRetries
		return nil
	}
}

// WithLogger sets a custom logger for the client.
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) error {
		if logger == nil {
			return errors.New("logger cannot be nil")
		}
		c.logger = logger
		return nil
	}
}

// WithUserAgent sets a custom user agent string.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		if userAgent == "" {
			return errors.New("user agent cannot be empty")
		}
		c.UserAgent = userAgent
		return nil
	}
}

// MustNewClient returns a new PipeOps API client and panics on error.
// This should only be used in init functions or when you are certain the URL is valid.
func MustNewClient(baseURL string, opts ...ClientOption) *Client {
	client, err := NewClient(baseURL, opts...)
	if err != nil {
		panic(err)
	}
	return client
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
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL %q: %w", urlStr, err)
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

// Do sends an API request and returns the API response with automatic retry logic.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context must be non-nil")
	}

	var resp *http.Response
	var err error

	// Retry loop
	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			// Calculate backoff delay with jitter
			waitDuration := c.calculateBackoff(attempt)

			safeURL := strings.ReplaceAll(req.URL.String(), "\n", "")
			safeURL = strings.ReplaceAll(safeURL, "\r", "")
			c.logger.Warn("Retrying request",
				"attempt", attempt,
				"max_attempts", c.retryConfig.MaxRetries,
				"wait_duration", waitDuration,
				"method", req.Method,
				"url", safeURL,
			)

			// Wait before retry, respecting context cancellation
			select {
			case <-time.After(waitDuration):
				// Continue with retry
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		// Clone request for retry (important for request body)
		reqClone := req.Clone(ctx)

		// Make the request
		resp, err = c.client.Do(reqClone)

		// Check if we should retry
		shouldRetry, checkErr := c.retryConfig.RetryPolicy(ctx, resp, err)

		if checkErr != nil {
			return nil, checkErr
		}

		if !shouldRetry {
			break
		}

		// If we have a response, drain and close the body before retrying
		if resp != nil {
			//nolint:errcheck // Best effort drain before retry
			io.Copy(io.Discard, resp.Body)
			//nolint:errcheck // Best effort close before retry
			safeURL := strings.ReplaceAll(req.URL.String(), "\n", "")
			safeURL = strings.ReplaceAll(safeURL, "\r", "")
			resp.Body.Close()
		}

		// Don't retry if this was the last attempt
		if attempt == c.retryConfig.MaxRetries {
			c.logger.Error("Max retries exceeded",
				"attempts", attempt+1,
				"method", req.Method,
				"url", safeURL,
			)
			break
		}
	}

	// Handle request error
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, fmt.Errorf("request failed after %d attempts: %w", c.retryConfig.MaxRetries+1, err)
	}

	// Defer closing response body
	defer func() {
		// Drain and close body to enable connection reuse
		//nolint:errcheck // Best effort drain for connection reuse
		io.Copy(io.Discard, resp.Body)
		//nolint:errcheck // Best effort close for connection reuse
		resp.Body.Close()
	}()

	// Check for API errors
	if err = CheckResponse(resp); err != nil {
		return resp, err
	}

	// Decode response if needed
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			if _, copyErr := io.Copy(w, resp.Body); copyErr != nil {
				return resp, fmt.Errorf("failed to copy response: %w", copyErr)
			}
		} else {
			// Read body for decoding
			bodyBytes, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				return resp, fmt.Errorf("failed to read response body: %w", readErr)
			}

			// Decode if we have content
			if len(bodyBytes) > 0 {
				if decErr := json.Unmarshal(bodyBytes, v); decErr != nil {
					return resp, fmt.Errorf("failed to decode response: %w", decErr)
				}
			}
		}
	}

	return resp, nil
}

// calculateBackoff calculates the backoff duration with exponential backoff and jitter.
func (c *Client) calculateBackoff(attempt int) time.Duration {
	// Exponential backoff: min * 2^(attempt-1)
	backoff := float64(c.retryConfig.RetryWaitMin) * math.Pow(2, float64(attempt-1))

	// Cap at max wait time
	if backoff > float64(c.retryConfig.RetryWaitMax) {
		backoff = float64(c.retryConfig.RetryWaitMax)
	}

	// Add jitter (±10% randomness to prevent thundering herd)
	jitter := backoff * 0.1 * (rand.Float64()*2 - 1)
	backoff += jitter

	return time.Duration(backoff)
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

	// Handle rate limiting specially
	if r.StatusCode == 429 {
		return parseRateLimitError(r)
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

// RateLimitError represents a rate limit error from the API.
type RateLimitError struct {
	Response   *http.Response
	RetryAfter time.Duration
	Limit      int
	Remaining  int
	Reset      time.Time
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded: retry after %v (limit: %d, remaining: %d)",
		e.RetryAfter, e.Limit, e.Remaining)
}

// parseRateLimitError parses rate limit information from response headers.
func parseRateLimitError(r *http.Response) *RateLimitError {
	err := &RateLimitError{
		Response: r,
	}

	// Parse Retry-After header (seconds or HTTP date)
	if retryAfter := r.Header.Get("Retry-After"); retryAfter != "" {
		if seconds, parseErr := time.ParseDuration(retryAfter + "s"); parseErr == nil {
			err.RetryAfter = seconds
		}
	}

	// Parse rate limit headers if available
	if limit := r.Header.Get("X-RateLimit-Limit"); limit != "" {
		//nolint:errcheck // Best effort parse, defaults used if parse fails
		fmt.Sscanf(limit, "%d", &err.Limit)
	}
	if remaining := r.Header.Get("X-RateLimit-Remaining"); remaining != "" {
		//nolint:errcheck // Best effort parse, defaults used if parse fails
		fmt.Sscanf(remaining, "%d", &err.Remaining)
	}

	// Default retry after if not specified
	if err.RetryAfter == 0 {
		err.RetryAfter = 60 * time.Second
	}

	return err
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
