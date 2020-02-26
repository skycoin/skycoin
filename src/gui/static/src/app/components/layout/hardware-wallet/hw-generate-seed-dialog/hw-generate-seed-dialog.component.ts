import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';

import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

/**
 * Allows to make a seedless device create a new random seed and use it. This modal window was
 * created for being oppenend by the hw wallet options modal window.
 */
@Component({
  selector: 'app-hw-generate-seed-dialog',
  templateUrl: './hw-generate-seed-dialog.component.html',
  styleUrls: ['./hw-generate-seed-dialog.component.scss'],
})
export class HwGenerateSeedDialogComponent extends HwDialogBaseComponent<HwGenerateSeedDialogComponent> {
  form: FormGroup;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwGenerateSeedDialogComponent>,
    private hwWalletService: HwWalletService,
    formBuilder: FormBuilder,
  ) {
    super(hwWalletService, dialogRef);

    this.form = formBuilder.group({
      words: [24, Validators.required],
    });
  }

  startOperation() {
    this.currentState = this.states.Processing;

    this.operationSubscription = this.hwWalletService.generateMnemonic(this.form.controls['words'].value).subscribe(
      () => {
        // Request the data and state of the hw wallet options modal window to be refreshed.
        this.data.requestOptionsComponentRefresh();
        this.closeModal();
      },
      err => this.processHwOperationError(err),
    );
  }
}
