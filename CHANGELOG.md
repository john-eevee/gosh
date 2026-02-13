# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2026-02-13

### Added
- **Authentication Preset System**: Save and reuse authentication credentials
  - Bearer token authentication (`gosh auth add bearer <name> token=<token>`)
  - Basic authentication (`gosh auth add basic <name> username=<user> password=<pass>`)
  - Custom header authentication (`gosh auth add custom <name> header=<header> value=<value>`)
  - `gosh auth` command for preset management (list, add, remove)
  - Secure storage in `.gosh/auth.yaml` with 0600 permissions
- GitHub Actions CI/CD pipeline with automated testing and releases
- Docker support with multi-stage builds
- Comprehensive code coverage reporting
- Security scanning with Gosec

### Changed
- Updated CLI parser to support `--auth` flag
- Enhanced request builder to apply authentication
- Improved README with authentication examples

## [0.1.0] - 2026-02-10

### Added
- **HTTPie CLI Alternative**: Full HTTP request execution with easy syntax
  - Support for GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
  - JSON response pretty-printing with automatic content-type detection
  - TTY-aware colored output (green/blue/yellow/red for status codes)
  - Response timing and size statistics

- **Saved Calls System**: Persist and recall frequently used requests
  - Save requests with `--save <name>`
  - Recall with `gosh recall <name>`
  - List all calls with `gosh list`
  - Delete calls with `gosh delete <name>`
  - YAML-based storage in `.gosh/calls/` directory

- **Path Templating**: Support for parameterized URLs
  - `{variable}` syntax in URLs
  - Interactive prompts for missing variables
  - CLI parameter override with `key=value` syntax
  - `--no-interactive` flag for non-TTY environments

- **Environment Variables**: Flexible configuration
  - `${VAR_NAME}` substitution syntax
  - Load from `.env` files
  - Load from `.gosh.yaml` environment sections
  - Load from `$XDG_CONFIG_HOME/gosh/config.yaml` global config

- **Workspace Awareness**: Project-level configuration
  - Auto-detect workspace root via `.gosh.yaml`
  - Fall back to git root
  - Per-workspace configuration and saved calls
  - Workspace config overrides global config

- **Configuration System**: Multiple config file support
  - `.gosh.yaml` for workspace defaults
  - `.env` for environment variables
  - Global config support

- **Advanced Features**:
  - Pipe support for stdin request bodies
  - `--dry` flag for validation without execution
  - `--info` flag showing full response details
  - Multiple headers with `-H KEY:VALUE` syntax
  - Query parameters with `==` syntax

### Technical Details
- Single executable binary (~15-20 MB)
- Minimal dependencies (yaml.v3, bubbletea, go-isatty, stdlib)
- 23 unit tests with comprehensive coverage
- No external CLI framework (hand-rolled parser)
- Uses stdlib `net/http` client

## Repository Structure

- `cmd/gosh/` - Entry point
- `internal/app/` - Main application orchestrator
- `internal/auth/` - Authentication preset management (v0.1.1+)
- `internal/cli/` - CLI parsing
- `internal/config/` - Configuration loading
- `internal/request/` - HTTP request handling
- `internal/storage/` - Saved calls persistence
- `internal/output/` - Response formatting
- `internal/ui/` - User interface utilities
- `pkg/version/` - Version information

## Installation

### From Source
```bash
git clone https://github.com/john-eevee/gosh.git
cd gosh
go build -o gosh ./cmd/gosh
```

### From Docker
```bash
docker pull ghcr.io/john-eevee/gosh:latest
docker run --rm ghcr.io/john-eevee/gosh:latest get https://api.example.com
```

## Contributing

Contributions are welcome! Please ensure:
- All tests pass (`make test`)
- Code is formatted (`make fmt`)
- No lint errors (`make lint`)
- Coverage is maintained or improved

## License

MIT License - see LICENSE file for details
