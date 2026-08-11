// Serves a skycoin-web dist directory and proxies /api to a node.
//
// `skycoin web --gui-dir` cannot be used for verification: it serves /assets
// from the binary's embedded copy, so the wasm cipher and the i18n JSON would
// come from the old bundle no matter what was built. Everything here comes from
// the directory under test.
const http = require('http');
const fs = require('fs');
const path = require('path');

const ROOT = path.resolve(process.argv[2]);
const PORT = Number(process.argv[3]);
const NODE = process.argv[4] || 'http://127.0.0.1:6420';
const nodeUrl = new URL(NODE);

const TYPES = {
  '.html': 'text/html', '.js': 'text/javascript', '.mjs': 'text/javascript',
  '.css': 'text/css', '.json': 'application/json', '.wasm': 'application/wasm',
  '.svg': 'image/svg+xml', '.png': 'image/png', '.jpg': 'image/jpeg',
  '.ico': 'image/x-icon', '.woff': 'font/woff', '.woff2': 'font/woff2',
  '.ttf': 'font/ttf', '.eot': 'application/vnd.ms-fontobject', '.map': 'application/json',
};

http.createServer((req, res) => {
  if (req.url.startsWith('/api')) {
    const p = http.request({
      hostname: nodeUrl.hostname, port: nodeUrl.port, path: req.url,
      method: req.method,
      // The node validates Host, Origin and Referer, so present the request as
      // if it came from the node's own origin. ci-scripts/ui-e2e.sh does the
      // same thing for the desktop wallet's dev proxy.
      headers: { ...req.headers, host: nodeUrl.host,
                 origin: NODE, referer: NODE + '/' },
    }, r => { res.writeHead(r.statusCode, r.headers); r.pipe(res); });
    p.on('error', e => { res.writeHead(502); res.end(String(e)); });
    req.pipe(p);
    return;
  }

  let rel = decodeURIComponent(req.url.split('?')[0]);
  if (rel === '/') rel = '/index.html';
  let file = path.join(ROOT, rel);
  if (!file.startsWith(ROOT)) { res.writeHead(403); res.end(); return; }
  if (!fs.existsSync(file) || fs.statSync(file).isDirectory()) {
    file = path.join(ROOT, 'index.html');   // SPA fallback
  }
  const body = fs.readFileSync(file);
  res.writeHead(200, { 'Content-Type': TYPES[path.extname(file)] || 'application/octet-stream' });
  res.end(body);
}).listen(PORT, '127.0.0.1', () => console.log('serving ' + ROOT + ' on ' + PORT + ' -> ' + NODE));
