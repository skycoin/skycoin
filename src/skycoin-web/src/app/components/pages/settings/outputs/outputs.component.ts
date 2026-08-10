import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Subscription, mergeMap } from 'rxjs';

import { SpendingService } from '../../../../services/wallet/spending.service';
import { Wallet } from '../../../../app.datatypes';
import { BaseCoin } from '../../../../coins/basecoin';
import { CoinService } from '../../../../services/coin.service';
import { openQrModal } from '../../../../utils';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';

@Component({
    selector: 'app-outputs',
    templateUrl: './outputs.component.html',
    styleUrls: ['./outputs.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class OutputsComponent implements OnInit, OnDestroy {
  wallets!: Wallet[] | null;
  currentCoin!: BaseCoin;
  showError = false;

  private subscription!: Subscription;
  private dataSubscription!: Subscription;
  private urlParams!: Params;

  constructor(
    private route: ActivatedRoute,
    private spendingService: SpendingService,
    private dialog: CustomMatDialogService,
    private coinService: CoinService,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit() {
    this.subscription = this.route.queryParams.pipe(mergeMap(params => {
      this.urlParams = params;
      return this.coinService.currentCoin;
    })).subscribe((coin: BaseCoin) => {
      this.wallets = null;
      this.currentCoin = coin;
      this.getWalletsOutputs();
      this.changeDetectorRef.markForCheck();
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
    this.closeDataSubscription();
  }

  showQr(event: any, address: any) {
    event.stopPropagation();
    openQrModal(this.dialog, address);
  }

  private getWalletsOutputs() {
    const address = this.urlParams['addr'];
    this.showError = false;

    this.closeDataSubscription();
    this.dataSubscription = this.spendingService.outputsWithWallets().subscribe(wallets => {
      if (wallets.length === 0) {
        this.wallets = [];
        return;
      }

      this.wallets = address
        ? this.getOutputsForSpecificAddress(wallets, address)
        : this.getOutputs(wallets);
      this.changeDetectorRef.markForCheck();
    },
    () => {
      this.showError = true;
      this.changeDetectorRef.markForCheck();
    });
  }

  private getOutputsForSpecificAddress(wallets: any, address: string) {
    const filteredWallets: Wallet[]  = wallets.filter((wallet: any) => {
      return wallet.addresses.find((addr: any) => {
        return addr.address === address;
      });
    }).map((wallet: any) => {
      return Object.assign({}, wallet);
    });

    return filteredWallets.map(wallet => {
      wallet.addresses = wallet.addresses.filter(addr => addr.address === address);
      return wallet;
    });
  }

  private getOutputs(wallets: any) {
    const copiedWallets = wallets.map((wallet: any) => Object.assign({}, wallet));

    return copiedWallets.filter((wallet: any) => {
      wallet.addresses = wallet.addresses.filter((addr: any) => addr.outputs.length > 0);
      return wallet.addresses.length > 0;
    });
  }

  private closeDataSubscription() {
    if (this.dataSubscription && !this.dataSubscription.closed) {
      this.dataSubscription.unsubscribe();
    }
  }
}
