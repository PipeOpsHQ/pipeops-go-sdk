# Admin Service

The Admin Service provides administrative operations (admin access required).

## Overview

```go
// Access the admin service (requires admin privileges)
adminService := client.Admin
```

## Methods

### List Users

List all users (admin only):

```go
users, _, err := client.Admin.ListUsers(ctx, &pipeops.AdminListUsersOptions{
    Page:  1,
    Limit: 50,
})
if err != nil {
    log.Fatalf("Failed to list users: %v", err)
}

for _, user := range users.Data.Users {
    fmt.Printf("- %s (%s)\n", user.Email, user.UUID)
}
```

### Get User

Get user details (admin only):

```go
user, _, err := client.Admin.GetUser(ctx, "user-uuid")
```

### Update User

Update user information (admin only):

```go
updated, _, err := client.Admin.UpdateUser(ctx, userUUID, &pipeops.UpdateUserRequest{
    IsActive: true,
    Role:     "user",
})
```

### Delete User

Delete a user account (admin only):

```go
_, err := client.Admin.DeleteUser(ctx, "user-uuid")
```

### Get Statistics

Get platform statistics:

```go
stats, _, err := client.Admin.GetStats(ctx)
if err != nil {
    log.Fatalf("Failed to get stats: %v", err)
}

fmt.Printf("Total Users: %d\n", stats.Data.TotalUsers)
fmt.Printf("Total Projects: %d\n", stats.Data.TotalProjects)
```

### Get System Health

Get system health metrics:

```go
health, _, err := client.Admin.GetSystemHealth(ctx)
```

### Get Audit Logs

Get audit logs:

```go
logs, _, err := client.Admin.GetAuditLogs(ctx)
```

### Broadcast Message

Send a broadcast message:

```go
_, err := client.Admin.BroadcastMessage(ctx, &pipeops.BroadcastRequest{
    Message: "System maintenance scheduled",
    Level:   "warning",
})
```

### Create Plan

Create a subscription plan:

```go
plan, _, err := client.Admin.CreatePlan(ctx, &pipeops.CreatePlanRequest{
    Name:        "Pro Plan",
    Description: "Professional tier",
    Price:       49.99,
    Features:    []string{"feature1", "feature2"},
})
```

### Subscribe User

Subscribe a user to a plan:

```go
_, err := client.Admin.SubscribeUser(ctx, &pipeops.SubscribeUserRequest{
    UserUUID: "user-uuid",
    PlanUUID: "plan-uuid",
})
```

## See Also

- [Users Service](users.md)
- [Billing Service](billing.md)
