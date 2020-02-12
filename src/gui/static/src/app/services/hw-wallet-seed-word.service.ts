import { Injectable } from '@angular/core';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

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
    return this.requestWordComponentInternal.openDialog(this.dialog, <MatDialogConfig> {
      isForHwWallet: true,
      wordNumber: 0,
      restoringSoftwareWallet: false,
    }).afterClosed().pipe(map(word => {
      return word;
    }));
  }
}
