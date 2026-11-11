# ui2 wallet parity

The node serves two wallet GUIs: the Angular wallet (`src/skycoin-web`) at the
root, and the server rendered one (`src/ui2`, templ and htmx) at `/ui2`. The
Angular wallet stays the default until `/ui2` does everything it does. This
table tracks how far that is, and each PR that ports an area updates it.

`/ui2` runs in two modes. With the WALLET API set (`-enable-gui`) the node
keeps the wallets ("node" below). Without it (`-web-wallet`) the wallets are
kept in the browser in the Angular wallet's own format, and `static/bw.js`
signs with the wasm cipher ("browser" below).

present: does what the Angular page does. partial: the core is there, the
notes say what is not. missing: not started. n/a: does not apply in that mode.

## Wallets

| feature | Angular | node | browser | notes |
|---|---|---|---|---|
| wallet list, total balance | wallets | present | present | |
| create, 12 or 24 word seed, deterministic or bip44 | create-wallet | present | present | |
| load from seed, address scan | create-wallet | present | present | |
| non standard seed | create-wallet | present | present | |
| bip39 seed passphrase on load | create-wallet | missing | present | |
| load temporarily | create-wallet | missing | missing | node `wallet/createTemp` |
| address list, balances, hide empty, filter | wallet-detail | present | present | |
| copy address | wallet-detail | present | present | |
| QR code of an address | qr-code | missing | missing | |
| exposed address marks, move coins | wallet-detail | present | present | |
| exposed addresses banner with a count | wallet-detail | present | present | |
| address history and outputs links | wallet-detail | missing | missing | |
| new address | wallet-options | present | missing | |
| scan for addresses | wallet-options | present | missing | |
| add bip44 account | wallet-detail | missing | missing | |
| rename | change-name | present | missing | |
| encrypt, decrypt | wallet-options | present | missing | |
| delete | wallet-options | n/a | missing | Angular only drops a node wallet from its list |
| account and chain xpubs | wallet-detail | present | present | |
| reset password with the seed | reset-password | missing | missing | |

## Sending

| feature | Angular | node | browser | notes |
|---|---|---|---|---|
| one destination, review, sign, broadcast | send-form | present | present | |
| wallet password at send | password-dialog | present | present | |
| change address | send-form-advanced | present | present | |
| warnings for spent, reused and own destinations | send-verify | present | present | a ui2 addition |
| several destinations | send-form-advanced | missing | missing | |
| pick source addresses or outputs | send-form-advanced | missing | missing | |
| hours allocation (auto share factor, manual) | send-form-advanced | missing | missing | ui2 uses auto, share 0.5 |
| fee presets slow, normal, fast | send-form | missing | missing | |
| send all available coins | send-form | missing | missing | |
| personal note | send-form | missing | missing | |
| fill from a `skycoin:` link | enter-link | missing | missing | |
| scan a QR code | qr-scan-dialog | missing | missing | |
| USD value | send-form, history | missing | missing | price source on hold |
| bulk send from CSV | bulk-send | missing | missing | |
| unsigned transaction | offline-tx | missing | missing | |
| sign a transaction | offline-tx | missing | missing | |
| broadcast a raw transaction | offline-tx | missing | missing | |
| joint transactions | joint-transaction | missing | missing | |

## History

| feature | Angular | node | browser | notes |
|---|---|---|---|---|
| history of one wallet | history | present | present | on the wallet page |
| history of every wallet | history | present | present | |
| filter by wallet or address | history | present | present | ui2 also searches by transaction ID |
| transaction detail | transaction-detail | present | present | inputs, outputs, hours sent and burned |
| notes on transactions | history | present | present | shared with the Angular wallet through localStorage |
| cached copy when the node does not answer | history | missing | missing | |

## Settings and status

| feature | Angular | node | browser | notes |
|---|---|---|---|---|
| sync progress, height, last block | header, blockchain | present | present | ui2 status page |
| coin supply | blockchain | missing | missing | |
| peers | network | partial | partial | no last sent, received and seen times |
| pending transactions list | pending-transactions | missing | missing | status shows the count |
| unspent outputs | outputs | missing | missing | |
| backup, show seed | backup | missing | missing | node needs INSECURE_WALLET_SEED |
| wallet directory, save addresses | backup | missing | n/a | |
| node URL per coin | nodes | n/a | missing | ui2 reads the node that serves it |
| ETH and BTC endpoints | nodes | missing | missing | |
| app password set, change, remove | app-password | n/a | partial | unlock only |
| lock screen, erase | app-lock | n/a | partial | unlock only |
| custom coins, ERC-20 | coins | missing | missing | |
| Skywallet | hardware-wallet | missing | missing | |
| onboarding wizard on first run | onboarding | missing | missing | the create page is there |
| node errors: unreachable, CSRF, no peers | header | partial | partial | unreachable only |
| language | select-language | present | present | |
| page outside a secure context warning | wallets | n/a | missing | |

## Multi-coin

| feature | Angular | node | browser | notes |
|---|---|---|---|---|
| coin picker | select-coin | missing | missing | |
| BTC wallets | electrum backend | missing | missing | |
| ETH and ERC-20 wallets | eth backend | missing | missing | |

## Out of scope

The buy page used swaplab.cc, whose API no longer answers.
