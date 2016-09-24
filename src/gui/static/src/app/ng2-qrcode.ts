import {Component, Inject, Input, ElementRef, OnInit} from 'angular2/core';

declare var QRCode: any;

@Component({
  selector: 'qrcode',
  template: ''
})
export class QRCodeComponent implements OnInit {
  @Input() qrdata: String = ''
  @Input() size: Number = 256
  @Input() level: String = 'M'
  @Input() colordark: String = '#000000'
  @Input() colorlight: String = '#ffffff'
  @Input() usesvg: Boolean = false

  constructor(
    private el: ElementRef
  ) { }

  ngOnInit() {
    try {
      if(this.qrdata===''){
        throw new Error("Empty QR Code data");
      }
      new QRCode(this.el.nativeElement, {
        text: this.qrdata,
        width: this.size,
        height: this.size,
        colorDark: this.colordark,
        colorLight: this.colorlight,
        useSVG: this.usesvg,
        correctLevel: QRCode.CorrectLevel[this.level.toString()]
      });
    }
    catch (e) {
      console.error("Error generating QR Code: " + e.message);
    }
  }
}
