import { of as observableOf, Subscription, Observable, of } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { BigNumber } from 'bignumber.js';

import { ApiService } from '../../../services/api/api.service';
import { ExplorerService } from '../../../services/explorer/explorer.service';
import { Transaction, GetBalanceResponse, GetUnconfirmedTransactionResponse, AddressTransactionsResponse } from 'app/app.datatypes';
import { PageBaseComponent } from '../page-base';
import { dataValidityTime } from 'app/app.config';

/**
 * Class for saving the responses returned by the server. This allows to avoid requesting the data
 * again when the user is checking different pages of the same address.
 */
export class CachedAddressDetails {
  /**
   * Address the data belongs to.
   */
  address: string;
  /**
   * Response obtained when the address transactions were requested.
   */
  transactionsResponse: AddressTransactionsResponse;
  /**
   * Response obtained when the address balance was requested.
   */
  balanceResponse: GetBalanceResponse;
  /**
   * Response obtained when the unconfirmed transactions were requested.
   */
  unconfirmedResponse: GetUnconfirmedTransactionResponse[];
}

/**
 * Page for showing info about a specific address.
 *
 * NOTE: This page may be reused between navigations, instead of being destroyed and recreated.
 * This helps to avoid realoading the same transactions again and again when an address has
 * so many transactions that it is needed to separate them in multiple pages. The main case in
 * which reloading data frequently can be very problematic is when checking the main address
 * of an exchange.
 */
@Component({
  selector: 'app-address-detail',
  templateUrl: './address-detail.component.html',
  styleUrls: ['./address-detail.component.scss']
})
export class AddressDetailComponent extends PageBaseComponent implements OnInit, OnDestroy {
  // Keys for persisting the server data, to be able to restore the state after navigation.
  private readonly persistentServerAddressResponseKey = 'serv-adr-response';
  private readonly persistentServerBalanceResponseKey = 'serv-blc-response';
  private readonly persistentServerUnconfirmedTxsResponseKey = 'serv-utx-response';

  /**
   * Var for saving the responses returned by the server. This allows to avoid requesting the data
   * again when the user is checking different pages of the same address.
   */
  private static cachedData: CachedAddressDetails;

  /**
   * Current address.
   */
  address: string;
  /**
   * Indicates if the data has been loaded.
   */
  dataLoaded = false;
  /**
   * How many transactions the address has.
   */
  totalTransactionsCount: number;
  /**
   * Total amount of coins received by the address.
   */
  totalReceived: BigNumber;
  /**
   * Total amount of coins sent from the address.
   */
  totalSent: BigNumber;
  /**
   * Current address coin balance.
   */
  balance: BigNumber;
  /**
   * Current address coin hour balance.
   */
  hoursBalance: BigNumber;
  /**
   * How many incoming coins are in unconfirmed transactions.
   */
  pendingIncomingCoins: BigNumber;
  /**
   * How many outgoing coins are in unconfirmed transactions.
   */
  pendingOutgoingCoins: BigNumber;
  /**
   * Total amount of pending coins (pendingIncomingCoins - pendingOutgoingCoins).
   */
  pendingCoins: BigNumber;
  /**
   * Total amount of pending coin hours.
   */
  pendingHours: BigNumber;
  /**
   * Transactions to be shown in the current page.
   */
  pageTransactions: any[];
  /**
   * Current page.
   */
  pageIndex = 0;
  /**
   * How many pages with transactions the current address has.
   */
  pageCount = 0;
  /**
   * Max number of transactions per page.
   */
  readonly pageSize = 25;
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
   * If true, the address has few transactions and all of them are stored in memory. If false,
   * only the transactions of the current page are in memory and no information about how
   * many coins the address has received and sent is calculated.
   */
  hasManyTransactions = false;

  // Subscriptions that will be cleaned when closing the page.
  private navParamsSubscription: Subscription;
  private operationSubscription: Subscription;

  constructor(
    private api: ApiService,
    private route: ActivatedRoute,
    public explorer: ExplorerService
  ) {
    super();
  }

  ngOnInit() {
    // Check the URL to detect changes in the requested address and page.
    this.navParamsSubscription = this.route.params.subscribe(params => {
      // Update the data.
      this.getData(params, true);
    });

    return super.ngOnInit();
  }

  ngOnDestroy() {
    // Clean all subscriptions.
    this.navParamsSubscription.unsubscribe();
    this.operationSubscription.unsubscribe();
  }

  /**
   * Gets the data of the requested address and page.
   */
  private getData(routeParams: Params, checkSavedData: boolean) {
    // Reset the cache var, but only if the address changed.
    if (!AddressDetailComponent.cachedData || AddressDetailComponent.cachedData.address !== routeParams['address']) {
      AddressDetailComponent.cachedData = new CachedAddressDetails();
      AddressDetailComponent.cachedData.address = routeParams['address'];
    }

    // Get the requested address and page.
    this.address = routeParams['address'];
    if (routeParams['page']) {
      this.pageIndex = parseInt(routeParams['page'], 10) - 1;
    } else {
      this.pageIndex = 0;
    }

    // Cancel any previous pending operation.
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }

    let oldSavedDataUsed = false;

    // Get the transactions.
    // Use saved data or get from the server. If there is no saved data, savedData is null.
    let savedData = checkSavedData ? this.getLocalValue(this.persistentServerAddressResponseKey) : null;
    let nextStep: Observable<any>;
    if (savedData) {
      // Reprocess the transactions, as recoveredTransactions is not saved correctly.
      (savedData.value as AddressTransactionsResponse).recoveredTransactions = this.explorer.processTransactionListFromServer(
        (savedData.value as AddressTransactionsResponse).originalTransactions,
        this.address, (savedData.value as AddressTransactionsResponse).addressHasManyTransactions
      );

      nextStep = of(savedData.value);
      oldSavedDataUsed = savedData.date < (new Date()).getTime() - dataValidityTime;
    } else {
      if (AddressDetailComponent.cachedData.transactionsResponse) {
        // If transactionsResponse still has a value, the address does not have many transactions
        // (so all are in memory) and only the page was changed, there is no need for loading
        // the transactions again.
        AddressDetailComponent.cachedData.transactionsResponse.currentPageIndex = this.pageIndex;
        nextStep = observableOf(AddressDetailComponent.cachedData.transactionsResponse);
      } else {
        nextStep = this.explorer.getTransactions(this.address, this.pageIndex + 1, this.pageSize);
      }
    }

    this.operationSubscription = nextStep.pipe(switchMap((response: AddressTransactionsResponse) => {
      if (!savedData) {
        this.saveLocalValue(this.persistentServerAddressResponseKey, response);
      }

      // Save the pagination information.
      this.totalTransactionsCount = response.totalTransactionsCount;
      this.hasManyTransactions = response.addressHasManyTransactions;
      this.pageCount = response.totalPages;
      this.pageIndex = response.currentPageIndex;

      // Save the response on the cache, if the address does not have many transactions.
      if (!this.hasManyTransactions) {
        AddressDetailComponent.cachedData.transactionsResponse = response;
      }

      if (!response.addressHasManyTransactions) {
        // Calculate the number of received and sent coins (counting the confirmed
        // transactions only).
        this.totalReceived = new BigNumber(0);
        response.recoveredTransactions.map(tx => this.totalReceived = this.totalReceived.plus(tx.balance.isGreaterThan(0) && tx.status ? tx.balance : 0));

        this.totalSent = new BigNumber(0);
        response.recoveredTransactions.map(tx => this.totalSent = this.totalSent.plus(tx.balance.isLessThan(0) && tx.status ? tx.balance : 0));
        this.totalSent = this.totalSent.negated();

        // Update the list of transactions that will be displayed in the UI.
        this.updateTransactions(response.recoveredTransactions);
      } else {
        this.pageTransactions = response.recoveredTransactions;
      }

      // Get the balance.
      // Use saved data or get from the server. If there is no saved data, savedData is null.
      savedData = checkSavedData ? this.getLocalValue(this.persistentServerBalanceResponseKey) : null;
      if (savedData) {
        oldSavedDataUsed = oldSavedDataUsed || (savedData.date < (new Date()).getTime() - dataValidityTime);
        return of(savedData.value);
      } else {
        if (AddressDetailComponent.cachedData.balanceResponse) {
          // If balanceResponse still has a value, only the page, and not the
          // address, was changed, so there is no need for loading the data again.
          return observableOf(AddressDetailComponent.cachedData.balanceResponse);
        } else {
          return this.api.getBalance(routeParams['address']);
        }
      }
    }), switchMap((response: GetBalanceResponse) => {
      if (!savedData) {
        this.saveLocalValue(this.persistentServerBalanceResponseKey, response);
      }

      // Save the response on the cache, if the address does not have many transactions.
      if (!this.hasManyTransactions) {
        AddressDetailComponent.cachedData.balanceResponse = response;
      }

      // Save the balance.
      this.balance = new BigNumber(response.confirmed.coins).dividedBy(1000000);
      this.hoursBalance = new BigNumber(response.confirmed.hours);

      // Get the pending transactions.
      // Use saved data or get from the server. If there is no saved data, savedData is null.
      savedData = checkSavedData ? this.getLocalValue(this.persistentServerUnconfirmedTxsResponseKey) : null;
      if (savedData) {
        oldSavedDataUsed = oldSavedDataUsed || (savedData.date < (new Date()).getTime() - dataValidityTime);
        return of(savedData.value);
      } else {
        if (AddressDetailComponent.cachedData.unconfirmedResponse) {
          // If unconfirmedResponse still has a value, only the page, and not the
          // address, was changed, so there is no need for loading the data again.
          return observableOf(AddressDetailComponent.cachedData.unconfirmedResponse);
        } else {
          return this.api.getUnconfirmedTransactions();
        }
      }
    })).subscribe((unconfirmed: GetUnconfirmedTransactionResponse[]) => {
      if (!savedData) {
        this.saveLocalValue(this.persistentServerUnconfirmedTxsResponseKey, unconfirmed);
      }

      // Save the response on the cache, if the address does not have many transactions.
      if (!this.hasManyTransactions) {
        AddressDetailComponent.cachedData.unconfirmedResponse = unconfirmed;
      }

      if (unconfirmed) {
        this.pendingIncomingCoins = new BigNumber(0);
        this.pendingOutgoingCoins = new BigNumber(0);
        this.pendingCoins = new BigNumber(0);
        this.pendingHours = new BigNumber(0);

        // Check al pending transactions.
        unconfirmed.forEach(tx => {
          // If there is an input with this address, all funds on it are considered going out.
          let coinsOut = new BigNumber(0);
          let hoursOut = new BigNumber(0);
          tx.transaction.inputs.forEach(input => {
            if (input.owner === routeParams['address']) {
              coinsOut = coinsOut.plus(input.coins);
              hoursOut = hoursOut.plus(input.calculated_hours);
            }
          });

          // If there is an output with this address, all funds on it are considered getting in.
          let coinsIn = new BigNumber(0);
          let hoursIn = new BigNumber(0);
          tx.transaction.outputs.forEach(output => {
            if (output.dst === routeParams['address']) {
              coinsIn = coinsIn.plus(output.coins);
              hoursIn = hoursIn.plus(output.hours);
            }
          });

          // Calculate the difference of funds entering and leaving.
          const coinsDifference = coinsIn.minus(coinsOut);
          const hoursDifference = hoursIn.minus(hoursOut);

          // Update the pending balance.
          this.pendingCoins = this.pendingCoins.plus(coinsDifference);
          this.pendingHours = this.pendingHours.plus(hoursDifference);

          if (coinsDifference.isGreaterThan(0)) {
            // If the total in this transaction is positive, update the amount of pending
            // incoming coins.
            this.pendingIncomingCoins = this.pendingIncomingCoins.plus(coinsDifference);
          } else if (coinsDifference.isLessThan(0)) {
            // If the total in this transaction is negative, update the amount of pending
            // outgoing coins.
            this.pendingOutgoingCoins = this.pendingOutgoingCoins.plus(coinsDifference.negated());
          }
        });
      }

      // If no transactions were found for the address.
      if (this.totalTransactionsCount < 1) {
        // Show that there are no transactions.
        this.longErrorMsg = 'addressDetail.withoutTransactions';
      }

      this.dataLoaded = true;

      // If old saved data was used, repeat the operation, ignoring the saved data.
      if (oldSavedDataUsed) {
        this.getData(routeParams, false);
      }

      //setTimeout(() => this.restoreScrollPosition(), 10);
    }, error => {
      if (!this.dataLoaded) {
        if (error.status >= 400 && error.status < 500) {
          // The address was not found.
          this.loadingMsg = 'general.noData';
          this.longErrorMsg = 'addressDetail.invalidAddress';
        } else {
          // Error loading the data.
          this.loadingMsg = 'general.shortLoadingErrorMsg';
          this.longErrorMsg = 'general.longLoadingErrorMsg';
        }
      }
    });
  }

  /**
   * Updates the list of transactions that will be displayed in the UI, to show only the
   * transactions of the current page.
   */
  private updateTransactions(alltransactions: Transaction[]) {
    // If the user request an invalid page, show a valid one.
    if (this.pageIndex > alltransactions.length / this.pageSize) {
      this.pageIndex = Math.floor(alltransactions.length / this.pageSize);
    }
    if (this.pageIndex < 0) {
      this.pageIndex = 0;
    }

    // Get the transaction of the current page.
    this.pageTransactions = [];
    for (let i = this.pageIndex * this.pageSize; i < (this.pageIndex + 1) * this.pageSize && i < alltransactions.length; i++) {
      this.pageTransactions.push(alltransactions[i]);
    }
  }
}
