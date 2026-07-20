import { useEffect, useState, useCallback } from 'react'
import { api } from '../api'
import CopyButton from './CopyButton'
import QrButton from './QrButton'

// Fallback confirmation threshold for the "n/N" progress display. The market
// reports its own value in each order-status (required_confirmations); this is
// only used if that field is absent (older market).
const REQUIRED_CONFIRMATIONS = 2

// Buy orders in these states are still in flight, so their on-chain progress
// (confirmations, payment tx) is worth polling live via get_order_status.
const IN_FLIGHT = new Set(['pending_payment', 'paid', 'confirmed'])

function Orders() {
  const [orders, setOrders] = useState([])
  // statuses maps order id -> live {status, confirmations, payment_tx_hash, paid_at}.
  const [statuses, setStatuses] = useState({})
  const [error, setError] = useState('')
  const [notice, setNotice] = useState('')
  const [busyId, setBusyId] = useState('')

  const load = useCallback(async () => {
    try {
      const data = await api.getOrders()
      const list = data.orders || []
      setOrders(list)
      setError('')

      // Poll live status only for the buyer's own in-flight buy orders; the
      // market only lets the buyer query an order's status.
      const tracked = list.filter((o) => o.type === 'buy' && IN_FLIGHT.has(o.status))
      const results = await Promise.all(
        tracked.map((o) =>
          api
            .getOrderStatus(o.id)
            .then((s) => [o.id, s])
            .catch(() => null)
        )
      )
      setStatuses((prev) => {
        const next = { ...prev }
        for (const r of results) {
          if (r) next[r[0]] = r[1]
        }
        return next
      })
    } catch (e) {
      setError(e.message)
    }
  }, [])

  useEffect(() => {
    load()
    const id = setInterval(load, 5000)
    return () => clearInterval(id)
  }, [load])

  // A buyer can cancel their own buy only before payment is observed. Once
  // canceled, the market blocks them from re-buying that same product.
  const cancelOrder = async (orderId) => {
    if (!window.confirm('Cancel this buy? You will not be able to buy this product again.')) {
      return
    }
    setBusyId(orderId)
    setError('')
    try {
      const resp = await api.cancelOrder(orderId)
      setNotice(resp.message || 'Order canceled.')
      setTimeout(() => setNotice(''), 6000)
      await load()
    } catch (e) {
      setError(e.message)
    } finally {
      setBusyId('')
    }
  }

  const canCancel = (order) => order.type === 'buy' && liveStatus(order) === 'pending_payment'

  // Themed status badges (all white text — no dark text on the dark UI).
  const getStatusBadge = (status) => {
    const statusMap = {
      pending_payment: { label: 'Pending Payment', class: 'pending_payment' },
      paid: { label: 'Paid', class: 'paid' },
      confirmed: { label: 'Confirmed', class: 'confirmed' },
      active: { label: 'Active', class: 'active' },
      completed: { label: 'Completed', class: 'completed' },
      cancelled: { label: 'Cancelled', class: 'cancelled' },
      canceled: { label: 'Cancelled', class: 'canceled' },
      expired: { label: 'Expired', class: 'expired' },
      pending: { label: 'Pending', class: 'pending' },
    }
    const info = statusMap[status] || { label: status, class: 'disabled' }
    return <span className={`badge ${info.class}`}>{info.label}</span>
  }

  const getTypeBadge = (type) =>
    type === 'buy' ? (
      <span className="badge buy">Buy</span>
    ) : (
      <span className="badge sell">Sell</span>
    )

  // For the buyer's pending order: exactly how much to pay and to which wallet.
  // Shown inline (with a copy button) so the payment target is never lost.
  const renderPay = (order) => {
    const status = liveStatus(order)
    if (order.type !== 'buy' || status !== 'pending_payment' || !order.seller_wallet) {
      return <span className="text-muted">—</span>
    }
    return (
      <div>
        <div><strong>{order.expected_payment_amount}</strong> {order.payment_currency}</div>
        <div className="copy-row mt-1">
          <code className="addr-box addr-sm" title={order.seller_wallet}>{order.seller_wallet}</code>
          <CopyButton text={order.seller_wallet} label="Copy" />
          <QrButton
            address={order.seller_wallet}
            amount={order.expected_payment_amount}
            coin={order.payment_currency}
            title="Pay the seller"
            hint="Pay the exact amount so the market can match your payment."
          />
        </div>
      </div>
    )
  }

  // Live view of an order: prefer the polled status over the list snapshot.
  const liveStatus = (order) => statuses[order.id]?.status || order.status

  // Confirmation progress cell for in-flight buy orders.
  const renderProgress = (order) => {
    const live = statuses[order.id]
    const status = liveStatus(order)
    if (order.type !== 'buy') return <span className="text-muted">—</span>
    if (status === 'completed') return <span className="text-success">Delivered</span>
    if (!live || !IN_FLIGHT.has(status)) return <span className="text-muted">—</span>
    if (status === 'pending_payment') {
      return <span className="text-muted">Awaiting payment…</span>
    }
    const confs = live.confirmations || 0
    // The market reports its own threshold; fall back to the default only if an
    // older market omits it.
    const required = live.required_confirmations || REQUIRED_CONFIRMATIONS
    const pct = Math.min(100, Math.round((confs / required) * 100))
    return (
      <div style={{ minWidth: '120px' }}>
        <div className="d-flex justify-content-between small mb-1">
          <span>Confirmations</span>
          <span>{confs}/{required}</span>
        </div>
        <div className="progress" style={{ height: '6px' }}>
          <div
            className="progress-bar bg-info"
            role="progressbar"
            style={{ width: `${pct}%` }}
            aria-valuenow={confs}
            aria-valuemin={0}
            aria-valuemax={required}
          />
        </div>
      </div>
    )
  }

  const renderTx = (order) => {
    const tx = statuses[order.id]?.payment_tx_hash
    if (!tx) return <span className="text-muted">—</span>
    return <code title={tx}>{shorten(tx)}</code>
  }

  return (
    <div>
      <div className="page-head"><h2>My Orders</h2></div>

      {notice && <div className="alert alert-info">{notice}</div>}
      {error && <div className="alert alert-danger">{error}</div>}

      {orders.length === 0 ? (
        <div className="alert alert-secondary">You don't have any orders yet.</div>
      ) : (
        <div className="panel table-wrap">
          <table className="table table-hover align-middle">
            <thead>
              <tr>
                <th>ID</th>
                <th>Type</th>
                <th>Amount</th>
                <th>Price</th>
                <th>Pay to (seller wallet)</th>
                <th>Status</th>
                <th>Progress</th>
                <th>Payment Tx</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {orders.map((order) => (
                <tr key={order.id}>
                  <td><code>{shorten(order.id)}</code></td>
                  <td>{getTypeBadge(order.type)}</td>
                  <td>{order.amount} {order.sell_coin}</td>
                  <td>{order.price} {order.payment_currency}</td>
                  <td>{renderPay(order)}</td>
                  <td>{getStatusBadge(liveStatus(order))}</td>
                  <td>{renderProgress(order)}</td>
                  <td>{renderTx(order)}</td>
                  <td>{order.created_at ? new Date(order.created_at).toLocaleString() : '-'}</td>
                  <td>
                    {canCancel(order) ? (
                      <button
                        className="btn btn-sm btn-outline-danger"
                        disabled={busyId === order.id}
                        onClick={() => cancelOrder(order.id)}
                      >
                        {busyId === order.id ? 'Canceling…' : 'Cancel'}
                      </button>
                    ) : (
                      <span className="text-muted">—</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

function shorten(id) {
  if (!id || id.length <= 12) return id
  return `${id.slice(0, 8)}…`
}

export default Orders
