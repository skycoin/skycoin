import { Component, OnDestroy, ViewChild, ChangeDetectorRef } from '@angular/core';
import { SubscriptionLike, combineLatest } from 'rxjs';
import { ActivatedRoute, Router } from '@angular/router';
import { FormGroup, FormBuilder, FormControl } from '@angular/forms';
import { map } from 'rxjs/operators';

import { ButtonComponent } from '../../layout/button/button.component';
import { MsgBarService } from '../../../services/msg-bar.service';
import { SoftwareWalletService } from '../../../services/wallet-operations/software-wallet.service';
import { WalletsAndAddressesService } from '../../../services/wallet-operations/wallets-and-addresses.service';
import { WalletBase } from '../../../services/wallet-operations/wallet-objects';

/**
 * Allows to use the seed to remove or change the password of an encrypted software wallet.
 * The URL for opening this page must have a param called "id", with the ID of the wallet
 * to which the password will be changed.
 */
@Component({
  selector: 'app-reset-password',
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.scss'],
})
export class ResetPasswordComponent implements OnDestroy {
  @ViewChild('resetButton') resetButton: ButtonComponent;

  form: FormGroup;
  // Allows to deactivate the form while the component is busy.
  busy = true;

  // Vars with the validation error messages.
  seedErrorMsg = '';
  passwordErrorMsg = '';

  private subscription: SubscriptionLike;
  private wallet: WalletBase;
  private done = false;
  private hideBarWhenClosing = true;

  constructor(
    public formBuilder: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private msgBarService: MsgBarService,
    private softwareWalletService: SoftwareWalletService,
    private walletsAndAddressesService: WalletsAndAddressesService,
    private changeDetector: ChangeDetectorRef,
  ) {
    this.initForm();
    // Get the wallets and route params.
    this.subscription = combineLatest([this.route.params, this.walletsAndAddressesService.allWallets]).pipe(map(result => {
      const params = result[0];
      const wallets = result[1];

      const wallet = wallets.find(w => w.id === params['id']);
      // Abort if the requested wallet does not exists or is temporal.
      if (!wallet || wallet.temporal) {
        setTimeout(() => this.router.navigate([''], {skipLocationChange: true}));

        return;
      }

      this.wallet = wallet;
      this.form.get('wallet').setValue(wallet.label);
      // Activate the form.
      this.busy = false;
    })).subscribe();
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
    if (this.hideBarWhenClosing) {
      this.msgBarService.hide();
    }
  }

  initForm() {
    this.form = new FormGroup({});
    this.form.addControl('wallet', new FormControl());
    this.form.addControl('seed', new FormControl(''));
    this.form.addControl('password', new FormControl(''));
    this.form.addControl('confirm', new FormControl(''));

    this.form.setValidators(this.validateForm.bind(this));
  }

  // Resets the wallet password.
  reset() {
    if (!this.form.valid || this.busy || this.done) {
      return;
    }

    this.busy = true;
    this.msgBarService.hide();
    this.resetButton.setLoading();

    this.softwareWalletService.resetPassword(this.wallet, this.form.value.seed, this.form.value.password !== '' ? this.form.value.password : null)
      .subscribe(() => {
        this.resetButton.setSuccess();
        this.resetButton.setDisabled();
        this.done = true;

        // Show a success msg and avoid closing it after closing this page.
        this.msgBarService.showDone('reset.done');
        this.hideBarWhenClosing = false;

        // Navigate from the page after a small delay.
        setTimeout(() => {
          this.router.navigate(['']);
        }, 2000);
      }, error => {
        // Reactivate the UI and show the error msg.
        this.busy = false;
        this.resetButton.resetState();
        this.msgBarService.showError(error);
      });

    // Avoids a problem with the change detection system.
    this.changeDetector.detectChanges();
  }

  /**
   * Validates the form and updates the vars with the validation errors.
   */
  validateForm() {
    this.seedErrorMsg = '';
    this.passwordErrorMsg = '';

    let valid = true;

    if (!this.form.get('seed').value) {
      valid = false;
      if (this.form.get('seed').touched) {
        this.seedErrorMsg = 'reset.seed-error-info';
      }
    }

    // Check if the 2 passwords entered by the user are equal.
    if (this.form.get('password').value !== this.form.get('confirm').value) {
      valid = false;
      if (this.form.get('confirm').touched) {
        this.passwordErrorMsg = 'reset.confirm-error-info';
      }
    }

    return valid ? null : { Invalid: true };
  }
}
