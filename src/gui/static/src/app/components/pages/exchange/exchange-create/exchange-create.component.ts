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
  readonly defaultFromCoin = 'BTC';
  readonly defaultFromAmount = '0.1';
  readonly toCoin = 'SKY';

  @ViewChild('exchangeButton') exchangeButton: ButtonComponent;
  @Output() submitted = new EventEmitter<ExchangeOrder>();
  form: FormGroup;
  tradingPairs: TradingPair[];
  activeTradingPair: TradingPair;
  agreement = false;
  subscription: ISubscription;

  get toAmount() {
    if (!this.activeTradingPair) {
      return 0;
    }

    return (this.form.get('fromAmount').value * this.activeTradingPair.price).toFixed(6);
  }

  get sendAmount() {
    const val = this.form.get('fromAmount').value;

    return isNaN(parseFloat(val)) ? 0 : val;
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

  private createForm() {
    this.form = this.formBuilder.group({
      fromCoin: [this.defaultFromCoin, Validators.required],
      fromAmount: [this.defaultFromAmount, Validators.required],
      toAddress: ['', Validators.required],
    }, {
      validator: this.validate.bind(this),
    });

    this.subscription = this.form.get('fromCoin').valueChanges.subscribe(() => {
      this.updateActiveTradingPair();
    });
  }

  private loadData() {
    this.exchangeService.tradingPairs().subscribe(pairs => {
      this.tradingPairs = [];

      pairs.forEach(pair => {
        if (pair.to === this.toCoin) {
          this.tradingPairs.push(pair);
        }
      });

      this.updateActiveTradingPair();
    });
  }

  private updateActiveTradingPair() {
    this.activeTradingPair = this.tradingPairs.find(p => {
      return p.from === this.defaultFromCoin || p.from === this.form.get('fromCoin').value;
    });
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
