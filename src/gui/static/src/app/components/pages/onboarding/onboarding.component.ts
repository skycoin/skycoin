import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { SubscriptionLike } from 'rxjs';
import { MatDialog } from '@angular/material/dialog';

import { LanguageData, LanguageService } from '../../../services/language.service';
import { WalletFormData } from '../wallets/create-wallet/create-wallet-form/create-wallet-form.component';
import { MsgBarService } from '../../../services/msg-bar.service';
import { OnboardingEncryptWalletComponent } from './onboarding-encrypt-wallet/onboarding-encrypt-wallet.component';
import { SelectLanguageComponent } from '../../layout/select-language/select-language.component';
import { WalletsAndAddressesService } from '../../../services/wallet-operations/wallets-and-addresses.service';

/**
 * Wizard for creating the first wallet.
 */
@Component({
  selector: 'app-onboarding',
  templateUrl: './onboarding.component.html',
  styleUrls: ['./onboarding.component.scss'],
})
export class OnboardingComponent implements OnInit, OnDestroy {
  @ViewChild('encryptForm') encryptForm: OnboardingEncryptWalletComponent;

  // Current stept to show.
  step = 1;
  // Data entered on the form of the first step.
  formData: WalletFormData;
  // Currently selected language.
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
    this.subscription = this.languageService.currentLanguage.subscribe(lang => this.language = lang);
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  // Called when the user finishes the first step.
  onLabelAndSeedCreated(data: WalletFormData) {
    this.formData = data,
    this.step = 2;
  }

  // Called when the user finishes the second step.
  onPasswordCreated(password: string|null) {
    // Create the wallet.
    this.walletsAndAddressesService.createSoftwareWallet(this.formData.loadTemporarily, this.formData.label, this.formData.seed, password).subscribe(() => {
      this.router.navigate(['/wallets']);
    }, e => {
      this.msgBarService.showError(e);
      // Make the form usable again.
      this.encryptForm.resetButton();
    });
  }

  // Return to step 1.
  onBack() {
    this.step = 1;
  }

  changelanguage() {
    SelectLanguageComponent.openDialog(this.dialog);
  }
}
