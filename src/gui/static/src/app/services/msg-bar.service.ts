import { Injectable } from '@angular/core';
import { MsgBarConfig, MsgBarComponent, MsgBarIcons, MsgBarColors } from '../components/layout/msg-bar/msg-bar.component';
import { parseResponseMessage } from '../utils/errors';

@Injectable()
export class MsgBarService {

  private msgBarComponentInternal: MsgBarComponent;
  set msgBarComponent(value: MsgBarComponent) {
    this.msgBarComponentInternal = value;
  }

  show(config: MsgBarConfig) {
    this.msgBarComponentInternal.config = config;
    this.msgBarComponentInternal.show();
  }

  hide() {
    this.msgBarComponentInternal.hide();
  }

  showError(body: string) {
    const config = new MsgBarConfig();
    config.text = parseResponseMessage(body);
    config.title = 'errors.error';
    config.icon = MsgBarIcons.error;
    config.color = MsgBarColors.red;

    this.show(config);
  }
}
