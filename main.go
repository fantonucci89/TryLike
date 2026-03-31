package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/fantonucci89/TryLike/internal/config"
	"github.com/fantonucci89/TryLike/internal/tui"
)

func main() {
	// Optional CLI arg: folder name to pre-fill search.
	var initialFilter string
	if len(os.Args) > 1 {
		initialFilter = os.Args[1]
	}

	// Load (or create) config.
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "trylike: failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Build and run the TUI.
	m := tui.New(cfg, initialFilter)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "trylike: %v\n", err)
		os.Exit(1)
	}

	if app, ok := finalModel.(tui.AppModel); ok && app.SelectedPath != "" {
		// fmt.Println(app.SelectedPath)
		spawnShell(app.SelectedPath)
	}
}

func spawnShell(dir string) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	if err := os.Chdir(dir); err != nil {
		os.Exit(1)
	}

	err := syscall.Exec(shell, []string{shell}, os.Environ())
	if err != nil {
		cmd := exec.Command(shell)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
