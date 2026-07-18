# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `ApplyCreateProjectDefaults` — prefer-client create defaults (PORT from network only if missing, source/environment/protocol/worker gaps only)
- `VolumeService` (`Client.Volumes`) for workspace PVC list/get/remount/delete/export against `/volumes`
- Real addon backup export methods on `AddOnService`: `ListAddonBackups`, `StartAddonBackupExport`, `GetAddonBackupExport`, `DownloadAddonBackupExport`
- Path contract tests for volumes and addon backups
- `docs/API_COVERAGE_PLAN.md` cascade plan (SDK → MCP → CLI)
- `GitOpsService` (`Client.GitOps`) for GitOps application CRUD, sync, sync-status, diff, and history against `/api/v1/gitops/applications`
- `ProjectGroupService` (`Client.ProjectGroups`) for project group plane P1 APIs against `/project-groups` (list/get/create/update/delete, members, topology, shared env, connect, redeploy, resolve, candidates)
- Path contract tests for GitOps and Project Groups services

### Fixed
- `Project.CustomDomainName` accepts both string and string-array JSON (project/fetch splits domains into an array).

### Changed
- `CreateProjectRequest` now matches control-plane `POST /project/create` (clusterUUID, environment_uuid, buildSettings, envVariables, networkSettings, workspace_uuid, …). Legacy `server_id` / `environment_id` / `build_command` fields are removed.
- `Project.CustomDomainName` type is `FlexibleCSVString` (string-compatible via `.String()` / `.First()`).
- GitHub Actions CI workflow for automated testing and linting
- GitHub Actions release workflow for automated releases
- GoReleaser configuration for release management
- Makefile for common development tasks
- golangci-lint configuration for code quality
- Release process documentation

### Changed
- `BackupService` methods now return a clear deprecation error; prefer addon backup export and volume export APIs

## [0.1.0] - Initial Release

### Added
- Complete PipeOps API coverage with 284 methods across 18 service modules
- OAuth 2.0 authorization code flow support
- Context support for all API methods
- Type-safe request/response structures
- Comprehensive examples and documentation
- Projects API (46 methods)
- Billing API (33 methods)
- Servers/Clusters API (22 methods)
- Cloud Providers API (17 methods)
- Teams API (11 methods)
- Admin API (20 methods)
- Add-Ons API (21 methods)
- Service Tokens API (5 methods)
- Environments API (8 methods)
- Authentication API (10 methods)
- Webhooks API (8 methods)
- Events & Survey API (23 methods)
- DeploymentWebhooks API (3 methods)
- Campaign API (7 methods)
- OAuth API (4 methods)
- OpenCost API (3 methods)
- Coupons API (2 methods)
- Workspaces API (6 methods)
- User Settings API (8 methods)

[Unreleased]: https://github.com/PipeOpsHQ/pipeops-go-sdk/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/PipeOpsHQ/pipeops-go-sdk/releases/tag/v0.1.0
