package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/imflawlezz/cursor-utils/installer/internal/installer"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case useSavedRootMsg:
		return m.continueFromSystem()
	case tea.KeyMsg:
		if key(msg) == "ctrl+c" {
			if m.applyCancel != nil {
				m.applyCancel()
				m.quitting = true
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit
		}
		if key(msg) == "?" && m.screen != screenPathInput && m.screen != screenProgress && m.screen != screenKeys {
			return m.openKeys(), nil
		}
	}

	switch m.screen {
	case screenUnsupported:
		return m.updateSimpleQuit(msg)
	case screenSystem:
		return m.updateSystem(msg)
	case screenPathInput:
		return m.updatePath(msg)
	case screenLoading:
		return m.updateLoading(msg)
	case screenVersion:
		return m.updateVersion(msg)
	case screenVersionList:
		return m.updateVersionList(msg)
	case screenManage:
		return m.updateManage(msg)
	case screenConfirmRemove:
		return m.updateConfirmRemove(msg)
	case screenDecision:
		return m.updateDecision(msg)
	case screenProgress:
		return m.updateProgress(msg)
	case screenResult:
		return m.updateResult(msg)
	case screenError:
		return m.updateError(msg)
	case screenKeys:
		return m.updateKeys(msg)
	}
	return m, nil
}

func (m Model) updateSimpleQuit(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key(k) {
	case "enter", "q", "esc":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updateSystem(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	opts := 2
	switch key(k) {
	case "q":
		return m, tea.Quit
	case "esc", "left", "h":
		if m.changingDir {
			m.changingDir = false
			m.screen = screenManage
			return m, nil
		}
	case "up", "k":
		m.sysCursor = clampIndex(m.sysCursor, opts, -1)
	case "down", "j":
		m.sysCursor = clampIndex(m.sysCursor, opts, 1)
	case "home", "pgup":
		m.sysCursor = 0
	case "end", "pgdown":
		m.sysCursor = opts - 1
	case "tab":
		m.sysCursor = clampIndex(m.sysCursor, opts, 1)
	case "shift+tab":
		m.sysCursor = clampIndex(m.sysCursor, opts, -1)
	case "enter", "right", "l":
		if !m.plat.DefaultOK {
			if m.sysCursor == 0 {
				m.screen = screenPathInput
				m.pathInput.Focus()
				return m, nil
			}
			return m, tea.Quit
		}
		if m.sysCursor == 0 {
			return m.continueFromSystem()
		}
		m.screen = screenPathInput
		if m.cursorRoot != "" {
			m.pathInput.SetValue(m.cursorRoot)
		} else if m.plat.DefaultDir != "" {
			m.pathInput.SetValue(m.plat.DefaultDir)
		}
		m.pathInput.Focus()
		return m, nil
	}
	return m, nil
}

func (m Model) updatePath(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if ok {
		switch key(k) {
		case "esc":
			m.pathInput.Blur()
			m.screen = screenSystem
			m.pathErr = ""
			return m, nil
		case "enter":
			val := m.pathInput.Value()
			if val == "" {
				val = m.pathInput.Placeholder
			}
			m.cursorRoot = val
			m.pathInput.Blur()
			return m.continueFromSystem()
		}
	}
	var cmd tea.Cmd
	m.pathInput, cmd = m.pathInput.Update(msg)
	return m, cmd
}

func (m Model) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key(msg) == "q" || key(msg) == "esc" {
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg.err
		switch m.load {
		case loadTags, loadManifest:
			m.errRetry = screenSystem
		case loadBundle:
			m.errRetry = screenVersion
		}
		m.screen = screenError
		return m, nil
	case manifestMsg:
		return m.afterManifest(msg.m)
	case tagsMsg:
		return m.afterTags(msg.tags)
	case bundleMsg:
		return m.afterBundle(msg.bundle)
	}
	return m, nil
}

func (m Model) updateVersion(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key(k) {
	case "q":
		return m, tea.Quit
	case "esc", "left", "h":
		m.screen = screenSystem
		return m, nil
	case "up", "k":
		if m.versionCursor > 0 {
			m.versionCursor--
		}
	case "down", "j":
		if m.versionCursor < 1 {
			m.versionCursor++
		}
	case "home", "pgup":
		m.versionCursor = 0
	case "end", "pgdown":
		m.versionCursor = 1
	case "tab":
		if m.versionCursor < 1 {
			m.versionCursor++
		}
	case "shift+tab":
		if m.versionCursor > 0 {
			m.versionCursor--
		}
	case "enter", "right", "l":
		if m.versionCursor == 0 {
			m.useLatest = true
			if m.latest == "" {
				m.screen = screenVersionList
				return m, nil
			}
			m.selectedTag = m.latest
			return m.loadSelectedVersion()
		}
		m.useLatest = false
		m.screen = screenVersionList
		m.listCursor = 0
		m.listOffset = 0
		return m, nil
	}
	return m, nil
}

func (m Model) updateVersionList(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	n := len(m.tags)
	switch key(k) {
	case "q":
		return m, tea.Quit
	case "esc", "left", "h":
		m.screen = screenVersion
		return m, nil
	case "up", "k":
		m.listCursor = clampIndex(m.listCursor, n, -1)
		if m.listCursor < m.listOffset {
			m.listOffset = m.listCursor
		}
	case "down", "j":
		m.listCursor = clampIndex(m.listCursor, n, 1)
		visible := 10
		if m.listCursor >= m.listOffset+visible {
			m.listOffset = m.listCursor - visible + 1
		}
	case "home":
		m.listCursor = 0
		m.listOffset = 0
	case "end":
		if n > 0 {
			m.listCursor = n - 1
			visible := 10
			if m.listCursor >= visible {
				m.listOffset = m.listCursor - visible + 1
			}
		}
	case "pgup":
		m.listCursor = clampIndex(m.listCursor, n, -10)
		if m.listCursor < m.listOffset {
			m.listOffset = m.listCursor
		}
	case "pgdown":
		m.listCursor = clampIndex(m.listCursor, n, 10)
		visible := 10
		if m.listCursor >= m.listOffset+visible {
			m.listOffset = m.listCursor - visible + 1
		}
	case "enter", "right", "l":
		if n == 0 {
			return m, nil
		}
		m.useLatest = false
		m.selectedTag = m.tags[m.listCursor]
		return m.loadSelectedVersion()
	}
	return m, nil
}

func (m Model) updateManage(msg tea.Msg) (tea.Model, tea.Cmd) {
	if err, ok := msg.(openErrMsg); ok {
		m.manageErr = err.err.Error()
		return m, nil
	}
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	m = m.ensureManageCursor()
	rows := m.manageRows()
	if len(rows) == 0 {
		return m, nil
	}
	row := rows[m.manageCursor]
	switch key(k) {
	case "q":
		return m, tea.Quit
	case "esc":
		m.screen = screenVersion
		return m, nil
	case "up", "k":
		m = m.moveManage(-1)
	case "down", "j":
		m = m.moveManage(1)
	case "left", "h":
		next, back := m.collapseOrBack()
		if back {
			next.screen = screenVersion
		}
		return next, nil
	case "right", "l":
		m = m.expandOrForward()
	case "tab":
		m = m.jumpSection(1)
	case "shift+tab":
		m = m.jumpSection(-1)
	case "home":
		m = m.jumpManageEdge(-1)
	case "end":
		m = m.jumpManageEdge(1)
	case "pgup":
		m = m.moveManagePage(-1)
	case "pgdown":
		m = m.moveManagePage(1)
	case " ":
		m.manageErr = ""
		switch row.Kind {
		case rowGroup:
			m = m.toggleGroup(row.GroupIdx)
		case rowFile:
			m = m.toggleFile(row.GroupIdx, row.FileIdx)
		}
		m = m.ensureManageCursor()
	case "enter":
		m.manageErr = ""
		switch row.Kind {
		case rowGroup:
			m = m.expandOrForward()
		case rowAction:
			return m.activateAction(row.Action)
		}
	}
	return m, nil
}

func (m Model) activateAction(action string) (Model, tea.Cmd) {
	switch action {
	case "Install / Update selected":
		return m.beginInstall()
	case "Remove selected":
		return m.beginRemove()
	case "Change version":
		m.bundle = nil
		m.plan = nil
		m.groups = nil
		m.screen = screenVersion
		return m, nil
	case "Change directory":
		m.changingDir = true
		m.savedRoot = false
		m.screen = screenSystem
		m.pathInput.SetValue(m.cursorRoot)
		return m, nil
	case "Open Cursor folder":
		m.manageErr = ""
		if m.cursorRoot == "" {
			m.manageErr = "No Cursor directory selected."
			return m, nil
		}
		return m, openFolderCmd(m.cursorRoot, m.plat.OS)
	case "Keybindings":
		return m.openKeys(), nil
	case "Quit":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updateConfirmRemove(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key(k) {
	case "q":
		return m, tea.Quit
	case "esc", "left", "h":
		m = m.restoreManagePlan()
		m.screen = screenManage
		return m, nil
	case "up", "k", "home", "pgup", "shift+tab":
		m.removeCursor = 0
	case "down", "j", "end", "pgdown", "tab":
		m.removeCursor = 1
	case "enter", "right", "l":
		if m.removeCursor == 0 {
			if len(m.pending) > 0 {
				m.screen = screenDecision
				return m, nil
			}
			return m.startApply(resultRemove)
		}
		m = m.restoreManagePlan()
		m.screen = screenManage
		return m, nil
	}
	return m, nil
}

func (m Model) updateDecision(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	opts := m.decisionOptions()
	switch key(k) {
	case "q":
		return m, tea.Quit
	case "esc", "left", "h":
		m = m.restoreManagePlan()
		m.screen = screenManage
		return m, nil
	case "up", "k":
		m.decCursor = clampIndex(m.decCursor, len(opts), -1)
	case "down", "j":
		m.decCursor = clampIndex(m.decCursor, len(opts), 1)
	case "home", "pgup":
		m.decCursor = 0
	case "end", "pgdown":
		if len(opts) > 0 {
			m.decCursor = len(opts) - 1
		}
	case "tab":
		m.decCursor = clampIndex(m.decCursor, len(opts), 1)
	case "shift+tab":
		m.decCursor = clampIndex(m.decCursor, len(opts), -1)
	case "enter", "right", "l":
		choice := opts[m.decCursor].choice
		if choice == installer.ChoiceCancel {
			m = m.restoreManagePlan()
			m.screen = screenManage
			return m, nil
		}
		op := m.pending[m.pendIndex]
		m.choices[op.RelPath] = choice
		m.pendIndex++
		m.decCursor = 0
		if m.pendIndex >= len(m.pending) {
			kind := resultInstall
			if m.plan != nil && m.plan.Kind == installer.PlanUninstall {
				kind = resultRemove
			}
			return m.startApply(kind)
		}
		return m, nil
	}
	return m, nil
}

type decisionOpt struct {
	label  string
	choice installer.Choice
}

func (m Model) decisionOptions() []decisionOpt {
	if m.pendIndex >= len(m.pending) {
		return nil
	}
	op := m.pending[m.pendIndex]
	switch op.Kind {
	case installer.OpConflict:
		return []decisionOpt{
			{"Skip", installer.ChoiceSkip},
			{"Overwrite and take ownership", installer.ChoiceOverwrite},
			{"Cancel", installer.ChoiceCancel},
		}
	case installer.OpModifiedUpdate:
		ver := m.selectedTag
		if ver == "" && m.plan != nil {
			ver = m.plan.ContentVersion
		}
		return []decisionOpt{
			{"Keep local version", installer.ChoiceKeepLocal},
			{"Replace with " + ver, installer.ChoiceOverwrite},
			{"Cancel", installer.ChoiceCancel},
		}
	case installer.OpModifiedRemove:
		return []decisionOpt{
			{"Keep local file", installer.ChoiceKeepLocal},
			{"Remove anyway", installer.ChoiceOverwrite},
			{"Cancel", installer.ChoiceCancel},
		}
	default:
		return []decisionOpt{{"Cancel", installer.ChoiceCancel}}
	}
}

func (m Model) updateProgress(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key(msg) {
		case "q", "esc":
			if m.applyCancel != nil {
				m.applyCancel()
			}
			return m, nil
		}
	case progressMsg:
		m.applyProgress(msg.ev)
		return m, listenApply(m.applyCh)
	case doneMsg:
		m.applyCancel = nil
		if m.quitting {
			return m, tea.Quit
		}
		if msg.err != nil {
			if msg.err.Error() == "installation was cancelled" || msg.err.Error() == "installation cancelled" {
				m.result = resultCancelled
				m.resultErr = msg.err
			} else {
				m.result = resultFailed
				m.resultErr = msg.err
			}
		} else {
			m.countDone()
		}
		m.screen = screenResult
		return m, nil
	}
	return m, nil
}

func (m *Model) applyProgress(ev installer.Event) {
	for i := range m.progress {
		if m.progress[i].Path == ev.Path {
			m.progress[i].Status = ev.Status
			return
		}
	}
	m.progress = append(m.progress, progressLine{Path: ev.Path, Status: ev.Status})
}

func (m *Model) countDone() {
	n := 0
	for _, p := range m.progress {
		if p.Status == installer.EventDone {
			n++
		}
	}
	m.filesDone = n
}

func (m Model) updateResult(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key(k) {
	case "q":
		return m, tea.Quit
	case "enter", "esc":
		return m.returnToManage()
	}
	return m, nil
}

func (m Model) updateError(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key(k) {
	case "q":
		return m, tea.Quit
	case "up", "k", "home", "pgup", "shift+tab":
		m.decCursor = 0
	case "down", "j", "end", "pgdown", "tab":
		m.decCursor = 1
	case "enter", "right", "l":
		if m.decCursor == 1 {
			return m, tea.Quit
		}
		m.err = nil
		m.screen = m.errRetry
		if m.errRetry == screenSystem {
			return m.continueFromSystem()
		}
		if m.errRetry == screenVersion {
			m.screen = screenVersion
			return m, nil
		}
		if m.errRetry == screenManage {
			return m.returnToManage()
		}
		return m, nil
	case "esc", "left", "h":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) openKeys() Model {
	m.keysFrom = m.screen
	m.screen = screenKeys
	return m
}

func (m Model) updateKeys(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key(k) {
	case "q":
		return m, tea.Quit
	case "esc", "enter", "left", "h", "?":
		if m.keysFrom == screenKeys {
			m.keysFrom = screenManage
		}
		m.screen = m.keysFrom
		return m, nil
	}
	return m, nil
}
