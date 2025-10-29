# PipeOps Go SDK

Welcome to the comprehensive documentation for the PipeOps Go SDK - a powerful, production-ready Go client library for interacting with the PipeOps Control Plane API.

## Features

- ✅ **Complete API Coverage** - All API endpoints covered across 18+ service modules
- 🔒 **Type-Safe** - Strongly typed request/response structures
- ⚡ **Context Support** - All methods support context for cancellation and timeouts
- 🔄 **Automatic Retries** - Built-in retry logic with exponential backoff
- 🚀 **Production-Ready** - Optimized HTTP client with connection pooling
- ⚙️ **Configurable** - Flexible configuration via functional options pattern
- 🛡️ **Rate Limit Handling** - Automatic detection and typed errors
- 🔐 **OAuth 2.0 Support** - Full OAuth 2.0 authorization code flow
- 📝 **Logging Support** - Optional logger interface for debugging
- ✨ **Well-Documented** - Comprehensive examples and documentation
- 🧪 **Tested** - Unit and integration tests included

## Quick Links

- [Installation](getting-started/installation.md) - Get started quickly
- [Quick Start Guide](getting-started/quickstart.md) - Your first API call
- [Authentication](authentication/overview.md) - Learn about authentication methods
- [API Services](api-services/overview.md) - Explore all available services
- [Examples](examples/complete-examples.md) - See real-world examples

## At a Glance

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
    // Create client
    client, _ := pipeops.NewClient("")
    
    // Login
    ctx := context.Background()
    loginResp, _, _ := client.Auth.Login(ctx, &pipeops.LoginRequest{
        Email:    "your-email@example.com",
        Password: "your-password",
    })
    
    // Set token
    client.SetToken(loginResp.Data.Token)
    
    // List projects
    projects, _, _ := client.Projects.List(ctx, nil)
    fmt.Printf("Found %d projects\n", len(projects.Data.Projects))
}
```

## SDK Services

The SDK is organized into specialized services for different API endpoints:

| Service | Description |
|---------|-------------|
| **Auth** | User authentication, signup, password management |
| **OAuth** | OAuth 2.0 authorization flows |
| **Projects** | Project creation, deployment, and management |
| **Servers** | Server/cluster provisioning and management |
| **Environments** | Environment configuration and management |
| **Teams** | Team collaboration and member management |
| **Workspaces** | Workspace organization and settings |
| **Billing** | Subscription, payment, and invoice management |
| **AddOns** | Marketplace add-on deployment |
| **Webhooks** | Webhook configuration and delivery management |
| **Users** | User profile and settings |
| **CloudProviders** | Cloud provider integration |
| **ServiceTokens** | Service account token management |

## Community and Support

- **GitHub Issues**: [Report bugs or request features](https://github.com/PipeOpsHQ/pipeops-go-sdk/issues)
- **Contributing**: See our [Contributing Guide](contributing.md)
- **Changelog**: View [Release Notes](changelog.md)

## License

This SDK is distributed under the terms specified in the [LICENSE](https://github.com/PipeOpsHQ/pipeops-go-sdk/blob/main/LICENSE) file.

---

**Ready to get started?** Head over to the [Installation Guide](getting-started/installation.md) to begin using the SDK.
