package pipeops

import "testing"

func TestApplyCreateProjectDefaults_PrefersClient(t *testing.T) {
	worker := true
	req := &CreateProjectRequest{
		Source:      "gitlab",
		Environment: "staging",
		EnvVariables: []CreateProjectEnvVar{
			{Key: "PORT", Value: "8080"},
			{Key: "APP_ENV", Value: "prod"},
		},
		BuildSettings: CreateProjectBuildSettings{Worker: &worker},
		NetworkSettings: []CreateProjectNetworkSetting{
			{Port: 3000, Protocol: "TCP"},
		},
		CommitURL:  "https://example.com/commit/abc",
		Repository: "https://example.com/repo",
	}
	ApplyCreateProjectDefaults(req)
	if req.Source != "gitlab" || req.Environment != "staging" {
		t.Fatalf("source/env overridden: %q %q", req.Source, req.Environment)
	}
	if len(req.EnvVariables) != 2 || req.EnvVariables[0].Value != "8080" {
		t.Fatalf("client PORT should be kept: %+v", req.EnvVariables)
	}
	if req.NetworkSettings[0].Protocol != "TCP" {
		t.Fatalf("protocol overridden: %q", req.NetworkSettings[0].Protocol)
	}
	if req.CommitURL != "https://example.com/commit/abc" {
		t.Fatalf("commitURL overridden: %q", req.CommitURL)
	}
	if req.BuildSettings.Worker == nil || !*req.BuildSettings.Worker {
		t.Fatal("worker should stay true")
	}
}

func TestApplyCreateProjectDefaults_FillsGapsOnly(t *testing.T) {
	req := &CreateProjectRequest{
		Repository: "https://github.com/acme/app",
		NetworkSettings: []CreateProjectNetworkSetting{
			{Port: 3000},
		},
	}
	ApplyCreateProjectDefaults(req)
	if req.Source != "github" || req.Environment != "development" {
		t.Fatalf("defaults: source=%q env=%q", req.Source, req.Environment)
	}
	if req.EnvVariables == nil {
		t.Fatal("envVariables nil")
	}
	found := false
	for _, e := range req.EnvVariables {
		if e.Key == "PORT" && e.Value == "3000" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected PORT from network: %+v", req.EnvVariables)
	}
	if req.NetworkSettings[0].Protocol != "HTTP" {
		t.Fatalf("protocol = %q", req.NetworkSettings[0].Protocol)
	}
	if req.CommitURL != req.Repository {
		t.Fatalf("commitURL soft default = %q", req.CommitURL)
	}
}
