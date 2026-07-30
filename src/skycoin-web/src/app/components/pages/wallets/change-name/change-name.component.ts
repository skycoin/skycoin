import { Component, Inject, OnInit, ViewChild, OnDestroy, ChangeDetectionStrategy } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

import { Wallet } from '../../../../app.datatypes';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MsgBarService } from '../../../../services/msg-bar.service';

@Component({
    selector: 'app-change-name',
    templateUrl: './change-name.component.html',
    styleUrls: ['./change-name.component.scss'],
    changeDetection: ChangeDetectionStrategy.Eager,
    standalone: false
})
export class ChangeNameComponent implements OnInit, OnDestroy {
  @ViewChild('button') button: ButtonComponent;
  form: UntypedFormGroup;

  constructor(
    @Inject(MAT_DIALOG_DATA) private data: Wallet,
    public dialogRef: MatDialogRef<ChangeNameComponent>,
    private formBuilder: UntypedFormBuilder,
    private walletService: WalletService,
    private msgBarService: MsgBarService,
  ) {}

  ngOnInit() {
    this.initForm();
  }

  ngOnDestroy() {
    this.msgBarService.hide();
  }

  closePopup() {
    this.dialogRef.close();
  }

  rename() {
    if (!this.form.valid || this.button.isLoading()) {
      return;
    }

    this.button.setLoading();
    this.data.label = this.form.value.label;
    this.walletService.saveWallets();
    this.dialogRef.close();

    setTimeout(() => this.msgBarService.showDone('common.changes-made'));
  }

  private initForm() {
    this.form = this.formBuilder.group({
      label: [this.data.label, Validators.required],
    });
  }
}
