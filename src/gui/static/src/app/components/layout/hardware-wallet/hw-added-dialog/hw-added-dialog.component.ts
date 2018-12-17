import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { WalletService } from '../../../../services/wallet.service';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

enum States {
  Initial,
  Finished,
  Failed,
}

@Component({
  selector: 'app-hw-added-dialog',
  templateUrl: './hw-added-dialog.component.html',
  styleUrls: ['./hw-added-dialog.component.scss'],
})
export class HwAddedDialogComponent extends HwDialogBaseComponent<HwAddedDialogComponent> {

  closeIfHwDisconnected = false;

  currentState: States = States.Initial;
  states = States;
  errorMsg = 'hardware-wallet.general.generic-error-internet';
  walletName: string;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwAddedDialogComponent>,
    private walletService: WalletService,
    hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
    this.operationSubscription = this.walletService.createHardwareWallet().subscribe(wallet => {
      this.walletService.updateWalletHasHwSecurityWarnings(wallet).subscribe(() => {
        this.walletName = wallet.label;
        this.currentState = States.Finished;
        this.data.requestOptionsComponentRefresh();
      });
    }, () => {
      this.currentState = States.Failed;
      this.data.requestOptionsComponentRefresh(this.errorMsg);
    });
  }
}
