import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';

import { AppConfig } from '../../../../../app.config';

/**
 * Modal window for displaying the seed of a wallet, for making a backup.
 */
@Component({
  selector: 'app-seed-modal',
  templateUrl: './seed-modal.component.html',
  styleUrls: ['./seed-modal.component.scss'],
})
export class SeedModalComponent {
  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog, seed: string): MatDialogRef<SeedModalComponent, any> {
    const config = new MatDialogConfig();
    config.data = seed;
    config.autoFocus = true;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(SeedModalComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: any,
    public dialogRef: MatDialogRef<SeedModalComponent>,
  ) {}
}
