package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fantonucci89/TryLike/internal/tui/styles"
)

// RenderMove renders the move dialog.
func RenderMove(width, height int, inputView, errMsg string) string {
	var sb strings.Builder

	sb.WriteString(styles.DialogTitleStyle.Render("Move Folder"))
	sb.WriteString("\n")
	sb.WriteString(styles.LabelStyle.Render("Destination path (within $HOME):"))
	sb.WriteString("\n")
	sb.WriteString(inputView)
	sb.WriteString("\n")

	if errMsg != "" {
		sb.WriteString("\n" + styles.ErrorStyle.Render(errMsg) + "\n")
	}

	bottomBar := renderHelpBar([]helpItem{
		{"Enter", "confirm"},
		{"Esc", "cancel"},
	}, width)

	dialogW := clamp(width-10, 40, 80)
	dialog := styles.DialogStyle.Width(dialogW).Render(sb.String())

	content := lipgloss.Place(width, height-lipgloss.Height(bottomBar), lipgloss.Center, lipgloss.Center, dialog)
	return pinToBottom(content, bottomBar, height)
}
