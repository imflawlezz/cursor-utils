package installer

import (
	"github.com/imflawlezz/cursor-utils/installer/internal/content"
	"github.com/imflawlezz/cursor-utils/installer/internal/manifest"
	"github.com/imflawlezz/cursor-utils/installer/internal/platform"
)

type PlanKind int

const (
	PlanInstall PlanKind = iota
	PlanUninstall
)

type OpKind int

const (
	OpAdd OpKind = iota
	OpUpdate
	OpRemove
	OpUnchanged
	OpConflict
	OpModifiedUpdate
	OpModifiedRemove
)

func (k OpKind) NeedsDecision() bool {
	switch k {
	case OpConflict, OpModifiedUpdate, OpModifiedRemove:
		return true
	default:
		return false
	}
}

type Choice int

const (
	ChoiceNone Choice = iota
	ChoiceSkip
	ChoiceOverwrite
	ChoiceKeepLocal
	ChoiceCancel
)

type Op struct {
	Kind        OpKind
	Component   string
	RelPath     string
	Data        []byte
	NewSHA      string
	ManifestSHA string
	DiskSHA     string
	Exists      bool
	TakenOver   bool
}

type Plan struct {
	Kind           PlanKind
	CursorRoot     string
	ContentVersion string
	Platform       platform.Info
	Existing       *manifest.Manifest
	Bundle         *content.Bundle
	Ops            []Op
}

func (p *Plan) NeedsDecision() []Op {
	var out []Op
	for _, op := range p.Ops {
		if op.Kind.NeedsDecision() {
			out = append(out, op)
		}
	}
	return out
}

func (p *Plan) FileCount() int {
	n := 0
	for _, op := range p.Ops {
		switch op.Kind {
		case OpAdd, OpUpdate, OpUnchanged, OpConflict, OpModifiedUpdate:
			n++
		}
	}
	return n
}

type ComponentSummary struct {
	ID            string
	DisplayName   string
	Available     bool
	Installed     bool
	InstalledVer  string
	TargetVer     string
	Action        string
	PendingOps    int
	ConflictCount int
	ModifiedCount int
}

func (p *Plan) ComponentSummaries() []ComponentSummary {
	type acc struct {
		available, installed                               bool
		add, update, remove, unchanged, conflict, modified int
	}
	stats := make(map[string]*acc)
	ensure := func(id string) *acc {
		s, ok := stats[id]
		if !ok {
			s = &acc{}
			stats[id] = s
		}
		return s
	}
	if p.Bundle != nil {
		for id, comp := range p.Bundle.Components {
			if len(comp.Files) > 0 {
				ensure(id).available = true
			}
		}
	}
	if p.Existing != nil {
		for id, rec := range p.Existing.Components {
			if len(rec.Files) > 0 {
				ensure(id).installed = true
			}
		}
	}
	for _, op := range p.Ops {
		s := ensure(op.Component)
		switch op.Kind {
		case OpAdd:
			s.add++
		case OpUpdate:
			s.update++
		case OpRemove:
			s.remove++
		case OpUnchanged:
			s.unchanged++
		case OpConflict:
			s.conflict++
		case OpModifiedUpdate:
			s.modified++
			s.update++
		case OpModifiedRemove:
			s.modified++
			s.remove++
		}
	}

	out := make([]ComponentSummary, 0, len(content.Registry))
	for _, spec := range content.Registry {
		id := string(spec.ID)
		s := stats[id]
		if s == nil {
			s = &acc{}
		}
		sum := ComponentSummary{
			ID:            id,
			DisplayName:   spec.DisplayName,
			Available:     s.available,
			Installed:     s.installed,
			InstalledVer:  "",
			TargetVer:     p.ContentVersion,
			ConflictCount: s.conflict,
			ModifiedCount: s.modified,
			PendingOps:    s.add + s.update + s.remove + s.conflict,
		}
		if p.Existing != nil {
			sum.InstalledVer = p.Existing.ComponentVersion(id)
		}
		switch {
		case !s.available && !s.installed:
			sum.Action = "Not available"
		case !s.available && s.installed:
			sum.Action = "Remove"
		case s.available && !s.installed:
			sum.Action = "Install"
		case s.conflict+s.modified+s.add+s.update+s.remove == 0:
			sum.Action = "Up to date"
		default:
			sum.Action = "Update"
		}
		out = append(out, sum)
	}
	return out
}

func (p *Plan) HasWork() bool {
	for _, op := range p.Ops {
		switch op.Kind {
		case OpAdd, OpUpdate, OpRemove, OpConflict, OpModifiedUpdate, OpModifiedRemove:
			return true
		}
	}
	return false
}

func (p *Plan) InstalledAnything() bool {
	return p.Existing != nil && !p.Existing.Empty()
}
