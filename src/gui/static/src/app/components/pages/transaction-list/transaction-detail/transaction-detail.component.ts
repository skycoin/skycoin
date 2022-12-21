import { Component, Inject } from '@angular/core';
import { MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA, MatLegacyDialogRef as MatDialogRef, MatLegacyDialog as MatDialog, MatLegacyDialogConfig as MatDialogConfig } from '@angular/material/legacy-dialog';

import { OldTransaction } from '../../../../services/wallet-operations/transaction-objects';

/**
 * Modal window for showing the details of a transaction from the transaction history.
 */
@Component({
  selector: 'app-transaction-detail',
  templateUrl: './transaction-detail.component.html',
  styleUrls: ['./transaction-detail.component.scss'],
})
export class TransactionDetailComponent {
  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog, transaction: OldTransaction): MatDialogRef<TransactionDetailComponent, any> {
    const config = new MatDialogConfig();
    config.data = transaction;
    config.autoFocus = false;
    config.width = '800px';

    return dialog.open(TransactionDetailComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) public transaction: OldTransaction,
    public dialogRef: MatDialogRef<TransactionDetailComponent>,
  ) {}

  closePopup() {
    this.dialogRef.close();
  }
}
