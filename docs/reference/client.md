# Client Reference

Complete reference for the Client type.

## Client Type

```go
type Client struct {
    // HTTP client used for requests
    client *http.Client
    
    // Base URL for API requests
    BaseURL *url.URL
    
    // User agent string
    UserAgent string
    
    // Services
    Auth         *AuthService
    OAuth        *OAuthService
    Projects     *ProjectService
    Servers      *ServerService
    // ... other services
}
```

## Constructor

### NewClient

```go
func NewClient(baseURL string, options ...ClientOption) (*Client, error)
```

Creates a new PipeOps API client.

**Parameters:**
- `baseURL` - API base URL (use "" for default)
- `options` - Configuration options

**Example:**
```go
client, err := pipeops.NewClient("",
    pipeops.WithTimeout(30*time.Second),
    pipeops.WithMaxRetries(3),
)
```

## Methods

### SetToken

```go
func (c *Client) SetToken(token string)
```

Set authentication token for requests.

### SetHTTPClient

```go
func (c *Client) SetHTTPClient(httpClient *http.Client)
```

Set custom HTTP client.

## Configuration Options

### WithTimeout

```go
func WithTimeout(timeout time.Duration) ClientOption
```

### WithMaxRetries

```go
func WithMaxRetries(maxRetries int) ClientOption
```

### WithHTTPClient

```go
func WithHTTPClient(httpClient *http.Client) ClientOption
```

### WithUserAgent

```go
func WithUserAgent(userAgent string) ClientOption
```

### WithLogger

```go
func WithLogger(logger Logger) ClientOption
```

## See Also

- [Configuration Guide](../getting-started/configuration.md)
