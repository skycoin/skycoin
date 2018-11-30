import { Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-hw-added-dialog',
  templateUrl: './hw-pin-help.component.html',
  styleUrls: ['./hw-pin-help.component.scss'],
})
export class HwPinHelpComponent {

  constructor(
    public dialogRef: MatDialogRef<HwPinHelpComponent>,
  ) { }
}
