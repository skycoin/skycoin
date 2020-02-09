import { Component, Input, OnDestroy } from '@angular/core';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { ChangeNameComponent, ChangeNameData } from '../change-name/change-name.component';
import { PasswordDialogComponent, PasswordDialogParams } from '../../../layout/password-dialog/password-dialog.component';
import { NumberOfAddressesComponent } from '../number-of-addresses/number-of-addresses';
import { TranslateService } from '@ngx-translate/core';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { Observable, SubscriptionLike } from 'rxjs';
import { copyTextToClipboard } from '../../../../utils/general-utils';
import { AppConfig } from '../../../../app.config';
import { Router } from '@angular/router';
import { HwConfirmAddressDialogComponent, AddressConfirmationParams } from '../../../layout/hardware-wallet/hw-confirm-address-dialog/hw-confirm-address-dialog.component';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { ApiService } from '../../../../services/api.service';
import { mergeMap, first } from 'rxjs/operators';
import { AddressOptionsComponent, AddressOptions } from './address-options/address-options.component';
import { ConfirmationParams, DefaultConfirmationButtons, ConfirmationComponent } from '../../../layout/confirmation/confirmation.component';
import { WalletsAndAddressesService } from '../../../../services/wallet-operations/wallets-and-addresses.service';
import { WalletWithBalance } from '../../../../services/wallet-operations/wallet-objects';
import { SoftwareWalletService } from '../../../../services/wallet-operations/software-wallet.service';
import { HardwareWalletService } from '../../../../services/wallet-operations/hardware-wallet.service';
import { HistoryService } from '../../../../services/wallet-operations/history.service';

@Component({
  selector: 'app-wallet-detail',
  templateUrl: './wallet-detail.component.html',
  styleUrls: ['./wallet-detail.component.scss'],
})
export class WalletDetailComponent implements OnDestroy {
  @Input() wallet: WalletWithBalance;

  confirmingIndex = null;
  workingWithAddresses = false;
  preparingToEdit = false;
  hideEmpty = false;

  private howManyAddresses: number;
  private editSubscription: SubscriptionLike;
  private confirmSubscription: SubscriptionLike;
  private txHistorySubscription: SubscriptionLike;

  constructor(
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    private hwWalletService: HwWalletService,
    private translateService: TranslateService,
    private router: Router,
    private apiService: ApiService,
    private walletsAndAddressesService: WalletsAndAddressesService,
    private softwareWalletService: SoftwareWalletService,
    private hardwareWalletService: HardwareWalletService,
    private historyService: HistoryService,
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
        .pipe(mergeMap(() => this.hardwareWalletService.getFeaturesAndUpdateData(this.wallet)))
        .subscribe(
          response => {
            this.continueEditWallet();
            this.preparingToEdit = false;

            if (response.walletNameUpdated) {
              this.msgBarService.showWarning('hardware-wallet.general.name-updated');
            }
          },
          err => {
            this.msgBarService.showError(err);
            this.preparingToEdit = false;
          },
        );
    } else {
      this.continueEditWallet();
    }
  }

  openAddressOptions() {
    if (this.workingWithAddresses) {
      return;
    }

    AddressOptionsComponent.openDialog(this.dialog).afterClosed().subscribe(result => {
      if (result === AddressOptions.new) {
        this.newAddress();
      } else if (result === AddressOptions.scan) {
        this.scanAddresses();
      }
    });
  }

  newAddress() {
    if (this.workingWithAddresses) {
      return;
    }

    if (this.wallet.isHardware && this.wallet.addresses.length >= AppConfig.maxHardwareWalletAddresses) {
      const confirmationParams: ConfirmationParams = {
        text: 'wallet.max-hardware-wallets-error',
        headerText: 'common.error',
        defaultButtons: DefaultConfirmationButtons.Close,
      };
      ConfirmationComponent.openDialog(this.dialog, confirmationParams);

      return;
    }

    this.msgBarService.hide();

    if (!this.wallet.isHardware) {
      const maxAddressesGap = 20;

      const eventFunction = (howManyAddresses, callback) => {
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
          this.txHistorySubscription = this.historyService.getTransactionsHistory(this.wallet).pipe(first()).subscribe(transactions => {
            const AddressesWithTxs = new Map<string, boolean>();

            transactions.forEach(transaction => {
              transaction.outputs.forEach(output => {
                if (!AddressesWithTxs.has(output.address)) {
                  AddressesWithTxs.set(output.address, true);
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
              const confirmationParams: ConfirmationParams = {
                text: 'wallet.add-many-confirmation',
                defaultButtons: DefaultConfirmationButtons.YesNo,
              };

              ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
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

      NumberOfAddressesComponent.openDialog(this.dialog, eventFunction);
    } else {
      this.howManyAddresses = 1;
      this.continueNewAddress();
    }
  }

  toggleEmpty() {
    this.hideEmpty = !this.hideEmpty;
  }

  deleteWallet() {
    this.msgBarService.hide();

    const confirmationParams: ConfirmationParams = {
      text: this.translateService.instant('wallet.delete-confirmation', {name: this.wallet.label}),
      checkboxText: 'wallet.delete-confirmation-check',
      defaultButtons: DefaultConfirmationButtons.YesNo,
    };

    ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
      if (confirmationResult) {
        this.walletsAndAddressesService.deleteHardwareWallet(this.wallet.id);

        this.walletsAndAddressesService.allWallets.pipe(first()).subscribe(wallets => {
          if (wallets.length === 0) {
            setTimeout(() => this.router.navigate(['/wizard']), 500);
          }
        });
      }
    });
  }

  toggleEncryption() {
    const params: PasswordDialogParams = {
      confirm: !this.wallet.encrypted,
      title: this.wallet.encrypted ? 'wallet.decrypt-button' : 'wallet.encrypt-button',
      description: this.wallet.encrypted ? 'wallet.decrypt-warning' : 'wallet.new.encrypt-warning',
      warning: this.wallet.encrypted,
      wallet: this.wallet.encrypted ? this.wallet : null,
    };

    PasswordDialogComponent.openDialog(this.dialog, params, false).componentInstance.passwordSubmit
      .subscribe(passwordDialog => {
        this.softwareWalletService.toggleEncryption(this.wallet, passwordDialog.password).subscribe(() => {
          passwordDialog.close();
          setTimeout(() => this.msgBarService.showDone('common.changes-made'));
        }, e => passwordDialog.error(e));
      });
  }

  confirmAddress(wallet, addressIndex, showCompleteConfirmation) {
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
      data.wallet = wallet;
      data.addressIndex = addressIndex;
      data.showCompleteConfirmation = showCompleteConfirmation;

      const config = new MatDialogConfig();
      config.width = '566px';
      config.autoFocus = false;
      config.data = data;
      this.dialog.open(HwConfirmAddressDialogComponent, config);

      this.confirmingIndex = null;
    }, err => {
      this.msgBarService.showError(err);
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

  private scanAddresses() {
    if (this.workingWithAddresses) {
      return;
    }

    this.workingWithAddresses = true;

    if (!this.wallet.isHardware && this.wallet.encrypted) {
      const dialogRef = PasswordDialogComponent.openDialog(this.dialog, { wallet: this.wallet });
      dialogRef.afterClosed().subscribe(() => this.workingWithAddresses = false);
      dialogRef.componentInstance.passwordSubmit.subscribe(passwordDialog => {
        this.walletsAndAddressesService.scanAddresses(this.wallet, passwordDialog.password).subscribe(result => {
          passwordDialog.close();

          setTimeout(() => {
            if (result) {
              this.msgBarService.showDone('wallet.scan-addresses.done-with-new-addresses');
            } else {
              this.msgBarService.showWarning('wallet.scan-addresses.done-without-new-addresses');
            }
          });
        }, error => {
          passwordDialog.error(error);
        });
      });
    } else {
      this.walletsAndAddressesService.scanAddresses(this.wallet).subscribe(result => {
        if (result) {
          this.msgBarService.showDone('wallet.scan-addresses.done-with-new-addresses');
        } else {
          this.msgBarService.showWarning('wallet.scan-addresses.done-without-new-addresses');
        }
        this.workingWithAddresses = false;
      }, err => {
        this.msgBarService.showError(err);
        this.workingWithAddresses = false;
      });
    }
  }

  private continueNewAddress() {
    this.workingWithAddresses = true;

    if (!this.wallet.isHardware && this.wallet.encrypted) {
      const dialogRef = PasswordDialogComponent.openDialog(this.dialog, { wallet: this.wallet });
      dialogRef.afterClosed().subscribe(() => this.workingWithAddresses = false);
      dialogRef.componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.walletsAndAddressesService.addAddressesToWallet(this.wallet, this.howManyAddresses, passwordDialog.password)
            .subscribe(() => {
              passwordDialog.close();
              setTimeout(() => this.msgBarService.showDone('common.changes-made'));
            }, error => passwordDialog.error(error));
        });
    } else {

      let procedure: Observable<any>;

      if (this.wallet.isHardware ) {
        procedure = this.hwWalletService.checkIfCorrectHwConnected(this.wallet.addresses[0].address).pipe(mergeMap(
          () => this.walletsAndAddressesService.addAddressesToWallet(this.wallet, this.howManyAddresses),
        ));
      } else {
        procedure = this.walletsAndAddressesService.addAddressesToWallet(this.wallet, this.howManyAddresses);
      }

      procedure.subscribe(() => {
        this.workingWithAddresses = false;
        this.msgBarService.showDone('common.changes-made');
      }, err => {
        this.msgBarService.showError(err);
        this.workingWithAddresses = false;
      });
    }
  }

  private continueEditWallet() {
    const data = new ChangeNameData();
    data.wallet = this.wallet;
    ChangeNameComponent.openDialog(this.dialog, data, false);
  }
}
