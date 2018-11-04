import { Component, OnDestroy } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { HwWalletService } from '../../../../services/hw-wallet.service';

export enum States {
  Initial,
  Finalized,
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

  constructor(
    public dialogRef: MatDialogRef<HwWipeDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    this.operationSubscription = this.hwWalletService.wipe().subscribe(
      () => {
        this.currentState = States.Finalized;
      },
      () => {
        this.currentState = States.Failed;
      },
    );
  }

  ngOnDestroy() {
    this.operationSubscription.unsubscribe();
  }
}
