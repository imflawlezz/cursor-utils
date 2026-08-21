package semver

import (
	"reflect"
	"testing"
)

func TestCanonical(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{"v0.1.0", "v0.1.0", true},
		{"0.1.0", "v0.1.0", true},
		{"v1.2.3-rc.1", "v1.2.3-rc.1", true},
		{"v1.0.0+build", "v1.0.0", true},
		{"latest", "", false},
		{"v1", "", false},
		{"", "", false},
		{"not-a-version", "", false},
		{"v1.2", "", false},
	}
	for _, tc := range cases {
		got, ok := Canonical(tc.in)
		if ok != tc.ok || got != tc.want {
			t.Errorf("Canonical(%q) = %q, %v; want %q, %v", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}

func TestSortNewestFirst(t *testing.T) {
	t.Parallel()
	in := []string{"v0.1.0", "v0.2.0", "not-semver", "v0.1.1", "0.3.0", "zzz", "v0.2.0-beta.1"}
	got := SortNewestFirst(in)
	want := []string{"0.3.0", "v0.2.0", "v0.2.0-beta.1", "v0.1.1", "v0.1.0", "zzz", "not-semver"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SortNewestFirst = %#v; want %#v", got, want)
	}
}

func TestLatestIgnoresInvalid(t *testing.T) {
	t.Parallel()
	got, ok := Latest([]string{"nightly", "v0.1.0", "weird"})
	if !ok || got != "v0.1.0" {
		t.Fatalf("Latest = %q, %v", got, ok)
	}
	if _, ok := Latest([]string{"nightly", "weird"}); ok {
		t.Fatal("expected no latest among invalid tags")
	}
}

func TestPreReleaseOrdering(t *testing.T) {
	t.Parallel()
	got, ok := Latest([]string{"v1.0.0-alpha", "v1.0.0"})
	if !ok || got != "v1.0.0" {
		t.Fatalf("Latest = %q, %v; want v1.0.0", got, ok)
	}
}

func TestContentTagsExcludeInstallerReleases(t *testing.T) {
	t.Parallel()
	in := []string{"v0.1.0", "installer-v1.0.0", "nightly", "0.9.0", "v0.2.0", "installer-v0.2.0"}
	got := SortNewestFirst(ContentTags(in))
	want := []string{"v0.2.0", "v0.1.0"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ContentTags = %#v; want %#v", got, want)
	}
	latest, ok := Latest(ContentTags(in))
	if !ok || latest != "v0.2.0" {
		t.Fatalf("latest = %q %v", latest, ok)
	}
}

func TestIsContentAndInstallerTag(t *testing.T) {
	t.Parallel()
	if !IsContentTag("v0.1.0") || !IsContentTag("v1.0.0-rc.1") {
		t.Fatal("content tags rejected")
	}
	for _, tag := range []string{"", "0.1.0", "nightly", "installer-v1.0.0", "v1", "installer-v"} {
		if IsContentTag(tag) {
			t.Errorf("expected non-content: %q", tag)
		}
	}
	if !IsInstallerTag("installer-v1.0.0") || !IsInstallerTag("installer-v0.1.0-rc.1") {
		t.Fatal("installer tags rejected")
	}
	if IsInstallerTag("v1.0.0") || IsInstallerTag("installer-v") || IsInstallerTag("installer-1.0.0") {
		t.Fatal("non-installer tag accepted")
	}
}
