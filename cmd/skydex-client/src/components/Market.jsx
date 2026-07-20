import { useEffect, useState, useCallback } from 'react'
import { api } from '../api'
import CopyButton from './CopyButton'
import QrButton from './QrButton'

function Market({ registered, marketReady, sellCoins = [], currencies = [], onNeedRegister }) {
  const [showCreateListing, setShowCreateListing] = useState(false)
  const [products, setProducts] = useState([])
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)
  // Deposit instructions from the last created listing. Kept until the user
  // dismisses it (never auto-hidden) so they don't lose the escrow address.
  const [deposit, setDeposit] = useState(null)
  // Payment instructions from the last buy. Kept until dismissed so the buyer
  // never loses how much to pay and to which seller wallet.
  const [payment, setPayment] = useState(null)

  const [newListing, setNewListing] = useState({
    sell_coin: '',
    amount: '',
    price: '',
    payment_currency: '',
  })

  const loadProducts = useCallback(async () => {
    try {
      const data = await api.getProducts()
      setProducts(data.products || [])
      setError('')
    } catch (e) {
      setError(e.message)
    }
  }, [])

  useEffect(() => {
    loadProducts()
    const id = setInterval(loadProducts, 5000)
    return () => clearInterval(id)
  }, [loadProducts])

  // Seed sensible defaults for the create form once the coin lists arrive.
  useEffect(() => {
    setNewListing((prev) => {
      const sell = prev.sell_coin || sellCoins[0] || ''
      const payOpts = [...new Set([...currencies, ...sellCoins.filter((c) => c !== sell)])]
      return {
        ...prev,
        sell_coin: sell,
        payment_currency: prev.payment_currency || payOpts[0] || '',
      }
    })
  }, [sellCoins, currencies])

  // Surface the market's "register first" as an actionable prompt, not a dead end.
  const handleErr = (e) => {
    if (/register/i.test(e.message) && onNeedRegister) onNeedRegister()
    setError(e.message)
  }

  const handleBuy = async (product) => {
    if (!registered) return onNeedRegister && onNeedRegister()
    setError('')
    setBusy(true)
    try {
      const order = await api.buyProduct(product.id)
      // Persist the payment instructions on-screen (with a copy button) instead
      // of a toast the buyer can miss — this is where they send their payment.
      setPayment({
        amount: order.amount,
        sell_coin: order.sell_coin,
        expected_payment_amount: order.expected_payment_amount,
        payment_currency: order.payment_currency,
        seller_wallet: order.seller_wallet,
        expires_at: order.expires_at,
      })
      loadProducts()
    } catch (e) {
      handleErr(e)
    } finally {
      setBusy(false)
    }
  }

  const openCreate = () => {
    if (!registered) return onNeedRegister && onNeedRegister()
    setShowCreateListing((v) => !v)
  }

  const handleCreateListing = async (e) => {
    e.preventDefault()
    setError('')
    setBusy(true)
    try {
      const resp = await api.createListing({
        sell_coin: newListing.sell_coin,
        amount: parseFloat(newListing.amount),
        price: parseFloat(newListing.price),
        payment_currency: newListing.payment_currency,
      })
      // Persist the deposit instructions on-screen (with a copy button) instead
      // of a toast the user can miss — this is the address they must fund.
      setDeposit({
        sell_coin: resp.sell_coin,
        amount: resp.amount,
        commission: resp.commission,
        expected_amount: resp.expected_amount,
        market_wallet: resp.market_wallet,
        expires_at: resp.expires_at,
      })
      setShowCreateListing(false)
      setNewListing({ sell_coin: sellCoins[0] || '', amount: '', price: '', payment_currency: '' })
    } catch (err) {
      handleErr(err)
    } finally {
      setBusy(false)
    }
  }

  // Payment options = external explorer coins plus any other fibercoin (you can't
  // pay in the coin you're selling).
  const paymentOptions = [...new Set([...currencies, ...sellCoins.filter((c) => c !== newListing.sell_coin)])]
  const canSell = sellCoins.length > 0

  return (
    <div>
      <div className="page-head">
        <h2>Market</h2>
        <button className="btn btn-connect" onClick={openCreate} disabled={!canSell}>
          {showCreateListing ? 'Cancel' : '+ New Sell Order'}
        </button>
      </div>

      {error && <div className="alert alert-danger">{error}</div>}

      {deposit && (
        <div className="panel deposit-box">
          <div className="deposit-head">
            <h4 className="panel-title mb-0">Fund your sell order</h4>
            <button type="button" className="deposit-close" onClick={() => setDeposit(null)} aria-label="Dismiss">✕</button>
          </div>
          <p className="mb-2">
            Send exactly <strong>{deposit.expected_amount} {deposit.sell_coin}</strong>{' '}
            ({deposit.amount} + {deposit.commission} commission) to this escrow address
            {deposit.expires_at ? <> before <strong>{new Date(deposit.expires_at).toLocaleTimeString()}</strong></> : ' within 15 minutes'}:
          </p>
          <div className="copy-row">
            <code className="addr-box">{deposit.market_wallet}</code>
            <CopyButton text={deposit.market_wallet} label="Copy address" />
            <QrButton
              address={deposit.market_wallet}
              amount={deposit.expected_amount}
              coin={deposit.sell_coin}
              title="Fund your sell order"
              hint="This is the market escrow address, not the buyer's."
            />
          </div>
          <p className="hint mt-2 mb-0">
            You can find this again anytime under <strong>My Listings</strong>. The offer becomes a
            purchasable product once the deposit is confirmed on-chain.
          </p>
        </div>
      )}

      {payment && (
        <div className="panel deposit-box">
          <div className="deposit-head">
            <h4 className="panel-title mb-0">Pay the seller to complete your purchase</h4>
            <button type="button" className="deposit-close" onClick={() => setPayment(null)} aria-label="Dismiss">✕</button>
          </div>
          <p className="mb-2">
            You are buying <strong>{payment.amount} {payment.sell_coin}</strong>. Pay exactly{' '}
            <strong>{payment.expected_payment_amount} {payment.payment_currency}</strong> to this seller wallet
            {payment.expires_at ? <> before <strong>{new Date(payment.expires_at).toLocaleTimeString()}</strong></> : ''}:
          </p>
          <div className="copy-row">
            <code className="addr-box">{payment.seller_wallet}</code>
            <CopyButton text={payment.seller_wallet} label="Copy address" />
            <QrButton
              address={payment.seller_wallet}
              amount={payment.expected_payment_amount}
              coin={payment.payment_currency}
              title="Pay the seller"
              hint="Pay the exact amount so the market can match your payment."
            />
          </div>
          <p className="hint mt-2 mb-0">
            You can find this again anytime under <strong>My Orders</strong>. Pay the exact amount so the
            market can match your payment; the {payment.sell_coin} is released to you after confirmation.
          </p>
        </div>
      )}

      {showCreateListing && (
        <div className="panel trade-builder mb-4">
          <h4 className="panel-title">New sell order</h4>
          <form onSubmit={handleCreateListing}>
            {/* "You sell" leg — the coin + amount you give up. */}
            <div className="trade-leg">
              <span className="leg-label">You sell</span>
              <input
                type="number"
                className="form-control leg-amount"
                placeholder="0.00"
                value={newListing.amount}
                onChange={(e) => setNewListing({ ...newListing, amount: e.target.value })}
                required
                step="0.0001"
                min="0"
                aria-label="amount to sell"
              />
              <select
                className="form-select leg-coin"
                value={newListing.sell_coin}
                onChange={(e) => {
                  const sell = e.target.value
                  const pay = newListing.payment_currency === sell ? '' : newListing.payment_currency
                  setNewListing({ ...newListing, sell_coin: sell, payment_currency: pay })
                }}
                required
                aria-label="sell coin"
              >
                {sellCoins.length === 0 && <option value="">No sell coins</option>}
                {sellCoins.map((c) => (
                  <option key={c} value={c}>{c}</option>
                ))}
              </select>
            </div>

            <div className="trade-arrow">↓</div>

            {/* "You get" leg — price sits next to the currency it's denominated in. */}
            <div className="trade-leg">
              <span className="leg-label">You get</span>
              <input
                type="number"
                className="form-control leg-amount"
                placeholder="0.00"
                value={newListing.price}
                onChange={(e) => setNewListing({ ...newListing, price: e.target.value })}
                required
                step="0.0001"
                min="0"
                aria-label="price"
              />
              <select
                className="form-select leg-coin"
                value={newListing.payment_currency}
                onChange={(e) => setNewListing({ ...newListing, payment_currency: e.target.value })}
                required
                aria-label="payment currency"
              >
                {paymentOptions.length === 0 && <option value="">No currencies</option>}
                {paymentOptions.map((c) => (
                  <option key={c} value={c}>{c}</option>
                ))}
              </select>
            </div>

            <div className="trade-summary">
              {newListing.amount && newListing.price && newListing.sell_coin && newListing.payment_currency ? (
                <>
                  Sell <strong>{newListing.amount} {newListing.sell_coin}</strong> for{' '}
                  <strong>{newListing.price} {newListing.payment_currency}</strong>. After creating this
                  order you have 15 minutes to deposit the {newListing.sell_coin} (amount + a small
                  commission) to the market escrow; the buyer’s payment goes straight to your wallet.
                </>
              ) : (
                <>Set what you sell and what you want for it. Payment can be an external coin (BTC/LTC) or another fibercoin.</>
              )}
            </div>

            <button
              type="submit"
              className="btn btn-connect"
              disabled={busy || paymentOptions.length === 0 || sellCoins.length === 0}
            >
              Create sell order
            </button>
          </form>
        </div>
      )}

      <h4 className="section-title">Available products</h4>
      {products.length === 0 ? (
        <div className="alert alert-secondary">
          {marketReady === false
            ? 'This market has no coins configured yet.'
            : 'No products available right now.'}
        </div>
      ) : (
        <div className="card-grid">
          {products.map((product) => (
            <div key={product.id} className="card product-card">
              <div className="product-amount">
                <strong>{product.amount} {product.sell_coin}</strong>
              </div>
              <div className="product-price">
                {product.price} {product.payment_currency}
              </div>
              <div className="text-muted small product-seller">
                Seller <code>{shorten(product.seller_pubkey)}</code>
              </div>
              <button className="btn btn-connect w-100" disabled={busy} onClick={() => handleBuy(product)}>
                Buy
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

function shorten(pk) {
  if (!pk || pk.length <= 14) return pk
  return `${pk.slice(0, 8)}…${pk.slice(-4)}`
}

export default Market
