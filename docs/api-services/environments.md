# Environments Service

The Environments Service manages environment configurations for projects.

## Overview

```go
// Access the environments service
environmentsService := client.Environments
```

## Methods

### List Environments

List all environments:

```go
environments, _, err := client.Environments.List(ctx)
if err != nil {
    log.Fatalf("Failed to list environments: %v", err)
}

for _, env := range environments.Data.Environments {
    fmt.Printf("- %s (%s)\n", env.Name, env.UUID)
}
```

### Get Environment

Get a specific environment:

```go
env, _, err := client.Environments.Get(ctx, "environment-uuid")
if err != nil {
    log.Fatalf("Failed to get environment: %v", err)
}

fmt.Printf("Environment: %s\n", env.Data.Environment.Name)
```

## Data Types

```go
type Environment struct {
    ID   string `json:"id,omitempty"`
    UUID string `json:"uuid,omitempty"`
    Name string `json:"name,omitempty"`
}
```

## See Also

- [Projects Service](projects.md)
- [Servers Service](servers.md)
