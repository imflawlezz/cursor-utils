package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/imflawlezz/cursor-utils/installer/internal/config"
	"github.com/imflawlezz/cursor-utils/installer/internal/installer"
)

var (
	highlight     = lipgloss.Color("#E8D5A3")
	titleStyle    = lipgloss.NewStyle().Bold(true)
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#B8B0A4"))
	labelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#D9D2C5"))
	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#2A2418")).
			Background(highlight)
	keyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#2A2418")).
			Background(highlight).
			Inline(true)
	dangerSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#3A1010")).
				Background(lipgloss.Color("#FFB8AE"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#C8E6B0"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#F0D48A"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F0A8A0"))
	authorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#E1C16E"))
	keyTextStyle = lipgloss.NewStyle().Bold(true).Foreground(highlight)
	statusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#C9C2B4"))
	bodyStyle    = lipgloss.NewStyle().Padding(1, 2)
)

func headerLine(version string) string {
	cfg := config.Default()
	if version == "" {
		version = cfg.Version
	}
	author := authorStyle.Render(cfg.GitHubOwner)
	author = osc8Link(cfg.AuthorURL(), author)
	return titleStyle.Render("cursor-utils installer") + " " +
		labelStyle.Render(version) + " by " + author
}

func osc8Link(url, text string) string {
	return "\x1b]8;;" + url + "\x1b\\" + text + "\x1b]8;;\x1b\\"
}

func checkMark(on, mixed bool) string {
	switch {
	case on:
		return "[x]"
	case mixed:
		return "[~]"
	default:
		return "[ ]"
	}
}

func fileInstallStatus(plan *installer.Plan, rel string) string {
	if plan == nil {
		return "not installed"
	}
	for _, op := range plan.Ops {
		if op.RelPath != rel {
			continue
		}
		switch op.Kind {
		case installer.OpAdd:
			return "not installed"
		case installer.OpUnchanged:
			return "installed"
		case installer.OpUpdate:
			return "installed · outdated"
		case installer.OpConflict:
			return "exists · not owned"
		case installer.OpModifiedUpdate:
			return "installed · modified"
		case installer.OpRemove, installer.OpModifiedRemove:
			return "installed · gone in this version"
		}
	}
	return "not installed"
}

func groupInstallStatus(plan *installer.Plan, id string, fileCount int) string {
	if plan == nil || fileCount == 0 {
		return "not installed"
	}
	var installed, outdated, modified, conflicts, missing int
	for _, op := range plan.Ops {
		if op.Component != id {
			continue
		}
		switch op.Kind {
		case installer.OpUnchanged:
			installed++
		case installer.OpUpdate:
			installed++
			outdated++
		case installer.OpModifiedUpdate:
			installed++
			modified++
		case installer.OpConflict:
			conflicts++
		case installer.OpAdd:
			missing++
		case installer.OpRemove, installer.OpModifiedRemove:
			installed++
		}
	}
	if installed == 0 && conflicts == 0 {
		return "not installed"
	}
	var parts []string
	if installed > 0 {
		line := fmt.Sprintf("%d/%d installed", installed, fileCount)
		if ver := componentInstalledVer(plan, id); ver != "" {
			line += " (" + ver + ")"
		}
		parts = append(parts, line)
	} else {
		parts = append(parts, "not installed")
	}
	if outdated > 0 {
		parts = append(parts, fmt.Sprintf("%d outdated", outdated))
	}
	if modified > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", modified))
	}
	if conflicts > 0 {
		parts = append(parts, fmt.Sprintf("%d not owned", conflicts))
	}
	if missing > 0 && installed > 0 {
		parts = append(parts, fmt.Sprintf("%d missing", missing))
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += " · " + parts[i]
	}
	return out
}

func componentInstalledVer(plan *installer.Plan, id string) string {
	if plan == nil {
		return ""
	}
	for _, sum := range plan.ComponentSummaries() {
		if sum.ID == id {
			return sum.InstalledVer
		}
	}
	return ""
}

func helpBar(items [][2]string) string {
	sep := dimStyle.Render(" / ")
	parts := make([]string, 0, len(items))
	for _, it := range items {
		var keys []string
		for _, k := range strings.Fields(it[0]) {
			keys = append(keys, keyStyle.Render(" "+k+" "))
		}
		parts = append(parts, strings.Join(keys, " ")+" "+dimStyle.Render(it[1]))
	}
	return strings.Join(parts, sep)
}
