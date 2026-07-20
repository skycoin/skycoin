import { useEffect, useId, useState } from 'react'
import { QRCodeSVG } from 'qrcode.react'
import CopyButton from './CopyButton'

// BIP-21 style URI schemes, keyed by coin symbol. Only coins whose scheme *and*
// `amount` parameter we are sure about live here: every entry below takes the
// amount as a decimal coin value (not wei/satoshi), so `<scheme>:<addr>?amount=<n>`
// is unambiguous. Anything missing falls back to a bare address QR, which every
// wallet can scan — an unrecognised URI, by contrast, scans as garbage.
//
// The market operator can enable arbitrary sell coins / payment currencies at
// runtime, so this map is a best-effort lookup and never a whitelist.
const URI_SCHEMES = {
  BTC: 'bitcoin',
  LTC: 'litecoin',
  DOGE: 'dogecoin',
  DASH: 'dash',
  BCH: 'bitcoincash',
  // Skycoin and the fibercoins share the address format and the skycoin: URI.
  SKY: 'skycoin',
}

// Amounts arrive from the market as plain JSON floats with no decimals metadata,
// and JS stringifies small ones in exponent form ("1e-7") — which a wallet parsing
// a BIP-21 amount will reject or, worse, misread. Force plain decimal notation and
// trim the trailing zeros that toFixed adds. 8 decimals covers every UTXO coin
// here (satoshi/litoshi) and Skycoin's 6-decimal droplets.
function formatAmount(n) {
  return n.toFixed(8).replace(/\.?0+$/, '')
}

// buildUri renders the payment URI for `coin`, or null when we have no scheme we
// trust for it (unknown coin, or no amount to encode).
export function buildUri(address, amount, coin) {
  const scheme = URI_SCHEMES[String(coin || '').toUpperCase()]
  const n = Number(amount)
  if (!scheme || !address || !Number.isFinite(n) || n <= 0) return null
  return `${scheme}:${address}?amount=${formatAmount(n)}`
}

// QrButton shows a small QR glyph next to an address. Clicking it opens a modal
// with a full-size, scannable QR so a phone wallet can pay without retyping the
// address — the transcription step this replaces is where funds get lost.
//
// When the coin has a known URI scheme the QR carries address *and* amount, so
// the scanning wallet pre-fills both; a toggle falls back to an address-only QR
// for wallets that don't understand the URI.
function QrButton({ address, amount, coin, title = 'Payment', hint }) {
  const [open, setOpen] = useState(false)
  // Tables render one QrButton per row, so the toggle needs a per-instance id to
  // keep its <label for=...> pointing at its own checkbox.
  const toggleId = useId()
  const uri = buildUri(address, amount, coin)
  // Only offer the URI form when we have one; otherwise the QR is the raw address.
  const [withAmount, setWithAmount] = useState(true)
  const value = uri && withAmount ? uri : address

  // Esc closes the modal, matching the backdrop click.
  useEffect(() => {
    if (!open) return undefined
    const onKey = (e) => {
      if (e.key === 'Escape') setOpen(false)
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [open])

  if (!address) return null

  return (
    <>
      <button
        type="button"
        className="btn btn-sm qr-btn"
        onClick={() => setOpen(true)}
        title="Show QR code"
        aria-label="Show QR code"
      >
        {/* Inline glyph rather than an icon font: the UI is embedded in the Go
            binary and must render with no network. */}
        <svg width="14" height="14" viewBox="0 0 16 16" aria-hidden="true" focusable="false">
          <path
            fill="currentColor"
            d="M0 0h7v7H0V0zm2 2v3h3V2H2zM0 9h7v7H0V9zm2 2v3h3v-3H2zM9 0h7v7H9V0zm2 2v3h3V2h-3zM9 9h2v2H9V9zm5 0h2v2h-2V9zm-3 3h2v2h-2v-2zm3 0h2v4h-4v-2h2v-2z"
          />
        </svg>
      </button>

      {open && (
        <div
          className="qr-backdrop"
          role="dialog"
          aria-modal="true"
          aria-label={`${title} QR code`}
          onClick={() => setOpen(false)}
        >
          {/* Stop propagation so clicks inside the card don't dismiss it. */}
          <div className="qr-modal" onClick={(e) => e.stopPropagation()}>
            <div className="deposit-head">
              <h4 className="panel-title mb-0">{title}</h4>
              <button
                type="button"
                className="deposit-close"
                onClick={() => setOpen(false)}
                aria-label="Close"
              >
                ✕
              </button>
            </div>

            {amount && coin && (
              <p className="qr-amount mb-2">
                Send exactly <strong>{amount} {coin}</strong>
              </p>
            )}

            {/* QR stays dark-on-white regardless of the app's navy theme —
                inverted codes are unreliable to scan. */}
            <div className="qr-canvas">
              <QRCodeSVG value={value} size={220} level="M" marginSize={2} bgColor="#FFFFFF" fgColor="#000000" />
            </div>

            {uri && (
              <div className="form-check qr-toggle">
                <input
                  className="form-check-input"
                  type="checkbox"
                  id={toggleId}
                  checked={withAmount}
                  onChange={(e) => setWithAmount(e.target.checked)}
                />
                <label className="form-check-label" htmlFor={toggleId}>
                  Include amount in QR ({coin} payment link)
                </label>
              </div>
            )}

            <div className="qr-addr">
              <code className="addr-box">{address}</code>
            </div>

            <div className="copy-row qr-actions">
              <CopyButton text={address} label="Copy address" />
              {uri && withAmount && <CopyButton text={uri} label="Copy payment link" />}
            </div>

            <p className="hint mt-2 mb-0">
              {uri && withAmount
                ? `Scan with a ${coin} wallet — the address and amount are filled in for you. If your wallet doesn't recognise the link, untick the box above for an address-only code.`
                : 'Scan to fill in the address. Enter the amount manually in your wallet.'}
              {hint ? ` ${hint}` : ''}
            </p>
          </div>
        </div>
      )}
    </>
  )
}

export default QrButton
