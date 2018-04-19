import { Component, Inject, OnDestroy } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material';

@Component({
  selector: 'app-seed-modal',
  templateUrl: './seed-modal.component.html',
  styleUrls: ['./seed-modal.component.css']
})
export class SeedModalComponent implements OnDestroy {

  constructor(
    public dialogRef: MatDialogRef<SeedModalComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {}

  ngOnDestroy() {
    this.data.seed = null;
  }

}
