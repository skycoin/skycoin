import { Component, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { WalletService } from '../../../../services/wallet.service';
import { ISubscription } from 'rxjs/Subscription';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { MessageIcons } from '../hw-message/hw-message.component';

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
export class HwAddedDialogComponent implements OnDestroy {

  msgIcons = MessageIcons;
  currentState: States = States.Initial;
  states = States;
  errorMsg = 'hardware-wallet.general.generic-error-internet';
  walletName: string;

  private operationSubscription: ISubscription;
  private hwConnectionSubscription: ISubscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) public notifyFinish: any,
    public dialogRef: MatDialogRef<HwAddedDialogComponent>,
    private walletService: WalletService,
    private hwWalletService: HwWalletService,
  ) {
    this.operationSubscription = this.walletService.createHardwareWallet().subscribe(wallet => {
      this.walletName = wallet.label;
      this.currentState = States.Finished;
      this.notifyFinish();
    }, () => {
      this.currentState = States.Failed;
      this.notifyFinish(this.errorMsg);
    });

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
