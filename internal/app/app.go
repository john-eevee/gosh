package app

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/gosh/internal/auth"
	"github.com/gosh/internal/cli"
	"github.com/gosh/internal/config"
	"github.com/gosh/internal/output"
	"github.com/gosh/internal/request"
	"github.com/gosh/internal/storage"
	"github.com/gosh/internal/ui"
	"github.com/mattn/go-isatty"
)

// App is the main application
type App struct {
	workspace *config.Workspace
	global    *config.GlobalConfig
	storage   *storage.Manager
	authMgr   *auth.Manager
	isTTY     bool
}

// NewApp creates a new app instance
func NewApp() (*App, error) {
	// Detect workspace
	workspace, err := config.DetectWorkspace()
	if err != nil {
		return nil, err
	}

	// Load global config
	global, err := config.LoadGlobalConfig()
	if err != nil {
		return nil, err
	}

	// Create storage manager
	storageMgr := storage.NewManager(workspace.Root)

	// Create auth manager
	authMgr := auth.NewManager(workspace.Root)
	if err := authMgr.Load(); err != nil {
		return nil, err
	}

	// Check if stdout is a TTY
	isTTY := isTerminal(os.Stdout)

	return &App{
		workspace: workspace,
		global:    global,
		storage:   storageMgr,
		authMgr:   authMgr,
		isTTY:     isTTY,
	}, nil
}

// Run runs the application
func (a *App) Run(args []string) error {
	parser := cli.NewParser(args)
	result, err := parser.Parse()
	if err != nil {
		return err
	}

	switch v := result.(type) {
	case *cli.ParsedRequest:
		return a.executeRequest(v)
	case *cli.RecallOptions:
		return a.executeRecall(v)
	case *cli.AuthCommand:
		return a.handleAuthCommand(v)
	case string:
		switch v {
		case "version":
			return a.printVersion()
		case "help":
			return a.printHelp()
		case "list":
			return a.listCalls()
		default:
			if len(v) > 7 && v[:7] == "delete:" {
				name := v[7:]
				return a.deleteCall(name)
			}
		}
	}

	return fmt.Errorf("unknown command")
}

// executeRequest executes an HTTP request
func (a *App) executeRequest(req *cli.ParsedRequest) error {
	// Apply default headers from workspace config
	if a.workspace.Config != nil && len(a.workspace.Config.DefaultHeaders) > 0 {
		// CLI headers override config defaults
		for key, val := range a.workspace.Config.DefaultHeaders {
			if _, exists := req.Headers[key]; !exists {
				req.Headers[key] = val
			}
		}
	}

	// Check for stdin body
	if req.HasStdinBody || !isTerminal(os.Stdin) {
		// Read from stdin
		stdinData, err := io.ReadAll(os.Stdin)
		if err == nil && len(stdinData) > 0 {
			req.Body = string(stdinData)
			req.HasStdinBody = true
		}
	}

	// Resolve environment variables in all parts
	req.URL = a.substituteEnvVars(req.URL)
	req.Headers = a.substituteEnvVarsInMap(req.Headers)
	req.Body = a.substituteEnvVars(req.Body)

	// Resolve template variables in URL
	tmpl := request.NewTemplate(req.URL)
	pathVars := tmpl.ExtractPathVars()

	resolvedPathVars := make(map[string]string)
	if len(pathVars) > 0 {
		for _, varName := range pathVars {
			var val string

			// Check if provided via CLI first
			if cliVal, ok := req.PathParams[varName]; ok {
				val = cliVal
			} else if !req.NoInteractive {
				// Prompt for missing path variables
				var err error
				val, err = ui.PromptForVariable(varName)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("missing required template variable: {%s}", varName)
			}

			resolvedPathVars[varName] = val
		}
		tmpl.SetPathVars(resolvedPathVars)
	}

	// Resolve environment variables
	tmpl.SetEnvVars(a.workspace.Env)

	// Resolve URL
	resolvedURL, err := tmpl.Resolve()
	if err != nil {
		return err
	}

	// Get timeout
	timeout := 30 * time.Second
	if a.global.Timeout != "" {
		parsedTimeout, err := time.ParseDuration(a.global.Timeout)
		if err == nil {
			timeout = parsedTimeout
		}
	}

	// Build request
	httpReq := &request.Request{
		Method:      req.Method,
		URL:         resolvedURL,
		Headers:     req.Headers,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Timeout:     timeout,
	}

	// Apply authentication if provided
	if req.Auth != "" {
		authPreset, err := a.authMgr.Get(req.Auth)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
		httpReq.Auth = authPreset
	}

	// If dry run, just save
	if req.Dry {
		if req.Save == "" {
			return fmt.Errorf("--dry requires --save to specify a name")
		}
		savedCall := storage.NewSavedCall(
			req.Save,
			req.Method,
			req.URL,
			req.Headers,
			req.QueryParams,
			req.Body,
		)
		if err := a.storage.Save(savedCall); err != nil {
			return err
		}
		fmt.Printf("Saved call: %s\n", req.Save)
		return nil
	}

	// Execute request
	executor := request.NewExecutor(timeout)
	resp, err := executor.Execute(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	// Save if requested
	if req.Save != "" {
		savedCall := storage.NewSavedCall(
			req.Save,
			req.Method,
			req.URL,
			req.Headers,
			req.QueryParams,
			req.Body,
		)
		if err := a.storage.Save(savedCall); err != nil {
			return err
		}
		fmt.Printf("Saved call: %s\n", req.Save)
	}

	// Format and output response
	formatter := output.NewFormatter(a.isTTY)
	output := formatter.FormatResponse(resp, req.Info)
	fmt.Print(output)

	return nil
}

// executeRecall executes a saved call
func (a *App) executeRecall(opts *cli.RecallOptions) error {
	// Load saved call
	savedCall, err := a.storage.Load(opts.Name)
	if err != nil {
		return err
	}

	// Apply overrides
	for _, val := range opts.ParameterOverride {
		savedCall.Body = val // Simplified; should handle proper merging
	}
	for key, val := range opts.Headers {
		savedCall.Headers[key] = val
	}

	// Create a ParsedRequest from saved call
	req := &cli.ParsedRequest{
		Method:      savedCall.Method,
		URL:         savedCall.URL,
		Headers:     savedCall.Headers,
		QueryParams: savedCall.QueryParams,
		Body:        savedCall.Body,
	}

	return a.executeRequest(req)
}

// listCalls lists all saved calls
func (a *App) listCalls() error {
	calls, err := a.storage.List()
	if err != nil {
		return err
	}

	if len(calls) == 0 {
		fmt.Println("No saved calls found")
		return nil
	}

	fmt.Println("Saved calls:")
	for _, call := range calls {
		fmt.Printf("  %s (%s %s)\n", call.Name, call.Method, call.URL)
	}

	return nil
}

// deleteCall deletes a saved call
func (a *App) deleteCall(name string) error {
	if err := a.storage.Delete(name); err != nil {
		return err
	}
	fmt.Printf("Deleted: %s\n", name)
	return nil
}

// printVersion prints version info
func (a *App) printVersion() error {
	// Will be implemented when we have version package
	fmt.Println("gosh version 0.1.0")
	return nil
}

// printHelp prints help text
func (a *App) printHelp() error {
	help := `gosh - HTTPie CLI alternative built with Go

Usage:
  gosh <METHOD> <URL> [OPTIONS] [HEADERS]

Commands:
  gosh <METHOD> <URL>     Execute an HTTP request
  gosh recall <name>      Execute a saved call
  gosh list              List all saved calls
  gosh delete <name>     Delete a saved call

Options:
  -H KEY:VALUE           Add a header
  -d DATA                Request body data
  --save NAME            Save the request
  --dry                  Parse without executing
  --info                 Show full response info
  --no-interactive       Don't prompt for variables
  --env ENVIRONMENT      Use specific environment
  --format FORMAT        Output format (json|raw)

Examples:
  gosh get https://api.example.com/users
  gosh post https://api.example.com/users -d '{"name":"John"}' -H Authorization:"Bearer xyz"
  gosh get https://api.example.com/users/{userId}
  gosh recall my-request userId=42
`
	fmt.Print(help)
	return nil
}

// isTerminal checks if a file is a terminal
func isTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd())
}

// substituteEnvVars substitutes environment variables in a string
func (a *App) substituteEnvVars(text string) string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		varName := match[2 : len(match)-1] // Extract variable name from ${...}
		if val, ok := a.workspace.Env[varName]; ok {
			return val
		}
		// Return original if not found
		return match
	})
}

// substituteEnvVarsInMap substitutes environment variables in all map values
func (a *App) substituteEnvVarsInMap(m map[string]string) map[string]string {
	result := make(map[string]string)
	for key, val := range m {
		result[key] = a.substituteEnvVars(val)
	}
	return result
}

// handleAuthCommand handles authentication preset management
func (a *App) handleAuthCommand(cmd *cli.AuthCommand) error {
	switch cmd.Subcommand {
	case "list":
		presets := a.authMgr.List()
		if len(presets) == 0 {
			fmt.Println("No authentication presets configured.")
			return nil
		}
		fmt.Println("Authentication Presets:")
		for name, preset := range presets {
			fmt.Printf("  %s (%s)\n", name, preset.Type)
		}
		return nil

	case "add":
		if cmd.Type == "" || cmd.Name == "" {
			return fmt.Errorf("auth add requires type and name")
		}

		preset := &auth.AuthPreset{
			Name: cmd.Name,
			Type: cmd.Type,
		}

		// Parse auth type specific flags
		switch cmd.Type {
		case "basic":
			if username, ok := cmd.Flags["username"]; ok {
				preset.Username = username
			} else if u, ok := cmd.Flags["u"]; ok {
				preset.Username = u
			}
			if password, ok := cmd.Flags["password"]; ok {
				preset.Password = password
			} else if p, ok := cmd.Flags["p"]; ok {
				preset.Password = p
			}
			if preset.Username == "" {
				return fmt.Errorf("basic auth requires --username or -u")
			}

		case "bearer":
			if token, ok := cmd.Flags["token"]; ok {
				preset.Token = token
			} else if t, ok := cmd.Flags["t"]; ok {
				preset.Token = t
			}
			if preset.Token == "" {
				return fmt.Errorf("bearer auth requires --token or -t")
			}

		case "custom":
			if header, ok := cmd.Flags["header"]; ok {
				preset.Header = header
			} else if h, ok := cmd.Flags["h"]; ok {
				preset.Header = h
			}
			if value, ok := cmd.Flags["value"]; ok {
				preset.Value = value
			} else if v, ok := cmd.Flags["v"]; ok {
				preset.Value = v
			}
			if prefix, ok := cmd.Flags["prefix"]; ok {
				preset.Prefix = prefix
			}
			if preset.Header == "" {
				return fmt.Errorf("custom auth requires --header or -h")
			}
			if preset.Value == "" {
				return fmt.Errorf("custom auth requires --value or -v")
			}

		default:
			return fmt.Errorf("unknown auth type: %s", cmd.Type)
		}

		if err := a.authMgr.Add(preset); err != nil {
			return err
		}
		fmt.Printf("Added auth preset: %s\n", cmd.Name)
		return nil

	case "remove":
		if cmd.Name == "" {
			return fmt.Errorf("auth remove requires a preset name")
		}
		if err := a.authMgr.Remove(cmd.Name); err != nil {
			return err
		}
		fmt.Printf("Removed auth preset: %s\n", cmd.Name)
		return nil

	default:
		return fmt.Errorf("unknown auth subcommand: %s", cmd.Subcommand)
	}
}
