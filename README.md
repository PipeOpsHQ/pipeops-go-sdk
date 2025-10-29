# PipeOps Go SDK

A comprehensive Go SDK for interacting with the PipeOps Control Plane API.

**Status:** ✅ **100% API Coverage** - 284 methods implementing all 262 unique endpoints from the Postman collection

## Features

- **Complete API Coverage**: All API endpoints covered across 18 service modules
- **Type-Safe**: Strongly typed request/response structures
- **Context Support**: All methods support context for cancellation and timeouts
- **OAuth 2.0**: Full OAuth 2.0 authorization code flow support
- **Flexible**: Custom HTTP client support
- **Well-Documented**: Comprehensive examples and documentation

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

## API Coverage

This SDK provides **100% coverage** of the PipeOps API with 284 methods across 18 service modules:

### Core Services (100+ methods)
- **Projects** (46 methods) - Full CRUD, logs, metrics, network policies, GitHub/GitLab integration, env vars, deployment controls, costs, observability
- **Billing** (33 methods) - Cards, subscriptions, invoices, usage, balance, credit, history, plans, trials, portal, workspace billing
- **Servers/Clusters** (22 methods) - CRUD, service tokens, agent operations (register, heartbeat, poll, tunnel status), cost allocation
- **Cloud Providers** (17 methods) - AWS, GCP, Azure, DigitalOcean, Huawei - full account management, cost calculators, OAuth flows

### Organization & Access (30+ methods)
- **Teams** (11 methods) - Create, update, invite members, list, get, delete, member management
- **Admin** (20 methods) - User administration, statistics, plan management, waitlist programs, bulk operations, subscriptions
- **Workspaces** (6 methods) - Create, list, get, billing email management
- **User Settings** (8 methods) - Preferences, notifications, profile management, delete profile

### Integrations & Extensions (74+ methods)
- **Add-Ons** (21 methods) - Submit, deploy, manage deployments, configurations, sessions, domains, bulk operations
- **Webhooks** (8 methods) - Full CRUD operations
- **Events & Survey** (23 methods) - Event management, surveys, partners, agreements, participants, profile
- **DeploymentWebhooks** (3 methods) - GitHub, GitLab, Bitbucket webhooks
- **Campaign** (7 methods) - Waitlist, hackathon management
- **OAuth** (4 methods) - Full OAuth 2.0 authorization code flow

### Additional Services (42+ methods)
- **Service Tokens** (5 methods) - Full CRUD for service account tokens
- **Environments** (8 methods) - CRUD, environment variable management
- **Authentication** (10 methods) - Login, signup, verification, password management, OAuth signin
- **OpenCost** (3 methods) - Cluster and project cost metrics
- **Coupons** (2 methods) - Create, retrieve coupons
- **Various** - Logs, notifications, templates, integrations, health checks, backups, security scans, audit logs, alerts

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
