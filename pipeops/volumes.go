package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// VolumeService handles workspace volume inventory and recovery APIs.
//
// Controller routes (JWT + team access; dual-auth SA may follow later):
//
//	GET    /volumes
//	GET    /volumes/:uuid
//	POST   /volumes/:uuid/remount
//	DELETE /volumes/:uuid
//	POST   /volumes/:uuid/export
//	GET    /volumes/:uuid/export
type VolumeService struct {
	client *Client
}

// Volume is the API representation of a workspace volume.
type Volume struct {
	UUID                   string   `json:"uuid,omitempty"`
	DisplayName            string   `json:"display_name,omitempty"`
	PVCName                string   `json:"pvc_name,omitempty"`
	MountPath              string   `json:"mount_path,omitempty"`
	SizeGB                 float32  `json:"size_gb,omitempty"`
	Status                 string   `json:"status,omitempty"`
	ClusterUUID            string   `json:"cluster_uuid,omitempty"`
	ClusterName            string   `json:"cluster_name,omitempty"`
	Namespace              string   `json:"namespace,omitempty"`
	OwnerType              string   `json:"owner_type,omitempty"`
	OwnerUUID              string   `json:"owner_uuid,omitempty"`
	OwnerName              string   `json:"owner_name,omitempty"`
	OwnerSessionID         string   `json:"owner_session_id,omitempty"`
	OwnerHref              string   `json:"owner_href,omitempty"`
	RetainedFromOwnerUUID  string   `json:"retained_from_owner_uuid,omitempty"`
	RetainedFromOwnerName  string   `json:"retained_from_owner_name,omitempty"`
	RetainedFromOwnerType  string   `json:"retained_from_owner_type,omitempty"`
	RetainedUntil          *string  `json:"retained_until,omitempty"`
	OriginalDeploymentName string   `json:"original_deployment_name,omitempty"`
	ExportStatus           string   `json:"export_status,omitempty"`
	ExportURL              string   `json:"export_url,omitempty"`
	ExportError            string   `json:"export_error,omitempty"`
	ExportFilename         string   `json:"export_filename,omitempty"`
	Actions                []string `json:"actions,omitempty"`
	CreatedAt              string   `json:"created_at,omitempty"`
	UpdatedAt              string   `json:"updated_at,omitempty"`
}

// VolumeSummary counts volumes by status.
type VolumeSummary struct {
	Mounted    int64 `json:"mounted"`
	Unattached int64 `json:"unattached"`
}

// VolumeListOptions filters and paginates volume list.
type VolumeListOptions struct {
	WorkspaceUUID string `url:"workspace_uuid,omitempty"`
	// Workspace is accepted as an alias by the controller.
	Workspace   string `url:"workspace,omitempty"`
	Status      string `url:"status,omitempty"`
	ClusterUUID string `url:"cluster_uuid,omitempty"`
	Limit       int    `url:"limit,omitempty"`
	Offset      int    `url:"offset,omitempty"`
}

// VolumeListResponse is GET /volumes.
type VolumeListResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Volumes []Volume      `json:"volumes"`
		Summary VolumeSummary `json:"summary"`
		Total   int64         `json:"total"`
		Limit   int           `json:"limit"`
		Offset  int           `json:"offset"`
	} `json:"data"`
}

// VolumeResponse is GET /volumes/:uuid.
type VolumeResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    Volume `json:"data"`
}

// RemountVolumeRequest remounts an unattached volume onto a live resource.
type RemountVolumeRequest struct {
	TargetType string `json:"target_type"` // project | addon
	TargetUUID string `json:"target_uuid"`
	MountPath  string `json:"mount_path,omitempty"`
}

// RemountVolumeResponse is POST /volumes/:uuid/remount.
type RemountVolumeResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Volume  Volume `json:"volume"`
		Message string `json:"message,omitempty"`
	} `json:"data"`
}

// VolumeExportResponse tracks async export status.
type VolumeExportResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		UUID        string `json:"uuid,omitempty"`
		Status      string `json:"status,omitempty"`
		DownloadURL string `json:"download_url,omitempty"`
		Filename    string `json:"filename,omitempty"`
		Error       string `json:"error,omitempty"`
		Message     string `json:"message,omitempty"`
	} `json:"data"`
}

// List returns workspace volumes.
// GET /volumes?workspace_uuid=
func (s *VolumeService) List(ctx context.Context, opts *VolumeListOptions) (*VolumeListResponse, *http.Response, error) {
	u := "volumes"
	if opts == nil {
		opts = &VolumeListOptions{}
	}
	if opts.WorkspaceUUID == "" && opts.Workspace == "" {
		if ws, _, err := firstWorkspaceUUID(ctx, s.client); err == nil {
			opts.WorkspaceUUID = ws
		}
	}
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(VolumeListResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Get returns one volume by UUID.
// GET /volumes/:uuid?workspace_uuid=
func (s *VolumeService) Get(ctx context.Context, volumeUUID string, opts *VolumeListOptions) (*VolumeResponse, *http.Response, error) {
	u := fmt.Sprintf("volumes/%s", volumeUUID)
	u, err := withVolumeWorkspaceQuery(ctx, s.client, u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(VolumeResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Remount schedules remounting an unattached volume onto a project or addon.
// POST /volumes/:uuid/remount?workspace_uuid=
func (s *VolumeService) Remount(ctx context.Context, volumeUUID string, body *RemountVolumeRequest, opts *VolumeListOptions) (*RemountVolumeResponse, *http.Response, error) {
	u := fmt.Sprintf("volumes/%s/remount", volumeUUID)
	u, err := withVolumeWorkspaceQuery(ctx, s.client, u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}

	out := new(RemountVolumeResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Delete permanently deletes a volume.
// DELETE /volumes/:uuid?workspace_uuid=
func (s *VolumeService) Delete(ctx context.Context, volumeUUID string, opts *VolumeListOptions) (*http.Response, error) {
	u := fmt.Sprintf("volumes/%s", volumeUUID)
	u, err := withVolumeWorkspaceQuery(ctx, s.client, u, opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// StartExport starts an async volume export.
// POST /volumes/:uuid/export?workspace_uuid=
func (s *VolumeService) StartExport(ctx context.Context, volumeUUID string, opts *VolumeListOptions) (*VolumeExportResponse, *http.Response, error) {
	u := fmt.Sprintf("volumes/%s/export", volumeUUID)
	u, err := withVolumeWorkspaceQuery(ctx, s.client, u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(VolumeExportResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// GetExport polls export status for a volume.
// GET /volumes/:uuid/export?workspace_uuid=
func (s *VolumeService) GetExport(ctx context.Context, volumeUUID string, opts *VolumeListOptions) (*VolumeExportResponse, *http.Response, error) {
	u := fmt.Sprintf("volumes/%s/export", volumeUUID)
	u, err := withVolumeWorkspaceQuery(ctx, s.client, u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(VolumeExportResponse)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

func withVolumeWorkspaceQuery(ctx context.Context, client *Client, path string, opts *VolumeListOptions) (string, error) {
	q := &VolumeListOptions{}
	if opts != nil {
		*q = *opts
	}
	if q.WorkspaceUUID == "" && q.Workspace == "" {
		if ws, _, err := firstWorkspaceUUID(ctx, client); err == nil {
			q.WorkspaceUUID = ws
		}
	}
	// Only workspace-related query params on detail endpoints.
	type wsOnly struct {
		WorkspaceUUID string `url:"workspace_uuid,omitempty"`
		Workspace     string `url:"workspace,omitempty"`
	}
	return addOptions(path, &wsOnly{WorkspaceUUID: q.WorkspaceUUID, Workspace: q.Workspace})
}
