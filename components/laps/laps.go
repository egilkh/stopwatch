package laps

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/egilkh/stopwatch/components/timer"
)

// Styles for the lap list component
var (
	lapNumberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("51"))

	lapTimeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// Model represents the lap list state
type Model struct {
	viewport viewport.Model
	laps     []timer.Lap
	ready    bool
	width    int
	height   int
}

// New creates a new lap list model
func New() Model {
	return Model{
		viewport: viewport.New(0, 0),
		laps:     []timer.Lap{},
	}
}

// Init initializes the lap list
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages for the lap list
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.viewport = viewport.New(m.width, m.height)
			m.ready = true
		} else {
			m.viewport.Width = m.width
			m.viewport.Height = m.height
		}

		m.updateContent()

	default:
		// Pass keyboard events to viewport for scrolling
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return m, cmd
}

// View renders the lap list
func (m Model) View() string {
	if !m.ready {
		return ""
	}

	return m.viewport.View()
}

// SetDimensions sets the viewport dimensions
func (m *Model) SetDimensions(width, height int) {
	m.width = width
	m.height = height

	if !m.ready {
		m.viewport = viewport.New(width, height)
		m.ready = true
	} else {
		m.viewport.Width = width
		m.viewport.Height = height
	}

	m.updateContent()
}

// SetLaps updates the lap list
func (m *Model) SetLaps(laps []timer.Lap) {
	m.laps = laps
	m.updateContent()
}

// AddLap adds a new lap and auto-scrolls to show it
func (m *Model) AddLap(lap timer.Lap) {
	m.laps = append(m.laps, lap)
	m.updateContent()
	m.viewport.GotoBottom()
}

// Clear clears all laps
func (m *Model) Clear() {
	m.laps = []timer.Lap{}
	m.updateContent()
	m.viewport.GotoTop()
}

// GetScrollInfo returns scroll information for display
func (m Model) GetScrollInfo() string {
	if m.viewport.TotalLineCount() <= m.viewport.Height {
		return ""
	}

	scrollPercent := int(float64(m.viewport.YOffset+m.viewport.Height) / float64(m.viewport.TotalLineCount()) * 100)
	if scrollPercent > 100 {
		scrollPercent = 100
	}

	return fmt.Sprintf("%d%% • Lines %d-%d of %d",
		scrollPercent,
		m.viewport.YOffset+1,
		min(m.viewport.YOffset+m.viewport.Height, m.viewport.TotalLineCount()),
		m.viewport.TotalLineCount())
}

// CanScroll returns whether the list can be scrolled
func (m Model) CanScroll() bool {
	return m.viewport.TotalLineCount() > m.viewport.Height
}

// updateContent updates the viewport content
func (m *Model) updateContent() {
	var b strings.Builder

	if len(m.laps) == 0 {
		b.WriteString(helpStyle.Render("No laps recorded yet.\nPress L to record a lap while running."))
	} else {
		for i, lap := range m.laps {
			if i > 0 {
				b.WriteString("\n")
			}

			b.WriteString(lapNumberStyle.Render(fmt.Sprintf("Lap %02d: ", lap.Number)))
			b.WriteString(lapTimeStyle.Render(timer.FormatDuration(lap.LapTime)))
			b.WriteString(" (Total: " + timer.FormatDuration(lap.TotalTime) + ")")
		}
	}

	m.viewport.SetContent(b.String())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
