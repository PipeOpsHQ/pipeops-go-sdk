# Workspaces Service

The Workspaces Service manages workspace organization and settings.

## Overview

```go
// Access the workspaces service
workspacesService := client.Workspaces
```

## Methods

### List Workspaces

List all workspaces:

```go
workspaces, _, err := client.Workspaces.List(ctx)
if err != nil {
    log.Fatalf("Failed to list workspaces: %v", err)
}

for _, workspace := range workspaces.Data.Workspaces {
    fmt.Printf("- %s\n", workspace.Name)
}
```

### Create Workspace

Create a new workspace:

```go
workspace, _, err := client.Workspaces.Create(ctx, &pipeops.CreateWorkspaceRequest{
    Name:        "Production Workspace",
    Description: "Production environment",
})
if err != nil {
    log.Fatalf("Failed to create workspace: %v", err)
}

fmt.Printf("Created workspace: %s\n", workspace.Data.Workspace.UUID)
```

### Get Workspace

Get workspace details:

```go
workspace, _, err := client.Workspaces.Get(ctx, "workspace-uuid")
if err != nil {
    log.Fatalf("Failed to get workspace: %v", err)
}

fmt.Printf("Workspace: %s\n", workspace.Data.Workspace.Name)
```

### Update Workspace

Update workspace information:

```go
updated, _, err := client.Workspaces.Update(ctx, workspaceUUID, &pipeops.UpdateWorkspaceRequest{
    Name:        "Updated Name",
    Description: "Updated description",
})
```

### Delete Workspace

Delete a workspace:

```go
_, err := client.Workspaces.Delete(ctx, "workspace-uuid")
```

### Set Billing Email

Set billing email for workspace:

```go
_, err := client.Workspaces.SetBillingEmail(ctx, workspaceUUID, &pipeops.SetBillingEmailRequest{
    Email: "billing@example.com",
})
```

## Data Types

```go
type Workspace struct {
    ID          string `json:"id,omitempty"`
    UUID        string `json:"uuid,omitempty"`
    Name        string `json:"name,omitempty"`
    Description string `json:"description,omitempty"`
}
```

## See Also

- [Projects Service](projects.md)
- [Teams Service](teams.md)
- [Billing Service](billing.md)
