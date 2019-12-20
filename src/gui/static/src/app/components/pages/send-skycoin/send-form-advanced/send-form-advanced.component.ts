import { SubscriptionLike } from 'rxjs';
import { first } from 'rxjs/operators';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild, ChangeDetectorRef } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { FormArray, FormBuilder, FormGroup } from '@angular/forms';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { ButtonComponent } from '../../../layout/button/button.component';
import { getHardwareWalletErrorMsg } from '../../../../utils/errors';
import { NavBarService } from '../../../../services/nav-bar.service';
import { SelectAddressComponent } from './select-address/select-address';
import { BigNumber } from 'bignumber.js';
import { ConfirmationData } from '../../../../app.datatypes';
import { BlockchainService } from '../../../../services/blockchain.service';
import { showConfirmationModal } from '../../../../utils';
import { AppService } from '../../../../services/app.service';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { TranslateService } from '@ngx-translate/core';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { PriceService } from '../../../../services/price.service';
import { SendFormComponent } from '../send-form/send-form.component';
import { ChangeNoteComponent } from '../send-preview/transaction-info/change-note/change-note.component';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { MultipleDestinationsDialogComponent } from '../../../layout/multiple-destinations-dialog/multiple-destinations-dialog.component';
import { FormSourceSelectionComponent, AvailableBalanceData } from '../form-parts/form-source-selection/form-source-selection.component';

@Component({
  selector: 'app-send-form-advanced',
  templateUrl: './send-form-advanced.component.html',
  styleUrls: ['./send-form-advanced.component.scss'],
})
export class SendFormAdvancedComponent implements OnInit, OnDestroy {
  @ViewChild('formSourceSelection', { static: false }) formSourceSelection: FormSourceSelectionComponent;
  @ViewChild('previewButton', { static: false }) previewButton: ButtonComponent;
  @ViewChild('sendButton', { static: false }) sendButton: ButtonComponent;
  @Input() formData: any;
  @Output() onFormSubmitted = new EventEmitter<any>();

  maxNoteChars = ChangeNoteComponent.MAX_NOTE_CHARS;
  form: FormGroup;
  availableBalance = new AvailableBalanceData();
  autoHours = true;
  autoOptions = false;
  autoShareValue = '0.5';
  previewTx: boolean;
  busy = false;
  doubleButtonActive = DoubleButtonActive;
  selectedCurrency = DoubleButtonActive.LeftButton;
  values: number[];
  price: number;
  totalCoins = new BigNumber(0);
  totalConvertedCoins = new BigNumber(0);
  totalHours = new BigNumber(0);

  private subscriptionsGroup: SubscriptionLike[] = [];
  private destinationSubscriptions: SubscriptionLike[] = [];
  private destinationHoursSubscriptions: SubscriptionLike[] = [];
  private syncCheckSubscription: SubscriptionLike;
  private processingSubscription: SubscriptionLike;

  constructor(
    public blockchainService: BlockchainService,
    public walletService: WalletService,
    private appService: AppService,
    private formBuilder: FormBuilder,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    private navbarService: NavBarService,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
    private priceService: PriceService,
    private changeDetector: ChangeDetectorRef,
  ) { }

  ngOnInit() {
    this.navbarService.showSwitch('send.simple', 'send.advanced', DoubleButtonActive.RightButton);

    this.form = this.formBuilder.group({
      changeAddress: [''],
      destinations: this.formBuilder.array(
        [this.createDestinationFormGroup()],
        this.validateDestinations.bind(this),
      ),
      note: [''],
    });
    this.subscriptionsGroup.push(this.priceService.price.subscribe(price => {
      this.price = price;
      this.updateValues();
    }));

    if (this.formData) {
      setTimeout(() => this.fillForm());
    }
  }

  ngOnDestroy() {
    if (this.processingSubscription && !this.processingSubscription.closed) {
      this.processingSubscription.unsubscribe();
    }
    this.closeSyncCheckSubscription();
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
    this.navbarService.hideSwitch();
    this.msgBarService.hide();
    this.destinationSubscriptions.forEach(s => s.unsubscribe());
    this.destinationHoursSubscriptions.forEach(s => s.unsubscribe());
  }

  sourceSelectionChanged() {
    this.availableBalance = this.formSourceSelection.availableBalance;
    this.form.get('destinations').updateValueAndValidity();
  }

  preview() {
    this.previewTx = true;
    this.checkBeforeSending();
    this.changeDetector.detectChanges();
  }

  send() {
    this.previewTx = false;
    this.checkBeforeSending();
  }

  changeActiveCurrency(value) {
    if (value !== this.selectedCurrency) {
      this.selectedCurrency = value;
      this.askIfConvertAmount();
      this.updateValues();
      (this.form.get('destinations') as FormArray).updateValueAndValidity();
    }
  }

  private askIfConvertAmount() {
    let validAmounts = 0;
    this.destControls.forEach((dest, i) => {
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

    const confirmationData: ConfirmationData = {
      text: this.translate.instant(validAmounts === 1 ? 'send.convert-confirmation' : 'send.convert-confirmation-plural', {from: fromText, to: toText}),
      headerText: 'confirmation.header-text',
      confirmButtonText: 'confirmation.confirm-button',
      cancelButtonText: 'confirmation.cancel-button',
    };

    showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
      if (confirmationResult) {
        this.convertAmounts();
      }
    });
  }

  private convertAmounts() {
    this.msgBarService.hide();

    let invalidValues = 0;
    let valuesWithPrecisionErrors = 0;
    this.destControls.forEach((dest, i) => {
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
          const recoveredValue = newValue.multipliedBy(this.price).decimalPlaces(SendFormComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
          if (!recoveredValue.isEqualTo(currentValue)) {
            valuesWithPrecisionErrors += 1;
          }

          dest.get('coins').setValue(newValue.toString());
        } else {
          const newValue = currentValue.multipliedBy(this.price).decimalPlaces(SendFormComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
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
      availableCoins = availableCoins.multipliedBy(this.price).decimalPlaces(SendFormComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
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
      availableCoins = availableCoins.decimalPlaces(SendFormComponent.MaxUsdDecimals, BigNumber.ROUND_FLOOR);
    }

    if (availableCoins.isLessThan(0)) {
      this.msgBarService.showError(this.translate.instant('send.no-coins-left'));
    } else {
      this.destControls[index].get('coins').setValue(availableCoins.toString());
    }
  }

  private updateValues() {
    if (!this.price) {
      this.values = null;

      return;
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

      if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
        const convertedValue = value.multipliedBy(this.price).decimalPlaces(2);

        this.totalCoins = this.totalCoins.plus(value);
        this.totalConvertedCoins = this.totalConvertedCoins.plus(convertedValue);

        this.values[i] = convertedValue.toNumber();
      } else {
        const convertedValue = value.dividedBy(this.price).decimalPlaces(this.blockchainService.currentMaxDecimals);

        this.totalCoins = this.totalCoins.plus(convertedValue);
        this.totalConvertedCoins = this.totalConvertedCoins.plus(value);

        this.values[i] = convertedValue.toNumber();
      }
    });

    if (!this.autoHours) {
      this.destControls.forEach((dest, i) => {
        const stringValue: string = dest.get('hours').value;
        const value = this.getAmount(stringValue, false);
        if (!value) {
          return;
        }

        this.totalHours = this.totalHours.plus(value);
      });
    }
  }

  private checkBeforeSending() {
    if (!this.form.valid || this.previewButton.isLoading() || this.sendButton.isLoading()) {
      return;
    }

    this.closeSyncCheckSubscription();
    this.syncCheckSubscription = this.blockchainService.synchronized.pipe(first()).subscribe(synchronized => {
      if (synchronized) {
        this.prepareTransaction();
      } else {
        this.showSynchronizingWarning();
      }
    });
  }

  private showSynchronizingWarning() {
    const confirmationData: ConfirmationData = {
      text: 'send.synchronizing-warning',
      headerText: 'confirmation.header-text',
      confirmButtonText: 'confirmation.confirm-button',
      cancelButtonText: 'confirmation.cancel-button',
    };

    showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
      if (confirmationResult) {
        this.prepareTransaction();
      }
    });
  }

  private prepareTransaction() {
    this.msgBarService.hide();
    this.previewButton.resetState();
    this.sendButton.resetState();

    const selectedSources = this.formSourceSelection.selectedSources;

    if (selectedSources.wallet.encrypted && !selectedSources.wallet.isHardware && !this.previewTx) {
      const config = new MatDialogConfig();
      config.data = {
        wallet: selectedSources.wallet,
      };

      this.dialog.open(PasswordDialogComponent, config).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.createTransaction(passwordDialog);
        });
    } else {
      if (!selectedSources.wallet.isHardware || this.previewTx) {
        this.createTransaction();
      } else {
        this.showBusy();
        this.processingSubscription = this.hwWalletService.checkIfCorrectHwConnected(selectedSources.wallet.addresses[0].address).subscribe(
          () => this.createTransaction(),
          err => this.showError(getHardwareWalletErrorMsg(this.translate, err)),
        );
      }
    }
  }

  addDestination() {
    const destinations = this.form.get('destinations') as FormArray;
    destinations.push(this.createDestinationFormGroup());
    this.updateValues();
  }

  removeDestination(index) {
    const destinations = this.form.get('destinations') as FormArray;
    destinations.removeAt(index);

    this.destinationSubscriptions[index].unsubscribe();
    this.destinationSubscriptions.splice(index, 1);
    this.destinationHoursSubscriptions[index].unsubscribe();
    this.destinationHoursSubscriptions.splice(index, 1);
    this.updateValues();
  }

  setShareValue(event) {
    this.autoShareValue = parseFloat(event.value).toFixed(2);
  }

  selectChangeAddress(event) {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.autoFocus = false;
    this.dialog.open(SelectAddressComponent, config).afterClosed().subscribe(response => {
      if (response) {
        this.form.get('changeAddress').setValue(response);
      }
    });
  }

  openMultipleDestinationsPopup() {

    let currentString = '';
    this.destControls.map((destControl, i) => {
      if ((destControl.get('address').value as string).trim().length > 0 ||
        (destControl.get('coins').value as string).trim().length > 0 ||
        (!this.autoHours && (destControl.get('hours').value as string).trim().length > 0)) {

          currentString += (destControl.get('address').value as string).replace(',', '');
          currentString += ', ' + (destControl.get('coins').value as string).replace(',', '');
          if (!this.autoHours) {
            currentString += ', ' + (destControl.get('hours').value as string).replace(',', '');
          }
          currentString += '\r\n';
      }
    });
    if (currentString.length > 0) {
      currentString = currentString.substr(0, currentString.length - 1);
    }

    const config = new MatDialogConfig();
    config.width = '566px';
    config.data = currentString;
    this.dialog.open(MultipleDestinationsDialogComponent, config).afterClosed().subscribe((response: string[][]) => {
      if (response) {
        while (this.destControls.length > 0) {
          (this.form.get('destinations') as FormArray).removeAt(0);
        }

        if (response.length > 0) {
          this.autoHours = response[0].length === 2;

          response.forEach((entry, i) => {
            this.addDestination();
            this.destControls[i].get('address').setValue(entry[0]);
            this.destControls[i].get('coins').setValue(entry[1]);
            if (!this.autoHours) {
              this.destControls[i].get('hours').setValue(entry[2]);
            }
          });
        } else {
          this.addDestination();
        }
      }
    });

    this.updateValues();
  }

  toggleOptions(event) {
    event.stopPropagation();
    event.preventDefault();

    this.autoOptions = !this.autoOptions;
  }

  setAutoHours(event) {
    this.autoHours = event.checked;
    this.form.get('destinations').updateValueAndValidity();
    this.updateValues();

    if (!this.autoHours) {
      this.autoOptions = false;
    }
  }

  private fillForm() {
    this.formSourceSelection.fill(this.formData);

    this.selectedCurrency = this.formData.form.currency;

    ['changeAddress', 'note'].forEach(name => {
      this.form.get(name).setValue(this.formData.form[name]);
    });

    for (let i = 0; i < this.formData.form.destinations.length - 1; i++) {
      this.addDestination();
    }

    this.destControls.forEach((destControl, i) => {
      ['address', 'hours'].forEach(name => {
        destControl.get(name).setValue(this.formData.form.destinations[i][name]);
      });
      destControl.get('coins').setValue(this.formData.form.destinations[i].originalAmount);
    });

    if (this.formData.form.hoursSelection.type === 'auto') {
      this.autoShareValue = this.formData.form.hoursSelection.share_factor;
      this.autoHours = true;
    } else {
      this.autoHours = false;
    }

    this.autoOptions = this.formData.form.autoOptions;

    this.updateValues();
  }

  get destControls() {
    return (this.form.get('destinations') as FormArray).controls;
  }

  private validateDestinations() {
    if (!this.form) {
      return { Required: true };
    }

    if (!this.formSourceSelection || !this.formSourceSelection.valid) {
      return { Invalid: true };
    }

    const invalidInput = this.destControls.find(control => {
      if (!control.get('address').value || (control.get('address').value as string).trim().length === 0) {
        return true;
      }

      const checkControls = ['coins'];

      if (!this.autoHours) {
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

    this.availableBalance = this.formSourceSelection.availableBalance;

    let destinationsCoins = new BigNumber(0);
    if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
      this.destControls.map(control => destinationsCoins = destinationsCoins.plus(control.value.coins));
    } else {
      this.updateValues();
      this.values.map(value => destinationsCoins = destinationsCoins.plus(value));
    }
    let destinationsHours = new BigNumber(0);
    if (!this.autoHours) {
      this.destControls.map(control => destinationsHours = destinationsHours.plus(control.value.hours));
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
        if (parts.length === 2 && parts[1].length > SendFormComponent.MaxUsdDecimals) {
          return null;
        }
      }
    } else if (name === 'hours') {
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

    this.destinationSubscriptions.push(group.get('coins').valueChanges.subscribe(value => {
      this.updateValues();
    }));

    this.destinationHoursSubscriptions.push(group.get('hours').valueChanges.subscribe(value => {
      this.updateValues();
    }));

    return group;
  }

  private showBusy() {
    if (this.previewTx) {
      this.previewButton.setLoading();
      this.sendButton.setDisabled();
    } else {
      this.sendButton.setLoading();
      this.previewButton.setDisabled();
    }
    this.busy = true;
    this.navbarService.disableSwitch();
  }

  private createTransaction(passwordDialog?: any) {
    this.showBusy();

    const selectedSources = this.formSourceSelection.selectedSources;

    const selectedAddresses = selectedSources.addresses && selectedSources.addresses.length > 0 ?
      selectedSources.addresses.map(addr => addr.address) : null;

    const selectedOutputs = selectedSources.unspentOutputs && selectedSources.unspentOutputs.length > 0 ?
      selectedSources.unspentOutputs.map(addr => addr.hash) : null;

      this.processingSubscription = this.walletService.createTransaction(
        selectedSources.wallet,
        selectedAddresses ? selectedAddresses : selectedSources.wallet.addresses.map(address => address.address),
        selectedOutputs,
        this.destinations,
        this.hoursSelection,
        this.form.get('changeAddress').value ? this.form.get('changeAddress').value : null,
        passwordDialog ? passwordDialog.password : null,
        this.previewTx,
      ).subscribe(transaction => {
        if (passwordDialog) {
          passwordDialog.close();
        }

        const note = this.form.value.note.trim();
        if (!this.previewTx) {
          this.processingSubscription = this.walletService.injectTransaction(transaction.encoded, note)
            .subscribe(noteSaved => {
              let showDone = true;
              if (note && !noteSaved) {
                this.msgBarService.showWarning(this.translate.instant('send.error-saving-note'));
                showDone = false;
              }

              this.showSuccess(showDone);
            }, error => this.showError(error));
        } else {
          let amount = new BigNumber('0');
          this.destinations.map(destination => amount = amount.plus(destination.coins));
          this.onFormSubmitted.emit({
            form: {
              wallet: selectedSources.wallet,
              addresses: selectedSources.addresses,
              changeAddress: this.form.get('changeAddress').value,
              destinations: this.destinations,
              hoursSelection: this.hoursSelection,
              autoOptions: this.autoOptions,
              allUnspentOutputs: this.formSourceSelection.unspentOutputsList,
              outputs: selectedSources.unspentOutputs,
              currency: this.selectedCurrency,
              note: note,
            },
            amount: amount,
            to: this.destinations.map(d => d.address),
            transaction,
          });
          this.busy = false;
          this.navbarService.enableSwitch();
        }
      }, error => {
        if (passwordDialog) {
          passwordDialog.error(error);
        }

        if (error && error.result) {
          this.showError(getHardwareWalletErrorMsg(this.translate, error));
        } else {
          this.showError(error);
        }
      });
  }

  private resetForm() {
    this.formSourceSelection.resetForm();
    this.form.get('changeAddress').setValue('');
    this.form.get('note').setValue('');

    while (this.destControls.length > 0) {
      (this.form.get('destinations') as FormArray).removeAt(0);
    }

    this.addDestination();

    this.autoHours = true;
    this.autoOptions = false;
    this.autoShareValue = '0.5';
    this.updateValues();
  }

  private get destinations() {
    return this.destControls.map((destControl, i) => {
      const destination = {
        address: ((destControl.get('address').value) as string).trim(),
        coins: ((this.selectedCurrency === DoubleButtonActive.LeftButton ? destControl.get('coins').value : this.values[i].toString()) as string).trim(),
        originalAmount: destControl.get('coins').value,
      };

      if (!this.autoHours) {
        destination['hours'] = destControl.get('hours').value;
      }

      return destination;
    });
  }

  private get hoursSelection() {
    let hoursSelection = {
      type: 'manual',
    };

    if (this.autoHours) {
      hoursSelection = <any> {
        type: 'auto',
        mode: 'share',
        share_factor: this.autoShareValue,
      };
    }

    return hoursSelection;
  }

  private closeSyncCheckSubscription() {
    if (this.syncCheckSubscription) {
      this.syncCheckSubscription.unsubscribe();
    }
  }

  private showSuccess(showDone: boolean) {
    this.busy = false;
    this.navbarService.enableSwitch();
    this.resetForm();

    if (showDone) {
      this.msgBarService.showDone('send.sent');
      this.sendButton.resetState();
    } else {
      this.sendButton.setSuccess();
      setTimeout(() => {
        this.sendButton.resetState();
      }, 3000);
    }
  }

  private showError(error) {
    this.busy = false;
    this.msgBarService.showError(error);
    this.navbarService.enableSwitch();
    this.previewButton.resetState().setEnabled();
    this.sendButton.resetState().setEnabled();
  }
}
