import { Component, Inject, ViewChild, OnDestroy } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { CreateWalletFormComponent } from './create-wallet-form/create-wallet-form.component';
import { SubscriptionLike } from 'rxjs';
import { BlockchainService } from '../../../../services/blockchain.service';
import { ConfirmationData } from '../../../../app.datatypes';
import { showConfirmationModal } from '../../../../utils';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { AppConfig } from '../../../../app.config';

export class CreateWalletParams {
  create: boolean;
}

@Component({
  selector: 'app-create-wallet',
  templateUrl: './create-wallet.component.html',
  styleUrls: ['./create-wallet.component.scss'],
})
export class CreateWalletComponent implements OnDestroy {
  @ViewChild('formControl', { static: false }) formControl: CreateWalletFormComponent;
  @ViewChild('createButton', { static: false }) createButton: ButtonComponent;
  @ViewChild('cancelButton', { static: false }) cancelButton: ButtonComponent;

  scan: Number;
  disableDismiss = false;

  private synchronized = true;
  private synchronizedSubscription: SubscriptionLike;

  public static openDialog(dialog: MatDialog, params: CreateWalletParams): MatDialogRef<CreateWalletComponent, any> {
    const config = new MatDialogConfig();
    config.data = params;
    config.autoFocus = true;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(CreateWalletComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) public data,
    public dialogRef: MatDialogRef<CreateWalletComponent>,
    private walletService: WalletService,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    blockchainService: BlockchainService,
  ) {
    this.synchronizedSubscription = blockchainService.synchronized.subscribe(value => this.synchronized = value);
  }

  ngOnDestroy() {
    this.synchronizedSubscription.unsubscribe();
    this.msgBarService.hide();
  }

  closePopup() {
    this.dialogRef.close();
  }

  createWallet() {
    if (!this.formControl.isValid || this.createButton.isLoading()) {
      return;
    }

    this.msgBarService.hide();

    if (this.synchronized || this.data.create) {
      this.continueCreating();
    } else {
      const confirmationData: ConfirmationData = {
        headerText: 'wallet.new.synchronizing-warning-title',
        text: 'wallet.new.synchronizing-warning-text',
        confirmButtonText: 'wallet.new.synchronizing-warning-continue',
        cancelButtonText: 'wallet.new.synchronizing-warning-cancel',
      };

      showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          this.continueCreating();
        }
      });
    }
  }

  private continueCreating() {
    const data = this.formControl.getData();

    this.createButton.resetState();
    this.createButton.setLoading();
    this.cancelButton.setDisabled();
    this.disableDismiss = true;

    this.walletService.create(data.label, data.seed, this.scan, data.password)
      .subscribe(() => {
        setTimeout(() => this.msgBarService.showDone('wallet.new.wallet-created'));
        this.dialogRef.close();
      }, e => {
        this.msgBarService.showError(e);
        this.createButton.resetState();
        this.cancelButton.disabled = false;
        this.disableDismiss = false;
      });
  }
}
