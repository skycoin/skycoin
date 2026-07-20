import { useEffect, useState, useCallback } from 'react'
import {
  api,
  loadTradeHistory,
  mergeTradeHistory,
  removeTradeHistoryItem,
  clearTradeHistory,
} from '../api'

// Terminal states: a trade that is done (settled, or ended without settling).
const TERMINAL = new Set(['completed', 'canceled', 'cancelled', 'expired'])

// History lists the user's finished trades — both buys and sells — that are no
// longer in flight. It is stored LOCALLY (localStorage) and accumulated over
// time, because the market deletes finished orders after its cleanup window
// (cleanup_days, default 3). So a trade stays in this list on this device even
// once the market has purged it. Active orders live under My Orders / My Listings.
function History({ marketPubKey, marketName }) {
  const [rows, setRows] = useState(() => sortByDate(loadTradeHistory()))

  const load = useCallback(async () => {
    try {
      const data = await api.getOrders()
      const terminal = (data.orders || []).filter((o) => TERMINAL.has(o.status))
      const merged = mergeTradeHistory(terminal, marketPubKey, marketName)
      setRows(sortByDate(merged))
    } catch {
      // Not connected / market unreachable — still show the locally saved history.
      setRows(sortByDate(loadTradeHistory()))
    }
  }, [marketPubKey, marketName])

  useEffect(() => {
    load()
    const id = setInterval(load, 8000)
    return () => clearInterval(id)
  }, [load])

  const handleClear = () => {
    if (window.confirm('Clear all locally saved trade history on this device?')) {
      setRows(sortByDate(clearTradeHistory()))
    }
  }

  const badge = (status) => {
    const cls =
      status === 'completed' ? 'enabled' : status === 'expired' ? 'disabled' : 'cancelled'
    return <span className={`badge ${cls}`}>{status}</span>
  }
  const typeBadge = (type) => (
    <span className={`badge ${type === 'buy' ? 'buy' : 'sell'}`}>{type}</span>
  )

  return (
    <div>
      <div className="page-head">
        <h2>History</h2>
        {rows.length > 0 && (
          <button type="button" className="link-btn" onClick={handleClear}>
            Clear history
          </button>
        )}
      </div>

      {rows.length === 0 ? (
        <div className="alert alert-secondary">No completed trades yet.</div>
      ) : (
        <>
          <div className="panel table-wrap">
            <table className="table table-hover align-middle">
              <thead>
                <tr>
                  <th>Type</th>
                  <th>Coin</th>
                  <th>Amount</th>
                  <th>Price</th>
                  <th>Market</th>
                  <th>Status</th>
                  <th>Date</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {rows.map((o) => (
                  <tr key={o.id}>
                    <td>{typeBadge(o.type)}</td>
                    <td>{o.sell_coin}</td>
                    <td>{o.amount} {o.sell_coin}</td>
                    <td>{o.price} {o.payment_currency}</td>
                    <td>
                      {o.market_name ? (
                        <span title={o.market_pk}>{o.market_name}</span>
                      ) : o.market_pk ? (
                        <code title={o.market_pk}>{shorten(o.market_pk)}</code>
                      ) : (
                        <span className="text-muted">—</span>
                      )}
                    </td>
                    <td>{badge(o.status)}</td>
                    <td>{o.created_at ? new Date(o.created_at).toLocaleString() : '-'}</td>
                    <td>
                      <button
                        type="button"
                        className="recent-del"
                        aria-label="Remove from history"
                        title="Remove from history"
                        onClick={() => setRows(sortByDate(removeTradeHistoryItem(o.id)))}
                      >
                        ✕
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="hint mt-2">
            History is saved locally on this device — the market keeps finished orders only
            temporarily, so this list persists even after they’re removed there.
          </div>
        </>
      )}
    </div>
  )
}

function sortByDate(list) {
  return [...list].sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
}

function shorten(pk) {
  if (!pk || pk.length <= 14) return pk
  return `${pk.slice(0, 8)}…${pk.slice(-4)}`
}

export default History
