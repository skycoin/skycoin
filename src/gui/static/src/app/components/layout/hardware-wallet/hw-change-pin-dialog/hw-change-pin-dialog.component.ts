import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { AppConfig } from '../../../../app.config';

enum States {
  Initial,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
  WrongPin,
  PinMismatch,
  DaemonError,
  Timeout,
}

@Component({
  selector: 'app-hw-change-pin-dialog',
  templateUrl: './hw-change-pin-dialog.component.html',
  styleUrls: ['./hw-change-pin-dialog.component.scss'],
})
export class HwChangePinDialogComponent extends HwDialogBaseComponent<HwChangePinDialogComponent> {

  changingExistingPin: boolean;
  currentState: States = States.Initial;
  states = States;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwChangePinDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);

    this.changingExistingPin = data.walletHasPin;

    this.operationSubscription = this.hwWalletService.getFeatures().flatMap(features => {
      if (!AppConfig.useHwWalletDaemon) {
        return this.hwWalletService.changePin(features.rawResponse.pinProtection);
      } else {
        return this.hwWalletService.changePin(features.rawResponse.pin_protection);
      }
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
