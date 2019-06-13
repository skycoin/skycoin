import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

enum States {
  Initial,
  Processing,
  ReturnedSuccess,
  ReturnedRefused,
  WrongPin,
  Failed,
  DaemonError,
  Timeout,
}

@Component({
  selector: 'app-hw-remove-pin-dialog',
  templateUrl: './hw-remove-pin-dialog.component.html',
  styleUrls: ['./hw-remove-pin-dialog.component.scss'],
})
export class HwRemovePinDialogComponent extends HwDialogBaseComponent<HwRemovePinDialogComponent> {

  currentState: States = States.Initial;
  states = States;
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
    this.currentState = States.Processing;

    this.operationSubscription = this.hwWalletService.removePin().subscribe(
      () => {
        this.data.requestOptionsComponentRefresh(null, true);
        this.currentState = States.ReturnedSuccess;
      },
      err => {
        if (err.result && err.result === OperationResults.FailedOrRefused) {
          this.currentState = States.ReturnedRefused;
        } else if (err.result && err.result === OperationResults.WrongPin) {
          this.currentState = States.WrongPin;
        } else if (err.result && err.result === OperationResults.DaemonError) {
          this.currentState = States.DaemonError;
        } else if (err.result && err.result === OperationResults.Timeout) {
          this.currentState = States.Timeout;
        } else if (err.result && err.result === OperationResults.Disconnected) {
          this.closeModal();
        } else {
          this.currentState = States.Failed;
        }
      },
    );
  }
}
