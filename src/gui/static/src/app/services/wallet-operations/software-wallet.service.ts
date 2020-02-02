import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Injectable } from '@angular/core';

import { ApiService } from '../api.service';
import { WalletsAndAddressesService } from './wallets-and-addresses.service';
import { WalletBase } from './wallet-objects';

/**
 * Allows to perform operations related to a software wallet.
 */
@Injectable()
export class SoftwareWalletService {
  constructor(
    private walletsAndAddressesService: WalletsAndAddressesService,
    private apiService: ApiService,
  ) { }

  /**
   * Allows to change the name or label which identifies a wallet.
   * @param wallet Wallet to modify.
   * @param label New name or label.
   * @returns The returned observable returns nothing, but it can fail in case of error.
   */
  renameWallet(wallet: WalletBase, label: string): Observable<void> {
    return this.apiService.post('wallet/update', { id: wallet.id, label: label }).pipe(map(() => {
      wallet.label = label;
      this.walletsAndAddressesService.informValuesUpdated(wallet);
    }));
  }

  /**
   * Makes an encrypted wallet to be unencrypted, or an unencrypted wallet to be encrypted.
   * @param wallet Wallet to modify.
   * @param password If the wallet is encrypted, the password of the wallet, to be able to
   * disable the encryptation. If the wallet is unencrypted, the password that will be used
   * for encrypting it.
   * @returns The returned observable returns nothing, but it can fail in case of error.
   */
  toggleEncryption(wallet: WalletBase, password: string): Observable<void> {
    return this.apiService.post('wallet/' + (wallet.encrypted ? 'decrypt' : 'encrypt'), { id: wallet.id, password }).pipe(map(w => {
      wallet.encrypted = w.meta.encrypted;
      this.walletsAndAddressesService.informValuesUpdated(wallet);
    }));
  }

  /**
   * Removes or changes the password of an encrypted wallet.
   * @param wallet Wallet to modify.
   * @param seed Seed of the wallet.
   * @param password New password for the wallet. If empty or null, the wallet will be
   * unencrypted after finishing the operation.
   */
  resetPassword(wallet: WalletBase, seed: string, password: string): Observable<void> {
    const params = new Object();
    params['id'] = wallet.id;
    params['seed'] = seed;
    if (password) {
      params['password'] = password;
    }

    return this.apiService.post('wallet/recover', params, {useV2: true}).pipe(map(w => {
      wallet.encrypted = w.data.meta.encrypted;
      this.walletsAndAddressesService.informValuesUpdated(wallet);
    }));
  }

  /**
   * Gets the seed of an encrypted wallet.
   * @param wallet Wallet to get the seed from.
   * @param password Wallet password.
   */
  getWalletSeed(wallet: WalletBase, password: string): Observable<string> {
    return this.apiService.post('wallet/seed', { id: wallet.id, password }).pipe(map(response => response.seed));
  }
}
