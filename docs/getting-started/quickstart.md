# Quick Start Guide

This guide will help you make your first API call with the PipeOps Go SDK in just a few minutes.

## Basic Example

Here's a complete example showing authentication and listing projects:

```go title="main.go"
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
    // Create a new client
    client, err := pipeops.NewClient("",
        pipeops.WithTimeout(30*time.Second),
        pipeops.WithMaxRetries(3),
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Login to authenticate
    loginResp, _, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
        Email:    "your-email@example.com",
        Password: "your-password",
    })
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    
    fmt.Printf("Login successful! User: %s\n", loginResp.Data.User.Email)
    
    // Set authentication token for subsequent requests
    client.SetToken(loginResp.Data.Token)
    
    // List all projects
    projects, _, err := client.Projects.List(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to list projects: %v", err)
    }
    
    fmt.Printf("\nFound %d projects:\n", len(projects.Data.Projects))
    for _, project := range projects.Data.Projects {
        fmt.Printf("  - %s (%s)\n", project.Name, project.UUID)
    }
}
```

## Step-by-Step Explanation

### 1. Create a Client

```go
client, err := pipeops.NewClient("",
    pipeops.WithTimeout(30*time.Second),
    pipeops.WithMaxRetries(3),
)
```

The client is your main interface to the PipeOps API. Pass an empty string to use the default API URL, or provide a custom URL for testing/staging environments.

**Configuration Options:**
- `WithTimeout()` - Set request timeout (default: 30s)
- `WithMaxRetries()` - Set maximum retry attempts (default: 3)
- `WithUserAgent()` - Set custom user agent
- `WithLogger()` - Add logging support
- `WithHTTPClient()` - Use a custom HTTP client

### 2. Authenticate

```go
loginResp, _, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
    Email:    "your-email@example.com",
    Password: "your-password",
})
```

All API calls (except authentication) require a valid authentication token. The `Login` method returns a token in the response.

### 3. Set the Token

```go
client.SetToken(loginResp.Data.Token)
```

After authentication, set the token on the client. It will be automatically included in all subsequent requests.

### 4. Make API Calls

```go
projects, _, err := client.Projects.List(ctx, nil)
```

Now you can call any API method. Most methods return three values:
1. Response data (specific to the endpoint)
2. Raw HTTP response (`*http.Response`)
3. Error (if any)

## Using Environment Variables

For better security, store credentials in environment variables:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
    client, _ := pipeops.NewClient("")
    
    // Get credentials from environment
    email := os.Getenv("PIPEOPS_EMAIL")
    password := os.Getenv("PIPEOPS_PASSWORD")
    
    if email == "" || password == "" {
        log.Fatal("Please set PIPEOPS_EMAIL and PIPEOPS_PASSWORD")
    }
    
    ctx := context.Background()
    loginResp, _, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
        Email:    email,
        Password: password,
    })
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    
    client.SetToken(loginResp.Data.Token)
    fmt.Println("Authentication successful!")
}
```

Run with:

```bash
export PIPEOPS_EMAIL="your-email@example.com"
export PIPEOPS_PASSWORD="your-password"
go run main.go
```

## Using an Existing Token

If you already have a token, skip the login step:

```go
client, _ := pipeops.NewClient("")

// Set your existing token
client.SetToken("your-existing-token")

// Make authenticated requests
ctx := context.Background()
projects, _, err := client.Projects.List(ctx, nil)
```

## Error Handling

Always check for errors in production code:

```go
projects, resp, err := client.Projects.List(ctx, nil)
if err != nil {
    // Check for specific error types
    if rateLimitErr, ok := err.(*pipeops.RateLimitError); ok {
        fmt.Printf("Rate limited. Retry after: %v\n", rateLimitErr.RetryAfter)
    } else {
        log.Printf("API error: %v\n", err)
    }
    return
}

// Check HTTP status
if resp.StatusCode != 200 {
    log.Printf("Unexpected status code: %d\n", resp.StatusCode)
}

// Use the data
fmt.Printf("Projects: %+v\n", projects)
```

## Context and Timeouts

Use contexts to control request cancellation and timeouts:

```go
// Timeout after 5 seconds
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

projects, _, err := client.Projects.List(ctx, nil)
if err == context.DeadlineExceeded {
    log.Println("Request timed out")
}
```

## Complete Working Example

See the [examples/basic](https://github.com/PipeOpsHQ/pipeops-go-sdk/tree/main/examples/basic) directory for a complete working example you can run locally.

## Next Steps

- [Configuration Guide](configuration.md) - Learn about all configuration options
- [Authentication](../authentication/overview.md) - Explore authentication methods
- [API Services](../api-services/overview.md) - Discover all available services
- [Error Handling](../advanced/error-handling.md) - Handle errors effectively
