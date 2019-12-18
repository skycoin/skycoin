import { Injectable } from '@angular/core';
import { MsgBarConfig, MsgBarComponent, MsgBarIcons, MsgBarColors } from '../components/layout/msg-bar/msg-bar.component';
import { parseResponseMessage } from '../utils/errors';
import { SubscriptionLike, of } from 'rxjs';
import { delay } from 'rxjs/operators';
import { AppConfig } from '../app.config';

@Injectable()
export class MsgBarService {

  private timeSubscription: SubscriptionLike;

  private msgBarComponentInternal: MsgBarComponent;
  set msgBarComponent(value: MsgBarComponent) {
    this.msgBarComponentInternal = value;
  }

  show(config: MsgBarConfig) {
    if (config.text === 'hardware-wallet.errors.daemon-connection' || config.text.indexOf('Problem connecting to the Skywallet Daemon') !== -1) {
      config.text = 'hardware-wallet.errors.daemon-connection-with-configurable-link';
      config.link = AppConfig.hwWalletDaemonDownloadUrl;
    }

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

  showError(body: string, duration = 20000) {
    const config = new MsgBarConfig();
    config.text = parseResponseMessage(body);
    config.title = 'errors.error';
    config.icon = MsgBarIcons.Error;
    config.color = MsgBarColors.Red;

    this.show(config);
    this.setTimer(duration);
  }

  showWarning(body: string, duration = 20000) {
    const config = new MsgBarConfig();
    config.text = parseResponseMessage(body);
    config.title = 'common.warning';
    config.icon = MsgBarIcons.Warning;
    config.color = MsgBarColors.Yellow;

    this.show(config);
    this.setTimer(duration);
  }

  showDone(body: string, duration = 10000) {
    const config = new MsgBarConfig();
    config.text = body;
    config.title = 'common.success';
    config.icon = MsgBarIcons.Done;
    config.color = MsgBarColors.Green;

    this.show(config);
    this.setTimer(duration);
  }

  private setTimer(duration = 10000) {
    if (this.timeSubscription) {
      this.timeSubscription.unsubscribe();
    }

    this.timeSubscription = of(1).pipe(delay(duration)).subscribe(() => this.hide());
  }
}
