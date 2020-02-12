import { Component, OnDestroy, ViewChild } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { SubscriptionLike } from 'rxjs';
import { MessageIcons } from './hw-message/hw-message.component';
import { HwWalletService } from '../../../services/hw-wallet.service';
import { ButtonComponent } from '../button/button.component';
import { processServiceError } from '../../../utils/errors';
import { AppConfig } from '../../../app.config';
import { OperationError, HWOperationResults } from '../../../utils/operation-error';

export class ResultProcessingResponse {
  text: String;
  link?: String;
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

  protected operationSubscription: SubscriptionLike;
  private hwConnectionSubscription: SubscriptionLike;

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

  protected processResult(result: OperationError) {
    if (result) {
      result = processServiceError(result);

      if (result.type === HWOperationResults.Disconnected && this.closeIfHwDisconnected) {
        this.closeModal();
      } else {
        this.showResult({
          text: result.translatableErrorMsg,
          icon: MessageIcons.Error,
        });
      }
    }
  }

  protected showResult(result: ResultProcessingResponse, focusButton = true) {
    if (result) {
      if (result.text === 'hardware-wallet.errors.daemon-connection' || result.text.indexOf('Problem connecting to the Skywallet Daemon') !== -1) {
        result.text = 'hardware-wallet.errors.daemon-connection-with-configurable-link';
        result.link = AppConfig.hwWalletDaemonDownloadUrl;
      }

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
