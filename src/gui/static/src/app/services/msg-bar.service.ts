import { Injectable } from '@angular/core';
import { SubscriptionLike, of } from 'rxjs';
import { delay } from 'rxjs/operators';

import { MsgBarConfig, MsgBarComponent, MsgBarIcons, MsgBarColors } from '../components/layout/msg-bar/msg-bar.component';
import { processServiceError } from '../utils/errors';
import { AppConfig } from '../app.config';
import { environment } from '../../environments/environment';
import { OperationError, HWOperationResults } from '../utils/operation-error';

/**
 * Allows to control the msg bar, which is a small horizontal bar used by the app
 * for showing notifications.
 */
@Injectable()
export class MsgBarService {

  private timeSubscription: SubscriptionLike;

  /**
   * Sets the component which will be used as the app msg bar.
   */
  set msgBarComponent(value: MsgBarComponent) {
    this.msgBarComponentInternal = value;
  }
  private msgBarComponentInternal: MsgBarComponent;

  /**
   * Hides the msg bar.
   */
  hide() {
    if (this.msgBarComponentInternal) {
      this.msgBarComponentInternal.hide();
    }
  }

  /**
   * Displays the msg bar for showing an error.
   * @param body Text to show or OperationError instance from which the error will be obtained.
   * @param duration How much time the msg bar will stay visible, in ms. The default
   * value is 20000.
   */
  showError(body: string | OperationError, duration = 20000) {
    const config = new MsgBarConfig();
    // Process the body param to make sure the correct error msg is obtained.
    config.text = processServiceError(body).translatableErrorMsg;
    config.title = 'common.error-title';
    config.icon = MsgBarIcons.Error;
    config.color = MsgBarColors.Red;

    // If showing a msg indicating an error connecting with the hw daemon, the msg is modified
    // to make sure it includes the download URL.
    if ((body as OperationError).type && (body as OperationError).type === HWOperationResults.DaemonConnectionError) {
      config.text = 'hardware-wallet.errors.daemon-connection-with-configurable-link';
      config.link = AppConfig.hwWalletDaemonDownloadUrl;
    }

    this.show(config);
    this.setTimer(duration);
  }

  /**
   * Displays the msg bar for showing a warning.
   * @param body Text to show .
   * @param duration How much time the msg bar will stay visible, in ms. The default
   * value is 20000.
   */
  showWarning(body: string, duration = 20000) {
    const config = new MsgBarConfig();
    config.text = processServiceError(body).translatableErrorMsg;
    config.title = 'common.warning-title';
    config.icon = MsgBarIcons.Warning;
    config.color = MsgBarColors.Yellow;

    this.show(config);
    this.setTimer(duration);
  }

  /**
   * Displays the msg bar for showing a confirmation.
   * @param body Text to show .
   * @param duration How much time the msg bar will stay visible, in ms. The default
   * value is 10000.
   */
  showDone(body: string, duration = 10000) {
    const config = new MsgBarConfig();
    config.text = body;
    config.title = 'common.done-title';
    config.icon = MsgBarIcons.Done;
    config.color = MsgBarColors.Green;

    this.show(config);
    this.setTimer(duration);
  }

  /**
   * Makes the msg bar to appear.
   * @param config Data about how to display the bar.
   */
  private show(config: MsgBarConfig) {
    if (this.msgBarComponentInternal) {
      this.msgBarComponentInternal.config = config;
      this.msgBarComponentInternal.show();
    }
  }

  /**
   * Starts a timer which will close the msg bar after some time. If a timer was still
   * running when calling this function, it is cancelled before starting the new one.
   * @param duration Timer duration, in ms.
   */
  private setTimer(duration = 10000) {
    if (this.timeSubscription) {
      this.timeSubscription.unsubscribe();
    }

    this.timeSubscription = of(1).pipe(delay(duration)).subscribe(() => this.hide());
  }
}
