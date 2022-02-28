import { throwError as observableThrowError, SubscriptionLike, concat, of } from 'rxjs';
import { Component, EventEmitter, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import * as moment from 'moment';
import { FormBuilder, FormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { retryWhen, delay, take, mergeMap } from 'rxjs/operators';

import { ButtonComponent } from '../../../layout/button/button.component';
import { ExchangeService, StoredExchangeOrder, TradingPair, ExchangeOrder } from '../../../../services/exchange.service';
import { SelectAddressComponent } from '../../../layout/select-address/select-address.component';
import { AppService } from '../../../../services/app.service';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { WalletsAndAddressesService } from '../../../../services/wallet-operations/wallets-and-addresses.service';

/**
 * Shows the form for creating an exchange order.
 */
@Component({
  selector: 'app-exchange-create',
  templateUrl: './exchange-create.component.html',
  styleUrls: ['./exchange-create.component.scss'],
})
export class ExchangeCreateComponent implements OnInit, OnDestroy {
  // Default coin the user will have to deposit.
  readonly defaultFromCoin = 'BTC';
  // Default amount of coins the user will have to deposit.
  readonly defaultFromAmount = '0.1';
  // Coin the user will receive.
  readonly toCoin = 'SKY';

  @ViewChild('exchangeButton') exchangeButton: ButtonComponent;
  // Event emited when the order has been created.
  @Output() submitted = new EventEmitter<StoredExchangeOrder>();

  form: FormGroup;
  tradingPairs: TradingPair[];
  // Currently selected trading pair
  activeTradingPair: TradingPair;
  problemGettingPairs = false;
  // If true, the form is shown deactivated.
  busy = false;

  // Vars with the validation error messages.
  coinErrorMsg = '';
  amountErrorMsg = '';
  addressErrorMsg = '';
  amountTooLow = false;
  amountTooHight = false;

  // If the user has acepted the agreement.
  private agreement = false;

  private subscriptionsGroup: SubscriptionLike[] = [];
  private exchangeSubscription: SubscriptionLike;
  private priceUpdateSubscription: SubscriptionLike;

  // Approximately how many coins will be received for the amount of coins the user will send,
  // as per the value entered on the form and the current price.
  get toAmount(): string {
    if (!this.activeTradingPair) {
      return '0';
    }

    const fromAmount = this.form.get('fromAmount').value;
    if (isNaN(fromAmount)) {
      return '0';
    } else {
      return (this.form.get('fromAmount').value * this.activeTradingPair.price).toFixed(this.appService.currentMaxDecimals);
    }
  }

  // How many coins the user will send, converted to a valid number.
  get sendAmount(): number {
    const val = this.form.get('fromAmount').value;

    return isNaN(parseFloat(val)) ? 0 : val;
  }

  constructor(
    private exchangeService: ExchangeService,
    private formBuilder: FormBuilder,
    private msgBarService: MsgBarService,
    private dialog: MatDialog,
    private appService: AppService,
    private translateService: TranslateService,
    private walletsAndAddressesService: WalletsAndAddressesService,
  ) { }

  ngOnInit() {
    this.createForm();
    this.loadData();
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
    this.removeExchangeSubscription();
    this.removePriceUpdateSubscription();
    this.msgBarService.hide();
    this.submitted.complete();
  }

  // Called when the user presses the checkbox for acepting the agreement.
  setAgreement(event) {
    this.agreement = event.checked;
    this.form.updateValueAndValidity();
  }

  // Opens the modal window for selecting one of the addresses the user has.
  selectAddress(event) {
    event.stopPropagation();
    event.preventDefault();

    SelectAddressComponent.openDialog(this.dialog).afterClosed().subscribe(address => {
      if (address) {
        this.form.get('toAddress').setValue(address);
      }
    });
  }

  // Creates the order.
  exchange() {
    if (!this.form.valid || this.busy) {
      return;
    }

    // Prepare the UI.
    this.busy = true;
    this.msgBarService.hide();
    this.exchangeButton.resetState();
    this.exchangeButton.setLoading();
    this.exchangeButton.setDisabled();

    const amount = parseFloat(this.form.get('fromAmount').value);

    const toAddress = (this.form.get('toAddress').value as string).trim();

    // Check if the address is valid.
    this.removeExchangeSubscription();
    this.exchangeSubscription = this.walletsAndAddressesService.verifyAddress(toAddress).subscribe(addressIsValid => {
      if (addressIsValid) {
        // Create the order.
        this.exchangeSubscription = this.exchangeService.exchange(
          this.activeTradingPair.pair,
          amount,
          toAddress,
          this.activeTradingPair.price,
        ).subscribe((order: ExchangeOrder) => {
          this.busy = false;
          // Emit the event.
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
          this.exchangeButton.resetState().setEnabled();
          this.msgBarService.showError(err);
        });
      } else {
        this.showInvalidAddress();
      }
    }, () => {
      this.showInvalidAddress();
    });
  }

  // Reactivates the form and shows a msg indicating that the address is invalid.
  private showInvalidAddress() {
    this.busy = false;

    this.exchangeButton.resetState().setEnabled();

    const errMsg = this.translateService.instant('exchange.invalid-address-error');
    this.msgBarService.showError(errMsg);
  }

  // Inits the form.
  private createForm() {
    this.form = this.formBuilder.group({
      fromCoin: [this.defaultFromCoin],
      fromAmount: [this.defaultFromAmount],
      toAddress: [''],
    });

    this.form.setValidators(this.validateForm.bind(this));

    this.subscriptionsGroup.push(this.form.get('fromCoin').valueChanges.subscribe(() => {
      this.updateActiveTradingPair();
    }));
  }

  // Loads the available trading pairs from the backend.
  private loadData() {
    this.subscriptionsGroup.push(this.exchangeService.tradingPairs()
      .pipe(retryWhen(errors => concat(errors.pipe(delay(2000), take(10)), observableThrowError(''))))
      .subscribe(pairs => {
        this.tradingPairs = [];

        // Use only the trading pairs which include the wallet coin.
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

  // Periodically updates the value of each trading pair indicating how many coins will be
  // received per coin sent.
  private updatePrices() {
    this.removePriceUpdateSubscription();

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

  // Updates the var with the currently selected trading pair.
  private updateActiveTradingPair() {
    this.activeTradingPair = this.tradingPairs.find(p => p.from === this.form.get('fromCoin').value);

    if (!this.activeTradingPair && this.tradingPairs.length > 0) {
      this.activeTradingPair = this.tradingPairs[0];
      this.form.get('fromCoin').setValue(this.activeTradingPair.from);
    }
  }

  /**
   * Validates the form and updates the vars with the validation errors.
   */
  validateForm() {
    this.coinErrorMsg = '';
    this.amountErrorMsg = '';
    this.addressErrorMsg = '';
    this.amountTooLow = false;
    this.amountTooHight = false;

    if (!this.activeTradingPair) {
      return null;
    }

    let valid = true;

    const fromAmount = this.form.get('fromAmount').value;

    // There must be a from amount.
    if (!fromAmount || isNaN(fromAmount)) {
      valid = false;
      if (this.form.get('fromAmount').touched) {
        this.amountErrorMsg = 'exchange.invalid-value-error-info';
      }
    } else {
      const parts = (fromAmount as string).split('.');

      // If there is a from amount, it must not have more than 6 decimals.
      if (parts.length > 1 && parts[1].length > 6) {
        valid = false;
        if (this.form.get('fromAmount').touched) {
          this.amountErrorMsg = 'exchange.invalid-value-error-info';
        }
      }
    }

    // If there is a from amount, it must be inside the limits.
    if (valid) {
      if (fromAmount < this.activeTradingPair.min) {
        this.amountTooLow = true;
        valid = false;
        if (this.form.get('fromAmount').touched) {
          this.amountErrorMsg = 'exchange.invalid-value-error-info';
        }
      }

      if (fromAmount > this.activeTradingPair.max) {
        this.amountTooHight = true;
        valid = false;
        if (this.form.get('fromAmount').touched) {
          this.amountErrorMsg = 'exchange.invalid-value-error-info';
        }
      }
    }

    // There must be a selected coin for the from amount.
    if (!this.form.get('fromCoin').value) {
      valid = false;
      if (this.form.get('fromCoin').touched) {
        this.coinErrorMsg = 'exchange.from-coin-error-info';
      }
    }

    // There must be a valid destination address.
    const address = this.form.get('toAddress').value as string;
    if (!address || address.length < 20) {
      valid = false;
      if (this.form.get('toAddress').touched) {
        this.addressErrorMsg = 'exchange.address-error-info';
      }
    }

    // The user must accept the agreement.
    if (!this.agreement) {
      valid = false;
    }

    return valid ? null : { Invalid: true };
  }

  private removeExchangeSubscription() {
    if (this.exchangeSubscription) {
      this.exchangeSubscription.unsubscribe();
    }
  }

  private removePriceUpdateSubscription() {
    if (this.priceUpdateSubscription) {
      this.priceUpdateSubscription.unsubscribe();
    }
  }
}
