package storage

import (
	"testing"
)

// TestManagerSaveWithInvalidPath tests saving when path is invalid
func TestManagerSaveWithInvalidPath(t *testing.T) {
	// Use a path that doesn't exist and can't be created
	mgr := NewManager("/invalid/path/that/does/not/exist")

	call := NewSavedCall("test", "GET", "https://api.example.com", nil, nil, "")
	err := mgr.Save(call)

	if err == nil {
		t.Fatalf("expected error for invalid path, got nil")
	}
}

// TestManagerLoadNonexistentCall tests loading call that doesn't exist
func TestManagerLoadNonexistentCall(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	_, err := mgr.Load("nonexistent-call")
	if err == nil {
		t.Fatalf("expected error for nonexistent call, got nil")
	}
}

// TestManagerDeleteNonexistentCall tests deleting call that doesn't exist
func TestManagerDeleteNonexistentCall(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	err := mgr.Delete("nonexistent-call")
	if err == nil {
		t.Fatalf("expected error for deleting nonexistent call, got nil")
	}
}

// TestManagerListWithEmptyDir tests listing when directory is empty
func TestManagerListWithEmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	calls, err := mgr.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(calls) != 0 {
		t.Errorf("expected empty list, got %d calls", len(calls))
	}
}

// TestManagerExistsWithNonexistentCall tests checking existence of nonexistent call
func TestManagerExistsWithNonexistentCall(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	exists := mgr.Exists("nonexistent")
	if exists {
		t.Fatalf("expected Exists to return false for nonexistent call")
	}
}

// TestManagerExistsWithExistingCall tests checking existence of existing call
func TestManagerExistsWithExistingCall(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	call := NewSavedCall("exists-test", "GET", "https://api.example.com", nil, nil, "")
	if err := mgr.Save(call); err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	exists := mgr.Exists("exists-test")
	if !exists {
		t.Fatalf("expected Exists to return true for saved call")
	}
}

// TestSaveAndLoadWithSpecialCharacters tests saving/loading with special chars in data
func TestSaveAndLoadWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	call := NewSavedCall(
		"special-chars",
		"POST",
		"https://api.example.com/users",
		map[string]string{
			"X-Custom": "value with spaces & special !@#$%",
		},
		nil,
		`{"name":"John \"Doe\"","email":"john@example.com"}`,
	)

	err := mgr.Save(call)
	if err != nil {
		t.Fatalf("failed to save call with special chars: %v", err)
	}

	loaded, err := mgr.Load("special-chars")
	if err != nil {
		t.Fatalf("failed to load call with special chars: %v", err)
	}

	if loaded.Body != call.Body {
		t.Errorf("body mismatch after load")
	}

	if loaded.Headers["X-Custom"] != "value with spaces & special !@#$%" {
		t.Errorf("header value not preserved with special chars")
	}
}

// TestManagerSaveOverwriteExisting tests overwriting existing call
func TestManagerSaveOverwriteExisting(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	call1 := NewSavedCall("test", "GET", "https://api.example.com/v1", nil, nil, "")
	if err := mgr.Save(call1); err != nil {
		t.Fatalf("failed to save call1: %v", err)
	}

	call2 := NewSavedCall("test", "POST", "https://api.example.com/v2", nil, nil, "body")
	err := mgr.Save(call2)
	if err != nil {
		t.Fatalf("failed to overwrite call: %v", err)
	}

	loaded, err := mgr.Load("test")
	if err != nil {
		t.Fatalf("failed to load overwritten call: %v", err)
	}

	if loaded.Method != "POST" {
		t.Errorf("method not updated on overwrite: got %q", loaded.Method)
	}

	if loaded.URL != "https://api.example.com/v2" {
		t.Errorf("URL not updated on overwrite")
	}
}
