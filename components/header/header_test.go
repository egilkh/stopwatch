package header

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	version := "1.0.0"
	before := time.Now()
	header := New(version)
	after := time.Now()

	if header.version != version {
		t.Errorf("Version should be %s, got %s", version, header.version)
	}

	if header.programStart.Before(before) || header.programStart.After(after) {
		t.Error("Program start time should be set to current time")
	}
}

func TestInit(t *testing.T) {
	header := New("1.0.0")
	cmd := header.Init()

	if cmd != nil {
		t.Error("Init should return nil command")
	}
}

func TestUpdate(t *testing.T) {
	header := New("1.0.0")

	// Test with various messages - header doesn't respond to any
	messages := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
		tea.WindowSizeMsg{Width: 80, Height: 24},
		time.Now(), // Random message type
	}

	for _, msg := range messages {
		newHeader, cmd := header.Update(msg)

		if cmd != nil {
			t.Error("Update should always return nil command")
		}

		if newHeader.version != header.version {
			t.Error("Update should not change header state")
		}

		if !newHeader.programStart.Equal(header.programStart) {
			t.Error("Update should not change program start time")
		}
	}
}

func TestView(t *testing.T) {
	version := "1.2.3"
	header := New(version)
	view := header.View()

	// Check version is in view
	if !strings.Contains(view, "v1.2.3") {
		t.Error("View should contain version with 'v' prefix")
	}

	// Check "Stopwatch" is in view
	if !strings.Contains(view, "Stopwatch") {
		t.Error("View should contain 'Stopwatch'")
	}

	// Check "Started:" is in view
	if !strings.Contains(view, "Started:") {
		t.Error("View should contain 'Started:'")
	}

	// Check date format (RFC3339) is used
	if !strings.Contains(view, "T") || !strings.Contains(view, ":") {
		t.Error("View should contain RFC3339 formatted time")
	}
}

func TestGetProgramStart(t *testing.T) {
	header := New("1.0.0")
	startTime := header.GetProgramStart()

	if !startTime.Equal(header.programStart) {
		t.Error("GetProgramStart should return the program start time")
	}

	// Verify it doesn't change
	time.Sleep(10 * time.Millisecond)
	startTime2 := header.GetProgramStart()

	if !startTime2.Equal(startTime) {
		t.Error("Program start time should not change")
	}
}

func TestViewFormatting(t *testing.T) {
	// Test with specific time for consistent formatting
	header := New("2.0.0")
	header.programStart = time.Date(2024, 12, 25, 15, 30, 45, 0, time.UTC)

	view := header.View()

	// Should contain the formatted date
	if !strings.Contains(view, "2024-12-25") {
		t.Error("View should contain formatted date")
	}

	// Should contain the formatted time
	if !strings.Contains(view, "15:30:45") {
		t.Error("View should contain formatted time")
	}
}

func TestHeaderImmutability(t *testing.T) {
	// Test that header state doesn't change after creation
	version := "1.0.0"
	header := New(version)
	originalStart := header.programStart
	originalVersion := header.version

	// Try various operations
	header.Init()
	header.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	header.View()
	header.GetProgramStart()

	// Verify nothing changed
	if header.version != originalVersion {
		t.Error("Version should not change")
	}

	if !header.programStart.Equal(originalStart) {
		t.Error("Program start time should not change")
	}
}
