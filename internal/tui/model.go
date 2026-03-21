package tui

import (
	"spt-scaffold/internal/config"
	"spt-scaffold/internal/scaffold"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// AppVersion is the displayed tool version.
const AppVersion = "v0.1.0"

// Screen states for the TUI state machine.
type state int

const (
	StateWelcome        state = iota
	StateModType        state = iota
	StateForm           state = iota
	StateConfirm        state = iota
	StateProgress       state = iota
	StateDone           state = iota
	StateError          state = iota
)

// tickMsg is sent by the blink ticker on the welcome screen.
type tickMsg struct{}

// sptVersionMsg carries the result of the NuGet versions fetch.
type sptVersionMsg struct {
	versions []string
	err      error
}

// fileCreatedMsg reports that a file was created during scaffolding.
type fileCreatedMsg struct{ name string }

// scaffoldDoneMsg signals that all files have been written.
type scaffoldDoneMsg struct{ err error }

// Model is the root Bubbletea model.
type Model struct {
	state state

	// Welcome
	blink bool
	termW int
	termH int

	// SPT versions (fetched from NuGet, descending order)
	sptVersions   []string
	sptVersionIdx int

	// Form
	formStep    int
	inputs      []textinput.Model
	formErrors  []string
	licenseIdx  int
	templateIdx int // index into the active template list (ServerTemplates or ClientTemplates)

	// Mod type + template selection
	modTypeIdx    int    // 0=Server, 1=Client
	modTypeHint   string // shown when user attempts to select WIP Client

	// Confirm
	cfg config.ModConfig

	// Progress
	spinner      spinner.Model
	createdFiles []string

	// Done
	doneMarkdown string
	doneVP       viewport.Model
	doneReady    bool

	// Error
	fatalErr error
}

// NewModel creates and returns the initial model.
func NewModel() Model {
	inputs := makeInputs()

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = spinnerStyle()

	return Model{
		state:      StateWelcome,
		blink:      true,
		inputs:     inputs,
		formErrors: make([]string, len(fields)),
		licenseIdx: 0,
		spinner:    sp,
	}
}

// Init starts any initial commands.
func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), textinput.Blink, fetchSptVersionsCmd())
}

// fetchSptVersionsCmd queries NuGet for all available SPT versions in the background.
func fetchSptVersionsCmd() tea.Cmd {
	return func() tea.Msg {
		versions, err := config.FetchSptVersions()
		return sptVersionMsg{versions: versions, err: err}
	}
}

// Update handles all messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termW = msg.Width
		m.termH = msg.Height
		if m.doneReady {
			m.doneVP.Width = msg.Width - 4
			m.doneVP.Height = msg.Height - 4
		}
		return m, nil
	case sptVersionMsg:
		if msg.err == nil && len(msg.versions) > 0 {
			m.sptVersions = msg.versions
		} else {
			m.sptVersions = []string{config.FallbackSptVersion}
		}
		return m, nil
	}

	switch m.state {
	case StateWelcome:
		return updateWelcome(m, msg)
	case StateModType:
		return updateModType(m, msg)
	case StateForm:
		return updateForm(m, msg)
	case StateConfirm:
		return updateConfirm(m, msg)
	case StateProgress:
		return updateProgress(m, msg)
	case StateDone:
		return updateDone(m, msg)
	case StateError:
		return updateError(m, msg)
	}
	return m, nil
}

// View renders the current screen.
func (m Model) View() string {
	switch m.state {
	case StateWelcome:
		return viewWelcome(m)
	case StateModType:
		return viewModType(m)
	case StateForm:
		return viewForm(m)
	case StateConfirm:
		return viewConfirm(m)
	case StateProgress:
		return viewProgress(m)
	case StateDone:
		return viewDone(m)
	case StateError:
		return viewError(m)
	}
	return ""
}

// activeTemplateList returns the correct template slice for the currently selected mod type,
// with the templateIdx clamped to a valid index.
func (m Model) activeTemplateList() []config.TemplateEntry {
	if config.ModTypes[m.modTypeIdx].Value == config.ModTypeClient {
		return config.ClientTemplates
	}
	return config.ServerTemplates
}

// selectedSptVersion returns the currently selected SPT version.
func (m Model) selectedSptVersion() string {
	if len(m.sptVersions) == 0 {
		return config.FallbackSptVersion
	}
	if m.sptVersionIdx < 0 || m.sptVersionIdx >= len(m.sptVersions) {
		return m.sptVersions[0]
	}
	return m.sptVersions[m.sptVersionIdx]
}

// buildConfig populates a ModConfig from the form inputs.
func (m Model) buildConfig() config.ModConfig {
	activeTemplates := m.activeTemplateList()

	tmplIdx := m.templateIdx
	if len(activeTemplates) == 0 || tmplIdx >= len(activeTemplates) {
		tmplIdx = 0
	}

	v := m.inputs[5].Value() // Version (fields[5])
	if v == "" {
		v = "1.0.0"
	}
	sptVer := m.selectedSptVersion()
	modType := config.ModTypes[m.modTypeIdx].Value
	modTemplate := activeTemplates[tmplIdx].Value

	return config.ModConfig{
		ModName:         m.inputs[3].Value(),
		Author:          m.inputs[4].Value(),
		Version:         v,
		SptVersion:      sptVer,
		SptVersionRange: config.SptVersionRange(sptVer),
		Desc:            m.inputs[7].Value(),
		RepoURL:         m.inputs[8].Value(),
		License:         config.Licenses[m.licenseIdx].SPDX,
		SptInstallPath:  m.inputs[2].Value(),
		ProjectGuid:     config.NewProjectGuid(),
		ModType:         modType,
		ModTemplate:     modTemplate,
	}
}

// runScaffold executes the scaffold generator and streams messages back.
func runScaffold(cfg config.ModConfig) tea.Cmd {
	return func() tea.Msg {
		ch := make(chan string, 10)
		errCh := make(chan error, 1)
		go func() {
			errCh <- scaffold.Generate(cfg, ch)
			close(ch)
		}()
		// drain — individual fileCreatedMsg are sent via scaffoldStreamCmd
		return scaffoldDoneMsg{err: <-errCh}
	}
}

// scaffoldStreamCmd streams individual file-created notifications.
func scaffoldStreamCmd(cfg config.ModConfig) tea.Cmd {
	return func() tea.Msg {
		return startScaffoldMsg{cfg: cfg}
	}
}

type startScaffoldMsg struct{ cfg config.ModConfig }
