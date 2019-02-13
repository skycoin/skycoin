import { Injectable, NgZone } from '@angular/core';
import { Http } from '@angular/http';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class PriceService {
  readonly PRICE_API_ID = 'sky-skycoin';

  price: Subject<number> = new BehaviorSubject<number>(null);

  constructor(
    private http: Http,
    private ngZone: NgZone,
  ) {
    this.ngZone.runOutsideAngular(() => {
      Observable.timer(0, 10 * 60 * 1000).subscribe(() => {
        this.http.get(`https://api.coinpaprika.com/v1/tickers/${this.PRICE_API_ID}?quotes=USD`)
          .map(response => response.json())
          .subscribe(response => this.ngZone.run(() => {
            this.price.next(response.quotes.USD.price);
          }));
      });
    });
  }
}
