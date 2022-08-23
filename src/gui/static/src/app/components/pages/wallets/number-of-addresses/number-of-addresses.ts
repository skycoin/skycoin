import { Component, OnInit, ViewChild, OnDestroy, Output, EventEmitter } from '@angular/core';
import { UntypedFormGroup, UntypedFormControl } from '@angular/forms';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';

import { ButtonComponent } from '../../../layout/button/button.component';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { AppConfig } from '../../../../app.config';

/**
 * Data sent when the user tries to add addresses with NumberOfAddressesComponent.
 */
export interface NumberOfAddressesEventData {
  /**
   * How many addresses the user wants to add.
   */
  howManyAddresses: number;
  /**
   * Callback function that must be used for informing the NumberOfAddressesComponent
   * instance that the preparations for adding the addresses have been finished.
   * @param close If the modal window must be closed.
   * @param endedWithError If the preparations ended because of an error.
   */
  callback(close: boolean, endedWithError?: boolean): void;
}

/**
 * Modal window for entering how many addresses to add to a wallet. It does not add the
 * addresses, but emits an event for informing the caller when the addresses must be created.
 */
@Component({
  selector: 'app-number-of-addresses',
  templateUrl: './number-of-addresses.html',
  styleUrls: ['./number-of-addresses.scss'],
})
export class NumberOfAddressesComponent implements OnInit, OnDestroy {
  // Confirmation button.
  @ViewChild('button') button: ButtonComponent;
  form: UntypedFormGroup;
  // Emits when the user request the addresses to be added.
  @Output() createRequested = new EventEmitter<NumberOfAddressesEventData>();

  // Vars with the validation error messages.
  inputErrorMsg = '';

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog): MatDialogRef<NumberOfAddressesComponent, any> {
    const config = new MatDialogConfig();
    config.autoFocus = true;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(NumberOfAddressesComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<NumberOfAddressesComponent>,
    private msgBarService: MsgBarService,
  ) {}

  ngOnInit() {
    this.form = new UntypedFormGroup({});
    this.form.addControl('quantity', new UntypedFormControl(1));

    this.form.setValidators(this.validateForm.bind(this));
  }

  ngOnDestroy() {
    this.msgBarService.hide();
    this.createRequested.complete();
  }

  closePopup() {
    this.dialogRef.close();
  }

  continue() {
    if (this.button.isLoading()) {
      return;
    }

    this.msgBarService.hide();
    this.button.setLoading();

    this.createRequested.emit({
      howManyAddresses: this.form.value.quantity,
      callback: (close, endedWithError = false) => {
        this.button.resetState();
        if (!endedWithError) {
          if (close) {
            this.closePopup();
          }
        } else {
          this.msgBarService.showError('wallet.add-addresses.error');
        }
      },
    });
  }

  /**
   * Validates the form and updates the vars with the validation errors.
   */
  validateForm() {
    this.inputErrorMsg = '';

    let valid = true;

    // The number must be an integer from 1 to 100.
    const value = this.form.get('quantity').value as number;
    if (!value || value < 1 || value > 100 || value !== Math.round(value)) {
      valid = false;
      if (this.form.get('quantity').touched) {
        this.inputErrorMsg = 'wallet.add-addresses.quantity-error-info';
      }
    }

    return valid ? null : { Invalid: true };
  }
}
