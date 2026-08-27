package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"

	"github.com/imflawlezz/cursor-utils/installer/internal/config"
	"github.com/imflawlezz/cursor-utils/installer/internal/github"
	"github.com/imflawlezz/cursor-utils/installer/internal/installer"
	"github.com/imflawlezz/cursor-utils/installer/internal/platform"
	"github.com/imflawlezz/cursor-utils/installer/internal/tui"
)

const help = `cursor-utils installer — install versioned Cursor content from GitHub tags.

Usage:
  cursor-utils           Start the interactive installer
  cursor-utils version   Print installer version
  cursor-utils help      Show this help

The installer manages only files it has recorded in its manifest.
It will not overwrite or delete your own Cursor files without asking.
`

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "--version", "-v":
			fmt.Printf("cursor-utils installer %s\n", config.Version)
			return 0
		case "help", "--help", "-h":
			fmt.Print(help)
			return 0
		default:
			fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
			fmt.Print(help)
			return 2
		}
	}

	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Fprintln(os.Stderr, "cursor-utils requires an interactive terminal.")
		fmt.Fprintln(os.Stderr, "Run it directly from a terminal.")
		return 1
	}

	plat, err := platform.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cursor-utils: %v\n", err)
		return 1
	}

	cfg := config.Default()
	eng := installer.New(cfg, github.New(cfg), plat)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	p := tea.NewProgram(
		tui.New(eng, plat),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithContext(ctx),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cursor-utils: %v\n", err)
		return 1
	}
	return 0
}
