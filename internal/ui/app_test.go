package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModelDefaults(t *testing.T) {
	m := NewModel()
	if m.ActiveView != DashboardView {
		t.Fatalf("ActiveView = %q, want %q", m.ActiveView, DashboardView)
	}
	if m.Width != 0 || m.Height != 0 {
		t.Fatalf("initial dimensions = (%d, %d), want (0, 0)", m.Width, m.Height)
	}
}

func TestModelHandlesWindowResize(t *testing.T) {
	m := NewModel()
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model, ok := updated.(Model)
	if !ok {
		t.Fatalf("Update() returned %T, want ui.Model", updated)
	}
	if model.Width != 100 || model.Height != 30 {
		t.Fatalf("after resize dimensions = (%d, %d), want (100, 30)", model.Width, model.Height)
	}
}

func TestModelViewIncludesUISections(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	view := m.View()
	if !strings.Contains(view, "Gobuild") {
		t.Fatalf("View() = %q, expected header text", view)
	}
	if !strings.Contains(view, "Projects") {
		t.Fatalf("View() = %q, expected projects pane", view)
	}
	if !strings.Contains(view, "Logs") {
		t.Fatalf("View() = %q, expected logs pane", view)
	}
}
