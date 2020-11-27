import { Component, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';

import { ButtonComponent } from '../../../layout/button/button.component';

/**
 * States in which OfflineDialogsBaseComponent and derived classes can be.
 */
export enum OfflineDialogsStates {
  Loading = 'Loading',
  ErrorLoading = 'ErrorLoading',
  ShowingForm = 'ShowingForm',
}

/**
 * Data for an element of the dropdown control.
 */
export interface OfflineDialogsDropdownElement {
  /**
   * Text to show on the UI.
   */
  name: string;
  /**
   * Actual value.
   */
  value: any;
}

/**
 * Base component for the modal windows used for operations related to creating transactions
 * with an offline wallet. It provides the UI and some functionality.
 */
@Component({
  template: '',
})
export class OfflineDialogsBaseComponent {
  @ViewChild('okButton') okButton: ButtonComponent;
  // Allows to deactivate the form while the component is busy.
  working = false;
  form: FormGroup;
  currentState: OfflineDialogsStates = OfflineDialogsStates.Loading;
  states = OfflineDialogsStates;
  validateForm = false;

  // Basic info for the UI. The values must be set by the subclasses.
  title = '';
  text = '';
  dropdownLabel = '';
  defaultDropdownText = '';
  inputLabel = '';
  // If it has a value, the textarea will be in read only mode.
  contents = '';
  cancelButtonText = '';
  okButtonText = '';

  // Vars with the validation error messages.
  dropdownErrorMsg = '';
  inputErrorMsg = '';

  // If not set, the dropdown control is not shown.
  dropdownElements: OfflineDialogsDropdownElement[];

  constructor(
    _formBuilder: FormBuilder,
  ) {
    this.form = _formBuilder.group({
      dropdown: [''],
      input: [''],
    });

    this.form.setValidators(this.validate.bind(this));
  }

  /**
   * Validates the form and updates the vars with the validation errors.
   */
  private validate() {
    this.dropdownErrorMsg = '';
    this.inputErrorMsg = '';

    let valid = true;

    if (!this.form.get('dropdown').value) {
      valid = false;
      if (this.form.get('dropdown').touched) {
        this.dropdownErrorMsg = 'offline-transactions.wallet-error-info';
      }
    }

    const inputValue = this.form.get('input').value as string;
    if (!inputValue || inputValue.length < 300 || !/^[0-9a-fA-F]+$/.test(inputValue)) {
      valid = false;
      if (this.form.get('input').touched) {
        this.inputErrorMsg = 'offline-transactions.tx-error-info';
      }
    }

    return valid ? null : { Invalid: true };
  }
}
