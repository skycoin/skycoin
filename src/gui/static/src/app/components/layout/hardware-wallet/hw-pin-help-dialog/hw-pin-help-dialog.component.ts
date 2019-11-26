import { Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-hw-added-dialog',
  templateUrl: './hw-pin-help-dialog.component.html',
  styleUrls: ['./hw-pin-help-dialog.component.scss'],
})
export class HwPinHelpDialogComponent {

  constructor(
    public dialogRef: MatDialogRef<HwPinHelpDialogComponent>,
  ) { }
}
