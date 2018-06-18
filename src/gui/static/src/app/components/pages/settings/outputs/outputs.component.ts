import { Component, OnDestroy } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { ActivatedRoute } from '@angular/router';
import { ISubscription } from 'rxjs/Subscription';
import { MatDialog, MatDialogConfig } from '@angular/material';
import { QrCodeComponent } from '../../../layout/qr-code/qr-code.component';

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
    private dialog: MatDialog,
  ) {
    route.queryParams.subscribe(params => this.loadData(params));
  }

  ngOnDestroy() {
    this.outputsSubscription.unsubscribe();
  }

  loadData(params) {
    const addr = params['addr'];

    this.outputsSubscription = this.walletService.outputsWithWallets().subscribe(wallets => {
      this.wallets = wallets
        .map(wallet => Object.assign({}, wallet))
        .map(wallet => {
          wallet.addresses = wallet.addresses.filter(address => {
            return addr ? address.address === addr : address.outputs.length > 0;
          });

          return wallet;
        })
        .filter(wallet => wallet.addresses.length > 0);
    });
  }

  showQrCode(event: any, address: string) {
    event.stopPropagation();

    const config = new MatDialogConfig();
    config.data = { address };
    this.dialog.open(QrCodeComponent, config);
  }
}
