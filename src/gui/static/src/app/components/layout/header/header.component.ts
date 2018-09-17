import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { PriceService } from '../../../services/price.service';
import { Subscription } from 'rxjs/Subscription';
import { WalletService } from '../../../services/wallet.service';
import { BlockchainService } from '../../../services/blockchain.service';
import { Observable } from 'rxjs/Observable';
import { ApiService } from '../../../services/api.service';
import { Http } from '@angular/http';
import { AppService } from '../../../services/app.service';
import 'rxjs/add/operator/skip';
import 'rxjs/add/operator/take';
import { shouldUpgradeVersion } from '../../../utils/semver';
import { TranslateService } from '@ngx-translate/core';
import { BigNumber } from 'bignumber.js';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
})
export class HeaderComponent implements OnInit, OnDestroy {
  @Input() headline: string;

  addresses = [];
  current: number;
  highest: number;
  percentage: number;
  querying = true;
  version: string;
  releaseVersion: string;
  updateAvailable: boolean;
  hasPendingTxs: boolean;
  price: number;

  private subscription: Subscription;
  private fetchVersionError: string;

  get loading() {
    return !this.current || !this.highest || this.current !== this.highest;
  }

  get coins() {
    let coins = new BigNumber('0');
    this.addresses.map(addr => coins = coins.plus(addr.coins));

    return coins.decimalPlaces(6).toString();
  }

  get hours() {
    let hours = new BigNumber('0');
    this.addresses.map(addr => hours = hours.plus(addr.hours));

    return hours.decimalPlaces(0).toString();
  }

  constructor(
    public appService: AppService,
    private apiService: ApiService,
    private blockchainService: BlockchainService,
    private priceService: PriceService,
    private walletService: WalletService,
    private http: Http,
    private translateService: TranslateService,
  ) {
    this.translateService.get('errors.fetch-version').subscribe(msg => {
      this.fetchVersionError = msg;
    });
  }

  ngOnInit() {
    this.subscription = this.blockchainService.progress
      .filter(response => !!response)
      .subscribe(response => {
        this.querying = false;
        this.highest = response.highest;
        this.current = response.current;
        this.percentage = this.current && this.highest ? (this.current / this.highest) : 0;
      });

    this.setVersion();

    this.subscription.add(this.priceService.price.subscribe(price => this.price = price));

    this.subscription.add(this.walletService.allAddresses().subscribe(addresses => {
      this.addresses = addresses.reduce((array, item) => {
        if (!array.find(addr => addr.address === item.address)) {
          array.push(item);
        }

        return array;
      }, []);
    }));

    this.subscription.add(this.walletService.pendingTransactions().subscribe(txs => {
      this.hasPendingTxs = txs.length > 0;
    }));
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  setVersion() {
    // Set build version
    setTimeout(() => {
      this.apiService.getVersion().first()
        .subscribe(output =>  {
          this.version = output.version;
          this.retrieveReleaseVersion();
        });
    }, 1000);
  }

  private retrieveReleaseVersion() {
    this.http.get('https://api.github.com/repos/skycoin/skycoin/tags')
      .map((res: any) => res.json())
      .catch((error: any) => Observable.throw(error || this.fetchVersionError))
      .subscribe(response =>  {
        this.releaseVersion = response.find(element => element['name'].indexOf('rc') === -1)['name'].substr(1);
        this.updateAvailable = shouldUpgradeVersion(this.version, this.releaseVersion);
      });
  }
}
