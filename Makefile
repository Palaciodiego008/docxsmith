.PHONY: build test clean install run-example help

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

# Show help
help:
	@echo "DocxSmith - Makefile commands:"
	@echo ""
	@echo "  make build           - Build the CLI tool"
	@echo "  make install         - Install the CLI tool"
	@echo "  make test            - Run tests"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo "  make run-example     - Run the basic usage example"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make lint            - Run linter"
	@echo "  make fmt             - Format code"
	@echo "  make tidy            - Tidy dependencies"
	@echo "  make create-test-doc - Create a sample test document"
	@echo "  make pre-commit      - Run all checks before commit"
	@echo "  make help            - Show this help message"
