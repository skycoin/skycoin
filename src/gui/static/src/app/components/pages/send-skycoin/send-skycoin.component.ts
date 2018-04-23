import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatSnackBar, MatSnackBarConfig } from '@angular/material/snack-bar';
import 'rxjs/add/operator/delay';
import 'rxjs/add/operator/filter';
import { ButtonComponent } from '../../layout/button/button.component';
import { PasswordDialogComponent } from '../../layout/password-dialog/password-dialog.component';
import { MatDialog } from '@angular/material';
import { parseResponseMessage } from '../../../utils/index';

@Component({
  selector: 'app-send-skycoin',
  templateUrl: './send-skycoin.component.html',
  styleUrls: ['./send-skycoin.component.scss']
})
export class SendSkycoinComponent implements OnInit, OnDestroy {
  @ViewChild('button') button: ButtonComponent;

  form: FormGroup;
  transactions = [];

  constructor(
    public formBuilder: FormBuilder,
    public walletService: WalletService,
    private snackbar: MatSnackBar,
    private dialog: MatDialog,
  ) {}

  ngOnInit() {
    this.initForm();
  }

  ngOnDestroy() {
    this.snackbar.dismiss();
  }

  send() {
    this.button.resetState();
    this.snackbar.dismiss();

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
    this.button.setLoading();

    this.walletService.sendSkycoin(
      this.form.value.wallet,
      this.form.value.address,
      Math.round(parseFloat(this.form.value.amount) * 1000000),
      passwordDialog ? passwordDialog.password : null
    )
      .delay(1000)
      .subscribe(
        () => {
          this.resetForm();
          this.button.setSuccess();
        },
        error => {
          const errorMessage = parseResponseMessage(error['_body']);
          const config = new MatSnackBarConfig();
          config.duration = 300000;
          this.snackbar.open(errorMessage, null, config);
          this.button.setError(errorMessage);
        }
      );

    if (passwordDialog) {
      passwordDialog.close();
    }
  }

  private initForm() {
    this.form = this.formBuilder.group({
      wallet: ['', Validators.required],
      address: ['', Validators.required],
      amount: ['', [Validators.required, Validators.min(0), Validators.max(0)]],
      notes: [''],
    });
    this.form.get('wallet').valueChanges.subscribe(value => {
      console.log(value);
      const balance = value && value.coins ? value.coins : 0;
      this.form.get('amount').setValidators([
        Validators.required,
        Validators.min(0),
        Validators.max(balance),
      ]);
      this.form.get('amount').updateValueAndValidity();
    });
  }

  private resetForm() {
    this.form.get('wallet').reset(undefined);
    this.form.get('address').reset(undefined);
    this.form.get('amount').reset(undefined);
  }
}
