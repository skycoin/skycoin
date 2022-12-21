import { Component, Inject, ViewChild, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { MatLegacyDialogRef as MatDialogRef, MatLegacyDialog as MatDialog, MatLegacyDialogConfig as MatDialogConfig } from '@angular/material/legacy-dialog';
import { MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA } from '@angular/material/legacy-dialog';
import { SubscriptionLike } from 'rxjs';

import { ButtonComponent } from '../../../layout/button/button.component';
import { CreateWalletFormComponent } from './create-wallet-form/create-wallet-form.component';
import { BlockchainService } from '../../../../services/blockchain.service';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { AppConfig } from '../../../../app.config';
import { ConfirmationParams, ConfirmationComponent, DefaultConfirmationButtons } from '../../../layout/confirmation/confirmation.component';
import { WalletsAndAddressesService } from '../../../../services/wallet-operations/wallets-and-addresses.service';

/**
 * Settings for CreateWalletComponent.
 */
export class CreateWalletParams {
  /**
   * If the modal window is for creating a new wallet (true) or for loading a wallet
   * using a seed (false).
   */
  create: boolean;
}

/**
 * Modal window for creating a new software wallet or loading a software wallet using a seed.
 */
@Component({
  selector: 'app-create-wallet',
  templateUrl: './create-wallet.component.html',
  styleUrls: ['./create-wallet.component.scss'],
})
export class CreateWalletComponent implements OnDestroy {
  @ViewChild('formControl') formControl: CreateWalletFormComponent;
  @ViewChild('createButton') createButton: ButtonComponent;
  @ViewChild('cancelButton') cancelButton: ButtonComponent;

  // If the normal ways for closing the modal window must be deactivated.
  disableDismiss = false;
  // Deactivates the form while the system is busy.
  busy = false;

  // If the blockchain is synchronized.
  private synchronized = true;
  private blockchainSubscription: SubscriptionLike;

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
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
    private changeDetector: ChangeDetectorRef,
    blockchainService: BlockchainService,
  ) {
    this.blockchainSubscription = blockchainService.progress.subscribe(response => this.synchronized = response.synchronized);
  }

  ngOnDestroy() {
    this.blockchainSubscription.unsubscribe();
    this.msgBarService.hide();
  }

  closePopup() {
    this.dialogRef.close();
  }

  // Checks if the blockchain is synchronized before creating the wallet. If it is synchronized,
  // it continues creating the wallet, if not, the user must confirm the operation first.
  checkAndCreateWallet() {
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

    this.changeDetector.detectChanges();
  }

  // Creates the wallet with the data entered on the form.
  private continueCreating() {
    this.busy = true;
    const data = this.formControl.getData();

    this.createButton.resetState();
    this.createButton.setLoading();
    this.cancelButton.setDisabled();
    this.disableDismiss = true;

    this.walletsAndAddressesService.createSoftwareWallet(data.loadTemporarily, data.label, data.seed, data.password)
      .subscribe(() => {
        this.busy = false;
        setTimeout(() => this.msgBarService.showDone('wallet.new.wallet-created'));
        this.dialogRef.close();
      }, e => {
        this.busy = false;
        this.msgBarService.showError(e);
        this.createButton.resetState();
        this.cancelButton.setEnabled();
        this.disableDismiss = false;
      });
  }
}
