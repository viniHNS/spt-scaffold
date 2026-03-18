package tui

import (
	"fmt"
	"time"

	"spt-scaffold/internal/config"
	"spt-scaffold/internal/scaffold"
	"spt-scaffold/internal/styles"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// fileStep is one file to generate with a delay.
type fileStep struct {
	name string
	done bool
}

// progressState tracks per-file generation.
type progressModel struct {
	steps   []fileStep
	current int
	cfg     config.ModConfig
	err     error
}

var generationSteps []string

// generateFileCmd generates one file and returns a msg.
func generateFileCmd(cfg config.ModConfig, idx int) tea.Cmd {
	return func() (msg tea.Msg) {
		defer func() {
			if r := recover(); r != nil {
				msg = scaffoldDoneMsg{err: fmt.Errorf("panic: %v", r)}
			}
		}()
		time.Sleep(500 * time.Millisecond) 
		name, err := scaffold.GenerateFile(cfg, idx)
		if err != nil {
			return scaffoldDoneMsg{err: err}
		}
		return fileCreatedMsg{name: name}
	}
}

func updateProgress(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case startScaffoldMsg:
		// Initialise file steps list and start generating the first file.
		names := scaffold.FileNames(msg.cfg)
		m.createdFiles = make([]string, 0, len(names))
		generationSteps = names
		return m, tea.Batch(
			m.spinner.Tick,
			generateFileCmd(msg.cfg, 0),
		)

	case fileCreatedMsg:
		m.createdFiles = append(m.createdFiles, msg.name)
		nextIdx := len(m.createdFiles)
		if nextIdx < len(generationSteps) {
			return m, tea.Batch(m.spinner.Tick, generateFileCmd(m.cfg, nextIdx))
		}
		// All done — build viewport for the done screen.
		md, _ := buildDoneMarkdown(m.cfg)
		m.doneMarkdown = md
		vpW := m.termW - 4
		if vpW < 40 {
			vpW = 40
		}
		vpH := m.termH - 4
		if vpH < 10 {
			vpH = 20
		}
		vp := viewport.New(vpW, vpH)
		vp.SetContent(buildDoneContent(m))
		m.doneVP = vp
		m.doneReady = true
		m.state = StateDone
		return m, nil

	case scaffoldDoneMsg:
		if msg.err != nil {
			m.fatalErr = msg.err
			m.state = StateError
		}
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func viewProgress(m Model) string {
	width := m.termW
	if width < 10 {
		width = 80
	}

	title := styles.ProgressTitle.Render("Generating project files…")

	spinLine := fmt.Sprintf("  %s  %s",
		m.spinner.View(),
		styles.Hint.Render("Working…"),
	)

	var fileLines []string
	for _, f := range m.createdFiles {
		fileLines = append(fileLines, styles.ProgressFile.Render("  ✓ "+f))
	}

	rows := []string{title, "", spinLine, ""}
	rows = append(rows, fileLines...)

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return lipgloss.NewStyle().
		Width(width).
		Padding(2, 2).
		Render(content)
}
