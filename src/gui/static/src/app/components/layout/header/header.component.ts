import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { PriceService } from '../../../price.service';
import { Subscription } from 'rxjs/Subscription';

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

  get balance() {
    if (this.price === null) return 'loading..';
    const balance = Math.round(this.coins * this.price * 100) / 100;
    return '$' + balance.toFixed(2) + ' ($' + (Math.round(this.price * 100) / 100) + ')';
  }

  constructor(
    private priceService: PriceService,
  ) {}

  ngOnInit() {
    this.priceSubscription = this.priceService.price.subscribe(price => this.price = price);
  }

  ngOnDestroy() {
    this.priceSubscription.unsubscribe();
  }

  get showLargeHeader(): boolean {
    return this.coins >= 0;
  }
}
