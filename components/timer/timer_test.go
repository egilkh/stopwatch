package timer

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	timer := New()

	if !timer.isRunning {
		t.Error("Timer should be running when created")
	}

	if timer.elapsedTime != 0 {
		t.Error("Elapsed time should be 0 when created")
	}

	if len(timer.laps) != 0 {
		t.Error("Should have no laps when created")
	}
}

func TestToggle(t *testing.T) {
	timer := New()

	// Should start running
	if !timer.isRunning {
		t.Fatal("Timer should be running initially")
	}

	// Stop the timer
	timer.Toggle()
	if timer.isRunning {
		t.Error("Timer should be stopped after toggle")
	}

	// Start the timer again
	timer.Toggle()
	if !timer.isRunning {
		t.Error("Timer should be running after second toggle")
	}
}

func TestReset(t *testing.T) {
	timer := New()

	// Add some elapsed time
	time.Sleep(10 * time.Millisecond)

	// Add a lap
	timer.AddLap()

	// Reset
	timer.Reset()

	if timer.isRunning {
		t.Error("Timer should not be running after reset")
	}

	if timer.elapsedTime != 0 {
		t.Error("Elapsed time should be 0 after reset")
	}

	if len(timer.laps) != 0 {
		t.Error("Laps should be cleared after reset")
	}

	if timer.lapStartTime != 0 {
		t.Error("Lap start time should be 0 after reset")
	}
}

func TestAddLap(t *testing.T) {
	timer := New()

	// Wait a bit to have some elapsed time
	time.Sleep(10 * time.Millisecond)

	// Add first lap
	timer.AddLap()

	if len(timer.laps) != 1 {
		t.Fatal("Should have 1 lap")
	}

	if timer.laps[0].Number != 1 {
		t.Error("First lap should be numbered 1")
	}

	if timer.laps[0].LapTime != timer.laps[0].TotalTime {
		t.Error("First lap time should equal total time")
	}

	// Wait and add second lap
	time.Sleep(10 * time.Millisecond)
	timer.AddLap()

	if len(timer.laps) != 2 {
		t.Fatal("Should have 2 laps")
	}

	if timer.laps[1].Number != 2 {
		t.Error("Second lap should be numbered 2")
	}

	if timer.laps[1].LapTime >= timer.laps[1].TotalTime {
		t.Error("Second lap time should be less than total time")
	}
}

func TestGetCurrentTime(t *testing.T) {
	timer := New()

	// Test while running
	time.Sleep(10 * time.Millisecond)
	currentTime := timer.GetCurrentTime()

	if currentTime <= 0 {
		t.Error("Current time should be positive while running")
	}

	// Test while stopped
	timer.Toggle() // Stop
	stoppedTime := timer.GetCurrentTime()
	time.Sleep(10 * time.Millisecond)
	stoppedTime2 := timer.GetCurrentTime()

	if stoppedTime != stoppedTime2 {
		t.Error("Time should not advance while stopped")
	}
}

func TestView(t *testing.T) {
	timer := New()
	view := timer.View()

	// Should show RUNNING status
	if !strings.Contains(view, "[RUNNING]") {
		t.Error("View should show RUNNING status")
	}

	// Should show time format
	if !strings.Contains(view, "00:") {
		t.Error("View should show formatted time")
	}

	// Test stopped view
	timer.Toggle()
	view = timer.View()

	if !strings.Contains(view, "[STOPPED]") {
		t.Error("View should show STOPPED status when stopped")
	}
}

func TestUpdate(t *testing.T) {
	timer := New()

	// Test tick message
	newTimer, cmd := timer.Update(TickMsg(time.Now()))

	if cmd == nil {
		t.Error("Update should return tick command")
	}

	// Test start/stop key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	newTimer, _ = timer.Update(msg)

	if newTimer.isRunning {
		t.Error("Timer should stop on 's' key")
	}

	// Test lap key while stopped (should not add lap)
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	oldLaps := len(newTimer.laps)
	newTimer, _ = newTimer.Update(msg)

	if len(newTimer.laps) != oldLaps {
		t.Error("Should not add lap while stopped")
	}

	// Start timer and test lap
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	newTimer, _ = newTimer.Update(msg)

	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	newTimer, _ = newTimer.Update(msg)

	if len(newTimer.laps) != 1 {
		t.Error("Should add lap while running")
	}

	// Test reset key
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	newTimer, _ = newTimer.Update(msg)

	if len(newTimer.laps) != 0 || newTimer.isRunning {
		t.Error("Reset should clear laps and stop timer")
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{0, "00:00.000"},
		{time.Second, "00:01.000"},
		{time.Minute, "01:00.000"},
		{time.Hour, "01:00:00.000"},
		{time.Minute + 23*time.Second + 456*time.Millisecond, "01:23.456"},
		{time.Hour + 23*time.Minute + 45*time.Second + 678*time.Millisecond, "01:23:45.678"},
	}

	for _, test := range tests {
		result := FormatDuration(test.duration)
		if result != test.expected {
			t.Errorf("FormatDuration(%v) = %s, expected %s", test.duration, result, test.expected)
		}
	}
}

func TestGetLaps(t *testing.T) {
	timer := New()

	// Initially no laps
	if len(timer.GetLaps()) != 0 {
		t.Error("Should have no laps initially")
	}

	// Add laps
	timer.AddLap()
	timer.AddLap()

	laps := timer.GetLaps()
	if len(laps) != 2 {
		t.Error("Should return all laps")
	}
}

func TestIsRunning(t *testing.T) {
	timer := New()

	if !timer.IsRunning() {
		t.Error("Should be running initially")
	}

	timer.Toggle()

	if timer.IsRunning() {
		t.Error("Should not be running after toggle")
	}
}
