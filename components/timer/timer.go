package timer

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles for the timer component
var (
	timerStyle = lipgloss.NewStyle().
			Bold(true)

	runningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	stoppedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
)

// Messages
type TickMsg time.Time

// Lap represents a single lap record
type Lap struct {
	Number    int
	TotalTime time.Duration
	LapTime   time.Duration
}

// Model represents the timer state
type Model struct {
	startTime    time.Time
	elapsedTime  time.Duration
	isRunning    bool
	lapStartTime time.Duration
	laps         []Lap
}

// New creates a new timer model
func New() Model {
	return Model{
		startTime: time.Now(),
		isRunning: true,
	}
}

// Init initializes the timer
func (m Model) Init() tea.Cmd {
	return tickCmd()
}

// Update handles messages for the timer
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		// Timer updates are handled by the parent
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "s", "S":
			m.Toggle()
		case "l", "L":
			if m.isRunning {
				m.AddLap()
			}
		case "r", "R":
			m.Reset()
		}
	}

	return m, nil
}

// View renders the timer
func (m Model) View() string {
	currentTime := m.GetCurrentTime()
	status := stoppedStyle.Render("[STOPPED]")
	if m.isRunning {
		status = runningStyle.Render("[RUNNING]")
	}
	return timerStyle.Render(FormatDuration(currentTime) + " " + status)
}

// Toggle starts or stops the timer
func (m *Model) Toggle() {
	if m.isRunning {
		m.elapsedTime = time.Since(m.startTime)
		m.isRunning = false
	} else {
		m.startTime = time.Now().Add(-m.elapsedTime)
		m.isRunning = true
		m.lapStartTime = m.elapsedTime
	}
}

// Reset resets the timer
func (m *Model) Reset() {
	m.elapsedTime = 0
	m.laps = []Lap{}
	m.isRunning = false
	m.lapStartTime = 0
}

// AddLap adds a new lap
func (m *Model) AddLap() {
	currentTime := m.GetCurrentTime()
	lapTime := currentTime - m.lapStartTime

	lapNumber := len(m.laps) + 1
	m.laps = append(m.laps, Lap{
		Number:    lapNumber,
		TotalTime: currentTime,
		LapTime:   lapTime,
	})

	m.lapStartTime = currentTime
}

// GetCurrentTime returns the current elapsed time
func (m Model) GetCurrentTime() time.Duration {
	if m.isRunning {
		return time.Since(m.startTime)
	}
	return m.elapsedTime
}

// GetLaps returns the laps
func (m Model) GetLaps() []Lap {
	return m.laps
}

// IsRunning returns whether the timer is running
func (m Model) IsRunning() bool {
	return m.isRunning
}

// tickCmd returns a command that sends a tick message
func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*25, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// FormatDuration formats a duration for display
func FormatDuration(d time.Duration) string {
	d = d.Round(time.Millisecond)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ms := d / time.Millisecond

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
	}
	return fmt.Sprintf("%02d:%02d.%03d", m, s, ms)
}
