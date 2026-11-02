package ui

// ViewMode tracks which screen is active within the TUI.
type ViewMode string

const (
	DashboardView ViewMode = "dashboard"
	LogsView      ViewMode = "logs"
	PluginView    ViewMode = "plugins"
	ConfigView    ViewMode = "config"
)

// Model stores the application state for the Bubble Tea UI.
type Model struct {
	ActiveView ViewMode
	Width      int
	Height     int
}

// NewModel builds the initial application state for the TUI.
func NewModel() Model {
	return Model{ActiveView: DashboardView}
}
