import { mergeMap } from 'rxjs/operators';
import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

/**
 * Allows to add a PIN code to the device or change the one it has. This modal window was
 * created for being oppenend by the hw wallet options modal window.
 */
@Component({
  selector: 'app-hw-change-pin-dialog',
  templateUrl: './hw-change-pin-dialog.component.html',
  styleUrls: ['./hw-change-pin-dialog.component.scss'],
})
export class HwChangePinDialogComponent extends HwDialogBaseComponent<HwChangePinDialogComponent> {
  // If true, the device already has a PIN code and the operation is for changing it. If
  // false, the device does not have a PIN code and the operation is for creating one.
  changingExistingPin: boolean;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwChangePinDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);

    this.changingExistingPin = data.walletHasPin;

    this.operationSubscription = this.hwWalletService.getFeatures().pipe(mergeMap(features => {
      // Update the title, just in case, as it should not be needed.
      this.changingExistingPin = features.rawResponse.pin_protection;

      return this.hwWalletService.changePin(features.rawResponse.pin_protection);
    })).subscribe(
      () => {
        this.showResult({
          text: 'hardware-wallet.general.completed',
          icon: this.msgIcons.Success,
        });
        // Request the hw wallet options modal window to refresh the security warnings.
        this.data.requestOptionsComponentRefresh(null, true);
      },
      err => this.processHwOperationError(err),
    );
  }
}
