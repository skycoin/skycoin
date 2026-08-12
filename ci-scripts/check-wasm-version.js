#!/usr/bin/env node
//
// Fails when a committed wasm was built from a dirty working tree.
//
// The wasm cipher is a build artifact kept in the repository, so it can never be
// checked by rebuilding and comparing: the commit it was compiled at is always
// an earlier one than the commit that carries it, and Go's output is not
// byte-reproducible across toolchain versions anyway.
//
// What is checkable is whether the thing was built from a commit at all. `go
// build` inside a git work tree stamps vcs.revision and vcs.modified into the
// binary, in plain text — no need to run the module or even parse wasm.
// vcs.modified=true means the tree had uncommitted changes, so no commit
// describes what is in the file and nobody can reproduce it. That is the state
// this rejects.
//
// A clean build is expected to report vcs.modified=false and a revision that is
// an ancestor of HEAD. The revision itself is not checked against HEAD, because
// committing the artifact necessarily moves HEAD past it.
//
// Both artifacts must carry a stamp. Upstream TinyGo writes none, so an
// unstamped file means the wasm was built with a toolchain that cannot say what
// it is — indistinguishable from any other build, and unreproducible for the
// same reason a dirty one is. github.com/0magnet/tinygo records them.
const fs = require('fs');
const path = require('path');
const zlib = require('zlib');

const ARTIFACTS = [
  'src/skycoin-lite/wasm-go/skycoin-lite.wasm.gz',
  'src/skycoin-lite/wasm-tinygo/skycoin-lite.wasm.gz',
];

const repoRoot = path.resolve(__dirname, '..');
let failed = false;

for (const rel of ARTIFACTS) {
  const file = path.join(repoRoot, rel);
  if (!fs.existsSync(file)) {
    console.error(`${rel}: missing — run "make build-wasm"`);
    failed = true;
    continue;
  }

  // The build stamp is a run of "key\tvalue" records in the binary's string
  // data; read it directly rather than depending on `go version -m`, which does
  // not recognise the wasm container. The artifact is committed gzipped, so
  // expand it in memory first.
  const contents = zlib.gunzipSync(fs.readFileSync(file)).toString('latin1');
  const read = key => {
    const found = contents.match(new RegExp(`build\\s${key}=([^\\s\\0]*)`));
    return found ? found[1] : null;
  };

  const modified = read('vcs\\.modified');
  const revision = read('vcs\\.revision');

  if (modified === null && revision === null) {
    console.error(
      `${rel}: no build stamp.\n` +
      `  Nothing records which commit this was built from, so it cannot be told\n` +
      `  apart from any other build. Upstream TinyGo writes no build information;\n` +
      `  rebuild with one that does:\n` +
      `    make build-wasm TINYGO=/path/to/0magnet/tinygo/build/tinygo`);
    failed = true;
    continue;
  }

  if (modified === 'true') {
    console.error(
      `${rel}: built from a dirty working tree.\n` +
      `  It reports vcs.modified=true, so no commit describes what is in it.\n` +
      `  Commit the source first, then rebuild it with "make build-wasm" and commit the result.`);
    failed = true;
    continue;
  }

  console.log(`${rel}: clean, built at ${(revision || 'unknown').slice(0, 12)}`);
}

process.exit(failed ? 1 : 0);
