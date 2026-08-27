#!/bin/sh
# Print the Keep a Changelog body for a version tag.
#   sh scripts/release-notes.sh v0.2.0
#   sh scripts/release-notes.sh installer-v1.2.0
set -eu

tag=${1:-${GORELEASER_CURRENT_TAG:-${GITHUB_REF_NAME:-}}}
[ -n "$tag" ] || {
	printf 'usage: %s <tag>\n' "$0" >&2
	exit 2
}

case "$tag" in
installer-v*)
	ver="installer-${tag#installer-v}"
	;;
v*)
	ver=${tag#v}
	;;
*)
	ver=$tag
	;;
esac

file=CHANGELOG.md
[ -f "$file" ] || {
	printf 'release-notes: missing %s\n' "$file" >&2
	exit 1
}

awk -v ver="$ver" '
	BEGIN { heading = "^## \\[" ver "\\]" }
	$0 ~ heading {
		found = 1
		next
	}
	found && /^## \[/ { exit }
	found && /^\[[^]]+\]:/ { exit }
	found {
		if (!started) {
			if ($0 ~ /^[ \t]*$/) next
			started = 1
		}
		print
	}
	END {
		if (!found) {
			printf "release-notes: no CHANGELOG.md section for %s\n", ver > "/dev/stderr"
			exit 1
		}
	}
' "$file"
