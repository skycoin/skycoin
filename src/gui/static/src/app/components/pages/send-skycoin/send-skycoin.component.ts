import { Component, OnInit, ViewChild } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { IntervalObservable } from 'rxjs/observable/IntervalObservable';
import { Router } from '@angular/router';
import 'rxjs/add/operator/delay';

@Component({
  selector: 'app-send-skycoin',
  templateUrl: './send-skycoin.component.html',
  styleUrls: ['./send-skycoin.component.css']
})
export class SendSkycoinComponent implements OnInit {
  @ViewChild('button') button;

  form: FormGroup;
  records = [];
  transactions = [];

  constructor(
    public formBuilder: FormBuilder,
    public walletService: WalletService,
    private router: Router,
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

  onActivate(response) {
    if (response.row && response.row.txid) {
      this.router.navigate(['/history', response.row.txid]);
    }
  }

  send() {
    this.button.setLoading();
    this.walletService.sendSkycoin(this.form.value.wallet_id, this.form.value.address, this.form.value.amount * 1000000)
      .delay(1000)
      .subscribe(
        response => {
          this.resetForm();
          this.button.setSuccess();
        },
        error => this.button.setError(error)
      );
  }

  private initForm() {
    this.form = this.formBuilder.group({
      wallet_id: ['', Validators.required],
      address: ['', Validators.required],
      amount: ['', Validators.required],
    });
  }

  private resetForm() {
    this.form.controls.wallet_id.reset(undefined);
    this.form.controls.address.reset(undefined);
    this.form.controls.amount.reset(undefined);
  }
}
