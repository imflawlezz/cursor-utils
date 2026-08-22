package tui

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/imflawlezz/cursor-utils/installer/internal/content"
	"github.com/imflawlezz/cursor-utils/installer/internal/installer"
	"github.com/imflawlezz/cursor-utils/installer/internal/manifest"
	"github.com/imflawlezz/cursor-utils/installer/internal/platform"
	"github.com/imflawlezz/cursor-utils/installer/internal/prefs"
)

type tagsMsg struct{ tags []string }
type bundleMsg struct{ bundle *content.Bundle }
type manifestMsg struct{ m *manifest.Manifest }
type errMsg struct{ err error }
type openErrMsg struct{ err error }
type progressMsg struct{ ev installer.Event }
type doneMsg struct{ err error }
type useSavedRootMsg struct{}

func openFolderCmd(dir, goos string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		switch goos {
		case "darwin":
			cmd = exec.Command("open", dir)
		case "linux":
			cmd = exec.Command("xdg-open", dir)
		case "windows":
			cmd = exec.Command("explorer", dir)
		default:
			return openErrMsg{fmt.Errorf("opening folders is not supported on this OS")}
		}
		if err := cmd.Start(); err != nil {
			return openErrMsg{err}
		}
		go cmd.Wait()
		return nil
	}
}

func fetchTagsCmd(eng *installer.Engine) tea.Cmd {
	return func() tea.Msg {
		tags, err := eng.ListTags(context.Background())
		if err != nil {
			return errMsg{err}
		}
		return tagsMsg{tags}
	}
}

func fetchBundleCmd(eng *installer.Engine, tag string) tea.Cmd {
	return func() tea.Msg {
		bundle, err := eng.Fetch(context.Background(), tag)
		if err != nil {
			return errMsg{err}
		}
		return bundleMsg{bundle}
	}
}

func loadManifestCmd(eng *installer.Engine, root string) tea.Cmd {
	return func() tea.Msg {
		m, err := eng.LoadManifest(root)
		if err != nil {
			return errMsg{err}
		}
		return manifestMsg{m}
	}
}

func listenApply(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return doneMsg{err: context.Canceled}
		}
		return msg
	}
}

func (m Model) startApply(kind resultKind) (Model, tea.Cmd) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan tea.Msg, 32)
	eng := m.engine
	plan := m.plan
	choices := m.choices
	go func() {
		err := eng.Apply(ctx, plan, choices, func(ev installer.Event) {
			select {
			case ch <- progressMsg{ev}:
			case <-ctx.Done():
			}
		})
		ch <- doneMsg{err}
		close(ch)
	}()
	m.applyCancel = cancel
	m.applyCh = ch
	m.screen = screenProgress
	m.result = kind
	m.progress = nil
	m.filesDone = 0
	m.resultVer = plan.ContentVersion
	m.resultComps = uniqueComponents(plan)
	return m, listenApply(ch)
}

func uniqueComponents(plan *installer.Plan) []string {
	if plan == nil {
		return nil
	}
	seen := map[string]struct{}{}
	var out []string
	for _, spec := range content.Registry {
		id := string(spec.ID)
		for _, op := range plan.Ops {
			if op.Component != id {
				continue
			}
			if _, ok := seen[id]; ok {
				break
			}
			seen[id] = struct{}{}
			out = append(out, spec.DisplayName)
		}
	}
	return out
}

func (m Model) continueFromSystem() (Model, tea.Cmd) {
	path := m.cursorRoot
	if path == "" {
		m.pathErr = "Choose a Cursor directory to continue."
		m.screen = screenPathInput
		m.pathInput.Focus()
		return m, nil
	}
	expanded, err := platform.ExpandUser(path, m.plat.Home)
	if err != nil {
		m.pathErr = err.Error()
		m.screen = screenPathInput
		m.pathInput.Focus()
		return m, nil
	}
	if err := platform.ValidateCursorRoot(expanded); err != nil {
		m.pathErr = err.Error()
		m.screen = screenPathInput
		m.pathInput.SetValue(expanded)
		m.pathInput.Focus()
		return m, nil
	}
	m.cursorRoot = expanded
	m.pathErr = ""
	m.changingDir = false
	_ = prefs.Save(m.plat.Home, prefs.File{CursorRoot: expanded})
	m.savedRoot = true
	m.screen = screenLoading
	m.load = loadManifest
	m.loadNote = "Reading installation state…"
	return m, loadManifestCmd(m.engine, m.cursorRoot)
}

func (m Model) afterManifest(man *manifest.Manifest) (Model, tea.Cmd) {
	m.existing = man
	if m.bundle != nil {
		return m.afterBundle(m.bundle)
	}
	m.screen = screenLoading
	m.load = loadTags
	m.loadNote = "Fetching versions…"
	return m, fetchTagsCmd(m.engine)
}

func (m Model) afterTags(tags []string) (Model, tea.Cmd) {
	m.tags = tags
	if latest, ok := m.engine.Latest(tags); ok {
		m.latest = latest
	} else {
		m.latest = ""
	}
	if m.useLatest && m.latest != "" {
		m.selectedTag = m.latest
	} else if m.selectedTag == "" && len(tags) > 0 {
		m.selectedTag = tags[0]
		m.useLatest = false
	}
	m.screen = screenVersion
	m.versionCursor = 0
	if !m.useLatest {
		m.versionCursor = 1
	}
	return m, nil
}

func (m Model) loadSelectedVersion() (Model, tea.Cmd) {
	tag := m.selectedTag
	if m.useLatest {
		if m.latest == "" {
			m.err = errNoLatest{}
			m.errRetry = screenVersion
			m.screen = screenError
			return m, nil
		}
		tag = m.latest
		m.selectedTag = tag
	}
	if tag == "" {
		m.screen = screenVersionList
		return m, nil
	}
	m.screen = screenLoading
	m.load = loadBundle
	m.loadNote = "Loading " + tag + "…"
	return m, fetchBundleCmd(m.engine, tag)
}

func (m Model) afterBundle(bundle *content.Bundle) (Model, tea.Cmd) {
	m.bundle = bundle
	plan, err := m.engine.PlanInstall(m.cursorRoot, bundle, m.existing)
	if err != nil {
		m.err = err
		m.errRetry = screenVersion
		m.screen = screenError
		return m, nil
	}
	m.plan = plan
	m.initSelection()
	m.manageCursor = 0
	m.manageErr = ""
	m.screen = screenManage
	m = m.snapManageCursor()
	return m, nil
}

type errNoLatest struct{}

func (errNoLatest) Error() string {
	return "No semantic version tags were found.\n\nChoose a specific version instead."
}

func (m Model) beginInstall() (Model, tea.Cmd) {
	sel := m.selection()
	if sel.Empty() {
		m.manageErr = "Select at least one component or file."
		return m, nil
	}
	plan, err := m.engine.PlanInstall(m.cursorRoot, m.bundle, m.existing)
	if err != nil {
		m.err = err
		m.errRetry = screenManage
		m.screen = screenError
		return m, nil
	}
	work := plan.Filter(sel)
	if len(work.Ops) == 0 {
		m.manageErr = "Nothing to install or update in the current selection."
		return m, nil
	}
	m.plan = work
	m.choices = map[string]installer.Choice{}
	m.pending = m.plan.NeedsDecision()
	m.pendIndex = 0
	m.decCursor = 0
	m.manageErr = ""
	if len(m.pending) > 0 {
		m.screen = screenDecision
		return m, nil
	}
	return m.startApply(resultInstall)
}

func (m Model) beginRemove() (Model, tea.Cmd) {
	sel := m.selection()
	if sel.Empty() {
		m.manageErr = "Select at least one component or file."
		return m, nil
	}
	if m.existing == nil || m.existing.Empty() {
		m.manageErr = "No cursor-utils content is installed."
		return m, nil
	}
	plan, err := m.engine.PlanUninstall(m.cursorRoot, m.existing)
	if err != nil {
		m.err = err
		m.errRetry = screenManage
		m.screen = screenError
		return m, nil
	}
	work := plan.Filter(sel)
	if len(work.Ops) == 0 {
		m.manageErr = "The selection does not include any installed cursor-utils files."
		return m, nil
	}
	m.plan = work
	m.choices = map[string]installer.Choice{}
	m.pending = m.plan.NeedsDecision()
	m.pendIndex = 0
	m.decCursor = 0
	m.manageErr = ""
	m.screen = screenConfirmRemove
	m.removeCursor = 1
	return m, nil
}

func (m Model) restoreManagePlan() Model {
	if m.bundle == nil || m.engine == nil {
		return m
	}
	plan, err := m.engine.PlanInstall(m.cursorRoot, m.bundle, m.existing)
	if err == nil {
		m.plan = plan
	}
	return m
}

func (m Model) snapManageCursor() Model {
	return m.ensureManageCursor()
}

func (m Model) returnToManage() (Model, tea.Cmd) {
	m.applyCancel = nil
	m.progress = nil
	m.manageErr = ""
	m.screen = screenLoading
	m.load = loadManifest
	m.loadNote = "Refreshing…"
	return m, loadManifestCmd(m.engine, m.cursorRoot)
}

func key(msg tea.KeyMsg) string {
	return strings.ToLower(msg.String())
}
