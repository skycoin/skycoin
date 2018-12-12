import { Component, OnDestroy, ViewChild, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { MessageIcons } from '../hw-message/hw-message.component';

enum States {
  Initial,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
  WrongWord,
}

@Component({
  selector: 'app-hw-restore-seed-dialog',
  templateUrl: './hw-restore-seed-dialog.component.html',
  styleUrls: ['./hw-restore-seed-dialog.component.scss'],
})
export class HwRestoreSeedDialogComponent implements OnDestroy {

  currentState: States = States.Initial;
  states = States;
  msgIcons = MessageIcons;

  private operationSubscription: ISubscription;
  private hwConnectionSubscription: ISubscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) public notifyFinish: any,
    public dialogRef: MatDialogRef<HwRestoreSeedDialogComponent>,
    private hwWalletService: HwWalletService,
  ) {
    this.operationSubscription = this.hwWalletService.recoverMnemonic().subscribe(
      () => {
        this.notifyFinish();
        this.currentState = States.ReturnedSuccess;
      },
      err => {
        if (err.result && err.result === OperationResults.FailedOrRefused) {
          this.currentState = States.ReturnedRefused;
        } else if (err.result && err.result === OperationResults.WrongWord) {
          this.currentState = States.WrongWord;
        } else {
          this.currentState = States.Failed;
        }
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

  closeModal() {
    this.dialogRef.close();
  }
}
