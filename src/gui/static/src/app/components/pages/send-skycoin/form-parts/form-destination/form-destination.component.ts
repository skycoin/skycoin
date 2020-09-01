import { SubscriptionLike } from 'rxjs';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { BigNumber } from 'bignumber.js';
import { TranslateService } from '@ngx-translate/core';

import { AppService } from '../../../../../services/app.service';
import { DoubleButtonActive } from '../../../../layout/double-button/double-button.component';
import { PriceService } from '../../../../../services/price.service';
import { MsgBarService } from '../../../../../services/msg-bar.service';
import { AvailableBalanceData } from '../../form-parts/form-source-selection/form-source-selection.component';
import { ConfirmationParams, ConfirmationComponent, DefaultConfirmationButtons } from '../../../../layout/confirmation/confirmation.component';
import { SendCoinsData } from '../../send-coins-form/send-coins-form.component';

/**
 * Data about the destinations entered by the user on FormDestinationComponent.
 */
export interface Destination {
  /**
   * Destination address.
   */
  address: string;
  /**
   * How many coins to send to the destination.
   */
  coins: string;
  /**
   * Original value entered by the user as the amount to send. It can be how many coins to send,
   * in which case the value will be the same as the one on the "coins" property, but it can
   * also contain a fiat value.
   */
  originalAmount: string;
  /**
   * How many hours to send to the destination.
   */
  hours?: string;
}

/**
 * Allows the user to set the destinations to were the coins will be sent, including how many
 * coins and hours to send to each one.
 */
@Component({
  selector: 'app-form-destination',
  templateUrl: './form-destination.component.html',
  styleUrls: ['./form-destination.component.scss'],
})
export class FormDestinationComponent implements OnInit, OnDestroy {
  /**
   * How many decimals the user can use when entering a value in usd. the UI can use a different
   * value when showing usd values (normally 2).
   */
  private static readonly MaxUsdDecimals = 6;

  // Balance available to send.
  @Input() availableBalance: AvailableBalanceData;
  // Allows to deactivate the form while the system is busy.
  @Input() busy: boolean;
  // Emits when there have been changes in the contents of the component, so the validation
  // status could have changed.
  @Output() onChanges = new EventEmitter<void>();
  // Emits when the user asks to open the modal window for bulk sending.
  @Output() onBulkRequested = new EventEmitter<void>();

  // If the manual hours field must be shown.
  private showHourFieldsInternal: boolean;
  @Input() set showHourFields(val: boolean) {
    if (val !== this.showHourFieldsInternal) {
      this.showHourFieldsInternal = val;
      if (this.form) {
        this.destControls.forEach(dest => {
          dest.get('hours').setValue('');
        });
      }
    }
  }
  get showHourFields(): boolean {
    return this.showHourFieldsInternal;
  }

  // If true, the form only allows to enter one destination and the amount of coins to send
  // to it. If false, the form allows multiple destinations, with their coins and hours.
  private showSimpleFormInternal: boolean;
  @Input() set showSimpleForm(val: boolean) {
    this.showSimpleFormInternal = val;

    if (this.form) {
      if (val) {
        this.form.get('address').setValidators(Validators.required);
      } else {
        this.form.get('address').clearValidators();
      }

      this.form.get('address').updateValueAndValidity();
      this.form.get('destinations').updateValueAndValidity();
    }
  }
  get showSimpleForm(): boolean {
    return this.showSimpleFormInternal;
  }

  form: FormGroup;
  doubleButtonActive = DoubleButtonActive;
  // Allows to know if the user is entering the values in the coin (left) or usd (right).
  selectedCurrency = DoubleButtonActive.LeftButton;
  // If the user is entering the values in USD, it contains the coins value of each destination.
  // if the user is entering the values in coins, it contains the USD value of each destination.
  values: BigNumber[];
  // Current USD price per coin.
  price: number;
  // Total amount of coins that will be sent to all destinations.
  totalCoins = new BigNumber(0);
  // Total usd value that will be sent to all destinations.
  totalFiat = new BigNumber(0);
  // Total amount of hours that will be sent to all destinations, if the manual hours are active.
  totalHours = new BigNumber(0);

  // Vars with the validation error messages.
  addressErrorMsgs: string[] = [];
  coinsErrorMsgs: string[] = [];
  hoursErrorMsgs: string[] = [];
  singleAddressErrorMsg = '';
  insufficientCoins = false;
  insufficientHours = false;

  // List for knowing which destination addresses were indentified as valid by the server,
  // by index.
  validAddressesList: boolean[];

  private priceSubscription: SubscriptionLike;
  private addressSubscription: SubscriptionLike;
  private destinationSubscriptions: SubscriptionLike[] = [];

  // Gets all the form field groups on the destinations array.
  get destControls() {
    return (this.form.get('destinations') as FormArray).controls;
  }

  constructor(
    private appService: AppService,
    private formBuilder: FormBuilder,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    private translate: TranslateService,
    private priceService: PriceService,
  ) { }

  ngOnInit() {
    this.form = this.formBuilder.group({
      address: [''],
      destinations: this.formBuilder.array([]),
    });
    this.form.setValidators(this.validateForm.bind(this));

    this.addDestination();

    // Inform when there are changes on the address field, shown on the simple form.
    this.addressSubscription = this.form.get('address').valueChanges.subscribe(() => {
      this.onChanges.emit();
    });

    // Keep the price updated.
    this.priceSubscription = this.priceService.price.subscribe(price => {
      this.price = price;
      this.updateValuesAndValidity();
    });
  }

  ngOnDestroy() {
    this.addressSubscription.unsubscribe();
    this.priceSubscription.unsubscribe();
    this.destinationSubscriptions.forEach(s => s.unsubscribe());
  }

  // Changes the currency in which the user enters the values on the UI.
  changeActiveCurrency(value: DoubleButtonActive) {
    if (value !== this.selectedCurrency) {
      this.selectedCurrency = value;
      this.askIfConvertAmount();
      this.updateValuesAndValidity();
      (this.form.get('destinations') as FormArray).updateValueAndValidity();
    }
  }

  // Must be called just after changing the currency in which the user enters the values on
  // the UI. It asks the user if the current values must be converted to the new currency.
  // If the user accepts, it calls the function to convert the values.
  private askIfConvertAmount() {
    // Before asking, check if there are valid values to convert.
    let validAmounts = 0;
    this.destControls.forEach(dest => {
      let value: string = dest.get('coins').value;
      value = value ? value.trim() : value;
      const currentValue = new BigNumber(value);

      if (!value || currentValue.isNaN()) {
        return;
      }

      validAmounts += 1;
    });
    if (validAmounts === 0) {
      return;
    }

    // Prepare the values for the confirmation modal.
    const usd = this.translate.instant('common.usd');
    const currentCoin = this.appService.coinName;
    let fromText: string;
    let toText: string;
    if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
      fromText = usd;
      toText = currentCoin;
    } else {
      fromText = currentCoin;
      toText = usd;
    }

    const confirmationParams: ConfirmationParams = {
      text: this.translate.instant(validAmounts === 1 ? 'send.convert-confirmation' : 'send.convert-confirmation-plural', {from: fromText, to: toText}),
      defaultButtons: DefaultConfirmationButtons.YesNo,
    };

    // Ask for confirmation.
    ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
      if (confirmationResult) {
        this.convertAmounts();
      }
    });
  }

  // If the component is set for the user to enter the values in usd, converts all the
  // destination values from the coin to usd. If it is set for the user to enter the
  // values in the coin, converts all the destination values from usd to the coin.
  private convertAmounts() {
    this.msgBarService.hide();

    // How many values were invalid numbers and it was not possible to convert them.
    let invalidValues = 0;
    // How many values were converted but had precision problems due to the amount of
    // decimal places.
    let valuesWithPrecisionErrors = 0;

    this.destControls.forEach(dest => {
      let value: string = dest.get('coins').value;
      value = value ? value.trim() : value;
      const currentValue = new BigNumber(value);

      if (value) {
        if (currentValue.isNaN()) {
          invalidValues += 1;

          return;
        }

        // Convert the value and check for precision errors.
        if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
          const newValue = currentValue.dividedBy(this.price).decimalPlaces(this.appService.currentMaxDecimals);
          const recoveredValue = newValue.multipliedBy(this.price).decimalPlaces(FormDestinationComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
          if (!recoveredValue.isEqualTo(currentValue)) {
            valuesWithPrecisionErrors += 1;
          }

          dest.get('coins').setValue(newValue.toString());
        } else {
          const newValue = currentValue.multipliedBy(this.price).decimalPlaces(FormDestinationComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
          const recoveredValue = newValue.dividedBy(this.price).decimalPlaces(this.appService.currentMaxDecimals);
          if (!recoveredValue.isEqualTo(currentValue)) {
            valuesWithPrecisionErrors += 1;
          }

          dest.get('coins').setValue(newValue.toString());
        }
      }
    });

    // Inform about any problem found during the procedure.
    if (invalidValues > 0 && valuesWithPrecisionErrors > 0) {
      this.msgBarService.showWarning(this.translate.instant('send.multiple-problems-warning'));
    } else if (invalidValues === 1) {
      this.msgBarService.showWarning(this.translate.instant('send.invaid-amount-warning'));
    } else if (invalidValues > 1) {
      this.msgBarService.showWarning(this.translate.instant('send.invaid-amounts-warning'));
    } else if (valuesWithPrecisionErrors === 1) {
      this.msgBarService.showWarning(this.translate.instant('send.precision-error-warning'));
    } else if (valuesWithPrecisionErrors > 1) {
      this.msgBarService.showWarning(this.translate.instant('send.precision-errors-warning'));
    }
  }

  // Assigns all the remaining coins to the destination corresponding to the provided index.
  assignAll(index: number) {
    this.msgBarService.hide();

    // If there are no available coins on the selected sources, show an error msg.
    if (this.availableBalance.availableCoins.isEqualTo(0)) {
      this.msgBarService.showError(this.translate.instant('send.no-wallet-selected-error'));

      return;
    }

    // Calculate the total available balance, in the currency being used on the UI.
    let availableBalance: BigNumber = this.availableBalance.availableCoins;
    if (this.selectedCurrency === DoubleButtonActive.RightButton) {
      availableBalance = availableBalance.multipliedBy(this.price).decimalPlaces(FormDestinationComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
    }

    // Subtract to the available balance all the values already asigned to the other destinations.
    this.destControls.forEach((dest, i) => {
      if (i !== index) {
        const value = this.getAmount((dest.get('coins').value as string).trim(), true);
        if (!value || value.isNaN()) {
          return;
        } else {
          availableBalance = availableBalance.minus(value);
        }
      }
    });

    // Limit the decimal places.
    if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
      availableBalance = availableBalance.decimalPlaces(this.appService.currentMaxDecimals, BigNumber.ROUND_FLOOR);
    } else {
      availableBalance = availableBalance.decimalPlaces(FormDestinationComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
    }

    // Use the value or show an error, if there are no enough coins.
    if (availableBalance.isLessThanOrEqualTo(0)) {
      this.msgBarService.showError(this.translate.instant('send.no-coins-left-error'));
    } else {
      this.destControls[index].get('coins').setValue(availableBalance.toString());
    }
  }

  // Updates the total and converted values and then updates the form validity.
  updateValuesAndValidity() {
    let inputInUsd = this.selectedCurrency !== DoubleButtonActive.LeftButton;
    let currentPrice = this.price;
    if (!this.price) {
      inputInUsd = false;
      currentPrice = 0;
    }

    // Reset the values.
    this.values = [];
    this.totalCoins = new BigNumber(0);
    this.totalFiat = new BigNumber(0);
    this.totalHours = new BigNumber(0);

    this.destControls.forEach((dest, i) => {
      // Update the coin values.
      let stringValue: string = dest.get('coins').value;
      let value = this.getAmount(stringValue, true);

      if (!value) {
        this.values[i] = new BigNumber(-1);
      } else {
        if (!inputInUsd) {
          // Calculate the value in USD.
          const convertedValue = value.multipliedBy(currentPrice).decimalPlaces(2);

          // Update the values.
          this.totalCoins = this.totalCoins.plus(value);
          this.totalFiat = this.totalFiat.plus(convertedValue);
          this.values[i] = convertedValue;
        } else {
          // Calculate the value in coins.
          const convertedValue = value.dividedBy(currentPrice).decimalPlaces(this.appService.currentMaxDecimals);

          // Update the values.
          this.totalCoins = this.totalCoins.plus(convertedValue);
          this.totalFiat = this.totalFiat.plus(value);
          this.values[i] = convertedValue;
        }
      }

      // Update the hour values.
      if (this.showHourFields) {
        stringValue = dest.get('hours').value;
        value = this.getAmount(stringValue, false);
        if (value) {
          this.totalHours = this.totalHours.plus(value);
        }
      }
    });

    // Update the form validity.
    setTimeout(() => {
      (this.form.get('destinations') as FormArray).updateValueAndValidity();
      this.onChanges.emit();
    });
  }

  // Adds a new empty destination to the form.
  addDestination() {
    const group = this.formBuilder.group({
      address: '',
      coins: '',
      hours: '',
    });

    this.destinationSubscriptions.push(group.valueChanges.subscribe(() => {
      // Inform when there are changes.
      this.updateValuesAndValidity();
    }));

    (this.form.get('destinations') as FormArray).push(group);
    this.addressErrorMsgs.push('');
    this.coinsErrorMsgs.push('');
    this.hoursErrorMsgs.push('');

    this.updateValuesAndValidity();
  }

  // Removes from the form the destination corresponding to the provided index.
  removeDestination(index) {
    const destinations = this.form.get('destinations') as FormArray;
    destinations.removeAt(index);

    // Remove the associated entry in the error arrays, if needed.
    if (this.validAddressesList && this.validAddressesList.length > index) {
      this.validAddressesList.splice(index, 1);
    }
    if (this.addressErrorMsgs && this.addressErrorMsgs.length > index) {
      this.addressErrorMsgs.splice(index, 1);
    }
    if (this.coinsErrorMsgs && this.coinsErrorMsgs.length > index) {
      this.coinsErrorMsgs.splice(index, 1);
    }
    if (this.hoursErrorMsgs && this.hoursErrorMsgs.length > index) {
      this.hoursErrorMsgs.splice(index, 1);
    }

    // Remove the subscription used to check the changes made to the fields of the destination.
    this.destinationSubscriptions[index].unsubscribe();
    this.destinationSubscriptions.splice(index, 1);

    this.updateValuesAndValidity();
  }

  requestBulkSend() {
    this.onBulkRequested.emit();
  }

  /**
   * Fills the form with the provided values.
   */
  fill(formData: SendCoinsData) {
    setTimeout(() => {
      this.selectedCurrency = formData.form.currency;

      for (let i = 0; i < formData.form.destinations.length - 1; i++) {
        this.addDestination();
      }

      this.destControls.forEach((destControl, i) => {
        ['address', 'hours'].forEach(name => {
          destControl.get(name).setValue(formData.form.destinations[i][name]);
        });
        destControl.get('coins').setValue(formData.form.destinations[i].originalAmount);

        if (this.showSimpleForm) {
          this.form.get('address').setValue(formData.form.destinations[i]['address']);
        }
      });

      this.updateValuesAndValidity();
    });
  }

  /**
   * Allows to change all the destinations.
   */
  setDestinations(newDestinations: Destination[]) {
    while (this.destControls.length > 0) {
      (this.form.get('destinations') as FormArray).removeAt(0);
    }

    newDestinations.forEach((destination, i) => {
      this.addDestination();
      this.destControls[i].get('address').setValue(destination.address);
      this.destControls[i].get('coins').setValue(destination.coins);
      if (destination.hours) {
        this.destControls[i].get('hours').setValue(destination.hours);
      }
    });
  }

  /**
   * Allows to know if the form is valid.
   */
  get valid(): boolean {
    return this.form.valid;
  }

  /**
   * Allows to know if the user is entering the values in coins (left) or usd (right).
   */
  get currentlySelectedCurrency(): DoubleButtonActive {
    return this.selectedCurrency;
  }

  /**
   * Allows to set a list indicating which addresses are valid, as validated by
   * the backend.
   * @param list Validity list. It must include if the address is valid for the index of
   * each destination. It can be null, to show all addresses as valid.
   */
  setValidAddressesList(list: boolean[]) {
    this.validAddressesList = list;

    if (this.validAddressesList && this.validAddressesList.length > this.destControls.length) {
      this.validAddressesList = this.validAddressesList.slice(0, this.destControls.length);
    }
  }

  /**
   * Allows to check if the address of a destination must be shown as valid, as validated by
   * the backend.
   * @param addressIndex Index of the address.
   */
  isAddressValid(addressIndex: number): boolean {
    if (this.validAddressesList && this.validAddressesList.length > addressIndex) {
      return this.validAddressesList[addressIndex];
    }

    return true;
  }

  /**
   * Gets the error msg that has to be shown for the coins field of a destination.
   * @param destinationIndex Index of the destination.
   */
  getCoinsErrorMsg(destinationIndex: number): string {
    if (destinationIndex < this.coinsErrorMsgs.length) {
      // Check if there is a validation error.
      if (this.coinsErrorMsgs[destinationIndex]) {
        return this.coinsErrorMsgs[destinationIndex];
      }

      // Check if the user is trying to send more coins than available, but only if
      // there is just one destination.
      if (this.destControls.length === 1 && this.insufficientCoins) {
        return 'send.insufficient-funds-error-info';
      }
    }

    return '';
  }

  /**
   * Gets the error msg that has to be shown for the hours field of a destination.
   * @param destinationIndex Index of the destination.
   */
  gethoursErrorMsg(destinationIndex: number): string {
    if (destinationIndex < this.hoursErrorMsgs.length) {
      // Check if there is a validation error.
      if (this.hoursErrorMsgs[destinationIndex]) {
        return this.hoursErrorMsgs[destinationIndex];
      }

      // Check if the user is trying to send more hours than available, but only if
      // there is just one destination.
      if (this.destControls.length === 1 && this.insufficientHours) {
        return 'send.insufficient-funds-error-info';
      }
    }

    return '';
  }

  /**
   * Returns all the destinations on the form. The hours are returned only if the form is showing
   * the fields for manually entering them.
   * @param cleanNumbers If true, the returned strings for the coins and hours will be cleaned
   * to be valid numbers. If false, function will just return exactly what the user wrote on
   * the form fields.
   */
  getDestinations(cleanNumbers: boolean): Destination[] {
    return this.destControls.map((destControl, i) => {
      // Get the string values.
      const destination = {
        address: this.showSimpleForm ? ((this.form.get('address').value) as string).trim() : ((destControl.get('address').value) as string).trim(),
        coins: ((this.selectedCurrency === DoubleButtonActive.LeftButton ? destControl.get('coins').value : this.values[i].toString()) as string).trim(),
        originalAmount: destControl.get('coins').value,
      };

      // Clean the values values.
      if (cleanNumbers) {
        destination.coins = new BigNumber(destination.coins).toString();
        destination.originalAmount = new BigNumber(destination.originalAmount).toString();
      }

      if (this.showHourFields) {
        destination['hours'] = destControl.get('hours').value;
        if (cleanNumbers) {
          destination['hours'] = new BigNumber(destination['hours']).toString();
        }
      }

      return destination;
    });
  }

  /**
   * Makes sure an errors array is the same size as "this.destControls.length" and sets all the
   * values to empty strings.
   */
  private resetErrorMsgsArray(array: string[]) {
    while (array.length > this.destControls.length) {
      array.pop();
    }
    while (array.length < this.destControls.length) {
      array.push('');
    }
    for (let i = 0; i < array.length; i++) {
      array[i] = '';
    }
  }

  /**
   * Validates the form and updates the vars with the validation errors.
   */
  validateForm() {
    this.singleAddressErrorMsg = '';
    this.resetErrorMsgsArray(this.addressErrorMsgs);
    this.resetErrorMsgsArray(this.coinsErrorMsgs);
    this.resetErrorMsgsArray(this.hoursErrorMsgs);
    this.insufficientCoins = false;
    this.insufficientHours = false;

    let valid = true;

    // Check the address field of the simple form.
    if (this.showSimpleForm) {
      const address = this.form.get('address').value as string;
      if (!address || address.length < 20) {
        valid = false;
        if (this.form.get('address').touched) {
          this.singleAddressErrorMsg = 'send.address-error-info';
        }
      }
    }

    // Check if there are invalid values.
    this.destControls.forEach((control, i) => {
      // Check the address, but not if showing the simple form.
      if (!this.showSimpleForm) {
        const address = control.get('address').value as string;
        if (!address || address.length < 20) {
          valid = false;
          if (control.get('address').touched) {
            this.addressErrorMsgs[i] = 'send.address-error-info';
          }
        }
      }

      // Check the coins.
      const coinsValue: string = control.get('coins').value;
      if (this.getAmount(coinsValue, true) === null) {
        valid = false;
        if (control.get('coins').touched) {
          this.coinsErrorMsgs[i] = 'send.invalid-value-error-info';
        }
      }

      // Check the hours, if showing the hours field.
      if (this.showHourFields) {
        const hoursValue: string = control.get('hours').value;
        if (this.getAmount(hoursValue, false) === null) {
          valid = false;
          if (control.get('hours').touched) {
            this.hoursErrorMsgs[i] = 'send.invalid-value-error-info';
          }
        }
      }
    });

    // Check how many coins and hours the user is trying to send.
    let destinationsCoins = new BigNumber(0);
    if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
      this.destControls.map(control => {
        const value = new BigNumber(control.get('coins').value);
        if (!value.isNaN()) {
          destinationsCoins = destinationsCoins.plus(value);
        }
      });
    } else {
      this.updateValuesAndValidity();
      this.values.map(value => {
        if (!value.isNaN()) {
          destinationsCoins = destinationsCoins.plus(value);
        }
      });
    }
    let destinationsHours = new BigNumber(0);
    if (this.showHourFields) {
      this.destControls.map(control => {
        const value = new BigNumber(control.get('hours').value);
        if (!value.isNaN()) {
          destinationsHours = destinationsHours.plus(control.get('hours').value);
        }
      });
    }

    // Fail if the user does not have enough coins or hours.
    if (destinationsCoins.isGreaterThan(this.availableBalance.availableCoins)) {
      this.insufficientCoins = true;
      valid = false;
    }
    if (destinationsHours.isGreaterThan(this.availableBalance.availableHours)) {
      this.insufficientHours = true;
      valid = false;
    }

    return valid ? null : { Invalid: true };
  }

  /**
   * Process a string and converts it to a BigNumber representing a coins or hours amount. If
   * the string is not a valid value, it returns null.
   * @param stringValue String to process.
   * @param checkingCoins If the function must treat the value as coins or hours while checking
   * for validity.
   */
  private getAmount(stringValue: string, checkingCoins: boolean): BigNumber {
    stringValue = stringValue ? stringValue.trim() : stringValue;
    const value = new BigNumber(stringValue);

    // Check for basic validity.
    if (!stringValue || value.isNaN()) {
      return null;
    }
    if (checkingCoins && value.isLessThanOrEqualTo(0)) {
      return null;
    }
    if (!checkingCoins && value.isLessThan(0)) {
      return null;
    }

    // Check the decimals.
    if (checkingCoins) {
      const parts = stringValue.split('.');

      if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
        if (parts.length === 2 && parts[1].length > this.appService.currentMaxDecimals) {
          return null;
        }
      } else {
        if (parts.length === 2 && parts[1].length > FormDestinationComponent.MaxUsdDecimals) {
          return null;
        }
      }
    } else {
      if (!value.isEqualTo(value.decimalPlaces(0))) {
        return null;
      }
    }

    return value;
  }

  resetForm() {
    this.form.get('address').setValue('');

    while (this.destControls.length > 0) {
      (this.form.get('destinations') as FormArray).removeAt(0);
    }

    this.addDestination();

    this.updateValuesAndValidity();
  }
}
