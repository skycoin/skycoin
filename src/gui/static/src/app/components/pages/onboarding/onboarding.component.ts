import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router } from '@angular/router';
import { WalletService } from '../../../services/wallet.service';
import { LanguageData, LanguageService } from '../../../services/language.service';
import { ISubscription } from 'rxjs/Subscription';
import { openChangeLanguageModal } from '../../../utils';
import { MatDialog } from '@angular/material';

@Component({
  selector: 'app-onboarding',
  templateUrl: './onboarding.component.html',
  styleUrls: ['./onboarding.component.scss'],
})
export class OnboardingComponent implements OnInit, OnDestroy {
  step = 1;
  label: string;
  seed: string;
  create: boolean;
  password: string|null;
  language: LanguageData;

  private subscription: ISubscription;

  constructor(
    private router: Router,
    private walletService: WalletService,
    private languageService: LanguageService,
    private dialog: MatDialog,
  ) { }

  ngOnInit() {
    this.subscription = this.languageService.currentLanguage
      .subscribe(lang => this.language = lang);
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  onLabelAndSeedCreated(data: [string, string, boolean]) {
    this.label = data[0];
    this.seed = data[1];
    this.create = data[2];

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
    openChangeLanguageModal(this.dialog)
      .subscribe(response => {
        if (response) {
          this.languageService.changeLanguage(response);
        }
      });
  }

  get fill() {
    return this.label ? { label: this.label, seed: this.seed, create: this.create } : null;
  }

  private createWallet() {
    this.walletService.create(this.label, this.seed, 100, this.password).subscribe(() => {
      this.router.navigate(['/wallets']);
    });
  }
}
