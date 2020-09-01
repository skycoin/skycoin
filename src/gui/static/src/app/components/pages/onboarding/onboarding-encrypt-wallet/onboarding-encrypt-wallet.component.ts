import { Component, EventEmitter, OnInit, Output, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup } from '@angular/forms';

import { ButtonComponent } from '../../../layout/button/button.component';

@Component({
  selector: 'app-onboarding-encrypt-wallet',
  templateUrl: './onboarding-encrypt-wallet.component.html',
  styleUrls: ['./onboarding-encrypt-wallet.component.scss'],
})
export class OnboardingEncryptWalletComponent implements OnInit {
  @ViewChild('button', { static: false }) button: ButtonComponent;
  @Output() onPasswordCreated = new EventEmitter<string|null>();
  @Output() onBack = new EventEmitter();
  form: FormGroup;

  // Vars with the validation error messages.
  password1ErrorMsg = '';
  password2ErrorMsg = '';

  constructor(
    private formBuilder: FormBuilder,
  ) { }

  ngOnInit() {
    this.initEncryptForm();
  }

  initEncryptForm() {
    this.form = this.formBuilder.group({
        password: new FormControl(''),
        confirm: new FormControl(''),
      },
    );

    this.form.setValidators(this.validateForm.bind(this));
  }

  setEncrypt(event) {
    event.checked ? this.form.enable() : this.form.disable();
  }

  emitCreatedPassword() {
    if ((this.form.enabled && !this.form.valid) || this.button.isLoading()) {
      return;
    }

    this.button.setLoading();

    this.onPasswordCreated.emit(this.form.enabled ? this.form.get('password').value : null);
  }

  emitBack() {
    this.onBack.emit();
  }

  resetButton() {
    this.button.resetState();
  }

  get isWorking() {
    return this.button ? this.button.isLoading() : false;
  }

  /**
   * Validates the form and updates the vars with the validation errors.
   */
  validateForm() {
    this.password1ErrorMsg = '';
    this.password2ErrorMsg = '';

    let valid = true;

    if (!this.form.get('password').value) {
      valid = false;
      if (this.form.get('password').touched) {
        this.password1ErrorMsg = 'password.password-error-info';
      }
    }

    if (!this.form.get('confirm').value) {
      valid = false;
      if (this.form.get('confirm').touched) {
        this.password2ErrorMsg = 'password.password-error-info';
      }
    }

    // If both password fields have a value, check if the 2 passwords entered by the user
    // are equal.
    if (valid && this.form.get('password').value !== this.form.get('confirm').value) {
      valid = false;
      this.password2ErrorMsg = 'password.confirm-error-info';
    }

    return valid ? null : { Invalid: true };
  }
}
