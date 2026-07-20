# SkyDEX — Architecture Design

## Overview
SkyDEX is built on top of the Skywire mesh network and ecosystem. To preserve
privacy, security, and decentralization, it is implemented as **two separate
applications (binaries)**, each running on its own Visor:

- **`skydex-client`** — the user-facing app (buyers and sellers).
- **`skydex-market`** — the backend engine (orderbook, escrow, blockchain verification)
  plus a small local operator UI for configuring and monitoring the market.

A client connects to a market over **dmsg** using the market's public key. Anyone can
run their own `skydex-market` to operate an independent market; the team also runs one
official market.

### Where the code lives (repo split)
SkyDEX follows the same split as skycoin-web: the **engine lives in the skycoin repo** and
carries no skywire dependency, and **skywire provides only the transport wrapper**.

- **skycoin** (`github.com/skycoin/skycoin`) — `src/skydex/*` is the engine (protocol,
  message framing, SQLite store, chain/escrow, background jobs, and the transport-agnostic
  server + client cores); `cmd/skydex-market` / `cmd/skydex-client` are standalone (TCP)
  reference commands plus the operator/auth UI and the React trading client. All trading
  logic lives here and builds/runs over plain TCP with no skywire import.
- **skywire** (`github.com/skycoin/skywire`) — `cmd/apps/skydex-{market,client}` are thin
  visor apps that supply the dmsg transport (appnet Listen/Dial), inject the authenticated
  public key as identity, and hand the connection to the skycoin engine. Plus
  `internal/skydex-{market,client}/app` (visor `app.Client` bridges) and
  `internal/skydex-e2e` (the two-visor dmsg wire test).

The engine is **transport-agnostic**: the server exposes `Accept(ctx, net.Listener,
IdentifyFunc)` + `ServeConn(conn, identity)`; the client exposes `NewConn(net.Conn)`. The
same code serves over a plain TCP listener (standalone) or a skywire appnet/dmsg listener
(the visor app). Over skywire the injected identity is the client's authenticated dmsg
public key; the operator UI reaches its host (the visor) through a small `Host` interface
(logger, operator-OTP publisher, market public key).

> **Scope note:** The market sells **SKY and any operator-added Skycoin fibercoin** (see §7
> `sell_coins`), paid for with an **external UTXO coin** (BTC/LTC, and others the operator
> enables) **or another fibercoin**. The commission is charged in **the sell coin** (see §2).
> SCH (Skycoin Coin Hours) is **not** traded and is **not** the commission unit — any earlier
> mention of trading SCH or charging the fee in SCH is obsolete. (Some examples below are
> written SKY-for-BTC for concreteness; they generalize to any configured pair.)

---

## 1. SkyDEX - Client (`skydex-client`)
Acts as the user interface for **buyers and sellers**. Every participant runs this app on
their own Visor. A single user can act as **both** buyer and seller at the same time.

### User Interface (UI)
- **Type:** Single Page Application (SPA), served in a web browser.
- **Address:** `http://localhost:8051` (the client app's `--addr`, default `:8051`)
- **Color palette:**
  - **Navy (primary background):** `#101F34`
  - **White (text/accent):** `#FFFFFF`
  - **Blue (action/highlight):** `#0273FF`
- **Status bar (top of page):**
  - Connection status to the market (connected / disconnected)
  - Public key of the connected market
  - Activity/processing indicator
- **Language:** English only. (The UI, the market, and all artifacts are in English.)
- **Embedding:** The built SPA is embedded into and served by the Go `skydex-client`
  binary via `embed.FS`. All UI dependencies are bundled locally — nothing is fetched from
  the internet at runtime.
- **Connect gate:** On open, the UI shows a single input pre-filled with the configured
  `--market-pk` value plus a **Connect** button. The client never connects automatically;
  only on Connect does it dial the market over dmsg. The main screen (below) appears once
  connected.
- **Local control API:** The Go binary exposes a localhost API the SPA calls:
  `GET /api/config` (default market pk), `GET /api/status`, `POST /api/connect` /
  `POST /api/disconnect`, and thin proxies that forward to the connected market over dmsg —
  `POST /api/register`, `GET /api/currencies`, `GET /api/products`, `POST /api/listings`,
  `GET /api/listings/mine`, `POST /api/listings/cancel`, `POST /api/buy`, `GET /api/orders`,
  `POST /api/orders/cancel`, `POST /api/order-status`.
  The SPA polls products/orders/listings every 5 s. For each of the buyer's in-flight buy
  orders it also polls `POST /api/order-status` to show live confirmation progress (n/2) and
  the payment tx hash as they appear.

### The main screen has 4 sections:

#### a) Market
- **Product list:** Displays items (SKY tokens) that sellers have listed for sale.
  (Phase 1 trades **SKY only**.)
- **Buyer flow:**
  1. The buyer picks an item from the available list (e.g. 10 SKY for 2 USD).
  2. **Product freeze:** As soon as the buyer selects an item, it is temporarily frozen —
     no other buyer can select it.
  3. The buyer confirms the seller's price.
  4. **15-minute window:** After confirmation, the buyer has exactly 15 minutes to pay.
  5. **Non-round-number validation:** The system shows the buyer a non-round amount
     (e.g. 2.00034123 USD) to pay. This makes the transaction easy to identify and verify
     precisely in the seller's wallet.
  6. **Confirmation & transfer:** Once the payment is observed in the seller's wallet and
     reaches the required confirmations, the SKY tokens are transferred to the buyer's
     wallet.
  7. **Cancel (buyer):** Before payment is observed the buyer can cancel the buy, which
     releases the product back to the market. A buyer who cancels **cannot buy that same
     product again** (`client.cancel_order`; the market rejects a re-buy with `BUYER_BLOCKED`).
     A cancel also records a `freeze_violation`, exactly like letting the order expire, so a
     buyer who repeatedly freezes then backs out is banned by the Ban Manager (§5).
- **Seller flow:**
  1. The seller creates a new listing (e.g. "I sell 10 SKY for X BTC").
  2. **Non-round-number validation:** The system tells the seller a non-round amount
     (e.g. 10.00045678 SKY) to transfer, so the deposit is easy to identify and verify
     precisely in the **market wallet**.
  3. **15-minute window:** After creating the listing, the seller has 15 minutes to
     transfer the SKY tokens to the **market wallet** (escrow).
  4. **Activation:** Once the transfer is confirmed into the market wallet, the listing
     becomes an active product visible to buyers.
  5. **Cancel (seller):** The seller can cancel their offer while it is pending (no deposit
     yet) or active (a product), **as long as no buyer has selected it** (`client.cancel_listing`;
     canceling a frozen/sold offer is rejected). A confirmed offer's escrowed SKY is returned
     immediately by the Return Scheduler (see §7.8/4).

#### b) My Listings (seller)
- The seller's own sell listings and their live lifecycle: **pending** (awaiting the SKY
  deposit) → **confirmed** (deposit detected on-chain) → **listed as a product**.
- For a still-pending listing it restates the exact non-round deposit amount and the market
  wallet to send it to; once confirmed it shows the deposit tx hash.
- Backed by `GET /api/listings/mine` (`client.get_listings`), polled every 5 s.

#### c) My Orders
- A full list of the user's buy and sell orders.
- Status of each order (pending payment, processing, completed, cancelled, expired).
- For the buyer's own in-flight buy orders, a **live progress bar** of confirmations (n/2)
  and the payment tx hash, refreshed by polling `POST /api/order-status` every 5 s.

#### d) Settings
- **Wallet addresses:**
  - **SKY** wallet (receive/send SKY tokens)
  - **BTC** (Bitcoin)
  - **BCH** (Bitcoin Cash)
  - **LTC** (Litecoin)
  - **USDT** — both **TRC20** and **ERC20**
- **Market public key:**
  - The client connects to the market over **dmsg**.
  - The user enters the target market's public key here to establish the secure link.

---

## 2. SkyDEX - Market (`skydex-market`)
The processing engine and "brain" of the exchange. Anyone can run their own instance on
their Visor to operate an independent market.

### Operator UI
The market ships a small **operator UI** (served locally by the market binary, embedded the
same way as the client) for the market owner to configure and monitor their market. It is
an operator tool, not a trading interface. It provides:

- **Configuration**
  - **Explorer URLs** per payment coin (BTC, BCH, LTC, USDT-ERC20, USDT-TRC20). An explorer
    left blank means that coin is **disabled** — buyers and sellers cannot trade SKY against
    it. (See "Explorer-gated currencies" in §3.)
  - **SKY fullnode** API URL (native SKY verification and transfers).
  - **Market SKY wallet** address (escrow wallet that holds sellers' SKY).
  - **Ban threshold** — the number of unproductive freezes before a temporary ban
    (`freeze_violations_limit`), and ban duration.
  - Commission (rate %, floor, cap), confirmations required, min/max trade size, listing/order
    expiry, cleanup days, per-coin explorers — all in `market_config`, editable in the UI.
- **Monitoring**
  - All active products (goods sellers listed) with their status (`active`, `frozen`/in-process,
    `sold`, …).
  - All orders and their live status.
  - Currently (temporarily) banned users and their remaining ban time.

The trading logic, blockchain interaction, escrow, and order management all run here; the
UI only reads/writes the market's own config and state.

The operator UI is served on the market's `--addr` (default `:8050`) and backed by a local
operator API (localhost only, separate from the client↔market dmsg protocol):
`GET/POST /api/config` (edit escrow wallet, SKY fullnode, explorers, ban rules, timeouts),
`GET /api/info` (market public key), and read-only `GET /api/products`, `/api/orders`,
`/api/bans`.

### Responsibilities
- **Orderbook management:**
  - Receive, store, and process buy/sell orders from clients.
  - Match orders and reply to clients.
- **Native blockchain integration:**
  - Direct, native connection to the **SKY** blockchain.
  - Execute token transfers between wallets when a trade completes.
- **Explorer integration (external chains):**
  - Connect to explorers for **USDT, BTC, LTC, BCH**.
  - Verify the off-network payments buyers make (the exact reason/mechanism is designed in
    later phases).
- **Escrow:**
  - The seller's SKY tokens are held in the market wallet until the trade completes.
  - **Immediate return policy:** In any of the following cases the tokens are returned to the
    seller's wallet as soon as the escrow can fund the refund transaction (the Return Scheduler
    retries if coin-hours are still accruing — see §7.8/4). No fixed delay.
    - The seller explicitly cancels the trade (at any point in the flow).
    - The product does not sell before its TTL/expiry.
    - Any other reason the trade fails to complete.
- **Commission:**
  - Charged to the **seller only**, in **SKY** (real, monetizable revenue — coin-hours are too
    low-value and their amount is uncontrollable, so they're not used as the fee unit).
  - **Formula:** `commission = clamp(amount_sky × rate%, min, cap)`, rounded up to SKY's 3-decimal
    precision. Defaults: `commission_rate_percent = 0.5` (0.5 %), `commission_min_sky = 0.001`
    (floor, so tiny trades still pay something), `commission_max_sky = 0` (no cap; a positive
    value caps it). All three are operator-tunable in the market UI.
  - **How it's collected:** the seller deposits `amount + commission` into escrow; the market
    delivers `amount` to the buyer and **retains the commission** in the escrow wallet. The
    commission is booked on each completed sale (`orders.commission_sky`, shown in the operator
    dashboard). The `CreateListingResponse` returns the breakdown (`amount_sky` + `commission_sky`
    = `expected_amount_sky`) so the client shows the seller exactly what to deposit. Refunds
    return the **full** deposited SKY (amount + commission) and are **not** charged commission
    (no sale occurred). The network burn fee (the ≥10 % coin-hour burn a SKY tx pays) is separate
    and paid from the escrowed SKY's accrued coin-hours.
- **Storage & privacy:**
  - **SQLite** as a lightweight embedded database.
  - **Auto-deletion:** All data related to a successful trade is automatically deleted from
    the database **3 days after completion** (to protect user privacy and security).

---

## 3. Communication Architecture & Protocol

### Client ↔ Market transport
- **Transport:** All client↔market communication happens over **dmsg** (Skywire's
  messaging layer).
- **Benefits of dmsg:**
  - End-to-end encrypted communication
  - Uses the Skywire mesh infrastructure
  - No public IP or centralized server required
  - Censorship/firewall resistant
- **Real-time updates:**
  - The client uses **simple polling**.
  - It periodically (e.g. every 5–10 s) asks the market for the status of its orders.

### Required confirmations
The market's threshold is the single source of truth: it defaults to `jobs.RequiredConfirmations`
(2) but is **operator-configurable** via `confirmations_required`, and is returned in each
`get_order_status` response (`required_confirmations`), so the client's progress bar shows `n/N`
without hardcoding `N` (it only falls back to a local default if an older market omits the field).
Likewise the client dials the market on `--market-port`, which config-gen sources from the same
`skyenv.SkydexMarketPort` constant as the market's routing port — so the client always dials
where the market listens.

The market waits for `confirmations_required` blockchain confirmations (default **2**, applied to
both the SKY deposit depth and the buyer's payment) before completing a trade. Two is acceptable
for small/medium, low-risk trades; an operator wanting more security (3–6 is typical for BTC) can
raise it in the market UI.

### Per-coin explorer configuration
The buyer pays in an external coin; SKY is always the product. A payment currency (from the
supported set `protocol.PaymentCurrencies` = BTC, BCH, LTC, DOGE, DASH — native coins
verifiable by address) becomes available **only when the operator configures an explorer
for it**. Per coin the operator sets, in the market UI:

- `explorer_<coin>_provider` — an explorer adapter the market implements (e.g. `esplora`),
  or empty to disable the coin;
- `explorer_<coin>_url` — explorer base URL (optional for coins with a built-in default —
  currently BTC/LTC — required for the rest);
- `explorer_<coin>_key` — optional API key.

The `esplora` adapter covers all supported UTXO coins (BTC, LTC, BCH, DOGE, DASH). BTC and LTC
ship with a default public endpoint and can be enabled with one tick; BCH/DOGE/DASH are
disabled by default and need the operator to supply an Esplora-compatible URL. The market UI
presents a per-coin **enable/disable checkbox** (enabling sets the provider, disabling clears
it while keeping the URL/key), so the operator can add, remove, or disable trade coins there;
enabling a coin with no default and no URL is rejected. Availability = a supported coin with a
configured provider. The client discovers the set via `client.get_currencies`; the market
rejects listings/buys in an unconfigured currency with `CURRENCY_UNAVAILABLE`. Tokens (USDT)
and privacy coins (XMR) are absent: tokens need per-contract endpoints, XMR can't be verified
by address.

### Payment verification model
No transaction hash is required from users. The market auto-confirms by scanning the recipient
address's recent transactions. Because SKY and the external payment coins have very different
on-chain precision, the two sides use different identifiers:

**SKY deposits (seller → shared escrow wallet).** Skycoin only accepts amounts to 3 decimals
(`MaxDropletPrecision = 3`), so a fine-grained unique-amount delta is impossible. Instead the
seller deposits the **exact round amount** and the deposit is matched by:
- **sender address** — the tx's input owner must be the seller's registered SKY address;
- **time window** — the tx's block time must fall in `[listing created, expires + grace]`
  (the lower bound blocks replay of an old, unrelated transaction);
- **exact amount** to the escrow wallet.

To keep this unambiguous, a seller may hold **only one pending listing at a time** (a second
`create_listing` is rejected until the first is deposited/canceled/expired) — so a single
deposit from that address maps to exactly one listing. The seller must deposit
`amount + commission` (see §2 commission); the market delivers `amount` and retains the
commission.

**External payments (buyer → seller's own address).** BTC/LTC/… have 8 decimals, so the buyer
is told to pay a **unique non-round amount** (e.g. `10.000000xy`), allocated to be unique among
concurrent pending payments to that seller wallet (allocate → check-DB → regenerate, under a
lock). A confirmed payment is matched by exact amount **within the order's time window**
(anti-replay of an old payment of a recycled amount); the buyer's registered address is a soft
corroborating signal (a tiebreaker), not a hard requirement, since wallets auto-select inputs.

**Precision & bounds.** `amount_sky` is normalized to SKY's spendable precision (**3 dp**);
`price` to the payment coin's (8 dp). An amount rounding to zero is rejected. A float-safety
ceiling (`amount_sky ≤ 1e9`, `price ≤ 1e7`) keeps `base × 10^dp` within float64's exact-integer
range (2^53). The operator can additionally set `min_trade_sky` / `max_trade_sky` (0 = unset).
External-coin amounts are summed in integer base units (satoshis) with a single final division,
avoiding float accumulation error.

### Blockchain backend (`src/skydex/chain`)
- **SKY deposit verification (implemented, live-validated):** a Skycoin node client reads the
  market wallet's confirmed transactions (`GET /api/v1/transactions?addrs=…&verbose=1`) and
  matches a deposit by exact amount + confirmation depth. This is a **read-only** call, so it
  works against the **public node `https://node.skycoin.com`** — no local fullnode required.
  Set `sky_fullnode_url = https://node.skycoin.com`. (Validated live in `sky_live_test.go`,
  `-tags livenet`.) The response read cap is 32 MiB so a busy escrow wallet's history isn't
  truncated.
- **SKY spend / delivery (implemented — local signing):** the public node's **wallet API is
  disabled** (`/api/v1/wallet/*` → 403), so the market **signs locally** in-process: it reads
  UTXOs via `GET /api/v1/outputs`, builds and signs the delivery/refund tx with `src/coin` +
  `src/cipher` + `src/transaction`, and broadcasts via `POST /api/v1/injectTransaction`. The
  escrow key is the per-coin `wallet_seed` (SKY's is `sky_wallet_seed`) and never leaves the
  market. Live-verified on mainnet against a self-hosted and the public SKY node.
- **External coins (pluggable adapters):** an `Explorer` interface + a per-currency
  **router** that dispatches each coin to the operator-configured adapter. Implemented
  adapters live in a registry (`ProvidersFor`, `SupportsProvider`).
  - **`esplora`** (implemented) — mempool.space (BTC) / litecoinspace (LTC); free, no key.
  - Future: `blockcypher` (DOGE/DASH), `etherscan`/`trongrid` (USDT), `3xpl` (many).
- The real backend is wired into the jobs when `sky_fullnode_url` is set; otherwise the
  market stays fully no-op. A coin with no configured explorer never auto-confirms.

### Markets per client
- **Single market:** A client connects to exactly **one** market at a time.
- The user can change the market public key in settings to switch to another market.

### Session management (single-session enforcement)
To avoid race conditions and keep data consistent, only one active session per public key
is allowed:
- **One session rule:** Each public key (the user's network identity) may be connected from
  **one device** at a time.
- **Reject duplicates:** If a public key is already connected to the market, any new
  connection attempt with the same public key is **rejected**.
- **Benefits:**
  - Prevents race conditions (e.g. the same user clicking one product from multiple tabs)
  - Simplifies session management
  - Ensures each user trades from a single point

---

## 4. Identity & Privacy
All users (buyers and sellers) are fully anonymous; no personally identifiable information
(PII) is stored.

- **Public-key identity:** Each user is identified solely by their Visor **public key**.
- **No authentication:** No name, email, phone, or other identity data is recorded.
- **Privacy-first:** The architecture aligns with Skywire's privacy-first philosophy —
  users trade without revealing their identity.

---

## 5. Violation & Ban System
To prevent abuse and ensure trades are serious, an automatic ban system is implemented:

1. **Freeze without purchase:** If a buyer freezes a product but does not pay within
   15 minutes, one "violation" is recorded against them.
2. **Three-strike rule:** If a buyer (by public key) freezes a product without completing
   the purchase **3 times**, they are **banned for 1 week**.
3. **Ban period:** During the 1-week ban, the user cannot buy any product.
4. **Storage:** The banned-users list is kept locally in the market and the user is removed
   automatically when the ban expires.

---

## 6. Design Philosophy — Why two apps?
Splitting into "client" and "market" serves these goals:
- **Real decentralization:** Anyone, anywhere can run `skydex-market` on their own
  Skywire node and operate an independent exchange with their own rules and fees.
- **Scalability:** Traffic and processing load is distributed across markets.
- **Official market:** Alongside personal markets, the team runs one (or more) official
  markets for public use.
- **Security:** Sensitive financial logic (blockchain and wallet interaction) is separated
  from the client UI and only runs on the market side.

---

## 7. Database Design
The market uses **SQLite** as an embedded database. Tables are designed to store the
minimum necessary data and to allow automatic deletion after 3 days.

### 7.1. `users`
Users (buyers and sellers) and their wallet addresses.

```sql
CREATE TABLE users (
    pubkey          TEXT PRIMARY KEY,           -- Visor public key (user identity)
    wallet_sky      TEXT NOT NULL,              -- SKY wallet address
    wallet_btc      TEXT,                       -- BTC wallet address
    wallet_bch      TEXT,                       -- BCH wallet address
    wallet_ltc      TEXT,                       -- LTC wallet address
    wallet_usdt_erc20 TEXT,                     -- USDT (ERC20) wallet address
    wallet_usdt_trc20 TEXT,                     -- USDT (TRC20) wallet address
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_updated ON users(updated_at);
```

### 7.2. `pending_listings`
When a seller creates a listing, it stays here until they transfer the SKY to the market
wallet.

```sql
CREATE TABLE pending_listings (
    id                  TEXT PRIMARY KEY,       -- UUID
    seller_pubkey       TEXT NOT NULL,          -- seller public key
    amount_sky          REAL NOT NULL,          -- SKY amount the buyer will receive
    expected_amount_sky REAL NOT NULL,          -- exact round amount to deposit = amount + commission
    price               REAL NOT NULL,          -- total price (e.g. 2 USD)
    payment_currency    TEXT NOT NULL,          -- BTC, BCH, LTC, USDT_ERC20, USDT_TRC20
    status              TEXT NOT NULL,          -- pending, confirmed, expired, cancelled
    expires_at          DATETIME NOT NULL,      -- expiry (15 min after creation)
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    confirmed_at        DATETIME,               -- when the transfer was confirmed
    tx_hash             TEXT,                   -- hash of the transfer to market
    FOREIGN KEY (seller_pubkey) REFERENCES users(pubkey)
);

CREATE INDEX idx_pending_listings_seller ON pending_listings(seller_pubkey);
CREATE INDEX idx_pending_listings_status ON pending_listings(status);
CREATE INDEX idx_pending_listings_expires ON pending_listings(expires_at);
```

### 7.3. `products`
Confirmed products, visible to buyers.

```sql
CREATE TABLE products (
    id                  TEXT PRIMARY KEY,       -- UUID
    seller_pubkey       TEXT NOT NULL,          -- seller public key
    amount_sky          REAL NOT NULL,          -- SKY amount for sale
    price               REAL NOT NULL,          -- total price
    payment_currency    TEXT NOT NULL,          -- payment currency
    status              TEXT NOT NULL,          -- active, frozen, sold, expired, cancelled
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    frozen_at           DATETIME,               -- when frozen (selected by a buyer)
    frozen_by           TEXT,                   -- public key of the buyer who froze it
    sold_at             DATETIME,               -- when sold
    FOREIGN KEY (seller_pubkey) REFERENCES users(pubkey)
);

CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_seller ON products(seller_pubkey);
```

### 7.4. `orders`
Created when a buyer selects a product.

```sql
CREATE TABLE orders (
    id                      TEXT PRIMARY KEY,       -- UUID
    product_id              TEXT NOT NULL,          -- selected product
    buyer_pubkey            TEXT NOT NULL,          -- buyer public key
    amount_sky              REAL NOT NULL,          -- SKY amount being bought
    price                   REAL NOT NULL,          -- total price
    payment_currency        TEXT NOT NULL,          -- payment currency
    expected_payment_amount REAL NOT NULL,          -- non-round amount buyer must pay (validation)
    seller_wallet           TEXT NOT NULL,          -- seller wallet to receive payment
    status                  TEXT NOT NULL,          -- pending_payment, paid, confirmed, completed, expired, cancelled
    expires_at              DATETIME NOT NULL,      -- expiry (15 min after freeze)
    created_at              DATETIME DEFAULT CURRENT_TIMESTAMP,
    paid_at                 DATETIME,               -- when payment was detected
    payment_tx_hash         TEXT,                   -- buyer payment tx hash
    confirmations           INTEGER DEFAULT 0,      -- current confirmations
    completed_at            DATETIME,               -- when completed (SKY sent to buyer)
    commission_sky          REAL DEFAULT 0,         -- SKY commission retained on the completed sale
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (buyer_pubkey) REFERENCES users(pubkey)
);

CREATE INDEX idx_orders_buyer ON orders(buyer_pubkey);
CREATE INDEX idx_orders_product ON orders(product_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_expires ON orders(expires_at);
```

### 7.5. `freeze_violations`
Records unproductive freezes for counting violations and enforcing bans.

```sql
CREATE TABLE freeze_violations (
    id              TEXT PRIMARY KEY,       -- UUID
    buyer_pubkey    TEXT NOT NULL,          -- offending buyer public key
    order_id        TEXT NOT NULL,          -- the expired order
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (buyer_pubkey) REFERENCES users(pubkey),
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

CREATE INDEX idx_violations_buyer ON freeze_violations(buyer_pubkey);
CREATE INDEX idx_violations_created ON freeze_violations(created_at);
```

### 7.6. `bans`
Banned users and ban expiry.

```sql
CREATE TABLE bans (
    pubkey      TEXT PRIMARY KEY,       -- banned user public key
    violations  INTEGER NOT NULL,       -- violations that led to the ban
    ban_until   DATETIME NOT NULL,      -- ban end (1 week)
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pubkey) REFERENCES users(pubkey)
);

CREATE INDEX idx_bans_until ON bans(ban_until);
```

### 7.7. `market_config`
Internal market configuration (wallet addresses, fee rate, etc.).

```sql
CREATE TABLE market_config (
    key     TEXT PRIMARY KEY,
    value   TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Initial values (defaults; all editable in the market UI)
INSERT INTO market_config (key, value) VALUES
    ('wallet_sky', ''),                  -- escrow SKY address (= sky_wallet_seed's first address)
    ('sky_wallet_seed', ''),             -- escrow hot-wallet seed; spends are signed locally
    ('sky_fullnode_url', ''),            -- Skycoin node base URL (reads + broadcast)
    ('commission_rate_percent', '0.5'),  -- commission: % of SKY sold
    ('commission_min_sky', '0.001'),     -- commission floor (SKY)
    ('commission_max_sky', '0'),         -- commission cap (SKY; 0 = none)
    ('confirmations_required', '2'),     -- blockchain confirmations before completing a trade
    ('min_trade_sky', '0'),              -- minimum listing amount (0 = none)
    ('max_trade_sky', '0'),              -- maximum listing amount (0 = safety ceiling only)
    ('listing_expiry_minutes', '15'),    -- pending-listing TTL
    ('order_expiry_minutes', '15'),      -- pending-order TTL
    ('cleanup_days', '3'),               -- delete completed trades after N days
    ('freeze_violations_limit', '3'),    -- violations before ban
    ('ban_duration_days', '7');          -- ban duration (days)
-- Per-coin explorer keys (explorer_<coin>_{provider,url,key}) are created on demand.
-- The wallet_sky/sky_wallet_seed/sky_fullnode_url keys are the legacy single-coin
-- escrow config; they now back the seeded SKY row of sell_coins (below) and remain
-- as a fallback for pre-fibercoin deployments.
```

#### `sell_coins` — the sell-side coin registry (fibercoins)

The market sells **SKY and any Skycoin fibercoin**. Because every Skycoin-family
coin shares the same address format, a user registers a **single** Skycoin address
(`users.wallet_sky`) that receives *every* sell coin they buy (and any refunds);
only the **escrow** side is per-coin. Each sell coin is a row with its own fullnode
URL and escrow hot wallet (seed + address) and confirmation depth. Operators
**add / edit / enable / disable / delete** coins in the UI — mirroring how payment
explorers are managed. SKY is seeded as the default row (bridging the legacy keys
above) but is otherwise an ordinary, removable row.

```sql
CREATE TABLE sell_coins (
    symbol        TEXT PRIMARY KEY,   -- ticker: SKY, MDL, …
    name          TEXT NOT NULL DEFAULT '',
    node_url      TEXT NOT NULL DEFAULT '',   -- this coin's fullnode base URL
    wallet_seed   TEXT NOT NULL DEFAULT '',   -- escrow hot-wallet seed (signed locally)
    wallet_addr   TEXT NOT NULL DEFAULT '',   -- escrow / deposit address
    confirmations INTEGER NOT NULL DEFAULT 1,
    enabled       INTEGER NOT NULL DEFAULT 1,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);
-- Seeded once: ('SKY','Skycoin', …, enabled=1). SKY's empty escrow fields fall
-- back to wallet_sky / sky_wallet_seed / sky_fullnode_url.
```

A listing/product/order records which coin it trades in a `sell_coin` column
(default `'SKY'`). The chain backend routes each on-chain action (deposit check,
delivery, refund, escrow-audit) to that coin's fullnode + escrow wallet. The escrow
audit runs **per sell coin**. Correspondingly the wire protocol carries a
`sell_coin` field and uses coin-generic amount fields (`amount`, `expected_amount`,
`commission`) rather than the old SKY-specific names.

#### Fibercoin-to-fibercoin trades (payment in a fibercoin)

The **payment currency** may itself be an enabled sell coin — so with two or more
sell coins configured you can trade **SKY/FIBER1, FIBER1/FIBER2**, etc. (a coin may
not be traded for itself). A fibercoin payment is a **peer-to-peer transfer** from
the buyer's Skycoin-family address to the seller's, verified on that coin's own
fullnode by **sender + time-window + exact amount** — the same mechanism as a
sell-coin deposit. Because the sender address (the buyer's registered address,
hard-matched) already disambiguates payments, the buyer sends a **FIXED** amount
(exactly `100 FIBER1`, *not* a unique `100.002`). This differs from external UTXO
payment coins (BTC/LTC/…), which are matched on a shared explorer-visible address
and therefore still need the unique non-round amount. Delivery of the sold coin
from escrow is unchanged; the market only *verifies* the payment leg. A single
buyer may not hold two identical fixed pending payments to the same seller (the
sender can no longer tell them apart); the second is rejected.

### 7.8. Background jobs (market)
The market runs several periodic background jobs:

#### 1. Escrow Checker (every 30 s)
- Check buyer payment transactions on explorers and the SKY blockchain.
- Inspect `orders` with `status = pending_payment` whose `expires_at` has not passed.
- Increment `confirmations` as new confirmations appear.
- Set status to `confirmed` once 2 confirmations are reached.
- After `confirmed`, transfer SKY to the buyer's wallet and set status to `completed`.

#### 2. Listing Checker (every 30 s)
- Check seller transfers into the market wallet.
- Inspect `pending_listings` with `status = pending`.
- On confirmed transfer, set status to `confirmed` and create a row in `products`.

#### 3. Expiry Handler (every 10 s)
- Handle expired orders:
  - `pending_listings` with `expires_at < now` → set status to `expired`.
  - `orders` with `status = pending_payment` and `expires_at < now` → set status to
    `expired`, release the product, and record a `freeze_violation`.
- Return frozen `products` whose order expired back to `active`.

> **Status:** All jobs — Escrow Checker, Listing Checker, Expiry Handler, Return Scheduler,
> Ban Manager, Cleanup Job and **Escrow Audit** (§7.8/7) — are implemented in
> `src/skydex/jobs` and run under a `Runner`. The blockchain-dependent checks
> (deposit/payment verification, SKY delivery and refunds) go through a `Chain` interface backed
> by a Skycoin node + per-coin explorers; SKY delivery/refund transactions are built and **signed
> locally** from `sky_wallet_seed` and broadcast via the node (the node's own wallet API is not
> used). Without a configured backend the market uses `NoopChain`, so the DB-side lifecycle runs
> but no trade settles on-chain.

#### 4. Return Scheduler (every 1 min)
- Handle unsold/cancelled listings — refunds are **immediate** (no delay).
- For each expired/cancelled listing: re-verify the seller's deposit is present in the market
  wallet on-chain (`DepositConfirmed`, which guards against refunding SKY that never arrived),
  then transfer the escrowed SKY back to the seller via `SendSKY`, record `return_tx_hash`, and
  set the listing status to `returned` (recorded once, so retries can't double-refund). If the
  escrow can't yet fund the refund transaction (its coin-hours are still accruing), `SendSKY`
  fails and the refund is retried on the next tick — so it settles as soon as it's fundable,
  without a fixed wait.
- The full deposit (amount + commission) is refunded; a cancelled/expired sale is not charged
  commission.
- Escrow-return bookkeeping lives on `pending_listings`: `closed_at` (when it went terminal),
  `returned_at`, and `return_tx_hash`. `products.listing_id` links a live product back to its
  listing so a seller cancel can find and deactivate the right product.

#### 5. Cleanup Job (every 1 h)
- Delete old data for privacy:
  - `orders` with `status = completed` and `completed_at < now - 3 days`.
  - `products` linked to deleted orders.
  - `freeze_violations` older than 2 weeks.
  - `pending_listings` with `status = expired` or `confirmed` older than 3 days.

#### 6. Ban Manager (every 1 min)
- Delete `bans` with `ban_until < now` (auto-release banned users).
- Count each user's `freeze_violations` in the last 24 hours.
- If violations >= 3, create a row in `bans`.

#### 7. Escrow Audit (every 5 min)
- Read the escrow wallet's live on-chain SKY balance (`EscrowBalance`) and compare it against
  the DB's outstanding obligations (`OutstandingEscrowSKY` = SKY for products awaiting delivery
  (`active`/`frozen`) + confirmed deposits awaiting refund). Log an **error on any shortfall**
  (on-chain below obligations → a delivery/refund could fail) or the healthy state with headroom
  (the headroom is the retained commission). Read-only; never moves funds; a no-op without a
  chain backend.

---

## 8. API Design
Client↔market communication uses a JSON-based message protocol over **dmsg**. All messages
are request/response, and polling is used for status updates.

### 8.1. Message structure
Every message has these fields:

```json
{
  "type": "<message_type>",
  "id": "<uuid>",
  "timestamp": 1234567890,
  "data": { ... }
}
```

- **type:** message type (e.g. `client.buy_product`)
- **id:** unique id to correlate request and response
- **timestamp:** send time (Unix timestamp)
- **data:** message payload

### 8.2. Client → Market messages

#### 1. Register/update user info
Phase 1 trades **SKY** (sold) for **BTC** or **LTC** (payment), so a user registers a SKY
address (to receive purchased SKY, or refunds) and — if selling — the BTC/LTC address to be
paid to. The market **validates each address's format** (base58check version+checksum, or
bech32/bech32m segwit; SKY base58check) and rejects a malformed one with `INVALID_WALLET`, so a
typo fails at registration rather than as a silent failed payout later. Other coins
(`wallet_bch`, `wallet_usdt_*`) are accepted but unused/unvalidated until they become supported.
```json
{
  "type": "client.register",
  "id": "uuid",
  "timestamp": 1234567890,
  "data": {
    "wallet_sky": "2i7rBKfQgBtBcbnKUr6FXEGFANXBNV1WYJw",
    "wallet_btc": "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq",
    "wallet_ltc": "ltc1q…"
  }
}
```

**Response:**
```json
{
  "type": "response",
  "id": "uuid",
  "timestamp": 1234567890,
  "status": "success",
  "data": { "message": "User registered successfully" }
}
```

#### 2. Get available payment currencies
```json
{ "type": "client.get_currencies", "id": "uuid", "timestamp": 1234567890, "data": {} }
```

**Response:** `data: { "currencies": ["BTC", "USDT_TRC20"] }` — only the coins whose explorer
the market has configured.

#### 3. Get active products
```json
{ "type": "client.get_products", "id": "uuid", "timestamp": 1234567890, "data": {} }
```

**Response:**
```json
{
  "type": "response",
  "id": "uuid",
  "timestamp": 1234567890,
  "status": "success",
  "data": {
    "products": [
      {
        "id": "product-uuid",
        "seller_pubkey": "02abc...",
        "amount_sky": 10.5,
        "price": 2.50,
        "payment_currency": "BTC",
        "created_at": "2026-07-10T12:00:00Z"
      }
    ]
  }
}
```

#### 4. Create a listing
```json
{
  "type": "client.create_listing",
  "id": "uuid",
  "timestamp": 1234567890,
  "data": { "amount_sky": 10.0, "price": 2.00, "payment_currency": "BTC" }
}
```

**Response:**
```json
{
  "type": "response",
  "id": "uuid",
  "timestamp": 1234567890,
  "status": "success",
  "data": {
    "listing_id": "listing-uuid",
    "amount_sky": 10,
    "commission_sky": 0.05,
    "expected_amount_sky": 10.05,
    "market_wallet": "sky1market...",
    "expires_at": "2026-07-10T12:15:00Z"
  }
}
```

#### 5. Buy a product (freeze)
```json
{
  "type": "client.buy_product",
  "id": "uuid",
  "timestamp": 1234567890,
  "data": { "product_id": "product-uuid" }
}
```

**Response (success):**
```json
{
  "type": "response",
  "id": "uuid",
  "timestamp": 1234567890,
  "status": "success",
  "data": {
    "order_id": "order-uuid",
    "amount_sky": 10.5,
    "price": 2.50,
    "payment_currency": "BTC",
    "expected_payment_amount": 2.00034123,
    "seller_wallet": "bc1seller...",
    "expires_at": "2026-07-10T12:15:00Z"
  }
}
```

**Response (error — product banned or frozen):**
```json
{
  "type": "response",
  "id": "uuid",
  "timestamp": 1234567890,
  "status": "error",
  "data": { "code": "PRODUCT_UNAVAILABLE", "message": "This product is no longer available" }
}
```

#### 6. Cancel a listing
```json
{
  "type": "client.cancel_listing",
  "id": "uuid",
  "timestamp": 1234567890,
  "data": { "listing_id": "listing-uuid" }
}
```

**Response:**
```json
{
  "type": "response",
  "id": "uuid",
  "timestamp": 1234567890,
  "status": "success",
  "data": { "message": "Offer canceled. Your SKY will be returned shortly." }
}
```

#### 7. Get the user's orders
```json
{ "type": "client.get_orders", "id": "uuid", "timestamp": 1234567890, "data": {} }
```

**Response:**
```json
{
  "type": "response",
  "id": "uuid",
  "timestamp": 1234567890,
  "status": "success",
  "data": {
    "orders": [
      {
        "id": "order-uuid",
        "type": "buy",
        "product_id": "product-uuid",
        "amount_sky": 10.5,
        "price": 2.50,
        "payment_currency": "BTC",
        "status": "pending_payment",
        "expires_at": "2026-07-10T12:15:00Z",
        "created_at": "2026-07-10T12:00:00Z"
      }
    ]
  }
}
```

#### 8. Get a specific order's status
```json
{
  "type": "client.get_order_status",
  "id": "uuid",
  "timestamp": 1234567890,
  "data": { "order_id": "order-uuid" }
}
```

**Response:**
```json
{
  "type": "response",
  "id": "uuid",
  "timestamp": 1234567890,
  "status": "success",
  "data": {
    "order_id": "order-uuid",
    "status": "confirmed",
    "confirmations": 2,
    "payment_tx_hash": "txhash...",
    "paid_at": "2026-07-10T12:05:00Z"
  }
}
```

### 8.3. Error codes

| Code | Meaning |
|------|---------|
| `INVALID_REQUEST` | Invalid or incomplete request |
| `USER_BANNED` | User is banned and cannot trade |
| `PRODUCT_NOT_FOUND` | Product not found |
| `PRODUCT_UNAVAILABLE` | Product unavailable (frozen or sold) |
| `LISTING_NOT_FOUND` | Listing not found |
| `ORDER_NOT_FOUND` | Order not found |
| `CURRENCY_UNAVAILABLE` | Payment currency not enabled at this market (no explorer configured) |
| `SESSION_CONFLICT` | User is connected from elsewhere |
| `INTERNAL_ERROR` | Internal market error |

### 8.4. dmsg connection protocol
1. **Connect:** The client opens a dmsg connection using the market's public key.
2. **Session management:** The market creates a session per active connection and records
   the user's public key.
3. **Single-session enforcement:** If a user with the same public key already has an active
   connection, the new connection is rejected (`SESSION_CONFLICT`).
4. **Timeout:** If a session sends no request for 5 minutes it is closed (a per-connection
   read deadline reset on every request), freeing the single-session slot and serving
   goroutine an abandoned connection would hold. The client polls well within this window, so
   an active user never trips it. A single response write is likewise bounded (30 s) so a peer
   that stops draining can't pin the serving goroutine. **Implemented** in `server.Serve`.
5. **Reconnection:** The user can reconnect with the same public key once the previous
   session has ended, creating a new session.

### 8.5. Polling for status updates
The client periodically (e.g. every 5 s) sends:
- `client.get_orders` to refresh the order list
- `client.get_order_status` for each active order
- `client.get_products` to refresh the product list

This is simple and reliable, and removes the need for WebSocket or SSE.

---

## 9. Repository Layout

The code is split across two repos (see "Where the code lives" in the Overview).

**skycoin** — `github.com/skycoin/skycoin` (the engine + UI; no skywire import):

```
skycoin/
├── src/skydex/                  # Engine (transport-agnostic; imported by skywire)
│   ├── protocol/                #   SHARED wire protocol: envelope, message types, codes
│   ├── message/                 #   length-prefixed frame codec (copy of skywire's pkg/skychat/message)
│   ├── db/                      #   SQLite layer + models + CRUD (8 tables + config, sell_coins)
│   ├── chain/                   #   Skycoin node client + per-coin explorer adapters + local signing
│   ├── jobs/                    #   background jobs under a Runner
│   ├── walletaddr/              #   per-coin wallet-address validation
│   ├── server/                  #   server core: Accept(net.Listener, IdentifyFunc) + ServeConn + handlers
│   └── market/                  #   client core: NewConn(net.Conn) + framed request/response (Do)
├── cmd/skydex-market/           # Standalone (TCP) market command + operator/auth UI
│   └── commands/                #   Run(ctx, cfg, net.Listener, IdentifyFunc, Host) · api/auth/ui · root.go (TCP) · static/
└── cmd/skydex-client/           # Standalone (TCP) client command + React SPA
    ├── commands/                #   Run(ctx, cfg, MarketDialer, log) · session/api · root.go (TCP) · static/ (go:embed)
    └── src/ · index.html · vite.config.js · package.json
```

**skywire** — `github.com/skycoin/skywire` (thin transport wrapper only):

```
skywire/
├── cmd/apps/skydex-market/      # Market visor app
│   ├── main.go                  #   cobra entrypoint
│   └── commands/root.go         #   RegisterApp("skydex-market", …); appnet listener + PK-identify → skydexmarket.Run
├── cmd/apps/skydex-client/      # Client visor app
│   ├── main.go                  #   cobra entrypoint
│   └── commands/skydex-client.go#   RegisterApp("skydex-client", …); appnetDialer (MarketDialer) → skydexclient.Run
├── internal/skydex-market/app/  # visor app.Client bridge (Host: Log/PublishOTP/PubKey)
├── internal/skydex-client/app/  # visor app.Client bridge
└── internal/skydex-e2e/         # two-visor dmsg wire test (drives the skycoin engine over a real mesh)
```

Transport model: the skywire **market wrapper** `Listen`s on `appnet.TypeDmsg` (via the
visor's `app.Client`), extracts the client's authenticated public key from the accepted
connection's `appnet.Addr`, and calls the engine's `Accept(...)` with that as the identity.
The skywire **client wrapper** `Dial`s the market's public key on `skyenv.SkydexMarketPort`
and hands the raw `net.Conn` to the engine's `market.NewConn`. Skywire supplies the
encrypted dmsg `net.Conn`; all framing (length-prefixed JSON `Envelope`s, `src/skydex/
message`) and protocol logic live in the engine. Standalone, the same engine runs over a
plain TCP listener/dialer. No bespoke dmsg server/client is implemented.

The `skydex-market` / `skydex-client` apps are registered via `launcher.RegisterApp` and
compiled into the merged `skywire` binary, so the visor launches them in-process exactly
like the other bundled apps (skychat, skynet, vpn-client, …).

> **Note:** `src/skydex/protocol` is the shared package imported by both engine sides and by
> the skywire wrappers; it lives with the engine because the market is the protocol
> authority. The market **operator UI** (config + monitoring) is implemented in
> `cmd/skydex-market/commands`, served locally by the market and embedded via `go:embed`.

---

## Testing & validation

Three tiers, from fully offline to fully live:

1. **Unit / integration (offline, default `go test`).** The engine's DB CRUD, protocol
   dispatch, business logic, background jobs, the operator/client HTTP APIs, and an Esplora
   adapter test against a mock server run in **skycoin** (`go test ./src/skydex/...
   ./cmd/skydex-...`). The in-process **two-visor dmsg E2E** runs in **skywire**
   (`internal/skydex-e2e`): the market on one dmsg identity and a seller + buyer on two
   others, over a real `dmsgtest` mesh (Noise + yamux), driving the full trade —
   list → deposit → promote → buy → pay → confirm → deliver → completed. The blockchain
   is a scripted stand-in so the whole lifecycle runs deterministically without a node.
   Skipped under `-short`.

2. **Live explorer validation (network, opt-in).** Validates the Esplora adapter against
   the real public APIs (mempool.space for BTC, litecoinspace.org for LTC): it derives a
   currently-confirmed on-chain payment and asserts the adapter parses the live response
   shape and matches the exact amount with a positive confirmation count.

   ```
   # from the skycoin repo:
   go test -tags livenet -run TestLiveEsplora -v ./src/skydex/chain/
   ```

3. **Fully-live trade (real infra — remaining step).** The one part that needs an
   operator environment: point the market at a **real Skycoin node** (deposit detection,
   `SendSKY` delivery/refunds — and later Coin Hours) and enable per-coin explorers, then
   run a real two-visor trade paying real BTC/LTC. The dmsg path and protocol are already
   proven by tier 1 and the explorer shapes by tier 2; tier 3 swaps the scripted chain for
   `chain.New(...)` against the node. This is the only piece that can't be self-served in CI.
