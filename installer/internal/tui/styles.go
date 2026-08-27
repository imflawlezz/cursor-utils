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
	keyTextStyle = lipgloss.NewStyle().Bold(true).Foreground(highlight)
	statusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#C9C2B4"))
	noticeStyle  = lipgloss.NewStyle().Foreground(highlight)
	borderStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#B8B0A4"))
)

func displayVersion(v string) string {
	v = strings.TrimSpace(v)
	if v == "" || v == "dev" {
		return "dev"
	}
	if !strings.HasPrefix(v, "v") {
		return "v" + v
	}
	return v
}

func headerTitle(version string) string {
	cfg := config.Default()
	return titleStyle.Render("cursor-utils installer") + " " +
		labelStyle.Render(displayVersion(version)) + " by " +
		styledHyperlink(cfg.AuthorURL(), cfg.GitHubOwner)
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

func fillSelected(line string, width int, destructive bool) string {
	if width < 1 {
		width = 1
	}
	line = padRight(truncateEnd(line, width), width)
	style := selectedStyle
	if destructive {
		style = dangerSelectedStyle
	}
	return style.Inline(true).MaxHeight(1).Render(line)
}

func optionLine(label string, selected, destructive bool, width int) string {
	if selected {
		return fillSelected("  > "+label, width, destructive)
	}
	return "    " + label
}

func truncateEnd(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= maxWidth {
		return s
	}
	if strings.Contains(s, "\x1b]8;") {
		return s
	}
	ellipsis := "…"
	budget := maxWidth - lipgloss.Width(ellipsis)
	if budget <= 0 {
		return ellipsis
	}
	var b strings.Builder
	w := 0
	for _, r := range s {
		rw := lipgloss.Width(string(r))
		if w+rw > budget {
			break
		}
		b.WriteRune(r)
		w += rw
	}
	return b.String() + ellipsis
}

func padRight(s string, width int) string {
	if width <= 0 {
		return ""
	}
	s = truncateEnd(s, width)
	if pad := width - lipgloss.Width(s); pad > 0 {
		s += strings.Repeat(" ", pad)
	}
	return s
}

func hyperlink(url, label string) string {
	if url == "" {
		return label
	}
	return "\x1b]8;;" + url + "\x1b\\" + label + "\x1b]8;;\x1b\\"
}

func styledHyperlink(url, label string) string {
	if url == "" {
		return label
	}
	return "\x1b[4;38;2;225;193;110m" + hyperlink(url, label) + "\x1b[0m"
}

func wrapFitted(s string, width int) string {
	if width < 8 || lipgloss.Width(s) <= width {
		return s
	}
	sep := dimStyle.Render(" / ")
	parts := strings.Split(s, sep)
	if len(parts) == 1 {
		return s
	}
	var lines []string
	var cur string
	for _, p := range parts {
		next := p
		if cur != "" {
			next = cur + sep + p
		}
		if lipgloss.Width(next) <= width {
			cur = next
			continue
		}
		if cur != "" {
			lines = append(lines, cur)
		}
		cur = p
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return strings.Join(lines, "\n")
}

func wrapWords(s string, width int) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if width < 8 {
		width = 8
	}
	words := strings.Fields(s)
	var lines []string
	var cur string
	for _, w := range words {
		if cur == "" {
			cur = w
			continue
		}
		next := cur + " " + w
		if lipgloss.Width(next) <= width {
			cur = next
			continue
		}
		lines = append(lines, cur)
		cur = w
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}

func wrapFrame(body, footer, version string, width, height int) string {
	if width < 2 {
		width = 80
	}
	inner := width - 2
	title := " " + headerTitle(version) + " "
	remain := inner - 1 - lipgloss.Width(title)
	if remain < 0 {
		remain = 0
	}
	side := borderStyle.Render("│")
	top := borderStyle.Render("┌─") + title + borderStyle.Render(strings.Repeat("─", remain)+"┐")
	bot := borderStyle.Render("└" + strings.Repeat("─", inner) + "┘")
	blank := side + strings.Repeat(" ", inner) + side
	bodyLines := splitViewLines(body)
	footerLines := splitViewLines(footer)

	minPad := 0
	if len(footerLines) > 0 {
		minPad = 1
	}
	pad := minPad
	if height > 0 {
		used := 4 + len(bodyLines) + len(footerLines)
		if extra := height - used; extra > pad {
			pad = extra
		}
	}

	var b strings.Builder
	b.WriteString(top + "\n")
	b.WriteString(blank + "\n")
	for _, line := range bodyLines {
		b.WriteString(side)
		b.WriteString(padRight("  "+line, inner))
		b.WriteString(side + "\n")
	}
	for i := 0; i < pad; i++ {
		b.WriteString(blank + "\n")
	}
	for _, line := range footerLines {
		b.WriteString(side)
		b.WriteString(padRight("  "+line, inner))
		b.WriteString(side + "\n")
	}
	b.WriteString(blank + "\n")
	b.WriteString(bot)
	return b.String()
}

func splitViewLines(s string) []string {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func selectMark(on, mixed bool) string {
	switch {
	case on:
		return "✓"
	case mixed:
		return "~"
	default:
		return " "
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
