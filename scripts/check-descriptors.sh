#!/usr/bin/env bash
# check-descriptors.sh checks that every non-test, non-generated Go file opens
# its package comment with "// Package <name> <path> c<layer>-<domain>-<system>",
# that <path> is the file's own path, and that one package carries one descriptor.
#
# See docs/code-descriptors.md. Only well-formedness and per-package agreement
# are checked. Which domain a package belongs to is a documentation choice, and
# the layer is already enforced by the Go compiler, which rejects import cycles.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 1

DESC='c[0-4]-(com|coin|app|ui)-[a-z0-9]+'

fail=0
files=0
declare -A pkgdesc

while IFS= read -r f; do
	[ -n "$f" ] || continue
	grep -qE '^// Code generated .* DO NOT EDIT\.$' "$f" && continue
	files=$((files + 1))
	pkg=$(awk '/^package [A-Za-z_]/ { print $2; exit }' "$f")
	line=$(grep -m1 -E "^// Package ${pkg} [^ ]+ ${DESC}\$" "$f")
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
	d=$(dirname "$f")
	if [ -z "${pkgdesc[$d]:-}" ]; then
		pkgdesc[$d]=$desc
	elif [ "${pkgdesc[$d]}" != "$desc" ]; then
		echo "MIXED DESCRIPTORS: $d carries ${pkgdesc[$d]} and $desc ($f), one package takes one descriptor"
		fail=1
	fi
done < <(git ls-files -- '*.go' ':!:vendor/**' ':!:*_test.go' ':!:**/node_modules/**')

echo "descriptors: checked ${files} files in ${#pkgdesc[@]} packages"
exit "$fail"
