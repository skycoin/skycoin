import { Component, ViewChild } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { ButtonComponent } from '../../../layout/button/button.component';

export enum OfflineDialogsStates {
  Loading,
  ErrorLoading,
  ShowingForm,
}

export interface OfflineDialogsDropdownElement {
  name: string;
  value: any;
}

@Component({
  template: '',
})
export class OfflineDialogsBaseComponent {
  @ViewChild('cancelButton', { static: false }) cancelButton: ButtonComponent;
  @ViewChild('okButton', { static: false }) okButton: ButtonComponent;
  working = false;
  form: FormGroup;
  currentState: OfflineDialogsStates = OfflineDialogsStates.Loading;
  states = OfflineDialogsStates;
  validateForm = false;

  title = '';
  text = '';
  dropdownLabel = '';
  defaultDropdownText = '';
  inputLabel = '';
  contents = '';
  cancelButtonText = '';
  okButtonText = '';

  dropdownElements: OfflineDialogsDropdownElement[];

  constructor(
    _formBuilder: FormBuilder,
  ) {
    this.form = _formBuilder.group({
      dropdown: ['', Validators.required],
      input: ['', Validators.compose([
        Validators.required,
        Validators.minLength(300),
        Validators.pattern('^[0-9a-fA-F]+$'),
      ])],
    });
  }
}
