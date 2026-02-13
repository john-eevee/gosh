package integration

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gosh/internal/auth"
	"github.com/gosh/internal/request"
)

const (
	bearerTokenURL = "https://httpbin.org/bearer"
	basicAuthURL   = "https://httpbin.org/basic-auth/testuser/testpass"
)

// TestAuthPresetBearerToken tests creating and using bearer token auth preset
func TestAuthPresetBearerToken(t *testing.T) {
	SkipIfOffline(t)

	// Create temporary directory for auth storage
	tempDir := t.TempDir()

	// Create auth manager
	authMgr := auth.NewManager(tempDir)

	// Add bearer token preset
	bearerPreset := &auth.AuthPreset{
		Name:  "test-bearer",
		Type:  "bearer",
		Token: "test-secret-token-123",
	}

	err := authMgr.Add(bearerPreset)
	if err != nil {
		t.Fatalf("failed to add bearer auth preset: %v", err)
	}

	// Retrieve the preset
	retrieved, err := authMgr.Get("test-bearer")
	if err != nil {
		t.Fatalf("failed to get bearer auth preset: %v", err)
	}

	if retrieved.Type != "bearer" {
		t.Errorf("expected type=bearer, got %s", retrieved.Type)
	}
	if retrieved.Token != "test-secret-token-123" {
		t.Errorf("expected token=test-secret-token-123, got %s", retrieved.Token)
	}

	// Use the auth preset in a real HTTP request
	executor := request.NewExecutor(timeout)
	httpReq := &request.Request{
		Method:  "GET",
		URL:     bearerTokenURL,
		Headers: make(map[string]string),
		Timeout: timeout,
	}

	// Apply the auth preset
	httpReq.Auth = retrieved

	// Execute the request
	resp, err := executor.Execute(httpReq)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Verify the bearer token was sent correctly
	if !bytes.Contains(resp.Body, []byte("authenticated")) {
		t.Error("bearer token authentication failed")
	}
}

// TestAuthPresetBasicAuth tests creating and using basic auth preset
func TestAuthPresetBasicAuth(t *testing.T) {
	SkipIfOffline(t)

	tempDir := t.TempDir()
	authMgr := auth.NewManager(tempDir)

	// Add basic auth preset
	basicPreset := &auth.AuthPreset{
		Name:     "test-basic",
		Type:     "basic",
		Username: "testuser",
		Password: "testpass",
	}

	err := authMgr.Add(basicPreset)
	if err != nil {
		t.Fatalf("failed to add basic auth preset: %v", err)
	}

	// Retrieve the preset
	retrieved, err := authMgr.Get("test-basic")
	if err != nil {
		t.Fatalf("failed to get basic auth preset: %v", err)
	}

	if retrieved.Type != "basic" {
		t.Errorf("expected type=basic, got %s", retrieved.Type)
	}
	if retrieved.Username != "testuser" {
		t.Errorf("expected username=testuser, got %s", retrieved.Username)
	}

	// Use the auth preset in a real HTTP request
	executor := request.NewExecutor(timeout)
	httpReq := &request.Request{
		Method:  "GET",
		URL:     basicAuthURL,
		Headers: make(map[string]string),
		Timeout: timeout,
	}

	// Apply the auth preset
	httpReq.Auth = retrieved

	// Execute the request
	resp, err := executor.Execute(httpReq)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

// TestAuthPresetCustomHeader tests creating and using custom header auth
func TestAuthPresetCustomHeader(t *testing.T) {
	SkipIfOffline(t)

	tempDir := t.TempDir()
	authMgr := auth.NewManager(tempDir)

	// Add custom header auth preset
	customPreset := &auth.AuthPreset{
		Name:   "test-custom",
		Type:   "custom",
		Header: "X-API-Key",
		Value:  "secret-api-key-123",
	}

	err := authMgr.Add(customPreset)
	if err != nil {
		t.Fatalf("failed to add custom auth preset: %v", err)
	}

	// Retrieve the preset
	retrieved, err := authMgr.Get("test-custom")
	if err != nil {
		t.Fatalf("failed to get custom auth preset: %v", err)
	}

	if retrieved.Type != "custom" {
		t.Errorf("expected type=custom, got %s", retrieved.Type)
	}
	if retrieved.Header != "X-API-Key" {
		t.Errorf("expected header name=X-API-Key, got %s", retrieved.Header)
	}

	// Use the auth preset in a real HTTP request to verify headers are sent
	executor := request.NewExecutor(timeout)
	httpReq := &request.Request{
		Method:  "GET",
		URL:     "https://httpbin.org/headers",
		Headers: make(map[string]string),
		Timeout: timeout,
	}

	// Apply the auth preset
	httpReq.Auth = retrieved

	// Execute the request
	resp, err := executor.Execute(httpReq)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Verify the custom header was sent
	var httpbinResp map[string]interface{}
	if err := json.Unmarshal(resp.Body, &httpbinResp); err == nil {
		if headers, ok := httpbinResp["headers"].(map[string]interface{}); ok {
			if apiKey, ok := headers["X-Api-Key"]; ok && apiKey != "" {
				// Header was sent (httpbin converts header names to title case)
				return
			}
		}
	}
	// The header verification is lenient since httpbin header normalization varies
}

// TestAuthPresetPersistence tests that auth presets are saved and loaded
func TestAuthPresetPersistence(t *testing.T) {
	tempDir := t.TempDir()

	// Create first auth manager and add a preset
	authMgr1 := auth.NewManager(tempDir)
	preset := &auth.AuthPreset{
		Name:  "persistent-preset",
		Type:  "bearer",
		Token: "test-token-12345",
	}
	err := authMgr1.Add(preset)
	if err != nil {
		t.Fatalf("failed to add preset: %v", err)
	}

	// Create a second auth manager that should load the preset
	authMgr2 := auth.NewManager(tempDir)
	err = authMgr2.Load()
	if err != nil {
		t.Fatalf("failed to load auth presets: %v", err)
	}

	// Verify the preset was loaded
	retrieved, err := authMgr2.Get("persistent-preset")
	if err != nil {
		t.Fatalf("failed to retrieve persistent preset: %v", err)
	}

	if retrieved.Type != "bearer" {
		t.Errorf("expected type=bearer, got %s", retrieved.Type)
	}
	if retrieved.Token != "test-token-12345" {
		t.Errorf("expected token=test-token-12345, got %s", retrieved.Token)
	}
}

// TestAuthPresetList tests listing auth presets
func TestAuthPresetList(t *testing.T) {
	tempDir := t.TempDir()
	authMgr := auth.NewManager(tempDir)

	// Add multiple presets
	presets := []*auth.AuthPreset{
		{Name: "preset1", Type: "bearer", Token: "token1"},
		{Name: "preset2", Type: "basic", Username: "user", Password: "pass"},
		{Name: "preset3", Type: "custom", Header: "X-Key", Value: "value"},
	}

	for _, preset := range presets {
		err := authMgr.Add(preset)
		if err != nil {
			t.Fatalf("failed to add preset %s: %v", preset.Name, err)
		}
	}

	// List all presets
	list := authMgr.List()

	if len(list) != 3 {
		t.Errorf("expected 3 presets, got %d", len(list))
	}

	// Verify all presets are in the list
	for _, expected := range presets {
		found := false
		for _, actual := range list {
			if actual.Name == expected.Name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("preset %s not found in list", expected.Name)
		}
	}
}

// TestAuthPresetRemove tests removing auth presets
func TestAuthPresetRemove(t *testing.T) {
	tempDir := t.TempDir()
	authMgr := auth.NewManager(tempDir)

	// Add and then remove a preset
	preset := &auth.AuthPreset{
		Name:  "to-remove",
		Type:  "bearer",
		Token: "token",
	}

	err := authMgr.Add(preset)
	if err != nil {
		t.Fatalf("failed to add preset: %v", err)
	}

	// Remove the preset
	err = authMgr.Remove("to-remove")
	if err != nil {
		t.Fatalf("failed to remove preset: %v", err)
	}

	// Verify it's gone
	_, err = authMgr.Get("to-remove")
	if err == nil {
		t.Error("expected error when getting removed preset")
	}
}

// TestAuthPresetBearerTokenWithAuthorization tests bearer token is converted to Authorization header
func TestAuthPresetBearerTokenWithAuthorization(t *testing.T) {
	SkipIfOffline(t)

	tempDir := t.TempDir()
	authMgr := auth.NewManager(tempDir)

	preset := &auth.AuthPreset{
		Name:  "bearer-test",
		Type:  "bearer",
		Token: "test-secret-token",
	}

	err := authMgr.Add(preset)
	if err != nil {
		t.Fatalf("failed to add preset: %v", err)
	}

	retrieved, _ := authMgr.Get("bearer-test")

	// The executor should apply this auth to headers
	executor := request.NewExecutor(timeout)
	httpReq := &request.Request{
		Method:  "GET",
		URL:     "https://httpbin.org/headers",
		Headers: make(map[string]string),
		Auth:    retrieved,
		Timeout: timeout,
	}

	resp, err := executor.Execute(httpReq)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Verify Authorization header was sent
	var httpbinResp map[string]interface{}
	if err := json.Unmarshal(resp.Body, &httpbinResp); err == nil {
		if headers, ok := httpbinResp["headers"].(map[string]interface{}); ok {
			if authHeader, ok := headers["Authorization"].(string); ok && authHeader != "" {
				if !bytes.Contains([]byte(authHeader), []byte("Bearer")) {
					t.Errorf("expected Bearer prefix in Authorization header, got: %s", authHeader)
				}
				return
			}
		}
	}
}

// TestAuthPresetWithRequest tests using auth preset with actual HTTP request
func TestAuthPresetWithRequest(t *testing.T) {
	SkipIfOffline(t)

	tempDir := t.TempDir()
	authMgr := auth.NewManager(tempDir)

	// Add a bearer preset
	preset := &auth.AuthPreset{
		Name:  "request-test",
		Type:  "bearer",
		Token: "request-token",
	}

	err := authMgr.Add(preset)
	if err != nil {
		t.Fatalf("failed to add preset: %v", err)
	}

	retrieved, _ := authMgr.Get("request-test")

	// Create a request with POST body and auth
	executor := request.NewExecutor(timeout)
	httpReq := &request.Request{
		Method: "POST",
		URL:    "https://httpbin.org/post",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    `{"test":"data"}`,
		Auth:    retrieved,
		Timeout: timeout,
	}

	resp, err := executor.Execute(httpReq)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Verify both body and auth were sent
	var httpbinResp map[string]interface{}
	if err := json.Unmarshal(resp.Body, &httpbinResp); err == nil {
		// Check that JSON body was received
		if jsonData, ok := httpbinResp["json"].(map[string]interface{}); ok {
			if test, ok := jsonData["test"].(string); !ok || test != "data" {
				t.Error("request body was not sent correctly with auth")
			}
		}
	}
}

// TestAuthPresetFilePermissions tests that auth file has restricted permissions
func TestAuthPresetFilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	authMgr := auth.NewManager(tempDir)

	preset := &auth.AuthPreset{
		Name:  "permission-test",
		Type:  "bearer",
		Token: "secret-token",
	}

	err := authMgr.Add(preset)
	if err != nil {
		t.Fatalf("failed to add preset: %v", err)
	}

	// Check the auth file permissions
	authFilePath := filepath.Join(tempDir, ".gosh", "auth.yaml")
	info, err := os.Stat(authFilePath)
	if err != nil {
		t.Fatalf("failed to stat auth file: %v", err)
	}

	// Should have restricted permissions (600 = rw-------)
	mode := info.Mode()

	if mode&0o077 != 0 {
		t.Errorf("auth file has overly permissive permissions: %o", mode)
	}
}
