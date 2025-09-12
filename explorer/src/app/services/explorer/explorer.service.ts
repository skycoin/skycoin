import { Observable, Subscription, of } from 'rxjs';
import { map, mergeMap, delay } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { BigNumber } from 'bignumber.js';

import { ApiService } from '../api/api.service';
import { Block, parseGenericBlock, parseGetUnconfirmedTransaction, Transaction, parseGetTransaction, AddressTransactionsResponse } from '../../app.datatypes';
import { CoinIdentifiers, namedAddresses } from '../../app.config';

/**
 * Allows to request information from the server. This service returns processed objects,
 * unlike ApiService, which returns almost the same data returned by the node.
 */
@Injectable()
export class ExplorerService {
  // Max transactions an address can have to be considered as not having many transactions.
  private readonly manyTransactionsCount = 100;

  // Variables with information about the node.
  private internalFullCoinName = ' ';
  private internalCoinName = ' ';
  private internalHoursName = ' ';
  private internalHoursNameSingular = ' ';
  private internalMaxDecimals = 6;

  // Map for getting the names associated with particular addresses.
  private namedAddressesMap = new Map<string, string>();

  /**
   * Full coin name returned by the node.
   */
  get fullCoinName(): string {
    return this.internalFullCoinName;
  }
  /**
   * Small coin name returned by the node.
   */
  get coinName(): string {
    return this.internalCoinName;
  }
  /**
   * Plural coin hours name returned by the node.
   */
  get hoursName(): string {
    return this.internalHoursName;
  }
  /**
   * Singular coin hours name returned by the node.
   */
  get hoursNameSingular(): string {
    return this.internalHoursNameSingular;
  }
  /**
   * Max number of decimal places the coin amounts can have.
   */
  get maxDecimals(): number {
    return this.internalMaxDecimals;
  }

  private nodeUrlSubscription: Subscription;
  private initializationSubscription: Subscription;

  constructor(
    private api: ApiService,
  ) { }

  /**
   * Initializes the service.
   */
  initialize() {
    namedAddresses.forEach(namedAddress => {
      this.namedAddressesMap.set(namedAddress.address, namedAddress.name);
    });

    this.getNodeInfo(0);
  }

  /**
   * Gets the basic info about the backend.
   *
   * @param delayMs Delay before starting to get the data.
   */
  private getNodeInfo(delayMs: number) {
    if (this.initializationSubscription) {
      this.initializationSubscription.unsubscribe();
    }

    this.initializationSubscription = of(0).pipe(delay(delayMs), mergeMap(() => this.api.getHealth())).subscribe(response => {
      // Get the information from the node if available.
      if (response.fiber && response.fiber.display_name) {
        this.internalFullCoinName = response.fiber.display_name;
      } else {
        this.internalFullCoinName = CoinIdentifiers.fullName;
      }

      if (response.fiber && response.fiber.ticker) {
        this.internalCoinName = response.fiber.ticker;
      } else {
        this.internalCoinName = CoinIdentifiers.coinName;
      }

      if (response.fiber && response.fiber.coin_hours_display_name) {
        this.internalHoursName = response.fiber.coin_hours_display_name;
      } else {
        this.internalHoursName = CoinIdentifiers.HoursName;
      }

      if (response.fiber && response.fiber.coin_hours_display_name_singular) {
        this.internalHoursNameSingular = response.fiber.coin_hours_display_name_singular;
      } else {
        this.internalHoursNameSingular = CoinIdentifiers.HoursNameSingular;
      }

      if (response.user_verify_transaction && response.user_verify_transaction.max_decimals) {
        this.internalMaxDecimals = response.user_verify_transaction.max_decimals;
      } else {
        this.internalMaxDecimals = 6;
      }
    }, () => {
      // Retry in case of error.
      this.getNodeInfo(3000);
    });
  }

  /**
   * Gets a block by its ID (sequence number).
   *
   * @param id Block ID (sequence number).
   */
  getBlock(id: number): Observable<Block> {
   return this.api.getBlockById(id).pipe(map(response => parseGenericBlock(response)));
  }

  /**
   * Gets an array with the blocks in a specific range.
   *
   * @param start Number (height) of the first block (inclusive).
   * @param end Number (height) of the last block (inclusive).
   */
  getBlocks(start: number, end: number): Observable<Block[]> {
    return this.api.getBlocks(start, end).pipe(
      map(response => response.blocks.map(block => parseGenericBlock(block)).sort((a, b) => b.id - a.id)));
  }

  /**
   * Gets a block by its hash.
   *
   * @param hash Block hash.
   */
  getBlockByHash(hash: string): Observable<Block> {
    return this.api.getBlockByHash(hash).pipe(map(response => parseGenericBlock(response)));
  }

  /**
   * Gets the list of transactions of a specific address. Depending on how many transactions
   * the address has, the function may return a paginated results or all the transactions.
   *
   * @param address Address to consult.
   * @param page Number of the desired results page.
   * @param pageSize Max number of transactions each results page can have.
   * @returns An object with the requested transactions and information about the pagination.
   * IMPORTANT: if the address has many transactions, the response will only contain the
   * transactions for the requested page (the response will indicate the address is considered
   * to have many transactions) and information about the balance variation after each transaction
   * will not be included. However, if the address does not have many transactions, the response
   * will include all the transactions the address has, which could be more than the requested max
   * per page, and the caller will be responsible for using only the transactions needed.
   */
  getTransactions(address: string, page: number, pageSize: number): Observable<AddressTransactionsResponse> {
    let transactionsCount: number;
    let hasManyTransactions: boolean;
    let currentPageIndex: number;
    let totalPages: number;

    // Request one page with just one transaction to know how many transactions the address has.
    return this.api.getAddressWithPagination(address, 1, 1).pipe(mergeMap(response => {
      transactionsCount = response.page_info.total_pages;
      totalPages = Math.ceil(transactionsCount / pageSize);

      // Make sure the page is not outside the boundaries.
      if (page < 1) {
        currentPageIndex = 0;
      } else if (page > totalPages) {
        currentPageIndex = totalPages - 1;
      } else {
        currentPageIndex = page - 1;
      }

      hasManyTransactions = transactionsCount > this.manyTransactionsCount;

      // Get the transaction list depending on how may transactions the address has.
      let nextStep: Observable<any>;
      if (!hasManyTransactions) {
        nextStep = this.api.getAddress(address);
      } else {
        nextStep = this.api.getAddressWithPagination(address, currentPageIndex + 1, pageSize).pipe(map(resp => resp.txns));
      }

      return nextStep;
    }), map(response => {
      const processedTransactions = this.processTransactionListFromServer(response, address, hasManyTransactions);

      return {
        totalTransactionsCount: transactionsCount,
        currentPageIndex: currentPageIndex,
        totalPages: totalPages,
        addressHasManyTransactions: hasManyTransactions,
        recoveredTransactions: processedTransactions,
        originalTransactions: response,
      } as AddressTransactionsResponse;
    }));
  }

  processTransactionListFromServer(transactionList: any, address: string, hasManyTransactions: boolean): Transaction[] {
    // Sort to get the lastest transactions last (it will be reversed below).
    let response = transactionList.sort((a, b) => a.txn.timestamp - b.txn.timestamp);

    // Process the response.
    let currentBalance = new BigNumber('0');
    response = response.map(rawTx => {
      const parsedTx = parseGetTransaction(rawTx, address);

      // Calculate the balance variation after every transaction, if all transactions
      // are in memory.
      if (!hasManyTransactions) {
        parsedTx.initialBalance = currentBalance;
        currentBalance = currentBalance.plus(parsedTx.balance);
        parsedTx.finalBalance = currentBalance;
      }

      return parsedTx;
    });

    return response.reverse();
  }

  /**
   * Gets a transaction by its hash.
   *
   * @param transactionId Transaction hash.
   */
  getTransaction(transactionId: string): Observable<Transaction> {
    return this.api.getTransaction(transactionId).pipe(
      map(response => parseGetTransaction(response)));
  }

  /**
   * Returns the name associated to an address, enclosed in parentheses. If no name has been
   * associated to the address, returns an empty string.
   */
  getAddressName(address: string): string {
    if (this.namedAddressesMap.has(address)) {
      return ' (' + this.namedAddressesMap.get(address) + ')';
    }

    return '';
  }
}
