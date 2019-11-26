import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService, TxData } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

@Component({
  selector: 'app-hw-confirm-tx-dialog',
  templateUrl: './hw-confirm-tx-dialog.component.html',
  styleUrls: ['./hw-confirm-tx-dialog.component.scss'],
})
export class HwConfirmTxDialogComponent extends HwDialogBaseComponent<HwConfirmTxDialogComponent> {

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: TxData[],
    public dialogRef: MatDialogRef<HwConfirmTxDialogComponent>,
    hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
  }
}
