import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild, ChangeDetectorRef } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialog, MatSnackBar, MatDialogConfig } from '@angular/material';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { ButtonComponent } from '../../../layout/button/button.component';
import { showSnackbarError, getHardwareWalletErrorMsg } from '../../../../utils/errors';
import { Subscription, ISubscription } from 'rxjs/Subscription';
import { NavBarService } from '../../../../services/nav-bar.service';
import { SelectAddressComponent } from './select-address/select-address';
import { BigNumber } from 'bignumber.js';
import { Output as UnspentOutput, Wallet, Address, ConfirmationData } from '../../../../app.datatypes';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/retryWhen';
import 'rxjs/add/operator/concat';
import { BlockchainService } from '../../../../services/blockchain.service';
import { showConfirmationModal } from '../../../../utils';
import { AppService } from '../../../../services/app.service';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { TranslateService } from '@ngx-translate/core';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { PriceService } from '../../../../services/price.service';
import { SendFormComponent } from '../send-form/send-form.component';
import { ChangeNoteComponent } from '../send-preview/transaction-info/change-note/change-note.component';

@Component({
  selector: 'app-send-form-advanced',
  templateUrl: './send-form-advanced.component.html',
  styleUrls: ['./send-form-advanced.component.scss'],
})
export class SendFormAdvancedComponent implements OnInit, OnDestroy {
  @ViewChild('previewButton') previewButton: ButtonComponent;
  @ViewChild('sendButton') sendButton: ButtonComponent;
  @Input() formData: any;
  @Output() onFormSubmitted = new EventEmitter<any>();

  maxNoteChars = ChangeNoteComponent.MAX_NOTE_CHARS;
  form: FormGroup;
  wallet: Wallet;
  addresses = [];
  allUnspentOutputs: UnspentOutput[] = [];
  unspentOutputs: UnspentOutput[] = [];
  loadingUnspentOutputs = false;
  availableCoins = new BigNumber(0);
  availableHours = new BigNumber(0);
  minimumFee = new BigNumber(0);
  autoHours = true;
  autoOptions = false;
  autoShareValue = '0.5';
  previewTx: boolean;
  busy = false;
  doubleButtonActive = DoubleButtonActive;
  selectedCurrency = DoubleButtonActive.LeftButton;
  values: number[];
  price: number;

  private subscriptions: Subscription;
  private getOutputsSubscriptions: ISubscription;
  private destinationSubscriptions: ISubscription[] = [];
  private syncCheckSubscription: ISubscription;
  private processingSubscription: ISubscription;

  constructor(
    public blockchainService: BlockchainService,
    public walletService: WalletService,
    private appService: AppService,
    private formBuilder: FormBuilder,
    private dialog: MatDialog,
    private snackbar: MatSnackBar,
    private navbarService: NavBarService,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
    private priceService: PriceService,
    private changeDetector: ChangeDetectorRef,
  ) { }

  ngOnInit() {
    this.navbarService.showSwitch('send.simple', 'send.advanced', DoubleButtonActive.RightButton);

    this.form = this.formBuilder.group({
      wallet: ['', Validators.required],
      addresses: [null],
      outputs: [null],
      changeAddress: [''],
      destinations: this.formBuilder.array(
        [this.createDestinationFormGroup()],
        this.validateDestinations.bind(this),
      ),
      note: [''],
    });

    this.subscriptions = this.form.get('wallet').valueChanges.subscribe(wallet => {
      this.wallet = wallet;

      this.closeGetOutputsSubscriptions();
      this.allUnspentOutputs = [];
      this.unspentOutputs = [];
      this.loadingUnspentOutputs = true;

      this.getOutputsSubscriptions = this.walletService.getWalletUnspentOutputs(wallet)
        .retryWhen(errors => errors.delay(1000).take(10).concat(Observable.throw('')))
        .subscribe(
          result => {
            this.loadingUnspentOutputs = false;
            this.allUnspentOutputs = result;
            this.unspentOutputs = this.filterUnspentOutputs();
          },
          () => this.loadingUnspentOutputs = false,
        );

      this.addresses = wallet.addresses.filter(addr => addr.coins > 0);
      this.form.get('addresses').setValue(null);
      this.form.get('outputs').setValue(null);

      this.updateAvailableBalance();
      this.form.get('destinations').updateValueAndValidity();
    });

    this.subscriptions.add(this.form.get('addresses').valueChanges.subscribe(() => {
      this.form.get('outputs').setValue(null);
      this.unspentOutputs = this.filterUnspentOutputs();

      this.updateAvailableBalance();
      this.form.get('destinations').updateValueAndValidity();
    }));

    this.subscriptions.add(this.form.get('outputs').valueChanges.subscribe(() => {
      this.updateAvailableBalance();
      this.form.get('destinations').updateValueAndValidity();
    }));

    this.subscriptions = this.priceService.price.subscribe(price => {
      this.price = price;
      this.updateValues();
    });

    if (this.formData) {
      this.fillForm();
    }
  }

  ngOnDestroy() {
    if (this.processingSubscription && !this.processingSubscription.closed) {
      this.processingSubscription.unsubscribe();
    }
    this.closeGetOutputsSubscriptions();
    this.closeSyncCheckSubscription();
    this.subscriptions.unsubscribe();
    this.navbarService.hideSwitch();
    this.snackbar.dismiss();
    this.destinationSubscriptions.forEach(s => s.unsubscribe());
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
    this.selectedCurrency = value;
    this.updateValues();
    (this.form.get('destinations') as FormArray).updateValueAndValidity();
  }

  private updateValues() {
    if (!this.price) {
      this.values = null;

      return;
    }

    this.values = [];

    this.destControls.forEach((dest, i) => {
      const value = dest.get('coins').value !== undefined ? dest.get('coins').value.replace(' ', '=') : '';

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

  private checkBeforeSending() {
    if (!this.form.valid || this.previewButton.isLoading() || this.sendButton.isLoading()) {
      return;
    }

    this.closeSyncCheckSubscription();
    this.syncCheckSubscription = this.blockchainService.synchronized.first().subscribe(synchronized => {
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
    this.snackbar.dismiss();
    this.previewButton.resetState();
    this.sendButton.resetState();

    if (this.form.value.wallet.encrypted && !this.form.value.wallet.isHardware && !this.previewTx) {
      const config = new MatDialogConfig();
      config.data = {
        wallet: this.form.value.wallet,
      };

      this.dialog.open(PasswordDialogComponent, config).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.createTransaction(passwordDialog);
        });
    } else {
      if (!this.form.value.wallet.isHardware || this.previewTx) {
        this.createTransaction();
      } else {
        this.showBusy();
        this.processingSubscription = this.hwWalletService.checkIfCorrectHwConnected((this.form.value.wallet as Wallet).addresses[0].address).subscribe(
          () => this.createTransaction(),
          err => this.showError(getHardwareWalletErrorMsg(this.hwWalletService, this.translate, err)),
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

  toggleOptions(event) {
    event.stopPropagation();
    event.preventDefault();

    this.autoOptions = !this.autoOptions;
  }

  setAutoHours(event) {
    this.autoHours = event.checked;
    this.form.get('destinations').updateValueAndValidity();

    if (!this.autoHours) {
      this.autoOptions = false;
    }
  }

  private fillForm() {
    this.addresses = this.formData.form.wallet.addresses;

    ['wallet', 'addresses', 'changeAddress', 'note'].forEach(name => {
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

    if (this.formData.form.allUnspentOutputs) {
      this.closeGetOutputsSubscriptions();

      this.allUnspentOutputs = this.formData.form.allUnspentOutputs;
      this.unspentOutputs = this.filterUnspentOutputs();

      this.form.get('outputs').setValue(this.formData.form.outputs);
    }

    this.selectedCurrency = this.formData.form.currency;
  }

  addressCompare(a, b) {
    return a && b && a.address === b.address;
  }

  outputCompare(a, b) {
    return a && b && a.hash === b.hash;
  }

  get destControls() {
    return (this.form.get('destinations') as FormArray).controls;
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
        const value = control.get(name).value !== undefined
          ? control.get(name).value.replace(' ', '=')
          : '';

        if (isNaN(value) || value.trim() === '') {
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
          if (value < 1 || parseInt(value, 10) !== parseFloat(value)) {
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

    this.destinationSubscriptions.push(group.get('coins').valueChanges.subscribe(value => {
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

    const selectedAddresses = this.form.get('addresses').value && (this.form.get('addresses').value as Address[]).length > 0 ?
      this.form.get('addresses').value.map(addr => addr.address) : null;

    const selectedOutputs = this.form.get('outputs').value && (this.form.get('outputs').value as UnspentOutput[]).length > 0 ?
      this.form.get('outputs').value.map(addr => addr.hash) : null;

      this.processingSubscription = this.walletService.createTransaction(
        this.form.value.wallet,
        selectedAddresses ? selectedAddresses : (this.form.value.wallet as Wallet).addresses.map(address => address.address),
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
              if (note && !noteSaved) {
                showSnackbarError(this.snackbar, this.translate.instant('send.error-saving-note'));
              }

              this.showSuccess();
            }, error => this.showError(error));
        } else {
          let amount = new BigNumber('0');
          this.destinations.map(destination => amount = amount.plus(destination.coins));

          this.onFormSubmitted.emit({
            form: {
              wallet: this.form.get('wallet').value,
              addresses: this.form.get('addresses').value,
              changeAddress: this.form.get('changeAddress').value,
              destinations: this.destinations,
              hoursSelection: this.hoursSelection,
              autoOptions: this.autoOptions,
              allUnspentOutputs: this.loadingUnspentOutputs ? null : this.allUnspentOutputs,
              outputs: this.form.get('outputs').value,
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
          this.showError(getHardwareWalletErrorMsg(this.hwWalletService, this.translate, error));
        } else {
          this.showError(error);
        }
      });
  }

  private resetForm() {
    this.form.get('wallet').setValue('', { emitEvent: false });
    this.form.get('addresses').setValue(null);
    this.form.get('outputs').setValue(null);
    this.form.get('changeAddress').setValue('');
    this.form.get('note').setValue('');

    this.wallet = null;

    while (this.destControls.length > 0) {
      (this.form.get('destinations') as FormArray).removeAt(0);
    }

    this.addDestination();

    this.autoHours = true;
    this.autoOptions = false;
    this.autoShareValue = '0.5';
  }

  private get destinations() {
    return this.destControls.map((destControl, i) => {
      const destination = {
        address: destControl.get('address').value,
        coins: this.selectedCurrency === DoubleButtonActive.LeftButton ? destControl.get('coins').value : this.values[i].toString(),
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

  private updateAvailableBalance() {
    if (this.form.get('wallet').value) {
      this.availableCoins = new BigNumber(0);
      this.availableHours = new BigNumber(0);

      const outputs: UnspentOutput[] = this.form.get('outputs').value;
      const addresses: Address[] = this.form.get('addresses').value;

      if (outputs && outputs.length > 0) {
        outputs.map(control => {
          this.availableCoins = this.availableCoins.plus(control.coins);
          this.availableHours = this.availableHours.plus(control.calculated_hours);
        });
      } else if (addresses && addresses.length > 0) {
        addresses.map(control => {
          this.availableCoins = this.availableCoins.plus(control.coins);
          this.availableHours = this.availableHours.plus(control.hours);
        });
      } else if (this.form.get('wallet').value) {
        const wallet: Wallet = this.form.get('wallet').value;
        this.availableCoins = wallet.coins;
        this.availableHours = wallet.hours;
      }

      if (this.availableCoins.isGreaterThan(0)) {
        const unburnedHoursRatio = new BigNumber(1).minus(new BigNumber(1).dividedBy(this.appService.burnRate));
        const sendableHours = this.availableHours.multipliedBy(unburnedHoursRatio).decimalPlaces(0, BigNumber.ROUND_FLOOR);
        this.minimumFee = this.availableHours.minus(sendableHours);
        this.availableHours = sendableHours;
      } else {
        this.minimumFee = new BigNumber(0);
        this.availableHours = new BigNumber(0);
      }
    }
  }

  private filterUnspentOutputs(): UnspentOutput[] {
    if (this.allUnspentOutputs.length === 0) {
      return [];
    } else if (!this.form.get('addresses').value || (this.form.get('addresses').value as Address[]).length === 0) {
      return this.allUnspentOutputs;
    } else {
      return this.allUnspentOutputs.filter(out => (this.form.get('addresses').value as Address[]).some(addr => addr.address === out.address));
    }
  }

  private closeGetOutputsSubscriptions() {
    this.loadingUnspentOutputs = false;

    if (this.getOutputsSubscriptions) {
      this.getOutputsSubscriptions.unsubscribe();
    }
  }

  private closeSyncCheckSubscription() {
    if (this.syncCheckSubscription) {
      this.syncCheckSubscription.unsubscribe();
    }
  }

  private showSuccess() {
    this.busy = false;
    this.navbarService.enableSwitch();
    this.sendButton.setSuccess();
    this.resetForm();

    setTimeout(() => {
      this.sendButton.resetState();
    }, 3000);
  }

  private showError(error) {
    this.busy = false;
    this.navbarService.enableSwitch();
    showSnackbarError(this.snackbar, error);
    this.previewButton.resetState().setEnabled();
    this.sendButton.resetState().setEnabled();
  }
}
