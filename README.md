# PipeOps Go SDK

A comprehensive Go SDK for interacting with the PipeOps Control Plane API.

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
- **Authentication** - Login, signup, password management
- **OAuth 2.0** - Authorization code flow, token exchange, user info
- **Projects** - Full CRUD operations with filtering
- **Servers/Clusters** - Server and cluster management
- **Environments** - Environment configuration

### Organization & Access
- **Teams** - Team creation, member management, invitations
- **Workspaces** - Workspace management and collaboration

### Add-Ons & Extensions
- **Add-Ons** - Browse, deploy, and manage add-ons
- **Webhooks** - Webhook configuration and management

### Billing & Usage
- **Billing** - Payment cards, subscriptions, invoices, usage tracking

### User Management
- **User Settings** - Preferences, notifications, profile management
- **Admin** - User administration, statistics, plan management (admin only)

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
