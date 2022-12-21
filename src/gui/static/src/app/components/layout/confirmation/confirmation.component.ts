import { Component, Inject } from '@angular/core';
import { MatLegacyDialogRef as MatDialogRef, MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA, MatLegacyDialog as MatDialog, MatLegacyDialogConfig as MatDialogConfig } from '@angular/material/legacy-dialog';

import { AppConfig } from '../../../app.config';

/**
 * Predefined button combinations for ConfirmationComponent.
 */
export enum DefaultConfirmationButtons {
  YesNo = 'YesNo',
  ContinueCancel = 'ContinueCancel',
  Close = 'Close',
}

/**
 * Settings for ConfirmationComponent.
 */
export interface ConfirmationParams {
  /**
   * Main text the modal window will show.
   */
  text: string;
  /**
   * Custom title for the modal window.
   */
  headerText?: string;
  /**
   * If set, the modal window will show a check box the user will have to acept before being able
   * to press the confirmation button.
   */
  checkboxText?: string;
  /**
   * Default buttons combination the modal window will show. If set, you will not be able to set
   * custom texts for the buttons.
   */
  defaultButtons?: DefaultConfirmationButtons;
  /**
   * Text for the confirmation button. If unset, the button is not shown.
   */
  confirmButtonText?: string;
  /**
   * Text for the cancel button. If unset, the button is not shown.
   */
  cancelButtonText?: string;
  /**
   * If true, the title of the modal window will be red.
   */
  redTitle?: boolean;
  /**
   * If true, the user will be forced to close the modal window using the confirm or cancel button.
   */
  disableDismiss?: boolean;
  /**
   * If set, the modal window will show the link under the main text.
   */
  linkText?: string;
  /**
   * Function that will be called when the user clicks the link.
   */
  linkFunction?(): void;
}

/**
 * Modal window used to ask the user to confirm an operation. If the user confirms the operation,
 * the modal window is closed and "true" is returned in the "afterClosed" event.
 */
@Component({
  selector: 'app-confirmation',
  templateUrl: './confirmation.component.html',
  styleUrls: ['./confirmation.component.scss'],
})
export class ConfirmationComponent {
  // If the user checked the checkbox.
  accepted = false;
  disableDismiss = false;

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog, confirmationParams: ConfirmationParams): MatDialogRef<ConfirmationComponent, any> {
    const config = new MatDialogConfig();
    config.data = confirmationParams;
    config.autoFocus = false;
    config.width = AppConfig.smallModalWidth;

    return dialog.open(ConfirmationComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<ConfirmationComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ConfirmationParams,
  ) {
    // Default header.
    if (!data.headerText) {
      data.headerText = 'confirmation.default-title';
    }

    // Use a default buttons combination, if requested.
    if (data.defaultButtons) {
      data.confirmButtonText = null;
      data.cancelButtonText = null;

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

  // Used by the checkbox.
  setAccept(event) {
    this.accepted = event.checked ? true : false;
  }
}
