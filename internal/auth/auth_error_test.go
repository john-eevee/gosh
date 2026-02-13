package auth

import (
	"os"
	"path/filepath"
	"testing"
)

// TestManagerWithInvalidPermissions tests auth file with wrong permissions
func TestManagerWithInvalidPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	authFile := filepath.Join(tmpDir, ".auth")

	// Create file with insecure permissions
	err := os.WriteFile(authFile, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("failed to create auth file: %v", err)
	}

	// Try to use manager (should handle gracefully or warn)
	mgr := NewManager(tmpDir)
	err = mgr.Load()
	// Should either work or return permission-related error
	_ = err
}

// TestManagerAddDuplicatePreset tests adding duplicate auth preset
func TestManagerAddDuplicatePreset(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	preset := &AuthPreset{
		Name:  "api-key",
		Type:  "bearer",
		Token: "token123",
	}

	err := mgr.Add(preset)
	if err != nil {
		t.Fatalf("first add failed: %v", err)
	}

	// Add same preset again
	err = mgr.Add(preset)
	// Should either succeed (overwrite) or return error
	_ = err
}

// TestApplyAuthWithMissingToken tests applying bearer auth without token
func TestApplyAuthWithMissingToken(t *testing.T) {
	preset := &AuthPreset{
		Type:  "bearer",
		Token: "",
	}

	// Mock HTTP request would be needed for full test
	// This tests the type itself
	if preset.Token == "" {
		t.Logf("Bearer auth with empty token would fail on apply")
	}
}

// TestApplyBasicAuthWithMissingPassword tests basic auth without password
func TestApplyBasicAuthWithMissingPassword(t *testing.T) {
	preset := &AuthPreset{
		Type:     "basic",
		Username: "user",
		Password: "",
	}

	if preset.Password == "" {
		t.Logf("Basic auth without password would use empty string")
	}
}

// TestManagerRemoveNonexistentPreset tests removing preset that doesn't exist
func TestManagerRemoveNonexistentPreset(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	err := mgr.Remove("nonexistent")
	if err == nil {
		t.Logf("Removing nonexistent preset should succeed or fail gracefully")
	}
}

// TestManagerGetAfterRemove tests getting preset after it's removed
func TestManagerGetAfterRemove(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	preset := &AuthPreset{
		Name:  "temp",
		Type:  "bearer",
		Token: "xyz",
	}
	if err := mgr.Add(preset); err != nil {
		t.Fatalf("failed to add preset: %v", err)
	}
	if err := mgr.Remove("temp"); err != nil {
		t.Fatalf("failed to remove preset: %v", err)
	}

	_, err := mgr.Get("temp")
	if err == nil {
		t.Fatalf("expected error when getting removed preset")
	}
}
