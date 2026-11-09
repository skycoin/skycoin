import { ChangeDetectionStrategy, ChangeDetectorRef, Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { UntypedFormArray, UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { MatDialogConfig } from '@angular/material/dialog';
import { Subscription, throwError, concat } from 'rxjs';
import { BigNumber } from 'bignumber.js';
import { retryWhen, delay, take, first, debounceTime } from 'rxjs';

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
import { BulkSendComponent, BulkDestination } from './bulk-send/bulk-send.component';
import { OfflineTxComponent, OfflineTxMode } from '../offline-tx/offline-tx.component';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { maxNoteChars } from '../../../../utils/tx-notes';

@Component({
    selector: 'app-send-form-advanced',
    templateUrl: './send-form-advanced.component.html',
    styleUrls: ['./send-form-advanced.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class SendFormAdvancedComponent implements OnInit, OnDestroy {
  @ViewChild('button') button!: ButtonComponent;
  @ViewChild('unsignedButton') unsignedButton!: ButtonComponent;
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
  // Builds an unsigned transaction from typed addresses, with no wallet loaded.
  manualUnsigned = false;
  modeButton = DoubleButtonActive.LeftButton;
  manualAddresses: string[] = [];
  maxNoteChars = maxNoteChars;

  private subscriptionsGroup: Subscription[] = [];
  private getOutputsSubscriptions!: Subscription;
  private unlockSubscription!: Subscription;
  private destinationSubscriptions: Subscription[] = [];
  // An encrypted server wallet's password, held only until the transaction is signed.
  private password: string | undefined;
  // The change address filled in for the selected wallet. A different value in
  // the field is the user's own and is kept when the wallet changes.
  private defaultChange = '';
  private defaultChangeSubscription?: Subscription;

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
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit() {
    this.navbarService.showSwitch('send.simple', 'send.advanced');

    this.form = this.formBuilder.group({
      wallet: ['', Validators.required],
      addresses: [null],
      outputs: [null],
      changeAddress: [''],
      note: [''],
      manualAddresses: [''],
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
            this.changeDetectorRef.markForCheck();
          },
          () => {
            this.loadingUnspentOutputs = false;
            this.changeDetectorRef.markForCheck();
          },
        );

      this.addresses = wallet.addresses.filter(addr => addr.balance!.isGreaterThan(0));
      this.form.get('addresses')!.setValue(null);
      this.form.get('outputs')!.setValue(null);
      this.fillDefaultChange(wallet);

      this.updateAvailableBalance();
      this.form.get('destinations')!.updateValueAndValidity();
      this.changeDetectorRef.markForCheck();
    }));

    this.subscriptionsGroup.push(this.form.get('manualAddresses')!.valueChanges.pipe(debounceTime(500)).subscribe(value => {
      this.loadManualOutputs(value);
      this.changeDetectorRef.markForCheck();
    }));

    this.subscriptionsGroup.push(this.form.get('addresses')!.valueChanges.subscribe(() => {
      this.form.get('outputs')!.setValue(null);
      this.unspentOutputs = this.filterUnspentOutputs();

      this.updateAvailableBalance();
      this.form.get('destinations')!.updateValueAndValidity();
      this.changeDetectorRef.markForCheck();
    }));

    this.subscriptionsGroup.push(this.form.get('outputs')!.valueChanges.subscribe(() => {
      this.updateAvailableBalance();
      this.form.get('destinations')!.updateValueAndValidity();
      this.changeDetectorRef.markForCheck();
    }));

    this.subscriptionsGroup.push(this.coinService.currentCoin
      .subscribe((coin: BaseCoin) => {
        this.resetForm();
        this.currentCoin = coin;
        this.changeDetectorRef.markForCheck();
      })
    );

    this.subscriptionsGroup.push(this.priceService.price.subscribe(price => {
      this.price = price;
      this.updateValues();
      this.changeDetectorRef.markForCheck();
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
    if (this.defaultChangeSubscription) {
      this.defaultChangeSubscription.unsubscribe();
    }
  }

  preview() {
    this.previewTx = true;
    this.unlockAndSend();
  }

  send() {
    this.previewTx = false;
    this.unlockAndSend();
  }

  /** Builds the transaction without signing it, for a wallet whose keys are elsewhere. */
  createUnsigned() {
    if (!this.form.valid || this.unsignedButton.isLoading()) {
      return;
    }

    this.msgBarService.hide();
    this.unsignedButton.setLoading();
    this.spendingService.createTransaction(
      this.manualUnsigned ? null! : this.form.get('wallet')!.value,
      this.manualUnsigned ? this.manualAddresses : this.selectedAddresses,
      this.selectedOutputs,
      this.destinations,
      this.hoursSelection,
      this.form.get('changeAddress')!.value ? this.form.get('changeAddress')!.value : null,
      undefined,
      true,
    ).subscribe(transaction => {
      this.unsignedButton.resetState();
      OfflineTxComponent.open(this.dialog, {
        mode: OfflineTxMode.Show,
        rawTx: (transaction as PreviewTransaction).encoded,
        unsigned: true,
      });
      this.changeDetectorRef.markForCheck();
    }, error => {
      this.unsignedButton.resetState();
      this.msgBarService.showError(error.message);
      this.changeDetectorRef.markForCheck();
    });
  }

  /**
   * Switches between sending from a wallet and building an unsigned transaction
   * from typed addresses, after the user confirms the latter.
   */
  changeFormType(value: DoubleButtonActive) {
    this.modeButton = value;
    if (value === DoubleButtonActive.LeftButton) {
      this.setManualUnsigned(false);

      return;
    }

    const confirmationData: ConfirmationData = {
      text: 'send.unsigned-confirmation',
      headerText: 'confirmation.header-text',
      confirmButtonText: 'confirmation.confirm-button',
      cancelButtonText: 'confirmation.cancel-button',
    };
    showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmed => {
      if (confirmed) {
        this.setManualUnsigned(true);
      } else {
        this.modeButton = DoubleButtonActive.LeftButton;
      }
      this.changeDetectorRef.markForCheck();
    });
  }

  private setManualUnsigned(manual: boolean) {
    this.manualUnsigned = manual;
    this.resetForm();
    this.form.get('wallet')!.setValidators(manual ? null : Validators.required);
    this.form.get('manualAddresses')!.setValidators(manual ? Validators.required : null);
    this.form.get('wallet')!.updateValueAndValidity({ emitEvent: false });
    this.form.get('manualAddresses')!.updateValueAndValidity({ emitEvent: false });
  }

  /** Loads the unspent outputs of the comma separated addresses typed in the manual mode. */
  private loadManualOutputs(value: string) {
    const addresses = (value || '').split(',').map(a => a.trim()).filter(a => a.length > 0);
    if (addresses.join(',') === this.manualAddresses.join(',')) {
      return;
    }

    this.manualAddresses = addresses;
    this.closeGetOutputsSubscriptions();
    this.allUnspentOutputs = [];
    this.unspentOutputs = [];
    this.form.get('outputs')!.setValue(null);
    if (addresses.length === 0) {
      this.updateAvailableBalance();
      this.form.get('destinations')!.updateValueAndValidity();

      return;
    }

    this.loadingUnspentOutputs = true;
    this.getOutputsSubscriptions = this.spendingService.getAddressesUnspentOutputs(addresses).subscribe(
      result => {
        this.loadingUnspentOutputs = false;
        this.allUnspentOutputs = result;
        this.unspentOutputs = this.filterUnspentOutputs();
        this.updateAvailableBalance();
        this.form.get('destinations')!.updateValueAndValidity();
        this.changeDetectorRef.markForCheck();
      },
      () => {
        this.loadingUnspentOutputs = false;
        this.msgBarService.showError('send.manual-addresses-error');
        this.changeDetectorRef.markForCheck();
      },
    );
  }

  changeActiveCurrency(value: any) {
    this.selectedCurrency = value;
    this.updateValues();
    (this.form.get('destinations') as UntypedFormArray).updateValueAndValidity();
  }

  /** Fills a destination with what the selected sources hold, less the other destinations. */
  assignAll(index: number) {
    this.msgBarService.hide();
    let available: BigNumber = this.availableCoins;
    if (this.selectedCurrency === DoubleButtonActive.RightButton) {
      available = available.multipliedBy(this.price).decimalPlaces(2, BigNumber.ROUND_FLOOR);
    }
    this.destControls.forEach((dest, i) => {
      const value = String(dest.get('coins')!.value ?? '').replace(/,/g, '').trim();
      if (i !== index && /^\d*\.?\d+$/.test(value)) {
        available = available.minus(value);
      }
    });
    const decimals = this.selectedCurrency === DoubleButtonActive.RightButton ? 2 : this.blockchainService.currentMaxDecimals;
    available = available.decimalPlaces(decimals, BigNumber.ROUND_FLOOR);

    if (available.isLessThanOrEqualTo(0)) {
      this.msgBarService.showError('send.no-coins-left-error');
    } else {
      this.destControls[index].get('coins')!.setValue(available.toString());
    }
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

    this.password = undefined;
    if (this.walletService.isServerWallet(wallet) && wallet.encrypted) {
      this.removeUnlockSubscription();
      this.unlockSubscription = PasswordDialogComponent.ask(this.dialog, { text: wallet.label }).subscribe(password => {
        if (password) {
          this.password = password;
          this.checkBeforeSending();
        }
        this.changeDetectorRef.markForCheck();
      });
    } else if (wallet.isHardware || this.walletService.isServerWallet(wallet)) {
      this.checkBeforeSending();
    } else if (!wallet.seed) {
      this.removeUnlockSubscription();

      this.unlockSubscription = openUnlockWalletModal(wallet, this.dialog).componentInstance
        .onWalletUnlocked.pipe(first()).subscribe(() => {
          this.checkBeforeSending();
          this.changeDetectorRef.markForCheck();
        });
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
      this.changeDetectorRef.markForCheck();
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
      this.changeDetectorRef.markForCheck();
    });
  }

  addDestination() {
    const destinations = this.form.get('destinations') as UntypedFormArray;
    destinations.push(this.createDestinationFormGroup());
    this.updateValues();
  }

  /** Replaces the destinations with ones typed as text, one per line. */
  openBulkSend() {
    const current = this.destControls
      .map(c => c.value)
      .filter(v => v.address || v.coins)
      .map(v => [v.address, v.coins].concat(this.autoHours ? [] : [v.hours]).join(','))
      .join('\n');

    BulkSendComponent.ask(this.dialog, current).subscribe(destinations => {
      if (destinations) {
        this.setDestinations(destinations);
      }
      this.changeDetectorRef.markForCheck();
    });
  }

  private setDestinations(destinations: BulkDestination[]) {
    while (this.destControls.length > 0) {
      this.removeDestination(0);
    }
    this.selectedCurrency = DoubleButtonActive.LeftButton;
    this.autoHours = destinations[0].hours === undefined;
    if (!this.autoHours) {
      this.autoOptions = false;
    }

    destinations.forEach((d, i) => {
      this.addDestination();
      const group = this.destControls[i];
      group.get('address')!.setValue(d.address);
      group.get('coins')!.setValue(d.coins);
      group.get('hours')!.setValue(d.hours || '');
    });
    this.form.get('destinations')!.updateValueAndValidity();
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

  /**
   * Fills the change address with the selected wallet's next unused address,
   * unless the field holds an address the user chose. Empty when the wallet
   * cannot say, as a node wallet, where the node picks a new address itself.
   */
  private fillDefaultChange(wallet: Wallet) {
    const control = this.form.get('changeAddress')!;
    const usersOwn = () => !!control.value && control.value !== this.defaultChange;
    if (usersOwn()) {
      return;
    }
    if (this.defaultChangeSubscription) {
      this.defaultChangeSubscription.unsubscribe();
    }
    this.defaultChangeSubscription = this.walletService.changeAddressFor(wallet).subscribe({
      next: address => {
        if (usersOwn()) {
          return;
        }
        this.defaultChange = address || '';
        control.setValue(this.defaultChange);
        this.changeDetectorRef.markForCheck();
      },
      error: () => undefined,
    });
  }

  selectChangeAddress(event: any) {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.autoFocus = false;
    this.dialog.open(SelectAddressComponent, config).afterClosed().subscribe(response => {
      if (response) {
        this.form.get('changeAddress')!.setValue(response);
      }
      this.changeDetectorRef.markForCheck();
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

    ['wallet', 'addresses', 'changeAddress', 'note'].forEach(name => {
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
      this.changeDetectorRef.markForCheck();
    }));

    return group;
  }

  private get selectedAddresses(): string[] | null {
    const addresses: Address[] | null = this.form.get('addresses')!.value;

    return addresses && addresses.length > 0 ? addresses.map(a => a.address) : null;
  }

  private get selectedOutputs(): string[] | null {
    const outputs: UnspentOutput[] | null = this.form.get('outputs')!.value;

    return outputs && outputs.length > 0 ? outputs.map(o => o.hash) : null;
  }

  private createTransaction() {
    if (this.previewTx) {
      this.button.setLoading();
    } else {
      this.button.setLoading();
    }

    const password = this.password;
    this.password = undefined;

    this.spendingService.createTransaction(
      this.form.get('wallet')!.value,
      this.selectedAddresses,
      this.selectedOutputs,
      this.destinations,
      this.hoursSelection,
      this.form.get('changeAddress')!.value ? this.form.get('changeAddress')!.value : null,
      password,
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
            note: this.form.get('note')!.value,
          },
          note: this.form.get('note')!.value,
          amount: amount,
          to: this.destinations.map(d => d.address),
          transaction,
        });
        this.changeDetectorRef.markForCheck();

        return undefined;
      })
      .then(() => {
        this.button.setSuccess();
        this.resetForm();

        setTimeout(() => {
          this.button.resetState();
          this.changeDetectorRef.markForCheck();
        }, 3000);
        this.changeDetectorRef.markForCheck();
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
    this.defaultChange = '';
    this.form.get('note')!.setValue('');
    this.form.get('manualAddresses')!.setValue('', { emitEvent: false });
    this.manualAddresses = [];
    this.closeGetOutputsSubscriptions();
    this.allUnspentOutputs = [];
    this.unspentOutputs = [];

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

    if (!this.form.get('wallet')!.value && !this.manualUnsigned) {
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
    } else if (this.manualUnsigned) {
      this.allUnspentOutputs.forEach(output => {
        this.availableCoins = this.availableCoins.plus(output.coins);
        this.availableHours = this.availableHours.plus(output.calculated_hours);
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
