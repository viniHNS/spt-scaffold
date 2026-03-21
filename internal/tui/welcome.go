package tui

import (
	"fmt"
	"strings"
	"time"

	"spt-scaffold/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const asciiArt = `
 ███████╗██████╗ ████████╗      ███████╗ ██████╗ █████╗ ███████╗███████╗ ██████╗ ██╗     ██████╗
 ██╔════╝██╔══██╗╚══██╔══╝      ██╔════╝██╔════╝██╔══██╗██╔════╝██╔════╝██╔═══██╗██║     ██╔══██╗
 ███████╗██████╔╝   ██║   █████╗███████╗██║     ███████║█████╗  █████╗  ██║   ██║██║     ██║  ██║
 ╚════██║██╔═══╝    ██║   ╚════╝╚════██║██║     ██╔══██║██╔══╝  ██╔══╝  ██║   ██║██║     ██║  ██║
 ███████║██║        ██║         ███████║╚██████╗██║  ██║██║     ██║     ╚██████╔╝███████╗██████╔╝
 ╚══════╝╚═╝        ╚═╝         ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝     ╚═╝      ╚═════╝ ╚══════╝╚═════╝ `

func tickCmd() tea.Cmd {
	return tea.Tick(600*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func updateWelcome(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.blink = !m.blink
		return m, tickCmd()
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.state = StateModType
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func viewWelcome(m Model) string {
	width := m.termW
	if width < 10 {
		width = 100
	}

	art := styles.Banner.Render(asciiArt)

	subtitle := styles.Subtitle.Render("Mod Scaffolder for SPT 4.0")

	var cursor string
	if m.blink {
		cursor = styles.Hint.Render("▋")
	} else {
		cursor = styles.Hint.Render(" ")
	}

	prompt := styles.Hint.Render("Press ") +
		lipgloss.NewStyle().Foreground(styles.ColorAmber).Bold(true).Render("ENTER") +
		styles.Hint.Render(" to begin  ") +
		cursor

	var sptLine string
	if len(m.sptVersions) == 0 {
		sptLine = styles.Hint.Render("Fetching SPT versions...")
	} else {
		sptLine = styles.Hint.Render("Latest SPT: ") +
			lipgloss.NewStyle().Foreground(styles.ColorAmberLight).Bold(true).Render(m.sptVersions[0])
	}

	divider := lipgloss.NewStyle().
		Foreground(styles.ColorBorder).
		Render(strings.Repeat("─", min(width-4, 100)))

	version := styles.Version.Render(fmt.Sprintf("spt-scaffold %s", AppVersion))

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		art,
		"",
		subtitle,
		"",
		sptLine,
		"",
		divider,
		"",
		prompt,
		"",
	)

	outer := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(content)

	versionLine := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Right).
		PaddingRight(2).
		Render(version)

	return outer + "\n" + versionLine
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
