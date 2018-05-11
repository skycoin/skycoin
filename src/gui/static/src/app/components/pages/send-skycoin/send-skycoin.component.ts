import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
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
    if (!this.form.valid || this.button.isLoading()) {
      return;
    }

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
    if (passwordDialog) {
      passwordDialog.close();
    }

    this.button.setLoading();

    this.walletService.createTransaction(
      this.form.value.wallet,
      this.form.value.address,
      this.form.value.amount,
      passwordDialog ? passwordDialog.password : null
    )
      .toPromise()
      .then(response => {
        return this.walletService.injectTransaction(response.encoded_transaction).toPromise();
      })
      .then(() => {
        this.resetForm();
        this.button.setSuccess();
        this.walletService.startDataRefreshSubscription();
      })
      .catch(error => {
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
      console.log(value);
      const balance = value && value.coins ? value.coins : 0;
      this.form.get('amount').setValidators([
        Validators.required,
        Validators.max(balance),
        this.validateAmount,
      ]);
      this.form.get('amount').updateValueAndValidity();
    });
  }

  private resetForm() {
    this.form.get('wallet').reset('');
    this.form.get('address').reset('');
    this.form.get('amount').reset('');
    this.form.get('notes').reset('');
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
