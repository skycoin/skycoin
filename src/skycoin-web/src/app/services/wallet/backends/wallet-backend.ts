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
   * backend encrypts it with password; the caller encrypts a browser wallet. A
   * temporary wallet is never stored, and the server backend loads it with createTemp.
   */
  create(label: string, seed: string, coinId: number, existing: Wallet[], walletType: string,
    seedPassphrase?: string, segwit?: boolean, password?: string, temporary?: boolean): Observable<Wallet>;

  /** Adds the next count addresses to wallet, on the given BIP44 account and chain if any. */
  addAddress(wallet: Wallet, accountIndex?: number, chainIndex?: number, count?: number): Observable<void>;

  /** Persists the list, where the backend has anything to persist. */
  save(wallets: Wallet[]): void;

  /**
   * The wallet's seed and seed passphrase, for a backup. password unlocks them
   * where the backend keeps them encrypted; null when it does not.
   */
  getSeed(wallet: Wallet, password: string | null): Observable<SeedBackup>;
}

/** What restores a wallet: the seed, and the bip39 seed passphrase when it has one. */
export interface SeedBackup {
  seed: string;
  passphrase?: string;
}
