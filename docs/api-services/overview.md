# API Services Overview

The PipeOps Go SDK is organized into specialized services, each handling a specific area of the PipeOps API. This modular approach makes it easy to discover and use the functionality you need.

## Available Services

| Service | Description | Documentation |
|---------|-------------|---------------|
| **Auth** | User authentication, signup, password management | [Guide](auth.md) |
| **OAuth** | OAuth 2.0 authorization flows | [OAuth Guide](../authentication/oauth.md) |
| **Projects** | Project creation, deployment, and management | [Guide](projects.md) |
| **Servers** | Server/cluster provisioning and management | [Guide](servers.md) |
| **Environments** | Environment configuration and management | [Guide](environments.md) |
| **Teams** | Team collaboration and member management | [Guide](teams.md) |
| **Workspaces** | Workspace organization and settings | [Guide](workspaces.md) |
| **Billing** | Subscription, payment, and invoice management | [Guide](billing.md) |
| **AddOns** | Marketplace add-on deployment | [Guide](addons.md) |
| **Webhooks** | Webhook configuration and delivery management | [Guide](webhooks.md) |
| **Users** | User profile and settings | [Guide](users.md) |
| **CloudProviders** | Cloud provider integration | [Guide](cloudproviders.md) |
| **ServiceTokens** | Service account token management | [Guide](servicetokens.md) |
| **Misc** | Miscellaneous utilities | [Guide](misc.md) |

## Service Architecture

All services follow a consistent pattern:

```go
// Access services through the client
client, _ := pipeops.NewClient("")
client.SetToken("your-token")

// Each service has its own methods
projects, _, err := client.Projects.List(ctx, nil)
servers, _, err := client.Servers.List(ctx)
teams, _, err := client.Teams.List(ctx)
```

## Common Patterns

### Context Usage

All service methods require a context parameter:

```go
ctx := context.Background()

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

projects, _, err := client.Projects.List(ctx, nil)
```

### Return Values

Most methods return three values:

```go
data, response, err := client.Projects.Get(ctx, "project-id")

// data: Typed response data
// response: Raw *http.Response
// err: Error if any
```

### Options and Filters

Many list methods accept options for filtering and pagination:

```go
projects, _, err := client.Projects.List(ctx, &pipeops.ProjectListOptions{
    WorkspaceID: "workspace-uuid",
    Page:        1,
    Limit:       20,
})
```

### Error Handling

Always check for errors:

```go
projects, resp, err := client.Projects.List(ctx, nil)
if err != nil {
    if resp != nil && resp.StatusCode == 401 {
        log.Println("Authentication required")
    }
    return err
}
```

## Quick Examples by Service

### Authentication

```go
// Login
loginResp, _, _ := client.Auth.Login(ctx, &pipeops.LoginRequest{
    Email:    "user@example.com",
    Password: "password",
})
client.SetToken(loginResp.Data.Token)
```

### Projects

```go
// List projects
projects, _, _ := client.Projects.List(ctx, nil)

// Get specific project
project, _, _ := client.Projects.Get(ctx, "project-uuid")

// Create project
newProject, _, _ := client.Projects.Create(ctx, &pipeops.CreateProjectRequest{
    Name:       "My App",
    ServerID:   "server-uuid",
    Repository: "https://github.com/user/repo",
})
```

### Servers

```go
// List servers
servers, _, _ := client.Servers.List(ctx)

// Get server details
server, _, _ := client.Servers.Get(ctx, "server-uuid")

// Create server
newServer, _, _ := client.Servers.Create(ctx, &pipeops.CreateServerRequest{
    Name:     "Production Server",
    Provider: "aws",
    Region:   "us-east-1",
})
```

### Teams

```go
// List teams
teams, _, _ := client.Teams.List(ctx)

// Create team
team, _, _ := client.Teams.Create(ctx, &pipeops.CreateTeamRequest{
    Name:        "Development Team",
    Description: "Core developers",
})

// Invite member
_, _ := client.Teams.InviteMember(ctx, teamUUID, &pipeops.InviteTeamMemberRequest{
    Email: "member@example.com",
    Role:  "developer",
})
```

### Billing

```go
// Get current balance
balance, _, _ := client.Billing.GetBalance(ctx)

// List invoices
invoices, _, _ := client.Billing.ListInvoices(ctx, nil)

// Add payment method
card, _, _ := client.Billing.AddCard(ctx, &pipeops.AddCardRequest{
    Token: "stripe-card-token",
})
```

### Webhooks

```go
// Create webhook
webhook, _, _ := client.Webhooks.Create(ctx, &pipeops.CreateWebhookRequest{
    URL:    "https://myapp.com/webhook",
    Events: []string{"project.deployed", "project.failed"},
})

// List webhooks
webhooks, _, _ := client.Webhooks.List(ctx)
```

## Service Method Categories

Services typically include methods for:

### CRUD Operations

- **Create** - Create new resources
- **Read** - Get or list resources
- **Update** - Modify existing resources
- **Delete** - Remove resources

Example:
```go
// Create
project, _, _ := client.Projects.Create(ctx, createReq)

// Read (Get)
project, _, _ := client.Projects.Get(ctx, projectUUID)

// Read (List)
projects, _, _ := client.Projects.List(ctx, nil)

// Update
project, _, _ := client.Projects.Update(ctx, projectUUID, updateReq)

// Delete
_, _ := client.Projects.Delete(ctx, projectUUID)
```

### Actions

Special operations on resources:

```go
// Deploy project
_, _ := client.Projects.Deploy(ctx, projectUUID)

// Restart project
_, _ := client.Projects.Restart(ctx, projectUUID)

// Sync deployment
_, _ := client.Projects.Sync(ctx, projectUUID)
```

### Nested Resources

Access related resources:

```go
// Get project logs
logs, _, _ := client.Projects.GetLogs(ctx, projectUUID, nil)

// Get project environment variables
envVars, _, _ := client.Projects.GetEnvVars(ctx, projectUUID)

// List project deployments
deployments, _, _ := client.Projects.ListDeployments(ctx, projectUUID, nil)
```

## Pagination

Many list methods support pagination:

```go
// Page 1, 20 items
projects, _, _ := client.Projects.List(ctx, &pipeops.ProjectListOptions{
    Page:  1,
    Limit: 20,
})

// Page 2
projects, _, _ := client.Projects.List(ctx, &pipeops.ProjectListOptions{
    Page:  2,
    Limit: 20,
})
```

## Filtering

Filter resources with options:

```go
// Filter projects by workspace
projects, _, _ := client.Projects.List(ctx, &pipeops.ProjectListOptions{
    WorkspaceID: "workspace-uuid",
})

// Filter by server
projects, _, _ := client.Projects.List(ctx, &pipeops.ProjectListOptions{
    ServerID: "server-uuid",
})
```

## Best Practices

### 1. Use Context Timeouts

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

projects, _, err := client.Projects.List(ctx, nil)
```

### 2. Check Errors

```go
projects, resp, err := client.Projects.List(ctx, nil)
if err != nil {
    log.Printf("API error: %v", err)
    if resp != nil {
        log.Printf("Status: %d", resp.StatusCode)
    }
    return err
}
```

### 3. Handle Rate Limits

```go
projects, _, err := client.Projects.List(ctx, nil)
if rateLimitErr, ok := err.(*pipeops.RateLimitError); ok {
    time.Sleep(rateLimitErr.RetryAfter)
    // Retry request
}
```

### 4. Reuse Client

```go
// ✅ Good - Create once, reuse
client, _ := pipeops.NewClient("")
client.SetToken(token)

// Use for multiple requests
projects, _, _ := client.Projects.List(ctx, nil)
servers, _, _ := client.Servers.List(ctx)

// ❌ Bad - Creating new client for each request
client1, _ := pipeops.NewClient("")
projects, _, _ := client1.Projects.List(ctx, nil)

client2, _ := pipeops.NewClient("")
servers, _, _ := client2.Servers.List(ctx)
```

## Next Steps

Explore detailed documentation for each service:

- [Projects Service](projects.md) - Comprehensive project management
- [Servers Service](servers.md) - Server and cluster management
- [Teams Service](teams.md) - Team collaboration
- [Billing Service](billing.md) - Payment and subscription management

Or learn about advanced features:

- [Error Handling](../advanced/error-handling.md)
- [Rate Limiting](../advanced/rate-limiting.md)
- [Logging](../advanced/logging.md)
