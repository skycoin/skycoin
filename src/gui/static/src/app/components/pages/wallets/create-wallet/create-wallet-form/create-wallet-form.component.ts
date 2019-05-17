import { Component, OnInit, OnDestroy, Input } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ISubscription } from 'rxjs/Subscription';
import { ApiService } from '../../../../../services/api.service';
import { Subject } from 'rxjs/Subject';
import 'rxjs/add/operator/switchMap';
import { Observable } from 'rxjs/Observable';

export class FormData {
  label: string;
  seed: string;
  password: string;
}

@Component({
  selector: 'app-create-wallet-form',
  templateUrl: './create-wallet-form.component.html',
  styleUrls: ['./create-wallet-form.component.scss'],
})
export class CreateWalletFormComponent implements OnInit, OnDestroy {
  @Input() create: boolean;
  @Input() whiteText: boolean;
  @Input() onboarding: boolean;

  form: FormGroup;
  normalSeed = true;
  customSeedAccepted = false;
  encrypt = true;

  private seed: Subject<string> = new Subject<string>();
  private statusSubscription: ISubscription;
  private seedValiditySubscription: ISubscription;

  constructor(
    private apiService: ApiService,
  ) { }

  ngOnInit() {
    if (!this.onboarding) {
      this.initForm();
    } else {
      this.initForm(false, null);
    }
  }

  ngOnDestroy() {
    this.statusSubscription.unsubscribe();
    this.seedValiditySubscription.unsubscribe();
  }

  get isValid(): boolean {
    return this.form.valid && (this.normalSeed || this.customSeedAccepted);
  }

  onCustomSeedAcceptance(event) {
    this.customSeedAccepted = event.checked;
  }

  setEncrypt(event) {
    this.encrypt = event.checked;
    this.form.updateValueAndValidity();
  }

  getData(): FormData {
    return {
      label: this.form.value.label,
      seed: this.form.value.seed,
      password: !this.onboarding && this.encrypt ? this.form.value.password : null,
    };
  }

  initForm(create: boolean = null, data: Object = null) {
    create = create !== null ? create : this.create;

    const validators = [];
    if (create) {
      validators.push(this.seedMatchValidator.bind(this));
    }
    if (!this.onboarding) {
      validators.push(this.validatePasswords.bind(this));
    }

    this.form = new FormGroup({}, validators);
    this.form.addControl('label', new FormControl(data ? data['label'] : '', [Validators.required]));
    this.form.addControl('seed', new FormControl(data ? data['seed'] : '', [Validators.required]));
    this.form.addControl('confirm_seed', new FormControl(data ? data['seed'] : ''));
    this.form.addControl('password', new FormControl());
    this.form.addControl('confirm_password', new FormControl());

    if (create && !data) {
      this.generateSeed(128);
    }

    if (data) {
      setTimeout(() => { this.seed.next(data['seed']); });
      this.customSeedAccepted = true;
    }

    if (this.statusSubscription && !this.statusSubscription.closed) {
      this.statusSubscription.unsubscribe();
    }
    this.statusSubscription = this.form.statusChanges.subscribe(() => {
      this.customSeedAccepted = false;
      this.seed.next(this.form.get('seed').value);
    });

    this.subscribeToSeedValidation();
  }

  generateSeed(entropy: number) {
    this.apiService.generateSeed(entropy).subscribe(seed => this.form.get('seed').setValue(seed));
  }

  private subscribeToSeedValidation() {
    if (this.seedValiditySubscription) {
      this.seedValiditySubscription.unsubscribe();
    }

    this.seedValiditySubscription = this.seed.asObservable().switchMap(seed => {
      if (!this.seedMatchValidator() || !this.create) {
        return this.apiService.post('wallet/seed/verify', {seed}, {}, true);
      } else {
        return Observable.of(0);
      }
    }).subscribe(() => {
      this.normalSeed = true;
    }, error => {
      if (error.status && error.status === 422) {
        this.normalSeed = false;
      } else {
        this.normalSeed = true;
      }
      this.subscribeToSeedValidation();
    });
  }

  private validatePasswords() {
    if (this.encrypt && this.form && this.form.get('password') && this.form.get('confirm_password')) {
      if (this.form.get('password').value) {
        if (this.form.get('password').value !== this.form.get('confirm_password').value) {
          return { NotEqual: true };
        }
      } else {
        return { Required: true };
      }
    }

    return null;
  }

  private seedMatchValidator() {
    if (this.form && this.form.get('seed') && this.form.get('confirm_seed')) {
      return this.form.get('seed').value === this.form.get('confirm_seed').value ? null : { NotEqual: true };
    } else {
      this.normalSeed = true;

      return { NotEqual: true };
    }
  }
}
