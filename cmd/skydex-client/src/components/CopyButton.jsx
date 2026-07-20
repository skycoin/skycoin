import { useState } from 'react'
import { copyToClipboard } from '../api'

// CopyButton copies `text` to the clipboard and briefly confirms. It works over
// plain http (visor) via the clipboard fallback in api.copyToClipboard.
function CopyButton({ text, className = 'btn btn-sm', label = 'Copy' }) {
  const [done, setDone] = useState(false)

  const onCopy = async () => {
    const ok = await copyToClipboard(text)
    if (ok) {
      setDone(true)
      setTimeout(() => setDone(false), 1500)
    }
  }

  return (
    <button type="button" className={className} onClick={onCopy} disabled={!text} title="Copy to clipboard">
      {done ? 'Copied ✓' : label}
    </button>
  )
}

export default CopyButton
