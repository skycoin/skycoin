import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { PriceService } from '../../../price.service';
import { Subscription } from 'rxjs/Subscription';
import { WalletService } from '../../../services/wallet.service';
import { BlockchainService } from '../../../services/blockchain.service';
import { Observable } from 'rxjs/Observable';
import { ApiService } from '../../../services/api.service';
import { Http } from '@angular/http';
import { AppService } from '../../../services/app.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit, OnDestroy {
  @Input() title: string;
  @Input() coins: number;
  @Input() hours: number;

  current: number;
  highest: number;
  percentage: number;
  querying = true;
  version: string;
  releaseVersion: string;
  updateAvailable: boolean;

  private price: number;
  private priceSubscription: Subscription;
  private walletSubscription: Subscription;

  get balance() {
    if (this.price === null) { return 'loading..'; }
    const balance = Math.round(this.coins * this.price * 100) / 100;
    return '$' + balance.toFixed(2) + ' ($' + (Math.round(this.price * 100) / 100) + ')';
  }

  get loading() {
    return !this.current || !this.highest || this.current !== this.highest;
  }

  constructor(
    public appService: AppService,
    private apiService: ApiService,
    private blockchainService: BlockchainService,
    private priceService: PriceService,
    private walletService: WalletService,
    private http: Http,
  ) {}

  ngOnInit() {
    this.setVersion();
    this.priceSubscription = this.priceService.price.subscribe(price => this.price = price);
    this.walletSubscription = this.walletService.all().subscribe(wallets => {
      this.coins = wallets.map(wallet => wallet.coins >= 0 ? wallet.coins : 0).reduce((a, b) => a + b, 0);
      this.hours = wallets.map(wallet => wallet.hours >= 0 ? wallet.hours : 0).reduce((a, b) => a + b, 0);
    });

    this.blockchainService.progress
      .filter(response => !!response)
      .subscribe(response => {
        this.querying = false;
        this.highest = response.highest;
        this.current = response.current;
        this.percentage = this.current && this.highest ? (this.current / this.highest) : 0;
      });
  }

  ngOnDestroy() {
    this.priceSubscription.unsubscribe();
    this.walletSubscription.unsubscribe();
  }

  setVersion() {
    // Set build version
    this.apiService.getVersion().first()
      .subscribe(output =>  {
        this.version = output.version;
        this.retrieveReleaseVersion();
      });
  }

  private higherVersion(first: string, second: string): boolean {
    const fa = first.split('.');
    const fb = second.split('.');
    for (let i = 0; i < 3; i++) {
      const na = Number(fa[i]);
      const nb = Number(fb[i]);
      if (na > nb || !isNaN(na) && isNaN(nb)) {
        return true;
      } else if (na < nb || isNaN(na) && !isNaN(nb)) {
        return false;
      }
    }
    return false;
  }

  private retrieveReleaseVersion() {
    this.http.get('https://api.github.com/repos/skycoin/skycoin/tags')
      .map((res: any) => res.json())
      .catch((error: any) => Observable.throw(error || 'Unable to fetch latest release version from github.'))
      .subscribe(response =>  {
        this.releaseVersion = response.find(element => element['name'].indexOf('rc') === -1)['name'].substr(1);
        this.updateAvailable = this.higherVersion(this.releaseVersion, this.version);
      });
  }
}
