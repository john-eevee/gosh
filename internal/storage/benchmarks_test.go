package storage

import (
	"testing"
)

// BenchmarkSaveCall benchmarks saving a call
func BenchmarkSaveCall(b *testing.B) {
	tmpDir := b.TempDir()
	mgr := NewManager(tmpDir)

	call := NewSavedCall(
		"test-call",
		"POST",
		"https://api.example.com/users",
		map[string]string{"Content-Type": "application/json"},
		map[string]string{"limit": "10"},
		`{"name":"John"}`,
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mgr.Save(call)
	}
}

// BenchmarkLoadCall benchmarks loading a call
func BenchmarkLoadCall(b *testing.B) {
	tmpDir := b.TempDir()
	mgr := NewManager(tmpDir)

	call := NewSavedCall(
		"test-call",
		"GET",
		"https://api.example.com/users",
		nil,
		nil,
		"",
	)
	_ = mgr.Save(call)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mgr.Load("test-call")
	}
}

// BenchmarkListCalls benchmarks listing calls
func BenchmarkListCalls(b *testing.B) {
	tmpDir := b.TempDir()
	mgr := NewManager(tmpDir)

	// Create multiple calls
	for i := 0; i < 10; i++ {
		call := NewSavedCall(
			"call-"+string(rune(i)),
			"GET",
			"https://api.example.com/endpoint",
			nil,
			nil,
			"",
		)
		_ = mgr.Save(call)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mgr.List()
	}
}

// BenchmarkDeleteCall benchmarks deleting a call
func BenchmarkDeleteCall(b *testing.B) {
	tmpDir := b.TempDir()
	mgr := NewManager(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		call := NewSavedCall(
			"call-to-delete",
			"GET",
			"https://api.example.com",
			nil,
			nil,
			"",
		)
		_ = mgr.Save(call)
		_ = mgr.Delete("call-to-delete")
	}
}

// BenchmarkExistsCall benchmarks checking if call exists
func BenchmarkExistsCall(b *testing.B) {
	tmpDir := b.TempDir()
	mgr := NewManager(tmpDir)

	call := NewSavedCall("test-call", "GET", "https://api.example.com", nil, nil, "")
	_ = mgr.Save(call)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.Exists("test-call")
	}
}
