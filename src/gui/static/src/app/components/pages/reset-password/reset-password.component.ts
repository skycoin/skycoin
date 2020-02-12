import { Component, OnDestroy, ViewChild } from '@angular/core';
import { SubscriptionLike,  combineLatest } from 'rxjs';
import { ButtonComponent } from '../../layout/button/button.component';
import { FormGroup, FormBuilder, FormControl, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { MsgBarService } from '../../../services/msg-bar.service';
import { SoftwareWalletService } from '../../../services/wallet-operations/software-wallet.service';
import { WalletsAndAddressesService } from '../../../services/wallet-operations/wallets-and-addresses.service';
import { WalletBase } from '../../../services/wallet-operations/wallet-objects';

@Component({
  selector: 'app-reset-password',
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.scss'],
})
export class ResetPasswordComponent implements OnDestroy {
  @ViewChild('resetButton', { static: false }) resetButton: ButtonComponent;

  form: FormGroup;
  busy = false;

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
  ) {
    this.initForm('');
    this.subscription = combineLatest(this.route.params, this.walletsAndAddressesService.allWallets, (params, wallets) => {
      const wallet = wallets.find(w => w.id === params['id']);
      if (!wallet) {
        setTimeout(() => this.router.navigate([''], {skipLocationChange: true}));

        return;
      }

      this.wallet = wallet;
      this.initForm(wallet.label);
    }).subscribe();
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
    if (this.hideBarWhenClosing) {
      this.msgBarService.hide();
    }
  }

  initForm(walletName: string) {
    const validators = [];
    validators.push(this.passwordMatchValidator.bind(this));

    this.form = new FormGroup({}, validators);
    this.form.addControl('wallet', new FormControl(walletName));
    this.form.addControl('seed', new FormControl('', [Validators.required]));
    this.form.addControl('password', new FormControl());
    this.form.addControl('confirm', new FormControl());
  }

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

        this.hideBarWhenClosing = false;
        this.msgBarService.showDone('reset.done');

        setTimeout(() => {
          this.router.navigate(['']);
        }, 2000);
      }, error => {
        this.busy = false;
        this.resetButton.resetState();
        this.msgBarService.showError(error);
      });
  }

  private passwordMatchValidator() {
    if (this.form && this.form.get('password') && this.form.get('confirm')) {
      return this.form.get('password').value === this.form.get('confirm').value ? null : { NotEqual: true };
    } else {
      return { NotEqual: true };
    }
  }
}
