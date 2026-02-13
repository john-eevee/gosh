package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// AuthType represents the type of authentication
type AuthType string

const (
	AuthTypeBasic  AuthType = "basic"
	AuthTypeBearer AuthType = "bearer"
	AuthTypeCustom AuthType = "custom"
)

// AuthPreset represents a saved authentication preset
type AuthPreset struct {
	Name     string   `yaml:"name"`
	Type     string   `yaml:"type"` // "basic", "bearer", "custom"
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Token    string   `yaml:"token"`
	Header   string   `yaml:"header"`  // Header name for custom auth
	Prefix   string   `yaml:"prefix"`  // Prefix for custom auth (e.g., "Bearer ", "token ")
	Value    string   `yaml:"value"`   // Value for custom auth
	Headers  []string `yaml:"headers"` // Additional headers to add
}

// AuthConfig holds all authentication presets
type AuthConfig struct {
	Presets map[string]*AuthPreset `yaml:"presets"`
}

// Apply adds the authentication to the given HTTP request
func (p *AuthPreset) Apply(req *http.Request) error {
	switch AuthType(strings.ToLower(p.Type)) {
	case AuthTypeBasic:
		return p.applyBasic(req)
	case AuthTypeBearer:
		return p.applyBearer(req)
	case AuthTypeCustom:
		return p.applyCustom(req)
	default:
		return fmt.Errorf("unknown auth type: %s", p.Type)
	}
}

// applyBasic applies HTTP Basic Authentication
func (p *AuthPreset) applyBasic(req *http.Request) error {
	if p.Username == "" {
		return fmt.Errorf("basic auth requires username")
	}
	// Password is optional (empty password is valid)

	credentials := base64.StdEncoding.EncodeToString(
		[]byte(p.Username + ":" + p.Password),
	)
	req.Header.Set("Authorization", "Basic "+credentials)
	return nil
}

// applyBearer applies Bearer token authentication
func (p *AuthPreset) applyBearer(req *http.Request) error {
	if p.Token == "" {
		return fmt.Errorf("bearer auth requires token")
	}
	req.Header.Set("Authorization", "Bearer "+p.Token)
	return nil
}

// applyCustom applies custom authentication headers
func (p *AuthPreset) applyCustom(req *http.Request) error {
	if p.Header == "" {
		return fmt.Errorf("custom auth requires header name")
	}

	headerName := p.Header
	headerValue := p.Value

	// Add prefix if specified
	if p.Prefix != "" {
		headerValue = p.Prefix + headerValue
	}

	req.Header.Set(headerName, headerValue)

	// Add any additional headers
	for _, h := range p.Headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			req.Header.Add(key, value)
		}
	}

	return nil
}
