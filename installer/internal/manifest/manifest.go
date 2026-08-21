package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/imflawlezz/cursor-utils/installer/internal/config"
)

const currentFormat = config.ManifestFormatVersion

type Manifest struct {
	FormatVersion    int                        `json:"formatVersion"`
	InstallerVersion string                     `json:"installerVersion"`
	Repository       string                     `json:"repository"`
	ContentVersion   string                     `json:"contentVersion"`
	InstalledAt      time.Time                  `json:"installedAt"`
	Platform         Platform                   `json:"platform"`
	CursorRoot       string                     `json:"cursorRoot"`
	Components       map[string]ComponentRecord `json:"components"`
}

type Platform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type ComponentRecord struct {
	Version string       `json:"version"`
	Files   []FileRecord `json:"files"`
}

type FileRecord struct {
	Path      string `json:"path"`
	SHA256    string `json:"sha256"`
	TakenOver bool   `json:"takenOver,omitempty"`
}

func New() *Manifest {
	return &Manifest{
		FormatVersion: currentFormat,
		Components:    make(map[string]ComponentRecord),
	}
}

func (m *Manifest) Clone() *Manifest {
	if m == nil {
		return nil
	}
	cp := *m
	cp.Components = make(map[string]ComponentRecord, len(m.Components))
	for id, rec := range m.Components {
		files := make([]FileRecord, len(rec.Files))
		copy(files, rec.Files)
		rec.Files = files
		cp.Components[id] = rec
	}
	return &cp
}

func (m *Manifest) File(rel string) (FileRecord, bool) {
	if m == nil {
		return FileRecord{}, false
	}
	for _, rec := range m.Components {
		for _, f := range rec.Files {
			if f.Path == rel {
				return f, true
			}
		}
	}
	return FileRecord{}, false
}

func (m *Manifest) Owns(rel string) bool {
	_, ok := m.File(rel)
	return ok
}

func (m *Manifest) AllFiles() []FileRecord {
	if m == nil {
		return nil
	}
	var out []FileRecord
	for _, rec := range m.Components {
		out = append(out, rec.Files...)
	}
	return out
}

func (m *Manifest) Empty() bool {
	return m == nil || len(m.AllFiles()) == 0
}

func (m *Manifest) ComponentVersion(id string) string {
	if m == nil {
		return ""
	}
	return m.Components[id].Version
}

func (m *Manifest) HasComponent(id string) bool {
	if m == nil {
		return false
	}
	rec, ok := m.Components[id]
	return ok && len(rec.Files) > 0
}

func (m *Manifest) PutFile(component string, rec FileRecord) {
	if m.Components == nil {
		m.Components = make(map[string]ComponentRecord)
	}
	cur := m.Components[component]
	replaced := false
	for i, f := range cur.Files {
		if f.Path == rec.Path {
			cur.Files[i] = rec
			replaced = true
			break
		}
	}
	if !replaced {
		cur.Files = append(cur.Files, rec)
	}
	m.Components[component] = cur
}

func (m *Manifest) RemoveFile(rel string) {
	if m == nil {
		return
	}
	for id, rec := range m.Components {
		files := rec.Files[:0]
		for _, f := range rec.Files {
			if f.Path != rel {
				files = append(files, f)
			}
		}
		rec.Files = files
		if len(rec.Files) == 0 {
			delete(m.Components, id)
		} else {
			m.Components[id] = rec
		}
	}
}

type FormatError struct {
	Message string
}

func (e *FormatError) Error() string { return e.Message }

func Load(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to read the cursor-utils manifest: %w", err)
	}
	return Parse(data)
}

func Parse(data []byte) (*Manifest, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, &FormatError{Message: "the cursor-utils manifest is empty"}
	}
	var raw struct {
		FormatVersion int `json:"formatVersion"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, &FormatError{Message: "the cursor-utils manifest is damaged and cannot be read"}
	}
	if raw.FormatVersion == 0 {
		return nil, &FormatError{Message: "the cursor-utils manifest is missing a format version"}
	}
	if raw.FormatVersion > currentFormat {
		return nil, &FormatError{Message: "this installation was created by a newer installer; update cursor-utils and try again"}
	}
	if raw.FormatVersion < currentFormat {
		migrated, err := migrate(data, raw.FormatVersion)
		if err != nil {
			return nil, err
		}
		data = migrated
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, &FormatError{Message: "the cursor-utils manifest is damaged and cannot be read"}
	}
	if m.Components == nil {
		m.Components = make(map[string]ComponentRecord)
	}
	m.FormatVersion = currentFormat
	return &m, nil
}

func migrate(data []byte, from int) ([]byte, error) {
	// Format 1 has no predecessor.
	return nil, &FormatError{Message: fmt.Sprintf("unsupported cursor-utils manifest format version %d", from)}
}

func WriteAtomic(path string, m *Manifest) error {
	if m == nil {
		return fmt.Errorf("manifest is nil")
	}
	m.FormatVersion = currentFormat
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to encode the cursor-utils manifest: %w", err)
	}
	data = append(data, '\n')
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("unable to create the cursor-utils manifest directory: %w", err)
	}
	tmp, err := os.CreateTemp(dir, "manifest.json.*.tmp")
	if err != nil {
		return fmt.Errorf("unable to create a temporary manifest file: %w", err)
	}
	tmpName := tmp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpName)
		}
	}()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("unable to write the cursor-utils manifest: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("unable to write the cursor-utils manifest: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("unable to write the cursor-utils manifest: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("unable to replace the cursor-utils manifest: %w", err)
	}
	cleanup = false
	return nil
}

func Remove(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to remove the cursor-utils manifest: %w", err)
	}
	dir := filepath.Dir(path)
	_ = os.Remove(dir) // no-op unless the directory is already empty
	return nil
}
