import { SubscriptionLike } from 'rxjs';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { BigNumber } from 'bignumber.js';
import { BlockchainService } from '../../../../../services/blockchain.service';
import { AppService } from '../../../../../services/app.service';
import { TranslateService } from '@ngx-translate/core';
import { DoubleButtonActive } from '../../../../layout/double-button/double-button.component';
import { PriceService } from '../../../../../services/price.service';
import { MsgBarService } from '../../../../../services/msg-bar.service';
import { AvailableBalanceData } from '../../form-parts/form-source-selection/form-source-selection.component';
import { ConfirmationParams, ConfirmationComponent, DefaultConfirmationButtons } from '../../../../layout/confirmation/confirmation.component';

export interface Destination {
  address: string;
  coins: string;
  originalAmount: string;
  hours?: string;
}

@Component({
  selector: 'app-form-destination',
  templateUrl: './form-destination.component.html',
  styleUrls: ['./form-destination.component.scss'],
})
export class FormDestinationComponent implements OnInit, OnDestroy {
  private static readonly MaxUsdDecimals = 6;

  @Input() availableBalance: AvailableBalanceData;
  @Input() busy: boolean;
  @Output() onChanges = new EventEmitter<void>();
  @Output() onBulkRequested = new EventEmitter<void>();

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
  selectedCurrency = DoubleButtonActive.LeftButton;
  values: number[];
  price: number;
  totalCoins = new BigNumber(0);
  totalConvertedCoins = new BigNumber(0);
  totalHours = new BigNumber(0);

  private priceSubscription: SubscriptionLike;
  private addressSubscription: SubscriptionLike;
  private destinationSubscriptions: SubscriptionLike[] = [];

  get destControls() {
    return (this.form.get('destinations') as FormArray).controls;
  }

  constructor(
    private blockchainService: BlockchainService,
    private appService: AppService,
    private formBuilder: FormBuilder,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    private translate: TranslateService,
    private priceService: PriceService,
  ) { }

  ngOnInit() {
    this.form = this.formBuilder.group({
      address: ['', this.showSimpleForm ? Validators.required : null],
      destinations: this.formBuilder.array(
        [this.createDestinationFormGroup()],
        this.validateDestinations.bind(this),
      ),
    });

    this.addressSubscription = this.form.get('address').valueChanges.subscribe(value => {
      this.onChanges.emit();
    });

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

  changeActiveCurrency(value) {
    if (value !== this.selectedCurrency) {
      this.selectedCurrency = value;
      this.askIfConvertAmount();
      this.updateValuesAndValidity();
      (this.form.get('destinations') as FormArray).updateValueAndValidity();
    }
  }

  private askIfConvertAmount() {
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

    ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
      if (confirmationResult) {
        this.convertAmounts();
      }
    });
  }

  private convertAmounts() {
    this.msgBarService.hide();

    let invalidValues = 0;
    let valuesWithPrecisionErrors = 0;
    this.destControls.forEach(dest => {
      let value: string = dest.get('coins').value;
      value = value ? value.trim() : value;
      const currentValue = new BigNumber(value);

      if (value) {
        if (!value || currentValue.isNaN()) {
          invalidValues += 1;

          return;
        }

        if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
          const newValue = currentValue.dividedBy(this.price).decimalPlaces(this.blockchainService.currentMaxDecimals);
          const recoveredValue = newValue.multipliedBy(this.price).decimalPlaces(FormDestinationComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
          if (!recoveredValue.isEqualTo(currentValue)) {
            valuesWithPrecisionErrors += 1;
          }

          dest.get('coins').setValue(newValue.toString());
        } else {
          const newValue = currentValue.multipliedBy(this.price).decimalPlaces(FormDestinationComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
          const recoveredValue = newValue.dividedBy(this.price).decimalPlaces(this.blockchainService.currentMaxDecimals);
          if (!recoveredValue.isEqualTo(currentValue)) {
            valuesWithPrecisionErrors += 1;
          }

          dest.get('coins').setValue(newValue.toString());
        }
      }
    });

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

  assignAll(index: number) {
    this.msgBarService.hide();

    let availableCoins: BigNumber = this.availableBalance.availableCoins;
    if (this.selectedCurrency === DoubleButtonActive.RightButton) {
      availableCoins = availableCoins.multipliedBy(this.price).decimalPlaces(FormDestinationComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
    }

    this.destControls.forEach((dest, i) => {
      if (i !== index) {
        const value = Number.parseFloat((dest.get('coins').value as string).trim());
        if (!value || isNaN(value)) {
          return;
        } else {
          availableCoins = availableCoins.minus(value);
        }
      }
    });

    if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
      availableCoins = availableCoins.decimalPlaces(this.blockchainService.currentMaxDecimals, BigNumber.ROUND_FLOOR);
    } else {
      availableCoins = availableCoins.decimalPlaces(FormDestinationComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
    }

    if (availableCoins.isLessThan(0)) {
      this.msgBarService.showError(this.translate.instant('send.no-coins-left-error'));
    } else {
      this.destControls[index].get('coins').setValue(availableCoins.toString());
    }
  }

  updateValuesAndValidity() {
    let inputInUsd = this.selectedCurrency !== DoubleButtonActive.LeftButton;
    let currentPrice = this.price;
    if (!this.price) {
      inputInUsd = false;
      currentPrice = 0;
    }

    this.values = [];
    this.totalCoins = new BigNumber(0);
    this.totalConvertedCoins = new BigNumber(0);
    this.totalHours = new BigNumber(0);

    this.destControls.forEach((dest, i) => {
      const stringValue: string = dest.get('coins').value;
      const value = this.getAmount(stringValue, true);
      if (!value) {
        this.values[i] = -1;

        return;
      }

      if (!inputInUsd) {
        const convertedValue = value.multipliedBy(currentPrice).decimalPlaces(2);

        this.totalCoins = this.totalCoins.plus(value);
        this.totalConvertedCoins = this.totalConvertedCoins.plus(convertedValue);

        this.values[i] = convertedValue.toNumber();
      } else {
        const convertedValue = value.dividedBy(currentPrice).decimalPlaces(this.blockchainService.currentMaxDecimals);

        this.totalCoins = this.totalCoins.plus(convertedValue);
        this.totalConvertedCoins = this.totalConvertedCoins.plus(value);

        this.values[i] = convertedValue.toNumber();
      }
    });

    this.destControls.forEach(dest => {
      const stringValue: string = dest.get('hours').value;
      const value = this.getAmount(stringValue, false);
      if (!value) {
        return;
      }

      this.totalHours = this.totalHours.plus(value);
    });

    setTimeout(() => {
      (this.form.get('destinations') as FormArray).updateValueAndValidity();
      this.onChanges.emit();
    });
  }

  addDestination() {
    const destinations = this.form.get('destinations') as FormArray;
    destinations.push(this.createDestinationFormGroup());
    this.updateValuesAndValidity();
  }

  removeDestination(index) {
    const destinations = this.form.get('destinations') as FormArray;
    destinations.removeAt(index);

    this.destinationSubscriptions[index].unsubscribe();
    this.destinationSubscriptions.splice(index, 1);
    this.updateValuesAndValidity();
  }

  requestBulkSend() {
    this.onBulkRequested.emit();
  }

  fill(formData: any) {
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

  get valid(): boolean {
    return this.form.valid;
  }

  get currentlySelectedCurrency(): DoubleButtonActive {
    return this.selectedCurrency;
  }

  getDestinations(includeHours: boolean, cleanNumbers: boolean): Destination[] {
    return this.destControls.map((destControl, i) => {
      const destination = {
        address: this.showSimpleForm ? ((this.form.get('address').value) as string).trim() : ((destControl.get('address').value) as string).trim(),
        coins: ((this.selectedCurrency === DoubleButtonActive.LeftButton ? destControl.get('coins').value : this.values[i].toString()) as string).trim(),
        originalAmount: destControl.get('coins').value,
      };

      if (cleanNumbers) {
        destination.coins = new BigNumber(destination.coins).toString();
        destination.originalAmount = new BigNumber(destination.originalAmount).toString();
      }

      if (includeHours) {
        destination['hours'] = destControl.get('hours').value;
        if (cleanNumbers) {
          destination['hours'] = new BigNumber(destination['hours']).toString();
        }
      }

      return destination;
    });
  }

  private validateDestinations() {
    if (!this.form) {
      return { Required: true };
    }

    const invalidInput = this.destControls.find(control => {
      if (!this.showSimpleForm && (!control.get('address').value || (control.get('address').value as string).trim().length === 0)) {
        return true;
      }

      const checkControls = ['coins'];

      if (this.showHourFields) {
        checkControls.push('hours');
      }

      return checkControls.map(name => {
        const stringValue: string = control.get(name).value;

        return this.getAmount(stringValue, name === 'coins') === null;
      }).find(e => e === true);
    });

    if (invalidInput) {
      return { Invalid: true };
    }

    let destinationsCoins = new BigNumber(0);
    if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
      this.destControls.map(control => destinationsCoins = destinationsCoins.plus(control.get('coins').value));
    } else {
      this.updateValuesAndValidity();
      this.values.map(value => destinationsCoins = destinationsCoins.plus(value));
    }
    let destinationsHours = new BigNumber(0);
    if (this.showHourFields) {
      this.destControls.map(control => destinationsHours = destinationsHours.plus(control.get('hours').value));
    }

    if (destinationsCoins.isGreaterThan(this.availableBalance.availableCoins) || destinationsHours.isGreaterThan(this.availableBalance.availableHours)) {
      return { Invalid: true };
    }

    return null;
  }

  private getAmount(stringValue: string, checkingCoins: boolean): BigNumber {
    stringValue = stringValue ? stringValue.trim() : stringValue;
    const value = new BigNumber(stringValue);

    if (!stringValue || value.isNaN() || value.isLessThanOrEqualTo(0)) {
      return null;
    }

    if (checkingCoins) {
      const parts = stringValue.split('.');

      if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
        if (parts.length === 2 && parts[1].length > this.blockchainService.currentMaxDecimals) {
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

  private createDestinationFormGroup() {
    const group = this.formBuilder.group({
      address: '',
      coins: '',
      hours: '',
    });

    this.destinationSubscriptions.push(group.valueChanges.subscribe(value => {
      this.updateValuesAndValidity();
    }));

    return group;
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
