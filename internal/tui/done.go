package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spt-scaffold/internal/config"
	"spt-scaffold/internal/styles"

	"github.com/charmbracelet/glamour"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updateDone(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q", "ctrl+c":
			return m, tea.Quit
		}
	}
	if m.doneReady {
		var cmd tea.Cmd
		m.doneVP, cmd = m.doneVP.Update(msg)
		return m, cmd
	}
	return m, nil
}

func viewDone(m Model) string {
	if !m.doneReady {
		return ""
	}

	width := m.termW
	if width < 10 {
		width = 80
	}

	scrollPct := fmt.Sprintf("%3.f%%", m.doneVP.ScrollPercent()*100)
	header := lipgloss.NewStyle().
		Width(width - 4).
		Align(lipgloss.Right).
		Render(styles.StepIndicator.Render(scrollPct))

	footer := lipgloss.NewStyle().
		Width(width - 4).
		Render(styles.QuitHint.Render("  ↑/↓ or j/k to scroll  •  Q to quit"))

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			header,
			m.doneVP.View(),
			footer,
		))
}

// buildDoneContent assembles the full text that goes inside the viewport.
func buildDoneContent(m Model) string {
	width := m.termW
	if width < 10 {
		width = 80
	}

	title := styles.DoneTitle.Render("✔  Project scaffolded successfully!")
	subline := styles.Subtitle.Render(fmt.Sprintf("  %s is ready.", m.cfg.ModName))
	tree := renderTree(m.cfg)

	nextSteps := ""
	if m.doneMarkdown != "" {
		r, err := glamour.NewTermRenderer(
			glamour.WithStylePath("dark"),
			glamour.WithWordWrap(min(width-8, 100)),
		)
		if err == nil {
			rendered, err := r.Render(m.doneMarkdown)
			if err == nil {
				nextSteps = rendered
			}
		}
	}

	rows := []string{
		title,
		subline,
		"",
		tree,
		"",
		nextSteps,
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderTree(cfg config.ModConfig) string {
	cwd, _ := os.Getwd()
	root := filepath.Join(cwd, cfg.ModName)

	lines := []string{
		styles.TreeDir.Render(root + "/"),
		styles.TreeLine.Render("  ├── " + cfg.ModName + ".csproj"),
		styles.TreeLine.Render("  ├── Mod.cs"),
		styles.TreeLine.Render("  ├── README.md"),
		styles.TreeLine.Render("  └── .gitignore"),
	}
	return strings.Join(lines, "\n")
}


func buildDoneMarkdown(cfg config.ModConfig) (string, error) {
	md := "## Next Steps\n\n" +
		"### Open in your IDE\n" +
		"- **Visual Studio**: open `" + cfg.ModName + ".csproj` → File > Open > Project/Solution\n" +
		"- **JetBrains Rider**: open the folder or `.csproj` directly\n\n" +
		"### Build the mod\n" +
		"```sh\n" +
		"cd " + cfg.ModName + "\n" +
		"dotnet build -c Release\n" +
		"```\n\n" +
		"### Resources\n\n" +
		"| Resource | Link |\n" +
		"|---|---|\n" +
		"| SPT Server (C#) Overview | [deepwiki.com/sp-tarkov](https://deepwiki.com/sp-tarkov/server-csharp/1-overview) |\n" +
		"| Server Mod Examples | [github.com/sp-tarkov](https://github.com/sp-tarkov/server-mod-examples) |\n" +
		"| SPT Wiki Modding Resources | [wiki.sp-tarkov.com](https://wiki.sp-tarkov.com/modding/Modding_Resources) |\n" +
		"| SPT Client Mod Examples | [Jehree/SPTClientModExamples](https://github.com/Jehree/SPTClientModExamples) |\n"
	return md, nil
}
