# One wallet GUI

The node GUI and skycoin-web are two Angular 22 apps that do the same job. This
note records how far apart they are and the plan for replacing both with one.

## Where they are now

| app | source | served by |
|---|---|---|
| node GUI | `src/gui/static` | the node, embedded by `src/gui/static.go` |
| skycoin-web | `src/skycoin-web` | `skycoin web`, and skywire's hypervisor dashboard, which copies the source in as its `#/wallet` route (`scripts/sync-skycoin-wallet.sh` in skywire) |

Measured on develop at the time of writing:

- The node GUI has 255 source files, skycoin-web 193 (specs left out).
- 99 paths exist in both. 77 of those differ in at least 40% of their lines,
  19 moderately, and 3 are near identical.
- 156 files exist only in the node GUI and 94 only in skycoin-web.

## Which one to keep

skycoin-web is the base, for three reasons.

1. It already has the wallet seam the merge needs. Each coin has a
   `serverWallets` flag: set, wallets live behind the node's wallet API (the
   same `wallets`, `wallet/create`, `wallet/newAddress`, `wallet/newAccount`,
   `wallet/xpub`, `wallet/encrypt`, `wallet/decrypt` and
   `wallet/transaction/sign` calls the node GUI makes); unset, they are kept in
   the browser and signed with skycoin-lite. With `--wallet-dir`, `skycoin web`
   runs the same `wallet.Service` the node does.
2. It is already shared by two hosts and is built to be embedded (its own
   module, routes and shell).
3. It has the multicoin, node selection and app password work the node GUI
   lacks.

## What moves from the node GUI

- sending to several destinations, and bulk sends from CSV
- creating and signing unsigned transactions
- the fuller hardware wallet support, firmware update included
- seed backup
- the network page and the exchange page
- the exposed address marks (an address that has sent coins)
- the BIP44 account view fixes and the password error handling (#3122, #3124),
  where skycoin-web has the same faults

## Steps

1. This inventory.
2. Settle the seam: one backend interface for server, browser and hardware
   wallets, where today the choice is spread through `wallet.service.ts`.
3. Port the features above into skycoin-web, one area per PR.
4. Build skycoin-web for the node: one coin, the node's own API, server
   wallets on. `src/gui/static.go` embeds that build, and `src/gui/static`
   goes.
5. Check every flow on the node, `skycoin web` (browser and `--wallet-dir`)
   and the hypervisor dashboard: create, load, encrypt and decrypt wallets,
   deterministic and BIP44, send, history and the hardware wallet.
