# gosh - HTTPie CLI Alternative in Go

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

## Installation

```bash
go build -o gosh ./cmd/gosh
# Move gosh to your PATH
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
cd gosh
go build -o gosh ./cmd/gosh
```

### Running Tests

```bash
go test ./...
```

## Limitations & Future Work

- Interactive prompts require TTY (use `--no-interactive` for piped input)
- Bubbletea UI integration for enhanced interactive mode
- Support for authentication presets (Basic, Bearer, Custom)
- Response caching and history
- Global saved calls (v2)
- Session management and cookies

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

MIT

## Contributing

Contributions are welcome! Please feel free to submit pull requests.
