package installer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/imflawlezz/cursor-utils/installer/internal/content"
	"github.com/imflawlezz/cursor-utils/installer/internal/manifest"
)

func (e *Engine) PlanInstall(root string, bundle *content.Bundle, existing *manifest.Manifest) (*Plan, error) {
	if bundle == nil || bundle.Empty() {
		return nil, fmt.Errorf("this version contains no installable Cursor content")
	}
	plan := &Plan{
		Kind:           PlanInstall,
		CursorRoot:     root,
		ContentVersion: bundle.Version,
		Platform:       e.plat,
		Existing:       existing,
		Bundle:         bundle,
	}

	seen := make(map[string]struct{})
	for _, file := range bundle.AllFiles() {
		seen[file.RelPath] = struct{}{}
		op, err := e.planIncoming(root, file, existing)
		if err != nil {
			return nil, err
		}
		plan.Ops = append(plan.Ops, op)
	}
	if existing != nil {
		for _, rec := range existing.AllFiles() {
			if _, ok := seen[rec.Path]; ok {
				continue
			}
			op, err := e.planObsolete(root, rec)
			if err != nil {
				return nil, err
			}
			plan.Ops = append(plan.Ops, op)
		}
	}
	return plan, nil
}

func (e *Engine) PlanUninstall(root string, existing *manifest.Manifest) (*Plan, error) {
	if existing == nil || existing.Empty() {
		return nil, fmt.Errorf("no cursor-utils content is installed in this directory")
	}
	plan := &Plan{
		Kind:           PlanUninstall,
		CursorRoot:     root,
		ContentVersion: existing.ContentVersion,
		Platform:       e.plat,
		Existing:       existing,
	}
	for _, rec := range existing.AllFiles() {
		op, err := e.planObsolete(root, rec)
		if err != nil {
			return nil, err
		}
		plan.Ops = append(plan.Ops, op)
	}
	return plan, nil
}

func (e *Engine) planIncoming(root string, file content.File, existing *manifest.Manifest) (Op, error) {
	op := Op{
		Kind:      OpAdd,
		Component: content.ComponentOf(file.RelPath),
		RelPath:   file.RelPath,
		Data:      file.Data,
		NewSHA:    file.SHA256,
	}
	full, err := content.SafeJoin(root, file.RelPath)
	if err != nil {
		return Op{}, err
	}
	st, err := os.Lstat(full)
	if err != nil {
		if os.IsNotExist(err) {
			if rec, ok := existing.File(file.RelPath); ok {
				op.ManifestSHA = rec.SHA256
				op.TakenOver = rec.TakenOver
			}
			op.Kind = OpAdd
			return op, nil
		}
		return Op{}, fmt.Errorf("cannot inspect %s: %w", file.RelPath, err)
	}
	if st.Mode()&os.ModeSymlink != 0 {
		return Op{}, fmt.Errorf("%s is a symlink; remove it manually before installing", file.RelPath)
	}
	if st.IsDir() {
		return Op{}, fmt.Errorf("%s exists as a directory", file.RelPath)
	}
	diskSHA, err := hashFile(full)
	if err != nil {
		return Op{}, err
	}
	op.Exists = true
	op.DiskSHA = diskSHA
	rec, owned := existing.File(file.RelPath)
	if !owned {
		op.Kind = OpConflict
		return op, nil
	}
	op.ManifestSHA = rec.SHA256
	op.TakenOver = rec.TakenOver
	if diskSHA != rec.SHA256 {
		op.Kind = OpModifiedUpdate
		return op, nil
	}
	if diskSHA == file.SHA256 {
		op.Kind = OpUnchanged
		return op, nil
	}
	op.Kind = OpUpdate
	return op, nil
}

func (e *Engine) planObsolete(root string, rec manifest.FileRecord) (Op, error) {
	op := Op{
		Kind:        OpRemove,
		Component:   content.ComponentOf(rec.Path),
		RelPath:     rec.Path,
		ManifestSHA: rec.SHA256,
		TakenOver:   rec.TakenOver,
	}
	if err := content.ValidateRelPath(rec.Path); err != nil {
		return Op{}, fmt.Errorf("manifest path %q is invalid: %w", rec.Path, err)
	}
	full, err := content.SafeJoin(root, rec.Path)
	if err != nil {
		return Op{}, err
	}
	st, err := os.Lstat(full)
	if err != nil {
		if os.IsNotExist(err) {
			op.Kind = OpRemove
			op.Exists = false
			return op, nil
		}
		return Op{}, fmt.Errorf("cannot inspect %s: %w", rec.Path, err)
	}
	if st.Mode()&os.ModeSymlink != 0 {
		return Op{}, fmt.Errorf("%s is a symlink; refusing to remove it automatically", rec.Path)
	}
	if st.IsDir() {
		return Op{}, fmt.Errorf("%s is a directory; refusing to remove it", rec.Path)
	}
	diskSHA, err := hashFile(full)
	if err != nil {
		return Op{}, err
	}
	op.Exists = true
	op.DiskSHA = diskSHA
	if rec.SHA256 != "" && diskSHA != rec.SHA256 {
		op.Kind = OpModifiedRemove
	}
	return op, nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("cannot read %s: %w", filepath.Base(path), err)
	}
	defer f.Close()
	h := sha256.New()
	n, err := io.Copy(h, io.LimitReader(f, content.MaxHashSize+1))
	if err != nil {
		return "", fmt.Errorf("cannot read %s: %w", filepath.Base(path), err)
	}
	if n > content.MaxHashSize {
		return "", fmt.Errorf("%s is too large to verify safely", filepath.Base(path))
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
