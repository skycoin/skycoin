import { Component, OnDestroy, OnInit } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { PriceService } from '../../../price.service';
import { Subscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-transaction-list',
  templateUrl: './transaction-list.component.html',
  styleUrls: ['./transaction-list.component.scss']
})
export class TransactionListComponent implements OnInit, OnDestroy {
  transactions: any[];

  private price: number;
  private priceSubscription: Subscription;

  constructor(
    private priceService: PriceService,
    private walletService: WalletService,
  ) { }

  ngOnInit() {
    this.priceSubscription = this.priceService.price.subscribe(price => this.price = price);
    this.walletService.transactions().subscribe(transactions => this.transactions = transactions);
  }

  ngOnDestroy() {
    this.priceSubscription.unsubscribe();
  }
}
