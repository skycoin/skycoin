import { useState } from 'react'
import { loadMarketHistory, removeMarketFromHistory, clearMarketHistory } from '../api'

// ConnectGate is the first screen the user sees. It shows the market public key
// (pre-filled from the --market-pk config value), lets the user confirm or edit
// it, and connects only when the user clicks Connect. The client never connects
// automatically. It also lists recently joined markets so the user can rejoin
// one with a click, remove entries, or clear the whole list.
function ConnectGate({ defaultPubKey, onConnected }) {
  const [pubKey, setPubKey] = useState(defaultPubKey || '')
  const [connecting, setConnecting] = useState(false)
  const [error, setError] = useState('')
  const [history, setHistory] = useState(() => loadMarketHistory())

  async function connectTo(rawPk) {
    setError('')
    const pk = (rawPk || '').trim()
    if (!pk) {
      setError('Please enter the market public key.')
      return
    }
    setConnecting(true)
    try {
      const res = await fetch('/api/connect', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ market_pk: pk }),
      })
      const data = await res.json()
      if (!res.ok) {
        setError(data.error || 'Failed to connect to the market.')
        return
      }
      onConnected(data.market_pk || pk, data.market_name || '')
    } catch (err) {
      setError('Failed to reach the market: ' + err.message)
    } finally {
      setConnecting(false)
    }
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    connectTo(pubKey)
  }

  const handleClear = () => {
    if (window.confirm('Clear all recent markets from this device?')) {
      setHistory(clearMarketHistory())
    }
  }

  return (
    <div className="app-container d-flex align-items-center justify-content-center">
      <div className="connect-card">
        <h3 className="mb-1 text-center">Skywire Exchange</h3>
        <p className="text-center text-muted mb-4">Connect to a market to start trading</p>

        <form onSubmit={handleSubmit}>
          <label className="form-label">Market Public Key</label>
          <input
            type="text"
            className="form-control mb-2"
            placeholder="Enter the skydex-market public key"
            value={pubKey}
            onChange={(e) => setPubKey(e.target.value)}
            disabled={connecting}
            autoFocus
          />

          {error && <div className="alert alert-danger py-2">{error}</div>}

          <button type="submit" className="btn btn-connect w-100" disabled={connecting}>
            {connecting ? 'Connecting…' : 'Connect'}
          </button>
        </form>

        {history.length > 0 && (
          <div className="recent-markets">
            <div className="recent-head">
              <span className="recent-title">Recent markets</span>
              <button type="button" className="link-btn" onClick={handleClear} disabled={connecting}>
                Clear history
              </button>
            </div>
            <ul className="recent-list">
              {history.map((m) => (
                <li key={m.pk} className="recent-item">
                  <button
                    type="button"
                    className="recent-connect"
                    onClick={() => connectTo(m.pk)}
                    disabled={connecting}
                    title={`Connect to ${m.pk}`}
                  >
                    <span className="recent-main">
                      {m.name && <span className="recent-name">{m.name}</span>}
                      <code className="recent-pk">{shortenPk(m.pk)}</code>
                    </span>
                    <span className="recent-time">{timeAgo(m.ts)}</span>
                  </button>
                  <button
                    type="button"
                    className="recent-del"
                    onClick={() => setHistory(removeMarketFromHistory(m.pk))}
                    disabled={connecting}
                    aria-label="Remove from history"
                    title="Remove from history"
                  >
                    ✕
                  </button>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  )
}

function shortenPk(pk) {
  if (!pk || pk.length <= 18) return pk
  return `${pk.slice(0, 10)}…${pk.slice(-6)}`
}

function timeAgo(ts) {
  if (!ts) return ''
  const s = Math.max(0, Math.floor((Date.now() - ts) / 1000))
  if (s < 60) return 'just now'
  const m = Math.floor(s / 60)
  if (m < 60) return `${m}m ago`
  const h = Math.floor(m / 60)
  if (h < 24) return `${h}h ago`
  const d = Math.floor(h / 24)
  return `${d}d ago`
}

export default ConnectGate
