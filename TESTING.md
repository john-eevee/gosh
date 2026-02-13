# Testing Guide for Gosh

This document outlines how to run, write, and understand the test suite for the Gosh project.

## Overview

Gosh has a comprehensive test suite consisting of:

- **Unit Tests**: 198 tests covering individual packages
  - App package: 30 tests (67.2% coverage)
  - Auth package: 23 tests (84.6% coverage)
  - CLI package: 49 tests (88.4% coverage)
  - Config package: 15 tests (85.1% coverage)
  - Output package: 25 tests (95.2% coverage)
  - Request package: 39 tests (93.9% coverage)
  - Storage package: 14 tests (84.0% coverage)
  - UI package: 13 tests (54.2% coverage)

- **Integration Tests**: 38 tests covering end-to-end workflows
  - HTTP integration tests: 20 tests (httpbin.org)
  - Auth preset workflow: 9 tests
  - Saved calls workflow: 9 tests

- **Performance Benchmarks**: 11 benchmarks for critical paths
  - Request builder and executor
  - CLI parsing (simple and complex)
  - Storage operations (save, load, list, delete)
  - Template resolution

**Total: 236 unit + integration tests + 11 benchmarks**
**Overall Coverage: 80.3% of statements**

## Quick Start

### Run All Tests

```bash
# Run all unit and integration tests
go test ./...

# Verbose output
go test ./... -v

# With race detector (recommended)
go test ./... -race
```

### Run Only Unit Tests

```bash
# Skip integration tests (faster, no network required)
go test ./internal/... ./cmd/...
go test ./internal/... ./cmd/... -v
```

### Run Only Integration Tests

```bash
# Requires internet connection (httpbin.org)
go test ./tests/integration/... -v
```

### Run Specific Test Package

```bash
# Test auth package
go test ./internal/auth -v

# Test config package
go test ./internal/config -v

# Test output package
go test ./internal/output -v

# Test storage package
go test ./internal/storage -v

# Test UI package
go test ./internal/ui -v
```

### Run Specific Test

```bash
# Run a single test
go test ./internal/auth -run TestBasicAuth -v

# Run tests matching a pattern
go test ./internal/config -run TestDetectWorkspace -v
```

## Test Categories

### Unit Tests

Unit tests validate individual functions and packages in isolation.

#### Auth Package (`internal/auth/types_test.go`)

Tests for authentication preset types:
- Bearer token authentication
- Basic authentication  
- Custom header authentication
- Case insensitivity
- Manager operations (add, remove, get, list)
- File persistence and loading

Run: `go test ./internal/auth -v`

#### Config Package (`internal/config/config_test.go`)

Tests for configuration management:
- Loading global config from XDG paths
- Loading workspace config (.gosh.yaml)
- Loading .env files with various formats
- Workspace detection (with .gosh.yaml, with .git)
- Nested directory resolution
- Priority handling (.gosh.yaml vs .git)

Run: `go test ./internal/config -v`

#### CLI Package (`internal/cli/parser_test.go`)

Tests for command-line argument parsing:
- HTTP method parsing
- URL parsing
- Header extraction
- Query parameter parsing
- Request body handling
- Various CLI flag combinations

Run: `go test ./internal/cli -v`

#### Output Package (`internal/output/formatter_test.go`)

Tests for response formatting:
- JSON pretty-printing
- Status code colorization (TTY vs non-TTY)
- Content-type detection
- Header formatting
- Multi-value header handling
- Non-JSON content (HTML, plain text)

Run: `go test ./internal/output -v`

#### Request Package (`internal/request/builder_test.go`, `template_test.go`)

Tests for HTTP request building:
- Request creation
- Template variable substitution
- Environment variable interpolation
- Edge cases and error handling

Run: `go test ./internal/request -v`

#### Storage Package (`internal/storage/manager_test.go`)

Tests for saved calls management:
- Saving and loading calls
- Call listing
- Call deletion
- File I/O operations

Run: `go test ./internal/storage -v`

#### UI Package (`internal/ui/prompt_test.go`)

Tests for interactive prompts:
- Prompt model rendering
- User input handling
- Whitespace trimming
- Special character handling

Run: `go test ./internal/ui -v`

### Integration Tests

Integration tests validate complete workflows using real HTTP requests to httpbin.org.

#### HTTP Integration Tests (`tests/integration/http_test.go`)

20 tests covering all HTTP functionality:
- HTTP methods (GET, POST, PUT, DELETE, PATCH, HEAD)
- Headers transmission and received headers
- Query parameters
- Request/response bodies
- Status codes (2xx, 3xx, 4xx, 5xx)
- Response timing measurement
- Response size calculation
- Redirects
- Large responses
- Timeout handling
- User-Agent headers
- Concurrent requests
- Empty bodies
- Special characters and Unicode

Run: `go test ./tests/integration -run "^TestGet|^TestPost|^TestPut" -v`

#### Auth Preset Integration Tests (`tests/integration/auth_test.go`)

9 tests for authentication workflows:
- Creating and using bearer token presets
- Creating and using basic auth presets
- Creating and using custom header presets
- Auth preset persistence (save/load)
- Listing auth presets
- Removing auth presets
- Authorization header format verification
- Using auth with actual HTTP requests
- File permission security (0600)

Run: `go test ./tests/integration -run "^TestAuthPreset" -v`

#### Saved Calls Integration Tests (`tests/integration/saved_calls_test.go`)

9 tests for saved call workflow:
- Saving basic calls
- Recalling/loading saved calls
- Listing all saved calls
- Deleting calls
- Saving calls with all fields (headers, query params, body)
- Using saved calls to execute actual HTTP requests
- Overwriting existing calls
- Handling empty call lists
- Special characters in call names

Run: `go test ./tests/integration -run "^TestSaveCall" -v`

## Running Tests with Different Options

### With Race Detector

Detects data races and concurrent access issues:

```bash
go test -race ./...
```

### With Coverage

Generate code coverage reports:

```bash
# Generate coverage file
go test -coverprofile=coverage.out ./...

# View coverage in HTML
go tool cover -html=coverage.out

# View coverage summary
go tool cover -func=coverage.out | tail -20
```

### With Timeout

Set timeout for all tests (useful for CI/CD):

```bash
# Unit tests (5 minutes)
go test -timeout=5m ./internal/... ./cmd/...

# Integration tests (10 minutes, requires network)
go test -timeout=10m ./tests/integration/...
```

### With Verbosity

Show more detailed output:

```bash
# Show all tests being run
go test -v ./...

# Show test output even for passing tests
go test -v -run TestName ./...
```

## CI/CD Integration

The project uses GitHub Actions for CI/CD. See `.github/workflows/ci.yml` for details.

### Test Stages

1. **Unit Tests** (5 min timeout, runs on all OS/Go versions)
   - Tested on: Ubuntu, macOS, Windows
   - Go versions: 1.21, 1.22

2. **Integration Tests** (10 min timeout, skips if offline)
   - Runs on Ubuntu only
   - Continues on error if network unavailable
   - Skips automatically if httpbin.org is unreachable

3. **Linting** (golangci-lint)
   - Runs on Ubuntu
   - Checks formatting with gofmt

4. **Coverage** (Codecov)
   - Generates coverage for unit tests only
   - Uploads to Codecov
   - Excludes integration tests from coverage metrics

5. **Security Scanning** (Gosec)
   - Runs on Ubuntu
   - Generates SARIF report for GitHub Security

## Integration Test Requirements

Integration tests require:

- **Network Access**: Connection to httpbin.org
- **Internet**: Active internet connection
- **Time**: ~20 seconds to run all integration tests

Integration tests will automatically skip if:

- httpbin.org is unreachable
- Network connection is unavailable
- Running in offline environments

All integration tests include `SkipIfOffline()` at the beginning.

## Writing New Tests

### Unit Test Template

```go
package mypackage

import "testing"

func TestFeatureName(t *testing.T) {
	// Arrange
	input := setupData()
	
	// Act
	result := FunctionUnderTest(input)
	
	// Assert
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
```

### Integration Test Template

```go
package integration

import "testing"

func TestFeatureIntegration(t *testing.T) {
	SkipIfOffline(t)  // Skip if network unavailable
	
	// Create resources
	executor := request.NewExecutor(timeout)
	
	// Execute HTTP request
	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	// Verify results
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
```

### Test Naming Convention

- Tests must start with `Test`
- Use descriptive names: `TestParseURLWithQuery` not `TestParse`
- One assertion per test when possible
- Test names should describe the scenario

### Table-Driven Tests

For testing multiple cases:

```go
func TestMultipleCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"case1", "input1", "expected1"},
		{"case2", "input2", "expected2"},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Function(test.input)
			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
}
```

## Common Testing Patterns

### Temporary Files

```go
tempDir := t.TempDir()  // Automatically cleaned up
filePath := filepath.Join(tempDir, "test.txt")
```

### Mocking stdin/stdout

```go
oldStdin := os.Stdin
defer func() { os.Stdin = oldStdin }()

tempFile, _ := os.CreateTemp("", "test-")
os.Stdin = tempFile
// ... run tests
tempFile.Close()
```

### Error Handling

```go
err := functionThatErrors()
if err != nil {
	t.Fatalf("expected no error, got %v", err)
}
```

### Assertions

```go
// Equality
if result != expected {
	t.Errorf("expected %v, got %v", expected, result)
}

// Presence in collection
found := false
for _, item := range items {
	if item == target {
		found = true
		break
	}
}
if !found {
	t.Error("expected item not found")
}

// String contains
if !bytes.Contains([]byte(result), []byte(substring)) {
	t.Errorf("expected %q to contain %q", result, substring)
}
```

## Troubleshooting

### Integration Tests Timeout

If integration tests timeout:

```bash
# Increase timeout
go test -timeout=20m ./tests/integration/...
```

### Network Connectivity Issues

```bash
# Check if httpbin.org is reachable
curl -s https://httpbin.org/status/200

# Integration tests will automatically skip if unreachable
```

### Coverage Reports

```bash
# Generate and view coverage
go test -coverprofile=coverage.out ./internal/... ./cmd/...
go tool cover -html=coverage.out

# View uncovered lines
go tool cover -html=coverage.out  # Look for red sections
```

### Race Conditions

```bash
# Run with race detector to find concurrent access issues
go test -race ./...
```

## Continuous Improvement

### Coverage Goals

- Minimum unit test coverage: 80%
- Current coverage: 80.3% overall (see breakdown above)
- Integration tests provide end-to-end validation
- All critical paths have benchmarks

### Performance Testing

Run benchmarks to measure performance:

```bash
# Run all benchmarks
go test -bench=. ./...

# Run benchmarks for specific package
go test -bench=. ./internal/request/...

# Include memory allocations
go test -bench=. -benchmem ./internal/request/...

# Run with specific count
go test -bench=. -count=5 ./internal/request/...
```

Available benchmarks:
- `BenchmarkBuilderBuild`: Request building performance
- `BenchmarkExecutorExecute`: HTTP request execution
- `BenchmarkTemplateResolve`: Template variable resolution
- `BenchmarkExtractPathVars`: Path variable extraction
- `BenchmarkExtractEnvVars`: Environment variable extraction
- `BenchmarkParseSimpleRequest`: Simple CLI parsing
- `BenchmarkParseComplexRequest`: Complex CLI parsing
- `BenchmarkSaveCall`: Saving HTTP calls to disk
- `BenchmarkLoadCall`: Loading HTTP calls from disk
- `BenchmarkListCalls`: Listing all saved calls
- `BenchmarkExistsCall`: Checking call existence

### Test Maintenance

- Update tests when changing functionality
- Keep integration tests in sync with actual API behaviors
- Document breaking changes in test files
- Run full test suite before committing
- Check benchmarks don't regress significantly

### Performance

- Unit tests should complete in < 5 minutes
- Integration tests should complete in < 10 minutes
- Use `-race` flag in CI to catch concurrency issues
- Profile slow tests with `go test -cpuprofile=cpu.prof`

## Related Documentation

- [Contributing Guidelines](CONTRIBUTING.md)
- [CI/CD Pipeline](CI_CD.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Changelog](CHANGELOG.md)
