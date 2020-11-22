import { Injectable } from '@angular/core';
import { MatDialog, MatDialogConfig } from '@angular/material';
import { HwPinDialogParams } from '../components/layout/hardware-wallet/hw-pin-dialog/hw-pin-dialog.component';
import { Observable } from 'rxjs/Observable';

export enum ChangePinStates {
  RequestingCurrentPin,
  RequestingNewPin,
  ConfirmingNewPin,
}

@Injectable()
export class HwWalletPinService {

  // Set on AppComponent to avoid a circular reference.
  private requestPinComponentInternal;
  set requestPinComponent(value) {
    this.requestPinComponentInternal = value;
  }

  // Values to be sent to HwPinDialogComponent
  changingPin: boolean;
  signingTx: boolean;
  changePinState: ChangePinStates;

  constructor(
    private dialog: MatDialog,
  ) {}

  resetValues() {
    this.changingPin = false;
    this.signingTx = false;
  }

  requestPin(): Observable<string> {
    return this.dialog.open(this.requestPinComponentInternal, <MatDialogConfig> {
      width: '350px',
      autoFocus: false,
      data : <HwPinDialogParams> {
        changingPin: this.changingPin,
        changePinState: this.changePinState,
        signingTx: this.signingTx,
      },
    }).afterClosed().map(pin => {
      if (this.changingPin) {
        if (this.changePinState === ChangePinStates.RequestingCurrentPin) {
          this.changePinState = ChangePinStates.RequestingNewPin;
        } else if (this.changePinState === ChangePinStates.RequestingNewPin) {
          this.changePinState = ChangePinStates.ConfirmingNewPin;
        }
      }

      return pin;
    });
  }
}
