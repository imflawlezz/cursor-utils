package tui

func clampIndex(cur, n, delta int) int {
	if n <= 0 {
		return 0
	}
	cur += delta
	if cur < 0 {
		return 0
	}
	if cur > n-1 {
		return n - 1
	}
	return cur
}

func (m Model) pageSize() int {
	n := m.height - 12
	if n < 5 {
		return 5
	}
	return n
}

func (m Model) ensureManageCursor() Model {
	rows := m.manageRows()
	if len(rows) == 0 {
		return m
	}
	if m.manageCursor >= len(rows) {
		m.manageCursor = len(rows) - 1
	}
	if m.manageCursor < 0 {
		m.manageCursor = 0
	}
	if rows[m.manageCursor].focusable() {
		return m
	}
	m = m.moveManage(1)
	rows = m.manageRows()
	if m.manageCursor >= 0 && m.manageCursor < len(rows) && rows[m.manageCursor].focusable() {
		return m
	}
	return m.moveManage(-1)
}

func (m Model) moveManage(delta int) Model {
	rows := m.manageRows()
	i := m.manageCursor + delta
	for i >= 0 && i < len(rows) {
		if rows[i].focusable() {
			m.manageCursor = i
			return m
		}
		i += delta
	}
	return m
}

func (m Model) moveManagePage(dir int) Model {
	steps := m.pageSize()
	if dir < 0 {
		steps = -steps
	}
	for i := 0; i < abs(steps); i++ {
		next := m.moveManage(sign(steps))
		if next.manageCursor == m.manageCursor {
			break
		}
		m = next
	}
	return m
}

func (m Model) jumpManageEdge(dir int) Model {
	rows := m.manageRows()
	if dir < 0 {
		for i := range rows {
			if rows[i].focusable() {
				m.manageCursor = i
				return m
			}
		}
		return m
	}
	for i := len(rows) - 1; i >= 0; i-- {
		if rows[i].focusable() {
			m.manageCursor = i
			return m
		}
	}
	return m
}

func (m Model) jumpSection(dir int) Model {
	m = m.ensureManageCursor()
	rows := m.manageRows()
	if len(rows) == 0 || m.manageCursor < 0 || m.manageCursor >= len(rows) {
		return m
	}
	curSec := rows[m.manageCursor].Section
	if dir > 0 {
		for i := m.manageCursor + 1; i < len(rows); i++ {
			if rows[i].focusable() && rows[i].Section != curSec {
				m.manageCursor = i
				return m
			}
		}
		for i := range rows {
			if rows[i].focusable() {
				m.manageCursor = i
				return m
			}
		}
		return m
	}
	start := m.manageCursor
	for start > 0 && rows[start-1].Section == curSec {
		start--
	}
	for i := start - 1; i >= 0; i-- {
		if rows[i].focusable() && rows[i].Section != curSec {
			sec := rows[i].Section
			j := i
			for j > 0 && rows[j-1].Section == sec {
				j--
			}
			for j < len(rows) && !rows[j].focusable() {
				j++
			}
			m.manageCursor = j
			return m
		}
	}
	if len(rows) == 0 {
		return m
	}
	lastSec := rows[len(rows)-1].Section
	j := len(rows) - 1
	for j > 0 && rows[j-1].Section == lastSec {
		j--
	}
	for j < len(rows) && !rows[j].focusable() {
		j++
	}
	m.manageCursor = j
	return m
}

func (m Model) focusGroup(idx int) Model {
	for i, r := range m.manageRows() {
		if r.Kind == rowGroup && r.GroupIdx == idx {
			m.manageCursor = i
			return m
		}
	}
	return m
}

func (m Model) collapseOrBack() (Model, bool) {
	m = m.ensureManageCursor()
	rows := m.manageRows()
	if m.manageCursor < 0 || m.manageCursor >= len(rows) {
		return m, false
	}
	row := rows[m.manageCursor]
	switch row.Kind {
	case rowFile:
		m.groups[row.GroupIdx].Expanded = false
		return m.focusGroup(row.GroupIdx), false
	case rowGroup:
		if m.groups[row.GroupIdx].Expanded {
			m.groups[row.GroupIdx].Expanded = false
			return m, false
		}
		return m, true
	default:
		return m, true
	}
}

func (m Model) expandOrForward() Model {
	m = m.ensureManageCursor()
	rows := m.manageRows()
	if m.manageCursor < 0 || m.manageCursor >= len(rows) {
		return m
	}
	row := rows[m.manageCursor]
	if row.Kind != rowGroup {
		return m
	}
	if !m.groups[row.GroupIdx].Expanded {
		m.groups[row.GroupIdx].Expanded = true
		return m
	}
	return m.moveManage(1)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func sign(n int) int {
	if n < 0 {
		return -1
	}
	return 1
}
