import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { SeedWordDialogParams } from '../components/layout/seed-word-dialog/seed-word-dialog.component';

/**
 * Allows to easily show the modal window used for entering a seed word requested
 * by the hw wallet.
 */
@Injectable()
export class HwWalletSeedWordService {

  // Set to SeedWordDialogComponent on AppComponent to avoid a circular reference.
  private requestWordComponentInternal;
  /**
   * Sets the class of the modal window used for entering a seed word requested
   * by the hw wallet.
   */
  set requestWordComponent(value) {
    this.requestWordComponentInternal = value;
  }

  constructor(
    private dialog: MatDialog,
  ) {}

  /**
   * Shows the modal window used for entering a seed word requested by the hw wallet.
   * @returns The word entered by the user, or null if the user cancelled the operation.
   */
  requestWord(): Observable<string> {
    return this.requestWordComponentInternal.openDialog(this.dialog, {
      reason: 'HWWalletOperation',
    } as SeedWordDialogParams).afterClosed().pipe(map(word => word));
  }
}
