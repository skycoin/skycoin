import { Component, Input } from '@angular/core';
import { Wallet } from '../../../../app.datatypes';
import { WalletService } from '../../../../services/wallet.service';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { ChangeNameComponent } from '../change-name/change-name.component';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';

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
    if (this.wallet.encrypted) {
      this.dialog.open(PasswordDialogComponent).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.walletService.addAddress(this.wallet, passwordDialog.password)
            .subscribe(() => passwordDialog.close(), () => passwordDialog.error());
        });
    } else {
      this.walletService.addAddress(this.wallet).subscribe();
    }
  }

  toggleEmpty() {
    this.wallet.hideEmpty = !this.wallet.hideEmpty;
  }

  toggleEncryption() {
    const config = new MatDialogConfig();
    config.data = {
      confirm: !this.wallet.encrypted,
      title: this.wallet.encrypted ? 'Decrypt Wallet' : 'Encrypt Wallet',
    };

    if (!this.wallet.encrypted) {
      config.data['description'] = 'We suggest that you encrypt each one of your wallets with a password. ' +
        'If you forget your password, you can reset it with your seed. ' +
        'Make sure you have your seed saved somewhere safe before encrypting your wallet.';
    }

    this.dialog.open(PasswordDialogComponent, config).componentInstance.passwordSubmit
      .subscribe(passwordDialog => {
        this.walletService.toggleEncryption(this.wallet, passwordDialog.password).subscribe(() => {
          passwordDialog.close();
        }, e => passwordDialog.error(e));
      });
  }

  copyAddress(address) {
    const selBox = document.createElement('textarea');

    selBox.style.position = 'fixed';
    selBox.style.left = '0';
    selBox.style.top = '0';
    selBox.style.opacity = '0';
    selBox.value = address.address;

    document.body.appendChild(selBox);
    selBox.focus();
    selBox.select();

    document.execCommand('copy');
    document.body.removeChild(selBox);

    address.copying = true;

    // wait for a while and then remove the 'copying' class
    setTimeout(function () {
      address.copying = false;
    }, 500);

  }
}
