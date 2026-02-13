package request

import (
	"testing"
)

func TestExtractPathVars(t *testing.T) {
	tmpl := NewTemplate("/users/{userId}/posts/{postId}")
	vars := tmpl.ExtractPathVars()

	if len(vars) != 2 {
		t.Errorf("expected 2 variables, got %d", len(vars))
	}

	if vars[0] != "userId" {
		t.Errorf("expected first var 'userId', got '%s'", vars[0])
	}

	if vars[1] != "postId" {
		t.Errorf("expected second var 'postId', got '%s'", vars[1])
	}
}

func TestExtractEnvVars(t *testing.T) {
	tmpl := NewTemplate("${API_BASE}/users/${API_VERSION}")
	vars := tmpl.ExtractEnvVars()

	if len(vars) != 2 {
		t.Errorf("expected 2 variables, got %d", len(vars))
	}

	if vars[0] != "API_BASE" {
		t.Errorf("expected first var 'API_BASE', got '%s'", vars[0])
	}

	if vars[1] != "API_VERSION" {
		t.Errorf("expected second var 'API_VERSION', got '%s'", vars[1])
	}
}

func TestResolvePathVars(t *testing.T) {
	tmpl := NewTemplate("/users/{userId}/posts/{postId}")
	tmpl.SetPathVars(map[string]string{
		"userId": "123",
		"postId": "456",
	})

	resolved, err := tmpl.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "/users/123/posts/456"
	if resolved != expected {
		t.Errorf("expected '%s', got '%s'", expected, resolved)
	}
}

func TestResolveEnvVars(t *testing.T) {
	tmpl := NewTemplate("${API_BASE}/users")
	tmpl.SetEnvVars(map[string]string{
		"API_BASE": "https://api.example.com",
	})

	resolved, err := tmpl.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "https://api.example.com/users"
	if resolved != expected {
		t.Errorf("expected '%s', got '%s'", expected, resolved)
	}
}

func TestResolveMixed(t *testing.T) {
	tmpl := NewTemplate("${API_BASE}/users/{userId}/posts/{postId}")
	tmpl.SetEnvVars(map[string]string{
		"API_BASE": "https://api.example.com",
	})
	tmpl.SetPathVars(map[string]string{
		"userId": "123",
		"postId": "456",
	})

	resolved, err := tmpl.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "https://api.example.com/users/123/posts/456"
	if resolved != expected {
		t.Errorf("expected '%s', got '%s'", expected, resolved)
	}
}

func TestResolveMissingPathVar(t *testing.T) {
	tmpl := NewTemplate("/users/{userId}")
	tmpl.SetPathVars(map[string]string{})

	_, err := tmpl.Resolve()
	if err == nil {
		t.Fatal("expected error for missing path variable")
	}
}

func TestResolveMissingEnvVar(t *testing.T) {
	tmpl := NewTemplate("${API_BASE}/users")
	tmpl.SetEnvVars(map[string]string{})

	_, err := tmpl.Resolve()
	if err == nil {
		t.Fatal("expected error for missing environment variable")
	}
}

func TestDuplicatePathVars(t *testing.T) {
	tmpl := NewTemplate("/users/{id}/posts/{id}")
	vars := tmpl.ExtractPathVars()

	// Should only return unique variables
	if len(vars) != 1 {
		t.Errorf("expected 1 unique variable, got %d", len(vars))
	}

	if vars[0] != "id" {
		t.Errorf("expected 'id', got '%s'", vars[0])
	}
}
