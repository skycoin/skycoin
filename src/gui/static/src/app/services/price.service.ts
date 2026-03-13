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

  private priceTickerId: string = null;
  private priceTickerSource: string = null;
  private configLoaded = false;

  constructor(
    private http: HttpClient,
    private ngZone: NgZone,
  ) {
    this.loadConfigAndStart();
  }

  /**
   * Loads price ticker config from the health API, then starts price updates.
   */
  private loadConfigAndStart() {
    if (environment.isInE2eMode) {
      this.priceInternal.next(1);
      return;
    }

    this.http.get('/api/v1/health').subscribe((response: any) => {
      if (response.fiber && response.fiber.price_ticker_id) {
        this.priceTickerId = response.fiber.price_ticker_id;
        this.priceTickerSource = response.fiber.price_ticker_source || 'coinpaprika';
      } else {
        // Fallback to hardcoded AppConfig for backwards compatibility
        this.priceTickerId = AppConfig.priceApiId;
        this.priceTickerSource = 'coinpaprika';
      }
      this.configLoaded = true;
      this.startDataRefreshSubscription(0);
    }, () => {
      // If health API fails, fall back to AppConfig
      this.priceTickerId = AppConfig.priceApiId;
      this.priceTickerSource = 'coinpaprika';
      this.configLoaded = true;
      this.startDataRefreshSubscription(0);
    });
  }

  /**
   * Makes the service start updating the data periodically. If this function was called
   * before, the previous updating procedure is cancelled.
   * @param delayMs Delay before starting to update the data.
   */
  private startDataRefreshSubscription(delayMs: number) {
    // If there is no API ID for getting the price, nothing is done.
    if (!this.priceTickerId) {
      return;
    }

    if (this.priceSubscription) {
      this.priceSubscription.unsubscribe();
    }

    this.ngZone.runOutsideAngular(() => {
      this.priceSubscription = of(0).pipe(delay(delayMs), mergeMap(() => {
        return this.fetchPrice();
      })).subscribe((price: number) => {
        this.ngZone.run(() => this.priceInternal.next(price));
        this.startDataRefreshSubscription(this.updatePeriod);
      }, () => {
        this.startDataRefreshSubscription(this.errorUpdatePeriod);
      });
    });
  }

  /**
   * Fetches the current price from the configured API source.
   */
  private fetchPrice(): Observable<number> {
    if (this.priceTickerSource === 'coingecko') {
      return this.http.get(
        `https://api.coingecko.com/api/v3/simple/price?ids=${this.priceTickerId}&vs_currencies=usd`
      ).pipe(mergeMap((response: any) => {
        return of(response[this.priceTickerId].usd);
      }));
    } else {
      // Default: coinpaprika
      return this.http.get(
        `https://api.coinpaprika.com/v1/tickers/${this.priceTickerId}?quotes=USD`
      ).pipe(mergeMap((response: any) => {
        return of(response.quotes.USD.price);
      }));
    }
  }
}
