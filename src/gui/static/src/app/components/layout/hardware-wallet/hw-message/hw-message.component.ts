import { Component, Input, Output, EventEmitter } from '@angular/core';
import { HwWalletTxRecipientData } from '../../../../services/hw-wallet.service';

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
  @Input() outputsList: HwWalletTxRecipientData[];
  @Input() lowerText: string;
  @Input() linkText: string;
  @Input() linkIsUrl = false;
  @Input() upperBigText: string;
  @Input() lowerBigText: string;
  @Input() lowerLightText: string;
  @Output() linkClicked = new EventEmitter();

  icons = MessageIcons;

  activateLink() {
    this.linkClicked.emit();
  }
}
