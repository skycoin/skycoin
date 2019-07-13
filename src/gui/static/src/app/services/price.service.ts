import { Injectable, NgZone } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Observable } from 'rxjs/Observable';
import { HttpClient } from '@angular/common/http';
import { ISubscription } from 'rxjs/Subscription';

@Injectable()
export class PriceService {
  readonly PRICE_API_ID = 'sky-skycoin';

  price: Subject<number> = new BehaviorSubject<number>(null);

  private readonly updatePeriod = 10 * 60 * 1000;
  private lastPriceSubscription: ISubscription;
  private timerSubscriptions: ISubscription[];

  constructor(
    private http: HttpClient,
    private ngZone: NgZone,
  ) {
    this.startTimer();
  }

  private startTimer(firstConnectionDelay = 0) {
    if (this.timerSubscriptions) {
      this.timerSubscriptions.forEach(sub => sub.unsubscribe());
    }

    this.timerSubscriptions = [];

    this.ngZone.runOutsideAngular(() => {
      this.timerSubscriptions.push(Observable.timer(this.updatePeriod, this.updatePeriod)
        .subscribe(() => {
          this.ngZone.run(() => !this.lastPriceSubscription ? this.loadPrice() : null );
        }));
    });

    this.timerSubscriptions.push(
      Observable.of(1).delay(firstConnectionDelay).subscribe(() => {
        this.ngZone.run(() => this.loadPrice());
      }));
  }

  private loadPrice() {
    if (!this.PRICE_API_ID) {
      return;
    }

    if (this.lastPriceSubscription) {
      this.lastPriceSubscription.unsubscribe();
    }

    this.lastPriceSubscription = this.http.get(`https://api.coinpaprika.com/v1/tickers/${this.PRICE_API_ID}?quotes=USD`)
      .subscribe((response: any) => {
        this.lastPriceSubscription = null;
        this.price.next(response.quotes.USD.price);
      },
      () => this.startTimer(30000));
  }
}
