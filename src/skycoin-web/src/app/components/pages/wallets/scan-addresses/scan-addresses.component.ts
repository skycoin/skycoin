import { Component, EventEmitter, Inject, OnInit, OnDestroy } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Subscription } from 'rxjs';
import { Observable } from 'rxjs';

import { Wallet } from '../../../../app.datatypes';
import { WalletService, ScanProgressData } from '../../../../services/wallet/wallet.service';
import { config } from '../../../../app.config';

@Component({
    selector: 'app-scan-addresses',
    templateUrl: './scan-addresses.component.html',
    styleUrls: ['./scan-addresses.component.scss'],
    standalone: false
})
export class ScanAddressesComponent implements OnInit, OnDestroy {
  progress = new ScanProgressData();
  showSlowMobileInfo = false;

  private subscriptionsGroup: Subscription[] = [];
  private slowInfoSubscription: Subscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) private data: Wallet,
    public dialogRef: MatDialogRef<ScanAddressesComponent>,
    private walletService: WalletService
  ) {}

  ngOnInit() {
    this.scan();
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
    this.removeSlowInfoSubscription();
  }

  closePopup(result = null) {
    this.dialogRef.close(result);
  }

  scan() {
    this.createSlowInfoSubscription();

    const onProgressChanged = new EventEmitter<ScanProgressData>();
    this.subscriptionsGroup.push(onProgressChanged.subscribe((progress: ScanProgressData) => {
      this.createSlowInfoSubscription();
      this.progress = progress;
    }));

    this.subscriptionsGroup.push(this.walletService.scanAddresses(this.data, onProgressChanged)
      .subscribe(
        () => this.closePopup(null),
        (error: Error) => this.closePopup(error)
      )
    );
  }

  private createSlowInfoSubscription() {
    this.removeSlowInfoSubscription();

    this.slowInfoSubscription = Observable.of(1).delay(config.timeBeforeSlowMobileInfo)
      .subscribe(() => this.showSlowMobileInfo = true);
  }

  private removeSlowInfoSubscription() {
    if (this.slowInfoSubscription) {
      this.slowInfoSubscription.unsubscribe();
    }
  }
}
