import { Component, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { MessageIcons } from '../hw-message/hw-message.component';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';

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
  templateUrl: './hw-change-pin-dialog.component.html',
  styleUrls: ['./hw-change-pin-dialog.component.scss'],
})
export class HwChangePinDialogComponent implements OnDestroy {

  changingExistingPin: boolean;
  currentState: States = States.Initial;
  states = States;
  msgIcons = MessageIcons;

  private operationSubscription: ISubscription;
  private hwConnectionSubscription: ISubscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwChangePinDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    this.changingExistingPin = data.walletHasPin;

    this.operationSubscription = this.hwWalletService.getFeatures().flatMap(features => {
      return this.hwWalletService.changePin(features.rawResponse.pinProtection);
    }).subscribe(
      () => {
        this.currentState = States.ReturnedSuccess;
        this.data.requestOptionsComponentRefresh(null, true);
      },
      err => {
        if (err.result && err.result === OperationResults.PinMismatch) {
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
