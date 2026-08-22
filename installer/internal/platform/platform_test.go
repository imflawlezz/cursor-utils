package platform

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDefaultDir(t *testing.T) {
	t.Parallel()
	dir, ok := DefaultDir("darwin", "/Users/ada")
	if !ok || dir != "/Users/ada/.cursor" {
		t.Fatalf("darwin: %q %v", dir, ok)
	}
	dir, ok = DefaultDir("linux", "/home/ada")
	if !ok || dir != "/home/ada/.cursor" {
		t.Fatalf("linux: %q %v", dir, ok)
	}
	dir, ok = DefaultDir("windows", `C:\Users\ada`)
	if !ok || dir != `C:\Users\ada\.cursor` {
		t.Fatalf("windows: %q %v", dir, ok)
	}
	if _, ok := DefaultDir("plan9", "/home/ada"); ok {
		t.Fatal("plan9 should not invent a default")
	}
	if _, ok := DefaultDir("linux", ""); ok {
		t.Fatal("empty home should not invent a default")
	}
}

func TestDisplayNameAndSupported(t *testing.T) {
	t.Parallel()
	if DisplayName("darwin") != "macOS" {
		t.Fatal(DisplayName("darwin"))
	}
	if !Supported("darwin") || !Supported("linux") || !Supported("windows") {
		t.Fatal("darwin, linux, and windows should be supported")
	}
	if Supported("plan9") {
		t.Fatal("plan9 should not be supported")
	}
}

func TestDetectRuntime(t *testing.T) {
	t.Parallel()
	info, err := Detect()
	if err != nil {
		t.Fatal(err)
	}
	if info.OS != runtime.GOOS {
		t.Fatalf("OS = %s", info.OS)
	}
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" || runtime.GOOS == "windows" {
		if !info.Supported || !info.DefaultOK {
			t.Fatalf("expected supported default: %+v", info)
		}
		if filepath.Base(info.DefaultDir) != ".cursor" {
			t.Fatalf("default dir = %s", info.DefaultDir)
		}
	}
}

func TestExpandUser(t *testing.T) {
	t.Parallel()
	got, err := ExpandUser("~/.cursor", "/Users/ada")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join("/Users/ada", ".cursor")
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestValidateCursorRoot(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	cursor := filepath.Join(root, ".cursor")
	if err := os.Mkdir(cursor, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := ValidateCursorRoot(cursor); err != nil {
		t.Fatal(err)
	}

	nested := filepath.Join(cursor, ".cursor")
	if err := ValidateCursorRoot(nested); err == nil {
		t.Fatal("expected nested .cursor to be rejected")
	}

	commands := filepath.Join(cursor, "commands")
	if err := os.Mkdir(commands, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := ValidateCursorRoot(commands); err == nil {
		t.Fatal("expected component directory to be rejected")
	}

	file := filepath.Join(root, "notdir")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateCursorRoot(file); err == nil {
		t.Fatal("expected file to be rejected")
	}

	missing := filepath.Join(root, "missing-parent", ".cursor")
	if err := ValidateCursorRoot(missing); err == nil {
		t.Fatal("expected missing parent to be rejected")
	}

	creatable := filepath.Join(root, "new-cursor")
	if err := ValidateCursorRoot(creatable); err != nil {
		t.Fatal(err)
	}
}
