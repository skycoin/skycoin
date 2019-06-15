import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { PriceService } from '../../../services/price.service';
import { Subscription, ISubscription } from 'rxjs/Subscription';
import { WalletService } from '../../../services/wallet.service';
import { BlockchainService } from '../../../services/blockchain.service';
import { AppService } from '../../../services/app.service';
import 'rxjs/add/operator/skip';
import 'rxjs/add/operator/take';
import { BigNumber } from 'bignumber.js';
import { NetworkService } from '../../../services/network.service';
import { AppConfig } from '../../../app.config';

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
  hasPendingTxs: boolean;
  price: number;
  synchronized = true;
  walletDownloadUrl = AppConfig.walletDownloadUrl;

  private subscription: Subscription;
  private synchronizedSubscription: ISubscription;
  // This should be deleted. View the comment in the constructor.
  // private fetchVersionError: string;

  get loading() {
    return !this.current || !this.highest || this.current !== this.highest || !this.coins || this.coins === 'NaN' || !this.hours || this.hours === 'NaN';
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
    public networkService: NetworkService,
    private blockchainService: BlockchainService,
    private priceService: PriceService,
    private walletService: WalletService,
  ) { }

  ngOnInit() {
    this.subscription = this.blockchainService.progress
      .filter(response => !!response)
      .subscribe(response => {
        this.querying = false;
        this.highest = response.highest;
        this.current = response.current;
        this.percentage = this.current && this.highest ? (this.current / this.highest) : 0;

        // Adding the code here prevents the warning from flashing if the wallet is synchronized. Also, adding the
        // subscription to this.subscription causes problems.
        if (!this.synchronizedSubscription) {
          this.synchronizedSubscription = this.blockchainService.synchronized.subscribe(value => this.synchronized = value);
        }
      });

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
      this.hasPendingTxs = txs.user.length > 0;
    }));
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
    if (this.synchronizedSubscription) {
      this.synchronizedSubscription.unsubscribe();
    }
  }
}
