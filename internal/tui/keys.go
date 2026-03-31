package tui

import tea "github.com/charmbracelet/bubbletea"

// Key binding helpers using bubbletea v1 KeyType constants.

func isQuit(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyCtrlC
}

func isDelete(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyCtrlD
}

func isRename(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyCtrlR
}

func isMove(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyCtrlM
}

func isCreate(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyCtrlN
}

func isConfirm(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyEnter
}

func isCancel(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyEsc
}

func isUp(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyUp
}

func isDown(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyDown
}

// isTypable returns true if the key is a printable rune (for triggering search focus).
func isTypable(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyRunes
}
