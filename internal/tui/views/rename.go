package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/tryl/internal/tui/styles"
)

// RenderRename renders the rename dialog.
func RenderRename(width, height int, inputView, errMsg string) string {
	var sb strings.Builder

	sb.WriteString(styles.DialogTitleStyle.Render("Rename Folder"))
	sb.WriteString("\n")
	sb.WriteString(styles.LabelStyle.Render("New name:"))
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

	dialogW := clamp(width-10, 30, 70)
	dialog := styles.DialogStyle.Width(dialogW).Render(sb.String())

	content := lipgloss.Place(width, height-lipgloss.Height(bottomBar), lipgloss.Center, lipgloss.Center, dialog)
	return pinToBottom(content, bottomBar, height)
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
