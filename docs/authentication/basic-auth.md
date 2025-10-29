# Basic Authentication

Basic authentication using email and password is the most straightforward way to authenticate with the PipeOps API.

## Login

Authenticate a user and receive an access token:

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
    client, _ := pipeops.NewClient("")
    ctx := context.Background()
    
    // Login with email and password
    resp, _, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
        Email:    "user@example.com",
        Password: "your-password",
    })
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    
    fmt.Printf("Login successful!\n")
    fmt.Printf("Token: %s\n", resp.Data.Token)
    fmt.Printf("User: %s %s\n", resp.Data.User.FirstName, resp.Data.User.LastName)
    
    // Set token for authenticated requests
    client.SetToken(resp.Data.Token)
}
```

### Request

```go
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

### Response

```go
type LoginResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Data    struct {
        Token string `json:"token"`
        User  User   `json:"user"`
    } `json:"data"`
}
```

## Signup

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

fmt.Printf("Account created for: %s\n", resp.Data.User.Email)
```

### Request

```go
type SignupRequest struct {
    Email     string `json:"email"`
    Password  string `json:"password"`
    FirstName string `json:"first_name,omitempty"`
    LastName  string `json:"last_name,omitempty"`
}
```

**Validation:**
- Email must be valid format
- Password must be at least 6 characters
- Email must not already exist

## Change Password

Change password for an authenticated user:

```go
// Must be authenticated first
client.SetToken(userToken)

resp, _, err := client.Auth.ChangePassword(ctx, &pipeops.ChangePasswordRequest{
    OldPassword: "current-password",
    NewPassword: "new-secure-password",
})
if err != nil {
    log.Fatalf("Password change failed: %v", err)
}

fmt.Println("Password changed successfully!")
```

### Request

```go
type ChangePasswordRequest struct {
    OldPassword string `json:"old_password"`
    NewPassword string `json:"new_password"`
}
```

## Password Reset Flow

### Step 1: Request Password Reset

Send a password reset email:

```go
resp, _, err := client.Auth.RequestPasswordReset(ctx, &pipeops.PasswordResetRequest{
    Email: "user@example.com",
})
if err != nil {
    log.Fatalf("Failed to request reset: %v", err)
}

fmt.Println("Password reset email sent!")
```

### Step 2: Verify Reset Token (Optional)

Verify the token from the email is valid:

```go
token := "reset-token-from-email"

resp, err := client.Auth.VerifyPasswordResetToken(ctx, token)
if err != nil {
    log.Fatalf("Invalid or expired token: %v", err)
}

fmt.Println("Token is valid!")
```

### Step 3: Reset Password

Complete the password reset:

```go
resp, err := client.Auth.ResetPassword(ctx, &pipeops.ResetPasswordRequest{
    Token:       "reset-token-from-email",
    NewPassword: "new-secure-password",
})
if err != nil {
    log.Fatalf("Password reset failed: %v", err)
}

fmt.Println("Password reset successful!")
```

## Email Verification

Activate a user's email address:

```go
resp, err := client.Auth.ActivateEmail(ctx, &pipeops.ActivateEmailRequest{
    Token: "activation-token-from-email",
})
if err != nil {
    log.Fatalf("Email activation failed: %v", err)
}

fmt.Println("Email activated successfully!")
```

## Two-Factor Authentication

If 2FA is enabled, verify the login with a code:

```go
// Initial login returns a pending status if 2FA is enabled
loginResp, _, _ := client.Auth.Login(ctx, &pipeops.LoginRequest{
    Email:    "user@example.com",
    Password: "password",
})

// Verify with 2FA code
resp, _, err := client.Auth.VerifyLogin(ctx, &pipeops.VerifyLoginRequest{
    Email: "user@example.com",
    Code:  "123456", // Code from authenticator app
})
if err != nil {
    log.Fatalf("2FA verification failed: %v", err)
}

// Set the token after successful verification
client.SetToken(resp.Data.Token)
```

## Complete Authentication Example

Here's a complete example with error handling:

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "time"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func authenticate(email, password string) (*pipeops.Client, error) {
    // Create client with timeout
    client, err := pipeops.NewClient("",
        pipeops.WithTimeout(30*time.Second),
        pipeops.WithMaxRetries(3),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create client: %w", err)
    }
    
    // Set request timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Validate input
    if email == "" || password == "" {
        return nil, errors.New("email and password are required")
    }
    
    // Attempt login
    resp, httpResp, err := client.Auth.Login(ctx, &pipeops.LoginRequest{
        Email:    email,
        Password: password,
    })
    if err != nil {
        // Check for specific errors
        if httpResp != nil {
            switch httpResp.StatusCode {
            case 401:
                return nil, errors.New("invalid email or password")
            case 429:
                return nil, errors.New("too many login attempts, please try again later")
            case 500:
                return nil, errors.New("server error, please try again")
            }
        }
        return nil, fmt.Errorf("login failed: %w", err)
    }
    
    // Set token
    client.SetToken(resp.Data.Token)
    
    log.Printf("Authenticated as: %s\n", resp.Data.User.Email)
    
    return client, nil
}

func main() {
    email := "user@example.com"
    password := "your-password"
    
    client, err := authenticate(email, password)
    if err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }
    
    // Use authenticated client
    ctx := context.Background()
    projects, _, err := client.Projects.List(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to list projects: %v", err)
    }
    
    fmt.Printf("Found %d projects\n", len(projects.Data.Projects))
}
```

## Best Practices

### 1. Use Environment Variables

```go
import "os"

email := os.Getenv("PIPEOPS_EMAIL")
password := os.Getenv("PIPEOPS_PASSWORD")

if email == "" || password == "" {
    log.Fatal("Set PIPEOPS_EMAIL and PIPEOPS_PASSWORD environment variables")
}
```

### 2. Implement Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, _, err := client.Auth.Login(ctx, loginReq)
```

### 3. Handle Errors Gracefully

```go
resp, httpResp, err := client.Auth.Login(ctx, loginReq)
if err != nil {
    // Don't expose internal errors to end users
    log.Printf("Login error: %v", err)
    
    // Return user-friendly message
    if httpResp != nil && httpResp.StatusCode == 401 {
        return errors.New("invalid credentials")
    }
    return errors.New("login failed, please try again")
}
```

### 4. Validate Input

```go
func validateCredentials(email, password string) error {
    if email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(email, "@") {
        return errors.New("invalid email format")
    }
    if len(password) < 6 {
        return errors.New("password must be at least 6 characters")
    }
    return nil
}
```

### 5. Secure Token Storage

```go
// For CLI applications
func saveToken(token string) error {
    home, _ := os.UserHomeDir()
    tokenFile := filepath.Join(home, ".pipeops", "token")
    
    // Create directory with restricted permissions
    os.MkdirAll(filepath.Dir(tokenFile), 0700)
    
    // Write token with restricted permissions
    return os.WriteFile(tokenFile, []byte(token), 0600)
}

func loadToken() (string, error) {
    home, _ := os.UserHomeDir()
    tokenFile := filepath.Join(home, ".pipeops", "token")
    
    data, err := os.ReadFile(tokenFile)
    if err != nil {
        return "", err
    }
    
    return string(data), nil
}
```

## Error Handling

Common authentication errors:

```go
resp, httpResp, err := client.Auth.Login(ctx, loginReq)
if err != nil {
    if httpResp != nil {
        switch httpResp.StatusCode {
        case 400:
            fmt.Println("Invalid request format")
        case 401:
            fmt.Println("Invalid email or password")
        case 403:
            fmt.Println("Account is disabled")
        case 429:
            fmt.Println("Too many login attempts")
        case 500:
            fmt.Println("Server error")
        default:
            fmt.Printf("Unexpected error: %d\n", httpResp.StatusCode)
        }
    }
    return err
}
```

## Next Steps

- [OAuth 2.0](oauth.md) - Learn about OAuth authentication
- [Service Tokens](../api-services/servicetokens.md) - Use service account tokens
- [Error Handling](../advanced/error-handling.md) - Handle errors effectively
