#!/usr/bin/env bash
# check-descriptors.sh checks that every non-test, non-generated Go file opens
# its package comment with "// Package <name> <path> c<layer>-<domain>-<system>",
# that <path> is the file's own path, and that one package carries one descriptor.
#
# See docs/code-descriptors.md. Only well-formedness and per-package agreement
# are checked. Which domain a package belongs to is a documentation choice, and
# the layer is already enforced by the Go compiler, which rejects import cycles.
# Written for bash 3.2 (macOS): no associative arrays.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 1

DESC='c[0-4]-(com|coin|app|ui)-[a-z0-9]+'

fail=0
files=0
# One "<dir> <descriptor>" line per checked file, for the per-package check.
pairs=$(mktemp) || exit 1
trap 'rm -f "$pairs"' EXIT

while IFS= read -r f; do
	[ -n "$f" ] || continue
	# tr drops the CR of a CRLF checkout (git on Windows).
	body=$(tr -d '\r' <"$f")
	grep -qE '^// Code generated .* DO NOT EDIT\.$' <<<"$body" && continue
	files=$((files + 1))
	pkg=$(awk '/^package [A-Za-z_]/ { print $2; exit }' <<<"$body")
	line=$(grep -m1 -E "^// Package ${pkg} [^ ]+ ${DESC}\$" <<<"$body")
	if [ -z "$line" ]; then
		echo "NO DESCRIPTOR: $f, add '// Package ${pkg} ${f} c<layer>-<domain>-<system>' (docs/code-descriptors.md)"
		fail=1
		continue
	fi
	path=$(awk '{ print $4 }' <<<"$line")
	desc=$(awk '{ print $5 }' <<<"$line")
	if [ "$path" != "$f" ]; then
		echo "WRONG PATH: $f names itself $path"
		fail=1
	fi
	echo "$(dirname "$f") $desc $f" >>"$pairs"
done < <(git ls-files -- '*.go' ':!:vendor/**' ':!:*_test.go' ':!:**/node_modules/**')

mixed=$(awk '
	!($1 in first) { first[$1] = $2; next }
	$2 != first[$1] && !seen[$1 " " $2]++ {
		print "MIXED DESCRIPTORS: " $1 " carries " first[$1] " and " $2 " (" $3 "), one package takes one descriptor"
	}' "$pairs")
if [ -n "$mixed" ]; then
	echo "$mixed"
	fail=1
fi

pkgs=$(awk '{ print $1 }' "$pairs" | sort -u | grep -c .)
echo "descriptors: checked ${files} files in ${pkgs} packages"
exit "$fail"
