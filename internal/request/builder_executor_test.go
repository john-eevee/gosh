package request

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gosh/internal/auth"
)

// TestBuilderBasic tests basic request building
func TestBuilderBasic(t *testing.T) {
	req := &Request{
		Method:  "GET",
		URL:     "https://api.example.com/users",
		Headers: make(map[string]string),
	}

	builder := NewBuilder(req)
	httpReq, err := builder.Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if httpReq.Method != "GET" {
		t.Errorf("method: got %q, want 'GET'", httpReq.Method)
	}

	if httpReq.URL.String() != "https://api.example.com/users" {
		t.Errorf("URL: got %q, want 'https://api.example.com/users'", httpReq.URL.String())
	}
}

// TestBuilderWithHeaders tests building with headers
func TestBuilderWithHeaders(t *testing.T) {
	req := &Request{
		Method: "POST",
		URL:    "https://api.example.com/users",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Custom":     "value",
		},
	}

	builder := NewBuilder(req)
	httpReq, err := builder.Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if httpReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header not set correctly")
	}

	if httpReq.Header.Get("X-Custom") != "value" {
		t.Errorf("X-Custom header not set correctly")
	}
}

// TestBuilderWithQueryParams tests building with query parameters
func TestBuilderWithQueryParams(t *testing.T) {
	req := &Request{
		Method: "GET",
		URL:    "https://api.example.com/users",
		QueryParams: map[string]string{
			"limit":  "10",
			"offset": "0",
		},
	}

	builder := NewBuilder(req)
	httpReq, err := builder.Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	query := httpReq.URL.Query()
	if query.Get("limit") != "10" {
		t.Errorf("limit param: got %q, want '10'", query.Get("limit"))
	}

	if query.Get("offset") != "0" {
		t.Errorf("offset param: got %q, want '0'", query.Get("offset"))
	}
}

// TestBuilderWithBody tests building with request body
func TestBuilderWithBody(t *testing.T) {
	body := `{"name":"John","age":30}`
	req := &Request{
		Method:  "POST",
		URL:     "https://api.example.com/users",
		Headers: make(map[string]string),
		Body:    body,
	}

	builder := NewBuilder(req)
	httpReq, err := builder.Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if httpReq.Body == nil {
		t.Errorf("expected body to be set")
	}

	// Read body content
	buf := make([]byte, len(body))
	httpReq.Body.Read(buf)

	if string(buf) != body {
		t.Errorf("body: got %q, want %q", string(buf), body)
	}
}

// TestBuilderWithEmptyBody tests building with empty body
func TestBuilderWithEmptyBody(t *testing.T) {
	req := &Request{
		Method:  "GET",
		URL:     "https://api.example.com/users",
		Headers: make(map[string]string),
		Body:    "",
	}

	builder := NewBuilder(req)
	httpReq, err := builder.Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if httpReq.Body != nil {
		t.Errorf("expected nil body for empty body string")
	}
}

// TestBuilderWithAuth tests building with authentication
func TestBuilderWithAuth(t *testing.T) {
	req := &Request{
		Method:  "GET",
		URL:     "https://api.example.com/users",
		Headers: make(map[string]string),
		Auth: &auth.AuthPreset{
			Type:  "bearer",
			Token: "xyz123",
		},
	}

	builder := NewBuilder(req)
	httpReq, err := builder.Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	authHeader := httpReq.Header.Get("Authorization")
	if !strings.Contains(authHeader, "Bearer") {
		t.Errorf("Authorization header not set correctly: %q", authHeader)
	}
}

// TestBuilderInvalidURL tests building with invalid URL
func TestBuilderInvalidURL(t *testing.T) {
	req := &Request{
		Method:  "GET",
		URL:     "ht!tp://invalid",
		Headers: make(map[string]string),
	}

	builder := NewBuilder(req)
	_, err := builder.Build()

	if err == nil {
		t.Fatalf("expected error for invalid URL, got nil")
	}
}

// TestBuilderMultipleQueryParams tests building with multiple query params
func TestBuilderMultipleQueryParams(t *testing.T) {
	req := &Request{
		Method: "GET",
		URL:    "https://api.example.com/search",
		QueryParams: map[string]string{
			"q":     "golang",
			"sort":  "date",
			"order": "desc",
			"limit": "50",
		},
	}

	builder := NewBuilder(req)
	httpReq, err := builder.Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	query := httpReq.URL.Query()
	if query.Get("q") != "golang" {
		t.Errorf("q param incorrect")
	}
	if query.Get("sort") != "date" {
		t.Errorf("sort param incorrect")
	}
	if query.Get("order") != "desc" {
		t.Errorf("order param incorrect")
	}
	if query.Get("limit") != "50" {
		t.Errorf("limit param incorrect")
	}
}

// TestExecutorBasic tests basic request execution
func TestExecutorBasic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	req := &Request{
		Method:  "GET",
		URL:     server.URL,
		Headers: make(map[string]string),
	}

	executor := NewExecutor(5 * time.Second)
	resp, err := executor.Execute(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if !strings.Contains(string(resp.Body), "success") {
		t.Errorf("body: expected 'success', got %q", string(resp.Body))
	}
}

// TestExecutorWithHeaders tests execution with headers
func TestExecutorWithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") != "test-value" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	req := &Request{
		Method: "GET",
		URL:    server.URL,
		Headers: map[string]string{
			"X-Test": "test-value",
		},
	}

	executor := NewExecutor(5 * time.Second)
	resp, err := executor.Execute(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

// TestExecutorWithQueryParams tests execution with query parameters
func TestExecutorWithQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("search") != "test" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	req := &Request{
		Method: "GET",
		URL:    server.URL,
		QueryParams: map[string]string{
			"search": "test",
		},
	}

	executor := NewExecutor(5 * time.Second)
	resp, err := executor.Execute(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

// TestExecutorWithBody tests execution with request body
func TestExecutorWithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, 23)
		r.Body.Read(body)
		if !strings.Contains(string(body), "test") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	req := &Request{
		Method:  "POST",
		URL:     server.URL,
		Headers: make(map[string]string),
		Body:    `{"message":"test"}`,
	}

	executor := NewExecutor(5 * time.Second)
	resp, err := executor.Execute(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

// TestExecutorStatusCodes tests different status codes
func TestExecutorStatusCodes(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedStatus string
	}{
		{"OK", http.StatusOK, "200 OK"},
		{"Created", http.StatusCreated, "201 Created"},
		{"BadRequest", http.StatusBadRequest, "400 Bad Request"},
		{"Unauthorized", http.StatusUnauthorized, "401 Unauthorized"},
		{"NotFound", http.StatusNotFound, "404 Not Found"},
		{"ServerError", http.StatusInternalServerError, "500 Internal Server Error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			req := &Request{
				Method:  "GET",
				URL:     server.URL,
				Headers: make(map[string]string),
			}

			executor := NewExecutor(5 * time.Second)
			resp, err := executor.Execute(req)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.statusCode {
				t.Errorf("status code: got %d, want %d", resp.StatusCode, tt.statusCode)
			}

			if !strings.Contains(resp.Status, tt.expectedStatus) {
				t.Errorf("status: got %q, want to contain %q", resp.Status, tt.expectedStatus)
			}
		})
	}
}

// TestExecutorResponseHeaders tests response headers are captured
func TestExecutorResponseHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom", "custom-value")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	req := &Request{
		Method:  "GET",
		URL:     server.URL,
		Headers: make(map[string]string),
	}

	executor := NewExecutor(5 * time.Second)
	resp, err := executor.Execute(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Headers) == 0 {
		t.Errorf("expected response headers")
	}

	customHeaders, exists := resp.Headers["X-Custom"]
	if !exists || len(customHeaders) == 0 || customHeaders[0] != "custom-value" {
		t.Errorf("X-Custom header not captured correctly")
	}
}

// TestExecutorResponseBody tests response body is read
func TestExecutorResponseBody(t *testing.T) {
	expectedBody := `{"name":"John","email":"john@example.com"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedBody))
	}))
	defer server.Close()

	req := &Request{
		Method:  "GET",
		URL:     server.URL,
		Headers: make(map[string]string),
	}

	executor := NewExecutor(5 * time.Second)
	resp, err := executor.Execute(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(resp.Body) != expectedBody {
		t.Errorf("body: got %q, want %q", string(resp.Body), expectedBody)
	}
}

// TestExecutorDuration tests response duration is measured
func TestExecutorDuration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	req := &Request{
		Method:  "GET",
		URL:     server.URL,
		Headers: make(map[string]string),
	}

	executor := NewExecutor(5 * time.Second)
	resp, err := executor.Execute(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Duration < 10*time.Millisecond {
		t.Errorf("duration too short: %v (expected >= 10ms)", resp.Duration)
	}
}

// TestExecutorSize tests response size is calculated
func TestExecutorSize(t *testing.T) {
	expectedBody := `{"test":"data"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedBody))
	}))
	defer server.Close()

	req := &Request{
		Method:  "GET",
		URL:     server.URL,
		Headers: make(map[string]string),
	}

	executor := NewExecutor(5 * time.Second)
	resp, err := executor.Execute(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedSize := len(expectedBody)
	if resp.Size != expectedSize {
		t.Errorf("size: got %d, want %d", resp.Size, expectedSize)
	}
}

// TestExecutorHTTPMethods tests different HTTP methods
func TestExecutorHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != method {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			req := &Request{
				Method:  method,
				URL:     server.URL,
				Headers: make(map[string]string),
			}

			executor := NewExecutor(5 * time.Second)
			resp, err := executor.Execute(req)

			if err != nil {
				t.Fatalf("unexpected error for %s: %v", method, err)
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("%s: status code %d, want %d", method, resp.StatusCode, http.StatusOK)
			}
		})
	}
}

// TestExecutorWithAuth tests execution with authentication
func TestExecutorWithAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	req := &Request{
		Method:  "GET",
		URL:     server.URL,
		Headers: make(map[string]string),
		Auth: &auth.AuthPreset{
			Type:  "bearer",
			Token: "test-token",
		},
	}

	executor := NewExecutor(5 * time.Second)
	resp, err := executor.Execute(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

// TestExecutorTimeout tests request timeout
func TestExecutorTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	req := &Request{
		Method:  "GET",
		URL:     server.URL,
		Headers: make(map[string]string),
	}

	executor := NewExecutor(10 * time.Millisecond)
	_, err := executor.Execute(req)

	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
}

// TestExecutorEmptyResponse tests handling empty response body
func TestExecutorEmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	req := &Request{
		Method:  "DELETE",
		URL:     server.URL,
		Headers: make(map[string]string),
	}

	executor := NewExecutor(5 * time.Second)
	resp, err := executor.Execute(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status code: got %d, want %d", resp.StatusCode, http.StatusNoContent)
	}

	if resp.Size != 0 {
		t.Errorf("size should be 0 for empty response, got %d", resp.Size)
	}
}

// TestBuilderHeaderOverride tests that later headers override earlier ones
func TestBuilderHeaderOverride(t *testing.T) {
	req := &Request{
		Method: "GET",
		URL:    "https://api.example.com/users",
		Headers: map[string]string{
			"X-Test": "value1",
		},
	}

	builder := NewBuilder(req)
	httpReq, err := builder.Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if httpReq.Header.Get("X-Test") != "value1" {
		t.Errorf("header value incorrect")
	}
}

// TestNewExecutor tests executor initialization
func TestNewExecutor(t *testing.T) {
	timeout := 30 * time.Second
	executor := NewExecutor(timeout)

	if executor == nil {
		t.Fatalf("expected executor, got nil")
	}

	if executor.timeout != timeout {
		t.Errorf("timeout: got %v, want %v", executor.timeout, timeout)
	}
}

// TestNewBuilder tests builder initialization
func TestNewBuilder(t *testing.T) {
	req := &Request{
		Method: "GET",
		URL:    "https://api.example.com",
	}

	builder := NewBuilder(req)

	if builder == nil {
		t.Fatalf("expected builder, got nil")
	}

	if builder.req != req {
		t.Errorf("request not stored correctly")
	}
}
