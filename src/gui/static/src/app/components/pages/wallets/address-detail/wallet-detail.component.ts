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


  copyAddress(address) {
    const selBox = document.createElement('textarea');

    selBox.style.position = 'fixed';
    selBox.style.left = '0';
    selBox.style.top = '0';
    selBox.style.opacity = '0';
    selBox.value = address;

    document.body.appendChild(selBox);
    selBox.focus();
    selBox.select();

    document.execCommand('copy');
    document.body.removeChild(selBox);

    const d = document.getElementsByClassName('click-to-copy');
    d[0].classList.toggle('copying');

    // wait for a while and then remove the 'copying' class
    setTimeout(function () {
      d[0].classList.toggle('copying');
    }, 500);

  }

  toggleClass() {
    const label = document.getElementsByClassName('copy-label');
    label[0].classList.toggle('hidden');
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
