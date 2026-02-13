# GOSH: HTTPie CLI Alternative - Implementation Plan

## 1. Project Structure
```
gosh/
├── cmd/
│   └── gosh/
│       └── main.go                 # Entry point, CLI bootstrapping
├── internal/
│   ├── app/
│   │   └── app.go                  # Main application orchestrator
│   ├── cli/
│   │   ├── parser.go               # Parse CLI flags and positional args
│   │   └── types.go                # CLI input data structures
│   ├── config/
│   │   ├── workspace.go            # Detect & load .gosh.yaml + git root
│   │   ├── global.go               # Load $XDG_CONFIG_HOME/gosh/config.yaml
│   │   ├── env.go                  # Load .env + environment substitution
│   │   └── types.go                # Config data structures
│   ├── request/
│   │   ├── builder.go              # Construct http.Request from parsed CLI
│   │   ├── executor.go             # Execute HTTP request & capture response
│   │   ├── template.go             # Template variable resolution ({var}, ${VAR})
│   │   └── types.go                # Request/Response types
│   ├── storage/
│   │   ├── manager.go              # Save/load/list saved calls
│   │   ├── path.go                 # Workspace storage paths
│   │   └── types.go                # SavedCall YAML structure
│   ├── output/
│   │   ├── formatter.go            # Format response (JSON pretty-print, headers, etc.)
│   │   ├── colorer.go              # Conditional coloring based on TTY
│   │   └── writer.go               # Write output to stdout
│   └── ui/
│       └── prompt.go               # Interactive prompts (bubbletea-based) for missing vars
├── pkg/
│   └── version.go                  # Version constant
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## 2. CLI Command Structure

```bash
# Request execution (primary)
gosh <METHOD> <URL> [OPTIONS] [HEADERS] [BODY_PARAMS]

# Saved call operations
gosh recall <name> [OVERRIDE_PARAMS]
gosh list [--workspace-only]
gosh delete <name>

# Utility
gosh --version
gosh --help
```

## 3. Command Examples & Syntax

```bash
# Simple GET
gosh get https://api.example.com/users

# With headers (one -H per header)
gosh post https://api.example.com/users \
  -H Authorization:"Bearer xyz" \
  -H Content-Type:"application/json"

# With body from stdin (raw string or piped)
echo '{"name":"John"}' | gosh post https://api.example.com/users

# With template paths (auto-prompts for missing {userId})
gosh get https://api.example.com/users/{userId}

# Override template at CLI
gosh get https://api.example.com/users/{userId} userId=42

# With environment variables
gosh get https://api.example.com/users/${API_BASE}/list
  # ${API_BASE} substituted from .env or .gosh.yaml environment

# Save request
gosh post https://api.example.com/users \
  -H Authorization:"Bearer xyz" \
  --save my-create-user

# Dry run (parse & validate, don't execute)
gosh post https://api.example.com/users \
  --dry \
  --save my-request

# Execute with full response info (status, headers, body, timing, size)
gosh get https://api.example.com/users --info

# Recall saved call with parameter override
gosh recall my-create-user \
  -H Authorization:"Bearer prod-token" \
  userId=99

# List all saved calls (workspace scope)
gosh list

# Delete saved call
gosh delete my-create-user
```

## 4. Data Structures & File Formats

**Saved Call** (`.gosh/calls/my-request.yaml`):
```yaml
name: my-request
method: POST
url: https://api.example.com/users
headers:
  Authorization: "Bearer ${API_TOKEN}"
  Content-Type: "application/json"
queryParams:
  limit: "10"
body: '{"name":"John","email":"john@example.com"}'
description: "Create a new user"
createdAt: "2026-02-12T22:30:00Z"
```

**Workspace Config** (`.gosh.yaml`):
```yaml
name: my-workspace
baseUrl: https://api.example.com
defaultHeaders:
  User-Agent: "gosh/1.0"
  Accept: "application/json"
environments:
  dev:
    API_TOKEN: "dev-token-123"
    API_BASE: "https://dev.example.com"
  prod:
    API_TOKEN: "${PROD_API_TOKEN}"  # From system env
    API_BASE: "https://api.example.com"
```

**Global Config** (`$XDG_CONFIG_HOME/gosh/config.yaml`):
```yaml
defaultEnvironment: dev
prettyPrint: true
timeout: 30s
userAgent: "gosh/1.0"
```

**.env File** (`.env` in workspace):
```
API_TOKEN=my-secret-token
API_BASE=https://dev.example.com
PROXY_URL=http://proxy.company.com:8080
```

## 5. Feature Breakdown by Priority

**Phase 1: Core (MVP)**
- [ ] CLI flag parsing (`-H`, `--save`, `--dry`, `--info`)
- [ ] Request builder (GET, POST, PUT, DELETE, PATCH)
- [ ] HTTP execution with response capture
- [ ] Path templating `{varName}` with interactive prompts
- [ ] Workspace detection (.gosh.yaml + git root fallback)
- [ ] Save calls to `.gosh/calls/` as YAML
- [ ] Load & execute saved calls with `recall`
- [ ] Basic response formatting (JSON pretty-print, status line)
- [ ] Conditional coloring (TTY detection)
- [ ] `--info` flag (full response with timing/size)

**Phase 2: Configuration & Environment**
- [ ] `.env` file loading
- [ ] Environment variable substitution `${VAR_NAME}`
- [ ] `.gosh.yaml` config loading (headers, base URL, environments)
- [ ] Global config at `$XDG_CONFIG_HOME/gosh/config.yaml`
- [ ] Config merge (workspace overrides global)
- [ ] `--env` flag to select environment context

**Phase 3: Request Operations**
- [ ] Pipe support (stdin → request body)
- [ ] `gosh list` to show saved calls
- [ ] `gosh delete <name>` to remove saved calls
- [ ] `recall` parameter override (merge mode)
- [ ] `--no-interactive` flag for non-TTY contexts

**Phase 4: Polish**
- [ ] Interactive bubbletea UI for template prompts
- [ ] Error messages (clear, actionable)
- [ ] Response output formatting (`--format=json|raw|text`)
- [ ] Tests & documentation
- [ ] Release v1.0

## 6. Technology Choices (Confirmed)

| Component | Tool | Rationale |
|-----------|------|-----------|
| HTTP Client | `net/http` (stdlib) | Lightweight, built-in, no external deps |
| CLI Parsing | Hand-rolled flag parser | Simple, minimal deps |
| TUI/Interactive | `bubbletea` | Modern, composable CLI UI |
| Config Format | YAML via `gopkg.in/yaml.v3` | Human-readable, nested structures |
| Templating | `text/template` | Stdlib, familiar to Go developers |
| Storage | YAML files in workspace | Simple, git-ignorable, version-controllable |

## 7. Workspace & Storage Paths

```
Current working directory or parent:
.gosh.yaml                          # Workspace config

$WORKSPACE/.gosh/
├── calls/
│   ├── my-request.yaml
│   ├── get-users.yaml
│   └── create-post.yaml
└── .env                             # Optional workspace env file

$XDG_DATA_HOME/gosh/
└── global/                          # Future: global saved calls (v2)

$XDG_CONFIG_HOME/gosh/
└── config.yaml                      # Global config
```

## 8. Implementation Phases & Deliverables

**Week 1: Foundation**
- Project structure & go.mod setup
- CLI flag parsing (basic flags: `-H`, `--save`, `--dry`, `--info`)
- HTTP request builder
- Response capture & basic formatting

**Week 2: Core Features**
- Path templating with interactive prompts (bubbletea)
- Workspace detection
- Saved call storage & retrieval
- `recall` command with parameter override

**Week 3: Configuration & Environment**
- `.gosh.yaml` loading
- `.env` support
- Global config loading
- Environment variable substitution

**Week 4: Polish & Release**
- Pipe support for bodies
- `list` and `delete` commands
- Output formatting options
- Tests & documentation
- v1.0 release

## 9. Configuration Details

### Defaults & Behaviors

**Request Method**: Case-insensitive (GET/get both work)
**Timeout**: Default 30 seconds, configurable via config
**Redirects**: Follow by default (up to 10 hops)
**Query Params**: Merge with URL query string if present
**Body Source Priority**: stdin > CLI argument

### Storage Conventions

- Workspace root determined by: `.gosh.yaml` presence OR git root
- Saved calls stored in `.gosh/calls/` within workspace
- All YAML files use workspace-relative paths
- Timestamps in SavedCall use RFC3339 format

### Environment Variable Substitution

- Path parameters: `{varName}` → interactive prompt or CLI override
- Environment variables: `${VAR_NAME}` → from .env, .gosh.yaml, or system env
- Substitution happens AFTER config loading, BEFORE request execution
- Unresolved variables cause clear error messages

## 10. Error Handling Strategy

- **Invalid URL**: "Invalid URL format: {url}"
- **Missing method**: "HTTP method required (GET, POST, etc.)"
- **Template var missing**: Interactive prompt, skip if `--no-interactive`
- **Config load error**: "Failed to load .gosh.yaml: {error}"
- **Network error**: "Request failed: {error} (code: {status_code})"
- **Save error**: "Failed to save call: {error}"

## 11. Success Criteria for Each Phase

**Phase 1 Complete When**:
- Can execute basic HTTP requests (GET, POST, etc.)
- Can save/recall requests
- Template variables work with interactive prompts
- --dry and --save flags function correctly
- Response displays with proper formatting

**Phase 2 Complete When**:
- .gosh.yaml loads and applies defaults
- .env files load
- Environment variable substitution works
- Global config loads when present
- --env flag selects environment context

**Phase 3 Complete When**:
- stdin pipes into request body
- `list` shows all saved calls
- `delete` removes saved calls
- `recall` accepts parameter overrides
- --no-interactive flag suppresses prompts

**Phase 4 Complete When**:
- Bubbletea UI provides smooth interactive prompts
- Error messages are clear and actionable
- All output formatting options work
- Tests cover critical paths
- Documentation is complete
