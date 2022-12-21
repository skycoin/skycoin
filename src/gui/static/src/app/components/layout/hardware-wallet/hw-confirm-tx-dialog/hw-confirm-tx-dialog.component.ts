import { Component, Inject } from '@angular/core';
import { MatLegacyDialogRef as MatDialogRef, MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA, MatLegacyDialogConfig as MatDialogConfig, MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';

import { HwWalletService, HwWalletTxRecipientData } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

/**
 * Allow the user to confirm a transaction before sending it.
 */
@Component({
  selector: 'app-hw-confirm-tx-dialog',
  templateUrl: './hw-confirm-tx-dialog.component.html',
  styleUrls: ['./hw-confirm-tx-dialog.component.scss'],
})
export class HwConfirmTxDialogComponent extends HwDialogBaseComponent<HwConfirmTxDialogComponent> {
  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog, recipientData: HwWalletTxRecipientData[]): MatDialogRef<HwConfirmTxDialogComponent, any> {
    const config = new MatDialogConfig();
    config.data = recipientData;
    config.autoFocus = false;
    config.width = '600px';

    return dialog.open(HwConfirmTxDialogComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: HwWalletTxRecipientData[],
    public dialogRef: MatDialogRef<HwConfirmTxDialogComponent>,
    hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
  }
}
