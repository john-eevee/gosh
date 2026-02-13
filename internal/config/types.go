package config

// GlobalConfig represents global gosh configuration
type GlobalConfig struct {
	DefaultEnvironment string `yaml:"defaultEnvironment"`
	PrettyPrint        bool   `yaml:"prettyPrint"`
	Timeout            string `yaml:"timeout"`
	UserAgent          string `yaml:"userAgent"`
}

// WorkspaceConfig represents workspace-level configuration
type WorkspaceConfig struct {
	Name           string                       `yaml:"name"`
	BaseURL        string                       `yaml:"baseUrl"`
	DefaultHeaders map[string]string            `yaml:"defaultHeaders"`
	Environments   map[string]map[string]string `yaml:"environments"`
}

// Workspace holds information about the current workspace
type Workspace struct {
	Root   string            // Root directory of workspace
	Config *WorkspaceConfig  // Loaded config
	Env    map[string]string // Merged environment variables
}
