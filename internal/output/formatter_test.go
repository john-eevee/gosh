package output

import (
	"strings"
	"testing"
	"time"

	"github.com/gosh/internal/request"
)

// TestNewFormatter tests creating a new formatter
func TestNewFormatter(t *testing.T) {
	// Test non-TTY formatter
	nonTTYFormatter := NewFormatter(false)
	if nonTTYFormatter == nil {
		t.Error("expected non-nil formatter")
	}

	// Test TTY formatter
	ttyFormatter := NewFormatter(true)
	if ttyFormatter == nil {
		t.Error("expected non-nil formatter")
	}
}

// TestFormatResponseStatusOnly tests formatting response with only status
func TestFormatResponseStatusOnly(t *testing.T) {
	formatter := NewFormatter(false)

	resp := &request.Response{
		StatusCode: 200,
		Headers:    make(map[string][]string),
		Body:       []byte{},
		Duration:   100 * time.Millisecond,
		Size:       0,
	}

	output := formatter.FormatResponse(resp, false)

	if !strings.Contains(output, "200") {
		t.Errorf("expected output to contain status code 200, got: %s", output)
	}
}

// TestFormatResponseWithInfo tests formatting response with headers and info
func TestFormatResponseWithInfo(t *testing.T) {
	formatter := NewFormatter(false)

	resp := &request.Response{
		StatusCode: 200,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
			"X-Custom":     {"custom-value"},
		},
		Body:     []byte(`{"key":"value"}`),
		Duration: 250 * time.Millisecond,
		Size:     15,
	}

	output := formatter.FormatResponse(resp, true)

	if !strings.Contains(output, "200") {
		t.Error("expected status code in output")
	}
	if !strings.Contains(output, "Headers:") {
		t.Error("expected Headers section in output")
	}
	if !strings.Contains(output, "Content-Type") {
		t.Error("expected Content-Type header in output")
	}
	if !strings.Contains(output, "Timing:") {
		t.Error("expected Timing info in output")
	}
	if !strings.Contains(output, "Size:") {
		t.Error("expected Size info in output")
	}
}

// TestFormatResponsePrettyPrintJSON tests JSON pretty-printing
func TestFormatResponsePrettyPrintJSON(t *testing.T) {
	formatter := NewFormatter(false)

	resp := &request.Response{
		StatusCode: 200,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body:     []byte(`{"name":"John","age":30,"email":"john@example.com"}`),
		Duration: 100 * time.Millisecond,
		Size:     50,
	}

	output := formatter.FormatResponse(resp, false)

	// Check if output contains formatted JSON (with indentation)
	if !strings.Contains(output, "name") {
		t.Error("expected JSON keys in output")
	}
	if !strings.Contains(output, "John") {
		t.Error("expected JSON values in output")
	}
	// Pretty-printed JSON should have newlines
	if !strings.Contains(output, "\n") {
		t.Error("expected formatted JSON with newlines")
	}
}

// TestFormatResponseInvalidJSON tests handling of invalid JSON
func TestFormatResponseInvalidJSON(t *testing.T) {
	formatter := NewFormatter(false)

	invalidJSON := []byte(`{invalid json content`)
	resp := &request.Response{
		StatusCode: 200,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body:     invalidJSON,
		Duration: 100 * time.Millisecond,
		Size:     len(invalidJSON),
	}

	output := formatter.FormatResponse(resp, false)

	// Should return raw body when JSON is invalid
	if !strings.Contains(output, "{invalid json content") {
		t.Error("expected raw body for invalid JSON")
	}
}

// TestFormatResponseNonJSONContent tests non-JSON response
func TestFormatResponseNonJSONContent(t *testing.T) {
	formatter := NewFormatter(false)

	htmlBody := []byte(`<html><body>Hello World</body></html>`)
	resp := &request.Response{
		StatusCode: 200,
		Headers: map[string][]string{
			"Content-Type": {"text/html"},
		},
		Body:     htmlBody,
		Duration: 100 * time.Millisecond,
		Size:     len(htmlBody),
	}

	output := formatter.FormatResponse(resp, false)

	// Should return raw HTML
	if !strings.Contains(output, "<html>") {
		t.Error("expected raw HTML in output")
	}
}

// TestFormatResponseEmptyBody tests response with empty body
func TestFormatResponseEmptyBody(t *testing.T) {
	formatter := NewFormatter(false)

	resp := &request.Response{
		StatusCode: 204,
		Headers:    make(map[string][]string),
		Body:       []byte{},
		Duration:   100 * time.Millisecond,
		Size:       0,
	}

	output := formatter.FormatResponse(resp, true)

	if !strings.Contains(output, "204") {
		t.Error("expected status code in output")
	}
}

// TestColorizeStatusGreen tests status coloring for 2xx codes
func TestColorizeStatusGreen(t *testing.T) {
	formatter := NewFormatter(true)

	statusCodes := []int{200, 201, 204, 299}
	for _, code := range statusCodes {
		output := formatter.colorizeStatus(code)

		// Should contain green ANSI code
		if !strings.Contains(output, "\033[32m") {
			t.Errorf("expected green color for status %d, got: %s", code, output)
		}
	}
}

// TestColorizeStatusBlue tests status coloring for 3xx codes
func TestColorizeStatusBlue(t *testing.T) {
	formatter := NewFormatter(true)

	statusCodes := []int{300, 301, 302, 399}
	for _, code := range statusCodes {
		output := formatter.colorizeStatus(code)

		// Should contain blue ANSI code
		if !strings.Contains(output, "\033[34m") {
			t.Errorf("expected blue color for status %d, got: %s", code, output)
		}
	}
}

// TestColorizeStatusYellow tests status coloring for 4xx codes
func TestColorizeStatusYellow(t *testing.T) {
	formatter := NewFormatter(true)

	statusCodes := []int{400, 401, 404, 499}
	for _, code := range statusCodes {
		output := formatter.colorizeStatus(code)

		// Should contain yellow ANSI code
		if !strings.Contains(output, "\033[33m") {
			t.Errorf("expected yellow color for status %d, got: %s", code, output)
		}
	}
}

// TestColorizeStatusRed tests status coloring for 5xx codes
func TestColorizeStatusRed(t *testing.T) {
	formatter := NewFormatter(true)

	statusCodes := []int{500, 502, 503, 599}
	for _, code := range statusCodes {
		output := formatter.colorizeStatus(code)

		// Should contain red ANSI code
		if !strings.Contains(output, "\033[31m") {
			t.Errorf("expected red color for status %d, got: %s", code, output)
		}
	}
}

// TestColorizeStatusNonTTY tests that non-TTY doesn't color
func TestColorizeStatusNonTTY(t *testing.T) {
	formatter := NewFormatter(false)

	statusCodes := []int{200, 404, 500}
	for _, code := range statusCodes {
		output := formatter.colorizeStatus(code)

		// Should not contain ANSI codes
		if strings.Contains(output, "\033[") {
			t.Errorf("expected no ANSI codes for non-TTY status %d, got: %s", code, output)
		}
		// Should just be the number
		if !strings.Contains(output, string(rune('0'+code/100))) {
			t.Errorf("expected status code number in output, got: %s", output)
		}
	}
}

// TestGetContentType tests extracting content type from headers
func TestGetContentType(t *testing.T) {
	formatter := NewFormatter(false)

	tests := []struct {
		name     string
		headers  map[string][]string
		expected string
	}{
		{
			name: "JSON content type",
			headers: map[string][]string{
				"Content-Type": {"application/json"},
			},
			expected: "application/json",
		},
		{
			name: "HTML content type",
			headers: map[string][]string{
				"Content-Type": {"text/html; charset=utf-8"},
			},
			expected: "text/html; charset=utf-8",
		},
		{
			name: "Case insensitive header name",
			headers: map[string][]string{
				"content-type": {"application/xml"},
			},
			expected: "application/xml",
		},
		{
			name: "Multiple header values",
			headers: map[string][]string{
				"Content-Type": {"application/json", "text/plain"},
			},
			expected: "application/json",
		},
		{
			name:     "Missing content type",
			headers:  map[string][]string{},
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := formatter.getContentType(test.headers)
			if output != test.expected {
				t.Errorf("expected %q, got %q", test.expected, output)
			}
		})
	}
}

// TestPrettyPrintJSONValid tests pretty-printing valid JSON
func TestPrettyPrintJSONValid(t *testing.T) {
	formatter := NewFormatter(false)

	jsonBody := []byte(`{"name":"John","age":30,"active":true}`)
	output := formatter.prettyPrintJSON(jsonBody)

	// Should contain all fields
	if !strings.Contains(output, "name") {
		t.Error("expected 'name' in pretty-printed JSON")
	}
	if !strings.Contains(output, "John") {
		t.Error("expected 'John' in pretty-printed JSON")
	}
	if !strings.Contains(output, "age") {
		t.Error("expected 'age' in pretty-printed JSON")
	}
	if !strings.Contains(output, "30") {
		t.Error("expected '30' in pretty-printed JSON")
	}

	// Should have indentation (newlines)
	if !strings.Contains(output, "\n") {
		t.Error("expected indented JSON with newlines")
	}
}

// TestPrettyPrintJSONInvalid tests pretty-printing invalid JSON
func TestPrettyPrintJSONInvalid(t *testing.T) {
	formatter := NewFormatter(false)

	jsonBody := []byte(`{invalid json}`)
	output := formatter.prettyPrintJSON(jsonBody)

	// Should return raw body for invalid JSON
	if output != string(jsonBody) {
		t.Errorf("expected raw body for invalid JSON, got: %s", output)
	}
}

// TestPrettyPrintJSONArray tests pretty-printing JSON array
func TestPrettyPrintJSONArray(t *testing.T) {
	formatter := NewFormatter(false)

	jsonBody := []byte(`[{"id":1},{"id":2}]`)
	output := formatter.prettyPrintJSON(jsonBody)

	if !strings.Contains(output, "id") {
		t.Error("expected 'id' in pretty-printed JSON array")
	}
	if !strings.Contains(output, "1") {
		t.Error("expected '1' in pretty-printed JSON array")
	}
}

// TestFormatBodyJSONContentType tests formatBody with JSON content type
func TestFormatBodyJSONContentType(t *testing.T) {
	formatter := NewFormatter(false)

	body := []byte(`{"test":"data"}`)
	headers := map[string][]string{
		"Content-Type": {"application/json"},
	}

	output := formatter.formatBody(body, headers)

	// Should be formatted
	if !strings.Contains(output, "\n") {
		t.Error("expected formatted JSON with indentation")
	}
}

// TestFormatBodyNonJSONContentType tests formatBody with non-JSON content type
func TestFormatBodyNonJSONContentType(t *testing.T) {
	formatter := NewFormatter(false)

	body := []byte(`This is plain text`)
	headers := map[string][]string{
		"Content-Type": {"text/plain"},
	}

	output := formatter.formatBody(body, headers)

	// Should be returned as-is
	if output != string(body) {
		t.Errorf("expected raw body, got: %s", output)
	}
}

// TestFormatBodyNoContentType tests formatBody without content type
func TestFormatBodyNoContentType(t *testing.T) {
	formatter := NewFormatter(false)

	body := []byte(`Some data`)
	headers := make(map[string][]string)

	output := formatter.formatBody(body, headers)

	// Should be returned as-is without JSON formatting
	if output != string(body) {
		t.Errorf("expected raw body, got: %s", output)
	}
}

// TestFormatResponseWithMultipleHeaders tests response with multiple header values
func TestFormatResponseWithMultipleHeaders(t *testing.T) {
	formatter := NewFormatter(false)

	resp := &request.Response{
		StatusCode: 200,
		Headers: map[string][]string{
			"Set-Cookie":   {"session=abc123", "token=xyz789"},
			"Content-Type": {"application/json"},
		},
		Body:     []byte(`{}`),
		Duration: 100 * time.Millisecond,
		Size:     2,
	}

	output := formatter.FormatResponse(resp, true)

	// Should include all header values
	if !strings.Contains(output, "Set-Cookie") {
		t.Error("expected Set-Cookie header in output")
	}
	if !strings.Contains(output, "session=abc123") {
		t.Error("expected first cookie value in output")
	}
	if !strings.Contains(output, "token=xyz789") {
		t.Error("expected second cookie value in output")
	}
}

// TestColorizeJSONNonTTY tests that colorizeJSON doesn't add colors for non-TTY
func TestColorizeJSONNonTTY(t *testing.T) {
	formatter := NewFormatter(false)

	jsonStr := `{"test": "value"}`
	output := formatter.colorizeJSON(jsonStr)

	// Non-TTY should return unchanged
	if output != jsonStr {
		t.Errorf("expected unchanged JSON for non-TTY, got: %s", output)
	}
}
