# Contributing with CI/CD

Guide for contributors on using the CI/CD pipeline.

## Quick Setup

```bash
# 1. Clone and setup
git clone https://github.com/Palaciodiego008/docxsmith.git
cd docxsmith

# 2. Install dependencies
go mod download

# 3. Setup git hooks (automatically run tests before commit)
make setup-hooks

# 4. Run tests
make test
```

## Development Workflow

### 1. Before You Start

```bash
# Ensure you have latest code
git pull origin main

# Create feature branch
git checkout -b feature/my-new-feature

# Verify everything works
make test
```

### 2. While Developing

```bash
# Format code frequently
make fmt

# Run tests after changes
make test

# Check specific package
go test ./pkg/docx -v
```

### 3. Before Committing

```bash
# Run full pre-commit checks
make pre-commit

# Or simulate full CI
make ci
```

The pre-commit hook will automatically run when you commit (if installed with `make setup-hooks`).

### 4. Push and Create PR

```bash
# Push changes
git push origin feature/my-new-feature

# Create PR on GitHub
# CI will automatically run all checks
```

## CI Pipeline Details

### What Runs on Every Push/PR

**Test Matrix (9 combinations):**
- OS: Ubuntu, macOS, Windows
- Go: 1.21, 1.22, 1.23

**Checks:**
- ✓ Unit tests with race detector
- ✓ Code formatting (gofmt)
- ✓ Static analysis (go vet)
- ✓ Linting (staticcheck)
- ✓ Security scan (gosec)
- ✓ Build verification
- ✓ Coverage reporting

**Expected Duration:** ~3-5 minutes

### Viewing CI Results

1. Go to your PR on GitHub
2. Scroll to bottom - see check status
3. Click "Details" to view logs
4. Fix any failures and push again

## Common CI Failures

### Format Check Failed

**Error:** "Go code is not formatted"

**Fix:**
```bash
make fmt
git add .
git commit --amend --no-edit
git push --force-with-lease
```

### Test Failed

**Error:** Test failures in CI

**Fix:**
```bash
# Run tests locally with same flags as CI
make test-race

# Fix the issue
# Run tests again
make test

# Commit fix
git commit -am "fix: resolve test failure"
git push
```

### Build Failed

**Error:** Build errors

**Fix:**
```bash
# Test build locally
go build ./cmd/docxsmith

# Check for missing imports
go mod tidy

# Fix and commit
git commit -am "fix: build errors"
git push
```

## Release Process

### Creating a Release

Only maintainers can create releases.

```bash
# 1. Update version
# Edit internal/cli/cli.go - update Version constant
# Edit CHANGELOG.md - move [Unreleased] to new version

# 2. Commit version bump
git commit -am "chore: bump version to v1.2.0"
git push origin main

# 3. Create and push tag
git tag v1.2.0
git push origin v1.2.0

# 4. GitHub Actions automatically:
#    - Runs all tests
#    - Builds binaries for 5 platforms
#    - Creates GitHub release
#    - Attaches binaries
#    - Generates release notes
```

### Release Artifacts

Automatic builds for:
- `docxsmith-linux-amd64`
- `docxsmith-linux-arm64`
- `docxsmith-darwin-amd64` (Intel Mac)
- `docxsmith-darwin-arm64` (M1/M2 Mac)
- `docxsmith-windows-amd64.exe`
- `checksums.txt` (SHA256)

## Local Testing

### Simulate CI Locally

```bash
# Run exact same checks as GitHub Actions
make ci
```

This runs:
1. Code formatting check
2. go vet
3. Tests with race detector
4. Build verification

### Test Specific Platforms

```bash
# Test on Linux
GOOS=linux GOARCH=amd64 go build ./cmd/docxsmith

# Test on macOS
GOOS=darwin GOARCH=amd64 go build ./cmd/docxsmith

# Test on Windows
GOOS=windows GOARCH=amd64 go build ./cmd/docxsmith
```

### Coverage Analysis

```bash
# Generate HTML coverage report
make test-coverage
open coverage.html

# Or CI-style coverage
make coverage-ci
```

## Pre-commit Hooks

### Installation

```bash
# One-time setup
make setup-hooks
```

### What It Does

Before every commit, automatically runs:
1. Format check (gofmt)
2. Static analysis (go vet)
3. Fast tests (with `-short` flag)
4. Build check

### Bypass Hook (Emergency Only)

```bash
# Skip pre-commit hook (not recommended)
git commit --no-verify -m "emergency fix"
```

## Code Quality Standards

### Required Before Merge

- ✅ All tests passing
- ✅ No race conditions
- ✅ Code formatted (gofmt)
- ✅ No vet warnings
- ✅ No lint warnings
- ✅ Build successful on all platforms
- ✅ Coverage maintained or improved
- ✅ No security issues

### Code Review Checklist

- [ ] Tests added for new features
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Code follows Go best practices
- [ ] No code duplication
- [ ] Error handling implemented
- [ ] Edge cases covered

## Troubleshooting

### Pre-commit Hook Not Working

```bash
# Verify installation
ls -la .git/hooks/pre-commit

# Reinstall
make setup-hooks

# Test manually
.git/hooks/pre-commit
```

### CI Passes Locally But Fails on GitHub

**Possible causes:**
- Different Go version
- Different OS
- Race conditions (only detected with `-race`)

**Solution:**
```bash
# Test with race detector locally
make test-race

# Test on specific Go version
go test -race ./...
```

### Dependabot PRs

Dependabot automatically creates PRs for dependency updates.

**What to do:**
1. Review the changes in the PR
2. CI will run automatically
3. If tests pass, merge
4. If tests fail, investigate and fix

## Best Practices

1. **Always run** `make ci` before pushing
2. **Install hooks** with `make setup-hooks`
3. **Write tests** for all new features
4. **Update docs** when adding features
5. **Keep commits atomic** - one feature per commit
6. **Write clear commit messages**
7. **Rebase** instead of merge when possible
8. **Run tests** after rebasing

## Continuous Improvement

### Adding New Tests

```go
// Use table-driven tests
func TestMyFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case 1", "input1", "output1"},
        {"case 2", "input2", "output2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Improving Coverage

```bash
# Find uncovered code
go test -coverprofile=coverage.out ./pkg/docx
go tool cover -html=coverage.out

# Add tests for uncovered lines
```

## Resources

- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [Go Testing Best Practices](https://golang.org/doc/tutorial/add-a-test)
- [Effective Go](https://golang.org/doc/effective_go)
- [CI/CD Documentation](./CI_CD.md)
