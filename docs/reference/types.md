# Types Reference

Common types used throughout the SDK.

## Timestamp

```go
type Timestamp struct {
    time.Time
}
```

Handles various date/time formats from the API.

## User

```go
type User struct {
    ID            string     `json:"id,omitempty"`
    UUID          string     `json:"uuid,omitempty"`
    Email         string     `json:"email,omitempty"`
    FirstName     string     `json:"first_name,omitempty"`
    LastName      string     `json:"last_name,omitempty"`
    IsActive      bool       `json:"is_active,omitempty"`
    EmailVerified bool       `json:"email_verified,omitempty"`
    CreatedAt     *Timestamp `json:"created_at,omitempty"`
    UpdatedAt     *Timestamp `json:"updated_at,omitempty"`
}
```

## Project

```go
type Project struct {
    ID            string     `json:"id,omitempty"`
    UUID          string     `json:"uuid,omitempty"`
    Name          string     `json:"name,omitempty"`
    Description   string     `json:"description,omitempty"`
    Status        string     `json:"status,omitempty"`
    ServerID      string     `json:"server_id,omitempty"`
    EnvironmentID string     `json:"environment_id,omitempty"`
    Repository    string     `json:"repository,omitempty"`
    Branch        string     `json:"branch,omitempty"`
    CreatedAt     *Timestamp `json:"created_at,omitempty"`
    UpdatedAt     *Timestamp `json:"updated_at,omitempty"`
}
```

## Server

```go
type Server struct {
    ID       string     `json:"id,omitempty"`
    UUID     string     `json:"uuid,omitempty"`
    Name     string     `json:"name,omitempty"`
    Provider string     `json:"provider,omitempty"`
    Region   string     `json:"region,omitempty"`
    Status   string     `json:"status,omitempty"`
    CreatedAt *Timestamp `json:"created_at,omitempty"`
}
```

## See Also

- [API Services](../api-services/overview.md)
