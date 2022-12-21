import { switchMap, delay, flatMap } from 'rxjs/operators';
import { Component, OnInit, OnDestroy, Input, Output, EventEmitter, ViewChild } from '@angular/core';
import { UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { SubscriptionLike, Subject, of } from 'rxjs';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { MatLegacyCheckbox as MatCheckbox } from '@angular/material/legacy-checkbox';

import { ApiService } from '../../../../../services/api.service';
import { WordAskedReasons } from '../../../../layout/seed-word-dialog/seed-word-dialog.component';
import { MsgBarService } from '../../../../../services/msg-bar.service';
import { ConfirmationParams, ConfirmationComponent, DefaultConfirmationButtons } from '../../../../layout/confirmation/confirmation.component';
import { OperationError } from '../../../../../utils/operation-error';
import { processServiceError } from '../../../../../utils/errors';
import { AssistedSeedFieldComponent } from './assisted-seed-field/assisted-seed-field.component';

/**
 * Data entered in an instance of CreateWalletFormComponent.
 */
export class WalletFormData {
  /**
   * If the form is for creating a new wallet (true) or loading a walled using a seed (false).
   */
  creatingNewWallet: boolean;
  /**
   * Label for the wallet.
   */
  label: string;
  /**
   * Seed for the wallet.
   */
  seed: string;
  /**
   * If set, the wallet must be encrypted with this password.
   */
  password: string;
  /**
   * If true, the seed was entered using the assisted mode.
   */
  enterSeedWithAssistance: boolean;
  /**
   * If creating a new wallet, the last automatically generated seed for the assisted mode.
   */
  lastAssistedSeed: string;
  /**
   * Last seed the user entered using the manual mode.
   */
  lastCustomSeed: string;
  /**
   * If creating a new wallet, how many words the automatically generated seed for the
   * assisted mode has. If loading a wallet, how many words the seed entered by the user with
   * the assisted mode has.
   */
  numberOfWords: number;
  /**
   * If the user entered a standard seed, if the manual mode was being used.
   */
  customSeedIsNormal: boolean;
  /**
   * If the advanced options panel was open.
   */
   advancedOptionsShown: boolean;
  /**
   * If the user selected that the wallet must be loaded temporarily.
   */
  loadTemporarily = false;
}

/**
 * Form for creating or loading a software wallet.
 */
@Component({
  selector: 'app-create-wallet-form',
  templateUrl: './create-wallet-form.component.html',
  styleUrls: ['./create-wallet-form.component.scss'],
})
export class CreateWalletFormComponent implements OnInit, OnDestroy {
  @ViewChild('temporalWalletCheck') temporalWalletCheck: MatCheckbox;
  // Component for entering the seed using the assisted mode.
  @ViewChild('assistedSeed') assistedSeed: AssistedSeedFieldComponent;

  // If the form is for creating a new wallet (true) or loading a walled using a seed (false).
  @Input() create: boolean;
  // If the form is being shown on the wizard (true) or not (false).
  @Input() onboarding: boolean;
  // Allows to deactivate the form while the system is busy.
  @Input() busy = false;
  // Emits when the user asks for the wallet ot be created.
  @Output() createRequested = new EventEmitter<void>();

  form: UntypedFormGroup;
  // If true, the user must enter the ssed using the asisted mode.
  enterSeedWithAssistance = true;
  // If the user entered a standard seed using the manual mode.
  customSeedIsNormal = true;
  // If the user entered a non-standard seed using the manual mode and confirmed to use it.
  customSeedAccepted = false;
  // If the user selected that the wallet must be created encrypted.
  encrypt = true;
  // If the advanced options must be shown on the UI.
  showAdvancedOptions = false;
  // If the user selected that the wallet must be loaded temporarily.
  loadTemporarily = false;
  // If creating a new wallet, the last automatically generated seed for the assisted mode.
  lastAssistedSeed = '';
  // How many words the last autogenerated seed for the assisted mode has, when creating
  // a new wallet.
  numberOfAutogeneratedWords = 0;
  // If the system is currently checking the custom seed entered by the user.
  checkingCustomSeed = false;

  // Vars with the validation error messages.
  labelErrorMsg = '';
  seed1ErrorMsg = '';
  seed2ErrorMsg = '';
  password1ErrorMsg = '';
  password2ErrorMsg = '';

  wordAskedReasons = WordAskedReasons;

  // Emits every time the seed should be checked again, to know if it is a standard seed.
  private seed: Subject<string> = new Subject<string>();

  private statusSubscription: SubscriptionLike;
  private seedValiditySubscription: SubscriptionLike;

  constructor(
    private apiService: ApiService,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
  ) { }

  ngOnInit() {
    if (!this.onboarding) {
      this.initForm();
    } else {
      this.initForm(false, null);
    }
  }

  ngOnDestroy() {
    this.msgBarService.hide();
    this.statusSubscription.unsubscribe();
    this.seedValiditySubscription.unsubscribe();
    this.createRequested.complete();
  }

  // Allows to know if the form is valid.
  get isValid(): boolean {
    // When entering the seed manually, the system must have finished checking the seed and the
    // seed must be normal or the user must confirm the usage of a custom seed. When using the
    // assisted mode, the user must enter the seed in the appropriate way.
    return this.form.valid && !this.checkingCustomSeed &&
      (
        (!this.enterSeedWithAssistance && (this.customSeedIsNormal || this.customSeedAccepted)) ||
        (this.enterSeedWithAssistance && this.assistedSeed.lastAssistedSeed)
      );
  }

  // Sets if the user has acepted to use a manually entered non-standard seed.
  onCustomSeedAcceptance(event) {
    this.customSeedAccepted = event.checked;
  }

  // Sets the user selection regarding whether the wallet must be encrypted or not.
  setEncrypt(event) {
    this.encrypt = event.checked;
    this.form.updateValueAndValidity();
  }

  // Shows or hides the advanced options.
  toggleAdvancedOptions() {
    // Stop if trying to close the options when a value has been changed.
    if (this.showAdvancedOptions === true && this.loadTemporarily) {
      this.msgBarService.showError('wallet.new.close-advanced-error');

      return;
    }

    this.showAdvancedOptions = !this.showAdvancedOptions;
  }

  // Sets the user selection regarding whether the wallet must be loaded temporarily or not.
  setTemporal(event) {
    if (event.checked) {
      this.temporalWalletCheck.checked = false;

      // Ask for confirmation before making the change.
      const confirmationParams: ConfirmationParams = {
        text: 'wallet.new.temporal-warning',
        headerText: 'common.warning-title',
        defaultButtons: DefaultConfirmationButtons.ContinueCancel,
        redTitle: true,
      };

      ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          this.loadTemporarily = true;
          this.form.updateValueAndValidity();
        }
      });
    } else {
      this.loadTemporarily = false;
      this.form.updateValueAndValidity();
    }
  }

  // Returns the data entered on the form.
  getData(): WalletFormData {
    return {
      creatingNewWallet: this.create,
      label: this.form.value.label,
      seed: this.enterSeedWithAssistance ? this.assistedSeed.lastAssistedSeed : this.form.value.seed,
      password: !this.onboarding && this.encrypt ? this.form.value.password : null,
      enterSeedWithAssistance: this.enterSeedWithAssistance,
      lastAssistedSeed: this.lastAssistedSeed,
      lastCustomSeed: this.form.value.seed,
      numberOfWords: !this.create ? this.form.value.number_of_words : this.numberOfAutogeneratedWords,
      customSeedIsNormal: this.customSeedIsNormal,
      advancedOptionsShown: this.showAdvancedOptions,
      loadTemporarily: this.loadTemporarily,
    };
  }

  // Switches between the assisted mode and the manual mode for entering the seed.
  changeSeedType() {
    this.msgBarService.hide();

    if (!this.enterSeedWithAssistance) {
      this.enterSeedWithAssistance = true;
      this.removeConfirmations();
    } else {
      // Ask for confirmation before making the change.
      const confirmationParams: ConfirmationParams = {
        text: this.create ? 'wallet.new.seed.custom-seed-warning-text' : 'wallet.new.seed.custom-seed-warning-text-recovering',
        headerText: 'common.warning-title',
        checkboxText: this.create ? 'common.generic-confirmation-check' : null,
        defaultButtons: DefaultConfirmationButtons.ContinueCancel,
        redTitle: true,
      };

      ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          this.enterSeedWithAssistance = false;
          this.removeConfirmations();
        }
      });
    }
  }

  /**
   * Inits or resets the form.
   * @param create If the form is for creating a new wallet (true) or loading a walled using
   * a seed (false). Use null to avoid changing the value set using the html tag.
   * @param data Data to populate the form.
   */
  initForm(create: boolean = null, data: WalletFormData = null) {
    this.msgBarService.hide();

    create = create !== null ? create : this.create;

    this.lastAssistedSeed = '';
    this.enterSeedWithAssistance = true;

    this.form = new UntypedFormGroup({});
    this.form.addControl('label', new UntypedFormControl(data ? data.label : ''));
    this.form.addControl('seed', new UntypedFormControl(data ? data.lastCustomSeed : ''));
    this.form.addControl('confirm_seed', new UntypedFormControl(data ? data.lastCustomSeed : ''));
    this.form.addControl('password', new UntypedFormControl());
    this.form.addControl('confirm_password', new UntypedFormControl());
    this.form.addControl('number_of_words', new UntypedFormControl(!this.create && data && data.numberOfWords ? data.numberOfWords : 12));

    this.form.setValidators(this.validateForm.bind(this));

    this.removeConfirmations(false);

    // Create a new random seed.
    if (create && !data) {
      this.generateSeed(128);
    }

    // Use the provided data.
    if (data) {
      this.enterSeedWithAssistance = data.enterSeedWithAssistance;
      this.lastAssistedSeed = data.lastAssistedSeed;
      this.assistedSeed.lastAssistedSeed = create ? data.lastAssistedSeed : data.seed;
      this.customSeedAccepted = true;
      this.customSeedIsNormal = data.customSeedIsNormal;

      this.showAdvancedOptions = data.advancedOptionsShown;
      this.loadTemporarily = data.loadTemporarily;

      if (this.create) {
        this.numberOfAutogeneratedWords = data.numberOfWords;
      }
    }

    if (this.statusSubscription && !this.statusSubscription.closed) {
      this.statusSubscription.unsubscribe();
    }
    this.statusSubscription = this.form.statusChanges.subscribe(() => {
      // Invalidate the custom seed confirmation if the data on the form is changed.
      this.customSeedAccepted = false;
      this.seed.next(this.form.get('seed').value);
    });

    this.subscribeToSeedValidation();
  }

  // Generates a new random seed for when creating a new wallet.
  generateSeed(entropy: number) {
    if (entropy === 128) {
      this.numberOfAutogeneratedWords = 12;
    } else {
      this.numberOfAutogeneratedWords = 24;
    }

    this.apiService.get('wallet/newSeed', { entropy: entropy }).subscribe(response => {
      this.lastAssistedSeed = response.seed;
      this.form.get('seed').setValue(response.seed);
      this.form.get('seed').markAsTouched();
      this.removeConfirmations();
    });
  }

  // Request the wallet to be created or loaded.
  requestCreation() {
    this.createRequested.emit();
  }

  /**
   * Removes the confirmations the user could have made for accepting the seed.
   * @param cleanSecondSeedField If true, the second field for manually entering a seed (the
   * one used for confirming the seed by entering it again) will be cleaned.
   */
  private removeConfirmations(cleanSecondSeedField = true) {
    this.customSeedAccepted = false;
    if (this.assistedSeed) {
      this.assistedSeed.lastAssistedSeed = null;
    }
    if (cleanSecondSeedField) {
      this.form.get('confirm_seed').setValue('');
    }
    this.form.updateValueAndValidity();
  }

  // Makes the component continually check if the user has manually entered a non-standard seed.
  private subscribeToSeedValidation() {
    if (this.seedValiditySubscription) {
      this.seedValiditySubscription.unsubscribe();
    }

    this.seedValiditySubscription = this.seed.asObservable().pipe(switchMap(seed => {
      // Verify the seed if it was entered manually and was confirmed.
      if (!this.enterSeedWithAssistance && seed.trim().length > 0 && (!this.create || this.seedsAreEqual())) {
        this.checkingCustomSeed = true;

        return of(0).pipe(delay(500), flatMap(() => this.apiService.post('wallet/seed/verify', {seed: seed}, {useV2: true})));
      } else {
        return of(0);
      }
    })).subscribe(result => {
      // The entered seed does not have problems if the backend (not the previous code) returned
      // a success response.
      this.customSeedIsNormal = result !== 0;
      this.checkingCustomSeed = false;
    }, (error: OperationError) => {
      this.checkingCustomSeed = false;
      // If the node said the seed is not standard, ask the user for confirmation before
      // allowing to use it.
      error = processServiceError(error);
      if (error && error.originalError && error.originalError.status === 422) {
        this.customSeedIsNormal = false;
      } else {
        // There was a problem performing the procedure.
        this.customSeedIsNormal = true;
        this.msgBarService.showWarning('wallet.new.seed-checking-error');
      }
      this.subscribeToSeedValidation();
    });
  }

  // Checks if the 2 custom seeds entered by the user are equal.
  private seedsAreEqual(): boolean {
    if (this.form && this.form.get('seed') && this.form.get('confirm_seed')) {
      return this.form.get('seed').value === this.form.get('confirm_seed').value;
    }

    this.customSeedIsNormal = true;

    return false;
  }

  /**
   * Validates the form and updates the vars with the validation errors.
   */
  validateForm() {
    this.labelErrorMsg = '';
    this.seed1ErrorMsg = '';
    this.seed2ErrorMsg = '';
    this.password1ErrorMsg = '';
    this.password2ErrorMsg = '';

    let valid = true;

    if (!this.form.get('label').value) {
      valid = false;
      if (this.form.get('label').touched) {
        this.labelErrorMsg = 'wallet.new.name-error-info';
      }
    }

    // Validate custom seeds.
    if (!this.enterSeedWithAssistance) {
      let enteredSeeds = true;
      if (!this.form.get('seed').value) {
        valid = false;
        enteredSeeds = false;
        this.customSeedIsNormal = true;
        if (this.form.get('seed').touched) {
          this.seed1ErrorMsg = 'wallet.new.seed-error-info';
        }
      }

      if (this.create) {
        if (!this.form.get('confirm_seed').value) {
          valid = false;
          enteredSeeds = false;
          this.customSeedIsNormal = true;
          if (this.form.get('confirm_seed').touched) {
            this.seed2ErrorMsg = 'wallet.new.seed-error-info';
          }
        }

        if (enteredSeeds) {
          if (!this.seed2ErrorMsg && !this.seedsAreEqual()) {
            valid = false;
            this.customSeedIsNormal = true;
            this.seed2ErrorMsg = 'wallet.new.confirm-seed-error-info';
          }
        }
      }
    }

    // Validate password.
    if (this.encrypt &&  !this.loadTemporarily && !this.onboarding) {
      let enteredPasswords = true;

      if (!this.form.get('password').value) {
        valid = false;
        enteredPasswords = false;
        if (this.form.get('password').touched) {
          this.password1ErrorMsg = 'password.password-error-info';
        }
      }

      if (!this.form.get('confirm_password').value) {
        valid = false;
        enteredPasswords = false;
        if (this.form.get('confirm_password').touched) {
          this.password2ErrorMsg = 'password.password-error-info';
        }
      }

      if (enteredPasswords) {
        if (!this.password2ErrorMsg && this.form.get('password').value !== this.form.get('confirm_password').value) {
          valid = false;
          this.password2ErrorMsg = 'password.confirm-error-info';
        }
      }
    }

    return valid ? null : { Invalid: true };
  }
}
