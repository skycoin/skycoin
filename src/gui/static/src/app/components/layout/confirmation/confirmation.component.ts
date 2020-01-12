import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { ConfirmationData } from '../../../app.datatypes';
import { AppConfig } from '../../../app.config';

@Component({
  selector: 'app-confirmation',
  templateUrl: './confirmation.component.html',
  styleUrls: ['./confirmation.component.scss'],
})
export class ConfirmationComponent {
  accepted = false;
  disableDismiss = false;

  public static openDialog(dialog: MatDialog, confirmationData: ConfirmationData): MatDialogRef<ConfirmationComponent, any> {
    const config = new MatDialogConfig();
    config.data = confirmationData;
    config.autoFocus = false;
    config.width = '450px';

    return dialog.open(ConfirmationComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<ConfirmationComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ConfirmationData,
  ) {
    this.disableDismiss = !!data.disableDismiss;
  }

  closeModal(isConfirmed: boolean) {
    this.dialogRef.close(isConfirmed);
  }

  setAccept(event) {
    this.accepted = event.checked ? true : false;
  }
}
