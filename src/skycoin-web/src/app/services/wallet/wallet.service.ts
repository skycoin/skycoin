import { Injectable, EventEmitter } from '@angular/core';
import { BehaviorSubject, Subscription, Observable, of } from 'rxjs';
import { catchError, mergeMap, map, filter, finalize, first, share } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';

import { CipherProvider } from '../cipher.provider';
import { Address, Wallet } from '../../app.datatypes';
import { BaseCoin } from '../../coins/basecoin';
import { CoinService } from '../coin.service';
import { ApiService } from '../api.service';
import { EncryptionService } from '../encryption.service';
import { SeedBackup, WalletBackend } from './backends/wallet-backend';
import { XPubKey } from './backends/server-wallet-backend';
import { ServerWalletBackend } from './backends/server-wallet-backend';
import { BrowserWalletBackend, ScanProgressData, bip44GapLimit, isBip44 } from './backends/browser-wallet-backend';

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
  // The scan of each wallet's used addresses under way, and when each last started.
  private readonly scans = new Map<string, Observable<number>>();
  private readonly lastScan = new Map<string, number>();
  private readonly autoScanInterval = 60 * 1000;

  constructor(
    private cipherProvider: CipherProvider,
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

    return this.browser.addAddress(wallet, accountIndex, chainIndex).pipe(map(() => {
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
    // Gateway coins keep their wallets on a server; this browser derives bip44
    // only for skycoin-style coins.
    const coin = this.coinService.coins.find(c => c.id === coinId);
    if (!backend.onServer && walletType === 'bip44' && coin && coin.isGateway()) {
      walletType = 'deterministic';
    }

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

  /** Stores a browser wallet unencrypted again, once password has decrypted it. */
  removeWalletPassword(wallet: Wallet, password: string): Observable<void> {
    return this.browser.removeWalletPassword(wallet, password).pipe(map(() => this.saveWallets()));
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

  /** The wallet's seed and seed passphrase, for a backup. password unlocks an encrypted one. */
  getSeed(wallet: Wallet, password: string | null): Observable<SeedBackup> {
    return (this.isServerWallet(wallet) ? this.server : this.browser).getSeed(wallet, password);
  }

  /** The folder the wallet service keeps the files in, or null for wallets in this browser. */
  walletFolder(): Observable<string | null> {
    return this.backend.onServer ? this.server.folder() : of(null);
  }

  getXPubKey(wallet: Wallet, path: string): Observable<XPubKey> {
    return this.server.getXPubKey(wallet, path);
  }

  /** True for a bip44 wallet kept in this browser, whose xpub is stored with it. */
  isBrowserBip44(wallet: Wallet): boolean {
    return !this.isServerWallet(wallet) && isBip44(wallet);
  }

  /** The first n external addresses a browser bip44 wallet's xpub derives. */
  xpubAddresses(wallet: Wallet, n: number): Observable<string[]> {
    return this.cipherProvider.bip44XPubAddresses(wallet.accounts![0].accountXpubKey!, 0, 0, n).pipe(
      map(addresses => addresses.map(a => a.address)));
  }

  /**
   * The address change goes to when a send names none: an unused address of a
   * browser wallet that is not in the wallet yet, the next change chain address
   * for bip44 or the next address for a deterministic wallet. Null leaves it to
   * the node, which for a wallet it keeps also picks a new address.
   */
  changeAddressFor(wallet: Wallet): Observable<string | null> {
    if (this.isServerWallet(wallet)) {
      return of(null);
    }

    return this.browser.peekChangeAddress(wallet);
  }

  /** Adds the change address a sent transaction paid to its browser wallet. */
  recordChangeAddress(wallet: Wallet, address: string): Observable<void> {
    if (this.isServerWallet(wallet)) {
      return of(undefined);
    }

    return this.browser.addChangeAddress(wallet, address).pipe(map(added => {
      if (added) {
        this.saveWallets();
      }
    }));
  }

  /** True for a bip44 or xpub wallet, whose used addresses can be found from its public keys. */
  canFindUsedAddresses(wallet: Wallet): boolean {
    if (this.currentCoin && this.currentCoin.isGateway()) {
      return false;
    }
    if (this.isServerWallet(wallet)) {
      return wallet.walletType === 'bip44' || wallet.walletType === 'xpub';
    }

    return isBip44(wallet);
  }

  /**
   * Adds the addresses the blockchain shows were used past the last ones a
   * bip44 or xpub wallet lists, until bip44GapLimit in a row are unused, as
   * reward payouts to new xpub indexes are. Emits how many it added.
   */
  findUsedAddresses(wallet: Wallet): Observable<number> {
    if (!this.canFindUsedAddresses(wallet)) {
      return of(0);
    }
    const key = this.scanKey(wallet);
    let scan = this.scans.get(key);
    if (!scan) {
      this.lastScan.set(key, Date.now());
      const run = this.isServerWallet(wallet)
        ? this.server.scanAddresses(wallet, bip44GapLimit)
        : this.cipherProvider.initialize().pipe(mergeMap(() => this.browser.discoverBip44Addresses(wallet)));
      scan = run.pipe(
        map(added => {
          if (!(this.wallets.value || []).includes(wallet)) {
            // The list was loaded again meanwhile and this object is gone, so scan the new one.
            this.lastScan.delete(key);
            Promise.resolve().then(() => this.findAllUsedAddresses());

            return 0;
          }
          if (added > 0) {
            this.saveWallets();
          }

          return added;
        }),
        finalize(() => this.scans.delete(key)),
        share());
      this.scans.set(key, scan);
    }

    return scan;
  }

  /**
   * Finds the used addresses of every wallet that can, at most once a minute
   * for each. The balance refresh runs it, which waits for the node to be in sync.
   */
  findAllUsedAddresses() {
    const now = Date.now();
    (this.wallets.value || []).filter(w => this.canFindUsedAddresses(w)).forEach(w => {
      const last = this.lastScan.get(this.scanKey(w));
      if (last === undefined || now - last >= this.autoScanInterval) {
        this.findUsedAddresses(w).pipe(catchError(() => of(0))).subscribe();
      }
    });
  }

  private scanKey(wallet: Wallet): string {
    return this.isServerWallet(wallet) ? 'file:' + wallet.filename : wallet.accounts![0].accountXpubKey!;
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
