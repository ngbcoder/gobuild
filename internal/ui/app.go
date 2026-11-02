package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1)
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	paneStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
)

// Init is the Bubble Tea startup hook.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles UI events and window resize messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			switch m.ActiveView {
			case DashboardView:
				m.ActiveView = LogsView
			case LogsView:
				m.ActiveView = DashboardView
			default:
				m.ActiveView = DashboardView
			}
		}
	}
	return m, nil
}

// View renders the dashboard shell for the TUI.
func (m Model) View() string {
	if m.Width == 0 {
		m.Width = 80
	}
	if m.Height == 0 {
		m.Height = 24
	}

	sidebarWidth := m.Width / 3
	logsWidth := m.Width - sidebarWidth - 6
	if logsWidth < 20 {
		logsWidth = 20
	}

	sidebar := paneStyle.Width(sidebarWidth).Render("Projects\n• demo-go\n• demo-rust\n• demo-node")
	logs := paneStyle.Width(logsWidth).Render("Logs\nBuild output appears here.\nReady for the next task.")
	header := headerStyle.Render("Gobuild")
	layout := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, logs)
	body := fmt.Sprintf("%s\n\n%s", header, layout)
	return appStyle.Width(m.Width).Height(m.Height).Render(body)
}
