package tui

import (
	"fmt"
	"path"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/imflawlezz/cursor-utils/installer/internal/installer"
)

func (m Model) View() string {
	if m.quitting && m.screen != screenProgress {
		return ""
	}
	return wrapFrame(m.screenBody(), m.helpLine(), m.version(), m.frameWidth(), m.height)
}

func (m Model) frameWidth() int {
	w := m.width
	if w <= 0 {
		return 80
	}
	return w
}

func (m Model) contentWidth() int {
	w := m.frameWidth() - 6
	if w < 8 {
		return 8
	}
	return w
}

func (m Model) helpLine() string {
	return wrapFitted(helpBar(footerFor(m.screen)), max(8, m.frameWidth()-6))
}

func (m Model) version() string {
	if m.engine != nil {
		return m.engine.Config().Version
	}
	return "dev"
}

func (m Model) screenBody() string {
	switch m.screen {
	case screenUnsupported:
		return m.viewUnsupported()
	case screenSystem:
		return m.viewSystem()
	case screenPathInput:
		return m.viewPath()
	case screenLoading:
		return m.viewLoading()
	case screenVersion:
		return m.viewVersion()
	case screenVersionList:
		return m.viewVersionList()
	case screenManage:
		return m.viewManage()
	case screenConfirmRemove:
		return m.viewConfirmRemove()
	case screenDecision:
		return m.viewDecision()
	case screenProgress:
		return m.viewProgress()
	case screenResult:
		return m.viewResult()
	case screenError:
		return m.viewError()
	case screenKeys:
		from := m.keysFrom
		if from == screenKeys {
			from = screenManage
		}
		return viewKeyHelp(from)
	default:
		return ""
	}
}

func (m Model) viewUnsupported() string {
	return errorStyle.Render(m.plat.DisplayName+" is not yet supported.") +
		"\n\nThis installer currently supports macOS, Linux, and Windows."
}

func (m Model) viewSystem() string {
	osLine := "OS: " + m.plat.DisplayName
	dir := m.cursorRoot
	if dir == "" {
		dir = m.plat.DefaultDir
	}
	if dir == "" {
		dir = "could not be detected"
	}

	cw := m.contentWidth()
	var b strings.Builder
	b.WriteString(labelStyle.Render("System") + "\n")
	b.WriteString(osLine + "\n")
	b.WriteString("Cursor directory: " + dir + "\n\n")

	if !m.plat.DefaultOK {
		for _, line := range wrapWords("The default Cursor configuration directory could not be determined. Choose the directory manually.", cw) {
			b.WriteString(line + "\n")
		}
		b.WriteString("\n")
		b.WriteString(m.options([]string{"Choose location", "Quit"}, m.sysCursor))
		return b.String()
	}
	b.WriteString("Is this correct?\n\n")
	b.WriteString(m.options([]string{"Continue", "Choose another location"}, m.sysCursor))
	return b.String()
}

func (m Model) viewPath() string {
	var b strings.Builder
	b.WriteString(labelStyle.Render("Cursor directory") + "\n\n")
	b.WriteString("Enter the Cursor configuration directory.\n")
	b.WriteString(dimStyle.Render("This is usually ~/.cursor") + "\n\n")
	b.WriteString(m.pathInput.View() + "\n")
	if m.pathErr != "" {
		b.WriteString("\n")
		for _, line := range wrapWords(m.pathErr, m.contentWidth()) {
			b.WriteString(errorStyle.Render(line) + "\n")
		}
	}
	return b.String()
}

func (m Model) viewLoading() string {
	note := m.loadNote
	if note == "" {
		note = "Working…"
	}
	return statusStyle.Render(note)
}

func (m Model) viewVersion() string {
	var b strings.Builder
	b.WriteString(labelStyle.Render("Content version") + "\n\n")
	latest := "Latest"
	if m.latest != "" {
		latest = "Latest (" + m.latest + ")"
	} else {
		latest = "Latest (none found)"
	}
	b.WriteString(m.radio([]string{latest, "Specific version"}, m.versionCursor))
	return b.String()
}

func (m Model) viewVersionList() string {
	cw := m.contentWidth()
	var b strings.Builder
	b.WriteString(labelStyle.Render("Select version") + "\n\n")
	if len(m.tags) == 0 {
		b.WriteString("No version tags were found in the repository.")
		return b.String()
	}
	visible := 10
	end := m.listOffset + visible
	if end > len(m.tags) {
		end = len(m.tags)
	}
	numW := len(fmt.Sprintf("%d", max(1, len(m.tags))))
	for i := m.listOffset; i < end; i++ {
		num := fmt.Sprintf("%*d. ", numW, i+1)
		line := num + m.tags[i]
		if i == m.listCursor {
			b.WriteString(fillSelected(line, cw, false) + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}
	return b.String()
}

func (m Model) viewManage() string {
	cw := m.contentWidth()
	var b strings.Builder
	ver := m.selectedTag
	if ver == "" && m.plan != nil {
		ver = m.plan.ContentVersion
	}
	b.WriteString("Content: " + ver + "\n")
	if m.cursorRoot != "" {
		b.WriteString("Cursor directory: " + compactPath(m.cursorRoot, m.plat.Home) + "\n")
	}
	b.WriteString("\n" + labelStyle.Render("Components") + "\n")
	if len(m.groups) == 0 {
		b.WriteString(dimStyle.Render("(empty — no components in this version)") + "\n")
	}

	rows := m.manageRows()
	for i, row := range rows {
		cursor := i == m.manageCursor
		switch row.Kind {
		case rowGroup:
			g := m.groups[row.GroupIdx]
			arrow := "▸"
			if g.Expanded {
				arrow = "▾"
			}
			mark := selectMark(g.allSelected(), g.anySelected() && !g.allSelected())
			left := "  " + arrow + " " + mark + " " + g.displayName()
			b.WriteString(m.suffixLine(left, "", cursor, false, cw) + "\n")
			status := fmt.Sprintf("%d files  ·  %s", len(g.Files), groupInstallStatus(m.plan, g.ID, len(g.Files)))
			b.WriteString(statusStyle.Render("      "+status) + "\n")
		case rowFile:
			g := m.groups[row.GroupIdx]
			f := g.Files[row.FileIdx]
			status := fileInstallStatus(m.plan, f.RelPath)
			mark := selectMark(f.Selected, false)
			left := "      " + mark + " " + fileLabel(f.RelPath)
			b.WriteString(m.suffixLine(left, status, cursor, false, cw) + "\n")
		case rowLabel:
			b.WriteString("\n" + labelStyle.Render(row.Action) + "\n")
		case rowSpacer:
			b.WriteString("\n")
		case rowAction:
			b.WriteString(optionLine(row.Action, cursor, row.Action == "Remove selected", cw) + "\n")
		}
	}
	if m.manageErr != "" {
		b.WriteString("\n")
		for _, line := range wrapWords(m.manageErr, cw) {
			b.WriteString(errorStyle.Render(line) + "\n")
		}
	}
	return b.String()
}

func (m Model) suffixLine(left, right string, cursor, destructive bool, width int) string {
	if width < 8 {
		width = 8
	}
	if right == "" {
		if cursor {
			return fillSelected(left, width, destructive)
		}
		return left
	}
	gap := "  "
	leftW := width - lipgloss.Width(gap) - lipgloss.Width(right)
	if leftW < 8 {
		right = truncateEnd(right, max(6, width/3))
		leftW = width - lipgloss.Width(gap) - lipgloss.Width(right)
		if leftW < 4 {
			leftW = 4
		}
	}
	line := padRight(truncateEnd(left, leftW), leftW) + gap + right
	if cursor {
		return fillSelected(line, width, destructive)
	}
	return padRight(truncateEnd(left, leftW), leftW) + gap + dimStyle.Render(right)
}

func (m Model) viewConfirmRemove() string {
	var paths []string
	if m.plan != nil {
		for _, op := range m.plan.Ops {
			if op.RelPath == "" {
				continue
			}
			paths = append(paths, op.RelPath)
		}
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Remove %d cursor-utils files?\n\n", len(paths)))
	b.WriteString(formatRemoveList(paths, 5))
	b.WriteString("\n")
	b.WriteString(m.optionsDanger([]string{"Remove", "Cancel"}, m.removeCursor, 0))
	return b.String()
}

func formatRemoveList(paths []string, maxN int) string {
	if maxN < 1 {
		maxN = 5
	}
	type grp struct {
		dir   string
		files []string
	}
	var groups []grp
	idx := map[string]int{}
	for _, p := range paths {
		dir := path.Dir(p)
		base := path.Base(p)
		if i, ok := idx[dir]; ok {
			groups[i].files = append(groups[i].files, base)
			continue
		}
		idx[dir] = len(groups)
		groups = append(groups, grp{dir: dir, files: []string{base}})
	}
	shown := 0
	var b strings.Builder
	for _, g := range groups {
		if shown >= maxN {
			break
		}
		b.WriteString("  " + g.dir + "/\n")
		for _, f := range g.files {
			if shown >= maxN {
				break
			}
			b.WriteString("    " + f + "\n")
			shown++
		}
	}
	if extra := len(paths) - shown; extra > 0 {
		b.WriteString(fmt.Sprintf("  …and %d more\n", extra))
	}
	return b.String()
}

func (m Model) viewDecision() string {
	if m.pendIndex >= len(m.pending) {
		return ""
	}
	op := m.pending[m.pendIndex]
	var b strings.Builder
	switch op.Kind {
	case installer.OpConflict:
		b.WriteString(op.RelPath + " already exists.\n\n")
		b.WriteString("This file is not owned by cursor-utils.\n\n")
	case installer.OpModifiedUpdate:
		b.WriteString(op.RelPath + " was modified after installation.\n\n")
	case installer.OpModifiedRemove:
		b.WriteString(op.RelPath + " was modified after installation.\n\n")
		b.WriteString("It is no longer in the selected version (or is being removed).\n\n")
	}
	labels := make([]string, 0)
	for _, o := range m.decisionOptions() {
		labels = append(labels, o.label)
	}
	b.WriteString(m.options(labels, m.decCursor))
	return b.String()
}

func (m Model) viewProgress() string {
	title := "Installing"
	if m.result == resultRemove {
		title = "Removing"
	}
	if m.resultVer != "" && m.result != resultRemove {
		title += " " + m.resultVer
	}
	var b strings.Builder
	b.WriteString(labelStyle.Render(title) + "\n\n")
	if len(m.progress) == 0 {
		b.WriteString(dimStyle.Render("Working…") + "\n")
		return b.String()
	}
	for _, line := range m.progress {
		mark := "●"
		style := noticeStyle
		switch line.Status {
		case installer.EventDone:
			mark = "✓"
			style = successStyle
		case installer.EventSkipped:
			mark = "–"
			style = dimStyle
		case installer.EventError:
			mark = "✕"
			style = errorStyle
		case installer.EventStarted:
			mark = "●"
			style = noticeStyle
		}
		b.WriteString("  " + style.Render(mark+" "+line.Path) + "\n")
	}
	return b.String()
}

func (m Model) viewResult() string {
	var b strings.Builder
	switch m.result {
	case resultCancelled:
		b.WriteString(warnStyle.Render("Cancelled") + "\n\n")
		b.WriteString("Existing Cursor files were left unchanged.")
	case resultFailed:
		b.WriteString(errorStyle.Render("Installation failed") + "\n\n")
		b.WriteString(friendly(m.resultErr) + "\n\n")
		b.WriteString("Existing Cursor files were left unchanged.")
	case resultRemove:
		b.WriteString(successStyle.Render("Removal complete") + "\n\n")
		b.WriteString(fmt.Sprintf("Files removed: %d\n", m.filesDone))
	default:
		b.WriteString(successStyle.Render("Installation complete") + "\n\n")
		b.WriteString("Version: " + m.resultVer + "\n")
		if len(m.resultComps) > 0 {
			b.WriteString("Components: " + strings.Join(m.resultComps, ", ") + "\n")
		}
		b.WriteString(fmt.Sprintf("Files installed: %d\n", m.filesDone))
	}
	return b.String()
}

func (m Model) viewError() string {
	var b strings.Builder
	b.WriteString(errorStyle.Render("Something went wrong") + "\n\n")
	b.WriteString(friendly(m.err) + "\n\n")
	b.WriteString(m.options([]string{"Retry", "Quit"}, m.decCursor))
	return b.String()
}

func (m Model) options(items []string, selected int) string {
	return m.optionsDanger(items, selected, -1)
}

func (m Model) optionsDanger(items []string, selected, danger int) string {
	cw := m.contentWidth()
	var b strings.Builder
	for i, item := range items {
		b.WriteString(optionLine(item, i == selected, i == danger, cw) + "\n")
	}
	return b.String()
}

func (m Model) radio(items []string, selected int) string {
	cw := m.contentWidth()
	var b strings.Builder
	for i, item := range items {
		mark := "○"
		if i == selected {
			mark = "●"
		}
		line := "  " + mark + " " + item
		if i == selected {
			line = fillSelected(line, cw, false)
		}
		b.WriteString(line + "\n")
	}
	return b.String()
}

func compactPath(p, home string) string {
	if p == "" {
		return p
	}
	if home != "" && (strings.HasPrefix(p, home+"/") || strings.HasPrefix(p, home+"\\")) {
		return "~" + p[len(home):]
	}
	return p
}
