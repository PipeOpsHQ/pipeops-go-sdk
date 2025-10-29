# Complete Examples

Real-world examples demonstrating SDK usage.

## CI/CD Integration

Automate deployments from CI/CD pipelines:

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
    // Get credentials from environment
    token := os.Getenv("PIPEOPS_TOKEN")
    projectUUID := os.Getenv("PROJECT_UUID")
    
    if token == "" || projectUUID == "" {
        log.Fatal("PIPEOPS_TOKEN and PROJECT_UUID required")
    }
    
    // Create client
    client, _ := pipeops.NewClient("")
    client.SetToken(token)
    
    ctx := context.Background()
    
    // Deploy project
    fmt.Println("Starting deployment...")
    deployment, _, err := client.Projects.Deploy(ctx, projectUUID)
    if err != nil {
        log.Fatalf("Deployment failed: %v", err)
    }
    
    fmt.Printf("Deployment started: %s\n", deployment.Data.DeploymentID)
    
    // Wait for deployment to complete
    // (Implementation depends on your needs)
    
    fmt.Println("Deployment complete!")
}
```

## Infrastructure Management

Manage servers and projects:

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
    
    // Create server
    server, _, err := client.Servers.Create(ctx, &pipeops.CreateServerRequest{
        Name:     "Production",
        Provider: "aws",
        Region:   "us-east-1",
    })
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }
    
    serverUUID := server.Data.Server.UUID
    fmt.Printf("Created server: %s\n", serverUUID)
    
    // Create project on server
    project, _, err := client.Projects.Create(ctx, &pipeops.CreateProjectRequest{
        Name:       "My App",
        ServerID:   serverUUID,
        Repository: "https://github.com/user/app",
        Branch:     "main",
    })
    if err != nil {
        log.Fatalf("Failed to create project: %v", err)
    }
    
    fmt.Printf("Created project: %s\n", project.Data.Project.UUID)
}
```

## Monitoring and Alerts

Monitor project logs and metrics:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
    client, _ := pipeops.NewClient("")
    client.SetToken("your-token")
    
    projectUUID := "project-uuid"
    ctx := context.Background()
    
    // Poll logs every 30 seconds
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        logs, _, err := client.Projects.GetLogs(ctx, projectUUID, &pipeops.LogsOptions{
            Limit:  100,
            Search: "error",
        })
        if err != nil {
            log.Printf("Failed to get logs: %v", err)
            continue
        }
        
        if len(logs.Data.Logs) > 0 {
            fmt.Printf("Found %d errors\n", len(logs.Data.Logs))
            // Send alert
        }
    }
}
```

## See Also

- [Basic Example](https://github.com/PipeOpsHQ/pipeops-go-sdk/tree/main/examples/basic)
- [OAuth Example](https://github.com/PipeOpsHQ/pipeops-go-sdk/tree/main/examples/oauth)
