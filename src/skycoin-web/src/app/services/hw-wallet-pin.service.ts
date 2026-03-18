import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Observable } from 'rxjs';
import { map } from 'rxjs';

import { HwPinDialogParams } from '../components/layout/hardware-wallet/hw-pin-dialog/hw-pin-dialog.component';

export enum ChangePinStates {
  RequestingCurrentPin = 'RequestingCurrentPin',
  RequestingNewPin = 'RequestingNewPin',
  ConfirmingNewPin = 'ConfirmingNewPin',
}

@Injectable()
export class HwWalletPinService {

  private requestPinComponentInternal;
  set requestPinComponent(value) {
    this.requestPinComponentInternal = value;
  }

  changingPin: boolean;
  signingTx: boolean;
  changePinState: ChangePinStates;

  constructor(
    private dialog: MatDialog,
  ) {}

  requestPin(): Observable<string> {
    return this.requestPinComponentInternal.openDialog(this.dialog, <HwPinDialogParams> {
      changingPin: this.changingPin,
      changePinState: this.changePinState,
      signingTx: this.signingTx,
    }).afterClosed().pipe(map(pin => {
      if (this.changingPin) {
        if (this.changePinState === ChangePinStates.RequestingCurrentPin) {
          this.changePinState = ChangePinStates.RequestingNewPin;
        } else if (this.changePinState === ChangePinStates.RequestingNewPin) {
          this.changePinState = ChangePinStates.ConfirmingNewPin;
        }
      }

      return pin;
    }));
  }
}
