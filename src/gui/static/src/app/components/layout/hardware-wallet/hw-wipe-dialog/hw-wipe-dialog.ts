import { Component, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { HwWalletService } from '../../../../services/hw-wallet.service';

export enum States {
  Initial,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
}

@Component({
  selector: 'app-hw-wipe-dialog',
  templateUrl: './hw-wipe-dialog.html',
  styleUrls: ['./hw-wipe-dialog.scss'],
})
export class HwWipeDialogComponent implements OnDestroy {

  currentState: States = States.Initial;
  states = States;

  private operationSubscription: ISubscription;
  private hwConnectionSubscription: ISubscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) public requestRecheck: any,
    public dialogRef: MatDialogRef<HwWipeDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    this.operationSubscription = this.hwWalletService.wipe().subscribe(
      response => {
        if (response.success) {
          this.requestRecheck();
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
