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
import { WalletService } from '../../../../services/wallet.service';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/switchMap';
import { ApiService } from '../../../../services/api.service';

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
  @Output() shownLast = new EventEmitter();
  form: FormGroup;
  tradingPairs: TradingPair[];
  activeTradingPair: TradingPair;

  private agreement = false;
  private subscription: ISubscription;
  private decimals = 0;

  get toAmount() {
    if (!this.activeTradingPair) {
      return 0;
    }

    return (this.form.get('fromAmount').value * this.activeTradingPair.price).toFixed(this.decimals);
  }

  get sendAmount() {
    const val = this.form.get('fromAmount').value;

    return isNaN(parseFloat(val)) ? 0 : val;
  }

  constructor(
    private exchangeService: ExchangeService,
    private walletService: WalletService,
    private apiService: ApiService,
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
      this.exchangeService.lastOrder = order;
      this.submitted.emit(order);
    }, err => {
      this.exchangeButton.resetState();
      this.exchangeButton.setEnabled();
      this.exchangeButton.setError(err);
      showSnackbarError(this.snackbar, err);
    });
  }

  hasLast() {
    return !!this.exchangeService.lastOrder;
  }

  showLast() {
    return this.shownLast.emit();
  }

  private createForm() {
    this.form = this.formBuilder.group({
      fromCoin: [this.defaultFromCoin, Validators.required],
      fromAmount: [this.defaultFromAmount, Validators.required],
      toAddress: ['', Validators.required],
    }, {
      validator: this.validate.bind(this),
      asyncValidator: this.validateAddress.bind(this),
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

    this.apiService.getHealth().subscribe(res => {
      this.decimals = res.user_verify_transaction.max_decimals;
    });
  }

  private updateActiveTradingPair() {
    this.activeTradingPair = this.tradingPairs.find(p => {
      return p.from === this.form.get('fromCoin').value;
    });
  }

  private validate(group: FormGroup) {
    if (!group || !this.activeTradingPair) {
      return null;
    }

    const fromAmount = group.get('fromAmount').value;

    if (isNaN(fromAmount)) {
      return { invalid: true };
    }

    if (fromAmount < this.activeTradingPair.min || fromAmount === '') {
      return { min: this.activeTradingPair.min };
    }

    if (fromAmount > this.activeTradingPair.max) {
      return { max: this.activeTradingPair.max };
    }

    const parts = (fromAmount as string).split('.');

    if (parts.length > 1 && parts[1].length > 6) {
      return { decimals: true };
    }

    if (!this.agreement) {
      return { agreement: true };
    }

    return null;
  }

  private validateAddress() {
    const address = this.form.get('toAddress').value;

    if (!address) {
      return Observable.create({ address: true });
    }

    return Observable
      .timer(500)
      .switchMap(() => this.walletService.verifyAddress(address))
      .map(res => res ? null : { address : true });
  }
}
