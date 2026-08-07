package pipeops

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
// Controller uses success/status interchangeably; Update returns data.team as a name string.
type TeamResponse struct {
	Status  string `json:"status"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Team Team   `json:"-"`
		UUID string `json:"uuid,omitempty"`
	} `json:"data"`
}

// UnmarshalJSON accepts data.team as either a Team object (fetch) or a name string (update).
func (r *TeamResponse) UnmarshalJSON(b []byte) error {
	type wire struct {
		Status  string `json:"status"`
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			Team json.RawMessage `json:"team"`
			UUID string          `json:"uuid,omitempty"`
		} `json:"data"`
	}
	var w wire
	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}
	r.Status = w.Status
	r.Success = w.Success
	r.Message = w.Message
	r.Data.UUID = w.Data.UUID
	if len(w.Data.Team) == 0 || string(w.Data.Team) == "null" {
		return nil
	}
	// Update endpoint returns data.team as the team name string.
	var name string
	if err := json.Unmarshal(w.Data.Team, &name); err == nil {
		r.Data.Team = Team{Name: name, UUID: w.Data.UUID}
		return nil
	}
	var team Team
	if err := json.Unmarshal(w.Data.Team, &team); err != nil {
		return err
	}
	r.Data.Team = team
	if r.Data.UUID == "" && team.UUID != "" {
		r.Data.UUID = team.UUID
	}
	return nil
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
	workspaceUUID, _, wsErr := firstWorkspaceUUID(ctx, s.client)
	if wsErr == nil {
		u, err := addOptions("team/fetch", &teamFetchOptions{WorkspaceUUID: workspaceUUID})
		if err != nil {
			return nil, nil, err
		}

		req, err := s.client.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return nil, nil, err
		}

		teamsResp := new(TeamsResponse)
		resp, err := s.client.Do(ctx, req, teamsResp)
		if err == nil {
			return teamsResp, resp, nil
		}
		if !isNotFound(err) {
			return nil, resp, err
		}
	}

	u := "team/fetch"

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
	u := fmt.Sprintf("team/fetch/%s", teamUUID)

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
	u := fmt.Sprintf("team/%s/delete", teamUUID)

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

// ListMembers lists members of a team.
// Controller has no GET /team/:uuid/members; members are embedded in GET /team/fetch/:uuid.
func (s *TeamService) ListMembers(ctx context.Context, teamUUID string) (*TeamMembersResponse, *http.Response, error) {
	u := fmt.Sprintf("team/fetch/%s", teamUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var raw map[string]interface{}
	resp, err := s.client.Do(ctx, req, &raw)
	if err != nil {
		return nil, resp, err
	}

	membersResp := &TeamMembersResponse{}
	if st, ok := raw["status"].(string); ok {
		membersResp.Status = st
	}
	if msg, ok := raw["message"].(string); ok {
		membersResp.Message = msg
	}
	membersResp.Data.Members = extractTeamMembersFromFetch(raw)
	return membersResp, resp, nil
}

// extractTeamMembersFromFetch pulls members from GET /team/fetch/:uuid payloads.
func extractTeamMembersFromFetch(raw map[string]interface{}) []TeamMember {
	data, ok := raw["data"].(map[string]interface{})
	if !ok || data == nil {
		return nil
	}
	teamObj, ok := data["team"].(map[string]interface{})
	if !ok || teamObj == nil {
		return nil
	}
	// Prefer TeamMembers / team_members arrays on the team object.
	for _, key := range []string{"TeamMembers", "team_members", "members"} {
		if arr, ok := teamObj[key].([]interface{}); ok {
			return mapInterfaceSliceToTeamMembers(arr)
		}
	}
	// TeamResourceUsers is a map[userID]details in some responses.
	if m, ok := teamObj["TeamResourceUsers"].(map[string]interface{}); ok {
		out := make([]TeamMember, 0, len(m))
		for _, v := range m {
			if row, ok := v.(map[string]interface{}); ok {
				out = append(out, teamMemberFromDetailsMap(row))
			}
		}
		return out
	}
	if m, ok := teamObj["team_resource_users"].(map[string]interface{}); ok {
		out := make([]TeamMember, 0, len(m))
		for _, v := range m {
			if row, ok := v.(map[string]interface{}); ok {
				out = append(out, teamMemberFromDetailsMap(row))
			}
		}
		return out
	}
	return nil
}

func mapInterfaceSliceToTeamMembers(arr []interface{}) []TeamMember {
	out := make([]TeamMember, 0, len(arr))
	for _, item := range arr {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, teamMemberFromDetailsMap(m))
	}
	return out
}

func teamMemberFromDetailsMap(m map[string]interface{}) TeamMember {
	tm := TeamMember{}
	if s, ok := m["uuid"].(string); ok {
		tm.UUID = s
	}
	if s, ok := m["UUID"].(string); ok && tm.UUID == "" {
		tm.UUID = s
	}
	if s, ok := m["role"].(string); ok {
		tm.Role = s
	}
	if s, ok := m["Role"].(string); ok && tm.Role == "" {
		tm.Role = s
	}
	if s, ok := m["email"].(string); ok {
		tm.Email = s
	}
	// Nested user object from TeamResourceUserDetails
	if user, ok := m["User"].(map[string]interface{}); ok {
		if s, ok := user["uuid"].(string); ok && tm.UUID == "" {
			tm.UUID = s
		}
		if s, ok := user["UUID"].(string); ok && tm.UUID == "" {
			tm.UUID = s
		}
		if s, ok := user["email"].(string); ok && tm.Email == "" {
			tm.Email = s
		}
		if s, ok := user["Email"].(string); ok && tm.Email == "" {
			tm.Email = s
		}
	}
	if user, ok := m["user"].(map[string]interface{}); ok {
		if s, ok := user["uuid"].(string); ok && tm.UUID == "" {
			tm.UUID = s
		}
		if s, ok := user["email"].(string); ok && tm.Email == "" {
			tm.Email = s
		}
	}
	return tm
}

// RemoveMember removes a member from a team.
// Controller: DELETE /team/:uuid/delete-member/:member_user_uuid
// memberUUID must be the member's user UUID (not email).
func (s *TeamService) RemoveMember(ctx context.Context, teamUUID, memberUserUUID string) (*http.Response, error) {
	u := fmt.Sprintf("team/%s/delete-member/%s", teamUUID, memberUserUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// UpdateMemberRoleRequest represents a request to update member role.
// Matches controller UpdateMemberPermissionInput (role + optional resource permissions).
type UpdateMemberRoleRequest struct {
	Role        string   `json:"role"`
	Permissions []string `json:"permissions,omitempty"` // not sent as-is; kept for MCP compat
	// AccessLevel optional: "workspace" | "resource"
	AccessLevel string `json:"access_level,omitempty"`
}

// UpdateMemberRole updates a team member's role.
// Controller: PUT /team/:uuid/update-member-permissions/:member_user_uuid
// memberUUID must be the member's user UUID (not email).
func (s *TeamService) UpdateMemberRole(ctx context.Context, teamUUID, memberUserUUID string, req *UpdateMemberRoleRequest) (*http.Response, error) {
	u := fmt.Sprintf("team/%s/update-member-permissions/%s", teamUUID, memberUserUUID)

	body := map[string]interface{}{
		"role":                    strings.TrimSpace(req.Role),
		"server_permissions":      []interface{}{},
		"project_permissions":     []interface{}{},
		"environment_permissions": []interface{}{},
		"addon_permissions":       []interface{}{},
	}
	if strings.TrimSpace(req.AccessLevel) != "" {
		body["access_level"] = strings.TrimSpace(req.AccessLevel)
	}

	httpReq, err := s.client.NewRequest(http.MethodPut, u, body)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	return resp, err
}

// AcceptInvitation accepts a team invitation.
func (s *TeamService) AcceptInvitation(ctx context.Context, inviteToken string) (*http.Response, error) {
	u := "team/accept-invite"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, &AcceptInviteRequest{InviteID: inviteToken})
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
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

// AcceptInviteRequest represents a request to accept a team invite.
type AcceptInviteRequest struct {
	InviteID    string `json:"invite_id"`
	InviteEmail string `json:"invite_email,omitempty"`
}

type teamFetchOptions struct {
	WorkspaceUUID string `url:"workspace_uuid"`
}
