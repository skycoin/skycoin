import { Component, OnInit, Inject } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MD_DIALOG_DATA, MdDialogRef } from '@angular/material';
import { Wallet } from '../../../../app.datatypes';

@Component({
  selector: 'app-change-name',
  templateUrl: './change-name.component.html',
  styleUrls: ['./change-name.component.css']
})
export class ChangeNameComponent implements OnInit {
  form: FormGroup;

  constructor(
    @Inject(MD_DIALOG_DATA) private data: Wallet,
    public dialogRef: MdDialogRef<ChangeNameComponent>,
    private formBuilder: FormBuilder,
    private walletService: WalletService,
  ) {}

  ngOnInit() {
    this.initForm();
  }

  closePopup() {
    this.dialogRef.close();
  }

  rename() {
    this.walletService.renameWallet(this.data, this.form.value.label)
      .subscribe(() => this.dialogRef.close(this.form.value.label));
  }

  private initForm() {
    this.form = this.formBuilder.group({
      label: [this.data.label, Validators.required],
    });
  }
}
