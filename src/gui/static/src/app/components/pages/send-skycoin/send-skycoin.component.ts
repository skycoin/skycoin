import { Component, OnInit, ViewChild } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatSnackBar, MatSnackBarConfig } from '@angular/material/snack-bar';
import 'rxjs/add/operator/delay';
import 'rxjs/add/operator/filter';
import { ButtonComponent } from '../../layout/button/button.component';
import { PasswordDialogComponent } from '../../layout/password-dialog/password-dialog.component';
import { MatDialog } from '@angular/material';

@Component({
  selector: 'app-send-skycoin',
  templateUrl: './send-skycoin.component.html',
  styleUrls: ['./send-skycoin.component.scss']
})
export class SendSkycoinComponent implements OnInit {
  @ViewChild('button') button: ButtonComponent;

  form: FormGroup;
  records = [];
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

  send() {
    this.button.setLoading();

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
    this.walletService.sendSkycoin(
      this.form.value.wallet,
      this.form.value.address,
      this.form.value.amount * 1000000,
      passwordDialog ? passwordDialog.password : null
    )
      .delay(1000)
      .subscribe(
        () => {
          this.resetForm();
          this.button.setSuccess();
        },
        error => {
          const config = new MatSnackBarConfig();
          config.duration = 300000;
          this.snackbar.open(error['_body'], null, config);
          this.button.setError(error);
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
    this.form.controls['wallet'].valueChanges.subscribe(value => {
      console.log(value);
      const balance = value && value.coins ? value.coins : 0;
      this.form.controls['amount'].setValidators([
        Validators.required,
        Validators.min(0),
        Validators.max(balance),
      ]);
      this.form.controls['amount'].updateValueAndValidity();
    });
  }

  private resetForm() {
    this.form.controls.wallet.reset(undefined);
    this.form.controls.address.reset(undefined);
    this.form.controls.amount.reset(undefined);
  }
}
