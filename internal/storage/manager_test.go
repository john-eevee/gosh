package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	mgr := NewManager(tmpDir)

	// Create a saved call
	call := NewSavedCall(
		"test-request",
		"POST",
		"https://api.example.com/users",
		map[string]string{"Authorization": "Bearer token"},
		map[string]string{},
		`{"name":"John"}`,
	)

	// Save it
	if err := mgr.Save(call); err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	// Load it back
	loaded, err := mgr.Load("test-request")
	if err != nil {
		t.Fatalf("failed to load call: %v", err)
	}

	if loaded.Name != "test-request" {
		t.Errorf("expected name 'test-request', got '%s'", loaded.Name)
	}

	if loaded.Method != "POST" {
		t.Errorf("expected method 'POST', got '%s'", loaded.Method)
	}

	if loaded.URL != "https://api.example.com/users" {
		t.Errorf("expected URL 'https://api.example.com/users', got '%s'", loaded.URL)
	}

	if loaded.Body != `{"name":"John"}` {
		t.Errorf("expected body '{\"name\":\"John\"}', got '%s'", loaded.Body)
	}
}

func TestList(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	// Save multiple calls
	calls := []struct {
		name   string
		method string
		url    string
	}{
		{"create-user", "POST", "https://api.example.com/users"},
		{"get-users", "GET", "https://api.example.com/users"},
		{"delete-user", "DELETE", "https://api.example.com/users/{id}"},
	}

	for _, c := range calls {
		call := NewSavedCall(c.name, c.method, c.url, map[string]string{}, map[string]string{}, "")
		if err := mgr.Save(call); err != nil {
			t.Fatalf("failed to save call %s: %v", c.name, err)
		}
	}

	// List calls
	list, err := mgr.List()
	if err != nil {
		t.Fatalf("failed to list calls: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("expected 3 calls, got %d", len(list))
	}
}

func TestDelete(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	// Save and delete
	call := NewSavedCall("to-delete", "GET", "https://api.example.com/test", map[string]string{}, map[string]string{}, "")
	if err := mgr.Save(call); err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	if !mgr.Exists("to-delete") {
		t.Fatal("call should exist after saving")
	}

	if err := mgr.Delete("to-delete"); err != nil {
		t.Fatalf("failed to delete call: %v", err)
	}

	if mgr.Exists("to-delete") {
		t.Fatal("call should not exist after deleting")
	}
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	call := NewSavedCall("test", "GET", "https://api.example.com", map[string]string{}, map[string]string{}, "")
	if err := mgr.Save(call); err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	if !mgr.Exists("test") {
		t.Fatal("expected call to exist")
	}

	if mgr.Exists("nonexistent") {
		t.Fatal("expected call to not exist")
	}
}

func TestGetCallPath(t *testing.T) {
	expected := filepath.Join("/workspace", ".gosh", "calls", "my-request.yaml")
	actual := GetCallPath("/workspace", "my-request")
	if actual != expected {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}
}

func TestGetCallsDir(t *testing.T) {
	tmpDir := t.TempDir()

	dir, err := GetCallsDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to get calls dir: %v", err)
	}

	// Verify directory was created
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Fatal("expected directory")
	}
}
