import { Component, Inject, OnInit, ViewChild } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatDialogRef } from '@angular/material';

declare const QRCode: any;

@Component({
  selector: 'app-qr-code',
  templateUrl: './qr-code.component.html',
  styleUrls: ['./qr-code.component.css'],
})
export class QrCodeComponent implements OnInit {
  @ViewChild('qr') qr: any;

  size = 300;
  level = 'M';
  colordark = '#000000';
  colorlight = '#ffffff';
  usesvg = false;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: any,
    public dialogRef: MatDialogRef<QrCodeComponent>,
  ) { }

  ngOnInit() {
    const qrcode = new QRCode(this.qr.nativeElement, {
      text: `skycoin:${this.data.address}`,
      width: this.size,
      height: this.size,
      colorDark: this.colordark,
      colorLight: this.colorlight,
      useSVG: this.usesvg,
      correctLevel: QRCode.CorrectLevel[this.level.toString()],
    });
  }
}
