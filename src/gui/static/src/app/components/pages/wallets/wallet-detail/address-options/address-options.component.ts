import { Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

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
