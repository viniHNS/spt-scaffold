package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"spt-scaffold/internal/config"
	"spt-scaffold/internal/scaffold"
	"spt-scaffold/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
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
		buildResourcesSection(),
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderTree(cfg config.ModConfig) string {
	cwd, _ := os.Getwd()
	rootPath := filepath.Join(cwd, cfg.ModName)

	lines := []string{
		styles.TreeDir.Render(rootPath + "/"),
	}

	names := scaffold.FileNames(cfg)
	root := buildNodeTree(names)

	// Sort root children: dirs first, then files alphabetically
	var children []*treeNode
	for _, c := range root.children {
		children = append(children, c)
	}
	sort.Slice(children, func(i, j int) bool {
		if children[i].isDir != children[j].isDir {
			return children[i].isDir
		}
		return children[i].name < children[j].name
	})

	for i, c := range children {
		lines = append(lines, formatTree(c, "  ", i == len(children)-1)...)
	}

	return strings.Join(lines, "\n")
}

type treeNode struct {
	name     string
	isDir    bool
	children map[string]*treeNode
}

func buildNodeTree(paths []string) *treeNode {
	root := &treeNode{children: make(map[string]*treeNode), isDir: true}
	for _, p := range paths {
		parts := strings.Split(filepath.ToSlash(p), "/")
		curr := root
		for i, part := range parts {
			if _, ok := curr.children[part]; !ok {
				curr.children[part] = &treeNode{
					name:     part,
					isDir:    i < len(parts)-1,
					children: make(map[string]*treeNode),
				}
			}
			curr = curr.children[part]
		}
	}
	return root
}

func formatTree(node *treeNode, prefix string, isLast bool) []string {
	var lines []string
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	nameFmt := node.name
	if node.isDir {
		nameFmt += "/"
		lines = append(lines, styles.TreeDir.Render(prefix+connector+nameFmt))
	} else {
		lines = append(lines, styles.TreeLine.Render(prefix+connector+nameFmt))
	}

	if node.isDir {
		newPrefix := prefix
		if !isLast {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}

		var children []*treeNode
		for _, c := range node.children {
			children = append(children, c)
		}
		sort.Slice(children, func(i, j int) bool {
			if children[i].isDir != children[j].isDir {
				return children[i].isDir
			}
			return children[i].name < children[j].name
		})

		for i, c := range children {
			lines = append(lines, formatTree(c, newPrefix, i == len(children)-1)...)
		}
	}

	return lines
}

func buildDoneMarkdown(cfg config.ModConfig) (string, error) {
	md := "## Next Steps\n\n" +
		"### Open in your IDE\n" +
		"- **Visual Studio**: open `" + cfg.ModName + ".sln` → File > Open > Project/Solution\n" +
		"- **JetBrains Rider**: open the folder or `" + cfg.ModName + ".sln` directly\n\n" +
		"### Build the mod\n" +
		"```sh\n" +
		"cd " + cfg.ModName + "\n" +
		"dotnet build -c Release\n" +
		"```\n"

	if cfg.ModType == config.ModTypeClient {
		dllPath := filepath.Join(cfg.SptInstallPath, "BepInEx", "plugins", cfg.ModName+".dll")
		md += "\nThe PostBuild target automatically copies the compiled DLL to:\n" +
			"`" + dllPath + "`\n"
	}

	return md, nil
}

func osc8Link(text, url string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, text)
}

func buildResourcesSection() string {
	type resource struct {
		name string
		text string
		url  string
	}
	resources := []resource{
		{"SPT Server (C#) Overview", "deepwiki.com/sp-tarkov", "https://deepwiki.com/sp-tarkov/server-csharp/1-overview"},
		{"Server Mod Examples", "github.com/sp-tarkov", "https://github.com/sp-tarkov/server-mod-examples"},
		{"SPT Wiki Modding Resources", "wiki.sp-tarkov.com", "https://wiki.sp-tarkov.com/modding/Modding_Resources"},
		{"SPT Client Mod Examples", "Jehree/SPTClientModExamples", "https://github.com/Jehree/SPTClientModExamples"},
	}

	heading := lipgloss.NewStyle().Foreground(styles.ColorAmber).Bold(true).Render("  Resources")
	nameStyle := lipgloss.NewStyle().Foreground(styles.ColorAmberLight).Width(30)
	linkStyle := lipgloss.NewStyle().Foreground(styles.ColorAmberDark).Underline(true)

	hint := lipgloss.NewStyle().Foreground(styles.ColorMuted).Italic(true).Render("  Ctrl+Click to open in browser")

	lines := []string{"", heading}
	for _, r := range resources {
		link := osc8Link(linkStyle.Render(r.text), r.url)
		lines = append(lines, fmt.Sprintf("  %s %s", nameStyle.Render(r.name), link))
	}
	lines = append(lines, "", hint)
	return strings.Join(lines, "\n")
}
