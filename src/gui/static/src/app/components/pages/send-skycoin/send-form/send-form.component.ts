import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import 'rxjs/add/operator/delay';
import 'rxjs/add/operator/filter';
import { ButtonComponent } from '../../../layout/button/button.component';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { MatDialog, MatSnackBar, MatDialogConfig } from '@angular/material';
import { showSnackbarError, getHardwareWalletErrorMsg } from '../../../../utils/errors';
import { ISubscription } from 'rxjs/Subscription';
import { NavBarService } from '../../../../services/nav-bar.service';
import { BigNumber } from 'bignumber.js';
import { Observable } from 'rxjs/Observable';
import { PreviewTransaction, Wallet } from '../../../../app.datatypes';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-send-form',
  templateUrl: './send-form.component.html',
  styleUrls: ['./send-form.component.scss'],
})
export class SendFormComponent implements OnInit, OnDestroy {
  @ViewChild('previewButton') previewButton: ButtonComponent;
  @ViewChild('sendButton') sendButton: ButtonComponent;
  @Input() formData: any;
  @Output() onFormSubmitted = new EventEmitter<any>();

  form: FormGroup;
  transactions = [];
  previewTx: boolean;
  busy = false;

  private formSubscription: ISubscription;
  private processingSubscription: ISubscription;

  constructor(
    public formBuilder: FormBuilder,
    public walletService: WalletService,
    private dialog: MatDialog,
    private snackbar: MatSnackBar,
    private navbarService: NavBarService,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
  ) {}

  ngOnInit() {
    this.navbarService.showSwitch('send.simple', 'send.advanced');
    this.initForm();
  }

  ngOnDestroy() {
    if (this.processingSubscription && !this.processingSubscription.closed) {
      this.processingSubscription.unsubscribe();
    }
    this.formSubscription.unsubscribe();
    this.navbarService.hideSwitch();
    this.snackbar.dismiss();
  }

  preview() {
    this.previewTx = true;
    this.unlockAndSend();
  }

  send() {
    this.previewTx = false;
    this.unlockAndSend();
  }

  private unlockAndSend() {
    if (!this.form.valid || this.previewButton.isLoading() || this.sendButton.isLoading()) {
      return;
    }

    this.snackbar.dismiss();
    this.previewButton.resetState();
    this.sendButton.resetState();

    if (this.form.value.wallet.encrypted && !this.form.value.wallet.isHardware) {
      const config = new MatDialogConfig();
      config.data = {
        wallet: this.form.value.wallet,
      };

      this.dialog.open(PasswordDialogComponent, config).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.createTransaction(passwordDialog);
        });
    } else {
      if (!this.form.value.wallet.isHardware) {
        this.createTransaction();
      } else {
        this.showBusy();
        this.processingSubscription = this.hwWalletService.checkIfCorrectHwConnected((this.form.value.wallet as Wallet).addresses[0].address).subscribe(
          () => this.createTransaction(),
          err => this.showError(getHardwareWalletErrorMsg(this.hwWalletService, this.translate, err)),
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
  }

  private createTransaction(passwordDialog?: any) {
    if (passwordDialog) {
      passwordDialog.close();
    }

    this.showBusy();

    let createTxRequest: Observable<PreviewTransaction>;

    if (!this.form.value.wallet.isHardware) {
      createTxRequest = this.walletService.createTransaction(
        this.form.value.wallet,
        null,
        [{
          address: this.form.value.address,
          coins: this.form.value.amount,
        }],
        {
          type: 'auto',
          mode: 'share',
          share_factor: '0.5',
        },
        null,
        passwordDialog ? passwordDialog.password : null,
      );
    } else {
      createTxRequest = this.walletService.createHwTransaction(
        this.form.value.wallet,
        this.form.value.address,
        new BigNumber(this.form.value.amount),
      );
    }

     this.processingSubscription = createTxRequest.subscribe(transaction => {
        if (!this.previewTx) {
          this.processingSubscription = this.walletService.injectTransaction(transaction.encoded)
            .subscribe(() => this.showSuccess(), error => this.showError(error));
        } else {
          this.onFormSubmitted.emit({
            form: {
              wallet: this.form.value.wallet,
              address: this.form.value.address,
              amount: this.form.value.amount,
            },
            amount: new BigNumber(this.form.value.amount),
            to: [this.form.value.address],
            transaction,
          });
          this.busy = false;
        }
      },
      error => {
        if (error && error.result) {
          this.showError(getHardwareWalletErrorMsg(this.hwWalletService, this.translate, error));
        } else {
          this.showError(error);
        }
      },
    );
  }

  private showSuccess() {
    this.busy = false;
    this.sendButton.setSuccess();
    this.resetForm();

    setTimeout(() => {
      this.sendButton.resetState();
    }, 3000);
  }

  private showError(error) {
    this.busy = false;
    showSnackbarError(this.snackbar, error);
    this.previewButton.resetState().setEnabled();
    this.sendButton.resetState().setEnabled();
  }

  private initForm() {
    this.form = this.formBuilder.group({
      wallet: ['', Validators.required],
      address: ['', Validators.required],
      amount: ['', Validators.required],
    });

    this.formSubscription = this.form.get('wallet').valueChanges.subscribe(value => {
      const balance = value && value.coins ? value.coins : 0;

      this.form.get('amount').setValidators([
        Validators.required,
        Validators.max(balance),
        this.validateAmount,
      ]);

      this.form.get('amount').updateValueAndValidity();
    });

    if (this.formData) {
      Object.keys(this.form.controls).forEach(control => {
        this.form.get(control).setValue(this.formData.form[control]);
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

    if (parts.length === 2 && parts[1].length > 6) {
      return { Invalid: true };
    }

    return null;
  }

  private resetForm() {
    this.form.get('wallet').setValue('');
    this.form.get('address').setValue('');
    this.form.get('amount').setValue('');
  }
}
