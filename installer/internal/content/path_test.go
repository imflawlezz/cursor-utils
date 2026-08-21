package content

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateRelPath(t *testing.T) {
	t.Parallel()
	ok := []string{
		"commands/commit.md",
		"commands/foo..md",
		"rules/core.mdc",
		"skills/foo/bar.md",
		"agents/helper.md",
		"hooks/event.json",
	}
	for _, p := range ok {
		if err := ValidateRelPath(p); err != nil {
			t.Errorf("%s: %v", p, err)
		}
	}
	bad := []string{
		"",
		"commands",
		"commands/",
		"/etc/passwd",
		"C:/Windows/x.md",
		"commands/../etc/passwd",
		"commands/foo/../../../etc/passwd",
		"../commands/x.md",
		"commands/foo/../../x.md",
		`commands\x.md`,
		"commands/.hidden.md",
		"commands/sub/.env",
		"README.md",
		"commands/foo/./bar.md",
		"..",
		"commands/foo.md\x00",
	}
	for _, p := range bad {
		if err := ValidateRelPath(p); err == nil {
			t.Errorf("expected reject: %q", p)
		}
	}
}

func TestSafeJoinRejectsEscape(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	full, err := SafeJoin(root, "commands/commit.md")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, "commands", "commit.md")
	if full != want {
		t.Fatalf("got %s want %s", full, want)
	}
	if !strings.HasPrefix(full, root) {
		t.Fatal("joined path escaped root")
	}
	if _, err := SafeJoin(root, "commands/../../outside.md"); err == nil {
		t.Fatal("expected escape to be rejected")
	}
	if _, err := SafeJoin(root, "/tmp/x.md"); err == nil {
		t.Fatal("expected absolute path to be rejected")
	}
}

func TestSHA256Stable(t *testing.T) {
	t.Parallel()
	a := SHA256([]byte("hello"))
	b := SHA256([]byte("hello"))
	c := SHA256([]byte("world"))
	if a != b || a == c || len(a) != 64 {
		t.Fatalf("sha mismatch %s %s %s", a, b, c)
	}
}
