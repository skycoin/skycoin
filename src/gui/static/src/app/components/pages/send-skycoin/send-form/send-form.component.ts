import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import 'rxjs/add/operator/delay';
import 'rxjs/add/operator/filter';
import { ButtonComponent } from '../../../layout/button/button.component';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { MatDialog, MatSnackBar } from '@angular/material';
import { showSnackbarError } from '../../../../utils/errors';
import { ISubscription } from 'rxjs/Subscription';
import { NavBarService } from '../../../../services/nav-bar.service';

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

  private subscription: ISubscription;

  constructor(
    public formBuilder: FormBuilder,
    public walletService: WalletService,
    private dialog: MatDialog,
    private snackbar: MatSnackBar,
    private navbarService: NavBarService,
  ) {}

  ngOnInit() {
    this.navbarService.showSwitch('send.simple', 'send.advanced');
    this.initForm();
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
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

    if (this.form.value.wallet.encrypted) {
      this.dialog.open(PasswordDialogComponent).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.createTransaction(passwordDialog);
        });
    } else {
      this.createTransaction();
    }
  }

  private createTransaction(passwordDialog?: any) {
    if (passwordDialog) {
      passwordDialog.close();
    }

    if (this.previewTx) {
      this.previewButton.setLoading();
      this.sendButton.setDisabled();
    } else {
      this.sendButton.setLoading();
      this.previewButton.setDisabled();
    }

    this.walletService.createTransaction(
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
    )
      .toPromise()
      .then(transaction => {
        if (!this.previewTx) {
          return this.walletService.injectTransaction(transaction.encoded).toPromise();
        }

        this.onFormSubmitted.emit({
          form: {
            wallet: this.form.value.wallet,
            address: this.form.value.address,
            amount: this.form.value.amount,
          },
          amount: this.form.value.amount,
          to: [this.form.value.address],
          transaction,
        });
      })
      .then(() => {
        this.sendButton.setSuccess();
        this.resetForm();

        setTimeout(() => {
          this.sendButton.resetState();
        }, 3000);
      })
      .catch(error => {
        showSnackbarError(this.snackbar, error);

        this.previewButton.resetState().setEnabled();
        this.sendButton.resetState().setEnabled();
      });
  }

  private initForm() {
    this.form = this.formBuilder.group({
      wallet: ['', Validators.required],
      address: ['', Validators.required],
      amount: ['', Validators.required],
    });

    this.subscription = this.form.get('wallet').valueChanges.subscribe(value => {
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
