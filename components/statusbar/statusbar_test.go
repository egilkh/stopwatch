package statusbar

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	statusBar := New()

	if statusBar.showResetMsg {
		t.Error("Should not show reset message initially")
	}

	if statusBar.showExportMsg {
		t.Error("Should not show export message initially")
	}

	if statusBar.canScroll {
		t.Error("Should not show scroll help initially")
	}
}

func TestInit(t *testing.T) {
	statusBar := New()
	cmd := statusBar.Init()

	if cmd != nil {
		t.Error("Init should return nil command")
	}
}

func TestUpdate(t *testing.T) {
	statusBar := New()

	// Test with various messages - statusbar doesn't actively respond to any
	messages := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
		tea.WindowSizeMsg{Width: 80, Height: 24},
		"random message",
	}

	for _, msg := range messages {
		newStatusBar, cmd := statusBar.Update(msg)

		if cmd != nil {
			t.Error("Update should always return nil command")
		}

		// Verify state hasn't changed
		if newStatusBar.showResetMsg != statusBar.showResetMsg {
			t.Error("Update should not change reset message state")
		}
	}
}

func TestViewHelp(t *testing.T) {
	statusBar := New()

	// Test basic help
	help := statusBar.ViewHelp()

	shortcuts := []string{"[S]tart/Stop", "[L]ap", "[R]eset", "[F]ile/Export", "[Q]uit"}
	for _, shortcut := range shortcuts {
		if !strings.Contains(help, shortcut) {
			t.Errorf("Help should contain %s", shortcut)
		}
	}

	// Should not show scroll help initially
	if strings.Contains(help, "Scroll") {
		t.Error("Should not show scroll help when canScroll is false")
	}

	// Test with scroll enabled
	statusBar.SetCanScroll(true)
	help = statusBar.ViewHelp()

	if !strings.Contains(help, "[↑/↓] Scroll") {
		t.Error("Should show scroll help when canScroll is true")
	}
}

func TestViewMessages(t *testing.T) {
	statusBar := New()

	// Initially no messages
	msg := statusBar.ViewMessages()
	if msg != "" {
		t.Error("Should return empty string when no messages")
	}

	// Test reset message
	statusBar.ShowResetMessage()
	msg = statusBar.ViewMessages()

	if !strings.Contains(msg, "Timer Reset") {
		t.Error("Should show reset message")
	}

	// Test export message
	statusBar = New() // Reset
	filename := "stopwatch_20240101_120000.json"
	statusBar.ShowExportMessage(filename)
	msg = statusBar.ViewMessages()

	if !strings.Contains(msg, "✓ Exported to") {
		t.Error("Should show export message with checkmark")
	}

	if !strings.Contains(msg, filename) {
		t.Error("Should show filename in export message")
	}
}

func TestShowResetMessage(t *testing.T) {
	statusBar := New()

	statusBar.ShowResetMessage()

	if !statusBar.showResetMsg {
		t.Error("Should set showResetMsg to true")
	}

	if statusBar.resetMsgTimer != 80 {
		t.Error("Should set reset timer to 80")
	}
}

func TestShowExportMessage(t *testing.T) {
	statusBar := New()
	filename := "test.json"

	statusBar.ShowExportMessage(filename)

	if !statusBar.showExportMsg {
		t.Error("Should set showExportMsg to true")
	}

	if statusBar.exportMsgTimer != 120 {
		t.Error("Should set export timer to 120")
	}

	if statusBar.exportFilename != filename {
		t.Error("Should store export filename")
	}
}

func TestSetCanScroll(t *testing.T) {
	statusBar := New()

	statusBar.SetCanScroll(true)
	if !statusBar.canScroll {
		t.Error("Should set canScroll to true")
	}

	statusBar.SetCanScroll(false)
	if statusBar.canScroll {
		t.Error("Should set canScroll to false")
	}
}

func TestTick(t *testing.T) {
	statusBar := New()

	// Test reset message timer
	statusBar.ShowResetMessage()
	initialTimer := statusBar.resetMsgTimer

	statusBar.Tick()

	if statusBar.resetMsgTimer != initialTimer-1 {
		t.Error("Should decrement reset timer")
	}

	// Test timer expiration
	for i := 0; i < 80; i++ {
		statusBar.Tick()
	}

	if statusBar.showResetMsg {
		t.Error("Should hide reset message when timer expires")
	}

	// Test export message timer
	statusBar.ShowExportMessage("test.json")
	initialTimer = statusBar.exportMsgTimer

	statusBar.Tick()

	if statusBar.exportMsgTimer != initialTimer-1 {
		t.Error("Should decrement export timer")
	}

	// Test both timers don't interfere
	statusBar.ShowResetMessage()
	statusBar.ShowExportMessage("test.json")
	resetTimer := statusBar.resetMsgTimer
	exportTimer := statusBar.exportMsgTimer

	statusBar.Tick()

	if statusBar.resetMsgTimer != resetTimer-1 {
		t.Error("Should decrement reset timer independently")
	}

	if statusBar.exportMsgTimer != exportTimer-1 {
		t.Error("Should decrement export timer independently")
	}
}

func TestHasMessage(t *testing.T) {
	statusBar := New()

	// Initially no message
	if statusBar.HasMessage() {
		t.Error("Should not have message initially")
	}

	// Test with reset message
	statusBar.ShowResetMessage()
	if !statusBar.HasMessage() {
		t.Error("Should have message when reset message shown")
	}

	// Test with export message
	statusBar = New()
	statusBar.ShowExportMessage("test.json")
	if !statusBar.HasMessage() {
		t.Error("Should have message when export message shown")
	}

	// Test with both messages
	statusBar.ShowResetMessage()
	if !statusBar.HasMessage() {
		t.Error("Should have message when both messages shown")
	}

	// Test after expiration
	statusBar = New()
	statusBar.ShowResetMessage()
	for i := 0; i < 81; i++ {
		statusBar.Tick()
	}

	if statusBar.HasMessage() {
		t.Error("Should not have message after expiration")
	}
}

func TestMessagePriority(t *testing.T) {
	statusBar := New()

	// Show both messages
	statusBar.ShowResetMessage()
	statusBar.ShowExportMessage("test.json")

	// Reset message should take priority in ViewMessages
	msg := statusBar.ViewMessages()
	if !strings.Contains(msg, "Timer Reset") {
		t.Error("Reset message should take priority")
	}

	// Let reset message expire
	for i := 0; i < 81; i++ {
		statusBar.Tick()
	}

	// Now export message should show
	msg = statusBar.ViewMessages()
	if !strings.Contains(msg, "Exported") {
		t.Error("Export message should show after reset expires")
	}
}
