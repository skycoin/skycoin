import { Component, OnDestroy, OnInit } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { ActivatedRoute } from '@angular/router';
import { Wallet } from '../../../app.datatypes';
import { Subscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-wallet-show',
  templateUrl: './wallet-show.component.html',
  styleUrls: ['./wallet-show.component.scss']
})
export class WalletShowComponent implements OnInit, OnDestroy {
  wallet: Wallet;
  private walletSubscription: Subscription;

  constructor(
    private route: ActivatedRoute,
    private walletService: WalletService,
  ) { }

  ngOnInit() {
    this.walletSubscription = this.route.params.switchMap(params => this.walletService.find(params.filename))
      .subscribe(wallet => this.wallet = wallet);
  }

  ngOnDestroy() {
    this.walletSubscription.unsubscribe();
  }

  newAddress() {
    this.walletService.addAddress(this.wallet).subscribe(() => console.log('completed'));
  }

  toggleEmpty() {
    this.wallet.hideEmpty = !this.wallet.hideEmpty;
  }
}
