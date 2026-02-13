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
}

// RecallOptions holds options for recall command
type RecallOptions struct {
	Name              string
	ParameterOverride map[string]string
	Headers           map[string]string
	Env               string
}
