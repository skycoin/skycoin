import { Component } from '@angular/core';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { AppConfig } from '../../../../../app.config';

export enum AddressOptions {
  new = 'new',
  scan = 'scan',
}

@Component({
  selector: 'app-address-options',
  templateUrl: './address-options.component.html',
  styleUrls: ['./address-options.component.scss'],
})
export class AddressOptionsComponent {

  addressOptions = AddressOptions;

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
