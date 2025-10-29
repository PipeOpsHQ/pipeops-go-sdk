# Authentication Overview

The PipeOps Go SDK supports multiple authentication methods to suit different use cases. All authenticated requests include an authentication token in the request headers.

## Authentication Methods

### 1. Email/Password Authentication

The most common method - authenticate with email and password to receive a token:

```go
loginResp, _, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
    Email:    "user@example.com",
    Password: "your-password",
})

client.SetToken(loginResp.Data.Token)
```

**Use Case:** User-facing applications, interactive tools

**See:** [Basic Authentication Guide](basic-auth.md)

### 2. OAuth 2.0

Use OAuth 2.0 authorization code flow for third-party integrations:

```go
// Generate authorization URL
authURL, _ := client.OAuth.Authorize(&pipeops.AuthorizeOptions{
    ClientID:     "your-client-id",
    RedirectURI:  "http://localhost:3000/callback",
    ResponseType: "code",
    Scope:        "user:read user:write",
})

// After user authorization, exchange code for token
token, _, _ := client.OAuth.ExchangeCodeForToken(ctx, &pipeops.TokenRequest{
    GrantType:    "authorization_code",
    Code:         authCode,
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
})

client.SetToken(token.AccessToken)
```

**Use Case:** Third-party integrations, delegated access

**See:** [OAuth 2.0 Guide](oauth.md)

### 3. Existing Token

If you already have a valid token, use it directly:

```go
client.SetToken("your-existing-token")
```

**Use Case:** Service accounts, automation scripts, token refresh flows

## Token Management

### Setting a Token

```go
// After authentication
client.SetToken(token)
```

### Token Lifecycle

Tokens returned by the API have a limited lifetime. Handle token expiration:

```go
// Check if request failed due to authentication
projects, resp, err := client.Projects.List(ctx, nil)
if err != nil && resp != nil && resp.StatusCode == 401 {
    // Token expired or invalid - re-authenticate
    loginResp, _, _ := client.Auth.Login(ctx, loginReq)
    client.SetToken(loginResp.Data.Token)
    
    // Retry the request
    projects, resp, err = client.Projects.List(ctx, nil)
}
```

### Token Storage

!!! warning "Security Best Practice"
    Never hardcode tokens in your source code. Use environment variables or secure storage.

```go
import "os"

// Load token from environment
token := os.Getenv("PIPEOPS_TOKEN")
if token != "" {
    client.SetToken(token)
}
```

## Authentication Flow Examples

### Interactive Application

```go
func authenticateUser(client *pipeops.Client, email, password string) (string, error) {
    ctx := context.Background()
    
    resp, _, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
        Email:    email,
        Password: password,
    })
    if err != nil {
        return "", fmt.Errorf("login failed: %w", err)
    }
    
    // Store token for subsequent requests
    client.SetToken(resp.Data.Token)
    
    return resp.Data.Token, nil
}
```

### Service Account

```go
func setupServiceClient() (*pipeops.Client, error) {
    client, _ := pipeops.NewClient("")
    
    // Load service token from environment
    token := os.Getenv("PIPEOPS_SERVICE_TOKEN")
    if token == "" {
        return nil, errors.New("PIPEOPS_SERVICE_TOKEN not set")
    }
    
    client.SetToken(token)
    return client, nil
}
```

### Token Refresh Pattern

```go
type AuthManager struct {
    client   *pipeops.Client
    email    string
    password string
    token    string
    mu       sync.RWMutex
}

func (am *AuthManager) GetToken(ctx context.Context) (string, error) {
    am.mu.RLock()
    token := am.token
    am.mu.RUnlock()
    
    if token != "" {
        return token, nil
    }
    
    // Need to authenticate
    return am.refreshToken(ctx)
}

func (am *AuthManager) refreshToken(ctx context.Context) (string, error) {
    am.mu.Lock()
    defer am.mu.Unlock()
    
    resp, _, err := am.client.Auth.Login(ctx, &pipeops.LoginRequest{
        Email:    am.email,
        Password: am.password,
    })
    if err != nil {
        return "", err
    }
    
    am.token = resp.Data.Token
    am.client.SetToken(am.token)
    
    return am.token, nil
}
```

## Security Best Practices

### 1. Protect Credentials

```go
// ✅ Good - Use environment variables
email := os.Getenv("PIPEOPS_EMAIL")
password := os.Getenv("PIPEOPS_PASSWORD")

// ❌ Bad - Hardcoded credentials
email := "user@example.com"
password := "hardcoded-password"
```

### 2. Use Context Timeouts

```go
// Set timeout for authentication
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, _, err := client.Auth.Login(ctx, loginReq)
```

### 3. Handle Authentication Errors

```go
resp, _, err := client.Auth.Login(ctx, loginReq)
if err != nil {
    // Don't expose sensitive error details to end users
    log.Printf("Authentication failed: %v", err)
    return errors.New("authentication failed")
}
```

### 4. Validate Input

```go
func login(email, password string) error {
    // Validate before making request
    if email == "" || password == "" {
        return errors.New("email and password required")
    }
    
    if !strings.Contains(email, "@") {
        return errors.New("invalid email format")
    }
    
    // Proceed with login
    resp, _, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
        Email:    email,
        Password: password,
    })
    // ...
}
```

### 5. Secure Token Storage

For persistent storage, use secure methods:

```go
import "github.com/zalando/go-keyring"

// Store token securely
func storeToken(token string) error {
    return keyring.Set("pipeops-app", "api-token", token)
}

// Retrieve token
func getToken() (string, error) {
    return keyring.Get("pipeops-app", "api-token")
}
```

## Authentication Errors

Common authentication errors and how to handle them:

| Error | Description | Solution |
|-------|-------------|----------|
| 401 Unauthorized | Invalid or expired token | Re-authenticate and get a new token |
| 400 Bad Request | Invalid credentials | Check email/password format |
| 429 Too Many Requests | Rate limited | Wait before retrying (check Retry-After header) |
| 500 Server Error | Server issue | Retry with exponential backoff |

## Next Steps

- [Basic Authentication](basic-auth.md) - Detailed email/password authentication guide
- [OAuth 2.0](oauth.md) - Complete OAuth integration guide
- [Service Tokens](../api-services/servicetokens.md) - Learn about service account tokens
