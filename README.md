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
- **Automatic Retries**: Built-in retry logic with exponential backoff for transient failures
- **Production-Ready HTTP Client**: Optimized connection pooling and timeouts
- **Configurable**: Flexible configuration options via functional options pattern
- **Rate Limit Handling**: Automatic detection and typed errors for rate limits
- **OAuth 2.0**: Full OAuth 2.0 authorization code flow support
- **Logging Support**: Optional logger interface for debugging
- **Well-Documented**: Comprehensive examples and documentation
- **Tested**: Unit and integration tests included

## Installation

```bash
go get github.com/PipeOpsHQ/pipeops-go-sdk
```

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
    // Create a new client with custom configuration
    client, err := pipeops.NewClient("https://api.pipeops.io",
        pipeops.WithTimeout(30*time.Second),  // Custom timeout
        pipeops.WithMaxRetries(3),             // Retry failed requests up to 3 times
    )
    if err != nil {
        log.Fatal(err)
    }
    
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
    projects, _, err := client.Projects.List(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d projects\n", len(projects.Data.Projects))
}
```

### Configuration Options

The SDK supports various configuration options through the functional options pattern:

```go
client, err := pipeops.NewClient("https://api.pipeops.io",
    // Set custom timeout (default: 30s)
    pipeops.WithTimeout(60*time.Second),
    
    // Set max retry attempts (default: 3)
    pipeops.WithMaxRetries(5),
    
    // Use custom HTTP client
    pipeops.WithHTTPClient(customHTTPClient),
    
    // Set custom user agent
    pipeops.WithUserAgent("my-app/1.0"),
    
    // Add custom logger for debugging
    pipeops.WithLogger(myLogger),
)
```

### Error Handling

The SDK provides typed errors for better error handling:

```go
projects, _, err := client.Projects.List(ctx, nil)
if err != nil {
    // Check for rate limit errors
    if rateLimitErr, ok := err.(*pipeops.RateLimitError); ok {
        fmt.Printf("Rate limited. Retry after: %v\n", rateLimitErr.RetryAfter)
        time.Sleep(rateLimitErr.RetryAfter)
        // Retry request
    }
    log.Fatal(err)
}
```

### Context and Timeouts

All API methods support context for cancellation and timeouts:

```go
// Set a timeout for the request
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

projects, _, err := client.Projects.List(ctx, nil)
if err != nil {
    if err == context.DeadlineExceeded {
        log.Println("Request timed out")
    }
    log.Fatal(err)
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

Comprehensive documentation is available in multiple formats:

### Online Documentation
- **[Full Documentation Site](https://pipeops-go-sdk.readthedocs.io/)** - Complete guides, API reference, and examples
- **[Getting Started Guide](docs/getting-started/quickstart.md)** - Quick start tutorial
- **[API Reference](https://godoc.org/github.com/PipeOpsHQ/pipeops-go-sdk/pipeops)** - Go package documentation

### Local Documentation
You can build and view the documentation locally using MkDocs:

```bash
# Quick start with helper script
./docs.sh install  # Install dependencies
./docs.sh serve    # Serve at http://127.0.0.1:8000
./docs.sh build    # Build static HTML

# Or use MkDocs directly
pip install -r docs/requirements.txt
mkdocs serve
mkdocs build
```

### Documentation Sections
- **[Installation](docs/getting-started/installation.md)** - Installation instructions
- **[Quick Start](docs/getting-started/quickstart.md)** - Get started in minutes
- **[Configuration](docs/getting-started/configuration.md)** - Client configuration options
- **[Authentication](docs/authentication/overview.md)** - Authentication methods and best practices
- **[API Services](docs/api-services/overview.md)** - Complete API service documentation
- **[Advanced Usage](docs/advanced/error-handling.md)** - Error handling, retries, logging, etc.
- **[Examples](docs/examples/complete-examples.md)** - Real-world usage examples

### Additional Resources
- [Working Code Examples](examples/) - Runnable example applications
- [PipeOps API Documentation](https://api.pipeops.io/docs) - Official REST API docs
- [Contributing Guide](CONTRIBUTING.md) - Contribution guidelines

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
