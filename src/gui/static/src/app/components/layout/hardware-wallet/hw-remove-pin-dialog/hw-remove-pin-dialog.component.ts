import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

@Component({
  selector: 'app-hw-remove-pin-dialog',
  templateUrl: './hw-remove-pin-dialog.component.html',
  styleUrls: ['./hw-remove-pin-dialog.component.scss'],
})
export class HwRemovePinDialogComponent extends HwDialogBaseComponent<HwRemovePinDialogComponent> {
  confirmed = false;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwRemovePinDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
  }

  setConfirmed(event) {
    this.confirmed = event.checked;
  }

  requestRemoval() {
    this.currentState = this.states.Processing;

    this.operationSubscription = this.hwWalletService.removePin().subscribe(
      () => {
        this.showResult({
          text: 'hardware-wallet.general.completed',
          icon: this.msgIcons.Success,
        });
        this.data.requestOptionsComponentRefresh(null, true);
      },
      err => this.processResult(err.result),
    );
  }
}
