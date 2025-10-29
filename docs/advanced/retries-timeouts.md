# Retries & Timeouts

Configure automatic retries and timeouts for robust API interactions.

## Timeout Configuration

### Client-Level Timeout

Set timeout when creating the client:

```go
client, err := pipeops.NewClient("",
    pipeops.WithTimeout(30*time.Second),
)
```

### Request-Level Timeout

Set timeout per request using context:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

projects, _, err := client.Projects.List(ctx, nil)
if errors.Is(err, context.DeadlineExceeded) {
    log.Println("Request timed out")
}
```

## Retry Configuration

### Client-Level Retries

Configure automatic retries:

```go
client, err := pipeops.NewClient("",
    pipeops.WithMaxRetries(5),
)
```

The SDK automatically retries:
- Network errors
- HTTP 5xx errors
- HTTP 429 rate limit errors

### Exponential Backoff

Retries use exponential backoff with jitter:

```
Attempt 1: Wait 100ms-500ms
Attempt 2: Wait 200ms-1s
Attempt 3: Wait 400ms-2s
...
```

### Custom Retry Logic

Implement custom retry logic:

```go
func retryWithBackoff(fn func() error, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        // Calculate backoff
        waitTime := time.Duration(math.Pow(2, float64(i))) * 100 * time.Millisecond
        jitter := time.Duration(rand.Int63n(int64(waitTime)))
        
        time.Sleep(waitTime + jitter)
    }
    
    return fmt.Errorf("max retries exceeded")
}

// Usage
err := retryWithBackoff(func() error {
    _, _, err := client.Projects.Deploy(ctx, projectUUID)
    return err
}, 5)
```

## Timeout Strategies

### Short Timeout for Health Checks

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

health, _, _ := client.Misc.HealthCheck(ctx)
```

### Long Timeout for Deployments

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

deployment, _, err := client.Projects.Deploy(ctx, projectUUID)
```

### Progressive Timeout

Increase timeout on retries:

```go
func requestWithProgressiveTimeout(attempt int) error {
    timeout := time.Duration(5+attempt*5) * time.Second
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    _, _, err := client.Projects.List(ctx, nil)
    return err
}
```

## Best Practices

### 1. Set Appropriate Timeouts

```go
// Quick reads
ctx, _ := context.WithTimeout(ctx, 10*time.Second)

// Long operations
ctx, _ := context.WithTimeout(ctx, 5*time.Minute)
```

### 2. Don't Retry Indefinitely

```go
// ✅ Good - Limited retries
client, _ := pipeops.NewClient("",
    pipeops.WithMaxRetries(5),
)

// ❌ Bad - Infinite retries
for {
    _, _, err := client.Projects.List(ctx, nil)
    if err == nil {
        break
    }
}
```

### 3. Use Exponential Backoff

```go
for i := 0; i < maxRetries; i++ {
    _, _, err := client.Projects.List(ctx, nil)
    if err == nil {
        break
    }
    
    // Exponential backoff
    time.Sleep(time.Duration(1<<i) * time.Second)
}
```

## See Also

- [Error Handling](error-handling.md)
- [Rate Limiting](rate-limiting.md)
