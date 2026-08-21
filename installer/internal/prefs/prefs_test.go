package prefs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	home := t.TempDir()
	got, err := Load(home)
	if err != nil || got != nil {
		t.Fatalf("empty home: %v %v", got, err)
	}
	if err := Save(home, File{CursorRoot: "/tmp/cursor"}); err != nil {
		t.Fatal(err)
	}
	got, err = Load(home)
	if err != nil {
		t.Fatal(err)
	}
	if got.CursorRoot != "/tmp/cursor" {
		t.Fatal(got.CursorRoot)
	}
	if filepath.Base(Path(home)) != ".cursor-utils.json" {
		t.Fatal(Path(home))
	}
}

func TestLoadCorrupt(t *testing.T) {
	home := t.TempDir()
	if err := os.WriteFile(Path(home), []byte("nope"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(home); err == nil {
		t.Fatal("expected error")
	}
}
