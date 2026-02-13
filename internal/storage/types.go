package storage

import "time"

// SavedCall represents a saved HTTP request
type SavedCall struct {
	Name        string            `yaml:"name"`
	Method      string            `yaml:"method"`
	URL         string            `yaml:"url"`
	Headers     map[string]string `yaml:"headers"`
	QueryParams map[string]string `yaml:"queryParams"`
	Body        string            `yaml:"body"`
	Description string            `yaml:"description"`
	CreatedAt   string            `yaml:"createdAt"`
}

// NewSavedCall creates a new saved call
func NewSavedCall(name, method, url string, headers, queryParams map[string]string, body string) *SavedCall {
	return &SavedCall{
		Name:        name,
		Method:      method,
		URL:         url,
		Headers:     headers,
		QueryParams: queryParams,
		Body:        body,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
}
