import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { WalletService } from '../../../../services/wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

enum States {
  Initial,
  Processing,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
}

@Component({
  selector: 'app-hw-wipe-dialog',
  templateUrl: './hw-wipe-dialog.component.html',
  styleUrls: ['./hw-wipe-dialog.component.scss'],
})
export class HwWipeDialogComponent extends HwDialogBaseComponent<HwWipeDialogComponent> {

  currentState: States = States.Initial;
  states = States;
  showDeleteFromList = true;
  deleteFromList = true;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwWipeDialogComponent>,
    private hwWalletService: HwWalletService,
    private walletService: WalletService,
  ) {
    super(hwWalletService, dialogRef);

    if (!data.wallet) {
      this.showDeleteFromList = false;
      this.deleteFromList = false;
    }
  }

  setDeleteFromList(event) {
    this.deleteFromList = event.checked;
  }

  requestWipe() {
    this.currentState = States.Processing;

    this.operationSubscription = this.hwWalletService.wipe().subscribe(
      () => {
        this.data.requestOptionsComponentRefresh();
        this.currentState = States.ReturnedSuccess;
        if (this.deleteFromList) {
          this.walletService.deleteHardwareWallet(this.data.wallet).subscribe();
        }
      },
      err => {
        if (err.result && err.result === OperationResults.FailedOrRefused) {
          this.currentState = States.ReturnedRefused;
        } else {
          this.currentState = States.Failed;
        }
      },
    );
  }
}
