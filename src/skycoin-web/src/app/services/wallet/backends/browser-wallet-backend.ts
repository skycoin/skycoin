import { EventEmitter } from '@angular/core';
import { Observable, forkJoin, from, of, throwError } from 'rxjs';
import { catchError, map, mergeMap } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { BigNumber } from 'bignumber.js';

import { Address, Bip44Account, Wallet } from '../../../app.datatypes';
import { defaultCoinId } from '../../../constants/coins-id.const';
import { convertAsciiToHexa } from '../../../utils/converters';
import { environment } from '../../../../environments/environment';
import { Bip44Address, CipherProvider, GenerateAddressResponse, pqSeedState } from '../../cipher.provider';
import { ApiService } from '../../api.service';
import { EncryptionService } from '../../encryption.service';
import { SeedBackup, WalletBackend } from './wallet-backend';

/** Unused addresses in a row that end a bip44 chain scan, the BIP44 gap limit. */
export const bip44GapLimit = 20;

/** True for a bip44 wallet kept in this browser. */
export function isBip44(wallet: Wallet): boolean {
  return wallet.walletType === 'bip44' && !wallet.filename && !!wallet.accounts && wallet.accounts.length > 0;
}

/**
 * Lists a browser bip44 wallet's addresses, external then change, as the
 * same objects its account holds, so keys and balances set on one show on both.
 */
function linkBip44Addresses(wallet: Wallet) {
  const account: Bip44Account = wallet.accounts![0];
  account.externalAddresses.forEach((a, i) => (a as Bip44Address).child_number = i);
  account.changeAddresses.forEach((a, i) => (a as Bip44Address).child_number = i);
  wallet.addresses = [...account.externalAddresses, ...account.changeAddresses];
}

export class ScanProgressData {
  addressesFound = 1;
  progress = 0;
}

/**
 * Wallets kept in this browser's localStorage, with keys derived by
 * skycoin-lite. The whole stored list may be encrypted under an app password.
 */
export class BrowserWalletBackend implements WalletBackend {
  readonly onServer = false;

  // The password the wallets in this browser are stored under, kept in memory
  // only, and the stored list it decrypted.
  private appPassword: string | null = null;
  private unlockedStoredWallets: string | null = null;
  private saves: Promise<void> = Promise.resolve();
  private readonly walletsStorageKey = 'wallets';
  private readonly lockedWalletsStorageKey = 'lockedWallets';

  constructor(
    private cipherProvider: CipherProvider,
    private translate: TranslateService,
    private apiService: ApiService,
    private encryptionService: EncryptionService,
  ) { }

  load(): Observable<Wallet[] | null> {
    if (this.appPassword === null && localStorage.getItem(this.lockedWalletsStorageKey)) {
      return of(null);
    }
    const storedWallets: string | null = this.unlockedStoredWallets ?? localStorage.getItem(this.walletsStorageKey);
    if (!storedWallets) {
      return of([]);
    }
    const wallets: Wallet[] = JSON.parse(storedWallets);
    wallets.forEach((wallet) => {
      if (!wallet.coinId) {
        wallet.coinId = defaultCoinId;
      }
      // A persisted plaintext seed means the wallet is usable immediately, so
      // clear the seed-confirmation gate — the re-entry prompt is only for
      // wallets whose seed we don't hold (encrypted-and-locked, or a legacy
      // addresses-only entry created before seeds were persisted).
      if (wallet.seed) {
        wallet.needSeedConfirmation = false;
      }
      if (isBip44(wallet)) {
        linkBip44Addresses(wallet);
      }
    });

    return of(wallets);
  }

  create(label: string, seed: string, coinId: number, existing: Wallet[], walletType = 'deterministic',
    seedPassphrase = ''): Observable<Wallet> {
    if (walletType === 'bip44') {
      return this.createBip44(label, seed, seedPassphrase, coinId, existing);
    }

    const isPQ = walletType === 'pq';

    return this.cipherProvider.generateAddress(this.initialSeedState(seed, walletType)).pipe(
      map((response: GenerateAddressResponse) => {
        const wallet: Wallet = {
          label: label,
          seed: seed,
          needSeedConfirmation: true,
          balance: new BigNumber('0'),
          hours: new BigNumber('0'),
          addresses: [response.address],
          nextSeed: response.nextSeed,
          coinId: coinId,
          walletType: isPQ ? 'pq' : 'deterministic',
        };
        if (existing.some((wlt: Wallet) =>
            wlt.addresses[0].address === wallet.addresses[0].address &&
            wlt.coinId === wallet.coinId)) {
          throw new Error(this.translate.instant('service.wallet.wallet-exists'));
        }

        return wallet;
      }));
  }

  /**
   * Makes a bip44 wallet: account 0 of coin type 8000, derived from a bip39
   * mnemonic and its seed passphrase the way the node's bip44 wallets are.
   */
  private createBip44(label: string, seed: string, seedPassphrase: string, coinId: number, existing: Wallet[]): Observable<Wallet> {
    return forkJoin([
      this.cipherProvider.bip44AccountXPub(seed, seedPassphrase),
      this.cipherProvider.bip44Addresses(seed, seedPassphrase, 0, 0, 1),
    ]).pipe(map(([xpub, external]) => {
      const wallet: Wallet = {
        label: label,
        seed: seed,
        seedPassphrase: seedPassphrase,
        needSeedConfirmation: true,
        balance: new BigNumber('0'),
        hours: new BigNumber('0'),
        addresses: [],
        coinId: coinId,
        walletType: 'bip44',
        accounts: [{ name: 'Account 0', index: 0, accountXpubKey: xpub, externalAddresses: external, changeAddresses: [] }],
      };
      linkBip44Addresses(wallet);
      if (existing.some(wlt => wlt.addresses[0] && wlt.addresses[0].address === wallet.addresses[0].address && wlt.coinId === coinId)) {
        throw new Error(this.translate.instant('service.wallet.wallet-exists'));
      }

      return wallet;
    }));
  }

  addAddress(wallet: Wallet, accountIndex?: number, chainIndex?: number): Observable<void> {
    if (isBip44(wallet)) {
      return this.addBip44Addresses(wallet, chainIndex === 1 ? 1 : 0, 1).pipe(map(() => undefined));
    }

    if (!wallet.seed || !wallet.nextSeed) {
      throw new Error(this.translate.instant('service.wallet.address-without-seed'));
    }

    return this.cipherProvider.generateAddress(wallet.nextSeed).pipe(
      map((response: GenerateAddressResponse) => {
        wallet.nextSeed = response.nextSeed;
        wallet.addresses.push(response.address);
      }));
  }

  save(wallets: Wallet[]): void {
    const strippedWallets: Wallet[] = [];
    wallets.forEach(wallet => {
      const strippedAddresses: Address[] = [];
      wallet.addresses.forEach(address => strippedAddresses.push({ address: address.address }));
      const stripped: Wallet = {
        coinId: wallet.coinId,
        needSeedConfirmation: wallet.needSeedConfirmation,
        label: wallet.label,
        addresses: strippedAddresses,
        isHardware: wallet.isHardware,
      };
      if (wallet.walletType === 'pq') {
        stripped.walletType = wallet.walletType;
      }
      if (isBip44(wallet)) {
        // The xpub is public and stays readable, so addresses can be derived
        // and scanned without the password.
        stripped.walletType = wallet.walletType;
        stripped.accounts = wallet.accounts!.map(a => ({
          name: a.name,
          index: a.index,
          accountXpubKey: a.accountXpubKey,
          externalAddresses: a.externalAddresses.map(e => ({ address: e.address })),
          changeAddresses: a.changeAddresses.map(e => ({ address: e.address })),
        }));
      }
      if (wallet.encryptedSeed) {
        stripped.encryptedSeed = wallet.encryptedSeed;
        if (wallet.encryptedNextSeed) {
          stripped.encryptedNextSeed = wallet.encryptedNextSeed;
        }
        if (wallet.encryptedSeedPassphrase) {
          stripped.encryptedSeedPassphrase = wallet.encryptedSeedPassphrase;
        }
      } else if (wallet.seed) {
        if (wallet.seedPassphrase) {
          stripped.seedPassphrase = wallet.seedPassphrase;
        }
        // No password → persist the plaintext seed so the wallet survives a
        // refresh WITHOUT forcing a seed re-entry. Set a password to encrypt
        // the seed at rest (then only encryptedSeed is stored). This is the
        // convenience default for the in-browser client-side wallet.
        stripped.seed = wallet.seed;
        if (wallet.nextSeed) {
          stripped.nextSeed = wallet.nextSeed;
        }
      }
      strippedWallets.push(stripped);
    });
    this.storeWallets(JSON.stringify(strippedWallets));
  }

  /**
   * The seed and seed passphrase. An encrypted wallet always needs its password,
   * even while its seed is in memory.
   */
  getSeed(wallet: Wallet, password: string | null): Observable<SeedBackup> {
    if (wallet.encryptedSeed) {
      return from(this.decryptSecrets(wallet, password || '')).pipe(
        map(secrets => ({ seed: secrets.seed, passphrase: secrets.passphrase || undefined })));
    }
    if (!wallet.seed) {
      return throwError(() => new Error(this.translate.instant('service.wallet.address-without-seed')));
    }

    return of({ seed: wallet.seed, passphrase: wallet.seedPassphrase || undefined });
  }

  /** Decrypts every encrypted secret of a wallet, or fails with 'Incorrect password'. */
  private async decryptSecrets(wallet: Wallet, password: string): Promise<{ seed: string, nextSeed?: string, passphrase?: string }> {
    try {
      const [seed, nextSeed, passphrase] = await Promise.all([
        this.encryptionService.decrypt(wallet.encryptedSeed!, password),
        wallet.encryptedNextSeed ? this.encryptionService.decrypt(wallet.encryptedNextSeed, password) : Promise.resolve(undefined),
        wallet.encryptedSeedPassphrase ? this.encryptionService.decrypt(wallet.encryptedSeedPassphrase, password) : Promise.resolve(undefined),
      ]);

      return { seed, nextSeed, passphrase };
    } catch {
      throw new Error('Incorrect password');
    }
  }

  get hasAppPassword(): boolean {
    return this.appPassword !== null;
  }

  // Writes the stored wallet list, encrypted when this browser has a password.
  // Writes are chained so a slower encryption never lands after a newer save.
  private storeWallets(serialized: string): Promise<void> {
    const password = this.appPassword;
    this.unlockedStoredWallets = password !== null ? serialized : null;
    this.saves = this.saves.then(async () => {
      if (password === null) {
        localStorage.setItem(this.walletsStorageKey, serialized);
        localStorage.removeItem(this.lockedWalletsStorageKey);
      } else {
        localStorage.setItem(this.lockedWalletsStorageKey, await this.encryptionService.encrypt(serialized, password));
        localStorage.removeItem(this.walletsStorageKey);
      }
    });

    return this.saves;
  }

  private currentStoredWallets(): string {
    return this.unlockedStoredWallets ?? localStorage.getItem(this.walletsStorageKey) ?? '[]';
  }

  /** True when there is a stored list under a password; false when there is nothing to unlock. */
  hasLockedStore(): boolean {
    return !!localStorage.getItem(this.lockedWalletsStorageKey);
  }

  /** Decrypts the stored list with password, so load returns it. */
  unlock(password: string): Promise<void> {
    const encrypted = localStorage.getItem(this.lockedWalletsStorageKey);

    return this.encryptionService.decrypt(encrypted!, password).then(stored => {
      this.appPassword = password;
      this.unlockedStoredWallets = stored;
    }, () => {
      throw new Error('Incorrect password');
    });
  }

  /** Stores this browser's wallets under a password, replacing any previous one. */
  setAppPassword(password: string): Promise<void> {
    const stored = this.currentStoredWallets();
    this.appPassword = password;

    return this.storeWallets(stored);
  }

  /** Stores this browser's wallets without a password again. */
  removeAppPassword(currentPassword: string): Promise<void> {
    if (currentPassword !== this.appPassword) {
      return Promise.reject(new Error('Incorrect password'));
    }
    const stored = this.currentStoredWallets();
    this.appPassword = null;

    return this.storeWallets(stored);
  }

  checkAppPassword(password: string): boolean {
    return this.appPassword !== null && password === this.appPassword;
  }

  /** Forgets the password and the wallets it guarded. This cannot be undone. */
  erase() {
    localStorage.removeItem(this.lockedWalletsStorageKey);
    localStorage.removeItem(this.walletsStorageKey);
    this.appPassword = null;
    this.unlockedStoredWallets = null;
  }

  /** Finds how many of the wallet's addresses have been used, keeping those. */
  scanAddresses(wallet: Wallet, onProgressChanged: EventEmitter<ScanProgressData>): Observable<void> {
    if (isBip44(wallet)) {
      return this.discoverBip44Addresses(wallet).pipe(map(() => {
        onProgressChanged.emit({ addressesFound: wallet.addresses.length, progress: 100 });
      }));
    }

    if (wallet.addresses.length !== 1) {
      throw new Error(this.translate.instant('service.wallet.invalid-wallet'));
    }

    const InitialNextSeed = wallet.nextSeed;

    return this.checkWalletAddresses(wallet, 0, onProgressChanged, InitialNextSeed!).pipe(
      map(lastIndexWithTxs => {
        const unnecessaryAddresses = wallet.addresses.length - 1 - lastIndexWithTxs;
        if (unnecessaryAddresses > 0) {
          wallet.addresses.splice(lastIndexWithTxs + 1, unnecessaryAddresses);
        }
      }),
      catchError(error => {
        if (wallet.addresses.length > 1) {
          wallet.addresses.splice(1, wallet.addresses.length - 1);
          wallet.nextSeed = InitialNextSeed;
        }
        return throwError(() => error);
      }));
  }

  unlockWallet(wallet: Wallet, seed: string, onProgressChanged: EventEmitter<number>): Observable<void> {
    if (isBip44(wallet)) {
      return this.unlockBip44(wallet, seed, wallet.seedPassphrase || '', onProgressChanged);
    }

    const currentSeed = this.initialSeedState(seed, wallet.walletType);

    return this.unlockWalletAddresses(currentSeed, wallet, 0, onProgressChanged).pipe(
      map((res: boolean) => {
        if (!res) {
          throw new Error(this.translate.instant('service.wallet.wrong-seed'));
        }

        wallet.seed = seed;
      }));
  }

  setWalletPassword(wallet: Wallet, password: string): Observable<void> {
    // Every secret is encrypted again under the new password, so each one must
    // be in memory; otherwise its old ciphertext would be dropped or left under
    // a password nobody uses anymore.
    if (!wallet.seed || (wallet.encryptedNextSeed && !wallet.nextSeed) ||
      (wallet.encryptedSeedPassphrase && wallet.seedPassphrase === undefined)) {
      return throwError(() => new Error('Wallet must be unlocked before setting a password'));
    }

    return new Observable<void>(observer => {
      Promise.all([
        this.encryptionService.encrypt(wallet.seed!, password),
        wallet.nextSeed ? this.encryptionService.encrypt(wallet.nextSeed, password) : Promise.resolve(null),
        wallet.seedPassphrase ? this.encryptionService.encrypt(wallet.seedPassphrase, password) : Promise.resolve(null),
      ]).then(([encSeed, encNextSeed, encPassphrase]) => {
        wallet.encryptedSeed = encSeed;
        wallet.encryptedNextSeed = encNextSeed || undefined;
        wallet.encryptedSeedPassphrase = encPassphrase || undefined;
        observer.next();
        observer.complete();
      }).catch(err => observer.error(err));
    });
  }

  unlockWalletWithPassword(wallet: Wallet, password: string, onProgressChanged: EventEmitter<number>,
    cleanSeed: (seed: string) => string): Observable<void> {
    if (!wallet.encryptedSeed) {
      return throwError(() => new Error('Wallet does not have an encrypted seed'));
    }

    return new Observable<void>(observer => {
      const decryptPromises = [
        this.encryptionService.decrypt(wallet.encryptedSeed!, password),
        wallet.encryptedNextSeed ? this.encryptionService.decrypt(wallet.encryptedNextSeed, password) : Promise.resolve(null),
        wallet.encryptedSeedPassphrase ? this.encryptionService.decrypt(wallet.encryptedSeedPassphrase, password) : Promise.resolve(''),
      ];

      Promise.all(decryptPromises).then(([seed, nextSeed, passphrase]) => {
        const clean = cleanSeed(seed!);
        if (isBip44(wallet)) {
          this.unlockBip44(wallet, clean, passphrase || '', onProgressChanged).subscribe(() => {
            observer.next();
            observer.complete();
          }, err => observer.error(err));

          return;
        }
        const currentSeed = this.initialSeedState(clean, wallet.walletType);

        this.unlockWalletAddresses(currentSeed, wallet, 0, onProgressChanged).subscribe(
          (res: boolean) => {
            if (!res) {
              observer.error(new Error(this.translate.instant('service.wallet.wrong-seed')));
              return;
            }
            wallet.seed = clean;
            // The stored next seed is only written with the password, so it can
            // be older than the addresses. The one derived above follows them.
            if (nextSeed && !wallet.nextSeed) {
              wallet.nextSeed = nextSeed;
            }
            observer.next();
            observer.complete();
          },
          err => observer.error(err)
        );
      }).catch(() => {
        observer.error(new Error('Incorrect password'));
      });
    });
  }

  /**
   * Stores the wallet unencrypted again. Its secrets are decrypted with password
   * and checked against its addresses first, and nothing changes unless all of
   * that succeeds, so a locked wallet cannot lose the seed it only has encrypted.
   */
  removeWalletPassword(wallet: Wallet, password: string): Observable<void> {
    if (!wallet.encryptedSeed) {
      return throwError(() => new Error('Wallet does not have an encrypted seed'));
    }

    return from(this.decryptSecrets(wallet, password)).pipe(
      mergeMap(secrets => this.checkSecrets(wallet, secrets.seed, secrets.passphrase || '').pipe(map(derivedNextSeed => {
        wallet.seed = secrets.seed;
        if (secrets.passphrase !== undefined) {
          wallet.seedPassphrase = secrets.passphrase;
        }
        if (!isBip44(wallet)) {
          wallet.nextSeed = secrets.nextSeed || derivedNextSeed;
        }
        wallet.encryptedSeed = undefined;
        wallet.encryptedNextSeed = undefined;
        wallet.encryptedSeedPassphrase = undefined;
      }))));
  }

  /**
   * Stores an encrypted wallet under a new password, or unencrypted when password
   * is empty, using its seed instead of the old password. Nothing changes unless
   * seed and passphrase give the wallet's stored xpub or addresses.
   */
  resetWalletPassword(wallet: Wallet, seed: string, passphrase: string, password: string | null): Observable<void> {
    if (!wallet.encryptedSeed) {
      return throwError(() => new Error('Wallet does not have an encrypted seed'));
    }
    const bip44 = isBip44(wallet);
    if (!bip44) {
      passphrase = '';
    }

    return this.checkSecrets(wallet, seed, passphrase).pipe(
      mergeMap(derivedNextSeed => {
        const nextSeed = bip44 ? undefined : derivedNextSeed;
        if (!password) {
          return of({ nextSeed, encrypted: null as (string | null)[] | null });
        }

        return from(Promise.all([
          this.encryptionService.encrypt(seed, password),
          nextSeed ? this.encryptionService.encrypt(nextSeed, password) : Promise.resolve(null),
          passphrase ? this.encryptionService.encrypt(passphrase, password) : Promise.resolve(null),
        ])).pipe(map(encrypted => ({ nextSeed, encrypted: encrypted as (string | null)[] | null })));
      }),
      map(({ nextSeed, encrypted }) => {
        wallet.seed = seed;
        if (bip44) {
          wallet.seedPassphrase = passphrase;
        } else {
          wallet.nextSeed = nextSeed;
        }
        wallet.encryptedSeed = encrypted ? encrypted[0]! : undefined;
        wallet.encryptedNextSeed = encrypted && encrypted[1] ? encrypted[1] : undefined;
        wallet.encryptedSeedPassphrase = encrypted && encrypted[2] ? encrypted[2] : undefined;
      }));
  }

  /**
   * Fails unless seed and passphrase give the wallet's stored xpub or addresses.
   * Emits the next seed after the last address of a deterministic or pq wallet,
   * and leaves the wallet as it was.
   */
  private checkSecrets(wallet: Wallet, seed: string, passphrase: string): Observable<string | undefined> {
    const wrongSeed = () => new Error(this.translate.instant('service.wallet.wrong-seed'));
    if (isBip44(wallet)) {
      return this.cipherProvider.bip44AccountXPub(seed, passphrase).pipe(map(xpub => {
        if (xpub !== wallet.accounts![0].accountXpubKey) {
          throw wrongSeed();
        }

        return undefined;
      }));
    }

    const copy: Wallet = { label: wallet.label, walletType: wallet.walletType, addresses: wallet.addresses.map(a => ({ address: a.address })) };
    if (copy.addresses.length === 0) {
      return throwError(wrongSeed);
    }

    return this.unlockWalletAddresses(this.initialSeedState(seed, wallet.walletType), copy, 0, new EventEmitter<number>()).pipe(
      map(ok => {
        if (!ok) {
          throw wrongSeed();
        }

        return copy.nextSeed;
      }));
  }

  /**
   * Derives the keys of every address of a bip44 wallet from its mnemonic,
   * failing when they do not give the stored xpub and addresses.
   */
  private unlockBip44(wallet: Wallet, seed: string, passphrase: string, onProgressChanged: EventEmitter<number>): Observable<void> {
    const account = wallet.accounts![0];

    return forkJoin([
      this.cipherProvider.bip44AccountXPub(seed, passphrase),
      this.cipherProvider.bip44Addresses(seed, passphrase, 0, 0, account.externalAddresses.length),
      this.cipherProvider.bip44Addresses(seed, passphrase, 1, 0, account.changeAddresses.length),
    ]).pipe(map(([xpub, external, change]) => {
      const matches = (stored: Address[], derived: Bip44Address[]) => stored.every((a, i) => a.address === derived[i].address);
      if (xpub !== account.accountXpubKey || !matches(account.externalAddresses, external) || !matches(account.changeAddresses, change)) {
        onProgressChanged.emit(0);
        throw new Error(this.translate.instant('service.wallet.wrong-seed'));
      }

      const setKeys = (stored: Address[], derived: Bip44Address[]) => stored.forEach((a, i) => {
        a.secret_key = derived[i].secret_key;
        a.public_key = derived[i].public_key;
      });
      setKeys(account.externalAddresses, external);
      setKeys(account.changeAddresses, change);
      wallet.seed = seed;
      wallet.seedPassphrase = passphrase;
      onProgressChanged.emit(100);
    }));
  }

  /** Appends num addresses to a chain of a bip44 wallet, derived from its xpub. */
  private addBip44Addresses(wallet: Wallet, chain: number, num: number): Observable<Address[]> {
    const account = wallet.accounts![0];
    const list = chain === 1 ? account.changeAddresses : account.externalAddresses;

    return this.cipherProvider.bip44XPubAddresses(account.accountXpubKey!, chain, list.length, num).pipe(map(derived => {
      list.push(...derived);
      linkBip44Addresses(wallet);

      return derived;
    }));
  }

  /** A fresh change address for a bip44 wallet, added to it. */
  nextChangeAddress(wallet: Wallet): Observable<string> {
    return this.addBip44Addresses(wallet, 1, 1).pipe(map(added => added[0].address));
  }

  /**
   * The address a send's change goes to unless the user names another, not
   * added to the wallet yet: the next change chain address of a bip44 wallet,
   * or the next address of a deterministic one. Null when there is none to
   * give, as for a hardware or post-quantum wallet or one whose keys are not
   * in memory, and then the node picks the change address.
   */
  peekChangeAddress(wallet: Wallet): Observable<string | null> {
    if (wallet.isHardware || wallet.walletType === 'pq') {
      return of(null);
    }
    if (isBip44(wallet)) {
      const account = wallet.accounts![0];

      return this.cipherProvider.bip44XPubAddresses(account.accountXpubKey!, 1, account.changeAddresses.length, 1).pipe(
        map(derived => derived[0].address));
    }
    if (!wallet.seed || !wallet.nextSeed) {
      return of(null);
    }

    // A next seed that does not derive leaves the change to the node, as before.
    return this.cipherProvider.generateAddress(wallet.nextSeed).pipe(
      map(response => response.address.address),
      catchError(() => of(null)),
    );
  }

  /**
   * Adds address to the wallet once a sent transaction paid it, if it is still
   * the one peekChangeAddress gives. Emits whether it was added.
   */
  addChangeAddress(wallet: Wallet, address: string): Observable<boolean> {
    return this.peekChangeAddress(wallet).pipe(mergeMap(next => {
      if (next !== address) {
        return of(false);
      }
      if (isBip44(wallet)) {
        return this.nextChangeAddress(wallet).pipe(map(added => added === address));
      }

      return this.addAddress(wallet).pipe(map(() => true));
    }));
  }

  /**
   * Adds the addresses the blockchain shows were used beyond the ones a bip44
   * wallet already has, on both chains, stopping after bip44GapLimit unused
   * ones in a row. Coins sent to an address the wallet never showed, as the
   * reward system does with addresses it derives from the xpub, are found so.
   * Emits how many addresses were added.
   */
  discoverBip44Addresses(wallet: Wallet): Observable<number> {
    const account = wallet.accounts![0];
    const scanChain = (chain: number, added: number): Observable<number> => {
      const list = chain === 1 ? account.changeAddresses : account.externalAddresses;
      const start = list.length;

      return this.cipherProvider.bip44XPubAddresses(account.accountXpubKey!, chain, start, bip44GapLimit).pipe(
        mergeMap(window => this.apiService.get('transactions', { addrs: window.map(a => a.address).join(',') } as any).pipe(
          map((transactions: any[]) => {
            const used = new Set<string>();
            (transactions || []).forEach(t => ((t.txn || t).outputs || []).forEach((o: any) => used.add(o.dst)));
            let last = -1;
            window.forEach((a, i) => {
              if (used.has(a.address)) {
                last = i;
              }
            });

            return { window, last };
          }))),
        mergeMap(({ window, last }) => {
          if (last < 0) {
            return of(added);
          }
          list.push(...window.slice(0, last + 1));
          linkBip44Addresses(wallet);

          return scanChain(chain, added + last + 1);
        }));
    };

    return scanChain(0, 0).pipe(mergeMap(added => scanChain(1, added)));
  }

  /** The state the first address of a wallet with this seed is derived from. */
  private initialSeedState(seed: string, walletType?: string): string {
    const seedHex = convertAsciiToHexa(seed);
    return walletType === 'pq' ? pqSeedState(seedHex, 0) : seedHex;
  }

  private checkWalletAddresses(wallet: Wallet, lastIndexWithTxs: number, onProgressChanged: EventEmitter<ScanProgressData>, nextSeed: string): Observable<number> {
    const minAdrressesToScan = environment.e2eTest ? 2 : 10;
    const maxAdrressesToScan = environment.e2eTest ? 2 : 100;

    return this.addAddress(wallet).pipe(
      mergeMap(() => this.apiService.get('transactions', { addrs: wallet.addresses[wallet.addresses.length - 1].address } as any)),
      mergeMap(transactions => {
        if (transactions && transactions.length > 0) {
          lastIndexWithTxs = wallet.addresses.length - 1;
          nextSeed = wallet.nextSeed!;
        } else {
          wallet.nextSeed = nextSeed;
        }

        onProgressChanged.emit({
          addressesFound: lastIndexWithTxs + 1,
          progress: (wallet.addresses.length - 1 - lastIndexWithTxs) / minAdrressesToScan * 100
        });

        if (lastIndexWithTxs + minAdrressesToScan === wallet.addresses.length - 1 || wallet.addresses.length === maxAdrressesToScan) {
          return of(lastIndexWithTxs);
        } else {
          return this.checkWalletAddresses(wallet, lastIndexWithTxs, onProgressChanged, nextSeed);
        }
      }));
  }

  private unlockWalletAddresses(currentSeed: string, wallet: Wallet, index: number, onProgressChanged: EventEmitter<number>): Observable<boolean> {
    return this.cipherProvider.generateAddress(currentSeed).pipe(
      mergeMap((response: GenerateAddressResponse) => {
        if (response.address.address !== wallet.addresses[index].address) {
          onProgressChanged.emit(0);
          return of(false);
        }

        onProgressChanged.emit((index + 1) / wallet.addresses.length * 100);
        wallet.nextSeed = response.nextSeed;
        wallet.addresses[index].public_key = response.address.public_key;
        wallet.addresses[index].secret_key = response.address.secret_key;
        index++;

        if (index === wallet.addresses.length) {
          return of(true);
        }

        return this.unlockWalletAddresses(response.nextSeed, wallet, index, onProgressChanged);
      }));
  }
}
