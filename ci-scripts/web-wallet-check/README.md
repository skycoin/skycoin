# Web wallet address-derivation check

Loads a known seed into `/src/skycoin-web` in a real browser and asserts the
address it derives, then checks that a fresh seed can be generated.

## Why this exists

The web wallet derives keys in the browser: `generateMnemonic()` from bip39
produces the seed, and `assets/scripts/skycoin-lite.wasm` derives addresses from
it. Neither path is covered by the unit suite, and a break in either is silent —
the UI renders, and the addresses are simply wrong or absent.

That is not hypothetical. Migrating this project to the esbuild builder broke
seed generation outright: bip39 calls `Buffer.from()` inside `mnemonicToSeed()`,
the webpack build supplied that global through `ProvidePlugin`, and esbuild has
no equivalent. The build succeeded, lint passed, and the unit tests passed. Only
clicking through the wallet surfaced `ReferenceError: Buffer is not defined`.

## What it checks

1. **Seed generation** — the onboarding form produces a 12 word mnemonic. This
   exercises bip39, and fails without the `Buffer` polyfill.
2. **Address derivation** — loading the seed below yields the expected address.
   This exercises the wasm cipher end to end.

The vector is the first entry of
`src/skycoin-lite/js/tests/test-fixtures/seed-0001.golden`:

    seed     work pride warrior taxi kick athlete good maze brief address shift creek
    address  2XDHP1JPEun347k6C559npPHT23V5GWeMkA

## Running it

Needs a skycoin node to talk to and a Chrome binary in `CHROME_BIN`.

    node ci-scripts/web-wallet-check/serve.js src/skycoin-web/src/gui/dist 8660 http://127.0.0.1:6420 &

    node ci-scripts/web-wallet-check/check-address-derivation.js 8660 \
      "work pride warrior taxi kick athlete good maze brief address shift creek" \
      2XDHP1JPEun347k6C559npPHT23V5GWeMkA

Exit status is 0 when both checks pass.

`serve.js` is used rather than `skycoin web --gui-dir` on purpose: that command
serves `/assets` from the binary's embedded copy, so the wasm cipher and the
i18n JSON would come from the old bundle no matter what was built, and the check
would pass against a broken build. Everything `serve.js` returns comes from the
directory under test. It also rewrites Host, Origin and Referer on proxied
requests, which the node's checks would otherwise reject with HTTP 403.

## Verified to fail

A check that cannot fail is worth nothing. Truncating
`assets/scripts/skycoin-lite.wasm` to 2 KB in a copy of the bundle makes this
script exit non-zero with

    [WASM] Failed to instantiate WASM module: CompileError ...
    expected address present: false

and removing the `Buffer` polyfill makes the seed generation check report zero
words.

## Not yet wired into CI

It needs a browser and a running node, and unlike the explorer suite it has no
pinned dataset of its own — it only needs the node for connectivity, not for
balances, so a node on the pinned blockchain-180 database is enough. Wiring it
into `ci.yml` alongside `test-explorer-e2e` is the obvious next step.
