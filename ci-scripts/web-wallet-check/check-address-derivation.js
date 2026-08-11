// Loads a known seed into the web wallet and checks the address it derives.
const p = require('puppeteer-core');
const sleep = ms => new Promise(r => setTimeout(r, ms));

const PORT = process.argv[2];
const SEED = process.argv[3];
const EXPECT = process.argv[4];

(async () => {
  const b = await p.launch({ executablePath: process.env.CHROME_BIN, headless: 'new',
    args: ['--no-sandbox', '--disable-dev-shm-usage'] });
  const pg = await b.newPage();
  const logs = [];
  pg.on('console', m => { const t = m.text(); if (/WASM|error/i.test(t)) logs.push(m.type() + ': ' + t.slice(0, 120)); });
  pg.on('pageerror', e => logs.push('PAGEERROR: ' + String(e).slice(0, 120)));

  await pg.goto('http://127.0.0.1:' + PORT + '/', { waitUntil: 'networkidle2', timeout: 60000 });
  await sleep(6000);

  // The "New" tab generates a seed through bip39, which is a different code
  // path from loading one and the one that fails without the Buffer global.
  const generated = await pg.evaluate(() => {
    const el = document.querySelector('[formcontrolname="seed"]');
    return el ? (el.value || '').trim() : '';
  });
  const wordCount = generated ? generated.split(/\s+/).filter(Boolean).length : 0;
  console.log('generated seed words:', wordCount);

  // Switch to the "Load" tab of the onboarding form.
  await pg.evaluate(() => {
    const el = [...document.querySelectorAll('*')]
      .filter(e => e.children.length === 0 && /^Load$/i.test((e.innerText || '').trim()))[0];
    if (el) el.click();
  });
  await sleep(1500);

  const filled = await pg.evaluate(seed => {
    const set = (name, v) => {
      const el = document.querySelector(`[formcontrolname="${name}"]`);
      if (!el) return false;
      const proto = el.tagName === 'TEXTAREA'
        ? window.HTMLTextAreaElement.prototype
        : window.HTMLInputElement.prototype;
      Object.getOwnPropertyDescriptor(proto, 'value').set.call(el, v);
      el.dispatchEvent(new Event('input', { bubbles: true }));
      el.dispatchEvent(new Event('blur', { bubbles: true }));
      return true;
    };
    return { label: set('label', 'harness'), seed: set('seed', seed) };
  }, SEED);
  await sleep(1500);

  // Submit.
  await pg.evaluate(() => {
    const btn = [...document.querySelectorAll('button, app-button')]
      .filter(e => /create|load/i.test((e.innerText || '').trim()))
      .pop();
    if (btn) btn.click();
  });
  await sleep(7000);

  const text = await pg.evaluate(() => document.body.innerText.replace(/\s+/g, ' '));
  // Reveal addresses if the wallet row needs expanding.
  await pg.evaluate(() => {
    const row = document.querySelector('.wallet, [class*=wallet-row], .-wallet');
    if (row) row.click();
  });
  await sleep(3000);
  const after = await pg.evaluate(() => document.body.innerText.replace(/\s+/g, ' '));

  const found = (after.includes(EXPECT) || text.includes(EXPECT)) && wordCount >= 12;
  const addrs = (after.match(/\b[1-9A-HJ-NP-Za-km-z]{26,35}\b/g) || []).filter(a => a.length >= 26);

  console.log('filled:', JSON.stringify(filled));
  console.log('expected address present:', found);
  console.log('addresses seen:', [...new Set(addrs)].slice(0, 4).join(', ') || '(none)');
  const RAW = /\b(change-coin|pending-txs|qr-code|hardware-wallet)\.[a-zA-Z0-9._-]+/g;
  console.log('raw hyphenated i18n keys:', (after.match(RAW) || []).length);
  console.log('body snippet:', after.slice(0, 200));
  console.log('--- logs ---'); [...new Set(logs)].slice(0, 5).forEach(l => console.log('  ' + l));
  await b.close();
  process.exit(found ? 0 : 1);
})().catch(e => { console.error('harness error:', e.message); process.exit(2); });
