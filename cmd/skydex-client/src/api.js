// Thin wrapper around the skydex-client local control API. Every call returns
// the parsed JSON body; on a non-2xx response it throws an Error carrying the
// server's message so callers can surface it.

async function request(method, url, body) {
  const opts = { method, headers: {} }
  if (body !== undefined) {
    opts.headers['Content-Type'] = 'application/json'
    opts.body = JSON.stringify(body)
  }
  const res = await fetch(url, opts)
  let data = {}
  try {
    data = await res.json()
  } catch {
    // empty/invalid body
  }
  if (!res.ok) {
    throw new Error(data.error || `request failed (${res.status})`)
  }
  return data
}

// Wallet addresses are stored locally as a convenience so the Settings form is
// pre-filled and the app can prompt first-time users to register. The market
// remains the source of truth; this is only a UX hint.
const WALLETS_KEY = 'exchange:wallets'

export function loadWallets() {
  try {
    return JSON.parse(localStorage.getItem(WALLETS_KEY)) || {}
  } catch {
    return {}
  }
}

export function saveWallets(w) {
  try {
    localStorage.setItem(WALLETS_KEY, JSON.stringify(w))
  } catch {
    // ignore storage errors (private mode, etc.)
  }
}

export function hasSavedWallets() {
  const w = loadWallets()
  return !!(w && w.SKY)
}

// Recent markets: a local, most-recent-first list of market public keys this
// client has connected to, so the user can rejoin from a list instead of
// re-typing the key. Purely a local convenience; the market is never contacted.
const MARKETS_KEY = 'exchange:markets'
const MARKETS_MAX = 12

export function loadMarketHistory() {
  try {
    const arr = JSON.parse(localStorage.getItem(MARKETS_KEY))
    return Array.isArray(arr) ? arr.filter((m) => m && m.pk) : []
  } catch {
    return []
  }
}

export function addMarketToHistory(pk, name) {
  const key = (pk || '').trim()
  if (!key) return loadMarketHistory()
  const existing = loadMarketHistory()
  const prev = existing.find((m) => m.pk === key)
  const rest = existing.filter((m) => m.pk !== key)
  // Keep a previously-learned name if this call doesn't carry one.
  const nm = (name || '').trim() || (prev && prev.name) || ''
  const next = [{ pk: key, name: nm, ts: Date.now() }, ...rest].slice(0, MARKETS_MAX)
  try {
    localStorage.setItem(MARKETS_KEY, JSON.stringify(next))
  } catch {
    // ignore storage errors (private mode, quota)
  }
  return next
}

export function removeMarketFromHistory(pk) {
  const next = loadMarketHistory().filter((m) => m.pk !== pk)
  try {
    localStorage.setItem(MARKETS_KEY, JSON.stringify(next))
  } catch {
    // ignore storage errors
  }
  return next
}

export function clearMarketHistory() {
  try {
    localStorage.removeItem(MARKETS_KEY)
  } catch {
    // ignore storage errors
  }
  return []
}

// Trade history is kept locally because the market deletes finished orders after
// its cleanup window (cleanup_days, default 3). We accumulate every terminal
// order we ever see into localStorage, keyed by order id, so the user's history
// persists on this device even after the market has purged it.
const HISTORY_KEY = 'exchange:history'

export function loadTradeHistory() {
  try {
    const arr = JSON.parse(localStorage.getItem(HISTORY_KEY))
    return Array.isArray(arr) ? arr.filter((o) => o && o.id) : []
  } catch {
    return []
  }
}

// mergeTradeHistory folds the given terminal orders into the local store (adding
// new ones, refreshing known ones), tagging each with the market it came from,
// and returns the merged list. Orders already recorded but no longer returned by
// the market (purged) are kept.
export function mergeTradeHistory(orders, marketPK, marketName) {
  const byId = new Map(loadTradeHistory().map((o) => [o.id, o]))
  for (const o of orders || []) {
    if (!o || !o.id) continue
    const prev = byId.get(o.id) || {}
    byId.set(o.id, {
      ...prev,
      ...o,
      market_pk: marketPK || prev.market_pk || '',
      market_name: marketName || prev.market_name || '',
      saved_at: prev.saved_at || Date.now(),
    })
  }
  const next = [...byId.values()]
  try {
    localStorage.setItem(HISTORY_KEY, JSON.stringify(next))
  } catch {
    // ignore storage errors
  }
  return next
}

export function removeTradeHistoryItem(id) {
  const next = loadTradeHistory().filter((o) => o.id !== id)
  try {
    localStorage.setItem(HISTORY_KEY, JSON.stringify(next))
  } catch {
    // ignore storage errors
  }
  return next
}

export function clearTradeHistory() {
  try {
    localStorage.removeItem(HISTORY_KEY)
  } catch {
    // ignore storage errors
  }
  return []
}

// Copy text to the clipboard, with a fallback for non-secure (plain http)
// contexts where navigator.clipboard is unavailable — the app is served over
// http on the visor, so the modern API is often missing. Returns true on
// success. Callers can surface their own "copied" feedback.
export async function copyToClipboard(text) {
  if (!text) return false
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text)
      return true
    }
  } catch {
    // fall through to the legacy execCommand path
  }
  try {
    const ta = document.createElement('textarea')
    ta.value = text
    ta.style.position = 'fixed'
    ta.style.top = '-1000px'
    ta.style.opacity = '0'
    document.body.appendChild(ta)
    ta.focus()
    ta.select()
    const ok = document.execCommand('copy')
    document.body.removeChild(ta)
    return ok
  } catch {
    return false
  }
}

export const api = {
  getConfig: () => request('GET', '/api/config'),
  getStatus: () => request('GET', '/api/status'),
  connect: (marketPK) => request('POST', '/api/connect', { market_pk: marketPK }),
  disconnect: () => request('POST', '/api/disconnect'),

  register: (wallets) => request('POST', '/api/register', wallets),
  getCurrencies: () => request('GET', '/api/currencies'),
  getProducts: () => request('GET', '/api/products'),
  getOrders: () => request('GET', '/api/orders'),
  getListings: () => request('GET', '/api/listings/mine'),
  createListing: (listing) => request('POST', '/api/listings', listing),
  cancelListing: (listingId) => request('POST', '/api/listings/cancel', { listing_id: listingId }),
  cancelOrder: (orderId) => request('POST', '/api/orders/cancel', { order_id: orderId }),
  buyProduct: (productId) => request('POST', '/api/buy', { product_id: productId }),
  getOrderStatus: (orderId) => request('POST', '/api/order-status', { order_id: orderId }),
}
