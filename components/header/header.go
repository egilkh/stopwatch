package header

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles for the header component
var (
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))
)

// Model represents the header state
type Model struct {
	version      string
	programStart time.Time
}

// New creates a new header model
func New(version string) Model {
	return Model{
		version:      version,
		programStart: time.Now(),
	}
}

// Init initializes the header
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages for the header
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// Header doesn't respond to any messages
	return m, nil
}

// View renders the header
func (m Model) View() string {
	header := fmt.Sprintf("Stopwatch v%s - Started: %s", m.version, m.programStart.Format(time.RFC3339))
	return headerStyle.Render(header)
}

// GetProgramStart returns the program start time
func (m Model) GetProgramStart() time.Time {
	return m.programStart
}
