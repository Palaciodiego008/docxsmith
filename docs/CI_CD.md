# CI/CD Pipeline - DocxSmith

DocxSmith includes a professional CI/CD setup with GitHub Actions, pre-commit hooks, and automated testing.

## GitHub Actions Workflows

### 1. CI Workflow (`.github/workflows/ci.yml`)

Runs on every push and pull request to `main` and `develop` branches.

**Jobs:**

#### Test Job
- **Matrix testing**: Tests across 3 OS (Ubuntu, macOS, Windows) and 3 Go versions (1.21, 1.22, 1.23)
- **Race detection**: Runs with `-race` flag to detect race conditions
- **Coverage**: Generates coverage report on Ubuntu + Go 1.23
- **Codecov integration**: Uploads coverage to Codecov.io

#### Build Job
- **Builds CLI**: Compiles the binary
- **Tests binary**: Runs basic commands to verify
- **Artifact upload**: Saves binary as GitHub artifact

#### Lint Job
- **go fmt check**: Ensures code is formatted
- **go vet**: Static analysis
- **staticcheck**: Advanced linting

#### Security Job
- **Gosec scanner**: Security vulnerability scanning
- **SARIF upload**: Results uploaded to GitHub Security tab

### 2. Release Workflow (`.github/workflows/release.yml`)

Triggers on git tags (`v*.*.*`).

**Automated Release Process:**
1. Runs all tests
2. Builds binaries for 5 platforms:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64 M1/M2)
   - Windows (amd64)
3. Creates SHA256 checksums
4. Creates GitHub release with:
   - All binaries attached
   - Auto-generated release notes
   - Changelog

**Usage:**
```bash
git tag v1.1.0
git push origin v1.1.0
# GitHub Actions automatically creates release
```

### 3. CodeQL Workflow (`.github/workflows/codeql.yml`)

Runs weekly and on main branch changes.

**Security Analysis:**
- Advanced code scanning
- Detects security vulnerabilities
- Results in GitHub Security tab
- Automated alerts for issues

## Dependabot

Auto-updates dependencies weekly (`.github/dependabot.yml`):
- Go modules
- GitHub Actions

**Automatic PRs created for:**
- Dependency updates
- Security patches
- Version upgrades

## Local Development

### Pre-commit Hooks

Install git hooks for local validation:

```bash
make setup-hooks
```

**What it runs before each commit:**
1. `gofmt` check
2. `go vet`
3. Tests (fast mode)
4. Build verification

### Manual CI Simulation

Run the same checks as GitHub Actions locally:

```bash
# Full CI pipeline
make ci

# Individual checks
make fmt          # Format code
make test-race    # Tests with race detector
make coverage-ci  # Coverage report
```

## Makefile Commands

### Testing Commands

```bash
make test           # Run all tests
make test-race      # Run with race detector (like CI)
make test-coverage  # Generate HTML coverage report
make coverage-ci    # Coverage in CI format (atomic mode)
```

### CI/CD Commands

```bash
make ci             # Simulate full GitHub Actions CI
make setup-hooks    # Install pre-commit hooks
make pre-commit     # Run pre-commit checks manually
```

### Code Quality

```bash
make fmt            # Format all code
make lint           # Run linter (requires golangci-lint)
make tidy           # Clean up go.mod
```

## Continuous Integration Flow

```
┌─────────────────┐
│  Developer      │
│  commits code   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Pre-commit     │◄─── Local validation
│  Hook runs      │     (fmt, vet, tests)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Push to        │
│  GitHub         │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  GitHub Actions │
│  CI Workflow    │
└────────┬────────┘
         │
         ├──► Test (3 OS × 3 Go versions = 9 jobs)
         ├──► Build (compile & verify)
         ├──► Lint (fmt, vet, staticcheck)
         └──► Security (gosec, codeql)
         │
         ▼
┌─────────────────┐
│  All checks     │
│  pass ✅        │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Ready to       │
│  merge/deploy   │
└─────────────────┘
```

## Release Flow

```
┌─────────────────┐
│  Tag version    │
│  git tag v1.0.0 │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Push tag       │
│  to GitHub      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Release        │
│  Workflow       │
└────────┬────────┘
         │
         ├──► Run all tests
         ├──► Build for 5 platforms
         ├──► Create checksums
         └──► Generate release notes
         │
         ▼
┌─────────────────┐
│  GitHub Release │
│  Published ✅   │
│  with binaries  │
└─────────────────┘
```

## Best Practices

### Before Committing

```bash
# 1. Format code
make fmt

# 2. Run tests
make test

# 3. Run full pre-commit checks
make pre-commit

# 4. (Optional) Simulate CI
make ci
```

### Before Creating PR

```bash
# Run full CI suite locally
make ci

# Check coverage
make coverage-ci
```

### Creating a Release

```bash
# 1. Update version in code
# 2. Update CHANGELOG.md
# 3. Commit changes
git commit -am "chore: bump version to v1.1.0"

# 4. Create and push tag
git tag v1.1.0
git push origin v1.1.0

# 5. GitHub Actions automatically creates release
```

## Monitoring

### GitHub Actions Dashboard

View CI/CD status:
- Go to repository → Actions tab
- See all workflow runs
- Check build artifacts
- Download binaries

### Code Coverage

View coverage reports:
- Codecov.io dashboard (if configured)
- Local: `make test-coverage` → open `coverage.html`
- CI output: Check Actions logs

### Security Alerts

Monitor security:
- Repository → Security tab
- CodeQL alerts
- Dependabot alerts
- Gosec findings

## Troubleshooting

### CI Failing Locally Passes

**Solution:**
```bash
# Run exact CI commands
make ci

# Check specific issues
go test -race ./...  # Race conditions
gofmt -s -l .       # Formatting
go vet ./...        # Static analysis
```

### Pre-commit Hook Not Running

**Solution:**
```bash
# Reinstall hooks
make setup-hooks

# Verify installation
ls -la .git/hooks/pre-commit

# Test manually
.git/hooks/pre-commit
```

### Build Failing on Specific OS

**Solution:**
- Check GitHub Actions matrix results
- Identify failing OS/Go version
- Test locally with matching version
- Fix platform-specific issues

## Configuration Files

| File | Purpose |
|------|---------|
| `.github/workflows/ci.yml` | Main CI pipeline |
| `.github/workflows/release.yml` | Release automation |
| `.github/workflows/codeql.yml` | Security scanning |
| `.github/dependabot.yml` | Dependency updates |
| `.githooks/pre-commit` | Local git hook |
| `Makefile` | Build automation |

## Integration with External Services

### Codecov (Coverage)

Add to repository settings:
```yaml
# codecov.yml
coverage:
  status:
    project:
      default:
        target: 70%
    patch:
      default:
        target: 80%
```

### Badge in README

```markdown
![CI](https://github.com/Palaciodiego008/docxsmith/workflows/CI/badge.svg)
![Coverage](https://codecov.io/gh/Palaciodiego008/docxsmith/branch/main/graph/badge.svg)
![Go Report](https://goreportcard.com/badge/github.com/Palaciodiego008/docxsmith)
```

## Local Setup for Contributors

```bash
# 1. Clone repository
git clone https://github.com/Palaciodiego008/docxsmith.git
cd docxsmith

# 2. Install dependencies
go mod download

# 3. Setup git hooks
make setup-hooks

# 4. Run tests
make test

# 5. Build
make build
```

## Automated Checks Summary

**Every Commit (Local):**
- ✓ Code formatting
- ✓ Go vet
- ✓ Unit tests
- ✓ Build verification

**Every Push (GitHub Actions):**
- ✓ Matrix testing (9 combinations)
- ✓ Race detection
- ✓ Coverage reporting
- ✓ Lint checks
- ✓ Security scanning
- ✓ Build artifacts

**Weekly (GitHub Actions):**
- ✓ CodeQL security scan
- ✓ Dependency updates (Dependabot)

**On Release (Git Tag):**
- ✓ Full test suite
- ✓ Multi-platform builds
- ✓ Checksum generation
- ✓ Automatic release notes
- ✓ Binary uploads

## Resources

- [GitHub Actions Docs](https://docs.github.com/actions)
- [Go Testing](https://golang.org/pkg/testing/)
- [Codecov](https://codecov.io/)
