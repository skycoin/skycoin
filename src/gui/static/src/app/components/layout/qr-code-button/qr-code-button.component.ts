import { Component, Input } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { QrCodeComponent, QrDialogConfig } from '../qr-code/qr-code.component';

@Component({
  selector: 'app-qr-code-button',
  templateUrl: './qr-code-button.component.html',
  styleUrls: ['./qr-code-button.component.scss'],
})
export class QrCodeButtonComponent {
  @Input() address: string;
  @Input() hideCoinRequestForm = false;
  @Input() ignoreCoinPrefix = false;

  constructor(
    private dialog: MatDialog,
  ) { }

  showQrCode(event) {
    event.stopPropagation();

    const config: QrDialogConfig = {
      address: this.address,
      hideCoinRequestForm: this.hideCoinRequestForm,
      ignoreCoinPrefix: this.ignoreCoinPrefix,
    };

    QrCodeComponent.openDialog(this.dialog, config);
  }
}
