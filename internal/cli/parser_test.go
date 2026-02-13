package cli

import (
	"testing"
)

func TestParseGetRequest(t *testing.T) {
	parser := NewParser([]string{"get", "https://api.example.com/users"})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req, ok := result.(*ParsedRequest)
	if !ok {
		t.Fatalf("expected ParsedRequest, got %T", result)
	}

	if req.Method != "GET" {
		t.Errorf("expected method GET, got %s", req.Method)
	}

	if req.URL != "https://api.example.com/users" {
		t.Errorf("expected URL https://api.example.com/users, got %s", req.URL)
	}
}

func TestParseWithHeaders(t *testing.T) {
	parser := NewParser([]string{
		"post", "https://api.example.com/users",
		"-H", "Authorization:Bearer token",
		"-H", "Content-Type:application/json",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if len(req.Headers) != 2 {
		t.Errorf("expected 2 headers, got %d", len(req.Headers))
	}

	if req.Headers["Authorization"] != "Bearer token" {
		t.Errorf("expected Authorization header 'Bearer token', got '%s'", req.Headers["Authorization"])
	}

	if req.Headers["Content-Type"] != "application/json" {
		t.Errorf("expected Content-Type header 'application/json', got '%s'", req.Headers["Content-Type"])
	}
}

func TestParseWithBody(t *testing.T) {
	parser := NewParser([]string{
		"post", "https://api.example.com/users",
		"-d", `{"name":"John"}`,
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Body != `{"name":"John"}` {
		t.Errorf("expected body '{\"name\":\"John\"}', got '%s'", req.Body)
	}
}

func TestParseWithSaveFlag(t *testing.T) {
	parser := NewParser([]string{
		"post", "https://api.example.com/users",
		"--save", "create-user",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Save != "create-user" {
		t.Errorf("expected Save 'create-user', got '%s'", req.Save)
	}
}

func TestParseWithPathParams(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users/{userId}",
		"userId=123",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.PathParams["userId"] != "123" {
		t.Errorf("expected userId param '123', got '%s'", req.PathParams["userId"])
	}
}

func TestParseRecall(t *testing.T) {
	parser := NewParser([]string{
		"recall", "create-user",
		"userId=456",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	opts, ok := result.(*RecallOptions)
	if !ok {
		t.Fatalf("expected RecallOptions, got %T", result)
	}

	if opts.Name != "create-user" {
		t.Errorf("expected name 'create-user', got '%s'", opts.Name)
	}

	if opts.ParameterOverride["userId"] != "456" {
		t.Errorf("expected userId override '456', got '%s'", opts.ParameterOverride["userId"])
	}
}

func TestParseWithQueryParams(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users",
		"limit==10",
		"offset==20",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.QueryParams["limit"] != "10" {
		t.Errorf("expected limit query param '10', got '%s'", req.QueryParams["limit"])
	}

	if req.QueryParams["offset"] != "20" {
		t.Errorf("expected offset query param '20', got '%s'", req.QueryParams["offset"])
	}
}

func TestParseInvalidMethod(t *testing.T) {
	parser := NewParser([]string{
		"invalid", "https://api.example.com/users",
	})
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error for invalid method")
	}
}

func TestParseList(t *testing.T) {
	parser := NewParser([]string{"list"})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "list" {
		t.Errorf("expected 'list', got %v", result)
	}
}

// TestParseDelete tests delete command parsing
func TestParseDelete(t *testing.T) {
	parser := NewParser([]string{"delete", "my-call"})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "delete:my-call"
	if result != expected {
		t.Errorf("expected %q, got %v", expected, result)
	}
}

// TestParseDeleteMissingName tests delete without name
func TestParseDeleteMissingName(t *testing.T) {
	parser := NewParser([]string{"delete"})
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error for missing call name")
	}
}

// TestParseAuthList tests auth list command
func TestParseAuthList(t *testing.T) {
	parser := NewParser([]string{"auth"})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd, ok := result.(*AuthCommand)
	if !ok {
		t.Fatalf("expected AuthCommand, got %T", result)
	}

	if cmd.Subcommand != "list" {
		t.Errorf("expected subcommand 'list', got %q", cmd.Subcommand)
	}
}

// TestParseAuthAdd tests auth add command
func TestParseAuthAdd(t *testing.T) {
	parser := NewParser([]string{
		"auth", "add", "bearer", "my-api",
		"token=abc123",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd := result.(*AuthCommand)
	if cmd.Subcommand != "add" {
		t.Errorf("expected subcommand 'add', got %q", cmd.Subcommand)
	}

	if cmd.Type != "bearer" {
		t.Errorf("expected type 'bearer', got %q", cmd.Type)
	}

	if cmd.Name != "my-api" {
		t.Errorf("expected name 'my-api', got %q", cmd.Name)
	}

	if cmd.Flags["token"] != "abc123" {
		t.Errorf("expected token 'abc123', got %q", cmd.Flags["token"])
	}
}

// TestParseAuthAddMissingArgs tests auth add with missing arguments
func TestParseAuthAddMissingArgs(t *testing.T) {
	parser := NewParser([]string{"auth", "add", "bearer"})
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error for missing arguments")
	}
}

// TestParseAuthRemove tests auth remove command
func TestParseAuthRemove(t *testing.T) {
	parser := NewParser([]string{"auth", "remove", "my-preset"})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd := result.(*AuthCommand)
	if cmd.Subcommand != "remove" {
		t.Errorf("expected subcommand 'remove', got %q", cmd.Subcommand)
	}

	if cmd.Name != "my-preset" {
		t.Errorf("expected name 'my-preset', got %q", cmd.Name)
	}
}

// TestParseAuthDelete tests auth delete alias for remove
func TestParseAuthDelete(t *testing.T) {
	parser := NewParser([]string{"auth", "delete", "my-preset"})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd := result.(*AuthCommand)
	if cmd.Subcommand != "remove" {
		t.Errorf("expected subcommand 'remove', got %q", cmd.Subcommand)
	}
}

// TestParseAuthRemoveMissingName tests auth remove without name
func TestParseAuthRemoveMissingName(t *testing.T) {
	parser := NewParser([]string{"auth", "remove"})
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error for missing preset name")
	}
}

// TestParseAuthInvalidSubcommand tests invalid auth subcommand
func TestParseAuthInvalidSubcommand(t *testing.T) {
	parser := NewParser([]string{"auth", "invalid"})
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error for invalid auth subcommand")
	}
}

// TestParseRecallWithHeaders tests recall with header overrides
func TestParseRecallWithHeaders(t *testing.T) {
	parser := NewParser([]string{
		"recall", "get-users",
		"-H", "Authorization:Bearer new-token",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	opts := result.(*RecallOptions)
	if opts.Name != "get-users" {
		t.Errorf("expected name 'get-users', got %q", opts.Name)
	}

	if opts.Headers["Authorization"] != "Bearer new-token" {
		t.Errorf("expected Authorization header, got %q", opts.Headers["Authorization"])
	}
}

// TestParseRecallMissingName tests recall without call name
func TestParseRecallMissingName(t *testing.T) {
	parser := NewParser([]string{"recall"})
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error for missing call name")
	}
}

// TestParseRecallWithEnv tests recall with env override
func TestParseRecallWithEnv(t *testing.T) {
	parser := NewParser([]string{
		"recall", "my-call",
		"--env=production",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	opts := result.(*RecallOptions)
	if opts.Env != "production" {
		t.Errorf("expected env 'production', got %q", opts.Env)
	}
}

// TestParseRecallWithEnvSeparate tests recall with env as separate arg
func TestParseRecallWithEnvSeparate(t *testing.T) {
	parser := NewParser([]string{
		"recall", "my-call",
		"--env", "staging",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	opts := result.(*RecallOptions)
	if opts.Env != "staging" {
		t.Errorf("expected env 'staging', got %q", opts.Env)
	}
}

// TestParseWithDry tests dry run flag
func TestParseWithDry(t *testing.T) {
	parser := NewParser([]string{
		"post", "https://api.example.com/users",
		"--dry",
		"--save", "create-user",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if !req.Dry {
		t.Errorf("expected Dry flag to be true")
	}
}

// TestParseWithInfo tests info flag
func TestParseWithInfo(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users",
		"--info",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if !req.Info {
		t.Errorf("expected Info flag to be true")
	}
}

// TestParseWithNoInteractive tests no-interactive flag
func TestParseWithNoInteractive(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users",
		"--no-interactive",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if !req.NoInteractive {
		t.Errorf("expected NoInteractive flag to be true")
	}
}

// TestParseWithEnv tests env flag
func TestParseWithEnv(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users",
		"--env=development",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Env != "development" {
		t.Errorf("expected env 'development', got %q", req.Env)
	}
}

// TestParseWithEnvSeparate tests env as separate argument
func TestParseWithEnvSeparate(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users",
		"--env", "production",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Env != "production" {
		t.Errorf("expected env 'production', got %q", req.Env)
	}
}

// TestParseWithFormat tests format flag
func TestParseWithFormat(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users",
		"--format=json",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Format != "json" {
		t.Errorf("expected format 'json', got %q", req.Format)
	}
}

// TestParseWithFormatSeparate tests format as separate argument
func TestParseWithFormatSeparate(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users",
		"--format", "raw",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Format != "raw" {
		t.Errorf("expected format 'raw', got %q", req.Format)
	}
}

// TestParseWithAuth tests auth flag
func TestParseWithAuth(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users",
		"--auth=my-api",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Auth != "my-api" {
		t.Errorf("expected auth 'my-api', got %q", req.Auth)
	}
}

// TestParseWithAuthSeparate tests auth as separate argument
func TestParseWithAuthSeparate(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users",
		"--auth", "my-preset",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Auth != "my-preset" {
		t.Errorf("expected auth 'my-preset', got %q", req.Auth)
	}
}

// TestParseHeaderWithEquals tests header parsing with = separator
func TestParseHeaderWithEquals(t *testing.T) {
	parser := NewParser([]string{
		"post", "https://api.example.com/users",
		"-H", "X-Custom=value",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Headers["X-Custom"] != "value" {
		t.Errorf("expected header value 'value', got %q", req.Headers["X-Custom"])
	}
}

// TestParseHeaderInline tests header flag without space
func TestParseHeaderInline(t *testing.T) {
	parser := NewParser([]string{
		"post", "https://api.example.com/users",
		"-HX-Custom:inline",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Headers["X-Custom"] != "inline" {
		t.Errorf("expected header value 'inline', got %q", req.Headers["X-Custom"])
	}
}

// TestParseBodyInline tests body flag without space
func TestParseBodyInline(t *testing.T) {
	parser := NewParser([]string{
		"post", "https://api.example.com/users",
		`-d{"test":"data"}`,
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Body != `{"test":"data"}` {
		t.Errorf("expected body, got %q", req.Body)
	}
}

// TestParseComplexRequest tests parsing complex request with many options
func TestParseComplexRequest(t *testing.T) {
	parser := NewParser([]string{
		"post", "https://api.example.com/users/{id}",
		"-H", "Authorization:Bearer token",
		"-H", "Content-Type:application/json",
		"-d", `{"name":"John"}`,
		"--save", "create-user",
		"--info",
		"id=123",
		"limit==10",
		"--auth=my-api",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Method != "POST" {
		t.Errorf("method: got %q", req.Method)
	}
	if req.URL != "https://api.example.com/users/{id}" {
		t.Errorf("URL: got %q", req.URL)
	}
	if len(req.Headers) != 2 {
		t.Errorf("headers count: expected 2, got %d", len(req.Headers))
	}
	if req.Body != `{"name":"John"}` {
		t.Errorf("body: got %q", req.Body)
	}
	if req.Save != "create-user" {
		t.Errorf("save: got %q", req.Save)
	}
	if !req.Info {
		t.Errorf("info flag not set")
	}
	if req.PathParams["id"] != "123" {
		t.Errorf("path param: got %q", req.PathParams["id"])
	}
	if req.QueryParams["limit"] != "10" {
		t.Errorf("query param: got %q", req.QueryParams["limit"])
	}
	if req.Auth != "my-api" {
		t.Errorf("auth: got %q", req.Auth)
	}
}

// TestParseUnexpectedArgument tests error on unexpected argument
func TestParseUnexpectedArgument(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users",
		"--unknown-flag",
	})
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error for unexpected argument")
	}
}

// TestParseMissingHeaderValue tests missing value for header flag
func TestParseMissingHeaderValue(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users",
		"-H",
	})
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error for missing header value")
	}
}

// TestParseMissingBodyValue tests missing value for body flag
func TestParseMissingBodyValue(t *testing.T) {
	parser := NewParser([]string{
		"post", "https://api.example.com/users",
		"-d",
	})
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error for missing body value")
	}
}

// TestParseInvalidHeaderFormat tests invalid header format
func TestParseInvalidHeaderFormat(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users",
		"-H", "InvalidHeaderNoSeparator",
	})
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error for invalid header format")
	}
}

// TestParseMissingMethodAndURL tests error when method and URL are missing
func TestParseMissingMethodAndURL(t *testing.T) {
	parser := NewParser([]string{"get"})
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}

// TestParseVersion tests version flag
func TestParseVersion(t *testing.T) {
	parser := NewParser([]string{"--version"})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "version" {
		t.Errorf("expected 'version', got %v", result)
	}
}

// TestParseVersionShort tests version short flag
func TestParseVersionShort(t *testing.T) {
	parser := NewParser([]string{"-v"})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "version" {
		t.Errorf("expected 'version', got %v", result)
	}
}

// TestParseHelp tests help flag
func TestParseHelp(t *testing.T) {
	parser := NewParser([]string{"--help"})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "help" {
		t.Errorf("expected 'help', got %v", result)
	}
}

// TestParseHelpShort tests help short flag
func TestParseHelpShort(t *testing.T) {
	parser := NewParser([]string{"-h"})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "help" {
		t.Errorf("expected 'help', got %v", result)
	}
}
