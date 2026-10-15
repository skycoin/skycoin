import { ChangeDetectionStrategy, ChangeDetectorRef, Component, ViewChild, inject } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';

import { WalletService } from '../../../services/wallet/wallet.service';
import { ButtonComponent } from '../button/button.component';
import { WALLET_BASE_PATH, walletPath } from '../../../wallet-path';

/**
 * Covers the wallet while the wallets in this browser are stored under a
 * password. Forgetting the password erases those wallets.
 */
@Component({
  selector: 'app-app-lock',
  templateUrl: './app-lock.component.html',
  styleUrls: ['./app-lock.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: false,
})
export class AppLockComponent {
  @ViewChild('unlockButton') unlockButton!: ButtonComponent;
  form: UntypedFormGroup;
  wrongPassword = false;
  confirmingErase = false;
  private walletBase = inject(WALLET_BASE_PATH);

  constructor(
    formBuilder: UntypedFormBuilder,
    private walletService: WalletService,
    private changeDetectorRef: ChangeDetectorRef,
    private router: Router,
  ) {
    this.form = formBuilder.group({ password: ['', Validators.required] });
  }

  unlock() {
    if (!this.form.valid || this.unlockButton.isLoading()) {
      return;
    }
    this.wrongPassword = false;
    this.unlockButton.setLoading();
    this.walletService.unlockApp(this.form.value.password).subscribe({
      error: () => {
        this.unlockButton.resetState();
        this.wrongPassword = true;
        this.changeDetectorRef.markForCheck();
      },
    });
  }

  erase() {
    this.walletService.eraseLockedWallets();
    this.router.navigate([walletPath(this.walletBase, '/wizard')]);
  }
}
