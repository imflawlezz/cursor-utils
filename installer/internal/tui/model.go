package tui

import (
	"context"
	"path"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/imflawlezz/cursor-utils/installer/internal/content"
	"github.com/imflawlezz/cursor-utils/installer/internal/installer"
	"github.com/imflawlezz/cursor-utils/installer/internal/manifest"
	"github.com/imflawlezz/cursor-utils/installer/internal/platform"
	"github.com/imflawlezz/cursor-utils/installer/internal/prefs"
)

type screen int

const (
	screenUnsupported screen = iota
	screenSystem
	screenPathInput
	screenLoading
	screenVersion
	screenVersionList
	screenManage
	screenConfirmRemove
	screenDecision
	screenProgress
	screenResult
	screenError
	screenKeys
)

type loadKind int

const (
	loadTags loadKind = iota
	loadBundle
	loadManifest
)

type resultKind int

const (
	resultInstall resultKind = iota
	resultRemove
	resultCancelled
	resultFailed
)

type progressLine struct {
	Path   string
	Status installer.EventStatus
}

type fileState struct {
	RelPath  string
	Selected bool
}

type groupState struct {
	ID       string
	Expanded bool
	Files    []fileState
}

func (g groupState) allSelected() bool {
	if len(g.Files) == 0 {
		return false
	}
	for _, f := range g.Files {
		if !f.Selected {
			return false
		}
	}
	return true
}

func (g groupState) anySelected() bool {
	for _, f := range g.Files {
		if f.Selected {
			return true
		}
	}
	return false
}

type manageSection int

const (
	secNone manageSection = iota
	secComponents
	secApply
	secSettings
	secApp
)

type manageRowKind int

const (
	rowGroup manageRowKind = iota
	rowFile
	rowAction
	rowLabel
	rowSpacer
)

type manageRow struct {
	Kind     manageRowKind
	GroupIdx int
	FileIdx  int
	Action   string
	Section  manageSection
}

func (r manageRow) focusable() bool {
	switch r.Kind {
	case rowGroup, rowFile, rowAction:
		return true
	default:
		return false
	}
}

type Model struct {
	engine *installer.Engine
	plat   platform.Info

	screen screen
	width  int
	height int

	cursorRoot  string
	pathInput   textinput.Model
	pathErr     string
	savedRoot   bool
	changingDir bool

	sysCursor int

	tags          []string
	latest        string
	useLatest     bool
	selectedTag   string
	versionCursor int
	listCursor    int
	listOffset    int

	bundle   *content.Bundle
	existing *manifest.Manifest
	plan     *installer.Plan
	groups   []groupState

	manageCursor int
	manageErr    string
	removeCursor int

	choices   map[string]installer.Choice
	pending   []installer.Op
	pendIndex int
	decCursor int

	progress    []progressLine
	applyCh     <-chan tea.Msg
	applyCancel context.CancelFunc
	result      resultKind
	resultErr   error
	filesDone   int
	resultVer   string
	resultComps []string

	err      error
	errRetry screen
	keysFrom screen
	load     loadKind
	loadNote string

	wantH     int
	wantW     int
	haveSize  bool
	userSized bool
	openURL   func(string) error

	quitting bool
}

func New(eng *installer.Engine, plat platform.Info) Model {
	ti := textinput.New()
	ti.Placeholder = plat.DefaultDir
	ti.Width = 64
	ti.Prompt = "> "
	ti.CharLimit = 512

	m := Model{
		engine:    eng,
		plat:      plat,
		useLatest: true,
		choices:   map[string]installer.Choice{},
		pathInput: ti,
	}
	if !plat.Supported {
		m.screen = screenUnsupported
		return m
	}
	m.screen = screenSystem
	if plat.DefaultOK {
		m.cursorRoot = plat.DefaultDir
	}
	if saved, err := prefs.Load(plat.Home); err == nil && saved != nil && saved.CursorRoot != "" {
		if platform.ValidateCursorRoot(saved.CursorRoot) == nil {
			m.cursorRoot = saved.CursorRoot
			m.savedRoot = true
		}
	}
	return m
}

func (m Model) Init() tea.Cmd {
	if m.savedRoot {
		return tea.Batch(sizePollCmd(), func() tea.Msg { return useSavedRootMsg{} })
	}
	return tea.Batch(sizePollCmd(), textinput.Blink)
}

func fileLabel(rel string) string {
	return path.Base(rel)
}

func (g groupState) displayName() string {
	return content.DisplayName(g.ID)
}
