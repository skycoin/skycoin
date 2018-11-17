import { Component, Input, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { ConfirmationData } from '../../../app.datatypes';

@Component({
  selector: 'app-confirmation',
  templateUrl: './confirmation.component.html',
  styleUrls: ['./confirmation.component.scss'],
})
export class ConfirmationComponent {
  accepted = false;

  constructor(
    public dialogRef: MatDialogRef<ConfirmationComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ConfirmationData,
  ) { }

  closeModal(isConfirmed: boolean) {
    this.dialogRef.close(isConfirmed);
  }

  setAccept(event) {
    this.accepted = event.checked ? true : false;
  }
}
