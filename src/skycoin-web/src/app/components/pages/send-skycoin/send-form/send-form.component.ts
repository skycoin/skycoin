import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Subscription } from 'rxjs';
import 'rxjs/add/operator/delay';
import 'rxjs/add/operator/filter';
import { BigNumber } from 'bignumber.js';
import { Observable } from 'rxjs';

import { WalletService } from '../../../../services/wallet/wallet.service';
import { SpendingService, HoursSelectionTypes } from '../../../../services/wallet/spending.service';
import { ButtonComponent } from '../../../layout/button/button.component';
import { Wallet, ConfirmationData } from '../../../../app.datatypes';
import { openUnlockWalletModal, showConfirmationModal } from '../../../../utils/index';
import { BaseCoin } from '../../../../coins/basecoin';
import { CoinService } from '../../../../services/coin.service';
import { BlockchainService } from '../../../../services/blockchain.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { config } from '../../../../app.config';
import { NavBarService } from '../../../../services/nav-bar.service';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { PriceService } from '../../../../services/price.service';
import { MsgBarService } from '../../../../services/msg-bar.service';

@Component({
  selector: 'app-send-form',
  templateUrl: './send-form.component.html',
  styleUrls: ['./send-form.component.scss'],
})
export class SendFormComponent implements OnInit, OnDestroy {

  public static readonly MaxUsdDecimal = 6;

  @ViewChild('button') button: ButtonComponent;
  @Input() formData: any;
  @Output() onFormSubmitted = new EventEmitter<any>();

  showSlowMobileInfo = false;
  form: FormGroup;
  wallets: Wallet[];
  currentCoin: BaseCoin;
  doubleButtonActive = DoubleButtonActive;
  selectedCurrency = DoubleButtonActive.LeftButton;
  value: number;
  valueGreaterThanBalance = false;
  price: number;

  private processSubscription: Subscription;
  private subscriptionsGroup: Subscription[] = [];
  private slowInfoSubscription: Subscription;

  constructor(
    public blockchainService: BlockchainService,
    private formBuilder: FormBuilder,
    private walletService: WalletService,
    private spendingService: SpendingService,
    private dialog: CustomMatDialogService,
    private coinService: CoinService,
    private navbarService: NavBarService,
    private msgBarService: MsgBarService,
    priceService: PriceService,
  ) {
    this.subscriptionsGroup.push(priceService.price.subscribe(price => {
      this.price = price;
      this.updateValue();
    }));
  }

  ngOnInit() {
    this.navbarService.showSwitch('send.simple', 'send.advanced');

    this.initForm();

    this.subscriptionsGroup.push(this.walletService.currentWallets
      .subscribe(wallets => this.wallets = wallets)
    );

    this.subscriptionsGroup.push(this.coinService.currentCoin
      .subscribe((coin: BaseCoin) => {
        this.resetForm();
        this.currentCoin = coin;
      })
    );

    if (this.formData) {
      Object.keys(this.form.controls).forEach(control => {
        if (this.form.get(control)) {
          this.form.get(control).setValue(this.formData.form[control]);
        }

        this.selectedCurrency = this.formData.form.currency;
      });
    }
  }

  ngOnDestroy() {
    this.removeSlowInfoSubscription();
    this.removeProcessSubscription();
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
    this.navbarService.hideSwitch();
    this.msgBarService.hide();
  }

  onVerify(event = null) {
    if (event) {
      event.preventDefault();
    }

    if (!this.form.valid || this.button.isLoading()) {
      return;
    }

    this.msgBarService.hide();
    this.button.resetState();

    const wallet = this.form.value.wallet;

    if (!wallet.seed) {
      this.removeProcessSubscription();

      this.processSubscription = openUnlockWalletModal(wallet, this.dialog).componentInstance
        .onWalletUnlocked.first().subscribe(() => this.checkBeforeSending());
    } else {
      this.checkBeforeSending();
    }
  }

  changeActiveCurrency(value) {
    this.selectedCurrency = value;
    this.updateValue();
    this.form.get('amount').updateValueAndValidity();
  }

  private updateValue() {
    if (!this.price) {
      this.value = null;

      return;
    }
    if (!this.form || this.validateAmount(this.form.get('amount') as FormControl) !== null || this.form.get('amount').value * 1 === 0) {
      this.value = -1;

      return;
    }

    const coinsInWallet = this.form.get('wallet').value && (this.form.get('wallet').value as Wallet).balance ?
      (this.form.get('wallet').value as Wallet).balance.toNumber() : -1;

    this.valueGreaterThanBalance = false;
    if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
      this.value = new BigNumber(this.form.get('amount').value).multipliedBy(this.price).decimalPlaces(2).toNumber();
      if (coinsInWallet > 0 && parseFloat(this.form.get('amount').value) > coinsInWallet) {
        this.valueGreaterThanBalance = true;
      }
    } else {
      this.value = new BigNumber(this.form.get('amount').value).dividedBy(this.price).decimalPlaces(this.blockchainService.currentMaxDecimals).toNumber();
      if (coinsInWallet > 0 && this.value > coinsInWallet) {
        this.valueGreaterThanBalance = true;
      }
    }
  }

  private resetForm() {
    this.form.get('wallet').setValue('', { emitEvent: false });
    this.form.get('address').setValue('');
    this.form.get('amount').setValue('');
    this.selectedCurrency = DoubleButtonActive.LeftButton;
  }

  private checkBeforeSending() {
    this.blockchainService.synchronized.first().subscribe(synchronized => {
      if (synchronized) {
        this.createTransaction(this.form.value.wallet);
      } else {
        this.showSynchronizingWarning();
      }
    });
  }

  private showSynchronizingWarning() {
    const confirmationData: ConfirmationData = {
      text: 'send.synchronizing-warning',
      headerText: 'confirmation.header-text',
      confirmButtonText: 'confirmation.confirm-button',
      cancelButtonText: 'confirmation.cancel-button',
    };

    showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
      if (confirmationResult) {
        this.createTransaction(this.form.value.wallet);
      }
    });
  }

  private createTransaction(wallet: Wallet) {
    this.button.setLoading();

    this.slowInfoSubscription = Observable.of(1).delay(config.timeBeforeSlowMobileInfo)
      .subscribe(() => this.showSlowMobileInfo = true);

    this.removeProcessSubscription();
    this.processSubscription = this.spendingService.createTransaction(
      wallet,
      null,
      null,
      [{
        address: this.form.value.address.replace(/\s/g, ''),
        coins: this.selectedCurrency === DoubleButtonActive.LeftButton ? new BigNumber(this.form.value.amount) : new BigNumber(this.value.toString()),
      }],
      {
        type: HoursSelectionTypes.Auto,
        ShareFactor: new BigNumber(0.5),
      },
      null
    ).subscribe(
        transaction => this.onTransactionCreated(transaction),
        error => this.onError(error)
      );
  }

  private onTransactionCreated(transaction) {
    this.showSlowMobileInfo = false;
    this.removeSlowInfoSubscription();
    this.onFormSubmitted.emit({
      form: {
        wallet: this.form.value.wallet,
        address: this.form.value.address,
        amount: this.form.value.amount,
        currency: this.selectedCurrency,
      },
      amount: new BigNumber(this.form.value.amount),
      to: [this.form.value.address],
      transaction,
    });
  }

  private onError(error) {
    this.showSlowMobileInfo = false;
    this.removeSlowInfoSubscription();
    this.msgBarService.showError(error.message);
    this.button.resetState();
  }

  private initForm() {
    this.form = this.formBuilder.group({
      wallet: ['', Validators.required],
      address: ['', Validators.required],
      amount: ['', [Validators.required]]
    });

    this.subscriptionsGroup.push(this.form.controls.wallet.valueChanges.subscribe(value => {
      const balance = value && value.balance ? value.balance : 0;

      this.form.controls.amount.setValidators([
        Validators.required,
        this.validateAmountWithValue.bind(this),
      ]);

      this.form.controls.amount.updateValueAndValidity();
    }));

    this.subscriptionsGroup.push(this.form.get('amount').valueChanges.subscribe(value => {
      this.updateValue();
    }));
  }

  private validateAmount(amountControl: FormControl) {
    if (!amountControl.value || isNaN(amountControl.value) || parseFloat(amountControl.value) <= 0) {
      return { Invalid: true };
    }

    const parts = amountControl.value.toString().split('.');

    if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
      if (parts.length === 2 && parts[1].length > this.blockchainService.currentMaxDecimals) {
        return { Invalid: true };
      }
    } else {
      if (parts.length === 2 && parts[1].length > SendFormComponent.MaxUsdDecimal) {
        return { Invalid: true };
      }
    }

    return null;
  }

  private validateAmountWithValue(amountControl: FormControl) {
    const firstValidation = this.validateAmount(amountControl);
    if (firstValidation) {
      return firstValidation;
    }

    const coinsInWallet = this.form.get('wallet').value && (this.form.get('wallet').value as Wallet).balance ?
      (this.form.get('wallet').value as Wallet).balance.toNumber() : -1;

    if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
      if (parseFloat(amountControl.value) > coinsInWallet) {
        return { Invalid: true };
      }
    } else {
      this.updateValue();
      if (this.value > coinsInWallet) {
        return { Invalid: true };
      }
    }

    return null;
  }

  private removeProcessSubscription() {
    if (this.processSubscription) {
      this.processSubscription.unsubscribe();
    }
  }

  private removeSlowInfoSubscription() {
    if (this.slowInfoSubscription) {
      this.slowInfoSubscription.unsubscribe();
    }
  }
}
