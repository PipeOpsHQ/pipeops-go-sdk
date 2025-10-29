# Cloud Providers Service

The Cloud Providers Service manages cloud provider integrations.

## Overview

```go
// Access the cloud providers service
cloudProvidersService := client.CloudProviders
```

## Methods

### List Cloud Providers

List available cloud providers:

```go
providers, _, err := client.CloudProviders.List(ctx)
if err != nil {
    log.Fatalf("Failed to list providers: %v", err)
}

for _, provider := range providers.Data.Providers {
    fmt.Printf("- %s: %s\n", provider.Name, provider.Status)
}
```

### Connect Provider

Connect a cloud provider account:

```go
connection, _, err := client.CloudProviders.Connect(ctx, &pipeops.ConnectProviderRequest{
    Provider:    "aws",
    Credentials: map[string]string{
        "access_key": "AKIAIOSFODNN7EXAMPLE",
        "secret_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
    },
    Region: "us-east-1",
})
```

### List Regions

List available regions for a provider:

```go
regions, _, err := client.CloudProviders.ListRegions(ctx, "aws")
```

## See Also

- [Servers Service](servers.md)
