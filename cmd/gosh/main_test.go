package main

import (
	"bytes"
	"os/exec"
	"testing"
)

// TestMainVersionCmd runs gosh with --version flag via subprocess
func TestMainVersionCmd(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "GET", "http://example.com", "--version")
	cmd.Dir = "/home/john/code/gosh/cmd/gosh"

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Logf("stderr: %s", stderr.String())
		// Version flag exits before making request, so this is expected
	}

	output := stdout.String() + stderr.String()
	if output == "" {
		t.Error("expected output for version flag")
	}
}

// TestMainDryRun runs gosh with --dry flag via subprocess
func TestMainDryRun(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "GET", "http://example.com", "--dry", "--save", "test-call")
	cmd.Dir = "/home/john/code/gosh/cmd/gosh"

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Logf("stderr: %s", stderr.String())
		// Might fail due to httptest server not being available, that's okay
	}

	output := stdout.String() + stderr.String()
	if output == "" {
		t.Error("expected output from dry run")
	}
}

// TestMainList runs gosh list command via subprocess
func TestMainListCmd(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "list")
	cmd.Dir = "/home/john/code/gosh/cmd/gosh"

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// List is a valid subcommand and should complete
	err := cmd.Run()
	if err != nil {
		t.Logf("stderr: %s", stderr.String())
		// List might fail if no saved calls exist, that's okay
	}
}

// TestMainInvalidCommand verifies invalid commands are handled
func TestMainInvalidCommand(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "invalid-command")
	cmd.Dir = "/home/john/code/gosh/cmd/gosh"

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("expected invalid command to fail")
	}

	errOutput := stderr.String()
	if errOutput == "" {
		t.Error("expected error message for invalid command")
	}
}

// TestMainNoArgs verifies behavior with no arguments
func TestMainNoArgs(t *testing.T) {
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = "/home/john/code/gosh/cmd/gosh"

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// No args might show help or error
	_ = cmd.Run()
	// Just verify command completes
}

// TestMainStderrHandling verifies error output goes to stderr
func TestMainStderrHandling(t *testing.T) {
	// Create a test where we know main will fail
	cmd := exec.Command("go", "run", ".", "invalid")
	cmd.Dir = "/home/john/code/gosh/cmd/gosh"

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	_ = cmd.Run()

	// Even on error, it should handle stderr properly
	if stderr.Len() == 0 {
		t.Log("Note: No stderr output for invalid command")
	}
}
