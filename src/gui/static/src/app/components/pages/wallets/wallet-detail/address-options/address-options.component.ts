import { Component } from '@angular/core';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';

import { AppConfig } from '../../../../../app.config';

/**
 * Options available in AddressOptionsComponent.
 */
export enum AddressOptions {
  new = 'new',
  scan = 'scan',
}

/**
 * Modal window with options for making operations with the addresses of a wallet. If the
 * user selects an option, the modal window is closed and a value from the AddressOptions enum
 * is returned in the "afterClosed" event.
 */
@Component({
  selector: 'app-address-options',
  templateUrl: './address-options.component.html',
  styleUrls: ['./address-options.component.scss'],
})
export class AddressOptionsComponent {
  addressOptions = AddressOptions;

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog): MatDialogRef<AddressOptionsComponent, any> {
    const config = new MatDialogConfig();
    config.autoFocus = false;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(AddressOptionsComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<AddressOptionsComponent>,
  ) { }

  closePopup() {
    this.dialogRef.close();
  }

  select(value: AddressOptions) {
    this.dialogRef.close(value);
  }
}
