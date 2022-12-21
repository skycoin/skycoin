import { Component } from '@angular/core';
import { MatLegacyDialogRef as MatDialogRef, MatLegacyDialog as MatDialog, MatLegacyDialogConfig as MatDialogConfig } from '@angular/material/legacy-dialog';

import { AppConfig } from '../../../../app.config';

/**
 * Shows instructions about how to use the PIN matrix.
 */
@Component({
  selector: 'app-hw-added-dialog',
  templateUrl: './hw-pin-help-dialog.component.html',
  styleUrls: ['./hw-pin-help-dialog.component.scss'],
})
export class HwPinHelpDialogComponent {
  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog): MatDialogRef<HwPinHelpDialogComponent, any> {
    const config = new MatDialogConfig();
    config.autoFocus = false;
    config.width = AppConfig.smallModalWidth;

    return dialog.open(HwPinHelpDialogComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<HwPinHelpDialogComponent>,
  ) { }
}
