import { Component, EventEmitter, Input, OnDestroy, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { SubscriptionLike } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { mergeMap } from 'rxjs/operators';

import { ButtonComponent } from '../../../layout/button/button.component';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { CopyRawTxData, CopyRawTxComponent } from '../offline-dialogs/implementations/copy-raw-tx.component';
import { ConfirmationParams, DefaultConfirmationButtons, ConfirmationComponent } from '../../../layout/confirmation/confirmation.component';
import { BalanceAndOutputsService } from '../../../../services/wallet-operations/balance-and-outputs.service';
import { SpendingService } from '../../../../services/wallet-operations/spending.service';
import { GeneratedTransaction } from '../../../../services/wallet-operations/transaction-objects';

/**
 * Shows the preview of a transaction before sending it to the network.
 */
@Component({
  selector: 'app-send-preview',
  templateUrl: './send-preview.component.html',
  styleUrls: ['./send-preview.component.scss'],
})
export class SendVerifyComponent implements OnDestroy {
  @ViewChild('sendButton', { static: false }) sendButton: ButtonComponent;
  @ViewChild('backButton', { static: false }) backButton: ButtonComponent;
  // Transaction which is going to be shown.
  @Input() transaction: GeneratedTransaction;
  // Emits when the preview must be removed from the UI and the form must be shown again. The
  // boolean value indicates if the form must be cleaned before showing it (true) or if it must
  // show the previously entered data again (false).
  @Output() onBack = new EventEmitter<boolean>();

  private sendSubscription: SubscriptionLike;

  constructor(
    private msgBarService: MsgBarService,
    private dialog: MatDialog,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
    private balanceAndOutputsService: BalanceAndOutputsService,
    private spendingService: SpendingService,
  ) {}

  ngOnDestroy() {
    this.msgBarService.hide();

    if (this.sendSubscription) {
      this.sendSubscription.unsubscribe();
    }
  }

  // Returns to the form.
  back() {
    this.onBack.emit(false);
  }

  // Sends the transaction.
  send() {
    if (this.sendButton.isLoading()) {
      return;
    }

    this.msgBarService.hide();
    this.sendButton.resetState();

    // If there is no wallet, the transaction is a manually created unsigned transaction, so
    // the raw transaction is shown in a modal window, so the user can sign it later,
    // instead of being sent.
    if (!this.transaction.wallet) {
      const data: CopyRawTxData = {
        rawTx: this.transaction.encoded,
        isUnsigned: true,
      };

      CopyRawTxComponent.openDialog(this.dialog, data).afterClosed().subscribe(() => {
        const confirmationParams: ConfirmationParams = {
          text: 'offline-transactions.copy-tx.reset-confirmation',
          defaultButtons: DefaultConfirmationButtons.YesNo,
        };

        // Ask the user if the form should be cleaned and shown again, to be able to create
        // a new transaction.
        ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
          if (confirmationResult) {
            this.onBack.emit(true);
          }
        });
      });

      return;
    }

    if (this.transaction.wallet.encrypted && !this.transaction.wallet.isHardware) {
      // If the wallet is encrypted, ask for the password and continue.
      PasswordDialogComponent.openDialog(this.dialog, { wallet: this.transaction.wallet }).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.finishSending(passwordDialog);
        });
    } else {
      if (!this.transaction.wallet.isHardware) {
        // If the wallet is not encrypted, continue.
        this.finishSending();
      } else {
        // If using a hw wallet, check the device first.
        this.showBusy();
        this.sendSubscription = this.hwWalletService.checkIfCorrectHwConnected(this.transaction.wallet.addresses[0].address).subscribe(
          () => this.finishSending(),
          err => this.showError(err),
        );
      }
    }
  }

  // Finish sending the transaction.
  private finishSending(passwordDialog?: any) {
    this.showBusy();

    const note = this.transaction.note.trim();

    // Sign the transaction.
    this.sendSubscription = this.spendingService.signTransaction(
      this.transaction.wallet,
      passwordDialog ? passwordDialog.password : null,
      this.transaction,
    ).pipe(mergeMap(encodedSignedTx => {
      // Close the password dialog, if it exists.
      if (passwordDialog) {
        passwordDialog.close();
      }

      // Send the transaction.
      return this.spendingService.injectTransaction(encodedSignedTx, note);
    })).subscribe(noteSaved => {
      // Show the final result.
      if (note && !noteSaved) {
        setTimeout(() => this.msgBarService.showWarning(this.translate.instant('send.saving-note-error')));
      } else {
        setTimeout(() => this.msgBarService.showDone('send.sent'));
      }

      this.balanceAndOutputsService.refreshBalance();

      this.onBack.emit(true);
    }, error => {
      if (passwordDialog) {
        passwordDialog.error(error);
      }

      this.showError(error);
    });
  }

  // Makes the UI to be shown busy.
  private showBusy() {
    this.sendButton.setLoading();
    this.backButton.setDisabled();
  }

  // Stops showing the UI busy and shows the error msg.
  private showError(error) {
    this.msgBarService.showError(error);
    this.sendButton.resetState();
    this.backButton.resetState().setEnabled();
  }
}
