import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

/**
 * Allows to remove the PIN code protection from a device. This modal window was created for
 * being oppenend by the hw wallet options modal window.
 */
@Component({
    selector: 'app-hw-remove-pin-dialog',
    templateUrl: './hw-remove-pin-dialog.component.html',
    styleUrls: ['./hw-remove-pin-dialog.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class HwRemovePinDialogComponent extends HwDialogBaseComponent<HwRemovePinDialogComponent> {
  // If the user has confirmed the operation with the checkbox.
  confirmed = false;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwRemovePinDialogComponent>,
    private hwWalletService: HwWalletService,
    changeDetectorRef: ChangeDetectorRef,
  ) {
    super(hwWalletService, dialogRef, changeDetectorRef);
  }

  setConfirmed(event: any) {
    this.confirmed = event.checked;
  }

  // Starts the operation.
  requestRemoval() {
    this.currentState = this.states.Processing;

    this.operationSubscription = this.hwWalletService.removePin().subscribe(
      () => {
        this.showResult({
          text: 'hardware-wallet.general.completed',
          icon: this.msgIcons.Success,
        });
        // Request the hw wallet options modal window to refresh the security warnings.
        this.data.requestOptionsComponentRefresh(undefined, true);
        this.changeDetectorRef.markForCheck();
      },
      err => {
        this.processHwOperationError(err);
        this.changeDetectorRef.markForCheck();
      },
    );
  }
}
