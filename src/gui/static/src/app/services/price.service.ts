import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class PriceService {
  readonly CMC_TICKER_ID = 1619;

  price: Subject<number> = new BehaviorSubject<number>(null);

  constructor(
    private http: Http,
  ) {
    Observable.timer(0, 10 * 60 * 1000).subscribe(() => {
      this.http.get(`https://api.coinmarketcap.com/v2/ticker/${this.CMC_TICKER_ID}/`)
        .map(response => response.json())
        .subscribe(response => this.price.next(response.data.quotes.USD.price));
    });
  }
}
