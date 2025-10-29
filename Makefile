.PHONY: help
help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.PHONY: deps
deps: ## Download and verify dependencies
	go mod download
	go mod verify

.PHONY: tidy
tidy: ## Tidy and vendor dependencies
	go mod tidy
	go mod verify

.PHONY: build
build: ## Build the project
	go build -v ./...

.PHONY: test
test: ## Run tests
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: test-short
test-short: ## Run short tests
	go test -v -short ./...

.PHONY: coverage
coverage: test ## Generate coverage report
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run --timeout=5m

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --fix --timeout=5m

.PHONY: check
check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)

.PHONY: clean
clean: ## Clean build artifacts
	rm -f coverage.txt coverage.html
	go clean -cache -testcache -modcache

.PHONY: install-tools
install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: release-snapshot
release-snapshot: ## Create a snapshot release (for testing)
	goreleaser release --snapshot --clean --skip=publish

.PHONY: release-test
release-test: ## Test release process without publishing
	goreleaser check
	goreleaser release --skip=publish --clean

.DEFAULT_GOAL := help
