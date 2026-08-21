package installer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/imflawlezz/cursor-utils/installer/internal/config"
	"github.com/imflawlezz/cursor-utils/installer/internal/content"
	"github.com/imflawlezz/cursor-utils/installer/internal/manifest"
	"github.com/imflawlezz/cursor-utils/installer/internal/platform"
	"github.com/imflawlezz/cursor-utils/installer/internal/semver"
)

type Engine struct {
	cfg      config.App
	provider content.Provider
	plat     platform.Info
}

func New(cfg config.App, provider content.Provider, plat platform.Info) *Engine {
	return &Engine{cfg: cfg, provider: provider, plat: plat}
}

func (e *Engine) Config() config.App { return e.cfg }

func (e *Engine) ListTags(ctx context.Context) ([]string, error) {
	tags, err := e.provider.ListTags(ctx)
	if err != nil {
		return nil, err
	}
	return semver.SortNewestFirst(semver.ContentTags(tags)), nil
}

func (e *Engine) Latest(tags []string) (string, bool) {
	return semver.Latest(semver.ContentTags(tags))
}

func (e *Engine) Fetch(ctx context.Context, tag string) (*content.Bundle, error) {
	if !semver.IsContentTag(tag) {
		return nil, fmt.Errorf("%s is not a content version (content tags look like v0.1.0)", tag)
	}
	return e.provider.Fetch(ctx, tag)
}

func (e *Engine) LoadManifest(root string) (*manifest.Manifest, error) {
	return manifest.Load(config.ManifestPath(root))
}

type EventStatus int

const (
	EventStarted EventStatus = iota
	EventDone
	EventSkipped
	EventError
)

type Event struct {
	Path   string
	Status EventStatus
	Err    error
}

func (e *Engine) Apply(ctx context.Context, plan *Plan, choices map[string]Choice, report func(Event)) error {
	if plan == nil {
		return fmt.Errorf("installation plan is missing")
	}
	if err := platform.ValidateCursorRoot(plan.CursorRoot); err != nil {
		return err
	}
	resolved, err := resolveOps(plan, choices)
	if err != nil {
		return err
	}

	tx := newTxn()
	var applyErr error
	defer func() {
		if applyErr != nil {
			tx.rollback()
		}
	}()

	newMan := e.buildResultManifest(plan, resolved)

	for _, step := range resolved {
		if err := ctx.Err(); err != nil {
			applyErr = fmt.Errorf("installation was cancelled")
			return applyErr
		}
		emit(report, Event{Path: step.RelPath, Status: EventStarted})
		if err := e.applyStep(tx, plan, step); err != nil {
			emit(report, Event{Path: step.RelPath, Status: EventError, Err: err})
			applyErr = err
			return applyErr
		}
		if step.skip {
			emit(report, Event{Path: step.RelPath, Status: EventSkipped})
		} else {
			emit(report, Event{Path: step.RelPath, Status: EventDone})
		}
	}

	manPath := config.ManifestPath(plan.CursorRoot)
	if newMan.Empty() {
		if err := ctx.Err(); err != nil {
			applyErr = fmt.Errorf("installation was cancelled")
			return applyErr
		}
		if err := manifest.Remove(manPath); err != nil {
			applyErr = err
			return applyErr
		}
		_ = os.Remove(config.ManifestDirPath(plan.CursorRoot))
		return nil
	}

	if err := manifest.WriteAtomic(manPath, newMan); err != nil {
		applyErr = err
		return applyErr
	}
	return nil
}

type resolvedOp struct {
	Op
	skip       bool
	takeOver   bool
	keepLocal  bool
	removeFile bool
	writeFile  bool
}

func resolveOps(plan *Plan, choices map[string]Choice) ([]resolvedOp, error) {
	out := make([]resolvedOp, 0, len(plan.Ops))
	for _, op := range plan.Ops {
		r := resolvedOp{Op: op}
		choice := choices[op.RelPath]
		if op.Kind.NeedsDecision() {
			switch choice {
			case ChoiceCancel:
				return nil, fmt.Errorf("installation cancelled")
			case ChoiceNone:
				return nil, fmt.Errorf("a decision is required for %s", op.RelPath)
			case ChoiceSkip:
				if op.Kind == OpConflict {
					r.skip = true
				} else {
					return nil, fmt.Errorf("invalid choice for %s", op.RelPath)
				}
			case ChoiceKeepLocal:
				if op.Kind == OpModifiedUpdate {
					r.skip = true
					r.keepLocal = true
				} else if op.Kind == OpModifiedRemove {
					r.skip = true
					r.keepLocal = true
				} else {
					return nil, fmt.Errorf("invalid choice for %s", op.RelPath)
				}
			case ChoiceOverwrite:
				switch op.Kind {
				case OpConflict:
					r.writeFile = true
					r.takeOver = true
				case OpModifiedUpdate:
					r.writeFile = true
				case OpModifiedRemove:
					r.removeFile = true
				default:
					return nil, fmt.Errorf("invalid choice for %s", op.RelPath)
				}
			default:
				return nil, fmt.Errorf("unknown choice for %s", op.RelPath)
			}
		} else {
			switch op.Kind {
			case OpAdd, OpUpdate:
				r.writeFile = true
			case OpRemove:
				r.removeFile = true
			case OpUnchanged:
				r.skip = true
			}
		}
		out = append(out, r)
	}
	return out, nil
}

func (e *Engine) applyStep(tx *txn, plan *Plan, step resolvedOp) error {
	full, err := content.SafeJoin(plan.CursorRoot, step.RelPath)
	if err != nil {
		return err
	}
	if step.writeFile {
		if step.takeOver && step.Exists {
			if err := backupFile(plan.CursorRoot, step.RelPath, full); err != nil {
				return err
			}
		}
		return tx.writeFile(full, step.Data)
	}
	if step.removeFile && step.Exists {
		if err := ensureStillOwned(full, step.Op); err != nil {
			return err
		}
		return tx.removeFile(full)
	}
	return nil
}

func ensureStillOwned(full string, op Op) error {
	st, err := os.Lstat(full)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if st.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing to remove symlink %s", op.RelPath)
	}
	if st.IsDir() {
		return fmt.Errorf("refusing to remove directory %s", op.RelPath)
	}
	if op.ManifestSHA == "" {
		return nil
	}
	got, err := hashFile(full)
	if err != nil {
		return err
	}
	if got != op.ManifestSHA && got != op.DiskSHA {
		return fmt.Errorf("%s changed during installation; aborting", op.RelPath)
	}
	return nil
}

func backupFile(root, rel, full string) error {
	if err := content.ValidateRelPath(rel); err != nil {
		return err
	}
	stamp := time.Now().UTC().Format("20060102T150405Z")
	base := filepath.Clean(config.BackupRoot(root))
	dest := filepath.Clean(filepath.Join(base, stamp, filepath.FromSlash(rel)))
	relToBase, err := filepath.Rel(base, dest)
	if err != nil || relToBase == ".." || strings.HasPrefix(relToBase, ".."+string(filepath.Separator)) {
		return fmt.Errorf("backup path for %s is invalid", rel)
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("unable to create a backup of %s: %w", rel, err)
	}
	data, err := os.ReadFile(full)
	if err != nil {
		return fmt.Errorf("unable to backup %s: %w", rel, err)
	}
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return fmt.Errorf("unable to backup %s: %w", rel, err)
	}
	return nil
}

func (e *Engine) buildResultManifest(plan *Plan, steps []resolvedOp) *manifest.Manifest {
	var m *manifest.Manifest
	if plan.Existing != nil {
		m = plan.Existing.Clone()
	} else {
		m = manifest.New()
	}
	m.InstallerVersion = e.cfg.Version
	m.Repository = e.cfg.RepositoryURL()
	m.InstalledAt = time.Now().UTC()
	m.Platform = manifest.Platform{OS: runtime.GOOS, Arch: runtime.GOARCH}
	if plan.Platform.OS != "" {
		m.Platform = manifest.Platform{OS: plan.Platform.OS, Arch: plan.Platform.Arch}
	}
	m.CursorRoot = plan.CursorRoot
	if plan.Kind != PlanUninstall && plan.ContentVersion != "" {
		m.ContentVersion = plan.ContentVersion
	}

	for _, step := range steps {
		switch {
		case step.removeFile:
			m.RemoveFile(step.RelPath)
		case step.Kind == OpConflict && step.skip:
			continue
		case step.Kind == OpModifiedRemove && step.keepLocal:
			m.RemoveFile(step.RelPath)
		default:
			rec := manifest.FileRecord{
				Path:      step.RelPath,
				TakenOver: step.TakenOver || step.takeOver,
			}
			switch {
			case step.writeFile:
				rec.SHA256 = step.NewSHA
				if rec.SHA256 == "" && step.Data != nil {
					rec.SHA256 = content.SHA256(step.Data)
				}
			case step.keepLocal || step.Kind == OpUnchanged || step.skip:
				if step.DiskSHA != "" {
					rec.SHA256 = step.DiskSHA
				} else {
					rec.SHA256 = step.ManifestSHA
				}
			default:
				rec.SHA256 = step.NewSHA
			}
			m.PutFile(step.Component, rec)
			if plan.Kind != PlanUninstall {
				comp := m.Components[step.Component]
				comp.Version = plan.ContentVersion
				m.Components[step.Component] = comp
			}
		}
	}
	return m
}

func emit(report func(Event), ev Event) {
	if report != nil {
		report(ev)
	}
}
