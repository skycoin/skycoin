// Stages the Go cipher test vectors where the karma builder can serve them.
//
// src/app/services/cipher.provider.lib.spec.ts checks the browser cipher —
// assets/scripts/skycoin-lite.wasm — against src/cipher/testsuite/testdata, the
// same golden files TestManyAddresses and friends use on the Go side. Checking
// both implementations against one set of vectors is the point: a copy kept in
// this project would drift, and a wallet that derives subtly different
// addresses than the node is the failure this is here to catch.
//
// The files cannot be referenced in place. Angular rejects an asset path
// outside the workspace root, which for this project is src/skycoin-web, so
// they are staged into a gitignored directory before each run.
const fs = require('fs');
const path = require('path');

const src = path.resolve(__dirname, '../../cipher/testsuite/testdata');
const dest = path.resolve(__dirname, '../e2e/test-fixtures');

if (!fs.existsSync(src)) {
  console.error(`cipher test vectors not found at ${src}`);
  process.exit(1);
}

fs.rmSync(dest, { recursive: true, force: true });
fs.mkdirSync(dest, { recursive: true });

const staged = fs.readdirSync(src).filter(f => f.endsWith('.golden'));
for (const f of staged) {
  fs.copyFileSync(path.join(src, f), path.join(dest, f));
}
console.log(`staged ${staged.length} cipher test vectors in ${path.relative(process.cwd(), dest)}`);
