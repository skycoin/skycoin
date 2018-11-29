import { Component, OnDestroy, ViewChild, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { MessageIcons } from '../hw-message/hw-message.component';

enum States {
  Initial,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
  WrongPin,
  PinMismatch,
}

@Component({
  selector: 'app-hw-change-pin-dialog',
  templateUrl: './hw-change-pin-dialog.html',
  styleUrls: ['./hw-change-pin-dialog.scss'],
})
export class HwChangePinDialogComponent implements OnDestroy {

  currentState: States = States.Initial;
  states = States;
  msgIcons = MessageIcons;

  private operationSubscription: ISubscription;
  private hwConnectionSubscription: ISubscription;

  constructor(
    public dialogRef: MatDialogRef<HwChangePinDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    this.operationSubscription = this.hwWalletService.changePin().subscribe(
      () => {
        this.currentState = States.ReturnedSuccess;
      },
      err => {
        if (err.rawResponse && typeof err.rawResponse === 'string' && (err.rawResponse as string).includes('PIN mismatch')) {
          this.currentState = States.PinMismatch;
        } else if (err.result && err.result === OperationResults.FailedOrRefused) {
          this.currentState = States.ReturnedRefused;
        } else if (err.result && err.result === OperationResults.WrongPin) {
          this.currentState = States.WrongPin;
        } else {
          this.currentState = States.Failed;
        }
      },
    );

    this.hwConnectionSubscription = this.hwWalletService.walletConnectedAsyncEvent.subscribe(connected => {
      if (!connected) {
        this.dialogRef.close();
      }
    });
  }

  ngOnDestroy() {
    this.operationSubscription.unsubscribe();
    this.hwConnectionSubscription.unsubscribe();
  }

  closeModal() {
    this.dialogRef.close();
  }
}
