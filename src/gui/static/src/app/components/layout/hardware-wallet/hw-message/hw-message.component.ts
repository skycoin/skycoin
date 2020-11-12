import { Component, Input, Output, EventEmitter } from '@angular/core';

import { HwWalletTxRecipientData } from '../../../../services/hw-wallet.service';

/**
 * Icons HwMessageComponent can show.
 */
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

/**
 * Generic control for showing messages. It includes an area for an icon on the left side
 * and various text with different format. All properties are optional.
 */
@Component({
  selector: 'app-hw-message',
  templateUrl: './hw-message.component.html',
  styleUrls: ['./hw-message.component.scss'],
})
export class HwMessageComponent {
  // Icon to show at the left.
  @Input() icon: MessageIcons = MessageIcons.None;
  // Text to show.
  @Input() text: string;
  // Link to show after the main text.
  @Input() linkText: string;
  // URL for the link. If no URL is set, the linkClicked event is dispatched when the user
  // clicks the link.
  @Input() linkIsUrl = false;
  // Outputs to show as a list. It is used to show the list of coins going out during a
  // transaction, so the user can check it before confirming the tx. It is shown under
  // the main text.
  @Input() outputsList: HwWalletTxRecipientData[];
  // Text shown after the main text, the link and the outputs list.
  @Input() lowerText: string;
  // Big text shown at the top of the main text. Can be used as a title.
  @Input() upperBigText: string;
  // Big text shown below the main text.
  @Input() lowerBigText: string;
  // Light text shown at the bottom of the component.
  @Input() lowerLightText: string;
  // Event dispatched when the user clicks the link if no URL for it was provided.
  @Output() linkClicked = new EventEmitter();

  icons = MessageIcons;

  activateLink() {
    this.linkClicked.emit();
  }
}
