import { Component, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { MessageIcons } from './hw-message/hw-message.component';
import { HwWalletService } from '../../../services/hw-wallet.service';

@Component({
  template: '',
})
export class HwDialogBaseComponent<T> implements OnDestroy {
  closeIfHwDisconnected = true;

  msgIcons = MessageIcons;

  protected operationSubscription: ISubscription;
  private hwConnectionSubscription: ISubscription;

  constructor(
    private _hwWalletService: HwWalletService,
    public _dialogRef: MatDialogRef<T>,
  ) {
    this.hwConnectionSubscription = this._hwWalletService.walletConnectedAsyncEvent.subscribe(connected => {
      this.hwConnectionChanged(connected);
      if (!connected && this.closeIfHwDisconnected) {
        this.closeModal();
      }
    });
  }

  ngOnDestroy() {
    if (this.operationSubscription && !this.operationSubscription.closed) {
      this.operationSubscription.unsubscribe();
    }
    this.hwConnectionSubscription.unsubscribe();
  }

  closeModal() {
    this._dialogRef.close();
  }

  hwConnectionChanged(connected: boolean) {

  }
}
