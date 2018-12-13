import {
  Component,
  EventEmitter,
  OnDestroy,
  OnInit,
  Output,
  ViewChild,
} from '@angular/core';
import { ButtonComponent } from '../../../layout/button/button.component';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ExchangeService } from '../../../../services/exchange.service';
import { ExchangeOrder, TradingPair } from '../../../../app.datatypes';
import { ISubscription } from 'rxjs/Subscription';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/merge';
import { MatDialog, MatDialogConfig, MatSnackBar } from '@angular/material';
import { showSnackbarError } from '../../../../utils/errors';
import { SelectAddressComponent } from '../../send-skycoin/send-form-advanced/select-address/select-address';

@Component({
  selector: 'app-exchange-create',
  templateUrl: './exchange-create.component.html',
  styleUrls: ['./exchange-create.component.scss'],
})
export class ExchangeCreateComponent implements OnInit, OnDestroy {
  @ViewChild('exchangeButton') exchangeButton: ButtonComponent;
  @Output() submitted = new EventEmitter<ExchangeOrder>();
  form: FormGroup;
  tradingPairs: TradingPair[];
  availableFrom = new Set<string>();
  availableTo = new Set<string>();
  activeTradingPair: TradingPair;
  agreement = false;
  subscription: ISubscription;

  readonly defaultFrom = 'BTC';
  readonly defaultFromAmount = '0.1';
  readonly defaultTo = 'SKY';

  get convertedAmount() {
    return (this.form.get('fromAmount').value * this.activeTradingPair.price).toFixed(6);
  }

  get sendAmount() {
    const val = this.form.get('fromAmount').value;

    return isNaN(parseFloat(val)) ? 0 : val;
  }

  get coinName() {
    return this.activeTradingPair ? this.activeTradingPair.to : this.defaultTo;
  }

  constructor(
    private exchangeService: ExchangeService,
    private formBuilder: FormBuilder,
    private snackbar: MatSnackBar,
    private dialog: MatDialog,
  ) { }

  ngOnInit() {
    this.createForm();
    this.loadData();
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
    this.snackbar.dismiss();
  }

  setAgreement(event) {
    this.agreement = event.checked;
    this.form.updateValueAndValidity();
  }

  selectAddress() {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.autoFocus = false;

    this.dialog.open(SelectAddressComponent, config).afterClosed().subscribe(address => {
      if (address) {
        this.form.get('toAddress').setValue(address);
      }
    });
  }

  exchange() {
    this.exchangeButton.resetState();
    this.exchangeButton.setLoading();
    this.exchangeButton.setDisabled();

    const amount = parseFloat(this.form.get('fromAmount').value);

    this.exchangeService.exchange(
      this.activeTradingPair.pair,
      amount,
      this.form.get('toAddress').value,
    ).subscribe((order: ExchangeOrder) => {
      this.exchangeService.setLastOrder(order);
      this.submitted.emit(order);
    }, err => {
      this.exchangeButton.resetState();
      this.exchangeButton.setEnabled();
      this.exchangeButton.setError(err);
      showSnackbarError(this.snackbar, err);
    });
  }

  private loadData() {
    this.exchangeService.tradingPairs().subscribe(pairs => {
      this.tradingPairs = pairs;

      this.tradingPairs.forEach(pair => {
        this.availableFrom.add(pair.from);

        if (pair.from === this.defaultFrom) {
          this.availableTo.add(pair.to);
        }
      });

      this.activeTradingPair = this.tradingPairs.find(p => {
        return p.pair === `${this.defaultFrom}/${this.defaultTo}`;
      });

      this.updateToAmount();
    });
  }

  private createForm() {
    this.form = this.formBuilder.group({
      from: [this.defaultFrom, Validators.required],
      fromAmount: [this.defaultFromAmount, Validators.required],
      to: [this.defaultTo, Validators.required],
      toAmount: ['', Validators.required],
      toAddress: ['', Validators.required],
    }, {
      validator: this.validate.bind(this),
    });

    this.subscription = Observable.merge(
      this.form.get('from').valueChanges,
      this.form.get('fromAmount').valueChanges,
      this.form.get('to').valueChanges,
    ).subscribe(() => {
      this.activeTradingPair = this.tradingPairs.find(p => {
        return p.pair === `${this.form.get('from').value}/${this.form.get('to').value}`;
      });

      if (!this.activeTradingPair) {
        this.activeTradingPair = this.tradingPairs.find(p => p.from === this.form.get('from').value);
        this.form.get('to').setValue(this.activeTradingPair.to);
      }

      this.updateToAmount();
    });
  }

  private updateToAmount() {
    this.form.get('toAmount').setValue(this.convertedAmount);
  }

  private validate(group: FormGroup) {
    if (!group || !this.activeTradingPair) {
      return null;
    }

    const fromAmount = group.get('fromAmount').value;

    if (isNaN(parseFloat(fromAmount))) {
      return { invalid: true };
    }

    if (fromAmount < this.activeTradingPair.min) {
      return { min: this.activeTradingPair.min };
    }

    if (fromAmount > this.activeTradingPair.max) {
      return { max: this.activeTradingPair.max };
    }

    if (!this.agreement) {
      return { agreement: true };
    }

    return null;
  }
}
