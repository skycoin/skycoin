import { Injectable, EventEmitter, Injector } from '@angular/core';
import { BehaviorSubject, Subscription, Observable, of, throwError } from 'rxjs';
import { mergeMap, map, filter, first, catchError } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { BigNumber } from 'bignumber.js';

import { CipherProvider, GenerateAddressResponse } from '../cipher.provider';
import { Address, Wallet } from '../../app.datatypes';
import { convertAsciiToHexa } from '../../utils/converters';
import { defaultCoinId } from '../../constants/coins-id.const';
import { BaseCoin } from '../../coins/basecoin';
import { CoinService } from '../coin.service';
import { ApiService } from '../api.service';
import { EncryptionService } from '../encryption.service';
import { environment } from '../../../environments/environment';

export class ScanProgressData {
  addressesFound = 1;
  progress = 0;
}

@Injectable()
export class WalletService {
  wallets: BehaviorSubject<Wallet[]> = new BehaviorSubject<Wallet[]>(null);

  private currentCoin: BaseCoin;
  private coinSubscription: Subscription;

  constructor(
    private cipherProvider: CipherProvider,
    private translate: TranslateService,
    private coinService: CoinService,
    private apiService: ApiService,
    private encryptionService: EncryptionService,
  ) {
    this.loadWallets();
    this.coinService.currentCoin.subscribe((coin) => this.currentCoin = coin);
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

  get addresses(): Observable<Address[]> {
    return this.currentWallets.pipe(map(wallets => wallets.reduce((array, wallet) => array.concat(wallet.addresses), [])));
  }

  addAddress(wallet: Wallet, saveWallet = true, accountIndex?: number, chainIndex?: number): Observable<void> {
    // Server-managed wallets generate addresses via the backend API
    if (this.currentCoin && this.currentCoin.serverWallets && wallet.filename) {
      return this.addServerAddress(wallet, accountIndex, chainIndex);
    }

    if (!wallet.seed || !wallet.nextSeed) {
      throw new Error(this.translate.instant('service.wallet.address-without-seed'));
    }

    return this.cipherProvider.generateAddress(wallet.nextSeed).pipe(
      map((response: GenerateAddressResponse) => {
        wallet.nextSeed = response.nextSeed;
        wallet.addresses.push(response.address);
        if (saveWallet) {
          this.saveWallets();
        }
      }));
  }

  private addServerAddress(wallet: Wallet, accountIndex?: number, chainIndex?: number): Observable<void> {
    const params: any = {
      id: wallet.filename,
      num: '1',
    };
    if (accountIndex !== undefined) {
      params.account = accountIndex.toString();
    }
    if (chainIndex !== undefined) {
      params.chain = chainIndex === 1 ? 'change' : 'external';
    }

    return this.apiService.post('wallet/newAddress', params).pipe(
      mergeMap(() => {
        // Reload wallet from server to get updated addresses and accounts
        return this.apiService.get('wallet', { id: wallet.filename }).pipe(
          map((response: any) => {
            wallet.addresses = (response.entries || []).map(e => ({ address: e.address }));
            if (response.meta.type === 'bip44' && response.accounts) {
              wallet.accounts = response.accounts.map(a => ({
                name: a.name,
                index: a.index,
                externalAddresses: (a.external_entries || []).map(e => ({ address: e.address })),
                changeAddresses: (a.change_entries || []).map(e => ({ address: e.address })),
              }));
            }
            this.saveWallets();
          }));
      }));
  }

  create(label: string, seed: string, coinId: number, save = true, walletType = 'deterministic', seedPassphrase?: string, segwit?: boolean): Observable<Wallet> {
    seed = this.getCleanSeed(seed);

    // When the backend has wallet management, create via the API so wallets persist to disk
    if (this.currentCoin && this.currentCoin.serverWallets) {
      return this.createServerWallet(label, seed, coinId, save, walletType, seedPassphrase, segwit);
    }

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

        if ((this.wallets.value || []).some((wlt: Wallet) =>
            wlt.addresses[0].address === wallet.addresses[0].address &&
            wlt.coinId === wallet.coinId)) {
          throw new Error(this.translate.instant('service.wallet.wallet-exists'));
        }

        if (save) {
          this.add(wallet);
        }

        return wallet;
      }));
  }

  private createServerWallet(label: string, seed: string, coinId: number, save: boolean, walletType: string, seedPassphrase?: string, segwit?: boolean): Observable<Wallet> {
    const params: any = {
      label: label || 'undefined',
      seed: seed,
      type: walletType,
    };

    if (walletType === 'bip44') {
      params.scan = 100;
    }

    if (seedPassphrase) {
      params['seed-passphrase'] = seedPassphrase;
    }

    // For Bitcoin BIP44 wallets, signal segwit preference to the backend
    if (segwit !== undefined) {
      params.segwit = segwit;
    }

    return this.apiService.post('wallet/create', params).pipe(
      map((response: any) => {
        const wallet: Wallet = {
          label: response.meta.label,
          balance: new BigNumber('0'),
          hours: new BigNumber('0'),
          addresses: (response.entries || []).map(e => ({ address: e.address })),
          coinId: coinId,
          walletType: response.meta.type || walletType,
        };

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
    if (wallet.addresses.length !== 1) {
      throw new Error(this.translate.instant('service.wallet.invalid-wallet'));
    }

    const InitialNextSeed = wallet.nextSeed;

    return this.checkWalletAddresses(wallet, 0, onProgressChanged, InitialNextSeed).pipe(
      map(lastIndexWithTxs => {
        const unnecessaryAddresses = wallet.addresses.length - 1 - lastIndexWithTxs;
        if (unnecessaryAddresses > 0) {
          wallet.addresses.splice(lastIndexWithTxs + 1, unnecessaryAddresses);
        }
        this.saveWallets();
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
    seed = this.getCleanSeed(seed);
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
        this.encryptionService.encrypt(wallet.seed, password),
        wallet.nextSeed ? this.encryptionService.encrypt(wallet.nextSeed, password) : Promise.resolve(null),
      ]).then(([encSeed, encNextSeed]) => {
        wallet.encryptedSeed = encSeed;
        if (encNextSeed) {
          wallet.encryptedNextSeed = encNextSeed;
        }
        this.saveWallets();
        observer.next();
        observer.complete();
      }).catch(err => observer.error(err));
    });
  }

  unlockWalletWithPassword(wallet: Wallet, password: string, onProgressChanged: EventEmitter<number>): Observable<void> {
    if (!wallet.encryptedSeed) {
      return throwError(() => new Error('Wallet does not have an encrypted seed'));
    }

    return new Observable<void>(observer => {
      const decryptPromises = [
        this.encryptionService.decrypt(wallet.encryptedSeed, password),
        wallet.encryptedNextSeed ? this.encryptionService.decrypt(wallet.encryptedNextSeed, password) : Promise.resolve(null),
      ];

      Promise.all(decryptPromises).then(([seed, nextSeed]) => {
        const cleanSeed = this.getCleanSeed(seed);
        const currentSeed = convertAsciiToHexa(cleanSeed);

        this.unlockWalletAddresses(currentSeed, wallet, 0, onProgressChanged).subscribe(
          (res: boolean) => {
            if (!res) {
              observer.error(new Error(this.translate.instant('service.wallet.wrong-seed')));
              return;
            }
            wallet.seed = cleanSeed;
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
    this.saveWallets();
  }

  saveWallets() {
    const currentWallets = this.wallets.value || [];
    // Only persist to localStorage for client-side wallets
    if (!(this.currentCoin && this.currentCoin.serverWallets)) {
      const strippedWallets: Wallet[] = [];
      currentWallets.forEach(wallet => {
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
        }
        if (wallet.encryptedNextSeed) {
          stripped.encryptedNextSeed = wallet.encryptedNextSeed;
        }
        strippedWallets.push(stripped);
      });
      localStorage.setItem('wallets', JSON.stringify(strippedWallets));
    }

    this.wallets.next(currentWallets);
  }

  private loadWallets() {
    // Wait for coins to load so we know whether to use server or client-side wallets.
    // This prevents the wallet wizard from opening before the current coin is known.
    this.coinService.coinsLoaded.pipe(first()).subscribe(() => {
      const coin = this.coinService.currentCoin.getValue();

      if (coin && coin.serverWallets) {
        this.loadWalletsFromServer();
      } else {
        this.loadWalletsFromLocalStorage();
      }

      // Re-fetch wallets when the user switches coins
      this.coinSubscription = this.coinService.currentCoin.pipe(
        filter((c: BaseCoin) => c !== null),
      ).subscribe((c) => {
        if (c.serverWallets) {
          this.loadWalletsFromServer();
        } else {
          this.loadWalletsFromLocalStorage();
        }
      });
    });
  }

  private loadWalletsFromLocalStorage() {
    const storedWallets: string = localStorage.getItem('wallets');
    if (storedWallets) {
      const wallets: Wallet[] = JSON.parse(storedWallets);
      wallets.filter(wallet => !wallet.coinId).forEach((wallet) => {
        wallet.coinId = defaultCoinId;
      });
      this.wallets.next(wallets);
    } else {
      this.wallets.next([]);
    }
  }

  private loadWalletsFromServer() {
    this.apiService.get('wallets').subscribe(
      (serverWallets: any[]) => {
        if (serverWallets && serverWallets.length > 0) {
          const wallets: Wallet[] = serverWallets.map(w => {
            const wlt: Wallet = {
              label: w.meta.label,
              addresses: (w.entries || []).map(e => ({ address: e.address })),
              coinId: this.currentCoin ? this.currentCoin.id : defaultCoinId,
              encrypted: w.meta.encrypted,
              walletType: w.meta.type || 'deterministic',
              filename: w.meta.filename,
            };

            // Parse BIP44 account structure
            if (w.meta.type === 'bip44' && w.accounts) {
              wlt.accounts = w.accounts.map(a => ({
                name: a.name,
                index: a.index,
                externalAddresses: (a.external_entries || []).map(e => ({ address: e.address })),
                changeAddresses: (a.change_entries || []).map(e => ({ address: e.address })),
              }));
            }

            return wlt;
          });
          this.wallets.next(wallets);
        } else {
          this.wallets.next([]);
        }
      },
      () => this.wallets.next([])
    );
  }

  getXPubKey(wallet: Wallet, path: string): Observable<string> {
    return this.apiService.get('wallet/xpub', {
      id: wallet.filename,
      path: path,
    }).pipe(map((response: any) => response.xpub_key));
  }

  addAccount(wallet: Wallet, name: string): Observable<Wallet> {
    return this.apiService.post('wallet/newAccount', {
      id: wallet.filename,
      name: name,
    }).pipe(map((response: any) => {
      // Update wallet with new account structure
      if (response.accounts) {
        wallet.accounts = response.accounts.map(a => ({
          name: a.name,
          index: a.index,
          externalAddresses: (a.external_entries || []).map(e => ({ address: e.address })),
          changeAddresses: (a.change_entries || []).map(e => ({ address: e.address })),
        }));
      }
      wallet.addresses = (response.entries || []).map(e => ({ address: e.address }));
      this.saveWallets();
      return wallet;
    }));
  }

  encryptWallet(wallet: Wallet, password: string): Observable<void> {
    return this.apiService.post('wallet/encrypt', {
      id: wallet.filename,
      password: password,
    }).pipe(map(() => {}));
  }

  decryptWallet(wallet: Wallet, password: string): Observable<void> {
    return this.apiService.post('wallet/decrypt', {
      id: wallet.filename,
      password: password,
    }).pipe(map(() => {}));
  }

  private getCleanSeed(seed: string): string {
    return seed.replace(/(\n|\r\n)$/, '');
  }

  private checkWalletAddresses(wallet: Wallet, lastIndexWithTxs: number, onProgressChanged: EventEmitter<ScanProgressData>, nextSeed: string): Observable<number> {
    const minAdrressesToScan = environment.e2eTest ? 2 : 10;
    const maxAdrressesToScan = environment.e2eTest ? 2 : 100;

    return this.addAddress(wallet, false).pipe(
      mergeMap(() => this.apiService.get('transactions', { addrs: wallet.addresses[wallet.addresses.length - 1].address })),
      mergeMap(transactions => {
        if (transactions && transactions.length > 0) {
          lastIndexWithTxs = wallet.addresses.length - 1;
          nextSeed = wallet.nextSeed;
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
