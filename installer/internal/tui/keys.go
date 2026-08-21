package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type keyBind struct {
	Action string
	Key    string
	Also   string
}

func footerFor(s screen) [][2]string {
	switch s {
	case screenUnsupported:
		return [][2]string{{"enter", "exit"}}
	case screenSystem:
		return [][2]string{{"↑↓", "up/down"}, {"enter", "confirm"}, {"?", "keys"}, {"q", "quit"}}
	case screenPathInput:
		return [][2]string{{"enter", "confirm"}, {"esc", "back"}}
	case screenLoading:
		return [][2]string{{"q", "quit"}}
	case screenVersion, screenVersionList:
		return [][2]string{{"↑↓", "up/down"}, {"enter", "confirm"}, {"esc", "back"}, {"?", "keys"}, {"q", "quit"}}
	case screenManage:
		return [][2]string{{"↑↓", "up/down"}, {"space", "select"}, {"enter", "confirm"}, {"?", "keys"}, {"q", "quit"}}
	case screenConfirmRemove, screenDecision:
		return [][2]string{{"↑↓", "up/down"}, {"enter", "confirm"}, {"esc", "back"}, {"?", "keys"}, {"q", "quit"}}
	case screenProgress:
		return [][2]string{{"ctrl+c", "quit"}}
	case screenResult:
		return [][2]string{{"enter", "confirm"}, {"q", "quit"}}
	case screenError:
		return [][2]string{{"↑↓", "up/down"}, {"enter", "confirm"}, {"esc", "back"}}
	case screenKeys:
		return [][2]string{{"esc", "back"}, {"q", "quit"}}
	default:
		return [][2]string{{"q", "quit"}}
	}
}

func keysFor(s screen) []keyBind {
	nav := []keyBind{
		{"Move up", "↑", "k"},
		{"Move down", "↓", "j"},
	}
	jump := []keyBind{
		{"First item", "Home", ""},
		{"Last item", "End", ""},
		{"Page up", "PgUp", ""},
		{"Page down", "PgDn", ""},
	}
	app := []keyBind{
		{"Back / cancel", "Esc", ""},
		{"Quit", "q", "Ctrl+C"},
		{"Keybindings", "?", ""},
	}
	switch s {
	case screenPathInput:
		return []keyBind{
			{"Confirm", "Enter", ""},
			{"Back", "Esc", ""},
		}
	case screenManage:
		out := append([]keyBind{}, nav...)
		out = append(out,
			keyBind{"Collapse / back", "←", "h"},
			keyBind{"Expand / forward", "→", "l"},
			keyBind{"Select / toggle", "Space", ""},
			keyBind{"Confirm / open", "Enter", ""},
			keyBind{"Next section", "Tab", ""},
			keyBind{"Previous section", "Shift+Tab", ""},
		)
		out = append(out, jump...)
		return append(out, app...)
	case screenVersionList:
		out := append([]keyBind{}, nav...)
		out = append(out, keyBind{"Confirm", "Enter", ""})
		out = append(out, jump...)
		return append(out, app...)
	case screenResult:
		return []keyBind{
			{"Continue", "Enter", ""},
			{"Quit", "q", "Ctrl+C"},
		}
	case screenProgress:
		return []keyBind{{"Quit", "Ctrl+C", "q"}}
	case screenLoading:
		return []keyBind{{"Quit", "q", "Ctrl+C"}}
	default:
		out := append([]keyBind{}, nav...)
		out = append(out, keyBind{"Confirm", "Enter", ""})
		if s != screenSystem && s != screenUnsupported {
			out = append(out, jump...)
		}
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
	b.WriteString(titleStyle.Render("Keybindings") + "\n\n")
	for i, row := range rows {
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
	return row.Key + ", " + row.Also
}
