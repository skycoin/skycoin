#!/usr/bin/env python3
"""Stage markdown from across the repo under docs/ for the MkDocs build.

MkDocs requires every source file to live under a single docs_dir. A lot of
Skycoin's documentation, however, lives next to the code it documents (the
top-level guides, src/api, cmd/skycoin-cli and assorted package/client
READMEs). This script copies those files into docs/ at build time and
rewrites their relative links so they resolve inside the rendered site.

The staged copies are gitignored (see docs/.gitignore) — the canonical copy
stays in its original location. Idempotent; safe to re-run before every
`mkdocs serve` / `mkdocs build`.

Link rewriting, per staged file (links are resolved against the file's
ORIGINAL location, then):
  - a link to another staged source -> a site-relative path to that page;
  - a link to any other file that exists in the repo, or any relative link
    we can't resolve in-site -> an absolute github.com/blob/develop URL;
  - root-absolute ("/path") links to existing repo files -> a github URL;
    otherwise left untouched;
  - http(s)/mailto/anchor-only links -> left untouched.

Used by scripts/docs-prepare.sh (CI + `make docs-serve` / `make docs-build`).
"""

import os
import posixpath
import re
import shutil
import sys

REPO_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DOCS_DIR = os.path.join(REPO_ROOT, "docs")
GITHUB_BLOB = "https://github.com/skycoin/skycoin/blob/develop/"

# original repo path -> destination path relative to docs/
MAPPING = {
    # Guides (top-level docs)
    "INSTALLATION.md": "guides/installation.md",
    "INTEGRATION.md": "guides/integration.md",
    "DEVELOPMENT.md": "guides/development.md",
    "DOCKER.md": "guides/docker.md",
    "FIBERCOIN-CONFIG.md": "guides/fibercoin-config.md",
    "RELEASE.md": "guides/release.md",
    "CHANGELOG.md": "guides/changelog.md",
    # REST API & CLI
    "src/api/README.md": "api/index.md",
    "cmd/skycoin-cli/README.md": "cli/index.md",
    # Packages
    "src/cipher/bip32/README.md": "packages/cipher-bip32.md",
    "src/cipher/bip39/README.md": "packages/cipher-bip39.md",
    "src/cipher/base58/README.md": "packages/cipher-base58.md",
    "src/cipher/encoder/README.md": "packages/cipher-encoder.md",
    "src/cipher/secp256k1-go/README.md": "packages/cipher-secp256k1.md",
    "src/daemon/pex/README.md": "packages/daemon-pex.md",
    "src/daemon/gnet/README.md": "packages/daemon-gnet.md",
    # Clients
    "src/gui/static/README.md": "clients/desktop.md",
    "src/skycoin-web/README.md": "clients/web.md",
    "src/skycoin-lite/README.md": "clients/lite.md",
    "explorer/README.md": "clients/explorer.md",
}

# Matches markdown inline link / image targets: ](url) or ](<url> "title").
LINK_RE = re.compile(r"\]\(\s*(<[^>]+>|[^)\s]+)(\s+\"[^\"]*\")?\s*\)")

SKIP_PREFIXES = ("http://", "https://", "mailto:", "tel:", "#")


def repo_exists(path):
    return os.path.isfile(os.path.join(REPO_ROOT, path))


def rewrite_links(text, src_repo_path, dest_rel):
    src_dir = posixpath.dirname(src_repo_path)
    dest_dir = posixpath.dirname(dest_rel)

    def repl(m):
        raw = m.group(1)
        title = m.group(2) or ""
        url = raw[1:-1] if raw.startswith("<") and raw.endswith(">") else raw

        if url.startswith(SKIP_PREFIXES) or not url:
            return m.group(0)

        path, _, anchor = url.partition("#")
        anchor = ("#" + anchor) if anchor else ""
        if not path:  # pure anchor
            return m.group(0)

        root_absolute = path.startswith("/")
        candidate = path[1:] if root_absolute else posixpath.normpath(
            posixpath.join(src_dir, path)
        )

        if candidate in MAPPING:
            new = posixpath.relpath(MAPPING[candidate], dest_dir or ".") + anchor
        elif repo_exists(candidate):
            new = GITHUB_BLOB + candidate + anchor
        elif root_absolute:
            return m.group(0)  # unknown absolute link — leave as-is
        else:
            # Unresolved relative link (dangling in the source too). Point at
            # where it would live on GitHub so the site has no dead in-tree
            # links; harmless if the target is also missing upstream.
            new = GITHUB_BLOB + candidate + anchor

        return "](" + new + title + ")"

    return LINK_RE.sub(repl, text)


def main():
    staged = []
    for src_rel, dest_rel in MAPPING.items():
        src = os.path.join(REPO_ROOT, src_rel)
        if not os.path.isfile(src):
            print(f"docs-prepare: WARNING: missing source {src_rel} "
                  f"(skipping {dest_rel})", file=sys.stderr)
            continue
        dest = os.path.join(DOCS_DIR, dest_rel)
        os.makedirs(os.path.dirname(dest), exist_ok=True)
        with open(src, encoding="utf-8") as fh:
            text = fh.read()
        text = rewrite_links(text, src_rel, dest_rel)
        with open(dest, "w", encoding="utf-8") as fh:
            fh.write(text)
        staged.append(dest_rel)

    print(f"docs-prepare: staged {len(staged)} pages under docs/ "
          "(guides, api, cli, packages, clients)")


if __name__ == "__main__":
    main()
