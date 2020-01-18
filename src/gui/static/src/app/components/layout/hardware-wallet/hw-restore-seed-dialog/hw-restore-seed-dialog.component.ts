import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';

@Component({
  selector: 'app-hw-restore-seed-dialog',
  templateUrl: './hw-restore-seed-dialog.component.html',
  styleUrls: ['./hw-restore-seed-dialog.component.scss'],
})
export class HwRestoreSeedDialogComponent extends HwDialogBaseComponent<HwRestoreSeedDialogComponent> {
  form: FormGroup;
  justCheckingSeed: boolean;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwRestoreSeedDialogComponent>,
    private hwWalletService: HwWalletService,
    formBuilder: FormBuilder,
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
          this.data.requestOptionsComponentRefresh();
          this.closeModal();
        } else {
          this.showResult({
            text: 'hardware-wallet.restore-seed.correct-seed',
            icon: this.msgIcons.Success,
          });
        }
      },
      err => this.processResult(err.result, 'hardware-wallet.errors.simple-error'),
    );
  }
}
