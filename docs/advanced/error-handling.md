# Error Handling

The SDK provides comprehensive error handling with typed errors for common scenarios.

## Basic Error Handling

Always check for errors returned by SDK methods:

```go
projects, resp, err := client.Projects.List(ctx, nil)
if err != nil {
    log.Printf("API error: %v", err)
    return err
}

// Use the data
fmt.Printf("Projects: %d\n", len(projects.Data.Projects))
```

## HTTP Status Codes

Check HTTP response status:

```go
projects, resp, err := client.Projects.List(ctx, nil)
if err != nil {
    if resp != nil {
        switch resp.StatusCode {
        case 400:
            log.Println("Bad request - check parameters")
        case 401:
            log.Println("Unauthorized - check authentication token")
        case 403:
            log.Println("Forbidden - insufficient permissions")
        case 404:
            log.Println("Not found")
        case 429:
            log.Println("Rate limited")
        case 500:
            log.Println("Server error")
        default:
            log.Printf("HTTP error: %d", resp.StatusCode)
        }
    }
    return err
}
```

## Rate Limit Errors

Handle rate limiting with typed errors:

```go
projects, _, err := client.Projects.List(ctx, nil)
if err != nil {
    if rateLimitErr, ok := err.(*pipeops.RateLimitError); ok {
        fmt.Printf("Rate limited. Retry after: %v\n", rateLimitErr.RetryAfter)
        
        // Wait and retry
        time.Sleep(rateLimitErr.RetryAfter)
        projects, _, err = client.Projects.List(ctx, nil)
    }
}
```

## Context Errors

Handle context cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

projects, _, err := client.Projects.List(ctx, nil)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("Request timed out")
    } else if errors.Is(err, context.Canceled) {
        log.Println("Request was canceled")
    }
    return err
}
```

## Retry Logic

Implement custom retry logic:

```go
func listProjectsWithRetry(client *pipeops.Client, ctx context.Context, maxRetries int) (*pipeops.ProjectsResponse, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        projects, resp, err := client.Projects.List(ctx, nil)
        if err == nil {
            return projects, nil
        }
        
        lastErr = err
        
        // Don't retry on client errors (4xx)
        if resp != nil && resp.StatusCode >= 400 && resp.StatusCode < 500 {
            return nil, err
        }
        
        // Wait before retry
        backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
        time.Sleep(backoff)
    }
    
    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

## Error Response Structure

API error responses follow this structure:

```go
type ErrorResponse struct {
    Response *http.Response
    Status   string `json:"status"`
    Message  string `json:"message"`
    Errors   map[string][]string `json:"errors,omitempty"`
}
```

## Validation Errors

Handle validation errors:

```go
_, _, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
    Email:    "invalid-email",
    Password: "pass",
})
if err != nil {
    if errResp, ok := err.(*pipeops.ErrorResponse); ok {
        if errResp.Errors != nil {
            for field, messages := range errResp.Errors {
                for _, msg := range messages {
                    fmt.Printf("%s: %s\n", field, msg)
                }
            }
        }
    }
}
```

## Best Practices

### 1. Always Check Errors

```go
// ✅ Good
projects, _, err := client.Projects.List(ctx, nil)
if err != nil {
    return err
}

// ❌ Bad
projects, _, _ := client.Projects.List(ctx, nil)
```

### 2. Provide Context

```go
projects, _, err := client.Projects.List(ctx, nil)
if err != nil {
    return fmt.Errorf("failed to list projects: %w", err)
}
```

### 3. Handle Specific Cases

```go
projects, resp, err := client.Projects.List(ctx, nil)
if err != nil {
    switch {
    case resp != nil && resp.StatusCode == 401:
        return refreshTokenAndRetry()
    case resp != nil && resp.StatusCode == 429:
        return retryAfterBackoff()
    default:
        return err
    }
}
```

## See Also

- [Retries & Timeouts](retries-timeouts.md)
- [Rate Limiting](rate-limiting.md)
- [Logging](logging.md)
