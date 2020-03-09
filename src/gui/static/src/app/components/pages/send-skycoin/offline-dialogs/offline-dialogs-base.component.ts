import { Component, ViewChild } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';

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
  @ViewChild('okButton', { static: false }) okButton: ButtonComponent;
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

  // If not set, the dropdown control is not shown.
  dropdownElements: OfflineDialogsDropdownElement[];

  constructor(
    _formBuilder: FormBuilder,
  ) {
    this.form = _formBuilder.group({
      dropdown: ['', Validators.required],
      input: ['', Validators.compose([
        Validators.required,
        // Mim for a transaction.
        Validators.minLength(300),
        // Hex characters only.
        Validators.pattern('^[0-9a-fA-F]+$'),
      ])],
    });
  }
}
