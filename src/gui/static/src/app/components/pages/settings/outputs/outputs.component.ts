import { Component, OnDestroy } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { ActivatedRoute } from '@angular/router';
import { SubscriptionLike } from 'rxjs';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { QrCodeComponent, QrDialogConfig } from '../../../layout/qr-code/qr-code.component';

@Component({
  selector: 'app-outputs',
  templateUrl: './outputs.component.html',
  styleUrls: ['./outputs.component.scss'],
})
export class OutputsComponent implements OnDestroy {
  wallets: any[]|null;

  private outputsSubscription: SubscriptionLike;
  private lastRouteParams: any;

  constructor(
    public walletService: WalletService,
    private route: ActivatedRoute,
    private dialog: MatDialog,
  ) {
    route.queryParams.subscribe(params => {
      this.wallets = null;
      this.lastRouteParams = params;
      this.walletService.startDataRefreshSubscription();
    });
    walletService.all().subscribe(() => this.loadData());
  }

  ngOnDestroy() {
    this.outputsSubscription.unsubscribe();
  }

  loadData() {
    const addr = this.lastRouteParams['addr'];

    this.outputsSubscription = this.walletService.outputsWithWallets().subscribe(wallets => {
      this.wallets = wallets
        .map(wallet => Object.assign({}, wallet))
        .map(wallet => {
          wallet.addresses = wallet.addresses.filter(address => {
            if (address.outputs.length > 0) {
              return addr ? address.address === addr : true;
            }
          });

          return wallet;
        })
        .filter(wallet => wallet.addresses.length > 0);
    });
  }

  showQrCode(event: any, address: string) {
    event.stopPropagation();

    const config: QrDialogConfig = { address };
    QrCodeComponent.openDialog(this.dialog, config);
  }
}
