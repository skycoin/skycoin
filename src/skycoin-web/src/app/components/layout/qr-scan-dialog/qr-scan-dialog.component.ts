import { ChangeDetectionStrategy, ChangeDetectorRef, Component, ElementRef, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Observable } from 'rxjs';

import { decodeQr } from '../../../utils/qr';

/** Frames decoded per second while the camera is open. */
const scansPerSecond = 5;

/**
 * Reads a QR code with the camera, or from a picked image. Closes with the
 * decoded text, or with undefined when cancelled. Use QrScanDialogComponent.scan.
 */
@Component({
    selector: 'app-qr-scan-dialog',
    templateUrl: './qr-scan-dialog.component.html',
    styleUrls: ['./qr-scan-dialog.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class QrScanDialogComponent implements OnInit, OnDestroy {
  @ViewChild('video', { static: true }) video!: ElementRef<HTMLVideoElement>;

  // A translation key saying why the camera is not showing, if it is not.
  cameraProblem = '';
  noCodeInImage = false;

  private stream: MediaStream | null = null;
  private timer: any = null;
  private canvas = document.createElement('canvas');

  static scan(dialog: MatDialog): Observable<string | undefined> {
    return dialog.open(QrScanDialogComponent, { width: '450px' }).afterClosed();
  }

  /** Whether this page may use a camera: plain http only counts on 127.0.0.1 and localhost. */
  static get cameraAllowed(): boolean {
    return !!(navigator.mediaDevices && navigator.mediaDevices.getUserMedia);
  }

  constructor(
    public dialogRef: MatDialogRef<QrScanDialogComponent>,
    private changeDetectorRef: ChangeDetectorRef,
  ) {}

  ngOnInit() {
    if (!QrScanDialogComponent.cameraAllowed) {
      this.cameraProblem = 'qr-scan.insecure';
      return;
    }
    navigator.mediaDevices.getUserMedia({ video: { facingMode: 'environment' }, audio: false }).then(stream => {
      this.stream = stream;
      const video = this.video.nativeElement;
      video.srcObject = stream;
      video.play().catch(() => {});
      this.timer = setInterval(() => {
        this.scanFrame();
        this.changeDetectorRef.markForCheck();
      }, 1000 / scansPerSecond);
      this.changeDetectorRef.markForCheck();
    }, () => {
      this.cameraProblem = 'qr-scan.no-camera';
      this.changeDetectorRef.markForCheck();
    });
  }

  ngOnDestroy() {
    if (this.timer) {
      clearInterval(this.timer);
    }
    if (this.stream) {
      this.stream.getTracks().forEach(t => t.stop());
    }
  }

  pickImage(event: Event) {
    const file = (event.target as HTMLInputElement).files?.[0];
    if (!file) {
      return;
    }
    const url = URL.createObjectURL(file);
    const img = new Image();
    img.onload = () => {
      decodeQr(img, this.canvas).then(text => {
        URL.revokeObjectURL(url);
        this.noCodeInImage = !text;
        if (text) {
          this.dialogRef.close(text);
        }
        this.changeDetectorRef.markForCheck();
      });
    };
    img.src = url;
  }

  private scanning = false;

  private scanFrame() {
    if (this.scanning) {
      return;
    }
    this.scanning = true;
    decodeQr(this.video.nativeElement, this.canvas).then(text => {
      this.scanning = false;
      if (text) {
        this.dialogRef.close(text);
      }
      this.changeDetectorRef.markForCheck();
    }, () => {
      this.scanning = false;
      this.changeDetectorRef.markForCheck();
    });
  }
}
