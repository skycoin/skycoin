/**
 * Fails if an OnPush component has an asynchronous callback that never marks
 * the view for checking.
 *
 * Under ChangeDetectionStrategy.OnPush a field assigned from a subscription or a
 * timer repaints nothing on its own — the view has to be marked dirty. Forget it
 * and there is no error and no warning, just a page that keeps showing stale
 * data. In a wallet that means a stale balance.
 *
 * The behavioural guards (see richlist.component.spec.ts) prove the contract for
 * the components they cover. This covers the rest: it is a cheap structural
 * sweep over every component in every front-end, so a newly added subscription
 * that forgets to mark fails the build rather than shipping.
 *
 * It is deliberately shallow. It only asks whether a mark appears somewhere in
 * the callback, not whether it is correctly placed, so it can be run over
 * hundreds of components without false alarms.
 *
 * Some callbacks genuinely touch no view state — they navigate, scroll, or hand
 * off to a method that marks for itself. Those opt out with a comment on or just
 * above the callback:
 *
 *     // change-detection: no view state — <reason>
 *
 * which keeps the exemption explicit and reviewable rather than silent.
 *
 * usage: node ci-scripts/check-onpush-marks.js <src-dir> [<src-dir> ...]
 */
const fs = require('fs');
const path = require('path');

const dirs = process.argv.slice(2);
if (!dirs.length) {
  console.error('usage: node ci-scripts/check-onpush-marks.js <src-dir> [<src-dir> ...]');
  process.exit(2);
}

const ASYNC_CALLEES = new Set(['subscribe', 'then', 'setTimeout', 'setInterval']);

function findTypescript(fromDir) {
  let dir = path.resolve(fromDir);
  for (let i = 0; i < 6; i++) {
    const candidate = path.join(dir, 'node_modules', 'typescript');
    if (fs.existsSync(candidate)) return require(candidate);
    dir = path.dirname(dir);
  }
  throw new Error('could not locate the typescript package from ' + fromDir);
}

function walk(dir, acc = []) {
  for (const e of fs.readdirSync(dir, { withFileTypes: true })) {
    const p = path.join(dir, e.name);
    if (e.isDirectory()) walk(p, acc);
    else if (e.name.endsWith('.ts') && !e.name.endsWith('.spec.ts')) acc.push(p);
  }
  return acc;
}

let checked = 0;
const failures = [];

for (const dir of dirs) {
  if (!fs.existsSync(dir)) { console.error('no such directory: ' + dir); process.exit(2); }
  const ts = findTypescript(dir);

  for (const file of walk(dir)) {
    const text = fs.readFileSync(file, 'utf8');
    if (!text.includes('ChangeDetectionStrategy.OnPush')) continue;
    checked++;

    const sf = ts.createSourceFile(file, text, ts.ScriptTarget.Latest, true);

    (function visit(node) {
      if (ts.isCallExpression(node)) {
        const e = node.expression;
        const name = ts.isPropertyAccessExpression(e) ? e.name.getText()
          : ts.isIdentifier(e) ? e.getText() : '';
        if (ASYNC_CALLEES.has(name)) {
          for (const arg of node.arguments) {
            if (!ts.isArrowFunction(arg) && !ts.isFunctionExpression(arg)) continue;
            const body = arg.body.getText();
            if (/\b(markForCheck|detectChanges)\s*\(/.test(body)) continue;

            const { line } = sf.getLineAndCharacterOfPosition(arg.getStart(sf));
            // Allow an explicit, documented exemption on the callback itself or
            // on one of the two lines above it.
            const lines = text.split('\n');
            const context = lines.slice(Math.max(0, line - 2), line + 1).join('\n') + body;
            if (/change-detection:\s*no view state/.test(context)) continue;

            failures.push(`${file}:${line + 1} — ${name}() callback never marks the view`);
          }
        }
      }
      ts.forEachChild(node, visit);
    })(sf);
  }
}

console.log(`checked ${checked} OnPush components`);
if (failures.length) {
  console.error(`\n${failures.length} unmarked asynchronous callback(s):`);
  failures.forEach(f => console.error('  ' + f));
  console.error('\nAdd this.<changeDetectorRef>.markForCheck() at the end of each, or');
  console.error('use the async pipe so Angular marks the view itself.');
  process.exit(1);
}
console.log('every asynchronous callback in an OnPush component marks the view.');
