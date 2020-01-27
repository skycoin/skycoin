import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { OfflineDialogsBaseComponent, OfflineDialogsStates } from '../offline-dialogs-base.component';
import { MsgBarService } from '../../../../../services/msg-bar.service';
import { FormBuilder } from '@angular/forms';
import { SubscriptionLike } from 'rxjs';
import { WalletService } from '../../../../../services/wallet.service';
import { AppConfig } from '../../../../../app.config';
import { BalanceAndOutputsService } from 'src/app/services/wallet-operations/balance-and-outputs.service';

@Component({
  selector: 'app-broadcast-raw-tx',
  templateUrl: '../offline-dialogs-base.component.html',
  styleUrls: ['../offline-dialogs-base.component.scss'],
})
export class BroadcastRawTxComponent extends OfflineDialogsBaseComponent implements OnInit, OnDestroy {
  title = 'offline-transactions.broadcast-tx.title';
  text = 'offline-transactions.broadcast-tx.text';
  inputLabel = 'offline-transactions.broadcast-tx.input-label';
  cancelButtonText = 'common.cancel-button';
  okButtonText = 'offline-transactions.broadcast-tx.send-button';
  validateForm = true;

  private operationSubscription: SubscriptionLike;

  public static openDialog(dialog: MatDialog): MatDialogRef<BroadcastRawTxComponent, any> {
    const config = new MatDialogConfig();
    config.autoFocus = true;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(BroadcastRawTxComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<BroadcastRawTxComponent>,
    private walletService: WalletService,
    private msgBarService: MsgBarService,
    private balanceAndOutputsService: BalanceAndOutputsService,
    formBuilder: FormBuilder,
  ) {
    super(formBuilder);

    this.currentState = OfflineDialogsStates.ShowingForm;
  }

  ngOnInit() {
    this.form.get('dropdown').setValue('dummy');
  }

  ngOnDestroy() {
    this.closeOperationSubscription();
  }

  cancelPressed() {
    this.dialogRef.close();
  }

  okPressed() {
    if (this.working) {
      return;
    }

    this.msgBarService.hide();
    this.working = true;
    this.okButton.setLoading();

    this.closeOperationSubscription();
    this.operationSubscription = this.walletService.injectTransaction(this.form.get('input').value, null).subscribe(response => {
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
