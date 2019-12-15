import { Component, EventEmitter, OnInit, Output, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
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

  constructor(
    private formBuilder: FormBuilder,
  ) { }

  ngOnInit() {
    this.initEncryptForm();
  }

  initEncryptForm() {
    this.form = this.formBuilder.group({
        password: new FormControl('', Validators.compose([Validators.required, Validators.minLength(2)])),
        confirm: new FormControl('',
          Validators.compose([
            Validators.required,
            Validators.minLength(2),
          ]),
        ),
      },
      {
        validator: this.passwordMatchValidator.bind(this),
      });
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

  private passwordMatchValidator(g: FormGroup) {
    return g.get('password').value === g.get('confirm').value
      ? null : { mismatch: true };
  }
}
