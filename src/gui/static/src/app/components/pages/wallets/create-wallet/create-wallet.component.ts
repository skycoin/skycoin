import { Component, Inject, ViewChild, OnDestroy } from '@angular/core';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { CreateWalletFormComponent } from './create-wallet-form/create-wallet-form.component';
import { SubscriptionLike } from 'rxjs';
import { BlockchainService } from '../../../../services/blockchain.service';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { AppConfig } from '../../../../app.config';
import { ConfirmationParams, ConfirmationComponent, DefaultConfirmationButtons } from '../../../layout/confirmation/confirmation.component';
import { WalletsAndAddressesService } from 'src/app/services/wallet-operations/wallets-and-addresses.service';

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

  disableDismiss = false;
  busy = false;

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
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    private walletsAndAddressesService: WalletsAndAddressesService,
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
    if (!this.formControl.isValid || this.busy) {
      return;
    }

    this.msgBarService.hide();

    if (this.synchronized || this.data.create) {
      this.continueCreating();
    } else {
      const confirmationParams: ConfirmationParams = {
        headerText: 'common.warning-title',
        text: 'wallet.new.synchronizing-warning-text',
        defaultButtons: DefaultConfirmationButtons.ContinueCancel,
        redTitle: true,
      };

      ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          this.continueCreating();
        }
      });
    }
  }

  private continueCreating() {
    this.busy = true;
    const data = this.formControl.getData();

    this.createButton.resetState();
    this.createButton.setLoading();
    this.cancelButton.setDisabled();
    this.disableDismiss = true;

    this.walletsAndAddressesService.createSoftwareWallet(data.label, data.seed, data.password)
      .subscribe(() => {
        this.busy = false;
        setTimeout(() => this.msgBarService.showDone('wallet.new.wallet-created'));
        this.dialogRef.close();
      }, e => {
        this.busy = false;
        this.msgBarService.showError(e);
        this.createButton.resetState();
        this.cancelButton.disabled = false;
        this.disableDismiss = false;
      });
  }
}
