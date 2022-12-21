import { Component, Input } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';

import { QrCodeComponent, QrDialogConfig } from '../qr-code/qr-code.component';

/**
 * Simple inline img button which opens the QR code modal window when pressed.
 */
@Component({
  selector: 'app-qr-code-button',
  templateUrl: './qr-code-button.component.html',
  styleUrls: ['./qr-code-button.component.scss'],
})
export class QrCodeButtonComponent {
  // Address the QR code modal window will show.
  @Input() address: string;
  // If true, the QR code modal window will not show the coin request form and the addreess
  // will not have the BIP21 prefix.
  @Input() showAddressOnly = false;

  constructor(
    private dialog: MatDialog,
  ) { }

  showQrCode(event) {
    event.stopPropagation();

    const config: QrDialogConfig = {
      address: this.address,
      showAddressOnly: this.showAddressOnly,
    };

    QrCodeComponent.openDialog(this.dialog, config);
  }
}
