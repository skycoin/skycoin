import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { PriceService } from '../../../price.service';
import { Subscription } from 'rxjs/Subscription';
import { WalletService } from '../../../services/wallet.service';
import { IntervalObservable } from 'rxjs/observable/IntervalObservable';
import { BlockchainService } from '../../../services/blockchain.service';

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

  private price: number;
  private priceSubscription: Subscription;
  private walletSubscription: Subscription;

  get balance() {
    if (this.price === null) return 'loading..';
    const balance = Math.round(this.coins * this.price * 100) / 100;
    return '$' + balance.toFixed(2) + ' ($' + (Math.round(this.price * 100) / 100) + ')';
  }

  get loading() {
    return !this.current || !this.highest || this.current != this.highest;
  }

  constructor(
    private blockchainService: BlockchainService,
    private priceService: PriceService,
    private walletService: WalletService,
  ) {}

  ngOnInit() {
    this.priceSubscription = this.priceService.price.subscribe(price => this.price = price);
    this.walletSubscription = this.walletService.all().subscribe(wallets => {
      this.coins = wallets.map(wallet => wallet.coins >= 0 ? wallet.coins : 0).reduce((a, b) => a + b, 0);
      this.hours = wallets.map(wallet => wallet.hours >= 0 ? wallet.hours : 0).reduce((a, b) => a + b, 0);
    });

    this.blockchainService.progress
      .filter(response => !!response)
      .subscribe(response => {
        this.highest = response.highest;
        this.current = response.current;
        this.percentage = this.current && this.highest ? (this.current / this.highest) : 0;
      });
  }

  ngOnDestroy() {
    this.priceSubscription.unsubscribe();
    this.walletSubscription.unsubscribe();
  }
}
