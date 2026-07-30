import { Injectable, NgZone } from '@angular/core';
import { BehaviorSubject, Observable, Subscription, ReplaySubject, of, forkJoin } from 'rxjs';
import { mergeMap, map, first, delay } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { ApiService } from '../api.service';
import { CoinService } from '../coin.service';
import { Address, Wallet, TotalBalance, Balance } from '../../app.datatypes';
import { WalletService } from './wallet.service';
import { isEqualOrSuperiorVersion } from '../../utils/semver';
import { GlobalsService } from '../globals.service';
import { BaseCoin } from '../../coins/basecoin';

export enum BalanceStates {
  Obtained,
  Error,
  Updating,
}

export class BalanceEvent {
  state!: BalanceStates;
  balance?: TotalBalance;
}

@Injectable()
export class BalanceService {
  lastBalancesUpdateTime: Date = new Date();
  totalBalance: ReplaySubject<BalanceEvent> = new ReplaySubject<BalanceEvent>();
  hasPendingTransactions: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  private canGetBalance = false;
  private schedulerSubscription!: Subscription;
  private currentCoin!: BaseCoin;

  private readonly shortUpdatePeriod = 10 * 1000;
  private readonly longUpdatePeriod = 300 * 1000;

  private get coinsMultiplier(): number {
    return this.currentCoin ? this.currentCoin.coinsMultiplier : 1000000;
  }

  constructor(
    private apiService: ApiService,
    private walletService: WalletService,
    private globalsService: GlobalsService,
    private coinService: CoinService,
    private _ngZone: NgZone
  ) {
    coinService.currentCoin.subscribe(coin => this.currentCoin = coin);
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

  private scheduleUpdate(delayMs: number) {
    this._ngZone.runOutsideAngular(() => {
      this.removeSubscription();

      this.schedulerSubscription = of(1).pipe(
        delay(delayMs),
        mergeMap(() => this.getBalance()),
      )
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
    return this.walletService.addresses.pipe(first(), mergeMap((addresses: Address[]) => {
      if (addresses.length === 0) {
        this.lastBalancesUpdateTime = new Date();
        this.sendTotalBalanceEvent({state: BalanceStates.Obtained, balance: { coins: new BigNumber('0'), hours: new BigNumber('0') }});
        return of(false);
      }

      return this.retrieveAddressesBalance(addresses).pipe(mergeMap((balance) => {
        return this.walletService.currentWallets.pipe(first(), map(wallets => this.calculateBalance(wallets, balance)));
      }));
    }));
  }

  private retrieveAddressesBalance(addresses: Address[]): Observable<Balance> {
    const formattedAddresses = addresses.map(a => a.address).join(',');

    if (this.currentCoin && this.currentCoin.isBitcoin()) {
      return this.apiService.get('btc/balance', { addrs: formattedAddresses } as any);
    }

    return this.globalsService.getValidNodeVersion().pipe(mergeMap(version => {
      if (isEqualOrSuperiorVersion(version, '0.25.0')) {
        return this.apiService.post('balance', { addrs: formattedAddresses });
      } else {
        return this.chunkedBalanceGet(addresses);
      }
    }));
  }

  // Splits addresses into chunks for GET requests to avoid URI length limits on older nodes.
  private chunkedBalanceGet(addresses: Address[]): Observable<Balance> {
    const chunks = this.chunkAddresses(addresses);

    if (chunks.length === 1) {
      return this.apiService.get('balance', { addrs: chunks[0] } as any);
    }

    return forkJoin(chunks.map(chunk => this.apiService.get('balance', { addrs: chunk } as any))).pipe(
      map((results: Balance[]) => {
        const merged: Balance = {
          confirmed: { coins: 0, hours: 0 },
          predicted: { coins: 0, hours: 0 },
          addresses: {},
        };
        results.forEach(r => {
          merged.confirmed.coins += r.confirmed.coins;
          merged.confirmed.hours += r.confirmed.hours;
          merged.predicted.coins += r.predicted.coins;
          merged.predicted.hours += r.predicted.hours;
          if (r.addresses) {
            Object.assign(merged.addresses, r.addresses);
          }
        });
        return merged;
      })
    );
  }

  private chunkAddresses(addresses: Address[], maxChars = 1800): string[] {
    const chunks: string[] = [];
    let current = '';
    addresses.forEach(a => {
      const addr = a.address;
      if (current.length > 0 && current.length + 1 + addr.length > maxChars) {
        chunks.push(current);
        current = addr;
      } else {
        current = current ? current + ',' + addr : addr;
      }
    });
    if (current) {
      chunks.push(current);
    }
    return chunks;
  }

  private calculateBalance(wallets: Wallet[], balance: Balance): boolean {
    const isBtc = this.currentCoin && this.currentCoin.isBitcoin();

    if (balance.addresses) {
      wallets.map((wallet: Wallet) => {
        wallet.balance = new BigNumber('0');
        wallet.hours = new BigNumber('0');

        wallet.addresses.map((address: Address) => {
          if (balance.addresses[address.address]) {
            address.balance = new BigNumber(balance.addresses[address.address].confirmed.coins).dividedBy(this.coinsMultiplier);
            address.hours = isBtc ? new BigNumber('0') : new BigNumber(balance.addresses[address.address].confirmed.hours || 0);
            wallet.balance = wallet.balance!.plus(address.balance);
            wallet.hours = wallet.hours!.plus(address.hours);
          }
        });
      });
    }

    this.lastBalancesUpdateTime = new Date();
    this.sendTotalBalanceEvent({
      state: BalanceStates.Obtained,
      balance: {
        coins: new BigNumber(balance.confirmed.coins).dividedBy(this.coinsMultiplier),
        hours: isBtc ? new BigNumber('0') : new BigNumber(balance.confirmed.hours || 0),
      }
    });
    return this.refreshPendingTransactions(balance);
  }

  private refreshPendingTransactions(balance: Balance) {
    const hasPendingTxs = balance.confirmed.coins !== balance.predicted.coins ||
      (!this.currentCoin?.isBitcoin() && balance.confirmed.hours !== balance.predicted.hours);

    this.hasPendingTransactions.next(hasPendingTxs);
    return hasPendingTxs;
  }
}
