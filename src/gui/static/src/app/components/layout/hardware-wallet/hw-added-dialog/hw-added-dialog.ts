import { Component, OnDestroy } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-hw-added-dialog',
  templateUrl: './hw-added-dialog.html',
  styleUrls: ['./hw-added-dialog.scss'],
})
export class HwAddedDialogComponent {

  constructor(
    public dialogRef: MatDialogRef<HwAddedDialogComponent>,
  ) { }
}
