#!/usr/bin/env bash
# scripts/docs-prepare.sh — entry point for staging documentation sources
# under docs/ before an MkDocs build. The actual copy + link-rewriting logic
# lives in scripts/docs-prepare.py (Python is always present: it runs MkDocs).
#
# Used by .github/workflows/docs.yml (CI) and `make docs-serve` /
# `make docs-build` (local preview).
set -euo pipefail
exec python3 "$(dirname "$0")/docs-prepare.py"
