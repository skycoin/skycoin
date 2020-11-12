import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { SubscriptionLike } from 'rxjs';

import { BalanceAndOutputsService } from '../../../../services/wallet-operations/balance-and-outputs.service';
import { WalletWithOutputs } from '../../../../services/wallet-operations/wallet-objects';

/**
 * Allows to see the list of unspent outputs of the registered wallets. The list can be
 * limited to one address by setting the"addr" param, on the URL, to the desired address.
 */
@Component({
  selector: 'app-outputs',
  templateUrl: './outputs.component.html',
  styleUrls: ['./outputs.component.scss'],
})
export class OutputsComponent implements OnDestroy {
  wallets: WalletWithOutputs[]|null;

  private outputsSubscription: SubscriptionLike;

  constructor(
    route: ActivatedRoute,
    private balanceAndOutputsService: BalanceAndOutputsService,
  ) {
    // Reload the data every time the url params change.
    route.queryParams.subscribe(params => {
      this.wallets = null;
      this.loadData(params);
    });
  }

  ngOnDestroy() {
    this.removeOutputsSubscription();
  }

  private loadData(lastRouteParams: Params) {
    const addr = lastRouteParams['addr'];

    this.removeOutputsSubscription();

    // Periodically get the list of wallets with the outputs.
    this.outputsSubscription = this.balanceAndOutputsService.outputsWithWallets.subscribe(wallets => {
      // The original response object is modified. No copy is created before doing this
      // because the data is only used by this page.
      this.wallets = wallets.map(wallet => {
        // Include only addresses with outputs or the requested address.
        wallet.addresses = wallet.addresses.filter(address => {
          if (address.outputs.length > 0) {
            return addr ? address.address === addr : true;
          }
        });

        return wallet;
      }).filter(wallet => wallet.addresses.length > 0);
    });
  }

  private removeOutputsSubscription() {
    if (this.outputsSubscription) {
      this.outputsSubscription.unsubscribe();
    }
  }
}
