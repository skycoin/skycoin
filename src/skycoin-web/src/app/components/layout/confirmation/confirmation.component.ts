import { Component, Input, Inject, ChangeDetectionStrategy } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ConfirmationData } from '../../../app.datatypes';

@Component({
    selector: 'app-confirmation',
    templateUrl: './confirmation.component.html',
    styleUrls: ['./confirmation.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class ConfirmationComponent {
  accepted = false;
  disableDismiss = false;

  constructor(
    public dialogRef: MatDialogRef<ConfirmationComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ConfirmationData
  ) {
    this.disableDismiss = !!data.disableDismiss;
  }

  closeModal(isConfirmed: boolean) {
    this.dialogRef.close(isConfirmed);
  }

  setAccept(event: any) {
    this.accepted = event.checked ? true : false;
  }
}
