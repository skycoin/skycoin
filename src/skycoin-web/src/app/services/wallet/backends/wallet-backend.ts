import { Observable } from 'rxjs';

import { Wallet } from '../../../app.datatypes';

/**
 * Where a coin's wallets are kept, and the operations that differ by where.
 * The server backend keeps them in a skycoin wallet service (the node, or
 * `skycoin web --wallet-dir`); the browser backend keeps them in this browser
 * and derives keys with skycoin-lite.
 */
export interface WalletBackend {
  /** True when the wallets live in a wallet service rather than this browser. */
  readonly onServer: boolean;

  /**
   * The stored wallets. The browser backend emits null while its store is
   * locked by an app password that has not been entered.
   */
  load(): Observable<Wallet[] | null>;

  /**
   * Makes a wallet. existing is the current list, for a duplicate check. The server
   * backend encrypts it with password; the caller encrypts a browser wallet.
   */
  create(label: string, seed: string, coinId: number, existing: Wallet[], walletType: string,
    seedPassphrase?: string, segwit?: boolean, password?: string): Observable<Wallet>;

  /** Adds the next address to wallet, on the given BIP44 account and chain if any. */
  addAddress(wallet: Wallet, accountIndex?: number, chainIndex?: number): Observable<void>;

  /** Persists the list, where the backend has anything to persist. */
  save(wallets: Wallet[]): void;

  /**
   * The wallet's seed, for a backup. password unlocks it where the backend keeps it
   * encrypted; null when it is not.
   */
  getSeed(wallet: Wallet, password: string | null): Observable<string>;
}
