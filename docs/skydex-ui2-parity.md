# skydex-client ui2 parity

The SkyDEX trading client serves two UIs. The React one (`cmd/skydex-client/src`,
built into `commands/static`) is at the root. The server rendered one
(`cmd/skydex-client/skydexui`, templ and htmx) is at `/ui2/`. The React UI
stays the default until `/ui2/` does everything it does. Each PR that ports an
area updates this table.

What the React UI keeps in the browser's localStorage, `/ui2/` keeps there too
under the same keys (`exchange:markets`, `exchange:wallets`,
`exchange:history`), so both UIs share it. The browser posts the list to the
server, which renders it.

present: does what the React UI does. partial: the core is there, the notes
say what is not. missing: not started.

| feature | React | ui2 | notes |
|---|---|---|---|
| connect screen, key pre-filled from `--market`, never auto-connects | ConnectGate | present | |
| connect error shown on the form | ConnectGate | present | |
| recent markets: rejoin, remove one, clear all | ConnectGate, api.js | present | |
| header: status dot, market name, short key, disconnect | Header | present | |
| market name refreshed while connected | App | present | polled with the readiness banner |
| tab bar | App | present | tabs are URLs, so a reload keeps the tab |
| "set your wallet addresses" banner | App | present | |
| "market isn't accepting trades" banner, polled every 8 s | App | present | |
| wallet addresses: Skycoin and one per payment coin, saved in the browser | Settings | present | |
| connected market key | Settings | present | |
| product grid, polled every 5 s | Market | missing | |
| buy, payment instructions kept until dismissed | Market | missing | |
| new sell order form, deposit instructions kept until dismissed | Market | missing | |
| buy or sell sends unregistered users to Settings | Market | missing | |
| my orders table, polled every 5 s | Orders | missing | |
| live status, confirmations n/N, payment tx of buys in flight | Orders | missing | |
| cancel a buy before payment | Orders | missing | |
| my listings table, lifecycle badges, polled every 5 s | MyListings | missing | |
| cancel a listing | MyListings | missing | |
| trade history kept in the browser, remove one, clear all | History | missing | |
| copy to clipboard, also over plain http | CopyButton | missing | |
| QR code with a payment URI and amount, address only toggle | QrButton | missing | |
| local times for expiry and creation | all tables | missing | |

## Beyond the React UI

- Writes from another site are refused, so no other page can drive the client.
- Forms work without JavaScript. Polling and the browser kept lists need it.
- Links both ways: "Classic UI" in the `/ui2/` header, "New UI (preview)" in
  the React UI's corner.
