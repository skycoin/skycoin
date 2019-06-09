import { Injectable } from '@angular/core';
import { MsgBarConfig, MsgBarComponent } from '../components/layout/msg-bar/msg-bar.component';

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
}
