package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/egilkh/stopwatch/components/header"
	"github.com/egilkh/stopwatch/components/laps"
	"github.com/egilkh/stopwatch/components/statusbar"
	"github.com/egilkh/stopwatch/components/timer"
)

const version = "1.0.0"

// Style for scroll info
var helpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241"))

// ExportData represents the JSON export structure
type ExportData struct {
	ExportTime  string      `json:"export_time"`
	StartTime   string      `json:"start_time"`
	TotalTime   string      `json:"total_time"`
	TotalTimeMs int64       `json:"total_time_ms"`
	LapCount    int         `json:"lap_count"`
	Laps        []ExportLap `json:"laps"`
}

// ExportLap represents a lap in the export
type ExportLap struct {
	Number      int    `json:"number"`
	LapTime     string `json:"lap_time"`
	LapTimeMs   int64  `json:"lap_time_ms"`
	TotalTime   string `json:"total_time"`
	TotalTimeMs int64  `json:"total_time_ms"`
}

// Model represents the application state
type Model struct {
	// Components
	header    header.Model
	timer     timer.Model
	lapList   laps.Model
	statusBar statusbar.Model

	// Application state
	width  int
	height int
	ready  bool
}

func initialModel() Model {
	return Model{
		header:    header.New(version),
		timer:     timer.New(),
		lapList:   laps.New(),
		statusBar: statusbar.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.header.Init(),
		m.timer.Init(),
		m.lapList.Init(),
		m.statusBar.Init(),
		tea.SetWindowTitle("Stopwatch"),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		// Calculate viewport size for lap list
		headerHeight := 5 // header + help + timer + spacing + blank line
		if m.statusBar.HasMessage() {
			headerHeight++
		}
		footerHeight := 1 // bottom padding

		viewportHeight := m.height - headerHeight - footerHeight
		if viewportHeight < 1 {
			viewportHeight = 1
		}

		// Update lap list dimensions
		m.lapList.SetDimensions(m.width, viewportHeight)

	case timer.TickMsg:
		// Update timer
		newTimer, cmd := m.timer.Update(msg)
		m.timer = newTimer
		cmds = append(cmds, cmd)

		// Update status bar timers
		m.statusBar.Tick()

		// Recalculate viewport if message state changed
		if m.ready {
			headerHeight := 5
			if m.statusBar.HasMessage() {
				headerHeight++
			}
			viewportHeight := m.height - headerHeight - 1
			if viewportHeight < 1 {
				viewportHeight = 1
			}
			m.lapList.SetDimensions(m.width, viewportHeight)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "s", "S":
			// Update timer
			newTimer, cmd := m.timer.Update(msg)
			m.timer = newTimer
			cmds = append(cmds, cmd)

		case "l", "L":
			if m.timer.IsRunning() {
				// Get current state before update
				oldLaps := m.timer.GetLaps()

				// Update timer (adds lap)
				newTimer, cmd := m.timer.Update(msg)
				m.timer = newTimer
				cmds = append(cmds, cmd)

				// Check if a new lap was added
				newLaps := m.timer.GetLaps()
				if len(newLaps) > len(oldLaps) {
					// Add the new lap to the lap list
					m.lapList.AddLap(newLaps[len(newLaps)-1])
				}
			}

		case "r", "R":
			// Update timer
			newTimer, cmd := m.timer.Update(msg)
			m.timer = newTimer
			cmds = append(cmds, cmd)

			// Clear lap list
			m.lapList.Clear()

			// Show reset message
			m.statusBar.ShowResetMessage()

			// Recalculate viewport
			headerHeight := 6 // includes reset message
			viewportHeight := m.height - headerHeight - 1
			if viewportHeight < 1 {
				viewportHeight = 1
			}
			m.lapList.SetDimensions(m.width, viewportHeight)

		case "f", "F":
			laps := m.timer.GetLaps()
			if len(laps) > 0 {
				filename := m.exportToJSON()
				m.statusBar.ShowExportMessage(filename)

				// Recalculate viewport
				headerHeight := 6 // includes export message
				viewportHeight := m.height - headerHeight - 1
				if viewportHeight < 1 {
					viewportHeight = 1
				}
				m.lapList.SetDimensions(m.width, viewportHeight)
			}

		default:
			// Pass keyboard events to lap list for scrolling
			newLapList, cmd := m.lapList.Update(msg)
			m.lapList = newLapList
			cmds = append(cmds, cmd)
		}

		// Update status bar scroll indicator
		m.statusBar.SetCanScroll(m.lapList.CanScroll())
	}

	// Update all components
	newHeader, cmd := m.header.Update(msg)
	m.header = newHeader
	cmds = append(cmds, cmd)

	newStatusBar, cmd := m.statusBar.Update(msg)
	m.statusBar = newStatusBar
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var b strings.Builder

	// Header
	b.WriteString(m.header.View() + "\n")

	// Status bar help
	b.WriteString(m.statusBar.ViewHelp() + "\n")

	// Timer
	b.WriteString(m.timer.View() + "\n")

	// Status messages
	if msg := m.statusBar.ViewMessages(); msg != "" {
		b.WriteString(msg + "\n")
	}

	// Blank line before laps
	b.WriteString("\n")

	// Lap list
	b.WriteString(m.lapList.View())

	// Scroll indicator
	if scrollInfo := m.lapList.GetScrollInfo(); scrollInfo != "" {
		b.WriteString("\n" + helpStyle.Render(scrollInfo))
	}

	return b.String()
}

func (m *Model) exportToJSON() string {
	// Generate filename with timestamp
	now := time.Now()
	filename := fmt.Sprintf("stopwatch_%s.json", now.Format("20060102_150405"))

	// Get data from components
	currentTime := m.timer.GetCurrentTime()
	laps := m.timer.GetLaps()
	programStart := m.header.GetProgramStart()

	// Create export data
	exportData := ExportData{
		ExportTime:  now.Format(time.RFC3339),
		StartTime:   programStart.Format(time.RFC3339),
		TotalTime:   timer.FormatDuration(currentTime),
		TotalTimeMs: currentTime.Milliseconds(),
		LapCount:    len(laps),
		Laps:        make([]ExportLap, len(laps)),
	}

	// Convert laps
	for i, lap := range laps {
		exportData.Laps[i] = ExportLap{
			Number:      lap.Number,
			LapTime:     timer.FormatDuration(lap.LapTime),
			LapTimeMs:   lap.LapTime.Milliseconds(),
			TotalTime:   timer.FormatDuration(lap.TotalTime),
			TotalTimeMs: lap.TotalTime.Milliseconds(),
		}
	}

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Sprintf("error-%s.json", now.Format("20060102_150405"))
	}

	// Write to file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Sprintf("error-%s.json", now.Format("20060102_150405"))
	}

	return filename
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
