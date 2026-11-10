import { Observable, of } from 'rxjs';
import { catchError, map, mergeMap } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { Address, Bip44Account, Wallet } from '../../../app.datatypes';
import { defaultCoinId } from '../../../constants/coins-id.const';
import { ApiService } from '../../api.service';
import { SeedBackup, WalletBackend } from './wallet-backend';

/** Wallets kept by a skycoin wallet service, through its wallet API. */
export class ServerWalletBackend implements WalletBackend {
  readonly onServer = true;

  constructor(private apiService: ApiService, private currentCoinId: () => number | undefined) { }

  load(): Observable<Wallet[]> {
    return this.apiService.get('wallets').pipe(
      map((serverWallets: any[]) => (serverWallets || []).map(w => toWallet(w, this.currentCoinId() ?? defaultCoinId))),
      catchError(() => of([])));
  }

  create(label: string, seed: string, coinId: number, _existing: Wallet[], walletType: string,
    seedPassphrase?: string, segwit?: boolean, password?: string, temporary = false): Observable<Wallet> {
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
    if (password && !temporary) {
      params.encrypt = true;
      params.password = password;
    }

    return this.apiService.post(temporary ? 'wallet/createTemp' : 'wallet/create', params).pipe(
      map((response: any) => ({
        ...toWallet(response, coinId),
        balance: new BigNumber('0'),
        hours: new BigNumber('0'),
      })));
  }

  addAddress(wallet: Wallet, accountIndex?: number, chainIndex?: number, count = 1): Observable<void> {
    const params: any = {
      id: wallet.filename,
      num: count.toString(),
    };
    if (accountIndex !== undefined) {
      params.account = accountIndex.toString();
    }
    if (chainIndex !== undefined) {
      params.chain = chainIndex === 1 ? 'change' : 'external';
    }

    return this.apiService.post('wallet/newAddress', params).pipe(mergeMap(() => this.reload(wallet)));
  }

  /**
   * Asks the wallet service to add the addresses used past the wallet's last
   * ones, on both chains of a bip44 wallet, until gap in a row are unused.
   * Emits how many it added.
   */
  scanAddresses(wallet: Wallet, gap: number): Observable<number> {
    const before = addressCount(wallet);

    return this.apiService.post('wallet/scan', { id: wallet.filename, num: gap.toString() }).pipe(
      // The reply lists only external addresses, so count from the wallet.
      mergeMap(() => this.reload(wallet)),
      map(() => addressCount(wallet) - before));
  }

  /** Reads the wallet's addresses and accounts again. */
  private reload(wallet: Wallet): Observable<void> {
    return this.apiService.get('wallet', { id: wallet.filename } as any).pipe(
      map((response: any) => {
        // The new address objects keep the balances already read, until the next refresh.
        const known = new Map<string, Address>();
        [wallet.addresses, ...(wallet.accounts || []).map(a => [...a.externalAddresses, ...a.changeAddresses])]
          .forEach(list => (list || []).forEach(a => known.set(a.address, a)));
        const keep = (a: Address) => {
          const old = known.get(a.address);
          if (old) {
            a.balance = old.balance;
            a.hours = old.hours;
          }
        };

        wallet.addresses = (response.entries || []).map((e: any) => ({ address: e.address }));
        wallet.addresses.forEach(keep);
        if (response.meta.type === 'bip44' && response.accounts) {
          wallet.accounts = parseAccounts(response.accounts);
          wallet.accounts.forEach(a => [...a.externalAddresses, ...a.changeAddresses].forEach(keep));
        }
      }));
  }

  save(_wallets: Wallet[]): void { }

  getSeed(wallet: Wallet, password: string | null): Observable<SeedBackup> {
    return this.apiService.post('wallet/seed', { id: wallet.filename, password: password || '' }).pipe(
      map((response: any) => ({ seed: response.seed, passphrase: response.seed_passphrase || undefined })));
  }

  /** The folder the wallet service keeps its wallet files in. */
  folder(): Observable<string> {
    return this.apiService.get('wallets/folderName').pipe(map((response: any) => response.address));
  }

  getXPubKey(wallet: Wallet, path: string): Observable<XPubKey> {
    return this.apiService.get('wallet/xpub', {
      id: wallet.filename,
      path: path,
    } as any).pipe(map((response: any) => ({ key: response.xpub_key, addresses: response.addresses || [] })));
  }

  /** Adds a BIP44 account. An encrypted wallet needs its password, as the keys come from the seed. */
  addAccount(wallet: Wallet, name: string, password?: string): Observable<Wallet> {
    const params: any = { id: wallet.filename, name: name };
    if (password) {
      params.password = password;
    }

    return this.apiService.post('wallet/newAccount', params).pipe(map((response: any) => {
      if (response.accounts) {
        wallet.accounts = parseAccounts(response.accounts);
      }
      wallet.addresses = (response.entries || []).map((e: any) => ({ address: e.address }));

      return wallet;
    }));
  }

  encrypt(wallet: Wallet, password: string): Observable<void> {
    return this.apiService.post('wallet/encrypt', {
      id: wallet.filename,
      password: password,
    }).pipe(map(() => {}));
  }

  /** Sets or removes the password of an encrypted wallet with its seed. Emits whether it is encrypted after. */
  recover(params: { id?: string, seed: string, seed_passphrase?: string, password?: string }): Observable<boolean> {
    return this.apiService.post('wallet/recover', params, {}, true).pipe(
      map((response: any) => !!(response && response.data ? response.data : response).meta.encrypted));
  }

  decrypt(wallet: Wallet, password: string): Observable<void> {
    return this.apiService.post('wallet/decrypt', {
      id: wallet.filename,
      password: password,
    }).pipe(map(() => {}));
  }
}

/** An xpub and the first addresses it derives. */
export interface XPubKey {
  key: string;
  addresses: string[];
}

function toWallet(w: any, coinId: number): Wallet {
  const wlt: Wallet = {
    label: w.meta.label,
    addresses: (w.entries || []).map((e: any) => ({ address: e.address })),
    coinId: coinId,
    encrypted: w.meta.encrypted,
    walletType: w.meta.type || 'deterministic',
    filename: w.meta.filename,
  };
  if (w.meta.temp) {
    wlt.temporary = true;
  }
  if (w.meta.type === 'bip44' && w.accounts) {
    wlt.accounts = parseAccounts(w.accounts);
  }

  return wlt;
}

/** The addresses of every account of a bip44 wallet, else the wallet's addresses. */
function addressCount(wallet: Wallet): number {
  if (wallet.accounts && wallet.accounts.length > 0) {
    return wallet.accounts.reduce((n, a) => n + a.externalAddresses.length + a.changeAddresses.length, 0);
  }

  return wallet.addresses.length;
}

function parseAccounts(accounts: any[]): Bip44Account[] {
  return accounts.map((a: any) => ({
    name: a.name,
    index: a.index,
    externalAddresses: (a.external_entries || []).map((e: any) => ({ address: e.address })),
    changeAddresses: (a.change_entries || []).map((e: any) => ({ address: e.address })),
  }));
}
