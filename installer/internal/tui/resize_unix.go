//go:build unix

package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func sizePollCmd() tea.Cmd { return nil }

func currentTermSize() (int, int, bool) { return 0, 0, false }

func sendTermResize(rows, cols int) {
	f, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()
	if cols > 0 {
		fmt.Fprintf(f, "\x1b[8;%d;%dt", rows, cols)
		return
	}
	fmt.Fprintf(f, "\x1b[8;%d;t", rows)
}
