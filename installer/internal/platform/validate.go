package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var componentDirNames = map[string]struct{}{
	"commands": {},
	"rules":    {},
	"skills":   {},
	"agents":   {},
	"hooks":    {},
}

// ValidateCursorRoot never creates directories.
func ValidateCursorRoot(path string) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path must be absolute")
	}
	cleaned := filepath.Clean(path)
	if err := rejectNestedCursor(cleaned); err != nil {
		return err
	}
	base := filepath.Base(cleaned)
	if _, ok := componentDirNames[base]; ok {
		return fmt.Errorf("%q looks like a Cursor %s directory; choose the Cursor root (usually ~/.cursor) instead", cleaned, base)
	}

	info, err := os.Lstat(cleaned)
	if err != nil {
		if os.IsNotExist(err) {
			parent := filepath.Dir(cleaned)
			pInfo, pErr := os.Stat(parent)
			if pErr != nil {
				if os.IsNotExist(pErr) {
					return fmt.Errorf("parent directory does not exist: %s", parent)
				}
				return fmt.Errorf("cannot access parent directory: %w", pErr)
			}
			if !pInfo.IsDir() {
				return fmt.Errorf("parent path is not a directory: %s", parent)
			}
			return nil
		}
		return fmt.Errorf("cannot access path: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		resolved, err := filepath.EvalSymlinks(cleaned)
		if err != nil {
			return fmt.Errorf("cannot resolve symlink: %w", err)
		}
		if err := ValidateCursorRoot(resolved); err != nil {
			return err
		}
		return nil
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", cleaned)
	}
	return nil
}

func rejectNestedCursor(path string) error {
	n := 0
	for _, elem := range splitPath(path) {
		if elem == ".cursor" {
			n++
		}
	}
	if n > 1 {
		return fmt.Errorf("refusing a .cursor directory nested inside another .cursor directory")
	}
	return nil
}

func splitPath(path string) []string {
	path = strings.ReplaceAll(path, `\`, "/")
	var out []string
	for _, e := range strings.Split(path, "/") {
		if e == "" || e == "." {
			continue
		}
		out = append(out, e)
	}
	return out
}

func LooksLikeComponentDir(path string) bool {
	_, ok := componentDirNames[filepath.Base(filepath.Clean(path))]
	return ok
}
