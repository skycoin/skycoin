import { Component, OnDestroy, OnInit } from '@angular/core';
import * as moment from 'moment';
import { Subscription } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { WalletService } from '../../../../services/wallet/wallet.service';
import { HistoryService } from '../../../../services/wallet/history.service';
import { NavBarService } from '../../../../services/nav-bar.service';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { Wallet } from '../../../../app.datatypes';
import { Observable } from 'rxjs';
import { BaseCoin } from '../../../../coins/basecoin';
import { CoinService } from '../../../../services/coin.service';
import { GlobalsService } from '../../../../services/globals.service';
import { isEqualOrSuperiorVersion } from '../../../../utils/semver';

@Component({
  selector: 'app-pending-transactions',
  templateUrl: './pending-transactions.component.html',
  styleUrls: ['./pending-transactions.component.scss']
})
export class PendingTransactionsComponent implements OnInit, OnDestroy {
  isLoading = false;
  transactions: any[] = [];
  currentCoin: BaseCoin;
  showError = false;

  private navbarSubscription: Subscription;
  private coinSubscription: Subscription;
  private dataSubscription: Subscription;

  constructor(
    private walletService: WalletService,
    private historyService: HistoryService,
    private navbarService: NavBarService,
    private coinService: CoinService,
    private globalsService: GlobalsService
  ) { }

  ngOnInit() {
    this.coinSubscription = this.coinService.currentCoin
      .subscribe((coin: BaseCoin) => {
        this.currentCoin = coin;
        this.navbarService.setActiveComponent();
      });

    this.navbarService.showSwitch('pending-txs.my', 'pending-txs.all', DoubleButtonActive.LeftButton);

    this.navbarSubscription = this.navbarService.activeComponent.subscribe(value => {
      this.loadTransactions(value);
    });
  }

  ngOnDestroy() {
    this.navbarSubscription.unsubscribe();
    this.coinSubscription.unsubscribe();

    this.closeDataSubscription();
    this.navbarService.hideSwitch();
  }

  private loadTransactions(value: number) {
    this.isLoading = true;
    this.transactions = [];
    this.showError = false;

    const showAllTransactions = value === DoubleButtonActive.RightButton;
    this.closeDataSubscription();
    this.dataSubscription = this.historyService.getAllPendingTransactions()
      .delay(32)
      .flatMap((transactions: any) => {
        return showAllTransactions ? Observable.of(transactions) : this.getWalletsTransactions(transactions);
      })
      .subscribe(transactions => {
        this.transactions = this.mapTransactions(transactions);
        this.isLoading = false;
      },
      () => this.showError = true);
  }

  private mapTransactions(transactions) {
    return transactions.map(transaction => {
      transaction.transaction.timestamp = moment(transaction.received).unix();
      return transaction.transaction;
    })
      .map(transaction => {
        transaction.amount = new BigNumber('0');
        transaction.hours = new BigNumber('0');
        transaction.outputs.map(output => {
          transaction.amount = transaction.amount.plus(output.coins);
          transaction.hours = transaction.hours.plus(output.hours);
        });

        return transaction;
      });
  }

  private getWalletsTransactions(transactions: any): Observable<any> {
    if (transactions.length === 0) {
      return Observable.of([]);
    }

    return this.globalsService.getValidNodeVersion().flatMap (version => {
      let allTransactions: Observable<any>;
      if (isEqualOrSuperiorVersion(version, '0.25.0')) {
        allTransactions = Observable.of(transactions);
      } else {
        allTransactions = this.getUpdatedTransactions(transactions);
      }

      return Observable.forkJoin(allTransactions, this.walletService.currentWallets.first(), (trans: any, wallets: Wallet[]) => {
        const walletAddresses = new Set<string>();
        wallets.forEach(wallet => {
          wallet.addresses.forEach(address => walletAddresses.add(address.address));
        });

        return trans.filter(tran => {
          if (isEqualOrSuperiorVersion(version, '0.25.0')) {
            return tran.transaction.inputs.some(input => walletAddresses.has(input.owner)) ||
            tran.transaction.outputs.some(output => walletAddresses.has(output.dst));
          } else {
            return tran.owner_addressses.some(address => walletAddresses.has(address)) ||
            tran.transaction.outputs.some(output => walletAddresses.has(output.dst));
          }
        });
      });
    });
  }

  private getUpdatedTransactions(transactions: any): Observable<any> {
    return Observable.forkJoin(transactions.map((transaction: any) => {
      return Observable.forkJoin(transaction.transaction.inputs
        .map(input => this.historyService.getTransactionDetails(input)
          .map(inputDetails => inputDetails.owner_address)))
        .map((addresses) => {
          transaction.owner_addressses = addresses;
          return transaction;
        });
    }));
  }

  private closeDataSubscription() {
    if (this.dataSubscription && !this.dataSubscription.closed) {
      this.dataSubscription.unsubscribe();
    }
  }
}
