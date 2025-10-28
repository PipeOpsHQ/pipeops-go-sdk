# PipeOps Go SDK

A comprehensive Go SDK for interacting with the PipeOps Control Plane API.

**Status:** 145+ endpoints implemented across 15 services (50% of 288 total API endpoints)

## Installation

```bash
go get github.com/PipeOpsHQ/pipeops-go-sdk
```

## Usage

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
    client := pipeops.NewClient("https://api.pipeops.io")
    
    // Login to get an authentication token
    ctx := context.Background()
    loginReq := &pipeops.LoginRequest{
        Email:    "your-email@example.com",
        Password: "your-password",
    }
    
    resp, _, err := client.Auth.Login(ctx, loginReq)
    if err != nil {
        log.Fatal(err)
    }
    
    // Set the token for authenticated requests
    client.SetToken(resp.Data.Token)
    
    // Now you can make authenticated API calls
    projects, _, _ := client.Projects.List(ctx, nil)
    fmt.Printf("Found %d projects\n", len(projects.Data.Projects))
}
```

## Features

This SDK provides access to the following PipeOps API categories:

### Core Services
- **Authentication** (4 endpoints) - Login, signup, password management
- **OAuth 2.0** (4 endpoints) - Authorization code flow, token exchange, user info
- **Projects** (16 endpoints) - Full CRUD, logs, metrics, GitHub integration, env vars, deployment controls, costs
- **Servers/Clusters** (14 endpoints) - CRUD, service tokens, agent operations, tunnel info, cost allocation
- **Environments** (6 endpoints) - CRUD, environment variable management

### Organization & Access
- **Teams** (6 endpoints) - Team creation, member management, invitations
- **Workspaces** (3 endpoints) - Workspace management and collaboration

### Add-Ons & Extensions
- **Add-Ons** (16 endpoints) - Submit, deploy, manage deployments, configurations, sessions, bulk operations
- **Webhooks** (5 endpoints) - Webhook configuration and management

### Billing & Usage
- **Billing** (22 endpoints) - Cards, subscriptions, invoices, usage, balance, credit, history, plans, trials, portal

### Cloud Providers
- **Cloud Providers** (13 endpoints) - AWS, GCP, Azure, DigitalOcean, Huawei account management

### User Management & Administration
- **User Settings** (5 endpoints) - Preferences, notifications, profile management
- **Admin** (16 endpoints) - User administration, statistics, plan management, waitlist programs, bulk operations

### Additional Services
- **Events & Survey** (15 endpoints) - Event management, surveys, partners, contact us, waitlist, dashboard

## OAuth 2.0 Support

The SDK includes full support for OAuth 2.0 authorization code flow:

```go
// Generate authorization URL
authURL, _ := client.OAuth.Authorize(&pipeops.AuthorizeOptions{
    ClientID:     "your-client-id",
    RedirectURI:  "http://localhost:3000/callback",
    ResponseType: "code",
    Scope:        "user:read user:write",
    State:        "random-state",
})

// Exchange code for token
token, _, _ := client.OAuth.ExchangeCodeForToken(ctx, &pipeops.TokenRequest{
    GrantType:    "authorization_code",
    Code:         authCode,
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
})

client.SetToken(token.AccessToken)

// Get user info
userInfo, _, _ := client.OAuth.GetUserInfo(ctx)
```

See `examples/oauth/` for a complete OAuth flow example.

## Documentation

For detailed API documentation, please refer to:
- [API Documentation](docs/README.md) - Comprehensive SDK documentation
- [Examples](examples/) - Working code examples
- [PipeOps API Documentation](https://api.pipeops.io/docs) - Official API docs

## Examples

- [Basic Usage](examples/basic/) - Authentication and basic API calls
- [OAuth Flow](examples/oauth/) - Complete OAuth 2.0 authorization example

## License

This SDK is distributed under the terms of the license specified in the LICENSE file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.
