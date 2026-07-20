import { useEffect, useState, useCallback } from 'react'
import Header from './components/Header'
import Market from './components/Market'
import MyListings from './components/MyListings'
import Orders from './components/Orders'
import History from './components/History'
import Settings from './components/Settings'
import ConnectGate from './components/ConnectGate'
import { api, hasSavedWallets, addMarketToHistory } from './api'

const TABS = [
  { id: 'market', label: 'Market' },
  { id: 'orders', label: 'My Orders' },
  { id: 'listings', label: 'My Listings' },
  { id: 'history', label: 'History' },
  { id: 'settings', label: 'Settings' },
]

function App() {
  const [activeTab, setActiveTab] = useState('market')
  const [isConnected, setIsConnected] = useState(false)
  const [marketPubKey, setMarketPubKey] = useState('')
  const [marketName, setMarketName] = useState('')
  const [defaultPubKey, setDefaultPubKey] = useState('')
  const [loading, setLoading] = useState(true)

  // registered is a UX hint (the market is the source of truth): true once the
  // user has saved wallet addresses at least once on this device.
  const [registered, setRegistered] = useState(hasSavedWallets())
  // ready reflects whether the connected market is actually accepting trades —
  // i.e. it has at least one enabled, fully-configured sell coin.
  const [ready, setReady] = useState(null) // null = unknown yet
  const [sellCoins, setSellCoins] = useState([])
  const [currencies, setCurrencies] = useState([])

  useEffect(() => {
    async function init() {
      try {
        const [cfg, status] = await Promise.all([
          fetch('/api/config').then((r) => r.json()),
          fetch('/api/status').then((r) => r.json()),
        ])
        setDefaultPubKey(cfg.market_pk || '')
        if (status.connected) {
          setIsConnected(true)
          setMarketPubKey(status.market_pk || '')
          setMarketName(status.market_name || '')
          if (status.market_pk) addMarketToHistory(status.market_pk, status.market_name)
        }
      } catch (e) {
        console.error('init failed', e)
      } finally {
        setLoading(false)
      }
    }
    init()
  }, [])

  // Poll the market's readiness (available sell coins / payment currencies) while
  // connected, so the UI can tell the user when a market has nothing configured.
  const refreshReadiness = useCallback(async () => {
    try {
      const d = await api.getCurrencies()
      const sc = d.sell_coins || []
      const cur = d.currencies || []
      setSellCoins(sc)
      setCurrencies(cur)
      setReady(sc.length > 0)
      // The market name rides along on get_currencies; keep it fresh and remember
      // it against the connected key in the recent-markets list.
      if (d.market_name) {
        setMarketName(d.market_name)
        if (marketPubKey) addMarketToHistory(marketPubKey, d.market_name)
      }
    } catch {
      // Leave readiness unknown on a transient error.
    }
  }, [marketPubKey])

  useEffect(() => {
    if (!isConnected) return undefined
    refreshReadiness()
    const id = setInterval(refreshReadiness, 8000)
    return () => clearInterval(id)
  }, [isConnected, refreshReadiness])

  function handleConnected(pk, name) {
    addMarketToHistory(pk, name)
    setMarketPubKey(pk)
    setMarketName(name || '')
    setIsConnected(true)
    setReady(null)
  }

  async function handleDisconnect() {
    try {
      await fetch('/api/disconnect', { method: 'POST' })
    } catch (e) {
      console.error('disconnect failed', e)
    }
    setIsConnected(false)
    setMarketPubKey('')
    setMarketName('')
    setReady(null)
  }

  if (loading) {
    return (
      <div className="app-container d-flex align-items-center justify-content-center">
        <div className="spinner-border text-light" role="status" />
      </div>
    )
  }

  if (!isConnected) {
    return <ConnectGate defaultPubKey={defaultPubKey} onConnected={handleConnected} />
  }

  const goToSettings = () => setActiveTab('settings')

  return (
    <div className="app-container">
      <Header isConnected={isConnected} marketPubKey={marketPubKey} marketName={marketName} onDisconnect={handleDisconnect} />

      <div className="tabbar">
        <div className="tabbar-inner">
          {TABS.map((t) => (
            <button
              key={t.id}
              className={`tab-link ${activeTab === t.id ? 'active' : ''}`}
              onClick={() => setActiveTab(t.id)}
            >
              {t.label}
            </button>
          ))}
        </div>
      </div>

      <div className="content">
        {/* Onboarding: prompt to set wallet addresses before trading. */}
        {!registered && activeTab !== 'settings' && (
          <div className="banner banner-info">
            <span>
              <strong>Set your wallet addresses first.</strong> You need a Skycoin address (and a
              BTC/LTC payout address to sell) before you can trade.
            </span>
            <button className="btn btn-sm" onClick={goToSettings}>Open Settings</button>
          </div>
        )}
        {/* Market readiness: nothing to trade until the operator configures a coin. */}
        {ready === false && (activeTab === 'market' || activeTab === 'listings') && (
          <div className="banner banner-warn">
            This market isn’t accepting trades yet — the operator hasn’t enabled any sell coin.
            Check back later.
          </div>
        )}

        {activeTab === 'market' && (
          <Market
            registered={registered}
            marketReady={ready}
            sellCoins={sellCoins}
            currencies={currencies}
            onNeedRegister={goToSettings}
          />
        )}
        {activeTab === 'orders' && <Orders />}
        {activeTab === 'listings' && <MyListings onNeedRegister={goToSettings} />}
        {activeTab === 'history' && <History marketPubKey={marketPubKey} marketName={marketName} />}
        {activeTab === 'settings' && (
          <Settings marketPubKey={marketPubKey} onRegistered={() => setRegistered(true)} />
        )}
      </div>
    </div>
  )
}

export default App
