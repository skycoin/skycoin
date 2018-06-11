import { Component, OnDestroy, OnInit } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import * as moment from 'moment';
import { ISubscription } from 'rxjs/Subscription';
import { NavBarService } from '../../../../services/nav-bar.service';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';

@Component({
  selector: 'app-pending-transactions',
  templateUrl: './pending-transactions.component.html',
  styleUrls: ['./pending-transactions.component.scss'],
})
export class PendingTransactionsComponent implements OnInit, OnDestroy {
  transactions: any[] = [];

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
      transaction.amount = transaction.outputs
        .map(output => output.coins >= 0 ? output.coins : 0)
        .reduce((a , b) => a + parseFloat(b), 0);

      return transaction;
    });
  }
}
