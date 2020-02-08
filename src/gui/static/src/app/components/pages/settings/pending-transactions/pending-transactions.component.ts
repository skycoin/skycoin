import { Component, OnDestroy, OnInit } from '@angular/core';
import * as moment from 'moment';
import { SubscriptionLike } from 'rxjs';
import { NavBarSwitchService } from '../../../../services/nav-bar-switch.service';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { BigNumber } from 'bignumber.js';
import { BalanceAndOutputsService } from '../../../../services/wallet-operations/balance-and-outputs.service';
import { HistoryService, PendingTransactionData } from '../../../../services/wallet-operations/history.service';

@Component({
  selector: 'app-pending-transactions',
  templateUrl: './pending-transactions.component.html',
  styleUrls: ['./pending-transactions.component.scss'],
})
export class PendingTransactionsComponent implements OnInit, OnDestroy {
  transactions: PendingTransactionData[] = null;

  private transactionsSubscription: SubscriptionLike;
  private navbarSubscription: SubscriptionLike;

  constructor(
    private navBarSwitchService: NavBarSwitchService,
    private balanceAndOutputsService: BalanceAndOutputsService,
    private historyService: HistoryService,
  ) {
    this.navbarSubscription = this.navBarSwitchService.activeComponent.subscribe(value => {
      this.startCheckingTransactions(value);
    });
  }

  ngOnInit() {
    this.navBarSwitchService.showSwitch('pending-txs.my-transactions-button', 'pending-txs.all-transactions-button');
  }

  ngOnDestroy() {
    this.removeTransactionsSubscription();
    this.navbarSubscription.unsubscribe();
    this.navBarSwitchService.hideSwitch();
  }

  private startCheckingTransactions(value) {
    this.transactions = null;

    this.removeTransactionsSubscription();

    // Currently gets the data only one time.
    this.transactionsSubscription = this.historyService.getPendingTransactions().subscribe(transactions => {
      this.transactions = value === DoubleButtonActive.LeftButton ? transactions.user : transactions.all;
    });

    // Due to some changes, must use a method for updating or getting the pending transactions, not this.
    this.balanceAndOutputsService.refreshBalance();
  }

  private removeTransactionsSubscription() {
    if (this.transactionsSubscription) {
      this.transactionsSubscription.unsubscribe();
    }
  }
}
