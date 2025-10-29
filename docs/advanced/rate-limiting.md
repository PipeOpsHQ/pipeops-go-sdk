# Rate Limiting

Handle API rate limits effectively with the SDK.

## Rate Limit Detection

The SDK automatically detects rate limit errors (HTTP 429):

```go
projects, _, err := client.Projects.List(ctx, nil)
if rateLimitErr, ok := err.(*pipeops.RateLimitError); ok {
    fmt.Printf("Rate limited. Retry after: %v\n", rateLimitErr.RetryAfter)
}
```

## Automatic Retry

The SDK automatically retries rate-limited requests:

```go
// Automatically handles rate limits with exponential backoff
client, _ := pipeops.NewClient("",
    pipeops.WithMaxRetries(5),
)

projects, _, err := client.Projects.List(ctx, nil)
```

## Manual Rate Limit Handling

Handle rate limits explicitly:

```go
projects, resp, err := client.Projects.List(ctx, nil)
if err != nil {
    if resp != nil && resp.StatusCode == 429 {
        // Get retry-after header
        retryAfter := resp.Header.Get("Retry-After")
        if retryAfter != "" {
            seconds, _ := strconv.Atoi(retryAfter)
            time.Sleep(time.Duration(seconds) * time.Second)
            
            // Retry
            projects, _, err = client.Projects.List(ctx, nil)
        }
    }
}
```

## Rate Limit Headers

Check rate limit headers:

```go
projects, resp, err := client.Projects.List(ctx, nil)
if resp != nil {
    limit := resp.Header.Get("X-RateLimit-Limit")
    remaining := resp.Header.Get("X-RateLimit-Remaining")
    reset := resp.Header.Get("X-RateLimit-Reset")
    
    fmt.Printf("Rate Limit: %s/%s (resets at %s)\n", remaining, limit, reset)
}
```

## Best Practices

### 1. Implement Backoff

```go
func makeRequestWithBackoff(ctx context.Context) error {
    backoff := 1 * time.Second
    maxBackoff := 60 * time.Second
    
    for {
        _, resp, err := client.Projects.List(ctx, nil)
        if err == nil {
            return nil
        }
        
        if resp != nil && resp.StatusCode == 429 {
            time.Sleep(backoff)
            backoff = time.Duration(math.Min(
                float64(backoff*2),
                float64(maxBackoff),
            ))
            continue
        }
        
        return err
    }
}
```

### 2. Use Token Bucket Pattern

```go
type RateLimiter struct {
    tokens chan struct{}
}

func NewRateLimiter(rate int) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, rate),
    }
    
    // Refill tokens
    go func() {
        ticker := time.NewTicker(time.Second / time.Duration(rate))
        for range ticker.C {
            select {
            case rl.tokens <- struct{}{}:
            default:
            }
        }
    }()
    
    return rl
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-rl.tokens:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### 3. Monitor Usage

```go
type RateLimitMonitor struct {
    requests  int
    limit     int
    remaining int
    reset     time.Time
}

func (m *RateLimitMonitor) Update(resp *http.Response) {
    m.requests++
    if resp != nil {
        // Parse headers
        limit, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Limit"))
        remaining, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
        
        m.limit = limit
        m.remaining = remaining
    }
}

func (m *RateLimitMonitor) ShouldWait() bool {
    return m.remaining < 10 // Conservative threshold
}
```

## See Also

- [Error Handling](error-handling.md)
- [Retries & Timeouts](retries-timeouts.md)
