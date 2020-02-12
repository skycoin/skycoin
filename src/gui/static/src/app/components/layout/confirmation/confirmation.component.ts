import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';

export enum DefaultConfirmationButtons {
  YesNo = 'YesNo',
  ContinueCancel = 'ContinueCancel',
  Close = 'Close',
}

export interface ConfirmationParams {
  text: string;
  headerText?: string;
  checkboxText?: string;
  defaultButtons?: DefaultConfirmationButtons;
  confirmButtonText?: string;
  cancelButtonText?: string;
  redTitle?: boolean;
  disableDismiss?: boolean;
  linkText?: string;
  linkFunction?(): void;
}

@Component({
  selector: 'app-confirmation',
  templateUrl: './confirmation.component.html',
  styleUrls: ['./confirmation.component.scss'],
})
export class ConfirmationComponent {
  accepted = false;
  disableDismiss = false;

  public static openDialog(dialog: MatDialog, confirmationParams: ConfirmationParams): MatDialogRef<ConfirmationComponent, any> {
    const config = new MatDialogConfig();
    config.data = confirmationParams;
    config.autoFocus = false;
    config.width = '450px';

    return dialog.open(ConfirmationComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<ConfirmationComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ConfirmationParams,
  ) {
    if (!data.headerText) {
      data.headerText = 'confirmation.default-title';
    }

    if (data.defaultButtons) {
      if (data.defaultButtons === DefaultConfirmationButtons.Close) {
        data.confirmButtonText = 'common.close-button';
      }

      if (data.defaultButtons === DefaultConfirmationButtons.YesNo) {
        data.confirmButtonText = 'confirmation.yes-button';
        data.cancelButtonText = 'confirmation.no-button';
      }

      if (data.defaultButtons === DefaultConfirmationButtons.ContinueCancel) {
        data.confirmButtonText = 'common.continue-button';
        data.cancelButtonText = 'common.cancel-button';
      }
    }

    this.disableDismiss = !!data.disableDismiss;
  }

  closeModal(isConfirmed: boolean) {
    this.dialogRef.close(isConfirmed);
  }

  setAccept(event) {
    this.accepted = event.checked ? true : false;
  }
}
