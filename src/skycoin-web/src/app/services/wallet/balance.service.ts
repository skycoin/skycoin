import { Injectable, NgZone } from '@angular/core';
import 'rxjs/add/operator/mergeMap';
import 'rxjs/add/operator/first';
import { BehaviorSubject } from 'rxjs';
import { Observable } from 'rxjs';
import { Subscription } from 'rxjs';
import { ReplaySubject } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { ApiService } from '../api.service';
import { Address, Wallet, TotalBalance, Balance } from '../../app.datatypes';
import { WalletService } from './wallet.service';
import { isEqualOrSuperiorVersion } from '../../utils/semver';
import { GlobalsService } from '../globals.service';

export enum BalanceStates {
  Obtained,
  Error,
  Updating,
}

export class BalanceEvent {
  state: BalanceStates;
  balance?: TotalBalance;
}

@Injectable()
export class BalanceService {
  lastBalancesUpdateTime: Date = new Date();
  totalBalance: ReplaySubject<BalanceEvent> = new ReplaySubject<BalanceEvent>();
  hasPendingTransactions: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  private canGetBalance = false;
  private schedulerSubscription: Subscription;

  private readonly coinsMultiplier = 1000000;
  private readonly shortUpdatePeriod = 10 * 1000;
  private readonly longUpdatePeriod = 300 * 1000;

  constructor(
    private apiService: ApiService,
    private walletService: WalletService,
    private globalsService: GlobalsService,
    private _ngZone: NgZone
  ) {
    walletService.wallets.subscribe(() => this.canGetBalance ? this.scheduleUpdate(0) : null);
  }

  startGettingBalances() {
    this.canGetBalance = true;
    this.scheduleUpdate(0);
  }

  stopGettingBalances() {
    this.removeSubscription();
    this.canGetBalance = false;
  }

  private removeSubscription() {
    if (!!this.schedulerSubscription && !this.schedulerSubscription.closed) {
      this.schedulerSubscription.unsubscribe();
    }
  }

  private scheduleUpdate(delay: number) {
    this._ngZone.runOutsideAngular(() => {
      this.removeSubscription();

      this.schedulerSubscription = Observable.of(1)
        .delay(delay)
        .flatMap(() => this.getBalance())
        .subscribe(
          hasPendingTxs => this.scheduleUpdate(hasPendingTxs ? this.shortUpdatePeriod : this.longUpdatePeriod),
          () => {
            this.scheduleUpdate(this.shortUpdatePeriod);
            this.sendTotalBalanceEvent({state: BalanceStates.Error});
          }
        );
    });
  }

  private sendTotalBalanceEvent(event: BalanceEvent) {
    this._ngZone.run(() => this.totalBalance.next(event));
  }

  private getBalance(): Observable<boolean> {
    this.sendTotalBalanceEvent({state: BalanceStates.Updating});
    return this.walletService.addresses.first().flatMap((addresses: Address[]) => {
      if (addresses.length === 0) {
        this.lastBalancesUpdateTime = new Date();
        this.sendTotalBalanceEvent({state: BalanceStates.Obtained, balance: { coins: new BigNumber('0'), hours: new BigNumber('0') }});
        return Observable.of(false);
      }

      return this.retrieveAddressesBalance(addresses).flatMap((balance) => {
        return this.walletService.currentWallets.first().map(wallets => this.calculateBalance(wallets, balance));
      });
    });
  }

  private retrieveAddressesBalance(addresses: Address[]): Observable<Balance> {
    const formattedAddresses = addresses.map(a => a.address).join(',');

    return this.globalsService.getValidNodeVersion().flatMap (version => {
      if (isEqualOrSuperiorVersion(version, '0.25.0')) {
        return this.apiService.post('balance', { addrs: formattedAddresses });
      } else {
        return this.apiService.get('balance', { addrs: formattedAddresses });
      }
    });
  }

  private calculateBalance(wallets: Wallet[], balance: Balance): boolean {
    if (balance.addresses) {
      wallets.map((wallet: Wallet) => {
        wallet.balance = new BigNumber('0');
        wallet.hours = new BigNumber('0');

        wallet.addresses.map((address: Address) => {
          if (balance.addresses[address.address]) {
            address.balance = new BigNumber(balance.addresses[address.address].confirmed.coins).dividedBy(this.coinsMultiplier);
            address.hours = new BigNumber(balance.addresses[address.address].confirmed.hours);
            wallet.balance = wallet.balance.plus(address.balance);
            wallet.hours = wallet.hours.plus(address.hours);
          }
        });
      });
    }

    this.lastBalancesUpdateTime = new Date();
    this.sendTotalBalanceEvent({
      state: BalanceStates.Obtained,
      balance: { coins: new BigNumber(balance.confirmed.coins).dividedBy(this.coinsMultiplier), hours: new BigNumber(balance.confirmed.hours) }
    });
    return this.refreshPendingTransactions(balance);
  }

  private refreshPendingTransactions(balance: Balance) {
    const hasPendingTxs = balance.confirmed.coins !== balance.predicted.coins ||
      balance.confirmed.hours !== balance.predicted.hours;

    this.hasPendingTransactions.next(hasPendingTxs);
    return hasPendingTxs;
  }
}
