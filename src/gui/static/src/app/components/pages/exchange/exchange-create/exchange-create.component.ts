import {
  Component,
  EventEmitter,
  OnDestroy,
  OnInit,
  Output,
  ViewChild,
} from '@angular/core';
import * as moment from 'moment';
import { ButtonComponent } from '../../../layout/button/button.component';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ExchangeService } from '../../../../services/exchange.service';
import { ExchangeOrder, TradingPair, StoredExchangeOrder } from '../../../../app.datatypes';
import { ISubscription } from 'rxjs/Subscription';
import 'rxjs/add/observable/merge';
import { MatDialog, MatDialogConfig } from '@angular/material';
import { SelectAddressComponent } from '../../send-skycoin/send-form-advanced/select-address/select-address';
import { WalletService } from '../../../../services/wallet.service';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/switchMap';
import { BlockchainService } from '../../../../services/blockchain.service';
import { TranslateService } from '@ngx-translate/core';
import { MsgBarService } from '../../../../services/msg-bar.service';

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
  @Output() submitted = new EventEmitter<StoredExchangeOrder>();
  form: FormGroup;
  tradingPairs: TradingPair[];
  activeTradingPair: TradingPair;
  problemGettingPairs = false;

  private agreement = false;
  private subscriptionsGroup: ISubscription[] = [];
  private exchangeSubscription: ISubscription;
  private priceUpdateSubscription: ISubscription;

  get toAmount() {
    if (!this.activeTradingPair) {
      return 0;
    }

    const fromAmount = this.form.get('fromAmount').value;
    if (isNaN(fromAmount)) {
      return 0;
    } else {
      return (this.form.get('fromAmount').value * this.activeTradingPair.price).toFixed(this.blockchainService.currentMaxDecimals);
    }
  }

  get sendAmount() {
    const val = this.form.get('fromAmount').value;

    return isNaN(parseFloat(val)) ? 0 : val;
  }

  constructor(
    private exchangeService: ExchangeService,
    private walletService: WalletService,
    private formBuilder: FormBuilder,
    private msgBarService: MsgBarService,
    private dialog: MatDialog,
    private blockchainService: BlockchainService,
    private translateService: TranslateService,
  ) { }

  ngOnInit() {
    this.createForm();
    this.loadData();
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
    this.removeExchangeSubscription();
    this.msgBarService.hide();

    if (this.priceUpdateSubscription) {
      this.priceUpdateSubscription.unsubscribe();
    }
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
    if (!this.form.valid || this.exchangeButton.isLoading()) {
      return;
    }

    this.exchangeButton.resetState();
    this.exchangeButton.setLoading();
    this.exchangeButton.setDisabled();

    const amount = parseFloat(this.form.get('fromAmount').value);

    const toAddress = (this.form.get('toAddress').value as string).trim();

    this.removeExchangeSubscription();
    this.exchangeSubscription = this.walletService.verifyAddress(toAddress).subscribe(addressIsValid => {
      if (addressIsValid) {
        this.exchangeSubscription = this.exchangeService.exchange(
          this.activeTradingPair.pair,
          amount,
          toAddress,
          this.activeTradingPair.price,
        ).subscribe((order: ExchangeOrder) => {
          this.submitted.emit({
            id: order.id,
            pair: order.pair,
            fromAmount: order.fromAmount,
            toAmount: order.toAmount,
            address: order.toAddress,
            timestamp: moment().unix(),
            price: this.activeTradingPair.price,
          });
        }, err => {
          this.exchangeButton.resetState();
          this.exchangeButton.setEnabled();
          this.exchangeButton.setError(err);
          this.msgBarService.showError(err);
        });
      } else {
        this.showInvalidAddress();
      }
    }, () => {
      this.showInvalidAddress();
    });
  }

  private showInvalidAddress() {
    this.exchangeButton.resetState();
    this.exchangeButton.setEnabled();

    const errMsg = this.translateService.instant('exchange.invalid-address');
    this.exchangeButton.setError(errMsg);
    this.msgBarService.showError(errMsg);
  }

  private createForm() {
    this.form = this.formBuilder.group({
      fromCoin: [this.defaultFromCoin, Validators.required],
      fromAmount: [this.defaultFromAmount, Validators.required],
      toAddress: ['', Validators.required],
    }, {
      validator: this.validate.bind(this),
    });

    this.subscriptionsGroup.push(this.form.get('fromCoin').valueChanges.subscribe(() => {
      this.updateActiveTradingPair();
    }));
  }

  private loadData() {
    this.subscriptionsGroup.push(this.exchangeService.tradingPairs()
      .retryWhen(errors => errors.delay(2000).take(10).concat(Observable.throw('')))
      .subscribe(pairs => {
        this.tradingPairs = [];

        pairs.forEach(pair => {
          if (pair.to === this.toCoin) {
            this.tradingPairs.push(pair);
          }
        });

        this.updateActiveTradingPair();
        this.updatePrices();
      }, () => {
        this.problemGettingPairs = true;
      }),
    );
  }

  private updatePrices() {
    this.priceUpdateSubscription = Observable.of(1).delay(60000).flatMap(() => this.exchangeService.tradingPairs())
      .retryWhen(errors => errors.delay(60000))
      .subscribe(pairs => {
        pairs.forEach(pair => {
          if (pair.to === this.toCoin) {
            const alreadySavedPair = this.tradingPairs.find(oldPair => oldPair.from === pair.from);
            if (alreadySavedPair) {
              alreadySavedPair.price = pair.price;
            }
          }
        });
        this.updatePrices();
      });
  }

  private updateActiveTradingPair() {
    this.activeTradingPair = this.tradingPairs.find(p => {
      return p.from === this.form.get('fromCoin').value;
    });

    if (!this.activeTradingPair && this.tradingPairs.length > 0) {
      this.activeTradingPair = this.tradingPairs[0];
      this.form.get('fromCoin').setValue(this.activeTradingPair.from);
    }
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

  private removeExchangeSubscription() {
    if (this.exchangeSubscription) {
      this.exchangeSubscription.unsubscribe();
    }
  }
}
