# Miscellaneous Service

The Misc Service provides various utility endpoints.

## Overview

```go
// Access the misc service
miscService := client.Misc
```

## Methods

### Health Check

Check API health:

```go
health, _, err := client.Misc.HealthCheck(ctx)
if err != nil {
    log.Fatalf("Health check failed: %v", err)
}

fmt.Printf("Status: %s\n", health.Status)
```

### Get Version

Get API version information:

```go
version, _, err := client.Misc.GetVersion(ctx)
fmt.Printf("API Version: %s\n", version.Data.Version)
```

## See Also

- [Admin Service](admin.md)
