package cli

import (
	"fmt"
	"strings"
)

// Parser handles command-line argument parsing
type Parser struct {
	Args []string
}

// NewParser creates a new CLI parser
func NewParser(args []string) *Parser {
	return &Parser{Args: args}
}

// Parse parses the CLI arguments and returns the appropriate command
func (p *Parser) Parse() (interface{}, error) {
	if len(p.Args) < 1 {
		return nil, fmt.Errorf("no command provided")
	}

	cmd := strings.ToLower(p.Args[0])

	switch cmd {
	case "recall":
		return p.parseRecall()
	case "list":
		return p.parseList()
	case "delete":
		return p.parseDelete()
	case "--version", "-v":
		return "version", nil
	case "--help", "-h":
		return "help", nil
	default:
		// Assume it's an HTTP method (GET, POST, etc.)
		return p.parseRequest()
	}
}

// parseRequest parses an HTTP request command
func (p *Parser) parseRequest() (*ParsedRequest, error) {
	if len(p.Args) < 2 {
		return nil, fmt.Errorf("method and URL required")
	}

	method := strings.ToUpper(p.Args[0])
	url := p.Args[1]

	// Validate method
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "DELETE": true,
		"PATCH": true, "HEAD": true, "OPTIONS": true, "TRACE": true, "CONNECT": true,
	}
	if !validMethods[method] {
		return nil, fmt.Errorf("invalid HTTP method: %s", method)
	}

	req := &ParsedRequest{
		Method:      method,
		URL:         url,
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
		PathParams:  make(map[string]string),
	}

	// Parse remaining arguments
	for i := 2; i < len(p.Args); i++ {
		arg := p.Args[i]

		switch {
		case arg == "--save":
			if i+1 >= len(p.Args) {
				return nil, fmt.Errorf("--save requires a name argument")
			}
			i++
			req.Save = p.Args[i]
		case arg == "--dry":
			req.Dry = true
		case arg == "--info":
			req.Info = true
		case arg == "--no-interactive":
			req.NoInteractive = true
		case strings.HasPrefix(arg, "--env="):
			req.Env = strings.TrimPrefix(arg, "--env=")
		case arg == "--env":
			if i+1 >= len(p.Args) {
				return nil, fmt.Errorf("--env requires a value")
			}
			i++
			req.Env = p.Args[i]
		case strings.HasPrefix(arg, "--format="):
			req.Format = strings.TrimPrefix(arg, "--format=")
		case arg == "--format":
			if i+1 >= len(p.Args) {
				return nil, fmt.Errorf("--format requires a value")
			}
			i++
			req.Format = p.Args[i]
		case strings.HasPrefix(arg, "-H"):
			// Header: -H key:value or -H key=value
			headerVal := arg
			if arg == "-H" {
				if i+1 >= len(p.Args) {
					return nil, fmt.Errorf("-H requires a value")
				}
				i++
				headerVal = p.Args[i]
			} else {
				headerVal = strings.TrimPrefix(arg, "-H")
			}

			// Parse header key:value
			parts := strings.SplitN(headerVal, ":", 2)
			if len(parts) != 2 {
				// Try equals sign
				parts = strings.SplitN(headerVal, "=", 2)
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid header format: %s (use key:value or key=value)", headerVal)
				}
			}
			req.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		case strings.HasPrefix(arg, "-d"):
			// Body data: -d '{"json":"data"}'
			bodyVal := arg
			if arg == "-d" {
				if i+1 >= len(p.Args) {
					return nil, fmt.Errorf("-d requires a value")
				}
				i++
				bodyVal = p.Args[i]
			} else {
				bodyVal = strings.TrimPrefix(arg, "-d")
			}
			req.Body = bodyVal
		default:
			// Could be query param, path param, or unknown
			if strings.Contains(arg, "==") {
				parts := strings.SplitN(arg, "==", 2)
				req.QueryParams[parts[0]] = parts[1]
			} else if strings.Contains(arg, "=") && !strings.HasPrefix(arg, "-") {
				// Path parameter: key=value
				parts := strings.SplitN(arg, "=", 2)
				req.PathParams[parts[0]] = parts[1]
			} else {
				// Unknown parameter
				return nil, fmt.Errorf("unexpected argument: %s", arg)
			}
		}
	}

	return req, nil
}

// parseRecall parses a recall command
func (p *Parser) parseRecall() (*RecallOptions, error) {
	if len(p.Args) < 2 {
		return nil, fmt.Errorf("recall requires a call name")
	}

	opts := &RecallOptions{
		Name:              p.Args[1],
		ParameterOverride: make(map[string]string),
		Headers:           make(map[string]string),
	}

	// Parse overrides
	for i := 2; i < len(p.Args); i++ {
		arg := p.Args[i]

		switch {
		case strings.HasPrefix(arg, "-H"):
			headerVal := arg
			if arg == "-H" {
				if i+1 >= len(p.Args) {
					return nil, fmt.Errorf("-H requires a value")
				}
				i++
				headerVal = p.Args[i]
			} else {
				headerVal = strings.TrimPrefix(arg, "-H")
			}

			parts := strings.SplitN(headerVal, ":", 2)
			if len(parts) != 2 {
				parts = strings.SplitN(headerVal, "=", 2)
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid header format: %s", headerVal)
				}
			}
			opts.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		case strings.HasPrefix(arg, "--env="):
			opts.Env = strings.TrimPrefix(arg, "--env=")
		case arg == "--env":
			if i+1 >= len(p.Args) {
				return nil, fmt.Errorf("--env requires a value")
			}
			i++
			opts.Env = p.Args[i]
		default:
			// Parameter override: key=value
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg, "=", 2)
				opts.ParameterOverride[parts[0]] = parts[1]
			}
		}
	}

	return opts, nil
}

// parseList parses a list command
func (p *Parser) parseList() (string, error) {
	return "list", nil
}

// parseDelete parses a delete command
func (p *Parser) parseDelete() (string, error) {
	if len(p.Args) < 2 {
		return "", fmt.Errorf("delete requires a call name")
	}
	// Return the name to delete as a special marker
	// We'll handle this in app logic
	return fmt.Sprintf("delete:%s", p.Args[1]), nil
}
