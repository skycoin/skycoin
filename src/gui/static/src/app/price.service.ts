import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import 'rxjs/add/observable/timer';

@Injectable()
export class PriceService {

  price: Subject<number> = new BehaviorSubject<number>(null);

  constructor(
    private http: Http,
  ) {
    setTimeout(() => {
      this.http.get('https://api.coinmarketcap.com/v1/ticker/skycoin/')
        .map(response => response.json()[0])
        .subscribe(data => this.price.next(data.price_usd));
    }, 1000);
  }

}
