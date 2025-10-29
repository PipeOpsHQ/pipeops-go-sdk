# Service Tokens Service

The Service Tokens Service manages service account tokens for API access.

## Overview

```go
// Access the service tokens service
serviceTokensService := client.ServiceTokens
```

## Methods

### Create Service Account Token

Create a new service account token:

```go
token, _, err := client.ServiceTokens.CreateServiceAccountToken(ctx, &pipeops.ServiceAccountTokenRequest{
    Name:        "CI/CD Token",
    Description: "For automated deployments",
    Scopes:      []string{"projects:write", "deployments:create"},
    ExpiresIn:   "30d",
})
if err != nil {
    log.Fatalf("Failed to create token: %v", err)
}

fmt.Printf("Token: %s\n", token.Data.Token)
fmt.Printf("Token ID: %s\n", token.Data.UUID)
```

### List Service Account Tokens

List all service account tokens:

```go
tokens, _, err := client.ServiceTokens.ListServiceAccountTokens(ctx)
if err != nil {
    log.Fatalf("Failed to list tokens: %v", err)
}

for _, token := range tokens.Data.Tokens {
    fmt.Printf("- %s (%s)\n", token.Name, token.UUID)
}
```

### Get Service Account Token

Get token details:

```go
token, _, err := client.ServiceTokens.GetServiceAccountToken(ctx, "token-uuid")
```

### Update Service Account Token

Update token metadata:

```go
updated, _, err := client.ServiceTokens.UpdateServiceAccountToken(ctx, tokenUUID, &pipeops.ServiceAccountTokenUpdateRequest{
    Name:        "Updated Name",
    Description: "Updated description",
})
```

### Revoke Service Account Token

Revoke a token:

```go
_, err := client.ServiceTokens.RevokeServiceAccountToken(ctx, "token-uuid")
if err != nil {
    log.Fatalf("Failed to revoke token: %v", err)
}

fmt.Println("Token revoked successfully")
```

## Complete Example

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
    client.SetToken("your-admin-token")
    
    ctx := context.Background()
    
    // Create service token for CI/CD
    token, _, err := client.ServiceTokens.CreateServiceAccountToken(ctx, &pipeops.ServiceAccountTokenRequest{
        Name:        "GitHub Actions",
        Description: "Token for GitHub Actions CI/CD",
        Scopes:      []string{"projects:write", "deployments:create"},
    })
    if err != nil {
        log.Fatalf("Failed to create token: %v", err)
    }
    
    fmt.Printf("Created token: %s\n", token.Data.UUID)
    fmt.Printf("Save this token securely: %s\n", token.Data.Token)
    
    // List all tokens
    tokens, _, err := client.ServiceTokens.ListServiceAccountTokens(ctx)
    if err != nil {
        log.Fatalf("Failed to list tokens: %v", err)
    }
    
    fmt.Printf("\nActive tokens:\n")
    for _, t := range tokens.Data.Tokens {
        fmt.Printf("- %s: %s\n", t.Name, t.CreatedAt)
    }
}
```

## Best Practices

1. **Scope Tokens**: Only grant necessary permissions
2. **Rotate Regularly**: Create new tokens periodically
3. **Store Securely**: Keep tokens in secure storage
4. **Monitor Usage**: Track token usage and revoke unused tokens

## See Also

- [Authentication Overview](../authentication/overview.md)
- [Servers Service](servers.md)
