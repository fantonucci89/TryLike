package tui

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/fantonucci89/TryLike/internal/config"
	"github.com/fantonucci89/TryLike/internal/fs"
	"github.com/fantonucci89/TryLike/internal/tui/views"
)

// viewState tracks which screen/dialog is active.
type viewState int

const (
	stateList          viewState = iota // main list + search
	stateDeleteConfirm                  // folders marked for delete, confirm prompt
	stateRename                         // rename dialog
	stateMove                           // move dialog
	stateCreate                         // create new folder dialog
)

// foldersLoadedMsg carries fresh folder data back to Update.
type foldersLoadedMsg struct {
	folders []fs.Folder
	err     error
}

// opDoneMsg signals that a filesystem operation completed.
type opDoneMsg struct{ err error }

// AppModel is the root Bubble Tea model.
// choices/cursor/selected mirror the spec from the instructions.
type AppModel struct {
	// ── Spec fields ────────────────────────────────────────────
	choices  []string         // folder names in the list (filtered or full)
	cursor   int              // index of cursor in choices
	selected map[int]struct{} // indices of folders marked for delete

	// ── Extended state ─────────────────────────────────────────
	state         viewState
	width, height int
	cfg           *config.Config

	allFolders    []fs.Folder // unfiltered folder list
	searchInput   textinput.Model
	searchFocused bool

	renameInput textinput.Model
	moveInput   textinput.Model

	// create view
	createNameInput textinput.Model
	createRepoInput textinput.Model
	createFocusRepo bool // which input is focused in create view

	// initial filter (from CLI arg)
	initialFilter string

	// feedback
	statusMsg string
	errMsg    string

	// Result: populated when the user presses Enter to open a folder.
	SelectedPath string
}

// New creates a fresh AppModel.
func New(cfg *config.Config, initialFilter string) AppModel {
	// Search input.
	si := textinput.New()
	si.Placeholder = "Search folders…"
	si.Prompt = "  "

	// Rename input.
	ri := textinput.New()
	ri.Placeholder = "New folder name"
	ri.Prompt = "  "
	ri.CharLimit = 128

	// Move input.
	mi := textinput.New()
	mi.Placeholder = "~/destination/path"
	mi.Prompt = "  "
	mi.CharLimit = 256

	// Create name input.
	cni := textinput.New()
	cni.Placeholder = "Folder name"
	cni.Prompt = "  "
	cni.CharLimit = 128

	// Create repo input.
	cri := textinput.New()
	cri.Placeholder = "https://github.com/user/repo (optional)"
	cri.Prompt = "  "
	cri.CharLimit = 256

	m := AppModel{
		selected:        make(map[int]struct{}),
		cfg:             cfg,
		searchInput:     si,
		renameInput:     ri,
		moveInput:       mi,
		createNameInput: cni,
		createRepoInput: cri,
		initialFilter:   initialFilter,
	}

	if initialFilter != "" {
		m.searchInput.SetValue(initialFilter)
	}

	return m
}

// ── Tea interface ───────────────────────────────────────────────────────────

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		loadFolders(m.cfg.BasePath),
		textinput.Blink,
	)
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case foldersLoadedMsg:
		if msg.err != nil {
			m.errMsg = fmt.Sprintf("Error loading folders: %v", msg.err)
			return m, nil
		}
		m.allFolders = msg.folders
		m.applyFilter()
		if m.initialFilter != "" {
			m.initialFilter = "" // consume once
		}
		return m, nil

	case opDoneMsg:
		if msg.err != nil {
			m.errMsg = fmt.Sprintf("Error: %v", msg.err)
		} else {
			m.errMsg = ""
		}
		// Reload folders after any mutating operation.
		return m, loadFolders(m.cfg.BasePath)

	case tea.KeyMsg:
		// Ctrl+C always quits.
		if isQuit(msg) {
			return m, tea.Quit
		}
		return m.handleKey(msg)
	}

	return m, nil
}

func (m AppModel) View() string {
	switch m.state {
	case stateRename:
		return views.RenderRename(
			m.width, m.height,
			m.renameInput.View(),
			m.errMsg,
		)
	case stateMove:
		return views.RenderMove(
			m.width, m.height,
			m.moveInput.View(),
			m.errMsg,
		)
	case stateCreate:
		return views.RenderCreate(
			m.width, m.height,
			m.createNameInput.View(),
			m.createRepoInput.View(),
			m.createFocusRepo,
			m.errMsg,
		)
	default: // stateList / stateDeleteConfirm
		return views.RenderList(
			m.width, m.height,
			m.searchInput.View(),
			m.searchFocused,
			m.choices,
			m.cursor,
			m.selected,
			m.state == stateDeleteConfirm,
			m.statusMsg,
			m.errMsg,
			m.cfg.BasePath,
		)
	}
}

// ── Key routing ─────────────────────────────────────────────────────────────

func (m AppModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateList:
		return m.handleListKey(msg)
	case stateDeleteConfirm:
		return m.handleDeleteConfirmKey(msg)
	case stateRename:
		return m.handleRenameKey(msg)
	case stateMove:
		return m.handleMoveKey(msg)
	case stateCreate:
		return m.handleCreateKey(msg)
	}
	return m, nil
}

// ── List state ───────────────────────────────────────────────────────────────

func (m AppModel) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// When search is focused, most keys go to the input.
	if m.searchFocused {
		switch {
		case isConfirm(msg) || isCancel(msg):
			m.searchFocused = false
			m.searchInput.Blur()
			return m, nil
		default:
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			m.applyFilter()
			return m, cmd
		}
	}

	// Global list keys.
	switch {
	case isConfirm(msg):
		if len(m.choices) > 0 {
			m.SelectedPath = m.folderPath(m.choices[m.cursor])
			return m, tea.Quit
		}
		return m, nil

	case isUp(msg):
		if m.cursor > 0 {
			m.cursor--
		}
		m.statusMsg = ""
		return m, nil

	case isDown(msg):
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
		m.statusMsg = ""
		return m, nil

	case isDelete(msg):
		if len(m.choices) == 0 {
			return m, nil
		}
		if _, ok := m.selected[m.cursor]; ok {
			delete(m.selected, m.cursor)
		} else {
			m.selected[m.cursor] = struct{}{}
		}
		if len(m.selected) > 0 {
			m.state = stateDeleteConfirm
		}
		return m, nil

	case isRename(msg):
		if len(m.choices) == 0 {
			return m, nil
		}
		m.renameInput.SetValue(m.choices[m.cursor])
		m.renameInput.CursorEnd()
		m.state = stateRename
		m.errMsg = ""
		return m, m.renameInput.Focus()

	case isMove(msg):
		if len(m.choices) == 0 {
			return m, nil
		}
		m.moveInput.Reset()
		m.state = stateMove
		m.errMsg = ""
		return m, m.moveInput.Focus()

	case isCreate(msg):
		m.createNameInput.Reset()
		m.createRepoInput.Reset()
		m.createFocusRepo = false
		m.state = stateCreate
		m.errMsg = ""
		return m, m.createNameInput.Focus()

	case isTypable(msg):
		// Focus the search input synchronously first (pointer receiver mutates
		// m.searchInput.focus = true), then forward the key through Update so
		// the rune is accepted immediately — not dropped on the next tick.
		m.searchFocused = true
		blinkCmd := m.searchInput.Focus() // sets focus=true right now
		var updateCmd tea.Cmd
		m.searchInput, updateCmd = m.searchInput.Update(msg)
		m.applyFilter()
		return m, tea.Batch(blinkCmd, updateCmd)
	}

	return m, nil
}

// ── Delete-confirm state ─────────────────────────────────────────────────────

func (m AppModel) handleDeleteConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case isConfirm(msg):
		var cmds []tea.Cmd
		for idx := range m.selected {
			if idx < len(m.choices) {
				name := m.choices[idx]
				path := m.folderPath(name)
				cmds = append(cmds, deleteFolder(path))
			}
		}
		m.selected = make(map[int]struct{})
		m.state = stateList
		m.statusMsg = "Folders deleted."
		return m, tea.Batch(cmds...)

	case isCancel(msg):
		m.selected = make(map[int]struct{})
		m.state = stateList
		return m, nil

	case isDelete(msg):
		if len(m.choices) == 0 {
			return m, nil
		}
		if _, ok := m.selected[m.cursor]; ok {
			delete(m.selected, m.cursor)
		} else {
			m.selected[m.cursor] = struct{}{}
		}
		if len(m.selected) == 0 {
			m.state = stateList
		}
		return m, nil

	case isUp(msg):
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil

	case isDown(msg):
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
		return m, nil
	}
	return m, nil
}

// ── Rename state ─────────────────────────────────────────────────────────────

func (m AppModel) handleRenameKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case isConfirm(msg):
		newName := strings.TrimSpace(m.renameInput.Value())
		if newName == "" {
			m.errMsg = "Name cannot be empty."
			return m, nil
		}
		oldName := m.choices[m.cursor]
		if oldName == newName {
			m.state = stateList
			m.renameInput.Blur()
			return m, nil
		}
		m.state = stateList
		m.renameInput.Blur()
		m.errMsg = ""
		return m, renameFolder(m.cfg.BasePath, oldName, newName)

	case isCancel(msg):
		m.state = stateList
		m.renameInput.Blur()
		m.errMsg = ""
		return m, nil

	default:
		var cmd tea.Cmd
		m.renameInput, cmd = m.renameInput.Update(msg)
		return m, cmd
	}
}

// ── Move state ───────────────────────────────────────────────────────────────

func (m AppModel) handleMoveKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case isConfirm(msg):
		dest := strings.TrimSpace(m.moveInput.Value())
		if dest == "" {
			m.errMsg = "Destination path cannot be empty."
			return m, nil
		}
		name := m.choices[m.cursor]
		src := m.folderPath(name)
		m.state = stateList
		m.moveInput.Blur()
		m.errMsg = ""
		return m, moveFolder(src, dest)

	case isCancel(msg):
		m.state = stateList
		m.moveInput.Blur()
		m.errMsg = ""
		return m, nil

	default:
		var cmd tea.Cmd
		m.moveInput, cmd = m.moveInput.Update(msg)
		return m, cmd
	}
}

// ── Create state ─────────────────────────────────────────────────────────────

func (m AppModel) handleCreateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case isCancel(msg):
		m.state = stateList
		m.createNameInput.Blur()
		m.createRepoInput.Blur()
		m.errMsg = ""
		return m, nil

	case isConfirm(msg):
		if !m.createFocusRepo {
			// Move focus to repo input on first Enter.
			m.createFocusRepo = true
			m.createNameInput.Blur()
			return m, m.createRepoInput.Focus()
		}
		// Second Enter = commit creation.
		name := strings.TrimSpace(m.createNameInput.Value())
		repo := strings.TrimSpace(m.createRepoInput.Value())
		if name == "" {
			m.errMsg = "Folder name cannot be empty."
			m.createFocusRepo = false
			m.createRepoInput.Blur()
			return m, m.createNameInput.Focus()
		}
		m.state = stateList
		m.createNameInput.Blur()
		m.createRepoInput.Blur()
		m.errMsg = ""
		return m, createFolder(m.cfg.BasePath, name, repo)

	case msg.Type == tea.KeyTab:
		m.createFocusRepo = !m.createFocusRepo
		if m.createFocusRepo {
			m.createNameInput.Blur()
			return m, m.createRepoInput.Focus()
		}
		m.createRepoInput.Blur()
		return m, m.createNameInput.Focus()

	default:
		var cmd tea.Cmd
		if m.createFocusRepo {
			m.createRepoInput, cmd = m.createRepoInput.Update(msg)
		} else {
			m.createNameInput, cmd = m.createNameInput.Update(msg)
		}
		return m, cmd
	}
}

// ── Helpers ──────────────────────────────────────────────────────────────────

// applyFilter rebuilds choices from allFolders based on search input value.
func (m *AppModel) applyFilter() {
	query := strings.ToLower(m.searchInput.Value())
	m.choices = nil
	for _, f := range m.allFolders {
		if query == "" || strings.Contains(strings.ToLower(f.Name), query) {
			m.choices = append(m.choices, f.Name)
		}
	}
	// Keep cursor in bounds.
	if m.cursor >= len(m.choices) {
		if len(m.choices) > 0 {
			m.cursor = len(m.choices) - 1
		} else {
			m.cursor = 0
		}
	}
	// Reset delete selections when filter changes (indices would be wrong).
	m.selected = make(map[int]struct{})
}

// folderPath returns the full path for a folder name.
func (m AppModel) folderPath(name string) string {
	for _, f := range m.allFolders {
		if f.Name == name {
			return f.Path
		}
	}
	return m.cfg.BasePath + "/" + name
}

// ── Commands ─────────────────────────────────────────────────────────────────

func loadFolders(basePath string) tea.Cmd {
	return func() tea.Msg {
		folders, err := fs.ListFolders(basePath)
		return foldersLoadedMsg{folders: folders, err: err}
	}
}

func deleteFolder(path string) tea.Cmd {
	return func() tea.Msg {
		return opDoneMsg{err: fs.Delete(path)}
	}
}

func renameFolder(basePath, oldName, newName string) tea.Cmd {
	return func() tea.Msg {
		return opDoneMsg{err: fs.Rename(basePath, oldName, newName)}
	}
}

func moveFolder(src, dest string) tea.Cmd {
	return func() tea.Msg {
		return opDoneMsg{err: fs.Move(src, dest)}
	}
}

func createFolder(basePath, name, repoURL string) tea.Cmd {
	return func() tea.Msg {
		var err error
		if repoURL != "" {
			err = fs.CreateFromGitHub(basePath, name, repoURL)
		} else {
			err = fs.CreateEmpty(basePath, name)
		}
		return opDoneMsg{err: err}
	}
}

// OpenConfigInEditor opens the config file in the configured editor.
func OpenConfigInEditor(cfg *config.Config) error {
	path, err := config.ConfigFilePath()
	if err != nil {
		return err
	}
	editor := cfg.Editor
	if editor == "" {
		editor = config.DefaultEditor
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = nil
	return cmd.Run()
}
