package laps

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/egilkh/stopwatch/components/timer"
)

func TestNew(t *testing.T) {
	lapList := New()

	if lapList.ready {
		t.Error("Lap list should not be ready before dimensions are set")
	}

	if len(lapList.laps) != 0 {
		t.Error("Should have no laps initially")
	}
}

func TestSetDimensions(t *testing.T) {
	lapList := New()

	lapList.SetDimensions(80, 20)

	if !lapList.ready {
		t.Error("Should be ready after setting dimensions")
	}

	if lapList.width != 80 {
		t.Error("Width should be set to 80")
	}

	if lapList.height != 20 {
		t.Error("Height should be set to 20")
	}

	if lapList.viewport.Width != 80 {
		t.Error("Viewport width should be set to 80")
	}

	if lapList.viewport.Height != 20 {
		t.Error("Viewport height should be set to 20")
	}
}

func TestAddLap(t *testing.T) {
	lapList := New()
	lapList.SetDimensions(80, 20)

	lap := timer.Lap{
		Number:    1,
		TotalTime: time.Second,
		LapTime:   time.Second,
	}

	lapList.AddLap(lap)

	if len(lapList.laps) != 1 {
		t.Error("Should have 1 lap after adding")
	}

	// Check viewport content updated
	content := lapList.viewport.View()
	if !strings.Contains(content, "Lap 01:") {
		t.Error("Viewport should show lap number")
	}
}

func TestClear(t *testing.T) {
	lapList := New()
	lapList.SetDimensions(80, 20)

	// Add some laps
	lap := timer.Lap{Number: 1, TotalTime: time.Second, LapTime: time.Second}
	lapList.AddLap(lap)
	lapList.AddLap(lap)

	// Clear
	lapList.Clear()

	if len(lapList.laps) != 0 {
		t.Error("Should have no laps after clear")
	}

	// Check viewport shows no laps message
	content := lapList.viewport.View()
	if !strings.Contains(content, "No laps recorded yet") {
		t.Error("Should show no laps message after clear")
	}
}

func TestView(t *testing.T) {
	lapList := New()

	// Test before ready
	view := lapList.View()
	if view != "" {
		t.Error("Should return empty string when not ready")
	}

	// Test after ready
	lapList.SetDimensions(80, 20)
	view = lapList.View()

	if !strings.Contains(view, "No laps recorded yet") {
		t.Error("Should show no laps message initially")
	}
}

func TestCanScroll(t *testing.T) {
	lapList := New()
	lapList.SetDimensions(80, 5) // Small height to test scrolling

	// Initially cannot scroll
	if lapList.CanScroll() {
		t.Error("Should not be able to scroll with no laps")
	}

	// Add many laps
	for i := 1; i <= 10; i++ {
		lap := timer.Lap{
			Number:    i,
			TotalTime: time.Duration(i) * time.Second,
			LapTime:   time.Second,
		}
		lapList.AddLap(lap)
	}

	// Now should be able to scroll
	if !lapList.CanScroll() {
		t.Error("Should be able to scroll with many laps")
	}
}

func TestGetScrollInfo(t *testing.T) {
	lapList := New()
	lapList.SetDimensions(80, 5)

	// No scroll info when can't scroll
	info := lapList.GetScrollInfo()
	if info != "" {
		t.Error("Should return empty string when cannot scroll")
	}

	// Add many laps
	for i := 1; i <= 10; i++ {
		lap := timer.Lap{
			Number:    i,
			TotalTime: time.Duration(i) * time.Second,
			LapTime:   time.Second,
		}
		lapList.AddLap(lap)
	}

	// Should have scroll info
	info = lapList.GetScrollInfo()
	if !strings.Contains(info, "%") {
		t.Error("Scroll info should contain percentage")
	}

	if !strings.Contains(info, "Lines") {
		t.Error("Scroll info should contain line information")
	}
}

func TestUpdate(t *testing.T) {
	lapList := New()

	// Test window size message
	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 30}
	newList, _ := lapList.Update(sizeMsg)

	if newList.width != 100 || newList.height != 30 {
		t.Error("Should update dimensions from window size message")
	}

	if !newList.ready {
		t.Error("Should be ready after window size message")
	}

	// Test key message for scrolling
	keyMsg := tea.KeyMsg{Type: tea.KeyDown}
	newList, _ = newList.Update(keyMsg)
	// Just verify it doesn't panic, actual scrolling tested in viewport
}

func TestUpdateContent(t *testing.T) {
	lapList := New()
	lapList.SetDimensions(80, 20)

	// Test with no laps
	lapList.updateContent()
	content := lapList.viewport.View()

	if !strings.Contains(content, "No laps recorded yet") {
		t.Error("Should show no laps message")
	}

	// Test with laps
	lap1 := timer.Lap{
		Number:    1,
		TotalTime: 5 * time.Second,
		LapTime:   5 * time.Second,
	}
	lap2 := timer.Lap{
		Number:    2,
		TotalTime: 8 * time.Second,
		LapTime:   3 * time.Second,
	}

	lapList.laps = []timer.Lap{lap1, lap2}
	lapList.updateContent()
	content = lapList.viewport.View()

	if !strings.Contains(content, "Lap 01:") {
		t.Error("Should show first lap")
	}

	if !strings.Contains(content, "Lap 02:") {
		t.Error("Should show second lap")
	}

	if !strings.Contains(content, "00:05.000") {
		t.Error("Should show formatted lap time")
	}

	if !strings.Contains(content, "Total: 00:08.000") {
		t.Error("Should show total time")
	}
}

func TestSetLaps(t *testing.T) {
	lapList := New()
	lapList.SetDimensions(80, 20)

	laps := []timer.Lap{
		{Number: 1, TotalTime: time.Second, LapTime: time.Second},
		{Number: 2, TotalTime: 2 * time.Second, LapTime: time.Second},
	}

	lapList.SetLaps(laps)

	if len(lapList.laps) != 2 {
		t.Error("Should have 2 laps after SetLaps")
	}

	content := lapList.viewport.View()
	if !strings.Contains(content, "Lap 01:") || !strings.Contains(content, "Lap 02:") {
		t.Error("Should show both laps in content")
	}
}
