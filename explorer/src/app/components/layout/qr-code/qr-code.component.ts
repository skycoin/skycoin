import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';

import { QrConfig } from 'app/app.config';

/**
 * Allows to access the QRcode generator located on src/js/qrcode.min.js.
 */
declare let QRCode: any; // eslint-disable-line @typescript-eslint/naming-convention

/**
 * Shows a QR code.
 */
@Component({
  selector: 'app-qr-code',
  templateUrl: './qr-code.component.html',
  styleUrls: ['./qr-code.component.scss']
})
export class QrCodeComponent implements OnInit {
  /**
   * Content of the QR code. The content of QrConfig.prefix will be added at the left.
   */
  @Input() set string(val: string) {
    this.stringInternal = val;
    this.updateQR();
  }
  private stringInternal: string;

  @ViewChild('qr', { static: true }) qr: ElementRef;

  // Size in pixels.
  size = 130;
  // Error correction level.
  level = 'M';
  // QR code color.
  colordark = '#000000';
  // Background code color.
  colorlight = '#ffffff';
  usesvg = false;

  ngOnInit() {
    this.updateQR();
  }

  private updateQR() {
    (this.qr.nativeElement as HTMLDivElement).innerHTML = '';
    const qr = new QRCode(this.qr.nativeElement, {
      text: QrConfig.prefix + this.stringInternal,
      width: this.size,
      height: this.size,
      colorDark: this.colordark,
      colorLight: this.colorlight,
      useSVG: this.usesvg,
      correctLevel: QRCode.CorrectLevel[this.level.toString()]
    });
  }
}
