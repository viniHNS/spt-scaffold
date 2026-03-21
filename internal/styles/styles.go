package styles

import "github.com/charmbracelet/lipgloss"

// Color palette — Tarkov amber/orange feel on dark backgrounds.
const (
	ColorAmber      = lipgloss.Color("#F59E0B")
	ColorAmberLight = lipgloss.Color("#FCD34D")
	ColorAmberDark  = lipgloss.Color("#B45309")
	ColorBg         = lipgloss.Color("#0D0D0D")
	ColorSurface    = lipgloss.Color("#1A1A1A")
	ColorBorder     = lipgloss.Color("#3D2B00")
	ColorMuted      = lipgloss.Color("#6B7280")
	ColorError      = lipgloss.Color("#EF4444")
	ColorSuccess    = lipgloss.Color("#22C55E")
	ColorWhite      = lipgloss.Color("#F9FAFB")
	ColorSubtle     = lipgloss.Color("#4B5563")
)

// App-wide container.
var App = lipgloss.NewStyle().
	Background(ColorBg).
	Foreground(ColorWhite)

// Banner / ASCII art title.
var Banner = lipgloss.NewStyle().
	Foreground(ColorAmber).
	Bold(true)

// Subtitle under the banner.
var Subtitle = lipgloss.NewStyle().
	Foreground(ColorAmberLight).
	Italic(true)

// Version tag (bottom-right).
var Version = lipgloss.NewStyle().
	Foreground(ColorMuted).
	Italic(true)

// Hint text (e.g. "Press ENTER to continue").
var Hint = lipgloss.NewStyle().
	Foreground(ColorAmberDark).
	Bold(false)

// FormBox is the outer border for form steps.
var FormBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorAmber).
	Padding(1, 3)

// Label is the field name label.
var Label = lipgloss.NewStyle().
	Foreground(ColorAmberLight).
	Bold(true)

// InputActive is the text input in active state.
var InputActive = lipgloss.NewStyle().
	Foreground(ColorWhite)

// ValidationError renders inline validation messages.
var ValidationError = lipgloss.NewStyle().
	Foreground(ColorError).
	Bold(true)

// CharCounter renders the live character counter.
var CharCounter = lipgloss.NewStyle().
	Foreground(ColorMuted)

// StepIndicator "Step X of 7".
var StepIndicator = lipgloss.NewStyle().
	Foreground(ColorMuted).
	Italic(true)

// ConfirmBox is the summary panel on the confirmation screen.
var ConfirmBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorAmberDark).
	Padding(1, 3)

// ConfirmKey is the key label [G] / [E].
var ConfirmKey = lipgloss.NewStyle().
	Foreground(ColorAmber).
	Bold(true)

// ConfirmAction is the action description text next to the key.
var ConfirmAction = lipgloss.NewStyle().
	Foreground(ColorWhite)

// FieldName is the label inside the confirm summary.
var FieldName = lipgloss.NewStyle().
	Foreground(ColorMuted).
	Width(18)

// FieldValue is the value inside the confirm summary.
var FieldValue = lipgloss.NewStyle().
	Foreground(ColorAmberLight)

// ProgressTitle is the heading on the progress screen.
var ProgressTitle = lipgloss.NewStyle().
	Foreground(ColorAmber).
	Bold(true)

// ProgressFile is a completed file line with checkmark.
var ProgressFile = lipgloss.NewStyle().
	Foreground(ColorSuccess)

// DoneTitle is the big success heading.
var DoneTitle = lipgloss.NewStyle().
	Foreground(ColorSuccess).
	Bold(true)

// TreeLine is styling for file tree entries.
var TreeLine = lipgloss.NewStyle().
	Foreground(ColorAmberLight)

// TreeDir is styling for directory entries in the tree.
var TreeDir = lipgloss.NewStyle().
	Foreground(ColorAmber).
	Bold(true)

// QuitHint "Press Q to quit".
var QuitHint = lipgloss.NewStyle().
	Foreground(ColorMuted).
	Italic(true)

// OutputPath displays the target path.
var OutputPath = lipgloss.NewStyle().
	Foreground(ColorAmberLight).
	Bold(true)

// SelectItem is an unselected list item.
var SelectItem = lipgloss.NewStyle().
	Foreground(ColorMuted)

// SelectItemActive is the highlighted list item.
var SelectItemActive = lipgloss.NewStyle().
	Foreground(ColorAmber).
	Bold(true)

// WipItem is a non-selectable list item displayed as coming soon.
var WipItem = lipgloss.NewStyle().
	Foreground(ColorMuted).
	Italic(true)
