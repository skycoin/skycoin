import { Component, Inject } from '@angular/core';
import { MatLegacyDialogRef as MatDialogRef, MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA, MatLegacyDialogConfig as MatDialogConfig, MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';

import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { HardwareWalletService } from '../../../../services/wallet-operations/hardware-wallet.service';
import { WalletBase } from '../../../../services/wallet-operations/wallet-objects';

/**
 * Settings for HwConfirmAddressDialogComponent.
 */
export class AddressConfirmationParams {
  /**
   * Wallet which contains the address to sonfirm.
   */
  wallet: WalletBase;
  /**
   * Index of the address inside the wallet.
   */
  addressIndex: number;
  /**
   * If true, the UI will show the complete confirmation text after the user confirms the address.
   * If false, a short text will be shown. The complete text should be displayed only the first
   * time the user confirms the address.
   */
  showCompleteConfirmation: boolean;
}

/**
 * Allows the user to confirm the desktop and hw wallets show the same address.
 */
@Component({
  selector: 'app-hw-confirm-address-dialog',
  templateUrl: './hw-confirm-address-dialog.component.html',
  styleUrls: ['./hw-confirm-address-dialog.component.scss'],
})
export class HwConfirmAddressDialogComponent extends HwDialogBaseComponent<HwConfirmAddressDialogComponent> {
  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog, confirmationParams: AddressConfirmationParams): MatDialogRef<HwConfirmAddressDialogComponent, any> {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.autoFocus = false;
    config.data = confirmationParams;

    return dialog.open(HwConfirmAddressDialogComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: AddressConfirmationParams,
    public dialogRef: MatDialogRef<HwConfirmAddressDialogComponent>,
    private hardwareWalletService: HardwareWalletService,
    hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);

    // Ask for confirmation and update the address and the wallet if the user confirms.
    this.operationSubscription = this.hardwareWalletService.confirmAddress(data.wallet, data.addressIndex).subscribe(
      () => {
        this.showResult({
          text: data.showCompleteConfirmation ? 'hardware-wallet.confirm-address.confirmation' : 'hardware-wallet.confirm-address.short-confirmation',
          icon: this.msgIcons.Success,
        });
      },
      err => this.processHwOperationError(err),
    );
  }
}
