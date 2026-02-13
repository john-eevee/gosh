package cli

// ParsedRequest holds all parsed CLI arguments
type ParsedRequest struct {
	Method       string
	URL          string
	Headers      map[string]string
	QueryParams  map[string]string
	Body         string
	PathParams   map[string]string // {var} style parameters
	HasStdinBody bool
	// Flags
	Save          string // Name to save as
	Dry           bool   // Don't execute
	Info          bool   // Show full response info
	NoInteractive bool   // Don't prompt for missing vars
	Env           string // Environment to use
	Format        string // Output format
	Auth          string // Authentication preset to use (format: "type:name")
}

// RecallOptions holds options for recall command
type RecallOptions struct {
	Name              string
	ParameterOverride map[string]string
	Headers           map[string]string
	Env               string
}

// AuthCommand holds auth subcommand details
type AuthCommand struct {
	Subcommand string            // "add", "remove", "list"
	Type       string            // Auth type: "basic", "bearer", "custom"
	Name       string            // Preset name
	Flags      map[string]string // Additional flags for add/remove
}
