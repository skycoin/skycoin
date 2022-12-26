import { Component, OnInit } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';

import { AppConfig } from '../../../../app.config';
import { parseRequestLink } from '../../../../utils/general-utils';

/**
 * Modal window for entering a link with info for populating the send form. If the user confirms
 * the operation, the modal window is closed and the link is returned in the "afterClosed" event.
 */
@Component({
  selector: 'app-enter-link',
  templateUrl: './enter-link.component.html',
  styleUrls: ['./enter-link.component.scss'],
})
export class EnterLinkComponent implements OnInit {
  form: UntypedFormGroup;

  // Vars with the validation error messages.
  inputErrorMsg = '';

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog): MatDialogRef<EnterLinkComponent, any> {
    const config = new MatDialogConfig();
    config.autoFocus = true;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(EnterLinkComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<EnterLinkComponent>,
    private formBuilder: UntypedFormBuilder,
  ) { }

  ngOnInit() {
    this.form = this.formBuilder.group({
      link: [''],
    });

    this.form.setValidators(this.validateForm.bind(this));
  }

  closePopup() {
    this.dialogRef.close();
  }

  /**
   * If the form is valid, closes the modal window and returns the link the user entered.
   */
  process() {
    if (!this.form.valid) {
      return;
    }

    this.dialogRef.close(this.form.get('link').value);
  }

  /**
   * Validates the form and updates the vars with the validation errors.
   */
  validateForm() {
    this.inputErrorMsg = '';

    let valid = true;

    // The link must be valid.
    if (!parseRequestLink(this.form.get('link').value)) {
      valid = false;
      if (this.form.get('link').touched) {
        this.inputErrorMsg = 'send.fill-with-link.link-error-info';
      }
    }

    return valid ? null : { Invalid: true };
  }
}
