# CI/CD Pipeline Documentation

This document describes the automated CI/CD pipeline for gosh.

## Repository

- **GitHub**: https://github.com/john-eevee/gosh
- **Docker Registry**: ghcr.io/john-eevee/gosh

## GitHub Actions Workflows

### 1. Continuous Integration (`.github/workflows/ci.yml`)

Runs on every push to `main` and `develop` branches, and on pull requests.

#### Jobs:

**Test** - Multi-platform testing
- Runs on: Ubuntu, macOS, Windows
- Go versions: 1.21, 1.22
- Steps:
  - Run unit tests with race detector
  - Run benchmarks
  - Build binary

**Lint** - Code quality checks
- Runs on: Ubuntu
- Steps:
  - Run golangci-lint
  - Check code formatting with gofmt

**Coverage** - Code coverage analysis
- Runs on: Ubuntu
- Steps:
  - Generate coverage report
  - Upload to Codecov
  - Display coverage summary

**Security** - Security scanning
- Runs on: Ubuntu
- Steps:
  - Run Gosec security scanner
  - Upload SARIF results to GitHub Security tab

#### Triggers:
```yaml
- Push to main or develop
- Pull request to main or develop
- Changes to .go files, go.mod, go.sum, or workflow files
```

### 2. Release Workflow (`.github/workflows/release.yml`)

Triggered when a tag matching `v*` is pushed.

#### Jobs:

**Build** - Cross-platform binary builds
- Builds for:
  - Linux: amd64, arm64
  - macOS: amd64, arm64
  - Windows: amd64
- Steps:
  - Build optimized binaries
  - Compress with tar.gz (Unix) or zip (Windows)
  - Upload as artifacts

**Release** - GitHub Release creation
- Steps:
  - Download all build artifacts
  - Generate release notes from CHANGELOG.md
  - Create GitHub Release with artifacts
  - Make available for download

**Docker** - Docker image building
- Steps:
  - Build Docker image with multi-stage build
  - Push to GitHub Container Registry (ghcr.io)
  - Tag with:
    - Full version: `v0.1.1`
    - Major.minor: `0.1`
    - Git SHA: `abc1234`

#### Triggers:
```yaml
- Tag push: git tag v0.1.1 && git push origin v0.1.1
```

### 3. CodeQL Analysis (`.github/workflows/codeql.yml`)

Security code scanning with GitHub's CodeQL.

#### Triggers:
- Every push to main/develop
- Weekly schedule (Sunday 4 AM UTC)
- All pull requests

### 4. Dependabot (`.github/dependabot.yml`)

Automated dependency updates.

#### Configuration:
- **Go Modules**: Weekly updates (Monday 3 AM UTC)
- **GitHub Actions**: Weekly updates (Monday 4 AM UTC)
- **Pull Request Limit**: 5 open PRs at a time
- **Auto-merge**: Not enabled (manual review required)

## Local Development

### Building Locally

```bash
# Install Go 1.21+
go version

# Clone repository
git clone https://github.com/john-eevee/gosh.git
cd gosh

# Download dependencies
make deps

# Run all checks (format, lint, test, build)
make all

# Or individually:
make fmt              # Format code with gofmt
make lint             # Run linters
make test             # Run tests
make coverage         # Generate coverage report
make coverage-html    # Generate HTML coverage
make build            # Build binary
make clean            # Remove artifacts
```

### Pre-commit Checklist

Before pushing, ensure:

```bash
# Format code
make fmt

# Run linters
make lint

# Run tests
make test

# Check coverage
make coverage

# Build successfully
make build
```

## Release Process

### Creating a Release

1. **Update Version**
   - Edit `pkg/version.go` (if version constant exists)
   - Update `CHANGELOG.md` with release notes

2. **Create Tag**
   ```bash
   git tag -a v0.2.0 -m "Release v0.2.0"
   git push origin v0.2.0
   ```

3. **GitHub Actions Workflow**
   - Workflow automatically triggers
   - Builds binaries for all platforms
   - Creates GitHub Release
   - Builds and pushes Docker image

4. **Verify Release**
   - Check https://github.com/john-eevee/gosh/releases
   - Download and test binaries
   - Verify Docker image: `docker pull ghcr.io/john-eevee/gosh:v0.2.0`

### Semantic Versioning

Follow [Semantic Versioning](https://semver.org/):
- `MAJOR.MINOR.PATCH` (e.g., `1.2.3`)
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes

## Docker Usage

### Pull Image
```bash
# Latest version
docker pull ghcr.io/john-eevee/gosh:latest

# Specific version
docker pull ghcr.io/john-eevee/gosh:v0.1.1

# Specific major.minor
docker pull ghcr.io/john-eevee/gosh:0.1
```

### Run Container
```bash
# Simple GET request
docker run --rm ghcr.io/john-eevee/gosh:latest \
  get https://api.example.com/users

# With environment variables
docker run --rm \
  -e API_TOKEN=my-token \
  ghcr.io/john-eevee/gosh:latest \
  get https://api.example.com/data

# Mount .gosh directory for saved calls
docker run --rm \
  -v ~/.gosh:/root/.gosh \
  ghcr.io/john-eevee/gosh:latest \
  list
```

### Build Locally
```bash
# Build Docker image
docker build -t gosh:dev .

# Run locally built image
docker run --rm gosh:dev get https://api.example.com
```

## GitHub Status Badges

Add to README.md:

```markdown
![CI](https://github.com/john-eevee/gosh/actions/workflows/ci.yml/badge.svg?branch=main)
![CodeQL](https://github.com/john-eevee/gosh/actions/workflows/codeql.yml/badge.svg?branch=main)
[![codecov](https://codecov.io/gh/john-eevee/gosh/branch/main/graph/badge.svg)](https://codecov.io/gh/john-eevee/gosh)
```

## Troubleshooting

### CI Pipeline Failures

1. **Test Failures**
   - Check `Actions` tab on GitHub
   - Look for failing test output
   - Run `make test` locally to reproduce

2. **Lint Failures**
   - Run `make fmt` to auto-format
   - Run `make lint` to see issues
   - Fix manually if needed

3. **Build Failures**
   - Check Go version compatibility
   - Verify dependencies: `go mod tidy`
   - Run `make build` locally

### Release Issues

1. **Docker Build Fails**
   - Check Dockerfile syntax
   - Verify Go compilation
   - Check GitHub Token permissions

2. **Binary Upload Fails**
   - Verify artifact paths match workflow
   - Check file permissions
   - Verify enough space in artifacts

## Performance

### CI Runtime

Typical workflow runtimes:
- **CI (Test + Lint + Coverage)**: ~3-5 minutes
- **CodeQL**: ~2-3 minutes
- **Release (All Platforms)**: ~8-12 minutes

### Cost

GitHub Actions provides:
- Free for public repositories
- 2,000 minutes/month for private repos

## Security

### Secrets Management
- GitHub Token: Used for releases, automatically provided
- No credentials stored in workflows
- Dependabot has limited permissions

### Branch Protection

Recommended settings for `main`:
- ✅ Require status checks to pass (CI, CodeQL)
- ✅ Require code reviews before merging
- ✅ Require approval from code owners
- ✅ Dismiss stale pull request approvals
- ✅ Require branches to be up to date

## Monitoring

### GitHub Actions Dashboard
- https://github.com/john-eevee/gosh/actions

### Workflow Badges
- CI: https://github.com/john-eevee/gosh/actions/workflows/ci.yml
- Release: https://github.com/john-eevee/gosh/actions/workflows/release.yml
- CodeQL: https://github.com/john-eevee/gosh/actions/workflows/codeql.yml

### Codecov
- https://codecov.io/gh/john-eevee/gosh

## Next Steps

Optional enhancements:

1. **Code Coverage Gates**: Require minimum coverage percentage
2. **Auto-merge**: Automatically merge Dependabot PRs
3. **SLSA Provenance**: Sign releases with SLSA
4. **Multi-arch Docker**: Add arm32v7, arm64v8 support
5. **Nightly Builds**: Schedule nightly releases
6. **E2E Testing**: Add integration tests against real APIs
7. **Performance Tracking**: Monitor binary size and startup time
8. **License Compliance**: Scan dependencies for license compliance

## References

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Dependabot Documentation](https://docs.github.com/en/code-security/dependabot)
- [CodeQL Documentation](https://codeql.github.com/docs)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Go Release Guidelines](https://golang.org/doc/release)
