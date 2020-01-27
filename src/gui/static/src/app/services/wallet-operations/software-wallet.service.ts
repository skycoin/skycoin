import { Observable } from 'rxjs';
import { tap, map } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { ApiService } from '../api.service';
import { WalletsAndAddressesService } from './wallets-and-addresses.service';
import { WalletBase } from './wallet-objects';

@Injectable()
export class SoftwareWalletService {
  constructor(
    private walletsAndAddressesService: WalletsAndAddressesService,
    private apiService: ApiService,
  ) { }

  renameWallet(wallet: WalletBase, label: string): Observable<WalletBase> {
    return this.apiService.post('wallet/update', { id: wallet.id, label: label }).pipe(tap(() => {
      wallet.label = label;
      this.walletsAndAddressesService.informValuesUpdated(wallet);
    }));
  }

  toggleEncryption(wallet: WalletBase, password: string): Observable<WalletBase> {
    return this.apiService.post('wallet/' + (wallet.encrypted ? 'decrypt' : 'encrypt'), { id: wallet.id, password }).pipe(tap(w => {
      wallet.encrypted = w.meta.encrypted;
      this.walletsAndAddressesService.informValuesUpdated(wallet);
    }));
  }

  resetPassword(wallet: WalletBase, seed: string, password: string): Observable<WalletBase> {
    const params = new Object();
    params['id'] = wallet.id;
    params['seed'] = seed;
    if (password) {
      params['password'] = password;
    }

    return this.apiService.post('wallet/recover', params, {}, true).pipe(tap(w => {
      wallet.encrypted = w.data.meta.encrypted;
      this.walletsAndAddressesService.informValuesUpdated(wallet);
    }));
  }

  getWalletSeed(wallet: WalletBase, password: string): Observable<string> {
    return this.apiService.post('wallet/seed', { id: wallet.id, password }).pipe(map(response => response.seed));
  }
}
