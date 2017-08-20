import { Component } from '@angular/core';
import { WalletService } from './services/wallet.service';
import { BlockchainService } from './services/blockchain.service';
import { IntervalObservable } from 'rxjs/observable/IntervalObservable';
import 'rxjs/add/operator/takeWhile';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {

  current: number;
  highest: number;
  percentage: number;

  constructor(
    public walletService: WalletService,
    private blockchainService: BlockchainService,
  ) {}

  ngOnInit() {
    IntervalObservable
      .create(3000)
      .flatMap(() => this.blockchainService.progress())
      .takeWhile(response => !response.current || response.current !== response.highest)
      .subscribe(response => {
        this.highest = response.highest;
        this.current = response.current;
        this.percentage = this.current && this.highest ? (this.current / this.highest * 100) : 0;
        console.log(response);
      }, error => console.log(error),
        () => this.completeLoading());
  }

  loading() {
    return !this.current || !this.highest || this.current != this.highest;
  }

  private completeLoading() {
    this.current = 999999999999;
    this.highest = 999999999999;
    this.walletService.refreshBalances();
  }
}
