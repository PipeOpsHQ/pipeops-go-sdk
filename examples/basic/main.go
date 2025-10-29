package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
	// Create a new client with custom configuration
	client, err := pipeops.NewClient("",
		pipeops.WithTimeout(60*time.Second),  // Custom timeout
		pipeops.WithMaxRetries(5),            // Retry up to 5 times
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Get credentials from environment variables
	email := os.Getenv("PIPEOPS_EMAIL")
	password := os.Getenv("PIPEOPS_PASSWORD")

	if email == "" || password == "" {
		log.Fatal("Please set PIPEOPS_EMAIL and PIPEOPS_PASSWORD environment variables")
	}

	// Login to get an authentication token with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	loginReq := &pipeops.LoginRequest{
		Email:    email,
		Password: password,
	}

	loginResp, _, err := client.Auth.Login(ctx, loginReq)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	fmt.Printf("Login successful! User: %s\n", loginResp.Data.User.Email)

	// Set the token for authenticated requests
	client.SetToken(loginResp.Data.Token)

	// List all projects
	projects, _, err := client.Projects.List(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list projects: %v", err)
	}

	fmt.Printf("\nFound %d projects:\n", len(projects.Data.Projects))
	for _, project := range projects.Data.Projects {
		fmt.Printf("  - %s (%s)\n", project.Name, project.UUID)
	}

	// List all servers
	servers, _, err := client.Servers.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list servers: %v", err)
	}

	fmt.Printf("\nFound %d servers:\n", len(servers.Data.Servers))
	for _, server := range servers.Data.Servers {
		fmt.Printf("  - %s (%s) - Provider: %s\n", server.Name, server.UUID, server.Provider)
	}

	// List all environments
	environments, _, err := client.Environments.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list environments: %v", err)
	}

	fmt.Printf("\nFound %d environments:\n", len(environments.Data.Environments))
	for _, env := range environments.Data.Environments {
		fmt.Printf("  - %s (%s)\n", env.Name, env.UUID)
	}
}
