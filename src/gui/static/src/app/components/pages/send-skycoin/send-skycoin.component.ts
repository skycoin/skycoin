import { Component, OnInit, ViewChild } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { IntervalObservable } from 'rxjs/observable/IntervalObservable';
import { Router } from '@angular/router';
import { MdDialogRef, MdSnackBar, MdSnackBarConfig } from '@angular/material';
import 'rxjs/add/operator/delay';
import 'rxjs/add/operator/filter';
import { ButtonComponent } from '../../layout/button/button.component';

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
    public dialogRef: MdDialogRef<SendSkycoinComponent>,
    private snackbar: MdSnackBar,
  ) {}

  ngOnInit() {
    this.initForm();
    IntervalObservable
      .create(2500)
      .filter(() => !!this.transactions.length)
      .flatMap(() => this.walletService.retrieveUpdatedTransactions(this.transactions))
      .subscribe(transactions => this.records = transactions);
    this.walletService.recent().subscribe(transactions => this.transactions = transactions);
  }

  closePopup() {
    this.dialogRef.close();
  }

  send() {
    this.button.setLoading();
    this.walletService.sendSkycoin(this.form.value.wallet, this.form.value.address, this.form.value.amount * 1000000)
      .delay(1000)
      .subscribe(
        () => {
          this.resetForm();
          this.button.setSuccess();
          this.dialogRef.close();
        },
        error => {
          const config = new MdSnackBarConfig();
          config.duration = 300000;
          this.snackbar.open(error['_body'], null, config);
          this.button.setError(error);
        }
      );
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
