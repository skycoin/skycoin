import { Injectable, NgZone } from '@angular/core';
import { MsgBarConfig, MsgBarComponent, MsgBarIcons, MsgBarColors } from '../components/layout/msg-bar/msg-bar.component';
import { parseResponseMessage } from '../utils/errors';
import { Subscription, of } from 'rxjs';
import { delay } from 'rxjs';

@Injectable()
export class MsgBarService {

  private timeSubscription: Subscription;

  private msgBarComponentInternal: MsgBarComponent;
  set msgBarComponent(value: MsgBarComponent) {
    this.msgBarComponentInternal = value;
  }

  constructor(
    private _ngZone: NgZone,
  ) {}

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

    this._ngZone.runOutsideAngular(() => {
      this.timeSubscription = of(1).pipe(delay(duration)).subscribe(() => this._ngZone.run(() => this.hide()));
    });
  }
}
