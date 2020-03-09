import { Component, EventEmitter, OnInit, Output, ViewChild, ChangeDetectorRef } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';

import { ButtonComponent } from '../../../layout/button/button.component';

/**
 * Shows the second step of the wizard, which allows the user to set the password.
 */
@Component({
  selector: 'app-onboarding-encrypt-wallet',
  templateUrl: './onboarding-encrypt-wallet.component.html',
  styleUrls: ['./onboarding-encrypt-wallet.component.scss'],
})
export class OnboardingEncryptWalletComponent implements OnInit {
  @ViewChild('button', { static: false }) button: ButtonComponent;
  // Emits when the user presses the button for going to the next step of the wizard, after
  // filling the form. Includes the password entered by the user, or null, if the user
  // selected not to encrypt the wallet.
  @Output() onPasswordCreated = new EventEmitter<string|null>();
  // Emits when the user presses the button for going back to the previous step of the wizard.
  @Output() onBack = new EventEmitter();
  form: FormGroup;

  constructor(
    private formBuilder: FormBuilder,
    private changeDetector: ChangeDetectorRef,
  ) { }

  ngOnInit() {
    this.form = this.formBuilder.group({
        password: new FormControl('', Validators.compose([Validators.required])),
        confirm: new FormControl('', Validators.compose([Validators.required])),
      },
      {
        validator: this.passwordMatchValidator.bind(this),
      });
  }

  // Called after pressing the checkbox for selecting if the wallet must be encrypted with
  // a password or not.
  setEncrypt(event) {
    event.checked ? this.form.enable() : this.form.disable();
  }

  // Emits an event for going to the next step of the wizard.
  emitCreatedPassword() {
    if ((this.form.enabled && !this.form.valid) || this.button.isLoading()) {
      return;
    }

    this.button.setLoading();

    this.onPasswordCreated.emit(this.form.enabled ? this.form.get('password').value : null);

    this.changeDetector.detectChanges();
  }

  // Emits an event for going to the previous step of the wizard.
  emitBack() {
    this.onBack.emit();
  }

  resetButton() {
    this.button.resetState();
  }

  // Allows to know if the app is processing and the form must be shown disabled.
  get isWorking() {
    return this.button ? this.button.isLoading() : false;
  }

  // Checks if both password fields match.
  private passwordMatchValidator(g: FormGroup) {
    return g.get('password').value === g.get('confirm').value ? null : { mismatch: true };
  }
}
