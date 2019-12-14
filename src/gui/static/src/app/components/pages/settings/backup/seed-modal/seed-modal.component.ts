import { Component, Inject, OnDestroy } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-seed-modal',
  templateUrl: './seed-modal.component.html',
  styleUrls: ['./seed-modal.component.css'],
})
export class SeedModalComponent implements OnDestroy {
  constructor(
    @Inject(MAT_DIALOG_DATA) public data: any,
    public dialogRef: MatDialogRef<SeedModalComponent>,
  ) {}

  ngOnDestroy() {
    this.data.seed = null;
  }
}
