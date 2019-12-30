import { throwError as observableThrowError, SubscriptionLike, of } from 'rxjs';
import { retryWhen, delay, first, mergeMap } from 'rxjs/operators';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { WalletService } from '../../../../../services/wallet.service';
import { FormBuilder, FormGroup } from '@angular/forms';
import { BigNumber } from 'bignumber.js';
import { Output as UnspentOutput, Wallet, Address } from '../../../../../app.datatypes';
import { AppService } from '../../../../../services/app.service';
import { HttpErrorResponse } from '@angular/common/http';

export class AvailableBalanceData {
  availableCoins = new BigNumber(0);
  availableHours = new BigNumber(0);
  minimumFee = new BigNumber(0);
  loading = false;
}

export interface SelectedSources {
  wallet: Wallet;
  addresses?: Address[];
  manualAddresses?: string[];
  unspentOutputs: UnspentOutput[];
}

export enum SourceSelectionModes {
  Wallet,
  All,
  Manual,
}

@Component({
  selector: 'app-form-source-selection',
  templateUrl: './form-source-selection.component.html',
  styleUrls: ['./form-source-selection.component.scss'],
})
export class FormSourceSelectionComponent implements OnInit, OnDestroy {
  @Input() busy: boolean;
  @Output() onSelectionChanged = new EventEmitter<void>();

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
  wallets: Wallet[];
  wallet: Wallet;
  addresses: Address[] = [];
  manualAddresses: string[] = [];
  allUnspentOutputs: UnspentOutput[] = [];
  unspentOutputs: UnspentOutput[] = [];
  loadingUnspentOutputs = false;
  errorLoadingManualOutputs = false;

  private subscriptionsGroup: SubscriptionLike[] = [];
  private getOutputsSubscriptions: SubscriptionLike;

  constructor(
    private walletService: WalletService,
    private appService: AppService,
    private formBuilder: FormBuilder,
  ) { }

  ngOnInit() {
    this.form = this.formBuilder.group({
      manualAddresses: [''],
      wallet: [''],
      addresses: [null],
      outputs: [null],
    });

    this.form.setValidators(this.validateForm.bind(this));

    this.subscriptionsGroup.push(this.form.get('manualAddresses').valueChanges.subscribe(() => {
      const manuallyEnteredAddresses = (this.form.get('manualAddresses').value as string);
      let manualAddresses: string[] = [];
      if (manuallyEnteredAddresses && manuallyEnteredAddresses.trim().length > 0) {
        manualAddresses = manuallyEnteredAddresses.split(',');
        manualAddresses = manualAddresses.map(address => address.trim());
        manualAddresses = manualAddresses.filter(address => address.length > 0);
      }

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

      this.loadingUnspentOutputs = false;

      if (addressesChanged || manualAddresses.length === 0) {
        this.manualAddresses = manualAddresses;

        this.closeGetOutputsSubscriptions();
        this.allUnspentOutputs = [];
        this.unspentOutputs = [];

        this.form.get('outputs').setValue(null);

        this.loadingUnspentOutputs = manualAddresses.length !== 0;
        this.errorLoadingManualOutputs = false;
        this.onSelectionChanged.emit();
      }

      if (manualAddresses.length !== 0 && addressesChanged) {
        this.getOutputsSubscriptions = of(1).pipe(delay(500), mergeMap(() => {
          return this.walletService.getOutputs((this.form.get('manualAddresses').value as string).replace(/ /g, '')).pipe(
            retryWhen((err) => {
              return err.pipe(mergeMap(response => {
                if (response instanceof HttpErrorResponse && response.status === 400) {
                  this.errorLoadingManualOutputs = true;

                  return observableThrowError(response);
                }

                return of(response);
              }), delay(1000));
            }),
          );
        })).subscribe(
          result => {
            this.loadingUnspentOutputs = false;
            this.allUnspentOutputs = result;
            this.unspentOutputs = this.filterUnspentOutputs();

            this.onSelectionChanged.emit();
          },
          () => {
            this.loadingUnspentOutputs = false;
            this.onSelectionChanged.emit();
          },
        );
      }
    }));

    this.subscriptionsGroup.push(this.form.get('wallet').valueChanges.subscribe(wallet => {
      this.wallet = wallet;

      this.closeGetOutputsSubscriptions();
      this.allUnspentOutputs = [];
      this.unspentOutputs = [];

      if (wallet && this.selectionMode === SourceSelectionModes.All) {
        this.loadingUnspentOutputs = true;
        this.getOutputsSubscriptions = this.walletService.getWalletUnspentOutputs(wallet).pipe(
          retryWhen(errors => errors.pipe(delay(1000))))
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
      this.form.get('addresses').setValue(null);
      this.form.get('outputs').setValue(null);

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

    this.subscriptionsGroup.push(this.walletService.all().pipe(first()).subscribe(wallets => {
      this.wallets = wallets;
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
    this.closeGetOutputsSubscriptions();
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
  }

  addressCompare(a, b) {
    return a && b && a.address === b.address;
  }

  outputCompare(a, b) {
    return a && b && a.hash === b.hash;
  }

  resetForm() {
    this.form.get('manualAddresses').setValue('');
    this.form.get('wallet').setValue('');
    this.form.get('addresses').setValue(null);
    this.form.get('outputs').setValue(null);

    this.wallet = null;
  }

  fill(formData: any) {
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

        this.form.get('manualAddresses').setValue(addressesString);
      } else {
        this.addresses = formData.form.wallet.addresses;

        ['wallet', 'addresses'].forEach(name => {
          if (formData.form[name]) {
            this.form.get(name).setValue(formData.form[name]);
          }
        });
      }

      if (formData.form.allUnspentOutputs) {
        this.closeGetOutputsSubscriptions();

        this.allUnspentOutputs = formData.form.allUnspentOutputs;
        this.unspentOutputs = this.filterUnspentOutputs();

        this.form.get('outputs').setValue(formData.form.outputs);
      }
    });
  }

  get addressessPlaceholder(): string {
    if (this.wallet) {
      return 'send.all-addresses';
    } else {
      return 'send.enter-wallet-for-outputs-and-addresses';
    }
  }

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

  get availableBalance(): AvailableBalanceData {
    const response = new AvailableBalanceData();

    if (this.selectionMode === SourceSelectionModes.Manual && this.unspentOutputs && this.unspentOutputs.length > 0) {
      const selectedOutputs: UnspentOutput[] = this.form.get('outputs').value;

      if (selectedOutputs && selectedOutputs.length > 0) {
        selectedOutputs.map(output => {
          response.availableCoins = response.availableCoins.plus(output.coins);
          response.availableHours = response.availableHours.plus(output.calculated_hours);
        });
      } else {
        this.unspentOutputs.forEach(output => {
          response.availableCoins = response.availableCoins.plus(output.coins);
          response.availableHours = response.availableHours.plus(output.calculated_hours);
        });
      }
    }

    if (this.selectionMode === SourceSelectionModes.Manual) {
      response.loading = this.loadingUnspentOutputs;
    }

    if (this.form.get('wallet').value) {
      const outputs: UnspentOutput[] = this.form.get('outputs').value;
      const addresses: Address[] = this.form.get('addresses').value;

      if (outputs && outputs.length > 0) {
        outputs.map(control => {
          response.availableCoins = response.availableCoins.plus(control.coins);
          response.availableHours = response.availableHours.plus(control.calculated_hours);
        });
      } else if (addresses && addresses.length > 0) {
        addresses.map(control => {
          response.availableCoins = response.availableCoins.plus(control.coins);
          response.availableHours = response.availableHours.plus(control.hours);
        });
      } else if (this.form.get('wallet').value) {
        const wallet: Wallet = this.form.get('wallet').value;
        response.availableCoins = wallet.coins;
        response.availableHours = wallet.hours;
      }
    }

    if (response.availableCoins.isGreaterThan(0)) {
      const unburnedHoursRatio = new BigNumber(1).minus(new BigNumber(1).dividedBy(this.appService.burnRate));
      const sendableHours = response.availableHours.multipliedBy(unburnedHoursRatio).decimalPlaces(0, BigNumber.ROUND_FLOOR);
      response.minimumFee = response.availableHours.minus(sendableHours);
      response.availableHours = sendableHours;
    }

    return response;
  }

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
        addresses: (this.form.get('addresses').value as Address[]),
        unspentOutputs: this.form.get('outputs').value,
      };
    }
  }

  get unspentOutputsList(): UnspentOutput[] {
    return this.loadingUnspentOutputs ? null : this.allUnspentOutputs;
  }

  private filterUnspentOutputs(): UnspentOutput[] {
    if (this.selectionMode === SourceSelectionModes.Manual) {
      return this.allUnspentOutputs;
    }

    if (this.allUnspentOutputs.length === 0) {
      return [];
    } else if (!this.form.get('addresses').value || (this.form.get('addresses').value as Address[]).length === 0) {
      return this.allUnspentOutputs;
    } else {
      const addressMap = new Map<string, boolean>();
      (this.form.get('addresses').value as Address[]).forEach(address => addressMap.set(address.address, true));

      return this.allUnspentOutputs.filter(out => addressMap.has(out.address));
    }
  }

  private closeGetOutputsSubscriptions() {
    this.loadingUnspentOutputs = false;

    if (this.getOutputsSubscriptions) {
      this.getOutputsSubscriptions.unsubscribe();
    }
  }

  private validateForm() {
    if (!this.form) {
      return { Required: true };
    }

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
