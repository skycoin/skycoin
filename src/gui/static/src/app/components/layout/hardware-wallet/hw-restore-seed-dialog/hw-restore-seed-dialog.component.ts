import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';

enum States {
  Initial,
  Processing,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
  WrongWord,
  WrongSeed,
  InvalidSeed,
}

@Component({
  selector: 'app-hw-restore-seed-dialog',
  templateUrl: './hw-restore-seed-dialog.component.html',
  styleUrls: ['./hw-restore-seed-dialog.component.scss'],
})
export class HwRestoreSeedDialogComponent extends HwDialogBaseComponent<HwRestoreSeedDialogComponent> {

  currentState: States = States.Initial;
  states = States;
  form: FormGroup;
  justCheckingSeed: boolean;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwRestoreSeedDialogComponent>,
    private hwWalletService: HwWalletService,
    formBuilder: FormBuilder,
  ) {
    super(hwWalletService, dialogRef);

    this.form = formBuilder.group({
      words: [24, Validators.required],
    });

    this.justCheckingSeed = !!this.data.wallet;
  }

  startOperation() {
    this.currentState = States.Processing;

    this.operationSubscription = this.hwWalletService.recoverMnemonic(this.form.controls['words'].value, this.justCheckingSeed).subscribe(
      () => {
        if (!this.justCheckingSeed) {
          this.data.requestOptionsComponentRefresh();
        }
        this.currentState = States.ReturnedSuccess;
      },
      err => {
        if (err.result && err.result === OperationResults.FailedOrRefused) {
          this.currentState = States.ReturnedRefused;
        } else if (err.result && err.result === OperationResults.WrongWord) {
          this.currentState = States.WrongWord;
        } else if (err.result && err.result === OperationResults.InvalidSeed) {
          this.currentState = States.InvalidSeed;
        } else if (err.result && err.result === OperationResults.WrongSeed) {
          this.currentState = States.WrongSeed;
        } else {
          this.currentState = States.Failed;
        }
      },
    );
  }
}
