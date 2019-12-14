import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild, ChangeDetectorRef } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import 'rxjs/add/operator/delay';
import 'rxjs/add/operator/filter';
import { ButtonComponent } from '../../../layout/button/button.component';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { getHardwareWalletErrorMsg } from '../../../../utils/errors';
import { ISubscription } from 'rxjs/Subscription';
import { NavBarService } from '../../../../services/nav-bar.service';
import { BigNumber } from 'bignumber.js';
import { Wallet, ConfirmationData } from '../../../../app.datatypes';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { TranslateService } from '@ngx-translate/core';
import { BlockchainService } from '../../../../services/blockchain.service';
import { showConfirmationModal } from '../../../../utils';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { PriceService } from '../../../../services/price.service';
import { ChangeNoteComponent } from '../send-preview/transaction-info/change-note/change-note.component';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { AppService } from '../../../../services/app.service';

@Component({
  selector: 'app-send-form',
  templateUrl: './send-form.component.html',
  styleUrls: ['./send-form.component.scss'],
})
export class SendFormComponent implements OnInit, OnDestroy {

  public static readonly MaxUsdDecimals = 6;

  @ViewChild('previewButton', { static: false }) previewButton: ButtonComponent;
  @ViewChild('sendButton', { static: false }) sendButton: ButtonComponent;
  @Input() formData: any;
  @Output() onFormSubmitted = new EventEmitter<any>();

  maxNoteChars = ChangeNoteComponent.MAX_NOTE_CHARS;
  form: FormGroup;
  transactions = [];
  previewTx: boolean;
  busy = false;
  doubleButtonActive = DoubleButtonActive;
  selectedCurrency = DoubleButtonActive.LeftButton;
  value: number;
  valueGreaterThanBalance = false;
  price: number;
  wallets: Wallet[];

  private subscriptionsGroup: ISubscription[] = [];
  private processingSubscription: ISubscription;
  private syncCheckSubscription: ISubscription;

  constructor(
    public formBuilder: FormBuilder,
    public blockchainService: BlockchainService,
    private walletService: WalletService,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    private navbarService: NavBarService,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
    private changeDetector: ChangeDetectorRef,
    private appService: AppService,
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
    this.subscriptionsGroup.push(this.walletService.all().first().subscribe(wallets => {
      this.wallets = wallets;

      if (wallets.length === 1) {
        this.form.get('wallet').setValue(wallets[0]);
      }
    }));
  }

  ngOnDestroy() {
    if (this.processingSubscription && !this.processingSubscription.closed) {
      this.processingSubscription.unsubscribe();
    }
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
    this.closeSyncCheckSubscription();
    this.navbarService.hideSwitch();
    this.msgBarService.hide();
  }

  preview() {
    this.previewTx = true;
    this.checkBeforeSending();
    this.changeDetector.detectChanges();
  }

  send() {
    this.previewTx = false;
    this.checkBeforeSending();
  }

  changeActiveCurrency(value) {
    if (value !== this.selectedCurrency) {
      this.selectedCurrency = value;
      this.askIfConvertAmount();
      this.updateValue();
      this.form.get('amount').updateValueAndValidity();
    }
  }

  private askIfConvertAmount() {

    if (!this.form.get('amount').value) {
      return;
    }
    const value = (this.form.get('amount').value as string).trim();
    const currentValue = new BigNumber((this.form.get('amount').value as string).trim());
    if (!value || currentValue.isNaN()) {
      return;
    }

    const usd = this.translate.instant('common.usd');
    const currentCoin = this.appService.coinName;
    let fromText: string;
    let toText: string;
    if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
      fromText = usd;
      toText = currentCoin;
    } else {
      fromText = currentCoin;
      toText = usd;
    }

    const confirmationData: ConfirmationData = {
      text: this.translate.instant('send.convert-confirmation', {from: fromText, to: toText}),
      headerText: 'confirmation.header-text',
      confirmButtonText: 'confirmation.confirm-button',
      cancelButtonText: 'confirmation.cancel-button',
    };

    showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
      if (confirmationResult) {
        this.convertAmount();
      }
    });
  }

  private convertAmount() {
    this.msgBarService.hide();

    if (this.form.get('amount').value) {
      const value = (this.form.get('amount').value as string).trim();
      const currentValue = new BigNumber(value);

      if (!value || currentValue.isNaN()) {
        this.msgBarService.showWarning(this.translate.instant('send.invaid-amount-warning'));

        return;
      }

      if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
        const newValue = currentValue.dividedBy(this.price).decimalPlaces(this.blockchainService.currentMaxDecimals);
        const recoveredValue = newValue.multipliedBy(this.price).decimalPlaces(SendFormComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
        if (!recoveredValue.isEqualTo(currentValue)) {
          this.msgBarService.showWarning(this.translate.instant('send.precision-error-warning'));
        }

        this.form.get('amount').setValue(newValue.toString());
      } else {
        const newValue = currentValue.multipliedBy(this.price).decimalPlaces(SendFormComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
        const recoveredValue = newValue.dividedBy(this.price).decimalPlaces(this.blockchainService.currentMaxDecimals);
        if (!recoveredValue.isEqualTo(currentValue)) {
          this.msgBarService.showWarning(this.translate.instant('send.precision-error-warning'));
        }

        this.form.get('amount').setValue(newValue.toString());
      }
    }
  }

  assignAll() {
    this.msgBarService.hide();

    let availableCoins: BigNumber = this.form.get('wallet').value && this.form.get('wallet').value.coins ? this.form.get('wallet').value.coins : new BigNumber(-1);
    if ((availableCoins as BigNumber).isEqualTo(-1)) {
      this.msgBarService.showError(this.translate.instant('send.no-wallet-selected'));

      return;
    }

    if (this.selectedCurrency === DoubleButtonActive.RightButton) {
      availableCoins = availableCoins.multipliedBy(this.price).decimalPlaces(SendFormComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
    }

    this.form.get('amount').setValue(availableCoins.toString());
  }

  private updateValue() {
    if (!this.price) {
      this.value = null;

      return;
    }
    if (!this.form || this.validateAmount(this.form.get('amount') as FormControl) !== null) {
      this.value = -1;

      return;
    }

    const coinsInWallet = this.form.get('wallet').value && this.form.get('wallet').value.coins ? this.form.get('wallet').value.coins : -1;

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

  private checkBeforeSending() {
    if (!this.form.valid || this.previewButton.isLoading() || this.sendButton.isLoading()) {
      return;
    }

    this.closeSyncCheckSubscription();
    this.syncCheckSubscription = this.blockchainService.synchronized.first().subscribe(synchronized => {
      if (synchronized) {
        this.prepareTransaction();
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
        this.prepareTransaction();
      }
    });
  }

  private prepareTransaction() {
    this.msgBarService.hide();
    this.previewButton.resetState();
    this.sendButton.resetState();

    if (this.form.value.wallet.encrypted && !this.form.value.wallet.isHardware && !this.previewTx) {
      const config = new MatDialogConfig();
      config.data = {
        wallet: this.form.value.wallet,
      };

      this.dialog.open(PasswordDialogComponent, config).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.createTransaction(passwordDialog);
        });
    } else {
      if (!this.form.value.wallet.isHardware || this.previewTx) {
        this.createTransaction();
      } else {
        this.showBusy();
        this.processingSubscription = this.hwWalletService.checkIfCorrectHwConnected((this.form.value.wallet as Wallet).addresses[0].address).subscribe(
          () => this.createTransaction(),
          err => this.showError(getHardwareWalletErrorMsg(this.translate, err)),
        );
      }
    }
  }

  private showBusy() {
    if (this.previewTx) {
      this.previewButton.setLoading();
      this.sendButton.setDisabled();
    } else {
      this.sendButton.setLoading();
      this.previewButton.setDisabled();
    }
    this.busy = true;
    this.navbarService.disableSwitch();
  }

  private createTransaction(passwordDialog?: any) {
    this.showBusy();

    this.processingSubscription = this.walletService.createTransaction(
      this.form.value.wallet,
      (this.form.value.wallet as Wallet).addresses.map(address => address.address),
      null,
      [{
        address: (this.form.value.address as string).trim(),
        coins: ((this.selectedCurrency === DoubleButtonActive.LeftButton ? this.form.value.amount : this.value.toString()) as string).trim(),
      }],
      {
        type: 'auto',
        mode: 'share',
        share_factor: '0.5',
      },
      null,
      passwordDialog ? passwordDialog.password : null,
      this.previewTx,
    ).subscribe(transaction => {
        if (passwordDialog) {
          passwordDialog.close();
        }

        const note = this.form.value.note.trim();
        if (!this.previewTx) {
          this.processingSubscription = this.walletService.injectTransaction(transaction.encoded, note)
            .subscribe(noteSaved => {
              let showDone = true;
              if (note && !noteSaved) {
                this.msgBarService.showWarning(this.translate.instant('send.error-saving-note'));
                showDone = false;
              }

              this.showSuccess(showDone);
            }, error => this.showError(error));
        } else {
          this.onFormSubmitted.emit({
            form: {
              wallet: this.form.value.wallet,
              address: (this.form.value.address as string).trim(),
              amount: this.form.value.amount,
              currency: this.selectedCurrency,
              note: note,
            },
            amount: new BigNumber(this.form.value.amount),
            to: [(this.form.value.address as string).trim()],
            transaction,
          });
          this.busy = false;
          this.navbarService.enableSwitch();
        }
      },
      error => {
        if (passwordDialog) {
          passwordDialog.error(error);
        }

        if (error && error.result) {
          this.showError(getHardwareWalletErrorMsg(this.translate, error));
        } else {
          this.showError(error);
        }
      },
    );
  }

  private showSuccess(showDone: boolean) {
    this.busy = false;
    this.navbarService.enableSwitch();
    this.resetForm();

    if (showDone) {
      this.msgBarService.showDone('send.sent');
      this.sendButton.resetState();
    } else {
      this.sendButton.setSuccess();
      setTimeout(() => {
        this.sendButton.resetState();
      }, 3000);
    }
  }

  private showError(error) {
    this.busy = false;
    this.msgBarService.showError(error);
    this.navbarService.enableSwitch();
    this.previewButton.resetState().setEnabled();
    this.sendButton.resetState().setEnabled();
  }

  private initForm() {
    this.form = this.formBuilder.group({
      wallet: ['', Validators.required],
      address: [''],
      amount: ['', Validators.required],
      note: [''],
    });

    this.form.get('address').setValidators([
      this.validateAddress.bind(this),
    ]);

    this.subscriptionsGroup.push(this.form.get('wallet').valueChanges.subscribe(value => {
      this.form.get('amount').setValidators([
        Validators.required,
        this.validateAmountWithValue.bind(this),
      ]);

      this.form.get('amount').updateValueAndValidity();
    }));

    this.subscriptionsGroup.push(this.form.get('amount').valueChanges.subscribe(value => {
      this.updateValue();
    }));

    if (this.formData) {
      Object.keys(this.form.controls).forEach(control => {
        if (this.form.get(control)) {
          this.form.get(control).setValue(this.formData.form[control]);
        }

        this.selectedCurrency = this.formData.form.currency;
      });
    }
  }

  private validateAddress(addressControl: FormControl) {
    if (!addressControl.value || (addressControl.value as string).trim().length === 0) {
      return { Required: true };
    }
  }

  private validateAmount(amountControl: FormControl) {
    let stringValue: string = amountControl.value;
    stringValue = stringValue ? stringValue.trim() : stringValue;
    const value = new BigNumber(stringValue);

    if (!stringValue || value.isNaN() || value.isLessThanOrEqualTo(0)) {
      return { Invalid: true };
    }

    const parts = stringValue.split('.');

    if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
      if (parts.length === 2 && parts[1].length > this.blockchainService.currentMaxDecimals) {
        return { Invalid: true };
      }
    } else {
      if (parts.length === 2 && parts[1].length > SendFormComponent.MaxUsdDecimals) {
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

    const coinsInWallet = this.form.get('wallet').value && this.form.get('wallet').value.coins ? this.form.get('wallet').value.coins : 0;
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

  private resetForm() {
    this.form.get('wallet').setValue('');
    this.form.get('address').setValue('');
    this.form.get('amount').setValue('');
    this.form.get('note').setValue('');
    this.selectedCurrency = DoubleButtonActive.LeftButton;
  }

  private closeSyncCheckSubscription() {
    if (this.syncCheckSubscription) {
      this.syncCheckSubscription.unsubscribe();
    }
  }
}
