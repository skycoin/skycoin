import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { WalletService } from '../../../../services/wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

@Component({
  selector: 'app-hw-wipe-dialog',
  templateUrl: './hw-wipe-dialog.component.html',
  styleUrls: ['./hw-wipe-dialog.component.scss'],
})
export class HwWipeDialogComponent extends HwDialogBaseComponent<HwWipeDialogComponent> {
  showDeleteFromList = true;
  deleteFromList = true;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwWipeDialogComponent>,
    private hwWalletService: HwWalletService,
    private walletService: WalletService,
  ) {
    super(hwWalletService, dialogRef);

    if (!data.wallet) {
      this.showDeleteFromList = false;
      this.deleteFromList = false;
    }
  }

  setDeleteFromList(event) {
    this.deleteFromList = event.checked;
  }

  requestWipe() {
    this.currentState = this.states.Processing;

    this.operationSubscription = this.hwWalletService.wipe().subscribe(
      () => {
        this.showResult({
          text: 'hardware-wallet.general.completed',
          icon: this.msgIcons.Success,
        });
        this.data.requestOptionsComponentRefresh();
        if (this.deleteFromList) {
          this.walletService.deleteHardwareWallet(this.data.wallet).subscribe();
        }
      },
      err => this.processResult(err),
    );
  }
}
