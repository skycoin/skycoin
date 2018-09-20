import { Component, OnDestroy, OnInit } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import * as moment from 'moment';
import { ISubscription } from 'rxjs/Subscription';
import { NavBarService } from '../../../../services/nav-bar.service';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { BigNumber } from 'bignumber.js';

@Component({
  selector: 'app-pending-transactions',
  templateUrl: './pending-transactions.component.html',
  styleUrls: ['./pending-transactions.component.scss'],
})
export class PendingTransactionsComponent implements OnInit, OnDestroy {
  transactions = null;

  private transactionsSubscription: ISubscription;
  private navbarSubscription: ISubscription;

  constructor(
    public walletService: WalletService,
    private navbarService: NavBarService,
  ) {
    this.navbarSubscription = this.navbarService.activeComponent.subscribe(value => {
      this.loadTransactions(value);
    });
  }

  ngOnInit() {
    this.navbarService.showSwitch('pending-txs.my', 'pending-txs.all');
  }

  ngOnDestroy() {
    this.transactionsSubscription.unsubscribe();
    this.navbarSubscription.unsubscribe();
    this.navbarService.hideSwitch();
  }

  private loadTransactions(value) {
    const method = value === DoubleButtonActive.LeftButton ? 'pendingTransactions' : 'allPendingTransactions';

    this.transactions = null;

    if (this.transactionsSubscription) {
      this.transactionsSubscription.unsubscribe();
    }

    if (method === 'pendingTransactions') {
      this.walletService.startDataRefreshSubscription();
    }

    this.transactionsSubscription = this.walletService[method]().subscribe(transactions => {
      this.transactions = this.mapTransactions(transactions);
    });
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
}
