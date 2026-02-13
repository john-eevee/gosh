package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gosh/internal/request"
)

// Formatter handles response formatting
type Formatter struct {
	isTTY bool
}

// NewFormatter creates a new formatter
func NewFormatter(isTTY bool) *Formatter {
	return &Formatter{isTTY: isTTY}
}

// FormatResponse formats the response for display
func (f *Formatter) FormatResponse(resp *request.Response, showInfo bool) string {
	var output strings.Builder

	// Status line
	statusColor := f.colorizeStatus(resp.StatusCode)
	output.WriteString(fmt.Sprintf("%s\n", statusColor))

	if showInfo {
		// Headers
		output.WriteString("\nHeaders:\n")
		for key, values := range resp.Headers {
			for _, val := range values {
				output.WriteString(fmt.Sprintf("  %s: %s\n", key, val))
			}
		}

		// Timing and size info
		output.WriteString(fmt.Sprintf("\nTiming: %v\n", resp.Duration))
		output.WriteString(fmt.Sprintf("Size: %d bytes\n", resp.Size))
	}

	// Body
	if len(resp.Body) > 0 {
		output.WriteString("\n")
		formattedBody := f.formatBody(resp.Body, resp.Headers)
		output.WriteString(formattedBody)
	}

	return output.String()
}

// formatBody attempts to pretty-print JSON, otherwise returns raw
func (f *Formatter) formatBody(body []byte, headers map[string][]string) string {
	// Check content type
	contentType := f.getContentType(headers)

	if strings.Contains(contentType, "application/json") {
		return f.prettyPrintJSON(body)
	}

	return string(body)
}

// prettyPrintJSON attempts to pretty-print JSON
func (f *Formatter) prettyPrintJSON(body []byte) string {
	var obj interface{}
	if err := json.Unmarshal(body, &obj); err != nil {
		// Not valid JSON, return as-is
		return string(body)
	}

	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return string(body)
	}

	if f.isTTY {
		return f.colorizeJSON(string(pretty))
	}

	return string(pretty)
}

// colorizeStatus returns colored status text if TTY, else plain
func (f *Formatter) colorizeStatus(statusCode int) string {
	status := fmt.Sprintf("%d", statusCode)
	if !f.isTTY {
		return status
	}

	switch {
	case statusCode < 300:
		return fmt.Sprintf("\033[32m%d\033[0m", statusCode) // Green
	case statusCode < 400:
		return fmt.Sprintf("\033[34m%d\033[0m", statusCode) // Blue
	case statusCode < 500:
		return fmt.Sprintf("\033[33m%d\033[0m", statusCode) // Yellow
	default:
		return fmt.Sprintf("\033[31m%d\033[0m", statusCode) // Red
	}
}

// colorizeJSON adds basic syntax coloring to JSON
func (f *Formatter) colorizeJSON(jsonStr string) string {
	// Simple approach: color keys and values differently
	// For production, consider a JSON-specific coloring library
	return jsonStr
}

// getContentType extracts content type from headers
func (f *Formatter) getContentType(headers map[string][]string) string {
	for key, values := range headers {
		if strings.EqualFold(key, "content-type") && len(values) > 0 {
			return values[0]
		}
	}
	return ""
}
