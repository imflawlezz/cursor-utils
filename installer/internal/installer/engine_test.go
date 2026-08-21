package installer

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/imflawlezz/cursor-utils/installer/internal/config"
	"github.com/imflawlezz/cursor-utils/installer/internal/content"
	"github.com/imflawlezz/cursor-utils/installer/internal/manifest"
	"github.com/imflawlezz/cursor-utils/installer/internal/platform"
)

type fakeProvider struct {
	tags    []string
	bundles map[string]*content.Bundle
}

func (f fakeProvider) ListTags(ctx context.Context) ([]string, error) {
	return append([]string{}, f.tags...), nil
}

func (f fakeProvider) Fetch(ctx context.Context, tag string) (*content.Bundle, error) {
	b, ok := f.bundles[tag]
	if !ok {
		return nil, os.ErrNotExist
	}
	return b, nil
}

func testEngine(t *testing.T, p content.Provider) *Engine {
	t.Helper()
	return New(config.Default(), p, platform.Info{OS: "darwin", Arch: "arm64", DisplayName: "macOS", Supported: true})
}

func bundle(t *testing.T, tag string, files map[string]string) *content.Bundle {
	t.Helper()
	b := content.NewBundle(tag)
	for p, body := range files {
		if err := b.Add(p, []byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	return b
}

func TestInstallUpdateUninstallProtectsUserFiles(t *testing.T) {
	root := t.TempDir()
	v1 := bundle(t, "v0.1.0", map[string]string{
		"commands/commit.md": "commit-v1",
		"commands/docs.md":   "docs-v1",
	})
	eng := testEngine(t, fakeProvider{})
	plan, err := eng.PlanInstall(root, v1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := eng.Apply(context.Background(), plan, nil, nil); err != nil {
		t.Fatal(err)
	}

	user := filepath.Join(root, "commands", "my-personal-command.md")
	if err := os.WriteFile(user, []byte("mine"), 0o644); err != nil {
		t.Fatal(err)
	}

	man, err := eng.LoadManifest(root)
	if err != nil || man == nil {
		t.Fatalf("manifest: %v %v", man, err)
	}
	if man.Owns("commands/my-personal-command.md") {
		t.Fatal("must not claim user file")
	}

	v2 := bundle(t, "v0.2.0", map[string]string{
		"commands/commit.md": "commit-v2",
		"commands/readme.md": "readme-v2",
	})
	plan, err = eng.PlanInstall(root, v2, man)
	if err != nil {
		t.Fatal(err)
	}
	kinds := opKinds(plan)
	if kinds["commands/commit.md"] != OpUpdate {
		t.Fatalf("commit: %v", kinds)
	}
	if kinds["commands/readme.md"] != OpAdd {
		t.Fatalf("readme: %v", kinds)
	}
	if kinds["commands/docs.md"] != OpRemove {
		t.Fatalf("docs should be removed: %v", kinds)
	}

	if err := eng.Apply(context.Background(), plan, nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(user); err != nil {
		t.Fatal("user file was touched")
	}
	if _, err := os.Stat(filepath.Join(root, "commands", "docs.md")); !os.IsNotExist(err) {
		t.Fatal("obsolete owned file should be removed")
	}
	got, _ := os.ReadFile(filepath.Join(root, "commands", "commit.md"))
	if string(got) != "commit-v2" {
		t.Fatalf("commit = %s", got)
	}

	man, err = eng.LoadManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	uplan, err := eng.PlanUninstall(root, man)
	if err != nil {
		t.Fatal(err)
	}
	if err := eng.Apply(context.Background(), uplan, nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(user); err != nil {
		t.Fatal("user file must survive uninstall")
	}
	if _, err := os.Stat(filepath.Join(root, "commands", "commit.md")); !os.IsNotExist(err) {
		t.Fatal("owned file should be gone")
	}
	if _, err := os.Stat(config.ManifestPath(root)); !os.IsNotExist(err) {
		t.Fatal("manifest should be removed")
	}
}

func TestConflictSkipAndOverwrite(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "commands"), 0o755); err != nil {
		t.Fatal(err)
	}
	existing := filepath.Join(root, "commands", "docs.md")
	if err := os.WriteFile(existing, []byte("user-docs"), 0o644); err != nil {
		t.Fatal(err)
	}
	eng := testEngine(t, fakeProvider{})
	b := bundle(t, "v0.1.0", map[string]string{
		"commands/docs.md":   "ours",
		"commands/commit.md": "commit",
	})
	plan, err := eng.PlanInstall(root, b, nil)
	if err != nil {
		t.Fatal(err)
	}
	if opKinds(plan)["commands/docs.md"] != OpConflict {
		t.Fatal(opKinds(plan))
	}

	if err := eng.Apply(context.Background(), plan, map[string]Choice{
		"commands/docs.md": ChoiceSkip,
	}, nil); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(existing)
	if string(got) != "user-docs" {
		t.Fatal("skipped file was overwritten")
	}
	man, err := eng.LoadManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if man.Owns("commands/docs.md") {
		t.Fatal("skip must not take ownership")
	}
	if !man.Owns("commands/commit.md") {
		t.Fatal("commit should be owned")
	}

	plan, err = eng.PlanInstall(root, b, man)
	if err != nil {
		t.Fatal(err)
	}
	if err := eng.Apply(context.Background(), plan, map[string]Choice{
		"commands/docs.md": ChoiceOverwrite,
	}, nil); err != nil {
		t.Fatal(err)
	}
	got, _ = os.ReadFile(existing)
	if string(got) != "ours" {
		t.Fatal("overwrite did not replace")
	}
	man, _ = eng.LoadManifest(root)
	rec, ok := man.File("commands/docs.md")
	if !ok || !rec.TakenOver {
		t.Fatal("overwrite should record takeover")
	}
	backups, err := os.ReadDir(config.BackupRoot(root))
	if err != nil || len(backups) == 0 {
		t.Fatal("expected a backup")
	}
}

func TestModifiedOwnedFileKeepAndReplace(t *testing.T) {
	root := t.TempDir()
	eng := testEngine(t, fakeProvider{})
	v1 := bundle(t, "v0.1.0", map[string]string{"commands/docs.md": "orig"})
	plan, _ := eng.PlanInstall(root, v1, nil)
	if err := eng.Apply(context.Background(), plan, nil, nil); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "commands", "docs.md")
	if err := os.WriteFile(path, []byte("edited"), 0o644); err != nil {
		t.Fatal(err)
	}
	man, _ := eng.LoadManifest(root)
	v2 := bundle(t, "v0.2.0", map[string]string{"commands/docs.md": "remote"})
	plan, err := eng.PlanInstall(root, v2, man)
	if err != nil {
		t.Fatal(err)
	}
	if opKinds(plan)["commands/docs.md"] != OpModifiedUpdate {
		t.Fatal(opKinds(plan))
	}

	if err := eng.Apply(context.Background(), plan, map[string]Choice{
		"commands/docs.md": ChoiceKeepLocal,
	}, nil); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != "edited" {
		t.Fatal("keep local overwrote")
	}
	man, _ = eng.LoadManifest(root)
	rec, _ := man.File("commands/docs.md")
	if rec.SHA256 != content.SHA256([]byte("edited")) {
		t.Fatal("manifest hash should match kept local file")
	}

	if err := os.WriteFile(path, []byte("edited-again"), 0o644); err != nil {
		t.Fatal(err)
	}
	man, _ = eng.LoadManifest(root)
	plan, _ = eng.PlanInstall(root, v2, man)
	if err := eng.Apply(context.Background(), plan, map[string]Choice{
		"commands/docs.md": ChoiceOverwrite,
	}, nil); err != nil {
		t.Fatal(err)
	}
	got, _ = os.ReadFile(path)
	if string(got) != "remote" {
		t.Fatal("replace did not apply")
	}
}

func TestFailedInstallRollsBack(t *testing.T) {
	root := t.TempDir()
	eng := testEngine(t, fakeProvider{})
	v1 := bundle(t, "v0.1.0", map[string]string{"commands/commit.md": "keep-me"})
	plan, _ := eng.PlanInstall(root, v1, nil)
	if err := eng.Apply(context.Background(), plan, nil, nil); err != nil {
		t.Fatal(err)
	}
	man, _ := eng.LoadManifest(root)
	v2 := bundle(t, "v0.2.0", map[string]string{
		"commands/commit.md": "new-commit",
		"commands/docs.md":   "new-docs",
	})
	plan, err := eng.PlanInstall(root, v2, man)
	if err != nil {
		t.Fatal(err)
	}
	blocker := filepath.Join(root, "commands", "docs.md")
	if err := os.Mkdir(blocker, 0o755); err != nil {
		t.Fatal(err)
	}
	err = eng.Apply(context.Background(), plan, nil, nil)
	if err == nil {
		t.Fatal("expected failure")
	}
	got, _ := os.ReadFile(filepath.Join(root, "commands", "commit.md"))
	if string(got) != "keep-me" {
		t.Fatalf("commit was not rolled back: %s", got)
	}
	man, err = eng.LoadManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if man.ContentVersion != "v0.1.0" {
		t.Fatalf("manifest should still be v0.1.0, got %s", man.ContentVersion)
	}
}

func TestChecksumVerificationAndOwnership(t *testing.T) {
	root := t.TempDir()
	eng := testEngine(t, fakeProvider{})
	body := "hello"
	b := bundle(t, "v0.1.0", map[string]string{"commands/commit.md": body})
	plan, _ := eng.PlanInstall(root, b, nil)
	if err := eng.Apply(context.Background(), plan, nil, nil); err != nil {
		t.Fatal(err)
	}
	man, _ := eng.LoadManifest(root)
	rec, _ := man.File("commands/commit.md")
	if rec.SHA256 != content.SHA256([]byte(body)) {
		t.Fatal(rec.SHA256)
	}
	plan, _ = eng.PlanInstall(root, b, man)
	if opKinds(plan)["commands/commit.md"] != OpUnchanged {
		t.Fatal(opKinds(plan))
	}
}

func TestCancelChoiceDoesNotModify(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "commands"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "commands", "docs.md"), []byte("user"), 0o644); err != nil {
		t.Fatal(err)
	}
	eng := testEngine(t, fakeProvider{})
	b := bundle(t, "v0.1.0", map[string]string{"commands/docs.md": "ours"})
	plan, _ := eng.PlanInstall(root, b, nil)
	err := eng.Apply(context.Background(), plan, map[string]Choice{"commands/docs.md": ChoiceCancel}, nil)
	if err == nil {
		t.Fatal("expected cancel")
	}
	got, _ := os.ReadFile(filepath.Join(root, "commands", "docs.md"))
	if string(got) != "user" {
		t.Fatal("cancel modified the file")
	}
	if _, err := os.Stat(config.ManifestPath(root)); !os.IsNotExist(err) {
		t.Fatal("manifest should not exist")
	}
}

func TestListTagsUsesProvider(t *testing.T) {
	eng := testEngine(t, fakeProvider{tags: []string{"v0.1.0", "v0.2.0", "nightly", "installer-v1.0.0", "0.9.0"}})
	tags, err := eng.ListTags(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 2 || tags[0] != "v0.2.0" || tags[1] != "v0.1.0" {
		t.Fatalf("content tags = %v", tags)
	}
	latest, ok := eng.Latest(tags)
	if !ok || latest != "v0.2.0" {
		t.Fatalf("latest = %s %v tags=%v", latest, ok, tags)
	}
}

func TestFetchRejectsInstallerTag(t *testing.T) {
	eng := testEngine(t, fakeProvider{
		bundles: map[string]*content.Bundle{
			"installer-v1.0.0": bundle(t, "installer-v1.0.0", map[string]string{"commands/commit.md": "x"}),
		},
	})
	if _, err := eng.Fetch(context.Background(), "installer-v1.0.0"); err == nil {
		t.Fatal("expected installer tag to be rejected")
	}
	if _, err := eng.Fetch(context.Background(), "nightly"); err == nil {
		t.Fatal("expected non-content tag to be rejected")
	}
}

func TestPartialInstallAndUninstallKeepUnselected(t *testing.T) {
	root := t.TempDir()
	eng := testEngine(t, fakeProvider{})
	v1 := bundle(t, "v0.1.0", map[string]string{
		"commands/commit.md": "c1",
		"commands/docs.md":   "d1",
	})
	plan, _ := eng.PlanInstall(root, v1, nil)
	if err := eng.Apply(context.Background(), plan, nil, nil); err != nil {
		t.Fatal(err)
	}
	v2 := bundle(t, "v0.2.0", map[string]string{
		"commands/commit.md": "c2",
		"commands/docs.md":   "d2",
	})
	man, _ := eng.LoadManifest(root)
	full, err := eng.PlanInstall(root, v2, man)
	if err != nil {
		t.Fatal(err)
	}
	partial := full.Filter(Selection{Files: map[string]bool{"commands/commit.md": true}})
	if err := eng.Apply(context.Background(), partial, nil, nil); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(filepath.Join(root, "commands", "commit.md"))
	if string(got) != "c2" {
		t.Fatalf("commit = %s", got)
	}
	got, _ = os.ReadFile(filepath.Join(root, "commands", "docs.md"))
	if string(got) != "d1" {
		t.Fatalf("unselected docs overwritten: %s", got)
	}
	man, _ = eng.LoadManifest(root)
	if !man.Owns("commands/docs.md") || !man.Owns("commands/commit.md") {
		t.Fatal("ownership dropped for unselected file")
	}

	uplan, err := eng.PlanUninstall(root, man)
	if err != nil {
		t.Fatal(err)
	}
	uplan = uplan.Filter(Selection{Files: map[string]bool{"commands/commit.md": true}})
	if err := eng.Apply(context.Background(), uplan, nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "commands", "commit.md")); !os.IsNotExist(err) {
		t.Fatal("selected file should be gone")
	}
	if _, err := os.Stat(filepath.Join(root, "commands", "docs.md")); err != nil {
		t.Fatal("unselected file should remain")
	}
	man, _ = eng.LoadManifest(root)
	if man.Owns("commands/commit.md") || !man.Owns("commands/docs.md") {
		t.Fatal(man.AllFiles())
	}
}

func opKinds(plan *Plan) map[string]OpKind {
	out := map[string]OpKind{}
	for _, op := range plan.Ops {
		out[op.RelPath] = op.Kind
	}
	return out
}

func TestManifestCorruptionSurfaces(t *testing.T) {
	root := t.TempDir()
	path := config.ManifestPath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{nope"), 0o644); err != nil {
		t.Fatal(err)
	}
	eng := testEngine(t, fakeProvider{})
	_, err := eng.LoadManifest(root)
	if err == nil {
		t.Fatal("expected corruption error")
	}
	var fe *manifest.FormatError
	if !isFormat(err) {
		t.Fatalf("%T %v", err, err)
	}
	_ = fe
}

func isFormat(err error) bool {
	_, ok := err.(*manifest.FormatError)
	return ok
}
