package installer

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/imflawlezz/cursor-utils/installer/internal/config"
	"github.com/imflawlezz/cursor-utils/installer/internal/github"
	"github.com/imflawlezz/cursor-utils/installer/internal/platform"
)

func TestLiveGitHubInstallUninstall(t *testing.T) {
	if os.Getenv("CURSOR_UTILS_LIVE") == "" {
		t.Skip("set CURSOR_UTILS_LIVE=1 to run against GitHub")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	cfg := config.Default()
	eng := New(cfg, github.New(cfg), platform.Info{OS: "darwin", Arch: "amd64", Supported: true})
	tags, err := eng.ListTags(ctx)
	if err != nil {
		t.Fatal(err)
	}
	latest, ok := eng.Latest(tags)
	if !ok {
		t.Fatalf("no semver tags in %v", tags)
	}
	bundle, err := eng.Fetch(ctx, latest)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	user := filepath.Join(root, "commands", "my-personal-command.md")
	if err := os.MkdirAll(filepath.Dir(user), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(user, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	plan, err := eng.PlanInstall(root, bundle, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := eng.Apply(ctx, plan, nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "commands", "commit.md")); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(user)
	if string(got) != "keep" {
		t.Fatal("user file overwritten")
	}
	man, err := eng.LoadManifest(root)
	if err != nil || man == nil {
		t.Fatal(err)
	}
	uplan, err := eng.PlanUninstall(root, man)
	if err != nil {
		t.Fatal(err)
	}
	if err := eng.Apply(ctx, uplan, nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "commands", "commit.md")); !os.IsNotExist(err) {
		t.Fatal("owned file remained")
	}
	if _, err := os.Stat(user); err != nil {
		t.Fatal("user file removed")
	}
}
