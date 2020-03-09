import { throwError as observableThrowError, SubscriptionLike, of } from 'rxjs';
import { retryWhen, delay, mergeMap, debounceTime } from 'rxjs/operators';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { BigNumber } from 'bignumber.js';
import { HttpErrorResponse } from '@angular/common/http';

import { AppService } from '../../../../../services/app.service';
import { BalanceAndOutputsService } from '../../../../../services/wallet-operations/balance-and-outputs.service';
import { WalletWithBalance, AddressWithBalance } from '../../../../../services/wallet-operations/wallet-objects';
import { Output as UnspentOutput } from '../../../../../services/wallet-operations/transaction-objects';
import { processServiceError } from '../../../../../utils/errors';
import { OperationError } from '../../../../../utils/operation-error';
import { SendCoinsData } from '../../send-coins-form/send-coins-form.component';

/**
 * Info about the balance which is available with the selections the user has
 * made in FormSourceSelectionComponent.
 */
export class AvailableBalanceData {
  /**
   * How many coins the user can send.
   */
  availableCoins = new BigNumber(0);
  /**
   * How many hours the user can send. It is not the total amount of hours in the selected
   * sources, because the fee for sending all the coins has been subtracted from the value.
   */
  availableHours = new BigNumber(0);
  /**
   * Minimum fee in hours the user must pay for sending all the available coins. If the user
   * does not send all the coins, the fee could be lower.
   */
  minimumFee = new BigNumber(0);
  /**
   * If the balance is still being loded. Only for the manual mode.
   */
  loading = false;
}

/**
 * Sources the user has selected using FormSourceSelectionComponent.
 */
export interface SelectedSources {
  /**
   * Selected wallet.
   */
  wallet: WalletWithBalance;
  /**
   * Optional selected addresses. The addresses are from the selected wallet.
   */
  addresses?: AddressWithBalance[];
  /**
   * Addresses the user has entered manually. Only for the manual mode.
   */
  manualAddresses?: string[];
  /**
   * Optional selected unspent outputs. The outputs are from the selected wallet and addresses
   * or the manually entered addresses.
   */
  unspentOutputs: UnspentOutput[];
}

/**
 * Modes in which FormSourceSelectionComponent can work.
 */
export enum SourceSelectionModes {
  /**
   * Simple mode in which the user can only select a wallet.
   */
  Wallet = 'Wallet',
  /**
   * Advanced mode in which the user can selected a wallet and then choose specific addresses and
   * outputs to have more control about the source from were the coins will be sent.
   */
  All = 'All',
  /**
   * Manual mode in which the user must enter the addresses manually, and have the option for
   * selecting specific outputs.
   */
  Manual = 'Manual',
}

/**
 * Allows the user to select the sources from were the coins will be sent. It has various modes
 * which allow different levels of control.
 */
@Component({
  selector: 'app-form-source-selection',
  templateUrl: './form-source-selection.component.html',
  styleUrls: ['./form-source-selection.component.scss'],
})
export class FormSourceSelectionComponent implements OnInit, OnDestroy {
  // Allows to deactivate the form while the system is busy.
  @Input() busy: boolean;
  // Event for informing when the user selection has changed or when there was a change in
  // the available balance.
  @Output() onSelectionChanged = new EventEmitter<void>();

  // Sets the mode in which the component works.
  private selectionModeInternal: SourceSelectionModes;
  @Input() set selectionMode(val: SourceSelectionModes) {
    this.selectionModeInternal = val;
    if (this.form) {
      this.resetForm();
      this.form.updateValueAndValidity();
    }
  }
  get selectionMode(): SourceSelectionModes {
    return this.selectionModeInternal;
  }

  sourceSelectionModes = SourceSelectionModes;
  form: FormGroup;
  // All available wallets.
  allWallets: WalletWithBalance[];
  // Wallet selected by the user.
  wallet: WalletWithBalance;
  // List with all addresses from the selected wallet.
  addresses: AddressWithBalance[] = [];
  // List of the addresses manually entered by the user, if in manual mode.
  manualAddresses: string[] = [];
  // List with all the outputs on the selected wallet or the manually entered addresses,
  // depending on the selected mode.
  allUnspentOutputs: UnspentOutput[] = [];
  // Filtered version of allUnspentOutputs shown on the UI. It does not contain the
  // outputs which are irrelevant for the other selected sources.
  unspentOutputs: UnspentOutput[] = [];
  // If the list of unspent outputs is being downloaded from the node.
  loadingUnspentOutputs = false;
  // If there was an error downloading the list of unspent outputs from the node.
  errorLoadingManualOutputs = false;

  private subscriptionsGroup: SubscriptionLike[] = [];
  private getOutputsSubscription: SubscriptionLike;

  constructor(
    private appService: AppService,
    private formBuilder: FormBuilder,
    private balanceAndOutputsService: BalanceAndOutputsService,
  ) { }

  ngOnInit() {
    this.form = this.formBuilder.group({
      manualAddresses: [''],
      wallet: [''],
      addresses: [null],
      outputs: [null],
    });

    this.form.setValidators(this.validateForm.bind(this));

    // When the user enters the addresses manually.
    this.subscriptionsGroup.push(this.form.get('manualAddresses').valueChanges.pipe(debounceTime(500)).subscribe(() => {
      // Separate all addresses.
      const manuallyEnteredAddresses = (this.form.get('manualAddresses').value as string);
      let manualAddresses: string[] = [];
      if (manuallyEnteredAddresses && manuallyEnteredAddresses.trim().length > 0) {
        manualAddresses = manuallyEnteredAddresses.split(',');
        manualAddresses = manualAddresses.map(address => address.trim());
        manualAddresses = manualAddresses.filter(address => address.length > 0);
      }

      // Check if there was a change.
      let addressesChanged = false;
      if (manualAddresses.length !== this.manualAddresses.length) {
        addressesChanged = true;
      } else {
        manualAddresses.forEach((address, i) => {
          if (this.manualAddresses[i] !== address) {
            addressesChanged = true;
          }
        });
      }

      // If a change was detected, save the list and reset the form before loading the data again.
      if (addressesChanged) {
        this.manualAddresses = manualAddresses;

        this.closeGetOutputsSubscription();
        this.allUnspentOutputs = [];
        this.unspentOutputs = [];

        this.form.get('outputs').setValue(null);

        this.loadingUnspentOutputs = manualAddresses.length !== 0;
        this.errorLoadingManualOutputs = false;
        // Inform about the changes in the available balance and loading status.
        this.onSelectionChanged.emit();
      }

      if (manualAddresses.length !== 0 && addressesChanged) {
        // Get the outputs of the entered addresses.
        this.getOutputsSubscription = this.balanceAndOutputsService.getOutputs((this.form.get('manualAddresses').value as string).replace(/ /g, '')).pipe(
          // Retry if there is an error, but not if the server returns 400 as response, which
          // means the user entered at least one invalid address.
          retryWhen((err) => {
            return err.pipe(mergeMap((response: OperationError) => {
              response = processServiceError(response);
              if (response.originalError && (response.originalError as HttpErrorResponse).status && (response.originalError as HttpErrorResponse).status === 400) {
                this.errorLoadingManualOutputs = true;

                return observableThrowError(response);
              }

              return of(response);
            }), delay(4000));
          }),
        ).subscribe(
          result => {
            this.loadingUnspentOutputs = false;
            this.allUnspentOutputs = result;
            this.unspentOutputs = this.filterUnspentOutputs();

            // Inform about the changes in the available balance and loading status.
            this.onSelectionChanged.emit();
          },
          () => {
            this.loadingUnspentOutputs = false;
            // Inform about the changes in the available balance and loading status.
            this.onSelectionChanged.emit();
          },
        );
      }
    }));

    // When the user changes the wallet using the dropdown.
    this.subscriptionsGroup.push(this.form.get('wallet').valueChanges.subscribe(wallet => {
      this.wallet = wallet;

      // Reset the form.
      this.closeGetOutputsSubscription();
      this.allUnspentOutputs = [];
      this.unspentOutputs = [];
      this.form.get('addresses').setValue(null);
      this.form.get('outputs').setValue(null);
      this.loadingUnspentOutputs = false;

      // Load the output list, if the form is showing a dropdown for selecting them.
      if (wallet && this.selectionMode === SourceSelectionModes.All) {
        this.loadingUnspentOutputs = true;
        this.getOutputsSubscription = this.balanceAndOutputsService.getWalletUnspentOutputs(wallet).pipe(
          retryWhen(errors => errors.pipe(delay(4000))))
          .subscribe(
            result => {
              this.loadingUnspentOutputs = false;
              this.allUnspentOutputs = result;
              this.unspentOutputs = this.filterUnspentOutputs();
            },
            () => this.loadingUnspentOutputs = false,
          );
      }

      if (wallet) {
        this.addresses = wallet.addresses.filter(addr => addr.coins > 0);
      } else {
        this.addresses = [];
      }

      this.onSelectionChanged.emit();
    }));

    this.subscriptionsGroup.push(this.form.get('addresses').valueChanges.subscribe(() => {
      this.form.get('outputs').setValue(null);
      this.unspentOutputs = this.filterUnspentOutputs();

      this.onSelectionChanged.emit();
    }));

    this.subscriptionsGroup.push(this.form.get('outputs').valueChanges.subscribe(() => {
      this.onSelectionChanged.emit();
    }));

    this.subscriptionsGroup.push(this.balanceAndOutputsService.walletsWithBalance.subscribe(wallets => {
      this.allWallets = wallets;
      if (wallets.length === 1) {
        setTimeout(() => {
          try {
            this.form.get('wallet').setValue(wallets[0]);
          } catch (e) { }
        });
      }
    }));
  }

  ngOnDestroy() {
    this.closeGetOutputsSubscription();
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
  }

  resetForm() {
    this.form.get('manualAddresses').setValue('');
    this.form.get('wallet').setValue('');
    this.form.get('addresses').setValue(null);
    this.form.get('outputs').setValue(null);

    this.wallet = null;
  }

  /**
   * Fills the form with the provided values.
   */
  fill(formData: SendCoinsData) {
    setTimeout(() => {
      if (this.selectionMode === SourceSelectionModes.Manual) {
        let addressesString = '';

        const addresses = (formData.form.manualAddresses as string[]);
        addresses.forEach((address, i) => {
          addressesString += address;
          if (i < addresses.length - 1) {
            addressesString += ', ';
          }
        });

        this.form.get('manualAddresses').setValue(addressesString, { emitEvent: false });
        this.manualAddresses = formData.form.manualAddresses;
      } else {
        this.addresses = formData.form.wallet.addresses;

        ['wallet', 'addresses'].forEach(name => {
          this.form.get(name).setValue(formData.form[name]);
        });
      }

      this.closeGetOutputsSubscription();
      this.allUnspentOutputs = formData.form.allUnspentOutputs;
      this.unspentOutputs = this.filterUnspentOutputs();
      this.form.get('outputs').setValue(formData.form.outputs);
    });
  }

  // Text shown in the addresses dropdown while no value is selected.
  get addressessPlaceholder(): string {
    if (this.wallet) {
      return 'send.all-addresses';
    } else {
      return 'send.enter-wallet-for-outputs-and-addresses';
    }
  }

  // Text shown in the outputs dropdown while no value is selected.
  get outputsPlaceholder(): string {
    if (this.selectionMode !== SourceSelectionModes.Manual) {
      if (this.wallet) {
        return 'send.all-outputs';
      } else {
        return 'send.enter-wallet-for-outputs-and-addresses';
      }
    } else {
      if (this.loadingUnspentOutputs) {
        return 'send.all-outputs';
      } else if (this.errorLoadingManualOutputs) {
        return 'send.invalid-addresses-for-outputs';
      } else if (!this.manualAddresses || this.manualAddresses.length === 0) {
        return 'send.enter-addresses-for-outputs';
      } else {
        return 'send.all-outputs';
      }
    }
  }

  get valid(): boolean {
    return this.form.valid;
  }

  /**
   * Allows to know the available balance with the selected sources.
   */
  get availableBalance(): AvailableBalanceData {
    const response = new AvailableBalanceData();

    // While in manual mode, the balance is obtained from the available or selected
    // unspent outputs.
    if (this.selectionMode === SourceSelectionModes.Manual && this.unspentOutputs && this.unspentOutputs.length > 0) {
      const selectedOutputs: UnspentOutput[] = this.form.get('outputs').value;

      if (selectedOutputs && selectedOutputs.length > 0) {
        selectedOutputs.map(output => {
          response.availableCoins = response.availableCoins.plus(output.coins);
          response.availableHours = response.availableHours.plus(output.hours);
        });
      } else {
        this.unspentOutputs.forEach(output => {
          response.availableCoins = response.availableCoins.plus(output.coins);
          response.availableHours = response.availableHours.plus(output.hours);
        });
      }
    }

    if (this.selectionMode === SourceSelectionModes.Manual) {
      response.loading = this.loadingUnspentOutputs;
    }

    // Only if not in manual mode (it is not possible to select a wallet in manual mode).
    if (this.form.get('wallet').value) {
      const outputs: UnspentOutput[] = this.form.get('outputs').value;
      const addresses: AddressWithBalance[] = this.form.get('addresses').value;

      // Get the balance from the selected outputs, the selected addresses, or the
      // selected wallet, depending on what the user has selected.
      if (outputs && outputs.length > 0) {
        outputs.map(control => {
          response.availableCoins = response.availableCoins.plus(control.coins);
          response.availableHours = response.availableHours.plus(control.hours);
        });
      } else if (addresses && addresses.length > 0) {
        addresses.map(control => {
          response.availableCoins = response.availableCoins.plus(control.coins);
          response.availableHours = response.availableHours.plus(control.hours);
        });
      } else if (this.form.get('wallet').value) {
        const wallet: WalletWithBalance = this.form.get('wallet').value;
        response.availableCoins = wallet.coins;
        response.availableHours = wallet.hours;
      }
    }

    // Calculate the max number of hours that can be sent.
    if (response.availableCoins.isGreaterThan(0)) {
      const unburnedHoursRatio = new BigNumber(1).minus(new BigNumber(1).dividedBy(this.appService.burnRate));
      const sendableHours = response.availableHours.multipliedBy(unburnedHoursRatio).decimalPlaces(0, BigNumber.ROUND_FLOOR);
      response.minimumFee = response.availableHours.minus(sendableHours);
      response.availableHours = sendableHours;
    }

    return response;
  }

  /**
   * Returns the sources the user has selected.
   */
  get selectedSources(): SelectedSources {
    if (this.selectionMode === SourceSelectionModes.Manual) {
      return {
        wallet: null,
        manualAddresses: this.manualAddresses,
        unspentOutputs: this.form.get('outputs').value,
      };
    } else {
      return {
        wallet: this.form.get('wallet').value,
        addresses: (this.form.get('addresses').value as AddressWithBalance[]),
        unspentOutputs: this.form.get('outputs').value,
      };
    }
  }

  /**
   * Returns the last list of unspent outputs obtained from the node.
   */
  get unspentOutputsList(): UnspentOutput[] {
    return this.loadingUnspentOutputs ? null : this.allUnspentOutputs;
  }

  /**
   * Filters the list of unspent outputs obtained from the node, to contain only the outputs
   * relevant for the currently selected sources, and saves it to be shown by the unspent
   * outputs dropdown.
   */
  private filterUnspentOutputs(): UnspentOutput[] {
    if (this.selectionMode === SourceSelectionModes.Manual) {
      return this.allUnspentOutputs;
    }

    // Use only the outputs of the selected addresses, if the user has selected one or more.
    if (this.allUnspentOutputs.length === 0) {
      return [];
    } else if (!this.form.get('addresses').value || (this.form.get('addresses').value as AddressWithBalance[]).length === 0) {
      return this.allUnspentOutputs;
    } else {
      const addressMap = new Map<string, boolean>();
      (this.form.get('addresses').value as AddressWithBalance[]).forEach(address => addressMap.set(address.address, true));

      return this.allUnspentOutputs.filter(out => addressMap.has(out.address));
    }
  }

  private closeGetOutputsSubscription() {
    this.loadingUnspentOutputs = false;

    if (this.getOutputsSubscription) {
      this.getOutputsSubscription.unsubscribe();
    }
  }

  /**
   * Validates the form.
   */
  private validateForm() {
    if (!this.form) {
      return { Required: true };
    }

    // The validation depends on the current mode.
    if (this.selectionMode === SourceSelectionModes.Manual) {
      if (!this.form.get('manualAddresses').value) {
        return { Invalid: true };
      }
    } else {
      if (!this.form.get('wallet').value) {
        return { Invalid: true };
      }
    }

    return null;
  }
}
