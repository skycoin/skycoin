# ui2 wallet parity

The node and `skycoin web` serve two wallet GUIs. The server rendered one
(`src/ui2`, templ and htmx) is the default, at the root. The Angular wallet
(`src/skycoin-web`, built with its base at `/angular/`) is kept at `/angular/`,
and the `/ui2/` paths the server rendered one was first served at redirect to
the root. This table tracks how far ui2 is from the Angular wallet, and each
PR that ports an area updates it.

ui2 runs in two modes. With the WALLET API set (`-enable-gui`) the node
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
| bip39 seed passphrase on create and load | create-wallet | present | present | |
| post-quantum (ML-DSA-44) wallet type | create-wallet | present | present | browser: addresses checked against the node, a post-quantum spend is not tested |
| load temporarily | create-wallet | present | partial | browser: kept encrypted in the tab's sessionStorage until it closes, where Angular keeps it unencrypted in memory until a reload |
| address list, balances, hide empty, filter | wallet-detail | present | present | |
| copy address | wallet-detail | present | present | |
| QR code of an address, with a requested amount | qr-code | present | present | drawn with the Angular wallet's qrcode.min.js |
| exposed address marks, move coins | wallet-detail | present | present | |
| exposed addresses banner with a count | wallet-detail | present | present | |
| address history and outputs links | wallet-detail | present | present | |
| new address | wallet-options | present | present | |
| scan for addresses | wallet-options | present | present | |
| add bip44 account, addresses per account | wallet-detail | present | missing | browser wallets keep one account, as the Angular wallet's browser backend does |
| rename | change-name | present | present | |
| encrypt, decrypt | wallet-options | present | partial | by design ui2 never stores a browser wallet decrypted: it sets or changes the password but does not remove it |
| delete | wallet-options | n/a | present | Angular only drops a node wallet from its list |
| account and chain xpubs | wallet-detail | present | present | |
| reset password with the seed | reset-password | present | present | |

## Sending

| feature | Angular | node | browser | notes |
|---|---|---|---|---|
| one destination, review, sign, broadcast | send-form | present | present | |
| wallet password at send | password-dialog | present | present | |
| change address | send-form-advanced | present | present | |
| warnings for spent, reused and own destinations | send-verify | present | present | a ui2 addition |
| several destinations | send-form-advanced | present | present | |
| pick source addresses or outputs | send-form-advanced | present | present | |
| hours allocation (auto share factor, manual) | send-form-advanced | present | present | |
| fee presets slow, normal, fast | send-form | present | n/a | BTC and ETH sends only, see multi-coin |
| send all available coins | send-form | present | present | |
| personal note | send-form | present | present | |
| fill from a `skycoin:` link | enter-link | present | present | |
| scan a QR code | qr-scan-dialog | present | present | camera or image, decoded with vendored jsQR 1.4.0; the camera path is untested here, the image path is |
| USD value | send-form, history | missing | missing | price source on hold |
| bulk send from CSV | bulk-send | present | present | lines pasted into the form |
| unsigned transaction | offline-tx | present | present | |
| sign a transaction | offline-tx | present | present | |
| broadcast a raw transaction | offline-tx | present | present | the transactions page needs no wallet to broadcast |
| joint transactions, create and co-sign | joint-transaction | present | present | |

## History

| feature | Angular | node | browser | notes |
|---|---|---|---|---|
| history of one wallet | history | present | present | on the wallet page |
| history of every wallet | history | present | present | |
| filter by wallet or address | history | present | present | ui2 also searches by transaction ID |
| transaction detail | transaction-detail | present | present | inputs, outputs, hours sent and burned |
| notes on transactions | history | present | present | shared with the Angular wallet through localStorage |
| cached copy when the node does not answer | history | present | present | the last good unsearched first page, kept in this browser |

## Settings and status

| feature | Angular | node | browser | notes |
|---|---|---|---|---|
| sync progress, height, last block | header, blockchain | present | present | ui2 status page |
| coin supply | blockchain | present | present | on the status page |
| peers | network | present | present | on the status page |
| pending transactions, mine or all | pending-transactions | present | present | |
| remove a pending transaction from the pool | pending-transactions | present | present | |
| unspent outputs | outputs | present | present | |
| backup, show seed | backup | present | present | node: needs INSECURE_WALLET_SEED; browser: on the wallet page, decrypted in bw.js |
| wallet directory, save addresses | backup | present | missing | |
| node URL per coin | nodes | n/a | missing | ui2 reads the nodes of the host that serves it; Angular can point the browser at another node |
| Ethereum RPC servers and indexer | nodes | present | n/a | skycoin web |
| app password set, change, remove | app-password | n/a | present | the wallets inside stay encrypted with their own passwords |
| lock screen, erase | app-lock | n/a | present | |
| custom coins: fiber coins and ERC-20 tokens | coins | present | n/a | skycoin web with --wallet-dir |
| Skywallet | hardware-wallet | missing | missing | |
| onboarding wizard on first run | onboarding | present | present | the wallet list opens on a create or load choice while there are none; the password is set in the create form, not a second step |
| node errors: unreachable, CSRF, no peers, syncing, pending | header | present | present | a strip under the top bar |
| language | select-language | present | present | |
| page outside a secure context warning | wallets | n/a | present | |

## Multi-coin

Multi-coin is skycoin web's: the coins it lists at `/api/v1/coins`, with Bitcoin
and Ethereum wallets kept by its wallet service (`--wallet-dir`). In ui2 the
"node" column is those server wallets; Angular keeps no Bitcoin or Ethereum
wallet in the browser either. The Bitcoin and Ethereum pages were checked only
with the repo's Go fakes (`cmd/skycoin-web/commands/ui2_coins_test.go`): a fake
`btc.Backend` and an httptest JSON-RPC node and indexer. A live check needs an
Electrum server or Bitcoin Core (testnet or regtest) and an Ethereum JSON-RPC
node with an Etherscan-style indexer, such as a local devnet.

| feature | Angular | node | browser | notes |
|---|---|---|---|---|
| coin picker | select-coin | present | n/a | |
| BTC wallets: create, load, legacy or segwit, balance, history, send with a fee preset | electrum backend | present | n/a | checked against a fake backend only |
| ETH and ERC-20 wallets: create, load, balance, history, send with a fee preset | eth backend | present | n/a | checked against a fake RPC node and indexer only |

## Out of scope

The buy page used swaplab.cc, whose API no longer answers.
