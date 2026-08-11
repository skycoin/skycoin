declare var Go: any;

/**
 * Loads assets/scripts/skycoin-lite.wasm and waits until it is usable.
 *
 * go.run() starts the Go program and only settles when that program exits, so
 * it cannot be awaited. The cipher becomes callable once main() has published
 * SkycoinCipher on window, which happens a tick or two after go.run() returns —
 * every spec that reached for it straight after initialising found undefined.
 *
 * wasm_exec.js, which defines Go, is listed under "scripts" on the test target
 * in angular.json.
 */
export async function loadCipherWasm(timeoutMs = 30000): Promise<void> {
  if (!(window as any).SkycoinCipher) {
    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(
      fetch('assets/scripts/skycoin-lite.wasm'), go.importObject);
    go.run(result.instance);
  }

  const deadline = Date.now() + timeoutMs;
  while (!(window as any).SkycoinCipher) {
    if (Date.now() > deadline) {
      throw new Error('the wasm cipher did not publish SkycoinCipher on window');
    }
    await new Promise(resolve => setTimeout(resolve, 20));
  }
}
