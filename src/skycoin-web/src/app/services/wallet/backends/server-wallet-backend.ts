import { Observable, of } from 'rxjs';
import { catchError, map, mergeMap } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { Bip44Account, Wallet } from '../../../app.datatypes';
import { defaultCoinId } from '../../../constants/coins-id.const';
import { ApiService } from '../../api.service';
import { WalletBackend } from './wallet-backend';

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
    seedPassphrase?: string, segwit?: boolean, password?: string): Observable<Wallet> {
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
    if (password) {
      params.encrypt = true;
      params.password = password;
    }

    return this.apiService.post('wallet/create', params).pipe(
      map((response: any) => ({
        ...toWallet(response, coinId),
        balance: new BigNumber('0'),
        hours: new BigNumber('0'),
      })));
  }

  addAddress(wallet: Wallet, accountIndex?: number, chainIndex?: number): Observable<void> {
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
      // Reload the wallet to get its updated addresses and accounts.
      mergeMap(() => this.apiService.get('wallet', { id: wallet.filename } as any)),
      map((response: any) => {
        wallet.addresses = (response.entries || []).map((e: any) => ({ address: e.address }));
        if (response.meta.type === 'bip44' && response.accounts) {
          wallet.accounts = parseAccounts(response.accounts);
        }
      }));
  }

  save(_wallets: Wallet[]): void { }

  getSeed(wallet: Wallet, password: string | null): Observable<string> {
    return this.apiService.post('wallet/seed', { id: wallet.filename, password: password || '' }).pipe(
      map((response: any) => response.seed));
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
  if (w.meta.type === 'bip44' && w.accounts) {
    wlt.accounts = parseAccounts(w.accounts);
  }

  return wlt;
}

function parseAccounts(accounts: any[]): Bip44Account[] {
  return accounts.map((a: any) => ({
    name: a.name,
    index: a.index,
    externalAddresses: (a.external_entries || []).map((e: any) => ({ address: e.address })),
    changeAddresses: (a.change_entries || []).map((e: any) => ({ address: e.address })),
  }));
}
