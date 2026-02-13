# Contributing to gosh

Thank you for your interest in contributing to gosh! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and constructive in all interactions. We're building a welcoming community.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/gosh.git
   cd gosh
   ```
3. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

1. **Ensure Go 1.21+** is installed
2. **Install dependencies**:
   ```bash
   go mod download
   ```
3. **Run tests to verify setup**:
   ```bash
   make test
   ```

## Making Changes

### Code Style

- Format code with `make fmt` (uses `gofmt`)
- Run `make vet` to check for issues
- Run `make lint` for comprehensive linting
- Follow Go idioms and conventions

### Testing

- Write tests for new functionality
- Ensure all tests pass: `make test`
- Aim for >80% code coverage: `make coverage`
- Run tests across platforms if possible

### Commit Messages

Follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Build, CI, dependency updates

Example:
```
feat(auth): add OAuth2 support

Add OAuth2 authentication preset type with PKCE support.
Implements token refresh and expiration handling.

Fixes #123
```

### Pull Requests

1. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request** on GitHub with:
   - Clear description of changes
   - Reference to related issues (#123)
   - Evidence of testing (test results, screenshots if UI changes)

3. **Address review comments**
   - Make requested changes
   - Push additional commits or use squash to organize
   - All commits will be squashed when merging

## Running CI Locally

Before pushing, verify the CI pipeline will pass:

```bash
# Run all checks
make all

# Or individually:
make fmt       # Format code
make vet       # Run go vet
make lint      # Run golangci-lint
make test      # Run tests
make coverage  # Check coverage
make build     # Build binary
```

## Adding Features

### New HTTP Method Support
Edit `internal/cli/parser.go` in the `validMethods` map.

### New Configuration Option
1. Update `internal/config/types.go`
2. Update config loading in `internal/config/global.go`
3. Update README and CHANGELOG

### New Authentication Type
1. Add handler to `internal/auth/types.go`
2. Write tests in `internal/auth/types_test.go`
3. Update CLI parser in `internal/cli/parser.go`
4. Update auth command handler in `internal/app/app.go`
5. Document in README

### New CLI Command
1. Add to parser in `internal/cli/parser.go`
2. Add handler in `internal/app/app.go`
3. Write tests
4. Update help text and README

## Documentation

- Update README.md for user-facing changes
- Update CHANGELOG.md for all notable changes
- Add code comments for non-obvious logic (explain "why", not "what")
- Keep documentation in sync with code

## Performance Considerations

- Avoid unnecessary allocations in hot paths
- Use sync.Pool for frequently created objects if needed
- Profile with `make bench` before/after optimizations
- Document performance trade-offs in comments

## Security

- Never hardcode credentials or secrets
- Validate all user input
- Follow OWASP guidelines for HTTP handling
- Report security issues privately to maintainers

## Release Process

Releases are automated via GitHub Actions when tags are created:

```bash
git tag v0.2.0
git push origin v0.2.0
```

This triggers:
1. Automated testing on multiple platforms
2. Binary building for all OS/architecture combinations
3. Docker image creation
4. GitHub Release creation with artifacts
5. Changelog generation

## Questions?

- Check existing issues and discussions
- Open a new discussion for questions
- Ask in pull requests if unclear

## Recognition

Contributors will be recognized in:
- CHANGELOG.md (major contributions)
- GitHub contributors page
- Release notes for their contributions

Thank you for contributing to gosh! 🎉
