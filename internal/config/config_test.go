package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLoadGlobalConfigNonexistent tests loading nonexistent global config
func TestLoadGlobalConfigNonexistent(t *testing.T) {
	// Use a temporary directory that doesn't contain a config file
	tempDir := t.TempDir()

	// Save original XDG_CONFIG_HOME
	origXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", origXDG)

	// Set XDG_CONFIG_HOME to temp directory
	os.Setenv("XDG_CONFIG_HOME", tempDir)

	config, err := LoadGlobalConfig()
	if err != nil {
		t.Fatalf("expected no error for nonexistent config, got %v", err)
	}

	if config == nil {
		t.Error("expected non-nil config")
	}
}

// TestLoadGlobalConfigValidFile tests loading valid global config file
func TestLoadGlobalConfigValidFile(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()

	// Save original XDG_CONFIG_HOME
	origXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", origXDG)

	// Create config directory and file
	configDir := filepath.Join(tempDir, "gosh")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	configContent := `defaultEnvironment: prod
prettyPrint: true
timeout: 30s
userAgent: gosh/0.1.1
`

	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	os.Setenv("XDG_CONFIG_HOME", tempDir)

	config, err := LoadGlobalConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if config.DefaultEnvironment != "prod" {
		t.Errorf("expected defaultEnvironment=prod, got %s", config.DefaultEnvironment)
	}
	if !config.PrettyPrint {
		t.Error("expected prettyPrint=true")
	}
	if config.Timeout != "30s" {
		t.Errorf("expected timeout=30s, got %s", config.Timeout)
	}
	if config.UserAgent != "gosh/0.1.1" {
		t.Errorf("expected userAgent=gosh/0.1.1, got %s", config.UserAgent)
	}
}

// TestLoadGlobalConfigInvalidYAML tests loading invalid YAML config
func TestLoadGlobalConfigInvalidYAML(t *testing.T) {
	tempDir := t.TempDir()

	origXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", origXDG)

	configDir := filepath.Join(tempDir, "gosh")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	// Invalid YAML with unclosed quote
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content:\n  - broken"), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	os.Setenv("XDG_CONFIG_HOME", tempDir)

	_, err := LoadGlobalConfig()
	if err == nil {
		t.Error("expected error for invalid YAML, got none")
	}
}

// TestLoadWorkspaceConfig tests loading workspace configuration
func TestLoadWorkspaceConfig(t *testing.T) {
	tempDir := t.TempDir()

	configPath := filepath.Join(tempDir, ".gosh.yaml")
	configContent := `name: myapi
baseUrl: https://api.example.com
defaultHeaders:
  X-API-Key: secret123
environments:
  dev:
    API_URL: http://localhost:3000
  prod:
    API_URL: https://api.example.com
`

	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write workspace config: %v", err)
	}

	config, err := LoadWorkspaceConfig(configPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if config.Name != "myapi" {
		t.Errorf("expected name=myapi, got %s", config.Name)
	}
	if config.BaseURL != "https://api.example.com" {
		t.Errorf("expected baseUrl=https://api.example.com, got %s", config.BaseURL)
	}
	if len(config.DefaultHeaders) != 1 {
		t.Errorf("expected 1 default header, got %d", len(config.DefaultHeaders))
	}
	if config.DefaultHeaders["X-API-Key"] != "secret123" {
		t.Errorf("expected X-API-Key=secret123, got %s", config.DefaultHeaders["X-API-Key"])
	}
	if len(config.Environments) != 2 {
		t.Errorf("expected 2 environments, got %d", len(config.Environments))
	}
}

// TestLoadWorkspaceConfigNonexistent tests loading nonexistent workspace config
func TestLoadWorkspaceConfigNonexistent(t *testing.T) {
	nonexistentPath := "/this/path/does/not/exist/.gosh.yaml"

	_, err := LoadWorkspaceConfig(nonexistentPath)
	if err == nil {
		t.Error("expected error for nonexistent config file")
	}
}

// TestLoadEnvFile tests loading .env file
func TestLoadEnvFile(t *testing.T) {
	tempDir := t.TempDir()

	envPath := filepath.Join(tempDir, ".env")
	envContent := `API_KEY=secret123
API_URL=https://api.example.com
DATABASE_URL=postgres://localhost/mydb
DEBUG=true
# This is a comment
EMPTY_VAR=
`

	if err := os.WriteFile(envPath, []byte(envContent), 0600); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	envVars, err := LoadEnvFile(envPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(envVars) != 5 {
		t.Errorf("expected 5 env vars, got %d", len(envVars))
	}
	if envVars["API_KEY"] != "secret123" {
		t.Errorf("expected API_KEY=secret123, got %s", envVars["API_KEY"])
	}
	if envVars["DATABASE_URL"] != "postgres://localhost/mydb" {
		t.Errorf("expected correct DATABASE_URL, got %s", envVars["DATABASE_URL"])
	}
	if envVars["EMPTY_VAR"] != "" {
		t.Errorf("expected EMPTY_VAR to be empty, got %s", envVars["EMPTY_VAR"])
	}
}

// TestLoadEnvFileWithSpaces tests loading .env file with spaces
func TestLoadEnvFileWithSpaces(t *testing.T) {
	tempDir := t.TempDir()

	envPath := filepath.Join(tempDir, ".env")
	envContent := `SPACED_KEY = spaced_value
	TABS	=	with_tabs
KEY_WITH_EQUALS=value=with=equals
`

	if err := os.WriteFile(envPath, []byte(envContent), 0600); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	envVars, err := LoadEnvFile(envPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if envVars["SPACED_KEY"] != "spaced_value" {
		t.Errorf("expected SPACED_KEY=spaced_value, got %s", envVars["SPACED_KEY"])
	}
	if envVars["TABS"] != "with_tabs" {
		t.Errorf("expected TABS=with_tabs, got %s", envVars["TABS"])
	}
	if envVars["KEY_WITH_EQUALS"] != "value=with=equals" {
		t.Errorf("expected KEY_WITH_EQUALS=value=with=equals, got %s", envVars["KEY_WITH_EQUALS"])
	}
}

// TestLoadEnvFileNonexistent tests loading nonexistent .env file
func TestLoadEnvFileNonexistent(t *testing.T) {
	nonexistentPath := "/this/path/does/not/exist/.env"

	_, err := LoadEnvFile(nonexistentPath)
	if err == nil {
		t.Error("expected error for nonexistent .env file")
	}
}

// TestLoadEnvFileEmpty tests loading empty .env file
func TestLoadEnvFileEmpty(t *testing.T) {
	tempDir := t.TempDir()

	envPath := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(envPath, []byte(""), 0600); err != nil {
		t.Fatalf("failed to write empty env file: %v", err)
	}

	envVars, err := LoadEnvFile(envPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(envVars) != 0 {
		t.Errorf("expected 0 env vars, got %d", len(envVars))
	}
}

// TestLoadEnvFileOnlyComments tests loading .env file with only comments
func TestLoadEnvFileOnlyComments(t *testing.T) {
	tempDir := t.TempDir()

	envPath := filepath.Join(tempDir, ".env")
	envContent := `# Comment 1
# Comment 2

# Comment 3
`

	if err := os.WriteFile(envPath, []byte(envContent), 0600); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	envVars, err := LoadEnvFile(envPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(envVars) != 0 {
		t.Errorf("expected 0 env vars from comments-only file, got %d", len(envVars))
	}
}

// TestDetectWorkspaceWithGoshYAML tests workspace detection with .gosh.yaml
func TestDetectWorkspaceWithGoshYAML(t *testing.T) {
	tempDir := t.TempDir()

	// Save original CWD
	origCWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer func() { _ = os.Chdir(origCWD) }()

	// Change to temp directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Create .gosh.yaml file
	goshPath := filepath.Join(tempDir, ".gosh.yaml")
	goshContent := `name: test-workspace
baseUrl: https://api.test.com
`
	if err := os.WriteFile(goshPath, []byte(goshContent), 0600); err != nil {
		t.Fatalf("failed to write .gosh.yaml: %v", err)
	}

	workspace, err := DetectWorkspace()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if workspace == nil {
		t.Fatal("expected non-nil workspace")
	}
	if workspace.Config == nil {
		t.Error("expected workspace config to be loaded")
	} else if workspace.Config.Name != "test-workspace" {
		t.Errorf("expected workspace name=test-workspace, got %s", workspace.Config.Name)
	}
}

// TestDetectWorkspaceWithEnvFile tests workspace detection with .env file
func TestDetectWorkspaceWithEnvFile(t *testing.T) {
	tempDir := t.TempDir()

	origCWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer func() { _ = os.Chdir(origCWD) }()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Create .env file
	envPath := filepath.Join(tempDir, ".env")
	envContent := `TEST_VAR=test_value
ANOTHER_VAR=another_value
`
	if err := os.WriteFile(envPath, []byte(envContent), 0600); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	workspace, err := DetectWorkspace()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(workspace.Env) != 2 {
		t.Errorf("expected 2 env vars, got %d", len(workspace.Env))
	}
	if workspace.Env["TEST_VAR"] != "test_value" {
		t.Errorf("expected TEST_VAR=test_value, got %s", workspace.Env["TEST_VAR"])
	}
}

// TestDetectWorkspaceWithGitRoot tests workspace detection with .git directory
func TestDetectWorkspaceWithGitRoot(t *testing.T) {
	tempDir := t.TempDir()

	origCWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer func() { _ = os.Chdir(origCWD) }()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Create .git directory
	gitPath := filepath.Join(tempDir, ".git")
	if err := os.MkdirAll(gitPath, 0700); err != nil {
		t.Fatalf("failed to create .git directory: %v", err)
	}

	workspace, err := DetectWorkspace()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if workspace == nil {
		t.Fatal("expected non-nil workspace")
	}
	// Root should be the directory with .git
	if !strings.Contains(workspace.Root, tempDir) {
		t.Errorf("expected workspace root to contain temp dir, got %s", workspace.Root)
	}
}

// TestDetectWorkspaceNestedDirectory tests workspace detection from nested directory
func TestDetectWorkspaceNestedDirectory(t *testing.T) {
	tempDir := t.TempDir()

	origCWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer func() { _ = os.Chdir(origCWD) }()

	// Create nested directory structure
	nestedDir := filepath.Join(tempDir, "sub", "dir")
	if err := os.MkdirAll(nestedDir, 0700); err != nil {
		t.Fatalf("failed to create nested directory: %v", err)
	}

	// Create .gosh.yaml in root
	goshPath := filepath.Join(tempDir, ".gosh.yaml")
	goshContent := `name: root-workspace
`
	if err := os.WriteFile(goshPath, []byte(goshContent), 0600); err != nil {
		t.Fatalf("failed to write .gosh.yaml: %v", err)
	}

	// Change to nested directory
	if err := os.Chdir(nestedDir); err != nil {
		t.Fatalf("failed to change to nested directory: %v", err)
	}

	workspace, err := DetectWorkspace()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if workspace.Config == nil {
		t.Error("expected workspace config to be loaded from parent directory")
	} else if workspace.Config.Name != "root-workspace" {
		t.Errorf("expected to find root workspace config, got name=%s", workspace.Config.Name)
	}
}

// TestDetectWorkspacePreferGoshYAMLOverGit tests that .gosh.yaml is preferred over .git
func TestDetectWorkspacePreferGoshYAMLOverGit(t *testing.T) {
	tempDir := t.TempDir()

	origCWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer func() { _ = os.Chdir(origCWD) }()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Create both .gosh.yaml and .git
	goshPath := filepath.Join(tempDir, ".gosh.yaml")
	goshContent := `name: gosh-workspace
`
	if err := os.WriteFile(goshPath, []byte(goshContent), 0600); err != nil {
		t.Fatalf("failed to write .gosh.yaml: %v", err)
	}

	gitPath := filepath.Join(tempDir, ".git")
	if err := os.MkdirAll(gitPath, 0700); err != nil {
		t.Fatalf("failed to create .git directory: %v", err)
	}

	workspace, err := DetectWorkspace()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should prefer .gosh.yaml
	if workspace.Config == nil {
		t.Error("expected workspace config to be loaded")
	} else if workspace.Config.Name != "gosh-workspace" {
		t.Errorf("expected .gosh.yaml to be preferred, got name=%s", workspace.Config.Name)
	}
}
