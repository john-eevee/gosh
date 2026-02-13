package cli

import (
	"testing"
)

func TestParseGetRequest(t *testing.T) {
	parser := NewParser([]string{"get", "https://api.example.com/users"})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req, ok := result.(*ParsedRequest)
	if !ok {
		t.Fatalf("expected ParsedRequest, got %T", result)
	}

	if req.Method != "GET" {
		t.Errorf("expected method GET, got %s", req.Method)
	}

	if req.URL != "https://api.example.com/users" {
		t.Errorf("expected URL https://api.example.com/users, got %s", req.URL)
	}
}

func TestParseWithHeaders(t *testing.T) {
	parser := NewParser([]string{
		"post", "https://api.example.com/users",
		"-H", "Authorization:Bearer token",
		"-H", "Content-Type:application/json",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if len(req.Headers) != 2 {
		t.Errorf("expected 2 headers, got %d", len(req.Headers))
	}

	if req.Headers["Authorization"] != "Bearer token" {
		t.Errorf("expected Authorization header 'Bearer token', got '%s'", req.Headers["Authorization"])
	}

	if req.Headers["Content-Type"] != "application/json" {
		t.Errorf("expected Content-Type header 'application/json', got '%s'", req.Headers["Content-Type"])
	}
}

func TestParseWithBody(t *testing.T) {
	parser := NewParser([]string{
		"post", "https://api.example.com/users",
		"-d", `{"name":"John"}`,
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Body != `{"name":"John"}` {
		t.Errorf("expected body '{\"name\":\"John\"}', got '%s'", req.Body)
	}
}

func TestParseWithSaveFlag(t *testing.T) {
	parser := NewParser([]string{
		"post", "https://api.example.com/users",
		"--save", "create-user",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.Save != "create-user" {
		t.Errorf("expected Save 'create-user', got '%s'", req.Save)
	}
}

func TestParseWithPathParams(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users/{userId}",
		"userId=123",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.PathParams["userId"] != "123" {
		t.Errorf("expected userId param '123', got '%s'", req.PathParams["userId"])
	}
}

func TestParseRecall(t *testing.T) {
	parser := NewParser([]string{
		"recall", "create-user",
		"userId=456",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	opts, ok := result.(*RecallOptions)
	if !ok {
		t.Fatalf("expected RecallOptions, got %T", result)
	}

	if opts.Name != "create-user" {
		t.Errorf("expected name 'create-user', got '%s'", opts.Name)
	}

	if opts.ParameterOverride["userId"] != "456" {
		t.Errorf("expected userId override '456', got '%s'", opts.ParameterOverride["userId"])
	}
}

func TestParseWithQueryParams(t *testing.T) {
	parser := NewParser([]string{
		"get", "https://api.example.com/users",
		"limit==10",
		"offset==20",
	})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := result.(*ParsedRequest)
	if req.QueryParams["limit"] != "10" {
		t.Errorf("expected limit query param '10', got '%s'", req.QueryParams["limit"])
	}

	if req.QueryParams["offset"] != "20" {
		t.Errorf("expected offset query param '20', got '%s'", req.QueryParams["offset"])
	}
}

func TestParseInvalidMethod(t *testing.T) {
	parser := NewParser([]string{
		"invalid", "https://api.example.com/users",
	})
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error for invalid method")
	}
}

func TestParseList(t *testing.T) {
	parser := NewParser([]string{"list"})
	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "list" {
		t.Errorf("expected 'list', got %v", result)
	}
}
