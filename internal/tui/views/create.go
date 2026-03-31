package views

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/user/tryl/internal/tui/styles"
	"strings"
)

// RenderCreate renders the create-folder dialog.
func RenderCreate(width, height int, nameView, repoView string, focusRepo bool, errMsg string) string {
	var sb strings.Builder

	sb.WriteString(styles.DialogTitleStyle.Render("Create New Folder"))
	sb.WriteString("\n")

	// Name field.
	nameLabel := styles.LabelStyle.Render("Folder name:")
	if !focusRepo {
		nameLabel = styles.HelpKeyStyle.Render("Folder name:")
	}
	sb.WriteString(nameLabel)
	sb.WriteString("\n")
	sb.WriteString(nameView)
	sb.WriteString("\n\n")

	// Repo field.
	repoLabel := styles.LabelStyle.Render("GitHub repo URL (optional):")
	if focusRepo {
		repoLabel = styles.HelpKeyStyle.Render("GitHub repo URL (optional):")
	}
	sb.WriteString(repoLabel)
	sb.WriteString("\n")
	sb.WriteString(repoView)
	sb.WriteString("\n")

	if errMsg != "" {
		sb.WriteString("\n" + styles.ErrorStyle.Render(errMsg) + "\n")
	}

	sb.WriteString(styles.HelpDescStyle.Render("\nTab to switch fields · Enter to advance/confirm · Esc to cancel"))

	bottomBar := renderHelpBar([]helpItem{
		{"Tab", "switch field"},
		{"Enter", "next / confirm"},
		{"Esc", "cancel"},
	}, width)

	dialogW := clamp(width-10, 40, 80)
	dialog := styles.DialogStyle.Width(dialogW).Render(sb.String())

	content := lipgloss.Place(width, height-lipgloss.Height(bottomBar), lipgloss.Center, lipgloss.Center, dialog)
	return pinToBottom(content, bottomBar, height)
}
