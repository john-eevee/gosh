package integration

import (
	"testing"
	"time"

	"github.com/gosh/internal/request"
	"github.com/gosh/internal/storage"
)

const (
	timeoutSavedCalls = 10 * time.Second
)

// TestSaveCallBasic tests saving a basic HTTP call
func TestSaveCallBasic(t *testing.T) {
	tempDir := t.TempDir()
	manager := storage.NewManager(tempDir)

	call := storage.NewSavedCall(
		"get-api",
		"GET",
		"https://api.example.com/users",
		map[string]string{"Authorization": "Bearer token"},
		map[string]string{"limit": "10"},
		"",
	)

	err := manager.Save(call)
	if err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	if !manager.Exists("get-api") {
		t.Error("expected saved call to exist")
	}
}

// TestRecallSavedCall tests loading a saved call
func TestRecallSavedCall(t *testing.T) {
	tempDir := t.TempDir()
	manager := storage.NewManager(tempDir)

	originalCall := storage.NewSavedCall(
		"post-api",
		"POST",
		"https://api.example.com/users",
		map[string]string{"Content-Type": "application/json"},
		map[string]string{},
		`{"name":"John","email":"john@example.com"}`,
	)
	originalCall.Description = "Create new user"

	err := manager.Save(originalCall)
	if err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	// Load the call back
	loaded, err := manager.Load("post-api")
	if err != nil {
		t.Fatalf("failed to load call: %v", err)
	}

	if loaded.Name != "post-api" {
		t.Errorf("expected name=post-api, got %s", loaded.Name)
	}
	if loaded.Method != "POST" {
		t.Errorf("expected method=POST, got %s", loaded.Method)
	}
	if loaded.URL != "https://api.example.com/users" {
		t.Errorf("expected correct URL, got %s", loaded.URL)
	}
	if loaded.Body != `{"name":"John","email":"john@example.com"}` {
		t.Errorf("expected correct body, got %s", loaded.Body)
	}
	if loaded.Description != "Create new user" {
		t.Errorf("expected description, got %s", loaded.Description)
	}
}

// TestListSavedCalls tests listing all saved calls
func TestListSavedCalls(t *testing.T) {
	tempDir := t.TempDir()
	manager := storage.NewManager(tempDir)

	// Save multiple calls
	calls := []*storage.SavedCall{
		storage.NewSavedCall("call1", "GET", "https://api.example.com/1", map[string]string{}, map[string]string{}, ""),
		storage.NewSavedCall("call2", "POST", "https://api.example.com/2", map[string]string{}, map[string]string{}, "body"),
		storage.NewSavedCall("call3", "DELETE", "https://api.example.com/3", map[string]string{}, map[string]string{}, ""),
	}

	for _, call := range calls {
		err := manager.Save(call)
		if err != nil {
			t.Fatalf("failed to save call %s: %v", call.Name, err)
		}
	}

	// List all calls
	list, err := manager.List()
	if err != nil {
		t.Fatalf("failed to list calls: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("expected 3 calls, got %d", len(list))
	}

	// Verify all calls are in the list
	callNames := make(map[string]bool)
	for _, call := range list {
		callNames[call.Name] = true
	}

	expectedNames := map[string]bool{"call1": true, "call2": true, "call3": true}
	for name := range expectedNames {
		if !callNames[name] {
			t.Errorf("expected call %s in list", name)
		}
	}
}

// TestDeleteSavedCall tests deleting a saved call
func TestDeleteSavedCall(t *testing.T) {
	tempDir := t.TempDir()
	manager := storage.NewManager(tempDir)

	call := storage.NewSavedCall(
		"to-delete",
		"GET",
		"https://api.example.com/delete",
		map[string]string{},
		map[string]string{},
		"",
	)

	err := manager.Save(call)
	if err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	if !manager.Exists("to-delete") {
		t.Error("expected call to exist before deletion")
	}

	// Delete the call
	err = manager.Delete("to-delete")
	if err != nil {
		t.Fatalf("failed to delete call: %v", err)
	}

	if manager.Exists("to-delete") {
		t.Error("expected call to not exist after deletion")
	}
}

// TestSaveCallWithAllFields tests saving call with all fields
func TestSaveCallWithAllFields(t *testing.T) {
	tempDir := t.TempDir()
	manager := storage.NewManager(tempDir)

	headers := map[string]string{
		"Authorization": "Bearer token123",
		"Content-Type":  "application/json",
		"X-Custom":      "value",
	}
	queryParams := map[string]string{
		"filter": "active",
		"sort":   "name",
		"limit":  "50",
	}
	body := `{"key":"value","nested":{"field":"data"}}`

	call := storage.NewSavedCall(
		"complex-call",
		"PUT",
		"https://api.example.com/resource/123",
		headers,
		queryParams,
		body,
	)
	call.Description = "Update resource with complex parameters"

	err := manager.Save(call)
	if err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	loaded, err := manager.Load("complex-call")
	if err != nil {
		t.Fatalf("failed to load call: %v", err)
	}

	if len(loaded.Headers) != 3 {
		t.Errorf("expected 3 headers, got %d", len(loaded.Headers))
	}
	if loaded.Headers["Authorization"] != "Bearer token123" {
		t.Errorf("expected Authorization header, got %s", loaded.Headers["Authorization"])
	}

	if len(loaded.QueryParams) != 3 {
		t.Errorf("expected 3 query params, got %d", len(loaded.QueryParams))
	}
	if loaded.QueryParams["filter"] != "active" {
		t.Errorf("expected filter=active, got %s", loaded.QueryParams["filter"])
	}

	if loaded.Body != body {
		t.Errorf("expected correct body, got %s", loaded.Body)
	}
}

// TestSavedCallWithHTTPExecution tests using a saved call to make a real HTTP request
func TestSavedCallWithHTTPExecution(t *testing.T) {
	SkipIfOffline(t)

	tempDir := t.TempDir()
	manager := storage.NewManager(tempDir)

	// Create and save a GET request
	call := storage.NewSavedCall(
		"get-httpbin",
		"GET",
		"https://httpbin.org/get",
		map[string]string{"X-Custom": "header-value"},
		map[string]string{"param": "value"},
		"",
	)

	err := manager.Save(call)
	if err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	// Load and execute the call
	loaded, err := manager.Load("get-httpbin")
	if err != nil {
		t.Fatalf("failed to load call: %v", err)
	}

	// Convert SavedCall to request.Request for execution
	executor := request.NewExecutor(timeout)
	httpReq := &request.Request{
		Method:      loaded.Method,
		URL:         loaded.URL,
		Headers:     loaded.Headers,
		QueryParams: loaded.QueryParams,
		Body:        loaded.Body,
		Timeout:     timeoutSavedCalls,
	}

	resp, err := executor.Execute(httpReq)
	if err != nil {
		t.Fatalf("failed to execute call: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

// TestSavedCallOverwrite tests that saving with same name overwrites
func TestSavedCallOverwrite(t *testing.T) {
	tempDir := t.TempDir()
	manager := storage.NewManager(tempDir)

	// Save initial call
	call1 := storage.NewSavedCall(
		"update-me",
		"GET",
		"https://api.example.com/v1",
		map[string]string{},
		map[string]string{},
		"",
	)

	err := manager.Save(call1)
	if err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	// Overwrite with new call
	call2 := storage.NewSavedCall(
		"update-me",
		"POST",
		"https://api.example.com/v2",
		map[string]string{"Content-Type": "application/json"},
		map[string]string{},
		`{"new":"body"}`,
	)

	err = manager.Save(call2)
	if err != nil {
		t.Fatalf("failed to save call: %v", err)
	}

	// Load and verify it's the new one
	loaded, err := manager.Load("update-me")
	if err != nil {
		t.Fatalf("failed to load call: %v", err)
	}

	if loaded.Method != "POST" {
		t.Errorf("expected method=POST, got %s", loaded.Method)
	}
	if loaded.URL != "https://api.example.com/v2" {
		t.Errorf("expected updated URL, got %s", loaded.URL)
	}
	if loaded.Body != `{"new":"body"}` {
		t.Errorf("expected updated body, got %s", loaded.Body)
	}
}

// TestSaveCallEmptyList tests listing when no calls exist
func TestSaveCallEmptyList(t *testing.T) {
	tempDir := t.TempDir()
	manager := storage.NewManager(tempDir)

	list, err := manager.List()
	if err != nil {
		t.Fatalf("expected no error listing empty calls, got %v", err)
	}

	if len(list) != 0 {
		t.Errorf("expected 0 calls, got %d", len(list))
	}
}

// TestSaveCallWithSpecialCharacters tests call names with special characters
func TestSaveCallWithSpecialCharacters(t *testing.T) {
	tempDir := t.TempDir()
	manager := storage.NewManager(tempDir)

	// Test with various call names
	callNames := []string{
		"get-users",
		"create_post",
		"delete.item",
	}

	for _, name := range callNames {
		call := storage.NewSavedCall(
			name,
			"GET",
			"https://api.example.com/test",
			map[string]string{},
			map[string]string{},
			"",
		)

		err := manager.Save(call)
		if err != nil {
			t.Fatalf("failed to save call with name %s: %v", name, err)
		}

		if !manager.Exists(name) {
			t.Errorf("expected call %s to exist", name)
		}

		loaded, err := manager.Load(name)
		if err != nil {
			t.Fatalf("failed to load call %s: %v", name, err)
		}

		if loaded.Name != name {
			t.Errorf("expected name=%s, got %s", name, loaded.Name)
		}
	}
}
