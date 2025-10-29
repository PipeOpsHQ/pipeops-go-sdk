# Configuration

The PipeOps Go SDK uses the functional options pattern for flexible configuration. This guide covers all available configuration options.

## Basic Configuration

```go
import (
    "time"
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

client, err := pipeops.NewClient("",
    pipeops.WithTimeout(30*time.Second),
    pipeops.WithMaxRetries(3),
)
```

## Configuration Options

### Base URL

Set a custom API endpoint (useful for testing or staging environments):

```go
// Use default production URL
client, _ := pipeops.NewClient("")

// Use custom URL
client, _ := pipeops.NewClient("https://staging-api.pipeops.io")
```

**Default:** `https://api.pipeops.io`

### Timeout

Set the maximum duration for requests:

```go
client, _ := pipeops.NewClient("",
    pipeops.WithTimeout(60*time.Second), // 60 second timeout
)
```

**Default:** 30 seconds

**Recommendation:** 
- Use shorter timeouts (10-30s) for interactive applications
- Use longer timeouts (60-120s) for batch operations

### Max Retries

Configure automatic retry behavior for failed requests:

```go
client, _ := pipeops.NewClient("",
    pipeops.WithMaxRetries(5), // Retry up to 5 times
)
```

**Default:** 3 retries

The SDK automatically retries:
- Network errors (connection timeouts, DNS failures)
- HTTP 5xx server errors
- HTTP 429 rate limit errors

Retries use exponential backoff with jitter.

### User Agent

Set a custom user agent string:

```go
client, _ := pipeops.NewClient("",
    pipeops.WithUserAgent("my-app/1.0.0"),
)
```

**Default:** `pipeops-go-sdk/1.0.0`

**Best Practice:** Include your application name and version for better tracking and debugging.

### Custom HTTP Client

Use your own configured HTTP client:

```go
import (
    "net/http"
    "time"
)

customClient := &http.Client{
    Timeout: 90 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        DisableCompression:  false,
        DisableKeepAlives:   false,
    },
}

client, _ := pipeops.NewClient("",
    pipeops.WithHTTPClient(customClient),
)
```

**Use Cases:**
- Custom TLS configuration
- Proxy support
- Connection pooling tuning
- Custom transport middleware

### Logger

Add logging for debugging and monitoring:

```go
import "log"

type MyLogger struct{}

func (l *MyLogger) Debug(msg string, keysAndValues ...interface{}) {
    log.Printf("[DEBUG] %s %v", msg, keysAndValues)
}

func (l *MyLogger) Info(msg string, keysAndValues ...interface{}) {
    log.Printf("[INFO] %s %v", msg, keysAndValues)
}

func (l *MyLogger) Warn(msg string, keysAndValues ...interface{}) {
    log.Printf("[WARN] %s %v", msg, keysAndValues)
}

func (l *MyLogger) Error(msg string, keysAndValues ...interface{}) {
    log.Printf("[ERROR] %s %v", msg, keysAndValues)
}

// Use the logger
client, _ := pipeops.NewClient("",
    pipeops.WithLogger(&MyLogger{}),
)
```

## Complete Configuration Example

Here's a production-ready configuration:

```go
package main

import (
    "crypto/tls"
    "log"
    "net/http"
    "time"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
    // Custom HTTP transport with TLS configuration
    transport := &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        TLSClientConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
        DisableCompression: false,
    }
    
    // Custom HTTP client
    httpClient := &http.Client{
        Timeout:   60 * time.Second,
        Transport: transport,
    }
    
    // Create SDK client with all options
    client, err := pipeops.NewClient("https://api.pipeops.io",
        pipeops.WithHTTPClient(httpClient),
        pipeops.WithTimeout(60*time.Second),
        pipeops.WithMaxRetries(5),
        pipeops.WithUserAgent("my-production-app/2.0.0"),
        pipeops.WithLogger(&MyLogger{}),
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    
    // Use the client
    log.Println("Client configured successfully")
}
```

## Environment-Based Configuration

Configure based on environment:

```go
package main

import (
    "os"
    "time"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func createClient() (*pipeops.Client, error) {
    // Base URL from environment
    baseURL := os.Getenv("PIPEOPS_API_URL")
    if baseURL == "" {
        baseURL = "https://api.pipeops.io" // Default
    }
    
    // Timeout from environment
    timeout := 30 * time.Second
    if timeoutStr := os.Getenv("PIPEOPS_TIMEOUT"); timeoutStr != "" {
        if d, err := time.ParseDuration(timeoutStr); err == nil {
            timeout = d
        }
    }
    
    return pipeops.NewClient(baseURL,
        pipeops.WithTimeout(timeout),
        pipeops.WithMaxRetries(3),
    )
}
```

## Proxy Configuration

Configure an HTTP proxy:

```go
import (
    "net/http"
    "net/url"
)

// Parse proxy URL
proxyURL, _ := url.Parse("http://proxy.example.com:8080")

// Create transport with proxy
transport := &http.Transport{
    Proxy: http.ProxyURL(proxyURL),
}

httpClient := &http.Client{
    Transport: transport,
}

// Use with SDK
client, _ := pipeops.NewClient("",
    pipeops.WithHTTPClient(httpClient),
)
```

## TLS Configuration

Custom TLS settings:

```go
import (
    "crypto/tls"
    "net/http"
)

// Custom TLS config
tlsConfig := &tls.Config{
    MinVersion:         tls.VersionTLS12,
    InsecureSkipVerify: false, // Always verify in production!
}

transport := &http.Transport{
    TLSClientConfig: tlsConfig,
}

httpClient := &http.Client{
    Transport: transport,
}

client, _ := pipeops.NewClient("",
    pipeops.WithHTTPClient(httpClient),
)
```

## Token Management

Setting and updating authentication tokens:

```go
// Set token after login
client.SetToken("your-auth-token")

// Update token when it changes
client.SetToken("new-refreshed-token")

// Token is automatically included in all requests
```

## Best Practices

### Production Configuration

```go
client, _ := pipeops.NewClient("",
    pipeops.WithTimeout(45*time.Second),      // Reasonable timeout
    pipeops.WithMaxRetries(5),                // Handle transient failures
    pipeops.WithUserAgent("app/1.0.0"),       // Track your application
    pipeops.WithLogger(productionLogger),     // Monitor API calls
)
```

### Development Configuration

```go
client, _ := pipeops.NewClient("https://staging-api.pipeops.io",
    pipeops.WithTimeout(60*time.Second),      // Longer timeout for debugging
    pipeops.WithMaxRetries(1),                // Fail fast for debugging
    pipeops.WithLogger(debugLogger),          // Verbose logging
)
```

### Testing Configuration

```go
client, _ := pipeops.NewClient("http://localhost:8080",
    pipeops.WithTimeout(5*time.Second),       // Short timeout
    pipeops.WithMaxRetries(0),                // No retries in tests
)
```

## Next Steps

- [Quick Start Guide](quickstart.md) - Make your first API call
- [Authentication](../authentication/overview.md) - Learn about authentication
- [Error Handling](../advanced/error-handling.md) - Handle errors and retries
