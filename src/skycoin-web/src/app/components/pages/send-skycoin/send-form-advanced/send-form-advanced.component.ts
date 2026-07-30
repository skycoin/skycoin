import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild, ChangeDetectionStrategy } from '@angular/core';
import { UntypedFormArray, UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { MatDialogConfig } from '@angular/material/dialog';
import { Subscription, throwError, concat } from 'rxjs';
import { BigNumber } from 'bignumber.js';
import { retryWhen, delay, take, first } from 'rxjs';

import { ButtonComponent } from '../../../layout/button/button.component';
import { NavBarService } from '../../../../services/nav-bar.service';
import { SelectAddressComponent } from './select-address/select-address';
import { Output as UnspentOutput, Wallet, Address, ConfirmationData, PreviewTransaction } from '../../../../app.datatypes';
import { BlockchainService } from '../../../../services/blockchain.service';
import { openUnlockWalletModal, showConfirmationModal } from '../../../../utils';
import { SpendingService, HoursSelection, HoursSelectionTypes, Destination } from '../../../../services/wallet/spending.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { BaseCoin } from '../../../../coins/basecoin';
import { CoinService } from '../../../../services/coin.service';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { PriceService } from '../../../../services/price.service';
import { SendFormComponent } from '../send-form/send-form.component';
import { MsgBarService } from '../../../../services/msg-bar.service';

@Component({
    selector: 'app-send-form-advanced',
    templateUrl: './send-form-advanced.component.html',
    styleUrls: ['./send-form-advanced.component.scss'],
    changeDetection: ChangeDetectionStrategy.Eager,
    standalone: false
})
export class SendFormAdvancedComponent implements OnInit, OnDestroy {
  @ViewChild('button') button!: ButtonComponent;
  @Input() formData: any;
  @Output() onFormSubmitted = new EventEmitter<any>();

  form!: UntypedFormGroup;
  wallet!: Wallet | null;
  addresses: Address[] = [];
  allUnspentOutputs: UnspentOutput[] = [];
  unspentOutputs: UnspentOutput[] = [];
  loadingUnspentOutputs = false;
  availableCoins = new BigNumber(0);
  availableHours = new BigNumber(0);
  minimumFee = new BigNumber(0);
  autoHours = true;
  autoOptions = false;
  autoShareValue = '0.5';
  previewTx!: boolean;
  currentCoin!: BaseCoin;
  doubleButtonActive = DoubleButtonActive;
  selectedCurrency = DoubleButtonActive.LeftButton;
  values!: number[];
  price!: number;

  private subscriptionsGroup: Subscription[] = [];
  private getOutputsSubscriptions!: Subscription;
  private unlockSubscription!: Subscription;
  private destinationSubscriptions: Subscription[] = [];

  constructor(
    public walletService: WalletService,
    private spendingService: SpendingService,
    private formBuilder: UntypedFormBuilder,
    private dialog: CustomMatDialogService,
    private navbarService: NavBarService,
    public blockchainService: BlockchainService,
    private coinService: CoinService,
    private priceService: PriceService,
    private msgBarService: MsgBarService,
  ) { }

  ngOnInit() {
    this.navbarService.showSwitch('send.simple', 'send.advanced');

    this.form = this.formBuilder.group({
      wallet: ['', Validators.required],
      addresses: [null],
      outputs: [null],
      changeAddress: [''],
      destinations: this.formBuilder.array(
        [this.createDestinationFormGroup()],
        this.validateDestinations.bind(this),
      ),
    });

    this.subscriptionsGroup.push(this.form.get('wallet')!.valueChanges.subscribe((wallet: Wallet) => {
      this.wallet = wallet;

      this.closeGetOutputsSubscriptions();
      this.allUnspentOutputs = [];
      this.unspentOutputs = [];
      this.loadingUnspentOutputs = true;

      this.getOutputsSubscriptions = this.spendingService.getWalletUnspentOutputs(wallet).pipe(
        retryWhen(errors => errors.pipe(delay(1000), take(10), o => concat(o, throwError(() => '')))))
        .subscribe(
          result => {
            this.loadingUnspentOutputs = false;
            this.allUnspentOutputs = result;
            this.unspentOutputs = this.filterUnspentOutputs();
          },
          () => this.loadingUnspentOutputs = false,
        );

      this.addresses = wallet.addresses.filter(addr => addr.balance!.isGreaterThan(0));
      this.form.get('addresses')!.setValue(null);
      this.form.get('outputs')!.setValue(null);

      this.updateAvailableBalance();
      this.form.get('destinations')!.updateValueAndValidity();
    }));

    this.subscriptionsGroup.push(this.form.get('addresses')!.valueChanges.subscribe(() => {
      this.form.get('outputs')!.setValue(null);
      this.unspentOutputs = this.filterUnspentOutputs();

      this.updateAvailableBalance();
      this.form.get('destinations')!.updateValueAndValidity();
    }));

    this.subscriptionsGroup.push(this.form.get('outputs')!.valueChanges.subscribe(() => {
      this.updateAvailableBalance();
      this.form.get('destinations')!.updateValueAndValidity();
    }));

    this.subscriptionsGroup.push(this.coinService.currentCoin
      .subscribe((coin: BaseCoin) => {
        this.resetForm();
        this.currentCoin = coin;
      })
    );

    this.subscriptionsGroup.push(this.priceService.price.subscribe(price => {
      this.price = price;
      this.updateValues();
    }));

    if (this.formData) {
      this.fillForm();
    }
  }

  ngOnDestroy() {
    this.closeGetOutputsSubscriptions();
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
    this.navbarService.hideSwitch();
    this.msgBarService.hide();
    this.destinationSubscriptions.forEach(s => s.unsubscribe());
  }

  preview() {
    this.previewTx = true;
    this.unlockAndSend();
  }

  send() {
    this.previewTx = false;
    this.unlockAndSend();
  }

  changeActiveCurrency(value: any) {
    this.selectedCurrency = value;
    this.updateValues();
    (this.form.get('destinations') as UntypedFormArray).updateValueAndValidity();
  }

  private updateValues() {
    if (!this.price) {
      this.values = null!;

      return;
    }

    this.values = [];

    this.destControls.forEach((dest, i) => {
      const value = dest.get('coins')!.value !== undefined ? dest.get('coins')!.value.replace(' ', '=') : '';

      if (isNaN(value) || value.trim() === '' || parseFloat(value) <= 0 || value * 1 === 0) {
        this.values[i] = -1;

        return;
      }

      const parts = value.split('.');
      if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
        if (parts.length === 2 && parts[1].length > this.blockchainService.currentMaxDecimals) {
          this.values[i] = -1;

          return;
        }
      } else {
        if (parts.length === 2 && parts[1].length > SendFormComponent.MaxUsdDecimal) {
          this.values[i] = -1;

          return;
        }
      }

      if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
        this.values[i] = new BigNumber(value).multipliedBy(this.price).decimalPlaces(2).toNumber();
      } else {
        this.values[i] = new BigNumber(value).dividedBy(this.price).decimalPlaces(this.blockchainService.currentMaxDecimals).toNumber();
      }
    });
  }

  unlockAndSend() {
    if (!this.form.valid || this.button.isLoading()) {
      return;
    }

    this.msgBarService.hide();
    this.button.resetState();

    const wallet = this.form.value.wallet;

    if (!wallet.seed) {
      this.removeUnlockSubscription();

      this.unlockSubscription = openUnlockWalletModal(wallet, this.dialog).componentInstance
        .onWalletUnlocked.pipe(first()).subscribe(() => this.checkBeforeSending());
    } else {
      this.checkBeforeSending();
    }
  }

  private checkBeforeSending() {
    this.blockchainService.synchronized.pipe(first()).subscribe(synchronized => {
      if (synchronized) {
        this.createTransaction();
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
        this.createTransaction();
      }
    });
  }

  addDestination() {
    const destinations = this.form.get('destinations') as UntypedFormArray;
    destinations.push(this.createDestinationFormGroup());
    this.updateValues();
  }

  removeDestination(index: any) {
    const destinations = this.form.get('destinations') as UntypedFormArray;
    destinations.removeAt(index);

    this.destinationSubscriptions[index].unsubscribe();
    this.destinationSubscriptions.splice(index, 1);
    this.updateValues();
  }

  setShareValue(value: number) {
    this.autoShareValue = value.toFixed(2);
  }

  selectChangeAddress(event: any) {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.autoFocus = false;
    this.dialog.open(SelectAddressComponent, config).afterClosed().subscribe(response => {
      if (response) {
        this.form.get('changeAddress')!.setValue(response);
      }
    });
  }

  toggleOptions(event: any) {
    event.stopPropagation();
    event.preventDefault();

    this.autoOptions = !this.autoOptions;
  }

  setAutoHours(event: any) {
    this.autoHours = event.checked;
    this.form.get('destinations')!.updateValueAndValidity();

    if (!this.autoHours) {
      this.autoOptions = false;
    }
  }

  private fillForm() {
    this.addresses = this.formData.form.wallet.addresses;

    ['wallet', 'addresses', 'changeAddress'].forEach(name => {
      this.form.get(name)!.setValue(this.formData.form[name]);
    });

    for (let i = 0; i < this.formData.form.destinations.length - 1; i++) {
      this.addDestination();
    }

    this.destControls.forEach((destControl, i) => {
      ['address', 'hours'].forEach(name => {
        destControl.get(name)!.setValue(this.formData.form.destinations[i][name]);
      });
      destControl.get('coins')!.setValue(this.formData.form.destinations[i].originalAmount);
    });

    if (this.formData.form.hoursSelection.type === HoursSelectionTypes.Auto) {
      this.autoShareValue = this.formData.form.hoursSelection.ShareFactor;
      this.autoHours = true;
    } else {
      this.autoHours = false;
    }

    this.autoOptions = this.formData.form.autoOptions;

    if (this.formData.form.allUnspentOutputs) {
      this.closeGetOutputsSubscriptions();

      this.allUnspentOutputs = this.formData.form.allUnspentOutputs;
      this.unspentOutputs = this.filterUnspentOutputs();

      this.form.get('outputs')!.setValue(this.formData.form.outputs);
    }

    this.selectedCurrency = this.formData.form.currency;
  }

  addressCompare(a: any, b: any) {
    return a && b && a.address === b.address;
  }

  outputCompare(a: any, b: any) {
    return a && b && a.hash === b.hash;
  }

  get destControls() {
    return (this.form.get('destinations') as UntypedFormArray).controls;
  }

  private validateDestinations() {
    if (!this.form) {
      return { Required: true };
    }

    const invalidInput = this.destControls.find(control => {
      const checkControls = ['coins'];

      if (!this.autoHours) {
        checkControls.push('hours');
      }

      return checkControls.map(name => {
        const value = control.get(name)!.value !== undefined
          ? String(control.get(name)!.value).replace(' ', '=')
          : '';

        if (isNaN((value as any)) || value.trim() === '') {
          return true;
        }

        if (parseFloat(value) <= 0) {
          return true;
        }

        if (name === 'coins') {
          const parts = value.split('.');

          if (this.selectedCurrency === DoubleButtonActive.LeftButton) {
            if (parts.length === 2 && parts[1].length > this.blockchainService.currentMaxDecimals) {
              return true;
            }
          } else {
            if (parts.length === 2 && parts[1].length > SendFormComponent.MaxUsdDecimal) {
              return true;
            }
          }
        } else if (name === 'hours') {
          if (Number(value) < 1 || parseInt(value, 10) !== parseFloat(value)) {
            return true;
          }
        }

        return false;
      }).find(e => e === true);
    });

    if (invalidInput) {
      return { Invalid: true };
    }

    this.updateAvailableBalance();

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

    if (destinationsCoins.isGreaterThan(this.availableCoins) || destinationsHours.isGreaterThan(this.availableHours)) {
      return { Invalid: true };
    }

    return null;
  }

  private createDestinationFormGroup() {
    const group = this.formBuilder.group({
      address: '',
      coins: '',
      hours: '',
    });

    this.destinationSubscriptions.push(group.get('coins')!.valueChanges.subscribe(value => {
      this.updateValues();
    }));

    return group;
  }

  private createTransaction() {
    if (this.previewTx) {
      this.button.setLoading();
    } else {
      this.button.setLoading();
    }

    const selectedAddresses = this.form.get('addresses')!.value && (this.form.get('addresses')!.value as Address[]).length > 0 ?
      this.form.get('addresses')!.value.map((addr: any) => addr.address) : null;

    const selectedOutputs = this.form.get('outputs')!.value && (this.form.get('outputs')!.value as UnspentOutput[]).length > 0 ?
      this.form.get('outputs')!.value.map((addr: any) => addr.hash) : null;

    this.spendingService.createTransaction(
      this.form.get('wallet')!.value,
      selectedAddresses,
      selectedOutputs,
      this.destinations,
      this.hoursSelection,
      this.form.get('changeAddress')!.value ? this.form.get('changeAddress')!.value : null,
    )
      .toPromise()
      .then(transaction => {
        if (!this.previewTx) {
          return this.spendingService.injectTransaction((transaction as PreviewTransaction).encoded).toPromise();
        }

        let amount = new BigNumber('0');
        this.destinations.map(destination => amount = amount.plus(destination.coins));

        this.onFormSubmitted.emit({
          form: {
            wallet: this.form.get('wallet')!.value,
            addresses: this.form.get('addresses')!.value,
            changeAddress: this.form.get('changeAddress')!.value,
            destinations: this.destinations,
            hoursSelection: this.hoursSelection,
            autoOptions: this.autoOptions,
            allUnspentOutputs: this.loadingUnspentOutputs ? null : this.allUnspentOutputs,
            outputs: this.form.get('outputs')!.value,
            currency: this.selectedCurrency,
          },
          amount: amount,
          to: this.destinations.map(d => d.address),
          transaction,
        });
      })
      .then(() => {
        this.button.setSuccess();
        this.resetForm();

        setTimeout(() => {
          this.button.resetState();
        }, 3000);
      })
      .catch(error => {
        this.msgBarService.showError(error.message);
        this.button.resetState();
      });
  }

  private resetForm() {
    this.form.get('wallet')!.setValue('', { emitEvent: false });
    this.form.get('addresses')!.setValue(null);
    this.form.get('outputs')!.setValue(null);
    this.form.get('changeAddress')!.setValue('');

    while (this.destControls.length > 0) {
      (this.form.get('destinations') as UntypedFormArray).removeAt(0);
    }

    this.addDestination();

    this.wallet = null;
    this.autoHours = true;
    this.autoOptions = false;
    this.autoShareValue = '0.5';
  }

  private get destinations(): Destination[] {
    return this.destControls.map((destControl, i) => {
      const destination = {
        address: destControl.get('address')!.value,
        coins: this.selectedCurrency === DoubleButtonActive.LeftButton ? new BigNumber(destControl.get('coins')!.value) : new BigNumber(this.values[i].toString()),
        originalAmount: destControl.get('coins')!.value,
      };

      if (!this.autoHours) {
        (destination as any)['hours'] = new BigNumber(destControl.get('hours')!.value);
      }

      return destination;
    });
  }

  private get hoursSelection(): HoursSelection {
    let hoursSelection = {
      type: HoursSelectionTypes.Manual,
    };

    if (this.autoHours) {
      hoursSelection = <any> {
        type: HoursSelectionTypes.Auto,
        ShareFactor: this.autoShareValue,
      };
    }

    return hoursSelection;
  }

  private updateAvailableBalance() {
    this.availableCoins = new BigNumber(0);
    this.availableHours = new BigNumber(0);
    this.minimumFee = new BigNumber(0);

    if (!this.form.get('wallet')!.value) {
      return;
    }

    const outputs: UnspentOutput[] = this.form.get('outputs')!.value;
    const addresses: Address[] = this.form.get('addresses')!.value;

    if (outputs && outputs.length > 0) {
      outputs.map(control => {
        this.availableCoins = this.availableCoins.plus(control.coins);
        this.availableHours = this.availableHours.plus(control.calculated_hours);
      });
    } else if (addresses && addresses.length > 0) {
      addresses.map(control => {
        this.availableCoins = this.availableCoins.plus(control.balance!);
        this.availableHours = this.availableHours.plus(control.hours!);
      });
    } else {
      const wallet: Wallet = this.form.get('wallet')!.value;
      this.availableCoins = wallet.balance!;
      this.availableHours = wallet.hours!;
    }

    const unburnedHoursRatio = new BigNumber(1).minus(new BigNumber(1).dividedBy(this.blockchainService.burnRate));
    const sendableHours = this.availableHours.multipliedBy(unburnedHoursRatio).decimalPlaces(0, BigNumber.ROUND_FLOOR);
    this.minimumFee = this.availableHours.minus(sendableHours);
    this.availableHours = sendableHours;
  }

  private filterUnspentOutputs(): UnspentOutput[] {
    if (this.allUnspentOutputs.length === 0) {
      return [];
    } else if (!this.form.get('addresses')!.value || (this.form.get('addresses')!.value as Address[]).length === 0) {
      return this.allUnspentOutputs;
    } else {
      return this.allUnspentOutputs.filter(out => (this.form.get('addresses')!.value as Address[]).some(addr => addr.address === out.address));
    }
  }

  private closeGetOutputsSubscriptions() {
    this.loadingUnspentOutputs = false;

    if (this.getOutputsSubscriptions) {
      this.getOutputsSubscriptions.unsubscribe();
    }
  }

  private removeUnlockSubscription() {
    if (this.unlockSubscription) {
      this.unlockSubscription.unsubscribe();
    }
  }
}
