import { Component, OnInit, ViewChild, OnDestroy, AfterViewInit } from '@angular/core';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { Observable } from 'rxjs';

import { WalletService } from '../../../../services/wallet/wallet.service';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { CoinService } from '../../../../services/coin.service';
import { BaseCoin } from '../../../../coins/basecoin';
import { LanguageService, LanguageData } from '../../../../services/language.service';
import { openChangeLanguageModal, showConfirmationModal, scanAddresses } from '../../../../utils';
import { CreateWalletFormComponent } from '../../wallets/create-wallet/create-wallet-form/create-wallet-form.component';
import { ConfirmationData, Wallet } from '../../../../app.datatypes';
import { BlockchainService } from '../../../../services/blockchain.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { config } from '../../../../app.config';
import { MsgBarService } from '../../../../services/msg-bar.service';

@Component({
    selector: 'app-onboarding-create-wallet',
    templateUrl: './onboarding-create-wallet.component.html',
    styleUrls: ['./onboarding-create-wallet.component.scss'],
    standalone: false
})
export class OnboardingCreateWalletComponent implements OnInit, AfterViewInit, OnDestroy {
  @ViewChild('formControl') formControl: CreateWalletFormComponent;
  @ViewChild('create') createButton;

  showSlowMobileInfo = false;
  showNewForm = true;
  doubleButtonActive = DoubleButtonActive.LeftButton;
  userHasWallets = false;
  creatingWallet = false;
  language: LanguageData;

  private slowInfoSubscription: Subscription;
  private subscription: Subscription;

  constructor(
    private dialog: CustomMatDialogService,
    private walletService: WalletService,
    private router: Router,
    private coinService: CoinService,
    private languageService: LanguageService,
    private blockchainService: BlockchainService,
    private translate: TranslateService,
    private msgBarService: MsgBarService,
  ) { }

  ngOnInit() {
    this.checkUserWallets();

    this.subscription = this.languageService.currentLanguage
      .subscribe(lang => this.language = lang);
  }

  ngAfterViewInit() {
    this.formControl.initForm(this.coinService.currentCoin.getValue());
  }

  ngOnDestroy() {
    this.removeSlowInfoSubscription();
    this.msgBarService.hide();
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }

  changeForm(newState: DoubleButtonActive) {
    newState === DoubleButtonActive.RightButton ? this.showNewForm = false : this.showNewForm = true;
    this.formControl.initForm(this.coinService.currentCoin.getValue(), this.showNewForm);
  }

  showSafe() {
    const data: ConfirmationData = {
      text: 'wizard.confirm.desc',
      headerText: 'wizard.confirm.title',
      checkboxText: 'wizard.confirm.checkbox',
      confirmButtonText: 'wizard.confirm.button',
      redTitle: true
    };

    showConfirmationModal(this.dialog, data).afterClosed().subscribe(result => {
      if (result) {
        this.createWallet();
      }
    });
  }

  loadWallet() {
    this.createWallet();
  }

  skip() {
    this.router.navigate(['/wallets'], { replaceUrl: true });
  }

  private showLanguageModal() {
    setTimeout(() => {
      openChangeLanguageModal(this.dialog, true)
        .subscribe(response => {
          if (response) {
            this.languageService.changeLanguage(response);
          }
          this.showDisclaimer();
        });
    }, 0);
  }

  changelanguage() {
    openChangeLanguageModal(this.dialog)
      .subscribe(response => {
        if (response) {
          this.languageService.changeLanguage(response);
        }
      });
  }

  private showDisclaimer() {
    const data: ConfirmationData = {
      text: 'onboarding.disclaimer.disclaimer-description',
      headerText: 'title.disclaimer',
      checkboxText: 'onboarding.disclaimer.disclaimer-check',
      confirmButtonText: 'onboarding.disclaimer.continue-button',
      disableDismiss: true,
    };

    showConfirmationModal(this.dialog, data);
  }

  private checkUserWallets() {
    this.walletService.haveWallets.first().subscribe(result => {
      if (!result) {
        this.userHasWallets = false;
        this.showLanguageModal();
      } else {
        this.userHasWallets = true;
      }
    });
  }

  private createWallet() {
    this.createButton.setLoading();
    this.creatingWallet = true;

    this.slowInfoSubscription = Observable.of(1).delay(config.timeBeforeSlowMobileInfo)
      .subscribe(() => this.showSlowMobileInfo = true);

    const data = this.formControl.getData();

    this.walletService.create(data.label, data.seed, data.coin.id, this.showNewForm, data.walletType, data.seedPassphrase)
      .subscribe(
        wallet => this.onCreateSuccess(wallet, data.coin),
        (error) => this.onCreateError(error.message)
      );
  }

  private onCreateSuccess(wallet: Wallet, coin: BaseCoin) {
    const initialCoin = this.coinService.currentCoin.value;
    this.coinService.changeCoin(coin);

    if (!this.showNewForm) {
      this.showSlowMobileInfo = false;
      this.removeSlowInfoSubscription();

      scanAddresses(this.dialog, wallet, this.blockchainService, this.translate).subscribe(
        response => this.processScanResponse(initialCoin, wallet, false, response),
        error => this.processScanResponse(initialCoin, wallet, true, error)
      );
    } else {
      this.finish();
    }
  }

  private processScanResponse(initialCoin: BaseCoin, wallet: Wallet, isError: boolean, response) {
    if (isError || response !== null) {
      this.coinService.changeCoin(initialCoin);
      this.onCreateError(response.message ? response.message : response.toString());
    } else {
      wallet.needSeedConfirmation = false;
      this.walletService.add(wallet);
      this.finish();
    }
  }

  private finish() {
    this.showSlowMobileInfo = false;
    this.removeSlowInfoSubscription();
    this.createButton.setSuccess();
    this.skip();
    this.creatingWallet = false;
  }

  private onCreateError(errorMesasge: string) {
    this.showSlowMobileInfo = false;
    this.removeSlowInfoSubscription();
    this.msgBarService.showError(errorMesasge);

    this.createButton.resetState();
    this.creatingWallet = false;
  }

  private removeSlowInfoSubscription() {
    if (this.slowInfoSubscription) {
      this.slowInfoSubscription.unsubscribe();
    }
  }
}
