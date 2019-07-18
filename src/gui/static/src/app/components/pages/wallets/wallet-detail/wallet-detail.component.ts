import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { Wallet, ConfirmationData } from '../../../../app.datatypes';
import { WalletService } from '../../../../services/wallet.service';
import { MatDialog, MatDialogConfig, MatDialogRef } from '@angular/material/dialog';
import { ChangeNameComponent, ChangeNameData } from '../change-name/change-name.component';
import { QrCodeComponent, QrDialogConfig } from '../../../layout/qr-code/qr-code.component';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { getHardwareWalletErrorMsg } from '../../../../utils/errors';
import { NumberOfAddressesComponent } from '../number-of-addresses/number-of-addresses';
import { TranslateService } from '@ngx-translate/core';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { Observable } from 'rxjs/Observable';
import { showConfirmationModal, copyTextToClipboard } from '../../../../utils';
import { AppConfig } from '../../../../app.config';
import { Router } from '@angular/router';
import { HwConfirmAddressDialogComponent, AddressConfirmationParams } from '../../../layout/hardware-wallet/hw-confirm-address-dialog/hw-confirm-address-dialog.component';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { ISubscription } from 'rxjs/Subscription';
import { ApiService } from '../../../../services/api.service';

@Component({
  selector: 'app-wallet-detail',
  templateUrl: './wallet-detail.component.html',
  styleUrls: ['./wallet-detail.component.scss'],
})
export class WalletDetailComponent implements OnDestroy {
  @Input() wallet: Wallet;

  confirmingIndex = null;
  creatingAddress = false;
  preparingToEdit = false;

  private howManyAddresses: number;
  private editSubscription: ISubscription;
  private confirmSubscription: ISubscription;
  private txHistorySubscription: ISubscription;

  constructor(
    private dialog: MatDialog,
    private walletService: WalletService,
    private msgBarService: MsgBarService,
    private hwWalletService: HwWalletService,
    private translateService: TranslateService,
    private router: Router,
    private apiService: ApiService,
  ) { }

  ngOnDestroy() {
    this.msgBarService.hide();
    if (this.editSubscription) {
      this.editSubscription.unsubscribe();
    }
    if (this.confirmSubscription) {
      this.confirmSubscription.unsubscribe();
    }
    if (this.txHistorySubscription) {
      this.txHistorySubscription.unsubscribe();
    }
  }

  editWallet() {
    this.msgBarService.hide();

    if (this.wallet.isHardware) {
      if (this.preparingToEdit) {
        return;
      }

      this.preparingToEdit = true;
      this.editSubscription = this.hwWalletService.checkIfCorrectHwConnected(this.wallet.addresses[0].address)
        .flatMap(() => this.walletService.getHwFeaturesAndUpdateData(this.wallet))
        .subscribe(
          () => {
            this.continueEditWallet();
            this.preparingToEdit = false;
          },
          err => {
            this.msgBarService.showError(getHardwareWalletErrorMsg(this.translateService, err));
            this.preparingToEdit = false;
          },
        );
    } else {
      this.continueEditWallet();
    }
  }

  newAddress() {
    if (this.creatingAddress) {
      return;
    }

    if (this.wallet.isHardware && this.wallet.addresses.length >= AppConfig.maxHardwareWalletAddresses) {
      const confirmationData: ConfirmationData = {
        text: 'wallet.max-hardware-wallets-error',
        headerText: 'errors.error',
        confirmButtonText: 'confirmation.close',
      };
      showConfirmationModal(this.dialog, confirmationData);

      return;
    }

    this.msgBarService.hide();

    if (!this.wallet.isHardware) {
      const maxAddressesGap = 20;

      const config = new MatDialogConfig();
      config.width = '566px';
      config.data = (howManyAddresses, callback) => {
        this.howManyAddresses = howManyAddresses;

        let lastWithBalance = 0;
        this.wallet.addresses.forEach((address, i) => {
          if (address.coins.isGreaterThan(0)) {
            lastWithBalance = i;
          }
        });

        if ((this.wallet.addresses.length - (lastWithBalance + 1)) + howManyAddresses < maxAddressesGap) {
          callback(true);
          this.continueNewAddress();
        } else {
          this.txHistorySubscription = this.apiService.getTransactions(this.wallet.addresses).first().subscribe(transactions => {
            const AddressesWithTxs = new Map<string, boolean>();

            transactions.forEach(transaction => {
              transaction.outputs.forEach(output => {
                if (!AddressesWithTxs.has(output.dst)) {
                  AddressesWithTxs.set(output.dst, true);
                }
              });
            });

            let lastWithTxs = 0;
            this.wallet.addresses.forEach((address, i) => {
              if (AddressesWithTxs.has(address.address)) {
                lastWithTxs = i;
              }
            });

            if ((this.wallet.addresses.length - (lastWithTxs + 1)) + howManyAddresses < maxAddressesGap) {
              callback(true);
              this.continueNewAddress();
            } else {
              const confirmationData: ConfirmationData = {
                text: 'wallet.add-many-confirmation',
                headerText: 'confirmation.header-text',
                confirmButtonText: 'confirmation.confirm-button',
                cancelButtonText: 'confirmation.cancel-button',
              };

              showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
                if (confirmationResult) {
                  callback(true);
                  this.continueNewAddress();
                } else {
                  callback(false);
                }
              });
            }
          }, () => callback(false, true));
        }
      };

      this.dialog.open(NumberOfAddressesComponent, config);
    } else {
      this.howManyAddresses = 1;
      this.continueNewAddress();
    }
  }

  toggleEmpty() {
    this.wallet.hideEmpty = !this.wallet.hideEmpty;
  }

  deleteWallet() {
    this.msgBarService.hide();

    const confirmationData: ConfirmationData = {
      text: this.translateService.instant('wallet.delete-confirmation', {name: this.wallet.label}),
      headerText: 'confirmation.header-text',
      checkboxText: 'wallet.delete-confirmation-check',
      confirmButtonText: 'confirmation.confirm-button',
      cancelButtonText: 'confirmation.cancel-button',
    };

    showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
      if (confirmationResult) {
        this.walletService.deleteHardwareWallet(this.wallet).subscribe(result => {
          if (result) {
            this.walletService.all().first().subscribe(wallets => {
              if (wallets.length === 0) {
                setTimeout(() => this.router.navigate(['/wizard']), 500);
              }
            });
          }
        });
      }
    });
  }

  toggleEncryption() {
    const config = new MatDialogConfig();
    config.data = {
      confirm: !this.wallet.encrypted,
      title: this.wallet.encrypted ? 'wallet.decrypt' : 'wallet.encrypt',
    };

    if (!this.wallet.encrypted) {
      config.data['description'] = 'wallet.new.encrypt-warning';
    } else {
      config.data['description'] = 'wallet.decrypt-warning';
      config.data['warning'] = true;
      config.data['wallet'] = this.wallet;
    }

    this.dialog.open(PasswordDialogComponent, config).componentInstance.passwordSubmit
      .subscribe(passwordDialog => {
        this.walletService.toggleEncryption(this.wallet, passwordDialog.password).subscribe(() => {
          passwordDialog.close();
          setTimeout(() => this.msgBarService.showDone('common.changes-made'));
        }, e => passwordDialog.error(e));
      });
  }

  confirmAddress(address, addressIndex, showCompleteConfirmation) {
    if (this.confirmingIndex !== null) {
      return;
    }

    this.confirmingIndex = addressIndex;
    this.msgBarService.hide();

    if (this.confirmSubscription) {
      this.confirmSubscription.unsubscribe();
    }

    this.confirmSubscription = this.hwWalletService.checkIfCorrectHwConnected(this.wallet.addresses[0].address).subscribe(response => {
      const data = new AddressConfirmationParams();
      data.address = address;
      data.addressIndex = addressIndex;
      data.showCompleteConfirmation = showCompleteConfirmation;

      const config = new MatDialogConfig();
      config.width = '566px';
      config.autoFocus = false;
      config.data = data;
      this.dialog.open(HwConfirmAddressDialogComponent, config);

      this.confirmingIndex = null;
    }, err => {
      this.msgBarService.showError(getHardwareWalletErrorMsg(this.translateService, err));
      this.confirmingIndex = null;
    });
  }

  copyAddress(event, address, duration = 500) {
    event.stopPropagation();

    if (address.copying) {
      return;
    }

    copyTextToClipboard(address.address);
    address.copying = true;

    setTimeout(function() {
      address.copying = false;
    }, duration);
  }

  showQrCode(event, address: string) {
    event.stopPropagation();

    const config: QrDialogConfig = { address };
    QrCodeComponent.openDialog(this.dialog, config);
  }

  private continueNewAddress() {
    this.creatingAddress = true;

    if (!this.wallet.isHardware && this.wallet.encrypted) {
      const config = new MatDialogConfig();
      config.data = {
        wallet: this.wallet,
      };

      const dialogRef = this.dialog.open(PasswordDialogComponent, config);
      dialogRef.afterClosed().subscribe(() => this.creatingAddress = false);
      dialogRef.componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.walletService.addAddress(this.wallet, this.howManyAddresses, passwordDialog.password)
            .subscribe(() => passwordDialog.close(), error => passwordDialog.error(error));
        });
    } else {

      let procedure: Observable<any>;

      if (this.wallet.isHardware ) {
        procedure = this.hwWalletService.checkIfCorrectHwConnected(this.wallet.addresses[0].address).flatMap(
          () => this.walletService.addAddress(this.wallet, this.howManyAddresses),
        );
      } else {
        procedure = this.walletService.addAddress(this.wallet, this.howManyAddresses);
      }

      procedure.subscribe(() => this.creatingAddress = false,
        err => {
          if (!this.wallet.isHardware ) {
            this.msgBarService.showError(err);
          } else {
            this.msgBarService.showError(getHardwareWalletErrorMsg(this.translateService, err));
          }
          this.creatingAddress = false;
        },
      );
    }
  }

  private continueEditWallet() {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.data = new ChangeNameData();
    (config.data as ChangeNameData).wallet = this.wallet;
    this.dialog.open(ChangeNameComponent, config);
  }
}
