import { Component, OnDestroy, ViewChild } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { MessageIcons } from './hw-message/hw-message.component';
import { HwWalletService, OperationResults } from '../../../services/hw-wallet.service';
import { ButtonComponent } from '../button/button.component';
import { getHardwareWalletErrorMsg } from '../../../utils/errors';

export class ResultProcessingResponse {
  text: String;
  icon: MessageIcons;
}

export enum States {
  Connecting,
  Initial,
  Processing,
  ShowingResult,
  Finished,
  Other,
}

@Component({
  template: '',
})
export class HwDialogBaseComponent<T> implements OnDestroy {
  @ViewChild('closeButton', { static: false }) closeButton: ButtonComponent;

  closeIfHwDisconnected = true;

  msgIcons = MessageIcons;
  currentState: States = States.Initial;
  states = States;
  result: ResultProcessingResponse;

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

  protected processResult(result: OperationResults, genericError: string = null) {
    if (result && result === OperationResults.Disconnected && this.closeIfHwDisconnected) {
      this.closeModal();
    } else if (result) {
      this.showResult({
        text: getHardwareWalletErrorMsg(null, {result: result}, genericError),
        icon: MessageIcons.Error,
      });
    }
  }

  protected showResult(result: ResultProcessingResponse, focusButton = true) {
    if (result) {
      this.currentState = States.ShowingResult;
      this.result = result;

      setTimeout(() => {
        if (this.closeButton && focusButton) {
          this.closeButton.focus();
        }
      });
    }
  }
}
