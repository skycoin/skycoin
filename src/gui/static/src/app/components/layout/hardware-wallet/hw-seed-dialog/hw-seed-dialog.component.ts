import { Component, OnDestroy, ViewChild, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { MessageIcons } from '../hw-message/hw-message.component';

enum States {
  Initial,
  Processing,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
}

@Component({
  selector: 'app-hw-seed-dialog',
  templateUrl: './hw-seed-dialog.component.html',
  styleUrls: ['./hw-seed-dialog.component.scss'],
})
export class HwSeedDialogComponent implements OnDestroy {
  form: FormGroup;
  currentState: States = States.Initial;
  states = States;
  msgIcons = MessageIcons;

  private operationSubscription: ISubscription;
  private hwConnectionSubscription: ISubscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) public notifyFinish: any,
    public dialogRef: MatDialogRef<HwSeedDialogComponent>,
    private hwWalletService: HwWalletService,
    private formBuilder: FormBuilder,
  ) {
    this.form = this.formBuilder.group({
      seed: ['', Validators.required],
    });

    this.hwConnectionSubscription = this.hwWalletService.walletConnectedAsyncEvent.subscribe(connected => {
      if (!connected) {
        this.closeModal();
      }
    });
  }

  closeModal() {
    this.dialogRef.close();
  }

  setSeed() {
    this.currentState = States.Processing;
    this.operationSubscription = this.hwWalletService.setMnemonic(this.form.value.seed).subscribe(
      () => {
        this.notifyFinish();
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

  ngOnDestroy() {
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }
    this.hwConnectionSubscription.unsubscribe();
  }
}
