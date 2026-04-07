package pipeops

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserService_GetProfile_UsesProfileDataRoute(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/profile/data" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/profile/data")
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"success": true,
			"message": "Profile data fetched succesfully",
			"data": {
				"email": "jane@example.com",
				"full_name": "Jane Doe",
				"avatar_url": "https://example.com/avatar.png",
				"email_verified": true,
				"password_changed_date": "2024-01-02T15:04:05Z",
				"oauth_user": "google",
				"temp_plan_id": 12,
				"payment_method": true,
				"charge_failed": false,
				"is_subscription_active": true,
				"is_subscription_active_date": "2024-01-03T15:04:05Z",
				"namespace": "jane-dev"
			}
		}`)); err != nil {
			t.Fatalf("write response error: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Users.GetProfile(context.Background())
	if err != nil {
		t.Fatalf("Users.GetProfile error: %v", err)
	}
	if resp.Status != "success" {
		t.Fatalf("Status = %q, want %q", resp.Status, "success")
	}

	got := resp.Data.User
	if got.Email != "jane@example.com" || got.FullName != "Jane Doe" || got.AvatarURL != "https://example.com/avatar.png" || !got.EmailVerified || got.OAuthUser != "google" || got.TempPlanID != 12 || !got.PaymentMethod || got.ChargeFailed || !got.IsSubscriptionActive || got.Namespace != "jane-dev" {
		t.Fatalf("user = %+v, want mapped profile-data fields", got)
	}
	if got.PasswordChangedDate == nil || got.PasswordChangedDate.IsZero() {
		t.Fatalf("PasswordChangedDate = %#v, want parsed timestamp", got.PasswordChangedDate)
	}
	if got.IsSubscriptionActiveDate == nil || got.IsSubscriptionActiveDate.IsZero() {
		t.Fatalf("IsSubscriptionActiveDate = %#v, want parsed timestamp", got.IsSubscriptionActiveDate)
	}
}

func TestUserService_GetProfile_AcceptsWrappedUserPayload(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"status":"success","message":"ok","data":{"user":{"email":"wrapped@example.com","full_name":"Wrapped User"}}}`)); err != nil {
			t.Fatalf("write response error: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, _, err := client.Users.GetProfile(context.Background())
	if err != nil {
		t.Fatalf("Users.GetProfile error: %v", err)
	}
	if got := resp.Data.User.Email; got != "wrapped@example.com" {
		t.Fatalf("email = %q, want %q", got, "wrapped@example.com")
	}
	if got := resp.Data.User.FullName; got != "Wrapped User" {
		t.Fatalf("full_name = %q, want %q", got, "Wrapped User")
	}
}
