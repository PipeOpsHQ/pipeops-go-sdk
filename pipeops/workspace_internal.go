package pipeops

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

func fetchWorkspaceList(ctx context.Context, client *Client) ([]workspaceListItem, *http.Response, error) {
	req, err := client.NewRequest(http.MethodGet, "workspace", nil)
	if err != nil {
		return nil, nil, err
	}

	rawResp := new(workspaceListEnvelope)
	resp, err := client.Do(ctx, req, rawResp)
	if err != nil {
		return nil, resp, err
	}

	var workspaces []workspaceListItem
	if len(rawResp.Data) == 0 {
		return workspaces, resp, nil
	}

	if err := json.Unmarshal(rawResp.Data, &workspaces); err == nil {
		return workspaces, resp, nil
	}

	var wrapped struct {
		Workspaces []workspaceListItem `json:"workspaces,omitempty"`
	}
	if err := json.Unmarshal(rawResp.Data, &wrapped); err == nil {
		return wrapped.Workspaces, resp, nil
	}

	return nil, resp, errors.New("failed to decode workspace list response")
}

func firstWorkspaceUUID(ctx context.Context, client *Client) (string, *http.Response, error) {
	workspaces, resp, err := fetchWorkspaceList(ctx, client)
	if err != nil {
		return "", resp, err
	}
	if len(workspaces) == 0 || workspaces[0].UUID == "" {
		return "", resp, errors.New("no workspaces found for the authenticated user")
	}
	return workspaces[0].UUID, resp, nil
}
