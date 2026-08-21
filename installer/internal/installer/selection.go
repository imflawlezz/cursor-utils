package installer

type Selection struct {
	Files          map[string]bool
	FullComponents map[string]bool
}

func (s Selection) Empty() bool {
	return len(s.Files) == 0 && len(s.FullComponents) == 0
}

func (s Selection) Includes(op Op) bool {
	if s.FullComponents[op.Component] {
		return true
	}
	return s.Files[op.RelPath]
}

func (p *Plan) Filter(sel Selection) *Plan {
	if p == nil {
		return nil
	}
	out := *p
	if sel.Empty() {
		out.Ops = nil
		return &out
	}
	ops := make([]Op, 0, len(p.Ops))
	for _, op := range p.Ops {
		if sel.Includes(op) {
			ops = append(ops, op)
		}
	}
	out.Ops = ops
	return &out
}
