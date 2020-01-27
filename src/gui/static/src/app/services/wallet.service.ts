import { throwError as observableThrowError, of, Subject, Observable, ReplaySubject, Subscription, BehaviorSubject } from 'rxjs';
import { concat, delay, retryWhen, first, take, mergeMap, catchError, map } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Address, NormalTransaction, PreviewTransaction, Wallet } from '../app.datatypes';
import { BigNumber } from 'bignumber.js';
import { HwWalletService, HwOutput, HwInput } from './hw-wallet.service';
import { TranslateService } from '@ngx-translate/core';
import { StorageService, StorageType } from './storage.service';
import { TxEncoder } from '../utils/tx-encoder';
import { BalanceAndOutputsService } from './wallet-operations/balance-and-outputs.service';

export interface PendingTransactions {
  user: any[];
  all: any[];
}

@Injectable()
export class WalletService {

  addresses: Address[];
  wallets: Subject<Wallet[]> = new ReplaySubject<Wallet[]>(1);
  pendingTxs: Subject<PendingTransactions> = new ReplaySubject<PendingTransactions>(1);
  dataRefreshSubscription: Subscription;

  initialLoadFailed: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  constructor(
    private apiService: ApiService,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
    private storageService: StorageService,
    private balanceAndOutputsService: BalanceAndOutputsService,
  ) {
    this.loadData();
    this.balanceAndOutputsService.refreshBalance();
  }

  addressesAsString(): Observable<string> {
    return this.allAddresses().pipe(map(addrs => addrs.map(addr => addr.address)), map(addrs => addrs.join(',')));
  }

  all(): Observable<Wallet[]> {
    return this.wallets.asObservable();
  }

  allAddresses(): Observable<any[]> {
    return this.all().pipe(map(wallets => wallets.reduce((array, wallet) => array.concat(wallet.addresses), [])));
  }

  saveHardwareWallets() {
    this.wallets.pipe(first()).subscribe(wallets => {
      const hardwareWallets: Wallet[] = [];

      wallets.map(wallet => {
        if (wallet.isHardware) {
          hardwareWallets.push(this.createHardwareWalletData(
            wallet.label,
            wallet.addresses.map(address => {
              return { address: address.address, confirmed: address.confirmed };
            }),
            wallet.hasHwSecurityWarnings,
            wallet.stopShowingHwSecurityPopup,
          ));
        }
      });

      this.hwWalletService.saveWalletsData(JSON.stringify(hardwareWallets));

      this.wallets.next(wallets);
    });
  }

  private createHardwareWalletData(label: string, addresses: {address: string, confirmed: boolean}[], hasHwSecurityWarnings: boolean, stopShowingHwSecurityPopup: boolean): Wallet {
    return {
      label: label,
      filename: '',
      hasHwSecurityWarnings: hasHwSecurityWarnings,
      stopShowingHwSecurityPopup: stopShowingHwSecurityPopup,
      coins: null,
      hours: null,
      addresses: addresses.map(address => {
        return {
          address: address.address,
          coins: null,
          hours: null,
          confirmed: address.confirmed,
        };
      }),
      encrypted: false,
      isHardware: true,
    };
  }

  private loadData(): void {
    let wallets: Wallet[] = [];
    let softwareWallets: Wallet[] = [];

    this.apiService.getWallets().pipe(first(), mergeMap(recoveredWallets => {
      softwareWallets = recoveredWallets;

      if (this.hwWalletService.hwWalletCompatibilityActivated) {
        return this.loadHardwareWallets(wallets);
      }

      return of(null);

    })).subscribe(() => {
      wallets = wallets.concat(softwareWallets);
      this.wallets.next(wallets);
    }, () => this.initialLoadFailed.next(true));
  }

  private loadHardwareWallets(wallets: Wallet[]): Observable<any> {
    return this.hwWalletService.getSavedWalletsData().pipe(map(storedWallets => {
      if (storedWallets) {
        const loadedWallets: Wallet[] = JSON.parse(storedWallets);
        loadedWallets.map(wallet => wallets.push(wallet));
      }

      return null;
    }));
  }
}
