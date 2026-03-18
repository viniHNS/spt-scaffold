package tui

import (
	"spt-scaffold/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updateError(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q", "enter", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func viewError(m Model) string {
	width := m.termW
	if width < 10 {
		width = 80
	}

	title := styles.ValidationError.Render("✗  Scaffold failed")

	errMsg := "(unknown error)"
	if m.fatalErr != nil {
		errMsg = m.fatalErr.Error()
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorError).
		Padding(1, 3).
		Width(width - 8).
		Render(errMsg)

	hint := styles.QuitHint.Render("  Press Q or ENTER to quit")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		box,
		"",
		hint,
	)

	return lipgloss.NewStyle().
		Width(width).
		Padding(2, 2).
		Render(content)
}
