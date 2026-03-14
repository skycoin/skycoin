import { Component, ElementRef, Inject, OnInit, ViewChild, OnDestroy } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef, MatDialogConfig } from '@angular/material/dialog';
import { UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { Subject, Subscription, debounceTime } from 'rxjs';

import { CoinService } from '../../../services/coin.service';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';
import { MsgBarService } from '../../../services/msg-bar.service';
import { ClipboardService } from '../../../services/clipboard.service';
import { Router } from '@angular/router';

declare var QRCode: any;

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
  showExtraAddressOptions?: boolean;
}

@Component({
    selector: 'app-qr-code',
    templateUrl: './qr-code.component.html',
    styleUrls: ['./qr-code.component.scss'],
    standalone: false
})
export class QrCodeComponent implements OnInit, OnDestroy {
  @ViewChild('qr') qr: ElementRef;

  form: UntypedFormGroup;
  currentQrContent: string;
  showForm = false;
  invalidCoins = false;
  invalidHours = false;

  private defaultQrConfig = new DefaultQrConfig();
  private subscriptionsGroup: Subscription[] = [];
  private updateQrEvent: Subject<boolean> = new Subject<boolean>();

  static openDialog(dialog: CustomMatDialogService, config: QrDialogConfig) {
    const dialogConfig = new MatDialogConfig();
    dialogConfig.data = config;
    dialogConfig.width = '390px';
    dialogConfig.autoFocus = false;
    dialog.open(QrCodeComponent, dialogConfig);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: QrDialogConfig,
    public dialogRef: MatDialogRef<QrCodeComponent>,
    public formBuilder: UntypedFormBuilder,
    private coinService: CoinService,
    private msgBarService: MsgBarService,
    private clipboardService: ClipboardService,
    private router: Router,
  ) { }

  ngOnInit() {
    this.initForm();
    this.updateQrContent();
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
    this.msgBarService.hide();
  }

  startShowingForm() {
    this.showForm = true;
  }

  copyText(text: string) {
    this.clipboardService.copy(text).then(() => this.msgBarService.showDone('qr.copied', 4000));
  }

  closePopup() {
    this.dialogRef.close();
  }

  goToDetail(path: string) {
    this.router.navigate([path], { queryParams: { addr: this.data.address }});
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
    this.currentQrContent = (!this.data.ignoreCoinPrefix ? (this.coinService.currentCoin.value.coinName.toLowerCase() + ':') : '') + this.data.address;

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
      if (Number.parseInt(hours).toString() === hours && Number.parseInt(hours) > 0) {
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
