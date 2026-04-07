package pipeops

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// UserService handles communication with the user settings related
// methods of the PipeOps API.
type UserService struct {
	client *Client
}

// UserSettings represents user settings.
type UserSettings struct {
	Notifications NotificationSettings `json:"notifications,omitempty"`
	Preferences   UserPreferences      `json:"preferences,omitempty"`
}

// NotificationSettings represents notification preferences.
type NotificationSettings struct {
	Email       bool `json:"email,omitempty"`
	Push        bool `json:"push,omitempty"`
	Deployments bool `json:"deployments,omitempty"`
	Billing     bool `json:"billing,omitempty"`
	Security    bool `json:"security,omitempty"`
}

// UserPreferences represents user preferences.
type UserPreferences struct {
	Theme    string `json:"theme,omitempty"`
	Language string `json:"language,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

// UserSettingsResponse represents user settings response.
type UserSettingsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Settings UserSettings `json:"settings"`
	} `json:"data"`
}

// GetSettings retrieves user settings.
func (s *UserService) GetSettings(ctx context.Context) (*UserSettingsResponse, *http.Response, error) {
	u := "user/settings"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	settingsResp := new(UserSettingsResponse)
	resp, err := s.client.Do(ctx, req, settingsResp)
	if err != nil {
		return nil, resp, err
	}

	return settingsResp, resp, nil
}

// UpdateSettingsRequest represents a request to update user settings.
type UpdateSettingsRequest struct {
	Notifications *NotificationSettings `json:"notifications,omitempty"`
	Preferences   *UserPreferences      `json:"preferences,omitempty"`
}

// UpdateSettings updates user settings.
func (s *UserService) UpdateSettings(ctx context.Context, req *UpdateSettingsRequest) (*UserSettingsResponse, *http.Response, error) {
	u := "user/settings"

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	settingsResp := new(UserSettingsResponse)
	resp, err := s.client.Do(ctx, httpReq, settingsResp)
	if err != nil {
		return nil, resp, err
	}

	return settingsResp, resp, nil
}

// UpdateNotificationSettingsRequest represents a request to update notification settings.
type UpdateNotificationSettingsRequest struct {
	Email       *bool `json:"email,omitempty"`
	Push        *bool `json:"push,omitempty"`
	Deployments *bool `json:"deployments,omitempty"`
	Billing     *bool `json:"billing,omitempty"`
	Security    *bool `json:"security,omitempty"`
}

// UpdateNotificationSettings updates notification settings.
func (s *UserService) UpdateNotificationSettings(ctx context.Context, req *UpdateNotificationSettingsRequest) (*UserSettingsResponse, *http.Response, error) {
	u := "user-settings/notification"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	settingsResp := new(UserSettingsResponse)
	resp, err := s.client.Do(ctx, httpReq, settingsResp)
	if err != nil {
		return nil, resp, err
	}

	return settingsResp, resp, nil
}

// ProfileResponse represents user profile response.
type ProfileResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		User User `json:"user"`
	} `json:"data"`
}

type profileEnvelope struct {
	Success bool            `json:"success,omitempty"`
	Status  string          `json:"status,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// GetProfile retrieves the current user's profile.
func (s *UserService) GetProfile(ctx context.Context) (*ProfileResponse, *http.Response, error) {
	u := "profile/data"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	rawResp := new(profileEnvelope)
	resp, err := s.client.Do(ctx, req, rawResp)
	if err != nil {
		return nil, resp, err
	}

	user, err := parseProfileUser(rawResp.Data)
	if err != nil {
		return nil, resp, err
	}

	profileResp := &ProfileResponse{
		Status:  coalesceNonEmpty(rawResp.Status, statusFromSuccess(rawResp.Success)),
		Message: rawResp.Message,
	}
	profileResp.Data.User = user

	return profileResp, resp, nil
}

func parseProfileUser(data json.RawMessage) (User, error) {
	if len(data) == 0 || string(data) == "null" {
		return User{}, nil
	}

	var wrapped struct {
		User *User `json:"user,omitempty"`
	}
	if err := json.Unmarshal(data, &wrapped); err == nil && wrapped.User != nil {
		return *wrapped.User, nil
	}

	var raw struct {
		ID                       jsonID     `json:"id,omitempty"`
		UUID                     string     `json:"uuid,omitempty"`
		Email                    string     `json:"email,omitempty"`
		FirstName                string     `json:"first_name,omitempty"`
		LastName                 string     `json:"last_name,omitempty"`
		FullName                 string     `json:"full_name,omitempty"`
		AvatarURL                string     `json:"avatar_url,omitempty"`
		EmailVerified            *bool      `json:"email_verified,omitempty"`
		PasswordChangedDate      *Timestamp `json:"password_changed_date,omitempty"`
		OAuthUser                string     `json:"oauth_user,omitempty"`
		TempPlanID               int        `json:"temp_plan_id,omitempty"`
		PaymentMethod            *bool      `json:"payment_method,omitempty"`
		ChargeFailed             *bool      `json:"charge_failed,omitempty"`
		IsSubscriptionActive     *bool      `json:"is_subscription_active,omitempty"`
		IsSubscriptionActiveDate *Timestamp `json:"is_subscription_active_date,omitempty"`
		Namespace                string     `json:"namespace,omitempty"`
		CreatedAt                *Timestamp `json:"created_at,omitempty"`
		UpdatedAt                *Timestamp `json:"updated_at,omitempty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return User{}, err
	}

	user := User{
		ID:                       raw.ID.String(),
		UUID:                     raw.UUID,
		Email:                    raw.Email,
		FirstName:                raw.FirstName,
		LastName:                 raw.LastName,
		FullName:                 raw.FullName,
		AvatarURL:                raw.AvatarURL,
		PasswordChangedDate:      raw.PasswordChangedDate,
		OAuthUser:                raw.OAuthUser,
		TempPlanID:               raw.TempPlanID,
		Namespace:                raw.Namespace,
		CreatedAt:                raw.CreatedAt,
		UpdatedAt:                raw.UpdatedAt,
		IsSubscriptionActiveDate: raw.IsSubscriptionActiveDate,
	}
	if user.FullName == "" {
		user.FullName = strings.TrimSpace(strings.Join([]string{raw.FirstName, raw.LastName}, " "))
	}
	if raw.EmailVerified != nil {
		user.EmailVerified = *raw.EmailVerified
	}
	if raw.PaymentMethod != nil {
		user.PaymentMethod = *raw.PaymentMethod
	}
	if raw.ChargeFailed != nil {
		user.ChargeFailed = *raw.ChargeFailed
	}
	if raw.IsSubscriptionActive != nil {
		user.IsSubscriptionActive = *raw.IsSubscriptionActive
	}

	return user, nil
}

// UpdateProfileRequest represents a request to update user profile.
type UpdateProfileRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Bio       string `json:"bio,omitempty"`
}

// UpdateProfile updates the current user's profile.
func (s *UserService) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) (*ProfileResponse, *http.Response, error) {
	u := "profile"

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	profileResp := new(ProfileResponse)
	resp, err := s.client.Do(ctx, httpReq, profileResp)
	if err != nil {
		return nil, resp, err
	}

	return profileResp, resp, nil
}

// ResetSecretToken resets the user's secret token (DEPRECATED).
func (s *UserService) ResetSecretToken(ctx context.Context) (*http.Response, error) {
	u := "user-settings/reset-secret"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// DeleteProfile initiates user profile deletion.
func (s *UserService) DeleteProfile(ctx context.Context) (*http.Response, error) {
	u := "user/delete-profile"

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// CancelProfileDeletion cancels a pending profile deletion request.
func (s *UserService) CancelProfileDeletion(ctx context.Context) (*http.Response, error) {
	u := "user/delete-profile/cancel"

	req, err := s.client.NewRequest(http.MethodPut, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}
