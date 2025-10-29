# Servers Service

The Servers Service manages server clusters, agents, and service tokens for infrastructure management.

## Overview

```go
// Access the servers service
serversService := client.Servers
```

## Methods

### List Servers

List all servers:

```go
servers, _, err := client.Servers.List(ctx)
if err != nil {
    log.Fatalf("Failed to list servers: %v", err)
}

fmt.Printf("Found %d servers\n", len(servers.Data.Servers))
for _, server := range servers.Data.Servers {
    fmt.Printf("- %s (%s) - Provider: %s\n", server.Name, server.UUID, server.Provider)
}
```

### Get Server

Get a specific server by UUID:

```go
server, _, err := client.Servers.Get(ctx, "server-uuid")
if err != nil {
    log.Fatalf("Failed to get server: %v", err)
}

fmt.Printf("Server: %s\n", server.Data.Server.Name)
fmt.Printf("Provider: %s\n", server.Data.Server.Provider)
fmt.Printf("Region: %s\n", server.Data.Server.Region)
fmt.Printf("Status: %s\n", server.Data.Server.Status)
```

### Create Server

Create a new server cluster:

```go
newServer, _, err := client.Servers.Create(ctx, &pipeops.CreateServerRequest{
    Name:     "Production Cluster",
    Provider: "aws",
    Region:   "us-east-1",
    Size:     "t3.medium",
})
if err != nil {
    log.Fatalf("Failed to create server: %v", err)
}

fmt.Printf("Created server: %s\n", newServer.Data.Server.UUID)
```

### Delete Server

Delete a server:

```go
_, err := client.Servers.Delete(ctx, "server-uuid")
if err != nil {
    log.Fatalf("Failed to delete server: %v", err)
}

fmt.Println("Server deleted successfully")
```

### Register Agent

Register a PipeOps agent on a cluster:

```go
agent, _, err := client.Servers.RegisterAgent(ctx, &pipeops.AgentRegisterRequest{
    ClusterUUID: "cluster-uuid",
    AgentToken:  "agent-token",
})
if err != nil {
    log.Fatalf("Failed to register agent: %v", err)
}

fmt.Printf("Agent registered: %s\n", agent.Data.AgentID)
```

### Agent Heartbeat

Send agent heartbeat to maintain connection:

```go
_, err := client.Servers.AgentHeartbeat(ctx, clusterUUID, &pipeops.AgentHeartbeatRequest{
    Status:  "healthy",
    Metrics: map[string]interface{}{
        "cpu":    45.2,
        "memory": 62.8,
    },
})
if err != nil {
    log.Fatalf("Heartbeat failed: %v", err)
}
```

### Get Agent Config

Get agent configuration:

```go
config, _, err := client.Servers.GetAgentConfig(ctx, "cluster-uuid")
if err != nil {
    log.Fatalf("Failed to get config: %v", err)
}

fmt.Printf("Agent config: %+v\n", config.Data)
```

### Get Cluster Connection

Get cluster connection details:

```go
connection, _, err := client.Servers.GetClusterConnection(ctx, "cluster-uuid")
if err != nil {
    log.Fatalf("Failed to get connection: %v", err)
}

fmt.Printf("Endpoint: %s\n", connection.Data.Endpoint)
```

### Create Service Token

Create a service account token:

```go
token, _, err := client.Servers.CreateServiceToken(ctx, &pipeops.ServiceTokenRequest{
    Name:        "CI/CD Token",
    Description: "For automated deployments",
    Scopes:      []string{"projects:write", "deployments:create"},
})
if err != nil {
    log.Fatalf("Failed to create token: %v", err)
}

fmt.Printf("Service token: %s\n", token.Data.Token)
```

### List Service Tokens

List all service tokens:

```go
tokens, _, err := client.Servers.ListServiceTokens(ctx)
if err != nil {
    log.Fatalf("Failed to list tokens: %v", err)
}

for _, token := range tokens.Data.Tokens {
    fmt.Printf("- %s (%s)\n", token.Name, token.UUID)
}
```

### Revoke Service Token

Revoke a service token:

```go
_, err := client.Servers.RevokeServiceToken(ctx, "token-uuid")
if err != nil {
    log.Fatalf("Failed to revoke token: %v", err)
}

fmt.Println("Token revoked")
```

## Data Types

### Server

```go
type Server struct {
    ID       string     `json:"id,omitempty"`
    UUID     string     `json:"uuid,omitempty"`
    Name     string     `json:"name,omitempty"`
    Provider string     `json:"provider,omitempty"`
    Region   string     `json:"region,omitempty"`
    Size     string     `json:"size,omitempty"`
    Status   string     `json:"status,omitempty"`
    CreatedAt *Timestamp `json:"created_at,omitempty"`
}
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
    client.SetToken("your-token")
    
    ctx := context.Background()
    
    // List all servers
    servers, _, err := client.Servers.List(ctx)
    if err != nil {
        log.Fatalf("Failed to list servers: %v", err)
    }
    
    fmt.Printf("Found %d servers:\n", len(servers.Data.Servers))
    for _, server := range servers.Data.Servers {
        fmt.Printf("- %s: %s on %s\n", server.Name, server.Provider, server.Region)
    }
    
    // Get details of first server
    if len(servers.Data.Servers) > 0 {
        serverUUID := servers.Data.Servers[0].UUID
        server, _, err := client.Servers.Get(ctx, serverUUID)
        if err != nil {
            log.Fatalf("Failed to get server: %v", err)
        }
        
        fmt.Printf("\nServer Details:\n")
        fmt.Printf("Name: %s\n", server.Data.Server.Name)
        fmt.Printf("Status: %s\n", server.Data.Server.Status)
    }
}
```

## See Also

- [Projects Service](projects.md) - Deploy projects to servers
- [Cloud Providers Service](cloudproviders.md) - Manage cloud provider connections
- [Service Tokens Service](servicetokens.md) - Service account management
