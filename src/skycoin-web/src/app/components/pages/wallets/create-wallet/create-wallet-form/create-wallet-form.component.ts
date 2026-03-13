import { Component, OnInit, OnDestroy, Input } from '@angular/core';
import { UntypedFormBuilder, UntypedFormControl, UntypedFormGroup, Validators } from '@angular/forms';
import * as Bip39 from 'bip39';
import { Subscription } from 'rxjs';

import { CoinService } from '../../../../../services/coin.service';
import { BaseCoin } from '../../../../../coins/basecoin';
import { Bip39WordListService } from '../../../../../services/bip39-word-list.service';
import { environment } from '../../../../../../environments/environment';

export class FormData {
  label: string;
  seed: string;
  coin: BaseCoin;
  walletType: string;
  seedPassphrase: string;
}

@Component({
    selector: 'app-create-wallet-form',
    templateUrl: './create-wallet-form.component.html',
    styleUrls: ['./create-wallet-form.component.scss'],
    standalone: false
})
export class CreateWalletFormComponent implements OnInit, OnDestroy {
  @Input() create: boolean;
  @Input() whiteText: boolean;
  @Input() showSlowMobileInfo: boolean;

  form: UntypedFormGroup;
  hasManyCoins: boolean;
  normalSeed = false;
  customSeedAccepted = false;
  isProduction = environment.production;

  private statusSubscription: Subscription;

  constructor(
    private formBuilder: UntypedFormBuilder,
    private coinService: CoinService,
    private bip39WordListService: Bip39WordListService
  ) { }

  ngOnInit() {
    this.hasManyCoins = this.coinService.coins.length > 1;
    this.initForm(this.coinService.currentCoin.getValue());
  }

  ngOnDestroy() {
    this.statusSubscription.unsubscribe();
  }

  get isValid(): boolean {
    return this.form.valid && (this.normalSeed || this.customSeedAccepted);
  }

  onCustomSeedAcceptance(event) {
    this.customSeedAccepted = event.checked;
  }

  getData(): FormData {
    return {
      label: this.form.value.label,
      seed: this.form.value.seed,
      coin: this.form.value.coin,
      walletType: this.form.value.wallet_type || 'deterministic',
      seedPassphrase: this.form.value.seed_passphrase || '',
    };
  }

  initForm(defaultCoin: BaseCoin, create: boolean = null) {
    create = create !== null ? create : this.create;

    this.form = this.formBuilder.group({
        label: new UntypedFormControl('', [ Validators.required ]),
        coin: new UntypedFormControl(defaultCoin, [ Validators.required ]),
        seed: new UntypedFormControl('', [ Validators.required ]),
        confirm_seed: new UntypedFormControl(),
        wallet_type: new UntypedFormControl('deterministic'),
        seed_passphrase: new UntypedFormControl(''),
      },
      {
        validator: create ? this.seedMatchValidator.bind(this) : null,
      }
    );

    if (create) {
      this.generateSeed(128);
    }

    this.statusSubscription = this.form.statusChanges.subscribe(() => {
      this.customSeedAccepted = false;
      this.normalSeed = this.validateSeed(this.form.get('seed').value);
    });
  }

  private generateSeed(entropy: number) {
    this.form.controls.seed.setValue(Bip39.generateMnemonic(entropy));
  }

  private validateSeed(seed: string): boolean {
    const processedSeed = seed.replace(/\r?\n|\r/g, ' ').replace(/ +/g, ' ').trim();
    if (seed !== processedSeed) {
      return false;
    }

    const words = seed.split(' ');
    const numberOfWords = words.length;
    if (numberOfWords !== 12 && numberOfWords !== 24) {
      return false;
    }

    for (let i = 0; i < numberOfWords; i++) {
      const validation = this.bip39WordListService.validateWord(words[i]);
      if (validation != null && !validation) {
        return false;
      }
    }

    return true;
  }

  private seedMatchValidator(formGroup: UntypedFormGroup) {
    return formGroup.get('seed').value === formGroup.get('confirm_seed').value ? null : { NotEqual: true };
  }
}
