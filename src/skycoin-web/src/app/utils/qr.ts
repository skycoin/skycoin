/** The largest frame edge handed to a decoder; a readable code survives it. */
const maxEdge = 640;

export interface ScannedPayment {
  address: string;
  amount?: string;
}

/**
 * Reads a scanned address, or a payment URI such as skycoin:ADDR?amount=1
 * or bitcoin:ADDR?amount=0.1 (BIP 21). An ethereum: URI keeps only its
 * address, since its value is in wei.
 */
export function parseScannedPayment(text: string): ScannedPayment {
  let rest = text.trim();
  const scheme = /^([a-z][a-z0-9+.-]*):/i.exec(rest);
  if (scheme) {
    rest = rest.substring(scheme[0].length).replace(/^\/\//, '');
  }
  const [path, query] = rest.split('?', 2);
  // EIP-681 adds @chain and /function after the address.
  const address = path.split('@')[0].split('/')[0];
  let amount: string | undefined;
  if (query && !(scheme && scheme[1].toLowerCase() === 'ethereum')) {
    const value = new URLSearchParams(query).get('amount');
    if (value && /^\d*\.?\d+$/.test(value)) {
      amount = value;
    }
  }

  return { address, amount };
}

/**
 * Decodes a QR code in an image or video frame, or returns null. The Android
 * app lends its own decoder as window.SkywireQr. Elsewhere jsQR, loaded on
 * first use, decodes in the page.
 */
export async function decodeQr(source: HTMLVideoElement | HTMLImageElement, canvas: HTMLCanvasElement): Promise<string | null> {
  const width = (source as HTMLVideoElement).videoWidth || (source as HTMLImageElement).naturalWidth || 0;
  const height = (source as HTMLVideoElement).videoHeight || (source as HTMLImageElement).naturalHeight || 0;
  if (!width || !height) {
    return null;
  }
  const scale = Math.min(1, maxEdge / Math.max(width, height));
  canvas.width = Math.max(1, Math.round(width * scale));
  canvas.height = Math.max(1, Math.round(height * scale));
  const ctx = canvas.getContext('2d', { willReadFrequently: true })!;
  ctx.drawImage(source, 0, 0, canvas.width, canvas.height);

  const host = (window as any).SkywireQr;
  if (host && typeof host.decode === 'function') {
    return host.decode(canvas.toDataURL('image/jpeg', 0.8)) || null;
  }
  const jsQR = (await import('jsqr')).default;
  const image = ctx.getImageData(0, 0, canvas.width, canvas.height);
  const code = jsQR(image.data, image.width, image.height);

  return code ? code.data : null;
}
