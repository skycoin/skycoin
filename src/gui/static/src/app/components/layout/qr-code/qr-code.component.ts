import { Component, Inject, ViewChild, OnDestroy, ElementRef, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { MatDialogRef } from '@angular/material/dialog';
import { FormGroup, FormBuilder } from '@angular/forms';
import { SubscriptionLike, Subject } from 'rxjs';
import { debounceTime } from 'rxjs/operators';
import { BigNumber } from 'bignumber.js';

import { copyTextToClipboard } from '../../../utils/general-utils';
import { AppConfig } from '../../../app.config';
import { MsgBarService } from '../../../services/msg-bar.service';
import { AppService } from '../../../services/app.service';

// Gives access to qrcode.js, imported from the resources folder.
declare const QRCode: any;

/**
 * Default QR code graphical config.
 */
class DefaultQrConfig {
  static readonly size = 180;
  static readonly level = 'M';
  static readonly colordark = '#000000';
  static readonly colorlight = '#ffffff';
  static readonly usesvg = false;
}

/**
 * Settings for QrCodeComponent.
 */
export interface QrDialogConfig {
  /**
   * Address the QR code will have.
   */
  address: string;
  /**
   * If true, the modal window will not show the coin request form and the addreess will not
   * have the BIP21 prefix.
   */
  showAddressOnly: boolean;
}

/**
 * Modal window used for showing QR codes.
 */
@Component({
  selector: 'app-qr-code',
  templateUrl: './qr-code.component.html',
  styleUrls: ['./qr-code.component.scss'],
})
export class QrCodeComponent implements OnInit, OnDestroy {
  @ViewChild('qrArea') qrArea: ElementRef;

  form: FormGroup;
  currentQrContent: string;
  formVisible = false;
  // For knowing if the form fields have errors.
  invalidCoins = false;
  invalidHours = false;

  private subscriptionsGroup: SubscriptionLike[] = [];
  // Emits every time the content of the QR code must be updated.
  private updateQrEvent: Subject<boolean> = new Subject<boolean>();

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
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
    private appService: AppService,
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

  showForm() {
    this.formVisible = true;
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

    // Each time a field is updated, update the content of the QR, but wait a prudential time.
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

  /**
   * Updates the content of the QR code.
   */
  private updateQrContent() {
    // If true, the QR only contains the address.
    if (this.data.showAddressOnly) {
      this.currentQrContent = this.data.address;
      this.updateQrCode();

      return;
    }

    // Add the BIP21 prefix.
    this.currentQrContent = AppConfig.uriSpecificatioPrefix.toLowerCase() + ':' + this.data.address;

    this.invalidCoins = false;
    this.invalidHours = false;

    let nextSeparator = '?';

    // Add the coins or alert if the value is not valid.
    if (this.form.get('coins').value) {
      const coins = new BigNumber(this.form.get('coins').value);
      if (!coins.isNaN() && coins.isGreaterThan(0) && coins.decimalPlaces() <= this.appService.currentMaxDecimals) {
        this.currentQrContent += nextSeparator + 'amount=' + coins.toString();
        nextSeparator = '&';
      } else {
        this.invalidCoins = true;
      }
    }

    // Add the hours or alert if the value is not valid.
    if (this.form.get('hours').value) {
      const hours = new BigNumber(this.form.get('hours').value);
      if (!hours.isNaN() && hours.isGreaterThan(0) && hours.decimalPlaces() === 0) {
        this.currentQrContent += nextSeparator + 'hours=' + hours.toString();
        nextSeparator = '&';
      } else {
        this.invalidHours = true;
      }
    }

    // Add the note.
    const note = this.form.get('note').value;
    if (note) {
      this.currentQrContent += nextSeparator + 'message=' + encodeURIComponent(note);
    }

    // Update the QR code image.
    this.updateQrCode();
  }

  private updateQrCode() {
    // Clean the area of the QR code.
    (this.qrArea.nativeElement as HTMLDivElement).innerHTML = '';

    // Creates a new QR code and adds it to the designated area.
    const qrCode = new QRCode(this.qrArea.nativeElement, {
      text: this.currentQrContent,
      width: DefaultQrConfig.size,
      height: DefaultQrConfig.size,
      colorDark: DefaultQrConfig.colordark,
      colorLight: DefaultQrConfig.colorlight,
      useSVG: DefaultQrConfig.usesvg,
      correctLevel: QRCode.CorrectLevel[DefaultQrConfig.level],
    });
  }
}
