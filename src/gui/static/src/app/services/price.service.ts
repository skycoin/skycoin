import { Injectable, NgZone } from '@angular/core';
import { Subject, BehaviorSubject, SubscriptionLike, timer, of } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { delay } from 'rxjs/operators';

@Injectable()
export class PriceService {
  readonly PRICE_API_ID = 'sky-skycoin';

  price: Subject<number> = new BehaviorSubject<number>(null);

  private readonly updatePeriod = 10 * 60 * 1000;
  private lastPriceSubscription: SubscriptionLike;
  private timerSubscriptions: SubscriptionLike[];

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
      this.timerSubscriptions.push(timer(this.updatePeriod, this.updatePeriod)
        .subscribe(() => {
          this.ngZone.run(() => !this.lastPriceSubscription ? this.loadPrice() : null );
        }));
    });

    this.timerSubscriptions.push(
      of(1).pipe(delay(firstConnectionDelay)).subscribe(() => {
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
