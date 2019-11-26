import { Component, Input, Output, EventEmitter } from '@angular/core';
import { TxData } from '../../../../services/hw-wallet.service';

export enum MessageIcons {
  None,
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
  @Input() outputsList: TxData[];
  @Input() lowerText: string;
  @Input() linkText: string;
  @Input() upperBigText: string;
  @Input() lowerBigText: string;
  @Input() lowerLightText: string;
  @Output() linkClicked = new EventEmitter();

  icons = MessageIcons;

  activateLink() {
    this.linkClicked.emit();
  }
}
