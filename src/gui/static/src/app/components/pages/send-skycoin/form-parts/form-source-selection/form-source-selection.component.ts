import { throwError as observableThrowError, SubscriptionLike } from 'rxjs';
import { retryWhen, delay, take, concat, first } from 'rxjs/operators';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { WalletService } from '../../../../../services/wallet.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { BigNumber } from 'bignumber.js';
import { Output as UnspentOutput, Wallet, Address } from '../../../../../app.datatypes';
import { BlockchainService } from '../../../../../services/blockchain.service';
import { AppService } from '../../../../../services/app.service';

export class AvailableBalanceData {
  availableCoins = new BigNumber(0);
  availableHours = new BigNumber(0);
  minimumFee = new BigNumber(0);
}

export interface SelectedSources {
  wallet: Wallet;
  addresses: Address[];
  unspentOutputs: UnspentOutput[];
}

@Component({
  selector: 'app-form-source-selection',
  templateUrl: './form-source-selection.component.html',
  styleUrls: ['./form-source-selection.component.scss'],
})
export class FormSourceSelectionComponent implements OnInit, OnDestroy {
  @Input() busy: boolean;
  @Output() onselectionChanged = new EventEmitter<void>();

  form: FormGroup;
  wallets: Wallet[];
  wallet: Wallet;
  addresses: Address[] = [];
  allUnspentOutputs: UnspentOutput[] = [];
  unspentOutputs: UnspentOutput[] = [];
  loadingUnspentOutputs = false;

  private subscriptionsGroup: SubscriptionLike[] = [];
  private getOutputsSubscriptions: SubscriptionLike;

  constructor(
    public blockchainService: BlockchainService,
    public walletService: WalletService,
    private appService: AppService,
    private formBuilder: FormBuilder,
  ) { }

  ngOnInit() {
    this.form = this.formBuilder.group({
      wallet: ['', Validators.required],
      addresses: [null],
      outputs: [null],
    });

    this.subscriptionsGroup.push(this.form.get('wallet').valueChanges.subscribe(wallet => {
      this.wallet = wallet;

      this.closeGetOutputsSubscriptions();
      this.allUnspentOutputs = [];
      this.unspentOutputs = [];
      this.loadingUnspentOutputs = true;

      this.getOutputsSubscriptions = this.walletService.getWalletUnspentOutputs(wallet).pipe(
        retryWhen(errors => errors.pipe(delay(1000), take(10), concat(observableThrowError('')))))
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

      this.onselectionChanged.emit();
    }));

    this.subscriptionsGroup.push(this.form.get('addresses').valueChanges.subscribe(() => {
      this.form.get('outputs').setValue(null);
      this.unspentOutputs = this.filterUnspentOutputs();

      this.onselectionChanged.emit();
    }));

    this.subscriptionsGroup.push(this.form.get('outputs').valueChanges.subscribe(() => {
      this.onselectionChanged.emit();
    }));

    this.subscriptionsGroup.push(this.walletService.all().pipe(first()).subscribe(wallets => {
      this.wallets = wallets;
      if (wallets.length === 1) {
        setTimeout(() => this.form.get('wallet').setValue(wallets[0]));
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
    this.form.get('wallet').setValue('', { emitEvent: false });
    this.form.get('addresses').setValue(null);
    this.form.get('outputs').setValue(null);

    this.wallet = null;
  }

  fill(formData: any) {
    this.addresses = formData.form.wallet.addresses;

    ['wallet', 'addresses'].forEach(name => {
      this.form.get(name).setValue(formData.form[name]);
    });

    if (formData.form.allUnspentOutputs) {
      this.closeGetOutputsSubscriptions();

      this.allUnspentOutputs = formData.form.allUnspentOutputs;
      this.unspentOutputs = this.filterUnspentOutputs();

      this.form.get('outputs').setValue(formData.form.outputs);
    }
  }

  get valid(): boolean {
    return this.form.valid;
  }

  get availableBalance(): AvailableBalanceData {
    const response = new AvailableBalanceData();

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

      if (response.availableCoins.isGreaterThan(0)) {
        const unburnedHoursRatio = new BigNumber(1).minus(new BigNumber(1).dividedBy(this.appService.burnRate));
        const sendableHours = response.availableHours.multipliedBy(unburnedHoursRatio).decimalPlaces(0, BigNumber.ROUND_FLOOR);
        response.minimumFee = response.availableHours.minus(sendableHours);
        response.availableHours = sendableHours;
      }
    }

    return response;
  }

  get selectedSources(): SelectedSources {
    return {
      wallet: this.form.get('wallet').value,
      addresses: this.form.get('addresses').value,
      unspentOutputs: this.form.get('outputs').value,
    };
  }

  get unspentOutputsList(): UnspentOutput[] {
    return this.loadingUnspentOutputs ? null : this.allUnspentOutputs;
  }

  private filterUnspentOutputs(): UnspentOutput[] {
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
}
