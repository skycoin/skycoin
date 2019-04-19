import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';
import { ExchangeOrder, TradingPair } from '../app.datatypes';
import { environment } from '../../environments/environment';

@Injectable()
export class ExchangeService {
  private readonly API_ENDPOINT = 'https://swaplab.cc/api/v3';
  private readonly API_KEY = 'w4bxe2tbf9beb72r';

  private _lastOrder: ExchangeOrder;

  set lastOrder(order) {
    this._lastOrder = order;
  }

  get lastOrder() {
    return this._lastOrder;
  }

  constructor(
    private http: HttpClient,
  ) { }

  tradingPairs(): Observable<TradingPair[]> {
    return this.post('trading_pairs').map(data => data.result);
  }

  exchange(pair: string, fromAmount: number, toAddress: string): Observable<ExchangeOrder> {
    return this.post('orders', { pair, fromAmount, toAddress }).map(data => data.result);
  }

  status(id: string): Observable<ExchangeOrder> {
    const statuses = [
      'user_waiting',
      'market_waiting_confirmations',
      'market_confirmed',
      'market_exchanged',
      'market_withdraw_waiting',
      'complete',
      'error', // move higher to see error state
    ];let index = 0;
    return Observable
      .timer(0, 3 * 1000)
      .flatMap(() => this.post('orders/status', { id }, { status: statuses[index++] }))
      .map(data => data.result);
  }

  isOrderFinished(order: ExchangeOrder) {
    return ['complete', 'error'].indexOf(order.status);
  }

  private post(url: string, body?: any, headers?: any): Observable<any> {
    return this.http.post(this.buildUrl(url), body, {
      responseType: 'json',
      headers: new HttpHeaders({
        'api-key': this.API_KEY,
        ...headers,
      }),
    });
  }

  private buildUrl(url: string) {
    if (environment.production || url === '/trading_pairs') {
      return `${this.API_ENDPOINT}/${url}`;
    }

    return `${this.API_ENDPOINT}sandbox/${url}`;
  }
}
