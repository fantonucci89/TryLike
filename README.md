# TryLike

A keyboard-driven TUI for managing a project workspace directory.

Navigate your folders with arrow keys, jump to any folder by typing to filter,
and open it in your shell with `Enter`. Create empty folders or clone a GitHub
repo directly into the workspace, rename, move, or bulk-delete with interactive
confirmation — all without leaving the terminal.

Built with [Go](https://go.dev), [Bubble Tea](https://github.com/charmbracelet/bubbletea),
[Bubbles](https://github.com/charmbracelet/bubbles), and
[Lip Gloss](https://github.com/charmbracelet/lipgloss).

---

## Features

- **Browse** all folders in your configured workspace path
- **Live search** — start typing any letter to filter the list instantly
- **Open** a folder in your shell by pressing `Enter` (requires a shell wrapper, see below)
- **Create** an empty folder or clone a GitHub repo into it
- **Rename** a folder with an interactive dialog
- **Move** a folder to any path within `$HOME`
- **Delete** one or more folders with a confirmation step
- Pinned bottom bar with context-sensitive key binding hints

---

## Installation

### From source

Requires Go 1.21+.

```bash
git clone https://github.com/fantonucci89/TryLike
cd TryLike
make install          # builds and copies to ~/.local/bin/tryl
```

Make sure `~/.local/bin` is in your `$PATH`.

### Pre-built binaries

Download a binary for your platform from the [Releases](https://github.com/fantonucci89/TryLike/releases) page,
make it executable, and place it somewhere in your `$PATH`:

```bash
chmod +x tryl-linux-amd64
mv tryl-linux-amd64 ~/.local/bin/tryl
```

Available targets:

| File | Platform |
|---|---|
| `tryl-linux-amd64` | Linux x86-64 |
| `tryl-linux-arm64` | Linux ARM64 |
| `tryl-darwin-amd64` | macOS Intel |
| `tryl-darwin-arm64` | macOS Apple Silicon |

---

## Usage

```bash
tryl                  # open the TUI
tryl my-project       # open with the search pre-filled with "my-project"
```

### Key bindings

| Key | Action |
|---|---|
| `↑` / `↓` | Move cursor up / down |
| Any letter | Focus the search bar and start filtering |
| `Enter` | Open the folder under the cursor (quit + cd) |
| `Ctrl+N` | Create a new folder |
| `Ctrl+R` | Rename the folder under the cursor |
| `Ctrl+M` | Move the folder under the cursor |
| `Ctrl+D` | Mark the folder under the cursor for deletion |
| `Esc` | Cancel the current action |
| `Ctrl+C` | Quit without opening a folder |

#### Delete flow

Press `Ctrl+D` to mark one or more folders. The bottom bar switches to:

| Key | Action |
|---|---|
| `Enter` | Confirm — permanently delete all marked folders |
| `Esc` | Cancel — unmark everything |
| `Ctrl+D` | Toggle the mark on the folder under the cursor |

#### Create flow

Press `Ctrl+N` to open the create dialog. Fill in a folder name, then either:

- Press `Enter` (or `Tab`) to skip to the repo URL field and confirm to **clone a GitHub repo**, or
- Leave the repo URL empty and press `Enter` to create an **empty folder**

---

## Configuration

On first run, `tryl` creates a config file at:

```
~/.config/tryl/config.toml
```

Default contents:

```toml
base_path = "/home/<you>/Work"
editor = "vim"
```

| Key | Description | Default |
|---|---|---|
| `base_path` | Absolute path to the workspace directory `tryl` manages | `~/Work` |
| `editor` | Editor used to open the config file | `vim` |

`~` is expanded automatically in `base_path`.

---

## Building

```bash
make build    # compile for the current host into ./tryl
make dist     # cross-compile all four release binaries into dist/
make clean    # remove dist/ and ./tryl
make install  # build + copy to ~/.local/bin/tryl
make help     # list all targets
```

---

## Project structure

```
tryl/
├── main.go                       # Entry point
├── Makefile
├── go.mod / go.sum
└── internal/
    ├── config/config.go          # Config file (~/.config/tryl/config.toml)
    ├── fs/fs.go                  # Filesystem operations
    └── tui/
        ├── model.go              # Bubble Tea model + state machine
        ├── keys.go               # Key binding helpers
        ├── styles/styles.go      # Lip Gloss styles
        └── views/
            ├── list.go           # Main list view
            ├── rename.go         # Rename dialog
            ├── move.go           # Move dialog
            └── create.go         # Create folder dialog
```

---

## License

MIT
