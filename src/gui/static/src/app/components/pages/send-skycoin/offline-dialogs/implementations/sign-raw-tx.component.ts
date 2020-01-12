import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { OfflineDialogsBaseComponent, OfflineDialogsStates } from '../offline-dialogs-base.component';
import { MsgBarService } from '../../../../../services/msg-bar.service';
import { FormBuilder } from '@angular/forms';
import { SubscriptionLike } from 'rxjs';
import { WalletService } from '../../../../../services/wallet.service';
import { first } from 'rxjs/operators';
import { CopyRawTxData, CopyRawTxComponent } from './copy-raw-tx.component';
import { Wallet } from '../../../../../app.datatypes';
import { PasswordDialogComponent } from '../../../../../components/layout/password-dialog/password-dialog.component';
import { parseResponseMessage } from '../../../../../utils/errors';
import { AppConfig } from '../../../../../app.config';

@Component({
  selector: 'app-sign-raw-tx',
  templateUrl: '../offline-dialogs-base.component.html',
  styleUrls: ['../offline-dialogs-base.component.scss'],
})
export class SignRawTxComponent extends OfflineDialogsBaseComponent implements OnInit, OnDestroy {
  title = 'offline-transactions.sign-tx.title';
  text = 'offline-transactions.sign-tx.text';
  dropdownLabel = 'offline-transactions.sign-tx.wallet-label';
  defaultDropdownText = 'offline-transactions.sign-tx.select-wallet';
  inputLabel = 'offline-transactions.sign-tx.input-label';
  cancelButtonText = 'offline-transactions.sign-tx.cancel';
  okButtonText = 'offline-transactions.sign-tx.sign';
  validateForm = true;

  private walletsSubscription: SubscriptionLike;
  private operationSubscription: SubscriptionLike;

  public static openDialog(dialog: MatDialog): MatDialogRef<SignRawTxComponent, any> {
    const config = new MatDialogConfig();
    config.autoFocus = true;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(SignRawTxComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<SignRawTxComponent>,
    private walletService: WalletService,
    private msgBarService: MsgBarService,
    private dialog: MatDialog,
    formBuilder: FormBuilder,
  ) {
    super(formBuilder);

    this.currentState = OfflineDialogsStates.Loading;
  }

  ngOnInit() {
    this.walletsSubscription = this.walletService.all().pipe(first()).subscribe(wallets => {
      if (wallets) {
        this.dropdownElements = [];

        wallets.forEach(wallet => {
          if (!wallet.isHardware) {
            this.dropdownElements.push({
              name: wallet.label,
              value: wallet,
            });
          }
        });

        this.currentState = OfflineDialogsStates.ShowingForm;

        setTimeout(() => {
          try {
            if (wallets.length === 1) {
              this.form.get('dropdown').setValue(wallets[0]);
            }
          } catch (e) { }
        });
      } else {
        this.currentState = OfflineDialogsStates.ErrorLoading;
      }
    }, () => this.currentState = OfflineDialogsStates.ErrorLoading);
  }

  ngOnDestroy() {
    this.walletsSubscription.unsubscribe();
    this.closeOperationSubscription();
  }

  cancelPressed() {
    this.dialogRef.close();
  }

  okPressed() {
    if (this.working) {
      return;
    }

    if ((this.form.get('dropdown').value as Wallet).encrypted) {
      PasswordDialogComponent.openDialog(this.dialog, { wallet: this.form.get('dropdown').value }).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          passwordDialog.close();
          this.signTransaction(passwordDialog.password);
        });
    } else {
      this.signTransaction(null);
    }
  }

  signTransaction(password: string) {
    this.msgBarService.hide();
    this.working = true;
    this.okButton.setLoading();

    this.closeOperationSubscription();
    this.operationSubscription = this.walletService.signTransaction(
      this.form.get('dropdown').value,
      password,
      null,
      this.form.get('input').value).subscribe(response => {

        this.working = false;
        this.okButton.resetState();

        this.msgBarService.showDone('offline-transactions.sign-tx.signed');
        this.cancelPressed();

        setTimeout(() => {
          const data: CopyRawTxData = {
            rawTx: response.encoded,
            isUnsigned: false,
          };

          CopyRawTxComponent.openDialog(this.dialog, data);
        }, 500);
    }, error => {
      this.working = false;
      this.okButton.resetState();

      const parsedErrorMsg = parseResponseMessage(error);
      if (parsedErrorMsg !== error) {
        this.msgBarService.showError(parsedErrorMsg);
      } else if (error && error.error && error.error.error && error.error.error.message) {
        this.msgBarService.showError(error.error.error.message);
      } else {
        this.msgBarService.showError('offline-transactions.sign-tx.error');
      }
    });
  }

  closeOperationSubscription() {
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }
  }
}
