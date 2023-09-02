import { Injectable, NgZone } from '@angular/core';
import { Subject, BehaviorSubject, of, Subscription, Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { delay, mergeMap } from 'rxjs/operators';

import { AppConfig } from '../app.config';
import { environment } from '../../environments/environment';

/**
 * Maintains updated and allows to known the current USD price of the coin.
 */
@Injectable()
export class PriceService {
  private readonly PRICE_API_ID = AppConfig.priceApiId;

  /**
   * Allows to know the current USD price of the coin.
   */
  get price(): Observable<number> {
    return this.priceInternal.asObservable();
  }
  private priceInternal: Subject<number> = new BehaviorSubject<number>(null);

  /**
   * Time interval in which periodic data updates will be made.
   */
  private readonly updatePeriod = 10 * 60 * 1000;
  /**
   * Time interval in which the periodic data updates will be restarted after an error.
   */
  private readonly errorUpdatePeriod = 30 * 1000;
  private priceSubscription: Subscription;

  constructor(
    private http: HttpClient,
    private ngZone: NgZone,
  ) {
    this.startDataRefreshSubscription(0);
  }

  /**
   * Makes the service start updating the data periodically. If this function was called
   * before, the previous updating procedure is cancelled.
   * @param delayMs Delay before starting to update the data.
   */
  private startDataRefreshSubscription(delayMs: number) {
    // If there is no API ID for getting the price, nothing is done.
    if (!this.PRICE_API_ID) {
      return;
    }

    if (this.priceSubscription) {
      this.priceSubscription.unsubscribe();
    }

    if (!environment.isInE2eMode) {
      this.ngZone.runOutsideAngular(() => {
        /* eslint-disable arrow-body-style */
        this.priceSubscription = of(0).pipe(delay(delayMs), mergeMap(() => {
          return this.http.get(`https://api.coinpaprika.com/v1/tickers/${this.PRICE_API_ID}?quotes=USD`);
        })).subscribe((response: any) => {
          this.ngZone.run(() => this.priceInternal.next(response.quotes.USD.price));
          this.startDataRefreshSubscription(this.updatePeriod);
        }, () => {
          this.startDataRefreshSubscription(this.errorUpdatePeriod);
        });
      });
    } else {
      // Set the price to 1 and stop making updates during e2e tests, to avoid potential
      // problems with the remote connection.
      this.priceInternal.next(1);
    }
  }
}
