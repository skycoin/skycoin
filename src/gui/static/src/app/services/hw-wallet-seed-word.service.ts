import { Injectable } from '@angular/core';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class HwWalletSeedWordService {

  // Set on AppComponent to avoid a circular reference.
  private requestWordComponentInternal;
  set requestWordComponent(value) {
    this.requestWordComponentInternal = value;
  }

  constructor(
    private dialog: MatDialog,
  ) {}

  requestWord(): Observable<string> {
    return this.dialog.open(this.requestWordComponentInternal, <MatDialogConfig> {
      width: '350px',
      data: {
        isForHwWallet: true,
        wordNumber: 0,
        restoringSoftwareWallet: false,
      },
    }).afterClosed().map(word => {
      return word;
    });
  }
}
