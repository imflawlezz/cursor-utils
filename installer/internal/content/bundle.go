package content

import (
	"context"
	"fmt"
)

type File struct {
	RelPath string
	Data    []byte
	SHA256  string
}

type Component struct {
	Spec  Spec
	Files []File
}

type Bundle struct {
	Version    string
	Components map[string]*Component
}

func NewBundle(version string) *Bundle {
	return &Bundle{
		Version:    version,
		Components: make(map[string]*Component),
	}
}

func (b *Bundle) Add(rel string, data []byte) error {
	if err := ValidateRelPath(rel); err != nil {
		return err
	}
	if len(data) > MaxFileSize {
		return fmt.Errorf("%s exceeds the maximum file size", rel)
	}
	id := ComponentOf(rel)
	spec, ok := Lookup(id)
	if !ok {
		return fmt.Errorf("unknown component %q", id)
	}
	comp, ok := b.Components[id]
	if !ok {
		comp = &Component{Spec: spec}
		b.Components[id] = comp
	}
	comp.Files = append(comp.Files, File{
		RelPath: rel,
		Data:    data,
		SHA256:  SHA256(data),
	})
	return nil
}

func (b *Bundle) Empty() bool {
	if b == nil {
		return true
	}
	for _, c := range b.Components {
		if len(c.Files) > 0 {
			return false
		}
	}
	return true
}

func (b *Bundle) File(rel string) (File, bool) {
	if b == nil {
		return File{}, false
	}
	id := ComponentOf(rel)
	comp, ok := b.Components[id]
	if !ok {
		return File{}, false
	}
	for _, f := range comp.Files {
		if f.RelPath == rel {
			return f, true
		}
	}
	return File{}, false
}

func (b *Bundle) AllFiles() []File {
	if b == nil {
		return nil
	}
	var out []File
	for _, spec := range Registry {
		comp, ok := b.Components[string(spec.ID)]
		if !ok {
			continue
		}
		out = append(out, comp.Files...)
	}
	return out
}

func (b *Bundle) HasComponent(id string) bool {
	if b == nil {
		return false
	}
	c, ok := b.Components[id]
	return ok && len(c.Files) > 0
}

// Provider fetches tagged content over HTTPS. It must not use Git or the TUI.
type Provider interface {
	ListTags(ctx context.Context) ([]string, error)
	Fetch(ctx context.Context, tag string) (*Bundle, error)
}
