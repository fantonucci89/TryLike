package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Color palette.
	colorPrimary = lipgloss.Color("63")  // blue-violet
	colorAccent  = lipgloss.Color("205") // pink
	colorMuted   = lipgloss.Color("240") // grey
	colorDanger  = lipgloss.Color("196") // red
	colorSuccess = lipgloss.Color("82")  // green
	colorFg      = lipgloss.Color("252") // near-white

	// App outer frame — horizontal margin only; vertical handled manually.
	AppStyle = lipgloss.NewStyle().
			Margin(0, 2)

	// Title / header bar.
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorFg).
			Background(colorPrimary).
			Padding(0, 1).
			MarginBottom(1)

	// Search bar wrapper.
	SearchBarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1).
			MarginBottom(1)

	SearchBarFocusedStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(colorAccent).
				Padding(0, 1).
				MarginBottom(1)

	// Folder list bordered container.
	FolderListStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1)

	// List item styles.
	ItemStyle = lipgloss.NewStyle().
			Foreground(colorFg)

	SelectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorAccent)

	DeleteMarkedStyle = lipgloss.NewStyle().
				Strikethrough(true).
				Foreground(colorDanger)

	// Cursor.
	CursorStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	// Bottom help bar — no top margin; pinned via spacer in the view.
	BottomBarStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(colorMuted).
			Padding(0, 1)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	// Dialog / overlay.
	DialogStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(1, 2).
			MarginTop(1)

	DialogTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorPrimary).
				MarginBottom(1)

	// Error / status.
	ErrorStyle = lipgloss.NewStyle().
			Foreground(colorDanger).
			Bold(true)

	StatusStyle = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Italic(true)

	// Label in dialogs.
	LabelStyle = lipgloss.NewStyle().
			Foreground(colorMuted)
)
