import { Component, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { HwWalletService } from '../../../../services/hw-wallet.service';

enum States {
  Initial,
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

  private operationSubscription: ISubscription;
  private hwConnectionSubscription: ISubscription;

  constructor(
    public dialogRef: MatDialogRef<HwBackupDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    this.operationSubscription = this.hwWalletService.backup().subscribe(
      response => {
        if (response.success) {
          this.currentState = States.ReturnedSuccess;
        } else {
          this.currentState = States.ReturnedRefused;
        }
      },
      () => {
        this.currentState = States.Failed;
      },
    );

    this.hwConnectionSubscription = this.hwWalletService.walletConnectedAsyncEvent.subscribe(connected => {
      if (!connected) {
        this.dialogRef.close();
      }
    });
  }

  ngOnDestroy() {
    this.operationSubscription.unsubscribe();
    this.hwConnectionSubscription.unsubscribe();
  }
}
