import { Component, Input, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { Observable, SubscriptionLike } from 'rxjs';
import { Router } from '@angular/router';
import { mergeMap, first } from 'rxjs/operators';

import { ChangeNameComponent, ChangeNameData } from '../change-name/change-name.component';
import { PasswordDialogComponent, PasswordDialogParams, PasswordSubmitEvent } from '../../../layout/password-dialog/password-dialog.component';
import { NumberOfAddressesComponent, NumberOfAddressesEventData } from '../number-of-addresses/number-of-addresses';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { copyTextToClipboard } from '../../../../utils/general-utils';
import { AppConfig } from '../../../../app.config';
import { HwConfirmAddressDialogComponent, AddressConfirmationParams } from '../../../layout/hardware-wallet/hw-confirm-address-dialog/hw-confirm-address-dialog.component';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { AddressOptionsComponent, AddressOptions } from './address-options/address-options.component';
import { ConfirmationParams, DefaultConfirmationButtons, ConfirmationComponent } from '../../../layout/confirmation/confirmation.component';
import { WalletsAndAddressesService } from '../../../../services/wallet-operations/wallets-and-addresses.service';
import { WalletWithBalance, AddressWithBalance } from '../../../../services/wallet-operations/wallet-objects';
import { SoftwareWalletService } from '../../../../services/wallet-operations/software-wallet.service';
import { HardwareWalletService } from '../../../../services/wallet-operations/hardware-wallet.service';
import { HistoryService } from '../../../../services/wallet-operations/history.service';
import { WalletsComponent } from '../wallets.component';

/**
 * Shows the option buttons and address list of a wallet on the wallet list.
 */
@Component({
  selector: 'app-wallet-detail',
  templateUrl: './wallet-detail.component.html',
  styleUrls: ['./wallet-detail.component.scss'],
})
export class WalletDetailComponent implements OnDestroy {
  @Input() wallet: WalletWithBalance;

  // Index of the address currently being confirmed. Used for showing the loading animation
  // on the UI.
  confirmingIndex = null;
  // If there is currently an operation with the addresses being done.
  workingWithAddresses = false;
  // If the preparations for renaming the wallet are being done.
  preparingToRename = false;
  // If all addresses without coins must be hidden on the address list.
  hideEmpty = false;
  // Allows to know which addresses are being copied, so the UI can show an indication.
  copying = new Map<string, boolean>();

  private renameSubscription: SubscriptionLike;
  private confirmSubscription: SubscriptionLike;
  private txHistorySubscription: SubscriptionLike;
  private scanSubscription: SubscriptionLike;
  private numberOfAddressesSubscription: SubscriptionLike;
  private addAddressesSubscription: SubscriptionLike;

  constructor(
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    private hwWalletService: HwWalletService,
    private translateService: TranslateService,
    private router: Router,
    private walletsAndAddressesService: WalletsAndAddressesService,
    private softwareWalletService: SoftwareWalletService,
    private hardwareWalletService: HardwareWalletService,
    private historyService: HistoryService,
  ) { }

  ngOnDestroy() {
    this.msgBarService.hide();
    if (this.renameSubscription) {
      this.renameSubscription.unsubscribe();
    }
    if (this.confirmSubscription) {
      this.confirmSubscription.unsubscribe();
    }
    if (this.txHistorySubscription) {
      this.txHistorySubscription.unsubscribe();
    }
    if (this.numberOfAddressesSubscription) {
      this.numberOfAddressesSubscription.unsubscribe();
    }
    if (this.scanSubscription) {
      this.scanSubscription.unsubscribe();
    }
    if (this.addAddressesSubscription) {
      this.addAddressesSubscription.unsubscribe();
    }
  }

  // Checks the wallet before opening the modal window for changing its label.
  renameWallet() {
    if (this.preparingToRename) {
      return;
    }

    if (WalletsComponent.busy) {
      this.msgBarService.showError('wallet.busy-error');

      return;
    }

    this.msgBarService.hide();

    if (this.wallet.isHardware) {
      this.preparingToRename = true;
      WalletsComponent.busy = true;

      // Check if the correct device is connected.
      this.renameSubscription = this.hwWalletService.checkIfCorrectHwConnected(this.wallet.addresses[0].address)
        // Check if the device still has the label this app knows.
        .pipe(mergeMap(() => this.hardwareWalletService.getFeaturesAndUpdateData(this.wallet)))
        .subscribe(
          response => {
            this.continueRenameWallet();
            this.preparingToRename = false;
            WalletsComponent.busy = false;

            // Inform if a different label was detected while checking the device.
            if (response.walletNameUpdated) {
              this.msgBarService.showWarning('hardware-wallet.general.name-updated');
            }
          },
          err => {
            this.msgBarService.showError(err);
            this.preparingToRename = false;
            WalletsComponent.busy = false;
          },
        );
    } else {
      // No checks needed for software wallets.
      this.continueRenameWallet();
    }
  }

  // Opens the modal window with options for making operations with the addresses.
  openAddressOptions() {
    if (this.workingWithAddresses) {
      return;
    }

    if (WalletsComponent.busy) {
      this.msgBarService.showError('wallet.busy-error');

      return;
    }

    AddressOptionsComponent.openDialog(this.dialog).afterClosed().subscribe(result => {
      if (result === AddressOptions.New) {
        this.newAddresses();
      } else if (result === AddressOptions.Scan) {
        this.scanAddresses();
      }
    });
  }

  // Adds addresses to the wallet. If the wallet is a software wallet, the user can select
  // how many addresses to add.
  private newAddresses() {
    if (this.workingWithAddresses) {
      return;
    }

    if (WalletsComponent.busy) {
      this.msgBarService.showError('wallet.busy-error');

      return;
    }

    // Don't allow more than the max number of addresses on a hw wallet.
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
      // Open the modal window for knowing how many addresses to add.
      const numberOfAddressesDialog = NumberOfAddressesComponent.openDialog(this.dialog);

      const maxAddressesGap = AppConfig.maxAddressesGap;
      // When the user requests the creation of the addresses, check if there will be a big
      // gap of unused addresses, before completing the operation.
      this.numberOfAddressesSubscription = numberOfAddressesDialog.componentInstance.createRequested.subscribe((eventData: NumberOfAddressesEventData) => {
        const howManyAddresses = eventData.howManyAddresses;
        const callback = eventData.callback;

        let lastWithBalance = 0;
        this.wallet.addresses.forEach((address, i) => {
          if (address.coins.isGreaterThan(0)) {
            lastWithBalance = i;
          }
        });

        // Try to use the current known balance to check if the new addresses will create
        // a gap of unused addresses bigger than the aceptable one. This is just a quick method
        // which is fast but could fail, as the code must detect a gap of unused addresses,
        // not one of addresses without balance.
        if ((this.wallet.addresses.length - (lastWithBalance + 1)) + howManyAddresses < maxAddressesGap) {
          callback(true);
          this.continueNewAddress(howManyAddresses);

          // If the previous check failed, use the real transaction history to be sure.
        } else {
          this.txHistorySubscription = this.historyService.getTransactionsHistory(this.wallet).pipe(first()).subscribe(transactions => {
            // Save which addresses have transaction history.
            const AddressesWithTxs = new Map<string, boolean>();
            transactions.forEach(transaction => {
              transaction.outputs.forEach(output => {
                if (!AddressesWithTxs.has(output.address)) {
                  AddressesWithTxs.set(output.address, true);
                }
              });
            });

            // Get the index of the last address with transaction history.
            let lastWithTxs = 0;
            this.wallet.addresses.forEach((address, i) => {
              if (AddressesWithTxs.has(address.address)) {
                lastWithTxs = i;
              }
            });

            if ((this.wallet.addresses.length - (lastWithTxs + 1)) + howManyAddresses < maxAddressesGap) {
              // Continue normally.
              callback(true);
              this.continueNewAddress(howManyAddresses);
            } else {
              // Tell the user that the gap could cause problems and ask for confirmation.
              const confirmationParams: ConfirmationParams = {
                text: 'wallet.add-many-confirmation',
                defaultButtons: DefaultConfirmationButtons.YesNo,
              };

              ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
                if (confirmationResult) {
                  callback(true);
                  this.continueNewAddress(howManyAddresses);
                } else {
                  callback(false);
                }
              });
            }
          }, () => callback(false, true));
        }
      });
    } else {
      // Hw wallets are limited to add one address at a time, for performance reasons.
      this.continueNewAddress(1);
    }
  }

  // Switches between showing and hiding the addresses without balance.
  toggleEmpty() {
    this.hideEmpty = !this.hideEmpty;
  }

  // Deletes a hw wallet.
  deleteHwWallet() {
    if (WalletsComponent.busy) {
      this.msgBarService.showError('wallet.busy-error');

      return;
    }

    this.msgBarService.hide();

    const confirmationParams: ConfirmationParams = {
      text: this.translateService.instant('wallet.delete-confirmation', {name: this.wallet.label}),
      checkboxText: 'wallet.delete-confirmation-check',
      defaultButtons: DefaultConfirmationButtons.YesNo,
    };

    // Ask for confirmation.
    ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
      if (confirmationResult) {
        this.walletsAndAddressesService.deleteHardwareWallet(this.wallet.id);

        // If there are no more wallets left, go to the wizard.
        this.walletsAndAddressesService.allWallets.pipe(first()).subscribe(wallets => {
          if (wallets.length === 0) {
            setTimeout(() => this.router.navigate(['/wizard']), 500);
          }
        });
      }
    });
  }

  // If the wallet is not encrypted, encrypts it. If the wallet is encrypted, removes
  // the encryption.
  toggleEncryption() {
    if (WalletsComponent.busy) {
      this.msgBarService.showError('wallet.busy-error');

      return;
    }

    const params: PasswordDialogParams = {
      confirm: !this.wallet.encrypted,
      title: this.wallet.encrypted ? 'wallet.decrypt-button' : 'wallet.encrypt-button',
      description: this.wallet.encrypted ? 'wallet.decrypt-warning' : 'wallet.new.encrypt-warning',
      warning: this.wallet.encrypted,
      wallet: this.wallet.encrypted ? this.wallet : null,
    };

    // Ask for the current password or the new one.
    PasswordDialogComponent.openDialog(this.dialog, params, false).componentInstance.passwordSubmit
      .subscribe(passwordDialog => {
        // Make the operation.
        this.softwareWalletService.toggleEncryption(this.wallet, passwordDialog.password).subscribe(() => {
          passwordDialog.close();
          setTimeout(() => this.msgBarService.showDone('common.changes-made'));
        }, e => passwordDialog.error(e));
      });
  }

  /**
   * Shows a modal window for the user to confirm if the address shown on the UI is equal to
   * the one stored on the hw wallet.
   * @param wallet Wallet with the address to be confirmed.
   * @param addressIndex Index of the address on the hw wallet.
   * @param showCompleteConfirmation Must be true if the address has not been confirmed yet, to
   * show a longer success message after the user confirms the address.
   */
  confirmAddress(wallet: WalletWithBalance, addressIndex: number, showCompleteConfirmation: boolean) {
    if (this.confirmingIndex !== null) {
      return;
    }

    if (WalletsComponent.busy) {
      this.msgBarService.showError('wallet.busy-error');

      return;
    }

    WalletsComponent.busy = true;
    this.confirmingIndex = addressIndex;
    this.msgBarService.hide();

    if (this.confirmSubscription) {
      this.confirmSubscription.unsubscribe();
    }

    // Check if the correct device is connected.
    this.confirmSubscription = this.hwWalletService.checkIfCorrectHwConnected(this.wallet.addresses[0].address).subscribe(() => {
      const data = new AddressConfirmationParams();
      data.wallet = wallet;
      data.addressIndex = addressIndex;
      data.showCompleteConfirmation = showCompleteConfirmation;

      HwConfirmAddressDialogComponent.openDialog(this.dialog, data);

      WalletsComponent.busy = false;
      this.confirmingIndex = null;
    }, err => {
      this.msgBarService.showError(err);
      WalletsComponent.busy = false;
      this.confirmingIndex = null;
    });
  }

  // Copies an address to the clipboard and sets it as being copied for the time set on
  // the "duration" param.
  copyAddress(event, address: AddressWithBalance, duration = 500) {
    event.stopPropagation();

    if (this.copying.has(address.address)) {
      return;
    }

    copyTextToClipboard(address.address);
    this.copying.set(address.address, true);

    setTimeout(() => {
      if (this.copying.has(address.address)) {
        this.copying.delete(address.address);
      }
    }, duration);
  }

  // Makes the preparations for asking the node to scan the addresses of the wallet again,
  // to add to it the addresses with transactions which have not been added to the addresses
  // list. Only for software wallets.
  private scanAddresses() {
    if (this.workingWithAddresses || this.wallet.isHardware) {
      return;
    }

    if (WalletsComponent.busy) {
      this.msgBarService.showError('wallet.busy-error');

      return;
    }

    this.workingWithAddresses = true;
    WalletsComponent.busy = true;

    // Ask for the password if the wallet is encrypted.
    if (this.wallet.encrypted) {
      const dialogRef = PasswordDialogComponent.openDialog(this.dialog, { wallet: this.wallet });

      dialogRef.afterClosed().subscribe(() => {
        this.workingWithAddresses = false;
        WalletsComponent.busy = false;
      });

      dialogRef.componentInstance.passwordSubmit.subscribe(passwordDialog => this.continueScanningAddresses(passwordDialog));
    } else {
      this.continueScanningAddresses();
    }
  }

  // Asks the node to scan the addresses of the wallet again.
  private continueScanningAddresses(passwordSubmitEvent?: PasswordSubmitEvent) {
    const password = passwordSubmitEvent ? passwordSubmitEvent.password : null;

    this.workingWithAddresses = true;
    WalletsComponent.busy = true;

    this.scanSubscription = this.walletsAndAddressesService.scanAddresses(this.wallet, password).subscribe(result => {
      if (passwordSubmitEvent) {
        passwordSubmitEvent.close();
      }

      setTimeout(() => {
        if (result) {
          this.msgBarService.showDone('wallet.scan-addresses.done-with-new-addresses');
        } else {
          this.msgBarService.showWarning('wallet.scan-addresses.done-without-new-addresses');
        }
      });

      this.workingWithAddresses = false;
      WalletsComponent.busy = false;
    }, err => {
      if (passwordSubmitEvent) {
        passwordSubmitEvent.error(err);
      } else {
        this.msgBarService.showError(err);
      }
      this.workingWithAddresses = false;
      WalletsComponent.busy = false;
    });
  }

  // Finish adding addresses to the wallet.
  private continueNewAddress(howManyAddresses: number) {
    this.workingWithAddresses = true;
    WalletsComponent.busy = true;

    if (!this.wallet.isHardware && this.wallet.encrypted) {
      // Ask for the password and continue.
      const dialogRef = PasswordDialogComponent.openDialog(this.dialog, { wallet: this.wallet });
      dialogRef.afterClosed().subscribe(() => {
        this.workingWithAddresses = false;
        WalletsComponent.busy = false;
      });
      dialogRef.componentInstance.passwordSubmit.subscribe(passwordDialog => {
        this.addAddressesSubscription = this.walletsAndAddressesService.addAddressesToWallet(this.wallet, howManyAddresses, passwordDialog.password).subscribe(() => {
          passwordDialog.close();
          setTimeout(() => this.msgBarService.showDone('common.changes-made'));
        }, error => passwordDialog.error(error));
      });
    } else {
      let procedure: Observable<any>;

      if (this.wallet.isHardware) {
        // Continue after checking the device.
        procedure = this.hwWalletService.checkIfCorrectHwConnected(this.wallet.addresses[0].address).pipe(mergeMap(
          () => this.walletsAndAddressesService.addAddressesToWallet(this.wallet, howManyAddresses),
        ));
      } else {
        procedure = this.walletsAndAddressesService.addAddressesToWallet(this.wallet, howManyAddresses);
      }

      this.addAddressesSubscription = procedure.subscribe(() => {
        this.workingWithAddresses = false;
        WalletsComponent.busy = false;
        this.msgBarService.showDone('common.changes-made');
      }, err => {
        this.msgBarService.showError(err);
        this.workingWithAddresses = false;
        WalletsComponent.busy = false;
      });
    }
  }

  // Opens the modal window for renaming the wallet.
  private continueRenameWallet() {
    const data = new ChangeNameData();
    data.wallet = this.wallet;
    ChangeNameComponent.openDialog(this.dialog, data, false);
  }
}
