import { Injectable } from '@angular/core';
import { MsgBarConfig, MsgBarComponent, MsgBarIcons, MsgBarColors } from '../components/layout/msg-bar/msg-bar.component';
import { processServiceError } from '../utils/errors';
import { SubscriptionLike, of } from 'rxjs';
import { delay } from 'rxjs/operators';
import { AppConfig } from '../app.config';
import { environment } from '../../environments/environment';
import { OperationError, HWOperationResults } from '../utils/operation-error';


@Injectable()
export class MsgBarService {

  private timeSubscription: SubscriptionLike;

  private msgBarComponentInternal: MsgBarComponent;
  set msgBarComponent(value: MsgBarComponent) {
    this.msgBarComponentInternal = value;
  }

  show(config: MsgBarConfig) {
    if (this.msgBarComponentInternal) {
      this.msgBarComponentInternal.config = config;
      this.msgBarComponentInternal.show();
    }
  }

  hide() {
    if (this.msgBarComponentInternal) {
      this.msgBarComponentInternal.hide();
    }
  }

  showError(body: string | OperationError, duration = 20000) {
    const config = new MsgBarConfig();
    config.text = processServiceError(body).translatableErrorMsg;
    config.title = 'common.error-title';
    config.icon = MsgBarIcons.Error;
    config.color = MsgBarColors.Red;

    if ((body as OperationError).type && (body as OperationError).type === HWOperationResults.DaemonError) {
      config.text = 'hardware-wallet.errors.daemon-connection-with-configurable-link';
      config.link = AppConfig.hwWalletDaemonDownloadUrl;
    }

    this.show(config);
    this.setTimer(duration);
  }

  showWarning(body: string, duration = 20000) {
    const config = new MsgBarConfig();
    config.text = processServiceError(body).translatableErrorMsg;
    config.title = 'common.warning-title';
    config.icon = MsgBarIcons.Warning;
    config.color = MsgBarColors.Yellow;

    this.show(config);
    this.setTimer(duration);
  }

  showDone(body: string, duration = 10000) {
    const config = new MsgBarConfig();
    config.text = body;
    config.title = 'common.done-title';
    config.icon = MsgBarIcons.Done;
    config.color = MsgBarColors.Green;

    this.show(config);
    this.setTimer(duration);
  }

  private setTimer(duration = 10000) {
    if (this.timeSubscription) {
      this.timeSubscription.unsubscribe();
    }

    if (environment.isInE2eMode) {
      duration = 500;
    }

    this.timeSubscription = of(1).pipe(delay(duration)).subscribe(() => this.hide());
  }
}
