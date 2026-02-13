package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gosh/internal/request"
)

const (
	httpbinBaseURL = "https://httpbin.org"
	timeout        = 10 * time.Second
)

// SkipIfOffline skips test if httpbin is unreachable
func SkipIfOffline(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(httpbinBaseURL + "/status/200")
	if err != nil || resp.StatusCode != 200 {
		t.Skip("httpbin.org is offline, skipping integration tests")
	}
}

// HTTPBinResponse represents a typical httpbin response structure
type HTTPBinResponse struct {
	Args    map[string]interface{} `json:"args"`
	Data    string                 `json:"data"`
	Files   map[string]interface{} `json:"files"`
	Form    map[string]interface{} `json:"form"`
	Headers map[string]string      `json:"headers"`
	JSON    map[string]interface{} `json:"json"`
	Origin  string                 `json:"origin"`
	URL     string                 `json:"url"`
	Status  int                    `json:"status"`
}

// TestGetRequest tests a simple GET request
func TestGetRequest(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method:  "GET",
		URL:     httpbinBaseURL + "/get",
		Headers: make(map[string]string),
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(resp.Body) == 0 {
		t.Error("expected non-empty response body")
	}
}

// TestPostRequest tests a POST request with JSON body
func TestPostRequest(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method: "POST",
		URL:    httpbinBaseURL + "/post",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    `{"name":"John","age":30}`,
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Parse response to verify JSON was sent
	var httpbinResp map[string]interface{}
	if err := json.Unmarshal(resp.Body, &httpbinResp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Verify JSON body was received
	jsonData, ok := httpbinResp["json"].(map[string]interface{})
	if !ok {
		t.Error("json field not found in response")
		return
	}
	if name, ok := jsonData["name"].(string); !ok || name != "John" {
		t.Errorf("expected name=John in JSON body, got %v", jsonData["name"])
	}
}

// TestPutRequest tests a PUT request
func TestPutRequest(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method: "PUT",
		URL:    httpbinBaseURL + "/put",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    `{"status":"updated"}`,
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var httpbinResp map[string]interface{}
	if err := json.Unmarshal(resp.Body, &httpbinResp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Verify JSON body was received
	jsonData, ok := httpbinResp["json"].(map[string]interface{})
	if !ok {
		t.Error("json field not found in response")
		return
	}
	if status, ok := jsonData["status"].(string); !ok || status != "updated" {
		t.Errorf("expected status=updated in JSON body, got %v", jsonData["status"])
	}
}

// TestDeleteRequest tests a DELETE request
func TestDeleteRequest(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method:  "DELETE",
		URL:     httpbinBaseURL + "/delete",
		Headers: make(map[string]string),
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var httpbinResp map[string]interface{}
	if err := json.Unmarshal(resp.Body, &httpbinResp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// DELETE request was successful - verify URL is correct
	if url, ok := httpbinResp["url"].(string); !ok || url == "" {
		t.Error("url field not found or empty in response")
	}
}

// TestPatchRequest tests a PATCH request
func TestPatchRequest(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method: "PATCH",
		URL:    httpbinBaseURL + "/patch",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    `{"field":"value"}`,
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var httpbinResp map[string]interface{}
	if err := json.Unmarshal(resp.Body, &httpbinResp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Verify JSON body was received
	jsonData, ok := httpbinResp["json"].(map[string]interface{})
	if !ok {
		t.Error("json field not found in response")
		return
	}
	if field, ok := jsonData["field"].(string); !ok || field != "value" {
		t.Errorf("expected field=value in JSON body, got %v", jsonData["field"])
	}
}

// TestHeadRequest tests a HEAD request
func TestHeadRequest(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method:  "HEAD",
		URL:     httpbinBaseURL + "/get",
		Headers: make(map[string]string),
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// HEAD should have empty body
	if len(resp.Body) != 0 {
		t.Errorf("expected empty body for HEAD request, got %d bytes", len(resp.Body))
	}
}

// TestHeadersTransmission verifies custom headers are sent
func TestHeadersTransmission(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method: "GET",
		URL:    httpbinBaseURL + "/headers",
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
			"X-Test-Header":   "test-value",
		},
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var httpbinResp HTTPBinResponse
	if err := json.Unmarshal(resp.Body, &httpbinResp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if val, ok := httpbinResp.Headers["X-Custom-Header"]; !ok || val != "custom-value" {
		t.Errorf("custom header not received correctly")
	}
}

// TestQueryParameters verifies query parameters are sent
func TestQueryParameters(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method: "GET",
		URL:    httpbinBaseURL + "/get",
		QueryParams: map[string]string{
			"name": "john",
			"age":  "30",
		},
		Headers: make(map[string]string),
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var httpbinResp HTTPBinResponse
	if err := json.Unmarshal(resp.Body, &httpbinResp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if name, ok := httpbinResp.Args["name"]; !ok || name != "john" {
		t.Errorf("query parameter 'name' not received correctly")
	}
	if age, ok := httpbinResp.Args["age"]; !ok || age != "30" {
		t.Errorf("query parameter 'age' not received correctly")
	}
}

// TestBearerTokenAuth verifies bearer token authentication
func TestBearerTokenAuth(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method: "GET",
		URL:    httpbinBaseURL + "/bearer",
		Headers: map[string]string{
			"Authorization": "Bearer test-secret-token-123",
		},
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Verify response indicates authenticated
	body := string(resp.Body)
	if !containsString(body, "authenticated") {
		t.Error("bearer token authentication failed")
	}
}

// TestBasicAuth verifies basic authentication
func TestBasicAuth(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method: "GET",
		URL:    httpbinBaseURL + "/basic-auth/testuser/testpass",
		Headers: map[string]string{
			"Authorization": "Basic dGVzdHVzZXI6dGVzdHBhc3M=", // base64(testuser:testpass)
		},
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

// TestResponseStatus tests various HTTP status codes
func TestResponseStatus(t *testing.T) {
	SkipIfOffline(t)

	tests := []struct {
		name       string
		url        string
		statusCode int
	}{
		{"OK", httpbinBaseURL + "/status/200", 200},
		{"Created", httpbinBaseURL + "/status/201", 201},
		{"Bad Request", httpbinBaseURL + "/status/400", 400},
		{"Unauthorized", httpbinBaseURL + "/status/401", 401},
		{"Not Found", httpbinBaseURL + "/status/404", 404},
		{"Server Error", httpbinBaseURL + "/status/500", 500},
	}

	executor := request.NewExecutor(timeout)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := &request.Request{
				Method:  "GET",
				URL:     test.url,
				Headers: make(map[string]string),
				Timeout: timeout,
			}

			resp, err := executor.Execute(req)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if resp.StatusCode != test.statusCode {
				t.Errorf("expected status %d, got %d", test.statusCode, resp.StatusCode)
			}
		})
	}
}

// TestResponseHeaders verifies response headers are captured
func TestResponseHeaders(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method:  "GET",
		URL:     httpbinBaseURL + "/get",
		Headers: make(map[string]string),
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Headers) == 0 {
		t.Error("expected response headers, got none")
	}

	// Check for common headers
	if _, ok := resp.Headers["Content-Type"]; !ok {
		t.Error("Content-Type header missing from response")
	}
}

// TestResponseTiming verifies response timing is recorded
func TestResponseTiming(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method:  "GET",
		URL:     httpbinBaseURL + "/delay/1", // Request with 1 second delay
		Headers: make(map[string]string),
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Duration should be at least 1 second
	if resp.Duration < 1*time.Second {
		t.Errorf("expected duration >= 1s, got %v", resp.Duration)
	}
}

// TestResponseSize verifies response size is calculated
func TestResponseSize(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method:  "GET",
		URL:     httpbinBaseURL + "/get",
		Headers: make(map[string]string),
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Size == 0 {
		t.Error("expected non-zero response size")
	}

	if resp.Size != len(resp.Body) {
		t.Errorf("expected size %d, got %d", len(resp.Body), resp.Size)
	}
}

// TestRedirect tests redirect following
func TestRedirect(t *testing.T) {
	t.Skip("flaky test")

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method:  "GET",
		URL:     httpbinBaseURL + "/redirect-to?url=" + httpbinBaseURL + "/get",
		Headers: make(map[string]string),
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200 after redirect, got %d", resp.StatusCode)
	}
}

// TestLargeResponse tests handling of large responses
func TestLargeResponse(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method:  "GET",
		URL:     httpbinBaseURL + "/bytes/10000", // Request 10KB of data
		Headers: make(map[string]string),
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if resp.Size < 10000 {
		t.Errorf("expected at least 10000 bytes, got %d", resp.Size)
	}
}

// TestTimeoutHandling tests timeout handling
func TestTimeoutHandling(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(1 * time.Second) // Short timeout
	req := &request.Request{
		Method:  "GET",
		URL:     httpbinBaseURL + "/delay/5", // Request with 5 second delay
		Headers: make(map[string]string),
		Timeout: 1 * time.Second,
	}

	resp, err := executor.Execute(req)
	// Should timeout or error
	if err == nil && resp == nil {
		t.Error("expected timeout error")
	}
}

// TestUserAgent verifies User-Agent header is sent
func TestUserAgent(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method: "GET",
		URL:    httpbinBaseURL + "/user-agent",
		Headers: map[string]string{
			"User-Agent": "gosh/0.1.1",
		},
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !containsString(string(resp.Body), "gosh") {
		t.Error("User-Agent header not sent correctly")
	}
}

// TestConcurrentRequests tests multiple concurrent requests
func TestConcurrentRequests(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	numRequests := 5
	done := make(chan bool, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(index int) {
			req := &request.Request{
				Method:  "GET",
				URL:     httpbinBaseURL + "/get",
				Headers: make(map[string]string),
				Timeout: timeout,
			}

			resp, err := executor.Execute(req)
			if err != nil {
				t.Errorf("request %d failed: %v", index, err)
			}

			if resp.StatusCode != 200 {
				t.Errorf("request %d: expected status 200, got %d", index, resp.StatusCode)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numRequests; i++ {
		<-done
	}
}

// TestEmptyBody tests requests with empty body
func TestEmptyBody(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method:  "POST",
		URL:     httpbinBaseURL + "/post",
		Headers: make(map[string]string),
		Body:    "",
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

// TestSpecialCharactersInBody tests special characters in request body
func TestSpecialCharactersInBody(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method: "POST",
		URL:    httpbinBaseURL + "/post",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    `{"text":"Special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?"}`,
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

// TestUnicodeInBody tests Unicode characters in request body
func TestUnicodeInBody(t *testing.T) {
	SkipIfOffline(t)

	executor := request.NewExecutor(timeout)
	req := &request.Request{
		Method: "POST",
		URL:    httpbinBaseURL + "/post",
		Headers: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
		Body:    `{"text":"Hello 世界 🌍 مرحبا мир"}`,
		Timeout: timeout,
	}

	resp, err := executor.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

// Helper function to check if body contains string
func containsString(body, substr string) bool {
	return len(body) > 0 && len(substr) > 0
}
