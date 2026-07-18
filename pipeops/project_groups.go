package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// ProjectGroupService handles unified project plane (project group) APIs.
//
// Controller routes (JWT + team access):
//
//	GET    /project-groups
//	POST   /project-groups
//	GET    /project-groups/resolve
//	GET    /project-groups/candidates
//	GET    /project-groups/:uuid
//	PATCH  /project-groups/:uuid
//	DELETE /project-groups/:uuid
//	GET    /project-groups/:uuid/topology
//	GET    /project-groups/:uuid/env
//	PUT    /project-groups/:uuid/env
//	POST   /project-groups/:uuid/env/inject
//	POST   /project-groups/:uuid/members
//	DELETE /project-groups/:uuid/members/:memberType/:memberUUID
//	POST   /project-groups/:uuid/connections
//	POST   /project-groups/:uuid/redeploy-apps
type ProjectGroupService struct {
	client *Client
}

// ProjectGroupWorkspaceOptions carries workspace query params used by most endpoints.
type ProjectGroupWorkspaceOptions struct {
	WorkspaceUUID string `url:"workspace_uuid,omitempty"`
	// Workspace is accepted as an alias by the controller.
	Workspace string `url:"workspace,omitempty"`
}

// ProjectGroupListOptions filters and paginates group list.
type ProjectGroupListOptions struct {
	WorkspaceUUID string `url:"workspace_uuid,omitempty"`
	Workspace     string `url:"workspace,omitempty"`
	Limit         int    `url:"limit,omitempty"`
	Offset        int    `url:"offset,omitempty"`
}

// ProjectGroupResolveOptions is GET /project-groups/resolve.
type ProjectGroupResolveOptions struct {
	WorkspaceUUID string `url:"workspace_uuid,omitempty"`
	Workspace     string `url:"workspace,omitempty"`
	MemberType    string `url:"member_type,omitempty"`
	MemberUUID    string `url:"member_uuid,omitempty"`
}

// ProjectGroupCandidatesOptions is GET /project-groups/candidates.
type ProjectGroupCandidatesOptions struct {
	WorkspaceUUID string `url:"workspace_uuid,omitempty"`
	Workspace     string `url:"workspace,omitempty"`
	GroupUUID     string `url:"group_uuid,omitempty"`
}

// ProjectGroupDetachOptions is DELETE .../members/... query options.
type ProjectGroupDetachOptions struct {
	WorkspaceUUID  string `url:"workspace_uuid,omitempty"`
	Workspace      string `url:"workspace,omitempty"`
	IncludeSession *bool  `url:"include_session,omitempty"`
}

// ProjectGroup is the API representation of a project group (UI: Project).
type ProjectGroup struct {
	UUID                   string               `json:"uuid,omitempty"`
	Name                   string               `json:"name,omitempty"`
	NameSlug               string               `json:"name_slug,omitempty"`
	WorkspaceUUID          string               `json:"workspace_uuid,omitempty"`
	DefaultClusterUUID     string               `json:"default_cluster_uuid,omitempty"`
	DefaultEnvironmentUUID string               `json:"default_environment_uuid,omitempty"`
	MemberCount            int                  `json:"member_count,omitempty"`
	Members                []ProjectGroupMember `json:"members,omitempty"`
	CreatedAt              string               `json:"created_at,omitempty"`
	UpdatedAt              string               `json:"updated_at,omitempty"`
}

// ProjectGroupMember is a service membership row.
type ProjectGroupMember struct {
	MemberType     string `json:"member_type,omitempty"`
	MemberUUID     string `json:"member_uuid,omitempty"`
	ServiceKind    string `json:"service_kind,omitempty"`
	DisplayOrder   int    `json:"display_order,omitempty"`
	Name           string `json:"name,omitempty"`
	Status         string `json:"status,omitempty"`
	ClusterUUID    string `json:"cluster_uuid,omitempty"`
	Environment    string `json:"environment,omitempty"`
	OwnerHref      string `json:"owner_href,omitempty"`
	OwnerSessionID string `json:"owner_session_id,omitempty"`
}

// CreateProjectGroupRequest creates an empty group.
type CreateProjectGroupRequest struct {
	Name                   string  `json:"name"`
	DefaultClusterUUID     *string `json:"default_cluster_uuid,omitempty"`
	DefaultEnvironmentUUID *string `json:"default_environment_uuid,omitempty"`
}

// UpdateProjectGroupRequest patches group metadata.
type UpdateProjectGroupRequest struct {
	Name                   *string `json:"name,omitempty"`
	DefaultClusterUUID     *string `json:"default_cluster_uuid,omitempty"`
	DefaultEnvironmentUUID *string `json:"default_environment_uuid,omitempty"`
}

// AttachProjectGroupMemberRequest attaches a service to a group.
type AttachProjectGroupMemberRequest struct {
	MemberType     string `json:"member_type"` // project | addon_deployment
	MemberUUID     string `json:"member_uuid"`
	IncludeSession *bool  `json:"include_session,omitempty"`
	Move           bool   `json:"move,omitempty"`
}

// AttachProjectGroupMemberResponse is returned after attach/move.
type AttachProjectGroupMemberResponse struct {
	AttachedMemberUUIDs   []string `json:"attached_member_uuids,omitempty"`
	IncludeSessionApplied bool     `json:"include_session_applied,omitempty"`
	GroupUUID             string   `json:"group_uuid,omitempty"`
}

// ProjectGroupTopologyNode is a plane service card.
type ProjectGroupTopologyNode struct {
	MemberType      string   `json:"member_type,omitempty"`
	MemberUUID      string   `json:"member_uuid,omitempty"`
	ServiceKind     string   `json:"service_kind,omitempty"`
	Name            string   `json:"name,omitempty"`
	Status          string   `json:"status,omitempty"`
	ClusterUUID     string   `json:"cluster_uuid,omitempty"`
	EnvironmentUUID string   `json:"environment_uuid,omitempty"`
	Namespace       string   `json:"namespace,omitempty"`
	OwnerHref       string   `json:"owner_href,omitempty"`
	OwnerSessionID  string   `json:"owner_session_id,omitempty"`
	InternalURL     string   `json:"internal_url,omitempty"`
	PublicURLs      []string `json:"public_urls,omitempty"`
	PrivateHostname string   `json:"private_hostname,omitempty"`
	ParentUUID      string   `json:"parent_uuid,omitempty"`
	PosX            float64  `json:"pos_x,omitempty"`
	PosY            float64  `json:"pos_y,omitempty"`
}

// ProjectGroupTopologyEdge connects two members or a volume.
type ProjectGroupTopologyEdge struct {
	Type       string `json:"type,omitempty"`
	FromUUID   string `json:"from_uuid,omitempty"`
	ToUUID     string `json:"to_uuid,omitempty"`
	Label      string `json:"label,omitempty"`
	Confidence string `json:"confidence,omitempty"`
}

// ProjectGroupTopologyVolume is a volume chip on the plane.
type ProjectGroupTopologyVolume struct {
	UUID        string  `json:"uuid,omitempty"`
	DisplayName string  `json:"display_name,omitempty"`
	PVCName     string  `json:"pvc_name,omitempty"`
	Status      string  `json:"status,omitempty"`
	OwnerType   string  `json:"owner_type,omitempty"`
	OwnerUUID   string  `json:"owner_uuid,omitempty"`
	SizeGB      float32 `json:"size_gb,omitempty"`
}

// ProjectGroupEnvironment is a group environment slot.
type ProjectGroupEnvironment struct {
	Slug        string `json:"slug,omitempty"`
	Name        string `json:"name,omitempty"`
	ClusterUUID string `json:"cluster_uuid,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	IsDefault   bool   `json:"is_default,omitempty"`
}

// ProjectGroupTopology is the plane payload.
type ProjectGroupTopology struct {
	Group              ProjectGroup                 `json:"group"`
	Nodes              []ProjectGroupTopologyNode   `json:"nodes,omitempty"`
	Edges              []ProjectGroupTopologyEdge   `json:"edges,omitempty"`
	Volumes            []ProjectGroupTopologyVolume `json:"volumes,omitempty"`
	UnattachedVolumes  []ProjectGroupTopologyVolume `json:"unattached_volumes,omitempty"`
	Warnings           []string                     `json:"warnings,omitempty"`
	TotalMemberCount   int                          `json:"total_member_count,omitempty"`
	VisibleMemberCount int                          `json:"visible_member_count,omitempty"`
	NestedNodeCount    int                          `json:"nested_node_count,omitempty"`
	Environments       []ProjectGroupEnvironment    `json:"environments,omitempty"`
	ActiveEnvironment  string                       `json:"active_environment,omitempty"`
}

// ResolveProjectGroupResponse maps a service id to its group.
type ResolveProjectGroupResponse struct {
	GroupUUID  string `json:"group_uuid,omitempty"`
	MemberType string `json:"member_type,omitempty"`
	MemberUUID string `json:"member_uuid,omitempty"`
}

// ConnectProjectGroupServicesRequest wires provider connection envs into a consumer.
type ConnectProjectGroupServicesRequest struct {
	ConsumerType string `json:"consumer_type"` // project
	ConsumerUUID string `json:"consumer_uuid"`
	ProviderType string `json:"provider_type"` // addon_deployment
	ProviderUUID string `json:"provider_uuid"`
	Overwrite    bool   `json:"overwrite,omitempty"`
	VariableSet  string `json:"variable_set,omitempty"`
}

// ConnectProjectGroupServicesResponse is returned after env wiring.
type ConnectProjectGroupServicesResponse struct {
	WrittenKeys      []string                 `json:"written_keys,omitempty"`
	SkippedKeys      []string                 `json:"skipped_keys,omitempty"`
	RestartTriggered bool                     `json:"restart_triggered,omitempty"`
	Edge             ProjectGroupTopologyEdge `json:"edge,omitempty"`
	Message          string                   `json:"message,omitempty"`
}

// ProjectGroupAttachCandidate is a project or addon that can be attached.
type ProjectGroupAttachCandidate struct {
	MemberType       string `json:"member_type,omitempty"`
	MemberUUID       string `json:"member_uuid,omitempty"`
	Name             string `json:"name,omitempty"`
	ServiceKind      string `json:"service_kind,omitempty"`
	Status           string `json:"status,omitempty"`
	SessionID        string `json:"session_id,omitempty"`
	CurrentGroupUUID string `json:"current_group_uuid,omitempty"`
	CurrentGroupName string `json:"current_group_name,omitempty"`
	InTargetGroup    bool   `json:"in_target_group,omitempty"`
}

// ProjectGroupAttachCandidates lists attachable services for the picker UI.
type ProjectGroupAttachCandidates struct {
	Projects []ProjectGroupAttachCandidate `json:"projects,omitempty"`
	Addons   []ProjectGroupAttachCandidate `json:"addons,omitempty"`
}

// ProjectGroupSharedEnvVar is a single shared key/value.
type ProjectGroupSharedEnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ProjectGroupSharedEnv is the group-level shared environment variables payload.
type ProjectGroupSharedEnv struct {
	Variables       []ProjectGroupSharedEnvVar `json:"variables,omitempty"`
	Injected        bool                       `json:"injected,omitempty"`
	WrittenKeys     []string                   `json:"written_keys,omitempty"`
	SkippedKeys     []string                   `json:"skipped_keys,omitempty"`
	ProjectsTouched []string                   `json:"projects_touched,omitempty"`
	AddonsTouched   []string                   `json:"addons_touched,omitempty"`
	RedeployQueued  []string                   `json:"redeploy_queued,omitempty"`
	Message         string                     `json:"message,omitempty"`
}

// UpsertProjectGroupSharedEnvRequest replaces the group shared env set.
type UpsertProjectGroupSharedEnvRequest struct {
	Variables      []ProjectGroupSharedEnvVar `json:"variables"`
	Inject         bool                       `json:"inject,omitempty"`
	Overwrite      bool                       `json:"overwrite,omitempty"`
	Redeploy       bool                       `json:"redeploy,omitempty"`
	KeepReferences bool                       `json:"keep_references,omitempty"`
}

// InjectProjectGroupSharedEnvRequest pushes stored group shared env into members.
type InjectProjectGroupSharedEnvRequest struct {
	Overwrite      bool     `json:"overwrite,omitempty"`
	Redeploy       bool     `json:"redeploy,omitempty"`
	MemberUUIDs    []string `json:"member_uuids,omitempty"`
	KeepReferences bool     `json:"keep_references,omitempty"`
}

// InjectProjectGroupSharedEnvResponse reports inject results.
type InjectProjectGroupSharedEnvResponse struct {
	WrittenKeys     []string `json:"written_keys,omitempty"`
	SkippedKeys     []string `json:"skipped_keys,omitempty"`
	ProjectsTouched []string `json:"projects_touched,omitempty"`
	AddonsTouched   []string `json:"addons_touched,omitempty"`
	RedeployQueued  []string `json:"redeploy_queued,omitempty"`
	Message         string   `json:"message,omitempty"`
}

// RedeployProjectGroupAppsResponse is bulk app redeploy result.
type RedeployProjectGroupAppsResponse struct {
	Queued  []string `json:"queued,omitempty"`
	Failed  []string `json:"failed,omitempty"`
	Message string   `json:"message,omitempty"`
}

// ProjectGroupListResponse is GET /project-groups.
type ProjectGroupListResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Groups []ProjectGroup `json:"groups"`
		Total  int64          `json:"total"`
		Limit  int            `json:"limit"`
		Offset int            `json:"offset"`
	} `json:"data"`
}

// ProjectGroupResponse is a single group envelope.
type ProjectGroupResponse struct {
	Success bool         `json:"success,omitempty"`
	Message string       `json:"message,omitempty"`
	Data    ProjectGroup `json:"data"`
}

// ProjectGroupAttachResponse is POST .../members.
type ProjectGroupAttachResponse struct {
	Success bool                             `json:"success,omitempty"`
	Message string                           `json:"message,omitempty"`
	Data    AttachProjectGroupMemberResponse `json:"data"`
}

// ProjectGroupTopologyResponse is GET .../topology.
type ProjectGroupTopologyResponse struct {
	Success bool                 `json:"success,omitempty"`
	Message string               `json:"message,omitempty"`
	Data    ProjectGroupTopology `json:"data"`
}

// ProjectGroupResolveResponse is GET /project-groups/resolve.
type ProjectGroupResolveResponse struct {
	Success bool                        `json:"success,omitempty"`
	Message string                      `json:"message,omitempty"`
	Data    ResolveProjectGroupResponse `json:"data"`
}

// ProjectGroupCandidatesResponse is GET /project-groups/candidates.
type ProjectGroupCandidatesResponse struct {
	Success bool                         `json:"success,omitempty"`
	Message string                       `json:"message,omitempty"`
	Data    ProjectGroupAttachCandidates `json:"data"`
}

// ProjectGroupSharedEnvResponse is GET/PUT .../env.
type ProjectGroupSharedEnvResponse struct {
	Success bool                  `json:"success,omitempty"`
	Message string                `json:"message,omitempty"`
	Data    ProjectGroupSharedEnv `json:"data"`
}

// ProjectGroupInjectSharedEnvResponse is POST .../env/inject.
type ProjectGroupInjectSharedEnvResponse struct {
	Success bool                                `json:"success,omitempty"`
	Message string                              `json:"message,omitempty"`
	Data    InjectProjectGroupSharedEnvResponse `json:"data"`
}

// ProjectGroupConnectResponse is POST .../connections.
type ProjectGroupConnectResponse struct {
	Success bool                                `json:"success,omitempty"`
	Message string                              `json:"message,omitempty"`
	Data    ConnectProjectGroupServicesResponse `json:"data"`
}

// ProjectGroupRedeployAppsResponse is POST .../redeploy-apps.
type ProjectGroupRedeployAppsResponse struct {
	Success bool                             `json:"success,omitempty"`
	Message string                           `json:"message,omitempty"`
	Data    RedeployProjectGroupAppsResponse `json:"data"`
}

func withProjectGroupWorkspace(ctx context.Context, client *Client, path string, opts *ProjectGroupWorkspaceOptions) (string, error) {
	q := &ProjectGroupWorkspaceOptions{}
	if opts != nil {
		*q = *opts
	}
	if q.WorkspaceUUID == "" && q.Workspace == "" {
		if ws, _, err := firstWorkspaceUUID(ctx, client); err == nil {
			q.WorkspaceUUID = ws
		}
	}
	return addOptions(path, q)
}

// List returns project groups for a workspace.
// GET /project-groups?workspace_uuid=&limit=&offset=
func (s *ProjectGroupService) List(ctx context.Context, opts *ProjectGroupListOptions) (*ProjectGroupListResponse, *http.Response, error) {
	if opts == nil {
		opts = &ProjectGroupListOptions{}
	}
	if opts.WorkspaceUUID == "" && opts.Workspace == "" {
		if ws, _, err := firstWorkspaceUUID(ctx, s.client); err == nil {
			opts.WorkspaceUUID = ws
		}
	}
	u, err := addOptions("project-groups", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectGroupListResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Get returns one project group by UUID.
// GET /project-groups/:uuid?workspace_uuid=
func (s *ProjectGroupService) Get(ctx context.Context, uuid string, opts *ProjectGroupWorkspaceOptions) (*ProjectGroupResponse, *http.Response, error) {
	u, err := withProjectGroupWorkspace(ctx, s.client, fmt.Sprintf("project-groups/%s", uuid), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectGroupResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Create creates an empty project group.
// POST /project-groups?workspace_uuid=
func (s *ProjectGroupService) Create(ctx context.Context, body *CreateProjectGroupRequest, opts *ProjectGroupWorkspaceOptions) (*ProjectGroupResponse, *http.Response, error) {
	u, err := withProjectGroupWorkspace(ctx, s.client, "project-groups", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectGroupResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Update patches project group metadata.
// PATCH /project-groups/:uuid?workspace_uuid=
func (s *ProjectGroupService) Update(ctx context.Context, uuid string, body *UpdateProjectGroupRequest, opts *ProjectGroupWorkspaceOptions) (*ProjectGroupResponse, *http.Response, error) {
	u, err := withProjectGroupWorkspace(ctx, s.client, fmt.Sprintf("project-groups/%s", uuid), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPatch, u, body)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectGroupResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Delete removes a project group.
// DELETE /project-groups/:uuid?workspace_uuid=
func (s *ProjectGroupService) Delete(ctx context.Context, uuid string, opts *ProjectGroupWorkspaceOptions) (*http.Response, error) {
	u, err := withProjectGroupWorkspace(ctx, s.client, fmt.Sprintf("project-groups/%s", uuid), opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// AttachMember attaches a project or addon to a group.
// POST /project-groups/:uuid/members?workspace_uuid=
func (s *ProjectGroupService) AttachMember(ctx context.Context, uuid string, body *AttachProjectGroupMemberRequest, opts *ProjectGroupWorkspaceOptions) (*ProjectGroupAttachResponse, *http.Response, error) {
	u, err := withProjectGroupWorkspace(ctx, s.client, fmt.Sprintf("project-groups/%s/members", uuid), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectGroupAttachResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// DetachMember detaches a member from a group.
// DELETE /project-groups/:uuid/members/:memberType/:memberUUID?workspace_uuid=
func (s *ProjectGroupService) DetachMember(ctx context.Context, uuid, memberType, memberUUID string, opts *ProjectGroupDetachOptions) (*http.Response, error) {
	q := &ProjectGroupDetachOptions{}
	if opts != nil {
		*q = *opts
	}
	if q.WorkspaceUUID == "" && q.Workspace == "" {
		if ws, _, err := firstWorkspaceUUID(ctx, s.client); err == nil {
			q.WorkspaceUUID = ws
		}
	}
	u, err := addOptions(fmt.Sprintf("project-groups/%s/members/%s/%s", uuid, memberType, memberUUID), q)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// GetTopology returns the plane topology for a group.
// GET /project-groups/:uuid/topology?workspace_uuid=
func (s *ProjectGroupService) GetTopology(ctx context.Context, uuid string, opts *ProjectGroupWorkspaceOptions) (*ProjectGroupTopologyResponse, *http.Response, error) {
	u, err := withProjectGroupWorkspace(ctx, s.client, fmt.Sprintf("project-groups/%s/topology", uuid), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectGroupTopologyResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// GetSharedEnv returns group-level shared environment variables.
// GET /project-groups/:uuid/env?workspace_uuid=
func (s *ProjectGroupService) GetSharedEnv(ctx context.Context, uuid string, opts *ProjectGroupWorkspaceOptions) (*ProjectGroupSharedEnvResponse, *http.Response, error) {
	u, err := withProjectGroupWorkspace(ctx, s.client, fmt.Sprintf("project-groups/%s/env", uuid), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectGroupSharedEnvResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// PutSharedEnv replaces the group shared env set.
// PUT /project-groups/:uuid/env?workspace_uuid=
func (s *ProjectGroupService) PutSharedEnv(ctx context.Context, uuid string, body *UpsertProjectGroupSharedEnvRequest, opts *ProjectGroupWorkspaceOptions) (*ProjectGroupSharedEnvResponse, *http.Response, error) {
	u, err := withProjectGroupWorkspace(ctx, s.client, fmt.Sprintf("project-groups/%s/env", uuid), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPut, u, body)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectGroupSharedEnvResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// InjectSharedEnv pushes stored group shared env into project members.
// POST /project-groups/:uuid/env/inject?workspace_uuid=
func (s *ProjectGroupService) InjectSharedEnv(ctx context.Context, uuid string, body *InjectProjectGroupSharedEnvRequest, opts *ProjectGroupWorkspaceOptions) (*ProjectGroupInjectSharedEnvResponse, *http.Response, error) {
	u, err := withProjectGroupWorkspace(ctx, s.client, fmt.Sprintf("project-groups/%s/env/inject", uuid), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectGroupInjectSharedEnvResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ConnectServices wires provider connection envs into a consumer project.
// POST /project-groups/:uuid/connections?workspace_uuid=
func (s *ProjectGroupService) ConnectServices(ctx context.Context, uuid string, body *ConnectProjectGroupServicesRequest, opts *ProjectGroupWorkspaceOptions) (*ProjectGroupConnectResponse, *http.Response, error) {
	u, err := withProjectGroupWorkspace(ctx, s.client, fmt.Sprintf("project-groups/%s/connections", uuid), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectGroupConnectResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// RedeployApps queues redeploys for application (project) members only.
// POST /project-groups/:uuid/redeploy-apps?workspace_uuid=
func (s *ProjectGroupService) RedeployApps(ctx context.Context, uuid string, opts *ProjectGroupWorkspaceOptions) (*ProjectGroupRedeployAppsResponse, *http.Response, error) {
	u, err := withProjectGroupWorkspace(ctx, s.client, fmt.Sprintf("project-groups/%s/redeploy-apps", uuid), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectGroupRedeployAppsResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ResolveMember maps a service id to its group (deep links).
// GET /project-groups/resolve?workspace_uuid=&member_type=&member_uuid=
func (s *ProjectGroupService) ResolveMember(ctx context.Context, opts *ProjectGroupResolveOptions) (*ProjectGroupResolveResponse, *http.Response, error) {
	if opts == nil {
		opts = &ProjectGroupResolveOptions{}
	}
	if opts.WorkspaceUUID == "" && opts.Workspace == "" {
		if ws, _, err := firstWorkspaceUUID(ctx, s.client); err == nil {
			opts.WorkspaceUUID = ws
		}
	}
	u, err := addOptions("project-groups/resolve", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectGroupResolveResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListCandidates lists attachable projects/addons for the picker UI.
// GET /project-groups/candidates?workspace_uuid=&group_uuid=
func (s *ProjectGroupService) ListCandidates(ctx context.Context, opts *ProjectGroupCandidatesOptions) (*ProjectGroupCandidatesResponse, *http.Response, error) {
	if opts == nil {
		opts = &ProjectGroupCandidatesOptions{}
	}
	if opts.WorkspaceUUID == "" && opts.Workspace == "" {
		if ws, _, err := firstWorkspaceUUID(ctx, s.client); err == nil {
			opts.WorkspaceUUID = ws
		}
	}
	u, err := addOptions("project-groups/candidates", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(ProjectGroupCandidatesResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}
