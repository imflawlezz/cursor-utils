package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	minTermRows = 22
	minTermCols = 72
)

type sizePollMsg struct{}

func (m *Model) resizeIfNeeded() tea.Cmd {
	if m.quitting || m.height <= 0 || m.userSized {
		return nil
	}
	rows := m.height
	cols := m.width
	if cols < minTermCols {
		cols = minTermCols
	}
	if rows < minTermRows {
		rows = minTermRows
	}
	if want := m.neededHeight(); want > rows {
		rows = want
	}
	if rows == m.height && cols == m.width {
		return nil
	}
	if rows == m.wantH && cols == m.wantW {
		return nil
	}
	m.wantH = rows
	m.wantW = cols
	return resizeTerm(rows, cols)
}

func (m Model) neededHeight() int {
	return countViewLines(wrapFrame(m.screenBody(), m.helpLine(), m.version(), m.frameWidth(), 0))
}

func countViewLines(s string) int {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

func resizeTerm(rows, cols int) tea.Cmd {
	if rows < 1 {
		return nil
	}
	return func() tea.Msg {
		sendTermResize(rows, cols)
		return nil
	}
}
