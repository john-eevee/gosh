package request

import (
	"io"
	"net/http"
	"time"
)

// Executor executes HTTP requests
type Executor struct {
	client  *http.Client
	timeout time.Duration
}

// NewExecutor creates a new request executor
func NewExecutor(timeout time.Duration) *Executor {
	return &Executor{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// Execute executes an HTTP request and returns the response
func (e *Executor) Execute(req *Request) (*Response, error) {
	builder := NewBuilder(req)
	httpReq, err := builder.Build()
	if err != nil {
		return nil, err
	}

	start := time.Now()
	httpResp, err := e.client.Do(httpReq)
	duration := time.Since(start)

	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	// Read response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		StatusCode: httpResp.StatusCode,
		Status:     httpResp.Status,
		Headers:    httpResp.Header,
		Body:       body,
		Duration:   duration,
		Size:       len(body),
	}

	return resp, nil
}
