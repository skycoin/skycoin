import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

import { BaseCoin } from '../coins/basecoin';
import { SkycoinCoin } from '../coins/skycoin.coin';
import { TestCoin } from '../coins/test.coin';
import { defaultCoinId } from '../constants/coins-id.const';
import { environment } from '../../environments/environment';
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

  private readonly correntCoinStorageKey = 'currentCoin';
  private readonly nodeUrlsStorageKey = 'nodeUrls';

  constructor(private translate: TranslateService) {
    this.loadAvailableCoins();
    this.loadNodeUrls();
    this.loadCurrentCoin();
    sessionStorage.setItem(this.correntCoinStorageKey, this.currentCoin.getValue().id.toString());
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
    const storedCoinId = sessionStorage.getItem(this.correntCoinStorageKey) || localStorage.getItem(this.correntCoinStorageKey);
    const coinId = storedCoinId ? +storedCoinId : defaultCoinId;
    const coin = this.coins.find((c: BaseCoin) => c.id === coinId);
    this.currentCoin.next(coin);
  }

  private saveCoin(coinId: number) {
    localStorage.setItem(this.correntCoinStorageKey, coinId.toString());
    sessionStorage.setItem(this.correntCoinStorageKey, coinId.toString());
  }

  private loadAvailableCoins() {
    this.coins.push(new SkycoinCoin());

    if (!environment.production) {
      this.coins.push(new TestCoin());
    }

    const IDs = new Map<number, boolean>();
    this.coins.forEach((value: BaseCoin) => {
      if (IDs[value.id]) {
        throw new Error('More than one coin with the same ID');
      }
      IDs[value.id] = true;
    });
  }
}
