# Sandboxes (Rexec BFF)

The Go SDK talks to the **PipeOps controller sandboxes BFF**, not Rexec directly.

- Base path: `GET/POST /api/v1/sandboxes` (alias of `/sandboxes`)
- Auth: user JWT (OAuth / dashboard session) **or** workspace service account `sat_*` with `api:read`/`api:write` or preset **`sandbox`**
- Tenancy: always pass `workspace_uuid` (also sent as `workspace` for dashboard parity)

## Quick start

```go
client, err := pipeops.NewClient("https://api.pipeops.io")
if err != nil {
    log.Fatal(err)
}
client.SetToken(os.Getenv("PIPEOPS_TOKEN")) // sat_* or user JWT

ws := &pipeops.SandboxWorkspaceOptions{WorkspaceUUID: "your-workspace-uuid"}

list, _, err := client.Sandboxes.List(ctx, ws)
// ...

created, _, err := client.Sandboxes.Create(ctx, ws, &pipeops.CreateSandboxRequest{
    Name:  "dev-box",
    Image: "ubuntu",
    Role:  "standard",
})

// Terminal / embed grant (short-lived)
session, _, err := client.Sandboxes.CreateSession(ctx, created.Data.ID, ws)
// session.Data.Token + session.Data.BaseURL — do not log

// Optional: long-lived Rexec API token for external tools (shown once)
mint, _, err := client.Sandboxes.MintAPIToken(ctx, ws, &pipeops.MintRexecAPITokenRequest{
    Name: "my-cli",
})
// store mint.Data.Token securely
```

## Lifecycle

| Method | HTTP |
|--------|------|
| `List` | `GET /api/v1/sandboxes` |
| `Get` | `GET /api/v1/sandboxes/:id` |
| `Create` | `POST /api/v1/sandboxes` |
| `Start` / `Stop` / `Delete` | `POST .../start`, `POST .../stop`, `DELETE .../:id` |
| `Restart` | stop then start (client helper) |
| `CreateSession` | `POST .../:id/session` |
| `Exec` | `POST .../:id/exec` |
| `ListFiles` | `GET .../:id/files` |
| `ReadFile` | `GET .../:id/files/content` |

## BYOS binding & usage

| Method | HTTP |
|--------|------|
| `GetRexecBinding` | `GET .../rexec-binding` |
| `UpsertRexecBinding` | `PUT .../rexec-binding` |
| `DeleteRexecBinding` | `DELETE .../rexec-binding` |
| `UsageDaily` | `GET .../usage/daily?from=&to=` (`YYYY-MM-DD`) |
| `MintAPIToken` | `POST .../api-token` |

## Not in this client

### Run a command inside a sandbox

```go
exec, _, err := client.Sandboxes.Exec(ctx, created.Data.ID, ws, &pipeops.ExecSandboxRequest{
	Command:        "uname -a && pwd",
	TimeoutSeconds: 60,
})
// exec.Data.Output, exec.Data.ExitCode, exec.Data.Stdout, exec.Data.Stderr
```

### List / read files

```go
list, _, err := client.Sandboxes.ListFiles(ctx, created.Data.ID, "/home/user", ws)
// list.Data.Files[i].Name, .Path, .IsDir, .Size

file, _, err := client.Sandboxes.ReadFile(ctx, created.Data.ID, "/home/user/app.go", ws)
// file.Data.Content, file.Data.Encoding ("utf-8" | "base64"), file.Data.Size
```

Upload and raw WebSocket terminal are not on the BFF path used here. Prefer `Exec` for one-shot agent commands, `ListFiles`/`ReadFile` for inspection, `CreateSession` + Rexec UI/terminal for interactive shells, or a direct Rexec client with a minted `rexec_*` token.
