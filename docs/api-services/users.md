# Users Service

The Users Service manages user profile and settings.

## Overview

```go
// Access the users service
usersService := client.Users
```

## Methods

### Get Profile

Get current user profile:

```go
profile, _, err := client.Users.GetProfile(ctx)
if err != nil {
    log.Fatalf("Failed to get profile: %v", err)
}

fmt.Printf("Name: %s %s\n", profile.Data.User.FirstName, profile.Data.User.LastName)
fmt.Printf("Email: %s\n", profile.Data.User.Email)
```

### Update Profile

Update user profile:

```go
updated, _, err := client.Users.UpdateProfile(ctx, &pipeops.UpdateProfileRequest{
    FirstName: "John",
    LastName:  "Doe",
    Phone:     "+1234567890",
})
if err != nil {
    log.Fatalf("Failed to update profile: %v", err)
}

fmt.Println("Profile updated successfully")
```

### Get Settings

Get user settings:

```go
settings, _, err := client.Users.GetSettings(ctx)
if err != nil {
    log.Fatalf("Failed to get settings: %v", err)
}

fmt.Printf("Theme: %s\n", settings.Data.Settings.Theme)
```

### Update Settings

Update user settings:

```go
_, err := client.Users.UpdateSettings(ctx, &pipeops.UpdateSettingsRequest{
    Theme:       "dark",
    Language:    "en",
    Timezone:    "UTC",
})
```

### Update Notification Settings

Update notification preferences:

```go
_, err := client.Users.UpdateNotificationSettings(ctx, &pipeops.UpdateNotificationSettingsRequest{
    EmailNotifications: true,
    DeploymentAlerts:   true,
    SecurityAlerts:     true,
})
```

### Reset Secret Token

Reset user's secret token:

```go
token, _, err := client.Users.ResetSecretToken(ctx)
if err != nil {
    log.Fatalf("Failed to reset token: %v", err)
}

fmt.Printf("New token: %s\n", token.Data.Token)
```

### Delete Profile

Request profile deletion:

```go
_, err := client.Users.DeleteProfile(ctx)
```

### Cancel Profile Deletion

Cancel a pending profile deletion:

```go
_, err := client.Users.CancelProfileDeletion(ctx)
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
    
    // Get current profile
    profile, _, err := client.Users.GetProfile(ctx)
    if err != nil {
        log.Fatalf("Failed to get profile: %v", err)
    }
    
    fmt.Printf("Current Profile:\n")
    fmt.Printf("Name: %s %s\n", 
        profile.Data.User.FirstName, 
        profile.Data.User.LastName)
    fmt.Printf("Email: %s\n", profile.Data.User.Email)
    
    // Update profile
    _, err = client.Users.UpdateProfile(ctx, &pipeops.UpdateProfileRequest{
        FirstName: "Updated",
        LastName:  "Name",
    })
    if err != nil {
        log.Fatalf("Failed to update: %v", err)
    }
    
    fmt.Println("\nProfile updated successfully")
}
```

## See Also

- [Auth Service](auth.md)
- [Teams Service](teams.md)
