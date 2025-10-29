# Contributing to PipeOps Go SDK

Thank you for your interest in contributing to the PipeOps Go SDK!

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/PipeOpsHQ/pipeops-go-sdk.git
cd pipeops-go-sdk
```

2. Install dependencies:
```bash
go mod download
# or use make
make deps
```

3. Run tests:
```bash
go test ./...
# or use make
make test
```

4. Format code:
```bash
go fmt ./...
# or use make
make fmt
```

5. Run linter:
```bash
go vet ./...
# or use make
make vet
```

6. Run golangci-lint (optional but recommended):
```bash
# Install golangci-lint if not already installed
make install-tools

# Run linter
make lint
```

## Available Make Commands

Run `make help` to see all available commands:

- `make deps` - Download and verify dependencies
- `make tidy` - Tidy and vendor dependencies
- `make build` - Build the project
- `make test` - Run tests with coverage
- `make test-short` - Run short tests
- `make coverage` - Generate HTML coverage report
- `make fmt` - Format code
- `make vet` - Run go vet
- `make lint` - Run golangci-lint
- `make lint-fix` - Run golangci-lint with auto-fix
- `make check` - Run all checks (fmt, vet, lint, test)
- `make clean` - Clean build artifacts
- `make install-tools` - Install development tools
- `make release-snapshot` - Create a snapshot release (for testing)
- `make release-test` - Test release process without publishing

## Adding New API Endpoints

1. Identify the appropriate service file in the `pipeops/` directory (e.g., `auth.go`, `projects.go`).

2. Add the request and response types:
```go
type MyNewRequest struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
}

type MyNewResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Data    struct {
        Result string `json:"result"`
    } `json:"data"`
}
```

3. Add the method to the service:
```go
func (s *MyService) MyNewMethod(ctx context.Context, req *MyNewRequest) (*MyNewResponse, *http.Response, error) {
    u := "my/endpoint"
    
    httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
    if err != nil {
        return nil, nil, err
    }
    
    resp := new(MyNewResponse)
    httpResp, err := s.client.Do(ctx, httpReq, resp)
    if err != nil {
        return nil, httpResp, err
    }
    
    return resp, httpResp, nil
}
```

4. Add tests for your new method in a `*_test.go` file.

5. Update documentation as needed.

## Code Style

- Follow standard Go conventions
- Use `gofmt` to format your code
- Run `go vet` to check for common errors
- Add comments for exported types and functions
- Keep functions focused and concise

## Pull Request Process

1. Fork the repository
2. Create a new branch for your feature or fix
3. Make your changes
4. Add tests if applicable
5. Ensure all tests pass with `make test`
6. Format your code with `make fmt`
7. Run linter with `make lint`
8. Submit a pull request

Pull requests will automatically trigger:
- CI workflow that runs tests across multiple Go versions
- Code formatting checks
- Linting with golangci-lint
- Coverage reporting

## Continuous Integration

This repository uses GitHub Actions for continuous integration. The following workflows are configured:

### CI Workflow
Runs on every push and pull request to `main` and `develop` branches:
- Tests on Go 1.21, 1.22, and 1.23
- Code formatting checks (`go fmt`)
- Static analysis (`go vet`)
- Linting with golangci-lint
- Coverage reporting to Codecov

### Release Workflow
Automatically triggered when a new version tag is pushed:
- Runs tests
- Creates a GitHub release with GoReleaser
- Generates changelog
- Creates source archives
- Updates release notes

See [RELEASE.md](RELEASE.md) for information on creating releases.

## Code Quality

We use several tools to maintain code quality:

- **golangci-lint**: Comprehensive Go linter with multiple checkers enabled
- **go vet**: Go's built-in static analyzer
- **gofmt**: Go's official code formatter

Configuration files:
- `.golangci.yml` - golangci-lint configuration
- `.goreleaser.yml` - GoReleaser configuration

## Dependency Management

This repository uses Dependabot to keep dependencies up to date:
- Go module dependencies are checked weekly
- GitHub Actions versions are checked weekly
- PRs are automatically created for updates

## Questions?

If you have questions about contributing, please open an issue or reach out to the maintainers.
