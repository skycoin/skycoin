import { EventEmitter } from '@angular/core';
import { Observable, from, of, throwError } from 'rxjs';
import { catchError, map, mergeMap } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { BigNumber } from 'bignumber.js';

import { Address, Wallet } from '../../../app.datatypes';
import { defaultCoinId } from '../../../constants/coins-id.const';
import { convertAsciiToHexa } from '../../../utils/converters';
import { environment } from '../../../../environments/environment';
import { CipherProvider, GenerateAddressResponse } from '../../cipher.provider';
import { ApiService } from '../../api.service';
import { EncryptionService } from '../../encryption.service';
import { WalletBackend } from './wallet-backend';

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
    });

    return of(wallets);
  }

  create(label: string, seed: string, coinId: number, existing: Wallet[]): Observable<Wallet> {
    return this.cipherProvider.generateAddress(convertAsciiToHexa(seed)).pipe(
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
          walletType: 'deterministic',
        };
        if (existing.some((wlt: Wallet) =>
            wlt.addresses[0].address === wallet.addresses[0].address &&
            wlt.coinId === wallet.coinId)) {
          throw new Error(this.translate.instant('service.wallet.wallet-exists'));
        }

        return wallet;
      }));
  }

  addAddress(wallet: Wallet): Observable<void> {
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
      if (wallet.encryptedSeed) {
        stripped.encryptedSeed = wallet.encryptedSeed;
        if (wallet.encryptedNextSeed) {
          stripped.encryptedNextSeed = wallet.encryptedNextSeed;
        }
      } else if (wallet.seed) {
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

  getSeed(wallet: Wallet, password: string | null): Observable<string> {
    if (wallet.seed) {
      return of(wallet.seed);
    }
    if (!wallet.encryptedSeed) {
      return throwError(() => new Error(this.translate.instant('service.wallet.address-without-seed')));
    }

    return from(this.encryptionService.decrypt(wallet.encryptedSeed, password || '')).pipe(
      catchError(() => throwError(() => new Error('Incorrect password'))));
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
    const currentSeed = convertAsciiToHexa(seed);

    return this.unlockWalletAddresses(currentSeed, wallet, 0, onProgressChanged).pipe(
      map((res: boolean) => {
        if (!res) {
          throw new Error(this.translate.instant('service.wallet.wrong-seed'));
        }

        wallet.seed = seed;
      }));
  }

  setWalletPassword(wallet: Wallet, password: string): Observable<void> {
    if (!wallet.seed) {
      return throwError(() => new Error('Wallet must be unlocked before setting a password'));
    }

    return new Observable<void>(observer => {
      Promise.all([
        this.encryptionService.encrypt(wallet.seed!, password),
        wallet.nextSeed ? this.encryptionService.encrypt(wallet.nextSeed, password) : Promise.resolve(null),
      ]).then(([encSeed, encNextSeed]) => {
        wallet.encryptedSeed = encSeed;
        if (encNextSeed) {
          wallet.encryptedNextSeed = encNextSeed;
        }
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
      ];

      Promise.all(decryptPromises).then(([seed, nextSeed]) => {
        const clean = cleanSeed(seed!);
        const currentSeed = convertAsciiToHexa(clean);

        this.unlockWalletAddresses(currentSeed, wallet, 0, onProgressChanged).subscribe(
          (res: boolean) => {
            if (!res) {
              observer.error(new Error(this.translate.instant('service.wallet.wrong-seed')));
              return;
            }
            wallet.seed = clean;
            if (nextSeed) {
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

  removeWalletPassword(wallet: Wallet): void {
    wallet.encryptedSeed = undefined;
    wallet.encryptedNextSeed = undefined;
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
