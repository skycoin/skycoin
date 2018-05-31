import { Component, OnInit, Inject, ViewChild } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Wallet } from '../../../../app.datatypes';
import { ButtonComponent } from '../../../layout/button/button.component';

@Component({
  selector: 'app-change-name',
  templateUrl: './change-name.component.html',
  styleUrls: ['./change-name.component.css'],
})
export class ChangeNameComponent implements OnInit {
  @ViewChild('button') button: ButtonComponent;
  form: FormGroup;

  constructor(
    public dialogRef: MatDialogRef<ChangeNameComponent>,
    @Inject(MAT_DIALOG_DATA) private data: Wallet,
    private formBuilder: FormBuilder,
    private walletService: WalletService,
  ) {}

  ngOnInit() {
    this.form = this.formBuilder.group({
      label: [this.data.label, Validators.required],
    });
  }

  closePopup() {
    this.dialogRef.close();
  }

  rename() {
    if (!this.form.valid || this.button.isLoading()) {
      return;
    }

    this.button.setLoading();

    this.walletService.renameWallet(this.data, this.form.value.label)
      .subscribe(() => this.dialogRef.close(this.form.value.label));
  }
}
