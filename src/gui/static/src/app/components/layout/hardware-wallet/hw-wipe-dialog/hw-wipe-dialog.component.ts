import { Component, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { MessageIcons } from '../hw-message/hw-message.component';
import { WalletService } from '../../../../services/wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';

enum States {
  Initial,
  Processing,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
}

@Component({
  selector: 'app-hw-wipe-dialog',
  templateUrl: './hw-wipe-dialog.component.html',
  styleUrls: ['./hw-wipe-dialog.component.scss'],
})
export class HwWipeDialogComponent implements OnDestroy {

  currentState: States = States.Initial;
  states = States;
  msgIcons = MessageIcons;
  showDeleteFromList = true;
  deleteFromList = true;

  private operationSubscription: ISubscription;
  private hwConnectionSubscription: ISubscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwWipeDialogComponent>,
    private hwWalletService: HwWalletService,
    private walletService: WalletService,
  ) {
    if (!data.wallet) {
      this.showDeleteFromList = false;
      this.deleteFromList = false;
    }

    this.hwConnectionSubscription = this.hwWalletService.walletConnectedAsyncEvent.subscribe(connected => {
      if (!connected) {
        this.closeModal();
      }
    });
  }

  ngOnDestroy() {
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }
    this.hwConnectionSubscription.unsubscribe();
  }

  setDeleteFromList(event) {
    this.deleteFromList = event.checked;
  }

  closeModal() {
    this.dialogRef.close();
  }

  requestWipe() {
    this.currentState = States.Processing;

    this.operationSubscription = this.hwWalletService.wipe().subscribe(
      () => {
        this.data.requestOptionsComponentRefresh();
        this.currentState = States.ReturnedSuccess;
        if (this.deleteFromList) {
          this.walletService.deleteHardwareWallet(this.data.wallet).subscribe();
        }
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
