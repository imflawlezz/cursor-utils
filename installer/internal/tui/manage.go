package tui

import (
	"github.com/imflawlezz/cursor-utils/installer/internal/content"
	"github.com/imflawlezz/cursor-utils/installer/internal/installer"
)

func (m *Model) initSelection() {
	prev := map[string]bool{}
	expanded := map[string]bool{}
	for _, g := range m.groups {
		expanded[g.ID] = g.Expanded
		for _, f := range g.Files {
			prev[f.RelPath] = f.Selected
		}
	}
	m.groups = nil
	if m.bundle == nil {
		return
	}
	for _, spec := range content.Registry {
		id := string(spec.ID)
		comp, ok := m.bundle.Components[id]
		if !ok || len(comp.Files) == 0 {
			continue
		}
		g := groupState{ID: id, Expanded: expanded[id]}
		for _, f := range comp.Files {
			sel := true
			if v, ok := prev[f.RelPath]; ok {
				sel = v
			}
			g.Files = append(g.Files, fileState{RelPath: f.RelPath, Selected: sel})
		}
		m.groups = append(m.groups, g)
	}
}

func (m Model) selection() installer.Selection {
	sel := installer.Selection{
		Files:          map[string]bool{},
		FullComponents: map[string]bool{},
	}
	for _, g := range m.groups {
		if g.allSelected() {
			sel.FullComponents[g.ID] = true
			continue
		}
		for _, f := range g.Files {
			if f.Selected {
				sel.Files[f.RelPath] = true
			}
		}
	}
	return sel
}

func (m Model) selectedCount() int {
	n := 0
	for _, g := range m.groups {
		for _, f := range g.Files {
			if f.Selected {
				n++
			}
		}
	}
	return n
}

func (m Model) manageRows() []manageRow {
	var rows []manageRow
	for i, g := range m.groups {
		rows = append(rows, manageRow{Kind: rowGroup, GroupIdx: i, Section: secComponents})
		if g.Expanded {
			for j := range g.Files {
				rows = append(rows, manageRow{Kind: rowFile, GroupIdx: i, FileIdx: j, Section: secComponents})
			}
		}
	}
	rows = append(rows, manageRow{Kind: rowLabel, Action: "Actions"})
	if m.selectedCount() > 0 {
		for _, a := range []string{"Install / Update selected", "Remove selected"} {
			rows = append(rows, manageRow{Kind: rowAction, Action: a, Section: secApply})
		}
	}
	rows = append(rows, manageRow{Kind: rowAction, Action: "Open Cursor folder", Section: secApply})
	rows = append(rows, manageRow{Kind: rowSpacer})
	for _, a := range []string{"Change version", "Change directory", "Repair TUI", "Keybindings", "Quit"} {
		sec := secSettings
		if a == "Keybindings" || a == "Quit" || a == "Repair TUI" {
			sec = secApp
		}
		rows = append(rows, manageRow{Kind: rowAction, Action: a, Section: sec})
	}
	return rows
}

func (m Model) toggleGroup(idx int) Model {
	if idx < 0 || idx >= len(m.groups) {
		return m
	}
	on := !m.groups[idx].allSelected()
	for i := range m.groups[idx].Files {
		m.groups[idx].Files[i].Selected = on
	}
	return m
}

func (m Model) toggleFile(g, f int) Model {
	if g < 0 || g >= len(m.groups) {
		return m
	}
	files := m.groups[g].Files
	if f < 0 || f >= len(files) {
		return m
	}
	m.groups[g].Files[f].Selected = !m.groups[g].Files[f].Selected
	return m
}
