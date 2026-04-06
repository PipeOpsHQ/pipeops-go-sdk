package pipeops

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCloudProviderService_DiscoveryEndpoints(t *testing.T) {
	t.Parallel()

	step := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		step++
		w.Header().Set("Content-Type", "application/json")
		switch step {
		case 1:
			if r.URL.Path != "/app/digital_ocean/regions" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/app/digital_ocean/regions")
			}
			if _, err := w.Write([]byte(`{"success":true,"data":{"digital_ocean":[{"title":"NYC1","value":"nyc1","code":"us"}]}}`)); err != nil {
				t.Fatalf("write response error: %v", err)
			}
		case 2:
			if r.URL.Path != "/app/aws/instance-categories" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/app/aws/instance-categories")
			}
			if _, err := w.Write([]byte(`{"success":true,"data":{"aws":{"instanceCategories":["General purpose"]}}}`)); err != nil {
				t.Fatalf("write response error: %v", err)
			}
		case 3:
			if r.URL.Path != "/app/digital_ocean/instance-types" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/app/digital_ocean/instance-types")
			}
			if got := r.URL.Query().Get("instanceClass"); got != "Basic" {
				t.Fatalf("instanceClass = %q, want %q", got, "Basic")
			}
			if got := r.URL.Query().Get("region"); got != "nyc1" {
				t.Fatalf("region = %q, want %q", got, "nyc1")
			}
			if _, err := w.Write([]byte(`{"success":true,"data":{"digital_ocean":{"Basic":[{"name":"s-1vcpu-1gb","vcpu":1}]}}}`)); err != nil {
				t.Fatalf("write response error: %v", err)
			}
		case 4:
			if r.URL.Path != "/app/gcp/server-templates" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/app/gcp/server-templates")
			}
			if _, err := w.Write([]byte(`{"success":true,"data":{"development":[{"uuid":"tpl1","package":"e2-small"}]}}`)); err != nil {
				t.Fatalf("write response error: %v", err)
			}
		default:
			t.Fatalf("unexpected call %d", step)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	regions, _, err := client.CloudProviders.ListRegions(context.Background(), "digital_ocean")
	if err != nil {
		t.Fatalf("ListRegions error: %v", err)
	}
	if len(regions.Data["digital_ocean"]) != 1 {
		t.Fatalf("len(regions) = %d, want 1", len(regions.Data["digital_ocean"]))
	}

	categories, _, err := client.CloudProviders.ListInstanceCategories(context.Background(), "aws")
	if err != nil {
		t.Fatalf("ListInstanceCategories error: %v", err)
	}
	if len(categories.Data["aws"].InstanceCategories) != 1 {
		t.Fatalf("len(instanceCategories) = %d, want 1", len(categories.Data["aws"].InstanceCategories))
	}

	instanceTypes, _, err := client.CloudProviders.ListInstanceTypes(context.Background(), "digital_ocean", &CloudProviderInstanceTypesOptions{InstanceClass: "Basic", Region: "nyc1"})
	if err != nil {
		t.Fatalf("ListInstanceTypes error: %v", err)
	}
	if len(instanceTypes.Data["digital_ocean"]["Basic"]) != 1 {
		t.Fatalf("len(instanceTypes) = %d, want 1", len(instanceTypes.Data["digital_ocean"]["Basic"]))
	}

	templates, _, err := client.CloudProviders.ListServerTemplates(context.Background(), "gcp")
	if err != nil {
		t.Fatalf("ListServerTemplates error: %v", err)
	}
	if len(templates.Data["development"]) != 1 {
		t.Fatalf("len(templates) = %d, want 1", len(templates.Data["development"]))
	}
}
