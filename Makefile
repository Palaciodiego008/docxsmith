.PHONY: build test clean install run-example help ci setup-hooks

# Build the CLI tool
build:
	@echo "Building DocxSmith CLI..."
	go build -o bin/docxsmith ./cmd/docxsmith

# Install the CLI tool
install:
	@echo "Installing DocxSmith CLI..."
	go install ./cmd/docxsmith

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run example
run-example:
	@echo "Running basic usage example..."
	cd examples && go run basic_usage.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f examples/*.docx

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Create a test document
create-test-doc:
	@echo "Creating test document..."
	@mkdir -p testdata
	go run cmd/docxsmith/main.go create -output testdata/sample.docx -text "This is a sample document for testing."
	@echo "Test document created at testdata/sample.docx"

# Run all checks before commit
pre-commit: fmt tidy test
	@echo "Pre-commit checks completed successfully!"

# Setup git hooks
setup-hooks:
	@echo "Setting up git hooks..."
	@mkdir -p .git/hooks
	@cp .githooks/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed successfully!"
	@echo "Pre-commit hook will run: fmt, vet, and tests"

# CI simulation (runs same checks as GitHub Actions)
ci: fmt tidy
	@echo "Running CI checks..."
	@echo "  ✓ Format check..."
	@test -z "$$(gofmt -s -l .)" || (echo "Code not formatted" && exit 1)
	@echo "  ✓ Go vet..."
	@go vet ./...
	@echo "  ✓ Tests with race detector..."
	@go test -race -timeout 5m ./...
	@echo "  ✓ Build check..."
	@go build -o /tmp/docxsmith ./cmd/docxsmith
	@rm -f /tmp/docxsmith
	@echo "✅ All CI checks passed!"

# Test with race detector (like GitHub Actions)
test-race:
	@echo "Running tests with race detector..."
	go test -race -v ./...

# Coverage with upload simulation
coverage-ci:
	@echo "Running coverage like CI..."
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out
	@echo "Coverage report: coverage.out"

# Show help
help:
	@echo "DocxSmith - Makefile commands:"
	@echo ""
	@echo "Build & Install:"
	@echo "  make build           - Build the CLI tool"
	@echo "  make install         - Install the CLI tool globally"
	@echo ""
	@echo "Testing:"
	@echo "  make test            - Run tests"
	@echo "  make test-race       - Run tests with race detector"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo "  make coverage-ci     - Coverage like GitHub Actions"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt             - Format code"
	@echo "  make lint            - Run linter"
	@echo "  make tidy            - Tidy dependencies"
	@echo "  make pre-commit      - Run all checks before commit"
	@echo ""
	@echo "CI/CD:"
	@echo "  make ci              - Simulate GitHub Actions CI"
	@echo "  make setup-hooks     - Install git pre-commit hooks"
	@echo ""
	@echo "Utilities:"
	@echo "  make run-example     - Run the basic usage example"
	@echo "  make create-test-doc - Create a sample test document"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make help            - Show this help message"
