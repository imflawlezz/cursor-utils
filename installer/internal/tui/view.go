package tui

import (
	"fmt"
	"path"
	"strings"

	"github.com/imflawlezz/cursor-utils/installer/internal/installer"
)

func (m Model) View() string {
	if m.quitting && m.screen != screenProgress {
		return ""
	}
	var body string
	switch m.screen {
	case screenUnsupported:
		body = m.viewUnsupported()
	case screenSystem:
		body = m.viewSystem()
	case screenPathInput:
		body = m.viewPath()
	case screenLoading:
		body = m.viewLoading()
	case screenVersion:
		body = m.viewVersion()
	case screenVersionList:
		body = m.viewVersionList()
	case screenManage:
		body = m.viewManage()
	case screenConfirmRemove:
		body = m.viewConfirmRemove()
	case screenDecision:
		body = m.viewDecision()
	case screenProgress:
		body = m.viewProgress()
	case screenResult:
		body = m.viewResult()
	case screenError:
		body = m.viewError()
	case screenKeys:
		from := m.keysFrom
		if from == screenKeys {
			from = screenManage
		}
		body = viewKeyHelp(from)
	}
	help := helpBar(footerFor(m.screen))
	w := m.width
	if w <= 0 {
		w = 80
	}
	header := headerLine(m.engine.Config().Version)
	content := header + "\n\n" + body + "\n\n" + help
	return bodyStyle.Width(w).Render(content)
}

func (m Model) viewUnsupported() string {
	return errorStyle.Render(m.plat.DisplayName+" is not yet supported.") +
		"\n\nThis installer currently supports macOS and Linux."
}

func (m Model) viewSystem() string {
	osLine := "  OS: " + m.plat.DisplayName
	dir := m.cursorRoot
	if dir == "" {
		dir = m.plat.DefaultDir
	}
	if dir == "" {
		dir = "could not be detected"
	}
	dirLine := "  Cursor directory: " + dir

	var b strings.Builder
	b.WriteString(labelStyle.Render("System") + "\n")
	b.WriteString(osLine + "\n")
	b.WriteString(dirLine + "\n\n")

	if !m.plat.DefaultOK {
		b.WriteString("The default Cursor configuration directory could not be\n")
		b.WriteString("determined. Choose the directory manually.\n\n")
		b.WriteString(options([]string{"Choose location", "Quit"}, m.sysCursor))
		return b.String()
	}
	b.WriteString("Is this correct?\n\n")
	b.WriteString(options([]string{"Continue", "Choose another location"}, m.sysCursor))
	return b.String()
}

func (m Model) viewPath() string {
	var b strings.Builder
	b.WriteString(labelStyle.Render("Cursor directory") + "\n\n")
	b.WriteString("Enter the Cursor configuration directory.\n")
	b.WriteString(dimStyle.Render("This is usually ~/.cursor") + "\n\n")
	b.WriteString(m.pathInput.View() + "\n")
	if m.pathErr != "" {
		b.WriteString("\n" + errorStyle.Render(m.pathErr) + "\n")
	}
	return b.String()
}

func (m Model) viewLoading() string {
	note := m.loadNote
	if note == "" {
		note = "Working…"
	}
	return note
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
	items := []string{latest, "Specific version"}
	b.WriteString(radio(items, m.versionCursor))
	return b.String()
}

func (m Model) viewVersionList() string {
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
	for i := m.listOffset; i < end; i++ {
		line := "    " + m.tags[i]
		if i == m.listCursor {
			line = selectedStyle.Render("  > " + m.tags[i])
		}
		b.WriteString(line + "\n")
	}
	return b.String()
}

func (m Model) viewManage() string {
	var b strings.Builder
	ver := m.selectedTag
	if ver == "" && m.plan != nil {
		ver = m.plan.ContentVersion
	}
	b.WriteString("Content: " + ver + "\n")
	if m.cursorRoot != "" {
		b.WriteString(dimStyle.Render("Cursor directory: "+m.cursorRoot) + "\n")
	}
	b.WriteString("\n" + labelStyle.Render("Components") + "\n\n")
	if len(m.groups) == 0 {
		b.WriteString(dimStyle.Render("  No components in this version.") + "\n\n")
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
			mark := checkMark(g.allSelected(), g.anySelected() && !g.allSelected())
			line := "  " + arrow + " " + mark + " " + g.displayName()
			if cursor {
				line = selectedStyle.Render(line)
			}
			b.WriteString(line + "\n")
			status := fmt.Sprintf("        %d files  ·  %s", len(g.Files), groupInstallStatus(m.plan, g.ID, len(g.Files)))
			b.WriteString(statusStyle.Render(status) + "\n")
		case rowFile:
			g := m.groups[row.GroupIdx]
			f := g.Files[row.FileIdx]
			mark := checkMark(f.Selected, false)
			status := fileInstallStatus(m.plan, f.RelPath)
			line := "      " + mark + " " + fileLabel(f.RelPath)
			if cursor {
				line = selectedStyle.Render(line + "  " + status)
			} else {
				b.WriteString(line)
				b.WriteString(statusStyle.Render("  "+status) + "\n")
				continue
			}
			b.WriteString(line + "\n")
		case rowLabel:
			b.WriteString("\n" + labelStyle.Render(row.Action) + "\n")
		case rowSpacer:
			b.WriteString("\n")
		case rowAction:
			line := "    " + row.Action
			if cursor {
				line = selectedStyle.Render("  > " + row.Action)
			}
			b.WriteString(line + "\n")
		}
	}
	if m.manageErr != "" {
		b.WriteString("\n" + errorStyle.Render(m.manageErr) + "\n")
	}
	return b.String()
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
	b.WriteString(optionsDanger([]string{"Remove", "Cancel"}, m.removeCursor, 0))
	return b.String()
}

func formatRemoveList(paths []string, max int) string {
	if max < 1 {
		max = 5
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
		if shown >= max {
			break
		}
		b.WriteString("  " + g.dir + "/\n")
		for _, f := range g.files {
			if shown >= max {
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
	b.WriteString(options(labels, m.decCursor))
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
	b.WriteString(title + "\n\n")
	if len(m.progress) == 0 {
		b.WriteString(dimStyle.Render("  Working…") + "\n")
		return b.String()
	}
	for _, line := range m.progress {
		mark := "●"
		style := selectedStyle
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
			style = selectedStyle
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
		b.WriteString(fmt.Sprintf("  Files removed: %d\n", m.filesDone))
	default:
		b.WriteString(successStyle.Render("Installation complete") + "\n\n")
		b.WriteString("  Version: " + m.resultVer + "\n")
		if len(m.resultComps) > 0 {
			b.WriteString("  Components: " + strings.Join(m.resultComps, ", ") + "\n")
		}
		b.WriteString(fmt.Sprintf("  Files installed: %d\n", m.filesDone))
	}
	b.WriteString("\nPress Enter to continue")
	return b.String()
}

func (m Model) viewError() string {
	var b strings.Builder
	b.WriteString(errorStyle.Render("Something went wrong") + "\n\n")
	b.WriteString(friendly(m.err) + "\n\n")
	b.WriteString(options([]string{"Retry", "Quit"}, m.decCursor))
	return b.String()
}

func options(items []string, selected int) string {
	return optionsDanger(items, selected, -1)
}

func optionsDanger(items []string, selected, danger int) string {
	var b strings.Builder
	for i, item := range items {
		if i == selected {
			style := selectedStyle
			if i == danger {
				style = dangerSelectedStyle
			}
			b.WriteString(style.Render("  > "+item) + "\n")
		} else {
			b.WriteString("    " + item + "\n")
		}
	}
	return b.String()
}

func radio(items []string, selected int) string {
	var b strings.Builder
	for i, item := range items {
		mark := "○"
		if i == selected {
			mark = "●"
		}
		line := "  " + mark + " " + item
		if i == selected {
			line = selectedStyle.Render(line)
		}
		b.WriteString(line + "\n")
	}
	return b.String()
}
