import { Injectable, EventEmitter } from '@angular/core';
import { BehaviorSubject, Subscription, Observable, of } from 'rxjs';
import { mergeMap, map, filter, first } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';

import { CipherProvider } from '../cipher.provider';
import { Address, Wallet } from '../../app.datatypes';
import { BaseCoin } from '../../coins/basecoin';
import { CoinService } from '../coin.service';
import { ApiService } from '../api.service';
import { EncryptionService } from '../encryption.service';
import { WalletBackend } from './backends/wallet-backend';
import { XPubKey } from './backends/server-wallet-backend';
import { ServerWalletBackend } from './backends/server-wallet-backend';
import { BrowserWalletBackend, ScanProgressData } from './backends/browser-wallet-backend';

export { ScanProgressData };

/**
 * The wallets of the current coin. Where they are kept is the coin's backend:
 * a skycoin wallet service when the coin has serverWallets, else this browser.
 */
@Injectable()
export class WalletService {
  wallets: BehaviorSubject<Wallet[]> = new BehaviorSubject<Wallet[]>(null as any);
  // True while this browser's wallets are stored under a password that has not
  // been entered yet. No wallet loads until it is.
  locked = new BehaviorSubject<boolean>(false);

  private currentCoin!: BaseCoin;
  private coinSubscription!: Subscription;
  private readonly server: ServerWalletBackend;
  private readonly browser: BrowserWalletBackend;

  constructor(
    cipherProvider: CipherProvider,
    private translate: TranslateService,
    private coinService: CoinService,
    apiService: ApiService,
    encryptionService: EncryptionService,
  ) {
    this.server = new ServerWalletBackend(apiService, () => this.currentCoin?.id);
    this.browser = new BrowserWalletBackend(cipherProvider, translate, apiService, encryptionService);
    this.loadWallets();
    this.coinService.currentCoin.subscribe((coin) => this.currentCoin = coin);
  }

  /** The backend that keeps the current coin's wallets. */
  private get backend(): WalletBackend {
    return this.currentCoin && this.currentCoin.serverWallets ? this.server : this.browser;
  }

  /** True when wallet is kept by a wallet service rather than this browser. */
  isServerWallet(wallet: Wallet): boolean {
    return this.backend.onServer && !!wallet.filename;
  }

  get haveWallets(): Observable<boolean> {
    return this.wallets.pipe(
      filter(wallets => wallets !== null),
      map(wallets => wallets.length > 0));
  }

  get currentWallets(): Observable<Wallet[]> {
    return this.wallets.pipe(
      filter(wallets => wallets !== null),
      mergeMap(wallets => this.coinService.currentCoin.pipe(
        filter((coin: BaseCoin) => coin !== null),
        map((coin: BaseCoin) => {
          if (coin.serverWallets) {
            // Server-managed wallets are fetched per-coin from the backend
            return wallets;
          }
          return wallets.filter(wallet => wallet.coinId === coin.id);
        })
      )),
      map(wallets => wallets ? wallets : []));
  }

  /** Every address of every wallet, including those of a BIP44 wallet's other accounts. */
  get addresses(): Observable<Address[]> {
    return this.currentWallets.pipe(map(wallets => {
      const byAddress = new Map<string, Address>();
      wallets.forEach(wallet => walletAddresses(wallet).forEach(a => {
        if (!byAddress.has(a.address)) {
          byAddress.set(a.address, a);
        }
      }));
      return Array.from(byAddress.values());
    }));
  }

  addAddress(wallet: Wallet, saveWallet = true, accountIndex?: number, chainIndex?: number): Observable<void> {
    if (this.isServerWallet(wallet)) {
      return this.server.addAddress(wallet, accountIndex, chainIndex).pipe(map(() => this.saveWallets()));
    }

    return this.browser.addAddress(wallet).pipe(map(() => {
      if (saveWallet) {
        this.saveWallets();
      }
    }));
  }

  /**
   * Makes a wallet. A non-empty password encrypts it before it is first stored:
   * the wallet service encrypts on create, the browser before saving.
   */
  create(label: string, seed: string, coinId: number, save = true, walletType = 'deterministic', seedPassphrase?: string,
    segwit?: boolean, password?: string): Observable<Wallet> {
    seed = this.getCleanSeed(seed);
    const backend = this.backend;

    return backend.create(label, seed, coinId, this.wallets.value || [], walletType, seedPassphrase, segwit, password).pipe(
      mergeMap(wallet => password && !backend.onServer
        ? this.browser.setWalletPassword(wallet, password).pipe(map(() => wallet))
        : of(wallet)),
      map(wallet => {
        if (save) {
          this.add(wallet);
        }

        return wallet;
      }));
  }

  add(wallet: Wallet) {
    this.wallets.value.push(wallet);
    this.saveWallets();
  }

  delete(wallet: Wallet) {
    const index = this.wallets.value.indexOf(wallet);
    if (index !== -1) {
      this.wallets.value.splice(index, 1);

      this.saveWallets();
    }
  }

  scanAddresses(wallet: Wallet, onProgressChanged: EventEmitter<ScanProgressData>): Observable<void> {
    return this.browser.scanAddresses(wallet, onProgressChanged).pipe(map(() => this.saveWallets()));
  }

  unlockWallet(wallet: Wallet, seed: string, onProgressChanged: EventEmitter<number>): Observable<void> {
    return this.browser.unlockWallet(wallet, this.getCleanSeed(seed), onProgressChanged);
  }

  /**
   * Derives the keys of a browser wallet whose seed is stored but whose keys
   * are not: they are never saved, so a reload drops them.
   */
  ensureKeys(wallet: Wallet): Observable<void> {
    if (wallet.isHardware || this.isServerWallet(wallet) || !wallet.seed || wallet.addresses.every(a => a.secret_key)) {
      return of(undefined);
    }

    return this.browser.unlockWallet(wallet, this.getCleanSeed(wallet.seed), new EventEmitter<number>());
  }

  setWalletPassword(wallet: Wallet, password: string): Observable<void> {
    return this.browser.setWalletPassword(wallet, password).pipe(map(() => this.saveWallets()));
  }

  unlockWalletWithPassword(wallet: Wallet, password: string, onProgressChanged: EventEmitter<number>): Observable<void> {
    return this.browser.unlockWalletWithPassword(wallet, password, onProgressChanged, seed => this.getCleanSeed(seed));
  }

  removeWalletPassword(wallet: Wallet): void {
    this.browser.removeWalletPassword(wallet);
    this.saveWallets();
  }

  saveWallets() {
    const currentWallets = this.wallets.value || [];
    this.backend.save(currentWallets);
    this.wallets.next(currentWallets);
  }

  get hasAppPassword(): boolean {
    return this.browser.hasAppPassword;
  }

  /** Decrypts the wallets stored in this browser and loads them. */
  unlockApp(password: string): Observable<void> {
    if (!this.browser.hasLockedStore()) {
      return of(undefined);
    }

    return new Observable<void>(observer => {
      this.browser.unlock(password).then(() => {
        this.locked.next(false);
        this.loadWalletsFrom(this.browser);
        observer.next();
        observer.complete();
      }, err => observer.error(err));
    });
  }

  /** Stores this browser's wallets under a password, replacing any previous one. */
  setAppPassword(password: string): Promise<void> {
    return this.browser.setAppPassword(password);
  }

  /** Stores this browser's wallets without a password again. */
  removeAppPassword(currentPassword: string): Promise<void> {
    return this.browser.removeAppPassword(currentPassword);
  }

  checkAppPassword(password: string): boolean {
    return this.browser.checkAppPassword(password);
  }

  /** Forgets the password and the wallets it guarded. This cannot be undone. */
  eraseLockedWallets() {
    this.browser.erase();
    this.locked.next(false);
    this.wallets.next([]);
  }

  /** The wallet's seed, for a backup. password unlocks an encrypted one. */
  getSeed(wallet: Wallet, password: string | null): Observable<string> {
    return (this.isServerWallet(wallet) ? this.server : this.browser).getSeed(wallet, password);
  }

  /** The folder the wallet service keeps the files in, or null for wallets in this browser. */
  walletFolder(): Observable<string | null> {
    return this.backend.onServer ? this.server.folder() : of(null);
  }

  getXPubKey(wallet: Wallet, path: string): Observable<XPubKey> {
    return this.server.getXPubKey(wallet, path);
  }

  addAccount(wallet: Wallet, name: string, password?: string): Observable<Wallet> {
    return this.server.addAccount(wallet, name, password).pipe(map(w => {
      this.saveWallets();

      return w;
    }));
  }

  encryptWallet(wallet: Wallet, password: string): Observable<void> {
    return this.server.encrypt(wallet, password);
  }

  decryptWallet(wallet: Wallet, password: string): Observable<void> {
    return this.server.decrypt(wallet, password);
  }

  private loadWallets() {
    // Wait for coins to load so we know whether to use server or client-side wallets.
    // This prevents the wallet wizard from opening before the current coin is known.
    this.coinService.coinsLoaded.pipe(first()).subscribe(() => {
      const coin = this.coinService.currentCoin.getValue();
      this.loadWalletsFrom(coin && coin.serverWallets ? this.server : this.browser);

      // Re-fetch wallets when the user switches coins
      this.coinSubscription = this.coinService.currentCoin.pipe(
        filter((c: BaseCoin) => c !== null),
      ).subscribe((c) => this.loadWalletsFrom(c.serverWallets ? this.server : this.browser));
    });
  }

  private loadWalletsFrom(backend: WalletBackend) {
    backend.load().subscribe(wallets => {
      if (wallets === null) {
        this.locked.next(true);
        return;
      }
      this.wallets.next(wallets);
    });
  }

  private getCleanSeed(seed: string): string {
    return seed.replace(/(\n|\r\n)$/, '');
  }
}

/**
 * A wallet's address objects, those of its BIP44 accounts included. The
 * accounts list their addresses as separate objects, which may repeat.
 */
export function walletAddresses(wallet: Wallet): Address[] {
  const accounts = (wallet.accounts || []).map(a => [...a.externalAddresses, ...a.changeAddresses]);
  return ([] as Address[]).concat(wallet.addresses, ...accounts);
}
