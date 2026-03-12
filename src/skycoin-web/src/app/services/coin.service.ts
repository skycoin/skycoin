import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, ReplaySubject } from 'rxjs';

import { BaseCoin } from '../coins/basecoin';
import { SkycoinCoin } from '../coins/skycoin.coin';
import { TranslateService } from '@ngx-translate/core';

export enum TemporarilyAllowCoinResult {
  OK = 0,
  AlreadyInUse = 1,
  Cancelled = 2,
}

@Injectable()
export class CoinService {

  currentCoin: BehaviorSubject<BaseCoin> = new BehaviorSubject<BaseCoin>(null);
  coins: BaseCoin[] = [];
  customNodeUrls: object;

  /**
   * Emits true once coins have been loaded (from server or fallback).
   */
  coinsLoaded: ReplaySubject<boolean> = new ReplaySubject<boolean>(1);

  private readonly currentCoinStorageKey = 'currentCoin';
  private readonly nodeUrlsStorageKey = 'nodeUrls';

  constructor(
    private translate: TranslateService,
    private http: HttpClient
  ) {
    this.loadNodeUrls();
    this.loadCoinsFromServer();
  }

  changeCoin(coin: BaseCoin) {
    if (coin.id !== this.currentCoin.value.id) {
      this.currentCoin.next(coin);
      this.saveCoin(coin.id);
    }
  }

  temporarilyAllowCoin(id: number, url: string): TemporarilyAllowCoinResult {
    if (window['isElectron']) {
      const params = {
        id: id,
        url: url,
        confirmationTitle: this.translate.instant('nodes.change.confirmation.title'),
        confirmationText: this.translate.instant('nodes.change.confirmation.text', {url: url}),
        confirmationOk: this.translate.instant('nodes.change.confirmation.ok'),
        confirmationCancel: this.translate.instant('nodes.change.confirmation.cancel'),
      };
      return window['ipcRenderer'].sendSync('temporarilyAllowCoinSync', params);
    } else {
      return TemporarilyAllowCoinResult.OK;
    }
  }

  removeTemporarilyAllowedCoin() {
    if (window['isElectron']) {
      window['ipcRenderer'].sendSync('removeTemporarilyAllowedCoinSync');
    }
  }

  changeNodeUrl(coinId: number, url: string) {
    if (!this.coins.find(coin => coin.id === coinId)) {
      return;
    }

    if (url.length > 0) {
      this.customNodeUrls[coinId.toString()] = url;
    } else {
      delete this.customNodeUrls[coinId.toString()];
    }

    if (!window['isElectron']) {
      localStorage.setItem(this.nodeUrlsStorageKey, JSON.stringify(this.customNodeUrls));
    } else {
      if (url !== '') {
        window['ipcRenderer'].sendSync('acceptTemporarilyAllowedCoinSync');
      } else {
        window['ipcRenderer'].sendSync('removeAllowedCoinSync', coinId);
      }
    }

    if (coinId === this.currentCoin.value.id) {
      this.currentCoin.next(this.currentCoin.value);
    }
  }

  private loadCoinsFromServer() {
    this.http.get('/api/v1/coins').subscribe(
      (serverCoins: any[]) => {
        if (serverCoins && serverCoins.length > 0) {
          this.coins = serverCoins.map(data => BaseCoin.fromServerData(data));
        } else {
          this.loadFallbackCoins();
        }
        this.finishCoinLoading();
      },
      () => {
        // Server not available — use fallback hardcoded coin
        this.loadFallbackCoins();
        this.finishCoinLoading();
      }
    );
  }

  private loadFallbackCoins() {
    this.coins = [new SkycoinCoin()];
  }

  private finishCoinLoading() {
    this.validateCoinIds();
    this.loadCurrentCoin();
    this.coinsLoaded.next(true);
  }

  private validateCoinIds() {
    const IDs = new Map<number, boolean>();
    this.coins.forEach((value: BaseCoin) => {
      if (IDs.has(value.id)) {
        throw new Error('More than one coin with the same ID');
      }
      IDs.set(value.id, true);
    });
  }

  private loadNodeUrls() {
    if (!window['isElectron']) {
      const savedUrls: object = JSON.parse(localStorage.getItem(this.nodeUrlsStorageKey));
      this.customNodeUrls = savedUrls ? savedUrls : {};
    } else {
      const savedUrls = window['ipcRenderer'].sendSync('loadNodeUrlsSync');
      this.customNodeUrls = savedUrls ? JSON.parse(savedUrls) : {};
    }
  }

  private loadCurrentCoin() {
    const storedCoinId = sessionStorage.getItem(this.currentCoinStorageKey) || localStorage.getItem(this.currentCoinStorageKey);
    let coin: BaseCoin;

    if (storedCoinId) {
      coin = this.coins.find((c: BaseCoin) => c.id === +storedCoinId);
    }

    // Fall back to first available coin
    if (!coin && this.coins.length > 0) {
      coin = this.coins[0];
    }

    if (coin) {
      this.currentCoin.next(coin);
      sessionStorage.setItem(this.currentCoinStorageKey, coin.id.toString());
    }
  }

  private saveCoin(coinId: number) {
    localStorage.setItem(this.currentCoinStorageKey, coinId.toString());
    sessionStorage.setItem(this.currentCoinStorageKey, coinId.toString());
  }
}
