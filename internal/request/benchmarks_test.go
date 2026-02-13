package request

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// BenchmarkBuilderBuild benchmarks request building
func BenchmarkBuilderBuild(b *testing.B) {
	req := &Request{
		Method: "POST",
		URL:    "https://api.example.com/users",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token123",
			"X-Custom":      "value",
		},
		QueryParams: map[string]string{
			"limit":  "10",
			"offset": "0",
		},
		Body: `{"name":"John","email":"john@example.com"}`,
	}

	builder := NewBuilder(req)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.Build()
	}
}

// BenchmarkExecutorExecute benchmarks request execution
func BenchmarkExecutorExecute(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	executor := NewExecutor(5 * time.Second)
	req := &Request{
		Method:  "GET",
		URL:     server.URL,
		Headers: make(map[string]string),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executor.Execute(req)
	}
}

// BenchmarkTemplateResolve benchmarks template resolution
func BenchmarkTemplateResolve(b *testing.B) {
	tmpl := NewTemplate("${API_BASE}/users/{userId}/posts/{postId}")
	tmpl.SetEnvVars(map[string]string{
		"API_BASE": "https://api.example.com",
	})
	tmpl.SetPathVars(map[string]string{
		"userId": "123",
		"postId": "456",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmpl.Resolve()
	}
}

// BenchmarkExtractPathVars benchmarks path variable extraction
func BenchmarkExtractPathVars(b *testing.B) {
	url := "/users/{userId}/posts/{postId}/comments/{commentId}"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmpl := NewTemplate(url)
		tmpl.ExtractPathVars()
	}
}

// BenchmarkExtractEnvVars benchmarks environment variable extraction
func BenchmarkExtractEnvVars(b *testing.B) {
	url := "${API_BASE}/api/${API_VERSION}/users/${USER_ID}"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmpl := NewTemplate(url)
		tmpl.ExtractEnvVars()
	}
}
