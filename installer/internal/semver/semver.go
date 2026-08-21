package semver

import (
	"sort"
	"strings"

	modsemver "golang.org/x/mod/semver"

	"github.com/imflawlezz/cursor-utils/installer/internal/config"
)

// Canonical is for ordering only; the original tag is left unchanged.
func Canonical(tag string) (string, bool) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return "", false
	}
	v := tag
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	if !modsemver.IsValid(v) {
		return "", false
	}
	// x/mod/semver treats "v1" as valid; we require MAJOR.MINOR.PATCH.
	core := strings.TrimPrefix(v, "v")
	if i := strings.IndexAny(core, "-+"); i >= 0 {
		core = core[:i]
	}
	if strings.Count(core, ".") < 2 {
		return "", false
	}
	return modsemver.Canonical(v), true
}

func Valid(tag string) bool {
	_, ok := Canonical(tag)
	return ok
}

// SortNewestFirst orders valid SemVer tags newest-first. Invalid tags follow
// reverse-lexicographically and are never chosen as latest.
func SortNewestFirst(tags []string) []string {
	type item struct {
		original  string
		canonical string
		valid     bool
	}
	items := make([]item, 0, len(tags))
	for _, t := range tags {
		it := item{original: t}
		if c, ok := Canonical(t); ok {
			it.canonical = c
			it.valid = true
		}
		items = append(items, it)
	}
	sort.SliceStable(items, func(i, j int) bool {
		a, b := items[i], items[j]
		if a.valid && b.valid {
			cmp := modsemver.Compare(a.canonical, b.canonical)
			if cmp != 0 {
				return cmp > 0
			}
			return a.original > b.original
		}
		if a.valid != b.valid {
			return a.valid
		}
		return a.original > b.original
	})
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.original
	}
	return out
}

func Latest(tags []string) (string, bool) {
	for _, t := range SortNewestFirst(tags) {
		if Valid(t) {
			return t, true
		}
	}
	return "", false
}

// IsInstallerTag is true for installer-vMAJOR.MINOR.PATCH, never for content.
func IsInstallerTag(tag string) bool {
	tag = strings.TrimSpace(tag)
	rest, ok := strings.CutPrefix(tag, config.InstallerTagPrefix)
	if !ok || rest == "" {
		return false
	}
	return Valid("v" + rest)
}

// IsContentTag is true for vMAJOR.MINOR.PATCH only. installer-v* tags are excluded.
func IsContentTag(tag string) bool {
	tag = strings.TrimSpace(tag)
	if IsInstallerTag(tag) {
		return false
	}
	if !strings.HasPrefix(tag, "v") {
		return false
	}
	return Valid(tag)
}

func ContentTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		if IsContentTag(t) {
			out = append(out, t)
		}
	}
	return out
}
