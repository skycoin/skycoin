import { throwError as observableThrowError,  Observable, ReplaySubject } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Blockchain, GetBlocksResponse, GetBlockchainMetadataResponse, GetUnconfirmedTransactionResponse,
    GetCurrentBalanceResponse, GenericBlockResponse, RichlistEntry, GetBalanceResponse, GetTransactionResponse, GetCoinSupplyResponse, GetSyncStateResponse } from '../../app.datatypes';

/**
 * Allows to request information from the server. This service returns almost the same data
 * returned by the node, unlike ApiService, which returns processed objects.
 *
 * This service uses the api of the intermediate server (explorer.go). The API is similar to the
 * Skycoin node API, but there are differences. There is more info clicking the "Explorer API"
 * link in the main page of the explorer. For even more info, seach the value of the "Internal
 * skycoin node path:" section in the Skycoin REST API docs:
 * https://github.com/skycoin/skycoin/blob/develop/src/api/README.md
 */
@Injectable()
export class ApiService {
  /**
   * Key used for saving the local node URL in the persistent storage and the window object.
   */
  private static readonly localNodeUrlKey = 'nodeUrl';

  /**
   * URL used to connect to the API.
   */
  private url = '/api/';
  /**
   * Subject for emitting every time the local node URL is changed or removed.
   */
  private localNodeUrlSubject: ReplaySubject<string> = new ReplaySubject<string>(1);

  constructor(
    private http: HttpClient
  ) { }

  /**
   * Initializes the service.
   */
  initialize() {
    // Get the URL of the local node that must be used as backend, if an URL was saved before.
    const localNodeUrl = localStorage.getItem(ApiService.localNodeUrlKey);
    if (localNodeUrl) {
      window[ApiService.localNodeUrlKey] = localNodeUrl;
    }

    this.localNodeUrlSubject.next(localNodeUrl);
  }

  /**
   * Sets the URL of the local node that must be used as backend. If the URL is null, any
   * previously saved URL is removed.
   */
  setNodeUrl(url: string) {
    window[ApiService.localNodeUrlKey] = url;

    if (url) {
      localStorage.setItem(ApiService.localNodeUrlKey, url);
    } else {
      localStorage.removeItem(ApiService.localNodeUrlKey);
    }

    this.localNodeUrlSubject.next(url);
  }

  /**
   * Emits every time the local node URL is changed or removed.
   */
  get localNodeUrl(): Observable<string> {
    return this.localNodeUrlSubject.asObservable();
  }

  /**
   * Get information about the state of the node.
   */
  getHealth(): Observable<any> {
    const url = !this.nodeUrl() ? 'health' : 'v1/health';

    return this.get(url);
  }

  /**
   * Gets the list of transactions of a specific address.
   *
   * @param address Address to consult.
   */
  getAddress(address: string): Observable<GetTransactionResponse[]> {
    const url = !this.nodeUrl() ? 'transactions' : 'v1/transactions';

    return this.get(url, { addrs: address, verbose: 1 });
  }

  /**
   * Gets the list of transactions of a specific address, using pagination.
   *
   * @param address Address to consult.
   */
  getAddressWithPagination(address: string, page: number, pageSize: number): Observable<any> {
    const url = !this.nodeUrl() ? 'paginatedTransactions' : 'v2/transactions';

    return this.get(url, { addrs: address, page: page, limit: pageSize, sort: 'desc', verbose: 1 })
      .pipe(map((response: any) => response.data));
  }

  /**
   * Gets the list of unconfirmed transactions.
   */
  getUnconfirmedTransactions(): Observable<GetUnconfirmedTransactionResponse[]> {
    const url = !this.nodeUrl() ? 'pendingTxs' : 'v1/pendingTxs';

    return this.get(url, { verbose: 1 });
  }

  /**
   * Gets a block by its ID (sequence number).
   *
   * @param id Block ID (sequence number).
   */
  getBlockById(id: number): Observable<GenericBlockResponse> {
    const url = !this.nodeUrl() ? 'block' : 'v1/block';

    return this.get(url, { seq: id, verbose: 1 });
  }

  /**
   * Gets a block by its hash.
   *
   * @param hash Block hash.
   */
  getBlockByHash(hash: string): Observable<GenericBlockResponse> {
    const url = !this.nodeUrl() ? 'block' : 'v1/block';

    return this.get(url, { hash: hash, verbose: 1 });
  }

  /**
   * Gets an array with the blocks in a specific range.
   *
   * @param startNumber Number (height) of the first block (inclusive).
   * @param endNumber Number (height) of the last block (inclusive).
   */
  getBlocks(startNumber: number, endNumber: number): Observable<GetBlocksResponse> {
    const url = !this.nodeUrl() ? 'blocks' : 'v1/blocks';

    return this.get(url, { start: startNumber, end: endNumber });
  }

  /**
   * Gets information about the state of the blockchain.
   */
  getBlockchainMetadata(): Observable<Blockchain> {
    const url = !this.nodeUrl() ? 'blockchain/metadata' : 'v1/blockchain/metadata';

    return this.get(url).pipe(
      map((res: GetBlockchainMetadataResponse) => ({
        blocks: res.head.seq,
      })));
  }

  /**
   * Gets information about the current coin supply.
   */
  getCoinSupply(): Observable<GetCoinSupplyResponse> {
    const url = !this.nodeUrl() ? 'coinSupply' : 'v1/coinSupply';

    return this.get(url);
  }

  /**
   * Gets the list of unspent outputs of an address.
   *
   * @param address Address to consult.
   */
  getCurrentBalance(address: string): Observable<GetCurrentBalanceResponse> {
    const url = !this.nodeUrl() ? 'currentBalance' : 'v1/outputs';

    return this.get(url, { addrs: address });
  }

  /**
   * Gets the balance of an address.
   *
   * @param address Address to consult.
   */
  getBalance(address: string): Observable<GetBalanceResponse> {
    const url = !this.nodeUrl() ? 'balance' : 'v1/balance';

    return this.get(url, { addrs: address });
  }

  /**
   * Gets a transaction by its hash.
   *
   * @param transactionId Transaction hash.
   */
  getTransaction(transactionId: string): Observable<GetTransactionResponse> {
    const url = !this.nodeUrl() ? 'transaction' : 'v1/transaction';

    return this.get(url, { txid: transactionId, verbose: 1 });
  }

  /**
   * Gets the list of unlocked addresses with most coins.
   */
  getRichlist(): Observable<RichlistEntry[]> {
    const url = !this.nodeUrl() ? 'richlist' : 'v1/richlist';

    return this.get(url).pipe(map((response: any) => response.richlist));
  }

  /**
   * Gets the sync state of the node.
   */
  getSyncState(): Observable<GetSyncStateResponse> {
    const url = !this.nodeUrl() ? 'blockchain/progress' : 'v1/blockchain/progress';

    return this.get(url);
  }

  // Old methods

  /**
   * Sends a GET request to the API.
   *
   * @param url URL segment (the URL of the API endpont after the "/api/" part).
   * @param options Arguments to send as URL params.
   */
  private get(url: string, options: object = null): any {
    return this.http.get(this.getUrl(url, options)).pipe(
      catchError((error: any) => observableThrowError(error || 'Server error'))
    );
  }

  /**
   * Process an object to use its properties for building a querystring, so it can be
   * sent in an API call.
   *
   * @param parameters Object with params and values to build the querystring.
   */
  private getQueryString(parameters: object = null): string {
    if (!parameters) {
      return '';
    }

    return Object.keys(parameters).reduce((array, key) => {
      array.push(key + '=' + encodeURIComponent(parameters[key]));
      return array;
    }, []).join('&');
  }

  /**
   * Gets the complete URL for making an API request.
   *
   * @param url URL segment (the URL of the API endpont after the "/api/" part).
   * @param options Arguments to send as URL params.
   */
  private getUrl(url: string, options: object = null): string {
    // Ensure that there is no a '/' at the beginning.
    if (url.startsWith('/')) {
      url = url.substr(1, url.length - 1);
    }

    let initialPart = this.nodeUrl();
    if (!initialPart) {
      initialPart = this.url;
    }

    return initialPart + url + '?' + this.getQueryString(options);
  }

  /**
   * Returns the URL of the node that the explorer must use as backend. If it does not return a
   * valid value, the explorer must use as backend the Go intermediate server included with it.
   */
  private nodeUrl() {
    return window[ApiService.localNodeUrlKey];
  }
}
