package tui

import (
	"regexp"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/imflawlezz/cursor-utils/installer/internal/config"
	"github.com/imflawlezz/cursor-utils/installer/internal/content"
	"github.com/imflawlezz/cursor-utils/installer/internal/installer"
	"github.com/imflawlezz/cursor-utils/installer/internal/platform"
)

var (
	ansiRE = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	osc8RE = regexp.MustCompile(`\x1b\]8;.*?(?:\x07|\x1b\\)`)
)

func stripANSI(s string) string {
	s = osc8RE.ReplaceAllString(s, "")
	return ansiRE.ReplaceAllString(s, "")
}

func testModel(plat platform.Info) Model {
	return New(installer.New(config.Default(), nil, plat), plat)
}

func TestSystemScreen(t *testing.T) {
	m := testModel(platform.Info{
		OS: "darwin", DisplayName: "macOS", Arch: "arm64",
		Home: "/Users/ada", DefaultDir: "/Users/ada/.cursor",
		DefaultOK: true, Supported: true,
	})
	view := m.View()
	if !strings.Contains(view, "macOS") || !strings.Contains(view, "/Users/ada/.cursor") {
		t.Fatalf("view = %s", view)
	}
	if !strings.Contains(view, "Continue") {
		t.Fatal(view)
	}
	if !strings.Contains(view, "cursor-utils installer") || !strings.Contains(view, "imflawlezz") {
		t.Fatal(view)
	}
	if !strings.Contains(view, "1.2.0") {
		t.Fatal("expected installer version 1.2.0")
	}
	if !strings.Contains(view, "https://github.com/imflawlezz") {
		t.Fatal("expected github profile hyperlink")
	}
	if !strings.Contains(view, "┌─") || !strings.Contains(stripANSI(view), "q  quit") {
		t.Fatalf("expected framed chrome:\n%s", view)
	}
}

func TestUnsupportedOS(t *testing.T) {
	m := testModel(platform.Info{OS: "plan9", DisplayName: "Plan 9", Supported: false})
	if m.screen != screenUnsupported {
		t.Fatalf("screen = %v", m.screen)
	}
	view := m.View()
	if !strings.Contains(view, "Plan 9 is not yet supported") {
		t.Fatal(view)
	}
}

func TestMissingDefaultRequiresManualChoice(t *testing.T) {
	m := testModel(platform.Info{
		OS: "linux", DisplayName: "Linux", Supported: true, DefaultOK: false,
	})
	view := m.View()
	if !strings.Contains(view, "could not be detected") {
		t.Fatal(view)
	}
	if !strings.Contains(view, "Choose location") {
		t.Fatal(view)
	}
}

func TestFileLabelKeepsExtension(t *testing.T) {
	if got := fileLabel("commands/commit.md"); got != "commit.md" {
		t.Fatalf("got %q", got)
	}
}

func TestManageShowsDiskStateAndOptions(t *testing.T) {
	m := testModel(platform.Info{
		OS: "darwin", DisplayName: "macOS", Arch: "arm64",
		Home: "/Users/ada", DefaultDir: "/Users/ada/.cursor",
		DefaultOK: true, Supported: true,
	})
	m.screen = screenManage
	m.selectedTag = "v0.1.0"
	m.cursorRoot = "/Users/ada/.cursor"
	m.groups = []groupState{{
		ID:       "commands",
		Expanded: true,
		Files: []fileState{
			{RelPath: "commands/commit.md", Selected: true},
			{RelPath: "commands/docs.md", Selected: false},
		},
	}}
	m.plan = &installer.Plan{
		ContentVersion: "v0.1.0",
		Ops: []installer.Op{
			{Kind: installer.OpUnchanged, Component: "commands", RelPath: "commands/commit.md"},
			{Kind: installer.OpAdd, Component: "commands", RelPath: "commands/docs.md"},
		},
	}
	view := m.View()
	for _, want := range []string{
		"commit.md", "docs.md", "installed", "not installed",
		"Actions", "Change directory", "Open Cursor folder", "Keybindings", "Repair TUI", "1/2 installed",
		"Install / Update selected", "Remove selected",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("missing %q in\n%s", want, view)
		}
	}
	if strings.Contains(view, "Options") {
		t.Fatal("Options header should be gone")
	}
	if strings.Contains(view, "[x]") || strings.Contains(view, "[ ]") {
		t.Fatalf("checkbox marks should be gone:\n%s", view)
	}
	plain := stripANSI(view)
	if !strings.Contains(plain, "✓") {
		t.Fatal("selected files should show ✓")
	}
	if !strings.Contains(plain, "✓ commit.md") {
		t.Fatalf("expected left-side checkbox before filename:\n%s", plain)
	}
	if strings.Contains(plain, "installed  ✓") || strings.Contains(plain, "installed ✓") {
		t.Fatalf("checkbox should not be on the right:\n%s", plain)
	}
	assertActionGap(t, stripANSI(view))
}

func TestManageHidesApplyWhenNothingSelected(t *testing.T) {
	m := testModel(platform.Info{
		OS: "darwin", DisplayName: "macOS", Supported: true, DefaultOK: true,
		Home: "/Users/ada",
	})
	m.screen = screenManage
	m.groups = []groupState{{
		ID: "commands",
		Files: []fileState{
			{RelPath: "commands/commit.md", Selected: false},
		},
	}}
	view := m.View()
	if strings.Contains(view, "Install / Update selected") || strings.Contains(view, "Remove selected") {
		t.Fatalf("apply actions should be hidden\n%s", view)
	}
	if !strings.Contains(view, "Change version") || !strings.Contains(view, "Quit") {
		t.Fatal(view)
	}
}

func TestManageKeysExpandSelectAndSection(t *testing.T) {
	m := testModel(platform.Info{
		OS: "darwin", DisplayName: "macOS", Supported: true, DefaultOK: true,
		DefaultDir: "/Users/ada/.cursor", Home: "/Users/ada",
	})
	m.screen = screenManage
	m.groups = []groupState{{
		ID:       "commands",
		Expanded: false,
		Files:    []fileState{{RelPath: "commands/commit.md", Selected: true}},
	}}
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	got, ok := next.(Model)
	if !ok || !got.groups[0].Expanded {
		t.Fatal("right should expand the group")
	}
	next, _ = got.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	got, ok = next.(Model)
	if !ok || got.groups[0].Files[0].Selected {
		t.Fatal("space should toggle selection")
	}
	next, _ = got.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got, ok = next.(Model)
	if !ok || !got.groups[0].Expanded {
		t.Fatal("enter on a group should open/expand, not toggle")
	}
	next, _ = got.Update(tea.KeyMsg{Type: tea.KeyTab})
	got, ok = next.(Model)
	if !ok {
		t.Fatal("tab should jump section")
	}
	rows := got.manageRows()
	if !rows[got.manageCursor].focusable() || rows[got.manageCursor].Section == secComponents {
		t.Fatalf("tab should leave components, cursor=%d kind=%v sec=%v", got.manageCursor, rows[got.manageCursor].Kind, rows[got.manageCursor].Section)
	}
}

func TestConfirmRemoveIsShort(t *testing.T) {
	m := testModel(platform.Info{
		OS: "darwin", DisplayName: "macOS", Supported: true, DefaultOK: true,
		Home: "/Users/ada",
	})
	m.screen = screenConfirmRemove
	m.plan = &installer.Plan{Ops: []installer.Op{
		{RelPath: "commands/commit.md"},
		{RelPath: "commands/changelog.md"},
		{RelPath: "commands/docs.md"},
	}}
	view := m.View()
	if !strings.Contains(view, "Remove 3 cursor-utils files?") {
		t.Fatalf("view = %s", view)
	}
	for _, want := range []string{"commands/", "commit.md", "changelog.md", "docs.md"} {
		if !strings.Contains(view, want) {
			t.Fatalf("missing %q in %s", want, view)
		}
	}
	if strings.Contains(view, "commands/commit.md") {
		t.Fatal("files should be grouped under the directory, not listed as full paths")
	}
	if strings.Contains(view, "Your other Cursor files") {
		t.Fatal("expected shorter confirm copy")
	}
}

func TestFormatRemoveListTruncatesAfterFive(t *testing.T) {
	paths := []string{
		"commands/a.md", "commands/b.md", "commands/c.md",
		"commands/d.md", "commands/e.md", "commands/f.md", "rules/g.md",
	}
	got := formatRemoveList(paths, 5)
	if !strings.Contains(got, "commands/") || strings.Contains(got, "rules/") {
		t.Fatalf("should group and stop after 5 files:\n%s", got)
	}
	if !strings.Contains(got, "a.md") || !strings.Contains(got, "e.md") {
		t.Fatalf("expected first 5 files:\n%s", got)
	}
	if strings.Contains(got, "f.md") || strings.Contains(got, "g.md") {
		t.Fatalf("sixth file should be truncated:\n%s", got)
	}
	if !strings.Contains(got, "…and 2 more") {
		t.Fatalf("expected remainder:\n%s", got)
	}
}

func TestKeybindingsScreen(t *testing.T) {
	m := testModel(platform.Info{
		OS: "darwin", DisplayName: "macOS", Supported: true, DefaultOK: true,
		DefaultDir: "/Users/ada/.cursor", Home: "/Users/ada",
	})
	m.screen = screenManage
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	got, ok := next.(Model)
	if !ok || got.screen != screenKeys {
		t.Fatal("? should open keybindings")
	}
	view := got.View()
	for _, want := range []string{"Keybindings", "Move", "List", "App", "Move up / down", "Toggle component / file", "Shift+Tab", "Ctrl+L"} {
		if !strings.Contains(view, want) {
			t.Fatalf("missing %q in\n%s", want, view)
		}
	}
	if strings.Contains(view, "Also") {
		t.Fatal("should not use an Also column")
	}
	next, _ = got.Update(tea.KeyMsg{Type: tea.KeyEsc})
	got, ok = next.(Model)
	if !ok || got.screen != screenManage {
		t.Fatal("esc should return to the previous screen")
	}

	sys := testModel(platform.Info{
		OS: "darwin", DisplayName: "macOS", Supported: true, DefaultOK: true,
		Home: "/Users/ada",
	})
	next, _ = sys.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	got, ok = next.(Model)
	if !ok || got.screen != screenKeys {
		t.Fatal("? should open keybindings from system")
	}
	view = got.View()
	if strings.Contains(view, "Toggle component / file") || strings.Contains(view, "Collapse / back") {
		t.Fatalf("system keys should not list manage-only bindings\n%s", view)
	}
	if !strings.Contains(view, "Confirm") {
		t.Fatal(view)
	}
}

func TestChangeDirectoryKeepsVersion(t *testing.T) {
	m := testModel(platform.Info{
		OS: "darwin", DisplayName: "macOS", Supported: true, DefaultOK: true,
		DefaultDir: "/Users/ada/.cursor", Home: "/Users/ada",
	})
	m.screen = screenManage
	m.selectedTag = "v0.1.0"
	m.bundle = content.NewBundle("v0.1.0")
	m.cursorRoot = "/Users/ada/.cursor"
	next, _ := m.activateAction("Change directory")
	if next.screen != screenSystem {
		t.Fatalf("screen = %v", next.screen)
	}
	if next.selectedTag != "v0.1.0" || next.bundle == nil {
		t.Fatal("changing directory should keep the selected version")
	}
	if !next.changingDir {
		t.Fatal("expected changingDir")
	}
	back, _ := next.Update(tea.KeyMsg{Type: tea.KeyEsc})
	got, ok := back.(Model)
	if !ok || got.screen != screenManage {
		t.Fatal("esc should return to manage without asking for a version")
	}
}

func TestQuitKey(t *testing.T) {
	m := testModel(platform.Info{
		OS: "darwin", DisplayName: "macOS", Supported: true, DefaultOK: true,
		DefaultDir: "/Users/ada/.cursor", Home: "/Users/ada",
	})
	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	_ = next
	if cmd == nil {
		t.Fatal("expected quit command")
	}
}

func assertActionGap(t *testing.T, got string) {
	t.Helper()
	lines := strings.Split(got, "\n")
	open := -1
	change := -1
	for i, line := range lines {
		if strings.Contains(line, "Open Cursor folder") {
			open = i
		}
		if strings.Contains(line, "Change version") {
			change = i
		}
	}
	if open < 0 || change < 0 {
		t.Fatalf("missing Open Cursor folder / Change version:\n%s", got)
	}
	if change != open+2 {
		t.Fatalf("Change version should sit one blank line under Open Cursor folder (open=%d change=%d):\n%s", open, change, got)
	}
	if strings.TrimSpace(strings.Trim(lines[open+1], " │")) != "" {
		t.Fatalf("expected blank line after Open Cursor folder, got %q", lines[open+1])
	}
}

func TestWrapFrameTitle(t *testing.T) {
	raw := wrapFrame("body", "", "1.2.0", 64, 0)
	if !strings.Contains(raw, "\x1b]8;;https://github.com/imflawlezz") {
		t.Fatalf("author should be a GitHub hyperlink:\n%q", raw)
	}
	got := stripANSI(raw)
	lines := strings.Split(got, "\n")
	if len(lines) < 3 {
		t.Fatalf("too few lines:\n%s", got)
	}
	top := lines[0]
	if !strings.HasPrefix(top, "┌─") || !strings.HasSuffix(top, "┐") {
		t.Fatalf("top border: %q", top)
	}
	if !strings.Contains(top, "cursor-utils installer") || !strings.Contains(top, "v1.2.0") {
		t.Fatalf("title missing: %q", top)
	}
	if !strings.Contains(top, "by imflawlezz") {
		t.Fatalf("author missing: %q", top)
	}
	x0, x1 := authorNameCells("1.2.0")
	if got := string([]rune(top)[x0:x1]); got != "imflawlezz" {
		t.Fatalf("author cells %d:%d = %q in %q", x0, x1, got, top)
	}
	if w := len([]rune(top)); w != 64 {
		t.Fatalf("top width %d want 64: %q", w, top)
	}
}

func TestWrapFrameFillsHeight(t *testing.T) {
	got := stripANSI(wrapFrame("body", "", "1.2.0", 56, 12))
	lines := strings.Split(got, "\n")
	if len(lines) != 12 {
		t.Fatalf("height %d want 12:\n%s", len(lines), got)
	}
}

func TestHelpAnchoredAtBottom(t *testing.T) {
	got := stripANSI(wrapFrame("body", "up/down  confirm", "1.2.0", 56, 16))
	lines := strings.Split(got, "\n")
	if len(lines) != 16 {
		t.Fatalf("height %d want 16:\n%s", len(lines), got)
	}
	if !strings.Contains(lines[len(lines)-3], "up/down") {
		t.Fatalf("help should sit above the bottom border:\n%s", got)
	}
}

func TestClickAuthorOpensGitHub(t *testing.T) {
	var opened string
	m := testModel(platform.Info{
		OS: "darwin", DisplayName: "macOS", Supported: true, DefaultOK: true,
		Home: "/Users/ada", DefaultDir: "/Users/ada/.cursor",
	})
	m.width = 80
	m.height = 24
	m.openURL = func(u string) error {
		opened = u
		return nil
	}
	x0, x1 := authorNameCells(m.version())
	if m.linkAt(x0, 0) != config.Default().AuthorURL() || m.linkAt(x1-1, 0) != config.Default().AuthorURL() {
		t.Fatalf("linkAt author %d-%d", x0, x1)
	}
	if m.linkAt(x0-1, 0) != "" || m.linkAt(x1, 0) != "" {
		t.Fatal("click outside author should miss")
	}
	m.Update(tea.MouseMsg{
		X:      x0,
		Y:      0,
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
	})
	if opened != config.Default().AuthorURL() {
		t.Fatalf("opened %q", opened)
	}
}

func TestRepairDisplayResetsResizeState(t *testing.T) {
	m := testModel(platform.Info{
		OS: "darwin", DisplayName: "macOS", Supported: true, DefaultOK: true,
		Home: "/Users/ada",
	})
	m.width = 100
	m.height = 36
	m.haveSize = true
	m.userSized = true
	m.wantH = 36
	m.wantW = 100
	next, cmd := m.repairDisplay()
	if next.userSized {
		t.Fatal("repair should allow auto-fit again")
	}
	if next.wantH != 0 || next.wantW != 0 {
		t.Fatalf("repair should forget last resize request, got %dx%d", next.wantW, next.wantH)
	}
	if cmd == nil {
		t.Fatal("expected clear screen + size query")
	}
}

func TestFillSelectedStaysOneLine(t *testing.T) {
	line := fillSelected("  1. "+strings.Repeat("Long Name ", 12)+"  installed  ✓", 40, false)
	plain := stripANSI(line)
	if strings.Contains(plain, "\n") {
		t.Fatalf("selected row wrapped:\n%q", plain)
	}
}
