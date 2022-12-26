import { Component } from '@angular/core';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';

import { AppConfig } from '../../../../../../app.config';

export enum DestinationTools {
  Bulk = 'bulk',
  Link = 'link',
}

/**
 * Modal window for selecting one of the tools for adding the destination info to the form. If
 * the user selects an option, the modal window is closed and the option from the
 * DestinationTools enum is returned in the "afterClosed" event.
 */
@Component({
  selector: 'app-destination-tools',
  templateUrl: './destination-tools.component.html',
  styleUrls: ['./destination-tools.component.scss'],
})
export class DestinationToolsComponent {
  destinationTools = DestinationTools;

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog): MatDialogRef<DestinationToolsComponent, DestinationTools> {
    const config = new MatDialogConfig();
    config.autoFocus = false;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(DestinationToolsComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<DestinationToolsComponent>,
  ) { }

  closePopup() {
    this.dialogRef.close();
  }

  /**
   * Closes the modal window and returns the selected option.
   */
  select(value: DestinationTools) {
    this.dialogRef.close(value);
  }
}
