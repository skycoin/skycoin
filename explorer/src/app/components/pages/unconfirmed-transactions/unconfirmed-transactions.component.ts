import { Component, OnInit, OnDestroy } from '@angular/core';
import { Observable, of, Subscription } from 'rxjs';

import { parseGetUnconfirmedTransaction, Transaction } from '../../../app.datatypes';
import { PageBaseComponent } from '../page-base';
import { dataValidityTime } from 'app/app.config';
import { ApiService } from 'app/services/api/api.service';

/**
 * Page for showing the list of unconfirmed (pending) transactions.
 */
@Component({
  selector: 'app-unconfirmed-transactions',
  templateUrl: './unconfirmed-transactions.component.html',
  styleUrls: ['./unconfirmed-transactions.component.scss']
})
export class UnconfirmedTransactionsComponent extends PageBaseComponent implements OnInit, OnDestroy {
  // Keys for persisting the server data, to be able to restore the state after navigation.
  private readonly persistentServerUnconfirmedTxsResponseKey = 'serv-utx-response';

  /**
   * Transactions list.
   */
  transactions: Transaction[];
  /**
   * Date of the oldest unconfirmed transaction.
   */
  leastRecent: number;
  /**
   * Date of the most recent unconfirmed transaction.
   */
  mostRecent: number;
  /**
   * Combined size (in bytes) of all unconfirmed transactions.
   */
  totalSize: number;
  /**
   * Small text to be shown in variaous parts (not in the loading control) while loading the data.
   * It may also contain small error messages.
   */
  loadingMsg = 'general.loadingMsg';
  /**
   * Error message to be shown in the loading control if there is a problem.
   */
  longErrorMsg: string;

  /**
   * Observable subscriptions that will be cleaned when closing the page.
   */
  private pageSubscriptions: Subscription[] = [];

  constructor(
    private api: ApiService,
  ) {
    super();
  }

  ngOnInit() {
    this.loadData(true);

    return super.ngOnInit();
  }

  private loadData(checkSavedData: boolean) {
    let oldSavedDataUsed = false;

    // Request the data.
    // Use saved data or get from the server. If there is no saved data, savedData is null.
    const savedData = checkSavedData ? this.getLocalValue(this.persistentServerUnconfirmedTxsResponseKey) : null;
    let nextOperation: Observable<any> = this.api.getUnconfirmedTransactions();
    if (savedData) {
      nextOperation = of(savedData.value);
      oldSavedDataUsed = savedData.date < (new Date()).getTime() - dataValidityTime;
    }

    this.pageSubscriptions.push(nextOperation.subscribe(transactions => {
      if (!savedData) {
        this.saveLocalValue(this.persistentServerUnconfirmedTxsResponseKey, transactions);
      }

      transactions = transactions.map(rawTx => parseGetUnconfirmedTransaction(rawTx));
      this.transactions = transactions;
      //this.transactions = transactions.map(rawTx => parseGetUnconfirmedTransaction(rawTx))
      // If there are unconfirmed transactions, calculate the values to be shown in the UI.
      if (transactions.length > 0) {
        const orderedList = transactions.sort((a, b) => b.timestamp - a.timestamp);
        this.mostRecent = orderedList[0].timestamp;
        this.leastRecent = orderedList[orderedList.length - 1].timestamp;
        this.totalSize = orderedList.map(tx => tx.length).reduce((sum, current) => sum + current);
      }

      // If old saved data was used, repeat the operation, ignoring the saved data.
      if (oldSavedDataUsed) {
        this.loadData(false);
      }
    }, () => {
      if (this.transactions === undefined) {
        // Error loading the data.
        this.loadingMsg = 'general.shortLoadingErrorMsg';
        this.longErrorMsg = 'general.longLoadingErrorMsg';
      }
    }));
  }

  ngOnDestroy() {
    // Clean all subscriptions.
    this.pageSubscriptions.forEach(sub => sub.unsubscribe());
  }
}
