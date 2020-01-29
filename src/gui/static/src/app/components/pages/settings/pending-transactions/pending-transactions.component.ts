import { Component, OnDestroy, OnInit } from '@angular/core';
import * as moment from 'moment';
import { SubscriptionLike } from 'rxjs';
import { NavBarService } from '../../../../services/nav-bar.service';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { BigNumber } from 'bignumber.js';
import { BalanceAndOutputsService } from '../../../../services/wallet-operations/balance-and-outputs.service';
import { HistoryService } from '../../../../services/wallet-operations/history.service';

@Component({
  selector: 'app-pending-transactions',
  templateUrl: './pending-transactions.component.html',
  styleUrls: ['./pending-transactions.component.scss'],
})
export class PendingTransactionsComponent implements OnInit, OnDestroy {
  transactions = null;

  private transactionsSubscription: SubscriptionLike;
  private navbarSubscription: SubscriptionLike;

  constructor(
    private navbarService: NavBarService,
    private balanceAndOutputsService: BalanceAndOutputsService,
    private historyService: HistoryService,
  ) {
    this.navbarSubscription = this.navbarService.activeComponent.subscribe(value => {
      this.startCheckingTransactions(value);
    });
  }

  ngOnInit() {
    this.navbarService.showSwitch('pending-txs.my-transactions-button', 'pending-txs.all-transactions-button');
  }

  ngOnDestroy() {
    this.removeTransactionsSubscription();
    this.navbarSubscription.unsubscribe();
    this.navbarService.hideSwitch();
  }

  private startCheckingTransactions(value) {
    this.transactions = null;

    this.removeTransactionsSubscription();

    // Currently gets the data only one time.
    this.transactionsSubscription = this.historyService.getPendingTransactions().subscribe(transactions => {
      this.transactions = this.mapTransactions(value === DoubleButtonActive.LeftButton ? transactions.user : transactions.all);
    });

    // Due to some changes, must use a method for updating or getting the pending transactions, not this.
    this.balanceAndOutputsService.refreshBalance();
  }

  private mapTransactions(transactions) {
    return transactions.map(transaction => {
      transaction.transaction.timestamp = moment(transaction.received).unix();

      return transaction.transaction;
    })
    .map(transaction => {
      let amount = new BigNumber('0');
      transaction.outputs.map(output => amount = amount.plus(output.coins));
      transaction.amount = amount.decimalPlaces(6).toString();

      let hours = new BigNumber('0');
      transaction.outputs.map(output => hours = hours.plus(output.hours));
      transaction.hours = hours.decimalPlaces(0).toString();

      return transaction;
    });
  }

  private removeTransactionsSubscription() {
    if (this.transactionsSubscription) {
      this.transactionsSubscription.unsubscribe();
    }
  }
}
