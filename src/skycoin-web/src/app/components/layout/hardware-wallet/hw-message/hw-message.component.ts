import { Component, Input, Output, EventEmitter, OnDestroy, ChangeDetectionStrategy } from '@angular/core';

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
    changeDetection: ChangeDetectionStrategy.Eager,
    standalone: false
})
export class HwMessageComponent implements OnDestroy {
  @Input() icon: MessageIcons = MessageIcons.None;
  @Input() text!: string;
  @Input() outputsList!: HwWalletTxRecipientData[];
  @Input() lowerText!: string;
  @Output() linkClicked = new EventEmitter();

  icons = MessageIcons;

  ngOnDestroy() {
    this.linkClicked.complete();
  }
}
