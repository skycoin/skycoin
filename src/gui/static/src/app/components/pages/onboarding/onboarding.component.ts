import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { LanguageData, LanguageService } from '../../../services/language.service';
import { SubscriptionLike } from 'rxjs';
import { MatDialog } from '@angular/material/dialog';
import { WalletFormData } from '../wallets/create-wallet/create-wallet-form/create-wallet-form.component';
import { MsgBarService } from '../../../services/msg-bar.service';
import { OnboardingEncryptWalletComponent } from './onboarding-encrypt-wallet/onboarding-encrypt-wallet.component';
import { SelectLanguageComponent } from '../../layout/select-language/select-language.component';
import { WalletsAndAddressesService } from '../../../services/wallet-operations/wallets-and-addresses.service';

@Component({
  selector: 'app-onboarding',
  templateUrl: './onboarding.component.html',
  styleUrls: ['./onboarding.component.scss'],
})
export class OnboardingComponent implements OnInit, OnDestroy {
  @ViewChild('encryptForm') encryptForm: OnboardingEncryptWalletComponent;

  step = 1;
  formData: WalletFormData;
  password: string|null;
  language: LanguageData;

  private subscription: SubscriptionLike;

  constructor(
    private router: Router,
    private languageService: LanguageService,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    private walletsAndAddressesService: WalletsAndAddressesService,
  ) { }

  ngOnInit() {
    this.subscription = this.languageService.currentLanguage
      .subscribe(lang => this.language = lang);
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  onLabelAndSeedCreated(data: WalletFormData) {
    this.formData = data,
    this.step = 2;
  }

  onPasswordCreated(password: string|null) {
    this.password = password;

    this.createWallet();
  }

  onBack() {
    this.step = 1;
  }

  changelanguage() {
    SelectLanguageComponent.openDialog(this.dialog);
  }

  get fill() {
    return this.formData;
  }

  private createWallet() {
    this.walletsAndAddressesService.createSoftwareWallet(this.formData.label, this.formData.seed, this.password).subscribe(() => {
      this.router.navigate(['/wallets']);
    }, e => {
      this.msgBarService.showError(e);
      this.encryptForm.resetButton();
    });
  }
}
