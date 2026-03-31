package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/tryl/internal/tui/styles"
)

// RenderList renders the main list view.
func RenderList(
	width, height int,
	searchView string,
	searchFocused bool,
	choices []string,
	cursor int,
	selected map[int]struct{},
	isDeleteConfirm bool,
	statusMsg, errMsg string,
	basePath string,
) string {
	// Horizontal margin applied by AppStyle (2 each side = 4 total).
	innerWidth := max(width-4, 10)

	// ── Title bar ────────────────────────────────────────────────────────────
	title := styles.TitleStyle.Width(innerWidth).Render("  tryl  —  " + basePath)

	// ── Search bar ───────────────────────────────────────────────────────────
	barStyle := styles.SearchBarStyle
	if searchFocused {
		barStyle = styles.SearchBarFocusedStyle
	}
	search := barStyle.Width(innerWidth - 2).Render(searchView)

	// ── Folder list items ────────────────────────────────────────────────────
	var listSB strings.Builder
	if len(choices) == 0 {
		listSB.WriteString(styles.LabelStyle.Render("No folders found."))
	} else {
		for i, name := range choices {
			_, marked := selected[i]

			var cursorStr string
			if i == cursor {
				cursorStr = styles.CursorStyle.Render("▶")
			} else {
				cursorStr = " "
			}

			var itemRendered string
			switch {
			case marked:
				itemRendered = styles.DeleteMarkedStyle.Render(name)
			case i == cursor:
				itemRendered = styles.SelectedItemStyle.Render(name)
			default:
				itemRendered = styles.ItemStyle.Render(name)
			}

			fmt.Fprintf(&listSB, "%s %s\n", cursorStr, itemRendered)
		}
	}
	// Remove trailing newline — the border style adds its own padding.
	listContent := strings.TrimRight(listSB.String(), "\n")
	folderList := styles.FolderListStyle.Width(innerWidth - 2).Render(listContent)

	// ── Status / error messages ───────────────────────────────────────────────
	var feedback string
	if errMsg != "" {
		feedback = styles.ErrorStyle.Render("  " + errMsg)
	} else if statusMsg != "" {
		feedback = styles.StatusStyle.Render("  " + statusMsg)
	}

	// ── Assemble main content block ───────────────────────────────────────────
	parts := []string{title, search, folderList}
	if feedback != "" {
		parts = append(parts, feedback)
	}
	content := styles.AppStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, parts...),
	)

	// ── Bottom bar ────────────────────────────────────────────────────────────
	var bottomBar string
	if isDeleteConfirm {
		bottomBar = renderHelpBar([]helpItem{
			{"Enter", "confirm delete"},
			{"Esc", "cancel"},
			{"Ctrl+D", "mark/unmark"},
			{"↑↓", "navigate"},
		}, width)
	} else {
		bottomBar = renderHelpBar([]helpItem{
			{"Enter", "open"},
			{"Ctrl+D", "delete"},
			{"Ctrl+R", "rename"},
			{"Ctrl+M", "move"},
			{"Ctrl+N", "new"},
			{"Ctrl+C", "quit"},
		}, width)
	}

	// ── Pin bottom bar to terminal bottom ─────────────────────────────────────
	return pinToBottom(content, bottomBar, height)
}

// pinToBottom places content at the top and bottomBar at the very last row(s)
// of the terminal, filling the gap with blank lines.
func pinToBottom(content, bottomBar string, termHeight int) string {
	if termHeight <= 0 {
		// Terminal size not yet known — just stack them.
		return lipgloss.JoinVertical(lipgloss.Left, content, bottomBar)
	}
	contentH := lipgloss.Height(content)
	barH := lipgloss.Height(bottomBar)
	gap := termHeight - contentH - barH
	if gap <= 0 {
		// No room — still render both without a spacer.
		return lipgloss.JoinVertical(lipgloss.Left, content, bottomBar)
	}
	return content + strings.Repeat("\n", gap) + bottomBar
}

type helpItem struct {
	key, desc string
}

func renderHelpBar(items []helpItem, width int) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		k := styles.HelpKeyStyle.Render(item.key)
		d := styles.HelpDescStyle.Render(" " + item.desc)
		parts = append(parts, k+d)
	}
	sep := styles.HelpDescStyle.Render("  ·  ")
	bar := strings.Join(parts, sep)
	w := max(width-4, 10)
	return styles.BottomBarStyle.Width(w).Render(bar)
}
