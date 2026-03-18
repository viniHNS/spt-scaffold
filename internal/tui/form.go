package tui

import (
	"fmt"
	"strings"

	"spt-scaffold/internal/config"
	"spt-scaffold/internal/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// fieldMeta describes one form step.
type fieldMeta struct {
	label        string
	placeholder  string
	defaultVal   string
	maxLen       int
	isList       bool // license picker
	isSptVersion bool // SPT version picker
}

// fields defines every form step in order.
// Steps 3 (isSptVersion) and 6 (isList) are pickers, not text inputs.
var fields = []fieldMeta{
	{label: "Mod Name", placeholder: "MyAwesomeMod"},
	{label: "Author", placeholder: "YourUsername"},
	{label: "Version", placeholder: "1.0.0", defaultVal: "1.0.0"},
	{label: "SPT Version", isSptVersion: true},
	{label: "Description", placeholder: "Short description (optional, max 120 chars)", maxLen: 120},
	{label: "Repository URL", placeholder: "https://github.com/you/mod (optional)"},
	{label: "License", isList: true},
}

func makeInputs() []textinput.Model {
	inputs := make([]textinput.Model, len(fields))
	for i, f := range fields {
		if f.isList || f.isSptVersion {
			continue
		}
		ti := textinput.New()
		ti.Placeholder = f.placeholder
		if f.defaultVal != "" {
			ti.SetValue(f.defaultVal)
		}
		if f.maxLen > 0 {
			ti.CharLimit = f.maxLen
		}
		inputs[i] = ti
	}
	return inputs
}

func spinnerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(styles.ColorAmber)
}

func textinputBlink() tea.Cmd {
	return textinput.Blink
}

// isTextStep returns true when the step uses a text input (not a picker).
func isTextStep(step int) bool {
	if step < 0 || step >= len(fields) {
		return false
	}
	return !fields[step].isList && !fields[step].isSptVersion
}

func validateStep(m Model, step int) error {
	switch step {
	case 0:
		return config.ValidateModName(m.inputs[0].Value())
	case 1:
		return config.ValidateAuthor(m.inputs[1].Value())
	case 2:
		v := m.inputs[2].Value()
		if v == "" {
			v = "1.0.0"
		}
		return config.ValidateSemver(v)
	case 3:
		return nil // SPT version picker is always valid
	case 4:
		return config.ValidateDescription(m.inputs[4].Value())
	case 5:
		return config.ValidateRepoURL(m.inputs[5].Value())
	case 6:
		return nil // license picker is always valid
	}
	return nil
}

func updateForm(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	step := m.formStep

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			if step <= 0 {
				return m, textinput.Blink
			}
			m.formErrors[step] = ""
			if isTextStep(step) {
				m.inputs[step].Blur()
			}
			m.formStep = step - 1
			if isTextStep(step - 1) {
				m.inputs[step-1].Focus()
			}
			return m, textinput.Blink

		case "enter":
			if step == 2 && m.inputs[2].Value() == "" {
				m.inputs[2].SetValue("1.0.0")
			}

			if err := validateStep(m, step); err != nil {
				m.formErrors[step] = err.Error()
				return m, nil
			}
			m.formErrors[step] = ""

			if step == len(fields)-1 {
				m.cfg = m.buildConfig()
				m.state = StateConfirm
				return m, nil
			}

			if isTextStep(step) {
				m.inputs[step].Blur()
			}
			m.formStep++
			next := m.formStep
			if isTextStep(next) {
				m.inputs[next].Focus()
			}
			return m, textinput.Blink

		case "up":
			switch {
			case fields[step].isSptVersion && m.sptVersionIdx > 0:
				m.sptVersionIdx--
			case fields[step].isList && m.licenseIdx > 0:
				m.licenseIdx--
			}
			return m, nil

		case "down":
			switch {
			case fields[step].isSptVersion && len(m.sptVersions) > 0 && m.sptVersionIdx < len(m.sptVersions)-1:
				m.sptVersionIdx++
			case fields[step].isList && m.licenseIdx < len(config.Licenses)-1:
				m.licenseIdx++
			}
			return m, nil
		}
	}

	// Route key events to the active text input.
	if isTextStep(step) {
		var cmd tea.Cmd
		m.inputs[step], cmd = m.inputs[step].Update(msg)
		return m, cmd
	}

	return m, nil
}

// sptVersionWindow returns the slice indices [start, end) to show in the picker.
// Keeps the selected item visible within a window of size `size`.
func sptVersionWindow(selected, total, size int) (int, int) {
	if total <= size {
		return 0, total
	}
	half := size / 2
	start := selected - half
	if start < 0 {
		start = 0
	}
	end := start + size
	if end > total {
		end = total
		start = end - size
	}
	return start, end
}

func viewForm(m Model) string {
	step := m.formStep
	f := fields[step]
	total := len(fields)

	width := m.termW
	if width < 10 {
		width = 80
	}

	stepStr := styles.StepIndicator.Render(fmt.Sprintf("Step %d of %d", step+1, total))
	header := lipgloss.NewStyle().Width(width - 4).Align(lipgloss.Right).Render(stepStr)

	label := styles.Label.Render(f.label)

	var inputView, extra string

	switch {
	case f.isSptVersion:
		if len(m.sptVersions) == 0 {
			inputView = styles.Hint.Render("  Fetching available versions...")
			extra = styles.Hint.Render("Please wait…")
		} else {
			const windowSize = 6
			start, end := sptVersionWindow(m.sptVersionIdx, len(m.sptVersions), windowSize)
			var lines []string
			if start > 0 {
				lines = append(lines, styles.Hint.Render(fmt.Sprintf("  ↑ %d more above", start)))
			}
			for i := start; i < end; i++ {
				v := m.sptVersions[i]
				label := v
				if i == 0 {
					label += "  " + lipgloss.NewStyle().Foreground(styles.ColorSuccess).Render("(latest)")
				}
				if i == m.sptVersionIdx {
					lines = append(lines, styles.SelectItemActive.Render("▶ "+label))
				} else {
					lines = append(lines, styles.SelectItem.Render("  "+label))
				}
			}
			remaining := len(m.sptVersions) - end
			if remaining > 0 {
				lines = append(lines, styles.Hint.Render(fmt.Sprintf("  ↓ %d more below", remaining)))
			}
			inputView = strings.Join(lines, "\n")
			extra = styles.Hint.Render("↑/↓ to select, ENTER to confirm, ESC to go back")
		}

	case f.isList:
		var lines []string
		for i, lic := range config.Licenses {
			if i == m.licenseIdx {
				lines = append(lines, styles.SelectItemActive.Render("▶ "+lic.Label))
			} else {
				lines = append(lines, styles.SelectItem.Render("  "+lic.Label))
			}
		}
		inputView = strings.Join(lines, "\n")
		extra = styles.Hint.Render("↑/↓ to select, ENTER to confirm")

	default:
		inputView = styles.InputActive.Render(m.inputs[step].View())
		if step == 4 { // Description
			count := len(m.inputs[4].Value())
			extra = styles.CharCounter.Render(fmt.Sprintf("%d / 120", count))
		} else {
			extra = styles.Hint.Render("ENTER to continue, ESC to go back")
		}
	}

	errLine := ""
	if m.formErrors[step] != "" {
		errLine = "\n" + styles.ValidationError.Render("✗ "+m.formErrors[step])
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		label,
		"",
		inputView,
		extra+errLine,
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
