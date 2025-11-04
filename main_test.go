package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/egilkh/stopwatch/components/timer"
)

func TestInitialModel(t *testing.T) {
	model := initialModel()

	// Verify components are initialized
	if model.header.GetProgramStart().IsZero() {
		t.Error("Header should have program start time")
	}

	if !model.timer.IsRunning() {
		t.Error("Timer should be running initially")
	}

	if model.ready {
		t.Error("Model should not be ready initially")
	}
}

func TestInit(t *testing.T) {
	model := initialModel()
	cmd := model.Init()

	if cmd == nil {
		t.Error("Init should return a batch command")
	}
}

func TestWindowSizeHandling(t *testing.T) {
	model := initialModel()

	// Send window size message
	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 40}
	newModelInterface, _ := model.Update(sizeMsg)
	m, ok := newModelInterface.(Model)
	if !ok {
		t.Fatal("Failed to assert Model type")
	}

	if !m.ready {
		t.Error("Model should be ready after window size")
	}

	if m.width != 100 || m.height != 40 {
		t.Error("Model should store window dimensions")
	}
}

func TestKeyHandling(t *testing.T) {
	model := initialModel()
	// Set window size first
	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 40}
	modelInterface, _ := model.Update(sizeMsg)
	model = modelInterface.(Model)

	// Test start/stop
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	newModel, _ := model.Update(keyMsg)
	m := newModel.(Model)

	if m.timer.IsRunning() {
		t.Error("Timer should stop on 's' key")
	}

	// Test lap (while stopped, should not add)
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	newModel, _ = m.Update(keyMsg)
	m = newModel.(Model)

	if len(m.timer.GetLaps()) != 0 {
		t.Error("Should not add lap while stopped")
	}

	// Start again and test lap
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	newModel, _ = m.Update(keyMsg)
	m = newModel.(Model)

	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	newModel, _ = m.Update(keyMsg)
	m = newModel.(Model)

	if len(m.timer.GetLaps()) != 1 {
		t.Error("Should add lap while running")
	}

	// Test reset
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	newModel, _ = m.Update(keyMsg)
	m = newModel.(Model)

	if len(m.timer.GetLaps()) != 0 {
		t.Error("Reset should clear laps")
	}

	if !m.statusBar.HasMessage() {
		t.Error("Reset should show message")
	}
}

func TestQuitHandling(t *testing.T) {
	model := initialModel()

	// Test 'q' key
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(keyMsg)

	if cmd == nil {
		t.Error("Should return quit command on 'q'")
	}

	// Test Ctrl+C
	keyMsg = tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd = model.Update(keyMsg)

	if cmd == nil {
		t.Error("Should return quit command on Ctrl+C")
	}
}

func TestTickHandling(t *testing.T) {
	model := initialModel()
	model.statusBar.ShowResetMessage()

	// Send tick messages
	tickMsg := timer.TickMsg(time.Now())

	for i := 0; i < 100; i++ {
		newModel, cmd := model.Update(tickMsg)
		model = newModel.(Model)

		if cmd == nil {
			t.Error("Tick should return a command")
		}
	}

	// After enough ticks, message should disappear
	if model.statusBar.HasMessage() {
		t.Error("Message should disappear after enough ticks")
	}
}

func TestView(t *testing.T) {
	model := initialModel()

	// Test view before ready
	view := model.View()
	if view != "Initializing..." {
		t.Error("Should show initializing message when not ready")
	}

	// Make ready
	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 40}
	newModelInterface, _ := model.Update(sizeMsg)
	model = newModelInterface.(Model)

	view = model.View()

	// Check all components are in view
	if !strings.Contains(view, "Stopwatch v"+version) {
		t.Error("View should contain header")
	}

	if !strings.Contains(view, "[S]tart/Stop") {
		t.Error("View should contain help text")
	}

	if !strings.Contains(view, "[RUNNING]") {
		t.Error("View should contain timer status")
	}

	if !strings.Contains(view, "No laps recorded yet") {
		t.Error("View should contain lap list placeholder")
	}
}

func TestExportFunctionality(t *testing.T) {
	model := initialModel()
	// Setup
	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 40}
	modelInterface, _ := model.Update(sizeMsg)
	model = modelInterface.(Model)

	// Add some laps properly through key events
	// First ensure timer is running
	if !model.timer.IsRunning() {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
		modelInterface, _ = model.Update(keyMsg)
		model = modelInterface.(Model)
	}

	// Add first lap
	time.Sleep(10 * time.Millisecond)
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	modelInterface, _ = model.Update(keyMsg)
	model = modelInterface.(Model)

	// Add second lap
	time.Sleep(10 * time.Millisecond)
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	modelInterface, _ = model.Update(keyMsg)
	model = modelInterface.(Model)

	// Test export with no laps shows no message
	modelNoLaps := initialModel()
	modelNoLapsInterface, _ := modelNoLaps.Update(sizeMsg)
	modelNoLaps = modelNoLapsInterface.(Model)
	keyMsgF := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
	newModel, _ := modelNoLaps.Update(keyMsgF)
	m := newModel.(Model)

	if m.statusBar.HasMessage() {
		t.Error("Should not show export message with no laps")
	}

	// Test export with laps
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
	newModel, _ = model.Update(keyMsg)
	m = newModel.(Model)

	if !m.statusBar.HasMessage() {
		t.Error("Should show export message after export")
	}

	// Check file was created (cleanup after)
	files, _ := os.ReadDir(".")
	found := false
	var exportFile string

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "stopwatch_") && strings.HasSuffix(file.Name(), ".json") {
			found = true
			exportFile = file.Name()
			break
		}
	}

	if !found {
		t.Error("Export file should be created")
	} else {
		// Verify file content
		data, err := os.ReadFile(exportFile)
		if err != nil {
			t.Errorf("Could not read export file: %v", err)
		}

		var export ExportData
		err = json.Unmarshal(data, &export)
		if err != nil {
			t.Errorf("Export file should be valid JSON: %v", err)
		}

		if export.LapCount != 2 {
			t.Error("Export should have 2 laps")
		}

		// Cleanup
		os.Remove(exportFile)
	}
}

func TestExportToJSON(t *testing.T) {
	model := initialModel()
	model.timer.AddLap()
	model.timer.AddLap()

	filename := model.exportToJSON()

	// Check filename format
	if !strings.HasPrefix(filename, "stopwatch_") {
		t.Error("Filename should start with 'stopwatch_'")
	}

	if !strings.HasSuffix(filename, ".json") {
		t.Error("Filename should end with '.json'")
	}

	// Verify file exists and is valid
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Could not read exported file: %v", err)
	}

	var export ExportData
	err = json.Unmarshal(data, &export)
	if err != nil {
		t.Errorf("Exported file should be valid JSON: %v", err)
	}

	// Verify content
	if export.LapCount != 2 {
		t.Error("Should export 2 laps")
	}

	if len(export.Laps) != 2 {
		t.Error("Should have 2 lap entries")
	}

	if export.Laps[0].Number != 1 {
		t.Error("First lap should be numbered 1")
	}

	// Cleanup
	os.Remove(filename)
}

func TestScrollingFunctionality(t *testing.T) {
	model := initialModel()
	// Set small height to test scrolling
	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 10}
	modelInterface, _ := model.Update(sizeMsg)
	model = modelInterface.(Model)

	// Add many laps
	for i := 0; i < 10; i++ {
		model.timer.AddLap()
		model.lapList.AddLap(model.timer.GetLaps()[i])
	}

	// Test scroll keys
	keyMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := model.Update(keyMsg)
	// Just verify it doesn't panic

	keyMsg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = newModel.Update(keyMsg)
	// Just verify it doesn't panic

	// Check scroll info in view
	m := newModel.(Model)
	view := m.View()

	if m.lapList.CanScroll() && !strings.Contains(view, "%") {
		t.Error("View should show scroll percentage when scrollable")
	}
}
