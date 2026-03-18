package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"spt-scaffold/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updateConfirm(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "g", "G":
			m.state = StateProgress
			return m, tea.Batch(m.spinner.Tick, scaffoldStreamCmd(m.cfg))
		case "e", "E":
			// Go back to form step 1.
			m.formStep = 0
			m.inputs[0].Focus()
			m.state = StateForm
			return m, textinputBlink()
		}
	}
	return m, nil
}

func viewConfirm(m Model) string {
	width := m.termW
	if width < 10 {
		width = 80
	}

	title := styles.DoneTitle.Render("Review your configuration")

	// Build summary table.
	row := func(name, value string) string {
		return lipgloss.JoinHorizontal(lipgloss.Top,
			styles.FieldName.Render(name+":"),
			styles.FieldValue.Render(value),
		)
	}

	cfg := m.cfg
	sptDisplay := cfg.SptVersion + "  " + lipgloss.NewStyle().Foreground(styles.ColorMuted).Render("(range: "+cfg.SptVersionRange+")")
	summary := lipgloss.JoinVertical(lipgloss.Left,
		row("Mod Name", cfg.ModName),
		row("Author", cfg.Author),
		row("Version", cfg.Version),
		row("SPT Compatibility", sptDisplay),
		row("Description", ifEmpty(cfg.Desc, "(none)")),
		row("Repository URL", ifEmpty(cfg.RepoURL, "(none)")),
		row("License", cfg.License),
	)

	// Output path.
	cwd, _ := os.Getwd()
	outPath := filepath.Join(cwd, cfg.ModName)
	pathLine := styles.Hint.Render("Output: ") + styles.OutputPath.Render(outPath)

	boxW := width - 8
	if boxW < 40 {
		boxW = 40
	}

	box := styles.ConfirmBox.Width(boxW).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			summary,
			"",
			pathLine,
		),
	)

	actions := lipgloss.JoinHorizontal(lipgloss.Top,
		styles.ConfirmKey.Render("[G]"),
		styles.ConfirmAction.Render(" Generate    "),
		styles.ConfirmKey.Render("[E]"),
		styles.ConfirmAction.Render(" Edit"),
	)

	full := lipgloss.JoinVertical(lipgloss.Left,
		box,
		"",
		"  "+actions,
	)

	return lipgloss.NewStyle().
		Width(width).
		Padding(2, 2).
		Render(full)
}

func ifEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

func outputPath(modName string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, modName)
}

func formatOutputPath(modName string) string {
	return fmt.Sprintf("%s", outputPath(modName))
}
