package content

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const (
	MaxFileSize    = 1 << 20 // 1 MiB per file
	MaxBundleFiles = 500
	MaxBundleBytes = 10 << 20 // 10 MiB
	MaxHashSize    = 5 << 20
)

func SHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// ValidateRelPath requires a clean /-separated path under a known component dir.
func ValidateRelPath(rel string) error {
	if rel == "" {
		return fmt.Errorf("empty path")
	}
	if !utf8.ValidString(rel) {
		return fmt.Errorf("path is not valid UTF-8")
	}
	if strings.ContainsRune(rel, 0) {
		return fmt.Errorf("path contains a null byte")
	}
	if strings.Contains(rel, `\`) {
		return fmt.Errorf("path contains a backslash")
	}
	if filepath.IsAbs(rel) || strings.HasPrefix(rel, "/") || strings.HasPrefix(rel, "//") {
		return fmt.Errorf("absolute path is not allowed")
	}
	if len(rel) >= 2 && rel[1] == ':' {
		return fmt.Errorf("absolute path is not allowed")
	}
	if strings.Contains(rel, ":") {
		return fmt.Errorf("path contains an unexpected colon")
	}
	cleaned := path.Clean(rel)
	if cleaned != rel {
		return fmt.Errorf("path is not clean")
	}
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return fmt.Errorf("path escapes the destination")
	}

	parts := strings.Split(rel, "/")
	if len(parts) < 2 {
		return fmt.Errorf("path must be inside a component directory")
	}
	if !Known(parts[0]) {
		return fmt.Errorf("unknown component %q", parts[0])
	}
	for i, p := range parts {
		if p == "" || p == "." || p == ".." {
			return fmt.Errorf("invalid path component")
		}
		if strings.ContainsAny(p, `<>:"|?*`) {
			return fmt.Errorf("path contains an invalid character")
		}
		if i > 0 && strings.HasPrefix(p, ".") {
			return fmt.Errorf("hidden files are not installed")
		}
	}
	return nil
}

// SafeJoin keeps the result inside root.
func SafeJoin(root, rel string) (string, error) {
	if err := ValidateRelPath(rel); err != nil {
		return "", err
	}
	if root == "" {
		return "", fmt.Errorf("cursor root is empty")
	}
	root = filepath.Clean(root)
	full := filepath.Join(root, filepath.FromSlash(rel))
	full = filepath.Clean(full)
	relToRoot, err := filepath.Rel(root, full)
	if err != nil {
		return "", err
	}
	if relToRoot == ".." || strings.HasPrefix(relToRoot, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes the Cursor directory")
	}
	if filepath.IsAbs(relToRoot) {
		return "", fmt.Errorf("path escapes the Cursor directory")
	}
	return full, nil
}

func ComponentOf(rel string) string {
	rel = strings.ReplaceAll(rel, `\`, "/")
	parts := strings.Split(rel, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}
