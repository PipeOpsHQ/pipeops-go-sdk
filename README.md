# PipeOps Go SDK

[![CI](https://github.com/PipeOpsHQ/pipeops-go-sdk/workflows/CI/badge.svg)](https://github.com/PipeOpsHQ/pipeops-go-sdk/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/PipeOpsHQ/pipeops-go-sdk)](https://goreportcard.com/report/github.com/PipeOpsHQ/pipeops-go-sdk)
[![GoDoc](https://godoc.org/github.com/PipeOpsHQ/pipeops-go-sdk?status.svg)](https://godoc.org/github.com/PipeOpsHQ/pipeops-go-sdk/pipeops)
[![License](https://img.shields.io/github/license/PipeOpsHQ/pipeops-go-sdk)](LICENSE)

A comprehensive Go SDK for interacting with the PipeOps Control Plane API.

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

# OAuth 2.0 Support

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

## Release Process

See [RELEASE.md](RELEASE.md) for information about creating releases and publishing the SDK.
