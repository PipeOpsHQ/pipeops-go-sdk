# GitHub Actions Workflows

This document describes the GitHub Actions workflows configured for this repository.

## Overview

The repository uses GitHub Actions for continuous integration, automated testing, and release management.

## Workflows

### 1. CI Workflow (`.github/workflows/ci.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

**Jobs:**

#### Test Job
- **Matrix Testing**: Runs tests on Go versions 1.21, 1.22, and 1.23
- **Steps:**
  1. Checkout code
  2. Set up Go environment
  3. Cache Go modules for faster builds
  4. Download and verify dependencies
  5. Build the project
  6. Run `go vet` for static analysis
  7. Check code formatting with `go fmt`
  8. Run tests with race detection and coverage
  9. Upload coverage to Codecov (Go 1.21 only)

#### Lint Job
- **Steps:**
  1. Checkout code
  2. Set up Go 1.21 environment
  3. Run golangci-lint with configured rules

**Purpose:** Ensures code quality and compatibility across multiple Go versions.

---

### 2. Release Workflow (`.github/workflows/release.yml`)

**Triggers:**
- Push of version tags matching pattern `v*.*.*` (e.g., `v1.0.0`, `v2.1.3`)

**Jobs:**

#### Release Job
- **Steps:**
  1. Checkout code with full git history
  2. Set up Go 1.21 environment
  3. Cache Go modules
  4. Download dependencies
  5. Run tests to ensure quality
  6. Run GoReleaser to create the release
  7. Upload release artifacts

**Outputs:**
- GitHub release with auto-generated changelog
- Source code archives (`.tar.gz`)
- Checksums file
- Release notes

**Purpose:** Automates the release process when version tags are pushed.

---

## GoReleaser Configuration (`.goreleaser.yml`)

GoReleaser is configured to:
- Skip binary builds (this is a library, not a standalone application)
- Generate organized changelogs grouped by type (features, fixes, etc.)
- Create source archives
- Calculate checksums
- Auto-detect pre-release versions (alpha, beta, rc)

### Changelog Groups
- **Features**: Commits starting with `feat:`
- **Bug Fixes**: Commits starting with `fix:`
- **Performance Improvements**: Commits starting with `perf:`
- **Refactors**: Commits starting with `refactor:`
- **Documentation**: Commits starting with `docs:`

---

## golangci-lint Configuration (`.golangci.yml`)

Enabled linters:
- `errcheck` - Checks for unchecked errors
- `gosimple` - Simplifies code
- `govet` - Reports suspicious constructs
- `ineffassign` - Detects ineffectual assignments
- `staticcheck` - Advanced static analysis
- `unused` - Checks for unused code
- `gofmt` - Checks code formatting
- `goimports` - Checks import formatting
- `misspell` - Finds commonly misspelled words
- `revive` - Fast, configurable linter
- `unconvert` - Removes unnecessary type conversions
- `unparam` - Reports unused function parameters
- `goconst` - Finds repeated strings
- `gocyclo` - Computes cyclomatic complexity
- `stylecheck` - Replacement for golint

---

## Dependabot Configuration (`.github/dependabot.yml`)

**Updates:**

### Go Modules
- **Frequency**: Weekly
- **Directory**: Root (`/`)
- **Open PRs limit**: 10
- **Labels**: `dependencies`, `go`

### GitHub Actions
- **Frequency**: Weekly
- **Directory**: Root (`/`)
- **Open PRs limit**: 5
- **Labels**: `dependencies`, `github-actions`

**Purpose:** Keeps dependencies and GitHub Actions up to date automatically.

---

## How to Use

### Running CI Locally

Before pushing code, you can run the same checks locally:

```bash
# Run all checks
make check

# Or run individually
make fmt      # Format code
make vet      # Run go vet
make lint     # Run golangci-lint
make test     # Run tests
```

### Creating a Release

1. Ensure all changes are merged to `main`
2. Create and push a version tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. The release workflow will automatically:
   - Run tests
   - Create a GitHub release
   - Generate changelog
   - Upload artifacts

See the release documentation for detailed release instructions.

### Testing Release Process

Test the release process without publishing:

```bash
# Check GoReleaser configuration
goreleaser check

# Create a snapshot release (local only)
make release-snapshot

# Test release without publishing
make release-test
```

---

## Badges

The following badges are available in the README:

- **CI Status**: Shows the status of the CI workflow
- **Go Report Card**: Code quality grade from goreportcard.com
- **GoDoc**: Documentation link
- **License**: Repository license information

---

## Troubleshooting

### CI Workflow Fails

1. Check the Actions tab for detailed logs
2. Run the same checks locally with `make check`
3. Ensure all dependencies are up to date with `go mod tidy`

### Linter Errors

1. Run `make lint` locally to see all errors
2. Use `make lint-fix` to auto-fix some issues
3. Check `.golangci.yml` for linter configuration

### Release Workflow Fails

1. Verify tag format is `v*.*.*` (e.g., `v1.0.0`)
2. Ensure all tests pass on `main` branch
3. Check GoReleaser config with `goreleaser check`

---

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GoReleaser Documentation](https://goreleaser.com/)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [Dependabot Documentation](https://docs.github.com/en/code-security/dependabot)
