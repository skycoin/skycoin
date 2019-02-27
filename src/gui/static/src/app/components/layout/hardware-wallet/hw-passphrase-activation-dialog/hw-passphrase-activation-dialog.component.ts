import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { HwPassphraseHelpDialogComponent } from '../hw-passphrase-help-dialog/hw-passphrase-help-dialog.component';

enum States {
  Initial,
  Processing,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
}

@Component({
  selector: 'app-hw-passphrase-activation-dialog',
  templateUrl: './hw-passphrase-activation-dialog.component.html',
  styleUrls: ['./hw-passphrase-activation-dialog.component.scss'],
})
export class HwPassphraseActivationDialogComponent extends HwDialogBaseComponent<HwPassphraseActivationDialogComponent> {

  currentState: States = States.Initial;
  states = States;
  activating: boolean;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwPassphraseActivationDialogComponent>,
    private hwWalletService: HwWalletService,
    private dialog: MatDialog,
  ) {
    super(hwWalletService, dialogRef);
    this.activating = !data.walletHasPassphrase;
  }

  changeConfig() {
    this.currentState = States.Processing;

    this.operationSubscription = this.hwWalletService.applySettings(this.activating).subscribe(
      () => {
        this.currentState = States.ReturnedSuccess;
        this.data.requestOptionsComponentRefresh();
      },
      err => {
        if (err.result && err.result === OperationResults.FailedOrRefused) {
          this.currentState = States.ReturnedRefused;
        } else {
          this.currentState = States.Failed;
        }
      },
    );
  }

  openHelp() {
    this.dialog.open(HwPassphraseHelpDialogComponent, <MatDialogConfig> {
      width: '450px',
    });
  }
}
