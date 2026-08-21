package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Info struct {
	OS          string
	DisplayName string
	Arch        string
	Home        string
	DefaultDir  string
	DefaultOK   bool
	Supported   bool
}

func Detect() (Info, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}
	dir, ok := DefaultDir(runtime.GOOS, home)
	if ok {
		dir = nativeClean(runtime.GOOS, dir)
	}
	return Info{
		OS:          runtime.GOOS,
		DisplayName: DisplayName(runtime.GOOS),
		Arch:        runtime.GOARCH,
		Home:        home,
		DefaultDir:  dir,
		DefaultOK:   ok && home != "",
		Supported:   Supported(runtime.GOOS),
	}, nil
}

func Supported(goos string) bool {
	return goos == "darwin" || goos == "linux"
}

func DisplayName(goos string) string {
	switch goos {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return goos
	}
}

// DefaultDir returns ok=false instead of guessing when goos is unknown or home is empty.
func DefaultDir(goos, home string) (string, bool) {
	home = strings.TrimSpace(home)
	if home == "" {
		return "", false
	}
	switch goos {
	case "darwin", "linux":
		return unixJoin(home, ".cursor"), true
	case "windows":
		// Path shape only; Supported("windows") is still false.
		return windowsJoin(home, ".cursor"), true
	default:
		return "", false
	}
}

func unixJoin(home, elem string) string {
	return strings.TrimRight(home, "/") + "/" + elem
}

func windowsJoin(home, elem string) string {
	return strings.TrimRight(home, `\/`) + `\` + elem
}

func nativeClean(goos, path string) string {
	if goos == "windows" {
		return path
	}
	return filepath.Clean(path)
}

// ExpandUser resolves a leading ~. Relative paths use the process cwd, not the Cursor root.
func ExpandUser(path, home string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}
	if path == "~" {
		if home == "" {
			return "", fmt.Errorf("home directory is unknown")
		}
		return filepath.Clean(home), nil
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		if home == "" {
			return "", fmt.Errorf("home directory is unknown")
		}
		return filepath.Clean(filepath.Join(home, path[2:])), nil
	}
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return abs, nil
}
