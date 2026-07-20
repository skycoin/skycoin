import { useEffect, useState, useCallback } from 'react'
import { api } from '../api'
import CopyButton from './CopyButton'
import QrButton from './QrButton'

// MyListings shows the seller their own pending sell listings and tracks each
// one's lifecycle live: pending (awaiting the SKY deposit) -> confirmed (deposit
// detected on-chain; the listing is now a purchasable product on the market).
function MyListings() {
  const [listings, setListings] = useState([])
  const [error, setError] = useState('')
  const [notice, setNotice] = useState('')
  const [busyId, setBusyId] = useState('')

  const load = useCallback(async () => {
    try {
      const data = await api.getListings()
      setListings(data.listings || [])
      setError('')
    } catch (e) {
      setError(e.message)
    }
  }, [])

  useEffect(() => {
    load()
    const id = setInterval(load, 5000)
    return () => clearInterval(id)
  }, [load])

  // The seller can cancel their offer while it is pending (no deposit yet) or
  // confirmed (an active product), as long as no buyer has selected it. A
  // confirmed offer's escrowed SKY is returned shortly after cancellation.
  const cancelListing = async (l) => {
    const msg =
      l.status === 'pending'
        ? 'Cancel this listing?'
        : `Cancel this offer? Your escrowed ${l.sell_coin || 'coin'} will be returned shortly after cancellation.`
    if (!window.confirm(msg)) return
    setBusyId(l.id)
    setError('')
    try {
      const resp = await api.cancelListing(l.id)
      setNotice(resp.message || 'Offer canceled.')
      setTimeout(() => setNotice(''), 8000)
      await load()
    } catch (e) {
      setError(e.message)
    } finally {
      setBusyId('')
    }
  }

  // The three lifecycle steps, with the current one highlighted. A confirmed
  // listing is already live on the market, so its reached step is "Listed as
  // product" (index 2), not "Confirmed".
  const steps = ['pending', 'confirmed', 'product']
  const stepIndex = (status) => (status === 'pending' ? 0 : 2)

  const renderLifecycle = (status) => {
    const active = stepIndex(status)
    const labels = ['Pending deposit', 'Confirmed', 'Listed as product']
    return (
      <div className="d-flex align-items-center gap-2 flex-wrap">
        {steps.map((_, i) => (
          <span key={i} className="d-flex align-items-center gap-2">
            <span
              className={`badge ${
                i < active ? 'bg-success' : i === active ? 'bg-info' : 'bg-secondary'
              }`}
            >
              {labels[i]}
            </span>
            {i < steps.length - 1 && <span className="text-muted">→</span>}
          </span>
        ))}
      </div>
    )
  }

  // Once a buyer selects a confirmed listing's product, the seller is committed;
  // surface that state (and it stops being cancelable).
  const productStateBadge = (l) => {
    if (l.status !== 'confirmed' || !l.product_status || l.product_status === 'active') return null
    const map = {
      frozen: { label: 'Buyer selected — awaiting payment', cls: 'pending' },
      sold: { label: 'Sold', cls: 'completed' },
      cancelled: { label: 'Cancelled', cls: 'cancelled' },
    }
    const info = map[l.product_status] || { label: l.product_status, cls: 'disabled' }
    return <span className={`badge ${info.cls}`}>{info.label}</span>
  }

  // What to show in the Actions column when a listing can't be canceled.
  const actionNote = (l) => {
    if (l.status === 'confirmed' && l.product_status === 'frozen') return 'Buyer paying…'
    if (l.status === 'confirmed' && l.product_status === 'sold') return 'Sold'
    return '—'
  }

  return (
    <div>
      <div className="page-head"><h2>My Listings</h2></div>

      {notice && <div className="alert alert-info">{notice}</div>}
      {error && <div className="alert alert-danger">{error}</div>}

      {listings.length === 0 ? (
        <div className="alert alert-secondary">
          You have no active sell listings. Create one from the Market tab.
        </div>
      ) : (
        <div className="panel table-wrap">
          <table className="table table-hover align-middle">
            <thead>
              <tr>
                <th>ID</th>
                <th>Coin</th>
                <th>Amount</th>
                <th>Price</th>
                <th>Deposit</th>
                <th>Send to (escrow address)</th>
                <th>Lifecycle</th>
                <th>Deposit Tx</th>
                <th>Expires</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {listings.map((l) => (
                <tr key={l.id}>
                  <td><code>{shorten(l.id)}</code></td>
                  <td>{l.sell_coin}</td>
                  <td>{l.amount} {l.sell_coin}</td>
                  <td>{l.price} {l.payment_currency}</td>
                  <td>
                    {l.status === 'pending' ? (
                      <span>
                        <strong>{l.expected_amount}</strong> {l.sell_coin}
                        {l.expected_amount > l.amount && (
                          <div className="small text-muted">
                            {l.amount} + {(l.expected_amount - l.amount).toFixed(3)} fee
                          </div>
                        )}
                      </span>
                    ) : (
                      <span className="text-muted">received</span>
                    )}
                  </td>
                  <td>
                    {l.status === 'pending' && l.market_wallet ? (
                      <div className="copy-row">
                        <code className="addr-box addr-sm" title={l.market_wallet}>{l.market_wallet}</code>
                        <CopyButton text={l.market_wallet} label="Copy" />
                        <QrButton
                          address={l.market_wallet}
                          amount={l.expected_amount}
                          coin={l.sell_coin}
                          title="Fund your sell order"
                          hint="This is the market escrow address, not the buyer's."
                        />
                      </div>
                    ) : (
                      <span className="text-muted">—</span>
                    )}
                  </td>
                  <td>
                    <div className="d-flex align-items-center gap-2 flex-wrap">
                      {renderLifecycle(l.status)}
                      {productStateBadge(l)}
                    </div>
                  </td>
                  <td>{l.tx_hash ? <code title={l.tx_hash}>{shorten(l.tx_hash)}</code> : <span className="text-muted">—</span>}</td>
                  <td>
                    {l.status === 'pending' && l.expires_at
                      ? new Date(l.expires_at).toLocaleTimeString()
                      : '-'}
                  </td>
                  <td>
                    {l.cancelable ? (
                      <button
                        className="btn btn-sm btn-outline-danger"
                        disabled={busyId === l.id}
                        onClick={() => cancelListing(l)}
                      >
                        {busyId === l.id ? 'Canceling…' : 'Cancel'}
                      </button>
                    ) : (
                      <span className="text-muted">{actionNote(l)}</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {listings.some((l) => l.status === 'pending') && (
        <div className="alert alert-info mt-3">
          For each pending listing, transfer the exact <strong>Deposit</strong> amount to the
          <strong> Send to</strong> escrow address before it expires. The offer becomes a
          purchasable product once the deposit is confirmed on-chain.
        </div>
      )}
    </div>
  )
}

function shorten(id) {
  if (!id || id.length <= 12) return id
  return `${id.slice(0, 8)}…`
}

export default MyListings
