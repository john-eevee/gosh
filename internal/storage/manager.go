package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Manager handles saving and loading saved calls
type Manager struct {
	workspaceRoot string
}

// NewManager creates a new storage manager
func NewManager(workspaceRoot string) *Manager {
	return &Manager{workspaceRoot: workspaceRoot}
}

// Save saves a call
func (m *Manager) Save(call *SavedCall) error {
	callsDir, err := GetCallsDir(m.workspaceRoot)
	if err != nil {
		return err
	}

	callPath := filepath.Join(callsDir, call.Name+".yaml")

	data, err := yaml.Marshal(call)
	if err != nil {
		return err
	}

	if err := os.WriteFile(callPath, data, 0644); err != nil {
		return err
	}

	return nil
}

// Load loads a saved call by name
func (m *Manager) Load(name string) (*SavedCall, error) {
	callPath := GetCallPath(m.workspaceRoot, name)

	data, err := os.ReadFile(callPath)
	if err != nil {
		return nil, fmt.Errorf("call not found: %s", name)
	}

	var call SavedCall
	if err := yaml.Unmarshal(data, &call); err != nil {
		return nil, err
	}

	return &call, nil
}

// List returns all saved calls
func (m *Manager) List() ([]*SavedCall, error) {
	callsDir, err := GetCallsDir(m.workspaceRoot)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(callsDir)
	if err != nil {
		return nil, err
	}

	var calls []*SavedCall
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		name := entry.Name()[:len(entry.Name())-5] // Remove .yaml

		call, err := m.Load(name)
		if err != nil {
			continue // Skip invalid files
		}

		calls = append(calls, call)
	}

	return calls, nil
}

// Delete deletes a saved call
func (m *Manager) Delete(name string) error {
	callPath := GetCallPath(m.workspaceRoot, name)

	if err := os.Remove(callPath); err != nil {
		return fmt.Errorf("failed to delete call: %w", err)
	}

	return nil
}

// Exists checks if a call exists
func (m *Manager) Exists(name string) bool {
	callPath := GetCallPath(m.workspaceRoot, name)
	_, err := os.Stat(callPath)
	return err == nil
}
