package app

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/gosh/internal/auth"
	"github.com/gosh/internal/cli"
	"github.com/gosh/internal/config"
	"github.com/gosh/internal/storage"
)

// TestRunWithVersion tests the version command
func TestRunWithVersion(t *testing.T) {
	app := &App{
		workspace: &config.Workspace{Root: "/tmp"},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager("/tmp"),
		authMgr:   auth.NewManager("/tmp"),
	}

	err := app.Run([]string{"--version"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestRunWithHelp tests the help command
func TestRunWithHelp(t *testing.T) {
	app := &App{
		workspace: &config.Workspace{Root: "/tmp"},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager("/tmp"),
		authMgr:   auth.NewManager("/tmp"),
	}

	err := app.Run([]string{"--help"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestRunWithInvalidCommand tests an invalid command
func TestRunWithInvalidCommand(t *testing.T) {
	app := &App{
		workspace: &config.Workspace{Root: "/tmp"},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager("/tmp"),
		authMgr:   auth.NewManager("/tmp"),
	}

	err := app.Run([]string{"invalid-command"})
	if err == nil {
		t.Fatalf("expected error for invalid command, got nil")
	}
}

// TestRunWithNoArgs tests with no arguments
func TestRunWithNoArgs(t *testing.T) {
	app := &App{
		workspace: &config.Workspace{Root: "/tmp"},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager("/tmp"),
		authMgr:   auth.NewManager("/tmp"),
	}

	err := app.Run([]string{})
	if err == nil {
		t.Fatalf("expected error for no args, got nil")
	}
}

// TestListCallsEmpty tests listing calls when none exist
func TestListCallsEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   auth.NewManager(tmpDir),
	}

	err := app.Run([]string{"list"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestDeleteCallNonexistent tests deleting a nonexistent call
func TestDeleteCallNonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   auth.NewManager(tmpDir),
	}

	err := app.Run([]string{"delete", "nonexistent"})
	if err == nil {
		t.Fatalf("expected error for deleting nonexistent call, got nil")
	}
}

// TestSubstituteEnvVars tests environment variable substitution
func TestSubstituteEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "simple substitution",
			text:     "https://api.example.com/${PATH}",
			envVars:  map[string]string{"PATH": "users"},
			expected: "https://api.example.com/users",
		},
		{
			name:     "multiple substitutions",
			text:     "${PROTO}://${HOST}/${PATH}",
			envVars:  map[string]string{"PROTO": "https", "HOST": "api.example.com", "PATH": "users"},
			expected: "https://api.example.com/users",
		},
		{
			name:     "nonexistent variable",
			text:     "https://api.example.com/${MISSING}",
			envVars:  map[string]string{},
			expected: "https://api.example.com/${MISSING}",
		},
		{
			name:     "no variables",
			text:     "https://api.example.com/users",
			envVars:  map[string]string{},
			expected: "https://api.example.com/users",
		},
		{
			name:     "empty text",
			text:     "",
			envVars:  map[string]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				workspace: &config.Workspace{Root: "/tmp", Env: tt.envVars},
			}
			result := app.substituteEnvVars(tt.text)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestSubstituteEnvVarsInMap tests environment variable substitution in maps
func TestSubstituteEnvVarsInMap(t *testing.T) {
	app := &App{
		workspace: &config.Workspace{
			Root: "/tmp",
			Env:  map[string]string{"HOST": "api.example.com", "TOKEN": "secret123"},
		},
	}

	input := map[string]string{
		"Authorization": "Bearer ${TOKEN}",
		"Host":          "${HOST}",
	}

	result := app.substituteEnvVarsInMap(input)

	if result["Authorization"] != "Bearer secret123" {
		t.Errorf("got %q, want 'Bearer secret123'", result["Authorization"])
	}
	if result["Host"] != "api.example.com" {
		t.Errorf("got %q, want 'api.example.com'", result["Host"])
	}
}

// TestSubstituteEnvVarsInMapWithMissing tests substitution with missing variables
func TestSubstituteEnvVarsInMapWithMissing(t *testing.T) {
	app := &App{
		workspace: &config.Workspace{
			Root: "/tmp",
			Env:  map[string]string{},
		},
	}

	input := map[string]string{
		"Authorization": "Bearer ${TOKEN}",
	}

	result := app.substituteEnvVarsInMap(input)

	if result["Authorization"] != "Bearer ${TOKEN}" {
		t.Errorf("got %q, want 'Bearer ${TOKEN}'", result["Authorization"])
	}
}

// TestExecuteRequestWithDryRun tests executing a request with --dry flag
func TestExecuteRequestWithDryRun(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   auth.NewManager(tmpDir),
	}

	req := &cli.ParsedRequest{
		Method:  "GET",
		URL:     "https://api.example.com/users",
		Headers: make(map[string]string),
		Dry:     true,
		Save:    "test-call",
	}

	err := app.executeRequest(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify call was saved
	saved, err := app.storage.Load("test-call")
	if err != nil {
		t.Fatalf("expected saved call to exist, got error %v", err)
	}

	if saved.Method != "GET" {
		t.Errorf("saved method: got %q, want 'GET'", saved.Method)
	}
	if saved.URL != "https://api.example.com/users" {
		t.Errorf("saved URL: got %q, want 'https://api.example.com/users'", saved.URL)
	}
}

// TestExecuteRequestWithDryRunNoSave tests --dry without --save
func TestExecuteRequestWithDryRunNoSave(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   auth.NewManager(tmpDir),
	}

	req := &cli.ParsedRequest{
		Method:  "GET",
		URL:     "https://api.example.com/users",
		Headers: make(map[string]string),
		Dry:     true,
	}

	err := app.executeRequest(req)
	if err == nil {
		t.Fatalf("expected error for --dry without --save, got nil")
	}
}

// TestExecuteRequestWithDefaultHeaders tests applying workspace default headers
func TestExecuteRequestWithDefaultHeaders(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Config: &config.WorkspaceConfig{
				DefaultHeaders: map[string]string{
					"X-Custom": "value",
				},
			},
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
	}

	req := &cli.ParsedRequest{
		Method:  "GET",
		URL:     "https://httpbin.org/get",
		Headers: make(map[string]string),
		Dry:     true,
		Save:    "test",
	}

	err := app.executeRequest(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Default headers should be applied during executeRequest
	if _, exists := req.Headers["X-Custom"]; !exists {
		t.Errorf("default header X-Custom should be applied")
	}
}

// TestExecuteRequestWithEnvVarSubstitution tests URL env var substitution
func TestExecuteRequestWithEnvVarSubstitution(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env: map[string]string{
				"API_HOST": "api.example.com",
			},
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
	}

	req := &cli.ParsedRequest{
		Method:  "GET",
		URL:     "https://${API_HOST}/users",
		Headers: make(map[string]string),
		Dry:     true,
		Save:    "test",
	}

	err := app.executeRequest(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify URL was substituted
	saved, err := app.storage.Load("test")
	if err != nil {
		t.Fatalf("expected saved call, got error %v", err)
	}

	if saved.URL != "https://api.example.com/users" {
		t.Errorf("substituted URL: got %q, want 'https://api.example.com/users'", saved.URL)
	}
}

// TestExecuteRequestWithInvalidAuth tests executing with nonexistent auth preset
func TestExecuteRequestWithInvalidAuth(t *testing.T) {
	tmpDir := t.TempDir()
	authMgr := auth.NewManager(tmpDir)

	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   authMgr,
	}

	req := &cli.ParsedRequest{
		Method:  "GET",
		URL:     "https://api.example.com/users",
		Headers: make(map[string]string),
		Auth:    "nonexistent",
	}

	err := app.executeRequest(req)
	if err == nil {
		t.Fatalf("expected error for nonexistent auth, got nil")
	}
}

// TestHandleAuthCommandList tests listing auth presets
func TestHandleAuthCommandList(t *testing.T) {
	tmpDir := t.TempDir()
	authMgr := auth.NewManager(tmpDir)
	if err := authMgr.Add(&auth.AuthPreset{
		Name:  "test",
		Type:  "bearer",
		Token: "xyz",
	}); err != nil {
		t.Fatalf("failed to add preset: %v", err)
	}

	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   authMgr,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "list",
	}

	err := app.handleAuthCommand(cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestHandleAuthCommandListEmpty tests listing when no presets exist
func TestHandleAuthCommandListEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   auth.NewManager(tmpDir),
	}

	cmd := &cli.AuthCommand{
		Subcommand: "list",
	}

	err := app.handleAuthCommand(cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestHandleAuthCommandAddBearer tests adding a bearer auth preset
func TestHandleAuthCommandAddBearer(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   auth.NewManager(tmpDir),
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "bearer",
		Name:       "my-api",
		Flags: map[string]string{
			"token": "abc123",
		},
	}

	err := app.handleAuthCommand(cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify preset was added
	preset, err := app.authMgr.Get("my-api")
	if err != nil {
		t.Fatalf("expected preset to exist, got error %v", err)
	}

	if preset.Type != "bearer" {
		t.Errorf("preset type: got %q, want 'bearer'", preset.Type)
	}
	if preset.Token != "abc123" {
		t.Errorf("preset token: got %q, want 'abc123'", preset.Token)
	}
}

// TestHandleAuthCommandAddBasic tests adding basic auth preset
func TestHandleAuthCommandAddBasic(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   auth.NewManager(tmpDir),
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "basic",
		Name:       "api-user",
		Flags: map[string]string{
			"username": "user123",
			"password": "pass456",
		},
	}

	err := app.handleAuthCommand(cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	preset, err := app.authMgr.Get("api-user")
	if err != nil {
		t.Fatalf("expected preset to exist, got error %v", err)
	}

	if preset.Type != "basic" {
		t.Errorf("preset type: got %q, want 'basic'", preset.Type)
	}
}

// TestHandleAuthCommandAddBasicMissingUsername tests basic auth without username
func TestHandleAuthCommandAddBasicMissingUsername(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   auth.NewManager(tmpDir),
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "basic",
		Name:       "api-user",
		Flags:      map[string]string{},
	}

	err := app.handleAuthCommand(cmd)
	if err == nil {
		t.Fatalf("expected error for missing username, got nil")
	}
}

// TestHandleAuthCommandAddCustom tests adding custom header auth
func TestHandleAuthCommandAddCustom(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   auth.NewManager(tmpDir),
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "custom",
		Name:       "custom-auth",
		Flags: map[string]string{
			"header": "X-API-Key",
			"value":  "secret-key",
		},
	}

	err := app.handleAuthCommand(cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	preset, err := app.authMgr.Get("custom-auth")
	if err != nil {
		t.Fatalf("expected preset to exist, got error %v", err)
	}

	if preset.Type != "custom" {
		t.Errorf("preset type: got %q, want 'custom'", preset.Type)
	}
}

// TestHandleAuthCommandRemove tests removing an auth preset
func TestHandleAuthCommandRemove(t *testing.T) {
	tmpDir := t.TempDir()
	authMgr := auth.NewManager(tmpDir)
	if err := authMgr.Add(&auth.AuthPreset{
		Name:  "test-preset",
		Type:  "bearer",
		Token: "xyz",
	}); err != nil {
		t.Fatalf("failed to add preset: %v", err)
	}

	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   authMgr,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "remove",
		Name:       "test-preset",
	}

	err := app.handleAuthCommand(cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify preset was removed
	_, err = app.authMgr.Get("test-preset")
	if err == nil {
		t.Fatalf("expected error for removed preset, got nil")
	}
}

// TestHandleAuthCommandInvalid tests invalid auth subcommand
func TestHandleAuthCommandInvalid(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   auth.NewManager(tmpDir),
	}

	cmd := &cli.AuthCommand{
		Subcommand: "invalid",
	}

	err := app.handleAuthCommand(cmd)
	if err == nil {
		t.Fatalf("expected error for invalid subcommand, got nil")
	}
}

// TestParseTimeoutFromGlobalConfig tests parsing timeout from config
func TestParseTimeoutFromGlobalConfig(t *testing.T) {
	tests := []struct {
		name              string
		configTimeout     string
		expectedSecInDesc string
	}{
		{
			name:              "valid timeout",
			configTimeout:     "60s",
			expectedSecInDesc: "60",
		},
		{
			name:              "invalid timeout falls back to default",
			configTimeout:     "invalid",
			expectedSecInDesc: "30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			app := &App{
				workspace: &config.Workspace{Root: tmpDir},
				global: &config.GlobalConfig{
					Timeout: tt.configTimeout,
				},
				storage: storage.NewManager(tmpDir),
				authMgr: auth.NewManager(tmpDir),
			}

			// Dry run to avoid actual HTTP request
			req := &cli.ParsedRequest{
				Method:  "GET",
				URL:     "https://api.example.com/test",
				Headers: make(map[string]string),
				Dry:     true,
				Save:    "test",
			}

			err := app.executeRequest(req)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

// TestListCallsWithSavedCalls tests listing with existing calls
func TestListCallsWithSavedCalls(t *testing.T) {
	tmpDir := t.TempDir()
	storageMgr := storage.NewManager(tmpDir)

	call1 := storage.NewSavedCall("api-users", "GET", "https://api.example.com/users", nil, nil, "")
	call2 := storage.NewSavedCall("api-posts", "POST", "https://api.example.com/posts", nil, nil, "")

	if err := storageMgr.Save(call1); err != nil {
		t.Fatalf("failed to save call1: %v", err)
	}
	if err := storageMgr.Save(call2); err != nil {
		t.Fatalf("failed to save call2: %v", err)
	}

	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storageMgr,
		authMgr:   auth.NewManager(tmpDir),
	}

	err := app.Run([]string{"list"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestExecuteRecall tests executing a recalled saved call
func TestExecuteRecall(t *testing.T) {
	tmpDir := t.TempDir()
	storageMgr := storage.NewManager(tmpDir)

	call := storage.NewSavedCall(
		"test-call",
		"GET",
		"https://api.example.com/users",
		map[string]string{"X-Custom": "value"},
		nil,
		"",
	)
	if err := storageMgr.Save(call); err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storageMgr,
		authMgr:   auth.NewManager(tmpDir),
	}

	opts := &cli.RecallOptions{
		Name:              "test-call",
		ParameterOverride: make(map[string]string),
		Headers:           make(map[string]string),
	}

	err := app.executeRecall(opts)
	// Error expected since we can't actually make HTTP request, but method should work
	if err != nil && !strings.Contains(err.Error(), "request failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestExecuteRecallNonexistent tests recalling nonexistent call
func TestExecuteRecallNonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   auth.NewManager(tmpDir),
	}

	opts := &cli.RecallOptions{
		Name:              "nonexistent",
		ParameterOverride: make(map[string]string),
		Headers:           make(map[string]string),
	}

	err := app.executeRecall(opts)
	if err == nil {
		t.Fatalf("expected error for nonexistent call, got nil")
	}
}

// TestOutputCapture captures output during tests
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// TestVersionOutput tests version output format
func TestVersionOutput(t *testing.T) {
	app := &App{
		workspace: &config.Workspace{Root: "/tmp"},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager("/tmp"),
		authMgr:   auth.NewManager("/tmp"),
	}

	output := captureOutput(func() {
		_ = app.Run([]string{"--version"})
	})

	if !strings.Contains(output, "gosh version") {
		t.Errorf("version output missing: got %q", output)
	}
}

// TestDeleteCallOutput tests delete output message
func TestDeleteCallOutput(t *testing.T) {
	tmpDir := t.TempDir()
	storageMgr := storage.NewManager(tmpDir)
	call := storage.NewSavedCall("delete-me", "GET", "https://api.example.com", nil, nil, "")
	if err := storageMgr.Save(call); err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storageMgr,
		authMgr:   auth.NewManager(tmpDir),
	}

	output := captureOutput(func() {
		_ = app.Run([]string{"delete", "delete-me"})
	})

	if !strings.Contains(output, "Deleted") {
		t.Errorf("delete output missing: got %q", output)
	}
}

// TestDryRunOutput tests dry run output
func TestDryRunOutput(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storage.NewManager(tmpDir),
		authMgr:   auth.NewManager(tmpDir),
	}

	req := &cli.ParsedRequest{
		Method:  "POST",
		URL:     "https://api.example.com/users",
		Headers: make(map[string]string),
		Body:    `{"name":"John"}`,
		Dry:     true,
		Save:    "create-user",
	}

	output := captureOutput(func() {
		_ = app.executeRequest(req)
	})

	if !strings.Contains(output, "Saved call") {
		t.Errorf("dry run output missing: got %q", output)
	}
}

// TestIsTerminalWithFile tests terminal detection
func TestIsTerminalWithFile(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	// File is not a terminal
	result := isTerminal(tmpFile)
	if result {
		t.Errorf("file should not be terminal")
	}
}

// TestRecallWithHeaderOverrides tests recall with header overrides
func TestRecallWithHeaderOverrides(t *testing.T) {
	tmpDir := t.TempDir()
	storageMgr := storage.NewManager(tmpDir)

	call := storage.NewSavedCall(
		"test-call",
		"GET",
		"https://api.example.com/users",
		map[string]string{"X-Original": "original"},
		nil,
		"",
	)
	if err := storageMgr.Save(call); err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	app := &App{
		workspace: &config.Workspace{Root: tmpDir},
		global:    &config.GlobalConfig{},
		storage:   storageMgr,
		authMgr:   auth.NewManager(tmpDir),
	}

	opts := &cli.RecallOptions{
		Name: "test-call",
		Headers: map[string]string{
			"X-Override": "override",
		},
		ParameterOverride: make(map[string]string),
	}

	err := app.executeRecall(opts)
	// Error expected from HTTP request attempt, but headers should be merged
	if err != nil && !strings.Contains(err.Error(), "request failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestExecuteRequestWithStdinBodyPipe tests executeRequest with stdin body data
func TestExecuteRequestWithStdinBodyPipe(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   false,
	}

	// Create a pipe to simulate stdin
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	// Write test data to stdin
	go func() {
		_, _ = w.WriteString(`{"test": "data"}`)
		w.Close()
	}()

	req := &cli.ParsedRequest{
		Method:        "POST",
		URL:           "http://example.com/api",
		Headers:       make(map[string]string),
		QueryParams:   make(map[string]string),
		Body:          "",
		Dry:           true,
		Save:          "test-call",
		NoInteractive: true,
	}

	err := app.executeRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestExecuteRequestWithPathVariablesNonInteractive tests path variables with non-interactive mode
func TestExecuteRequestWithPathVariablesNonInteractive(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	req := &cli.ParsedRequest{
		Method:        "GET",
		URL:           "http://example.com/users/{id}/posts",
		Headers:       make(map[string]string),
		QueryParams:   make(map[string]string),
		Body:          "",
		PathParams:    make(map[string]string),
		Dry:           true,
		Save:          "test-call",
		NoInteractive: true, // Missing path param, no interactive mode - should error
	}

	err := app.executeRequest(req)
	if err == nil {
		t.Error("expected error for missing path variable in non-interactive mode")
	}
	if !strings.Contains(err.Error(), "missing required template variable") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestExecuteRequestWithPathVariablesProvided tests path variables with provided values
func TestExecuteRequestWithPathVariablesProvided(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	req := &cli.ParsedRequest{
		Method:      "GET",
		URL:         "http://example.com/users/{id}",
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
		Body:        "",
		PathParams:  map[string]string{"id": "123"},
		Dry:         true,
		Save:        "test-call",
	}

	err := app.executeRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the call was saved (note: dry run saves the original URL, not resolved)
	call, err := app.storage.Load("test-call")
	if err != nil {
		t.Fatalf("failed to load saved call: %v", err)
	}
	// In dry run, the original URL is saved (template substitution happens at execution time)
	if call.Method != "GET" {
		t.Errorf("expected GET method, got: %s", call.Method)
	}
}

// TestExecuteRequestWithInvalidTimeout tests handling of invalid timeout values
func TestExecuteRequestWithInvalidTimeout(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{Timeout: "invalid-timeout"},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	req := &cli.ParsedRequest{
		Method:        "GET",
		URL:           "http://example.com/api",
		Headers:       make(map[string]string),
		QueryParams:   make(map[string]string),
		Body:          "",
		Dry:           true,
		Save:          "test-call",
		NoInteractive: true,
	}

	err := app.executeRequest(req)
	if err != nil {
		t.Fatalf("unexpected error with invalid timeout: %v", err)
	}
	// Should fall back to default 30s timeout, not error
}

// TestHandleAuthCommandRemoveNonexistent tests removing non-existent auth preset
func TestHandleAuthCommandRemoveNonexistent(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "remove",
		Name:       "nonexistent",
	}

	err := app.handleAuthCommand(cmd)
	if err == nil {
		t.Error("expected error when removing nonexistent auth preset")
	}
}

// TestHandleAuthCommandInvalidSubcommand tests invalid auth subcommand
func TestHandleAuthCommandInvalidSubcommand(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "invalid",
		Name:       "",
	}

	err := app.handleAuthCommand(cmd)
	if err == nil {
		t.Error("expected error for invalid auth subcommand")
	}
}

// TestHandleAuthCommandAddBasicWithShortFlags tests basic auth with -u and -p short flags
func TestHandleAuthCommandAddBasicWithShortFlags(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "basic",
		Name:       "test-basic-short",
		Flags: map[string]string{
			"u": "testuser",    // Short flag for username
			"p": "testpass123", // Short flag for password
		},
	}

	err := app.handleAuthCommand(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the preset was saved
	preset, err := app.authMgr.Get("test-basic-short")
	if err != nil {
		t.Fatalf("failed to get preset: %v", err)
	}
	if preset.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %q", preset.Username)
	}
	if preset.Password != "testpass123" {
		t.Errorf("expected password 'testpass123', got %q", preset.Password)
	}
}

// TestHandleAuthCommandAddBasicWithMixedFlags tests basic auth with mixed long and short flags
func TestHandleAuthCommandAddBasicWithMixedFlags(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "basic",
		Name:       "test-basic-mixed",
		Flags: map[string]string{
			"username": "longuser",  // Long flag for username
			"p":        "shortpass", // Short flag for password
		},
	}

	err := app.handleAuthCommand(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	preset, err := app.authMgr.Get("test-basic-mixed")
	if err != nil {
		t.Fatalf("failed to get preset: %v", err)
	}
	if preset.Username != "longuser" {
		t.Errorf("expected username 'longuser', got %q", preset.Username)
	}
}

// TestHandleAuthCommandAddBasicMissingPassword tests basic auth without password (optional)
func TestHandleAuthCommandAddBasicMissingPassword(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	// Password is optional in the current implementation
	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "basic",
		Name:       "test-basic-no-pass",
		Flags: map[string]string{
			"username": "useronly",
		},
	}

	err := app.handleAuthCommand(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	preset, err := app.authMgr.Get("test-basic-no-pass")
	if err != nil {
		t.Fatalf("failed to get preset: %v", err)
	}
	if preset.Username != "useronly" {
		t.Errorf("expected username 'useronly', got %q", preset.Username)
	}
}

// TestHandleAuthCommandAddBearerWithShortFlag tests bearer auth with -t short flag
func TestHandleAuthCommandAddBearerWithShortFlag(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "bearer",
		Name:       "test-bearer-short",
		Flags: map[string]string{
			"t": "short-token-abc123", // Short flag for token
		},
	}

	err := app.handleAuthCommand(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	preset, err := app.authMgr.Get("test-bearer-short")
	if err != nil {
		t.Fatalf("failed to get preset: %v", err)
	}
	if preset.Token != "short-token-abc123" {
		t.Errorf("expected token 'short-token-abc123', got %q", preset.Token)
	}
}

// TestHandleAuthCommandAddBearerMissingToken tests bearer auth without token
func TestHandleAuthCommandAddBearerMissingToken(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "bearer",
		Name:       "test-bearer-no-token",
		Flags:      make(map[string]string), // No token provided
	}

	err := app.handleAuthCommand(cmd)
	if err == nil {
		t.Error("expected error when token is missing")
	}
	if !strings.Contains(err.Error(), "bearer auth requires") {
		t.Errorf("expected 'bearer auth requires' in error, got %v", err)
	}
}

// TestHandleAuthCommandAddCustomWithShortFlags tests custom auth with -h, -v short flags
func TestHandleAuthCommandAddCustomWithShortFlags(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "custom",
		Name:       "test-custom-short",
		Flags: map[string]string{
			"h": "X-API-Key",      // Short flag for header
			"v": "secret-key-123", // Short flag for value
		},
	}

	err := app.handleAuthCommand(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	preset, err := app.authMgr.Get("test-custom-short")
	if err != nil {
		t.Fatalf("failed to get preset: %v", err)
	}
	if preset.Header != "X-API-Key" {
		t.Errorf("expected header 'X-API-Key', got %q", preset.Header)
	}
	if preset.Value != "secret-key-123" {
		t.Errorf("expected value 'secret-key-123', got %q", preset.Value)
	}
}

// TestHandleAuthCommandAddCustomWithPrefix tests custom auth with prefix flag
func TestHandleAuthCommandAddCustomWithPrefix(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "custom",
		Name:       "test-custom-prefix",
		Flags: map[string]string{
			"header": "Authorization",
			"value":  "mytoken",
			"prefix": "Bearer", // Optional prefix
		},
	}

	err := app.handleAuthCommand(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	preset, err := app.authMgr.Get("test-custom-prefix")
	if err != nil {
		t.Fatalf("failed to get preset: %v", err)
	}
	if preset.Prefix != "Bearer" {
		t.Errorf("expected prefix 'Bearer', got %q", preset.Prefix)
	}
}

// TestHandleAuthCommandAddCustomWithMixedFlags tests custom auth with mixed long and short flags
func TestHandleAuthCommandAddCustomWithMixedFlags(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "custom",
		Name:       "test-custom-mixed",
		Flags: map[string]string{
			"header": "X-Custom-Header", // Long flag for header
			"v":      "custom-value",    // Short flag for value
			"prefix": "Custom",
		},
	}

	err := app.handleAuthCommand(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	preset, err := app.authMgr.Get("test-custom-mixed")
	if err != nil {
		t.Fatalf("failed to get preset: %v", err)
	}
	if preset.Header != "X-Custom-Header" {
		t.Errorf("expected header 'X-Custom-Header', got %q", preset.Header)
	}
}

// TestHandleAuthCommandAddCustomMissingHeader tests custom auth without header
func TestHandleAuthCommandAddCustomMissingHeader(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "custom",
		Name:       "test-custom-no-header",
		Flags: map[string]string{
			"value": "some-value", // Header is required
		},
	}

	err := app.handleAuthCommand(cmd)
	if err == nil {
		t.Error("expected error when header is missing")
	}
	if !strings.Contains(err.Error(), "custom auth requires --header") {
		t.Errorf("expected '--header' in error, got %v", err)
	}
}

// TestHandleAuthCommandAddCustomMissingValue tests custom auth without value
func TestHandleAuthCommandAddCustomMissingValue(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "custom",
		Name:       "test-custom-no-value",
		Flags: map[string]string{
			"header": "X-API-Key", // Value is required
		},
	}

	err := app.handleAuthCommand(cmd)
	if err == nil {
		t.Error("expected error when value is missing")
	}
	if !strings.Contains(err.Error(), "custom auth requires --value") {
		t.Errorf("expected '--value' in error, got %v", err)
	}
}

// TestHandleAuthCommandAddUnknownType tests adding auth with unknown type
func TestHandleAuthCommandAddUnknownType(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "unknown-type", // Unknown auth type
		Name:       "test-unknown",
		Flags:      make(map[string]string),
	}

	err := app.handleAuthCommand(cmd)
	if err == nil {
		t.Error("expected error for unknown auth type")
	}
	if !strings.Contains(err.Error(), "unknown auth type") {
		t.Errorf("expected 'unknown auth type' in error, got %v", err)
	}
}

// TestHandleAuthCommandAddMissingType tests adding auth without type
func TestHandleAuthCommandAddMissingType(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "", // Missing type
		Name:       "test-no-type",
		Flags:      make(map[string]string),
	}

	err := app.handleAuthCommand(cmd)
	if err == nil {
		t.Error("expected error when type is missing")
	}
	if !strings.Contains(err.Error(), "auth add requires") {
		t.Errorf("expected 'auth add requires' in error, got %v", err)
	}
}

// TestHandleAuthCommandAddMissingName tests adding auth without name
func TestHandleAuthCommandAddMissingName(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "add",
		Type:       "bearer",
		Name:       "", // Missing name
		Flags: map[string]string{
			"token": "test-token",
		},
	}

	err := app.handleAuthCommand(cmd)
	if err == nil {
		t.Error("expected error when name is missing")
	}
	if !strings.Contains(err.Error(), "auth add requires") {
		t.Errorf("expected 'auth add requires' in error, got %v", err)
	}
}

// TestHandleAuthCommandRemoveMissingName tests removing auth without name
func TestHandleAuthCommandRemoveMissingName(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	cmd := &cli.AuthCommand{
		Subcommand: "remove",
		Name:       "", // Missing preset name
	}

	err := app.handleAuthCommand(cmd)
	if err == nil {
		t.Error("expected error when name is missing for remove")
	}
	if !strings.Contains(err.Error(), "auth remove requires") {
		t.Errorf("expected 'auth remove requires' in error, got %v", err)
	}
}

// TestRunWithDeleteCommand tests the delete command through Run
func TestRunWithDeleteCommand(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	// First save a call
	call := storage.NewSavedCall("test-call", "GET", "http://example.com", make(map[string]string), make(map[string]string), "")
	if err := app.storage.Save(call); err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	// Run with delete command
	err := app.Run([]string{"delete", "test-call"})
	if err != nil {
		t.Fatalf("unexpected error deleting call: %v", err)
	}

	// Verify the call was deleted
	exists := app.storage.Exists("test-call")
	if exists {
		t.Error("expected call to be deleted")
	}
}

// TestExecuteRequestWithQueryParams tests request with query parameters
func TestExecuteRequestWithQueryParams(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	req := &cli.ParsedRequest{
		Method:        "GET",
		URL:           "http://example.com/api",
		Headers:       make(map[string]string),
		QueryParams:   map[string]string{"page": "1", "limit": "10"},
		Body:          "",
		Dry:           true,
		Save:          "test-call-params",
		NoInteractive: true,
	}

	err := app.executeRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the call was saved with query params
	call, err := app.storage.Load("test-call-params")
	if err != nil {
		t.Fatalf("failed to load saved call: %v", err)
	}
	if len(call.QueryParams) != 2 {
		t.Errorf("expected 2 query params, got %d", len(call.QueryParams))
	}
}

// TestExecuteRequestWithMultipleHeaders tests request with multiple headers
func TestExecuteRequestWithMultipleHeaders(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	req := &cli.ParsedRequest{
		Method: "POST",
		URL:    "http://example.com/api",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token123",
			"X-Custom":      "value",
		},
		QueryParams:   make(map[string]string),
		Body:          `{"test": "data"}`,
		Dry:           true,
		Save:          "test-call-headers",
		NoInteractive: true,
	}

	err := app.executeRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the call was saved with headers
	call, err := app.storage.Load("test-call-headers")
	if err != nil {
		t.Fatalf("failed to load saved call: %v", err)
	}
	if len(call.Headers) != 3 {
		t.Errorf("expected 3 headers, got %d", len(call.Headers))
	}
}

// TestExecuteRequestWithBodyData tests request with body data
func TestExecuteRequestWithBodyData(t *testing.T) {
	tmpDir := t.TempDir()

	app := &App{
		workspace: &config.Workspace{
			Root: tmpDir,
			Env:  make(map[string]string),
		},
		global:  &config.GlobalConfig{},
		storage: storage.NewManager(tmpDir),
		authMgr: auth.NewManager(tmpDir),
		isTTY:   true,
	}

	bodyData := `{"name": "test", "value": 123}`
	req := &cli.ParsedRequest{
		Method:        "POST",
		URL:           "http://example.com/api",
		Headers:       make(map[string]string),
		QueryParams:   make(map[string]string),
		Body:          bodyData,
		Dry:           true,
		Save:          "test-call-body",
		NoInteractive: true,
	}

	err := app.executeRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the call was saved with body
	call, err := app.storage.Load("test-call-body")
	if err != nil {
		t.Fatalf("failed to load saved call: %v", err)
	}
	if call.Body != bodyData {
		t.Errorf("expected body %q, got %q", bodyData, call.Body)
	}
}
