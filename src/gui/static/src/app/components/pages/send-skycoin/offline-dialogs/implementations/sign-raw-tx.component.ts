import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { FormBuilder } from '@angular/forms';
import { SubscriptionLike } from 'rxjs';
import { first } from 'rxjs/operators';

import { OfflineDialogsBaseComponent, OfflineDialogsStates } from '../offline-dialogs-base.component';
import { MsgBarService } from '../../../../../services/msg-bar.service';
import { CopyRawTxData, CopyRawTxComponent } from './copy-raw-tx.component';
import { PasswordDialogComponent } from '../../../../../components/layout/password-dialog/password-dialog.component';
import { AppConfig } from '../../../../../app.config';
import { SpendingService } from '../../../../../services/wallet-operations/spending.service';
import { WalletsAndAddressesService } from '../../../../../services/wallet-operations/wallets-and-addresses.service';
import { WalletBase } from '../../../../../services/wallet-operations/wallet-objects';
import { OfflineTxPreviewComponent } from '../offline-tx-preview/offline-tx-preview.component';

/**
 * Allows to sign an unsigned raw tx. For it to be able to sign the transaction, all the
 * inputs must belong to one (only one) of the registered software wallets. After signing
 * the transaction, it opens a new modal window for showing the signed raw tx.
 */
@Component({
  selector: 'app-sign-raw-tx',
  templateUrl: '../offline-dialogs-base.component.html',
  styleUrls: ['../offline-dialogs-base.component.scss'],
})
export class SignRawTxComponent extends OfflineDialogsBaseComponent implements OnInit, OnDestroy {
  // Configure the UI.
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

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
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
    // Get the wallet list.
    this.walletsSubscription = this.walletsAndAddressesService.allWallets.pipe(first()).subscribe(wallets => {
      if (wallets) {
        this.dropdownElements = [];

        // Create a list with the software wallets, for the dropdown control.
        wallets.forEach(wallet => {
          if (!wallet.isHardware) {
            this.dropdownElements.push({
              name: wallet.label,
              value: wallet,
            });
          }
        });

        // Fill the dropdown control.
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
    this.msgBarService.hide();
  }

  cancelPressed() {
    this.dialogRef.close();
  }

  okPressed() {
    if (this.working) {
      return;
    }

    // TODO: reactivate after the problems with the 'transaction/verify' API endpoint are solved.
    // this.showPreview();

    // TODO: remove after the problems with the 'transaction/verify' API endpoint are solved.
    // Get the wallet password, if needed, and start signing the transaction.
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

  // Decodes the transaction, shows the preview and continues if the user confirms.
  // TODO: reactivate after the problems with the 'transaction/verify' API endpoint are solved.
  private showPreview() {
    this.msgBarService.hide();
    this.working = true;
    this.okButton.setLoading();

    this.closeOperationSubscription();
    this.operationSubscription = this.spendingService.decodeTransaction(
      this.form.get('input').value, true,
    ).subscribe(result => {
      this.working = false;
      this.okButton.resetState();

      // Open the preview modal window.
      OfflineTxPreviewComponent.openDialog(this.dialog, result).afterClosed().subscribe(r => {
        if (r) {
          // Get the wallet password, if needed, and start signing the transaction.
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
      });
    }, error => {
      this.working = false;
      this.okButton.resetState();

      this.msgBarService.showError(error);
    });
  }

  // Signs the transaction with the selected wallet.
  private signTransaction(password: string) {
    this.msgBarService.hide();
    this.working = true;
    this.okButton.setLoading();

    this.closeOperationSubscription();
    this.operationSubscription = this.spendingService.signTransaction(
      this.form.get('dropdown').value,
      password,
      null,
      this.form.get('input').value).subscribe(encodedSignedTx => {
        this.cancelPressed();
        setTimeout(() => this.msgBarService.showDone('offline-transactions.sign-tx.signed'));

        // After a short delay, open the copy modal window, to see the signed transaction.
        setTimeout(() => {
          const data: CopyRawTxData = {
            rawTx: encodedSignedTx,
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
