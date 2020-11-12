import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { FormBuilder } from '@angular/forms';
import { SubscriptionLike } from 'rxjs';

import { OfflineDialogsBaseComponent, OfflineDialogsStates } from '../offline-dialogs-base.component';
import { MsgBarService } from '../../../../../services/msg-bar.service';
import { AppConfig } from '../../../../../app.config';
import { BalanceAndOutputsService } from '../../../../../services/wallet-operations/balance-and-outputs.service';
import { SpendingService } from '../../../../../services/wallet-operations/spending.service';

/**
 * Allows to send a signed raw transaction to the network, to spend the coins.
 */
@Component({
  selector: 'app-broadcast-raw-tx',
  templateUrl: '../offline-dialogs-base.component.html',
  styleUrls: ['../offline-dialogs-base.component.scss'],
})
export class BroadcastRawTxComponent extends OfflineDialogsBaseComponent implements OnInit, OnDestroy {
  // Configure the UI.
  title = 'offline-transactions.broadcast-tx.title';
  text = 'offline-transactions.broadcast-tx.text';
  inputLabel = 'offline-transactions.broadcast-tx.input-label';
  cancelButtonText = 'common.cancel-button';
  okButtonText = 'offline-transactions.broadcast-tx.send-button';
  validateForm = true;

  private operationSubscription: SubscriptionLike;

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog): MatDialogRef<BroadcastRawTxComponent, any> {
    const config = new MatDialogConfig();
    config.autoFocus = true;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(BroadcastRawTxComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<BroadcastRawTxComponent>,
    private msgBarService: MsgBarService,
    private balanceAndOutputsService: BalanceAndOutputsService,
    private spendingService: SpendingService,
    formBuilder: FormBuilder,
  ) {
    super(formBuilder);

    this.currentState = OfflineDialogsStates.ShowingForm;
  }

  ngOnInit() {
    // Needed for making the form validation work.
    this.form.get('dropdown').setValue('dummy');
  }

  ngOnDestroy() {
    this.closeOperationSubscription();
  }

  cancelPressed() {
    this.dialogRef.close();
  }

  // Sends the transaction to the network.
  okPressed() {
    if (this.working) {
      return;
    }

    this.msgBarService.hide();
    this.working = true;
    this.okButton.setLoading();

    this.closeOperationSubscription();
    this.operationSubscription = this.spendingService.injectTransaction(this.form.get('input').value, null).subscribe(() => {
      this.balanceAndOutputsService.refreshBalance();

      this.msgBarService.showDone('offline-transactions.broadcast-tx.sent');
      this.cancelPressed();
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
