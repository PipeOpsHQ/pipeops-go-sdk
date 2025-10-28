package pipeops

import (
	"context"
	"net/http"
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
	u := "user/settings/notification"

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

// GetProfile retrieves the current user's profile.
func (s *UserService) GetProfile(ctx context.Context) (*ProfileResponse, *http.Response, error) {
	u := "profile"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	profileResp := new(ProfileResponse)
	resp, err := s.client.Do(ctx, req, profileResp)
	if err != nil {
		return nil, resp, err
	}

	return profileResp, resp, nil
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
