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
	isList       bool // license picker (step 9)
	isSptVersion bool // SPT version picker (step 6)
	isModType    bool // mod type picker (step 0)
	isTemplate   bool // template picker (step 1)
	isSptPath    bool // SPT install path text input (step 2, client only)
}

// fields defines every form step in order.
// Steps 0 (isModType), 1 (isTemplate), 6 (isSptVersion), 9 (isList) are pickers, not text inputs.
// Step 2 (isSptPath) is a text input shown only to client mod users.
var fields = []fieldMeta{
	{label: "Mod Type", isModType: true},                                                            // step 0
	{label: "Template", isTemplate: true},                                                           // step 1
	{label: "SPT Install Path", placeholder: `C:\SPT`, isSptPath: true},                            // step 2 (client only)
	{label: "Mod Name", placeholder: "MyAwesomeMod"},                                               // step 3
	{label: "Author", placeholder: "YourUsername"},                                                  // step 4
	{label: "Version", placeholder: "1.0.0", defaultVal: "1.0.0"},                                  // step 5
	{label: "SPT Version", isSptVersion: true},                                                      // step 6
	{label: "Description", placeholder: "Short description (optional, max 120 chars)", maxLen: 120}, // step 7
	{label: "Repository URL", placeholder: "https://github.com/you/mod (optional)"},                // step 8
	{label: "License", isList: true},                                                                // step 9
}

// Form step indices — update these if the fields slice order changes.
const (
	stepTemplate = 1 // isTemplate picker
	stepSptPath  = 2 // isSptPath text input (client only)
	stepModName  = 3 // first text input
)

func makeInputs() []textinput.Model {
	inputs := make([]textinput.Model, len(fields))
	for i, f := range fields {
		if f.isList || f.isSptVersion || f.isModType || f.isTemplate {
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
	f := fields[step]
	return !f.isList && !f.isSptVersion && !f.isModType && !f.isTemplate
}

func validateStep(m Model, step int) error {
	switch step {
	case 2:
		return config.ValidateSptInstallPath(m.inputs[2].Value())
	case 3:
		return config.ValidateModName(m.inputs[3].Value())
	case 4:
		return config.ValidateAuthor(m.inputs[4].Value())
	case 5:
		v := m.inputs[5].Value()
		if v == "" {
			v = "1.0.0"
		}
		return config.ValidateSemver(v)
	case 7:
		return config.ValidateDescription(m.inputs[7].Value())
	case 8:
		return config.ValidateRepoURL(m.inputs[8].Value())
	}
	return nil // steps 0, 1, 6, 9: pickers; step 2 for server: unreachable via navigation
}

func updateForm(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	step := m.formStep

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			isServer := config.ModTypes[m.modTypeIdx].Value == config.ModTypeServer
			firstStep := stepTemplate

			if step <= firstStep {
				if isTextStep(step) {
					m.inputs[step].Blur()
				}
				m.formErrors[step] = ""
				m.state = StateModType
				return m, nil
			}
			m.formErrors[step] = ""
			if isTextStep(step) {
				m.inputs[step].Blur()
			}

			if step == stepModName && isServer {
				m.formStep = stepTemplate
			} else {
				m.formStep = step - 1
			}

			prev := m.formStep
			if isTextStep(prev) {
				m.inputs[prev].Focus()
			}
			return m, textinput.Blink

		case "enter":
			// Version default fill — step 5
			if step == 5 && m.inputs[5].Value() == "" {
				m.inputs[5].SetValue("1.0.0")
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

			// Server skips SPT Install Path (step 2) when advancing from Template (step 1)
			if step == stepTemplate && config.ModTypes[m.modTypeIdx].Value == config.ModTypeServer {
				m.formStep = stepModName
			} else {
				m.formStep++
			}

			next := m.formStep
			if isTextStep(next) {
				m.inputs[next].Focus()
			}
			return m, textinput.Blink

		case "up":
			switch {
			case fields[step].isModType:
				newIdx := m.modTypeIdx - 1
				if newIdx >= 0 && !config.ModTypes[newIdx].Disabled {
					m.modTypeIdx = newIdx
					m.templateIdx = 0
				}
			case fields[step].isTemplate && m.templateIdx > 0:
				m.templateIdx--
			case fields[step].isSptVersion && m.sptVersionIdx > 0:
				m.sptVersionIdx--
			case fields[step].isList && m.licenseIdx > 0:
				m.licenseIdx--
			}
			return m, nil

		case "down":
			switch {
			case fields[step].isModType:
				newIdx := m.modTypeIdx + 1
				if newIdx < len(config.ModTypes) && !config.ModTypes[newIdx].Disabled {
					m.modTypeIdx = newIdx
					m.templateIdx = 0
				}
			case fields[step].isTemplate:
				tmpl := m.activeTemplateList()
				if m.templateIdx < len(tmpl)-1 {
					m.templateIdx++
				}
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
	case f.isModType:
		var lines []string
		for i, mt := range config.ModTypes {
			if mt.Disabled {
				lines = append(lines, lipgloss.NewStyle().Foreground(styles.ColorMuted).
					Render("  "+mt.Label+"  (coming soon)"))
			} else if i == m.modTypeIdx {
				lines = append(lines, styles.SelectItemActive.Render("▶ "+mt.Label))
			} else {
				lines = append(lines, styles.SelectItem.Render("  "+mt.Label))
			}
		}
		inputView = strings.Join(lines, "\n")
		extra = styles.Hint.Render("↑/↓ to select, ENTER to confirm")

	case f.isTemplate:
		tmpl := m.activeTemplateList()
		var lines []string
		for i, t := range tmpl {
			if i == m.templateIdx {
				lines = append(lines, styles.SelectItemActive.Render("▶ "+t.Label))
			} else {
				lines = append(lines, styles.SelectItem.Render("  "+t.Label))
			}
		}
		inputView = strings.Join(lines, "\n")
		tmplIdx := m.templateIdx
		if len(tmpl) == 0 {
			break
		}
		if tmplIdx >= len(tmpl) {
			tmplIdx = len(tmpl) - 1
		}
		extra = styles.Hint.Render(tmpl[tmplIdx].Desc)

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
		if fields[step].maxLen > 0 {
			count := len(m.inputs[step].Value())
			extra = styles.CharCounter.Render(fmt.Sprintf("%d / %d", count, fields[step].maxLen))
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
