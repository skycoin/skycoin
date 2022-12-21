import { Injectable } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { HwPinDialogParams } from '../components/layout/hardware-wallet/hw-pin-dialog/hw-pin-dialog.component';

/**
 * Diferent modes in which the modal window used for requesting the hw wallet PIN can be while
 * setting or changing the PIN.
 */
export enum ChangePinStates {
  RequestingCurrentPin = 'RequestingCurrentPin',
  RequestingNewPin = 'RequestingNewPin',
  ConfirmingNewPin = 'ConfirmingNewPin',
}

/**
 * Allows to easily show the modal window used for requesting the hw wallet PIN.
 */
@Injectable()
export class HwWalletPinService {

  // Set on AppComponent to avoid a circular reference.
  private requestPinComponentInternal;
  /**
   * Sets the class of the modal window used for entering the hw wallet PIN.
   */
  set requestPinComponent(value) {
    this.requestPinComponentInternal = value;
  }

  // Values to be sent to HwPinDialogComponent to configure the modal window the next
  // time it is openned. The values are public to make it posible to change them from
  // different parts of the code, as there are multiple parts which need to configure
  // the modal window but it is automatically opened by the daemon service.
  /**
   * If the modal window will be openned the next time for setting or changing the PIN
   * (true) or just for checking the current PIN (false).
   */
  changingPin: boolean;
  /**
   * If the modal window will be openned the next time for signing a transaction or not.
   */
  signingTx: boolean;
  /**
   * State in which the modal window will be shown the next time if changingPin is true.
   */
  changePinState: ChangePinStates;

  constructor(
    private dialog: MatDialog,
  ) {}

  /**
   * Opens the modal window for the user to enter the PIN.
   * @returns The PIN entered by the user, or null if the user cancelled the operation.
   */
  requestPin(): Observable<string> {
    return this.requestPinComponentInternal.openDialog(this.dialog, {
      changingPin: this.changingPin,
      changePinState: this.changePinState,
      signingTx: this.signingTx,
    } as HwPinDialogParams).afterClosed().pipe(map(pin => {
      if (this.changingPin) {
        // If setting or changing the PIN, automatically change the state to the one corresponding
        // to the next step.
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
