import { switchMap } from 'rxjs/operators';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { of, Subscription } from 'rxjs';

import { Transaction } from '../../../app.datatypes';
import { ExplorerService } from '../../../services/explorer/explorer.service';
import { PageBaseComponent } from '../page-base';
import { dataValidityTime } from 'app/app.config';

/**
 * Page for showing info about a specific transaction.
 */
@Component({
  selector: 'app-transaction-detail',
  templateUrl: './transaction-detail.component.html',
  styleUrls: ['./transaction-detail.component.scss']
})
export class TransactionDetailComponent extends PageBaseComponent implements OnInit, OnDestroy {
  // Keys for persisting the server data, to be able to restore the state after navigation.
  private readonly persistentServerTransactionResponseKey = 'serv-trx-response';

  /**
   * Current transaction.
   */
  transaction: Transaction;
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
    private explorer: ExplorerService,
    private route: ActivatedRoute,
  ) {
    super();
  }

  ngOnInit() {
    // Get the URL params and request the data.
    this.pageSubscriptions.push(this.route.params.pipe(
      switchMap((params: Params) => this.explorer.getTransaction(params['txid']))
    ).subscribe(
      transaction => this.transaction = transaction,
      error => {
        if (error.status >= 400 && error.status < 500) {
          // The transaction was not found.
          this.loadingMsg = 'general.noData';
          this.longErrorMsg = 'transactionDetail.canNotFind';
        } else {
          // Error loading the data.
          this.loadingMsg = 'general.shortLoadingErrorMsg';
          this.longErrorMsg = 'general.longLoadingErrorMsg';
        }
      }
    ));

    return super.ngOnInit();
  }

  private loadData(checkSavedData: boolean) {
    let oldSavedDataUsed = false;

    let savedData;

    // Get the URL params and request the data.
    this.pageSubscriptions.push(this.route.params.pipe(
      switchMap((params: Params) => {
        // Get the transaction data.
        // Use saved data or get from the server. If there is no saved data, savedData is null.
        savedData = checkSavedData ? this.getLocalValue(this.persistentServerTransactionResponseKey) : null;
        if (savedData) {
          oldSavedDataUsed = savedData.date < (new Date()).getTime() - dataValidityTime;
          return of(savedData.value);
        } else {
          return this.explorer.getTransaction(params['txid']);
        }
      })
    ).subscribe(
      transaction => {
        if (!savedData) {
          this.saveLocalValue(this.persistentServerTransactionResponseKey, transaction);
        }
        this.transaction = transaction;

        // If old saved data was used, repeat the operation, ignoring the saved data.
        if (oldSavedDataUsed) {
          this.loadData(false);
        }
      },
      error => {
        if (this.transaction === undefined) {
          if (error.status >= 400 && error.status < 500) {
            // The transaction was not found.
            this.loadingMsg = 'general.noData';
            this.longErrorMsg = 'transactionDetail.canNotFind';
          } else {
            // Error loading the data.
            this.loadingMsg = 'general.shortLoadingErrorMsg';
            this.longErrorMsg = 'general.longLoadingErrorMsg';
          }
        }
      }
    ));
  }

  ngOnDestroy() {
    // Clean all subscriptions.
    this.pageSubscriptions.forEach(sub => sub.unsubscribe());
  }
}
