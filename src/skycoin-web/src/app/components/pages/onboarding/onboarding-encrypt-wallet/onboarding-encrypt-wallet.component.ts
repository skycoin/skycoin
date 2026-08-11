import { Component, OnInit, ViewChild, Output, EventEmitter, ChangeDetectionStrategy } from '@angular/core';
import { UntypedFormBuilder, UntypedFormControl, UntypedFormGroup, Validators } from '@angular/forms';

import { ButtonComponent } from '../../../layout/button/button.component';

@Component({
    selector: 'app-onboarding-encrypt-wallet',
    templateUrl: './onboarding-encrypt-wallet.component.html',
    styleUrls: ['./onboarding-encrypt-wallet.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class OnboardingEncryptWalletComponent implements OnInit {
  @ViewChild('button') button!: ButtonComponent;
  @Output() onPasswordCreated = new EventEmitter<string|null>();
  @Output() onBack = new EventEmitter();

  form!: UntypedFormGroup;

  constructor(
    private formBuilder: UntypedFormBuilder,
  ) { }

  get isWorking(): boolean {
    // The template binds this through [ngClass], so it is read during the first
    // change detection pass — before a non-static @ViewChild is resolved.
    return !!this.button && this.button.isLoading();
  }

  ngOnInit() {
    this.initEncryptForm();
  }

  initEncryptForm() {
    this.form = this.formBuilder.group({
        password: new UntypedFormControl('', Validators.compose([Validators.required, Validators.minLength(2)])),
        confirm: new UntypedFormControl('',
          Validators.compose([
            Validators.required,
            Validators.minLength(2),
          ]),
        ),
      },
      {
        validator: this.passwordMatchValidator.bind(this),
      });

    this.form.disable();
  }

  setEncrypt(event: any) {
    event.checked ? this.form.enable() : this.form.disable();
  }

  emitCreatedPassword() {
    if ((this.form.enabled && !this.form.valid) || this.button.isLoading()) {
      return;
    }

    this.button.setLoading();
    this.onPasswordCreated.emit(this.form.enabled ? this.form.get('password')!.value : null);
  }

  emitBack() {
    this.onBack.emit();
  }

  private passwordMatchValidator(g: UntypedFormGroup) {
    return g.get('password')!.value === g.get('confirm')!.value
      ? null : { mismatch: true };
  }
}
