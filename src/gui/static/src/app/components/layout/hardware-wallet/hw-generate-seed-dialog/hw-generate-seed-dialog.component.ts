import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

enum States {
  Initial,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
}

@Component({
  selector: 'app-hw-generate-seed-dialog',
  templateUrl: './hw-generate-seed-dialog.component.html',
  styleUrls: ['./hw-generate-seed-dialog.component.scss'],
})
export class HwGenerateSeedDialogComponent extends HwDialogBaseComponent<HwGenerateSeedDialogComponent> {

  currentState: States = States.Initial;
  states = States;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwGenerateSeedDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
    this.operationSubscription = this.hwWalletService.generateMnemonic().subscribe(
      () => {
        this.data.requestOptionsComponentRefresh();
        this.currentState = States.ReturnedSuccess;
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
