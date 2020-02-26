import { Component } from '@angular/core';

/**
 * Icons the msg bar can show.
 */
export enum MsgBarIcons {
  Error = 'error',
  Done = 'done',
  Warning = 'warning',
}

// The enum has the names of the classes used for displaying the colors.
/**
 * Colors the msg bar can show.
 */
export enum MsgBarColors {
  Red = 'red-background',
  Green = 'green-background',
  Yellow = 'yellow-background',
}

/**
 * Settings for the msg bar.
 */
export class MsgBarConfig {
  title?: string;
  text: string;
  /**
   * If set, the link will be shown after the text. Must be a valid URL.
   */
  link?: string;
  icon?: MsgBarIcons;
  color?: MsgBarColors;
}

/**
 * Small horizontal bar used by the app for showing notifications. It was made to be created
 * and added to the UI just after starting the app and being controlled by a service.
 */
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
      // Gives the illusion of the bar reapering if it was already being shown.
      setTimeout(() => this.visible = true, 32);
    } else {
      this.visible = true;
    }
  }

  hide() {
    this.visible = false;
  }
}
