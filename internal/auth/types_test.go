package auth

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	preset := &AuthPreset{
		Name:     "test-basic",
		Type:     "basic",
		Username: "john",
		Password: "secret123",
	}

	err := preset.Apply(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	auth := req.Header.Get("Authorization")
	if auth == "" {
		t.Fatal("Authorization header not set")
	}

	if auth != "Basic am9objpzZWNyZXQxMjM=" {
		t.Errorf("expected 'Basic am9objpzZWNyZXQxMjM=', got '%s'", auth)
	}
}

func TestBasicAuthWithoutPassword(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	preset := &AuthPreset{
		Name:     "test-basic",
		Type:     "basic",
		Username: "john",
		Password: "",
	}

	err := preset.Apply(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	auth := req.Header.Get("Authorization")
	if auth != "Basic am9objo=" {
		t.Errorf("expected 'Basic am9objo=', got '%s'", auth)
	}
}

func TestBasicAuthMissingUsername(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	preset := &AuthPreset{
		Name:     "test-basic",
		Type:     "basic",
		Username: "",
		Password: "secret",
	}

	err := preset.Apply(req)
	if err == nil {
		t.Fatal("expected error for missing username")
	}

	if err.Error() != "basic auth requires username" {
		t.Errorf("expected 'basic auth requires username', got '%v'", err)
	}
}

func TestBearerAuth(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	preset := &AuthPreset{
		Name:  "test-bearer",
		Type:  "bearer",
		Token: "my-secret-token-123",
	}

	err := preset.Apply(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	auth := req.Header.Get("Authorization")
	if auth != "Bearer my-secret-token-123" {
		t.Errorf("expected 'Bearer my-secret-token-123', got '%s'", auth)
	}
}

func TestBearerAuthMissingToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	preset := &AuthPreset{
		Name:  "test-bearer",
		Type:  "bearer",
		Token: "",
	}

	err := preset.Apply(req)
	if err == nil {
		t.Fatal("expected error for missing token")
	}

	if err.Error() != "bearer auth requires token" {
		t.Errorf("expected 'bearer auth requires token', got '%v'", err)
	}
}

func TestCustomAuth(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	preset := &AuthPreset{
		Name:   "test-custom",
		Type:   "custom",
		Header: "X-API-Key",
		Value:  "my-api-key",
		Prefix: "",
	}

	err := preset.Apply(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	auth := req.Header.Get("X-API-Key")
	if auth != "my-api-key" {
		t.Errorf("expected 'my-api-key', got '%s'", auth)
	}
}

func TestCustomAuthWithPrefix(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	preset := &AuthPreset{
		Name:   "test-custom",
		Type:   "custom",
		Header: "Authorization",
		Value:  "my-token",
		Prefix: "Token ",
	}

	err := preset.Apply(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	auth := req.Header.Get("Authorization")
	if auth != "Token my-token" {
		t.Errorf("expected 'Token my-token', got '%s'", auth)
	}
}

func TestCustomAuthMissingHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	preset := &AuthPreset{
		Name:   "test-custom",
		Type:   "custom",
		Header: "",
		Value:  "my-value",
	}

	err := preset.Apply(req)
	if err == nil {
		t.Fatal("expected error for missing header name")
	}

	if err.Error() != "custom auth requires header name" {
		t.Errorf("expected 'custom auth requires header name', got '%v'", err)
	}
}

func TestUnknownAuthType(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	preset := &AuthPreset{
		Name: "test-unknown",
		Type: "unknown",
	}

	err := preset.Apply(req)
	if err == nil {
		t.Fatal("expected error for unknown auth type")
	}

	if err.Error() != "unknown auth type: unknown" {
		t.Errorf("expected 'unknown auth type: unknown', got '%v'", err)
	}
}

func TestAuthTypeCase(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	// Test uppercase
	preset := &AuthPreset{
		Name:  "test-bearer",
		Type:  "BEARER",
		Token: "test-token",
	}

	err := preset.Apply(req)
	if err != nil {
		t.Fatalf("expected no error for uppercase type, got %v", err)
	}

	if req.Header.Get("Authorization") != "Bearer test-token" {
		t.Error("bearer auth failed with uppercase type")
	}
}

func TestManagerAdd(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	preset := &AuthPreset{
		Name:     "test-basic",
		Type:     "basic",
		Username: "john",
		Password: "secret",
	}

	err := m.Add(preset)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify saved to disk
	retrieved, err := m.Get("test-basic")
	if err != nil {
		t.Fatalf("expected no error retrieving preset, got %v", err)
	}

	if retrieved.Username != "john" {
		t.Errorf("expected username 'john', got '%s'", retrieved.Username)
	}
}

func TestManagerAddEmptyName(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	preset := &AuthPreset{
		Name:     "",
		Type:     "basic",
		Username: "john",
	}

	err := m.Add(preset)
	if err == nil {
		t.Fatal("expected error for empty name")
	}

	if err.Error() != "preset name cannot be empty" {
		t.Errorf("expected 'preset name cannot be empty', got '%v'", err)
	}
}

func TestManagerAddEmptyType(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	preset := &AuthPreset{
		Name:     "test",
		Type:     "",
		Username: "john",
	}

	err := m.Add(preset)
	if err == nil {
		t.Fatal("expected error for empty type")
	}

	if err.Error() != "preset type cannot be empty" {
		t.Errorf("expected 'preset type cannot be empty', got '%v'", err)
	}
}

func TestManagerRemove(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	preset := &AuthPreset{
		Name:     "test-basic",
		Type:     "basic",
		Username: "john",
	}

	if err := m.Add(preset); err != nil {
		t.Fatalf("expected no error adding preset: %v", err)
	}

	err := m.Remove("test-basic")
	if err != nil {
		t.Fatalf("expected no error removing preset, got %v", err)
	}

	_, err = m.Get("test-basic")
	if err == nil {
		t.Fatal("expected error getting removed preset")
	}
}

func TestManagerRemoveNonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	err := m.Remove("nonexistent")
	if err == nil {
		t.Fatal("expected error removing nonexistent preset")
	}

	if err.Error() != "auth preset not found: nonexistent" {
		t.Errorf("expected 'auth preset not found: nonexistent', got '%v'", err)
	}
}

func TestManagerGetNonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	_, err := m.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error getting nonexistent preset")
	}

	if err.Error() != "auth preset not found: nonexistent" {
		t.Errorf("expected 'auth preset not found: nonexistent', got '%v'", err)
	}
}

func TestManagerLoadFromDisk(t *testing.T) {
	tmpDir := t.TempDir()

	// Create initial manager and add preset
	m1 := NewManager(tmpDir)
	preset := &AuthPreset{
		Name:  "test-bearer",
		Type:  "bearer",
		Token: "my-token",
	}
	if err := m1.Add(preset); err != nil {
		t.Fatalf("failed to add preset: %v", err)
	}

	// Create new manager and load from disk
	m2 := NewManager(tmpDir)
	err := m2.Load()
	if err != nil {
		t.Fatalf("expected no error loading, got %v", err)
	}

	retrieved, err := m2.Get("test-bearer")
	if err != nil {
		t.Fatalf("expected no error getting preset, got %v", err)
	}

	if retrieved.Token != "my-token" {
		t.Errorf("expected token 'my-token', got '%s'", retrieved.Token)
	}
}

func TestManagerLoadNonexistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	// Should not error if file doesn't exist
	err := m.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}

	// Presets should be empty
	presets := m.List()
	if len(presets) != 0 {
		t.Errorf("expected 0 presets, got %d", len(presets))
	}
}

func TestManagerList(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	p1 := &AuthPreset{Name: "preset1", Type: "basic", Username: "user1"}
	p2 := &AuthPreset{Name: "preset2", Type: "bearer", Token: "token2"}

	if err := m.Add(p1); err != nil {
		t.Fatalf("failed to add preset1: %v", err)
	}
	if err := m.Add(p2); err != nil {
		t.Fatalf("failed to add preset2: %v", err)
	}

	presets := m.List()
	if len(presets) != 2 {
		t.Errorf("expected 2 presets, got %d", len(presets))
	}

	if _, exists := presets["preset1"]; !exists {
		t.Error("preset1 not found in list")
	}

	if _, exists := presets["preset2"]; !exists {
		t.Error("preset2 not found in list")
	}
}

func TestPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	preset := &AuthPreset{
		Name:  "test",
		Type:  "bearer",
		Token: "secret-token",
	}

	if err := m.Add(preset); err != nil {
		t.Fatalf("expected no error adding preset: %v", err)
	}

	// Check that config file has restrictive permissions (0600)
	fileInfo, err := os.Stat(m.configPath)
	if err != nil {
		t.Fatalf("expected auth config file to exist, got %v", err)
	}

	mode := fileInfo.Mode().Perm()
	if mode != 0600 {
		t.Errorf("expected file permissions 0600, got %o", mode)
	}
}

func TestConfigDirCreation(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir)

	preset := &AuthPreset{
		Name:  "test",
		Type:  "bearer",
		Token: "token",
	}

	err := m.Add(preset)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that .gosh directory was created
	goshDir := filepath.Join(tmpDir, ".gosh")
	if _, err := os.Stat(goshDir); os.IsNotExist(err) {
		t.Error("expected .gosh directory to be created")
	}
}
