# gosh - HTTPie CLI Alternative in Go

[![CI](https://github.com/john-eevee/gosh/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/john-eevee/gosh/actions/workflows/ci.yml)
[![CodeQL](https://github.com/john-eevee/gosh/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/john-eevee/gosh/actions/workflows/codeql.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/docker-latest-blue?logo=docker)](https://github.com/john-eevee/gosh/pkgs/container/gosh)

A lightweight, environment-aware HTTP CLI tool built with Go. Make HTTP requests from the command line with support for saved calls, templated paths, environment variables, and workspace-specific configurations.

## Features

- **Simple HTTP Requests**: GET, POST, PUT, DELETE, PATCH with easy syntax
- **Saved Calls**: Save and recall frequently used requests with `--save` and `gosh recall`
- **Path Templating**: Support for templated paths with `{variable}` syntax and interactive prompts
- **Environment Variables**: Substitute environment variables with `${VAR_NAME}` syntax
- **Workspace Aware**: Automatic detection of workspace roots and per-workspace configurations
- **Configuration Files**: 
  - `.gosh.yaml` for workspace-specific defaults
  - `.env` for environment variables
  - `$XDG_CONFIG_HOME/gosh/config.yaml` for global settings
- **Response Formatting**: Pretty-print JSON responses with automatic content-type detection
- **Pipe Support**: Read request bodies from stdin
- **TTY-Aware Coloring**: Automatic color output when connected to a terminal
- **Authentication Presets**: Save and reuse Bearer tokens, Basic auth, and custom authentication headers

## Installation

### From Source
```bash
git clone https://github.com/john-eevee/gosh.git
cd gosh
go build -o gosh ./cmd/gosh
sudo mv gosh /usr/local/bin/
```

### From GitHub Releases
Download pre-built binaries for your platform from [releases](https://github.com/john-eevee/gosh/releases):

```bash
# macOS
curl -L https://github.com/john-eevee/gosh/releases/download/v0.1.1/gosh-macos-amd64 -o gosh
chmod +x gosh
sudo mv gosh /usr/local/bin/

# Linux
curl -L https://github.com/john-eevee/gosh/releases/download/v0.1.1/gosh-linux-amd64 -o gosh
chmod +x gosh
sudo mv gosh /usr/local/bin/

# Windows
# Download from https://github.com/john-eevee/gosh/releases
```

### From Docker
```bash
docker pull ghcr.io/john-eevee/gosh:latest
docker run --rm ghcr.io/john-eevee/gosh:latest get https://api.example.com
```

### Using Homebrew (coming soon)
```bash
brew install john-eevee/gosh/gosh
```

## Quick Start

### Basic Requests

```bash
# Simple GET
gosh get https://api.example.com/users

# GET with headers
gosh get https://api.example.com/users -H Authorization:"Bearer token-123"

# POST with JSON body
gosh post https://api.example.com/users -d '{"name":"John","email":"john@example.com"}'

# Multiple headers
gosh post https://api.example.com/data \
  -H Authorization:"Bearer xyz" \
  -H Content-Type:"application/json" \
  -d '{"key":"value"}'
```

### With Template Variables

```bash
# URL with template variable - prompts for missing values
gosh get https://api.example.com/users/{userId}

# Provide template variable from CLI
gosh get https://api.example.com/users/{userId} userId=123

# Multiple template variables
gosh get https://api.example.com/users/{userId}/posts/{postId} userId=1 postId=42
```

### Environment Variables

Create a `.env` file in your workspace:

```bash
API_TOKEN=my-secret-token
API_BASE=https://dev.example.com
```

Then use in requests:

```bash
gosh get https://api.example.com/data -H Authorization:"Bearer ${API_TOKEN}"
gosh get ${API_BASE}/users
```

### Saving Requests

```bash
# Save the request for later use
gosh post https://api.example.com/users \
  -H Authorization:"Bearer token" \
  -d '{"name":"John"}' \
  --save create-user

# List all saved calls
gosh list

# Execute saved call
gosh recall create-user

# Execute saved call with parameter override
gosh recall create-user userId=456 -H Authorization:"Bearer different-token"

# Delete saved call
gosh delete create-user
```

### Dry Run & Validation

```bash
# Parse and save without executing
gosh post https://api.example.com/users \
  -d '{"name":"John"}' \
  --dry \
  --save my-request
```

### Response Information

```bash
# Show full response info (status, headers, body, timing, size)
gosh get https://api.example.com/users --info

# Regular output (just body)
gosh get https://api.example.com/users
```

### Pipe Support

```bash
# Read body from stdin
cat user.json | gosh post https://api.example.com/users

# Pipe curl response
curl https://example.com/data | gosh post https://api.example.com/transform
```

### Authentication Presets

Save and reuse authentication credentials with the `gosh auth` command:

#### Bearer Token Authentication

```bash
# Add a bearer token preset
gosh auth add bearer myapi token=my-secret-token-123

# Use it in requests
gosh get https://api.example.com/data --auth myapi

# List all presets
gosh auth list

# Remove a preset
gosh auth remove myapi
```

#### Basic Authentication

```bash
# Add basic auth (username + password)
gosh auth add basic prod-api username=john password=secret123

# Use it in requests
gosh get https://api.example.com/protected --auth prod-api
```

#### Custom Authentication

```bash
# Add custom header-based authentication
gosh auth add custom myapi header=X-API-Key value=secret-key-12345

# With prefix (e.g., "Token " prefix)
gosh auth add custom github-api header=Authorization value=ghp_token123 prefix="token "

# Use in requests
gosh get https://api.example.com/resources --auth myapi
```

#### Auth Command Reference

```bash
# Add a preset
gosh auth add <type> <name> [options]
  type: basic, bearer, or custom
  options depend on type:
    basic:   username=USER password=PASS
    bearer:  token=TOKEN
    custom:  header=HEADER value=VALUE [prefix=PREFIX]

# List all presets
gosh auth list

# Remove a preset
gosh auth remove <name>

# Use preset in request
gosh <METHOD> <URL> --auth <name>
```

## Workspace Configuration

### `.gosh.yaml` (Workspace Config)

Create a `.gosh.yaml` file in your project root to set workspace defaults:

```yaml
name: my-api-workspace
baseUrl: https://api.example.com
defaultHeaders:
  User-Agent: "gosh/1.0"
  Accept: "application/json"
environments:
  dev:
    API_TOKEN: "dev-token-123"
    API_BASE: "https://dev.example.com"
  prod:
    API_TOKEN: "prod-token-456"
    API_BASE: "https://api.example.com"
```

### `.env` (Environment Variables)

Create a `.env` file for local environment variables:

```
API_TOKEN=my-token
API_BASE=https://api.example.com
TIMEOUT=30s
```

### Global Config

Create `$XDG_CONFIG_HOME/gosh/config.yaml` (or `~/.config/gosh/config.yaml`):

```yaml
defaultEnvironment: dev
prettyPrint: true
timeout: 30s
userAgent: "gosh/1.0"
```

## Command Syntax

### HTTP Requests

```bash
gosh <METHOD> <URL> [OPTIONS] [HEADERS] [BODY]

OPTIONS:
  -H KEY:VALUE              Add header (can be multiple)
  -d DATA                   Request body
  --save NAME               Save request as NAME
  --dry                     Parse without executing
  --info                    Show full response info
  --no-interactive          Don't prompt for missing variables
  --env ENVIRONMENT         Use specific environment context
  --format json|raw|text    Output format
  --auth PRESET             Use authentication preset

TEMPLATE SYNTAX:
  {varName}                 Path/URL variable (interactive prompt)
  ${ENV_VAR}                Environment variable substitution
  key=value                 Path parameter override
```

### Saved Calls

```bash
gosh recall <name> [OVERRIDES]
gosh list
gosh delete <name>
```

### Authentication

```bash
gosh auth add <type> <name> [options]
gosh auth list
gosh auth remove <name>
```

## Examples

### Create and Execute API Request

```bash
# 1. Create a request
gosh post https://jsonplaceholder.typicode.com/posts \
  -d '{"title":"My Post","body":"Hello World","userId":1}' \
  --save new-post

# 2. View saved calls
gosh list

# 3. Execute saved call
gosh recall new-post

# 4. View details
gosh recall new-post --info
```

### Using Templates and Env Vars

```bash
# Set up workspace
cat > .gosh.yaml << EOF
name: my-api
defaultHeaders:
  Authorization: "Bearer ${API_TOKEN}"
  Accept: "application/json"
EOF

cat > .env << EOF
API_TOKEN=secret-token-123
API_BASE=https://api.example.com
EOF

# Use templated request
gosh get 'https://api.example.com/users/{userId}' userId=42

# With environment variable
gosh get https://api.example.com/posts -H Authorization:"Bearer ${API_TOKEN}"
```

### Piping Data

```bash
# Transform data through multiple requests
curl https://source.api.com/data | gosh post https://dest.api.com/import

# Use jq to format and pipe
echo '{"user":"john"}' | jq . | gosh post https://api.example.com/users
```

## Storage & Workspace Detection

Saved calls are stored in the workspace root:

```
your-project/
├── .gosh.yaml              # Workspace config
├── .env                    # Environment variables
└── .gosh/
    └── calls/
        ├── create-user.yaml
        ├── list-posts.yaml
        └── update-post.yaml
```

Workspace detection:
1. Looks for `.gosh.yaml` in current directory or parent directories
2. Falls back to git root (`.git` directory)
3. Uses current directory if no workspace marker found

## Color Output

Output is automatically colored when connected to a terminal:
- **Green**: 2xx status codes
- **Blue**: 3xx status codes  
- **Yellow**: 4xx status codes
- **Red**: 5xx status codes

Disable colors by piping output: `gosh get url | cat`

## Exit Codes

- `0`: Success
- `1`: Error (invalid arguments, network error, etc.)

## Development

### Building from Source

```bash
git clone https://github.com/john-eevee/gosh.git
cd gosh
make build
# Binary is in bin/gosh
```

### Running Tests

```bash
make test          # Run all tests
make coverage      # Generate coverage report
make lint          # Run linters
make all           # Format, lint, test, and build
```

### Using Make

```bash
make help          # Show all available targets
make install       # Install to $GOBIN
make fmt           # Format code
make vet           # Run go vet
make clean         # Remove artifacts
```

See [Makefile](Makefile) for more targets and options.

## CI/CD Pipeline

This project uses GitHub Actions for:
- **Automated Testing**: Runs on multiple platforms and Go versions
- **Code Quality**: Linting and formatting checks
- **Security**: Gosec vulnerability scanning and CodeQL analysis
- **Coverage**: Code coverage tracking with Codecov
- **Releases**: Automated multi-platform binary builds and Docker images

See [CI_CD.md](CI_CD.md) for detailed pipeline documentation.

## Limitations & Future Work

Current limitations:
- Interactive prompts require TTY (use `--no-interactive` for piped input)

Future enhancements:
- Enhanced bubbletea UI for interactive mode
- Global saved calls and credentials management
- Response caching and history
- Session management and cookie handling
- OAuth2/OIDC authentication support
- WebSocket support
- GraphQL query builder

## Comparison to HTTPie

| Feature | gosh | HTTPie |
|---------|------|--------|
| Language | Go | Python |
| Single Binary | ✅ | ❌ |
| Workspace Config | ✅ | ❌ |
| Path Templating | ✅ | ❌ |
| Saved Calls | ✅ | ❌ |
| Environment Vars | ✅ | ✅ |
| Response Formatting | ✅ | ✅ |
| Pipe Support | ✅ | ✅ |

## License

MIT - See [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Code of Conduct

Please see [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for our community standards.

## Support

- 📚 [Documentation](README.md)
- 📋 [CI/CD Documentation](CI_CD.md)
- 🐛 [Issues](https://github.com/john-eevee/gosh/issues)
- 💬 [Discussions](https://github.com/john-eevee/gosh/discussions)
- 🐳 [Docker Registry](https://github.com/john-eevee/gosh/pkgs/container/gosh)
