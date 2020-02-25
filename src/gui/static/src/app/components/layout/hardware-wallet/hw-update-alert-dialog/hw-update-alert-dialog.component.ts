import { Component } from '@angular/core';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';

import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

/**
 * Modal window for alerting the user that there is a firmware update available. If the user
 * selects to install the update, the modal window is closed and "true" is returned in the
 * "afterClosed" envent.
 */
@Component({
  selector: 'app-hw-update-alert-dialog',
  templateUrl: './hw-update-alert-dialog.component.html',
  styleUrls: ['./hw-update-alert-dialog.component.scss'],
})
export class HwUpdateAlertDialogComponent extends HwDialogBaseComponent<HwUpdateAlertDialogComponent> {
  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog): MatDialogRef<HwUpdateAlertDialogComponent, any> {
    const config = new MatDialogConfig();
    config.autoFocus = false;
    config.width = '450px';

    return dialog.open(HwUpdateAlertDialogComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<HwUpdateAlertDialogComponent>,
    hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
  }

  update() {
    this._dialogRef.close(true);
  }
}
