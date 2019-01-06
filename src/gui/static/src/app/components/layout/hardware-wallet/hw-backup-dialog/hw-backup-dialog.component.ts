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
  Failed,
}

@Component({
  selector: 'app-hw-backup-dialog',
  templateUrl: './hw-backup-dialog.component.html',
  styleUrls: ['./hw-backup-dialog.component.scss'],
})
export class HwBackupDialogComponent extends HwDialogBaseComponent<HwBackupDialogComponent> {

  currentState: States = States.Initial;
  states = States;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwBackupDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
  }

  requestBackup() {
    this.currentState = States.Processing;

    this.operationSubscription = this.hwWalletService.backup().subscribe(
      () => {
        this.currentState = States.ReturnedSuccess;
        this.data.requestOptionsComponentRefresh(null, true);
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
