package tui

import (
	"os/exec"
	"runtime"

	"github.com/charmbracelet/lipgloss"

	"github.com/imflawlezz/cursor-utils/installer/internal/config"
)

func authorNameCells(version string) (x0, x1 int) {
	cfg := config.Default()
	prefix := "┌─ cursor-utils installer " + displayVersion(version) + " by "
	x0 = lipgloss.Width(prefix)
	x1 = x0 + lipgloss.Width(cfg.GitHubOwner)
	return
}

func (m Model) linkAt(x, y int) string {
	x0, x1 := authorNameCells(m.version())
	if y == 0 && x >= x0 && x < x1 {
		return config.Default().AuthorURL()
	}
	return ""
}

func openURL(raw string) error {
	if raw == "" {
		return nil
	}
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", raw).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", raw).Start()
	default:
		return exec.Command("xdg-open", raw).Start()
	}
}
