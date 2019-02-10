import { Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

@Component({
  selector: 'app-hw-confirm-tx-dialog',
  templateUrl: './hw-confirm-tx-dialog.component.html',
  styleUrls: ['./hw-confirm-tx-dialog.component.scss'],
})
export class HwConfirmTxDialogComponent extends HwDialogBaseComponent<HwConfirmTxDialogComponent> {

  constructor(
    public dialogRef: MatDialogRef<HwConfirmTxDialogComponent>,
    hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
  }
}
