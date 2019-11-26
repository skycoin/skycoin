import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

@Component({
  selector: 'app-hw-backup-dialog',
  templateUrl: './hw-backup-dialog.component.html',
  styleUrls: ['./hw-backup-dialog.component.scss'],
})
export class HwBackupDialogComponent extends HwDialogBaseComponent<HwBackupDialogComponent> {
  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwBackupDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
  }

  requestBackup() {
    this.currentState = this.states.Processing;

    this.operationSubscription = this.hwWalletService.backup().subscribe(
      () => {
        this.showResult({
          text: 'hardware-wallet.general.completed',
          icon: this.msgIcons.Success,
        });
        this.data.requestOptionsComponentRefresh(null, true);
      },
      err => this.processResult(err.result),
    );
  }
}
