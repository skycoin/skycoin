import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { UntypedFormGroup, UntypedFormBuilder, Validators } from '@angular/forms';

import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

/**
 * Allows to load a seed on a seedless device or to let the user to enter a seed to check
 * if it is equal to the one on the device (if a wallet object is included in the params). This
 * modal window was created for being oppenend by the hw wallet options modal window.
 */
@Component({
  selector: 'app-hw-restore-seed-dialog',
  templateUrl: './hw-restore-seed-dialog.component.html',
  styleUrls: ['./hw-restore-seed-dialog.component.scss'],
})
export class HwRestoreSeedDialogComponent extends HwDialogBaseComponent<HwRestoreSeedDialogComponent> {
  form: UntypedFormGroup;
  // If true, the seed entered by the user will not be loaded on the device, the operation will
  // just compare it to the one on the device.
  justCheckingSeed: boolean;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwRestoreSeedDialogComponent>,
    private hwWalletService: HwWalletService,
    formBuilder: UntypedFormBuilder,
  ) {
    super(hwWalletService, dialogRef);

    this.form = formBuilder.group({
      words: [24, Validators.required],
    });

    this.justCheckingSeed = !!this.data.wallet;
  }

  startOperation() {
    this.currentState = this.states.Processing;

    this.operationSubscription = this.hwWalletService.recoverMnemonic(this.form.controls['words'].value, this.justCheckingSeed).subscribe(
      () => {
        if (!this.justCheckingSeed) {
          // Request the data and state of the hw wallet options modal window to be refreshed.
          this.data.requestOptionsComponentRefresh();
          this.closeModal();
        } else {
          this.showResult({
            text: 'hardware-wallet.restore-seed.correct-seed',
            icon: this.msgIcons.Success,
          });
        }
      },
      err => this.processHwOperationError(err),
    );
  }
}
