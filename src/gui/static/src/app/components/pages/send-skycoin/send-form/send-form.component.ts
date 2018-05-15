import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import 'rxjs/add/operator/delay';
import 'rxjs/add/operator/filter';
import { ButtonComponent } from '../../../layout/button/button.component';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { MatDialog, MatSnackBar, MatSnackBarConfig } from '@angular/material';
import { parseResponseMessage } from '../../../../utils/errors';

@Component({
  selector: 'app-send-form',
  templateUrl: './send-form.component.html',
  styleUrls: ['./send-form.component.scss'],
})
export class SendFormComponent implements OnInit {
  @ViewChild('button') button: ButtonComponent;
  @Input() formData: any;
  @Output() onFormSubmitted = new EventEmitter<any>();

  form: FormGroup;
  transactions = [];

  constructor(
    public formBuilder: FormBuilder,
    public walletService: WalletService,
    private dialog: MatDialog,
    private snackbar: MatSnackBar,
  ) {}

  ngOnInit() {
    this.initForm();
  }

  send() {
    if (!this.form.valid || this.button.isLoading()) {
      return;
    }

    this.snackbar.dismiss();
    this.button.resetState();

    if (this.form.value.wallet.encrypted) {
      this.dialog.open(PasswordDialogComponent).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this._send(passwordDialog);
        });
    } else {
      this._send();
    }
  }

  private _send(passwordDialog?: any) {
    if (passwordDialog) {
      passwordDialog.close();
    }

    this.button.setLoading();

    this.walletService.createTransaction(
      this.form.value.wallet,
      this.form.value.address,
      this.form.value.amount,
      passwordDialog ? passwordDialog.password : null,
    )
      .subscribe(transaction => {
        this.onFormSubmitted.emit({
          wallet: this.form.value.wallet,
          address: this.form.value.address,
          amount: this.form.value.amount,
          notes: this.form.value.notes,
          transaction,
        });
      }, error => {
        const errorMessage = parseResponseMessage(error['_body']);
        const config = new MatSnackBarConfig();
        config.duration = 300000;
        this.snackbar.open(errorMessage, null, config);
        this.button.setError(errorMessage);
      });
  }

  private initForm() {
    this.form = this.formBuilder.group({
      wallet: ['', Validators.required],
      address: ['', Validators.required],
      amount: ['', Validators.required],
      notes: [''],
    });

    this.form.get('wallet').valueChanges.subscribe(value => {
      const balance = value && value.coins ? value.coins : 0;

      this.form.get('amount').setValidators([
        Validators.required,
        Validators.max(balance),
        this.validateAmount,
      ]);

      this.form.get('amount').updateValueAndValidity();
    });

    if (this.formData) {
      Object.keys(this.form.controls).forEach(control => {
        this.form.get(control).setValue(this.formData[control]);
      });
    }
  }

  private validateAmount(amountControl: FormControl) {
    if (isNaN(amountControl.value)) {
      return { Invalid: true };
    }

    if (parseFloat(amountControl.value) <= 0) {
      return { Invalid: true };
    }

    const parts = amountControl.value.split('.');

    if (parts.length === 2 && parts[1].length > 6) {
      return { Invalid: true };
    }

    return null;
  }
}
