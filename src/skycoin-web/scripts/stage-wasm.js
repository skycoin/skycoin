// Stages the skycoin-lite wasm cipher for `ng serve` and `ng test`.
//
// The repository keeps exactly one copy of each build, in the Go package that
// embeds it: src/skycoin-lite/wasm-go and src/skycoin-lite/wasm-tinygo. The
// server picks one by build tag (cmd/skycoin-web/commands/wasm_go.go) and
// serves it from memory at /assets/scripts/, ahead of the static handler — so
// the production bundle neither needs nor gets a copy.
//
// The Angular dev server and karma do not go through that server, and Angular
// will not read an asset from outside the workspace root, so both get a staged
// copy here. It is gitignored and rewritten on every run, which is what keeps
// it from drifting the way the old committed copy did.
//
// The Go build is staged rather than the TinyGo one because a standard `go
// build` of the wallet embeds the Go build; wasm_exec.js differs between the
// two toolchains and has to match the wasm beside it.
const fs = require('fs');
const path = require('path');
const zlib = require('zlib');

const src = path.resolve(__dirname, '../../skycoin-lite/wasm-go');
const dest = path.resolve(__dirname, '../src/assets/scripts');
// The wasm is committed gzipped, because it is a large artifact that changes on
// every rebuild. The Go server hands that straight to the browser with
// Content-Encoding: gzip, but the Angular dev server and karma serve static
// files as-is, so the staged copy is expanded here.
const files = [
  { from: 'skycoin-lite.wasm.gz', to: 'skycoin-lite.wasm', gzipped: true },
  { from: 'wasm_exec.js', to: 'wasm_exec.js', gzipped: false },
];

for (const f of files) {
  const from = path.join(src, f.from);
  if (!fs.existsSync(from)) {
    console.error(`${f.from} not found in ${src} — run "make build-wasm" first.`);
    process.exit(1);
  }

  const body = fs.readFileSync(from);
  fs.writeFileSync(path.join(dest, f.to), f.gzipped ? zlib.gunzipSync(body) : body);
}
console.log(`staged ${files.map(f => f.to).join(' and ')} from ${path.relative(process.cwd(), src)}`);
