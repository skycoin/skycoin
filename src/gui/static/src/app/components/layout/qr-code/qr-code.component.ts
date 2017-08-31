import { Component, Inject, Input, ElementRef, OnInit, ViewChild } from '@angular/core';
import { MD_DIALOG_DATA } from '@angular/material';

declare var QRCode: any;

@Component({
  selector: 'app-qr-code',
  templateUrl: './qr-code.component.html',
  styleUrls: ['./qr-code.component.css']
})
export class QrCodeComponent implements OnInit {
  @ViewChild('qr') qr: any;

  size: number = 300;
  level: string = 'M';
  colordark: string = '#000000';
  colorlight: string = '#ffffff';
  usesvg: boolean = false;

  constructor(
    @Inject(MD_DIALOG_DATA) private data: any,
    private el: ElementRef
  ) { }

  ngOnInit() {
    new QRCode(this.qr.nativeElement, {
      text: this.data.address,
      width: this.size,
      height: this.size,
      colorDark: this.colordark,
      colorLight: this.colorlight,
      useSVG: this.usesvg,
      correctLevel: QRCode.CorrectLevel[this.level.toString()]
    });
  }
}
