package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type keyBind struct {
	Action  string
	Key     string
	Also    string
	Section string
}

func footerFor(s screen) [][2]string {
	switch s {
	case screenUnsupported:
		return [][2]string{{"enter", "exit"}}
	case screenSystem:
		return [][2]string{{"enter", "confirm"}, {"q", "quit"}}
	case screenPathInput:
		return [][2]string{{"enter", "confirm"}, {"esc", "back"}}
	case screenLoading:
		return [][2]string{{"q", "quit"}}
	case screenVersion, screenVersionList:
		return [][2]string{{"enter", "confirm"}, {"esc", "back"}, {"q", "quit"}}
	case screenManage:
		return [][2]string{{"space", "select"}, {"enter", "confirm"}, {"q", "quit"}}
	case screenConfirmRemove, screenDecision:
		return [][2]string{{"enter", "confirm"}, {"esc", "back"}, {"q", "quit"}}
	case screenProgress:
		return [][2]string{{"ctrl+c", "quit"}}
	case screenResult:
		return [][2]string{{"enter", "confirm"}, {"q", "quit"}}
	case screenError:
		return [][2]string{{"enter", "confirm"}, {"q", "quit"}}
	case screenKeys:
		return [][2]string{{"esc", "back"}, {"q", "quit"}}
	default:
		return [][2]string{{"q", "quit"}}
	}
}

func keysFor(s screen) []keyBind {
	nav := []keyBind{
		{Section: "Move", Action: "Move up / down", Key: "↑ ↓", Also: "Up / Down"},
		{Section: "Move", Action: "Change value / page", Key: "← →", Also: "Left / Right"},
		{Section: "Move", Action: "Next section", Key: "Tab", Also: "Shift+Tab"},
	}
	app := []keyBind{
		{Section: "App", Action: "Open folder", Key: "o"},
		{Section: "App", Action: "Repair TUI", Key: "Ctrl+L"},
		{Section: "App", Action: "Keybindings", Key: "?"},
		{Section: "App", Action: "Quit", Key: "q"},
	}
	switch s {
	case screenPathInput:
		return []keyBind{
			{Section: "Edit", Action: "Confirm", Key: "Enter"},
			{Section: "App", Action: "Back", Key: "Escape"},
		}
	case screenManage:
		out := append([]keyBind{}, nav...)
		out = append(out,
			keyBind{Section: "List", Action: "Toggle component / file", Key: "Space"},
			keyBind{Section: "List", Action: "Expand / confirm", Key: "Enter"},
			keyBind{Section: "List", Action: "Collapse / back", Key: "←", Also: "h"},
		)
		return append(out, app...)
	case screenVersionList:
		out := append([]keyBind{}, nav[:2]...)
		out = append(out, keyBind{Section: "List", Action: "Confirm", Key: "Enter"})
		return append(out, app...)
	case screenResult:
		return []keyBind{
			{Section: "App", Action: "Continue", Key: "Enter"},
			{Section: "App", Action: "Repair TUI", Key: "Ctrl+L"},
			{Section: "App", Action: "Quit", Key: "q", Also: "Ctrl+C"},
		}
	case screenProgress:
		return []keyBind{{Section: "App", Action: "Quit", Key: "Ctrl+C", Also: "q"}}
	case screenLoading:
		return []keyBind{{Section: "App", Action: "Quit", Key: "q", Also: "Ctrl+C"}}
	default:
		out := append([]keyBind{}, nav...)
		out = append(out, keyBind{Section: "List", Action: "Confirm", Key: "Enter"})
		return append(out, app...)
	}
}

func viewKeyHelp(from screen) string {
	rows := keysFor(from)
	if len(rows) == 0 {
		rows = keysFor(screenManage)
	}
	keyW := 0
	rendered := make([]string, len(rows))
	for i, row := range rows {
		rendered[i] = plainKeys(row)
		if w := lipgloss.Width(rendered[i]); w > keyW {
			keyW = w
		}
	}
	var b strings.Builder
	b.WriteString(labelStyle.Render("Keybindings") + "\n")
	prev := ""
	for i, row := range rows {
		if row.Section != prev {
			b.WriteString("\n")
			if row.Section != "" {
				b.WriteString(dimStyle.Render(row.Section) + "\n")
			}
			prev = row.Section
		}
		keys := rendered[i]
		pad := keyW + 4 - lipgloss.Width(keys)
		if pad < 2 {
			pad = 2
		}
		b.WriteString("  " + keyTextStyle.Render(keys) + strings.Repeat(" ", pad) + row.Action + "\n")
	}
	return b.String()
}

func plainKeys(row keyBind) string {
	if row.Also == "" {
		return row.Key
	}
	return row.Key + "  " + row.Also
}
