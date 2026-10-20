import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Observable } from 'rxjs';

/**
 * Asks for one seed word when the hardware wallet requests it, while restoring or
 * checking a seed. The device shows which word it wants.
 */
@Injectable()
export class HwWalletSeedWordService {
  // Set by WalletShellComponent, which avoids a circular import.
  private requestWordComponentInternal: any;
  set requestWordComponent(value: any) {
    this.requestWordComponentInternal = value;
  }

  constructor(private dialog: MatDialog) {}

  /** The word entered, or undefined when cancelled. */
  requestWord(): Observable<string | undefined> {
    return this.requestWordComponentInternal.openDialog(this.dialog).afterClosed();
  }
}
