import { Injectable, NgZone } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, of, Subscription, Observable } from 'rxjs';
import { delay, mergeMap } from 'rxjs/operators';

import { CoinService } from './coin.service';
import { BaseCoin } from '../coins/basecoin';

@Injectable()
export class PriceService {

  price = new BehaviorSubject<number>(null);

  private readonly updatePeriod = 10 * 60 * 1000;
  private readonly errorUpdatePeriod = 30 * 1000;
  private priceTickerId: string | null = null;
  private priceTickerSource: string = 'coinpaprika';
  private priceSubscription: Subscription;

  constructor(
    private http: HttpClient,
    private ngZone: NgZone,
    private coinService: CoinService
  ) {
    this.coinService.currentCoin.subscribe((coin: BaseCoin) => {
      if (!coin) {
        return;
      }
      this.priceTickerId = coin.priceTickerId;
      this.priceTickerSource = coin.priceTickerSource || 'coinpaprika';
      this.loadConfigAndStart();
    });
  }

  private loadConfigAndStart() {
    // Price ticker config is already provided by the coin discovery.
    // Just start the refresh subscription.
    this.startDataRefreshSubscription(0);
  }

  private startDataRefreshSubscription(delayMs: number) {
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
        this.ngZone.run(() => this.price.next(price));
        this.startDataRefreshSubscription(this.updatePeriod);
      }, () => {
        this.startDataRefreshSubscription(this.errorUpdatePeriod);
      });
    });
  }

  private fetchPrice(): Observable<number> {
    if (this.priceTickerSource === 'coingecko') {
      return this.http.get(
        `https://api.coingecko.com/api/v3/simple/price?ids=${this.priceTickerId}&vs_currencies=usd`
      ).pipe(mergeMap((response: any) => {
        return of(response[this.priceTickerId].usd);
      }));
    } else {
      return this.http.get(
        `https://api.coinpaprika.com/v1/tickers/${this.priceTickerId}?quotes=USD`
      ).pipe(mergeMap((response: any) => {
        return of(response.quotes.USD.price);
      }));
    }
  }
}
