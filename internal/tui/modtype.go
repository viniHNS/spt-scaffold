package tui

import (
	"fmt"

	"spt-scaffold/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updateModType(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = StateWelcome
			m.modTypeHint = ""
			return m, nil
		case "up":
			if m.modTypeIdx > 0 {
				m.modTypeIdx--
			}
			return m, nil
		case "down":
			if m.modTypeIdx < 1 {
				m.modTypeIdx++
			}
			return m, nil
		case "enter":
			m.modTypeHint = ""
			m.formStep = stepTemplate
			m.templateIdx = 0
			m.state = StateForm
			return m, nil
		}
	}
	return m, nil
}

func viewModType(m Model) string {
	width := m.termW
	if width < 10 {
		width = 80
	}

	stepStr := styles.StepIndicator.Render("Step 1 of 2")
	header := lipgloss.NewStyle().Width(width - 4).Align(lipgloss.Right).Render(stepStr)

	label := styles.Label.Render("Mod Type")

	var serverRow, clientRow string
	if m.modTypeIdx == 0 {
		serverRow = styles.SelectItemActive.Render("▶ Server Mod")
	} else {
		serverRow = styles.SelectItem.Render("  Server Mod")
	}
	if m.modTypeIdx == 1 {
		clientRow = styles.SelectItemActive.Render("▶ Client Mod")
	} else {
		clientRow = styles.SelectItem.Render("  Client Mod")
	}

	list := lipgloss.JoinVertical(lipgloss.Left, serverRow, clientRow)

	hintLines := styles.Hint.Render("↑/↓ to select, ENTER to confirm, ESC to go back")
	if m.modTypeHint != "" {
		hintLines = fmt.Sprintf("%s\n%s",
			styles.Hint.Render(m.modTypeHint),
			hintLines,
		)
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		label,
		"",
		list,
		"",
		hintLines,
	)

	boxW := width - 8
	if boxW < 40 {
		boxW = 40
	}

	box := styles.FormBox.Width(boxW).Render(body)

	full := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		box,
	)

	return lipgloss.NewStyle().
		Width(width).
		Padding(2, 2).
		Render(full)
}
