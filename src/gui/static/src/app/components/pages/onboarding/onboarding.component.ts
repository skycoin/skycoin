import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { WalletService } from '../../../services/wallet.service';

@Component({
  selector: 'app-onboarding',
  templateUrl: './onboarding.component.html',
  styleUrls: ['./onboarding.component.scss']
})
export class OnboardingComponent {

  step = 1;
  label: string;
  seed: string;
  password: string|null;

  constructor(
    private router: Router,
    private walletService: WalletService,
  ) { }

  onLabelAndSeedCreated(data: [string, string]) {
    this.label = data[0];
    this.seed = data[1];

    this.step = 2;
  }

  onPasswordCreated(password: string|null) {
    this.password = password;

    this.createWallet();
  }

  onBack() {
    this.step = 1;
  }

  get fill() {
    return this.label ? { label: this.label, seed: this.seed } : null;
  }

  private createWallet() {
    this.walletService.create(this.label, this.seed, 100, this.password).subscribe(() => {
      this.router.navigate(['/wallets']);
    });
  }

}
