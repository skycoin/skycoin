function Header({ isConnected, marketPubKey, marketName, onDisconnect }) {
  return (
    <header className="header">
      <div className="brand">
        Skywire <span>Exchange</span>
      </div>

      <div className="header-right">
        <div className="status-indicator">
          <span className={`status-dot ${isConnected ? 'connected' : ''}`} />
          <span className="status-text">{isConnected ? 'Connected' : 'Disconnected'}</span>
        </div>

        {marketName && <span className="market-name" title={marketPubKey}>{marketName}</span>}

        {marketPubKey && (
          <span className="market-pubkey" title={marketPubKey}>
            {marketPubKey.slice(0, 8)}…{marketPubKey.slice(-6)}
          </span>
        )}

        {isConnected && (
          <button className="btn btn-outline-light btn-sm" onClick={onDisconnect}>
            Disconnect
          </button>
        )}
      </div>
    </header>
  )
}

export default Header
