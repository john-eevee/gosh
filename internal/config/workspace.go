package config

import (
	"os"
	"path/filepath"
)

// DetectWorkspace detects and loads the current workspace
func DetectWorkspace() (*Workspace, error) {
	// Start from current directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Look for .gosh.yaml or git root
	root, err := findWorkspaceRoot(cwd)
	if err != nil {
		return nil, err
	}

	workspace := &Workspace{
		Root: root,
		Env:  make(map[string]string),
	}

	// Try loading .gosh.yaml
	configPath := filepath.Join(root, ".gosh.yaml")
	if _, err := os.Stat(configPath); err == nil {
		config, err := LoadWorkspaceConfig(configPath)
		if err != nil {
			return nil, err
		}
		workspace.Config = config
	}

	// Load .env file if it exists
	envPath := filepath.Join(root, ".env")
	if _, err := os.Stat(envPath); err == nil {
		envVars, err := LoadEnvFile(envPath)
		if err != nil {
			return nil, err
		}
		workspace.Env = envVars
	}

	return workspace, nil
}

// findWorkspaceRoot finds the workspace root by looking for .gosh.yaml or git root
func findWorkspaceRoot(startPath string) (string, error) {
	current := startPath

	// Look for .gosh.yaml
	for {
		goshPath := filepath.Join(current, ".gosh.yaml")
		if _, err := os.Stat(goshPath); err == nil {
			return current, nil
		}

		// Look for .git directory
		gitPath := filepath.Join(current, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			// Reached root directory, return current directory as workspace
			return startPath, nil
		}
		current = parent
	}
}
