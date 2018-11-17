import { Component, OnDestroy, ViewChild, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { HwWalletService } from '../../../../services/hw-wallet.service';

enum States {
  Initial,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
}

@Component({
  selector: 'app-hw-generate-seed-dialog',
  templateUrl: './hw-generate-seed-dialog.html',
  styleUrls: ['./hw-generate-seed-dialog.scss'],
})
export class HwGenerateSeedDialogComponent implements OnDestroy {

  currentState: States = States.Initial;
  states = States;

  private operationSubscription: ISubscription;
  private hwConnectionSubscription: ISubscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) public requestRecheck: any,
    public dialogRef: MatDialogRef<HwGenerateSeedDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    this.operationSubscription = this.hwWalletService.generateMnemonic().subscribe(
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
