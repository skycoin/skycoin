import { Component } from '@angular/core';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';

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
    config.width = '450px';

    return dialog.open(HwPinHelpDialogComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<HwPinHelpDialogComponent>,
  ) { }
}
