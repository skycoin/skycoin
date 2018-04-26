import { Component } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-outputs',
  templateUrl: './outputs.component.html',
  styleUrls: ['./outputs.component.scss']
})
export class OutputsComponent {

  wallets: any[];

  constructor(
    public walletService: WalletService,
    private route: ActivatedRoute,
  ) {
    route.queryParams.subscribe(params => this.loadData(params));
  }

  loadData(params) {
    const addr = params['addr'];

    this.walletService.outputsWithWallets().subscribe(wallets => {
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
