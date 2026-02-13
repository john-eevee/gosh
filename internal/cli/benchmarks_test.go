package cli

import (
	"testing"
)

// BenchmarkParseSimpleRequest benchmarks simple request parsing
func BenchmarkParseSimpleRequest(b *testing.B) {
	args := []string{"get", "https://api.example.com/users"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(args)
		_, _ = parser.Parse()
	}
}

// BenchmarkParseComplexRequest benchmarks complex request parsing
func BenchmarkParseComplexRequest(b *testing.B) {
	args := []string{
		"post", "https://api.example.com/users/{userId}",
		"-H", "Authorization:Bearer token",
		"-H", "Content-Type:application/json",
		"-H", "X-Custom:value",
		"-d", `{"name":"John","email":"john@example.com"}`,
		"--save", "create-user",
		"--info",
		"userId=123",
		"limit==10",
		"offset==20",
		"--auth=my-api",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(args)
		_, _ = parser.Parse()
	}
}

// BenchmarkParseWithManyHeaders benchmarks parsing with many headers
func BenchmarkParseWithManyHeaders(b *testing.B) {
	args := []string{
		"get", "https://api.example.com/users",
		"-H", "Authorization:Bearer token",
		"-H", "Content-Type:application/json",
		"-H", "X-Custom-1:value1",
		"-H", "X-Custom-2:value2",
		"-H", "X-Custom-3:value3",
		"-H", "X-Custom-4:value4",
		"-H", "X-Custom-5:value5",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(args)
		_, _ = parser.Parse()
	}
}

// BenchmarkParseRecall benchmarks recall command parsing
func BenchmarkParseRecall(b *testing.B) {
	args := []string{
		"recall", "get-users",
		"-H", "Authorization:Bearer new-token",
		"userId=456",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(args)
		_, _ = parser.Parse()
	}
}

// BenchmarkParseAuthCommand benchmarks auth command parsing
func BenchmarkParseAuthCommand(b *testing.B) {
	args := []string{
		"auth", "add", "bearer", "my-api",
		"token=abc123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(args)
		_, _ = parser.Parse()
	}
}
