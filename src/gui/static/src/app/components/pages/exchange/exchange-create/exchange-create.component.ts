import { throwError as observableThrowError, SubscriptionLike, Observable, of } from 'rxjs';
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
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { SelectAddressComponent } from '../../../layout/select-address/select-address.component';
import { WalletService } from '../../../../services/wallet.service';
import { BlockchainService } from '../../../../services/blockchain.service';
import { TranslateService } from '@ngx-translate/core';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { retryWhen, delay, take, concat, mergeMap } from 'rxjs/operators';

@Component({
  selector: 'app-exchange-create',
  templateUrl: './exchange-create.component.html',
  styleUrls: ['./exchange-create.component.scss'],
})
export class ExchangeCreateComponent implements OnInit, OnDestroy {
  readonly defaultFromCoin = 'BTC';
  readonly defaultFromAmount = '0.1';
  readonly toCoin = 'SKY';

  @ViewChild('exchangeButton', { static: false }) exchangeButton: ButtonComponent;
  @Output() submitted = new EventEmitter<StoredExchangeOrder>();
  form: FormGroup;
  tradingPairs: TradingPair[];
  activeTradingPair: TradingPair;
  problemGettingPairs = false;
  busy = false;

  private agreement = false;
  private subscriptionsGroup: SubscriptionLike[] = [];
  private exchangeSubscription: SubscriptionLike;
  private priceUpdateSubscription: SubscriptionLike;

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

  selectAddress(event) {
    event.stopPropagation();
    event.preventDefault();

    SelectAddressComponent.openDialog(this.dialog).afterClosed().subscribe(address => {
      if (address) {
        this.form.get('toAddress').setValue(address);
      }
    });
  }

  exchange() {
    if (!this.form.valid || this.busy) {
      return;
    }

    this.busy = true;
    this.msgBarService.hide();

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
          this.busy = false;
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
          this.busy = false;
          this.exchangeButton.resetState();
          this.exchangeButton.setEnabled();
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
    this.busy = false;

    this.exchangeButton.resetState();
    this.exchangeButton.setEnabled();

    const errMsg = this.translateService.instant('exchange.invalid-address');
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
      .pipe(retryWhen(errors => errors.pipe(delay(2000), take(10), concat(observableThrowError('')))))
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
    this.priceUpdateSubscription = of(1).pipe(delay(60000), mergeMap(() => this.exchangeService.tradingPairs()),
      retryWhen(errors => errors.pipe(delay(60000))))
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
