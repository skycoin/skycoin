import { Component, Inject, ViewChild, OnDestroy, ElementRef, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { MatDialogRef } from '@angular/material/dialog';
import { FormGroup, FormBuilder } from '@angular/forms';
import { SubscriptionLike, Subject } from 'rxjs';
import { copyTextToClipboard } from '../../../utils';
import { AppConfig } from '../../../app.config';
import { MsgBarService } from '../../../services/msg-bar.service';
import { debounceTime } from 'rxjs/operators';

declare const QRCode: any;

class DefaultQrConfig {
  readonly size = 180;
  readonly level = 'M';
  readonly colordark = '#000000';
  readonly colorlight = '#ffffff';
  readonly usesvg = false;
}

export interface QrDialogConfig {
  address: string;
  ignoreCoinPrefix?: boolean;
  hideCoinRequestForm?: boolean;
}

@Component({
  selector: 'app-qr-code',
  templateUrl: './qr-code.component.html',
  styleUrls: ['./qr-code.component.scss'],
})
export class QrCodeComponent implements OnInit, OnDestroy {
  @ViewChild('qr', { static: false }) qr: ElementRef;

  form: FormGroup;
  currentQrContent: string;
  showForm = false;
  invalidCoins = false;
  invalidHours = false;

  private defaultQrConfig = new DefaultQrConfig();
  private subscriptionsGroup: SubscriptionLike[] = [];
  private updateQrEvent: Subject<boolean> = new Subject<boolean>();

  static openDialog(dialog: MatDialog, config: QrDialogConfig) {
    const dialogConfig = new MatDialogConfig();
    dialogConfig.data = config;
    dialogConfig.width = '390px';
    dialog.open(QrCodeComponent, dialogConfig);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: QrDialogConfig,
    public dialogRef: MatDialogRef<QrCodeComponent>,
    public formBuilder: FormBuilder,
    private msgBarService: MsgBarService,
  ) { }

  ngOnInit() {
    setTimeout(() => {
      this.initForm();
      this.updateQrContent();
    });
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
    this.msgBarService.hide();
  }

  startShowingForm() {
    this.showForm = true;
  }

  copyText(text) {
    copyTextToClipboard(text);
    this.msgBarService.showDone('common.copied', 4000);
  }

  private initForm() {
    this.form = this.formBuilder.group({
      coins: [''],
      hours: [''],
      note: [''],
    });

    this.subscriptionsGroup.push(this.form.get('coins').valueChanges.subscribe(this.reportValueChanged.bind(this)));
    this.subscriptionsGroup.push(this.form.get('hours').valueChanges.subscribe(this.reportValueChanged.bind(this)));
    this.subscriptionsGroup.push(this.form.get('note').valueChanges.subscribe(this.reportValueChanged.bind(this)));

    this.subscriptionsGroup.push(this.updateQrEvent.pipe(debounceTime(500)).subscribe(() => {
      this.updateQrContent();
    }));
  }

  private reportValueChanged() {
    this.updateQrEvent.next(true);
  }

  private updateQrContent() {
    this.currentQrContent = (!this.data.ignoreCoinPrefix ? (AppConfig.uriSpecificatioPrefix.toLowerCase() + ':') : '') + this.data.address;

    this.invalidCoins = false;
    this.invalidHours = false;

    let nextSeparator = '?';

    const coins = this.form.get('coins').value;
    if (coins) {
      if (Number.parseFloat(coins).toString() === coins && Number.parseFloat(coins) > 0) {
        this.currentQrContent += nextSeparator + 'amount=' + this.form.get('coins').value;
        nextSeparator = '&';
      } else {
        this.invalidCoins = true;
      }
    }

    const hours = this.form.get('hours').value;
    if (hours) {
      if (Number.parseInt(hours, 10).toString() === hours && Number.parseInt(hours, 10) > 0) {
        this.currentQrContent += nextSeparator + 'hours=' + this.form.get('hours').value;
        nextSeparator = '&';
      } else {
        this.invalidHours = true;
      }
    }

    const note = this.form.get('note').value;
    if (note) {
      this.currentQrContent += nextSeparator + 'message=' + encodeURIComponent(note);
    }

    this.updateQrCode();
  }

  private updateQrCode() {
    (this.qr.nativeElement as HTMLDivElement).innerHTML = '';

    const qrcode = new QRCode(this.qr.nativeElement, {
      text: this.currentQrContent,
      width: this.defaultQrConfig.size,
      height: this.defaultQrConfig.size,
      colorDark: this.defaultQrConfig.colordark,
      colorLight: this.defaultQrConfig.colorlight,
      useSVG: this.defaultQrConfig.usesvg,
      correctLevel: QRCode.CorrectLevel[this.defaultQrConfig.level],
    });
  }
}
