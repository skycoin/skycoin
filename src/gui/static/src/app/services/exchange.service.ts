import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';
import {
  ExchangeOrder,
  StoredExchangeOrder,
  TradingPair,
} from '../app.datatypes';
import { environment } from '../../environments/environment';
import { StorageService, StorageType } from './storage.service';
import * as moment from 'moment';

@Injectable()
export class ExchangeService {
  private readonly API_ENDPOINT = 'https://swaplab.cc/api/v3';
  private readonly API_KEY = 'w4bxe2tbf9beb72r';
  private readonly STORAGE_KEY = 'exchange-orders';

  private _lastOrder: ExchangeOrder;

  set lastOrder(order) {
    this._lastOrder = order;
  }

  get lastOrder() {
    return this._lastOrder;
  }

  constructor(
    private http: HttpClient,
    private storageService: StorageService,
  ) { }

  tradingPairs(): Observable<TradingPair[]> {
    return this.post('trading_pairs').map(data => data.result);
  }

  exchange(pair: string, fromAmount: number, toAddress: string): Observable<ExchangeOrder> {
    return this.post('orders', { pair, fromAmount, toAddress })
      .map(data => data.result)
      .do(result => this.storeOrder(result));
  }

  status(id: string): Observable<ExchangeOrder> {
    let lastKnownStatus = null;

    return Observable
      .timer(0, 30 * 1000)
      .flatMap(() => this.post('orders/status', { id }))
      .map(data => data.result)
      .do(response => lastKnownStatus = response)
      .filter(response => response !== null)
      .retryWhen((error$) => {
        return error$
          .do(error => {
            if (!(error.error instanceof ProgressEvent)) {
              throw error;
            }
          })
          .delay(3000);
      });

  }

  history() {
    return this.storageService.get(StorageType.CLIENT, this.STORAGE_KEY)
      .map((res) => JSON.parse(res.data));
  }

  isOrderFinished(order: ExchangeOrder) {
    return ['complete', 'error'].indexOf(order.status) > -1;
  }

  private post(url: string, body?: any, headers?: any): Observable<any> {
    return this.http.post(this.buildUrl(url), body, {
      responseType: 'json',
      headers: new HttpHeaders({
        'api-key': this.API_KEY,
        'Accept': 'application/json',
        ...headers,
      }),
    });
  }

  private buildUrl(url: string) {
    if (environment.production || url === 'trading_pairs') {
      return `${this.API_ENDPOINT}/${url}`;
    }

    return `${this.API_ENDPOINT}sandbox/${url}`;
  }

  private storeOrder(order: ExchangeOrder) {
    this.history().subscribe(
      (oldOrders: StoredExchangeOrder[]) => {
        this.storeOrderEntry(oldOrders, order);
      },
      () => {
        this.storeOrderEntry([], order);
      },
    );
  }

  private storeOrderEntry(orders: StoredExchangeOrder[], order: ExchangeOrder) {
    orders.push({
      id: order.id,
      pair: order.pair,
      fromAmount: order.fromAmount,
      timestamp: moment().unix(),
    });

    const data = JSON.stringify(orders);

    this.storageService.store(StorageType.CLIENT, this.STORAGE_KEY, data).subscribe();
  }
}
