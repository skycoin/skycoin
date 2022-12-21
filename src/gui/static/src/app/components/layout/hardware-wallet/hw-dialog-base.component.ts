import { Component, OnDestroy, ViewChild } from '@angular/core';
import { MatLegacyDialogRef as MatDialogRef } from '@angular/material/legacy-dialog';
import { SubscriptionLike } from 'rxjs';

import { MessageIcons } from './hw-message/hw-message.component';
import { HwWalletService } from '../../../services/hw-wallet.service';
import { ButtonComponent } from '../button/button.component';
import { processServiceError } from '../../../utils/errors';
import { AppConfig } from '../../../app.config';
import { OperationError, HWOperationResults } from '../../../utils/operation-error';

/**
 * Data to show on the UI to inform the user about the result of an operation.
 */
export class ResultData {
  /**
   * Text to show.
   */
  text: string;
  /**
   * Link to show after the text. Must be a valid URL.
   */
  link?: string;
  /**
   * Text to show for the link, in case the raw value set in the "link" var is not what the
   * UI must show. If not set, the value of the "link" var is used.
   */
  linkText?: string;
  /**
   * If true, the route set in the "link" var is used as an internal link, so it is not openned
   * in a new tab.
   */
  linkIsInternal?: boolean;
  /**
   * Icon to show.
   */
  icon: MessageIcons;
}

/**
 * Different states in which an implementation of HwDialogBaseComponent can be. Each
 * implementation is responsible for using each state as appropiate, but the
 * ShowingResult state should be used for displaying in the UI the result of the
 * requested oepration.
 */
export enum States {
  Connecting,
  Initial,
  Processing,
  ShowingResult,
  Finished,
  Other,
}

/**
 * Base class for the modal windows related to the hw wallet. It includes code for simplifying
 * several actions which are common for the hw wallet modal windows, like managing states,
 * operations with the close button, closing the window if the device is disconnected and more.
 * The type expected by this class is just the class implementing it.
 */
@Component({
  template: '',
})
export class HwDialogBaseComponent<T> implements OnDestroy {
  // Reference to the close button. For it to work the implementation must have "#closeButton"
  // added to the close button tag.
  @ViewChild('closeButton') closeButton: ButtonComponent;

  // If true, the modal window will be automatically closed if the device is disconnected.
  closeIfHwDisconnected = true;

  // Allows to access from the HTML files the icons the msg areas can show.
  msgIcons = MessageIcons;
  // Current state of the window. Must be interpreted by the implementation.
  currentState: States = States.Initial;
  // Allows to access from the HTML files the states in which the window can be.
  states = States;
  // Result to show on the UI when the state indicates a result must be shown.
  result: ResultData;

  // Add operation subscriptions to this var to close them automatically when closing the window.
  protected operationSubscription: SubscriptionLike;
  private hwConnectionSubscription: SubscriptionLike;

  constructor(
    private _hwWalletService: HwWalletService,
    public _dialogRef: MatDialogRef<T>,
  ) {
    // Inform connection events and close the window if needed.
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

  /**
   * Called every time the connection state of the device changes.
   * @param connected If the device is connected (true) or not (false).
   */
  hwConnectionChanged(connected: boolean) {

  }

  /**
   * Process the result of an operation which finished in an error. It closes the modal window if
   * appropiate or prepares the error for being displayed on the UI.
   * @param result Result obtained after finishing the operation.
   */
  protected processHwOperationError(result: OperationError) {
    if (result) {
      result = processServiceError(result);

      if (result.type === HWOperationResults.Disconnected && this.closeIfHwDisconnected) {
        this.closeModal();
      } else {
        // Prepares the error for being displayed on the UI.
        this.showResult({
          text: result.translatableErrorMsg,
          icon: MessageIcons.Error,
        });
      }
    }
  }

  /**
   * Process and saves a result to be shown on the UI and also changes the state of the modal
   * window to the one indicationg a result must be shown to the user.
   * @param result Result to process.
   * @param focusButton If true, the close button of the modal window is focused.
   */
  protected showResult(result: ResultData, focusButton = true) {
    if (result) {
      // If there was an error connecting with the daemon, the link to download the daemon is
      // added to the elements which will be displayed.
      if (result.text === 'hardware-wallet.errors.daemon-connection' || result.text.indexOf('Problem connecting to the Skywallet Daemon') !== -1) {
        result.text = 'hardware-wallet.errors.daemon-connection-with-configurable-link';
        result.link = AppConfig.hwWalletDaemonDownloadUrl;
      }

      // If the operation was cancelled for inactivity, a link to open the page for wiping the
      // Skywallet is added.
      if (result.text === 'hardware-wallet.errors.timeout') {
        result.text = 'hardware-wallet.errors.timeout-with-configurable-link';
        result.linkText = 'force-skywallet-wipe.title';
        result.link = '#/skywallet-wipe';
        result.linkIsInternal = true;
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
