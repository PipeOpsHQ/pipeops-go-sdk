package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// TeamService handles communication with the team related
// methods of the PipeOps API.
type TeamService struct {
	client *Client
}

// Team represents a PipeOps team.
type Team struct {
	ID          string     `json:"id,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	OwnerID     string     `json:"owner_id,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`
}

// TeamsResponse represents a list of teams response.
type TeamsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Teams []Team `json:"teams"`
	} `json:"data"`
}

// TeamResponse represents a single team response.
type TeamResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Team Team `json:"team"`
	} `json:"data"`
}

// CreateTeamRequest represents a request to create a team.
type CreateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Create creates a new team.
func (s *TeamService) Create(ctx context.Context, req *CreateTeamRequest) (*TeamResponse, *http.Response, error) {
	u := "team/create"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	teamResp := new(TeamResponse)
	resp, err := s.client.Do(ctx, httpReq, teamResp)
	if err != nil {
		return nil, resp, err
	}

	return teamResp, resp, nil
}

// UpdateTeamRequest represents a request to update a team.
type UpdateTeamRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Update updates a team.
func (s *TeamService) Update(ctx context.Context, teamUUID string, req *UpdateTeamRequest) (*TeamResponse, *http.Response, error) {
	u := fmt.Sprintf("team/%s/update", teamUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, nil, err
	}

	teamResp := new(TeamResponse)
	resp, err := s.client.Do(ctx, httpReq, teamResp)
	if err != nil {
		return nil, resp, err
	}

	return teamResp, resp, nil
}

// InviteTeamMemberRequest represents a request to invite a team member.
type InviteTeamMemberRequest struct {
	Email       string   `json:"email"`
	Role        string   `json:"role,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// InviteTeamMemberResponse represents a team member invite response.
type InviteTeamMemberResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		InviteID string `json:"invite_id,omitempty"`
	} `json:"data"`
}

// InviteMember invites a new member to the team.
func (s *TeamService) InviteMember(ctx context.Context, teamUUID string, req *InviteTeamMemberRequest) (*InviteTeamMemberResponse, *http.Response, error) {
	u := fmt.Sprintf("team/%s/invite", teamUUID)

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	inviteResp := new(InviteTeamMemberResponse)
	resp, err := s.client.Do(ctx, httpReq, inviteResp)
	if err != nil {
		return nil, resp, err
	}

	return inviteResp, resp, nil
}

// List lists all teams for the authenticated user.
func (s *TeamService) List(ctx context.Context) (*TeamsResponse, *http.Response, error) {
	u := "team"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	teamsResp := new(TeamsResponse)
	resp, err := s.client.Do(ctx, req, teamsResp)
	if err != nil {
		return nil, resp, err
	}

	return teamsResp, resp, nil
}

// Get fetches a team by UUID.
func (s *TeamService) Get(ctx context.Context, teamUUID string) (*TeamResponse, *http.Response, error) {
	u := fmt.Sprintf("team/%s", teamUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	teamResp := new(TeamResponse)
	resp, err := s.client.Do(ctx, req, teamResp)
	if err != nil {
		return nil, resp, err
	}

	return teamResp, resp, nil
}

// Delete deletes a team.
func (s *TeamService) Delete(ctx context.Context, teamUUID string) (*http.Response, error) {
	u := fmt.Sprintf("team/%s", teamUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// TeamMember represents a team member.
type TeamMember struct {
	ID          string     `json:"id,omitempty"`
	UUID        string     `json:"uuid,omitempty"`
	Email       string     `json:"email,omitempty"`
	Role        string     `json:"role,omitempty"`
	Permissions []string   `json:"permissions,omitempty"`
	JoinedAt    *Timestamp `json:"joined_at,omitempty"`
}

// TeamMembersResponse represents team members response.
type TeamMembersResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Members []TeamMember `json:"members"`
	} `json:"data"`
}

// ListMembers lists all members of a team.
func (s *TeamService) ListMembers(ctx context.Context, teamUUID string) (*TeamMembersResponse, *http.Response, error) {
	u := fmt.Sprintf("team/%s/members", teamUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	membersResp := new(TeamMembersResponse)
	resp, err := s.client.Do(ctx, req, membersResp)
	if err != nil {
		return nil, resp, err
	}

	return membersResp, resp, nil
}

// RemoveMember removes a member from a team.
func (s *TeamService) RemoveMember(ctx context.Context, teamUUID, memberUUID string) (*http.Response, error) {
	u := fmt.Sprintf("team/%s/members/%s", teamUUID, memberUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// UpdateMemberRoleRequest represents a request to update member role.
type UpdateMemberRoleRequest struct {
	Role        string   `json:"role"`
	Permissions []string `json:"permissions,omitempty"`
}

// UpdateMemberRole updates a team member's role.
func (s *TeamService) UpdateMemberRole(ctx context.Context, teamUUID, memberUUID string, req *UpdateMemberRoleRequest) (*http.Response, error) {
	u := fmt.Sprintf("team/%s/members/%s/role", teamUUID, memberUUID)

	httpReq, err := s.client.NewRequest(http.MethodPut, u, req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// AcceptInvitation accepts a team invitation.
func (s *TeamService) AcceptInvitation(ctx context.Context, inviteToken string) (*http.Response, error) {
	u := fmt.Sprintf("team/invite/accept/%s", inviteToken)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// RejectInvitation rejects a team invitation.
func (s *TeamService) RejectInvitation(ctx context.Context, inviteToken string) (*http.Response, error) {
	u := fmt.Sprintf("team/invite/reject/%s", inviteToken)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}
