import { ChangeDetectionStrategy, ChangeDetectorRef, Component, EventEmitter, Inject, OnDestroy, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Subscription, of } from 'rxjs';
import { delay } from 'rxjs';

import { Wallet } from '../../../../app.datatypes';
import { WalletService, ScanProgressData } from '../../../../services/wallet/wallet.service';
import { config } from '../../../../app.config';

@Component({
    selector: 'app-scan-addresses',
    templateUrl: './scan-addresses.component.html',
    styleUrls: ['./scan-addresses.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class ScanAddressesComponent implements OnInit, OnDestroy {
  progress = new ScanProgressData();
  showSlowMobileInfo = false;

  private subscriptionsGroup: Subscription[] = [];
  private slowInfoSubscription!: Subscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) private data: Wallet,
    public dialogRef: MatDialogRef<ScanAddressesComponent>,
    private walletService: WalletService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {}

  ngOnInit() {
    this.scan();
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
    this.removeSlowInfoSubscription();
  }

  closePopup(result: Error | null = null) {
    this.dialogRef.close(result);
  }

  scan() {
    this.createSlowInfoSubscription();

    const onProgressChanged = new EventEmitter<ScanProgressData>();
    this.subscriptionsGroup.push(onProgressChanged.subscribe((progress: ScanProgressData) => {
      this.createSlowInfoSubscription();
      this.progress = progress;
      this.changeDetectorRef.markForCheck();
    }));

    this.subscriptionsGroup.push(this.walletService.scanAddresses(this.data, onProgressChanged)
      .subscribe(
        () => {
          this.closePopup(null);
          this.changeDetectorRef.markForCheck();
        },
        (error: Error) => {
          this.closePopup(error);
          this.changeDetectorRef.markForCheck();
        }
      )
    );
  }

  private createSlowInfoSubscription() {
    this.removeSlowInfoSubscription();

    this.slowInfoSubscription = of(1).pipe(delay(config.timeBeforeSlowMobileInfo))
      .subscribe(() => {
        this.showSlowMobileInfo = true;
        this.changeDetectorRef.markForCheck();
      });
  }

  private removeSlowInfoSubscription() {
    if (this.slowInfoSubscription) {
      this.slowInfoSubscription.unsubscribe();
    }
  }
}
