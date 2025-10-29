# Authentication Service

The Authentication Service handles user authentication, account management, and password operations.

## Overview

```go
// Access the authentication service
authService := client.Auth
```

**Note:** This page documents the Auth service API methods. For authentication guides and best practices, see:
- [Authentication Overview](../authentication/overview.md)
- [Basic Authentication Guide](../authentication/basic-auth.md)

## Methods

### Login

Authenticate with email and password:

```go
resp, _, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
    Email:    "user@example.com",
    Password: "your-password",
})
if err != nil {
    log.Fatalf("Login failed: %v", err)
}

// Set token for authenticated requests
client.SetToken(resp.Data.Token)
fmt.Printf("Logged in as: %s\n", resp.Data.User.Email)
```

### Signup

Create a new user account:

```go
resp, _, err := client.Auth.Signup(ctx, &pipeops.SignupRequest{
    Email:     "newuser@example.com",
    Password:  "secure-password",
    FirstName: "John",
    LastName:  "Doe",
})
if err != nil {
    log.Fatalf("Signup failed: %v", err)
}

fmt.Printf("Account created: %s\n", resp.Data.User.Email)
```

### Change Password

Change the authenticated user's password:

```go
// Must be authenticated
client.SetToken(userToken)

resp, _, err := client.Auth.ChangePassword(ctx, &pipeops.ChangePasswordRequest{
    OldPassword: "current-password",
    NewPassword: "new-password",
})
if err != nil {
    log.Fatalf("Password change failed: %v", err)
}

fmt.Println("Password changed successfully")
```

### Request Password Reset

Send a password reset email:

```go
resp, _, err := client.Auth.RequestPasswordReset(ctx, &pipeops.PasswordResetRequest{
    Email: "user@example.com",
})
if err != nil {
    log.Fatalf("Failed to request reset: %v", err)
}

fmt.Println("Password reset email sent")
```

### Reset Password

Complete password reset with token from email:

```go
resp, err := client.Auth.ResetPassword(ctx, &pipeops.ResetPasswordRequest{
    Token:       "reset-token-from-email",
    NewPassword: "new-password",
})
if err != nil {
    log.Fatalf("Password reset failed: %v", err)
}

fmt.Println("Password reset successful")
```

### Verify Password Reset Token

Check if a reset token is valid:

```go
resp, err := client.Auth.VerifyPasswordResetToken(ctx, "token-from-email")
if err != nil {
    log.Println("Token is invalid or expired")
} else {
    fmt.Println("Token is valid")
}
```

### Activate Email

Activate user email with token:

```go
resp, err := client.Auth.ActivateEmail(ctx, &pipeops.ActivateEmailRequest{
    Token: "activation-token-from-email",
})
if err != nil {
    log.Fatalf("Email activation failed: %v", err)
}

fmt.Println("Email activated successfully")
```

### Verify Login (2FA)

Verify login with two-factor authentication code:

```go
resp, _, err := client.Auth.VerifyLogin(ctx, &pipeops.VerifyLoginRequest{
    Email: "user@example.com",
    Code:  "123456", // 2FA code
})
if err != nil {
    log.Fatalf("2FA verification failed: %v", err)
}

client.SetToken(resp.Data.Token)
fmt.Println("2FA verification successful")
```

### OAuth Signup

Initiate OAuth signup with a provider:

```go
resp, err := client.Auth.OAuthSignup(ctx, "google")
if err != nil {
    log.Fatalf("OAuth signup failed: %v", err)
}

// Redirects to OAuth provider
```

### OAuth Callback

Handle OAuth callback after authorization:

```go
resp, _, err := client.Auth.OAuthCallback(ctx, "google")
if err != nil {
    log.Fatalf("OAuth callback failed: %v", err)
}

client.SetToken(resp.Data.Token)
fmt.Printf("Authenticated via OAuth: %s\n", resp.Data.User.Email)
```

## Data Types

### LoginRequest

```go
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

### SignupRequest

```go
type SignupRequest struct {
    Email     string `json:"email"`
    Password  string `json:"password"`
    FirstName string `json:"first_name,omitempty"`
    LastName  string `json:"last_name,omitempty"`
}
```

### User

```go
type User struct {
    ID            string     `json:"id,omitempty"`
    UUID          string     `json:"uuid,omitempty"`
    Email         string     `json:"email,omitempty"`
    FirstName     string     `json:"first_name,omitempty"`
    LastName      string     `json:"last_name,omitempty"`
    IsActive      bool       `json:"is_active,omitempty"`
    EmailVerified bool       `json:"email_verified,omitempty"`
    CreatedAt     *Timestamp `json:"created_at,omitempty"`
    UpdatedAt     *Timestamp `json:"updated_at,omitempty"`
}
```

## Complete Example

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
    ctx := context.Background()
    
    // Get credentials from environment
    email := os.Getenv("PIPEOPS_EMAIL")
    password := os.Getenv("PIPEOPS_PASSWORD")
    
    // Login
    loginResp, _, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
        Email:    email,
        Password: password,
    })
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    
    fmt.Printf("Successfully logged in as: %s\n", loginResp.Data.User.Email)
    fmt.Printf("User ID: %s\n", loginResp.Data.User.UUID)
    fmt.Printf("Email Verified: %t\n", loginResp.Data.User.EmailVerified)
    
    // Set token for future requests
    client.SetToken(loginResp.Data.Token)
    
    // Now you can make authenticated API calls
    projects, _, _ := client.Projects.List(ctx, nil)
    fmt.Printf("Found %d projects\n", len(projects.Data.Projects))
}
```

## See Also

- [Authentication Overview](../authentication/overview.md)
- [Basic Authentication](../authentication/basic-auth.md)
- [OAuth 2.0](../authentication/oauth.md)
- [Users Service](users.md) - Manage user profile and settings
