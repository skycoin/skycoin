import { Component, Inject, ChangeDetectionStrategy } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogConfig, MatDialog } from '@angular/material/dialog';

import { HwWalletTxRecipientData } from '../../../../services/hw-wallet.service';
import { MessageIcons } from '../hw-message/hw-message.component';

@Component({
    selector: 'app-hw-confirm-tx-dialog',
    templateUrl: './hw-confirm-tx-dialog.component.html',
    styleUrls: ['./hw-confirm-tx-dialog.component.scss'],
    changeDetection: ChangeDetectionStrategy.Eager,
    standalone: false
})
export class HwConfirmTxDialogComponent {
  msgIcons = MessageIcons;

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
  ) { }
}
