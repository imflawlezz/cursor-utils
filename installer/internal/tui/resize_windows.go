//go:build windows

package tui

import (
	"fmt"
	"os"
	"time"
	"unsafe"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/sys/windows"
)

var (
	kernel32                       = windows.NewLazySystemDLL("kernel32.dll")
	procSetConsoleScreenBufferSize = kernel32.NewProc("SetConsoleScreenBufferSize")
	procSetConsoleWindowInfo       = kernel32.NewProc("SetConsoleWindowInfo")
)

func sizePollCmd() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(time.Time) tea.Msg {
		return sizePollMsg{}
	})
}

func currentTermSize() (int, int, bool) {
	f, err := os.OpenFile("CONOUT$", os.O_RDWR, 0)
	if err != nil {
		return 0, 0, false
	}
	defer f.Close()
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(windows.Handle(f.Fd()), &info); err != nil {
		return 0, 0, false
	}
	w := int(info.Window.Right - info.Window.Left + 1)
	h := int(info.Window.Bottom - info.Window.Top + 1)
	if w < 1 || h < 1 {
		return 0, 0, false
	}
	return w, h, true
}

func sendTermResize(rows, cols int) {
	if rows < 1 {
		return
	}
	seq := fmt.Sprintf("\x1b[8;%d;t", rows)
	if cols > 0 {
		seq = fmt.Sprintf("\x1b[8;%d;%dt", rows, cols)
	}
	_, _ = os.Stdout.WriteString(seq)
	f, err := os.OpenFile("CONOUT$", os.O_RDWR, 0)
	if err != nil {
		return
	}
	defer f.Close()
	h := windows.Handle(f.Fd())
	enableVT(h)
	_, _ = f.WriteString(seq)
	resizeConsole(h, rows, cols)
}

func enableVT(h windows.Handle) {
	var mode uint32
	if err := windows.GetConsoleMode(h, &mode); err != nil {
		return
	}
	_ = windows.SetConsoleMode(h, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}

func resizeConsole(h windows.Handle, rows, cols int) {
	if cols < 1 {
		cols = minTermCols
	}
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(h, &info); err != nil {
		return
	}
	bufX, bufY := info.Size.X, info.Size.Y
	if int(bufX) < cols {
		bufX = int16(cols)
	}
	if int(bufY) < rows {
		bufY = int16(rows)
	}
	if bufX != info.Size.X || bufY != info.Size.Y {
		if err := setConsoleScreenBufferSize(h, windows.Coord{X: bufX, Y: bufY}); err != nil {
			return
		}
	}
	rect := windows.SmallRect{
		Left:   0,
		Top:    0,
		Right:  int16(cols - 1),
		Bottom: int16(rows - 1),
	}
	_ = setConsoleWindowInfo(h, &rect)
}

func setConsoleScreenBufferSize(h windows.Handle, size windows.Coord) error {
	packed := uintptr(uint32(uint16(size.X)) | uint32(uint16(size.Y))<<16)
	r1, _, err := procSetConsoleScreenBufferSize.Call(uintptr(h), packed)
	if r1 == 0 {
		return err
	}
	return nil
}

func setConsoleWindowInfo(h windows.Handle, rect *windows.SmallRect) error {
	r1, _, err := procSetConsoleWindowInfo.Call(uintptr(h), 1, uintptr(unsafe.Pointer(rect)))
	if r1 == 0 {
		return err
	}
	return nil
}
