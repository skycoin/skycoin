import { Component, Input } from '@angular/core';

export enum MessageIcons {
  Spinner,
  Success,
  Error,
  Usb,
  HardwareWallet,
  Warning,
  Confirm,
}

@Component({
  selector: 'app-hw-message',
  templateUrl: './hw-message.component.html',
  styleUrls: ['./hw-message.component.scss'],
})
export class HwMessageComponent {
  @Input() icon: MessageIcons = MessageIcons.HardwareWallet;
  @Input() text: string;
  @Input() upperBigText: string;
  @Input() lowerBigText: string;

  icons = MessageIcons;
}
