import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { PriceService } from '../../../price.service';
import { Subscription } from 'rxjs/Subscription';
import { WalletService } from '../../../services/wallet.service';
import { MdDialog } from '@angular/material';
import { Router } from '@angular/router';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit, OnDestroy {
  @Input() title: string;
  @Input() coins: number;
  @Input() hours: number;

  private price: number;
  private priceSubscription: Subscription;
  private walletSubscription: Subscription;

  get balance() {
    if (this.price === null) return 'loading..';
    const balance = Math.round(this.coins * this.price * 100) / 100;
    return '$' + balance.toFixed(2) + ' ($' + (Math.round(this.price * 100) / 100) + ')';
  }

  constructor(
    private priceService: PriceService,
    private walletService: WalletService,
  ) {}

  ngOnInit() {
    this.priceSubscription = this.priceService.price.subscribe(price => this.price = price);
    this.walletSubscription = this.walletService.all().subscribe(wallets => {
      this.coins = wallets.map(wallet => wallet.coins >= 0 ? wallet.coins : 0).reduce((a , b) => a + b, 0);
      this.hours = wallets.map(wallet => wallet.hours >= 0 ? wallet.hours : 0).reduce((a , b) => a + b, 0);
    })
  }

  ngOnDestroy() {
    this.priceSubscription.unsubscribe();
    this.walletSubscription.unsubscribe();
  }
}
