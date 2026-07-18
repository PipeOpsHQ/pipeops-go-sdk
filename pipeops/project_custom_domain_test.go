package pipeops

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFlexibleCSVString_UnmarshalStringAndArray(t *testing.T) {
	var fromString FlexibleCSVString
	if err := json.Unmarshal([]byte(`"https://a.example.com"`), &fromString); err != nil {
		t.Fatal(err)
	}
	if fromString.String() != "https://a.example.com" {
		t.Fatalf("string form = %q", fromString)
	}

	var fromArray FlexibleCSVString
	if err := json.Unmarshal([]byte(`["https://a.example.com","https://b.example.com"]`), &fromArray); err != nil {
		t.Fatal(err)
	}
	if fromArray.String() != "https://a.example.com,https://b.example.com" {
		t.Fatalf("array form = %q", fromArray)
	}
	if fromArray.First() != "https://a.example.com" {
		t.Fatalf("First = %q", fromArray.First())
	}
	if got := fromArray.All(); len(got) != 2 || got[1] != "https://b.example.com" {
		t.Fatalf("All = %#v", got)
	}

	var fromNull FlexibleCSVString
	if err := json.Unmarshal([]byte(`null`), &fromNull); err != nil {
		t.Fatal(err)
	}
	if fromNull.String() != "" {
		t.Fatalf("null = %q", fromNull)
	}
}

func TestProjectGet_CustomDomainNameArray(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := NewClient(server.URL, WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}

	mux.HandleFunc("/workspace", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    []map[string]string{{"UUID": "ws-1"}},
		}); err != nil {
			t.Fatal(err)
		}
	})

	// Mirrors controller project/fetch: CustomDomainName is a string array.
	mux.HandleFunc("/project/fetch/", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "ok",
			"data": map[string]interface{}{
				"project": map[string]interface{}{
					"UUID":             "p1",
					"Name":             "pipeops-hello-app",
					"CustomDomainName": []string{"https://pipeops-hello-app.example.com", ""},
					"public_url":       "https://pipeops-hello-app.example.com",
				},
			},
		}); err != nil {
			t.Fatal(err)
		}
	})

	resp, _, err := client.Projects.Get(context.Background(), "pipeops-hello-app", &ProjectGetOptions{WorkspaceUUID: "ws-1"})
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if !strings.Contains(resp.Data.Project.CustomDomainName.String(), "pipeops-hello-app.example.com") {
		t.Fatalf("CustomDomainName = %q", resp.Data.Project.CustomDomainName)
	}
	if resp.Data.Project.PublicURL != "https://pipeops-hello-app.example.com" {
		t.Fatalf("PublicURL = %q", resp.Data.Project.PublicURL)
	}
}
