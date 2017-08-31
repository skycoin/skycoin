import { Component, Input } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { WalletModel } from '../../../../models/wallet.model';
import { MdDialog, MdDialogConfig } from '@angular/material';
import { ChangeNameComponent } from '../change-name/change-name.component';
import { QrCodeComponent } from '../../../layout/qr-code/qr-code.component';

@Component({
  selector: 'app-wallet-detail',
  templateUrl: './wallet-detail.component.html',
  styleUrls: ['./wallet-detail.component.css']
})
export class WalletDetailComponent {
  @Input() wallet: WalletModel;

  constructor(
    public walletService: WalletService,
    private dialog: MdDialog,
  ) {}

  addAddress() {
    this.walletService.addAddress(this.wallet).subscribe(output => this.wallet.entries.push(output));
  }

  renameWallet() {
    const config = new MdDialogConfig();
    config.width = '500px';
    config.data = this.wallet;
    this.dialog.open(ChangeNameComponent, config).afterClosed().subscribe(result => {
      if (result) {
        this.wallet.meta.label = result;
      }
    });
  }

  showQr(address) {
    const config = new MdDialogConfig();
    config.data = address;
    this.dialog.open(QrCodeComponent, config);
  }
}
