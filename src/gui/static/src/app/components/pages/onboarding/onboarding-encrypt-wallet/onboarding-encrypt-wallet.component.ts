import { Component, EventEmitter, OnInit, Output, ViewChild, ChangeDetectorRef, OnDestroy, Input } from '@angular/core';
import { FormBuilder, FormControl, FormGroup } from '@angular/forms';

import { ButtonComponent } from '../../../layout/button/button.component';

/**
 * Shows the second step of the wizard, which allows the user to set the password.
 */
@Component({
  selector: 'app-onboarding-encrypt-wallet',
  templateUrl: './onboarding-encrypt-wallet.component.html',
  styleUrls: ['./onboarding-encrypt-wallet.component.scss'],
})
export class OnboardingEncryptWalletComponent implements OnInit, OnDestroy {
  @ViewChild('button') button: ButtonComponent;
  // Emits when the user presses the button for going to the next step of the wizard, after
  // filling the form. Includes the password entered by the user, or null, if the user
  // selected not to encrypt the wallet.
  @Output() passwordCreated = new EventEmitter<string|null>();
  // Emits when the user presses the button for going back to the previous step of the wizard.
  @Output() goBack = new EventEmitter();
  form: FormGroup;

  // If creating a temporal wallet.
  creatingTemporal_ = false;
  @Input() set creatingTemporal(val: boolean) {
    this.creatingTemporal_ = val;
    setTimeout(() => {
      if (this.form) {
        this.form.updateValueAndValidity();
      }
    });
  }
  get creatingTemporal(): boolean {
    return this.creatingTemporal_;
  }

  // Vars with the validation error messages.
  password1ErrorMsg = '';
  password2ErrorMsg = '';

  constructor(
    private formBuilder: FormBuilder,
    private changeDetector: ChangeDetectorRef,
  ) { }

  ngOnInit() {
    this.form = this.formBuilder.group({
        password: new FormControl(''),
        confirm: new FormControl(''),
      },
    );

    this.form.setValidators(this.validateForm.bind(this));
  }

  ngOnDestroy() {
    this.passwordCreated.complete();
    this.goBack.complete();
  }

  // Called after pressing the checkbox for selecting if the wallet must be encrypted with
  // a password or not.
  setEncrypt(event) {
    if (event.checked) {
      this.form.enable();
    } else {
      this.form.disable();
    }
  }

  // Emits an event for going to the next step of the wizard.
  emitCreatedPassword() {
    if ((!this.creatingTemporal && (this.form.enabled && !this.form.valid)) || this.button.isLoading()) {
      return;
    }

    this.button.setLoading();

    this.passwordCreated.emit(!this.creatingTemporal && this.form.enabled ? this.form.get('password').value : null);

    this.changeDetector.detectChanges();
  }

  // Emits an event for going to the previous step of the wizard.
  emitBack() {
    this.goBack.emit();
  }

  // Returns the continue button to its initial state.
  resetButton() {
    this.button.resetState();
  }

  // Allows to know if the app is processing and the form must be shown disabled.
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
