import { Component, Input, OnInit, OnDestroy, Renderer2, ViewChild, NgZone } from '@angular/core';
import 'rxjs/add/observable/interval';
import { Subscription } from 'rxjs';
import { Overlay } from '@angular/cdk/overlay';
import { Observable } from 'rxjs';

import { BalanceService, BalanceStates } from '../../../../services/wallet/balance.service';
import { CoinService } from '../../../../services/coin.service';
import { BaseCoin } from '../../../../coins/basecoin';
import { openChangeCoinModal, openChangeLanguageModal, getTimeSinceLastBalanceUpdate } from '../../../../utils';
import { LanguageService, LanguageData } from '../../../../services/language.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';

@Component({
    selector: 'app-top-bar',
    templateUrl: './top-bar.component.html',
    styleUrls: ['./top-bar.component.scss'],
    standalone: false
})
export class TopBarComponent implements OnInit, OnDestroy {
  @Input() headline: string;

  timeSinceLastBalanceUpdate = 0;
  balanceObtained = false;
  problemUpdatingBalance: boolean;
  updatingBalance = false;
  currentCoin: BaseCoin;
  language: LanguageData;
  hasManyCoins: boolean;

  private subscriptionsGroup: Subscription[] = [];

  constructor(private balanceService: BalanceService,
              private coinService: CoinService,
              private dialog: CustomMatDialogService,
              private overlay: Overlay,
              private renderer: Renderer2,
              private languageService: LanguageService,
              private _ngZone: NgZone) {
  }

  ngOnInit() {
    this.subscriptionsGroup.push(this.languageService.currentLanguage
      .subscribe(lang => this.language = lang));

    this.hasManyCoins = this.coinService.coins.length > 1;

    this.subscriptionsGroup.push(
      this.coinService.currentCoin.subscribe((coin: BaseCoin) => {
        this.currentCoin = coin;
        this.balanceObtained = false;
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
      })
    );

    this._ngZone.runOutsideAngular(() => {
      this.subscriptionsGroup.push(
        Observable.interval(5000).subscribe(() => this._ngZone.run(() => this.timeSinceLastBalanceUpdate = getTimeSinceLastBalanceUpdate(this.balanceService)))
      );
    });
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
  }

  refresBalance() {
    this.balanceService.startGettingBalances();
  }

  changeCoin() {
    openChangeCoinModal(this.dialog, this.renderer, this.overlay)
      .subscribe(response => {
        if (response) {
          this.coinService.changeCoin(response);
        }
      });
  }

  changelanguage() {
    openChangeLanguageModal(this.dialog)
      .subscribe(response => {
        if (response) {
          this.languageService.changeLanguage(response);
        }
      });
  }
}
