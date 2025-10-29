# Projects Service

The Projects Service manages application projects, deployments, and related operations.

## Overview

```go
// Access the projects service
projectsService := client.Projects
```

## Methods

### List Projects

List all projects with optional filtering:

```go
// List all projects
projects, _, err := client.Projects.List(ctx, nil)
if err != nil {
    log.Fatalf("Failed to list projects: %v", err)
}

fmt.Printf("Found %d projects\n", len(projects.Data.Projects))
for _, project := range projects.Data.Projects {
    fmt.Printf("- %s (%s)\n", project.Name, project.UUID)
}
```

With filters:

```go
projects, _, err := client.Projects.List(ctx, &pipeops.ProjectListOptions{
    WorkspaceID: "workspace-uuid",
    ServerID:    "server-uuid",
    Page:        1,
    Limit:       20,
})
```

### Get Project

Get a specific project by UUID:

```go
project, _, err := client.Projects.Get(ctx, "project-uuid")
if err != nil {
    log.Fatalf("Failed to get project: %v", err)
}

fmt.Printf("Project: %s\n", project.Data.Project.Name)
fmt.Printf("Status: %s\n", project.Data.Project.Status)
fmt.Printf("Repository: %s\n", project.Data.Project.Repository)
```

### Create Project

Create a new project:

```go
newProject, _, err := client.Projects.Create(ctx, &pipeops.CreateProjectRequest{
    Name:          "My Application",
    Description:   "Production application",
    ServerID:      "server-uuid",
    EnvironmentID: "environment-uuid",
    Repository:    "https://github.com/user/repo",
    Branch:        "main",
    BuildCommand:  "npm run build",
    StartCommand:  "npm start",
    Port:          3000,
    Framework:     "nodejs",
    EnvVars: map[string]interface{}{
        "NODE_ENV": "production",
        "API_KEY":  "secret-key",
    },
})
if err != nil {
    log.Fatalf("Failed to create project: %v", err)
}

fmt.Printf("Created project: %s\n", newProject.Data.Project.UUID)
```

### Update Project

Update an existing project:

```go
updated, _, err := client.Projects.Update(ctx, projectUUID, &pipeops.UpdateProjectRequest{
    Name:         "Updated Name",
    Description:  "Updated description",
    BuildCommand: "yarn build",
    StartCommand: "yarn start",
    Port:         8080,
})
if err != nil {
    log.Fatalf("Failed to update project: %v", err)
}

fmt.Println("Project updated successfully")
```

### Delete Project

Delete a project:

```go
_, err := client.Projects.Delete(ctx, "project-uuid")
if err != nil {
    log.Fatalf("Failed to delete project: %v", err)
}

fmt.Println("Project deleted successfully")
```

### Deploy Project

Trigger a deployment:

```go
deployment, _, err := client.Projects.Deploy(ctx, "project-uuid")
if err != nil {
    log.Fatalf("Deployment failed: %v", err)
}

fmt.Printf("Deployment started: %s\n", deployment.Data.DeploymentID)
```

### Get Project Logs

Retrieve project logs:

```go
logs, _, err := client.Projects.GetLogs(ctx, projectUUID, &pipeops.LogsOptions{
    Limit:  100,
    Search: "error",
})
if err != nil {
    log.Fatalf("Failed to get logs: %v", err)
}

for _, logEntry := range logs.Data.Logs {
    fmt.Printf("Log: %v\n", logEntry)
}
```

### Get Environment Variables

Get project environment variables:

```go
envVars, _, err := client.Projects.GetEnvVars(ctx, "project-uuid")
if err != nil {
    log.Fatalf("Failed to get env vars: %v", err)
}

for key, value := range envVars.Data.EnvVars {
    fmt.Printf("%s=%s\n", key, value)
}
```

### Update Environment Variables

Update project environment variables:

```go
_, err := client.Projects.UpdateEnvVars(ctx, projectUUID, &pipeops.UpdateEnvVarsRequest{
    EnvVars: map[string]string{
        "DATABASE_URL": "postgresql://...",
        "REDIS_URL":    "redis://...",
    },
})
if err != nil {
    log.Fatalf("Failed to update env vars: %v", err)
}

fmt.Println("Environment variables updated")
```

### Restart Project

Restart a project:

```go
_, err := client.Projects.Restart(ctx, "project-uuid")
if err != nil {
    log.Fatalf("Failed to restart: %v", err)
}

fmt.Println("Project restarted")
```

### Get GitHub Branches

Get available branches from a GitHub repository:

```go
branches, _, err := client.Projects.GetGitHubBranches(ctx, &pipeops.GitHubBranchesRequest{
    Repository: "https://github.com/user/repo",
})
if err != nil {
    log.Fatalf("Failed to get branches: %v", err)
}

for _, branch := range branches.Data.Branches {
    fmt.Printf("Branch: %s\n", branch)
}
```

### Update Domain

Update project domain:

```go
domain, _, err := client.Projects.UpdateDomain(ctx, projectUUID, &pipeops.DomainRequest{
    Domain: "myapp.com",
})
if err != nil {
    log.Fatalf("Failed to update domain: %v", err)
}

fmt.Printf("Domain updated: %s\n", domain.Data.Domain)
```

## Data Types

### Project

```go
type Project struct {
    ID            string     `json:"id,omitempty"`
    UUID          string     `json:"uuid,omitempty"`
    Name          string     `json:"name,omitempty"`
    Description   string     `json:"description,omitempty"`
    Status        string     `json:"status,omitempty"`
    ServerID      string     `json:"server_id,omitempty"`
    EnvironmentID string     `json:"environment_id,omitempty"`
    WorkspaceID   string     `json:"workspace_id,omitempty"`
    Repository    string     `json:"repository,omitempty"`
    Branch        string     `json:"branch,omitempty"`
    BuildCommand  string     `json:"build_command,omitempty"`
    StartCommand  string     `json:"start_command,omitempty"`
    Port          int        `json:"port,omitempty"`
    Framework     string     `json:"framework,omitempty"`
    CreatedAt     *Timestamp `json:"created_at,omitempty"`
    UpdatedAt     *Timestamp `json:"updated_at,omitempty"`
}
```

### ProjectListOptions

```go
type ProjectListOptions struct {
    WorkspaceID string `url:"workspace_id,omitempty"`
    ServerID    string `url:"server_id,omitempty"`
    Page        int    `url:"page,omitempty"`
    Limit       int    `url:"limit,omitempty"`
}
```

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
    client, _ := pipeops.NewClient("")
    
    // Authenticate
    loginResp, _, _ := client.Auth.Login(ctx, &pipeops.LoginRequest{
        Email:    "user@example.com",
        Password: "password",
    })
    client.SetToken(loginResp.Data.Token)
    
    ctx := context.Background()
    
    // Create a new project
    project, _, err := client.Projects.Create(ctx, &pipeops.CreateProjectRequest{
        Name:          "My Web App",
        ServerID:      "server-uuid",
        EnvironmentID: "env-uuid",
        Repository:    "https://github.com/user/webapp",
        Branch:        "main",
        BuildCommand:  "npm run build",
        StartCommand:  "npm start",
        Port:          3000,
    })
    if err != nil {
        log.Fatalf("Failed to create project: %v", err)
    }
    
    projectUUID := project.Data.Project.UUID
    fmt.Printf("Created project: %s\n", projectUUID)
    
    // Deploy the project
    deployment, _, err := client.Projects.Deploy(ctx, projectUUID)
    if err != nil {
        log.Fatalf("Deployment failed: %v", err)
    }
    
    fmt.Printf("Deployment started: %s\n", deployment.Data.DeploymentID)
    
    // Get logs
    logs, _, err := client.Projects.GetLogs(ctx, projectUUID, nil)
    if err != nil {
        log.Fatalf("Failed to get logs: %v", err)
    }
    
    fmt.Printf("Retrieved %d log entries\n", len(logs.Data.Logs))
}
```

## See Also

- [Servers Service](servers.md) - Manage servers for projects
- [Environments Service](environments.md) - Configure project environments
- [Webhooks Service](webhooks.md) - Set up deployment webhooks
