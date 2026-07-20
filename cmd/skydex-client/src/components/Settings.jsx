import { useEffect, useState } from 'react'
import { api, loadWallets, saveWallets } from '../api'

// Settings is where a user registers the addresses the market uses for them: the
// required Skycoin address (receives every purchased sell coin + refunds) and a
// payout address for each payment coin the market accepts — including any custom
// coin the operator added. Values are remembered locally and pre-filled here.
function Settings({ marketPubKey, onRegistered }) {
  const saved = loadWallets()
  const [sky, setSky] = useState(saved.SKY || '')
  const [payouts, setPayouts] = useState({})
  const [currencies, setCurrencies] = useState([])
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState(null)

  useEffect(() => {
    api
      .getCurrencies()
      .then((d) => {
        const cur = d.currencies || []
        setCurrencies(cur)
        setPayouts(() => {
          const init = {}
          cur.forEach((c) => {
            init[c] = saved[c] || ''
          })
          return init
        })
      })
      .catch(() => {})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const setPayout = (cur, val) => setPayouts((p) => ({ ...p, [cur]: val }))

  const handleSave = async () => {
    setMessage(null)
    if (!sky.trim()) {
      setMessage({ type: 'danger', text: 'A Skycoin address is required.' })
      return
    }
    const wallets = {}
    Object.entries(payouts).forEach(([cur, addr]) => {
      if (addr.trim()) wallets[cur] = addr.trim()
    })
    setSaving(true)
    try {
      await api.register({ wallet_sky: sky.trim(), wallets })
      saveWallets({ SKY: sky.trim(), ...wallets })
      setMessage({ type: 'success', text: 'Addresses saved. You can now trade.' })
      if (onRegistered) onRegistered()
    } catch (e) {
      setMessage({ type: 'danger', text: e.message })
    } finally {
      setSaving(false)
    }
  }

  return (
    <div>
      <div className="page-head">
        <h2>Settings</h2>
      </div>

      <div className="panel">
        <h4 className="panel-title">Your wallet addresses</h4>
        <p className="text-muted mb-3">
          Your <strong>Skycoin</strong> address receives every purchased sell coin (SKY and any
          fibercoin share the format) and any refund. To <strong>sell</strong>, add a payout address
          for each payment coin you’ll accept — that’s where buyers pay you directly.
        </p>

        <div className="field-grid">
          <div>
            <label className="form-label">Skycoin address <span className="req">*</span></label>
            <input
              type="text"
              className="form-control"
              placeholder="Skycoin-family address (receives SKY + all fibercoins)"
              value={sky}
              onChange={(e) => setSky(e.target.value)}
            />
          </div>
          {currencies.map((cur) => (
            <div key={cur}>
              <label className="form-label">{cur} payout address</label>
              <input
                type="text"
                className="form-control"
                placeholder={`${cur} address (to be paid when selling for ${cur})`}
                value={payouts[cur] || ''}
                onChange={(e) => setPayout(cur, e.target.value)}
              />
            </div>
          ))}
        </div>
        {currencies.length === 0 && (
          <div className="hint mt-2">This market has no external payment coins enabled yet.</div>
        )}

        {message && <div className={`alert alert-${message.type} mt-3 mb-0`}>{message.text}</div>}

        <button className="btn btn-connect mt-3" onClick={handleSave} disabled={saving}>
          {saving ? 'Saving…' : 'Save addresses'}
        </button>
      </div>

      <div className="panel">
        <h4 className="panel-title">Market connection</h4>
        <label className="form-label">Connected market public key</label>
        <input type="text" className="form-control" value={marketPubKey} readOnly />
        <div className="hint mt-1">Use Disconnect in the top bar to switch markets.</div>
      </div>
    </div>
  )
}

export default Settings
