package statusbar

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles for the status bar component
var (
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	resetMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("226"))

	exportMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42"))
)

// Model represents the status bar state
type Model struct {
	showResetMsg   bool
	resetMsgTimer  int
	showExportMsg  bool
	exportMsgTimer int
	exportFilename string
	canScroll      bool
}

// New creates a new status bar model
func New() Model {
	return Model{}
}

// Init initializes the status bar
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages for the status bar
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// Status bar responds to tick messages to update timers
	switch msg.(type) {
	case tea.KeyMsg:
		// Reset any temporary messages on key press
		// (handled by parent)
	}

	return m, nil
}

// View renders the status bar (help text only)
func (m Model) ViewHelp() string {
	help := "[S]tart/Stop [L]ap [R]eset [F]ile/Export [Q]uit"
	if m.canScroll {
		help += " [↑/↓] Scroll"
	}
	return helpStyle.Render(help)
}

// ViewMessages renders any temporary messages
func (m Model) ViewMessages() string {
	if m.showResetMsg {
		return resetMessageStyle.Render("--- Timer Reset ---")
	}
	if m.showExportMsg {
		return exportMessageStyle.Render(fmt.Sprintf("✓ Exported to %s", m.exportFilename))
	}
	return ""
}

// ShowResetMessage shows the reset message
func (m *Model) ShowResetMessage() {
	m.showResetMsg = true
	m.resetMsgTimer = 80 // 2 seconds at 25ms tick
}

// ShowExportMessage shows the export message
func (m *Model) ShowExportMessage(filename string) {
	m.exportFilename = filename
	m.showExportMsg = true
	m.exportMsgTimer = 120 // 3 seconds at 25ms tick
}

// SetCanScroll updates whether scroll help should be shown
func (m *Model) SetCanScroll(canScroll bool) {
	m.canScroll = canScroll
}

// Tick decrements message timers
func (m *Model) Tick() {
	if m.showResetMsg {
		m.resetMsgTimer--
		if m.resetMsgTimer <= 0 {
			m.showResetMsg = false
		}
	}
	if m.showExportMsg {
		m.exportMsgTimer--
		if m.exportMsgTimer <= 0 {
			m.showExportMsg = false
		}
	}
}

// HasMessage returns whether any message is being displayed
func (m Model) HasMessage() bool {
	return m.showResetMsg || m.showExportMsg
}
