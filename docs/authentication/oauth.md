# OAuth 2.0 Authentication

The PipeOps Go SDK provides full support for OAuth 2.0 authorization code flow, enabling third-party applications to access the PipeOps API on behalf of users.

## OAuth 2.0 Flow Overview

1. **Authorization** - Direct user to authorization URL
2. **User Consent** - User grants permission
3. **Callback** - Receive authorization code
4. **Token Exchange** - Exchange code for access token
5. **API Access** - Use access token to make API calls

## Complete OAuth Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

var (
    clientID     = "your-oauth-client-id"
    clientSecret = "your-oauth-client-secret"
    redirectURI  = "http://localhost:3000/callback"
)

func main() {
    client, _ := pipeops.NewClient("")
    
    // Step 1: Generate authorization URL
    authURL, err := client.OAuth.Authorize(&pipeops.AuthorizeOptions{
        ClientID:     clientID,
        RedirectURI:  redirectURI,
        ResponseType: "code",
        Scope:        "user:read user:write projects:read",
        State:        "random-state-string", // CSRF protection
    })
    if err != nil {
        log.Fatalf("Failed to generate auth URL: %v", err)
    }
    
    fmt.Printf("Visit this URL to authorize:\n%s\n", authURL)
    
    // Step 2: Set up callback handler
    http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
        // Step 3: Extract authorization code
        code := r.URL.Query().Get("code")
        state := r.URL.Query().Get("state")
        
        // Verify state to prevent CSRF
        if state != "random-state-string" {
            http.Error(w, "Invalid state", http.StatusBadRequest)
            return
        }
        
        if code == "" {
            http.Error(w, "No code received", http.StatusBadRequest)
            return
        }
        
        // Step 4: Exchange code for token
        ctx := context.Background()
        token, _, err := client.OAuth.ExchangeCodeForToken(ctx, &pipeops.TokenRequest{
            GrantType:    "authorization_code",
            Code:         code,
            ClientID:     clientID,
            ClientSecret: clientSecret,
            RedirectURI:  redirectURI,
        })
        if err != nil {
            log.Printf("Token exchange failed: %v", err)
            http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
            return
        }
        
        // Step 5: Use the access token
        client.SetToken(token.AccessToken)
        
        // Get user info
        userInfo, _, err := client.OAuth.GetUserInfo(ctx)
        if err != nil {
            log.Printf("Failed to get user info: %v", err)
            http.Error(w, "Failed to get user info", http.StatusInternalServerError)
            return
        }
        
        fmt.Fprintf(w, "Authorized! User: %s (%s)", userInfo.Data.Name, userInfo.Data.Email)
    })
    
    // Start server
    log.Println("Starting server on :3000")
    log.Fatal(http.ListenAndServe(":3000", nil))
}
```

## Step 1: Generate Authorization URL

Create a URL to redirect the user for authorization:

```go
authURL, err := client.OAuth.Authorize(&pipeops.AuthorizeOptions{
    ClientID:     "your-client-id",
    RedirectURI:  "http://localhost:3000/callback",
    ResponseType: "code",
    Scope:        "user:read user:write",
    State:        "random-secure-string",
})
if err != nil {
    log.Fatalf("Failed to generate URL: %v", err)
}

// Redirect user to authURL
fmt.Printf("Authorization URL: %s\n", authURL)
```

### Parameters

| Parameter | Description | Required |
|-----------|-------------|----------|
| `ClientID` | Your OAuth application ID | Yes |
| `RedirectURI` | Callback URL after authorization | Yes |
| `ResponseType` | Set to "code" for authorization code flow | Yes |
| `Scope` | Space-separated list of permissions | No |
| `State` | Random string for CSRF protection | Recommended |

### Available Scopes

- `user:read` - Read user profile information
- `user:write` - Modify user profile
- `projects:read` - Read project information
- `projects:write` - Create and modify projects
- `servers:read` - Read server information
- `servers:write` - Create and modify servers

## Step 2: Handle Callback

After user authorization, handle the callback:

```go
http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
    // Extract parameters
    code := r.URL.Query().Get("code")
    state := r.URL.Query().Get("state")
    errorParam := r.URL.Query().Get("error")
    
    // Check for errors
    if errorParam != "" {
        errorDesc := r.URL.Query().Get("error_description")
        http.Error(w, fmt.Sprintf("Authorization error: %s - %s", errorParam, errorDesc), http.StatusBadRequest)
        return
    }
    
    // Verify state (CSRF protection)
    if state != expectedState {
        http.Error(w, "State mismatch", http.StatusBadRequest)
        return
    }
    
    // Proceed with token exchange...
})
```

## Step 3: Exchange Code for Token

Exchange the authorization code for an access token:

```go
ctx := context.Background()

token, _, err := client.OAuth.ExchangeCodeForToken(ctx, &pipeops.TokenRequest{
    GrantType:    "authorization_code",
    Code:         code,
    ClientID:     clientID,
    ClientSecret: clientSecret,
    RedirectURI:  redirectURI,
})
if err != nil {
    log.Fatalf("Token exchange failed: %v", err)
}

fmt.Printf("Access Token: %s\n", token.AccessToken)
fmt.Printf("Expires In: %d seconds\n", token.ExpiresIn)
fmt.Printf("Refresh Token: %s\n", token.RefreshToken)
```

### Token Response

```go
type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`      // "Bearer"
    ExpiresIn    int    `json:"expires_in"`      // Seconds until expiration
    RefreshToken string `json:"refresh_token"`   // Use to get new access token
    Scope        string `json:"scope"`           // Granted scopes
}
```

## Step 4: Use Access Token

Set the token and make API calls:

```go
// Set the access token
client.SetToken(token.AccessToken)

// Make authenticated API calls
ctx := context.Background()
projects, _, err := client.Projects.List(ctx, nil)
if err != nil {
    log.Fatalf("Failed to list projects: %v", err)
}

fmt.Printf("Found %d projects\n", len(projects.Data.Projects))
```

## Get User Information

Retrieve information about the authenticated user:

```go
userInfo, _, err := client.OAuth.GetUserInfo(ctx)
if err != nil {
    log.Fatalf("Failed to get user info: %v", err)
}

fmt.Printf("User ID: %s\n", userInfo.Data.Sub)
fmt.Printf("Email: %s\n", userInfo.Data.Email)
fmt.Printf("Name: %s\n", userInfo.Data.Name)
fmt.Printf("Email Verified: %t\n", userInfo.Data.EmailVerified)
```

### UserInfo Response

```go
type UserInfo struct {
    Sub           string `json:"sub"`              // User ID
    Email         string `json:"email"`
    EmailVerified bool   `json:"email_verified"`
    Name          string `json:"name"`
    GivenName     string `json:"given_name"`
    FamilyName    string `json:"family_name"`
    Picture       string `json:"picture"`
}
```

## Refresh Tokens

Use a refresh token to get a new access token:

```go
newToken, _, err := client.OAuth.ExchangeCodeForToken(ctx, &pipeops.TokenRequest{
    GrantType:    "refresh_token",
    RefreshToken: storedRefreshToken,
    ClientID:     clientID,
    ClientSecret: clientSecret,
})
if err != nil {
    log.Fatalf("Token refresh failed: %v", err)
}

// Update the token
client.SetToken(newToken.AccessToken)

// Store new refresh token
storedRefreshToken = newToken.RefreshToken
```

## Token Management

### Store Tokens Securely

```go
import (
    "encoding/json"
    "os"
)

type TokenStorage struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresAt    int64  `json:"expires_at"`
}

func saveTokens(token *pipeops.TokenResponse) error {
    storage := TokenStorage{
        AccessToken:  token.AccessToken,
        RefreshToken: token.RefreshToken,
        ExpiresAt:    time.Now().Unix() + int64(token.ExpiresIn),
    }
    
    data, _ := json.Marshal(storage)
    return os.WriteFile("tokens.json", data, 0600)
}

func loadTokens() (*TokenStorage, error) {
    data, err := os.ReadFile("tokens.json")
    if err != nil {
        return nil, err
    }
    
    var storage TokenStorage
    err = json.Unmarshal(data, &storage)
    return &storage, err
}
```

### Automatic Token Refresh

```go
type OAuthClient struct {
    client       *pipeops.Client
    accessToken  string
    refreshToken string
    expiresAt    time.Time
    mu           sync.RWMutex
}

func (oc *OAuthClient) GetToken(ctx context.Context) (string, error) {
    oc.mu.RLock()
    // Check if token is still valid
    if time.Now().Before(oc.expiresAt.Add(-5 * time.Minute)) {
        token := oc.accessToken
        oc.mu.RUnlock()
        return token, nil
    }
    oc.mu.RUnlock()
    
    // Need to refresh
    return oc.refreshAccessToken(ctx)
}

func (oc *OAuthClient) refreshAccessToken(ctx context.Context) (string, error) {
    oc.mu.Lock()
    defer oc.mu.Unlock()
    
    // Double-check after acquiring lock
    if time.Now().Before(oc.expiresAt.Add(-5 * time.Minute)) {
        return oc.accessToken, nil
    }
    
    // Refresh the token
    newToken, _, err := oc.client.OAuth.ExchangeCodeForToken(ctx, &pipeops.TokenRequest{
        GrantType:    "refresh_token",
        RefreshToken: oc.refreshToken,
        ClientID:     clientID,
        ClientSecret: clientSecret,
    })
    if err != nil {
        return "", err
    }
    
    // Update stored tokens
    oc.accessToken = newToken.AccessToken
    oc.refreshToken = newToken.RefreshToken
    oc.expiresAt = time.Now().Add(time.Duration(newToken.ExpiresIn) * time.Second)
    oc.client.SetToken(newToken.AccessToken)
    
    return oc.accessToken, nil
}
```

## Security Best Practices

### 1. Use State Parameter

Always use a state parameter to prevent CSRF attacks:

```go
import "crypto/rand"
import "encoding/hex"

func generateState() string {
    b := make([]byte, 16)
    rand.Read(b)
    return hex.EncodeToString(b)
}

state := generateState()
// Store state in session for verification

authURL, _ := client.OAuth.Authorize(&pipeops.AuthorizeOptions{
    State: state,
    // ... other params
})
```

### 2. Use HTTPS

Always use HTTPS for redirect URIs in production:

```go
// ✅ Good
redirectURI := "https://myapp.com/callback"

// ❌ Bad (only for local development)
redirectURI := "http://localhost:3000/callback"
```

### 3. Secure Token Storage

Never store tokens in plain text:

```go
// Use secure storage mechanisms
import "github.com/zalando/go-keyring"

func storeToken(token string) error {
    return keyring.Set("myapp", "oauth-token", token)
}

func getToken() (string, error) {
    return keyring.Get("myapp", "oauth-token")
}
```

### 4. Validate Redirect URI

Ensure the redirect URI matches exactly:

```go
const registeredRedirectURI = "https://myapp.com/callback"

if redirectURI != registeredRedirectURI {
    return errors.New("invalid redirect URI")
}
```

## Error Handling

```go
token, resp, err := client.OAuth.ExchangeCodeForToken(ctx, tokenReq)
if err != nil {
    if resp != nil {
        switch resp.StatusCode {
        case 400:
            log.Println("Invalid request - check parameters")
        case 401:
            log.Println("Invalid client credentials")
        case 403:
            log.Println("Access denied")
        default:
            log.Printf("Token exchange failed: %v", err)
        }
    }
    return err
}
```

## Complete Example

See the [examples/oauth](https://github.com/PipeOpsHQ/pipeops-go-sdk/tree/main/examples/oauth) directory for a complete working OAuth implementation.

## Next Steps

- [Basic Authentication](basic-auth.md) - Learn about email/password auth
- [API Services](../api-services/overview.md) - Explore available API services
- [Error Handling](../advanced/error-handling.md) - Handle errors effectively
