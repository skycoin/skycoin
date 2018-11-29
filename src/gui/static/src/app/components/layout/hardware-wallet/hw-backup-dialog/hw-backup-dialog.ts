import { Component, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { MessageIcons } from '../hw-message/hw-message.component';

enum States {
  Initial,
  Processing,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
}

@Component({
  selector: 'app-hw-backup-dialog',
  templateUrl: './hw-backup-dialog.html',
  styleUrls: ['./hw-backup-dialog.scss'],
})
export class HwBackupDialogComponent implements OnDestroy {

  currentState: States = States.Initial;
  states = States;
  msgIcons = MessageIcons;

  private operationSubscription: ISubscription;
  private hwConnectionSubscription: ISubscription;

  constructor(
    public dialogRef: MatDialogRef<HwBackupDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    this.hwConnectionSubscription = this.hwWalletService.walletConnectedAsyncEvent.subscribe(connected => {
      if (!connected) {
        this.dialogRef.close();
      }
    });
  }

  ngOnDestroy() {
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }
    this.hwConnectionSubscription.unsubscribe();
  }

  closeModal() {
    this.dialogRef.close();
  }

  requestBackup() {
    this.currentState = States.Processing;

    this.operationSubscription = this.hwWalletService.backup().subscribe(
      () => {
        this.currentState = States.ReturnedSuccess;
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
}
