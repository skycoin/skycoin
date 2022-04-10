import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';

import { AppConfig } from '../../../../../app.config';
import { DecodedTransaction } from '../../../../../services/wallet-operations/transaction-objects';

/**
 * Allows to see the preview of a raw transaction. If the user confirms the operation,
 * the modal window is closed and "true" is returned in the "afterClosed" event.
 */
@Component({
  selector: 'app-offline-tx-preview',
  templateUrl: 'offline-tx-preview.component.html',
  styleUrls: ['offline-tx-preview.component.scss'],
})
export class OfflineTxPreviewComponent {
  tx: DecodedTransaction;

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog, data: DecodedTransaction): MatDialogRef<OfflineTxPreviewComponent, any> {
    const config = new MatDialogConfig();
    config.data = data;
    config.autoFocus = true;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(OfflineTxPreviewComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) data: DecodedTransaction,
    public dialogRef: MatDialogRef<OfflineTxPreviewComponent>,
  ) {
    this.tx = data;
  }

  cancelPressed() {
    this.dialogRef.close();
  }

  okPressed() {
    this.dialogRef.close(true);
  }
}
