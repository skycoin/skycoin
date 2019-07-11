import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild, ChangeDetectorRef } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import 'rxjs/add/operator/delay';
import 'rxjs/add/operator/filter';
import { ButtonComponent } from '../../../layout/button/button.component';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { MatDialog, MatDialogConfig } from '@angular/material';
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

@Component({
  selector: 'app-send-form',
  templateUrl: './send-form.component.html',
  styleUrls: ['./send-form.component.scss'],
})
export class SendFormComponent implements OnInit, OnDestroy {

  public static readonly MaxUsdDecimal = 6;

  @ViewChild('previewButton') previewButton: ButtonComponent;
  @ViewChild('sendButton') sendButton: ButtonComponent;
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

  private subscriptionsGroup: ISubscription[] = [];
  private processingSubscription: ISubscription;
  private syncCheckSubscription: ISubscription;

  constructor(
    public formBuilder: FormBuilder,
    public blockchainService: BlockchainService,
    public walletService: WalletService,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    private navbarService: NavBarService,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
    private changeDetector: ChangeDetectorRef,
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
        address: this.form.value.address,
        coins: this.selectedCurrency === DoubleButtonActive.LeftButton ? this.form.value.amount : this.value.toString(),
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
              if (note && !noteSaved) {
                this.msgBarService.showError(this.translate.instant('send.error-saving-note'));
              }

              this.showSuccess();
            }, error => this.showError(error));
        } else {
          this.onFormSubmitted.emit({
            form: {
              wallet: this.form.value.wallet,
              address: this.form.value.address,
              amount: this.form.value.amount,
              currency: this.selectedCurrency,
              note: note,
            },
            amount: new BigNumber(this.form.value.amount),
            to: [this.form.value.address],
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

  private showSuccess() {
    this.busy = false;
    this.navbarService.enableSwitch();
    this.sendButton.setSuccess();
    this.resetForm();

    setTimeout(() => {
      this.sendButton.resetState();
    }, 3000);
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
      address: ['', Validators.required],
      amount: ['', Validators.required],
      note: [''],
    });

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

  private validateAmount(amountControl: FormControl) {
    if (isNaN(amountControl.value.replace(' ', '='))) {
      return { Invalid: true };
    }

    if (parseFloat(amountControl.value) <= 0) {
      return { Invalid: true };
    }

    const parts = amountControl.value.split('.');

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
