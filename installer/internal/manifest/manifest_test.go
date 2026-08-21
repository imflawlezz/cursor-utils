package manifest

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRoundTrip(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, ".cursor-utils", "manifest.json")
	m := New()
	m.InstallerVersion = "0.1.0"
	m.Repository = "https://github.com/imflawlezz/cursor-utils"
	m.ContentVersion = "v0.1.0"
	m.InstalledAt = time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC)
	m.Platform = Platform{OS: "darwin", Arch: "arm64"}
	m.CursorRoot = dir
	m.PutFile("commands", FileRecord{Path: "commands/commit.md", SHA256: "abc"})

	if err := WriteAtomic(path, m); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.ContentVersion != "v0.1.0" || !got.Owns("commands/commit.md") {
		t.Fatalf("%+v", got)
	}
	if rec, ok := got.File("commands/commit.md"); !ok || rec.SHA256 != "abc" {
		t.Fatalf("%+v", rec)
	}
	if got.Owns("commands/personal.md") {
		t.Fatal("must not own unrelated files")
	}
}

func TestLoadMissing(t *testing.T) {
	t.Parallel()
	got, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil || got != nil {
		t.Fatalf("got %v %v", got, err)
	}
}

func TestParseCorruptionAndVersions(t *testing.T) {
	t.Parallel()
	if _, err := Parse([]byte("not json")); err == nil {
		t.Fatal("expected corrupt json error")
	}
	if _, err := Parse([]byte(`{"formatVersion":99}`)); err == nil {
		t.Fatal("expected future version error")
	}
	if _, err := Parse([]byte(`{"formatVersion":0}`)); err == nil {
		t.Fatal("expected missing version error")
	}
	if _, err := Parse([]byte(`{"formatVersion":1,"components":null}`)); err != nil {
		t.Fatal(err)
	}
}

func TestWriteAtomicReplacesWithoutTruncatingFirst(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.json")
	first := New()
	first.ContentVersion = "v0.1.0"
	if err := WriteAtomic(path, first); err != nil {
		t.Fatal(err)
	}
	second := New()
	second.ContentVersion = "v0.2.0"
	if err := WriteAtomic(path, second); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if got.ContentVersion != "v0.2.0" {
		t.Fatal(got.ContentVersion)
	}
}

func TestRemoveFileAndEmpty(t *testing.T) {
	t.Parallel()
	m := New()
	m.PutFile("commands", FileRecord{Path: "commands/a.md", SHA256: "1"})
	m.PutFile("commands", FileRecord{Path: "commands/b.md", SHA256: "2"})
	m.RemoveFile("commands/a.md")
	if m.Owns("commands/a.md") || !m.Owns("commands/b.md") {
		t.Fatal(m.AllFiles())
	}
	m.RemoveFile("commands/b.md")
	if !m.Empty() {
		t.Fatal("expected empty")
	}
}
