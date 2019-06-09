import { Component } from '@angular/core';

export enum MsgBarIcons {
  error = 'error',
}

export enum MsgBarColors {
  red = 'red-background',
}

export class MsgBarConfig {
  title?: string;
  text: string;
  icon?: MsgBarIcons;
  color?: MsgBarColors;
}

@Component({
  selector: 'app-msg-bar',
  templateUrl: './msg-bar.component.html',
  styleUrls: ['./msg-bar.component.scss'],
})
export class MsgBarComponent {
  config = new MsgBarConfig();
  visible = false;

  constructor() { }

  show() {
    this.visible = true;
  }

  hide() {
    this.visible = false;
  }
}
