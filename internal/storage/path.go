package storage

import (
	"os"
	"path/filepath"
)

// GetCallsDir returns the directory where saved calls are stored
func GetCallsDir(workspaceRoot string) (string, error) {
	callsDir := filepath.Join(workspaceRoot, ".gosh", "calls")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(callsDir, 0755); err != nil {
		return "", err
	}

	return callsDir, nil
}

// GetCallPath returns the full path to a saved call file
func GetCallPath(workspaceRoot, name string) string {
	return filepath.Join(workspaceRoot, ".gosh", "calls", name+".yaml")
}
