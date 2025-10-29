package pipeops

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		wantErr bool
	}{
		{
			name:    "valid URL",
			baseURL: "https://api.pipeops.io",
			wantErr: false,
		},
		{
			name:    "empty URL uses default",
			baseURL: "",
			wantErr: false,
		},
		{
			name:    "invalid URL",
			baseURL: "ht!tp://invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.baseURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	timeout := 10 * time.Second
	client, err := NewClient("", WithTimeout(timeout))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client.client.Timeout != timeout {
		t.Errorf("timeout = %v, want %v", client.client.Timeout, timeout)
	}
}

func TestWithMaxRetries(t *testing.T) {
	maxRetries := 5
	client, err := NewClient("", WithMaxRetries(maxRetries))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client.retryConfig.MaxRetries != maxRetries {
		t.Errorf("maxRetries = %v, want %v", client.retryConfig.MaxRetries, maxRetries)
	}
}

func TestClient_DoWithRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			// Fail first 2 attempts with 503
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		// Succeed on 3rd attempt
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, WithMaxRetries(3))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req, err := client.NewRequest(http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}

	ctx := context.Background()
	resp, err := client.Do(ctx, req, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	if attempts != 3 {
		t.Errorf("attempts = %v, want 3", attempts)
	}
}

func TestClient_DoWithRateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Limit", "100")
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, WithMaxRetries(0))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req, err := client.NewRequest(http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}

	ctx := context.Background()
	_, err = client.Do(ctx, req, nil)
	
	if err == nil {
		t.Fatal("Do() expected error for rate limit")
	}

	rateLimitErr, ok := err.(*RateLimitError)
	if !ok {
		t.Fatalf("expected RateLimitError, got %T", err)
	}

	if rateLimitErr.Limit != 100 {
		t.Errorf("Limit = %v, want 100", rateLimitErr.Limit)
	}
}

func TestClient_SetToken(t *testing.T) {
	client, err := NewClient("")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	token := "test-token-123"
	client.SetToken(token)

	if client.token != token {
		t.Errorf("token = %v, want %v", client.token, token)
	}
}

func TestMustNewClient(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("MustNewClient panicked as expected for invalid URL")
		}
	}()

	// Should not panic
	client := MustNewClient("")
	if client == nil {
		t.Error("MustNewClient() returned nil")
	}

	// Should panic
	_ = MustNewClient("ht!tp://invalid")
	t.Error("MustNewClient() should have panicked")
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "no @ symbol",
			email:   "userexample.com",
			wantErr: true,
		},
		{
			name:    "no domain",
			email:   "user@",
			wantErr: true,
		},
		{
			name:    "no username",
			email:   "@example.com",
			wantErr: true,
		},
		{
			name:    "with spaces",
			email:   "  user@example.com  ",
			wantErr: false, // Should be trimmed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
