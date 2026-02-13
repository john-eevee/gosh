package auth

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Manager handles authentication preset operations
type Manager struct {
	configPath string
	presets    map[string]*AuthPreset
}

// NewManager creates a new auth manager for a workspace
func NewManager(workspaceRoot string) *Manager {
	return &Manager{
		configPath: filepath.Join(workspaceRoot, ".gosh", "auth.yaml"),
		presets:    make(map[string]*AuthPreset),
	}
}

// Load loads authentication presets from disk
func (m *Manager) Load() error {
	// If file doesn't exist, that's OK
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read auth config: %w", err)
	}

	var config AuthConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse auth config: %w", err)
	}

	if config.Presets != nil {
		m.presets = config.Presets
	}

	return nil
}

// Save saves authentication presets to disk
func (m *Manager) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create auth config directory: %w", err)
	}

	config := AuthConfig{
		Presets: m.presets,
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal auth config: %w", err)
	}

	if err := os.WriteFile(m.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write auth config: %w", err)
	}

	return nil
}

// Get retrieves a preset by name
func (m *Manager) Get(name string) (*AuthPreset, error) {
	preset, exists := m.presets[name]
	if !exists {
		return nil, fmt.Errorf("auth preset not found: %s", name)
	}
	return preset, nil
}

// Add adds or updates an authentication preset
func (m *Manager) Add(preset *AuthPreset) error {
	if preset.Name == "" {
		return fmt.Errorf("preset name cannot be empty")
	}
	if preset.Type == "" {
		return fmt.Errorf("preset type cannot be empty")
	}

	m.presets[preset.Name] = preset
	return m.Save()
}

// Remove removes an authentication preset
func (m *Manager) Remove(name string) error {
	if _, exists := m.presets[name]; !exists {
		return fmt.Errorf("auth preset not found: %s", name)
	}
	delete(m.presets, name)
	return m.Save()
}

// List returns all authentication presets
func (m *Manager) List() map[string]*AuthPreset {
	return m.presets
}
