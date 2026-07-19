# API coverage plan: go-sdk → MCP → CLI

**Date:** 2026-07-18  
**Goal:** Close gaps between dual-auth product APIs and automation surfaces in a fixed cascade.

```text
Controller (source of truth)
    → pipeops-go-sdk  (typed clients + path contracts)
        → pipeops-mcp   (tools over SDK)
            → pipeops-cli (commands over SDK / thin wrappers)
```

Never invent paths in MCP/CLI without an SDK method. Never expose admin-only or SA-denylisted surfaces for platform tokens without an explicit product decision.

---

## 1. Current inventory (approx.)

| Layer | Surface size | Notes |
|-------|--------------|--------|
| **Controller dual-auth** | Project, cluster/server, environment, workspace (limited), team (limited), addons, external registry, service tokens | JWT session **or** `sat_*` with `api:*` |
| **go-sdk** | Broad: projects, servers, addons, billing, workspaces, teams, auth, VCS helpers, **+ `additional.go` speculative services** | Quality uneven: many methods exist; some paths wrong/stale |
| **MCP** | ~89 tools | Strong on projects/addons/billing/tokens/VCS/registries; thin on networks/metrics/migrate/groups/volumes/gitops |
| **CLI** | Thin product surface | Strong on agent/login/workspace/list/env; weak vs MCP (no metrics, network, migrate, gitops, volumes, project groups) |

### Dual-auth route groups (automation-relevant)

| Group | Dual-auth | SA notes |
|-------|-----------|----------|
| `/project/*` | Yes | Env dump GET denied for SA |
| `/cluster/*` | Yes | kubeconfig/credentials denied |
| `/environment/*` | Yes | env dump/mutators restricted |
| `/workspace/*` | Yes | create/delete/SSO/billing-email denied for SA |
| `/team/*` | Yes | invite/membership mutations denied |
| `/addons/*` (product) | Yes | Overview fixed for team/SA |
| `/api/v1/external-registry` | Yes | |
| `/api/v1/service-account-tokens` | User JWT for CRUD; tokens call dual-auth APIs | |
| `/project-groups/*` | Product UI primarily | **No SDK** |
| `/volumes/*` | Product | **No SDK** |
| `/api/v1/gitops/*` | Product | **No SDK** |
| `/terraform/*` | Product | **No SDK** |
| Addon backups export | Under `/addons/deployments/.../backups` | **No real SDK** (`additional.go` uses wrong paths) |

---

## 2. Gap categories

### A. Missing in SDK (controller has real APIs)

| Domain | Controller (examples) | SDK today | Priority for automation |
|--------|----------------------|-----------|-------------------------|
| **Project groups / plane** | `/project-groups` list/create/topology/env inject/redeploy | None | **P1** (platform feature) |
| **Workspace volumes** | `/volumes` list/get/remount/export | None | **P1** |
| **GitOps** | `/api/v1/gitops/applications` CRUD + sync | None | **P1** |
| **Addon backup export** | `GET/POST .../addons/deployments/:id/backups*` | Stub `backups/projects/...` in `additional.go` | **P1** |
| **Terraform** | `/terraform` CRUD | None | **P2** |
| **Project overview** | `GET project/overview/:uuid` (`projectURL`) | Partial via fetch `public_url` | **P0 quality** (already partially fixed) |
| **Cluster nodes/events/metrics** | `/cluster/nodes|events|metrics|insight|overview` | Partial via servers | **P2** |
| **PR previews** | dedicated routes | None | **P2** |
| **Object store** | if dual-auth | None | **P3** |
| **Real audit log** | product audit tables | Stub `AuditLogService` | **P2** |

### B. In SDK but wrong / speculative (`additional.go`)

Treat as **debt**, not coverage:

- `BackupService` paths ≠ addon backup export API  
- Alerts / templates / security scan / generic audit may not match controller  
- Prefer delete or rewrite against real routes before MCP/CLI bind to them  

### C. In SDK but missing from MCP tools

High-value **ProjectService** methods not exposed as tools (non-exhaustive):

| SDK method | MCP tool today | CLI today |
|------------|----------------|-----------|
| `GetNetworkSettings` / update port / network policies | No | No |
| `UpdateDomain` / `DeleteCustomDomain` / SSL check | Partial (addon domain only) | No |
| `MigrateProject` | No | No |
| `GetCosts` / metrics family | No | No |
| `BulkDelete` | No | No |
| `GetPodsFromLabel` / runtime logs by pod | Partial logs only | logs only |
| `DeployFromImage` | Yes | Weak/create path |
| Addon `UpdateDeployment` / `SyncDeployment` / bulk delete | No | Partial delete only |
| Cluster cost allocation | Yes (`get_server_cost_allocation`) | Partial |

### D. In MCP but thin/missing in CLI

| MCP tool area | CLI |
|---------------|-----|
| Full project CRUD + deploy/restart/stop | Partial (`project get/create/update/delete/deploy`, list) |
| Env get/set | Yes (`project env`, mask/`--reveal`) |
| VCS list/search/link | No dedicated commands |
| External registries | No |
| Cloud provider discovery | No |
| Billing/subscribe/cards | No (intentional for many CLIs) |
| Service tokens | Yes (`token`) |
| Addon deploy/list/get | Yes (partial) |
| Workspaces CRUD | Yes (partial) |

---

## 3. Design rules for the cascade

1. **SDK first**  
   - Path strings match controller routes (verify with integration or postman route tests).  
   - Always pass `workspace` / `workspace_uuid` where dual-auth requires it.  
   - Prefer typed responses; avoid silent wrong unmarshals (see project URL lesson).  
   - Document SA denylist: method exists but returns 401 for SA on secret paths.

2. **MCP second**  
   - One tool per user intent (not one tool per HTTP method).  
   - Annotations: `readOnlyHint`, `destructiveHint`, `idempotentHint`.  
   - Scope tags for future OAuth: `pipeops:read` vs `projects:write` / `deployments:write` / `addons:write` / `billing:write`.  
   - Do **not** expose: team invite/delete, workspace delete, env secret dump without JWT policy, kubeconfig, SA token create with `tokens:admin` without confirmation.

3. **CLI third**  
   - Commands map 1:1 to high-traffic MCP/SDK workflows operators need in a terminal.  
   - Skip billing card mutation and OAuth-admin surfaces unless product asks.  
   - Always workspace-scoped (`workspace select` + query injection).

---

## 4. Phased delivery plan

### Phase 0 — Hygiene (1 PR, SDK only)

**Goal:** Trustworthy baseline before adding surface area.

| Work | Detail |
|------|--------|
| Audit `additional.go` | Mark/remove dead paths; fix or delete Backup/Alert/Audit stubs |
| Contract tests | Extend `routes_postman_test` / path fixtures for dual-auth groups |
| Workspace injection | Ensure all dual-auth methods send `workspace` / `workspace_uuid` (addons overview pattern) |
| Docs | Generate method → path matrix from code |

**Exit:** No SDK method points at a known-wrong path.

### Phase 1 — SDK “ops core” (P0/P1 product gaps)

**Add or harden packages:**

| Package | Methods (minimum) | Controller source |
|---------|-------------------|-------------------|
| `volumes` | List, Get, Remount, StartExport, GetExport | `/volumes` |
| `gitops` | List/Get/Create/Update/Delete app, TriggerSync, SyncStatus, Diff, History | `/api/v1/gitops` |
| `projectgroups` | List, Get, Create, Topology, Activity, RedeployApps, AttachMember, Connect, Env catalog/inject (read-safe first) | `/project-groups` |
| `addons` (extend) | ListBackups, StartBackupExport, GetExport, DownloadExport | `/addons/deployments/:id/backups*` |
| `projects` (harden) | Ensure Get returns `public_url`; Overview helper if needed | fetch + overview |

**Also fix:**

- Project Get unmarshaling of `public_url` / `CustomDomainName` (partially done).  
- Addon overview workspace (done).  
- Delete/replace false `BackupService` paths.

**Exit:** Integration smoke: list volumes, list gitops apps, list project groups, list addon backups with `sat_*` + JWT.

### Phase 2 — MCP tools on Phase 1 SDK

| Tool set | Examples | Annotations |
|----------|----------|-------------|
| Volumes | `list_volumes`, `get_volume`, `export_volume` | read / write |
| GitOps | `list_gitops_apps`, `get_gitops_app`, `sync_gitops_app`, `get_gitops_diff` | read / write |
| Project groups | `list_project_groups`, `get_project_group`, `get_project_group_topology`, `redeploy_project_group` | read / write |
| Addon backups | `list_addon_backups`, `export_addon_backup`, `get_addon_backup_export` | read / write |
| Project network/domain | `get_project_network`, `update_project_domain`, `delete_project_domain` | write careful |
| Project migrate | `migrate_project` | destructiveHint |

Wire tools **only** through SDK methods. Reuse workspace resolution helpers already in MCP.

**Exit:** MCP can complete: “show me volumes in this workspace”, “sync gitops app X”, “export postgres addon backup”.

### Phase 3 — CLI on Phase 2 tools/SDK

| Command area | Commands |
|--------------|----------|
| `pipeops volume` | `list`, `get`, `export` |
| `pipeops gitops` | `list`, `get`, `sync`, `diff` |
| `pipeops group` / `plane` | `list`, `get`, `topology`, `redeploy` |
| `pipeops addon backup` | `list`, `export`, `status` |
| `pipeops project` extend | `domain`, `network`, `migrate`, ensure `get` shows URL |
| Optional | `pipeops project metrics` (read-only) |

Skip CLI for: billing cards, team invites, workspace delete, token admin unless requested.

**Exit:** Operator runbooks work without dashboard for deploys, volumes, gitops, backups.

### Phase 4 — Secondary SDK/MCP (P2)

| Area | Notes |
|------|-------|
| Terraform module CRUD | `/terraform` |
| Cluster insight/nodes/events | Ops debugging |
| PR preview APIs | If dual-auth approved |
| Real audit log API | Align with `audit_log` product API, not stub |
| Metrics tools in MCP | CPU/memory/network — often large payloads; paginate |

### Phase 5 — Explicit non-goals (unless product reopens)

| Surface | Reason |
|---------|--------|
| Admin `/admin/*` | Not platform SA |
| Observability boards | SA denylist until hard-scope |
| DB/Mongo/Redis studio | Interactive, credential-heavy |
| Kubeconfig download | SA denylist |
| Team invite/accept | SA denylist |
| Full billing checkout UI flows | Prefer portal URL already in MCP |

---

## 5. Suggested PR stack (Graphite-friendly)

```text
sdk/0-hygiene-additional-go
  └─ sdk/1-volumes
      └─ sdk/2-gitops
          └─ sdk/3-project-groups
              └─ sdk/4-addon-backups
                  └─ mcp/1-tools-volumes-gitops-groups-backups
                      └─ mcp/2-tools-project-network-domain-migrate
                          └─ cli/1-volume-gitops-group
                              └─ cli/2-addon-backup-project-extend
```

Versioning:

- Tag SDK after Phase 1 (`v0.13.x`).  
- Bump MCP + CLI deps to that tag before shipping tools/commands.

---

## 6. Quality bar (every PR)

| Check | Requirement |
|-------|-------------|
| Path contract | Unit test or route fixture for method → HTTP path + method |
| Workspace | Query/body includes bound workspace; SA hard-scope compatible |
| Errors | Surface 401/403/404 without swallowing |
| Secrets | No env/credential dumps for SA; CLI mask/`--reveal` pattern where values shown |
| Docs | README/API_REFERENCE snippet for new methods/tools/commands |
| SA denylist | Tests that denied paths stay denied |

---

## 7. Immediate next actions (recommended order)

1. **SDK Phase 0 hygiene** — audit `additional.go`, fix/remove false Backup paths.  
2. **SDK Phase 1** — volumes + gitops + project groups + real addon backups (four focused PRs).  
3. **MCP Phase 2** — tools for those four domains + annotations.  
4. **CLI Phase 3** — operator commands for the same four.  
5. Only then metrics/terraform/PR-preview.

---

## 8. Success metrics

| Metric | Target |
|--------|--------|
| Dual-auth product domains with first-class SDK package | volumes, gitops, project-groups, addon-backups |
| MCP tools covering those domains | ≥ 12 new tools |
| CLI commands for those domains | ≥ 10 new subcommands |
| Known-wrong SDK paths | 0 |
| SA smoke suite | list projects/servers/addons/volumes/groups with `sat_*` |

---

## Appendix — cascade dependency

```text
Missing controller API? → open controller dual-auth issue first (not SDK fiction).
Missing SDK method? → block MCP tool and CLI command.
SDK method exists but wrong path? → fix SDK before any consumer PR.
MCP tool without CLI? → OK if low terminal demand (e.g. billing cards).
CLI without MCP? → avoid; prefer MCP parity for AI + human ops.
```
