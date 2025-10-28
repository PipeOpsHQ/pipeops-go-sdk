package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// EventService handles communication with the events related
// methods of the PipeOps API.
type EventService struct {
	client *Client
}

// Event represents a system event.
type Event struct {
	ID          string     `json:"id,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	Type        string     `json:"type,omitempty"`
	Description string     `json:"description,omitempty"`
	Enabled     bool       `json:"enabled,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
}

// EventsResponse represents a list of events response.
type EventsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Events []Event `json:"events"`
	} `json:"data"`
}

// EventResponse represents a single event response.
type EventResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Event Event `json:"event"`
	} `json:"data"`
}

// ListEvents lists all events.
func (s *EventService) ListEvents(ctx context.Context) (*EventsResponse, *http.Response, error) {
	u := "user-settings/events"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	eventsResp := new(EventsResponse)
	resp, err := s.client.Do(ctx, req, eventsResp)
	if err != nil {
		return nil, resp, err
	}

	return eventsResp, resp, nil
}

// ToggleEventRequest represents a request to toggle an event.
type ToggleEventRequest struct {
	Enabled bool `json:"enabled"`
}

// ToggleEvent toggles an event on/off.
func (s *EventService) ToggleEvent(ctx context.Context, eventUUID string, req *ToggleEventRequest) (*EventResponse, *http.Response, error) {
	u := fmt.Sprintf("user-settings/events/toggle/%s", eventUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	eventResp := new(EventResponse)
	resp, err := s.client.Do(ctx, httpReq, eventResp)
	if err != nil {
		return nil, resp, err
	}

	return eventResp, resp, nil
}

// GetResourceEvents gets resource usage events.
func (s *EventService) GetResourceEvents(ctx context.Context) (*EventsResponse, *http.Response, error) {
	u := "user-settings/resource/events"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	eventsResp := new(EventsResponse)
	resp, err := s.client.Do(ctx, req, eventsResp)
	if err != nil {
		return nil, resp, err
	}

	return eventsResp, resp, nil
}

// UpdateResourceEventRequest represents a request to update resource event settings.
type UpdateResourceEventRequest struct {
	Enabled   bool   `json:"enabled"`
	Threshold int    `json:"threshold,omitempty"`
	EventType string `json:"event_type,omitempty"`
}

// UpdateResourceEvent updates resource usage event settings.
func (s *EventService) UpdateResourceEvent(ctx context.Context, eventUUID string, req *UpdateResourceEventRequest) (*EventResponse, *http.Response, error) {
	u := fmt.Sprintf("user-settings/resource/events/%s", eventUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	eventResp := new(EventResponse)
	resp, err := s.client.Do(ctx, httpReq, eventResp)
	if err != nil {
		return nil, resp, err
	}

	return eventResp, resp, nil
}

// SurveyService handles communication with the survey related
// methods of the PipeOps API.
type SurveyService struct {
	client *Client
}

// Survey represents a survey.
type Survey struct {
	ID        string     `json:"id,omitempty"`
	UUID      string     `json:"uuid,omitempty"`
	RoleID    string     `json:"role_id,omitempty"`
	Answers   []string   `json:"answers,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
}

// SurveyResponse represents a survey response.
type SurveyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Survey Survey `json:"survey"`
	} `json:"data"`
}

// CreateSurveyRequest represents a request to create a survey.
type CreateSurveyRequest struct {
	RoleID  string   `json:"role_id"`
	Answers []string `json:"answers"`
}

// CreateSurvey creates an onboarding survey.
func (s *SurveyService) CreateSurvey(ctx context.Context, req *CreateSurveyRequest) (*SurveyResponse, *http.Response, error) {
	u := "user/create-onboarding-survey"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	surveyResp := new(SurveyResponse)
	resp, err := s.client.Do(ctx, httpReq, surveyResp)
	if err != nil {
		return nil, resp, err
	}

	return surveyResp, resp, nil
}

// SurveyRolesResponse represents survey roles response.
type SurveyRolesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Roles []map[string]interface{} `json:"roles"`
	} `json:"data"`
}

// GetSurveyRoles gets available survey roles.
func (s *SurveyService) GetSurveyRoles(ctx context.Context) (*SurveyRolesResponse, *http.Response, error) {
	u := "user/survey/roles"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	rolesResp := new(SurveyRolesResponse)
	resp, err := s.client.Do(ctx, req, rolesResp)
	if err != nil {
		return nil, resp, err
	}

	return rolesResp, resp, nil
}

// SurveyQuestionsResponse represents survey questions response.
type SurveyQuestionsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Questions []map[string]interface{} `json:"questions"`
	} `json:"data"`
}

// GetRoleQuestions gets questions for a survey role.
func (s *SurveyService) GetRoleQuestions(ctx context.Context, roleID string) (*SurveyQuestionsResponse, *http.Response, error) {
	u := fmt.Sprintf("user/survey/roles/%s", roleID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	questionsResp := new(SurveyQuestionsResponse)
	resp, err := s.client.Do(ctx, req, questionsResp)
	if err != nil {
		return nil, resp, err
	}

	return questionsResp, resp, nil
}

// SurveyDiscoveriesResponse represents survey discoveries response.
type SurveyDiscoveriesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Discoveries []string `json:"discoveries"`
	} `json:"data"`
}

// GetSurveyDiscoveries gets survey discovery options.
func (s *SurveyService) GetSurveyDiscoveries(ctx context.Context) (*SurveyDiscoveriesResponse, *http.Response, error) {
	u := "user/survey/discoveries"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	discoveriesResp := new(SurveyDiscoveriesResponse)
	resp, err := s.client.Do(ctx, req, discoveriesResp)
	if err != nil {
		return nil, resp, err
	}

	return discoveriesResp, resp, nil
}

// PartnerService handles communication with the partners related
// methods of the PipeOps API.
type PartnerService struct {
	client *Client
}

// Partner represents a partner.
type Partner struct {
	ID          string     `json:"id,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
}

// PartnerResponse represents a partner response.
type PartnerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Partner Partner `json:"partner"`
	} `json:"data"`
}

// PartnersResponse represents a list of partners response.
type PartnersResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Partners []Partner `json:"partners"`
	} `json:"data"`
}

// CreatePartnerRequest represents a request to create a partner.
type CreatePartnerRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Create creates a new partner.
func (s *PartnerService) Create(ctx context.Context, req *CreatePartnerRequest) (*PartnerResponse, *http.Response, error) {
	u := "partners"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	partnerResp := new(PartnerResponse)
	resp, err := s.client.Do(ctx, httpReq, partnerResp)
	if err != nil {
		return nil, resp, err
	}

	return partnerResp, resp, nil
}

// UpdatePartnerRequest represents a request to update a partner.
type UpdatePartnerRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Update updates a partner.
func (s *PartnerService) Update(ctx context.Context, partnerUUID string, req *UpdatePartnerRequest) (*PartnerResponse, *http.Response, error) {
	u := fmt.Sprintf("partners/%s", partnerUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	partnerResp := new(PartnerResponse)
	resp, err := s.client.Do(ctx, httpReq, partnerResp)
	if err != nil {
		return nil, resp, err
	}

	return partnerResp, resp, nil
}

// Get gets a partner by UUID.
func (s *PartnerService) Get(ctx context.Context, partnerUUID string) (*PartnerResponse, *http.Response, error) {
	u := fmt.Sprintf("partners/%s", partnerUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	partnerResp := new(PartnerResponse)
	resp, err := s.client.Do(ctx, req, partnerResp)
	if err != nil {
		return nil, resp, err
	}

	return partnerResp, resp, nil
}

// List lists all partners.
func (s *PartnerService) List(ctx context.Context) (*PartnersResponse, *http.Response, error) {
	u := "partners"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	partnersResp := new(PartnersResponse)
	resp, err := s.client.Do(ctx, req, partnersResp)
	if err != nil {
		return nil, resp, err
	}

	return partnersResp, resp, nil
}

// MiscService handles miscellaneous API methods.
type MiscService struct {
	client *Client
}

// ContactUsRequest represents a contact us request.
type ContactUsRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// ContactUsResponse represents a contact us response.
type ContactUsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ContactUs sends a contact us message.
func (s *MiscService) ContactUs(ctx context.Context, req *ContactUsRequest) (*ContactUsResponse, *http.Response, error) {
	u := "misc/contact_us"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	contactResp := new(ContactUsResponse)
	resp, err := s.client.Do(ctx, httpReq, contactResp)
	if err != nil {
		return nil, resp, err
	}

	return contactResp, resp, nil
}

// JoinWaitlistRequest represents a join waitlist request.
type JoinWaitlistRequest struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// JoinWaitlistResponse represents a join waitlist response.
type JoinWaitlistResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// JoinWaitlist joins the waitlist.
func (s *MiscService) JoinWaitlist(ctx context.Context, req *JoinWaitlistRequest) (*JoinWaitlistResponse, *http.Response, error) {
	u := "misc/join_waitlist"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	waitlistResp := new(JoinWaitlistResponse)
	resp, err := s.client.Do(ctx, httpReq, waitlistResp)
	if err != nil {
		return nil, resp, err
	}

	return waitlistResp, resp, nil
}

// DashboardDataResponse represents dashboard data response.
type DashboardDataResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Dashboard map[string]interface{} `json:"dashboard"`
	} `json:"data"`
}

// GetDashboardData gets dashboard data.
func (s *MiscService) GetDashboardData(ctx context.Context) (*DashboardDataResponse, *http.Response, error) {
	u := "app/data"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	dashboardResp := new(DashboardDataResponse)
	resp, err := s.client.Do(ctx, req, dashboardResp)
	if err != nil {
		return nil, resp, err
	}

	return dashboardResp, resp, nil
}

// PartnerAgreementService handles partner agreement related methods.
type PartnerAgreementService struct {
	client *Client
}

// PartnerAgreement represents a partner agreement.
type PartnerAgreement struct {
	ID          string     `json:"id,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	PartnerUUID string     `json:"partner_uuid,omitempty"`
	Name        string     `json:"name,omitempty"`
	Terms       string     `json:"terms,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
}

// PartnerAgreementRequest represents a partner agreement request.
type PartnerAgreementRequest struct {
	PartnerUUID string `json:"partner_uuid"`
	Name        string `json:"name"`
	Terms       string `json:"terms,omitempty"`
}

// PartnerAgreementResponse represents a partner agreement response.
type PartnerAgreementResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Agreement PartnerAgreement `json:"agreement"`
	} `json:"data"`
}

// PartnerAgreementsResponse represents partner agreements response.
type PartnerAgreementsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Agreements []PartnerAgreement `json:"agreements"`
	} `json:"data"`
}

// CreateAgreement creates a partner agreement.
func (s *PartnerAgreementService) CreateAgreement(ctx context.Context, req *PartnerAgreementRequest) (*PartnerAgreementResponse, *http.Response, error) {
	u := "partners/agreements"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	agreementResp := new(PartnerAgreementResponse)
	resp, err := s.client.Do(ctx, httpReq, agreementResp)
	if err != nil {
		return nil, resp, err
	}

	return agreementResp, resp, nil
}

// UpdateAgreement updates a partner agreement.
func (s *PartnerAgreementService) UpdateAgreement(ctx context.Context, agreementUUID string, req *PartnerAgreementRequest) (*PartnerAgreementResponse, *http.Response, error) {
	u := fmt.Sprintf("partners/agreements/%s", agreementUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	agreementResp := new(PartnerAgreementResponse)
	resp, err := s.client.Do(ctx, httpReq, agreementResp)
	if err != nil {
		return nil, resp, err
	}

	return agreementResp, resp, nil
}

// GetAgreement gets a partner agreement by UUID.
func (s *PartnerAgreementService) GetAgreement(ctx context.Context, agreementUUID string) (*PartnerAgreementResponse, *http.Response, error) {
	u := fmt.Sprintf("partners/agreements/%s", agreementUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	agreementResp := new(PartnerAgreementResponse)
	resp, err := s.client.Do(ctx, req, agreementResp)
	if err != nil {
		return nil, resp, err
	}

	return agreementResp, resp, nil
}

// ListAgreements lists all partner agreements.
func (s *PartnerAgreementService) ListAgreements(ctx context.Context) (*PartnerAgreementsResponse, *http.Response, error) {
	u := "partners/agreements"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	agreementsResp := new(PartnerAgreementsResponse)
	resp, err := s.client.Do(ctx, req, agreementsResp)
	if err != nil {
		return nil, resp, err
	}

	return agreementsResp, resp, nil
}

// PartnerParticipantService handles partner participant related methods.
type PartnerParticipantService struct {
	client *Client
}

// ParticipantUploadRequest represents a participant upload request.
type ParticipantUploadRequest struct {
	Data map[string]interface{} `json:"data"`
}

// UploadParticipants uploads participants for an agreement.
func (s *PartnerParticipantService) UploadParticipants(ctx context.Context, agreementID string, req *ParticipantUploadRequest) (*http.Response, error) {
	u := fmt.Sprintf("partners/agreements/%s/uploads", agreementID)

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// VerifyCodeResponse represents verification code response.
type VerifyCodeResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Valid bool `json:"valid"`
	} `json:"data"`
}

// VerifyProgramCode verifies a program verification code.
func (s *PartnerParticipantService) VerifyProgramCode(ctx context.Context, code string) (*VerifyCodeResponse, *http.Response, error) {
	u := "partners/participants/verify?code=" + code

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	verifyResp := new(VerifyCodeResponse)
	resp, err := s.client.Do(ctx, req, verifyResp)
	if err != nil {
		return nil, resp, err
	}

	return verifyResp, resp, nil
}

// ProfileService handles profile related methods.
type ProfileService struct {
	client *Client
}

// DeleteProfile deletes the user profile.
func (s *ProfileService) DeleteProfile(ctx context.Context) (*http.Response, error) {
	u := "user/delete-profile"

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// CancelProfileDeletion cancels a pending profile deletion.
func (s *ProfileService) CancelProfileDeletion(ctx context.Context) (*http.Response, error) {
	u := "user/delete-profile/cancel"

	req, err := s.client.NewRequest(http.MethodPut, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}
