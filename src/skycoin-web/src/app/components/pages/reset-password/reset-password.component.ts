import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, ViewChild, inject } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup, ValidationErrors, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription, combineLatest } from 'rxjs';

import { Wallet } from '../../../app.datatypes';
import { WALLET_BASE_PATH, walletPath } from '../../../wallet-path';
import { MsgBarService } from '../../../services/msg-bar.service';
import { WalletService, walletId } from '../../../services/wallet/wallet.service';
import { isBip44 } from '../../../services/wallet/backends/browser-wallet-backend';
import { ButtonComponent } from '../../layout/button/button.component';

/**
 * Uses the seed to change or remove the password of an encrypted wallet, for
 * when the password is forgotten. The route's id param names the wallet (walletId).
 */
@Component({
  selector: 'app-reset-password',
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: false,
})
export class ResetPasswordComponent implements OnDestroy {
  @ViewChild('resetButton') resetButton!: ButtonComponent;

  form: UntypedFormGroup;
  wallet: Wallet | null = null;
  busy = false;
  private done = false;
  private subscription: Subscription;
  private walletBase = inject(WALLET_BASE_PATH);

  constructor(
    formBuilder: UntypedFormBuilder,
    route: ActivatedRoute,
    private router: Router,
    private walletService: WalletService,
    private msgBarService: MsgBarService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    this.form = formBuilder.group({
      seed: ['', Validators.required],
      passphrase: [''],
      password: [''],
      confirm: [''],
    }, { validators: this.matches });

    this.subscription = combineLatest([route.params, walletService.currentWallets]).subscribe(([params, wallets]) => {
      // A reset under way saves the wallets, which must not send the page away.
      if (this.done || this.busy) {
        return;
      }
      const wallet = wallets.find(w => !w.isHardware && walletId(w) === params['id'] && walletService.isEncrypted(w));
      if (!wallet) {
        this.router.navigate([walletPath(this.walletBase, '/wallets')]);

        return;
      }
      this.wallet = wallet;
      this.changeDetectorRef.markForCheck();
    });
  }

  /** A browser bip44 wallet keeps its seed passphrase encrypted too, so it must be entered again. */
  get needsPassphrase(): boolean {
    return !!this.wallet && (isBip44(this.wallet) || (this.wallet.walletType === 'bip44' && !!this.wallet.filename));
  }

  get mismatch(): boolean {
    return !!this.form.errors?.['mismatch'] && this.form.get('confirm')!.dirty;
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  reset() {
    if (!this.wallet || !this.form.valid || this.busy || this.done) {
      return;
    }

    this.busy = true;
    this.msgBarService.hide();
    this.resetButton.setLoading();
    const { seed, passphrase, password } = this.form.value;
    this.walletService.resetPassword(this.wallet, seed.trim(), passphrase || '', password || null).subscribe(() => {
      this.done = true;
      this.resetButton.setSuccess();
      this.resetButton.setDisabled();
      this.msgBarService.showDone('reset.done');
      setTimeout(() => {
        this.router.navigate([walletPath(this.walletBase, '/wallets')]);
        this.changeDetectorRef.markForCheck();
      }, 2000);
      this.changeDetectorRef.markForCheck();
    }, error => {
      this.busy = false;
      this.resetButton.resetState();
      this.msgBarService.showError(error && error.message ? error.message : String(error));
      this.changeDetectorRef.markForCheck();
    });
  }

  private matches(group: AbstractControl): ValidationErrors | null {
    return group.get('password')!.value === group.get('confirm')!.value ? null : { mismatch: true };
  }
}
