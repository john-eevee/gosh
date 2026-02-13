package request

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Builder constructs an http.Request from Request details
type Builder struct {
	req *Request
}

// NewBuilder creates a new request builder
func NewBuilder(req *Request) *Builder {
	return &Builder{req: req}
}

// Build constructs an http.Request
func (b *Builder) Build() (*http.Request, error) {
	// Parse URL and add query parameters
	u, err := url.Parse(b.req.URL)
	if err != nil {
		return nil, err
	}

	// Add query parameters if any
	if len(b.req.QueryParams) > 0 {
		q := u.Query()
		for key, val := range b.req.QueryParams {
			q.Add(key, val)
		}
		u.RawQuery = q.Encode()
	}

	// Create request with body
	var bodyReader io.Reader
	if b.req.Body != "" {
		bodyReader = strings.NewReader(b.req.Body)
	}

	httpReq, err := http.NewRequest(b.req.Method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, val := range b.req.Headers {
		httpReq.Header.Set(key, val)
	}

	// Apply authentication if provided
	if b.req.Auth != nil {
		if err := b.req.Auth.Apply(httpReq); err != nil {
			return nil, err
		}
	}

	return httpReq, nil
}
