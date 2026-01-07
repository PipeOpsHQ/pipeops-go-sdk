# OpenCode Guidelines

## Build, Test, and Lint Commands
1. **Build the project:**
   ```bash
   make build
   ```
2. **Run all tests:**
   ```bash
   make test
   ```
3. **Run a single test file:**
   ```bash
   go test ./path/to/test_file.go
   ```
4. **Lint the codebase:**
   ```bash
   make lint
   ```
5. **Check formatting:**
   ```bash
   make fmt
   ```

## Code Style Guidelines
1. **Imports:**
   - Standard library imports first, followed by third-party, and then local imports.
   - Use `goimports` for sorting and grouping imports.
2. **Formatting:**
   - Use `go fmt` to ensure consistent formatting.
   - Avoid aligning struct fields or variable declarations; use idiomatic Go style.
3. **Types and Naming:**
   - Follow `CamelCase` for exported identifiers and `camelCase` for internal variables.
   - Use descriptive names; avoid abbreviations.
4. **Error Handling:**
   - Check and handle all errors.
   - Return errors rather than panicking, except in truly exceptional circumstances.
   - Wrap errors with `fmt.Errorf` when additional context is helpful.
5. **Testing:**
   - Place tests in `*_test.go` files.
   - Use `t.Run` for subtests and table-driven approaches where applicable.
   - Strive for high test coverage.
6. **Go Version Compatibility:**
   - Ensure compatibility with Go 1.21, 1.22, and 1.23 as defined in CI workflows.
7. **CI Workflows:**
   - All changes must pass GitHub Actions workflows, including tests and linting.

## Contribution Workflow
1. Open a feature branch for changes.
2. Add/modify tests relevant to your changes.
3. Ensure `make test`, `make lint`, and `make fmt` pass.
4. Follow the commit message convention: `<type>: <description>`.
5. Submit a pull request with a clear description.
