import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';
import { ExchangeOrder, TradingPair } from '../app.datatypes';
import { environment } from '../../environments/environment';

@Injectable()
export class ExchangeService {
  private readonly API_ENDPOINT = 'https://swaplab.cc/api/v3';
  private readonly API_KEY = 'w4bxe2tbf9beb72r';

  private lastOrder: ExchangeOrder;

  constructor(
    private http: HttpClient,
  ) { }

  tradingPairs(): Observable<TradingPair[]> {
    return this.post('/trading_pairs').map(data => data.result);
  }

  exchange(pair: string, fromAmount: number, toAddress: string): Observable<ExchangeOrder> {
    return this.post('/orders', { pair, fromAmount, toAddress }).map(data => data.result);
  }

  status(id: string): Observable<ExchangeOrder> {
    return Observable
      .timer(0, 60 * 1000)
      .flatMap(() => this.post('/orders/status', { id }))
      .map(data => data.result)
      .takeWhile((order: ExchangeOrder) => order.status !== 'complete');
  }

  setLastOrder(order) {
    this.lastOrder = order;
  }

  getLastOrder() {
    return this.lastOrder;
  }

  private post(url: string, body?: any): Observable<any> {
    return this.http.post(this.buildUrl(url), body, {
      responseType: 'json',
      headers: new HttpHeaders({
        'api-key': this.API_KEY,
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
