# PipeOps Go SDK Documentation

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Authentication](#authentication)
- [Core Concepts](#core-concepts)
- [API Services](#api-services)
- [Error Handling](#error-handling)
- [Advanced Usage](#advanced-usage)
- [Reference](#reference)

## Installation

```bash
go get github.com/PipeOpsHQ/pipeops-go-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
    // Create a new client
    client := pipeops.NewClient("")
    
    // Login
    ctx := context.Background()
    loginResp, _, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
        Email:    "your-email@example.com",
        Password: "your-password",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Set token for authenticated requests
    client.SetToken(loginResp.Data.Token)
    
    // List projects
    projects, _, err := client.Projects.List(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d projects\n", len(projects.Data.Projects))
}
```

## Authentication

The SDK supports token-based authentication. After logging in, the token is automatically included in all subsequent requests.

```go
// Login to get a token
loginResp, _, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
    Email:    "your-email@example.com",
    Password: "your-password",
})

// Set the token
client.SetToken(loginResp.Data.Token)

// Or use an existing token
client.SetToken("your-existing-token")
```

## Core Concepts

### Client

The `Client` is the main entry point for interacting with the PipeOps API. It manages:
- HTTP communication
- Authentication tokens
- Service instances

### Services

The SDK is organized into services, each handling a specific area of the API:

- `Auth` - Authentication and user management
- `Projects` - Project management
- `Servers` - Server/cluster management
- `Environments` - Environment management
- `Teams` - Team and member management
- `Workspaces` - Workspace management
- `Billing` - Billing and subscription management
- `AddOns` - Add-on deployment and management
- `Webhooks` - Webhook configuration
- `Users` - User settings
- `Admin` - Administrative operations

### Context

All API methods require a `context.Context` parameter for cancellation and timeout control:

```go
ctx := context.Background()

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

projects, _, err := client.Projects.List(ctx, nil)
```

## API Services

### Authentication Service

```go
// Login
loginResp, _, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
    Email:    "email@example.com",
    Password: "password",
})

// Signup
signupResp, _, err := client.Auth.Signup(ctx, &pipeops.SignupRequest{
    Email:     "email@example.com",
    Password:  "password",
    FirstName: "John",
    LastName:  "Doe",
})

// Request password reset
resetResp, _, err := client.Auth.RequestPasswordReset(ctx, &pipeops.PasswordResetRequest{
    Email: "email@example.com",
})

// Change password
changeResp, _, err := client.Auth.ChangePassword(ctx, &pipeops.ChangePasswordRequest{
    OldPassword: "old-password",
    NewPassword: "new-password",
})
```

### Project Service

```go
// List all projects
projects, _, err := client.Projects.List(ctx, nil)

// List projects with filters
projects, _, err := client.Projects.List(ctx, &pipeops.ProjectListOptions{
    WorkspaceID: "workspace-uuid",
    Page:        1,
    Limit:       10,
})

// Get a specific project
project, _, err := client.Projects.Get(ctx, "project-uuid")

// Create a project
newProject, _, err := client.Projects.Create(ctx, &pipeops.CreateProjectRequest{
    Name:          "My Project",
    Description:   "Project description",
    ServerID:      "server-uuid",
    EnvironmentID: "environment-uuid",
    Repository:    "https://github.com/user/repo",
    Branch:        "main",
})

// Update a project
updatedProject, _, err := client.Projects.Update(ctx, "project-uuid", &pipeops.UpdateProjectRequest{
    Name:        "Updated Name",
    Description: "Updated description",
})

// Delete a project
_, err := client.Projects.Delete(ctx, "project-uuid")
```

### Server Service

```go
// List all servers
servers, _, err := client.Servers.List(ctx)

// Get a specific server
server, _, err := client.Servers.Get(ctx, "server-uuid")
```

### Environment Service

```go
// List all environments
environments, _, err := client.Environments.List(ctx)

// Get a specific environment
environment, _, err := client.Environments.Get(ctx, "environment-uuid")
```

## Error Handling

The SDK returns detailed error information:

```go
projects, resp, err := client.Projects.List(ctx, nil)
if err != nil {
    if errResp, ok := err.(*pipeops.ErrorResponse); ok {
        fmt.Printf("Status Code: %d\n", errResp.Response.StatusCode)
        fmt.Printf("Error Message: %s\n", errResp.Message)
    } else {
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

## Advanced Usage

### Custom HTTP Client

You can provide a custom HTTP client for advanced configuration:

```go
import (
    "net/http"
    "time"
)

httpClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        10,
        IdleConnTimeout:     30 * time.Second,
        DisableCompression:  true,
    },
}

client := pipeops.NewClient("")
client.SetHTTPClient(httpClient)
```

### Custom Base URL

For testing or using a different API endpoint:

```go
client := pipeops.NewClient("https://staging-api.pipeops.io")
```

### Request/Response Inspection

Access the raw HTTP response:

```go
projects, httpResp, err := client.Projects.List(ctx, nil)
if httpResp != nil {
    fmt.Printf("Status: %d\n", httpResp.StatusCode)
    fmt.Printf("Headers: %v\n", httpResp.Header)
}
```

## Reference

### Client Configuration

```go
type Client struct {
    BaseURL   *url.URL
    UserAgent string
    // ... other fields
}

// Create with defaults
client := pipeops.NewClient("")

// Create with custom base URL
client := pipeops.NewClient("https://api.pipeops.io")

// Set token
client.SetToken("your-token")

// Set custom HTTP client
client.SetHTTPClient(customHTTPClient)
```

### Common Types

```go
// Timestamp - handles various date/time formats
type Timestamp struct {
    time.Time
}

// User
type User struct {
    ID            string
    UUID          string
    Email         string
    FirstName     string
    LastName      string
    IsActive      bool
    EmailVerified bool
    CreatedAt     *Timestamp
    UpdatedAt     *Timestamp
}

// Project
type Project struct {
    ID            string
    UUID          string
    Name          string
    Description   string
    Status        string
    ServerID      string
    EnvironmentID string
    WorkspaceID   string
    Repository    string
    Branch        string
    CreatedAt     *Timestamp
    UpdatedAt     *Timestamp
}
```

## Examples

See the [examples](../examples) directory for more detailed usage examples.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines on contributing to this SDK.

## License

This SDK is distributed under the MIT License. See [LICENSE](../LICENSE) for more information.
