import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { OfflineDialogsBaseComponent, OfflineDialogsStates } from '../offline-dialogs-base.component';
import { MsgBarService } from '../../../../../services/msg-bar.service';
import { FormBuilder } from '@angular/forms';
import { SubscriptionLike } from 'rxjs';
import { first } from 'rxjs/operators';
import { CopyRawTxData, CopyRawTxComponent } from './copy-raw-tx.component';
import { PasswordDialogComponent } from '../../../../../components/layout/password-dialog/password-dialog.component';
import { AppConfig } from '../../../../../app.config';
import { SpendingService } from '../../../../../services/wallet-operations/spending.service';
import { WalletsAndAddressesService } from '../../../../../services/wallet-operations/wallets-and-addresses.service';
import { WalletBase } from '../../../../../services/wallet-operations/wallet-objects';

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
  cancelButtonText = 'common.cancel-button';
  okButtonText = 'offline-transactions.sign-tx.sign-button';
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
    private msgBarService: MsgBarService,
    private dialog: MatDialog,
    private spendingService: SpendingService,
    private walletsAndAddressesService: WalletsAndAddressesService,
    formBuilder: FormBuilder,
  ) {
    super(formBuilder);

    this.currentState = OfflineDialogsStates.Loading;
  }

  ngOnInit() {
    this.walletsSubscription = this.walletsAndAddressesService.allWallets.pipe(first()).subscribe(wallets => {
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

    if ((this.form.get('dropdown').value as WalletBase).encrypted) {
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
    this.operationSubscription = this.spendingService.signTransaction(
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

      this.msgBarService.showError(error);
    });
  }

  closeOperationSubscription() {
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }
  }
}
