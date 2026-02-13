package request

import (
	"fmt"
	"regexp"
	"strings"
)

// Template handles path and variable templating
type Template struct {
	text     string
	pathVars map[string]string
	envVars  map[string]string
}

// NewTemplate creates a new template
func NewTemplate(text string) *Template {
	return &Template{
		text:     text,
		pathVars: make(map[string]string),
		envVars:  make(map[string]string),
	}
}

// SetPathVars sets path variables (for {var} substitution)
func (t *Template) SetPathVars(vars map[string]string) {
	t.pathVars = vars
}

// SetEnvVars sets environment variables (for ${VAR} substitution)
func (t *Template) SetEnvVars(vars map[string]string) {
	t.envVars = vars
}

// ExtractPathVars extracts all path variables from the template
func (t *Template) ExtractPathVars() []string {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(t.text, -1)

	var vars []string
	seen := make(map[string]bool)
	for _, match := range matches {
		varName := match[1]
		if !seen[varName] {
			vars = append(vars, varName)
			seen[varName] = true
		}
	}

	return vars
}

// ExtractEnvVars extracts all environment variables from the template
func (t *Template) ExtractEnvVars() []string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(t.text, -1)

	var vars []string
	seen := make(map[string]bool)
	for _, match := range matches {
		varName := match[1]
		if !seen[varName] {
			vars = append(vars, varName)
			seen[varName] = true
		}
	}

	return vars
}

// Resolve resolves all variables in the template
func (t *Template) Resolve() (string, error) {
	result := t.text

	// Substitute environment variables first
	envVars := t.ExtractEnvVars()
	for _, varName := range envVars {
		val, ok := t.envVars[varName]
		if !ok {
			return "", fmt.Errorf("environment variable not found: %s", varName)
		}
		placeholder := fmt.Sprintf("${%s}", varName)
		result = strings.ReplaceAll(result, placeholder, val)
	}

	// Check for unresolved path variables
	pathVars := t.ExtractPathVars()
	unresolved := []string{}
	for _, varName := range pathVars {
		if _, ok := t.pathVars[varName]; !ok {
			unresolved = append(unresolved, varName)
		}
	}

	if len(unresolved) > 0 {
		return "", fmt.Errorf("missing template variables: %v", unresolved)
	}

	// Substitute path variables
	for _, varName := range pathVars {
		val := t.pathVars[varName]
		placeholder := fmt.Sprintf("{%s}", varName)
		result = strings.ReplaceAll(result, placeholder, val)
	}

	return result, nil
}
