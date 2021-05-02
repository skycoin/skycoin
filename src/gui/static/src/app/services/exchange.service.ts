import { throwError as observableThrowError, Observable, BehaviorSubject, SubscriptionLike, of } from 'rxjs';
import { map, mergeMap, retryWhen, delay, catchError, tap } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import * as moment from 'moment';

import { StorageService, StorageType } from './storage.service';
import { environment } from '../../environments/environment';
import { processServiceError } from '../utils/errors';
import { OperationError } from '../utils/operation-error';

/**
 * Represents a trading pair acepted by Swaplab.
 */
export class TradingPair {
  /**
   * Coin to deposit.
   */
  from: string;
  /**
   * Coin that will be received.
   */
  to: string;
  /**
   * How many coins will be received per deposited coin.
   */
  price: number;
  /**
   * Name of the trading pair. Works as the ID of the trading pair.
   */
  pair: string;
  /**
   * Minimum number of coins that can be deposited per order.
   */
  min: number;
  /**
   * Maximum number of coins that can be deposited per order.
   */
  max: number;
}

/**
 * Response returned by the service when creating or checking an order.
 */
export class ExchangeOrder {
  pair: string;
  fromAmount: number|null;
  toAmount: number;
  toAddress: string;
  toTag: string|null;
  refundAddress: string|null;
  refundTag: string|null;
  id: string;
  exchangeAddress: string;
  exchangeTag: string|null;
  toTx?: string|null;
  status: string;
  message?: string;
}

/**
 * Data of an order saved in the persistent storage.
 */
export class StoredExchangeOrder {
  /**
   * ID of the order.
   */
  id: string;
  /**
   * Name of the coins pair.
   */
  pair: string;
  /**
   * How many coins the user must sent.
   */
  fromAmount: number;
  /**
   * Approximately how many coins the user will receive. The amount can change and will be
   * final only after the order has been completed.
   */
  toAmount: number;
  /**
   * Address where the user will receive the coins.
   */
  address: string;
  /**
   * Unix date indicating when the order was saved.
   */
  timestamp: number;
  /**
   * Approximately how many coins the user will receive per deposited coin, at the time
   * the order was created.
   */
  price: number;
}

/**
 * Allows to work with the Swaplab integration.
 */
@Injectable()
export class ExchangeService {
  /**
   * URL for connecting with the backend.
   */
  private readonly API_ENDPOINT = environment.production ? 'https://swaplab.cc/api/v3' : '/swaplab/api/v3';
  /**
   * Key used for saving the old orders in the persistent storage.
   */
  private readonly STORAGE_KEY = 'exchange-orders';
  /**
   * Key used for saving in the persistent storage the data of the last consulted order.
   */
  private readonly LAST_VIEWED_STORAGE_KEY = 'last-viewed-order';
  /**
   * Key for accessing the API.
   */
  private readonly API_KEY = environment.swaplab.apiKey;
  /**
   * If the service must work in test mode. If true, the service will call the sandbox API
   * endpoints and return false states for the orders.
   */
  private readonly TEST_MODE = environment.swaplab.activateTestMode;

  /**
   * Allows to know when the process of loading the last viewed order has been finished.
   * It does not guarantee that an order was loaded, as maybe the user has never
   * checked the state of an order.
   */
  lastViewedOrderLoaded: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  private saveLastViewedSubscription: SubscriptionLike;

  /**
   * Last order the user has checked. Must be updated manually. It is used to show the status
   * of that order again when opening the exchange section, so it must be erased if it is not
   * longer appropriate to show the status again immediately after opening the exchange section.
   */
  set lastViewedOrder(order: StoredExchangeOrder) {
    this._lastViewedOrder = order;

    if (this.saveLastViewedSubscription) {
      this.saveLastViewedSubscription.unsubscribe();
    }
    this.saveLastViewedSubscription = this.storageService.store(StorageType.CLIENT, this.LAST_VIEWED_STORAGE_KEY, JSON.stringify(order)).subscribe();
  }
  get lastViewedOrder(): StoredExchangeOrder {
    return this._lastViewedOrder;
  }
  private _lastViewedOrder: StoredExchangeOrder;

  constructor(
    private http: HttpClient,
    private storageService: StorageService,
  ) {
    // Load the data of the last consulted order.
    storageService.get(StorageType.CLIENT, this.LAST_VIEWED_STORAGE_KEY).subscribe(result => {
      this.lastViewedOrder = JSON.parse(result.data);
      this.lastViewedOrderLoaded.next(true);
    }, () => {
      this.lastViewedOrderLoaded.next(true);
    });
  }

  /**
   * Gets the complete list of trading pairs registered in the exchange service.
   */
  tradingPairs(): Observable<TradingPair[]> {
    return this.post('trading_pairs').pipe(map(data => data.result));
  }

  /**
   * Creates and saves a new exchange order.
   * @param pair Name of the pair the order will use.
   * @param fromAmount How many coins the user will send.
   * @param toAddress Address where the user wants to receive the coins.
   * @param price Price at the moment the order was created.
   */
  exchange(pair: string, fromAmount: number, toAddress: string, price: number): Observable<ExchangeOrder> {
    let response: ExchangeOrder;

    return this.post('orders', { pair, fromAmount, toAddress }).pipe(
      mergeMap(data => {
        response = data.result;

        return this.storeOrder(response, price);
      }), map(() => response));
  }

  /**
   * Checks on the backend the status of a previously created order.
   * @param id ID of the order.
   * @param devForceState If provided, the function will return the provided value as the state
   * of the order.
   * @returns Object with the current state of the order. If the service is in test mode, a
   * simulated state is returned.
   */
  status(id: string, devForceState?: string): Observable<ExchangeOrder> {
    if (this.TEST_MODE && !devForceState) {
      devForceState = 'user_waiting';
    }

    return this.post('orders/status', { id }, this.TEST_MODE ? { status: devForceState } : null).pipe(
      // Retry after a delay, unless the service says that the order does not exist.
      retryWhen((err) => {
        return err.pipe(mergeMap((response: OperationError) => {
          if (response.originalError && response.originalError.status && response.originalError.status === 404) {
            return observableThrowError(response);
          }

          return of(response);
        }), delay(3000));
      }), map(data => data.result));
  }

  /**
   * Returns the list of saved orders.
   */
  history(): Observable<StoredExchangeOrder[]> {
    return this.storageService.get(StorageType.CLIENT, this.STORAGE_KEY).pipe(
      map((res) => JSON.parse(res.data)));
  }

  /**
   * Allows to know if the status of an order indicates that the order has
   * been terminated.
   */
  isOrderFinished(order: ExchangeOrder): boolean {
    return ['complete', 'error', 'user_deposit_timeout'].indexOf(order.status) > -1;
  }

  /**
   * Sends a POST request to the service API.
   * @param url URL to send the request to, without the "http://x:x/api/vx/" part.
   * @param body Object with the key/value pairs to be sent to the backend.
   * @param headers Additional headers to send.
   */
  private post(url: string, body?: object, headers?: object): Observable<any> {
    return this.http.post(this.buildUrl(url), body, {
      responseType: 'json',
      headers: new HttpHeaders({
        'api-key': this.API_KEY,
        'Accept': 'application/json',
        ...headers,
      }),
    }).pipe(catchError((error: any) => observableThrowError(processServiceError(error))));
  }

  /**
   * Sanitizes the URL and adds the sandbox part if the service is running in test mode.
   */
  private buildUrl(url: string): string {
    if (!this.TEST_MODE || url === 'trading_pairs') {
      return `${this.API_ENDPOINT}/${url}`;
    }

    return `${this.API_ENDPOINT}sandbox/${url}`;
  }

  /**
   * Adds an order to the list of saved orders.
   * @param order Order to save.
   * @param price Price at the time the order was created.
   */
  private storeOrder(order: ExchangeOrder, price: number): Observable<any> {
    return this.history().pipe(
      // If there are no previous orders, add the new order to an empty array,
      catchError((err: OperationError) => {
        try {
          if (err.originalError && err.originalError.status && err.originalError.status === 404) {
            return of([]);
          }
        } catch (e) {}

        return observableThrowError(err);
      }),
      mergeMap((oldOrders: StoredExchangeOrder[]) => {
        const newOrder: StoredExchangeOrder = {
          id: order.id,
          pair: order.pair,
          fromAmount: order.fromAmount,
          toAmount: order.toAmount,
          address: order.toAddress,
          timestamp: moment().unix(),
          price: price,
        };

        oldOrders.push(newOrder);
        const data = JSON.stringify(oldOrders);
        oldOrders.pop();

        return this.storageService.store(StorageType.CLIENT, this.STORAGE_KEY, data).pipe(
          tap(() => oldOrders.push(newOrder)));
      }));
  }
}
