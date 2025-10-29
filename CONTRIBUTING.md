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
```

3. Run tests:
```bash
go test ./...
```

4. Format code:
```bash
go fmt ./...
```

5. Run linter:
```bash
go vet ./...
```

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
5. Ensure all tests pass
6. Format your code
7. Submit a pull request

## Questions?

If you have questions about contributing, please open an issue or reach out to the maintainers.
