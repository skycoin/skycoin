import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';

import { MsgBarService } from '../../../services/msg-bar.service';
import { AppConfig } from '../../../app.config';
import { Destination } from '../../pages/send-skycoin/form-parts/form-destination/form-destination.component';

/**
 * Modal window which allows the user to enter multiple destinations with just a text
 * string, when creating a transaction. The format of the string is just each destination
 * in a new line, each one with the address, coins and hours (optional), separated by commas,
 * in that order. If a destination has hours, all destinations must have hours. Destinations
 * can have invalid addresses, coins and hours, the only thing this component checks is if
 * there is a string value for each property of the destinations. In the "afterClosed" event
 * it returns an array with all the destinations as instances of BasicDestinationData, or
 * null, if the operation was cancelled.
 */
@Component({
  selector: 'app-multiple-destinations-dialog',
  templateUrl: './multiple-destinations-dialog.component.html',
  styleUrls: ['./multiple-destinations-dialog.component.scss'],
})
export class MultipleDestinationsDialogComponent implements OnInit, OnDestroy {
  form: FormGroup;

  // Vars with the validation error messages.
  inputErrorMsg = '';

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog, content: string): MatDialogRef<MultipleDestinationsDialogComponent, any> {
    const config = new MatDialogConfig();
    config.data = content;
    config.autoFocus = true;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(MultipleDestinationsDialogComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<MultipleDestinationsDialogComponent>,
    @Inject(MAT_DIALOG_DATA) private data: string,
    private formBuilder: FormBuilder,
    private msgBarService: MsgBarService,
  ) { }

  ngOnInit() {
    this.form = this.formBuilder.group({
      data: [this.data],
    });

    this.form.setValidators(this.validateForm.bind(this));
  }

  ngOnDestroy() {
    this.msgBarService.hide();
  }

  /**
   * Process the text entered by the user. If there are no errors, it closes the modal window and
   * returns the list will all the destinations.
   */
  processData() {
    try {
      // Check if empty.
      if ((this.form.value.data as string).trim().length === 0) {
        this.msgBarService.showError('send.bulk-send.no-data-error');

        return;
      }

      // Get all destinations.
      let entries = (this.form.value.data as string).split(/\r?\n/);
      if (!entries || entries.length === 0) {
        this.msgBarService.showError('send.bulk-send.no-data-error');

        return;
      }

      // Remove empty lines.
      entries = entries.filter(entry => entry.trim().length > 0);

      // Destination must have 2 or 3 (if including hours) parts.
      const firstElementParts = entries[0].split(',').length;
      if (firstElementParts !== 2 && firstElementParts !== 3) {
        this.msgBarService.showError('send.bulk-send.invalid-data-error');

        return;
      }

      // Separate all the parts each destination.
      const processedEntries: Destination[] = [];
      let consistentNumberOfParts = true;
      entries.forEach(entry => {
        const entryDataParts = entry.split(',');
        const data: Destination = {
          address: entryDataParts[0].trim(),
          coins: entryDataParts[1].trim(),
          originalAmount: null,
        };
        data.hours = entryDataParts.length === 3 ? entryDataParts[2].trim() : undefined;
        processedEntries.push(data);

        // Check if this has the same number of parts as the first one.
        if (entryDataParts.length !== firstElementParts) {
          consistentNumberOfParts = false;
        }
      });

      // Do not allow a mix of some destinations with hours and others without them.
      if (!consistentNumberOfParts) {
        this.msgBarService.showError('send.bulk-send.inconsistent-data-error');

        return;
      }

      this.dialogRef.close(processedEntries);
    } catch (e) {
      this.msgBarService.showError('send.bulk-send.invalid-data-error');
    }
  }

  /**
   * Validates the form and updates the vars with the validation errors.
   */
  validateForm() {
    this.inputErrorMsg = '';

    let valid = true;

    if (!this.form.get('data').value || !this.form.get('data').value.trim()) {
      valid = false;
      if (this.form.get('data').touched) {
        this.inputErrorMsg = 'send.bulk-send.data-error-info';
      }
    }

    return valid ? null : { Invalid: true };
  }
}
