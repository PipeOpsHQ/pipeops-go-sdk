# Errors Reference

Error types used by the SDK.

## ErrorResponse

```go
type ErrorResponse struct {
    Response *http.Response
    Status   string              `json:"status"`
    Message  string              `json:"message"`
    Errors   map[string][]string `json:"errors,omitempty"`
}
```

General API error response.

## RateLimitError

```go
type RateLimitError struct {
    Response   *http.Response
    RetryAfter time.Duration
}
```

Returned when rate limited (HTTP 429).

## Usage

```go
projects, resp, err := client.Projects.List(ctx, nil)
if err != nil {
    // Check for rate limit
    if rateLimitErr, ok := err.(*pipeops.RateLimitError); ok {
        time.Sleep(rateLimitErr.RetryAfter)
        // Retry
    }
    
    // Check for general error
    if errResp, ok := err.(*pipeops.ErrorResponse); ok {
        log.Printf("API error: %s", errResp.Message)
    }
}
```

## See Also

- [Error Handling](../advanced/error-handling.md)
