import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnDestroy, OnInit } from '@angular/core';
import { UntypedFormBuilder, UntypedFormControl, UntypedFormGroup, Validators } from '@angular/forms';
import { generateMnemonic } from 'bip39';
import { Subscription } from 'rxjs';

import { CoinService } from '../../../../../services/coin.service';
import { NodeHealthService, CoinHealth } from '../../../../../services/node-health.service';
import { BaseCoin } from '../../../../../coins/basecoin';
import { Bip39WordListService } from '../../../../../services/bip39-word-list.service';
import { environment } from '../../../../../../environments/environment';

export class FormData {
  label!: string;
  seed!: string;
  coin!: BaseCoin;
  walletType!: string;
  seedPassphrase!: string;
  segwit!: boolean;
  // Optional password entered at creation. When set, the new wallet's seed is
  // encrypted right after it's created; blank leaves it stored unencrypted.
  password!: string;
}

@Component({
    selector: 'app-create-wallet-form',
    templateUrl: './create-wallet-form.component.html',
    styleUrls: ['./create-wallet-form.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class CreateWalletFormComponent implements OnInit, OnDestroy {
  @Input() create!: boolean;
  @Input() whiteText!: boolean;
  @Input() showSlowMobileInfo!: boolean;

  form!: UntypedFormGroup;
  hasManyCoins!: boolean;
  normalSeed = false;
  customSeedAccepted = false;
  isProduction = environment.production;
  showWalletType = false;
  isBitcoinCoin = false;

  // Backend connection health of the currently-selected coin, so the create
  // screen can gate the coin choice by whether its node/electrum is actually
  // reachable (shown even for the default coin when no wallet exists yet).
  selectedCoinHealth: CoinHealth | null = null;
  checkingCoinHealth = false;

  private statusSubscription!: Subscription;
  private healthSubscription!: Subscription;
  private coinsLoadedSub!: Subscription;

  constructor(
    private formBuilder: UntypedFormBuilder,
    private coinService: CoinService,
    private nodeHealthService: NodeHealthService,
    private bip39WordListService: Bip39WordListService,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit() {
    // Build the form only once the coin list has loaded. Coins arrive
    // asynchronously (server request), so reading coins/currentCoin
    // synchronously here would see an empty list — hiding both the coin
    // selector (hasManyCoins) and the wallet-type selector, and leaving a
    // bare deterministic-only create form. coinsLoaded is a ReplaySubject(1):
    // it fires immediately if coins are already loaded, else when they arrive.
    this.coinsLoadedSub = this.coinService.coinsLoaded.subscribe(() => {
      this.hasManyCoins = this.coinService.coins.length > 1;
      const coin = this.coinService.currentCoin.getValue() || this.coinService.coins[0] || null;
      this.initForm(coin);
      this.changeDetectorRef.markForCheck();
    });
  }

  ngOnDestroy() {
    if (this.coinsLoadedSub) {
      this.coinsLoadedSub.unsubscribe();
    }
    if (this.statusSubscription) {
      this.statusSubscription.unsubscribe();
    }
    if (this.healthSubscription) {
      this.healthSubscription.unsubscribe();
    }
  }

  /**
   * Probes the given coin's backend so the create form can show whether the
   * node/electrum for the chosen coin is reachable. Wallet creation itself is
   * client-side, so an unavailable backend does not block creation — it only
   * warns the user the wallet won't sync until the connection is restored.
   */
  private probeCoinHealth(coin: BaseCoin) {
    if (this.healthSubscription) {
      this.healthSubscription.unsubscribe();
    }
    this.selectedCoinHealth = null;
    this.checkingCoinHealth = true;
    this.healthSubscription = this.nodeHealthService.check(coin).subscribe(health => {
      this.selectedCoinHealth = health;
      this.checkingCoinHealth = false;
      this.changeDetectorRef.markForCheck();
    });
  }

  /** True when we have a probe result and the selected coin's backend is down. */
  get selectedCoinUnavailable(): boolean {
    return !!this.selectedCoinHealth && !this.selectedCoinHealth.available;
  }

  get isValid(): boolean {
    return this.form.valid && (this.normalSeed || this.customSeedAccepted);
  }

  onCustomSeedAcceptance(event: any) {
    this.customSeedAccepted = event.checked;
  }

  /**
   * The wallet type new wallets should default to for the given coin. bip44
   * (HD, multicoin-capable) whenever it is actually honored — Bitcoin (always,
   * segwit/bech32) and any server-managed coin (the node derives real bip44
   * accounts + xpub) — otherwise deterministic, for a pure client-side non-BTC
   * coin whose in-browser create path only ever produces deterministic wallets.
   */
  private defaultTypeForCoin(coin: BaseCoin): string {
    if (!coin) {
      return 'deterministic';
    }
    return (coin.isBitcoin() || coin.serverWallets) ? 'bip44' : 'deterministic';
  }

  getData(): FormData {
    return {
      label: this.form.value.label,
      seed: this.form.value.seed,
      coin: this.form.value.coin,
      walletType: this.form.value.wallet_type || 'deterministic',
      seedPassphrase: this.form.value.seed_passphrase || '',
      segwit: !!this.form.value.segwit,
      password: this.form.value.password || '',
    };
  }

  initForm(defaultCoin: BaseCoin, create: boolean | null = null) {
    create = create !== null ? create : this.create;

    this.isBitcoinCoin = defaultCoin ? defaultCoin.isBitcoin() : false;
    // Show the wallet-TYPE selector (deterministic / bip44 / segwit) whenever a
    // coin is selected. The type is independent of serverWallets (server-side vs
    // browser custody): a browser wallet (serverWallets=false, e.g. the skywire
    // embedded wallet) still creates deterministic OR bip44 (and Bitcoin bip44
    // can be segwit). Gating on serverWallets wrongly hid every type option for
    // client-side wallets. The individual options are already coin-gated in the
    // template (deterministic only for non-BTC; segwit only for BTC bip44).
    this.showWalletType = !!defaultCoin;
    this.probeCoinHealth(defaultCoin);

    // Default new wallets to bip44 (HD) wherever bip44 is actually honored:
    // Bitcoin (always bip44 — segwit/bech32 by default) and any server-managed
    // coin, whose node performs real bip44 derivation (HD accounts + xpub). Fall
    // back to deterministic only for a pure client-side non-BTC coin, where the
    // in-browser create path always produces a deterministic wallet regardless
    // of the selected type — so a bip44 default there would mislead.
    const defaultWalletType = this.defaultTypeForCoin(defaultCoin);

    this.form = this.formBuilder.group({
        label: new UntypedFormControl('', [ Validators.required ]),
        coin: new UntypedFormControl(defaultCoin, [ Validators.required ]),
        seed: new UntypedFormControl('', [ Validators.required ]),
        confirm_seed: new UntypedFormControl(),
        wallet_type: new UntypedFormControl(defaultWalletType),
        seed_passphrase: new UntypedFormControl(''),
        segwit: new UntypedFormControl(true),
        password: new UntypedFormControl(''),
      },
      {
        validator: create ? this.seedMatchValidator.bind(this) : null,
      }
    );

    if (create) {
      this.generateSeed(128);
    }

    // Update coin-dependent state when coin selection changes
    this.form.get('coin')!.valueChanges.subscribe((coin: BaseCoin) => {
      if (coin) {
        this.isBitcoinCoin = coin.isBitcoin();
        this.showWalletType = true; // type selector is coin-type-gated in the template, not serverWallets-gated
        this.probeCoinHealth(coin);
        this.form.get('wallet_type')!.setValue(this.defaultTypeForCoin(coin));
        if (this.isBitcoinCoin) {
          this.form.get('segwit')!.setValue(true);
        }
      }
      this.changeDetectorRef.markForCheck();
    });

    this.statusSubscription = this.form.statusChanges.subscribe(() => {
      this.customSeedAccepted = false;
      this.normalSeed = this.validateSeed(this.form.get('seed')!.value);
      this.changeDetectorRef.markForCheck();
    });
  }

  generateSeed(entropy: number) {
    this.form.controls.seed.setValue(generateMnemonic(entropy));
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
    return formGroup.get('seed')!.value === formGroup.get('confirm_seed')!.value ? null : { NotEqual: true };
  }
}
