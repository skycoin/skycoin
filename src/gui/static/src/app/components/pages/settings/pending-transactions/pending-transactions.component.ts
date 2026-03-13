import { Component, OnDestroy, OnInit } from '@angular/core';
import { SubscriptionLike, of } from 'rxjs';
import { delay, mergeMap } from 'rxjs/operators';
import { MatDialog } from '@angular/material/dialog';

import { NavBarSwitchService } from '../../../../services/nav-bar-switch.service';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { HistoryService, PendingTransactionData } from '../../../../services/wallet-operations/history.service';
import { ConfirmationComponent, DefaultConfirmationButtons } from '../../../layout/confirmation/confirmation.component';

/**
 * Allows to see the list of pending transactions. It uses the nav bar to know if it must show
 * all pending tx or just the pending tx affecting the user.
 */
@Component({
    selector: 'app-pending-transactions',
    templateUrl: './pending-transactions.component.html',
    styleUrls: ['./pending-transactions.component.scss'],
    standalone: false
})
export class PendingTransactionsComponent implements OnInit, OnDestroy {
  // Transactions to show on the UI.
  transactions: PendingTransactionData[] = null;

  private transactionsSubscription: SubscriptionLike;
  private navbarSubscription: SubscriptionLike;

  private selectedNavbarOption: DoubleButtonActive;

  // Time interval in which periodic data updates will be made.
  private readonly updatePeriod = 10 * 1000;
  // Time interval in which the periodic data updates will be restarted after an error.
  private readonly errorUpdatePeriod = 2 * 1000;

  constructor(
    private navBarSwitchService: NavBarSwitchService,
    private historyService: HistoryService,
    private dialog: MatDialog,
  ) {
    this.navbarSubscription = this.navBarSwitchService.activeComponent.subscribe(value => {
      this.selectedNavbarOption = value;
      this.transactions = null;
      this.startDataRefreshSubscription(0);
    });
  }

  ngOnInit() {
    this.navBarSwitchService.showSwitch('pending-txs.my-transactions-button', 'pending-txs.all-transactions-button');
  }

  ngOnDestroy() {
    this.navbarSubscription.unsubscribe();
    this.removeTransactionsSubscription();
    this.navBarSwitchService.hideSwitch();
  }

  deleteTransaction(txid: string) {
    ConfirmationComponent.openDialog(this.dialog, {
      text: 'pending-txs.delete-confirm',
      defaultButtons: DefaultConfirmationButtons.YesNo,
      redTitle: true,
    }).afterClosed().subscribe(result => {
      if (result) {
        this.historyService.deletePendingTransaction(txid).subscribe(() => {
          this.transactions = null;
          this.startDataRefreshSubscription(0);
        });
      }
    });
  }

  /**
   * Makes the page start updating the data periodically. If this function was called before,
   * the previous updating procedure is cancelled.
   * @param delayMs Delay before starting to update the data.
   */
  private startDataRefreshSubscription(delayMs: number) {
    this.removeTransactionsSubscription();

    this.transactionsSubscription = of(0).pipe(delay(delayMs), mergeMap(() => this.historyService.getPendingTransactions())).subscribe(transactions => {
      this.transactions = this.selectedNavbarOption === DoubleButtonActive.LeftButton ? transactions.user : transactions.all;

      // Update again after a delay.
      this.startDataRefreshSubscription(this.updatePeriod);
    }, () => this.startDataRefreshSubscription(this.errorUpdatePeriod));
  }

  private removeTransactionsSubscription() {
    if (this.transactionsSubscription) {
      this.transactionsSubscription.unsubscribe();
    }
  }
}
