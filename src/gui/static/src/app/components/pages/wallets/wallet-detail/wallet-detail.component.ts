import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { Wallet } from '../../../../app.datatypes';
import { WalletService } from '../../../../services/wallet.service';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { ChangeNameComponent } from '../change-name/change-name.component';
import { QrCodeComponent } from '../../../layout/qr-code/qr-code.component';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { MatSnackBar } from '@angular/material';
import { showSnackbarError } from '../../../../utils/errors';
import { TranslateService } from '@ngx-translate/core';
import { NumberOfAddressesComponent } from '../number-of-addresses/number-of-addresses';

@Component({
  selector: 'app-wallet-detail',
  templateUrl: './wallet-detail.component.html',
  styleUrls: ['./wallet-detail.component.scss'],
})
export class WalletDetailComponent implements OnInit, OnDestroy {
  @Input() wallet: Wallet;

  private encryptionWarning: string;
  private HowManyAddresses: number;

  constructor(
    private dialog: MatDialog,
    private walletService: WalletService,
    private snackbar: MatSnackBar,
    private translateService: TranslateService,
  ) { }

  ngOnInit() {
    this.translateService.get('wallet.new.encrypt-warning').subscribe(msg => {
      this.encryptionWarning = msg;
    });
  }

  ngOnDestroy() {
    this.snackbar.dismiss();
  }

  editWallet() {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.data = this.wallet;
    this.dialog.open(ChangeNameComponent, config);
  }

  newAddress() {
    this.snackbar.dismiss();

    const config = new MatDialogConfig();
    config.width = '566px';
    this.dialog.open(NumberOfAddressesComponent, config).afterClosed()
      .subscribe(response => {
        if (response) {
          this.HowManyAddresses = response;
          this.continueNewAddress();
        }
      });
  }

  toggleEmpty() {
    this.wallet.hideEmpty = !this.wallet.hideEmpty;
  }

  toggleEncryption() {
    const config = new MatDialogConfig();
    config.data = {
      confirm: !this.wallet.encrypted,
      title: this.wallet.encrypted ? 'wallet.decrypt' : 'wallet.encrypt',
    };

    if (!this.wallet.encrypted) {
      config.data['description'] = this.encryptionWarning;
    } else {
      config.data['wallet'] = this.wallet;
    }

    this.dialog.open(PasswordDialogComponent, config).componentInstance.passwordSubmit
      .subscribe(passwordDialog => {
        this.walletService.toggleEncryption(this.wallet, passwordDialog.password).subscribe(() => {
          passwordDialog.close();
        }, e => passwordDialog.error(e));
      });
  }

  copyAddress(event, address, duration = 500) {
    event.stopPropagation();

    if (address.copying) {
      return;
    }

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

    setTimeout(function() {
      address.copying = false;
    }, duration);
  }

  showQrCode(event, address: string) {
    event.stopPropagation();

    const config = new MatDialogConfig();
    config.data = { address };
    this.dialog.open(QrCodeComponent, config);
  }

  private continueNewAddress() {
    if (this.wallet.encrypted) {
      const config = new MatDialogConfig();
      config.data = {
        wallet: this.wallet,
      };

      this.dialog.open(PasswordDialogComponent, config).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.walletService.addAddress(this.wallet, this.HowManyAddresses, passwordDialog.password)
            .subscribe(() => passwordDialog.close(), () => passwordDialog.error());
        });
    } else {
      this.walletService.addAddress(this.wallet, this.HowManyAddresses).subscribe(null, err => showSnackbarError(this.snackbar, err));
    }
  }
}
