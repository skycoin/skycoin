import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';

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
        this.showResult({
          text: 'hardware-wallet.general.completed',
          icon: this.msgIcons.Success,
        });
        this.data.requestOptionsComponentRefresh();
      },
      err => this.processResult(err.result),
    );
  }
}
