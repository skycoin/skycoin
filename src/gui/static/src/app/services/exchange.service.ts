import { throwError as observableThrowError, Observable, BehaviorSubject, SubscriptionLike, of } from 'rxjs';
import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpErrorResponse } from '@angular/common/http';
import {
  ExchangeOrder,
  StoredExchangeOrder,
  TradingPair,
} from '../app.datatypes';
import { StorageService, StorageType } from './storage.service';
import * as moment from 'moment';
import { ApiService } from './api.service';
import { environment } from '../../environments/environment';
import { map, mergeMap, retryWhen, delay, catchError, tap } from 'rxjs/operators';

@Injectable()
export class ExchangeService {
  private readonly API_ENDPOINT = 'https://swaplab.cc/api/v3';
  private readonly STORAGE_KEY = 'exchange-orders';
  private readonly LAST_VIEWED_STORAGE_KEY = 'last-viewed-order';
  private readonly API_KEY = environment.swaplab.apiKey;
  private readonly TEST_MODE = environment.swaplab.activateTestMode;

  lastViewedOrderLoaded: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  private saveLastViewedSubscription: SubscriptionLike;
  private _lastViewedOrder: StoredExchangeOrder;

  set lastViewedOrder(order) {
    this._lastViewedOrder = order;

    if (this.saveLastViewedSubscription) {
      this.saveLastViewedSubscription.unsubscribe();
    }
    this.saveLastViewedSubscription = this.storageService.store(StorageType.CLIENT, this.LAST_VIEWED_STORAGE_KEY, JSON.stringify(order)).subscribe();
  }

  get lastViewedOrder() {
    return this._lastViewedOrder;
  }

  constructor(
    private http: HttpClient,
    private storageService: StorageService,
    private apiService: ApiService,
  ) {
    storageService.get(StorageType.CLIENT, this.LAST_VIEWED_STORAGE_KEY).subscribe(result => {
      this.lastViewedOrder = JSON.parse(result.data);
      this.lastViewedOrderLoaded.next(true);
    }, () => {
      this.lastViewedOrderLoaded.next(true);
    });
  }

  tradingPairs(): Observable<TradingPair[]> {
    return this.post('trading_pairs').pipe(map(data => data.result));
  }

  exchange(pair: string, fromAmount: number, toAddress: string, price: number): Observable<ExchangeOrder> {
    let response: ExchangeOrder;

    return this.post('orders', { pair, fromAmount, toAddress }).pipe(
      mergeMap(data => {
        response = data.result;

        return this.storeOrder(response, price);
      }), map(() => response));
  }

  status(id: string, devForceState?: string): Observable<ExchangeOrder> {
    if (this.TEST_MODE && !devForceState) {
      devForceState = 'user_waiting';
    }

    return this.post('orders/status', { id }, this.TEST_MODE ? { status: devForceState } : null).pipe(
      retryWhen((err) => {
        return err.pipe(mergeMap(response => {
          if (response instanceof HttpErrorResponse && response.status === 404) {
            return observableThrowError(response);
          }

          return of(response);
        }), delay(3000));
      }), map(data => data.result));
  }

  history() {
    return this.storageService.get(StorageType.CLIENT, this.STORAGE_KEY).pipe(
      map((res) => JSON.parse(res.data)));
  }

  isOrderFinished(order: ExchangeOrder) {
    return ['complete', 'error', 'user_deposit_timeout'].indexOf(order.status) > -1;
  }

  private post(url: string, body?: any, headers?: any): Observable<any> {
    return this.http.post(this.buildUrl(url), body, {
      responseType: 'json',
      headers: new HttpHeaders({
        'api-key': this.API_KEY,
        'Accept': 'application/json',
        ...headers,
      }),
    }).pipe(catchError((error: any) => this.apiService.processConnectionError(error)));
  }

  private buildUrl(url: string) {
    if (!this.TEST_MODE || url === 'trading_pairs') {
      return `${this.API_ENDPOINT}/${url}`;
    }

    return `${this.API_ENDPOINT}sandbox/${url}`;
  }

  private storeOrder(order: ExchangeOrder, price: number) {
    return this.history().pipe(
      catchError((err: HttpErrorResponse) => {
        try {
          if (err.status && err.status === 404) {
            return of([]);
          }
        } catch (e) {}

        return observableThrowError(err);
      }),
      mergeMap((oldOrders: StoredExchangeOrder[]) => this.storeOrderEntry(oldOrders, order, price)));
  }

  private storeOrderEntry(orders: StoredExchangeOrder[], order: ExchangeOrder, price: number): Observable<any> {
    const newOrder = {
      id: order.id,
      pair: order.pair,
      fromAmount: order.fromAmount,
      toAmount: order.toAmount,
      address: order.toAddress,
      timestamp: moment().unix(),
      price: price,
    };

    orders.push(newOrder);
    const data = JSON.stringify(orders);
    orders.pop();

    return this.storageService.store(StorageType.CLIENT, this.STORAGE_KEY, data).pipe(
      tap(() => orders.push(newOrder)));
  }
}
