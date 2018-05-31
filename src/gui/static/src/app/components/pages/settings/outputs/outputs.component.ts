import { Component, OnDestroy } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { ActivatedRoute } from '@angular/router';
import { ISubscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-outputs',
  templateUrl: './outputs.component.html',
  styleUrls: ['./outputs.component.scss'],
})
export class OutputsComponent implements OnDestroy {
  wallets: any[];

  private outputsSubscription: ISubscription;

  constructor(
    public walletService: WalletService,
    private route: ActivatedRoute,
  ) {
    route.queryParams.subscribe(params => this.loadData(params));
  }

  ngOnDestroy() {
    this.outputsSubscription.unsubscribe();
  }

  loadData(params) {
    const addr = params['addr'];

    this.outputsSubscription = this.walletService.outputsWithWallets().subscribe(wallets => {
      if (addr) {
        wallets = wallets.filter(wallet => {
          return wallet.addresses.find(address => address.address === addr);
        }).map(wallet => {
          wallet.addresses = wallet.addresses.filter(address => address.address === addr);

          return wallet;
        });
      }

      this.wallets = wallets;
    });
  }
}
