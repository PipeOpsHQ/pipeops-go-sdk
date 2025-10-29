# PipeOps Go SDK Documentation

Welcome to the comprehensive documentation for the PipeOps Go SDK. This documentation is organized into several sections to help you get started quickly and dive deep into advanced features.

## 📚 Documentation Structure

### Getting Started
- **[Installation](getting-started/installation.md)** - Install the SDK in your Go project
- **[Quick Start](getting-started/quickstart.md)** - Your first API call in 5 minutes
- **[Configuration](getting-started/configuration.md)** - Configure the client for your needs

### Authentication
- **[Overview](authentication/overview.md)** - Authentication methods and best practices
- **[Basic Authentication](authentication/basic-auth.md)** - Email/password authentication
- **[OAuth 2.0](authentication/oauth.md)** - OAuth 2.0 integration guide

### API Services
Complete documentation for all API services:

- **[Overview](api-services/overview.md)** - Service architecture and patterns
- **[Auth](api-services/auth.md)** - User authentication and account management
- **[Projects](api-services/projects.md)** - Project deployment and management
- **[Servers](api-services/servers.md)** - Server/cluster management
- **[Environments](api-services/environments.md)** - Environment configuration
- **[Teams](api-services/teams.md)** - Team collaboration
- **[Workspaces](api-services/workspaces.md)** - Workspace organization
- **[Billing](api-services/billing.md)** - Subscription and payment management
- **[Add-ons](api-services/addons.md)** - Marketplace add-on deployment
- **[Webhooks](api-services/webhooks.md)** - Webhook configuration
- **[Users](api-services/users.md)** - User profile and settings
- **[Cloud Providers](api-services/cloudproviders.md)** - Cloud provider integration
- **[Service Tokens](api-services/servicetokens.md)** - Service account tokens
- **[Miscellaneous](api-services/misc.md)** - Utility endpoints

### Advanced Usage
- **[Error Handling](advanced/error-handling.md)** - Handle errors effectively
- **[Retries & Timeouts](advanced/retries-timeouts.md)** - Configure retries and timeouts
- **[Rate Limiting](advanced/rate-limiting.md)** - Handle API rate limits
- **[Logging](advanced/logging.md)** - Add logging for debugging
- **[Custom HTTP Client](advanced/custom-http-client.md)** - Use custom HTTP clients

### Examples
- **[Complete Examples](examples/complete-examples.md)** - Real-world usage examples
- **[Common Patterns](examples/common-patterns.md)** - Best practices and patterns

### API Reference
- **[Client](reference/client.md)** - Client type reference
- **[Types](reference/types.md)** - Common type definitions
- **[Errors](reference/errors.md)** - Error types and handling

## 🚀 Quick Start

Here's a minimal example to get you started:

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

## 📖 Building Documentation

### View Locally

Install MkDocs and serve the documentation locally:

```bash
# Install dependencies
pip install mkdocs mkdocs-material pymdown-extensions

# Serve documentation at http://127.0.0.1:8000
mkdocs serve
```

### Build Static Site

Build the documentation as static HTML:

```bash
mkdocs build
```

The built site will be in the `site/` directory.

## 🔗 Additional Resources

- **[GitHub Repository](https://github.com/PipeOpsHQ/pipeops-go-sdk)** - Source code and issues
- **[Go Package Documentation](https://godoc.org/github.com/PipeOpsHQ/pipeops-go-sdk/pipeops)** - GoDoc reference
- **[PipeOps API Docs](https://api.pipeops.io/docs)** - Official REST API documentation
- **[Working Examples](../examples/)** - Runnable code examples

## 🤝 Contributing

We welcome contributions to improve the documentation! Please see [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

## 📄 License

This SDK is distributed under the terms specified in the [LICENSE](../LICENSE) file.
