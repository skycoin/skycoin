// Zeroes the build time Angular writes into the service worker manifest, so
// the committed bundle is the same on every build and CI can compare it.
// The worker decides updates by the content hashes, not by this field.
const fs = require('fs');
const path = require('path');

const manifest = path.resolve(__dirname, '../src/gui/dist/ngsw.json');
if (fs.existsSync(manifest)) {
  const text = fs.readFileSync(manifest, 'utf8');
  fs.writeFileSync(manifest, text.replace(/"timestamp": \d+/, '"timestamp": 0'));
}
