# Teams Service

The Teams Service manages team collaboration and member management.

## Overview

```go
// Access the teams service
teamsService := client.Teams
```

## Methods

### List Teams

List all teams:

```go
teams, _, err := client.Teams.List(ctx)
if err != nil {
    log.Fatalf("Failed to list teams: %v", err)
}

for _, team := range teams.Data.Teams {
    fmt.Printf("- %s\n", team.Name)
}
```

### Create Team

Create a new team:

```go
team, _, err := client.Teams.Create(ctx, &pipeops.CreateTeamRequest{
    Name:        "Development Team",
    Description: "Core development team",
})
if err != nil {
    log.Fatalf("Failed to create team: %v", err)
}

fmt.Printf("Created team: %s\n", team.Data.Team.UUID)
```

### Get Team

Get team details:

```go
team, _, err := client.Teams.Get(ctx, "team-uuid")
if err != nil {
    log.Fatalf("Failed to get team: %v", err)
}

fmt.Printf("Team: %s\n", team.Data.Team.Name)
```

### Update Team

Update team information:

```go
updated, _, err := client.Teams.Update(ctx, teamUUID, &pipeops.UpdateTeamRequest{
    Name:        "Updated Team Name",
    Description: "Updated description",
})
```

### Delete Team

Delete a team:

```go
_, err := client.Teams.Delete(ctx, "team-uuid")
```

### Invite Member

Invite a member to the team:

```go
_, err := client.Teams.InviteMember(ctx, teamUUID, &pipeops.InviteTeamMemberRequest{
    Email: "member@example.com",
    Role:  "developer",
})
if err != nil {
    log.Fatalf("Failed to invite member: %v", err)
}

fmt.Println("Invitation sent")
```

### List Members

List team members:

```go
members, _, err := client.Teams.ListMembers(ctx, "team-uuid")
if err != nil {
    log.Fatalf("Failed to list members: %v", err)
}

for _, member := range members.Data.Members {
    fmt.Printf("- %s (%s)\n", member.Email, member.Role)
}
```

### Remove Member

Remove a member from team:

```go
_, err := client.Teams.RemoveMember(ctx, teamUUID, memberUUID)
```

### Update Member Role

Update a member's role:

```go
_, err := client.Teams.UpdateMemberRole(ctx, teamUUID, memberUUID, &pipeops.UpdateMemberRoleRequest{
    Role: "admin",
})
```

### Accept Invitation

Accept a team invitation:

```go
_, err := client.Teams.AcceptInvitation(ctx, "invite-token")
```

### Reject Invitation

Reject a team invitation:

```go
_, err := client.Teams.RejectInvitation(ctx, "invite-token")
```

## Data Types

```go
type Team struct {
    ID          string `json:"id,omitempty"`
    UUID        string `json:"uuid,omitempty"`
    Name        string `json:"name,omitempty"`
    Description string `json:"description,omitempty"`
}
```

## See Also

- [Workspaces Service](workspaces.md)
- [Users Service](users.md)
