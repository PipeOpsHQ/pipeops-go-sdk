# Add-ons Service

The Add-ons Service manages marketplace add-on deployments.

## Overview

```go
// Access the add-ons service
addonsService := client.AddOns
```

## Methods

### List Add-ons

List available add-ons:

```go
addons, _, err := client.AddOns.List(ctx)
if err != nil {
    log.Fatalf("Failed to list add-ons: %v", err)
}

for _, addon := range addons.Data.AddOns {
    fmt.Printf("- %s: %s\n", addon.Name, addon.Description)
}
```

### Get Add-on

Get specific add-on details:

```go
addon, _, err := client.AddOns.Get(ctx, "addon-uuid")
if err != nil {
    log.Fatalf("Failed to get add-on: %v", err)
}

fmt.Printf("Add-on: %s\n", addon.Data.AddOn.Name)
fmt.Printf("Version: %s\n", addon.Data.AddOn.Version)
```

### Deploy Add-on

Deploy an add-on:

```go
deployment, _, err := client.AddOns.Deploy(ctx, &pipeops.DeployAddOnRequest{
    AddOnUUID:   "addon-uuid",
    ProjectUUID: "project-uuid",
    Config: map[string]interface{}{
        "memory": "512Mi",
        "cpu":    "250m",
    },
})
if err != nil {
    log.Fatalf("Deployment failed: %v", err)
}

fmt.Printf("Deployed: %s\n", deployment.Data.Deployment.UUID)
```

### List Deployments

List add-on deployments:

```go
deployments, _, err := client.AddOns.ListDeployments(ctx)
if err != nil {
    log.Fatalf("Failed to list deployments: %v", err)
}

for _, deployment := range deployments.Data.Deployments {
    fmt.Printf("- %s: %s\n", deployment.Name, deployment.Status)
}
```

### Get Deployment

Get deployment details:

```go
deployment, _, err := client.AddOns.GetDeployment(ctx, "deployment-uuid")
```

### Update Deployment

Update deployment configuration:

```go
updated, _, err := client.AddOns.UpdateDeployment(ctx, deploymentUUID, &pipeops.UpdateDeploymentRequest{
    Config: map[string]interface{}{
        "replicas": 3,
    },
})
```

### Delete Deployment

Delete an add-on deployment:

```go
_, err := client.AddOns.DeleteDeployment(ctx, "deployment-uuid")
```

### Sync Deployment

Sync deployment state:

```go
_, err := client.AddOns.SyncDeployment(ctx, "deployment-uuid")
```

### List Categories

List add-on categories:

```go
categories, _, err := client.AddOns.ListCategories(ctx)
```

### Submit Add-on

Submit an add-on to marketplace:

```go
submission, _, err := client.AddOns.SubmitAddOn(ctx, &pipeops.AddOnSubmissionRequest{
    Name:        "My Add-on",
    Description: "A useful add-on",
    Version:     "1.0.0",
    Category:    "database",
})
```

## See Also

- [Projects Service](projects.md)
