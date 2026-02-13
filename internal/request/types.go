package request

import "time"

// Request holds HTTP request details
type Request struct {
	Method      string
	URL         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        string
	Timeout     time.Duration
}

// Response holds the HTTP response
type Response struct {
	StatusCode int
	Status     string
	Headers    map[string][]string
	Body       []byte
	Duration   time.Duration
	Size       int // Size in bytes
}
