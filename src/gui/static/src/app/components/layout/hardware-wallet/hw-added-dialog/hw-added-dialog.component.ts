import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { WalletService } from '../../../../services/wallet.service';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { Wallet } from '../../../../app.datatypes';
import { HwPassphraseHelpDialogComponent } from '../hw-passphrase-help-dialog/hw-passphrase-help-dialog.component';

enum States {
  Initial,
  PassphraseWarning,
  Adding,
  Finished,
  Failed,
}

@Component({
  selector: 'app-hw-added-dialog',
  templateUrl: './hw-added-dialog.component.html',
  styleUrls: ['./hw-added-dialog.component.scss'],
})
export class HwAddedDialogComponent extends HwDialogBaseComponent<HwAddedDialogComponent> {

  closeIfHwDisconnected = true;

  currentState: States = States.Initial;
  states = States;
  errorMsg = 'hardware-wallet.general.generic-error-internet';
  wallet: Wallet;
  form: FormGroup;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwAddedDialogComponent>,
    private walletService: WalletService,
    hwWalletService: HwWalletService,
    private formBuilder: FormBuilder,
    private dialog: MatDialog,
  ) {
    super(hwWalletService, dialogRef);
    this.operationSubscription = hwWalletService.getFeatures().subscribe(features => {
      if (features.rawResponse.passphraseProtection) {
        this.currentState = States.PassphraseWarning;
      } else {
        this.startAdding();
      }
    });
  }

  cancel() {
    this.data.requestOptionsComponentRefresh('hardware-wallet.added.canceled');
    this.closeModal();
  }

  startAdding() {
    this.currentState = States.Adding;

    this.operationSubscription = this.walletService.createHardwareWallet().subscribe(wallet => {
      this.walletService.updateWalletHasHwSecurityWarnings(wallet).subscribe(() => {
        this.wallet = wallet;

        this.form = this.formBuilder.group({
          label: [wallet.label, Validators.required],
        });

        this.closeIfHwDisconnected = false;
        this.currentState = States.Finished;
        this.data.requestOptionsComponentRefresh();
      });
    }, () => {
      this.currentState = States.Failed;
      this.data.requestOptionsComponentRefresh(this.errorMsg);
    });
  }

  saveNameAndCloseModal() {
    this.wallet.label = this.form.value.label;
    this.walletService.saveHardwareWallets();
    this.closeModal();
  }

  openHelp() {
    this.dialog.open(HwPassphraseHelpDialogComponent, <MatDialogConfig> {
      width: '450px',
    });
  }
}
