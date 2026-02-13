package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// PromptModel is a simple text input model for bubbletea
type PromptModel struct {
	varName string
	input   string
	done    bool
}

// Init initializes the model
func (m *PromptModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *PromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.done = true
			return m, tea.Quit
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			m.input += msg.String()
		}
	}
	return m, nil
}

// View renders the model
func (m *PromptModel) View() string {
	return fmt.Sprintf("Enter value for {%s}: %s", m.varName, m.input)
}

// PromptForVariable prompts for a missing template variable
func PromptForVariable(varName string) (string, error) {
	// For now, use simple stdin prompt
	// Can be enhanced with bubbletea later
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter value for {%s}: ", varName)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// PromptInteractively prompts for multiple variables using bubbletea
func PromptInteractively(variables []string) (map[string]string, error) {
	result := make(map[string]string)

	for _, varName := range variables {
		val, err := PromptForVariable(varName)
		if err != nil {
			return nil, err
		}
		result[varName] = val
	}

	return result, nil
}
