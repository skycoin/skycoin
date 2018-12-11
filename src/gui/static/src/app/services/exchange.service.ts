import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';
import { TradingPair } from '../app.datatypes';

@Injectable()
export class ExchangeService {
  private readonly API_ENDPOINT = 'https://swaplab.cc/api/v3';
  private readonly API_KEY = 'w4bxe2tbf9beb72r';

  constructor(
    private http: HttpClient,
  ) { }

  tradingPairs(): Observable<TradingPair[]> {
    return this.post('/trading_pairs').map(data => data.result);
  }

  private post(url: string, body?: any): Observable<any> {
    return this.http.post(`${this.API_ENDPOINT}/${url}`, body, {
      responseType: 'json',
      headers: new HttpHeaders({
        'api-key': this.API_KEY,
      }),
    });
  }
}
