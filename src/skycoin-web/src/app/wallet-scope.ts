/**
 * The class that scopes the wallet's styles (wallet-styles.scss). The shell
 * carries it, and so does every dialog, menu and select pane the wallet opens,
 * since those render outside the shell.
 */
export const WALLET_SCOPE_CLASS = 'skycoin-wallet';

/** classes plus WALLET_SCOPE_CLASS, for a pane or backdrop's class option. */
export function withWalletScope(classes?: string | string[]): string[] {
  return [...(classes ? ([] as string[]).concat(classes) : []), WALLET_SCOPE_CLASS];
}
