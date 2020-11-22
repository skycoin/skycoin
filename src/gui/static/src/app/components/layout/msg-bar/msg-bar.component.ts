import { Component } from '@angular/core';

export enum MsgBarIcons {
  Error = 'error',
  Done = 'done',
  Warning = 'warning',
}

export enum MsgBarColors {
  Red = 'red-background',
  Green = 'green-background',
  Yellow = 'yellow-background',
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
    if (this.visible) {
      this.visible = false;
      setTimeout(() => this.visible = true, 32);
    } else {
      this.visible = true;
    }
  }

  hide() {
    this.visible = false;
  }
}
