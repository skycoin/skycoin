import { Component, Input } from '@angular/core';
import { Wallet } from '../../../../app.datatypes';
import { WalletService } from '../../../../services/wallet.service';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { ChangeNameComponent } from '../change-name/change-name.component';

@Component({
  selector: 'app-wallet-detail',
  templateUrl: './wallet-detail.component.html',
  styleUrls: ['./wallet-detail.component.scss']
})
export class WalletDetailComponent {
  @Input() wallet: Wallet;

  constructor(
    private dialog: MatDialog,
    private walletService: WalletService,
  ) { }

  editWallet() {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.data = this.wallet;
    this.dialog.open(ChangeNameComponent, config);
  }

  newAddress() {
    this.walletService.addAddress(this.wallet).subscribe();
  }

  toggleEmpty() {
    this.wallet.hideEmpty = !this.wallet.hideEmpty;
  }

  copyAddress(address, i) {
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
    d[i].classList.toggle('copying');

    // wait for a while and then remove the 'copying' class
    setTimeout(function () {
      d[i].classList.toggle('copying');
    }, 500);

  }

  toggleClass(i) {
    const label = document.getElementsByClassName('copy-label');
    label[i].classList.toggle('hidden');
  }
}
