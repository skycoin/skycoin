import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, NgZone, OnDestroy, OnInit, Renderer2, ViewChild } from '@angular/core';
import { Subscription, interval } from 'rxjs';

import { BalanceService, BalanceStates } from '../../../../services/wallet/balance.service';
import { CoinService } from '../../../../services/coin.service';
import { BaseCoin } from '../../../../coins/basecoin';
import { openChangeLanguageModal, getTimeSinceLastBalanceUpdate } from '../../../../utils';
import { LanguageService, LanguageData } from '../../../../services/language.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { CipherProvider } from '../../../../services/cipher.provider';

@Component({
    selector: 'app-top-bar',
    templateUrl: './top-bar.component.html',
    styleUrls: ['./top-bar.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class TopBarComponent implements OnInit, OnDestroy {
  @Input() headline!: string;

  timeSinceLastBalanceUpdate = 0;
  balanceObtained = false;
  problemUpdatingBalance!: boolean;
  updatingBalance = false;
  currentCoin!: BaseCoin;
  language!: LanguageData;
  hasManyCoins!: boolean;
  availableCoins: BaseCoin[] = [];

  /**
   * Identifies the wasm cipher this wallet is running, as a short commit.
   *
   * The cipher is a committed build artifact rather than something compiled
   * alongside the page, so the wallet's own version does not identify it. A
   * "modified" suffix means it was built from a working tree with uncommitted
   * changes, so no commit describes it — CI rejects that, and seeing it here
   * means the wallet is running something nobody can reproduce.
   */
  get cipherVersion(): string {
    const build = this.cipherProvider.buildInfo;
    if (!build) {
      return '';
    }

    return build.commit.slice(0, 7) + (build.modified ? ' (modified)' : '');
  }

  private subscriptionsGroup: Subscription[] = [];

  constructor(private balanceService: BalanceService,
              private coinService: CoinService,
              private dialog: CustomMatDialogService,
              private renderer: Renderer2,
              private languageService: LanguageService,
              private _ngZone: NgZone,
              private changeDetectorRef: ChangeDetectorRef,
              private cipherProvider: CipherProvider,) {
  }

  ngOnInit() {
    this.subscriptionsGroup.push(this.languageService.currentLanguage
      .subscribe(lang => {
        this.language = lang;
        this.changeDetectorRef.markForCheck();
      }));

    this.hasManyCoins = this.coinService.coins.length > 1;
    this.availableCoins = this.coinService.coins;

    this.subscriptionsGroup.push(
      this.coinService.currentCoin.subscribe((coin: BaseCoin) => {
        this.currentCoin = coin;
        this.balanceObtained = false;
        this.changeDetectorRef.markForCheck();
      })
    );

    this.subscriptionsGroup.push(
      this.balanceService.totalBalance.subscribe(balance => {
        if (balance) {
          if (balance.state === BalanceStates.Obtained) {
            this.balanceObtained = true;
          }
          this.timeSinceLastBalanceUpdate = getTimeSinceLastBalanceUpdate(this.balanceService);
          this.problemUpdatingBalance = balance.state === BalanceStates.Error;
          this.updatingBalance = balance.state === BalanceStates.Updating;
        }
        this.changeDetectorRef.markForCheck();
      })
    );

    this._ngZone.runOutsideAngular(() => {
      this.subscriptionsGroup.push(
        interval(5000).subscribe(() => {
          this._ngZone.run(() => this.timeSinceLastBalanceUpdate = getTimeSinceLastBalanceUpdate(this.balanceService));
          this.changeDetectorRef.markForCheck();
        })
      );
    });
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
  }

  refresBalance() {
    this.balanceService.startGettingBalances();
  }

  selectCoin(coin: BaseCoin) {
    if (coin && coin.id !== this.currentCoin.id) {
      this.coinService.changeCoin(coin);
    }
  }

  changelanguage() {
    openChangeLanguageModal(this.dialog)
      .subscribe(response => {
        if (response) {
          this.languageService.changeLanguage(response);
        }
        this.changeDetectorRef.markForCheck();
      });
  }
}
