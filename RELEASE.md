# Release Process

This document describes the release process for the PipeOps Go SDK.

## Overview

The PipeOps Go SDK uses semantic versioning and automated releases through GitHub Actions and GoReleaser.

## Semantic Versioning

We follow [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality in a backwards compatible manner
- **PATCH** version for backwards compatible bug fixes

## Release Workflow

### 1. Prepare for Release

Before creating a release, ensure:

1. All changes are merged to the `main` branch
2. All CI checks pass
3. Documentation is up to date
4. `CHANGELOG.md` is updated (if maintained manually)

### 2. Auto-generated Release Tags

Commits merged into the `main` branch automatically trigger the **Auto Tag Release** workflow. It creates the next semantic version tag and pushes it to the repository (default bump is `patch`, `feat` commits bump `minor`). Use `[skip release]` or `[skip tag]` in the commit subject to short-circuit auto-tagging when needed.

Manual tagging still works when you need a specific version:

```bash
# For a new minor version
git tag v1.1.0

# For a new patch version
git tag v1.0.1

# Push the tag to trigger the release workflow
git push origin v1.1.0
```

### 3. Automated Release Process

When you push a tag matching `v*.*.*`, the following happens automatically:

1. **CI Validation**: The release workflow runs all tests
2. **GoReleaser**: Creates a GitHub release with:
   - Auto-generated changelog
   - Source code archives
   - Release notes
3. **GitHub Release**: A new release is published on GitHub

### 4. Verify the Release

After the release workflow completes:

1. Check the [Releases page](https://github.com/PipeOpsHQ/pipeops-go-sdk/releases)
2. Verify the changelog and release notes
3. Test installation: `go get github.com/PipeOpsHQ/pipeops-go-sdk@v1.1.0`

## Pre-release Versions

To create a pre-release version (alpha, beta, rc):

```bash
# Alpha release
git tag v1.1.0-alpha.1
git push origin v1.1.0-alpha.1

# Beta release
git tag v1.1.0-beta.1
git push origin v1.1.0-beta.1

# Release candidate
git tag v1.1.0-rc.1
git push origin v1.1.0-rc.1
```

Pre-releases are automatically marked as "pre-release" on GitHub.

## Testing the Release Process

To test the release process without actually publishing:

```bash
# Test GoReleaser configuration
goreleaser check

# Create a snapshot build (local only)
goreleaser release --snapshot --clean --skip=publish

# Or use the Makefile
make release-snapshot
```

## Rollback a Release

If you need to rollback a release:

1. Delete the tag locally and remotely:

   ```bash
   git tag -d v1.1.0
   git push --delete origin v1.1.0
   ```

2. Delete the GitHub release from the [Releases page](https://github.com/PipeOpsHQ/pipeops-go-sdk/releases)

3. Create a new patch version with the fix

## Go Module Proxy

After publishing a release, the Go module proxy automatically indexes it. Users can install specific versions:

```bash
# Latest version
go get github.com/PipeOpsHQ/pipeops-go-sdk

# Specific version
go get github.com/PipeOpsHQ/pipeops-go-sdk@v1.1.0

# Latest patch for a minor version
go get github.com/PipeOpsHQ/pipeops-go-sdk@v1.1
```

## Release Checklist

- [ ] All tests pass
- [ ] Code is properly formatted (`go fmt ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] Documentation is updated
- [ ] Examples work with the new version
- [ ] Version tag follows semantic versioning
- [ ] Tag is pushed to GitHub
- [ ] Release workflow completes successfully
- [ ] GitHub release is created
- [ ] Release can be installed via `go get`

## Troubleshooting

### Release workflow fails

1. Check the Actions tab for error details
2. Verify the tag format is correct (`v*.*.*`)
3. Ensure all tests pass on the main branch
4. Check GoReleaser configuration with `goreleaser check`

### Module not found after release

The Go module proxy may take a few minutes to index a new release. If users can't install immediately:

1. Wait 5-10 minutes for proxy indexing
2. Try with a proxy refresh: `GOPROXY=https://proxy.golang.org go get github.com/PipeOpsHQ/pipeops-go-sdk@v1.1.0`
3. Clear Go cache: `go clean -modcache`

## Continuous Integration

All pull requests and commits to `main` automatically run:

- Tests across multiple Go versions (1.21, 1.22, 1.23)
- Code formatting checks
- Linting with golangci-lint
- Code coverage reporting

This ensures that the `main` branch is always in a releasable state.
